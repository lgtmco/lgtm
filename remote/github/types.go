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
