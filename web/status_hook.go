package web

import (
	"regexp"

	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	"github.com/lgtmco/lgtm/store"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func StatusHook(c *gin.Context) {
	hook, err := remote.GetStatusHook(c, c.Request)
	if err != nil {
		log.Errorf("Error parsing hook. %s", err)
		c.String(500, "Error parsing hook. %s", err)
		return
	}
	if hook == nil {
		c.String(200, "pong")
		return
	}

	repo, err := store.GetRepoSlug(c, hook.Repo.Slug)
	if err != nil {
		log.Errorf("Error getting repository %s. %s", hook.Repo.Slug, err)
		c.String(404, "Repository not found.")
		return
	}
	user, err := store.GetUser(c, repo.UserID)
	if err != nil {
		log.Errorf("Error getting repository owner %s. %s", repo.Slug, err)
		c.String(404, "Repository owner not found.")
		return
	}

	rcfile, _ := remote.GetContents(c, user, repo, ".lgtm")
	config, err := model.ParseConfig(rcfile)
	if err != nil {
		log.Errorf("Error parsing .lgtm file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing .lgtm file. %s.", err)
		return
	}

	//todo check the statuses of all of the checks on the branches for this commit
	//todo if all of the statuses are success, then merge and create a tag for the version
	//todo to create the version, need to scan the comments on the pull request to see if anyone specified a version #
	//todo if so, use the largest specified version #. if not, increment the last version version # for the release
	/*
	comments, err := remote.GetComments(c, user, repo, hook.Issue.Number)
	if err != nil {
		log.Errorf("Error retrieving comments for %s pr %d. %s", repo.Slug, hook.Issue.Number, err)
		c.String(500, "Error retrieving comments. %s.", err)
		return
	}
	approvers := getApprovers(config, maintainer, hook.Issue, comments)
	approved := len(approvers) >= config.Approvals
	err = remote.SetStatus(c, user, repo, hook.Issue.Number, approved)
	if err != nil {
		log.Errorf("Error setting status for %s pr %d. %s", repo.Slug, hook.Issue.Number, err)
		c.String(500, "Error setting status. %s.", err)
		return
	}

	*/
	log.Debugf("processed status for %s. received %v ", repo.Slug, hook)

	c.IndentedJSON(200, gin.H{
		//"approvers":   maintainer.People,
		"settings":    config,
		//"approved":    approved,
		//"approved_by": approvers,
	})
}
