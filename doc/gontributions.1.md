% gontributions(1) Version 0.7 | contributions saver and lister

NAME
====

**gontributions** — OSS contributions saver and lister

SYNOPSIS
========

| **gontributions** \[OPTIONS] \[SUBCOMMAND]

DESCRIPTION
===========

This program lists your open source contributions. They get printed on your terminal when running the program. Additionally an HTML page gets created. There are some templates in the templates directory. Feel free to add your own and create a pull request to this repository.

Sometimes upstream repositories disappear. And with it your contributions. Initially I wrote this program as a way to back up all the repositories I contributed to and get a list of my contributions. Either because I was proud of them, or because I might wanted to look at them again. In any case, I just like backups. From there the program grew, to create a small overview in HTML format so you could show the projects you contribute to to others.

Options
-------

\--config

:   Set which config file to use

\--template

: Set which template to use. If it contains a dot it will be treated as a path to a user defined template. No dot means using an internal template

\--output

: Define name of the generated HTMl file

\--no-pull

: Don't update VCS repositories


SUBCOMMANDS
=====

*exconf*

: Show an example configuration file

*list-templates*

: List all built-in templates

*show-template*

: Display the content of a template

CONFIGURATION
===========

```
gontributions exconf > example.json
$EDITOR example.json
```

BUGS
====

See GitHub Issues: <https://github.com/jubalh/gontributions/issues>

AUTHOR
======

gontributions is written by by Michael Vetter <jubalh@iodoru.org>.
