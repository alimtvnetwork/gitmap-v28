package vscodeworkspace

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
		return apperror.WrapSimple(err, fmt.Sprintf("vscode-workspace: failed to write temp file %q", tmpPath))
	}

	if err := os.WriteFile(tmpPath, bytesOut, 0o644); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("vscode-workspace: failed to write temp file %q", tmpPath))
	}

	if err := os.Rename(tmpPath, outPath); err != nil {
		_ = os.Remove(tmpPath)

		return apperror.WrapSimple(err, fmt.Sprintf("vscode-workspace: failed to commit %q", outPath))
	}

	return nil
}
