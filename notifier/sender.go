package notifier

import "golang.org/x/net/context"

//go:generate mockery -name Sender -output mock -case=underscore

// Sender defines a notification provider that is capable of sending out
// notifications to a list of maintainers or reviewers. An example provider
// might be a Slack or GitHub bot.
type Sender interface {
	Send(*Notification) error
}

// Send sends a notification to the list of maintainers indicating a commit is
// ready for their review and possible approval.
func Send(c context.Context, n *Notification) error {
	return FromContext(c).Send(n)
}
