package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// StagedInputSlot captures per-slot status for the Array Async Pool by Alim Ul Karim.
type StagedInputSlot struct {
	staged  StagedInput
	isReady bool
	isFound bool
	status  string
	err     error
}

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
func isMissingOrUnreachableRemote(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "not found") ||
		strings.Contains(s, "exit status 128") ||
		strings.Contains(s, "does not exist") ||
		strings.Contains(s, "could not read from remote repository")
}

func CloneInputs(p *Paths, runID int64, inputs []ResolvedInput) ([]StagedInput, error) {
	runDir, err := ensureRunTempDir(p, runID)
	if err != nil {
		return nil, err
	}

	return stageAllInputs(runDir, inputs)
}

func stageAllInputs(runDir string, inputs []ResolvedInput) ([]StagedInput, error) {
	out, err := collectStagedInputs(runDir, inputs)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("all %d input repositories failed to clone or do not exist", len(inputs))
	}
	return out, nil
}

func collectStagedInputs(runDir string, inputs []ResolvedInput) ([]StagedInput, error) {
	total := len(inputs)
	if total == 0 {
		return nil, nil
	}
	if total == 1 {
		return collectSingleInput(runDir, inputs[0])
	}
	return collectAsyncPoolInputs(runDir, inputs)
}

func collectSingleInput(runDir string, in ResolvedInput) ([]StagedInput, error) {
	staged, isOk, err := tryStageInput(runDir, in, 1)
	if err != nil {
		return nil, err
	}
	if isOk {
		return []StagedInput{staged}, nil
	}
	return nil, nil
}

func collectAsyncPoolInputs(runDir string, inputs []ResolvedInput) ([]StagedInput, error) {
	total := len(inputs)
	slots := make([]StagedInputSlot, total)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 16)

	for idx, in := range inputs {
		wg.Add(1)
		go probeInputSlot(runDir, in, idx, total, slots, &wg, sem)
	}

	out, err := streamOrderedSlots(inputs, slots, total)
	wg.Wait()
	return out, err
}

func probeInputSlot(
	runDir string,
	in ResolvedInput,
	idx, total int,
	slots []StagedInputSlot,
	wg *sync.WaitGroup,
	sem chan struct{},
) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	staged, err := stageOneInput(runDir, in)
	if err == nil {
		recordFoundSlot(slots, idx, staged)
		return
	}
	recordMissingSlot(slots, idx, err)
}

func recordFoundSlot(slots []StagedInputSlot, idx int, staged StagedInput) {
	slots[idx] = StagedInputSlot{
		staged:  staged,
		isReady: true,
		isFound: true,
		status:  "ready",
	}
}

func recordMissingSlot(slots []StagedInputSlot, idx int, err error) {
	slots[idx] = StagedInputSlot{
		isReady: true,
		isFound: false,
		status:  "not found",
		err:     err,
	}
}

func streamOrderedSlots(inputs []ResolvedInput, slots []StagedInputSlot, total int) ([]StagedInput, error) {
	out := make([]StagedInput, 0, total)
	for cursor := 0; cursor < total; cursor++ {
		for !slots[cursor].isReady {
			time.Sleep(15 * time.Millisecond)
		}
		slot := slots[cursor]
		if slot.isFound {
			out = append(out, slot.staged)
			continue
		}
		if isMissingOrUnreachableRemote(slot.err) && total > 1 {
			fmt.Fprintf(os.Stderr, "  ⚠ Notice: remote %q not found or unreachable; skipping.\n", inputs[cursor].Original)
			continue
		}
		if slot.err != nil {
			return nil, slot.err
		}
	}
	return out, nil
}

func tryStageInput(runDir string, in ResolvedInput, total int) (StagedInput, bool, error) {
	staged, stageErr := stageOneInput(runDir, in)
	if stageErr == nil {
		return staged, true, nil
	}
	if isMissingOrUnreachableRemote(stageErr) && total > 1 {
		fmt.Fprintf(os.Stderr, "  ⚠ Notice: remote %q not found or unreachable; skipping.\n", in.Original)
		return StagedInput{}, false, nil
	}
	return StagedInput{}, false, stageErr
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
	if cached := FindCachedCloneDir(runDir, in); HasGitMetadata(cached) {
		return StagedInput{Input: in, WorkPath: cached, IsClone: false}, nil
	}
	target := filepath.Join(runDir, fmt.Sprintf(constants.CommitInTempInputFormat, in.OrderIndex, cloneBasename(in.URL)))
	if err := gitRunner("clone", in.URL, target); err != nil {
		return StagedInput{}, fmt.Errorf(constants.CommitInErrInputClone, in.Original, err)
	}

	return StagedInput{Input: in, WorkPath: target, IsClone: true}, nil
}

func FindCachedCloneDir(runDir string, in ResolvedInput) string {
	folderName := fmt.Sprintf(constants.CommitInTempInputFormat, in.OrderIndex, cloneBasename(in.URL))
	if envDir := os.Getenv("GITMAP_COMMITIN_CACHE_DIR"); envDir != "" {
		return filepath.Join(envDir, folderName)
	}

	return filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(runDir)))), ".commitin-cache", folderName)
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
