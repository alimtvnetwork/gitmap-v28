// Package cmd — chromeprofile_resolve_test.go: edge cases for
// resolving user-supplied profile identifiers (dir name OR display
// name) to an on-disk path, plus the human-readable summary used in
// error banners.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupFakeChromeRoot(t *testing.T, profiles map[string]string) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", root)
	cache := map[string]map[string]any{}
	for dir, display := range profiles {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o700); err != nil {
			t.Fatal(err)
		}
		cache[dir] = map[string]any{"name": display}
	}
	state := map[string]any{
		"profile": map[string]any{"info_cache": cache},
	}
	raw, _ := json.Marshal(state)
	if err := os.WriteFile(filepath.Join(root, "Local State"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestResolveChromeProfileByDirectoryName(t *testing.T) {
	root := setupFakeChromeRoot(t, map[string]string{"Profile 15": "Lovable"})
	res, ok := resolveChromeProfile("Profile 15")
	if !ok || res.Path != filepath.Join(root, "Profile 15") {
		t.Fatalf("dir lookup failed: %+v ok=%v", res, ok)
	}
	if res.DisplayName != "Lovable" || res.Dir != "Profile 15" {
		t.Fatalf("display enrichment lost: %+v", res)
	}
}

func TestResolveChromeProfileByDisplayNameCaseInsensitive(t *testing.T) {
	setupFakeChromeRoot(t, map[string]string{"Profile 15": "Lovable"})
	res, ok := resolveChromeProfile("  LOVABLE  ")
	if !ok || res.Dir != "Profile 15" || res.DisplayName != "Lovable" {
		t.Fatalf("display lookup failed: %+v ok=%v", res, ok)
	}
}

func TestResolveChromeProfileAbsolutePathPassthrough(t *testing.T) {
	dir := t.TempDir()
	res, ok := resolveChromeProfile(dir)
	if !ok || res.Path != dir {
		t.Fatalf("abs passthrough: %+v ok=%v", res, ok)
	}
}

func TestResolveChromeProfileAbsolutePathMissingReturnsFalse(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "ghost")
	res, ok := resolveChromeProfile(missing)
	if ok {
		t.Fatalf("missing abs should be !ok: %+v", res)
	}
}

func TestResolveChromeProfileUnknownReturnsFalse(t *testing.T) {
	setupFakeChromeRoot(t, map[string]string{"Profile 1": "Work"})
	if res, ok := resolveChromeProfile("nope"); ok {
		t.Fatalf("unknown should be !ok: %+v", res)
	}
}

func TestResolveChromeProfileDirThinWrapper(t *testing.T) {
	root := setupFakeChromeRoot(t, map[string]string{"Default": "Personal"})
	got, ok := resolveChromeProfileDir("Personal")
	if !ok || got != filepath.Join(root, "Default") {
		t.Fatalf("wrapper: got=%q ok=%v", got, ok)
	}
}

func TestChromeProfileDestinationCarriesDisplayName(t *testing.T) {
	setupFakeChromeRoot(t, map[string]string{"Profile 2": "Side"})
	res := chromeProfileDestination("Profile 2")
	if res.DisplayName != "Side" || res.Dir != "Profile 2" {
		t.Fatalf("destination enrichment: %+v", res)
	}
}

func TestChromeProfileSummaryFormats(t *testing.T) {
	cases := []struct {
		name string
		in   chromeProfileResolution
		want string
	}{
		{"display+dir", chromeProfileResolution{Dir: "Profile 15", DisplayName: "Lovable"}, "Lovable (dir: Profile 15)"},
		{"dir only", chromeProfileResolution{Dir: "Default"}, "Default"},
		{"display equals dir", chromeProfileResolution{Dir: "Default", DisplayName: "Default"}, "Default"},
		{"input fallback", chromeProfileResolution{Input: "raw"}, "raw"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := chromeProfileSummary(tc.in); got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestAvailableChromeProfileNamesFiltersNonProfileDirs(t *testing.T) {
	setupFakeChromeRoot(t, map[string]string{
		"Default":   "Personal",
		"Profile 1": "Work",
	})
	root := os.Getenv("GITMAP_CHROME_USER_DATA")
	_ = os.MkdirAll(filepath.Join(root, "Crashpad"), 0o700)
	_ = os.WriteFile(filepath.Join(root, "NotADir"), []byte(""), 0o600)

	names := availableChromeProfileNames()
	if len(names) != 2 {
		t.Fatalf("filter: got %v", names)
	}
	joined := strings.Join(names, ",")
	if !strings.Contains(joined, "Default") || !strings.Contains(joined, "Profile 1") {
		t.Fatalf("missing expected entries: %v", names)
	}
}

func TestReadChromeLocalStateMissingReturnsNil(t *testing.T) {
	t.Setenv("GITMAP_CHROME_USER_DATA", t.TempDir())
	if s := readChromeLocalState(); s != nil {
		t.Fatalf("expected nil, got %+v", s)
	}
}

func TestReadChromeLocalStateMalformedReturnsNil(t *testing.T) {
	root := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", root)
	_ = os.WriteFile(filepath.Join(root, "Local State"), []byte("{not-json"), 0o600)
	if s := readChromeLocalState(); s != nil {
		t.Fatalf("expected nil for malformed, got %+v", s)
	}
}
