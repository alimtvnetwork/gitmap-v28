// Package cmd — lowercasefix_ops.go implements scanning and two-step git renaming.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExecuteLowerCaseFix orchestrates scanning, renaming, and committing.
func ExecuteLowerCaseFix(opts LowerCaseFixOptions) error {
	cwd, _ := os.Getwd()
	pairs, err := findRenameCandidates(cwd, opts.Patterns)
	if err != nil {
		return err
	}
	if len(pairs) == 0 {
		fmt.Printf("%s No uppercase files matching patterns found to rename.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}

	renderRenamePreview(pairs)
	if opts.IsDryRun {
		fmt.Printf("\n%sℹ [dry-run] %d file(s) would be renamed to lowercase.%s\n\n", constants.ColorYellow, len(pairs), constants.ColorReset)
		return nil
	}

	return applyRenamesAndCommit(pairs, opts)
}

func applyRenamesAndCommit(pairs []RenamePair, opts LowerCaseFixOptions) error {
	for _, p := range pairs {
		if err := performSafeGitRename(p); err != nil {
			return err
		}
		fmt.Printf("  %s %s → %s\n", constants.ColorGreen+"✓"+constants.ColorReset, p.OldBase, p.NewBase)
	}

	if opts.IsCommit {
		return commitRenames(pairs, opts.CommitMessage)
	}
	return nil
}

func performSafeGitRename(p RenamePair) error {
	tmpPath := p.OldPath + ".tmp-lcf"
	if err := runGitMvForLcf(p.OldPath, tmpPath); err != nil {
		return fallbackRename(p, tmpPath, err)
	}
	return runGitMvForLcf(tmpPath, p.NewPath)
}

func fallbackRename(p RenamePair, tmpPath string, gitErr error) error {
	if err := os.Rename(p.OldPath, tmpPath); err != nil {
		return fmt.Errorf("git mv failed (%v) and rename to tmp failed: %w", gitErr, err)
	}
	if err := os.Rename(tmpPath, p.NewPath); err != nil {
		return fmt.Errorf("rename tmp to %s failed: %w", p.NewPath, err)
	}
	return nil
}

func runGitMvForLcf(src, dst string) error {
	cmd := exec.Command("git", "mv", src, dst)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git mv %s to %s failed: %w (output: %s)", src, dst, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func renderRenamePreview(pairs []RenamePair) {
	fmt.Printf("\n%s⚡ Lowercase File Rename Candidates:%s\n\n", constants.ColorCyan, constants.ColorReset)
	for _, p := range pairs {
		rel, _ := filepath.Rel(".", p.OldPath)
		fmt.Printf("  %s•%s %s %s(%s → %s)%s\n",
			constants.ColorYellow, constants.ColorReset,
			rel, constants.ColorDim, p.OldBase, p.NewBase, constants.ColorReset)
	}
}
