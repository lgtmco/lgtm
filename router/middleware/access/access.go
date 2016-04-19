package access

import (
	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/router/middleware/session"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func RepoAdmin(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
		user  = session.User(c)
	)

	perm, err := cache.GetPerm(c, user, owner, name)
	if err != nil {
		log.Errorf("Cannot find repository %s/%s. %s", owner, name, err)
		c.String(404, "Not Found")
		c.Abort()
		return
	}
	if !perm.Admin {
		log.Errorf("User %s does not have Admin access to repository %s/%s", user.Login, owner, name)
		c.String(403, "Insufficient privileges")
		c.Abort()
		return
	}
	log.Debugf("User %s granted Admin access to %s/%s", user.Login, owner, name)
	c.Next()
}

func RepoPull(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
		user  = session.User(c)
	)

	perm, err := cache.GetPerm(c, user, owner, name)
	if err != nil {
		log.Errorf("Cannot find repository %s/%s. %s", owner, name, err)
		c.String(404, "Not Found")
		c.Abort()
		return
	}
	if !perm.Pull {
		log.Errorf("User %s does not have Pull access to repository %s/%s", user.Login, owner, name)
		c.String(404, "Not Found")
		c.Abort()
		return
	}
	log.Debugf("User %s granted Pull access to %s/%s", user.Login, owner, name)
	c.Next()
}
