package util

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestBinaryInstalled(t *testing.T) {
	if BinaryInstalled("git") == false {
		t.Errorf("Failed: git should be installed on this computer")
	}
}

func TestFileExists(t *testing.T) {
	if FileExists("testfile") {
		t.Error("Failed: File should not exist")
		return
	}
	err := ioutil.WriteFile("testfile", []byte("foo"), 0644)
	if err != nil {
		t.Error("Unexpected error:", err)
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
		t.Error("Failed")
	}
	expected = "repo.git"
	actual = LocalRepoName("git@somewhere.com:repo.git")
	if actual != expected {
		t.Error("Failed")
	}
}
