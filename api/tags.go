package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/router/middleware/session"
	"github.com/lgtmco/lgtm/store"
)

func GetTags(c *gin.Context) {
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
	tags, err := cache.GetTags(c, user, repo)
	if err != nil {
		log.Errorf("Error getting remote tag list. %s", err)
		c.String(500, "Error getting remote tag list")
		return
	}

	// copy the slice since we don't
	// want any nasty data races if the slice came from the cache.
	tagsc := make(model.TagList, len(tags))
	copy(tagsc, tags)

	c.JSON(200, tagsc)
}
