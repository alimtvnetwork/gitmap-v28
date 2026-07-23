package cmd

import "testing"

// TestParseAddLFSInstallFlagsDefaults verifies the zero-flag invocation
// produces a no-dry-run config — the default path that actually writes
// to disk.
func TestParseAddLFSInstallFlagsDefaults(t *testing.T) {
	got := parseAddLFSInstallFlags(nil)
	if got.dryRun {
		t.Errorf("dryRun: want false by default, got true")
	}
}

// TestParseAddLFSInstallFlagsDryRun verifies the --dry-run flag flips
// the bool. This is the safety net users rely on before committing.
func TestParseAddLFSInstallFlagsDryRun(t *testing.T) {
	got := parseAddLFSInstallFlags([]string{"--dry-run"})
	if !got.dryRun {
		t.Errorf("dryRun: want true with --dry-run, got false")
	}
}

// TestAddLFSInstallTagIsStable locks in the marker tag so an accidental
// rename doesn't orphan blocks already on disk in users' repos.
func TestAddLFSInstallTagIsStable(t *testing.T) {
	if addLFSInstallTag != "lfs/common" {
		t.Errorf("marker tag drifted: want %q, got %q", "lfs/common", addLFSInstallTag)
	}
}
