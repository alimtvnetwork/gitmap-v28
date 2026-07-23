package clonenext

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// LoadBatchFromCSV — header + column tests
// ---------------------------------------------------------------------------

func TestLoadBatchFromCSV_HeaderlessSinglePath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "list.csv")
	writeFile(t, path, "./my-repo\n")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1", len(got))
	}
	if !strings.HasSuffix(got[0], "my-repo") {
		t.Errorf("got %q, want suffix %q", got[0], "my-repo")
	}
}

func TestLoadBatchFromCSV_WithHeaderAndExtraColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "list.csv")
	writeFile(t, path, "repo,version,note\n./alpha,v++,first\n./beta,v3,second\n")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d rows, want 2", len(got))
	}
}

func TestLoadBatchFromCSV_NamedPathColumnNotFirst(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "list.csv")
	writeFile(t, path, "name,path,note\nalpha,./alpha-dir,some note\nbeta,./beta-dir,\n")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d rows, want 2", len(got))
	}
	if !strings.HasSuffix(got[0], "alpha-dir") {
		t.Errorf("row 0 = %q, want suffix alpha-dir", got[0])
	}
}

func TestLoadBatchFromCSV_EmptyReturnsSentinel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "list.csv")
	writeFile(t, path, "repo\n\n  \n")

	_, err := LoadBatchFromCSV(path)
	if !errors.Is(err, ErrBatchEmpty) {
		t.Fatalf("err = %v, want ErrBatchEmpty", err)
	}
}

// ---------------------------------------------------------------------------
// BOM + line endings
// ---------------------------------------------------------------------------

func TestLoadBatchFromCSV_BOMStripped(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bom.csv")
	// UTF-8 BOM + "repo" header
	bom := "\xEF\xBB\xBFrepo,note\n./r1,bom test\n"
	writeFile(t, path, bom)

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV (BOM): %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1 (BOM header was not recognized?)", len(got))
	}
}

func TestLoadBatchFromCSV_WindowsCRLF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "crlf.csv")
	writeFile(t, path, "repo,note\r\n./a,one\r\n./b,two\r\n")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV (CRLF): %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d rows, want 2", len(got))
	}
}

func TestLoadBatchFromCSV_BareCR(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cr.csv")
	// Classic macOS CR-only line endings
	writeFile(t, path, "repo\r./a\r./b\r")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV (bare CR): %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d rows, want 2 (bare CR not normalized?)", len(got))
	}
}

// ---------------------------------------------------------------------------
// Missing optional columns (ragged rows)
// ---------------------------------------------------------------------------

func TestLoadBatchFromCSV_MissingOptionalColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ragged.csv")
	// Header has 3 columns; data rows may have fewer.
	writeFile(t, path, "repo,version,note\n./a\n./b,v2\n./c,v3,some note\n")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV (ragged): %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d rows, want 3", len(got))
	}
}

func TestLoadBatchFromCSV_PathColumnMissedInShortRow(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "short.csv")
	// Path is column 2 ("path") but some rows only have 1 cell.
	writeFile(t, path, "name,tag,path\nalpha,v1,./a\nbeta\n")

	got, err := LoadBatchFromCSV(path)
	if err != nil {
		t.Fatalf("LoadBatchFromCSV: %v", err)
	}
	// beta row is too short to have a path — should be silently skipped.
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1 (short row should have been skipped)", len(got))
	}
}

// ---------------------------------------------------------------------------
// WalkBatchFromDir
// ---------------------------------------------------------------------------

func TestWalkBatchFromDir_OnlyGitDirsIncluded(t *testing.T) {
	root := t.TempDir()
	mkRepo(t, filepath.Join(root, "zeta"))
	mkRepo(t, filepath.Join(root, "alpha"))
	if err := os.MkdirAll(filepath.Join(root, "not-a-repo"), 0o755); err != nil {
		t.Fatalf("mkdir non-repo: %v", err)
	}

	got, err := WalkBatchFromDir(root)
	if err != nil {
		t.Fatalf("WalkBatchFromDir: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d repos, want 2", len(got))
	}
	if !strings.HasSuffix(got[0], "alpha") || !strings.HasSuffix(got[1], "zeta") {
		t.Errorf("got %v, want sorted [alpha, zeta]", got)
	}
}

func TestWalkBatchFromDir_NoReposReturnsSentinel(t *testing.T) {
	root := t.TempDir()

	_, err := WalkBatchFromDir(root)
	if !errors.Is(err, ErrBatchEmpty) {
		t.Fatalf("err = %v, want ErrBatchEmpty", err)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mkRepo(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(path, ".git"), 0o755); err != nil {
		t.Fatalf("mkRepo %s: %v", path, err)
	}
}
