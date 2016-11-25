package github

import "github.com/lgtmco/lgtm/model"

// commentHook represents a subset of the issue_comment payload.
type commentHook struct {
	Issue struct {
		Link   string `json:"html_url"`
		Number int    `json:"number"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`

		PullRequest struct {
			Link string `json:"html_url"`
		} `json:"pull_request"`
	} `json:"issue"`

	Comment struct {
		Body string `json:"body"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"comment"`

	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Desc     string `json:"description"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login  string `json:"login"`
			Type   string `json:"type"`
			Avatar string `json:"avatar_url"`
		} `json:"owner"`
	} `json:"repository"`
}

func (ch *commentHook) toHook() *model.Hook {
	return &model.Hook{
		Issue: &model.Issue{
			Number: ch.Issue.Number,
			Author: ch.Issue.User.Login,
		},
		Repo: &model.Repo{
			Owner: ch.Repository.Owner.Login,
			Name:  ch.Repository.Name,
			Slug:  ch.Repository.FullName,
		},
		Comment: &model.Comment{
			Body:   ch.Comment.Body,
			Author: ch.Comment.User.Login,
		},
	}
}
