package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/lgtmco/lgtm/model"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func TestUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logrus.SetOutput(ioutil.Discard)

	g := goblin.Goblin(t)

	g.Describe("User endpoint", func() {
		g.It("Should return the authenticated user", func() {

			e := gin.New()
			e.NoRoute(GetUser)
			e.Use(func(c *gin.Context) {
				c.Set("user", fakeUser)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/", nil)
			e.ServeHTTP(w, r)

			want, _ := json.Marshal(fakeUser)
			got := strings.TrimSpace(w.Body.String())
			g.Assert(got).Equal(string(want))
			g.Assert(w.Code).Equal(200)
		})
	})
}

var (
	fakeUser  = &model.User{Login: "octocat"}
	fakeTeams = []*model.Team{
		{Login: "drone"},
		{Login: "docker"},
	}
)
