package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/shared/httputil"
	"golang.org/x/oauth2"
	"errors"
)

// name of the status message posted to GitHub
const context = "approvals/lgtm"

type Github struct {
	URL    string
	API    string
	Client string
	Secret string
	Scopes []string
}

func (g *Github) GetUser(res http.ResponseWriter, req *http.Request) (*model.User, error) {

	var config = &oauth2.Config{
		ClientID:     g.Client,
		ClientSecret: g.Secret,
		RedirectURL:  fmt.Sprintf("%s/login", httputil.GetURL(req)),
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", g.URL),
			TokenURL: fmt.Sprintf("%s/login/oauth/access_token", g.URL),
		},
		Scopes: g.Scopes,
	}

	// get the oauth code from the incoming request. if no code is present
	// redirec the user to GitHub login to retrieve a code.
	var code = req.FormValue("code")
	if len(code) == 0 {
		state := fmt.Sprintln(time.Now().Unix())
		http.Redirect(res, req, config.AuthCodeURL(state), http.StatusSeeOther)
		return nil, nil
	}

	// exchanges the oauth2 code for an access token
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("Error exchanging token. %s", err)
	}

	// get the currently authenticated user details for the access token
	client := setupClient(g.API, token.AccessToken)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, fmt.Errorf("Error fetching user. %s", err)
	}

	return &model.User{
		Login:  *user.Login,
		Token:  token.AccessToken,
		Avatar: *user.AvatarURL,
	}, nil
}

func (g *Github) GetUserToken(token string) (string, error) {
	client := setupClient(g.API, token)
	user, _, err := client.Users.Get("")
	if err != nil {
		return "", fmt.Errorf("Error fetching user. %s", err)
	}
	return *user.Login, nil
}

func (g *Github) GetTeams(user *model.User) ([]*model.Team, error) {
	client := setupClient(g.API, user.Token)
	orgs, _, err := client.Organizations.List("", &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, fmt.Errorf("Error fetching teams. %s", err)
	}
	teams := []*model.Team{}
	for _, org := range orgs {
		team := model.Team{
			Login:  *org.Login,
			Avatar: *org.AvatarURL,
		}
		teams = append(teams, &team)
	}
	return teams, nil
}

func (g *Github) GetMembers(user *model.User, team string) ([]*model.Member, error) {
	client := setupClient(g.API, user.Token)
	teams, _, err := client.Organizations.ListTeams(team, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, fmt.Errorf("Error accessing team list. %s", err)
	}
	var id int
	for _, team := range teams {
		if strings.ToLower(*team.Name) == "maintainers" {
			id = *team.ID
			break
		}
	}
	if id == 0 {
		return nil, fmt.Errorf("Error finding approvers team. %s", err)
	}
	opts := github.OrganizationListTeamMembersOptions{}
	opts.PerPage = 100
	teammates, _, err := client.Organizations.ListTeamMembers(id, &opts)
	if err != nil {
		return nil, fmt.Errorf("Error fetching team members. %s", err)
	}
	var members []*model.Member
	for _, teammate := range teammates {
		members = append(members, &model.Member{
			Login: *teammate.Login,
		})
	}
	return members, nil
}

func (g *Github) GetRepo(user *model.User, owner, name string) (*model.Repo, error) {
	client := setupClient(g.API, user.Token)
	repo_, _, err := client.Repositories.Get(owner, name)
	if err != nil {
		return nil, fmt.Errorf("Error fetching repository. %s", err)
	}
	return &model.Repo{
		Owner:   owner,
		Name:    name,
		Slug:    *repo_.FullName,
		Link:    *repo_.HTMLURL,
		Private: *repo_.Private,
	}, nil
}

func (g *Github) GetPerm(user *model.User, owner, name string) (*model.Perm, error) {
	client := setupClient(g.API, user.Token)
	repo, _, err := client.Repositories.Get(owner, name)
	if err != nil {
		return nil, fmt.Errorf("Error fetching repository. %s", err)
	}
	m := &model.Perm{}
	m.Admin = (*repo.Permissions)["admin"]
	m.Push = (*repo.Permissions)["push"]
	m.Pull = (*repo.Permissions)["pull"]
	return m, nil
}

