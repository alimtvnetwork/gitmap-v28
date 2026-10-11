package cmdmergeai

import (
	"testing"
)

const sampleUserText = `
alimtvnetwork/macro-mate-scripts-fixer-logo-icon
https://github.com/alimtvnetwork/macro-mate-scripts-fixer-logo-icon

alimtvnetwork/logo-1
https://github.com/alimtvnetwork/logo-1

logo-1/Projects/01-MacroMate/icons-image/Logo-Animated.gif at main · alimtvnetwork/logo-1
https://github.com/alimtvnetwork/logo-1/blob/main/Projects/01-MacroMate/icons-image/Logo-Animated.gif

alimtvnetwork/iconcraft-studio
https://github.com/alimtvnetwork/iconcraft-studio

alimtvnetwork/brand-builder-studio
https://github.com/alimtvnetwork/brand-builder-studio

alimtvnetwork/brand-weaver
https://github.com/alimtvnetwork/brand-weaver

alimtvnetwork/brand-forge-studio
https://github.com/alimtvnetwork/brand-forge-studio

alimtvnetwork/macro-mate-logo-icon-27f183d0
https://github.com/alimtvnetwork/macro-mate-logo-icon-27f183d0
`

func TestParseSourceManifestContent_UserSample(t *testing.T) {
	destHint, resolved, skipped, err := ParseSourceManifestContent(sampleUserText)
	if err != nil {
		t.Fatalf("unexpected error parsing manifest content: %v", err)
	}

	if destHint != "" {
		t.Errorf("expected empty destHint, got %s", destHint)
	}

	if len(resolved) != 7 {
		t.Fatalf("expected 7 resolved repos, got %d", len(resolved))
	}

	if len(skipped) != 1 {
		t.Fatalf("expected 1 skipped item, got %d", len(skipped))
	}

	// Verify the skipped item is the logo-1 animated gif file
	s := skipped[0]
	if !s.IsFile {
		t.Errorf("expected skipped item to be a file URL")
	}
	if s.Repo != "alimtvnetwork/logo-1" {
		t.Errorf("expected skipped repo to be alimtvnetwork/logo-1, got %s", s.Repo)
	}

	// Verify all 7 expected repos are present in order
	expectedRepos := []string{
		"alimtvnetwork/macro-mate-scripts-fixer-logo-icon",
		"alimtvnetwork/logo-1",
		"alimtvnetwork/iconcraft-studio",
		"alimtvnetwork/brand-builder-studio",
		"alimtvnetwork/brand-weaver",
		"alimtvnetwork/brand-forge-studio",
		"alimtvnetwork/macro-mate-logo-icon-27f183d0",
	}

	for i, exp := range expectedRepos {
		if resolved[i].CanonicalRepo != exp {
			t.Errorf("index %d: expected %s, got %s", i, exp, resolved[i].CanonicalRepo)
		}
	}
}

func TestParseSourceManifestContent_EmptyHeadersAndDeduplication(t *testing.T) {
	raw := `
# Destination: icon-studio-v1
https://github.com/alimtvnetwork/repo-a
https://github.com/alimtvnetwork/repo-b
https://github.com/alimtvnetwork/repo-a
https://github.com/alimtvnetwork/repo-b/blob/main/readme.md
`
	destHint, resolved, skipped, err := ParseSourceManifestContent(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if destHint != "icon-studio-v1" {
		t.Errorf("expected destHint 'icon-studio-v1', got '%s'", destHint)
	}

	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved repos, got %d", len(resolved))
	}

	if len(skipped) != 2 {
		t.Fatalf("expected 2 skipped items (1 duplicate repo + 1 file URL), got %d", len(skipped))
	}

	if resolved[0].Title != "alimtvnetwork/repo-a" {
		t.Errorf("expected title inferred from URL, got %s", resolved[0].Title)
	}
}

