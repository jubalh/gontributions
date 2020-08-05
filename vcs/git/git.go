package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jubalh/gontributions/util"
)

// Git is the type to work with git repos.
type Git struct {
	workingDirectory string
}

// NewGit creates the directory in which we will download the repos.
func NewGit() Git {
	return Git{workingDirectory: "repos-git"}
}

// GetWD returns the working directory.
func (g Git) GetWD() string {
	return g.workingDirectory
}

// CloneRepo takes an url and the directory where it should be cloned to
// in then checks out the directory there in a folder with the same name as
// the last part as the URL
func (g Git) CloneRepo(url string, wd string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("git", "clone", url, local)
	cmd.Dir = wd
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

// UpdateRepo takes an URL and a working directory
func (g Git) UpdateRepo(url string, wd string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("git", "pull")
	cmd.Dir = filepath.Join(wd, local)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()

	fmt.Print(cmdOutput) //TODO: do we want it?

	if err != nil {
		cerr := &util.RepoError{}
		cerr.Update = true
		cerr.Err = err
		return cerr
	}

	return nil
}

// Count returns how often email occurs in the log for
// the git repository at url.
func (g Git) Count(path string, email string) (count int, err error) {
	authorSwitch := "--author=" + email
	cmd := exec.Command("git", "log", "--pretty=tformat:%s", authorSwitch)
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
