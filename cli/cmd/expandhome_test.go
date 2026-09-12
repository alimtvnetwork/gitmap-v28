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
	"runtime"
	"strings"
	"testing"
)

func TestExpandHome(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		validate func(t *testing.T, got string)
	}{
		{"bare tilde expands to home", "~", func(t *testing.T, got string) {
			if strings.HasPrefix(got, "~") {
				t.Fatalf("expandHome(~) = %q, want non-tilde path", got)
			}
		}},
		{"forward slash subpath expands", "~/x", func(t *testing.T, got string) {
			if strings.HasPrefix(got, "~") || !strings.HasSuffix(got, "x") {
				t.Fatalf("expandHome(~/x) = %q, want path ending in x", got)
			}
		}},
		{"backslash subpath expands", `~\x`, func(t *testing.T, got string) {
			if strings.HasPrefix(got, "~") || !strings.HasSuffix(got, "x") {
				t.Fatalf(`expandHome(~\x) = %q, want path ending in x`, got)
			}
		}},
		{"tilde-user form is not expanded", "~foo", func(t *testing.T, got string) {
			if got != "~foo" {
				t.Fatalf("expandHome(~foo) = %q, want ~foo", got)
			}
		}},
		{"tilde-user with subpath is not expanded", "~foo/bar", func(t *testing.T, got string) {
			if got != "~foo/bar" {
				t.Fatalf("expandHome(~foo/bar) = %q, want ~foo/bar", got)
			}
		}},
		{"plain relative path passes through", "x", func(t *testing.T, got string) {
			if got != "x" {
				t.Fatalf("expandHome(x) = %q, want x", got)
			}
		}},
		{"empty string passes through", "", func(t *testing.T, got string) {
			if got != "" {
				t.Fatalf("expandHome(\"\") = %q, want \"\"", got)
			}
		}},
		{"absolute path passes through", absoluteFixture(), func(t *testing.T, got string) {
			if got != absoluteFixture() {
				t.Fatalf("expandHome(%q) = %q, want %q", absoluteFixture(), got, absoluteFixture())
			}
		}},
		{"tilde in middle is not touched", "a/~/b", func(t *testing.T, got string) {
			if got != "a/~/b" {
				t.Fatalf("expandHome(a/~/b) = %q, want a/~/b", got)
			}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := expandHome(tc.in)
			tc.validate(t, got)
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
