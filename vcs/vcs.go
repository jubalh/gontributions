package vcs

import (
	"path/filepath"

	"github.com/jubalh/gontributions/util"
)

// VCS defines an interface for the types Git, OBS etc.
// They all need to clone and update a repo.
// get the working dir and count the commits based on a given email.
type VCS interface {
	CloneRepo(url string, wd string) error
	UpdateRepo(url string, wd string) error
	GetWD() string
	Count(path string, email string) (count int, err error)
}

// GetLatestRepo either clones a new repo or updates an existing one
// into the 'repos' directory.
func GetLatestRepo(url string, v VCS) (err error) {
	local := util.LocalRepoName(url)

	if util.FileExists(filepath.Join(v.GetWD(), local)) {
		err = v.UpdateRepo(url, v.GetWD())
	} else {
		err = v.CloneRepo(url, v.GetWD())
	}
	return
}
