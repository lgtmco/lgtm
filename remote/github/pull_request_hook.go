package github

import "github.com/lgtmco/lgtm/model"

// pullRequestHook represents a subset of the pull_request payload.
type pullRequestHook struct {
	Action string `json:"action"`
	Number int    `json:"number"`

	PullRequest struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`

	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}

func (prh *pullRequestHook) toHook() *model.Hook {
	return &model.Hook{
		Issue: &model.Issue{
			Number: prh.Number,
			Author: prh.PullRequest.User.Login,
		},
		Repo: &model.Repo{
			Owner: prh.Repository.Owner.Login,
			Name:  prh.Repository.Name,
			Slug:  prh.Repository.FullName,
		},
	}
}

func (prh *pullRequestHook) analysable() bool {
	return prh.Action == "opened" ||
		prh.Action == "synchronize"
}
