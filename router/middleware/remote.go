package middleware

import (
	"strings"

	"github.com/lgtmco/lgtm/remote/github"

	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
)

const (
	DefaultURL   = "https://github.com"
	DefaultAPI   = "https://api.github.com/"
	DefaultScope = "user:email,read:org,public_repo"
)

var (
	server = envflag.String("GITHUB_URL", DefaultURL, "")
	client = envflag.String("GITHUB_CLIENT", "", "")
	secret = envflag.String("GITHUB_SECRET", "", "")
	scope  = envflag.String("GITHUB_SCOPE", DefaultScope, "")
)

func Remote() gin.HandlerFunc {
	remote := &github.Github{
		API:    DefaultAPI,
		URL:    *server,
		Client: *client,
		Secret: *secret,
		Scopes: strings.Split(*scope, ","),
	}
	if remote.URL != DefaultURL {
		remote.URL = strings.TrimSuffix(remote.URL, "/")
		remote.API = remote.URL + "/api/v3/"
	}
	return func(c *gin.Context) {
		c.Set("remote", remote)
		c.Next()
	}
}
