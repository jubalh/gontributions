//go:generate go-bindata -pkg main -o default-templates-bindata.go templates/
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
	"time"

	"gopkg.in/urfave/cli.v1"

	"github.com/jubalh/gontributions/gontrib"
	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs/mediawiki"
	"github.com/jubalh/gontributions/vcs/obs"
)

// TemplateFill is there so it can easily be extended without breaking old templates/layouts
type TemplateFill struct {
	Contributions []gontrib.Contribution
}

const (
	templateFolderName = "templates"
	templatesFolderEnv = "GONTRIB_TEMPLATES_PTH"
)

// loadConfig loads a json configuration from filename
// and creates a Configuration from it.
func loadConfig(filename string) (gontribs gontrib.Configuration, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gontribs)
	return
}

func putdate() string {
	return time.Now().Local().Format("2006-01-02")
}

// fillTemplate puts the information of a Contribution
// into a template.
func fillTemplate(contributions []gontrib.Contribution, tempContent string, writer io.Writer) {
	tf := TemplateFill{Contributions: contributions}

	funcMap := template.FuncMap{
		"putdate": putdate,
	}

	t, err := template.New("string-template").Funcs(funcMap).Parse(tempContent)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = t.Execute(writer, tf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main will set and parse the cli options.
func main() {
	app := cli.NewApp()

	app.Name = "gontributions"
	app.Usage = "contributions lister"
	app.Author = "Michael Vetter"
	app.Version = "v0.5.2"
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
		cli.BoolFlag{
			Name:  "no-pull",
			Usage: "Don't update VCS repositories",
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

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Run will handle the functionallity.
func run(ctx *cli.Context) error {
	// Load specified json configuration file
	configPath := ctx.GlobalString("config")
	configuration, err := loadConfig(configPath)
	if err != nil {
		// if config cant be loaded because default one is used
		// (set in StringFlag) and is not available, then show the usage.
		if !ctx.IsSet("config") {
			cli.ShowAppHelp(ctx)
			return nil
		}
		return cli.NewExitError(err.Error(), 1)
	}

	// Get users template selection
	templateName := ctx.GlobalString("template")

	var templateData string

	// Get Template as templateData string
	templatesPath := os.Getenv(templatesFolderEnv)
	if templatesPath == "" {
		// Use asset
		data, err := Asset(filepath.Join(templateFolderName, templateName))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		templateData = string(data)
	} else {
		// Use template from user defined folder
		absoluteTemplatePath := filepath.Join(templatesPath, templateName)
		if !util.FileExists(absoluteTemplatePath) {
			var s string
			fmt.Sprintf(s, "Template file %s does not exist\n", absoluteTemplatePath)
			return cli.NewExitError(s, 1)
		}
		data, err := ioutil.ReadFile(absoluteTemplatePath)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		templateData = string(data)
	}

	gontrib.PullSources = !ctx.GlobalBool("no-pull")

	contributions, err := gontrib.ScanContributions(configuration)
	if err != nil {
		util.PrintInfo(nil, err.Error(), util.PI_ERROR)
		return cli.NewExitError(err.Error(), 1)
	}

	outputPath := ctx.GlobalString("output")
	f, err := os.Create(outputPath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	fillTemplate(contributions, templateData, writer)
	if err := writer.Flush(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	util.PrintInfoF(nil, "\nReport saved in: %s", util.PI_INFO, outputPath)

	errorfile, err := os.Open("errors.log")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer errorfile.Close()
	fi, err := errorfile.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if fi.Size() > 0 {
		util.PrintInfoF(nil, "Some contributions could not be checked. See: errors.log", util.PI_ERROR)
	} else {
		os.Remove("errors.log")
	}

	return nil
}

// Create an example configuration file which the user can
// adapt to his own needs.
func cmdExconf(c *cli.Context) error {
	configuration := gontrib.Configuration{
		Emails: []string{"jubalh@openmailbox.org", "g.bluehut@gmail.com"},
		Projects: []gontrib.Project{
			{Name: "nudoku", Description: "Ncurses sudoku game", Gitrepos: []string{"https://github.com/jubalh/nudoku"}},
			{Name: "profanity", Description: "Ncurses based XMPP client", URL: "http://profanity.im/", Gitrepos: []string{"https://github.com/boothj5/profanity"}},
			{Name: "Funtoo", Description: "Linux distribution", URL: "http://funtoo.org/", Gitrepos: []string{"https://github.com/funtoo/ego", "https://github.com/funtoo/metro"}, MediaWikis: []mediawiki.MediaWiki{{BaseUrl: "http://funtoo.org", User: "jubalh"}}},
			{Name: "openSUSE", Description: "Linux distribution", URL: "http://opensuse.org/", Obs: []obs.OpenBuildService{{Apiurl: "https://api.opensuse.org", Repo: "utilities/vifm"}}},
			{Name: "prosody", Description: "A modern XMPP communication server", URL: "http://prosody.im/", Hgrepos: []string{"https://hg.prosody.im/prosody-modules"}},
		},
	}

	text, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(string(text))
	return nil
}
