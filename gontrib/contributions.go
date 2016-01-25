package gontrib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs/git"
	"github.com/jubalh/gontributions/vcs/mediawiki"
	"github.com/jubalh/gontributions/vcs/obs"
)

// Project hold all important information
// about a project.
type Project struct {
	Name        string
	Description string
	URL         string
	Gitrepos    []string
	MediaWikis  []mediawiki.MediaWiki
	Obs         []obs.OpenBuildService
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
	os.Mkdir("repos-ob", 0755)

	for _, project := range configuration.Projects {
		var sumCount int
		for _, repo := range project.Gitrepos {
			util.PrintInfo("Working on "+repo, util.PI_TASK)
			git.GetLatestGitRepo(repo)
			for _, email := range configuration.Emails {
				path := filepath.Join("repos-git", util.LocalRepoName(repo))
				gitCount, err := git.CountCommits(path, email)
				if err != nil {
					fmt.Println(err)
				}

				if gitCount != 0 {
					util.PrintInfoF("%s: %d commits", util.PI_RESULT, email, gitCount)
					sumCount += gitCount
				}

			}
		}

		for _, wiki := range project.MediaWikis {
			util.PrintInfoF("Working on MediaWiki %s as %s", util.PI_TASK, wiki.BaseUrl, wiki.User)

			wikiCount, err := mediawiki.GetUserEdits(wiki.BaseUrl, wiki.User)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			if wikiCount != 0 {
				util.PrintInfoF("%d edits", util.PI_RESULT, wikiCount)
				sumCount += wikiCount
			}
		}

	Loop_obs:
		for _, obsEntry := range project.Obs {
			util.PrintInfo("Working on "+obsEntry.Repo, util.PI_TASK)

			err := obs.GetLatestRepo(obsEntry)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			for _, email := range configuration.Emails {
				obsCount, err := obs.CountCommits("repos-obs"+"/"+obsEntry.Repo, email)
				if err != nil {
					if err == obs.ErrNoChangesFileFound {
						util.PrintInfo("No .changes file found", util.PI_RESULT)
						break Loop_obs
					}
					fmt.Fprintln(os.Stderr, err)
				}

				if obsCount != 0 {
					util.PrintInfoF("%s: %d changes", util.PI_RESULT, email, obsCount)
					sumCount += obsCount
				}
			}
		}

		if sumCount > 0 {
			c := Contribution{project, sumCount}
			contributions = append(contributions, c)
		}
	}
	return contributions
}
