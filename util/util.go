package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	PI_INFO = iota
	PI_TASK
	PI_RESULT
)

// FileExists returns true if the file or directory exists
// otherwise it will return false.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// LocalRepoName returns the last part of an URL behind the slash.
// If the url is 'https://github.com/golang/go' it will return 'go'.
func LocalRepoName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

// PrintInfo prints text with a marker to easily spot it.
func PrintInfo(text string, mode int) {
	var pre string
	switch {
	case mode == PI_INFO:
		pre = ""
	case mode == PI_TASK:
		pre = "* "
	case mode == PI_RESULT:
		pre = "==> "
	}
	fmt.Println(pre + text)
}
