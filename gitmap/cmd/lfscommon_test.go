package cmd

import "testing"

// TestLFSCommonPatternsMatchSpec locks in the curated default extension
// list so accidental edits (typos, removals, reordering) are caught by
// CI before they reach users.
func TestLFSCommonPatternsMatchSpec(t *testing.T) {
	want := []string{
		"*.pptx", "*.ppt", "*.eps", "*.psd", "*.ttf", "*.wott",
		"*.svg", "*.ai", "*.jpg", "*.bmp", "*.png", "*.zip",
		"*.gz", "*.tar", "*.rar", "*.7z", "*.mp4", "*.aep",
	}

	if len(lfsCommonPatterns) != len(want) {
		t.Fatalf("lfsCommonPatterns length: want %d, got %d", len(want), len(lfsCommonPatterns))
	}

	for i, p := range want {
		if lfsCommonPatterns[i] != p {
			t.Errorf("pattern[%d]: want %q, got %q", i, p, lfsCommonPatterns[i])
		}
	}
}

// TestLFSCommonPatternsAreUnique guarantees no duplicate entries — every
// pattern would otherwise appear twice in .gitattributes after `git lfs
// track`.
func TestLFSCommonPatternsAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range lfsCommonPatterns {
		if seen[p] {
			t.Errorf("duplicate pattern in lfsCommonPatterns: %q", p)
		}
		seen[p] = true
	}
}
