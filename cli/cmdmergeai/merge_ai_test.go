package cmdmergeai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMergeAIArgs(t *testing.T) {
	// Comma separated
	dest, sources, err := parseMergeAIArgs([]string{"my-target", "url1,url2,url3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "my-target" {
		t.Fatalf("expected dest 'my-target', got '%s'", dest)
	}
	if len(sources) != 3 {
		t.Fatalf("expected 3 sources, got %d", len(sources))
	}

	// Space separated
	dest2, sources2, err2 := parseMergeAIArgs([]string{"dest-dir", "u1", "u2"})
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if dest2 != "dest-dir" || len(sources2) != 2 {
		t.Fatalf("expected dest-dir and 2 sources, got %s and %d", dest2, len(sources2))
	}
}

func TestStageAndSequenceFilesCollision(t *testing.T) {
	tempDest, _ := os.MkdirTemp("", "merge-ai-dest-*")
	defer os.RemoveAll(tempDest)

	tempSrc1, _ := os.MkdirTemp("", "merge-ai-src1-*")
	defer os.RemoveAll(tempSrc1)

	tempSrc2, _ := os.MkdirTemp("", "merge-ai-src2-*")
	defer os.RemoveAll(tempSrc2)

	// Create shared main.go in both
	_ = os.WriteFile(filepath.Join(tempSrc1, "main.go"), []byte("package main // from src1"), 0644)
	_ = os.WriteFile(filepath.Join(tempSrc2, "main.go"), []byte("package main // from src2"), 0644)

	// Create unique files
	_ = os.WriteFile(filepath.Join(tempSrc1, "unique1.txt"), []byte("unique 1"), 0644)
	_ = os.WriteFile(filepath.Join(tempSrc2, "unique2.txt"), []byte("unique 2"), 0644)

	sources := []SourceRepoConfig{
		{
			RepoURL:      "https://github.com/mock/src1.git",
			LocalStaging: tempSrc1,
			FolderTree:   []string{"main.go", "unique1.txt"},
		},
		{
			RepoURL:      "https://github.com/mock/src2.git",
			LocalStaging: tempSrc2,
			FolderTree:   []string{"main.go", "unique2.txt"},
		},
	}

	collisions, uniquePlaced, errStage := stageAndSequenceFiles(tempDest, sources)
	if errStage != nil {
		t.Fatalf("stageAndSequenceFiles failed: %v", errStage)
	}

	if len(collisions) != 1 {
		t.Fatalf("expected 1 collision for main.go, got %d", len(collisions))
	}
	if collisions[0].CanonicalPath != "main.go" {
		t.Fatalf("expected collision on main.go, got %s", collisions[0].CanonicalPath)
	}
	if len(collisions[0].Variants) != 2 {
		t.Fatalf("expected 2 variants (01 and 02), got %d", len(collisions[0].Variants))
	}

	// Verify collision files exist on disk
	if !fileExists(filepath.Join(tempDest, "01_main.go")) {
		t.Fatalf("expected 01_main.go to exist in dest")
	}
	if !fileExists(filepath.Join(tempDest, "02_main.go")) {
		t.Fatalf("expected 02_main.go to exist in dest")
	}

	// Verify unique files exist
	if !fileExists(filepath.Join(tempDest, "unique1.txt")) {
		t.Fatalf("expected unique1.txt to exist in dest")
	}
	if !fileExists(filepath.Join(tempDest, "unique2.txt")) {
		t.Fatalf("expected unique2.txt to exist in dest")
	}
	if uniquePlaced != 2 {
		t.Fatalf("expected 2 unique files placed, got %d", uniquePlaced)
	}
}

func TestWriteManifestAndInstruction(t *testing.T) {
	tempDest, _ := os.MkdirTemp("", "merge-ai-manifest-*")
	defer os.RemoveAll(tempDest)

	var manifest MergeAIManifest
	manifest.Attributes.Tool = "gitmap-merge-ai"
	manifest.Data.Destination.Slug = "test-monorepo"

	if err := WriteMergeAIManifest(tempDest, manifest); err != nil {
		t.Fatalf("WriteMergeAIManifest failed: %v", err)
	}
	if !fileExists(filepath.Join(tempDest, "merge-ai-manifest.json")) {
		t.Fatalf("merge-ai-manifest.json does not exist")
	}

	collisions := []FileCollisionRecord{
		{
			CanonicalPath: "main.go",
			Variants: []FileCollisionVariant{
				{Sequence: "01", File: "01_main.go", SourceRepo: "src1", Commit: "abc"},
				{Sequence: "02", File: "02_main.go", SourceRepo: "src2", Commit: "xyz"},
			},
		},
	}

	if err := WriteMergeAIInstruction(tempDest, collisions); err != nil {
		t.Fatalf("WriteMergeAIInstruction failed: %v", err)
	}
	if !fileExists(filepath.Join(tempDest, "instruction.md")) {
		t.Fatalf("instruction.md does not exist")
	}
}
