package vcs

import (
	"bytes"
	"fmt"
	"github.com/jubalh/gontributions/util"
	"os/exec"
	"strings"
)

func GetLatestGitRepo(url string) {
	local := util.LocalRepoName(url)
	if util.FileExists("repos/" + local) {
		updateRepo(url)
	} else {
		cloneRepo(url)
	}
}

func cloneRepo(url string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("git", "clone", url, local)
	cmd.Dir = "repos"
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Print(cmdOutput)
	return nil
}

func updateRepo(url string) error {
	local := util.LocalRepoName(url)

	cmd := exec.Command("git", "update")
	cmd.Dir = "repos/" + local
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Print(cmdOutput)
	return nil
}

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
