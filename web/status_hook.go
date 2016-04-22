package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-version"
	"regexp"
)

type StatusResponse struct {
	SHA     *string `json:"sha,omitempty"`
	Version *string `json:"version,omitempty"`
}

func processStatusHook(c *gin.Context, hook *model.StatusHook) {
	repo, user, err := getRepoAndUser(c, hook.Repo.Slug)
	if err != nil {
		return
	}

	config, maintainer, err := getConfigAndMaintainers(c, user, repo)
	if err != nil {
		return
	}

	if !config.DoMerge {
		c.IndentedJSON(200, gin.H{
		})
		return
	}

	merged := map[string]StatusResponse{}

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

			merged[v.Title] = StatusResponse{
				SHA: sha,
			}

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

			foundVersion := getMaxVersionComment(config, maintainer, v.Issue,comments)

			if foundVersion != nil && foundVersion.GreaterThan(maxVer) {
				maxVer = foundVersion
			} else {
				maxParts := maxVer.Segments()
				maxVer, _ = version.NewVersion(fmt.Sprintf("%d.%d.%d", maxParts[0], maxParts[1], maxParts[2] + 1))
			}

			err = remote.Tag(c, user, repo, maxVer, sha)
			if err != nil {
				log.Warnf("Unable to tag branch %s: %s", v.Title, err)
				continue
			}
			verStr := maxVer.String()
			result := merged[v.Title]
			result.Version = &verStr
		}
	}
	log.Debugf("processed status for %s. received %v ", repo.Slug, hook)

	c.IndentedJSON(200, gin.H{
		"merged":    merged,
	})
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
