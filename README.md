# gontributions

*open source contributions lister*

This program lists your open source contributions.
They get printed on your terminal when running the program.
Additionally an HTML page gets created. There are some templates in the `templates` directory. Feel free to add your own and create a pull request to this repository.

You need to create a json file containing your open source contributions. The format is pretty easy. It is best to create an example configuration file and adapt it.

Create an example json configuration and adapt the example to your needs:

```
gontributions exconf > example_conf.json
vi example_conf.json
```

Get an overview of your open source contributions:

```
gontributions --config example_conf.json
xdg-open output.html
```

Choose another template for your overview

```
gontributions --config gontrib.json --template detailed.html --output another.html
xdg-open another.html
```

If you don't specify a configuration file it will automatically look for `gontrib.json`.

Per default it expects a `templates/` directory to exist in which it will look for the template specified via the `--template` switch, or `default.html` if not specified.

The **GONTRIB_TEMPLATES_PTH** environment variable can be used to change the path in which it will look for this.

## Features
Search for commits in:
- [x] Git
- [ ] Subversion
- [ ] Mercurial
- [ ] Bazaar
- [x] Open Build Service

- [x] support to count Wiki edits

Display:
- [x] project name
- [x] description
- [x] project URL
- [x] number of commits/contributions

