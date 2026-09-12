package release

import (
	"testing"
)

func TestParseVersionFromCommit(t *testing.T) {
	cases := map[string]string{
		"bump version to 1.2.0":                  "1.2.0",
		"chore(release): bump version to v1.2.0": "v1.2.0",
		"chore(release): v1.2.0":                 "v1.2.0",
		"release v1.14.0":                        "v1.14.0",
		"release: v6.200.0 description":          "v6.200.0",
		"release(all): v1.2.0":                   "v1.2.0",
		"release(pkg): 1.2.0":                    "1.2.0",
		"bump to v1.2.0":                         "v1.2.0",
		"bump: v1.2.0":                           "v1.2.0",
		"v1.2.0":                                 "v1.2.0",
		"release v1.2.3-rc.1":                    "v1.2.3-rc.1",
	}

	for msg, expected := range cases {
		v, isFound := ParseVersionFromCommit(msg)
		if !isFound || v != expected {
			t.Errorf("For %q, got (%q, %v)", msg, v, isFound)
		}
	}
}

func TestParseVersionFromCommit_NotFound(t *testing.T) {
	cases := []string{
		"fix bug",
		"fix: support go 1.22 in CI",
		"feat: upgrade to v1.2.0 library",
		"docs: mention v1.2.0 in readme",
	}

	for _, msg := range cases {
		v, isFound := ParseVersionFromCommit(msg)
		if isFound || v != "" {
			t.Errorf("Expected not found for %q, got (%q, %v)", msg, v, isFound)
		}
	}
}
