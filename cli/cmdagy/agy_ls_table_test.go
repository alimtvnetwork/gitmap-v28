// Package cmdagy — agy_ls_table_test.go tests table truncation, formatting, and folder grouping.
package cmdagy

import (
	"path/filepath"
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
		{"ai-empathy-prompt-tuner", 15, "ai-emp...-tuner"},
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

func makeTestProject(id, name, path string) AgyProject {
	if path == "" || path == "—" {
		return AgyProject{ID: id, Name: name}
	}

	return AgyProject{
		ID:   id,
		Name: name,
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{
					GitFolder: &AgyGitFolder{
						FolderURI: path,
					},
				},
			},
		},
	}
}

func TestResolveProjectParentFolder(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"", "[Global / Config]"},
		{"—", "[Global / Config]"},
		{"d:/repos/project1", "d:/repos"},
		{"/home/user/project", "/home/user"},
	}

	for _, tt := range tests {
		p := makeTestProject("1", "test", tt.path)
		got := resolveProjectParentFolder(p)
		wantClean := filepath.Clean(tt.want)
		if got != tt.want && got != wantClean {
			t.Errorf("resolveProjectParentFolder(%q) = %q; want %q", tt.path, got, tt.want)
		}
	}
}

func TestGroupProjectsByRootFolder(t *testing.T) {
	projects := []AgyProject{
		makeTestProject("1", "repo1", "d:/work/repo1"),
		makeTestProject("2", "repo2", "d:/work/repo2"),
		makeTestProject("3", "tool", "d:/tools/tool"),
	}

	folders, groupMap := groupProjectsByRootFolder(projects)
	if len(folders) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(folders))
	}

	if len(groupMap[folders[0]])+len(groupMap[folders[1]]) != 3 {
		t.Errorf("total projects in groupMap must be 3")
	}
}
