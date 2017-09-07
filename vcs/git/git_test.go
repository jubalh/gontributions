package git

import (
	"os"
	"path/filepath"
	"testing"
)

const reponame = "dummy-git-repo"

var (
	repoURL            string
	absoluteTargetPath string
	absoluteRepoPath   string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	repoURL = "https://github.com/jubalh/testrepo"
	absoluteTargetPath = filepath.Join(wd, "testdata")
	absoluteRepoPath = filepath.Join(absoluteTargetPath, reponame)
}

func TestCloneRepo(t *testing.T) {
	setup()
	defer teardown()

	t.Logf("Cloning\n\tfrom:%s\n\tto:%s\n", repoURL, absoluteTargetPath)

	rd := RepoData{repoURL, absoluteTargetPath, reponame}

	if err := cloneRepo(rd); err != nil {
		t.Error("Error: ", err)
		// Does not make any sense to continue if we can not reliably clone a repo.
		t.FailNow()
	}
}

func TestCountCommits(t *testing.T) {
	setup()
	defer teardown()

	rd := RepoData{repoURL, absoluteTargetPath, reponame}
	err := cloneRepo(rd)
	if err != nil {
		t.Error("Error: ", err)
	}

	countCommit(t, "jubalh@openmailbox.org", absoluteRepoPath, 1)
	countCommit(t, "bilbo@shire.ch", absoluteRepoPath, 0)
}

func countCommit(t *testing.T, email string, url string, expected int) {
	count, err := CountCommits(url, email)
	if err != nil {
		t.Error("Unexpected error: ", err)
	}
	if count != expected {
		t.Errorf("CountCommits returned: %d, expected: %d", count, expected)
	}
}

func setup() {
	err := os.MkdirAll(absoluteTargetPath, 0755)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	if err := os.RemoveAll(absoluteTargetPath); err != nil {
		panic(err)
	}
}
