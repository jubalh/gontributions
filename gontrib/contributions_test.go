package gontrib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jubalh/gontributions/vcs/mediawiki"
	"github.com/jubalh/gontributions/vcs/obs"
)

func TestScanContributions(t *testing.T) {
	PullSources = true
	configuration := Configuration{
		Emails: []string{"jubalh@openmailbox.org", "g.bluehut@gmail.com"},
		Projects: []Project{
			{Name: "nudoku", Description: "Ncurses sudoku game", Gitrepos: []string{"https://github.com/jubalh/nudoku"}},
			{Name: "profanity", Description: "Ncurses based XMPP client", URL: "http://profanity.im/", Gitrepos: []string{"https://github.com/boothj5/profanity"}},
			{Name: "Funtoo", Description: "Linux distribution", URL: "http://funtoo.org/", Gitrepos: []string{"https://github.com/funtoo/ego", "https://github.com/funtoo/metro"}, MediaWikis: []mediawiki.MediaWiki{{BaseUrl: "http://funtoo.org", User: "jubalh"}}},
			{Name: "FuntooMediaOnly", Description: "Linux distribution", URL: "http://funtoo.org/", MediaWikis: []mediawiki.MediaWiki{{BaseUrl: "http://funtoo.org", User: "jubalh"}}},
			{Name: "openSUSE", Description: "Linux distribution", URL: "http://opensuse.org/", Obs: []obs.OpenBuildService{{Apiurl: "https://api.opensuse.org", Repo: "utilities/vifm"}}},
		},
	}

	contributions, err := ScanContributions(configuration)
	if err != nil {
		t.Errorf("Got an error scanning contributions: " + err.Error())
	}
	os.RemoveAll("repos-git")
	os.RemoveAll("repos-obs")
	os.RemoveAll("repos-hg")

	files, err := filepath.Glob("repos-*")
	if err != nil {
		t.Errorf("Got an error looking for repos- directories " + err.Error())
	}
	if len(files) >= 1 {
		s := strings.Join(files, ", ")
		t.Errorf("There are 'repo-' directories still present: " + s)
	}

	for i, contribution := range contributions {
		if i == 0 && contribution.Count < 80 {
			t.Errorf("Got wrong number of contributions for " + contribution.Project.Name)
		} else if i == 1 && contribution.Count < 9 {
			t.Errorf("Got wrong number of contributions for " + contribution.Project.Name)
		} else if i == 2 && contribution.Count < 99 {
			t.Errorf("Got wrong number of contributions for " + contribution.Project.Name)
		} else if i == 3 && contribution.Count < 93 {
			t.Errorf("Got wrong number of contributions for " + contribution.Project.Name)
		}
	}

}
