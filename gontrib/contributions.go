package gontrib

import (
	"fmt"
	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs"
	"os"
)

type Project struct {
	Name        string
	Description string
	Gitrepos    []string
}

type Configuration struct {
	Emails   []string
	Projects []Project
}

type Contribution struct {
	Project     string
	Count       int
	Description string
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
			vcs.GetLatestGitRepo(repo, false)
			for _, email := range configuration.Emails {
				count, err := vcs.CountCommits(repo, email)
				if err != nil {
					fmt.Println(err)
				}

				s := fmt.Sprintf("%s: %d commits", email, count)
				util.PrintInfo(s)

				sumCount += count
			}
		}
		if sumCount > 0 {
			c := Contribution{project.Name, sumCount, project.Description}
			contributions = append(contributions, c)
		}
	}
	return contributions
}
