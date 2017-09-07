package vcs

import (
	"path/filepath"

	"github.com/jubalh/gontributions/util"
)

type VCS interface {
	CloneRepo(url string, wd string) error
	UpdateRepo(url string, wd string) error
	GetWD() string
}

// GetLatestRepo either clones a new repo or updates an existing one
// into the 'repos' directory.
func GetLatestRepo(url string, v VCS) (err error) {
	var local string

	if util.FileExists(filepath.Join(v.GetWD(), local)) {
		err = v.UpdateRepo(url, v.GetWD())
	} else {
		err = v.CloneRepo(url, v.GetWD())
	}
	return
}
