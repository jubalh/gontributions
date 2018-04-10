package debian

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/jubalh/gontributions/util"
)

type Debian struct {
	workingDirectory string
}

func NewDebian() Debian {
	return Debian{workingDirectory: "repos-debian"}
}

func (d Debian) GetWD() string {
	return d.workingDirectory
}

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
