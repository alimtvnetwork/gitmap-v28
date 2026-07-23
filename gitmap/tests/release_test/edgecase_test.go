package release_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// --- Pre-release version tests ---

// TestPreRelease_ParseAndString verifies parsing and round-tripping
// of pre-release version strings (e.g. v1.0.0-rc.1, v2.0.0-beta).
func TestPreRelease_ParseAndString(t *testing.T) {
	cases := []struct {
		input      string
		major      int
		minor      int
		patch      int
		preRelease string
		str        string
	}{
		{"v1.0.0-rc.1", 1, 0, 0, "rc.1", "v1.0.0-rc.1"},
		{"v2.3.0-beta", 2, 3, 0, "beta", "v2.3.0-beta"},
		{"v0.1.0-alpha.2", 0, 1, 0, "alpha.2", "v0.1.0-alpha.2"},
		{"1.0.0-rc.1", 1, 0, 0, "rc.1", "v1.0.0-rc.1"},
	}

	for _, tc := range cases {
		v, err := release.Parse(tc.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tc.input, err)
			continue
		}
		if v.Major != tc.major || v.Minor != tc.minor || v.Patch != tc.patch {
			t.Errorf("Parse(%q) = %d.%d.%d, want %d.%d.%d",
				tc.input, v.Major, v.Minor, v.Patch, tc.major, tc.minor, tc.patch)
		}
		if v.PreRelease != tc.preRelease {
			t.Errorf("Parse(%q).PreRelease = %q, want %q", tc.input, v.PreRelease, tc.preRelease)
		}
		if !v.IsPreRelease() {
			t.Errorf("Parse(%q).IsPreRelease() should be true", tc.input)
		}
		if v.String() != tc.str {
			t.Errorf("Parse(%q).String() = %q, want %q", tc.input, v.String(), tc.str)
		}
	}
}

// TestPreRelease_StableGreaterThanPreRelease verifies that a stable
// version is always greater than the same version with a pre-release tag.
func TestPreRelease_StableGreaterThanPreRelease(t *testing.T) {
	stable, _ := release.Parse("v1.0.0")
	rc, _ := release.Parse("v1.0.0-rc.1")

	if !stable.GreaterThan(rc) {
		t.Error("stable v1.0.0 should be greater than v1.0.0-rc.1")
	}
	if rc.GreaterThan(stable) {
		t.Error("v1.0.0-rc.1 should NOT be greater than stable v1.0.0")
	}
}

