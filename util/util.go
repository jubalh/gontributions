package util

import (
	"fmt"
	"os"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func LocalRepoName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func PrintInfo(text string) {
	fmt.Println("==> ", text)
}
