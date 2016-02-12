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

## Contribute
There are some features that would be nice to have, but i don't need them myself. I created [issues](https://github.com/jubalh/gontributions/issues) about them.
Feel free to work on them or implement your own ideas.

I try to have the `master` branch clear and only have working versions there. Please work on the `develop` branch and create a new `feature/count-my-hair` for every feature you implement.

## FAQ
**Does it have to be Open Source?**

You could also add private repositories and just publish the finished report. Noone will learn about the location of the repos, your precious commits or your username. Still you would have a nice overview of what you did.

Furthermore: if someone deletes a repository you contributed to, you still have your local copy if you use *gontributions* which will count into your contributions report. If we would just query some remote server you don't have any control over it. When it's deleted it's just gone. And your work lost.

**But there is 'GitHub Contributions'!**

Right. There is. But:

- It only displays contributions in a certain timeframe, you have no control over it.
- It works only for GitHub. What about your contributions hosted at GitLab, BitBucket or on sourceforge?
- What about your non-git Contributions? Like wiki edits, subversion commits etc.

**But there is Ohloh/Open Hub!**

Right. There is. But:

- The service hangs often so you need to ping admins to restart the scanning tool.
- You rely on a service you don't have much control over
- See about deleted repos above.

**What else is cool about it?**

Quite a few people use their personal website as their resume.  
With this tool you can add an overview of the work you did and you can adjust the look and feel to the rest of your website.
If you rely on Open Hub it will always look like Open Hub, if you create your own *template* for *gontributions* your report has totally your touch.

**Can I see such a report please?**

Glad you ask. Of course you can!  
Take a look at [mine](http://iodoru.org/gontrib).
