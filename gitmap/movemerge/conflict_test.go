package movemerge

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolver_PreferLeft(t *testing.T) {
	r := NewResolver(PreferLeft, nil, io.Discard)
	c, err := r.Resolve("x", FileMeta{}, FileMeta{})
	if err != nil || c != ChoiceLeft {
		t.Errorf("got c=%v err=%v", c, err)
	}
}

func TestResolver_PreferNewerPicksLatestMtime(t *testing.T) {
	dir := t.TempDir()
	older := writeWithMTime(t, filepath.Join(dir, "a"), time.Now().Add(-time.Hour))
	newer := writeWithMTime(t, filepath.Join(dir, "b"), time.Now())
	r := NewResolver(PreferNewer, nil, io.Discard)
	c, _ := r.Resolve("rel", FileMeta{Info: older}, FileMeta{Info: newer})
	if c != ChoiceRight {
		t.Errorf("PreferNewer: want Right (RIGHT is newer), got %v", c)
	}
}

func TestResolver_StickyAllLeft(t *testing.T) {
	in := bytes.NewBufferString("A\n")
	r := NewResolver(PreferNone, in, io.Discard)
	first, _ := r.Resolve("x", FileMeta{Info: dummyInfo(t)}, FileMeta{Info: dummyInfo(t)})
	if first != ChoiceLeft {
		t.Fatalf("first answer: got %v", first)
	}
	second, _ := r.Resolve("y", FileMeta{Info: dummyInfo(t)}, FileMeta{Info: dummyInfo(t)})
	if second != ChoiceLeft {
		t.Errorf("sticky All-Left didn't persist: got %v", second)
	}
}

func TestResolver_QuitKey(t *testing.T) {
	in := bytes.NewBufferString("Q\n")
	r := NewResolver(PreferNone, in, io.Discard)
	c, _ := r.Resolve("x", FileMeta{Info: dummyInfo(t)}, FileMeta{Info: dummyInfo(t)})
	if c != ChoiceQuit {
		t.Errorf("Q: got %v", c)
	}
}

// writeWithMTime creates a file and returns its FileInfo with set mtime.
func writeWithMTime(t *testing.T, path string, mt time.Time) os.FileInfo {
	t.Helper()
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mt, mt); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	return info
}

// dummyInfo returns any os.FileInfo for tests that don't depend on it.
func dummyInfo(t *testing.T) os.FileInfo {
	return writeWithMTime(t, filepath.Join(t.TempDir(), "z"), time.Now())
}
