package checkpoint

import (
	"path/filepath"
	"testing"
)

func TestCheckpointRoundtrip(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")
	f, err := Open(dir, "input-1", 100)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if f.IsDone("abc") {
		t.Fatal("fresh checkpoint should not report sha as done")
	}
	if err := f.MarkDone("abc"); err != nil {
		t.Fatalf("MarkDone: %v", err)
	}
	g, err := Open(dir, "input-1", 101)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if !g.IsDone("abc") {
		t.Fatal("reopened checkpoint must report previously-done sha")
	}
	if g.IsDone("def") {
		t.Fatal("unrelated sha must not be flagged done")
	}
}

func TestCheckpointTolerantOfCorruptFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")
	f, err := Open(dir, "x", 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.MarkDone("sha1"); err != nil {
		t.Fatal(err)
	}
	// Corrupt the file then reopen — must not error.
	if err := writeFile(f.path, []byte("{not json")); err != nil {
		t.Fatal(err)
	}
	g, err := Open(dir, "x", 2)
	if err != nil {
		t.Fatalf("reopen on corrupt: %v", err)
	}
	if g.IsDone("sha1") {
		t.Fatal("corrupt reset must clear done set")
	}
}

// writeFile is a tiny indirection so the test stays free of os import
// repetition; checkpoint.go owns the os.WriteFile in production paths.
func writeFile(p string, b []byte) error {
	return osWriteFile(p, b)
}
