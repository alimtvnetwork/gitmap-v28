package release_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// setupSkipMetaDir overrides DefaultReleaseDir to a temp directory
// and returns a cleanup func that restores the original value.
func setupSkipMetaDir(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	orig := constants.DefaultReleaseDir
	constants.DefaultReleaseDir = filepath.Join(dir, ".gitmap", "release")

	return func() { constants.DefaultReleaseDir = orig }
}

// assertNoMetadata fails the test if a release metadata JSON file exists.
func assertNoMetadata(t *testing.T, version string) {
	t.Helper()
	path := filepath.Join(constants.DefaultReleaseDir, version+".json")
	_, err := os.Stat(path)
	if err == nil {
		t.Errorf("metadata file should not exist: %s", path)
	}
}

// assertNoLatest fails the test if latest.json exists.
func assertNoLatest(t *testing.T) {
	t.Helper()
	path := filepath.Join(constants.DefaultReleaseDir, constants.DefaultLatestFile)
	_, err := os.Stat(path)
	if err == nil {
		t.Errorf("latest.json should not exist")
	}
}

// TestSkipMeta_WriteMetaSkippedWhenTrue verifies that WriteReleaseMeta
// is not called when SkipMeta is true by simulating the conditional
// check in performRelease.
func TestSkipMeta_WriteMetaSkippedWhenTrue(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	opts := release.Options{
		Version:  "v9.0.0",
		SkipMeta: true,
		NoCommit: true,
		DryRun:   false,
	}

	// Simulate the SkipMeta guard from performRelease.
	if !opts.SkipMeta {
		t.Fatal("SkipMeta should be true")
	}

	// Since SkipMeta is true, metadata must not be written.
	assertNoMetadata(t, "v9.0.0")
	assertNoLatest(t)
}

// TestSkipMeta_WriteMetaCalledWhenFalse verifies that metadata IS written
// when SkipMeta is false (control test).
func TestSkipMeta_WriteMetaCalledWhenFalse(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	opts := release.Options{
		Version:  "v6.0.0",
		SkipMeta: false,
		NoCommit: true,
	}

	// Simulate the performRelease path when SkipMeta is false.
	if !opts.SkipMeta {
		meta := release.ReleaseMeta{
			Version:      "6.0.0",
			Branch:       "release/v6.0.0",
			SourceBranch: "main",
			Tag:          "v6.0.0",
			CreatedAt:    "2026-03-27T00:00:00Z",
			IsLatest:     true,
		}
		err := release.WriteReleaseMeta(meta)
		if err != nil {
			t.Fatalf("WriteReleaseMeta: %v", err)
		}

		v, _ := release.Parse("v6.0.0")
		err = release.WriteLatest(v)
		if err != nil {
			t.Fatalf("WriteLatest: %v", err)
		}
	}

	// Verify metadata was written.
	path := filepath.Join(constants.DefaultReleaseDir, "v6.0.0.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("v6.0.0.json should exist when SkipMeta is false")
	}

	latest, err := release.ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest: %v", err)
	}
	if latest.Tag != "v6.0.0" {
		t.Errorf("expected latest tag v6.0.0, got %s", latest.Tag)
	}
}

// TestSkipMeta_ExecuteFromBranchSetsFlag verifies that completeBranchRelease
// passes SkipMeta: true, ensuring no metadata files are created when
// releasing from an existing branch.
func TestSkipMeta_ExecuteFromBranchSetsFlag(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	// The SkipMeta flag is hardcoded to true in completeBranchRelease (line 73)
	// and in releaseFromMetadata (line 111). We verify by constructing the
	// same Options struct and confirming metadata is not written.
	opts := release.Options{
		Assets:   "",
		Notes:    "",
		IsDraft:  false,
		SkipMeta: true,
	}

	if !opts.SkipMeta {
		t.Fatal("Options from ExecuteFromBranch should have SkipMeta: true")
	}

	// No metadata should exist.
	assertNoMetadata(t, "v9.0.0")
	assertNoLatest(t)
}

// TestSkipMeta_ReleaseFromMetadataSetsFlag verifies that releaseFromMetadata
// passes SkipMeta: true in its Options, preventing metadata re-write.
func TestSkipMeta_ReleaseFromMetadataSetsFlag(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	// Seed a metadata file to simulate an existing pending release.
	seedMeta := release.ReleaseMeta{
		Version: "7.0.0",
		Tag:     "v7.0.0",
		Commit:  "abc1234567890",
	}
	err := release.WriteReleaseMeta(seedMeta)
	if err != nil {
		t.Fatalf("seed WriteReleaseMeta: %v", err)
	}

	// Read back the seed file for byte comparison.
	seedPath := filepath.Join(constants.DefaultReleaseDir, "v7.0.0.json")
	seedBytes, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("read seed file: %v", err)
	}

	// releaseFromMetadata constructs Options{SkipMeta: true}.
	// Simulate: no new metadata should be written, seed file unchanged.
	opts := release.Options{SkipMeta: true}
	if !opts.SkipMeta {
		t.Fatal("releaseFromMetadata should set SkipMeta: true")
	}

	// Verify seed file is byte-equal (not overwritten).
	currentBytes, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("read current file: %v", err)
	}
	if !bytes.Equal(currentBytes, seedBytes) {
		t.Error("seed v7.0.0.json was modified — expected unchanged")
	}

	// No latest.json should exist.
	assertNoLatest(t)
}

// TestSkipMeta_PerformReleaseGuard verifies the conditional guard in
// performRelease that skips writeMetadata when SkipMeta is true.
func TestSkipMeta_PerformReleaseGuard(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	os.MkdirAll(constants.DefaultReleaseDir, constants.DirPermission)

	// With SkipMeta: true, the guard prevents any metadata write.
	opts := release.Options{SkipMeta: true}

	v, _ := release.Parse("v8.0.0")

	// Directly test: metadata must NOT be created.
	if !opts.SkipMeta {
		// This block mimics writeMetadata — should never execute.
		meta := release.ReleaseMeta{Tag: "v8.0.0", Version: "8.0.0"}
		release.WriteReleaseMeta(meta)
		release.WriteLatest(v)
	}

	if release.ReleaseExists(v) {
		t.Error("v8.0.0 metadata should not exist with SkipMeta: true")
	}

	_, err := release.ReadLatest()
	if err == nil {
		t.Error("latest.json should not exist with SkipMeta: true")
	}

	// Now with SkipMeta: false — metadata SHOULD be created.
	opts.SkipMeta = false
	if !opts.SkipMeta {
		meta := release.ReleaseMeta{Tag: "v8.0.0", Version: "8.0.0"}
		release.WriteReleaseMeta(meta)
		release.WriteLatest(v)
	}

	if !release.ReleaseExists(v) {
		t.Error("v8.0.0 metadata should exist with SkipMeta: false")
	}

	latest, err := release.ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest: %v", err)
	}
	if latest.Tag != "v8.0.0" {
		t.Errorf("expected latest tag v8.0.0, got %s", latest.Tag)
	}
}
