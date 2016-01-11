package main

import (
	"encoding/json"
	"fmt"
	"github.com/jubalh/gontributions/gontrib"
	"html/template"
	"io/ioutil"
	"os"
)

func loadConfig() (gontrib.Configuration, error) {
	contribs := gontrib.Configuration{}

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

func fillTemplate(contributions []gontrib.Contribution) {
	t, err := template.ParseFiles("templates/detailed.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(os.Stdout, contributions)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	configuration, err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}

	contributions := gontrib.ScanContributions(configuration)

	fillTemplate(contributions)
}
