package runlog

import "testing"

// TestResolveReleaseBranchNameDefaultOn locks the spec default: the
// flag is OFF, tag is annotated semver, real run → branch IS created.
func TestResolveReleaseBranchNameDefaultOn(t *testing.T) {
	got := ResolveReleaseBranchName("v1.2.3", true, false, false)
	if got != "release/v1.2.3" {
		t.Fatalf("default ON: got %q, want release/v1.2.3", got)
	}
}

// TestResolveReleaseBranchNameSuppressedByFlag locks the §08 T4 case
// from the acceptance matrix: --no-release-branch wins even on an
// annotated version tag.
func TestResolveReleaseBranchNameSuppressedByFlag(t *testing.T) {
	if got := ResolveReleaseBranchName("v2.0.0", true, true, false); got != "" {
		t.Fatalf("--no-release-branch should suppress, got %q", got)
	}
}

// TestResolveReleaseBranchNameNonVersionTagNeverGetsBranch covers the
// "tag is annotated but not SemVer" case (e.g. `nightly`).
func TestResolveReleaseBranchNameNonVersionTagNeverGetsBranch(t *testing.T) {
	for _, name := range []string{"nightly", "release-1.0", "v1.2", ""} {
		if got := ResolveReleaseBranchName(name, true, false, false); got != "" {
			t.Errorf("non-version tag %q: got %q, want \"\"", name, got)
		}
	}
}

// TestResolveReleaseBranchNameLightweightNeverGetsBranch locks the
// strict-semver contract: a lightweight tag whose name happens to
// look like `v1.2.3` MUST NOT get a release branch.
func TestResolveReleaseBranchNameLightweightNeverGetsBranch(t *testing.T) {
	for _, name := range []string{"v1.2.3", "1.2.3", "v2.0.0-rc.1", "v3.0.0+build.7"} {
		if got := ResolveReleaseBranchName(name, false, false, false); got != "" {
			t.Errorf("lightweight %q: got %q, want \"\" (annotated-only)", name, got)
		}
	}
}

// TestResolveReleaseBranchNameDryRunSuppresses covers spec §9.4 R6:
// dry-run sets MirroredReleaseBranch to NULL.
func TestResolveReleaseBranchNameDryRunSuppresses(t *testing.T) {
	if got := ResolveReleaseBranchName("v1.0.0", true, false, true); got != "" {
		t.Fatalf("dry-run should suppress, got %q", got)
	}
}

// TestResolveReleaseBranchNameFlagBeatsEverything proves the flag is
// the highest-priority signal — even with a perfectly-valid annotated
// version tag and a real run, --no-release-branch returns "".
func TestResolveReleaseBranchNameFlagBeatsEverything(t *testing.T) {
	if got := ResolveReleaseBranchName("v1.0.0-rc.1", true, true, false); got != "" {
		t.Fatalf("flag should beat valid version tag, got %q", got)
	}
}
