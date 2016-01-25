package mediawiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MediaWiki holds the base URL of the wiki page to which later the
// API call will get appended and the username to the wiki.
type MediaWiki struct {
	BaseUrl string
	User    string
}

// User field in MediaWikis response
type User struct {
	Id    int `json:"userid"`
	Name  string
	Edits int `json:"editcount"`
}

// Query field in MediaWikis response
type Query struct {
	Users []User `json:"users"`
}

//MediaWikis response
type Response struct {
	Query Query
}

// GetUserEdits calls wikiUrl MediaWiki API to retrieve the number of edits
// the user username has done.
func GetUserEdits(wikiUrl string, username string) (int, error) {
	request := fmt.Sprintf("%s/api.php?action=query&list=users&format=json&usprop=editcount&ususers=%s", wikiUrl, username)
	resp, err := http.Get(request)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var r Response
	json.Unmarshal(body, &r)

	return r.Query.Users[0].Edits, nil
}
