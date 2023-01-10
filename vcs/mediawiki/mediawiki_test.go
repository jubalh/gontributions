package mediawiki

import "testing"

// This Wiki is almost dead.
// And LXQt has it's own Wiki now so I won't edit this anymore.
// This makes it save as a testing environment since the edits won't change.
const wikiURL = "https://www.funtoo.org"
const invalidURL = "hobbit://mordor.me"
const noURL = "://.."
const wikiUsername = "Jubalh"
const expectedEdits = 93
const expErrorInvalidURL = "Not a valid URL"
const expErrorNoHTTPGet = "Not able to HTTP GET"

func TestGetUserEdits(t *testing.T) {
	t.Logf("Querying LXDE wiki for user 'Jubalh'")
	count, err := GetUserEdits(wikiURL, wikiUsername)
	if err != nil {
		t.Error("Error: ", err)
		t.FailNow()
	}
	if count != expectedEdits {
		t.Errorf("GetUserEdits returned: %d, expected: %d", count, expectedEdits)
	}
}

func TestGetUserEditsInvalidURL(t *testing.T) {
	t.Logf("Querying invalid URL")
	_, err := GetUserEdits(invalidURL, wikiUsername)

	if err == nil {
		t.Errorf("Should have failed because used invalid URL: %s", invalidURL)
		t.FailNow()
	}

	if err.Error() != expErrorNoHTTPGet {
		t.Errorf("Should have failed because with '%s' but got: '%s'", expErrorNoHTTPGet, err.Error())
		t.FailNow()
	}
}

func TestGetUserEditsNoHTTPGet(t *testing.T) {
	t.Logf("Parse a non-URL")
	_, err := GetUserEdits(noURL, wikiUsername)

	if err == nil {
		t.Errorf("Should have failed because used URL: %s", invalidURL)
		t.FailNow()
	}

	if err.Error() != expErrorInvalidURL {
		t.Errorf("Should have failed because with '%s' but got: '%s'", expErrorInvalidURL, err.Error())
		t.FailNow()
	}
}
