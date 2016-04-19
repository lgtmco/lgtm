package cache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bradrydzewski/lgtm/model"
	"github.com/bradrydzewski/lgtm/remote"
	"github.com/bradrydzewski/lgtm/remote/mock"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func TestHelper(t *testing.T) {

	g := goblin.Goblin(t)

	g.Describe("Cache helpers", func() {

		var c *gin.Context
		var r *mock.Remote

		g.BeforeEach(func() {
			c = new(gin.Context)
			ToContext(c, Default())

			r = new(mock.Remote)
			remote.ToContext(c, r)
		})

		g.It("Should get permissions from remote", func() {
			r.On("GetPerm", fakeUser, fakeRepo.Owner, fakeRepo.Name).Return(fakePerm, nil).Once()
			p, err := GetPerm(c, fakeUser, fakeRepo.Owner, fakeRepo.Name)
			g.Assert(p).Equal(fakePerm)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get permissions from cache", func() {
			key := fmt.Sprintf("perms:%s:%s/%s",
				fakeUser.Login,
				fakeRepo.Owner,
				fakeRepo.Name,
			)

			Set(c, key, fakePerm)
			r.On("GetPerm", fakeUser, fakeRepo.Owner, fakeRepo.Name).Return(nil, fakeErr).Once()
			p, err := GetPerm(c, fakeUser, fakeRepo.Owner, fakeRepo.Name)
			g.Assert(p).Equal(fakePerm)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get permissions error", func() {
			r.On("GetPerm", fakeUser, fakeRepo.Owner, fakeRepo.Name).Return(nil, fakeErr).Once()
			p, err := GetPerm(c, fakeUser, fakeRepo.Owner, fakeRepo.Name)
			g.Assert(p == nil).IsTrue()
			g.Assert(err).Equal(fakeErr)
		})

		g.It("Should set and get repos", func() {

			r.On("GetRepos", fakeUser).Return(fakeRepos, nil).Once()
			p, err := GetRepos(c, fakeUser)
			g.Assert(p).Equal(fakeRepos)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get repos", func() {
			key := fmt.Sprintf("repos:%s",
				fakeUser.Login,
			)

			Set(c, key, fakeRepos)
			r.On("GetRepos", fakeUser).Return(nil, fakeErr).Once()
			p, err := GetRepos(c, fakeUser)
			g.Assert(p).Equal(fakeRepos)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get repos error", func() {
			r.On("GetRepos", fakeUser).Return(nil, fakeErr).Once()
			p, err := GetRepos(c, fakeUser)
			g.Assert(p == nil).IsTrue()
			g.Assert(err).Equal(fakeErr)
		})

		g.It("Should set and get teams", func() {
			r.On("GetTeams", fakeUser).Return(fakeTeams, nil).Once()
			p, err := GetTeams(c, fakeUser)
			g.Assert(p).Equal(fakeTeams)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get teams", func() {
			key := fmt.Sprintf("teams:%s",
				fakeUser.Login,
			)

			Set(c, key, fakeTeams)
			r.On("GetTeams", fakeUser).Return(nil, fakeErr).Once()
			p, err := GetTeams(c, fakeUser)
			g.Assert(p).Equal(fakeTeams)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get team error", func() {
			r.On("GetTeams", fakeUser).Return(nil, fakeErr).Once()
			p, err := GetTeams(c, fakeUser)
			g.Assert(p == nil).IsTrue()
			g.Assert(err).Equal(fakeErr)
		})

		g.It("Should set and get members", func() {
			r.On("GetMembers", fakeUser, "drone").Return(fakeMembers, nil).Once()
			p, err := GetMembers(c, fakeUser, "drone")
			g.Assert(p).Equal(fakeMembers)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get members", func() {
			key := "members:drone"

			Set(c, key, fakeMembers)
			r.On("GetMembers", fakeUser, "drone").Return(nil, fakeErr).Once()
			p, err := GetMembers(c, fakeUser, "drone")
			g.Assert(p).Equal(fakeMembers)
			g.Assert(err).Equal(nil)
		})

		g.It("Should get member error", func() {
			r.On("GetMembers", fakeUser, "drone").Return(nil, fakeErr).Once()
			p, err := GetMembers(c, fakeUser, "drone")
			g.Assert(p == nil).IsTrue()
			g.Assert(err).Equal(fakeErr)
		})
	})
}

var (
	fakeErr   = errors.New("Not Found")
	fakeUser  = &model.User{Login: "octocat"}
	fakePerm  = &model.Perm{Pull: true, Push: true, Admin: true}
	fakeRepo  = &model.Repo{Owner: "octocat", Name: "Hello-World"}
	fakeRepos = []*model.Repo{
		{Owner: "octocat", Name: "Hello-World"},
		{Owner: "octocat", Name: "hello-world"},
		{Owner: "octocat", Name: "Spoon-Knife"},
	}
	fakeTeams = []*model.Team{
		{Login: "drone"},
		{Login: "docker"},
	}
	fakeMembers = []*model.Member{
		{Login: "octocat"},
	}
)
