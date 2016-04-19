package model

import "testing"

func TestParseMaintainer(t *testing.T) {
	var files = []string{maintainerFile, maintainerFileEmail, maintainerFileSimple, maintainerFileMixed, maintainerFileToml}
	for _, file := range files {
		parsed, err := ParseMaintainerStr(file)
		if err != nil {
			t.Error(err)
			return
		}
		if len(parsed.People) != len(people) {
			t.Errorf("Wanted %d maintainers, got %d", len(people), len(parsed.People))
			return
		}
		for _, want := range people {
			got, ok := parsed.People[want.Login]
			if !ok {
				t.Errorf("Wanted user %s in file", want.Login)
			} else if want.Login != got.Login {
				t.Errorf("Wanted login %s, got %s", want.Login, got.Login)
			}
		}
	}
}

var people = []Person{
	{Login: "bradrydzewski"},
	{Login: "mattnorris"},
}

var maintainerFile = `
Brad Rydzewski <brad.rydzewski@mail.com> (@bradrydzewski)
Matt Norris <matt.norris@mail.com> (@mattnorris)
`

var maintainerFileEmail = `
bradrydzewski <brad.rydzewski@mail.com>
mattnorris <matt.norris@mail.com>
`

// simple format with usernames only. includes
// spaces and comments.
var maintainerFileSimple = `
bradrydzewski
mattnorris`

// simple format with usernames only. includes
// spaces and comments.
var maintainerFileMixed = `
bradrydzewski
Matt Norris <matt.norris@mail.com> (@mattnorris)
`

// advanced toml format for the maintainers file.
var maintainerFileToml = `
[org]
	[org.core]
		people = [
			"mattnorris",
			"bradrydzewski",
		]

[people]

	[people.bradrydzewski]
	name = "Brad Rydzewski"
	email = "brad.rydzewski@mail.com"
	login = "bradrydzewski"

	[people.mattnorris]
	name = "Matt Norris"
	email = "matt.norris@mail.com"
	login = "mattnorris"
`
