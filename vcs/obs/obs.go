package obs

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jubalh/gontributions/util"
)

// ErrNoChangesFileFound is an error used when there is no .changes file
// in an OpenBuildService directory.
var ErrNoChangesFileFound = errors.New("No .changes file found")

// OpenBuildService holds the API URL to the open build service instance and the
// repository name
type OpenBuildService struct {
	Apiurl string
	Repo   string
}

// GetLatestRepo gets the newest version of an OpenBuildService
// repository.
func GetLatestRepo(info OpenBuildService) error {
	var err error
	if util.FileExists(filepath.Join("repos-obs", info.Repo)) {
		err = updateRepo(info)
	} else {
		err = checkoutRepo(info)
	}
	return err
}

// checkoutRepo checks out a new OBS repo.
func checkoutRepo(info OpenBuildService) error {
	cmd := exec.Command("osc", "-A", info.Apiurl, "co", info.Repo)
	cmd.Dir = "repos-obs"
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	return err
}

// updateRepo gets the latest changes from an existing OBS repo.
func updateRepo(info OpenBuildService) error {
	cmd := exec.Command("osc", "up")
	cmd.Dir = filepath.Join("repos-obs", info.Repo)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	return err
}

// CountCommits returns the number of commits in the OpenBuildService
// repository saved at path for email.
// It returns ErrNoChangesFileFound error in case it couldnt locate any
// .changes file And forwards other errors that might occur.
func CountCommits(path string, email string) (count int, err error) {
	var changesFiles []string

	//search for a file ending in '.changes'
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".osc" {
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, ".changes") {
			changesFiles = append(changesFiles, path)
		}
		return err
	})

	if err != nil && err != filepath.SkipDir {
		return
	}

	if len(changesFiles) == 0 {
		err = ErrNoChangesFileFound //TODO: can we add the filename here?
		return
	}

	// Read every .changes file
	for i := 0; i < len(changesFiles); i++ {
		changesFile := changesFiles[i]

		var changes []byte
		changes, err = ioutil.ReadFile(changesFile)
		if err != nil {
			return
		}
		// for now just count how often the mail occurs
		count += strings.Count(string(changes), email)
	}

	return count, nil
}
