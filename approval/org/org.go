package org

import (
	"regexp"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/approval"
)

func init() {
	approval.Register("org", Org)
}

// Org is a helper function that analyzes the list of comments
// and returns the list of approvers.
// The rules for Org are:
// - At most one approver from each team
// - If SelfApprovalOff is true, no member of the team of the creator of the Pull Request is allowed to approve the PR
// - If a person appears on more than one team, they only count once, for the first team in which they appear
// (not solving the optimal grouping problem yet)
func Org(config *model.Config, maintainer *model.Maintainer, issue *model.Issue, comments []*model.Comment, p approval.Processor) {
	//groups that have already approved
	approvergm := map[string]bool{}

	//org that the author belongs to
	authorOrg := ""

	//key is person, value is map of the orgs for that person
	orgMap := map[string]map[string]bool{}

	//key is org name, value is org
	for k, v := range maintainer.Org {
		//value is name of person in the org
		for _, name := range v.People {
			if name == issue.Author {
				authorOrg = name
			}
			m, ok := orgMap[name]
			if !ok {
				m := map[string]bool{}
				orgMap[name] = m
			}
			m[k] = true
		}
	}

	matcher, err := regexp.Compile(config.Pattern)
	if err != nil {
		// this should never happen
		return
	}

	for _, comment := range comments {
		// verify the comment matches the approval pattern
		if !matcher.MatchString(comment.Body) {
			continue
		}
		//get the orgs for the current comment's author
		curOrgs, ok := orgMap[comment.Author]
		if !ok {
			// the user must be a valid maintainer of the project
			continue
		}
		for curOrg := range curOrgs {
			// your group cannot lgtm your own pull request
			if config.SelfApprovalOff && curOrg == authorOrg {
				continue
			}
			// your group cannot approve twice
			if approvergm[curOrg] {
				continue
			}
			//found one!
			approvergm[curOrg] = true
			p(maintainer, comment)
		}
	}
	return
}
