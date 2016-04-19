package model

type Hook struct {
	Repo    *Repo
	Issue   *Issue
	Comment *Comment
}
