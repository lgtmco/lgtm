package remote

//go:generate mockery -name Remote -output mock -case=underscore

import (
	"net/http"

	"github.com/lgtmco/lgtm/model"
	"golang.org/x/net/context"
)

type Remote interface {
	// GetUser authenticates a user with the remote system.
	GetUser(http.ResponseWriter, *http.Request) (*model.User, error)

	// GetUserToken authenticates a user with the remote system using
	// the remote systems OAuth token.
	GetUserToken(string) (string, error)

	// GetTeams gets a team list from the remote system.
	GetTeams(*model.User) ([]*model.Team, error)

	// GetMembers gets a team member list from the remote system.
	GetMembers(*model.User, string) ([]*model.Member, error)

	// GetRepo gets a repository from the remote system.
	GetRepo(*model.User, string, string) (*model.Repo, error)

	// GetPerm gets a repository permission from the remote system.
	GetPerm(*model.User, string, string) (*model.Perm, error)

	// GetRepo gets a repository list from the remote system.
	GetRepos(*model.User) ([]*model.Repo, error)

	// SetHook adds a webhook to the remote repository.
	SetHook(*model.User, *model.Repo, string) error

	// DelHook deletes a webhook from the remote repository.
	DelHook(*model.User, *model.Repo, string) error

	// GetComments gets pull request comments from the remote system.
	GetComments(*model.User, *model.Repo, int) ([]*model.Comment, error)

	// GetContents gets the file contents from the remote system.
	GetContents(*model.User, *model.Repo, string) ([]byte, error)

	// SetStatus adds or updates the pull request status in the remote system.
	SetStatus(*model.User, *model.Repo, int, bool) error

	// GetHook gets the hook from the http Request.
	GetHook(r *http.Request) (*model.Hook, error)

	// GetStatusHook gets the status hook from the http Request.
	GetStatusHook(r *http.Request) (*model.StatusHook, error)

	// GetPRHook gets the pull request hook from the http Request.
	GetPRHook(r *http.Request) (*model.PRHook, error)

	// GetPushHook gets the push hook from the http Request.
	GetPushHook(r *http.Request) (*model.PushHook, error)

	// GetBranchStatus returns overall status for the named branch from the remote system
	GetBranchStatus(*model.User, *model.Repo, string) (*model.BranchStatus, error)

	// MergePR merges the named pull request from the remote system
	MergePR(u *model.User, r *model.Repo, pullRequest model.PullRequest, approvers []*model.Person) (*string, error)

	// GetMaxExistingTag finds the highest version across all tags
	ListTags(u *model.User, r *model.Repo) ([]model.Tag, error)

	// Tag applies a tag with the specified version to the specified sha
	Tag(u *model.User, r *model.Repo, version *string, sha *string) error

	// GetPullRequestsForCommit returns all pull requests associated with a commit SHA
	GetPullRequestsForCommit(u *model.User, r *model.Repo, sha *string) ([]model.PullRequest, error)

	// UpdatePRsForCommit sets the commit's status to pending for LGTM if it is already on an open Pull Request
	UpdatePRsForCommit(u *model.User, r *model.Repo, sha *string) (bool, error)
}

// GetUser authenticates a user with the remote system.
func GetUser(c context.Context, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	return FromContext(c).GetUser(w, r)
}

// GetUserToken authenticates a user with the remote system using
// the remote systems OAuth token.
func GetUserToken(c context.Context, token string) (string, error) {
	return FromContext(c).GetUserToken(token)
}

// GetTeams gets a team list from the remote system.
func GetTeams(c context.Context, u *model.User) ([]*model.Team, error) {
	return FromContext(c).GetTeams(u)
}

// GetMembers gets a team members list from the remote system.
func GetMembers(c context.Context, u *model.User, team string) ([]*model.Member, error) {
	return FromContext(c).GetMembers(u, team)
}

// GetRepo gets a repository from the remote system.
func GetRepo(c context.Context, u *model.User, owner, name string) (*model.Repo, error) {
	return FromContext(c).GetRepo(u, owner, name)
}

// GetPerm gets a repository permission from the remote system.
func GetPerm(c context.Context, u *model.User, owner, name string) (*model.Perm, error) {
	return FromContext(c).GetPerm(u, owner, name)
}

// GetRepos gets a repository list from the remote system.
func GetRepos(c context.Context, u *model.User) ([]*model.Repo, error) {
	return FromContext(c).GetRepos(u)
}

// GetComments gets pull request comments from the remote system.
func GetComments(c context.Context, u *model.User, r *model.Repo, num int) ([]*model.Comment, error) {
	return FromContext(c).GetComments(u, r, num)
}

// GetContents gets the file contents from the remote system.
func GetContents(c context.Context, u *model.User, r *model.Repo, path string) ([]byte, error) {
	return FromContext(c).GetContents(u, r, path)
}

// SetHook adds a webhook to the remote repository.
func SetHook(c context.Context, u *model.User, r *model.Repo, hook string) error {
	return FromContext(c).SetHook(u, r, hook)
}

// DelHook deletes a webhook from the remote repository.
func DelHook(c context.Context, u *model.User, r *model.Repo, hook string) error {
	return FromContext(c).DelHook(u, r, hook)
}

// SetStatus adds or updates the pull request status in the remote system.
func SetStatus(c context.Context, u *model.User, r *model.Repo, num int, ok bool) error {
	return FromContext(c).SetStatus(u, r, num, ok)
}

// GetHook gets the hook from the http Request.
func GetHook(c context.Context, r *http.Request) (*model.Hook, error) {
	return FromContext(c).GetHook(r)
}

// GetStatusHook gets the status hook from the http Request.
func GetStatusHook(c context.Context, r *http.Request) (*model.StatusHook, error) {
	return FromContext(c).GetStatusHook(r)
}

// GetPRHook gets the pull request hook from the http Request.
func GetPRHook(c context.Context, r *http.Request) (*model.PRHook, error) {
	return FromContext(c).GetPRHook(r)
}

// GetPushHook gets the push hook from the http Request.
func GetPushHook(c context.Context, r *http.Request) (*model.PushHook, error) {
	return FromContext(c).GetPushHook(r)
}

// GetBranchStatus gets the overal status for a branch from the remote repository.
func GetBranchStatus(c context.Context, u *model.User, r *model.Repo, branch string) (*model.BranchStatus, error) {
	return FromContext(c).GetBranchStatus(u, r, branch)
}

func MergePR(c context.Context, u *model.User, r *model.Repo, pullRequest model.PullRequest, approvers []*model.Person) (*string, error) {
	return FromContext(c).MergePR(u, r, pullRequest, approvers)
}

func ListTags(c context.Context, u *model.User, r *model.Repo) ([]model.Tag, error) {
	return FromContext(c).ListTags(u, r)
}

func Tag(c context.Context, u *model.User, r *model.Repo, version *string, sha *string) error {
	return FromContext(c).Tag(u, r, version, sha)
}

func GetPullRequestsForCommit(c context.Context, u *model.User, r *model.Repo, sha *string) ([]model.PullRequest, error) {
	return FromContext(c).GetPullRequestsForCommit(u, r, sha)

}

func UpdatePRsForCommit(c context.Context, u *model.User, r *model.Repo, sha *string) (bool, error) {
	return FromContext(c).UpdatePRsForCommit(u, r, sha)
}