package git

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	absoluteBaseRepoPath string
	absoluteTargetPath   string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	absoluteBaseRepoPath = filepath.Join(wd, "testdata", "dummy-git-repo")
	absoluteTargetPath = filepath.Join(wd, "testdata", "repos")
}

func TestCloneRepo(t *testing.T) {
	setup()
	defer teardown()

	t.Logf("Cloning\n\tfrom:%s\n\tto:%s\n", absoluteBaseRepoPath, absoluteTargetPath)

	rd := RepoData{absoluteBaseRepoPath, absoluteTargetPath, "dummy-git-repo"}

	if err := cloneRepo(rd); err != nil {
		t.Error("Error: ", err)
		// Does not make any sense to continue if we can not reliably clone a repo.
		t.FailNow()
	}
}

func TestCountCommits(t *testing.T) {
	setup()
	defer teardown()

	rd := RepoData{absoluteBaseRepoPath, absoluteTargetPath, "dummy-git-repo"}
	err := cloneRepo(rd)
	if err != nil {
		t.Error("Error: ", err)
	}

	countCommit(t, "jubalh@openmailbox.org", rd.url, 1)
	countCommit(t, "bilbo@shire.ch", rd.url, 0)
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
