package gontrib

import (
	"fmt"
	"os"

	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs/git"
	"github.com/jubalh/gontributions/vcs/mediawiki"
)

// MediaWiki holds the base URL of the wiki page to which later the
// API call will get appended and the username to the wiki.
type MediaWiki struct {
	BaseUrl string
	User    string
}

// Project hold all important information
// about a project.
type Project struct {
	Name        string
	Description string
	URL         string
	Gitrepos    []string
	MediaWikis  []MediaWiki
}

// Configuration holds the users E-Mail adresses
// and Projects he contributed to.
type Configuration struct {
	Emails   []string
	Projects []Project
}

// Contribution hols the Projects the user
// contributed to and a the ammounts of contributions
type Contribution struct {
	Project Project
	Count   int
}

// ScanContributions takes a Configuration containing a list of emails
// and a list of projects and returns a list of Contributions
// which contain of a project, how often a user contributed to it and
// description of the project.
func ScanContributions(configuration Configuration) []Contribution {
	contributions := []Contribution{}

	os.Mkdir("repos", 0755)

	for _, project := range configuration.Projects {
		var sumCount int
		for _, repo := range project.Gitrepos {
			util.PrintInfo("Working on " + repo)
			git.GetLatestGitRepo(repo, false)
			for _, email := range configuration.Emails {
				count, err := git.CountCommits(repo, email)
				if err != nil {
					fmt.Println(err)
				}

				s := fmt.Sprintf("%s: %d commits", email, count)
				util.PrintInfo(s)

				sumCount += count
			}
		}
		for _, wiki := range project.MediaWikis {
			wikiCount, err := mediawiki.GetUserEdits(wiki.BaseUrl, wiki.User)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			s := fmt.Sprintf("%d edits on %s MediaWiki as %s", wikiCount, wiki.BaseUrl, wiki.User)
			util.PrintInfo(s)

			sumCount += wikiCount
		}
		if sumCount > 0 {
			c := Contribution{project, sumCount}
			contributions = append(contributions, c)
		}
	}
	return contributions
}
