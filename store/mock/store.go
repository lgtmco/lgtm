package mock

import "github.com/bradrydzewski/lgtm/store"
import "github.com/stretchr/testify/mock"

type Store struct {
	mock.Mock
}

func (_m *Store) Users() store.UserStore {
	ret := _m.Called()

	var r0 store.UserStore
	if rf, ok := ret.Get(0).(func() store.UserStore); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(store.UserStore)
	}

	return r0
}
func (_m *Store) Repos() store.RepoStore {
	ret := _m.Called()

	var r0 store.RepoStore
	if rf, ok := ret.Get(0).(func() store.RepoStore); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(store.RepoStore)
	}

	return r0
}
