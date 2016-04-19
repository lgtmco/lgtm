package datastore

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/lgtmco/lgtm/model"
)

func Test_userstore(t *testing.T) {
	db := openTest()
	defer db.Close()
	s := From(db)

	g := goblin.Goblin(t)
	g.Describe("User", func() {

		// before each test be sure to purge the package
		// table data from the database.
		g.BeforeEach(func() {
			db.Exec("DELETE FROM users")
		})

		g.It("Should Update a User", func() {
			user := model.User{
				Login: "joe",
				Email: "foo@bar.com",
				Token: "e42080dddf012c718e476da161d21ad5",
			}
			err1 := s.CreateUser(&user)
			err2 := s.UpdateUser(&user)
			getuser, err3 := s.GetUser(user.ID)
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 == nil).IsTrue()
			g.Assert(user.ID).Equal(getuser.ID)
		})

		g.It("Should Add a new User", func() {
			user := model.User{
				Login: "joe",
				Email: "foo@bar.com",
				Token: "e42080dddf012c718e476da161d21ad5",
			}
			err := s.CreateUser(&user)
			g.Assert(err == nil).IsTrue()
			g.Assert(user.ID != 0).IsTrue()
		})

		g.It("Should Get a User", func() {
			user := model.User{
				Login:  "joe",
				Token:  "f0b461ca586c27872b43a0685cbc2847",
				Secret: "976f22a5eef7caacb7e678d6c52f49b1",
				Email:  "foo@bar.com",
				Avatar: "b9015b0857e16ac4d94a0ffd9a0b79c8",
			}

			s.CreateUser(&user)
			getuser, err := s.GetUser(user.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(user.ID).Equal(getuser.ID)
			g.Assert(user.Login).Equal(getuser.Login)
			g.Assert(user.Token).Equal(getuser.Token)
			g.Assert(user.Secret).Equal(getuser.Secret)
			g.Assert(user.Email).Equal(getuser.Email)
			g.Assert(user.Avatar).Equal(getuser.Avatar)
		})

		g.It("Should Get a User By Login", func() {
			user := model.User{
				Login: "joe",
				Email: "foo@bar.com",
				Token: "e42080dddf012c718e476da161d21ad5",
			}
			s.CreateUser(&user)
			getuser, err := s.GetUserLogin(user.Login)
			g.Assert(err == nil).IsTrue()
			g.Assert(user.ID).Equal(getuser.ID)
			g.Assert(user.Login).Equal(getuser.Login)
		})

		g.It("Should Enforce Unique User Login", func() {
			user1 := model.User{
				Login: "joe",
				Email: "foo@bar.com",
				Token: "e42080dddf012c718e476da161d21ad5",
			}
			user2 := model.User{
				Login: "joe",
				Email: "foo@bar.com",
				Token: "ab20g0ddaf012c744e136da16aa21ad9",
			}
			err1 := s.CreateUser(&user1)
			err2 := s.CreateUser(&user2)
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsFalse()
		})

		g.It("Should Del a User", func() {
			user := model.User{
				Login: "joe",
				Email: "foo@bar.com",
				Token: "e42080dddf012c718e476da161d21ad5",
			}
			s.CreateUser(&user)
			_, err1 := s.GetUser(user.ID)
			err2 := s.DeleteUser(&user)
			_, err3 := s.GetUser(user.ID)
			g.Assert(err1 == nil).IsTrue()
			g.Assert(err2 == nil).IsTrue()
			g.Assert(err3 == nil).IsFalse()
		})
	})
}
