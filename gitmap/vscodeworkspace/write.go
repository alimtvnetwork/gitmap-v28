package vscodeworkspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// WriteAtomic encodes ws and commits it to outPath via temp-file +
// rename. Same atomic-rename pattern vscodepm.Sync uses for
// projects.json — partial writes are never observed, even if the
// process is killed mid-flush.
func WriteAtomic(outPath string, ws Workspace) error {
	bytesOut, err := Encode(ws)
	if err != nil {
		return err
	}

	tmpPath := outPath + constants.VSCodeWorkspaceTempSuffix
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf(constants.ErrVSCodeWorkspaceWriteTemp, tmpPath, err)
	}
	if err := os.WriteFile(tmpPath, bytesOut, 0o644); err != nil {
		return fmt.Errorf(constants.ErrVSCodeWorkspaceWriteTemp, tmpPath, err)
	}
	if err := os.Rename(tmpPath, outPath); err != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf(constants.ErrVSCodeWorkspaceRename, outPath, err)
	}

	return nil
}
