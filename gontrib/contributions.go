package gontrib

import (
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

// scanGit is a helper function for ScanContributions which takes care of the git part
func scanGit(project Project, emails []string, contributions []Contribution) (int, error) {
	var sum int
	for _, repo := range project.Gitrepos {
		util.PrintInfo("Working on "+repo, util.PI_TASK)
		git.GetLatestRepo(repo)
		for _, email := range emails {
			path := filepath.Join("repos-git", util.LocalRepoName(repo))
			gitCount, err := git.CountCommits(path, email)
			if err != nil {
				return 0, err
			}

			if gitCount != 0 {
				util.PrintInfoF("%s: %d commits", util.PI_RESULT, email, gitCount)
				sum += gitCount
			}
		}
	}
	return sum, nil
}

// scanWiki is a helper function for ScanContributions which takes care of the MediaWiki part
func scanWiki(project Project, emails []string, contributions []Contribution) int {
	var sum int
	for _, wiki := range project.MediaWikis {
		util.PrintInfoF("Working on MediaWiki %s as %s", util.PI_TASK, wiki.BaseUrl, wiki.User)

		wikiCount, err := mediawiki.GetUserEdits(wiki.BaseUrl, wiki.User)
		if err != nil {
			util.PrintInfo(err.Error(), util.PI_ERROR)
		}

		if wikiCount != 0 {
			util.PrintInfoF("%d edits", util.PI_RESULT, wikiCount)
			sum += wikiCount
		}
	}
	return sum
}

// scanOBS is a helper function for ScanContributions which takes care of the OBS part
func scanOBS(project Project, emails []string, contributions []Contribution) (int, error) {
	var sum int
Loop_obs:
	for _, obsEntry := range project.Obs {
		util.PrintInfo("Working on "+obsEntry.Repo, util.PI_TASK)

		err := obs.GetLatestRepo(obsEntry)
		if err != nil {
			return 0, err
		}
		for _, email := range emails {
			obsCount, err := obs.CountCommits("repos-obs"+"/"+obsEntry.Repo, email)
			if err != nil {
				if err == obs.ErrNoChangesFileFound {
					util.PrintInfo("No .changes file found", util.PI_RESULT) // TODO: mild error?
					break Loop_obs
				}
				util.PrintInfo(err.Error(), util.PI_ERROR) // TODO: return?
			}

			if obsCount != 0 {
				util.PrintInfoF("%s: %d changes", util.PI_RESULT, email, obsCount)
				sum += obsCount
			}
		}
	}
	return sum, nil
}

// ScanContributions takes a Configuration containing a list of emails
// and a list of projects and returns a list of Contributions
// which contain of a project, how often a user contributed to it and
// description of the project.
func ScanContributions(configuration Configuration) ([]Contribution, error) {
	contributions := []Contribution{}

	os.Mkdir("repos-git", 0755)
	os.Mkdir("repos-obs", 0755)

	for _, project := range configuration.Projects {
		var sumCount int

		sum, err := scanGit(project, configuration.Emails, contributions)
		if err != nil {
			return nil, err
		}
		sumCount += sum

		sum = scanWiki(project, configuration.Emails, contributions)
		sumCount += sum

		sum, err = scanOBS(project, configuration.Emails, contributions)
		if err != nil {
			return nil, err
		}
		sumCount += sum

		if sumCount > 0 {
			c := Contribution{project, sumCount}
			contributions = append(contributions, c)
		}
	}

	return contributions, nil
}
