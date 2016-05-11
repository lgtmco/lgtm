package org

import (
	"testing"
	"github.com/lgtmco/lgtm/model"
)

var maintainerToml = `
[people]
	[people.bob]
		login = "bob"
	[people.fred]
		login = "fred"
	[people.jon]
		login = "jon"
	[people.ralph]
		login = "ralph"
	[people.george]
		login = "george"

[org]
[org.cap]
people = [
            "bob",
            "fred",
            "jon"
]

[org.iron]
people = [
            "ralph",
            "george"
]
`

func TestOrg(t *testing.T) {
	config := &model.Config {
		Pattern: `(?i)LGTM\s*(\S*)`,
		SelfApprovalOff: true,
	}
	m, err := model.ParseMaintainerStr(maintainerToml)
	if err != nil {
		t.Fatal(err)
	}
	issue := &model.Issue{
		Author: "jon",
	}
	comments := []*model.Comment {
		{
			Body: "lgtm",
			Author: "bob",
		},
		{
			Body: "lgtm",
			Author: "qwerty",
		},
		{
			Body: "not an approval",
			Author: "ralph",
		},
		{
			Body: "lgtm",
			Author: "george",
		},
		{
			Body: "lgtm",
			Author: "ralph",
		},
	}
	people := []string{}
	Org(config, m, issue, comments, func(m *model.Maintainer, c *model.Comment) {
		people = append(people, c.Author)
	})
	if len(people) != 1 {
		t.Errorf("Expected one person, had %d", len(people))
	}
	if people[0] != "george" {
		t.Errorf("Expected george, had %s", people[0])
	}
}