package debian

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/jubalh/gontributions/util"
)

// Debian is the type to work with Debian contributions
type Debian struct {
	workingDirectory string
}

// NewDebian creates the directory in which we will download the changelogs
func NewDebian() Debian {
	return Debian{workingDirectory: "repos-debian"}
}

// GetWD returns the working directory.
func (d Debian) GetWD() string {
	return d.workingDirectory
}

// CloneRepo downloads the changelog via wget
func (d Debian) CloneRepo(pkg string, wd string) error {
	url := "http://metadata.ftp-master.debian.org/changelogs/main/" + pkg[0:1] + "/" + pkg + "/testing_changelog"
	cmd := exec.Command("wget", url, "-O", pkg)
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

// UpdateRepo downloads the changelog again
func (d Debian) UpdateRepo(pkg string, wd string) error {
	return d.CloneRepo(pkg, wd)
}

// Count returns how often email occurs in the log for
// the git repository at url.
func (d Debian) Count(path string, email string) (count int, err error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	count = strings.Count((string(f)), email)
	return
}
