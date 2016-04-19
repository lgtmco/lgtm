package remote

// Account represents a user or team account.
type Account struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar"`
	Kind   string `json:"type"`
}

// Issue represents an issue or pull request.
type Issue struct {
	Number int    `json:"issue"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Comment represents a user comment on an issue
// or pull request.
type Comment struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}
