package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lgtmco/lgtm/model"

	cache "github.com/lgtmco/lgtm/cache/mock"
	remote "github.com/lgtmco/lgtm/remote/mock"

	"github.com/Sirupsen/logrus"
	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func TestTeams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logrus.SetOutput(ioutil.Discard)

	g := goblin.Goblin(t)

	g.Describe("Team endpoint", func() {
		g.It("Should return the team list", func() {
			cache := new(cache.Cache)
			cache.On("Get", "teams:octocat").Return(fakeTeams, nil).Once()

			e := gin.New()
			e.NoRoute(GetTeams)
			e.Use(func(c *gin.Context) {
				c.Set("user", fakeUser)
				c.Set("cache", cache)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/", nil)
			e.ServeHTTP(w, r)

			// the user is appended to the team list so we retrieve a full list of
			// accounts to which the user has access.
			teams := append(fakeTeams, &model.Team{
				Login: fakeUser.Login,
			})

			want, _ := json.Marshal(teams)
			got := strings.TrimSpace(w.Body.String())
			g.Assert(got).Equal(string(want))
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return a 500 error", func() {
			remote := new(remote.Remote)
			cache := new(cache.Cache)
			cache.On("Get", "teams:octocat").Return(nil, fmt.Errorf("Not Found")).Once()
			remote.On("GetTeams", fakeUser).Return(nil, fmt.Errorf("Not Found")).Once()

			e := gin.New()
			e.NoRoute(GetTeams)
			e.Use(func(c *gin.Context) {
				c.Set("user", fakeUser)
				c.Set("cache", cache)
				c.Set("remote", remote)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/", nil)
			e.ServeHTTP(w, r)

			got := strings.TrimSpace(w.Body.String())
			g.Assert(got).Equal("Error getting team list")
			g.Assert(w.Code).Equal(500)
		})
	})
}
