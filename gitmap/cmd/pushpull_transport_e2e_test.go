package cmd

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestApplyTransportFlag_E2E exercises the end-to-end transport
// rewrite using the real `git` binary against a temp bare repo.
// Skipped when git isn't available on PATH.
func TestApplyTransportFlag_E2E(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
	tmp := t.TempDir()
	bare := filepath.Join(tmp, "origin.git")
	work := filepath.Join(tmp, "work")
	mustGit(t, tmp, "init", "--bare", bare)
	mustGit(t, tmp, "clone", bare, work)

	// Seed remote.origin.url with a synthetic HTTPS URL so the
	// converter has a recognisable host/owner/repo to flip.
	mustGit(t, work, "remote", "set-url", "origin", "https://github.com/acme/widgets.git")

	t.Run("https_to_ssh_persists", func(t *testing.T) {
		changed, oldURL, newURL, err := ApplyTransportFlag(work, true, false)
		if err != nil {
			t.Fatalf("ApplyTransportFlag: %v", err)
		}
		if !changed {
			t.Fatalf("expected changed=true; got false (old=%s new=%s)", oldURL, newURL)
		}
		want := "git@github.com:acme/widgets.git"
		if got := readOrigin(t, work); got != want {
			t.Fatalf("origin url = %q want %q", got, want)
		}
	})

	t.Run("ssh_to_https_persists", func(t *testing.T) {
		changed, _, _, err := ApplyTransportFlag(work, false, true)
		if err != nil {
			t.Fatalf("ApplyTransportFlag: %v", err)
		}
		if !changed {
			t.Fatalf("expected changed=true")
		}
		want := "https://github.com/acme/widgets.git"
		if got := readOrigin(t, work); got != want {
			t.Fatalf("origin url = %q want %q", got, want)
		}
	})

	t.Run("noop_when_already_correct", func(t *testing.T) {
		// origin is HTTPS from previous subtest — asking for HTTPS
		// again should report no change and not rewrite.
		changed, _, _, err := ApplyTransportFlag(work, false, true)
		if err != nil {
			t.Fatalf("ApplyTransportFlag: %v", err)
		}
		if changed {
			t.Fatalf("expected changed=false for idempotent call")
		}
	})

	t.Run("ssh_wins_when_both_set", func(t *testing.T) {
		changed, _, newURL, err := ApplyTransportFlag(work, true, true)
		if err != nil {
			t.Fatalf("ApplyTransportFlag: %v", err)
		}
		if !changed || !strings.HasPrefix(newURL, "git@") {
			t.Fatalf("ssh did not win: changed=%v newURL=%s", changed, newURL)
		}
	})

	t.Run("unrecognised_url_fails_open", func(t *testing.T) {
		mustGit(t, work, "remote", "set-url", "origin", "file:///tmp/local-bare")
		changed, _, _, err := ApplyTransportFlag(work, true, false)
		if err != nil {
			t.Fatalf("ApplyTransportFlag should fail-open, got err: %v", err)
		}
		if changed {
			t.Fatalf("expected changed=false for unrecognised url")
		}
	})
}

// TestExtractTransportFlags covers every accepted alias spelling and
// confirms unrelated args pass through untouched.
func TestExtractTransportFlags(t *testing.T) {
	cases := []struct {
		name      string
		in        []string
		wantSSH   bool
		wantHTTPS bool
		wantRest  []string
	}{
		{"none", []string{"origin", "main"}, false, false, []string{"origin", "main"}},
		{"double_ssh", []string{"--ssh"}, true, false, []string{}},
		{"single_ssh", []string{"-ssh"}, true, false, []string{}},
		{"short_sh", []string{"--sh"}, true, false, []string{}},
		{"double_https", []string{"--https", "origin"}, false, true, []string{"origin"}},
		{"short_ht_then_ref", []string{"--ht", "main"}, false, true, []string{"main"}},
		{"both", []string{"--ssh", "--https"}, true, true, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, h, rest := extractTransportFlags(tc.in)
			if s != tc.wantSSH || h != tc.wantHTTPS {
				t.Fatalf("flags = (%v,%v), want (%v,%v)", s, h, tc.wantSSH, tc.wantHTTPS)
			}
			if strings.Join(rest, " ") != strings.Join(tc.wantRest, " ") {
				t.Fatalf("rest = %v, want %v", rest, tc.wantRest)
			}
		})
	}
}

func mustGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func readOrigin(t *testing.T, dir string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		t.Fatalf("read origin: %v", err)
	}
	return strings.TrimSpace(string(out))
}
