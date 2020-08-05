package mediawiki

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// MediaWiki holds the base URL of the wiki page to which later the
// API call will get appended and the username to the wiki.
type MediaWiki struct {
	BaseURL string
	User    string
}

// User field in MediaWikis response
type User struct {
	ID    int `json:"userid"`
	Name  string
	Edits int `json:"editcount"`
}

// Query field in MediaWikis response
type Query struct {
	Users []User `json:"users"`
}

// GetUserEdits calls wikiURL MediaWiki API to retrieve the number of edits
// the user username has done.
func GetUserEdits(wikiURL string, username string) (count int, err error) {
	wikiURL, err := url.Parse(wikiURL)
	if err != nil {
		return 0, errors.New("Not a valid URL")
	}
	wikiURL.Path += "/api.php"
	parameters := url.Values{}
	parameters.Add("action", "query")
	parameters.Add("list", "users")
	parameters.Add("format", "json")
	parameters.Add("usprop", "editcount")
	parameters.Add("ususers", username)
	wikiURL.RawQuery = parameters.Encode()

	resp, err := http.Get(wikiURL.String())
	if err != nil {
		return 0, errors.New("Not able to HTTP GET")
	}
	defer resp.Body.Close()

	// MediaWikis response as anon struct
	var response struct {
		Query Query
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return 0, errors.New("Not able to decode JSON")
	}

	if len(response.Query.Users) > 0 {
		count = response.Query.Users[0].Edits
	} else {
		err = errors.New("Did not get a 'user' returned")
	}

	return
}
