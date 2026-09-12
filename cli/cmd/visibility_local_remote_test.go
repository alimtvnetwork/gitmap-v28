// Package cmd — visibility_local_remote_test.go: regression coverage
// for the file:// / local-path warn-and-skip path added so CI fixtures
// backed by local bare repos never fail `make-public` with exit code 4
// (ExitVisBadProvider).
package cmd

import (
	"testing"
)

// TestIsLocalRemote_ClassifiesFileAndPathSchemes locks in the set of
// URL shapes treated as local (warn-skip) vs remote (classify normally).
func TestIsLocalRemote_ClassifiesFileAndPathSchemes(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"file:///tmp/fixture.git", true},
		{"file://C:/repos/x.git", true},
		{"FILE:///tmp/x", true},
		{"/srv/git/x.git", true},
		{"C:/repos/x.git", true},
		{"D:\\repos\\x.git", true},
		{"https://github.com/o/r.git", false},
		{"git@github.com:o/r.git", false},
		{"ssh://git@gitlab.com/o/r.git", false},
		{"", false},
	}

	for _, c := range cases {
		if got := isLocalRemote(c.url); got != c.want {
			t.Errorf("isLocalRemote(%q) = %v, want %v", c.url, got, c.want)
		}
	}
}
