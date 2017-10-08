package bzr

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jubalh/gontributions/util"
)

type Bzr struct {
	workingDirectory string
}

func NewBzr() Bzr {
	return Bzr{workingDirectory: "repos-bzr"}
}

func (b Bzr) GetWD() string {
	return b.workingDirectory
}

// CloneRepo takes an url and the directory where it should be cloned to
// in then checks out the directory there in a folder with the same name as
// the last part as the URL
func (b Bzr) CloneRepo(url string, wd string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("bzr", "branch", url, local)
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
func (b Bzr) UpdateRepo(url string, wd string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("bzr", "pull")
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
// the Bazaar repository at url.
func (b Bzr) Count(path string, email string) (count int, err error) {
	authorSwitch := "--match-author=" + email
	cmd := exec.Command("bzr", "log", "--short", authorSwitch)
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
