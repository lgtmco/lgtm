package session

import (
	"net/http"

	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/shared/token"
	"github.com/lgtmco/lgtm/store"

	"github.com/gin-gonic/gin"
)

func User(c *gin.Context) *model.User {
	v, ok := c.Get("user")
	if !ok {
		return nil
	}
	u, ok := v.(*model.User)
	if !ok {
		return nil
	}
	return u
}

func UserMust(c *gin.Context) {
	user := User(c)
	switch {
	case user == nil:
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
	default:
		c.Next()
	}
}

func SetUser(c *gin.Context) {
	var user *model.User

	// authenticates the user via an authentication cookie
	// or an auth token.
	t, err := token.ParseRequest(c.Request, func(t *token.Token) (string, error) {
		var err error
		user, err = store.GetUserLogin(c, t.Text)
		return user.Secret, err
	})

	if err == nil {
		c.Set("user", user)

		// if this is a session token (ie not the API token)
		// this means the user is accessing with a web browser,
		// so we should implement CSRF protection measures.
		if t.Kind == token.SessToken {
			err = token.CheckCsrf(c.Request, func(t *token.Token) (string, error) {
				return user.Secret, nil
			})
			// if csrf token validation fails, exit immediately
			// with a not authorized error.
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
	}
	c.Next()
}
