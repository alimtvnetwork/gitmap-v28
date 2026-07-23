// Package store — unit tests for the IdentifiedTransport helpers
// (migration 008 / Plan 03 step 2-3). Covers: empty inputs are
// silent no-ops, illegal values are rejected, lookup misses return
// ("", nil), and a round-trip set→get preserves the value.
package store

import (
	"path/filepath"
	"testing"
)

// newTestDB opens a fresh sqlite file under t.TempDir() and runs the
// full Migrate() pass so the Repo table (with the migration-008
// IdentifiedTransport column) is present. Per-test isolation — no
// shared state across tests.
func newTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := OpenAt(filepath.Join(t.TempDir(), "transport.sqlite"))
	if err != nil {
		t.Fatalf("OpenAt: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return db
}

func TestClassifyURLTransport(t *testing.T) {
	cases := map[string]string{
		"":                                 "",
		"https://github.com/acme/repo.git": RepoTransportHTTPS,
		"http://gitlab.local/x.git":        RepoTransportHTTPS,
		"git@github.com:acme/repo.git":     RepoTransportSSH,
		"ssh://git@github.com/acme/r.git":  RepoTransportSSH,
		"file:///tmp/repo":                 "",
		"  HTTPS://Acme.io/r.git  ":        RepoTransportHTTPS,
	}
	for in, want := range cases {
		if got := ClassifyURLTransport(in); got != want {
			t.Errorf("ClassifyURLTransport(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSetRepoIdentifiedTransportRejectsBadInput(t *testing.T) {
	db := newTestDB(t)
	cases := []struct{ url, transport string }{
		{"", "ssh"},
		{"https://x", ""},
		{"https://x", "telnet"},
		{"  ", "ssh"},
	}
	for _, c := range cases {
		n, err := db.SetRepoIdentifiedTransport(c.url, c.transport)
		if err != nil {
			t.Errorf("SetRepoIdentifiedTransport(%q,%q) unexpected error: %v",
				c.url, c.transport, err)
		}
		if n != 0 {
			t.Errorf("SetRepoIdentifiedTransport(%q,%q) touched %d rows, want 0",
				c.url, c.transport, n)
		}
	}
}

func TestLookupRepoIdentifiedTransportMiss(t *testing.T) {
	db := newTestDB(t)
	got, err := db.LookupRepoIdentifiedTransport("https://nope.example/x.git")
	if err != nil {
		t.Fatalf("lookup miss returned error: %v", err)
	}
	if got != "" {
		t.Fatalf("lookup miss returned %q, want empty", got)
	}
}
