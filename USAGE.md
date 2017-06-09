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
gontributions ships with two example templates: `default.html` and `detailed.html`.

The user can also create his own templates. For this he *needs* to set the environment variable `GONTRIB_TEMPLATES_PTH`.

gontributions will then check for the template file in that folder.
Here is an example:

```
gontributions exconf > example.json
mkdir my-awesome-templates && cd my-awesome-templates
vi superman.html
GONTRIB_TEMPLATES_PTH=my-awesome-templates gontributions --config example.json --template superman.html
```
