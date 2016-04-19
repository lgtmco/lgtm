package mock

import "github.com/lgtmco/lgtm/notifier"
import "github.com/stretchr/testify/mock"

type Sender struct {
	mock.Mock
}

func (_m *Sender) Send(_a0 *notifier.Notification) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*notifier.Notification) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
