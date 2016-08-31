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

	repo_, _, err := client.Repositories.Get(repo.Owner, repo.Name)
	if err != nil {
		return err
	}

	old, err := GetHook(client, repo.Owner, repo.Name, link)
	if err == nil && old != nil {
		client.Repositories.DeleteHook(repo.Owner, repo.Name, *old.ID)
	}

	_, err = CreateHook(client, repo.Owner, repo.Name, link)
	if err != nil {
		log.Debugf("Error creating the webhook at %s. %s", link, err)
		return err
	}

	in := new(Branch)
	in.Protection.Enabled = true
	in.Protection.Checks.Enforcement = "non_admins"
	in.Protection.Checks.Contexts = []string{context}

	client_ := NewClientToken(g.API, user.Token)
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
	checks := []string{}
	for _, check := range branch.Protection.Checks.Contexts {
		if check != context {
			checks = append(checks, check)
		}
	}
	branch.Protection.Checks.Contexts = checks
	return client_.BranchProtect(repo.Owner, repo.Name, *repo_.DefaultBranch, branch)
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

func (g *Github) SetStatus(u *model.User, r *model.Repo, num, granted, required int) error {
	client := setupClient(g.API, u.Token)

	pr, _, err := client.PullRequests.Get(r.Owner, r.Name, num)
	if err != nil {
		return err
	}

	status := "success"
	desc := "this commit looks good"

	if granted < required {
		status = "pending"
		desc = fmt.Sprintf("%d of %d required approvals granted", granted, required)
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
	eventType := r.Header.Get("X-Github-Event")
	if eventType == "issue_comment" {
		return getCommentHook(r)
	} else if eventType == "pull_request" {
		return getPRHook(r)
	} else {
		return nil, nil
	}
}

func getCommentHook(r *http.Request) (*model.Hook, error) {
	data := commentHook{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	if len(data.Issue.PullRequest.Link) == 0 {
		return nil, nil
	}
	return data.toHook(), nil
}

func getPRHook(r *http.Request) (*model.Hook, error) {
	data := pullRequestHook{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.analysable() {
		return nil, nil
	}
	return data.toHook(), nil
}
