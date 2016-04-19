package notifier

import "golang.org/x/net/context"

const key = "sender"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Sender associated with this context.
func FromContext(c context.Context) Sender {
	return c.Value(key).(Sender)
}

// ToContext adds the Sender to this context if it supports
// the Setter interface.
func ToContext(c Setter, s Sender) {
	c.Set(key, s)
}
