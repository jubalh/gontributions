package mediawiki

import "testing"

// This Wiki is almost dead.
// And LXQt has it's own Wiki now so I won't edit this anymore.
// This makes it save as a testing environment since the edits won't change.
const wikiURL = "http://wiki.lxde.org/en"
const wikiUsername = "Jubalh"
const expectedEdits = 8

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
