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
	GetMembers(*model.User, string, string) ([]*model.Member, error)

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
	SetStatus(*model.User, *model.Repo, int, int, int) error

	// GetHook gets the hook from the http Request.
	GetHook(r *http.Request) (*model.Hook, error)
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
func GetMembers(c context.Context, u *model.User, owner string, maintainers string) ([]*model.Member, error) {
	return FromContext(c).GetMembers(u, owner, maintainers)
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
func SetStatus(c context.Context, u *model.User, r *model.Repo, num, granted, required int) error {
	return FromContext(c).SetStatus(u, r, num, granted, required)
}

// GetHook gets the hook from the http Request.
func GetHook(c context.Context, r *http.Request) (*model.Hook, error) {
	return FromContext(c).GetHook(r)
}
