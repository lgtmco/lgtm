package web

import (
	"github.com/gin-gonic/gin"
	log "github.com/Sirupsen/logrus"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	"github.com/lgtmco/lgtm/approval"
)

func processCommentHook(c *gin.Context, hook *model.Hook) {
	repo, user, err := getRepoAndUser(c, hook.Repo.Slug)
	if err != nil {
		return
	}

	config, maintainer, err := getConfigAndMaintainers(c, user, repo)
	if err != nil {
		return
	}

	comments, err := getComments(c, user, repo, hook.Issue.Number)
	if err != nil {
		return
	}

	alg, err := approval.Lookup(config.ApprovalAlg)
	if err != nil {
		log.Errorf("Error getting approval algorithm %s. %s", config.ApprovalAlg, err)
		c.String(500, "Error getting approval algorithm %s. %s", config.ApprovalAlg, err)
		return
	}
	approvers := alg(config, maintainer, hook.Issue, comments)
	approved := len(approvers) >= config.Approvals
	err = remote.SetStatus(c, user, repo, hook.Issue.Number, approved)
	if err != nil {
		log.Errorf("Error setting status for %s pr %d. %s", repo.Slug, hook.Issue.Number, err)
		c.String(500, "Error setting status. %s.", err)
		return
	}

	log.Debugf("processed comment for %s. received %d of %d approvals", repo.Slug, len(approvers), config.Approvals)

	c.IndentedJSON(200, gin.H{
		"approvers":   maintainer.People,
		"settings":    config,
		"approved":    approved,
		"approved_by": approvers,
	})
}
