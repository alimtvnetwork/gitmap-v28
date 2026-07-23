// Package cmd — chromeprofile_resolve.go: resolves user-supplied Chrome
// profile identifiers (directory name like "Profile 1" OR display name
// like "Lovable" shown in Chrome's profile picker) to an on-disk path.
//
// Chrome stores the human-readable display name in
// <UserData>/Local State under profile.info_cache[<dir>].name. This
// file reads that index so users can pass the same name they see in
// Chrome instead of guessing "Profile N".
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// chromeLocalState is the minimal shape we need from Chrome's
// "Local State" JSON. Extra fields are ignored.
type chromeLocalState struct {
	Profile struct {
		InfoCache map[string]struct {
			Name string `json:"name"`
		} `json:"info_cache"`
	} `json:"profile"`
}

// chromeProfileEntry pairs a profile directory with its display name.
type chromeProfileEntry struct {
	Dir         string
	DisplayName string
}

type chromeProfileResolution struct {
	Input       string
	Path        string
	Dir         string
	DisplayName string
}

// readChromeLocalState loads <UserData>/Local State, returning nil on
// any I/O or parse error so callers degrade gracefully to dir-only.
func readChromeLocalState() *chromeLocalState {
	path := filepath.Join(chromeUserDataDir(), "Local State")
	raw, err := os.ReadFile(path) //nolint:gosec // user-data path
	if err != nil {
		return nil
	}
	var s chromeLocalState
	if json.Unmarshal(raw, &s) != nil {
		return nil
	}
	return &s
}

// chromeProfileEntries returns every profile directory that exists on
// disk, paired with its display name when known.
func chromeProfileEntries() []chromeProfileEntry {
	dirs := availableChromeProfileNames()
	state := readChromeLocalState()
	out := make([]chromeProfileEntry, 0, len(dirs))
	for _, d := range dirs {
		display := ""
		if state != nil {
			if info, ok := state.Profile.InfoCache[d]; ok {
				display = info.Name
			}
		}
		out = append(out, chromeProfileEntry{Dir: d, DisplayName: display})
	}
	return out
}

// resolveChromeProfileDir maps a user-supplied identifier to an
// absolute on-disk path. Resolution order:
//  1. Absolute path → returned as-is.
//  2. Literal directory name ("Default", "Profile 1") → joined.
//  3. Display name from Local State (case-insensitive) → joined.
//
// Returns (path, ok). When ok is false, callers should print the
// available list and exit ExitChromeProfileNotFound.
func resolveChromeProfile(name string) (chromeProfileResolution, bool) {
	if filepath.IsAbs(name) {
		res := chromeProfileFromPath(name, name)
		return res, chromeProfilePathExists(name)
	}
	direct := filepath.Join(chromeUserDataDir(), name)
	if chromeProfilePathExists(direct) {
		return chromeProfileFromPath(name, direct), true
	}
	return resolveChromeProfileDisplayName(name, direct)
}

func resolveChromeProfileDir(name string) (string, bool) {
	res, ok := resolveChromeProfile(name)
	return res.Path, ok
}

func chromeProfileDestination(name string) chromeProfileResolution {
	return chromeProfileFromPath(name, chromeProfilePath(name))
}

func chromeProfileFromPath(input, path string) chromeProfileResolution {
	dir := filepath.Base(path)
	return chromeProfileResolution{Input: input, Path: path, Dir: dir, DisplayName: chromeProfileDisplayName(dir)}
}

func chromeProfileDisplayName(dir string) string {
	state := readChromeLocalState()
	if state == nil {
		return ""
	}
	if info, ok := state.Profile.InfoCache[dir]; ok {
		return info.Name
	}
	return ""
}

func resolveChromeProfileDisplayName(name, direct string) (chromeProfileResolution, bool) {
	state := readChromeLocalState()
	if state == nil {
		return chromeProfileResolution{Input: name, Path: direct, Dir: filepath.Base(direct)}, false
	}
	want := strings.ToLower(strings.TrimSpace(name))
	for dir, info := range state.Profile.InfoCache {
		if strings.ToLower(strings.TrimSpace(info.Name)) == want {
			p := filepath.Join(chromeUserDataDir(), dir)
			if chromeProfilePathExists(p) {
				return chromeProfileResolution{Input: name, Path: p, Dir: dir, DisplayName: info.Name}, true
			}
		}
	}
	return chromeProfileResolution{Input: name, Path: direct, Dir: filepath.Base(direct)}, false
}

func chromeProfileSummary(p chromeProfileResolution) string {
	if p.DisplayName != "" && p.DisplayName != p.Dir {
		return fmt.Sprintf("%s (dir: %s)", p.DisplayName, p.Dir)
	}
	if p.Dir != "" {
		return p.Dir
	}
	return p.Input
}

// printAvailableChromeProfilesWithDisplay writes a "did you mean…"
// stderr block that includes the display name alongside the directory
// so users can match what Chrome's profile picker shows.
func printAvailableChromeProfilesWithDisplay() {
	root := chromeUserDataDir()
	entries := chromeProfileEntries()
	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "  available profiles under %s: (none found)\n", root)
		return
	}
	fmt.Fprintf(os.Stderr, "  available profiles under %s:\n", root)
	for _, e := range entries {
		if e.DisplayName != "" {
			fmt.Fprintf(os.Stderr, "    - %s  (display: %q)\n", e.Dir, e.DisplayName)
			continue
		}
		fmt.Fprintf(os.Stderr, "    - %s\n", e.Dir)
	}
}
