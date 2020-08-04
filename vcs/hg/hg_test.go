package hg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jubalh/gontributions/util"
)

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
	repoURL = "https://bitbucket.org/Mojachiee/testrepo"
	absoluteTargetPath = filepath.Join(wd, "testdata")
	absoluteRepoPath = filepath.Join(absoluteTargetPath, util.LocalRepoName(repoURL))
}

func TestCloneRepo(t *testing.T) {
	setup()
	defer teardown()

	t.Logf("Cloning\n\tfrom:%s\n\tto:%s\n", repoURL, absoluteTargetPath)

	h := NewHg()
	if err := h.CloneRepo(repoURL, absoluteTargetPath); err != nil {
		t.Error("Error: ", err)
		// Does not make any sense to continue if we can not reliably clone a repo.
		t.FailNow()
	}
}

func TestCountCommits(t *testing.T) {
	setup()
	defer teardown()

	h := NewHg()
	if err := h.CloneRepo(repoURL, absoluteTargetPath); err != nil {
		t.Error("Error: ", err)
	}

	countCommit(t, h, "joe.stephenson36@gmail.com", absoluteRepoPath, 2)
	countCommit(t, h, "bilbo@shire.ch", absoluteRepoPath, 0)
}

func countCommit(t *testing.T, h *Hg, email string, path string, expected int) {
	count, err := h.Count(path, email)
	if err != nil {
		t.Error("Unexpected error: ", err)
	}
	if count != expected {
		t.Errorf("Count returned: %d, expected: %d", count, expected)
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
