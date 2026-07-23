package vscodepm

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestDetectTagsRecognizesMarkers seeds a temp dir with several known
// marker files and confirms DetectTags emits them in canonical order.
func TestDetectTagsRecognizesMarkers(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "package.json"), "{}")
	mustWrite(t, filepath.Join(root, "go.mod"), "module x\n")
	mustWrite(t, filepath.Join(root, "Dockerfile"), "FROM scratch\n")
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	got := DetectTags(root)
	want := []string{
		constants.AutoTagGit,
		constants.AutoTagNode,
		constants.AutoTagGo,
		constants.AutoTagDocker,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectTags() = %v, want %v", got, want)
	}
}

// TestDetectTagsEmptyDir returns nil for a directory with no markers.
func TestDetectTagsEmptyDir(t *testing.T) {
	if got := DetectTags(t.TempDir()); got != nil {
		t.Fatalf("DetectTags(empty) = %v, want nil", got)
	}
}

// TestDetectTagsMissingPath returns nil for a non-existent path.
func TestDetectTagsMissingPath(t *testing.T) {
	if got := DetectTags(filepath.Join(t.TempDir(), "does-not-exist")); got != nil {
		t.Fatalf("DetectTags(missing) = %v, want nil", got)
	}
}

// TestUnionTagsPreservesUserOrderAndDedupes confirms the additive UNION
// keeps the user's tag order first and never duplicates.
func TestUnionTagsPreservesUserOrderAndDedupes(t *testing.T) {
	existing := []string{"work", "favorite"}
	incoming := []string{"git", "work", "node"}

	got := unionTags(existing, incoming)
	want := []string{"work", "favorite", "git", "node"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unionTags() = %v, want %v", got, want)
	}
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
