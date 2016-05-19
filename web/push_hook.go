package web

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
)

func processPushHook(c *gin.Context, pushHook *model.PushHook) {
	repo, user, err := getRepoAndUser(c, pushHook.Repo.Slug)
	if err != nil {
		return
	}
	updated, err := remote.UpdatePRsForCommit(c, user, repo, &pushHook.SHA)
	if err != nil {
		log.Errorf("Error setting status. %s", err)
		c.String(500, "Error setting status. %s", err)
		return
	}
	if updated {
		c.IndentedJSON(200, gin.H{
			"commit": pushHook.SHA,
			"status": "pending",
		})
	} else {
		c.IndentedJSON(200, gin.H{})
	}
}
