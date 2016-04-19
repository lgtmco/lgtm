package datastore

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/lgtmco/lgtm/model"
)

func Test_repostore(t *testing.T) {
	db := openTest()
	defer db.Close()

	s := From(db)
	g := goblin.Goblin(t)
	g.Describe("Repo", func() {

		// before each test be sure to purge the package
		// table data from the database.
		g.BeforeEach(func() {
			db.Exec("DELETE FROM repos")
			db.Exec("DELETE FROM users")
		})

		g.It("Should Set a Repo", func() {
			repo := model.Repo{
				UserID: 1,
				Slug:   "bradrydzewski/drone",
				Owner:  "bradrydzewski",
				Name:   "drone",
			}
			err1 := s.CreateRepo(&repo)
			err2 := s.UpdateRepo(&repo)
			getrepo, err3 := s.GetRepo(repo.ID)
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 == nil).IsTrue()
			g.Assert(repo.ID).Equal(getrepo.ID)
		})

		g.It("Should Add a Repo", func() {
			repo := model.Repo{
				UserID: 1,
				Slug:   "bradrydzewski/drone",
				Owner:  "bradrydzewski",
				Name:   "drone",
			}
			err := s.CreateRepo(&repo)
			g.Assert(err == nil).IsTrue()
			g.Assert(repo.ID != 0).IsTrue()
		})

		g.It("Should Get a Repo by ID", func() {
			repo := model.Repo{
				UserID:  1,
				Slug:    "bradrydzewski/drone",
				Owner:   "bradrydzewski",
				Name:    "drone",
				Link:    "https://github.com/octocat/hello-world",
				Private: true,
			}
			s.CreateRepo(&repo)
			getrepo, err := s.GetRepo(repo.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(repo.ID).Equal(getrepo.ID)
			g.Assert(repo.UserID).Equal(getrepo.UserID)
			g.Assert(repo.Owner).Equal(getrepo.Owner)
			g.Assert(repo.Name).Equal(getrepo.Name)
			g.Assert(repo.Private).Equal(getrepo.Private)
			g.Assert(repo.Link).Equal(getrepo.Link)
		})

		g.It("Should Get a Repo by Slug", func() {
			repo := model.Repo{
				UserID: 1,
				Slug:   "bradrydzewski/drone",
				Owner:  "bradrydzewski",
				Name:   "drone",
			}
			s.CreateRepo(&repo)
			getrepo, err := s.GetRepoSlug(repo.Slug)
			g.Assert(err == nil).IsTrue()
			g.Assert(repo.ID).Equal(getrepo.ID)
			g.Assert(repo.UserID).Equal(getrepo.UserID)
			g.Assert(repo.Owner).Equal(getrepo.Owner)
			g.Assert(repo.Name).Equal(getrepo.Name)
		})

		g.It("Should Get a Multiple Repos", func() {
			repo1 := &model.Repo{
				UserID: 1,
				Owner:  "foo",
				Name:   "bar",
				Slug:   "foo/bar",
			}
			repo2 := &model.Repo{
				UserID: 2,
				Owner:  "octocat",
				Name:   "fork-knife",
				Slug:   "octocat/fork-knife",
			}
			repo3 := &model.Repo{
				UserID: 2,
				Owner:  "octocat",
				Name:   "hello-world",
				Slug:   "octocat/hello-world",
			}
			s.CreateRepo(repo1)
			s.CreateRepo(repo2)
			s.CreateRepo(repo3)

			repos, err := s.GetRepoMulti("octocat/fork-knife", "octocat/hello-world")
			g.Assert(err == nil).IsTrue()
			g.Assert(len(repos)).Equal(2)
			g.Assert(repos[0].ID).Equal(repo2.ID)
			g.Assert(repos[1].ID).Equal(repo3.ID)
		})

		g.It("Should Delete a Repo", func() {
			repo := model.Repo{
				UserID: 1,
				Slug:   "bradrydzewski/drone",
				Owner:  "bradrydzewski",
				Name:   "drone",
			}
			s.CreateRepo(&repo)
			_, err1 := s.GetRepo(repo.ID)
			err2 := s.DeleteRepo(&repo)
			_, err3 := s.GetRepo(repo.ID)
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 == nil).IsFalse()
		})

		g.It("Should Enforce Unique Repo Name", func() {
			repo1 := model.Repo{
				UserID: 1,
				Slug:   "bradrydzewski/drone",
				Owner:  "bradrydzewski",
				Name:   "drone",
			}
			repo2 := model.Repo{
				UserID: 2,
				Slug:   "bradrydzewski/drone",
				Owner:  "bradrydzewski",
				Name:   "drone",
			}
			err1 := s.CreateRepo(&repo1)
			err2 := s.CreateRepo(&repo2)
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsFalse()
		})
	})
}
