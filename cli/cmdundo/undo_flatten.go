package cmdundo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo/constants"
)

// RunUndoFlatten restores a flattened repository from its marked-for-deletion archive.
func RunUndoFlatten(repoName string) error {
	if isHelpRequest(repoName) {
		fmt.Print(constants.MsgFlattenUsage)

		return nil
	}

	trimmed, err := resolveDefaultRepoName(repoName)
	if err != nil {
		return err
	}

	return executeFlattenRestore(trimmed)
}

func isHelpRequest(repoName string) bool {
	switch repoName {
	case constants.FlagHelpShort, constants.FlagHelpLong, constants.FlagHelpWord:
		return true
	default:
		return false
	}
}

func executeFlattenRestore(repoName string) error {
	absPath, err := filepath.Abs(repoName)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.OpResolveRepoPathFmt, repoName))
	}

	parentDir := filepath.Dir(absPath)
	canonicalName := strings.TrimSuffix(filepath.Base(absPath), constants.SuffixMarkedForDel)
	archiveName := canonicalName + constants.SuffixMarkedForDel

	return performFlattenRestore(parentDir, canonicalName, archiveName)
}

func performFlattenRestore(parentDir, canonicalName, archiveName string) error {
	resolvedParent, canonicalPath, archivePath, err := locateArchivePaths(parentDir, canonicalName, archiveName)
	if err != nil {
		return err
	}

	ensureWorkingDirNotInTarget(canonicalPath, resolvedParent)
	if err := removeExistingCanonicalPath(canonicalPath); err != nil {
		return err
	}

	return moveArchiveToCanonical(archivePath, canonicalPath, canonicalName, archiveName)
}

func ensureWorkingDirNotInTarget(canonicalPath, resolvedParent string) {
	cwd, _ := os.Getwd()
	if isInsideOrSameDir(cwd, canonicalPath) {
		_ = os.Chdir(resolvedParent)
	}
}

func moveArchiveToCanonical(archivePath, canonicalPath, canonicalName, archiveName string) error {
	fmt.Printf(constants.MsgFlattenRestoringFmt, canonicalName, archiveName)
	if err := os.Rename(archivePath, canonicalPath); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.OpRestoreArchiveFmt, archivePath, canonicalPath))
	}

	fmt.Printf(constants.MsgFlattenRestoredSuccessFmt, constants.ColorGreen, canonicalName, archiveName, constants.ColorReset)
	fmt.Printf(constants.MsgFlattenRestoredPathFmt, canonicalPath)

	return nil
}

func resolveDefaultRepoName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed != "" && trimmed != constants.SentinelCurrentDir {
		return trimmed, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", apperror.WrapSimple(err, constants.OpResolveCwd)
	}

	return filepath.Base(cwd), nil
}

func locateArchivePaths(parentDir, canonicalName, archiveName string) (string, string, string, error) {
	canonicalPath := filepath.Join(parentDir, canonicalName)
	archivePath := filepath.Join(parentDir, archiveName)
	if hasArchivePath(archivePath) {
		return parentDir, canonicalPath, archivePath, nil
	}

	cwd, _ := os.Getwd()
	cwdArchive := filepath.Join(cwd, archiveName)
	if hasArchivePath(cwdArchive) {
		return cwd, filepath.Join(cwd, canonicalName), cwdArchive, nil
	}

	return "", "", "", apperror.NewNotFoundError(fmt.Sprintf(constants.ErrArchiveNotFoundFmt, archiveName, canonicalName, archivePath))
}

func hasArchivePath(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func removeExistingCanonicalPath(canonicalPath string) error {
	if !hasCanonicalPath(canonicalPath) {
		return nil
	}

	fmt.Printf(constants.MsgFlattenRemovingTreeFmt, canonicalPath)
	if err := os.RemoveAll(canonicalPath); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.OpRemoveExistingFlattenedFmt, canonicalPath))
	}

	return nil
}

func hasCanonicalPath(canonicalPath string) bool {
	_, err := os.Stat(canonicalPath)

	return !os.IsNotExist(err)
}

// RunUndoFlattenCLI is a CLI adapter for RunUndoFlatten that accepts args.
func RunUndoFlattenCLI(args []string) error {
	if hasHelpFlag(args) {
		return RunUndoFlatten(constants.FlagHelpShort)
	}

	target := extractFlattenTarget(args)

	return RunUndoFlatten(target)
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if isHelpRequest(arg) {
			return true
		}
	}

	return false
}

func extractFlattenTarget(args []string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
	}

	return ""
}

func isInsideOrSameDir(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}