func (g *Github) GetRepos(u *model.User) ([]*model.Repo, error) {
	client := setupClient(g.API, u.Token)
	all, err := GetUserRepos(client)
	if err != nil {
		return nil, err
	}

	repos := []*model.Repo{}
	for _, repo := range all {
		// only list repositories that I can admin
		if repo.Permissions == nil || (*repo.Permissions)["admin"] == false {
			continue
		}
		repos = append(repos, &model.Repo{
			Owner:   *repo.Owner.Login,
			Name:    *repo.Name,
			Slug:    *repo.FullName,
			Link:    *repo.HTMLURL,
			Private: *repo.Private,
		})
	}

	return repos, nil
}

func (g *Github) SetHook(user *model.User, repo *model.Repo, link string) error {
	client := setupClient(g.API, user.Token)

	old, err := GetHook(client, repo.Owner, repo.Name, link)
	if err == nil && old != nil {
		client.Repositories.DeleteHook(repo.Owner, repo.Name, *old.ID)
	}

	_, err = CreateHook(client, repo.Owner, repo.Name, link)
	if err != nil {
		log.Debugf("Error creating the webhook at %s. %s", link, err)
		return err
	}

	repo_, _, err := client.Repositories.Get(repo.Owner, repo.Name)
	if err != nil {
		return err
	}

	in := new(Branch)
	in.Protection.Enabled = true
	in.Protection.Checks.Enforcement = "non_admins"

	/*
		JCB 04/21/16 confirmed with Github support -- must specify all existing contexts when
		 adding a new one, otherwise the other contexts will be removed.
	 */
	client_ := NewClientToken(g.API, user.Token)
	branch, _ := client_.Branch(repo.Owner, repo.Name, *repo_.DefaultBranch)

	in.Protection.Checks.Contexts = append(buildOtherContextSlice(branch), context)

	err = client_.BranchProtect(repo.Owner, repo.Name, *repo_.DefaultBranch, in)
	if err != nil {
		if g.URL == "https://github.com" {
			return err
		}
		log.Warnf("Error configuring protected branch for %s/%s@%s. %s", repo.Owner, repo.Name, *repo_.DefaultBranch, err)
	}
	return nil
}

func (g *Github) DelHook(user *model.User, repo *model.Repo, link string) error {
	client := setupClient(g.API, user.Token)

	hook, err := GetHook(client, repo.Owner, repo.Name, link)
	if err != nil {
		return err
	} else if hook == nil {
		return nil
	}
	_, err = client.Repositories.DeleteHook(repo.Owner, repo.Name, *hook.ID)
	if err != nil {
		return err
	}

	repo_, _, err := client.Repositories.Get(repo.Owner, repo.Name)
	if err != nil {
		return err
	}

	client_ := NewClientToken(g.API, user.Token)
	branch, _ := client_.Branch(repo.Owner, repo.Name, *repo_.DefaultBranch)
	if len(branch.Protection.Checks.Contexts) == 0 {
		return nil
	}

	branch.Protection.Checks.Contexts = buildOtherContextSlice(branch)

	return client_.BranchProtect(repo.Owner, repo.Name, *repo_.DefaultBranch, branch)
}

// buildOtherContextSlice returns all contexts besides the one for LGTM
func buildOtherContextSlice(branch *Branch) []string {
	checks := []string{}
	for _, check := range branch.Protection.Checks.Contexts {
		if check != context {
			checks = append(checks, check)
		}
	}
	return checks
}

func (g *Github) GetComments(u *model.User, r *model.Repo, num int) ([]*model.Comment, error) {
	client := setupClient(g.API, u.Token)

	opts := github.IssueListCommentsOptions{Direction: "desc", Sort: "created"}
	opts.PerPage = 100
	comments_, _, err := client.Issues.ListComments(r.Owner, r.Name, num, &opts)
	if err != nil {
		return nil, err
	}
	comments := []*model.Comment{}
	for _, comment := range comments_ {
		comments = append(comments, &model.Comment{
			Author: *comment.User.Login,
			Body:   *comment.Body,
		})
	}
	return comments, nil
}

