package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgtmco/lgtm/api"
	"github.com/lgtmco/lgtm/router/middleware/access"
	"github.com/lgtmco/lgtm/router/middleware/header"
	"github.com/lgtmco/lgtm/router/middleware/session"
	"github.com/lgtmco/lgtm/web"
	"github.com/lgtmco/lgtm/web/static"
	"github.com/lgtmco/lgtm/web/template"
)

func Load(middleware ...gin.HandlerFunc) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	e.SetHTMLTemplate(template.Template())
	e.StaticFS("/static", static.FileSystem())

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)
	e.Use(session.SetUser)

	e.GET("/api/user", session.UserMust, api.GetUser)
	e.GET("/api/user/teams", session.UserMust, api.GetTeams)
	e.GET("/api/user/repos", session.UserMust, api.GetRepos)
	e.GET("/api/repos/:owner/:repo", session.UserMust, access.RepoPull, api.GetRepo)
	e.POST("/api/repos/:owner/:repo", session.UserMust, access.RepoAdmin, api.PostRepo)
	e.DELETE("/api/repos/:owner/:repo", session.UserMust, access.RepoAdmin, api.DeleteRepo)
	e.GET("/api/repos/:owner/:repo/maintainers", session.UserMust, access.RepoPull, api.GetMaintainer)
	e.GET("/api/repos/:owner/:repo/maintainers/:org", session.UserMust, access.RepoPull, api.GetMaintainerOrg)
	e.GET("/api/repos/:owner/:repo/tags", session.UserMust, access.RepoPull, api.GetTags)

	e.POST("/hook", web.Hook)
	e.GET("/login", web.Login)
	e.POST("/login", web.LoginToken)
	e.GET("/logout", web.Logout)
	e.NoRoute(web.Index)

	return e
}
