package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	repoPath        = "testdata/dummy-git-repo/"
	clonedReposPath = "testdata/repos/"
)

var (
	absoluteRepoPath   string
	absoluteTargetPath string
)

// mwm: The function "func init(){}" gets called by convention
// when the variables and constants are defined, but before
// any other method
func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	absoluteRepoPath = filepath.Join(wd, "testdata", "dummy-git-repo")
	absoluteTargetPath = filepath.Join(wd, "testdata", "repos")
}

// mwm: Use a TestMain func for setup and teardown,
// That is what it's there for
func TestMain(t *testing.M) {

	// mwm: It is better to sanitize the environment
	// than to throw errors
	if err := os.RemoveAll(absoluteTargetPath); err != nil {
		// We can panic, since there is something _seriously_ wrong
		panic(err)
	}
	// mwm: Always use the principle of least privilege
	os.Mkdir(absoluteTargetPath, 0755)
	fmt.Println(os.Getwd())
	// Run all the tests
	t.Run()

	if err := os.RemoveAll(clonedReposPath); err != nil {
		panic(err)
	}
}

func TestCloneRepo(t *testing.T) {

	t.Logf("Cloning\n\tfrom:%s\n\tto:%s", absoluteRepoPath, absoluteTargetPath)

	rd := RepoData{absoluteRepoPath, absoluteTargetPath, "dummy-git-repo"}

	if err := cloneRepo(rd); err != nil {
		t.Error("Error: ", err)
		// Does not make any sense to continue if we can not reliably clone a repo.
		t.FailNow()
	}

}

// Should I test this too?
// If yes I will need this ugly 'test' parameter :-(
func TestGetLatestGitRepo(t *testing.T) {

	err := GetLatestGitRepo(absoluteRepoPath, true)
	if err != nil {
		t.Error("Error: ", err)
	}

}

func TestCountCommits(t *testing.T) {

	rd := RepoData{"../testdata/dummy-git-repo", "repos", "dummy-git-repo"}
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
		t.Error("CountCommits returned: %d, expected: %d", count, expected)
	}
}
