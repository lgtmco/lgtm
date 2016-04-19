package notifier

// Notification represents a notification that we are sending to a list of
// maintainers indicating a commit is ready for their review and, hopefully,
// approval.
type Notification struct {
	Reviewers []*Reviewer
	Commit    *Commit
}

// Reviewer represents a repository maintainer or contributor that is being
// notified of a commit to review.
type Reviewer struct {
	Login string
	Email string
}

// Commit represents the commit for which we are notifiying the maintainers.
type Commit struct {
	Repo    string
	Message string
	Author  string
	Link    string
}
