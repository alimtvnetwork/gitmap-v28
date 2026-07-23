package cmd

// Locks the stricter expandHome contract documented in scanresolve.go:
// only the literal "~", "~/...", or "~\..." forms expand. Anything else
// (including "~foo", which some shells treat as "user foo's home")
// passes through verbatim because Go has no portable cross-platform
// resolver for the ~user form on Windows.
//
// Regression context: an earlier sshgenutil.go declared a looser version
// of expandHome that also expanded bare-prefix forms. The duplicate was
// removed (v3.76.1) and this test prevents the looser semantics from
// sneaking back in.

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skipf("no home dir available on this platform: %v", err)
	}

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"bare tilde expands to home", "~", home},
		{"forward slash subpath expands", "~/x", filepath.Join(home, "x")},
		{"backslash subpath expands", `~\x`, filepath.Join(home, "x")},
		{"tilde-user form is not expanded", "~foo", "~foo"},
		{"tilde-user with subpath is not expanded", "~foo/bar", "~foo/bar"},
		{"plain relative path passes through", "x", "x"},
		{"empty string passes through", "", ""},
		{"absolute path passes through", absoluteFixture(), absoluteFixture()},
		{"tilde in middle is not touched", "a/~/b", "a/~/b"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := expandHome(tc.in)
			if got != tc.want {
				t.Fatalf("expandHome(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// absoluteFixture returns a platform-appropriate absolute path so the
// "passes through" case is meaningful on both Windows and *nix runners.
func absoluteFixture() string {
	if runtime.GOOS == "windows" {
		return `C:\tmp\x`
	}

	return "/tmp/x"
}
