package obs

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	absoluteTargetPath string
	repos              [2]OpenBuildService
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	absoluteTargetPath = filepath.Join(wd, "repos-obs")
	repos[0].Apiurl = "https://api.opensuse.org"
	repos[0].Repo = "utilities/vifm"
}

func TestCheckoutAndUpdateRepo(t *testing.T) {
	setup()
	defer teardown()
	checkosc(t)

	t.Logf("Checking out\n\t%s\n\tto:%s\n", repos[0].Repo, absoluteTargetPath)

	if err := checkoutRepo(repos[0]); err != nil {
		t.Error("Error: ", err)
		t.FailNow()
	}

	t.Logf("Updating %s\n", repos[0].Repo)

	if err := updateRepo(repos[0]); err != nil {
		t.Error("Error: ", err)
		t.FailNow()
	}
}

func TestGetLatestRepoAndCountCommits(t *testing.T) {
	setup()
	defer teardown()
	checkosc(t)

	t.Logf("Running first time GetLatestRepo, will checkout")

	if err := GetLatestRepo(repos[0]); err != nil {
		t.Error("Error: ", err)
		t.FailNow()
	}

	t.Logf("Running second time GetLatestRepo, will update")

	if err := GetLatestRepo(repos[0]); err != nil {
		t.Error("Error: ", err)
		t.FailNow()
	}

	expected := 8
	// Old email, wont use it anymore
	count, err := CountCommits(filepath.Join(absoluteTargetPath, repos[0].Repo), "g.bluehut@gmail.com")
	if err != nil {
		t.Error("Error: ", err)
		t.FailNow()
	}
	if count != expected {
		t.Errorf("Count returned: %d, expected: %d", count, expected)
	}
}

func setup() {
	err := os.MkdirAll(absoluteTargetPath, 0755)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	if err := os.RemoveAll(absoluteTargetPath); err != nil {
		panic(err)
	}
}