func (g *Github) GetContents(u *model.User, r *model.Repo, path string) ([]byte, error) {
	client := setupClient(g.API, u.Token)
	content, _, _, err := client.Repositories.GetContents(r.Owner, r.Name, path, nil)
	if err != nil {
		return nil, err
	}
	return content.Decode()
}

func (g *Github) SetStatus(u *model.User, r *model.Repo, num int, ok bool) error {
	client := setupClient(g.API, u.Token)

	pr, _, err := client.PullRequests.Get(r.Owner, r.Name, num)
	if err != nil {
		return err
	}

	status := "pending"
	desc := "this commit is pending approval"
	if ok {
		status = "success"
		desc = "this commit looks good"
	}

	data := github.RepoStatus{
		Context:     github.String(context),
		State:       github.String(status),
		Description: github.String(desc),
	}

	_, _, err = client.Repositories.CreateStatus(r.Owner, r.Name, *pr.Head.SHA, &data)
	return err
}

func (g *Github) GetHook(r *http.Request) (*model.Hook, error) {

	// only process comment hooks
	if r.Header.Get("X-Github-Event") != "issue_comment" {
		return nil, nil
	}

	data := commentHook{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Issue.PullRequest.Link) == 0 {
		return nil, nil
	}

	hook := new(model.Hook)
	hook.Issue = new(model.Issue)
	hook.Issue.Number = data.Issue.Number
	hook.Issue.Author = data.Issue.User.Login
	hook.Repo = new(model.Repo)
	hook.Repo.Owner = data.Repository.Owner.Login
	hook.Repo.Name = data.Repository.Name
	hook.Repo.Slug = data.Repository.FullName
	hook.Comment = new(model.Comment)
	hook.Comment.Body = data.Comment.Body
	hook.Comment.Author = data.Comment.User.Login

	return hook, nil
}

