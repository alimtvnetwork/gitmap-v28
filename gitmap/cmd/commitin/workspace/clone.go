package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// StagedInput is the on-disk result of CloneInputs for one entry.
// WorkPath is the path the walker reads from. For local folders this
// is the original AbsPath (read-only walk, no copy needed). For URLs
// and versioned siblings we always have a fresh path under TempRoot.
type StagedInput struct {
	Input    ResolvedInput
	WorkPath string
	IsClone  bool // true when WorkPath was created by us under TempRoot
}

// CloneInputs implements spec §3.1 stage 08. For each ResolvedInput:
//   - GitUrl              → git clone into <TempRoot>/<runId>/<idx>-<basename>
//   - VersionedSibling    → git clone (local path) into temp (full history)
//   - LocalFolder         → reuse AbsPath in place; NO copy
//
// runID anchors the temp subtree so concurrent runs never collide.
// The directory is created on first call; subsequent calls reuse it.
func CloneInputs(p *Paths, runID int64, inputs []ResolvedInput) ([]StagedInput, error) {
	runDir, err := ensureRunTempDir(p, runID)
	if err != nil {
		return nil, err
	}
	out := make([]StagedInput, 0, len(inputs))
	for _, in := range inputs {
		staged, stageErr := stageOneInput(runDir, in)
		if stageErr != nil {
			return nil, stageErr
		}
		out = append(out, staged)
	}
	return out, nil
}

// ensureRunTempDir creates <TempRoot>/<runId>/ once.
func ensureRunTempDir(p *Paths, runID int64) (string, error) {
	dir := filepath.Join(p.TempRoot, fmt.Sprintf("%d", runID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir temp run dir %s: %w", dir, err)
	}
	return dir, nil
}

// stageOneInput dispatches per Kind. Each branch returns a populated
// StagedInput or wraps the underlying error in the spec §2.7 message.
func stageOneInput(runDir string, in ResolvedInput) (StagedInput, error) {
	switch in.Kind {
	case constants.CommitInInputKindLocalFolder:
		return stageLocalFolder(in)
	case constants.CommitInInputKindGitUrl:
		return stageRemoteUrl(runDir, in)
	case constants.CommitInInputKindVersionedSibling:
		return stageVersionedSibling(runDir, in)
	}
	return StagedInput{}, fmt.Errorf(constants.CommitInErrInputOpen, in.Original, fmt.Errorf("unknown kind %q", in.Kind))
}

// stageLocalFolder verifies the directory exists (read-only walk).
func stageLocalFolder(in ResolvedInput) (StagedInput, error) {
	info, err := os.Stat(in.AbsPath)
	if err != nil {
		return StagedInput{}, fmt.Errorf(constants.CommitInErrInputOpen, in.Original, err)
	}
	if !info.IsDir() {
		return StagedInput{}, fmt.Errorf(constants.CommitInErrInputOpen, in.Original, fmt.Errorf("not a directory"))
	}
	return StagedInput{Input: in, WorkPath: in.AbsPath}, nil
}

// stageRemoteUrl runs `git clone <url> <runDir>/<idx>-<basename>`.
func stageRemoteUrl(runDir string, in ResolvedInput) (StagedInput, error) {
	target := filepath.Join(runDir, fmt.Sprintf(constants.CommitInTempInputFormat, in.OrderIndex, cloneBasename(in.URL)))
	if err := gitRunner("clone", in.URL, target); err != nil {
		return StagedInput{}, fmt.Errorf(constants.CommitInErrInputClone, in.Original, err)
	}
	return StagedInput{Input: in, WorkPath: target, IsClone: true}, nil
}

// stageVersionedSibling clones the local sibling so the walker sees a
// pristine history with no working-tree pollution.
func stageVersionedSibling(runDir string, in ResolvedInput) (StagedInput, error) {
	target := filepath.Join(runDir, fmt.Sprintf(constants.CommitInTempInputFormat, in.OrderIndex, filepath.Base(in.AbsPath)))
	if err := gitRunner("clone", in.AbsPath, target); err != nil {
		return StagedInput{}, fmt.Errorf(constants.CommitInErrInputClone, in.Original, err)
	}
	return StagedInput{Input: in, WorkPath: target, IsClone: true}, nil
}
