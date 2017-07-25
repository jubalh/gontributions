package hg

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
func GetLatestRepo(url string) (err error) {
	var local string

	local = util.LocalRepoName(url)

	rd := RepoData{url: url, workingDirectory: "repos-hg", localName: local}

	if util.FileExists(filepath.Join("repos-hg", local)) {
		err = updateRepo(rd)
	} else {
		err = cloneRepo(rd)
	}
	return
}

// cloneRepo takes a RepoData struct and clones the repository
// specified in rd.
func cloneRepo(rd RepoData) error {
	//fmt.Printf("Running 'hg clone %s %s' in %s\n", rd.url, rd.localName, rd.workingDirectory)

	cmd := exec.Command("hg", "clone", rd.url, rd.localName)
	cmd.Dir = rd.workingDirectory
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()

	fmt.Print(cmdOutput) //TODO: do we want it?

	if err != nil {
		cerr := &util.RepoError{}
		cerr.Clone = true
		cerr.Err = err
		return cerr
	}

	return nil
}

// updateRepo takes a RepoData struct and updates the repository
// specified in rd.
func updateRepo(rd RepoData) error {
	// we need a 'hg pull -u' for 'hg pull' and 'hg update' here
	cmd := exec.Command("hg", "pull", "-u")
	cmd.Dir = filepath.Join(rd.workingDirectory, rd.localName)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()

	fmt.Print(cmdOutput) //TODO: do we want it?

	if err != nil {
		cerr := &util.RepoError{}
		cerr.Clone = true
		cerr.Err = err
		return cerr
	}

	return nil
}

// CountCommits returns how often email occurs in the log for
// the hg repository at url.
func CountCommits(path string, email string) (count int, err error) {
	cmd := exec.Command("hg", "log", "--template", "{author}\n", "-u", email)
	cmd.Dir = path
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()

	if err != nil {
		return
	}

	count = strings.Count((string(cmdOutput.Bytes())), "\n")
	return
}
