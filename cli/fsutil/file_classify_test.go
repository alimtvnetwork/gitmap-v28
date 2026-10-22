package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifyFileTargets_Empty(t *testing.T) {
	cf := ClassifyFileTargets("", "", "")
	if len(cf.Read) != 0 || len(cf.Edit) != 0 || len(cf.Created) != 0 || cf.Total != 0 {
		t.Fatalf("expected empty classification, got %+v", cf)
	}
}

func TestClassifyFileTargets_Prefixed(t *testing.T) {
	raw := `["read:docs/spec.md", "edit:cli/main.go", "created:cli/new.go"]`
	cf := ClassifyFileTargets(raw, "", "")

	if len(cf.Read) != 1 || cf.Read[0] != "docs/spec.md" {
		t.Fatalf("unexpected read files: %+v", cf.Read)
	}
	if len(cf.Edit) != 1 || cf.Edit[0] != "cli/main.go" {
		t.Fatalf("unexpected edit files: %+v", cf.Edit)
	}
	if len(cf.Created) != 1 || cf.Created[0] != "cli/new.go" {
		t.Fatalf("unexpected created files: %+v", cf.Created)
	}
	if cf.Total != 3 {
		t.Fatalf("expected total 3, got %d", cf.Total)
	}
}

func TestClassifyFileTargets_DefaultAction(t *testing.T) {
	raw := `cli/store/db.go, cli/store/repo.go`
	cf := ClassifyFileTargets(raw, "READ", "")

	if len(cf.Read) != 2 {
		t.Fatalf("expected 2 read files from default action READ, got: %+v", cf.Read)
	}
	if cf.Total != 2 {
		t.Fatalf("expected total 2, got %d", cf.Total)
	}
}

func TestClassifyFileTargets_DiskFallback(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	_ = os.WriteFile(existingFile, []byte("hello"), 0o644)

	raw := `["existing.txt", "nonexistent.txt"]`
	cf := ClassifyFileTargets(raw, "", tmpDir)

	if len(cf.Edit) != 1 || cf.Edit[0] != "existing.txt" {
		t.Fatalf("expected existing.txt to be classified as edit, got: %+v", cf.Edit)
	}
	if len(cf.Created) != 1 || cf.Created[0] != "nonexistent.txt" {
		t.Fatalf("expected nonexistent.txt to be classified as created, got: %+v", cf.Created)
	}
}
