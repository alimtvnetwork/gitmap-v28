package cmddiff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/diff"
)

func TestDispatchDiffUnrecognized(t *testing.T) {
	handled, err := DispatchDiff("other")
	if handled || err != nil {
		t.Errorf("expected false, nil for unrecognized command, got handled=%v, err=%v", handled, err)
	}
}

func TestDispatchCompareUnrecognized(t *testing.T) {
	handled, err := DispatchCompare("other")
	if handled || err != nil {
		t.Errorf("expected false, nil for unrecognized command, got handled=%v, err=%v", handled, err)
	}
}

func TestIsCompareCommand(t *testing.T) {
	if !isCompareCommand(constants.CmdCompare) {
		t.Errorf("expected %s to be recognized as compare command", constants.CmdCompare)
	}
	if !isCompareCommand(constants.CmdCompareAlias) {
		t.Errorf("expected %s to be recognized as compare command", constants.CmdCompareAlias)
	}
	if isCompareCommand("invalid") {
		t.Error("expected invalid command to not be recognized as compare command")
	}
}

func TestGuardDiffPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "diffguard_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ep1 := diff.Endpoint{WorkingDir: tempDir}
	ep2 := diff.Endpoint{WorkingDir: tempDir}
	if err := guardDiffPaths(ep1, ep2); err == nil {
		t.Error("expected guardDiffPaths to fail when left and right are identical directories")
	}

	subDir := filepath.Join(tempDir, "sub")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("failed to create sub dir: %v", err)
	}
	ep3 := diff.Endpoint{WorkingDir: subDir}
	if err := guardDiffPaths(ep1, ep3); err != nil {
		t.Errorf("expected guardDiffPaths to pass for different directories, got: %v", err)
	}
}
