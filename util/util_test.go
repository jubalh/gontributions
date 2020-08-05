package util

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRepoError(t *testing.T) {
	err := &RepoError{
		Clone: true,
		Err:   errors.New("an error")}
	s := err.Error()
	if s != "an error" {
		t.Errorf("Expected 'an error' but got '%s'", s)
	}
}

func TestBinaryInstalled(t *testing.T) {
	if BinaryInstalled("git") == false {
		t.Error("Failed: git should be installed on this computer")
	}
}

func TestFileExists(t *testing.T) {
	if FileExists("testfile") {
		t.Error("Failed: File should not exist")
		return
	}
	err := ioutil.WriteFile("testfile", []byte("foo"), 0644)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !FileExists("testfile") {
		t.Error("Failed: File should exist")
	}
	os.Remove("testfile")
}

func TestLocalRepoName(t *testing.T) {
	expected := "nudoku"
	actual := LocalRepoName("https://github.com/jubalh/nudoku")
	if actual != expected {
		t.Errorf("Expected '%s' got '%s'", expected, actual)
	}
	expected = "repo.git"
	actual = LocalRepoName("git@somewhere.com:repo.git")
	if actual != expected {
		t.Errorf("Expected '%s' got '%s'", expected, actual)
	}
}

func TestPrintInfo(t *testing.T) {
	var b bytes.Buffer

	msg := "message"

	PrintInfo(&b, msg, PiInfo)
	actual := strings.TrimSpace(b.String())

	if actual != msg {
		t.Errorf("Expected '%s' got '%s'", msg, actual)
	}
}
