package release

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// setupTempDir overrides DefaultReleaseDir to a temp directory and returns a cleanup func.
func setupTempDir(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	orig := constants.DefaultReleaseDir
	constants.DefaultReleaseDir = filepath.Join(dir, ".gitmap", "release")
	return func() { constants.DefaultReleaseDir = orig }
}

func TestWriteAndReadReleaseMeta(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	meta := ReleaseMeta{
		Version:      "1.2.3",
		Branch:       "release/v1.2.3",
		SourceBranch: "main",
		Commit:       "abc123",
		Tag:          "v1.2.3",
		Assets:       []string{"./dist"},
		IsDraft:      false,
		IsPreRelease: false,
		CreatedAt:    "2026-03-05T12:00:00Z",
		IsLatest:     true,
	}

	err := WriteReleaseMeta(meta)
	if err != nil {
		t.Fatalf("WriteReleaseMeta error: %v", err)
	}

	v, _ := Parse("v1.2.3")
	path := filepath.Join(constants.DefaultReleaseDir, v.String()+constants.ExtJSON)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("metadata file not created at %s", path)
	}
}

func TestReleaseExistsTrueAfterWrite(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, _ := Parse("v1.0.0")
	if ReleaseExists(v) {
		t.Error("should not exist before write")
	}

	meta := ReleaseMeta{Tag: "v1.0.0", Version: "1.0.0"}
	err := WriteReleaseMeta(meta)
	if err != nil {
		t.Fatalf("WriteReleaseMeta error: %v", err)
	}

	if ReleaseExists(v) == false {
		t.Error("should exist after write")
	}
}

func TestReleaseExistsFalseForMissing(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, _ := Parse("v9.9.9")
	if ReleaseExists(v) {
		t.Error("non-existent release should return false")
	}
}

func TestWriteAndReadLatest(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	os.MkdirAll(constants.DefaultReleaseDir, constants.DirPermission)

	v, _ := Parse("v1.0.0")
	err := WriteLatest(v)
	if err != nil {
		t.Fatalf("WriteLatest error: %v", err)
	}

	latest, err := ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest error: %v", err)
	}

	if latest.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", latest.Version)
	}
	if latest.Tag != "v1.0.0" {
		t.Errorf("expected tag v1.0.0, got %s", latest.Tag)
	}
	if latest.Branch != "release/v1.0.0" {
		t.Errorf("expected branch release/v1.0.0, got %s", latest.Branch)
	}
}

func TestWriteLatestOnlyUpgrades(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	os.MkdirAll(constants.DefaultReleaseDir, constants.DirPermission)

	v2, _ := Parse("v2.0.0")
	WriteLatest(v2)

	v1, _ := Parse("v1.0.0")
	WriteLatest(v1)

	latest, err := ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest error: %v", err)
	}

	if latest.Version != "2.0.0" {
		t.Errorf("latest should remain 2.0.0, got %s", latest.Version)
	}
}

func TestWriteLatestUpgradesHigher(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	os.MkdirAll(constants.DefaultReleaseDir, constants.DirPermission)

	v1, _ := Parse("v1.0.0")
	WriteLatest(v1)

	v2, _ := Parse("v2.0.0")
	WriteLatest(v2)

	latest, _ := ReadLatest()
	if latest.Version != "2.0.0" {
		t.Errorf("expected 2.0.0, got %s", latest.Version)
	}
}

func TestReadLatestMissingFile(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	_, err := ReadLatest()
	if err == nil {
		t.Error("expected error when latest.json does not exist")
	}
}
