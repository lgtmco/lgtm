package model

type User struct {
	ID     int64  `json:"id"      meddler:"user_id,pk"`
	Login  string `json:"login"   meddler:"user_login"`
	Email  string `json:"email"   meddler:"user_email"`
	Token  string `json:"-"       meddler:"user_token"`
	Avatar string `json:"avatar"  meddler:"user_avatar"`
	Secret string `json:"-"       meddler:"user_secret"`
}
