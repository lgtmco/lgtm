package model

import (
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/ianschenck/envflag"
)

type Config struct {
	Approvals       int    `json:"approvals"         toml:"approvals"`
	Pattern         string `json:"pattern"           toml:"pattern"`
	Team            string `json:"team"              toml:"team"`
	SelfApprovalOff bool   `json:"self_approval_off" toml:"self_approval_off"`

	re *regexp.Regexp
}

var (
	approvals = envflag.Int("LGTM_APPROVALS", 2, "")
	pattern = envflag.String("LGTM_PATTERN", "(?i)LGTM", "")
	team = envflag.String("LGTM_TEAM", "MAINTAINERS", "")
	selfApprovalOff = envflag.Bool("LGTM_SELF_APPROVAL_OFF", false, "")
)

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
		c.Approvals = *approvals
	}
	if len(c.Pattern) == 0 {
		c.Pattern = *pattern
	}
	if len(c.Team) == 0 {
		c.Team = *team
	}
	if c.SelfApprovalOff == false {
		c.SelfApprovalOff = *selfApprovalOff
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
