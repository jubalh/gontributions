# Usage

## MediaWiki

*Example configuration*

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
