package hg

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jubalh/gontributions/util"
)

type Hg struct {
	workingDirectory string
}

func NewHg() *Hg {
	return &Hg{workingDirectory: "repos-hg"}
}

func (h Hg) GetWD() string {
	return h.workingDirectory
}

// CloneRepo takes an url and the directory where it should be cloned to
// in then checks out the directory there in a folder with the same name as
// the last part as the URL
func (h Hg) CloneRepo(url string, wd string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("hg", "clone", url, local)
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
func (h Hg) UpdateRepo(url string, wd string) error {
	local := util.LocalRepoName(url)

	// we need a 'hg pull -u' for 'hg pull' and 'hg update' here
	cmd := exec.Command("hg", "pull", "-u")
	cmd.Dir = filepath.Join(wd, local)
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
