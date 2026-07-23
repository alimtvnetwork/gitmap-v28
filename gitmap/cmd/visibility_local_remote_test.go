// Package cmd — visibility_local_remote_test.go: regression coverage
// for the file:// / local-path warn-and-skip path added so CI fixtures
// backed by local bare repos never fail `make-public` with exit code 4
// (ExitVisBadProvider).
//
// Two layers of coverage:
//
//  1. Pure-func table test for isLocalRemote — cheap, exhaustive.
//  2. Subprocess exec of resolveProviderAndSlugOrExit — proves the
//     wired-in behavior (stderr message + exit code 0, NOT 4) without
//     spawning the full CLI binary.
package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

// TestResolveProviderAndSlug_LocalRemote_ExitsZero re-execs this test
// binary with GITMAP_TEST_RESOLVE_LOCAL=<url> and asserts the child
// exits with ExitVisOK (0), never ExitVisBadProvider (4), and emits
// the local-skip stderr message. This is the exact behavior CI relies
// on for file:// fixture remotes.
func TestResolveProviderAndSlug_LocalRemote_ExitsZero(t *testing.T) {
	if url := os.Getenv("GITMAP_TEST_RESOLVE_LOCAL"); url != "" {
		// Child process — must os.Exit itself.
		resolveProviderAndSlugOrExit(url)
		// If control returns here, the guard did NOT fire; force a
		// non-zero exit distinct from 0 and 4 so the parent flags it.
		os.Exit(99)
	}

	urls := []string{
		"file:///tmp/gitmap-fixture.git",
		"/var/tmp/local-bare.git",
		"C:/repos/local-bare.git",
	}
	for _, url := range urls {
		url := url
		t.Run(url, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestResolveProviderAndSlug_LocalRemote_ExitsZero")
			cmd.Env = append(os.Environ(), "GITMAP_TEST_RESOLVE_LOCAL="+url)
			out, err := cmd.CombinedOutput()

			exitCode := 0
			if err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					exitCode = ee.ExitCode()
				} else {
					t.Fatalf("exec failed: %v\noutput:\n%s", err, out)
				}
			}

			if exitCode == constants.ExitVisBadProvider {
				t.Fatalf("local remote %q wrongly rejected with ExitVisBadProvider (4)\noutput:\n%s", url, out)
			}
			if exitCode != constants.ExitVisOK {
				t.Fatalf("local remote %q exited with %d, want ExitVisOK (0)\noutput:\n%s", url, exitCode, out)
			}
			if !strings.Contains(string(out), "skipping local remote") {
				t.Errorf("expected local-skip stderr message, got:\n%s", out)
			}
		})
	}
}
