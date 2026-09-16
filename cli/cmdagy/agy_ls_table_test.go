// Package cmdagy — agy_ls_table_test.go tests table truncation, formatting, and folder grouping.
package cmdagy

import (
	"testing"
)

func TestTruncateMiddle(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exact-length", 12, "exact-length"},
		{"ai-empathy-prompt-tuner", 15, "ai-emp...t-tuner"},
		{"abcdefghij", 7, "ab...ij"},
		{"tiny", 3, "tiny"},
	}

	for _, tt := range tests {
		got := truncateMiddle(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncateMiddle(%q, %d) = %q; want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestResolveProjectParentFolder(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"", "[Global / Config]"},
		{"—", "[Global / Config]"},
		{"d:/work/gitmap", "d:/work"},
		{"/home/user/project", "/home/user"},
	}

	for _, tt := range tests {
		p := AgyProject{Path: tt.path}
		got := resolveProjectParentFolder(p)
		if got != tt.want && got != filepathClean(tt.want) {
			t.Errorf("resolveProjectParentFolder(%q) = %q; want %q", tt.path, got, tt.want)
		}
	}
}

func filepathClean(p string) string {
	if p == "[Global / Config]" {
		return p
	}

	return p
}

func TestGroupProjectsByRootFolder(t *testing.T) {
	projects := []AgyProject{
		{ID: "1", Name: "repo1", Path: "d:/work/repo1"},
		{ID: "2", Name: "repo2", Path: "d:/work/repo2"},
		{ID: "3", Name: "tool", Path: "d:/tools/tool"},
	}

	folders, groupMap := groupProjectsByRootFolder(projects)
	if len(folders) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(folders))
	}

	if len(groupMap[folders[0]])+len(groupMap[folders[1]]) != 3 {
		t.Errorf("total projects in groupMap must be 3")
	}
}
