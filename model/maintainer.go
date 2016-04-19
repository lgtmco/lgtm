package model

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// Person represets an individual in the MAINTAINERS file.
type Person struct {
	Name  string `json:"name"  toml:"name"`
	Email string `json:"email" toml:"email"`
	Login string `json:"login" toml:"login"`
}

// Org represents a group, team or subset of users.
type Org struct {
	People []string `json:"people" toml:"people"`
}

// Maintainer represents a MAINTAINERS file.
type Maintainer struct {
	People map[string]*Person `json:"people"    toml:"people"`
	Org    map[string]*Org    `json:"org"       toml:"org"`
}

// ParseMaintainer parses a projects MAINTAINERS file and returns
// the list of maintainers.
func ParseMaintainer(data []byte) (*Maintainer, error) {
	return ParseMaintainerStr(string(data))
}

// ParseMaintainerStr parses a projects MAINTAINERS file in string
// format and returns the list of maintainers.
func ParseMaintainerStr(data string) (*Maintainer, error) {
	m, err := parseMaintainerToml(data)
	if err != nil {
		m, err = parseMaintainerText(data)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// FromOrg returns a new Maintainer file with a subset of people
// that are part of the specified org.
func FromOrg(from *Maintainer, name string) (*Maintainer, error) {
	m := new(Maintainer)
	m.Org = map[string]*Org{}
	m.People = map[string]*Person{}
	var members []string

	switch {
	case from.Org == nil:
		return nil, fmt.Errorf("No organization section")
	case from.People == nil:
		return nil, fmt.Errorf("No people section")
	case len(from.People) == 0:
		return nil, fmt.Errorf("No people section")
	}
	org, ok := from.Org[name]
	if !ok {
		return nil, fmt.Errorf("No organization section for %s", name)
	}

	for _, login := range org.People {
		person, ok := from.People[login]
		if !ok {
			continue
		}
		m.People[login] = person
		members = append(members, person.Login)
	}
	m.Org["core"] = &Org{members}
	return m, nil
}

func parseMaintainerToml(data string) (*Maintainer, error) {
	m := new(Maintainer)
	_, err := toml.Decode(data, m)
	if err != nil {
		return nil, err
	}
	if m.People == nil {
		return nil, fmt.Errorf("Invalid Toml format. Missing people section.")
	}
	// if the person is defined in the file, but the Login field is
	// empty, we can use the map key as the Login value. This is mainly
	// here to support Docker projects, which use GitHub instead of Login
	for k, v := range m.People {
		if len(v.Login) == 0 {
			v.Login = k
		}
	}
	return m, nil
}

func parseMaintainerText(data string) (*Maintainer, error) {
	m := new(Maintainer)
	m.People = map[string]*Person{}

	buf := bytes.NewBufferString(data)
	reader := bufio.NewReader(buf)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		item := parseln(string(line))
		if len(item) == 0 {
			continue
		}

		person := parseLogin(item)
		if person == nil {
			person = parseLoginMeta(item)
		}
		if person == nil {
			person = parseLoginEmail(item)
		}
		if person == nil {
			return nil, fmt.Errorf("Invalid file format.")
		}

		m.People[person.Login] = person
	}
	return m, nil
}

func parseln(s string) string {
	if s == "" || string(s[0]) == "#" {
		return ""
	}
	index := strings.Index(s, " #")
	if index > -1 {
		s = strings.TrimSpace(s[0:index])
	}
	return s
}

// regular expression determines if a line in the maintainers
// file only has the single GitHub username and no other metadata.
var reLogin = regexp.MustCompile("^\\w[\\w-]+$")

// regular expression determines if a line in the maintainers
// file has the username and metadata.
var reLoginMeta = regexp.MustCompile("(.+) <(.+)> \\(@(.+)\\)")

// regular expression determines if a line in the maintainers
// file has the username and email.
var reLoginEmail = regexp.MustCompile("(.+) <(.+)>")

func parseLoginMeta(line string) *Person {
	matches := reLoginMeta.FindStringSubmatch(line)
	if len(matches) != 4 {
		return nil
	}
	return &Person{
		Name:  strings.TrimSpace(matches[1]),
		Email: strings.TrimSpace(matches[2]),
		Login: strings.TrimSpace(matches[3]),
	}
}

func parseLoginEmail(line string) *Person {
	matches := reLoginEmail.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil
	}
	return &Person{
		Login: strings.TrimSpace(matches[1]),
		Email: strings.TrimSpace(matches[2]),
	}
}

func parseLogin(line string) *Person {
	line = strings.TrimSpace(line)
	if !reLogin.MatchString(line) {
		return nil
	}
	return &Person{
		Login: line,
	}
}
