package api

import (
	"fmt"

	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	"github.com/lgtmco/lgtm/router/middleware/session"
	"github.com/lgtmco/lgtm/shared/httputil"
	"github.com/lgtmco/lgtm/shared/token"
	"github.com/lgtmco/lgtm/store"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// GetRepos gets the active repository list.
func GetRepos(c *gin.Context) {
	user := session.User(c)
	repos, err := cache.GetRepos(c, user)
	if err != nil {
		log.Errorf("Error getting remote repository list. %s", err)
		c.String(500, "Error getting remote repository list")
		return
	}

	// copy the slice since we are going to mutate it and don't
	// want any nasty data races if the slice came from the cache.
	repoc := make([]*model.Repo, len(repos))
	copy(repoc, repos)

	repom, err := store.GetRepoIntersectMap(c, repos)
	if err != nil {
		log.Errorf("Error getting active repository list. %s", err)
		c.String(500, "Error getting active repository list")
		return
	}

	// merges the slice of active and remote repositories favoring
	// and swapping in local repository information when possible.
	for i, repo := range repoc {
		repo_, ok := repom[repo.Slug]
		if ok {
			repoc[i] = repo_
		}
	}
	c.IndentedJSON(200, repoc)
}

// GetRepo gets the repository by slug.
func GetRepo(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
	)
	repo, err := store.GetRepoOwnerName(c, owner, name)
	if err != nil {
		log.Errorf("Error getting repository %s. %s", name, err)
		c.String(404, "Error getting repository %s", name)
		return
	}
	c.JSON(200, repo)
}

// PostRepo activates a new repository.
func PostRepo(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
		user  = session.User(c)
	)

	// verify repo doesn't already exist
	if _, err := store.GetRepoOwnerName(c, owner, name); err == nil {
		c.AbortWithStatus(409)
		c.String(409, "Error activating a repository that is already active.")
		return
	}

	repo, err := remote.GetRepo(c, user, owner, name)
	if err != nil {
		c.String(404, "Error finding repository in GitHub. %s")
		return
	}
	repo.UserID = user.ID
	repo.Secret = model.Rand()

	// creates a token to authorize the link callback url
	t := token.New(token.HookToken, repo.Slug)
	sig, err := t.Sign(repo.Secret)
	if err != nil {
		c.String(500, "Error activating repository. %s")
		return
	}

	// create the hook callback url
	link := fmt.Sprintf(
		"%s/hook?access_token=%s",
		httputil.GetURL(c.Request),
		sig,
	)
	err = remote.SetHook(c, user, repo, link)
	if err != nil {
		c.String(500, "Error creating hook. %s", err)
		return
	}

	statusLink := fmt.Sprintf(
		"%s/status_hook?access_token=%s",
		httputil.GetURL(c.Request),
		sig,
	)
	err = remote.SetStatusHook(c, user, repo, statusLink)
	if err != nil {
		c.String(500, "Error creating status hook. %s", err)
		return
	}

	err = store.CreateRepo(c, repo)
	if err != nil {
		c.String(500, "Error activating the repository. %s", err)
		return
	}
	c.IndentedJSON(200, repo)
}

// DeleteRepo deletes a repository configuration.
func DeleteRepo(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
		user  = session.User(c)
	)
	repo, err := store.GetRepoOwnerName(c, owner, name)
	if err != nil {
		log.Errorf("Error getting repository %s. %s", name, err)
		c.AbortWithStatus(404)
		return
	}
	err = store.DeleteRepo(c, repo)
	if err != nil {
		log.Errorf("Error deleting repository %s. %s", name, err)
		c.AbortWithStatus(500)
		return
	}
	link := fmt.Sprintf(
		"%s/hook",
		httputil.GetURL(c.Request),
	)
	err = remote.DelHook(c, user, repo, link)
	if err != nil {
		log.Errorf("Error deleting repository hook for %s. %s", name, err)
	}
	c.String(200, "")
}
