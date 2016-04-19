package api

import (
	"github.com/bradrydzewski/lgtm/cache"
	"github.com/bradrydzewski/lgtm/model"
	"github.com/bradrydzewski/lgtm/router/middleware/session"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// GetTeams gets the list of user teams.
func GetTeams(c *gin.Context) {
	user := session.User(c)
	teams, err := cache.GetTeams(c, user)
	if err != nil {
		log.Errorf("Error getting team list. %s", err)
		c.String(500, "Error getting team list")
		return
	}
	teams = append(teams, &model.Team{
		Login:  user.Login,
		Avatar: user.Avatar,
	})
	c.JSON(200, teams)
}
