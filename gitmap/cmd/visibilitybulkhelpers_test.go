package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestOwnerFromSlug(t *testing.T) {
	cases := map[string]string{
		"alice/repo":      "alice",
		"alice/repo-v12":  "alice",
		"bare":            "bare",
		"":                "",
		"org/sub/project": "org",
	}
	for in, want := range cases {
		if got := ownerFromSlug(in); got != want {
			t.Fatalf("ownerFromSlug(%q)=%q want %q", in, got, want)
		}
	}
}

func TestProviderHost(t *testing.T) {
	if got := providerHost(constants.ProviderGitLab); got != constants.HostGitLab {
		t.Fatalf("gitlab host=%q want %q", got, constants.HostGitLab)
	}

	if got := providerHost(constants.ProviderGitHub); got != constants.HostGitHub {
		t.Fatalf("github host=%q want %q", got, constants.HostGitHub)
	}

	if got := providerHost("unknown"); got != constants.HostGitHub {
		t.Fatalf("unknown provider should default to github, got %q", got)
	}
}
