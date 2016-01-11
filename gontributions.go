package main

import (
	"encoding/json"
	"fmt"
	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs"
	"html/template"
	"io/ioutil"
	"os"
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
	Project string
	Count   int
}

func loadConfig() (Configuration, error) {
	contribs := Configuration{}

	s, err := ioutil.ReadFile("config.json")
	if err != nil {
		return contribs, err
	}
	err = json.Unmarshal(s, &contribs)
	if err != nil {
		return contribs, err
	}

	return contribs, nil
}

func main() {
	configuration, err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}

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
		c := Contribution{project.Name, sumCount}
		contributions = append(contributions, c)
	}
	temp(contributions)
}

func temp(contributions []Contribution) {
	t, err := template.ParseFiles("templates/default.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(os.Stdout, contributions)
	if err != nil {
		fmt.Println(err)
	}
}
