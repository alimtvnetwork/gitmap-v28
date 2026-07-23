package cmd

// Focused tests for the --backend / --name filters added to
// `startup-list`. The filter logic lives in startuplistfilter.go
// and is OS-agnostic (operates on already-collected Entry values),
// so these tests run on every CI runner.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// fixtureStartupEntries returns one entry for each shape startup.List can
// produce so a single slice covers Linux/.desktop, macOS/.plist,
// Windows/registry, and Windows/startup-folder.
func fixtureStartupEntries() []startup.Entry {
	return []startup.Entry{
		{Name: "gitmap-watch.desktop",
			Path: "/home/me/.config/autostart/gitmap-watch.desktop",
			Exec: "/usr/local/bin/gitmap watch"},
		{Name: "gitmap.watch.plist",
			Path: "/Users/me/Library/LaunchAgents/gitmap.watch.plist",
			Exec: "/usr/local/bin/gitmap watch"},
		{Name: "gitmap-watch",
			Path: `HKCU\Software\Microsoft\Windows\CurrentVersion\Run\gitmap-watch`,
			Exec: `C:\gitmap.exe watch`},
		{Name: "gitmap-watch",
			Path: `HKLM\Software\Microsoft\Windows\CurrentVersion\Run\gitmap-watch`,
			Exec: `C:\gitmap.exe watch`},
		{Name: "gitmap-watch.lnk",
			Path: `C:\Users\me\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\gitmap-watch.lnk`,
			Exec: ""},
	}
}

// TestFilterStartupList_NoFilters returns every entry untouched and
// in the original order — important because the renderer iterates
// the slice as-is and we don't want a hidden reorder.
func TestFilterStartupList_NoFilters(t *testing.T) {
	in := fixtureStartupEntries()
	got := filterStartupList(in, "", "")
	if len(got) != len(in) {
		t.Fatalf("len = %d, want %d", len(got), len(in))
	}
	for i := range in {
		if got[i].Path != in[i].Path {
			t.Errorf("entry %d reordered: got %q want %q", i, got[i].Path, in[i].Path)
		}
	}
}

// TestFilterStartupList_BackendRegistry keeps only the entry whose
// Path starts with `HKCU\` — the discriminator runValuePath emits
// in startup/winregistry_windows.go.
func TestFilterStartupList_BackendRegistry(t *testing.T) {
	got := filterStartupList(fixtureStartupEntries(), "registry", "")
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1; got %#v", len(got), got)
	}
	if got[0].Name != "gitmap-watch" {
		t.Errorf("name = %q, want gitmap-watch", got[0].Name)
	}
}

// TestFilterStartupList_BackendStartupFolder keeps only the .lnk
// entry — the `.lnk` extension is the cross-platform discriminator
// the .desktop / .plist entries can never share.
func TestFilterStartupList_BackendStartupFolder(t *testing.T) {
	got := filterStartupList(fixtureStartupEntries(), "startup-folder", "")
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1; got %#v", len(got), got)
	}
	if got[0].Name != "gitmap-watch.lnk" {
		t.Errorf("name = %q, want gitmap-watch.lnk", got[0].Name)
	}
}

// TestFilterStartupList_NameMatchesAcrossOSes confirms a single
// --name value (`watch`) matches all four OS-specific entry shapes
// produced by Add — proving the prefix/suffix stripping in
// logicalEntryName covers every Add code path.
func TestFilterStartupList_NameMatchesAcrossOSes(t *testing.T) {
	got := filterStartupList(fixtureStartupEntries(), "", "watch")
	if len(got) != 5 {
		t.Fatalf("len = %d, want 5; got %#v", len(got), got)
	}
}

// TestFilterStartupList_BackendRegistryHKLM keeps only the entry
// whose Path starts with `HKLM\` — the discriminator runValuePathFor
// emits for the machine-wide registry-hklm backend.
func TestFilterStartupList_BackendRegistryHKLM(t *testing.T) {
	got := filterStartupList(fixtureStartupEntries(), "registry-hklm", "")
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1; got %#v", len(got), got)
	}
	if !strings.HasPrefix(got[0].Path, `HKLM\`) {
		t.Errorf("path = %q, want HKLM-rooted", got[0].Path)
	}
}

// TestFilterStartupList_BackendAndName confirms the AND semantics
// matching the user-facing contract `--backend=registry --name=watch`
// returns ONLY the registry entry for that name.
func TestFilterStartupList_BackendAndName(t *testing.T) {
	got := filterStartupList(fixtureStartupEntries(), "registry", "watch")
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1; got %#v", len(got), got)
	}
	if got[0].Name != "gitmap-watch" {
		t.Errorf("name = %q, want gitmap-watch", got[0].Name)
	}
}

// TestFilterStartupList_NameNoMatch returns an empty (non-nil)
// slice when nothing matches — important for renderers that
// distinguish "filtered to zero" from "List failed".
func TestFilterStartupList_NameNoMatch(t *testing.T) {
	got := filterStartupList(fixtureStartupEntries(), "", "does-not-exist")
	if got == nil {
		t.Fatal("got nil slice; want empty non-nil")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

// TestParseStartupListFlags_BackendValidation covers the four
// outcomes: empty (no filter, ok), valid registry, valid
// startup-folder, invalid (rejected). Same shape as the existing
// json-indent table test for visual parity.
func TestParseStartupListFlags_BackendValidation(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"empty", []string{}, false},
		{"registry", []string{"--backend=registry"}, false},
		{"registry-hklm", []string{"--backend=registry-hklm"}, false},
		{"startup-folder", []string{"--backend=startup-folder"}, false},
		{"unknown", []string{"--backend=hkcu"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseStartupListFlags(tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("args=%v: wantErr=%v got %v", tc.args, tc.wantErr, err)
			}
		})
	}
}

// TestParseStartupListFlags_NameAcceptsAnyString confirms --name is
// passed through as-is (no validation) — empty / typical / weird
// values all parse, because filtering a name that doesn't exist is
// a legitimate user query that should produce zero rows, not an
// error.
func TestParseStartupListFlags_NameAcceptsAnyString(t *testing.T) {
	for _, name := range []string{"", "watch", "gitmap-watch", "name with spaces"} {
		opts, err := parseStartupListFlags([]string{"--name=" + name})
		if err != nil {
			t.Fatalf("name=%q: err = %v", name, err)
		}
		if opts.name != name {
			t.Errorf("name=%q: opts.name = %q", name, opts.name)
		}
	}
}
