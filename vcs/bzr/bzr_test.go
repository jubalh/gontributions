package bzr

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
	repoURL = "lp:my-first-project"
	absoluteTargetPath = filepath.Join(wd, "testdata")
	absoluteRepoPath = filepath.Join(absoluteTargetPath, util.LocalRepoName(repoURL))
}

func TestCloneRepo(t *testing.T) {
	setup()
	defer teardown()

	t.Logf("Cloning\n\tfrom:%s\n\tto:%s\n", repoURL, absoluteTargetPath)

	b := NewBzr()
	if err := b.CloneRepo(repoURL, absoluteTargetPath); err != nil {
		t.Error("Error: ", err)
		// Does not make any sense to continue if we can not reliably clone a repo.
		t.FailNow()
	}
}

func TestCountCommits(t *testing.T) {
	setup()
	defer teardown()

	b := NewBzr()
	if err := b.CloneRepo(repoURL, absoluteTargetPath); err != nil {
		t.Error("Error: ", err)
	}

	countCommit(t, b, "aarjan.baskota@gmail.com", absoluteRepoPath, 6)
	countCommit(t, b, "bilbo@shire.ch", absoluteRepoPath, 0)
}

func countCommit(t *testing.T, b Bzr, email string, path string, expected int) {
	count, err := b.Count(path, email)
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
