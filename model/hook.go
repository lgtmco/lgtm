package model

type Hook struct {
	Kind         string
	IssueComment *IssueCommentHook
	Status       *StatusHook
}

type IssueCommentHook struct {
	Repo    *Repo
	Issue   *Issue
	Comment *Comment
}
