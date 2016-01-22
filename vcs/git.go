package vcs

import (
	"bytes"
	"fmt"
	"github.com/jubalh/gontributions/util"
	"os/exec"
	"strings"
)

// RepoData holds the URL of a repository,
// the working directory where to execute the commands
// in, and the name of the local repositor.
type RepoData struct {
	url              string
	workingDirectory string
	localName        string
}

// GetLatestGitRepo either clones a new repo or updates an existing one
// into the 'repos' directory.
func GetLatestGitRepo(url string, isTest bool) error {
	var err error
	var local string

	if isTest {
		local = "dummy-git-repo"
	} else {
		local = util.LocalRepoName(url)
	}

	rd := RepoData{url: url, workingDirectory: "repos", localName: local}

	if util.FileExists("repos/" + local) {
		err = updateRepo(rd)
	} else {
		err = cloneRepo(rd)
	}
	return err
}

// cloneRepo takes a RepoData struct and clones the repository
// specified in rd.
func cloneRepo(rd RepoData) error {
	cmd := exec.Command("git", "clone", rd.url, rd.localName)
	cmd.Dir = rd.workingDirectory
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Print(cmdOutput)
	return nil
}

// updateRepo takes a RepoData struct and updates the repository
// specified in rd.
func updateRepo(rd RepoData) error {
	cmd := exec.Command("git", "update")
	cmd.Dir = rd.workingDirectory + "/" + rd.localName
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Print(cmdOutput)
	return nil
}

// CountCommits returns how often email occurs in the log for
// the git repository at url.
func CountCommits(url string, email string) (int, error) {
	local := util.LocalRepoName(url)

	authorSwitch := "--author=" + email
	cmd := exec.Command("git", "log", "--pretty=tformat:%s", authorSwitch)
	cmd.Dir = "repos/" + local
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	//fmt.Println(strings.Join(cmd.Args, " ") + " in " + cmd.Dir)
	s := (string(cmdOutput.Bytes()))
	return strings.Count(s, "\n"), nil
}
