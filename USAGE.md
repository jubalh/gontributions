# Usage

## MediaWiki

**Example configuration**

```
{
    "Projects": [
        {
            "Name": "Wikipedia",
            "Description": "Contributions to Wikipedia, the free encyclopedia",
            "URL": "http://wikipedia.org/",
            "MediaWikis": [
                {
                    "BaseUrl": "http://en.wikipedia.org/w",
                    "User": "WhisperToMe"
                }
            ]
        }
    ]
}
```

Above is an example minimal configuration that just checks MediaWiki contributions, in this case for Wikipedia.

*Note:* Some instances of MediaWiki have their `api.php` directly at `someurl.example`. The ones that follow the recommendation of MediaWiki however have it at `someurl.example/w`. Thus you need to find out where the api of the MediaWiki instance which you want to query lies and set this as your `BaseUrl`.

## Own templates
gontributions ships with three example templates: `default`, `detailed` and `fancy`.

The user can also create his own templates. For this he just has to specify the path to them via `--template`.
If it contains no dot *gontributions* will check in it internal assets. If it contains a dot it is treated as a user template

Here is an example:

```
gontributions exconf > example.json
gontributions --config example.json --template fancy
```

vs

```
gontributions exconf > example.json
gontributions show-template --name fancy > superman-template.html
vi superman-template.html
gontributions --config example.json --template superman-template.html
```

You can use `list-templates` to show you which templates are built in, and `show-template --name default` to print it on standard output.
