package vcs

import (
	"github.com/jubalh/gontributions/util"
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	setup(t)

	rd := RepoData{"../testdata/dummy-git-repo", "repos", "dummy-git-repo"}
	err := cloneRepo(rd)
	if err != nil {
		t.Error("Error: ", err)
	}

	cleanUp(t)
}

// Should I test this too?
// If yes I will need this ugly 'test' parameter :-(
func TestGetLatestGitRepo(t *testing.T) {
	setup(t)

	err := GetLatestGitRepo("../testdata/dummy-git-repo", true)
	if err != nil {
		t.Error("Error: ", err)
	}

	cleanUp(t)
}

func TestCountCommits(t *testing.T) {
	setup(t)

	rd := RepoData{"../testdata/dummy-git-repo", "repos", "dummy-git-repo"}
	err := cloneRepo(rd)
	if err != nil {
		t.Error("Error: ", err)
	}

	countCommit(t, "jubalh@openmailbox.org", rd.url, 1)
	countCommit(t, "bilbo@shire.ch", rd.url, 0)

	cleanUp(t)
}

func setup(t *testing.T) {
	os.Mkdir("repos", 0777)
	if util.FileExists("repos/dummy-git-repo") {
		t.Error("Dummy repo should not exist")
	}
}

func cleanUp(t *testing.T) {
	if util.FileExists("repos") {
		err := os.RemoveAll("repos")
		if err != nil {
			t.Error("Unexpected error: ", err)
		}
	}
}

func countCommit(t *testing.T, email string, url string, expected int) {
	count, err := CountCommits(url, email)
	if err != nil {
		t.Error("Unexpected error: ", err)
	}
	if count != expected {
		t.Error("CountCommits returned: %d, expected: %d", count, expected)
	}
}
