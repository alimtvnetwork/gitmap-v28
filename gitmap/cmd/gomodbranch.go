package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// requireInsideWorkTree exits if the current directory is not inside a git repo.
func requireInsideWorkTree() error {
	if gitutil.IsInsideWorkTree() {
		return nil
	}

	fmt.Fprint(os.Stderr, constants.ErrGoModNotRepo)

	return apperror.NewSimple("fatal error", "E9000")
}

// deriveSlug sanitizes a module path into a branch-safe slug.
func deriveSlug(modulePath string) string {
	slug := strings.ReplaceAll(modulePath, "/", "-")
	slug = strings.ReplaceAll(slug, ".", "-")
	slug = strings.ReplaceAll(slug, "@", "-")

	return slug
}

func setupGoModBranch(branch string) {
	if err := ensureBranchNotExists(branch); err != nil {
		cliexit.HandleGeneralError(err)
	}
	if err := createBranchAtHead(branch); err != nil {
		cliexit.HandleGeneralError(err)
	}
}

// createGoModBranches creates backup and feature branches from current HEAD.
func createGoModBranches(slug string) (string, string) {
	backupBranch := constants.GoModBackupPrefix + slug
	featureBranch := constants.GoModFeaturePrefix + slug
	setupGoModBranch(backupBranch)
	setupGoModBranch(featureBranch)
	checkoutBranch(featureBranch)

	return backupBranch, featureBranch
}

// ensureBranchNotExists aborts if the branch already exists.
func ensureBranchNotExists(branch string) error {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, constants.GitBranchListFlag, branch)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	if len(strings.TrimSpace(string(out))) > 0 {
		return apperror.NewSimple(constants.ErrGoModBranchExists, "E9000")
	}

	return nil
}

// createBranchAtHead creates a branch at the current HEAD without checking it out.
func createBranchAtHead(branch string) error {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, branch)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return apperror.NewSimple(constants.ErrGoModBranchExists, "E9000")
	}

	return nil
}

// checkoutBranch checks out the given branch.
func checkoutBranch(branch string) {
	cmd := exec.Command(constants.GitBin, constants.GitCheckout, branch)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not checkout branch %s: %v\n", branch, err)
	}
}

// goModCurrentBranch returns the name of the current branch.
func goModCurrentBranch() string {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return constants.DefaultBranch
	}

	return strings.TrimSpace(string(out))
}

// isWorkTreeDirty checks if there are uncommitted changes.
func isWorkTreeDirty() bool {
	cmd := exec.Command(constants.GitBin, constants.GitStatus, constants.GitStatusShort)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

func stageAllChanges() {
	stageCmd := exec.Command(constants.GitBin, constants.GitAdd, constants.GitAddAll)
	stageCmd.Stderr = os.Stderr
	if err := stageCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not stage changes: %v\n", err)
	}
}

// commitGoModChanges stages and commits all changes.
func commitGoModChanges(oldPath, newPath string, fileCount int) {
	stageAllChanges()
	msg := fmt.Sprintf(constants.GoModCommitMsgFmt, oldPath, newPath, fileCount)
	commitCmd := exec.Command(constants.GitBin, constants.GitCommit, constants.GitCommitMsg, msg)
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		cliexit.HandleGeneralError(apperror.WrapSimple(err, constants.ErrGoModCommitFailed))
	}
}

// mergeGoModBranch checks out the original branch and merges the feature branch.
func mergeGoModBranch(originalBranch, featureBranch, newPath string) {
	checkoutBranch(originalBranch)

	mergeMsg := fmt.Sprintf(constants.GoModMergeMsgFmt, newPath)
	cmd := exec.Command(constants.GitBin, constants.GitMerge, constants.GitMergeNoFF, constants.GitCommitMsg, mergeMsg, featureBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cliexit.HandleGeneralError(apperror.NewSimple(constants.ErrGoModMergeConflict, "E9000"))
	}
}

// goModTidy runs go mod tidy in the current directory.
func goModTidy() error {
	cmd := exec.Command(constants.GoBin, "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
