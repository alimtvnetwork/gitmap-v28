package cmd

import "testing"

// TestIsSSHURL — boundary classifier guard. Critical because the
// reclone fix in clonefixrepofoldertransport.go hinges on a binary
// SSH/non-SSH split; a wrong classification re-introduces the
// silent HTTPS downgrade the audit identified.
func TestIsSSHURL(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"git@github.com:owner/repo.git", true},
		{"ssh://git@github.com/owner/repo.git", true},
		{"SSH://git@github.com/owner/repo.git", true},
		{"https://github.com/owner/repo.git", false},
		{"http://example.com/x.git", false},
		{"", false},
		{"   ", false},
	}
	for _, c := range cases {
		got := isSSHURL(c.in)
		if got != c.want {
			t.Fatalf("isSSHURL(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// TestPreferExistingFolderTransport_NoDotGit returns the positional
// URL untouched when the destination has no `.git/` — the fresh
// clone case must not be perturbed.
func TestPreferExistingFolderTransport_NoDotGit(t *testing.T) {
	tmp := t.TempDir()
	got := preferExistingFolderTransport("https://github.com/owner/repo.git", tmp)
	if got != "https://github.com/owner/repo.git" {
		t.Fatalf("expected URL untouched, got %q", got)
	}
}

// TestRewriteToMatchExisting_SSHFromHTTPS proves the rewrite arm
// that closes the silent-downgrade bug: an SSH-origin folder + an
// HTTPS positional must produce an SSH URL.
func TestRewriteToMatchExisting_SSHFromHTTPS(t *testing.T) {
	got := rewriteToMatchExisting("https://github.com/owner/repo.git", true)
	want := "git@github.com:owner/repo.git"
	if got != want {
		t.Fatalf("rewriteToMatchExisting(https→ssh) = %q, want %q", got, want)
	}
}

// TestRewriteToMatchExisting_HTTPSFromSSH proves the symmetric arm:
// HTTPS-origin folder + SSH positional must rewrite to HTTPS.
func TestRewriteToMatchExisting_HTTPSFromSSH(t *testing.T) {
	got := rewriteToMatchExisting("git@github.com:owner/repo.git", false)
	want := "https://github.com/owner/repo.git"
	if got != want {
		t.Fatalf("rewriteToMatchExisting(ssh→https) = %q, want %q", got, want)
	}
}
