package cmdmergeai

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

func TestParseMergeAIArgs(t *testing.T) {
	// Comma separated
	dest, sources, _, dryRun, err := parseMergeAIArgs([]string{"my-target", "https://github.com/a/u1,https://github.com/b/u2,https://github.com/c/u3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "my-target" {
		t.Fatalf("expected dest 'my-target', got '%s'", dest)
	}
	if len(sources) != 3 {
		t.Fatalf("expected 3 sources, got %d", len(sources))
	}
	if dryRun {
		t.Errorf("expected dryRun false")
	}

	// Space separated
	dest2, sources2, _, _, err2 := parseMergeAIArgs([]string{"dest-dir", "https://github.com/a/u1", "https://github.com/b/u2"})
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if dest2 != "dest-dir" || len(sources2) != 2 {
		t.Fatalf("expected dest-dir and 2 sources, got %s and %d", dest2, len(sources2))
	}

	// Named flags with --dry-run
	dest3, sources3, _, dryRun3, err3 := parseMergeAIArgs([]string{"-d", "icon-studio-v1", "-s", "https://github.com/a/u1", "--dry-run"})
	if err3 != nil {
		t.Fatalf("unexpected error with named flags: %v", err3)
	}
	if dest3 != "icon-studio-v1" || len(sources3) != 1 || !dryRun3 {
		t.Fatalf("expected icon-studio-v1, 1 source, dryRun=true, got %s, %d, %v", dest3, len(sources3), dryRun3)
	}

	// Help flag
	destHelp, _, _, _, errHelp := parseMergeAIArgs([]string{"--help"})
	if errHelp != nil {
		t.Fatalf("unexpected error on help: %v", errHelp)
	}
	if destHelp != "" {
		t.Fatalf("expected empty dest on help")
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
	if !fsutil.FileExists(filepath.Join(tempDest, "01_main.go")) {
		t.Fatalf("expected 01_main.go to exist in dest")
	}
	if !fsutil.FileExists(filepath.Join(tempDest, "02_main.go")) {
		t.Fatalf("expected 02_main.go to exist in dest")
	}

	// Verify unique files exist
	if !fsutil.FileExists(filepath.Join(tempDest, "unique1.txt")) {
		t.Fatalf("expected unique1.txt to exist in dest")
	}
	if !fsutil.FileExists(filepath.Join(tempDest, "unique2.txt")) {
		t.Fatalf("expected unique2.txt to exist in dest")
	}
	if uniquePlaced != 2 {
		t.Fatalf("expected 2 unique files placed, got %d", uniquePlaced)
	}
}

func TestStageAndSequenceFilesCollision_ThreeRepos(t *testing.T) {
	tempDest, _ := os.MkdirTemp("", "merge-ai-dest3-*")
	defer os.RemoveAll(tempDest)

	tempSrc1, _ := os.MkdirTemp("", "merge-ai-src1-*")
	defer os.RemoveAll(tempSrc1)

	tempSrc2, _ := os.MkdirTemp("", "merge-ai-src2-*")
	defer os.RemoveAll(tempSrc2)

	tempSrc3, _ := os.MkdirTemp("", "merge-ai-src3-*")
	defer os.RemoveAll(tempSrc3)

	_ = os.WriteFile(filepath.Join(tempSrc1, "readme.md"), []byte("# Repo 1"), 0644)
	_ = os.WriteFile(filepath.Join(tempSrc2, "readme.md"), []byte("# Repo 2"), 0644)
	_ = os.WriteFile(filepath.Join(tempSrc3, "readme.md"), []byte("# Repo 3"), 0644)

	sources := []SourceRepoConfig{
		{
			RepoURL:      "https://github.com/mock/src1.git",
			LocalStaging: tempSrc1,
			FolderTree:   []string{"readme.md"},
		},
		{
			RepoURL:      "https://github.com/mock/src2.git",
			LocalStaging: tempSrc2,
			FolderTree:   []string{"readme.md"},
		},
		{
			RepoURL:      "https://github.com/mock/src3.git",
			LocalStaging: tempSrc3,
			FolderTree:   []string{"readme.md"},
		},
	}

	collisions, _, err := stageAndSequenceFiles(tempDest, sources)
	if err != nil {
		t.Fatalf("stageAndSequenceFiles failed: %v", err)
	}

	if len(collisions) != 1 {
		t.Fatalf("expected 1 collision, got %d", len(collisions))
	}
	if len(collisions[0].Variants) != 3 {
		t.Fatalf("expected 3 variants for readme.md, got %d", len(collisions[0].Variants))
	}

	if !fsutil.FileExists(filepath.Join(tempDest, "01_readme.md")) {
		t.Errorf("expected 01_readme.md on disk")
	}
	if !fsutil.FileExists(filepath.Join(tempDest, "02_readme.md")) {
		t.Errorf("expected 02_readme.md on disk")
	}
	if !fsutil.FileExists(filepath.Join(tempDest, "03_readme.md")) {
		t.Errorf("expected 03_readme.md on disk")
	}
}

func TestWriteManifestAndInstruction(t *testing.T) {
	tempDest, _ := os.MkdirTemp("", "merge-ai-manifest-*")
	defer os.RemoveAll(tempDest)

	var manifest MergeAIManifest
	manifest.Attributes.Tool = "gitmap-merge-ai"
	manifest.Data.Destination.Slug = "test-monorepo"

	entry := SourceRepoSequenceEntry{
		SequenceOrder: 1,
		Title:         "Sample Repo",
		RepoURL:       "https://github.com/mock/repo.git",
		Branch:        "main",
		CommitMetadata: CommitMetadata{
			SHA:     "1234567890abcdef",
			Author:  "Test Author <test@example.com>",
			Date:    "2026-10-11T12:00:00Z",
			Message: "Initial commit",
			Branch:  "main",
		},
	}
	manifest.Data.RepoSequence = append(manifest.Data.RepoSequence, entry)

	if err := WriteMergeAIManifest(tempDest, manifest); err != nil {
		t.Fatalf("WriteMergeAIManifest failed: %v", err)
	}
	manifestPath := filepath.Join(tempDest, "merge-ai-manifest.json")
	if !fsutil.FileExists(manifestPath) {
		t.Fatalf("merge-ai-manifest.json does not exist")
	}

	manifestBytes, errRead := os.ReadFile(manifestPath)
	if errRead != nil {
		t.Fatalf("failed to read manifest: %v", errRead)
	}
	manifestContent := string(manifestBytes)
	if !strings.Contains(manifestContent, "1234567890abcdef") || !strings.Contains(manifestContent, "Test Author") {
		t.Errorf("manifest does not contain expected commit metadata: %s", manifestContent)
	}

	collisions := []FileCollisionRecord{
		{
			CanonicalPath: "main.go",
			Variants: []FileCollisionVariant{
				{Sequence: "01", File: "01_main.go", SourceRepo: "src1", Commit: "abc"},
				{Sequence: "02", File: "02_main.go", SourceRepo: "src2", Commit: "xyz"},
			},
		},
		{
			CanonicalPath: "assets/logo.png",
			Variants: []FileCollisionVariant{
				{Sequence: "01", File: "assets/01_logo.png", SourceRepo: "src1", Commit: "abc"},
				{Sequence: "02", File: "assets/02_logo.png", SourceRepo: "src2", Commit: "xyz"},
			},
		},
	}

	skipped := []SkippedSourceItem{
		{
			RawInput: "https://github.com/mock/file.gif",
			Title:    "Logo GIF",
			Repo:     "mock/src1",
			Reason:   "file URL from repository 'mock/src1' which is already included",
			IsFile:   true,
			FilePath: "file.gif",
		},
	}

	if err := WriteMergeAIInstruction(tempDest, collisions, skipped...); err != nil {
		t.Fatalf("WriteMergeAIInstruction failed: %v", err)
	}
	instPath := filepath.Join(tempDest, "instruction.md")
	if !fsutil.FileExists(instPath) {
		t.Fatalf("instruction.md does not exist")
	}

	instBytes, _ := os.ReadFile(instPath)
	instContent := string(instBytes)
	if !strings.Contains(instContent, "Folder Structure") {
		t.Errorf("instruction.md missing Folder Structure advice")
	}
	if !strings.Contains(instContent, "Media Asset") {
		t.Errorf("instruction.md missing Media Asset classification")
	}
	if !strings.Contains(instContent, "side-by-side") {
		t.Errorf("instruction.md missing side-by-side preservation advice")
	}
}
