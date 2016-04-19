package store

import (
	"path"

	"github.com/lgtmco/lgtm/model"

	"golang.org/x/net/context"
)

//go:generate mockery -name Store -output mock -case=underscore

// Store defines a data storage abstraction for managing structured data
// in the system.
type Store interface {
	// GetUser gets a user by unique ID.
	GetUser(int64) (*model.User, error)

	// GetUserLogin gets a user by unique Login name.
	GetUserLogin(string) (*model.User, error)

	// CreateUser creates a new user account.
	CreateUser(*model.User) error

	// UpdateUser updates a user account.
	UpdateUser(*model.User) error

	// DeleteUser deletes a user account.
	DeleteUser(*model.User) error

	// GetRepo gets a repo by unique ID.
	GetRepo(int64) (*model.Repo, error)

	// GetRepoSlug gets a repo by its full name.
	GetRepoSlug(string) (*model.Repo, error)

	// GetRepoMulti gets a list of multiple repos by their full name.
	GetRepoMulti(...string) ([]*model.Repo, error)

	// GetRepoOwner gets a list by owner.
	GetRepoOwner(string) ([]*model.Repo, error)

	// CreateRepo creates a new repository.
	CreateRepo(*model.Repo) error

	// UpdateRepo updates a user repository.
	UpdateRepo(*model.Repo) error

	// DeleteRepo deletes a user repository.
	DeleteRepo(*model.Repo) error
}

// GetUser gets a user by unique ID.
func GetUser(c context.Context, id int64) (*model.User, error) {
	return FromContext(c).GetUser(id)
}

// GetUserLogin gets a user by unique Login name.
func GetUserLogin(c context.Context, login string) (*model.User, error) {
	return FromContext(c).GetUserLogin(login)
}

// CreateUser creates a new user account.
func CreateUser(c context.Context, user *model.User) error {
	return FromContext(c).CreateUser(user)
}

// UpdateUser updates a user account.
func UpdateUser(c context.Context, user *model.User) error {
	return FromContext(c).UpdateUser(user)
}

// DeleteUser deletes a user account.
func DeleteUser(c context.Context, user *model.User) error {
	return FromContext(c).DeleteUser(user)
}

// GetRepo gets a repo by unique ID.
func GetRepo(c context.Context, id int64) (*model.Repo, error) {
	return FromContext(c).GetRepo(id)
}

// GetRepoSlug gets a repo by its full name.
func GetRepoSlug(c context.Context, slug string) (*model.Repo, error) {
	return FromContext(c).GetRepoSlug(slug)
}

// GetRepoOwnerName gets a repo by its owner and name.
func GetRepoOwnerName(c context.Context, owner, name string) (*model.Repo, error) {
	return GetRepoSlug(c, path.Join(owner, name))
}

// GetRepoMulti gets a list of multiple repos by their full name.
func GetRepoMulti(c context.Context, slug ...string) ([]*model.Repo, error) {
	return FromContext(c).GetRepoMulti(slug...)
}

// GetRepoOwner gets a repo list by account.
func GetRepoOwner(c context.Context, owner string) ([]*model.Repo, error) {
	return FromContext(c).GetRepoOwner(owner)
}

// GetRepoIntersect gets a repo list by account login.
func GetRepoIntersect(c context.Context, repos []*model.Repo) ([]*model.Repo, error) {
	slugs := make([]string, len(repos))
	for i, repo := range repos {
		slugs[i] = repo.Slug
	}
	return GetRepoMulti(c, slugs...)
}

// GetRepoIntersectMap gets a repo set by account login where the key is
// the repository slug and the value is the repository struct.
func GetRepoIntersectMap(c context.Context, repos []*model.Repo) (map[string]*model.Repo, error) {
	repos, err := GetRepoIntersect(c, repos)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*model.Repo, len(repos))
	for _, repo := range repos {
		set[repo.Slug] = repo
	}
	return set, nil
}

// CreateRepo creates a new repository.
func CreateRepo(c context.Context, repo *model.Repo) error {
	return FromContext(c).CreateRepo(repo)
}

// UpdateRepo updates a user repository.
func UpdateRepo(c context.Context, repo *model.Repo) error {
	return FromContext(c).UpdateRepo(repo)
}

// DeleteRepo deletes a user repository.
func DeleteRepo(c context.Context, repo *model.Repo) error {
	return FromContext(c).DeleteRepo(repo)
}
