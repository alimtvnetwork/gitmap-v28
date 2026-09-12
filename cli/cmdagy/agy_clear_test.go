package cmdagy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAgyClearProtectsPinnedProjects(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("HOME", tmp)

	pinnedPath := filepath.Join(tmp, "pinned-repo")
	_ = os.MkdirAll(pinnedPath, 0755)

	pinned, err := addPinnedProjectTarget(pinnedPath)
	if err != nil {
		t.Fatalf("failed to pin project: %v", err)
	}

	// Now pretend the directory is missing (deleted)
	_ = os.RemoveAll(pinnedPath)

	projects := []AgyProject{
		{
			ID:   pinned.ID,
			Name: "pinned-repo",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{
						GitFolder: &AgyGitFolder{
							FolderURI: "file:///" + filepath.ToSlash(pinnedPath),
						},
					},
				},
			},
		},
		{
			ID:   "unpinned-missing",
			Name: "unpinned-missing",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{
						GitFolder: &AgyGitFolder{
							FolderURI: "file:///" + filepath.ToSlash(filepath.Join(tmp, "unpinned-missing")),
						},
					},
				},
			},
		},
	}

	targets := selectClearTargets(projects)
	for _, target := range targets {
		if target.ID == pinned.ID {
			t.Errorf("agy clear targeted pinned project %s for removal!", pinned.ID)
		}
	}

	if len(targets) != 1 || targets[0].ID != "unpinned-missing" {
		t.Errorf("expected only unpinned-missing in targets, got: %+v", targets)
	}
}
