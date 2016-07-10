package github

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string  { return e.Message }
func (e Error) String() string { return e.Message }

type Branch struct {
	Protection struct {
		Enabled bool `json:"enabled"`
		Checks  struct {
			Enforcement string   `json:"enforcement_level"`
			Contexts    []string `json:"contexts"`
		} `json:"required_status_checks"`
	} `json:"protection"`
}

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

	Repository Repository `json:"repository"`
}

type statusHook struct {
	SHA   string `json:"sha"`
	State string `json:"state"`

	Branches []struct {
		Name   string `json:"name"`
		Commit struct {
			SHA string `json:"sha"`
			URL string `json:"url"`
		} `json:"commit"`
	} `json:"branches"`

	Repository Repository `json:"repository"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Desc     string `json:"description"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login  string `json:"login"`
		Type   string `json:"type"`
		Avatar string `json:"avatar_url"`
	} `json:"owner"`
}
