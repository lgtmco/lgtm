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

func Hook(c *gin.Context) {
	hook, err := remote.GetHook(c, c.Request)
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

	// THIS IS COMPLETELY DUPLICATED IN THE API SECTION. NOT IDEAL
	file, err := remote.GetContents(c, user, repo, "MAINTAINERS")
	if err != nil {
		log.Debugf("no MAINTAINERS file for %s. Checking for team members.", repo.Slug)
		members, merr := cache.GetMembers(c, user, repo.Owner, config.Team)
		if merr != nil {
			log.Errorf("Error getting repository %s. %s", repo.Slug, err)
			log.Errorf("Error getting org members %s. %s", repo.Owner, merr)
			c.String(404, "MAINTAINERS file not found. %s", err)
			return
		} else {
			for _, member := range members {
				file = append(file, member.Login...)
				file = append(file, '\n')
			}
		}
	}

	maintainer, err := model.ParseMaintainer(file)
	if err != nil {
		log.Errorf("Error parsing MAINTAINERS file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing MAINTAINERS file. %s.", err)
		return
	}

	comments, err := remote.GetComments(c, user, repo, hook.Issue.Number)
	if err != nil {
		log.Errorf("Error retrieving comments for %s pr %d. %s", repo.Slug, hook.Issue.Number, err)
		c.String(500, "Error retrieving comments. %s.", err)
		return
	}
	approvers := getApprovers(config, maintainer, hook.Issue, comments)
	approved := len(approvers) >= config.Approvals
	err = remote.SetStatus(c, user, repo, hook.Issue.Number, len(approvers), config.Approvals)
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

// getApprovers is a helper function that analyzes the list of comments
// and returns the list of approvers.
func getApprovers(config *model.Config, maintainer *model.Maintainer, issue *model.Issue, comments []*model.Comment) []*model.Person {
	approverm := map[string]bool{}
	approvers := []*model.Person{}

	matcher, err := regexp.Compile(config.Pattern)
	if err != nil {
		// this should never happen
		return approvers
	}

	for _, comment := range comments {
		// cannot lgtm your own pull request
		if config.SelfApprovalOff && comment.Author == issue.Author {
			continue
		}
		// the user must be a valid maintainer of the project
		person, ok := maintainer.People[comment.Author]
		if !ok {
			continue
		}
		// the same author can't approve something twice
		if _, ok := approverm[comment.Author]; ok {
			continue
		}
		// verify the comment matches the approval pattern
		if matcher.MatchString(comment.Body) {
			approverm[comment.Author] = true
			approvers = append(approvers, person)
		}
	}

	return approvers
}
