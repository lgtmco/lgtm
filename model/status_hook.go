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
	Name         string
	BranchStatus string
	Mergeable    bool
	BaseName     string
}

type PRHook struct {
	Number int
	Update bool
	Repo   *Repo
}

type PushHook struct {
	SHA  string
	Repo *Repo
}