func (g *Github) GetStatusHook(r *http.Request) (*model.StatusHook, error) {

	// only process comment hooks
	if r.Header.Get("X-Github-Event") != "status" {
		return nil, nil
	}

	data := statusHook{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	log.Debug(data)

	if data.State != "success" {
		return nil, nil
	}

	hook := new(model.StatusHook)

	hook.SHA = data.SHA

	hook.Repo = new(model.Repo)
	hook.Repo.Owner = data.Repository.Owner.Login
	hook.Repo.Name = data.Repository.Name
	hook.Repo.Slug = data.Repository.FullName

	log.Debug(*hook)

	return hook, nil
}

func (g *Github) GetPRHook(r *http.Request) (*model.PRHook, error) {

	// only process comment hooks
	if r.Header.Get("X-Github-Event") != "pull_request" {
		return nil, nil
	}

	data := prHook{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	log.Debug(data)

	if data.Action != "opened" {
		return nil, nil
	}

	hook := new(model.PRHook)

	hook.Number = data.Number
	hook.Repo = new(model.Repo)
	hook.Repo.Owner = data.Repository.Owner.Login
	hook.Repo.Name = data.Repository.Name
	hook.Repo.Slug = data.Repository.FullName

	log.Debug(*hook)

	return hook, nil
}

func (g *Github) GetPullRequestsForCommit(u *model.User, r *model.Repo, sha *string) ([]model.PullRequest, error) {
	client := setupClient(g.API, u.Token)
	log.Debug("sha == ", sha, *sha)
	issues, _, err := client.Search.Issues(fmt.Sprintf("%s&type=pr", *sha), &github.SearchOptions{
		TextMatch: false,
	})
	if err != nil {
		return nil, err
	}
	out := make([]model.PullRequest, len(issues.Issues))
	for k, v := range issues.Issues {
		pr, _, err := client.PullRequests.Get(r.Owner, r.Name, *v.Number)
		if err != nil {
			return nil, err
		}

		mergeable := true
		if pr.Mergeable != nil {
			mergeable = *pr.Mergeable
		}

		status, _, err := client.Repositories.GetCombinedStatus(r.Owner, r.Name, *sha, nil)
		if err != nil {
			return nil, err
		}

		log.Debug("current issue ==", v)
		log.Debug("current pr ==", *pr)
		log.Debug("combined status ==", *status)

		combinedState := *status.State
		if combinedState == "success" {
			log.Debug("overall status is success -- checking to see if all status checks returned success")
			for _, v := range status.Statuses {
				log.Debugf("status check %s returned %s", *v.Context, *v.State)
				if *v.State != "success" {
					log.Debugf("setting combined status check to %s", *v.State)
					combinedState = *v.State
				}
			}
		}
		out[k] = model.PullRequest{
			Issue: model.Issue{
				Number: *v.Number,
				Title: *v.Title,
				Author: *v.User.Login,
			},
			Branch:  model.Branch{
				Name: *pr.Head.Ref,
				BranchStatus: combinedState,
				Mergeable: mergeable,

			},
		}
	}
	return out, nil
}

func (g *Github) GetBranchStatus(u *model.User, r *model.Repo, branch string) (*model.BranchStatus, error) {
	client := setupClient(g.API, u.Token)
	statuses, _, err := client.Repositories.GetCombinedStatus(r.Owner, r.Name, branch, nil)
	if err != nil {
		return nil, err
	}

	return (*model.BranchStatus)(statuses.State), nil
}

func (g *Github) MergePR(u *model.User, r *model.Repo, pullRequest model.PullRequest, approvers []*model.Person) (*string, error) {
	client := setupClient(g.API, u.Token)

	msg := "Merged by LGTM\n"
	if len(approvers) > 0 {
		apps := "Approved by:\n"
		for _, v := range approvers {
			//Brad Rydzewski <brad.rydzewski@mail.com> (@bradrydzewski)
			if len(v.Name) > 0 {
				apps += fmt.Sprintf("%s", v.Name)
			}
			if len(v.Email) > 0 {
				apps += fmt.Sprintf(" <%s>", v.Email)
			}
			if len(v.Login) > 0 {
				apps += fmt.Sprintf(" (@%s)", v.Login)
			}
			apps += "\n"
		}
		msg += apps
	}
	result, _, err := client.PullRequests.Merge(r.Owner, r.Name, pullRequest.Number, msg)
	if err != nil {
		return nil, err
	}

	if !(*result.Merged) {
		return nil, errors.New(*result.Message)
	}
	return result.SHA, nil
}

func (g *Github) ListTags(u *model.User, r *model.Repo) ([]model.Tag, error) {
	client := setupClient(g.API, u.Token)

	tags, _, err := client.Repositories.ListTags(r.Owner, r.Name, nil)
	if err != nil {
		return nil, err
	}
	out := make([]model.Tag, len(tags))
	for k, v := range tags {
		out[k] = model.Tag(*v.Name)
	}
	return out, nil
}

func (g *Github) Tag(u *model.User, r *model.Repo, version *string, sha *string) error {
	client := setupClient(g.API, u.Token)

	t := time.Now()
	tag, _, err := client.Git.CreateTag(r.Owner, r.Name, &github.Tag{
		Tag:     version,
		SHA:     sha,
		Message: github.String("Tagged by LGTM"),
		Tagger: &github.CommitAuthor{
			Date:  &t,
			Name:  github.String("LGTM"),
			Email: github.String("LGTM@lgtm.co"),
		},
		Object: &github.GitObject{
			SHA:  sha,
			Type: github.String("commit"),
		},
	})

	if err != nil {
		return err
	}
	_, _, err = client.Git.CreateRef(r.Owner, r.Name, &github.Reference{
		Ref: github.String("refs/tags/" + *version),
		Object: &github.GitObject{
			SHA: tag.SHA,
		},
	})

	return err
}