// TestPreRelease_MetadataPreservesPreRelease verifies that writing and
// reading metadata for a pre-release version preserves the suffix.
func TestPreRelease_MetadataPreservesPreRelease(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	meta := release.ReleaseMeta{
		Version: "1.0.0-rc.1",
		Tag:     "v1.0.0-rc.1",
		Branch:  "release/v1.0.0-rc.1",
	}
	if err := release.WriteReleaseMeta(meta); err != nil {
		t.Fatalf("WriteReleaseMeta: %v", err)
	}

	path := filepath.Join(constants.DefaultReleaseDir, "v1.0.0-rc.1.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}

	var got release.ReleaseMeta
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Version != "1.0.0-rc.1" {
		t.Errorf("version = %q, want %q", got.Version, "1.0.0-rc.1")
	}
	if got.Tag != "v1.0.0-rc.1" {
		t.Errorf("tag = %q, want %q", got.Tag, "v1.0.0-rc.1")
	}
}

// --- Bump resolution tests ---

// TestBump_AllLevels verifies that major, minor, and patch bumps
// produce correct version increments.
func TestBump_AllLevels(t *testing.T) {
	base, _ := release.Parse("v2.5.3")

	cases := []struct {
		level    string
		expected string
	}{
		{constants.BumpMajor, "v3.0.0"},
		{constants.BumpMinor, "v2.6.0"},
		{constants.BumpPatch, "v2.5.4"},
	}

	for _, tc := range cases {
		result, err := release.Bump(base, tc.level)
		if err != nil {
			t.Errorf("Bump(%s) error: %v", tc.level, err)
			continue
		}
		if result.String() != tc.expected {
			t.Errorf("Bump(%s) = %s, want %s", tc.level, result.String(), tc.expected)
		}
		if result.IsPreRelease() {
			t.Errorf("Bump(%s) should not produce a pre-release", tc.level)
		}
	}
}

// TestBump_InvalidLevel verifies that an unknown bump level returns an error.
func TestBump_InvalidLevel(t *testing.T) {
	base, _ := release.Parse("v1.0.0")
	_, err := release.Bump(base, "huge")
	if err == nil {
		t.Fatal("expected error for invalid bump level")
	}
}

// TestBump_FromPreRelease verifies that bumping a pre-release version
// strips the pre-release suffix and increments correctly.
func TestBump_FromPreRelease(t *testing.T) {
	pre, _ := release.Parse("v1.0.0-rc.1")

	patched, err := release.Bump(pre, constants.BumpPatch)
	if err != nil {
		t.Fatalf("Bump patch: %v", err)
	}
	if patched.String() != "v1.0.1" {
		t.Errorf("patch bump from pre-release = %s, want v1.0.1", patched.String())
	}
	if patched.IsPreRelease() {
		t.Error("patch bump should strip pre-release suffix")
	}

	minor, err := release.Bump(pre, constants.BumpMinor)
	if err != nil {
		t.Fatalf("Bump minor: %v", err)
	}
	if minor.String() != "v1.1.0" {
		t.Errorf("minor bump from pre-release = %s, want v1.1.0", minor.String())
	}
}

// TestBump_FromZero verifies bump behavior starting at v0.0.0.
func TestBump_FromZero(t *testing.T) {
	zero, _ := release.Parse("v0.0.0")

	patch, _ := release.Bump(zero, constants.BumpPatch)
	if patch.String() != "v0.0.1" {
		t.Errorf("patch from zero = %s, want v0.0.1", patch.String())
	}

	minor, _ := release.Bump(zero, constants.BumpMinor)
	if minor.String() != "v0.1.0" {
		t.Errorf("minor from zero = %s, want v0.1.0", minor.String())
	}

	major, _ := release.Bump(zero, constants.BumpMajor)
	if major.String() != "v1.0.0" {
		t.Errorf("major from zero = %s, want v1.0.0", major.String())
	}
}

// --- Parse edge cases ---

// TestParse_PartialVersions verifies that single and two-segment versions
// are zero-padded correctly.
func TestParse_PartialVersions(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"v1", "v1.0.0"},
		{"v2.3", "v2.3.0"},
		{"v10.20.30", "v10.20.30"},
	}

	for _, tc := range cases {
		v, err := release.Parse(tc.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tc.input, err)
			continue
		}
		if v.String() != tc.expected {
			t.Errorf("Parse(%q).String() = %q, want %q", tc.input, v.String(), tc.expected)
		}
	}
}

// TestParse_InvalidInputs verifies that malformed version strings
// produce errors.
func TestParse_InvalidInputs(t *testing.T) {
	invalids := []string{"v", "", "vabc", "v1.2.3.4", "v1.x.0"}
	for _, input := range invalids {
		_, err := release.Parse(input)
		if err == nil {
			t.Errorf("Parse(%q) should return error", input)
		}
	}
}

// TestGreaterThan_Ordering verifies version comparison across all segments.
func TestGreaterThan_Ordering(t *testing.T) {
	cases := []struct {
		a, b    string
		greater bool
	}{
		{"v2.0.0", "v1.0.0", true},
		{"v1.1.0", "v1.0.0", true},
		{"v1.0.1", "v1.0.0", true},
		{"v1.0.0", "v1.0.0", false},
		{"v1.0.0", "v2.0.0", false},
		{"v1.0.0", "v1.0.0-rc.1", true},
	}

	for _, tc := range cases {
		a, _ := release.Parse(tc.a)
		b, _ := release.Parse(tc.b)
		got := a.GreaterThan(b)
		if got != tc.greater {
			t.Errorf("%s.GreaterThan(%s) = %v, want %v", tc.a, tc.b, got, tc.greater)
		}
	}
}

// --- Multi-release sequence test ---

