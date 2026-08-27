package release

import (
	"testing"
)

func TestParseVersionFromCommit(t *testing.T) {
	cases := map[string]string{
		"bump version to 1.2.0":                  "1.2.0",
		"chore(release): bump version to v1.2.0": "v1.2.0",
		"release v1.14.0":                        "v1.14.0",
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
	v, isFound := ParseVersionFromCommit("fix bug")
	if isFound || v != "" {
		t.Errorf("Expected not found, got %v", v)
	}
}
