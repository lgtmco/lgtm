package web

import (
	"regexp"

	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	"github.com/lgtmco/lgtm/store"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-version"
)

func Hook(c *gin.Context) {
	hook, err := remote.GetHook(c, c.Request)
	if err != nil {
		log.Errorf("Error parsing hook. %s", err)
		c.String(500, "Error parsing hook. %s", err)
		return
	}
	if hook != nil {
		processCommentHook(c, hook)
	}

	statusHook, err := remote.GetStatusHook(c, c.Request)
	if err != nil {
		log.Errorf("Error parsing status hook. %s", err)
		c.String(500, "Error parsing status hook. %s", err)
		return
	}
	if statusHook != nil {
		processStatusHook(c, statusHook)
	}

	if hook == nil && statusHook == nil {
		c.String(200, "pong")
		return
	}
}

func processStatusHook(c *gin.Context, hook *model.StatusHook) {
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

	config, maintainer, err := getConfigAndMaintainers(c, user, repo)
	if err != nil {
		return
	}

	if !config.DoMerge {
		c.IndentedJSON(200, gin.H{})
		return
	}

	merged := map[string]string{}
	vers := map[string]string{}

	pullRequests, err := remote.GetPullRequestsForCommit(c, user, hook.Repo, &hook.SHA)
	log.Debugf("sha for commit is %s, pull requests are: %s", hook.SHA, pullRequests)

	if err != nil {
		log.Errorf("Error while getting pull requests for commit %s %s", hook.SHA, err)
		c.String(500, "Error while getting pull requests for commit %s %s", hook.SHA, err)
		return
	}
	//check the statuses of all of the checks on the branches for this commit
	for _, v := range pullRequests {
		//if all of the statuses are success, then merge and create a tag for the version
		if v.Branch.BranchStatus == "success" && v.Branch.Mergeable {
			sha, err := remote.MergePR(c, user, hook.Repo, v)
			if err != nil {
				log.Warnf("Unable to merge pull request %s: %s", v.Title, err)
				continue
			} else {
				log.Debugf("Merged pull request %s", v.Title)
			}

			merged[v.Title] = *sha

			if !config.DoVersion {
				continue
			}

			// to create the version, need to scan the comments on the pull request to see if anyone specified a version #
			// if so, use the largest specified version #. if not, increment the last version version # for the release
			maxVer, err := remote.GetMaxExistingTag(c, user, hook.Repo)
			if err != nil {
				log.Warnf("Unable to find the max version tag for %s/%s: %s", hook.Repo.Owner, hook.Repo.Name, err)
				continue
			}

			comments, err := getComments(c, user, repo, v.Number)
			if err != nil {
				log.Warnf("Unable to find the comments for pull request %s: %s", v.Title, err)
				continue
			}

			foundVersion := getMaxVersionComment(config, maintainer, v.Issue, comments)

			if foundVersion != nil && foundVersion.GreaterThan(maxVer) {
				maxVer = foundVersion
			} else {
				maxParts := maxVer.Segments()
				maxVer, _ = version.NewVersion(fmt.Sprintf("%d.%d.%d", maxParts[0], maxParts[1], maxParts[2]+1))
			}

			err = remote.Tag(c, user, repo, maxVer, sha)
			if err != nil {
				log.Warnf("Unable to tag branch %s: %s", v.Title, err)
				continue
			}
			vers[v.Title] = maxVer.String()
		}
	}
	log.Debugf("processed status for %s. received %v ", repo.Slug, hook)

	c.IndentedJSON(200, gin.H{
		"merged":   merged,
		"versions": vers,
	})
}

func processCommentHook(c *gin.Context, hook *model.Hook) {

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

	config, maintainer, err := getConfigAndMaintainers(c, user, repo)
	if err != nil {
		return
	}

	comments, err := getComments(c, user, repo, hook.Issue.Number)
	if err != nil {
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

	log.Debugf("processed comment for %s. received %d of %d approvals", repo.Slug, len(approvers), config.Approvals)

	c.IndentedJSON(200, gin.H{
		"approvers":   maintainer.People,
		"settings":    config,
		"approved":    approved,
		"approved_by": approvers,
	})
}

func getConfigAndMaintainers(c *gin.Context, user *model.User, repo *model.Repo) (*model.Config, *model.Maintainer, error) {
	rcfile, _ := remote.GetContents(c, user, repo, ".lgtm")
	config, err := model.ParseConfig(rcfile)
	if err != nil {
		log.Errorf("Error parsing .lgtm file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing .lgtm file. %s.", err)
		return nil, nil, err
	}

	// THIS IS COMPLETELY DUPLICATED IN THE API SECTION. NOT IDEAL
	file, err := remote.GetContents(c, user, repo, "MAINTAINERS")
	if err != nil {
		log.Debugf("no MAINTAINERS file for %s. Checking for team members.", repo.Slug)
		members, merr := cache.GetMembers(c, user, repo.Owner)
		if merr != nil {
			log.Errorf("Error getting repository %s. %s", repo.Slug, err)
			log.Errorf("Error getting org members %s. %s", repo.Owner, merr)
			c.String(404, "MAINTAINERS file not found. %s", err)
			return nil, nil, err
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
		return nil, nil, err
	}
	return config, maintainer, nil
}

func getComments(c *gin.Context, user *model.User, repo *model.Repo, num int) ([]*model.Comment, error) {
	comments, err := remote.GetComments(c, user, repo, num)
	if err != nil {
		log.Errorf("Error retrieving comments for %s pr %d. %s", repo.Slug, num, err)
		c.String(500, "Error retrieving comments. %s.", err)
		return nil, err
	}
	return comments, nil
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

// getMaxVersionComment is a helper function that analyzes the list of comments
// and returns the maximum version found in a comment.
func getMaxVersionComment(config *model.Config, maintainer *model.Maintainer, issue model.Issue, comments []*model.Comment) *version.Version {
	approverm := map[string]bool{}
	approvers := []*model.Person{}

	matcher, err := regexp.Compile(config.Pattern)
	if err != nil {
		// this should never happen
		return nil
	}

	var maxVersion *version.Version

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
		m := matcher.FindStringSubmatch(comment.Body)
		if len(m) > 0 {
			approverm[comment.Author] = true
			approvers = append(approvers, person)

			if len(m) > 1 {
				//has a version
				curVersion, err := version.NewVersion(m[1])
				if err != nil {
					continue
				}
				if maxVersion == nil || curVersion.GreaterThan(maxVersion) {
					maxVersion = curVersion
				}
			}
		}
	}

	return maxVersion
}
