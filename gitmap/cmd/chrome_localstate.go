// Package cmd — pure-Go Local State parser shared by chrome which/list/etc.
// Kept separate from chrome_which.go so it is trivially unit-testable
// without exec'ing Chrome or hitting the user's real profile directory.
package cmd

import "encoding/json"

// ChromeLocalState is the slice of Chrome's Local State JSON gitmap reads.
type ChromeLocalState struct {
	LastUsed   string
	LastActive []string
	Profiles   map[string]ChromeProfileEntry // dir name → entry
}

// ChromeProfileEntry mirrors profile.info_cache[<dir>] fields gitmap surfaces.
type ChromeProfileEntry struct {
	DirName     string
	DisplayName string
	GAIAName    string
	UserName    string
	IsActive    bool
}

// ParseChromeLocalState decodes Local State JSON into a flat, test-friendly shape.
// Robust to: missing keys, extra keys, empty info_cache, last_used pointing at a
// directory not present in info_cache (returns the entry with empty display name).
func ParseChromeLocalState(raw []byte) (ChromeLocalState, error) {
	var doc struct {
		Profile struct {
			LastUsed       string   `json:"last_used"`
			LastActiveProf []string `json:"last_active_profiles"`
			InfoCache      map[string]struct {
				Name     string `json:"name"`
				GAIAName string `json:"gaia_name"`
				UserName string `json:"user_name"`
			} `json:"info_cache"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return ChromeLocalState{}, err
	}
	active := map[string]bool{}
	for _, d := range doc.Profile.LastActiveProf {
		active[d] = true
	}
	out := ChromeLocalState{
		LastUsed:   doc.Profile.LastUsed,
		LastActive: doc.Profile.LastActiveProf,
		Profiles:   map[string]ChromeProfileEntry{},
	}
	for dir, info := range doc.Profile.InfoCache {
		out.Profiles[dir] = ChromeProfileEntry{
			DirName:     dir,
			DisplayName: info.Name,
			GAIAName:    info.GAIAName,
			UserName:    info.UserName,
			IsActive:    active[dir],
		}
	}
	// Ensure last_used dir is queryable even if absent from info_cache.
	if out.LastUsed != "" {
		if _, ok := out.Profiles[out.LastUsed]; !ok {
			out.Profiles[out.LastUsed] = ChromeProfileEntry{DirName: out.LastUsed}
		}
	}
	return out, nil
}

// DisplayNameFor returns the human-readable name for a profile directory,
// falling back to the directory name when no entry exists.
func (s ChromeLocalState) DisplayNameFor(dir string) string {
	if e, ok := s.Profiles[dir]; ok && e.DisplayName != "" {
		return e.DisplayName
	}
	return dir
}
