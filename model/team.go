package model

type Team struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar"`
}

type Member struct {
	Login string `json:"login"`
}
