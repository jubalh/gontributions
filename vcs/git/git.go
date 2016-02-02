package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jubalh/gontributions/util"
)

// RepoData holds the URL of a repository,
// the working directory where to execute the commands
// in, and the name of the local repositor.
type RepoData struct {
	url              string
	workingDirectory string
	localName        string
}

// GetLatestRepo either clones a new repo or updates an existing one
// into the 'repos' directory.
func GetLatestRepo(url string) error {
	var err error
	var local string

	local = util.LocalRepoName(url)

	rd := RepoData{url: url, workingDirectory: "repos-git", localName: local}

	if util.FileExists(filepath.Join("repos-git", local)) {
		err = updateRepo(rd)
	} else {
		err = cloneRepo(rd)
	}
	return err
}

// cloneRepo takes a RepoData struct and clones the repository
// specified in rd.
func cloneRepo(rd RepoData) error {
	//fmt.Printf("Running 'git clone %s %s' in %s\n", rd.url, rd.localName, rd.workingDirectory)

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
	cmd := exec.Command("git", "pull")
	cmd.Dir = filepath.Join(rd.workingDirectory, rd.localName)
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
func CountCommits(path string, email string) (int, error) {
	authorSwitch := "--author=" + email
	cmd := exec.Command("git", "log", "--pretty=tformat:%s", authorSwitch)
	cmd.Dir = path
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
