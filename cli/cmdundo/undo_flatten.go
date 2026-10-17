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
	if repoName == "-h" || repoName == "--help" || repoName == "help" {
		fmt.Printf("Usage: gitmap repo undo-flatten [repo]\n\nRestore a flattened repository from its <repo>-marked-for-deletion backup archive.\n")
		return nil
	}

	trimmed, err := resolveDefaultRepoName(repoName)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(trimmed)
	if err != nil {
		return fmt.Errorf("resolve repository path %q: %w", trimmed, err)
	}

	parentDir := filepath.Dir(absPath)
	baseName := filepath.Base(absPath)
	canonicalName := strings.TrimSuffix(baseName, "-marked-for-deletion")
	archiveName := canonicalName + "-marked-for-deletion"

	resolvedParent, canonicalPath, archivePath, err := locateArchivePaths(parentDir, canonicalName, archiveName)
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	if isInsideOrSameDir(cwd, canonicalPath) {
		_ = os.Chdir(resolvedParent)
	}

	if err := removeExistingCanonicalPath(canonicalPath); err != nil {
		return err
	}

	fmt.Printf("Restoring %s from %s...\n", canonicalName, archiveName)
	if err := os.Rename(archivePath, canonicalPath); err != nil {
		return fmt.Errorf("restore archive %s to %s: %w", archivePath, canonicalPath, err)
	}

	fmt.Printf("\n%s✓ Repository %s successfully restored from %s!%s\n", constants.ColorGreen, canonicalName, archiveName, constants.ColorReset)
	fmt.Printf("  Restored Path: %s\n\n", canonicalPath)

	return nil
}

func resolveDefaultRepoName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed != "" && trimmed != "." {
		return trimmed, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve current working directory: %w", err)
	}
	return filepath.Base(cwd), nil
}

func locateArchivePaths(parentDir, canonicalName, archiveName string) (string, string, string, error) {
	canonicalPath := filepath.Join(parentDir, canonicalName)
	archivePath := filepath.Join(parentDir, archiveName)
	if _, err := os.Stat(archivePath); err == nil {
		return parentDir, canonicalPath, archivePath, nil
	}
	cwd, _ := os.Getwd()
	cwdArchive := filepath.Join(cwd, archiveName)
	if _, cwdErr := os.Stat(cwdArchive); cwdErr == nil {
		return cwd, filepath.Join(cwd, canonicalName), cwdArchive, nil
	}
	return "", "", "", fmt.Errorf("archive %q not found for repository %q (expected %s)", archiveName, canonicalName, archivePath)
}

func removeExistingCanonicalPath(canonicalPath string) error {
	if _, err := os.Stat(canonicalPath); os.IsNotExist(err) {
		return nil
	}
	fmt.Printf("Removing flattened working tree %s...\n", canonicalPath)
	if err := os.RemoveAll(canonicalPath); err != nil {
		return fmt.Errorf("remove existing flattened repository %s: %w", canonicalPath, err)
	}
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
