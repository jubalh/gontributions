package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jubalh/gontributions/gontrib"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

func loadConfig(filename string) (gontrib.Configuration, error) {
	contribs := gontrib.Configuration{}

	s, err := ioutil.ReadFile(filename)
	if err != nil {
		return contribs, err
	}
	err = json.Unmarshal(s, &contribs)
	if err != nil {
		return contribs, err
	}

	return contribs, nil
}

func fillTemplate(contributions []gontrib.Contribution, templateName string, writer io.Writer) {
	t, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		fmt.Println(err)
	}
	//err = t.Execute(os.Stdout, contributions)
	err = t.Execute(writer, contributions)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	app := cli.NewApp()

	app.Name = "gontributions"
	app.Usage = "contributions lister"
	app.Author = "Michael Vetter"
	app.Version = "0.1"
	app.Email = "jubalh@openmailbox.org"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "config.json",
			Usage: "Set which config file to use",
		},
		cli.StringFlag{
			Name:  "template",
			Value: "default.html",
			Usage: "Set which template to use",
		},
		cli.StringFlag{
			Name:  "output",
			Value: "output.html",
			Usage: "Define name of the generated HTMl file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "exconf",
			Usage:  "Create an example config file",
			Action: cmdExconf,
		},
	}

	app.Action = run

	app.Run(os.Args)
}

func run(cli *cli.Context) {
	configPath := cli.GlobalString("config")
	configuration, err := loadConfig(configPath)
	if err != nil {
		fmt.Println(err)
	}

	contributions := gontrib.ScanContributions(configuration)

	templateName := cli.GlobalString("template")

	outputPath := cli.GlobalString("output")
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	fillTemplate(contributions, templateName, writer)
	writer.Flush()
}

func cmdExconf(c *cli.Context) {
	f, err := os.Create("example_conf.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)

	configuration := gontrib.Configuration{
		Emails: []string{"jubalh@openmailbox.org", "g.bluehut@gmail.com"},
		Projects: []gontrib.Project{
			{Name: "nudoku", Description: "Ncurses sudoku game", Gitrepos: []string{"https://github.com/jubalh/nudoku"}},
			{Name: "profanity", Description: "Ncurses based XMPP client", Gitrepos: []string{"https://github.com/boothj5/profanity"}},
		},
	}

	text, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	writer.Write(text)
	writer.Flush()
}
