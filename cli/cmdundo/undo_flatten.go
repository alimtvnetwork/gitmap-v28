package cmdundo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunUndoFlatten restores a flattened repository from its -marked-for-deletion archive.
func RunUndoFlatten(repoName string) error {
	trimmed := strings.TrimSpace(repoName)
	if trimmed == "" || trimmed == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve current working directory: %w", err)
		}

		trimmed = filepath.Base(cwd)
	}

	absPath, err := filepath.Abs(trimmed)
	if err != nil {
		return fmt.Errorf("resolve repository path %q: %w", trimmed, err)
	}

	parentDir := filepath.Dir(absPath)
	baseName := filepath.Base(absPath)
	canonicalName := strings.TrimSuffix(baseName, "-marked-for-deletion")
	archiveName := canonicalName + "-marked-for-deletion"

	canonicalPath := filepath.Join(parentDir, canonicalName)
	archivePath := filepath.Join(parentDir, archiveName)

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		cwdArchive := filepath.Join(cwd, archiveName)
		if _, cwdErr := os.Stat(cwdArchive); cwdErr == nil {
			archivePath = cwdArchive
			canonicalPath = filepath.Join(cwd, canonicalName)
			parentDir = cwd
		} else {
			return fmt.Errorf("archive %q not found for repository %q (expected %s)", archiveName, canonicalName, archivePath)
		}
	}

	cwd, _ := os.Getwd()
	if isInsideOrSameDir(cwd, canonicalPath) {
		_ = os.Chdir(parentDir)
	}

	if _, err := os.Stat(canonicalPath); err == nil {
		fmt.Printf("Removing flattened working tree %s...\n", canonicalPath)
		if err := os.RemoveAll(canonicalPath); err != nil {
			return fmt.Errorf("remove existing flattened repository %s: %w", canonicalPath, err)
		}
	}

	fmt.Printf("Restoring %s from %s...\n", canonicalName, archiveName)
	if err := os.Rename(archivePath, canonicalPath); err != nil {
		return fmt.Errorf("restore archive %s to %s: %w", archivePath, canonicalPath, err)
	}

	fmt.Printf("\n%s✓ Repository %s successfully restored from %s!%s\n", constants.ColorGreen, canonicalName, archiveName, constants.ColorReset)
	fmt.Printf("  Restored Path: %s\n\n", canonicalPath)

	return nil
}

// RunUndoFlattenCLI is a CLI adapter for RunUndoFlatten that accepts args.
func RunUndoFlattenCLI(args []string) error {
	target := ""
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") && target == "" {
			target = arg
		}
	}

	return RunUndoFlatten(target)
}

func isInsideOrSameDir(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}
