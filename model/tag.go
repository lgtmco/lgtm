package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-version"
)

type Tag string

type TagList []Tag

func (tl TagList) GetMaxTag() (Tag, *version.Version) {
	//find the previous largest semver value
	var maxVer *version.Version
	var maxTag Tag

	for _, tag := range tl {
		curVer, err := version.NewVersion(string(tag))
		if err != nil {
			continue
		}
		if maxVer == nil || curVer.GreaterThan(maxVer) {
			maxVer = curVer
			maxTag = tag
		}
	}

	if maxVer == nil {
		maxVer, _ = version.NewVersion("v0.0.0")
	}
	log.Debugf("maxVer found is %s", maxVer.String())
	return maxTag, maxVer
}
