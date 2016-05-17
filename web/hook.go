package web

import (
	"github.com/lgtmco/lgtm/cache"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/lgtmco/lgtm/store"
)

func Hook(c *gin.Context) {
	hook, err := remote.GetHook(c, c.Request)
	if err != nil {
		log.Errorf("Error parsing hook. %s", err)
		c.String(500, "Error parsing hook. %s", err)
		return
	}
	if hook != nil {
		processCommentHook(c, hook)
	}

	statusHook, err := remote.GetStatusHook(c, c.Request)
	if err != nil {
		log.Errorf("Error parsing status hook. %s", err)
		c.String(500, "Error parsing status hook. %s", err)
		return
	}
	if statusHook !=  nil {
		processStatusHook(c, statusHook)
	}

	prHook, err := remote.GetPRHook(c, c.Request)
	if prHook != nil {
		processPRHook(c, prHook)
	}

	if err != nil {
		log.Errorf("Error parsing pull request hook. %s", err)
		c.String(500, "Error parsing pull request hook. %s", err)
		return
	}

	pushHook, err := remote.GetPushHook(c, c.Request)
	if err != nil {
		log.Errorf("Error parsing push hook. %s", err)
		c.String(500, "Error parsing push hook. %s", err)
		return
	}
	if pushHook != nil {
		processPushHook(c, pushHook)
	}

	if hook == nil && statusHook == nil && prHook == nil && pushHook == nil {
		c.String(200, "pong")
		return
	}
}

func processPRHook(c *gin.Context, prHook *model.PRHook) {
	repo, user, err := getRepoAndUser(c, prHook.Repo.Slug)
	if err != nil {
		return
	}
	err = remote.SetStatus(c, user, repo, prHook.Number, false)
	if err != nil {
		log.Errorf("Error setting status. %s", err)
		c.String(500, "Error setting status. %s", err)
		return
	}

	c.IndentedJSON(200, gin.H{
		"number":   prHook.Number,
		"approved":    false,
	})
}

func getRepoAndUser(c *gin.Context, slug string) (*model.Repo, *model.User, error){
	repo, err := store.GetRepoSlug(c, slug)
	if err != nil {
		log.Errorf("Error getting repository %s. %s", slug, err)
		c.String(404, "Repository not found.")
		return nil, nil, err
	}
	user, err := store.GetUser(c, repo.UserID)
	if err != nil {
		log.Errorf("Error getting repository owner %s. %s", repo.Slug, err)
		c.String(404, "Repository owner not found.")
		return nil, nil, err
	}
	return repo, user, err
}

func getConfigAndMaintainers(c *gin.Context, user *model.User, repo *model.Repo) (*model.Config, *model.Maintainer, error) {
	rcfile, _ := remote.GetContents(c, user, repo, ".lgtm")
	config, err := model.ParseConfig(rcfile)
	if err != nil {
		log.Errorf("Error parsing .lgtm file for %s. %s", repo.Slug, err)
		c.String(500, "Error parsing .lgtm file. %s.", err)
		return nil, nil, err
	}

	// THIS IS COMPLETELY DUPLICATED IN THE API SECTION. NOT IDEAL
	file, err := remote.GetContents(c, user, repo, "MAINTAINERS")
	if err != nil {
		log.Debugf("no MAINTAINERS file for %s. Checking for team members.", repo.Slug)
		members, merr := cache.GetMembers(c, user, repo.Owner)
		if merr != nil {
			log.Errorf("Error getting repository %s. %s", repo.Slug, err)
			log.Errorf("Error getting org members %s. %s", repo.Owner, merr)
			c.String(404, "MAINTAINERS file not found. %s", err)
			return nil, nil, err
		} else {
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
		return nil, nil, err
	}
	return config, maintainer, nil
}

func getComments(c *gin.Context, user *model.User, repo *model.Repo, num int) ([]*model.Comment, error) {
	comments, err := remote.GetComments(c, user, repo, num)
	if err != nil {
		log.Errorf("Error retrieving comments for %s pr %d. %s", repo.Slug, num, err)
		c.String(500, "Error retrieving comments. %s.", err)
		return nil, err
	}
	return comments, nil
}

