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
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"

	"github.com/jubalh/gontributions/gontrib"
	"github.com/jubalh/gontributions/util"
	"github.com/jubalh/gontributions/vcs/mediawiki"
	"github.com/jubalh/gontributions/vcs/obs"
)

// Summary hold the overall nr of contributions and projects
// we do it like this so we can calculate here, in templates its not supported
type Summary struct {
	TotalContributions int
	ProjectCount       int
}

// TemplateFill is there so it can easily be extended without breaking old templates/layouts
type TemplateFill struct {
	Contributions []gontrib.Contribution
	Summary       Summary
}

const (
	templateAssetFolderName = "templates"
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

// putdate() returns the current date.
// Meant for the template.
func putdate() string {
	return time.Now().Local().Format("2006-01-02")
}

// fillTemplate puts the information of a Contribution
// into a template.
func fillTemplate(contributions []gontrib.Contribution, tempContent string, writer io.Writer) {
	tf := TemplateFill{Contributions: contributions}

	for _, contribution := range contributions {
		tf.Summary.TotalContributions += contribution.Count
	}
	tf.Summary.ProjectCount = len(contributions)

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
			Value: "default",
			Usage: "Set which template to use. If it contains a dot it will be treated as a path to a user defined template. No dot means using an internal template",
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
		{
			Name:   "list-templates",
			Usage:  "List all built-in templates",
			Action: cmdListTemplates,
		},
		{
			Name:   "show-template",
			Usage:  "Display the content of a template",
			Action: cmdShowTemplate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "default",
					Usage: "Display the content of a template",
				},
			},
		},
	}

	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// getTemplateOrExitError gets the template, if the templatename doesnt contain
// a dot, from the internal assets. if it contains a dot it is treated as a user
// defined template and is treated like a path to the file.
// it returns the template as a string
func getTemplateOrExitError(ctx *cli.Context) (string, error) {
	// Get users template selection
	templateName := ctx.GlobalString("template")

	// if it does not contain a dot it is an internal template, from the assets
	if strings.Contains(templateName, ".") == false {
		// Use asset
		data, err := Asset(filepath.Join(templateAssetFolderName, templateName) + ".html")
		if err != nil {
			return "", cli.NewExitError(err.Error(), 1)
		}
		return string(data), nil
	} else {
		// user defined template, use templateName as path to file
		if !util.FileExists(templateName) {
			s := fmt.Sprintf("Template file %s does not exist\n", templateName)
			return "", cli.NewExitError(s, 1)
		}
		data, err := ioutil.ReadFile(templateName)
		if err != nil {
			return "", cli.NewExitError(err.Error(), 1)
		}
		return string(data), nil
	}
}

// notifyOnErrors; in case of a pull, update or other error regarding the checkout and
// processing of contributions; will print a message and the path to the error log
// otherwise if will remove any existing 'errors.log' (old errors from last run).
func notifyOnErrors() {
	errorfile, err := os.Open("errors.log")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer errorfile.Close()
	fi, err := errorfile.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if fi.Size() > 0 {
		util.PrintInfoF(os.Stderr, "Some contributions could not be checked. See: errors.log", util.PI_ERROR)
	} else {
		os.Remove("errors.log")
	}
}

// Run will handle the functionality.
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

	// get template
	templateData, err := getTemplateOrExitError(ctx)
	if err != nil {
		return err
	}

	gontrib.PullSources = !ctx.GlobalBool("no-pull")

	// scan
	contributions, err := gontrib.ScanContributions(configuration)
	if err != nil {
		util.PrintInfo(os.Stderr, err.Error(), util.PI_ERROR)
		return cli.NewExitError(err.Error(), 1)
	}

	// define output
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

	util.PrintInfoF(os.Stdout, "\nReport saved in: %s", util.PI_INFO, outputPath)

	notifyOnErrors()

	return nil
}

// Create an example configuration file which the user can
// adapt to his own needs.
func cmdExconf(c *cli.Context) error {
	configuration := gontrib.Configuration{
		Emails: []string{"jubalh@openmailbox.org", "g.bluehut@gmail.com"},
		Projects: []gontrib.Project{
			{Name: "nudoku", Description: "Ncurses sudoku game", Gitrepos: []string{"https://github.com/jubalh/nudoku"}, Tags: []string{"C", "Terminal", "Game"}, Role: "maintainer"},
			{Name: "profanity", Description: "Ncurses based XMPP client", URL: "http://profanity.im/", Gitrepos: []string{"https://github.com/boothj5/profanity"}, Tags: []string{"C", "XMPP", "Terminal"}},
			{Name: "Funtoo", Description: "Linux distribution", URL: "http://funtoo.org/", Gitrepos: []string{"https://github.com/funtoo/ego", "https://github.com/funtoo/metro"}, MediaWikis: []mediawiki.MediaWiki{{BaseUrl: "http://funtoo.org", User: "jubalh"}}},
			{Name: "openSUSE", Description: "Linux distribution", URL: "http://opensuse.org/", Obs: []obs.OpenBuildService{{Apiurl: "https://api.opensuse.org", Repo: "utilities/vifm"}}},
			{Name: "prosody", Description: "A modern XMPP communication server", URL: "http://prosody.im/", Hgrepos: []string{"https://hg.prosody.im/prosody-modules"}, Tags: []string{"Lua", "XMPP"}},
		},
	}

	text, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(string(text))
	return nil
}

func cmdListTemplates(c *cli.Context) error {
	internalTemplates, err := AssetDir(templateAssetFolderName)
	if err != nil {
		return err
	}

	for _, t := range internalTemplates {
		fmt.Println(t[:strings.Index(t, ".")])
	}
	return nil
}

func cmdShowTemplate(c *cli.Context) error {
	template := c.String("name")

	data, err := Asset(filepath.Join(templateAssetFolderName, template) + ".html")
	if err != nil {
		return err
	}
	fmt.Print(string(data))

	return nil
}