func TestParseURLDetails(t *testing.T) {
	cases := []struct {
		url      string
		repo     string
		isFile   bool
		filePath string
	}{
		{
			url:      "https://github.com/alimtvnetwork/logo-1/blob/main/Projects/01-MacroMate/icons-image/Logo-Animated.gif",
			repo:     "alimtvnetwork/logo-1",
			isFile:   true,
			filePath: "Projects/01-MacroMate/icons-image/Logo-Animated.gif",
		},
		{
			url:      "https://github.com/alimtvnetwork/logo-1.git",
			repo:     "alimtvnetwork/logo-1",
			isFile:   false,
			filePath: "",
		},
		{
			url:      "alimtvnetwork/brand-weaver",
			repo:     "alimtvnetwork/brand-weaver",
			isFile:   false,
			filePath: "",
		},
	}

	for _, c := range cases {
		repo, _, isFile, filePath, err := ParseURLDetails(c.url)
		if err != nil {
			t.Errorf("unexpected error for %s: %v", c.url, err)
			continue
		}
		if repo != c.repo {
			t.Errorf("expected repo %s, got %s", c.repo, repo)
		}
		if isFile != c.isFile {
			t.Errorf("expected isFile %v, got %v", c.isFile, isFile)
		}
		if filePath != c.filePath {
			t.Errorf("expected filePath %s, got %s", c.filePath, filePath)
		}
	}
}

func TestParseSourceManifestContent_FileURLBeforeRepo(t *testing.T) {
	// Tests that if the file URL is listed BEFORE the repository URL,
	// the file URL is still skipped because the repo is included!
	raw := `
logo-1 animated gif
https://github.com/alimtvnetwork/logo-1/blob/main/Projects/01-MacroMate/icons-image/Logo-Animated.gif

alimtvnetwork/logo-1
https://github.com/alimtvnetwork/logo-1
`
	_, resolved, skipped, err := ParseSourceManifestContent(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resolved) != 1 {
		t.Fatalf("expected 1 resolved repo, got %d", len(resolved))
	}
	if resolved[0].CanonicalRepo != "alimtvnetwork/logo-1" {
		t.Errorf("expected resolved repo alimtvnetwork/logo-1, got %s", resolved[0].CanonicalRepo)
	}
	if resolved[0].IsFromFile {
		t.Errorf("expected IsFromFile to be false for the full repository")
	}

	if len(skipped) != 1 {
		t.Fatalf("expected 1 skipped file URL, got %d", len(skipped))
	}
	if !skipped[0].IsFile {
		t.Errorf("expected skipped item to be IsFile")
	}
}

func TestParseSourceManifestContent_GenericLabelsAndMarkdownHeadings(t *testing.T) {
	raw := `
# Destination: icon-studio-2026

## My MacroMate Repo
https://github.com/alimtvnetwork/macro-mate-scripts-fixer-logo-icon

Title:
https://github.com/alimtvnetwork/iconcraft-studio

Repo:
https://github.com/alimtvnetwork/brand-weaver
`
	destHint, resolved, _, err := ParseSourceManifestContent(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if destHint != "icon-studio-2026" {
		t.Errorf("expected destHint icon-studio-2026, got %s", destHint)
	}

	if len(resolved) != 3 {
		t.Fatalf("expected 3 resolved repos, got %d", len(resolved))
	}

	if resolved[0].Title != "My MacroMate Repo" {
		t.Errorf("expected markdown heading title 'My MacroMate Repo', got '%s'", resolved[0].Title)
	}

	// For Title: and Repo: generic labels, title is inferred from URL
	if resolved[1].Title != "alimtvnetwork/iconcraft-studio" {
		t.Errorf("expected inferred title for generic label, got '%s'", resolved[1].Title)
	}
	if resolved[2].Title != "alimtvnetwork/brand-weaver" {
		t.Errorf("expected inferred title for generic label, got '%s'", resolved[2].Title)
	}
}

func TestParseMergeAIArgs_MultiWordDestination(t *testing.T) {
	dest, _, _, err := parseMergeAIArgs([]string{"Icon", "Studio", "v1", "merge-ai-icon-studio.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "Icon Studio v1" {
		t.Errorf("expected dest 'Icon Studio v1', got '%s'", dest)
	}

	dest2, _, _, err2 := parseMergeAIArgs([]string{"Icon", "Studio", "2026", "merge-ai-icon-studio.txt"})
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if dest2 != "Icon Studio 2026" {
		t.Errorf("expected dest 'Icon Studio 2026', got '%s'", dest2)
	}
}

