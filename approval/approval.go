package approval

import (
	"github.com/lgtmco/lgtm/model"
	"regexp"
	"strings"
	"fmt"
)

// Func takes in the information needed to figure out which approvers were in the PR comments
// and returns a slice of the approvers that were found
type Func func (*model.Config, *model.Maintainer, *model.Issue, []*model.Comment) []*model.Person

var approvalMap = map[string]Func{}

func Register(name string, f Func) error {
	if _, ok := approvalMap[strings.ToLower(name)]; ok {
		return fmt.Errorf("Approval Algorithm %s is already registered.",name)
	}
	approvalMap[strings.ToLower(name)] = f
	return nil
}

func Lookup(name string) (Func, error) {
	f, ok := approvalMap[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("Unknown Approval Algorithm %s",name)
	}
	return f, nil
}

func init() {
	Register("simple", Simple)
}

// Simple is a helper function that analyzes the list of comments
// and returns the list of approvers.
func Simple(config *model.Config, maintainer *model.Maintainer, issue *model.Issue, comments []*model.Comment) []*model.Person {
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


