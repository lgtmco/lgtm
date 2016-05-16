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

	approvers, err := buildApprovers(c, user, repo, config, maintainer, hook.Issue)
	if err != nil {
		return
	}

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

func buildApprovers(c *gin.Context, user *model.User, repo *model.Repo, config *model.Config, maintainer *model.Maintainer, issue *model.Issue) ([]*model.Person, error) {
	comments, err := getComments(c, user, repo, issue.Number)
	if err != nil {
		log.Errorf("Error getting comments for %s/%s/%d", repo.Owner, repo.Name, issue.Number)
		c.String(500, "Error getting comments for %s/%s/%d", repo.Owner, repo.Name, issue.Number)
		return nil, err
	}

	alg, err := approval.Lookup(config.ApprovalAlg)
	if err != nil {
		log.Errorf("Error getting approval algorithm %s. %s", config.ApprovalAlg, err)
		c.String(500, "Error getting approval algorithm %s. %s", config.ApprovalAlg, err)
		return nil, err
	}
	approvers := getApprovers(config, maintainer, issue, comments, alg)
	return approvers, nil
}

func getApprovers(config *model.Config, maintainer *model.Maintainer, issue *model.Issue, comments []*model.Comment, matcher approval.Func) []*model.Person {
	approvers := []*model.Person{}
	matcher(config, maintainer, issue, comments, func(maintainer *model.Maintainer, comment *model.Comment) {
		approvers = append(approvers, maintainer.People[comment.Author])
	})
	return approvers
}
