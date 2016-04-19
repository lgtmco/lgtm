package model

type Repo struct {
	ID      int64  `json:"id,omitempty"       meddler:"repo_id,pk"`
	UserID  int64  `json:"-"                  meddler:"repo_user_id"`
	Owner   string `json:"owner"              meddler:"repo_owner"`
	Name    string `json:"name"               meddler:"repo_name"`
	Slug    string `json:"slug"               meddler:"repo_slug"`
	Link    string `json:"link_url"           meddler:"repo_link"`
	Private bool   `json:"private"            meddler:"repo_private"`
	Secret  string `json:"-"                  meddler:"repo_secret"`
}

type Perm struct {
	Pull  bool
	Push  bool
	Admin bool
}
