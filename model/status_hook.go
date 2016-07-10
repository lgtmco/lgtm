package model

type StatusHook struct {
	SHA  string
	Repo *Repo
}

type PullRequest struct {
	Issue
	Branch Branch
}

type Branch struct {
	Name      string
	Status    string
	Mergeable bool
}