// TestMultiRelease_SequentialVersions verifies that three sequential
// releases produce correct metadata and latest.json always points to
// the highest version.
func TestMultiRelease_SequentialVersions(t *testing.T) {
	_, releaseDir, cleanup := initE2ERepo(t)
	defer cleanup()

	versions := []string{"v1.0.0", "v1.1.0", "v1.2.0"}

	for _, ver := range versions {
		err := release.Execute(release.Options{
			Version:  ver,
			DryRun:   false,
			NoCommit: true,
			SkipMeta: false,
		})
		if err != nil {
			t.Fatalf("Execute(%s): %v", ver, err)
		}
	}

	// All three metadata files should exist.
	for _, ver := range versions {
		metaPath := filepath.Join(releaseDir, ver+".json")
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			t.Errorf("%s.json should exist", ver)
		}
	}

	// latest.json should point to the highest version.
	latest, err := release.ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest: %v", err)
	}
	if latest.Tag != "v1.2.0" {
		t.Errorf("latest tag = %s, want v1.2.0", latest.Tag)
	}

	// All tags should exist locally.
	for _, ver := range versions {
		if !release.TagExistsLocally(ver) {
			t.Errorf("tag %s should exist", ver)
		}
	}
}

// TestMultiRelease_MixedBumps verifies a release sequence using
// different bump levels produces correct versions.
func TestMultiRelease_MixedBumps(t *testing.T) {
	base, _ := release.Parse("v1.0.0")

	// Simulate: patch → minor → major
	patched, _ := release.Bump(base, constants.BumpPatch)
	if patched.String() != "v1.0.1" {
		t.Fatalf("step 1: got %s, want v1.0.1", patched.String())
	}

	minored, _ := release.Bump(patched, constants.BumpMinor)
	if minored.String() != "v1.1.0" {
		t.Fatalf("step 2: got %s, want v1.1.0", minored.String())
	}

	majored, _ := release.Bump(minored, constants.BumpMajor)
	if majored.String() != "v2.0.0" {
		t.Fatalf("step 3: got %s, want v2.0.0", majored.String())
	}
}

// TestMultiRelease_OutOfOrderBlocked verifies that releasing a version
// lower than an existing tag still succeeds (tags are independent),
// but metadata files coexist.
func TestMultiRelease_OutOfOrderMetadata(t *testing.T) {
	cleanup := setupSkipMetaDir(t)
	defer cleanup()

	// Write v2.0.0 first, then v1.5.0.
	for _, ver := range []string{"v2.0.0", "v1.5.0"} {
		meta := release.ReleaseMeta{
			Version: ver[1:], // strip v
			Tag:     ver,
		}
		if err := release.WriteReleaseMeta(meta); err != nil {
			t.Fatalf("WriteReleaseMeta(%s): %v", ver, err)
		}
	}

	// Both files should exist.
	for _, ver := range []string{"v2.0.0", "v1.5.0"} {
		path := filepath.Join(constants.DefaultReleaseDir, ver+".json")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s.json should exist", ver)
		}
	}
}

// TestMultiRelease_PreReleaseToStable verifies releasing rc then stable
// for the same version series.
func TestMultiRelease_PreReleaseToStable(t *testing.T) {
	_, releaseDir, cleanup := initE2ERepo(t)
	defer cleanup()

	// Release rc first.
	err := release.Execute(release.Options{
		Version:  "v3.0.0-rc.1",
		DryRun:   false,
		NoCommit: true,
		SkipMeta: false,
	})
	if err != nil {
		t.Fatalf("Execute(rc): %v", err)
	}

	// Commit something to allow a new release.
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "prep stable")
	cmd.Dir, _ = os.Getwd()
	if err := cmd.Run(); err != nil {
		t.Fatalf("prep commit: %v", err)
	}

	// Release stable.
	err = release.Execute(release.Options{
		Version:  "v3.0.0",
		DryRun:   false,
		NoCommit: true,
		SkipMeta: false,
	})
	if err != nil {
		t.Fatalf("Execute(stable): %v", err)
	}

	// Both metadata files should exist.
	rcPath := filepath.Join(releaseDir, "v3.0.0-rc.1.json")
	stablePath := filepath.Join(releaseDir, "v3.0.0.json")
	if _, err := os.Stat(rcPath); os.IsNotExist(err) {
		t.Error("v3.0.0-rc.1.json should exist")
	}
	if _, err := os.Stat(stablePath); os.IsNotExist(err) {
		t.Error("v3.0.0.json should exist")
	}

	// latest.json should point to stable.
	latest, err := release.ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest: %v", err)
	}
	if latest.Tag != "v3.0.0" {
		t.Errorf("latest = %s, want v3.0.0", latest.Tag)
	}
}
