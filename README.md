# gontributions

*open source contributions lister*

This program lists your open source contributions.
They get printed on your terminal when running the program.
Additionally an HTML page gets created. There are some templates in the `templates` directory. Feel free to add your own and create a pull request to this repository.

You need to create a json file containing your open source contributions. The format is pretty easy. It is best to create an example configuration file and adapt it.

Create an example json configuration and adapt the example to your needs:

```
gontributions exconf
vi example_conf.json
```

Get an overview of your open source contributions:

```
gontributions --config example_conf.json
xdg-open output.html
```

Choose another template for your overview

```
gontributions --config config.json --template detailed.html --output another.html
xdg-open another.html
```

## Features
Search for commits in:
- git

Print out:
- project name
- description
- number of commits

## Todo
Support more version control systems:
- svn
- hg
- bzr
- obs

Add support to count Wiki edits
