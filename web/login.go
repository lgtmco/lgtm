package web

import (
	"net/http"
	"time"

	"github.com/bradrydzewski/lgtm/model"
	"github.com/bradrydzewski/lgtm/remote"
	"github.com/bradrydzewski/lgtm/shared/httputil"
	"github.com/bradrydzewski/lgtm/shared/token"
	"github.com/bradrydzewski/lgtm/store"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Login attempts to authorize a user via GitHub oauth2. If the user does not
// yet exist, and new account is created. Upon successful login the user is
// redirected to the main screen.
func Login(c *gin.Context) {
	// render the error page if the login fails. Without this block
	// we would encounter an infinite number of redirects.
	if err := c.Query("error"); len(err) != 0 {
		c.HTML(500, "error.html", gin.H{"error": err})
		return
	}

	// when dealing with redirects we may need
	// to adjust the content type. I cannot, however,
	// rememver why, so need to revisit this line.
	c.Writer.Header().Del("Content-Type")

	tmpuser, err := remote.GetUser(c, c.Writer, c.Request)
	if err != nil {
		log.Errorf("cannot authenticate user. %s", err)
		c.Redirect(303, "/login?error=oauth_error")
		return
	}
	// this will happen when the user is redirected by
	// the remote provide as part of the oauth dance.
	if tmpuser == nil {
		return
	}

	// get the user from the database
	u, err := store.GetUserLogin(c, tmpuser.Login)
	if err != nil {

		// create the user account
		u = &model.User{}
		u.Login = tmpuser.Login
		u.Token = tmpuser.Token
		u.Avatar = tmpuser.Avatar
		u.Secret = model.Rand()

		// insert the user into the database
		if err := store.CreateUser(c, u); err != nil {
			log.Errorf("cannot insert %s. %s", u.Login, err)
			c.Redirect(303, "/login?error=internal_error")
			return
		}
	}

	// update the user meta data and authorization
	// data and cache in the datastore.
	u.Token = tmpuser.Token
	u.Avatar = tmpuser.Avatar

	if err := store.UpdateUser(c, u); err != nil {
		log.Errorf("cannot update %s. %s", u.Login, err)
		c.Redirect(303, "/login?error=internal_error")
		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	token := token.New(token.SessToken, u.Login)
	tokenstr, err := token.SignExpires(u.Secret, exp)
	if err != nil {
		log.Errorf("cannot create token for %s. %s", u.Login, err)
		c.Redirect(303, "/login?error=internal_error")
		return
	}

	httputil.SetCookie(c.Writer, c.Request, "user_sess", tokenstr)
	c.Redirect(303, "/")
}

// LoginToken authenticates a user with their GitHub token and
// returns an LGTM API token in the response.
func LoginToken(c *gin.Context) {
	access := c.Query("access_token")
	login, err := remote.GetUserToken(c, access)
	if err != nil {
		c.String(403, "Unable to authenticate user. %s", err)
		return
	}
	user, err := store.GetUserLogin(c, login)
	if err != nil {
		c.String(404, "Unable to authenticate user %s. Not registered.", user.Login)
		return
	}
	exp := time.Now().Add(time.Hour * 72).Unix()
	token := token.New(token.UserToken, user.Login)
	tokenstr, err := token.SignExpires(user.Secret, exp)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, &tokenPayload{
		Access:  tokenstr,
		Expires: exp - time.Now().Unix(),
	})
}

type tokenPayload struct {
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	Expires int64  `json:"expires_in,omitempty"`
}

// Logout terminates the session for the currently authenticated user,
// deleting all session cookies, and redirecting back to the main page.
func Logout(c *gin.Context) {
	httputil.DelCookie(c.Writer, c.Request, "user_sess")
	c.HTML(200, "logout.html", gin.H{})
}
