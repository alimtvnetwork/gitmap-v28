// Package fsutil — exists_test.go pins the contract divergence between
// the three existence predicates. The cmd package previously hit a
// redeclaration build break when two files defined fileExists with
// different semantics; the centralization in fsutil only pays off if
// each variant's contract is tested and stable. If a future contributor
// tries to "simplify" by collapsing two variants, these tests fail and
// surface the intent before the merge lands.
package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExistsContract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "marker")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if FileExists("") {
		t.Fatal("FileExists(\"\") must be false (empty short-circuit)")
	}
	if FileExists(dir) {
		t.Fatal("FileExists(dir) must be false (strict variant rejects directories)")
	}
	if !FileExists(file) {
		t.Fatal("FileExists(file) must be true")
	}
	if FileExists(filepath.Join(dir, "missing")) {
		t.Fatal("FileExists(missing) must be false")
	}
}

func TestFileOrDirExistsContract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "marker")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if FileOrDirExists("") {
		t.Fatal("FileOrDirExists(\"\") must be false (empty short-circuit; debug-dump contract)")
	}
	if !FileOrDirExists(dir) {
		t.Fatal("FileOrDirExists(dir) must be true (loose variant accepts directories)")
	}
	if !FileOrDirExists(file) {
		t.Fatal("FileOrDirExists(file) must be true")
	}
	if FileOrDirExists(filepath.Join(dir, "missing")) {
		t.Fatal("FileOrDirExists(missing) must be false")
	}
}

func TestDirExistsContract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "marker")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if DirExists("") {
		t.Fatal("DirExists(\"\") must be false (empty short-circuit)")
	}
	if DirExists(file) {
		t.Fatal("DirExists(file) must be false (strict variant rejects files)")
	}
	if !DirExists(dir) {
		t.Fatal("DirExists(dir) must be true")
	}
	if DirExists(filepath.Join(dir, "missing")) {
		t.Fatal("DirExists(missing) must be false")
	}
}
