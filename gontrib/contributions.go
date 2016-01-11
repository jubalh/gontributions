package gontrib

import (
	"fmt"
	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs"
)

type project struct {
	Name        string
	Description string
	Gitrepos    []string
}

type Configuration struct {
	Emails   []string
	Projects []project
}

type Contribution struct {
	Project     string
	Count       int
	Description string
}

func ScanContributions(configuration Configuration) []Contribution {
	contributions := []Contribution{}

	for _, project := range configuration.Projects {
		var sumCount int
		for _, repo := range project.Gitrepos {
			util.PrintInfo("Working on " + repo)
			vcs.GetLatestGitRepo(repo)
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
		c := Contribution{project.Name, sumCount, project.Description}
		contributions = append(contributions, c)
	}
	return contributions
}
