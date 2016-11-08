package api

import (
	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	"github.com/lgtmco/lgtm/router/middleware/session"
	"github.com/lgtmco/lgtm/store"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// GetMaintainer gets the MAINTAINER configuration file.
func GetMaintainer(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
		user  = session.User(c)
	)
	repo, err := store.GetRepoOwnerName(c, owner, name)
	if err != nil {
		log.Errorf("Error getting repository %s. %s", name, err)
		c.AbortWithStatus(404)
		return
	}

	rcfile, _ := remote.GetContents(c, user, repo, ".lgtm")
	config, err := model.ParseConfig(rcfile)
	if err != nil {
		log.Errorf("Error parsing .lgtm file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing .lgtm file. %s.", err)
		return
	}

	file, err := remote.GetContents(c, user, repo, "MAINTAINERS")
	if err != nil {
		log.Debugf("no MAINTAINERS file for %s. Checking for team members.", repo.Slug)
		members, merr := cache.GetMembers(c, user, repo.Owner, config.Team)
		if merr != nil {
			log.Errorf("Error getting repository %s. %s", repo.Slug, err)
			log.Errorf("Error getting org members %s. %s", repo.Owner, merr)
			c.String(404, "MAINTAINERS file not found. %s", err)
			return
		} else {
			log.Printf("found %v members", len(members))
			for _, member := range members {
				file = append(file, member.Login...)
				file = append(file, '\n')
			}
		}
	}

	maintainer, err := model.ParseMaintainer(file)
	if err != nil {
		log.Errorf("Error parsing MAINTAINERS file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing MAINTAINERS file. %s.", err)
		return
	}
	c.JSON(200, maintainer)
}

// GetMaintainer gets the MAINTAINER configuration file and returns
// a subset of the file with members belonging to the specified organization.
func GetMaintainerOrg(c *gin.Context) {
	var (
		owner = c.Param("owner")
		name  = c.Param("repo")
		team  = c.Param("org")
		user  = session.User(c)
	)
	repo, err := store.GetRepoOwnerName(c, owner, name)
	if err != nil {
		log.Errorf("Error getting repository %s. %s", name, err)
		c.AbortWithStatus(404)
		return
	}
	file, err := remote.GetContents(c, user, repo, "MAINTAINERS")
	if err != nil {
		log.Errorf("Error getting repository %s. %s", repo.Slug, err)
		c.String(404, "MAINTAINERS file not found. %s", err)
		return
	}
	maintainer, err := model.ParseMaintainer(file)
	if err != nil {
		log.Errorf("Error parsing MAINTAINERS file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing MAINTAINERS file. %s.", err)
		return
	}
	subset, err := model.FromOrg(maintainer, team)
	if err != nil {
		log.Errorf("Error getting subset of MAINTAINERS file for %s/%s. %s", repo.Slug, team, err)
		c.String(500, "Error getting subset of MAINTAINERS file. %s.", err)
		return
	}
	c.JSON(200, subset)
}
