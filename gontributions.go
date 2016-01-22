package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/jubalh/gontributions/gontrib"
)

const (
	templateFolderName = "templates"
	templatesFolderEnv = "GONTRIB_TEMPLATES_PTH"
)

// loadConfig loads a json configuration from filename
// and creates a Configuration from it.
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

// fillTemplate puts the information of a Contribution
// into a template.
func fillTemplate(contributions []gontrib.Contribution, templatePath string, writer io.Writer) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
	}
	//err = t.Execute(os.Stdout, contributions)
	err = t.Execute(writer, contributions)
	if err != nil {
		fmt.Println(err)
	}
}

// Main will set and parse the cli options.
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
			Value: "gontrib.json",
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
			Usage:  "Show an example configuration file",
			Action: cmdExconf,
		},
	}

	app.Action = run

	app.Run(os.Args)
}

// Run will handle the functionallity.
func run(cli *cli.Context) {
	// Load specified json configuration file
	configPath := cli.GlobalString("config")
	configuration, err := loadConfig(configPath)
	if err != nil {
		fmt.Println(err)
	}

	contributions := gontrib.ScanContributions(configuration)

	// Get users template selection
	templateName := cli.GlobalString("template")

	// Build path to template
	templatesPath := os.Getenv(templatesFolderEnv)
	if templatesPath == "" {
		templatesPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		templatesPath = filepath.Join(templatesPath, templateFolderName)
	}
	absoluteTemplatePath := filepath.Join(templatesPath, templateName)

	outputPath := cli.GlobalString("output")
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	fillTemplate(contributions, absoluteTemplatePath, writer)
	writer.Flush()
}

// Create an example configuration file which the user can
// adapt to his own needs.
func cmdExconf(c *cli.Context) {
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
	fmt.Println(string(text))
}
