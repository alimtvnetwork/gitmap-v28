// Package cmd — agy_test.go provides unit tests for Antigravity commands and helpers.
package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseFolderUri(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{
			"file:///D%3A%2Fwp-work%2Friseup-asia%2F03-aukgo%2Fextendcore",
			filepath.FromSlash("D:/wp-work/riseup-asia/03-aukgo/extendcore"),
		},
		{
			"file:///C%3A%2FUsers%2FAlim%2Fproject",
			filepath.FromSlash("C:/Users/Alim/project"),
		},
	}
	for _, tc := range cases {
		got := parseFolderUri(tc.input)
		if !strings.EqualFold(got, tc.want) {
			t.Errorf("parseFolderUri(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestShortProjectId(t *testing.T) {
	if got := shortProjectId("0349c4d0-5a91-4f3e-800f-81fd53fc724f"); got != "0349c4d0" {
		t.Errorf("expected 0349c4d0, got %s", got)
	}
	if got := shortProjectId("short"); got != "short" {
		t.Errorf("expected short, got %s", got)
	}
}

func TestFormatRelativeTime(t *testing.T) {
	if got := formatRelativeTime(""); got != "—" {
		t.Errorf("expected —, got %s", got)
	}
	nowStr := time.Now().Format(time.RFC3339Nano)
	if got := formatRelativeTime(nowStr); got != "just now" {
		t.Errorf("expected 'just now', got %s", got)
	}
	pastStr := time.Now().Add(-2 * time.Hour).Format(time.RFC3339Nano)
	if got := formatRelativeTime(pastStr); got != "2h ago" {
		t.Errorf("expected '2h ago', got %s", got)
	}
}

func TestFilterAgyProjects(t *testing.T) {
	projects := []AgyProject{
		{ID: "1", Name: "Alpha"},
		{ID: "2", Name: "Beta"},
		{ID: "3", Name: "Gamma"},
	}
	agyLsFilter = "alp"
	filtered := filterAgyProjects(projects)
	agyLsFilter = ""

	if len(filtered) != 1 || filtered[0].Name != "Alpha" {
		t.Fatalf("expected [Alpha], got %v", filtered)
	}
}

func TestSortAgyProjects(t *testing.T) {
	projects := []AgyProject{
		{ID: "1", Name: "Zeta", UpdatedAt: "2026-09-01T00:00:00Z"},
		{ID: "2", Name: "Alpha", UpdatedAt: "2026-09-02T00:00:00Z"},
	}
	sortAgyProjects(projects, "name")
	if projects[0].Name != "Alpha" {
		t.Errorf("expected Alpha first, got %s", projects[0].Name)
	}
	sortAgyProjects(projects, "time")
	if projects[0].Name != "Alpha" {
		t.Errorf("expected Alpha first by time, got %s", projects[0].Name)
	}
}

func TestClearTargetsSelection(t *testing.T) {
	tempDir := t.TempDir()
	validDir := filepath.Join(tempDir, "valid")
	_ = os.MkdirAll(validDir, 0755)

	projects := []AgyProject{
		{ID: "outside-of-project", Name: "Outside"},
		{
			ID:   "valid-proj",
			Name: "Valid",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{GitFolder: &AgyGitFolder{FolderURI: "file:///" + filepath.ToSlash(validDir)}},
				},
			},
		},
		{
			ID:   "missing-proj",
			Name: "Missing",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{GitFolder: &AgyGitFolder{FolderURI: "file:///mock/nonexistent/path"}},
				},
			},
		},
	}

	targets := selectClearTargets(projects)
	if len(targets) != 1 || targets[0].ID != "missing-proj" {
		t.Fatalf("expected [missing-proj], got %v", targets)
	}
}

func TestFindAgyDuplicates(t *testing.T) {
	projects := []AgyProject{
		{
			ID:        "proj-1",
			Name:      "RepoA",
			UpdatedAt: "2026-09-01T12:00:00Z",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{GitFolder: &AgyGitFolder{FolderURI: "file:///mock/repos/repo-a"}},
				},
			},
		},
		{
			ID:        "proj-2-dup",
			Name:      "RepoA-old",
			UpdatedAt: "2026-08-01T12:00:00Z",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{GitFolder: &AgyGitFolder{FolderURI: "file:///mock/repos/repo-a"}},
				},
			},
		},
	}

	dups := findAgyDuplicates(projects, "")
	if len(dups) != 1 || dups[0].ID != "proj-2-dup" {
		t.Fatalf("expected [proj-2-dup] to be duplicate, got %v", dups)
	}

	exceptDups := findAgyDuplicates(projects, "proj-2-dup")
	if len(exceptDups) != 0 {
		t.Fatalf("expected 0 duplicates with --except, got %v", exceptDups)
	}
}

func TestIsAgyProjectExcepted(t *testing.T) {
	proj := AgyProject{
		ID:   "abc-123",
		Name: "my-service",
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{GitFolder: &AgyGitFolder{FolderURI: "file:///mock/repos/my-service"}},
			},
		},
	}

	if !isAgyProjectExcepted(proj, "abc-123") {
		t.Errorf("expected match on ID")
	}
	if !isAgyProjectExcepted(proj, "my-service") {
		t.Errorf("expected match on Name")
	}
	if !isAgyProjectExcepted(proj, "D:/repos/my-service") {
		t.Errorf("expected match on Path")
	}
	if isAgyProjectExcepted(proj, "other-service") {
		t.Errorf("expected no match on other-service")
	}
}


