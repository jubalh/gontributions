package util

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	PI_INFO = iota
	PI_TASK
	PI_RESULT
	PI_ERROR
	PI_MILD_ERROR
)

// FileExists returns true if the file or directory exists
// otherwise it will return false.
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// BinaryInstalled checks whether binary exists in any path in
// the PATH environment variable.
func BinaryInstalled(binary string) bool {
	env := os.Getenv("PATH")
	paths := strings.Split(env, ":")
	if runtime.GOOS == "windows" {
		paths = strings.Split(env, ";")
	}
	for _, fpath := range paths {
		pathToBinary := path.Join(fpath, binary)
		if runtime.GOOS == "windows" {
			pathToBinary = path.Join(fpath, binary+".exe")
		}
		if FileExists(pathToBinary) {
			return true
		}
	}
	return false
}

// LocalRepoName returns the last part of an URL behind the slash.
// If the url is 'https://github.com/golang/go' it will return 'go'.
// If the url is 'git@golang.com:go' it will also return 'go'.
func LocalRepoName(url string) string {
	parts := strings.Split(url, "/")
	possibleName := parts[len(parts)-1]
	parts = strings.Split(possibleName, ":")
	return parts[len(parts)-1]
}

// PrintInfo prints text with a marker to easily spot it.
func PrintInfo(text string, mode int) {
	var pre string
	switch {
	case mode == PI_INFO:
		pre = ""
	case mode == PI_TASK:
		pre = "\033[36m*\033[39m "
	case mode == PI_RESULT:
		pre = "\033[35m==>\033[39m "
	case mode == PI_ERROR:
		pre = "\033[31mError:\033[39m "
	case mode == PI_MILD_ERROR:
		pre = "\033[33m==\033[39m "
	}
	fmt.Println(pre + text)
}

// PrintInfo prints formatable text with a marker to easily spot it.
func PrintInfoF(text string, mode int, args ...interface{}) {
	PrintInfo(fmt.Sprintf(text, args...), mode)
}
