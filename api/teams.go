package api

import (
	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/router/middleware/session"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// GetTeams gets the list of user teams.
func GetTeams(c *gin.Context) {
	user := session.User(c)
	teams, err := cache.GetTeams(c, user)
	if err != nil {
		logrus.Errorf("Error getting teams for user %s. %s", user.Login, err)
		c.String(500, "Error getting team list")
		return
	}
	teams = append(teams, &model.Team{
		Login:  user.Login,
		Avatar: user.Avatar,
	})
	c.JSON(200, teams)
}
