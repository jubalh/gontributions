package vcs

import (
	"os"
	"testing"
)

var (
	repoPath        = "testdata/dummy-git-repo/"
	clonedReposPath = "testdata/repos/"
)

// Use a TestMain func for setup and teardowm
func TestMain(t *testing.M) {

	// It is better to sanitize the environment
	// Than to throw errors
	if err := os.RemoveAll(clonedReposPath); err != nil {
		// We can panic, since there is something _seriously_ wrong
		panic(err)
	}
	// Always use the principle of least privilege
	os.Mkdir("testdata/repos", 0755)
	
	// Run all the tests
	t.Run()
	
	if err := os.RemoveAll(clonedReposPath); err != nil {
		panic(err)
	}
}

func TestCloneRepo(t *testing.T) {

	rd := RepoData{"../testdata/dummy-git-repo", "repos", "dummy-git-repo"}
	err := cloneRepo(rd)
	if err != nil {
		t.Error("Error: ", err)
	}

}

// Should I test this too?
// If yes I will need this ugly 'test' parameter :-(
func TestGetLatestGitRepo(t *testing.T) {

	err := GetLatestGitRepo("../testdata/dummy-git-repo", true)
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
