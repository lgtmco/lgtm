package model

import (
	"regexp"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Approvals       int    `json:"approvals"         toml:"approvals"`
	Pattern         string `json:"pattern"           toml:"pattern"`
	Team            string `json:"team"              toml:"team"`
	SelfApprovalOff bool   `json:"self_approval_off" toml:"self_approval_off"`
	DoMerge         bool   `json:"do_merge" toml:"do_merge"`
	DoVersion       bool   `json:"do_version" toml:"do_version"`

	re *regexp.Regexp
}

// ParseConfig parses a projects .lgtm file
func ParseConfig(data []byte) (*Config, error) {
	return ParseConfigStr(string(data))
}

// ParseConfigStr parses a projects .lgtm file in string format.
func ParseConfigStr(data string) (*Config, error) {
	c := new(Config)
	_, err := toml.Decode(data, c)
	if err != nil {
		return nil, err
	}
	if c.Approvals == 0 {
		c.Approvals = 2
	}
	if len(c.Pattern) == 0 {
		c.Pattern = `(?i)LGTM\s*(\S*)`
	}
	if len(c.Team) == 0 {
		c.Team = "MAINTAINERS"
	}
	c.re, err = regexp.Compile(c.Pattern)
	return c, err
}

// IsMatch returns true if the text matches the regular
// epxression pattern.
func (c *Config) IsMatch(text string) bool {
	if c.re == nil {
		// this should never happen
		return false
	}
	return c.re.MatchString(text)
}
