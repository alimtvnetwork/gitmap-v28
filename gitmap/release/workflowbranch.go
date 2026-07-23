package release

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ExecuteFromBranch runs the release workflow from an existing release branch.
func ExecuteFromBranch(branchName, assetsPath, notes string, isDraft, dryRun, noCommit, yes bool) error {
	version, err := extractVersionFromBranch(branchName)
	if err != nil {
		return err
	}

	err = validateExistingBranch(branchName, version)
	if err != nil {
		return err
	}

	fmt.Printf(constants.MsgReleaseBranchStart, branchName)

	if dryRun {
		return printDryRun(version, branchName, version.String(), branchName, Options{
			Assets: assetsPath, Notes: notes, IsDraft: isDraft, DryRun: true,
		})
	}

	return completeBranchRelease(version, branchName, assetsPath, notes, isDraft, noCommit, yes)
}

// extractVersionFromBranch parses the version from a release branch name.
func extractVersionFromBranch(branchName string) (Version, error) {
	prefix := constants.ReleaseBranchPrefix
	if len(branchName) <= len(prefix) {
		return Version{}, fmt.Errorf(constants.ErrReleaseInvalidVersion, branchName)
	}

	versionStr := branchName[len(prefix):]

	return Parse(versionStr)
}

// validateExistingBranch checks the branch exists and tag doesn't.
func validateExistingBranch(branchName string, v Version) error {
	if BranchExists(branchName) {
		return checkDuplicate(v)
	}

	return fmt.Errorf(constants.ErrReleaseBranchNotFound, branchName)
}

// completeBranchRelease checks out the branch and runs tag/push/release.
func completeBranchRelease(v Version, branchName, assetsPath, notes string, isDraft, noCommit, yes bool) error {
	originalBranch, branchErr := CurrentBranchName()
	if branchErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine current branch: %v\n", branchErr)
	}

	err := CheckoutBranch(branchName)
	if err != nil {
		return fmt.Errorf("checkout branch: %w", err)
	}

	tag := v.String()
	err = CreateTag(tag, constants.ReleaseTagPrefix+tag)
	if err != nil {
		Rollback("", tag, originalBranch)

		return fmt.Errorf("create tag: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	opts := Options{Assets: assetsPath, Notes: notes, IsDraft: isDraft, SkipMeta: true}

	err = pushAndFinalize(v, branchName, tag, branchName, opts)
	if err != nil {
		Rollback("", tag, originalBranch)

		return err
	}

	err = returnToBranch(originalBranch)
	if err != nil {
		return err
	}

	if !noCommit {
		AutoCommit(v.String(), false, yes)
	} else {
		fmt.Print(constants.MsgAutoCommitSkipped)
	}

	return nil
}

// ExecutePending finds all release branches without tags and releases them.
// Also discovers unreleased versions from .gitmap/release/v*.json metadata files.
func ExecutePending(assetsPath, notes string, isDraft, dryRun, noCommit, yes bool) error {
	branches, err := listReleaseBranches()
	if err != nil {
		return fmt.Errorf("could not list release branches: %w", err)
	}

	pending := filterPendingBranches(branches)
	metaPending := discoverMetadataPending(pending)

	total := len(pending) + len(metaPending)
	if total == 0 {
		fmt.Println(constants.MsgReleasePendingNone)

		return nil
	}

	fmt.Printf(constants.MsgReleasePendingFound, total)

	if len(metaPending) > 0 {
		fmt.Printf(constants.MsgPendingMetaFound, len(metaPending))
	}

	err = releasePendingBranches(pending, assetsPath, notes, isDraft, dryRun, noCommit, yes)
	if err != nil {
		return err
	}

	return releasePendingFromMetadata(metaPending, assetsPath, notes, isDraft, dryRun)
}

// releasePendingBranches iterates and releases each pending branch.
func releasePendingBranches(pending []string, assetsPath, notes string, isDraft, dryRun, noCommit, yes bool) error {
	for _, branchName := range pending {
		err := ExecuteFromBranch(branchName, assetsPath, notes, isDraft, dryRun, noCommit, yes)
		if err != nil {
			fmt.Printf(constants.MsgReleasePendingFailed, branchName, err)
			continue
		}
	}

	return nil
}

// listReleaseBranches returns all local branches matching release/v*.
func listReleaseBranches() ([]string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, constants.GitBranchListFlag, constants.ReleaseBranchPrefix+constants.GitTagGlob)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseBranchLines(string(out)), nil
}

// parseBranchLines extracts branch names from git branch output.
func parseBranchLines(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var branches []string

	for _, line := range lines {
		branch := strings.TrimSpace(line)
		branch = strings.TrimPrefix(branch, "* ")
		if len(branch) > 0 {
			branches = append(branches, branch)
		}
	}

	return branches
}

// filterPendingBranches returns branches whose version tag does not exist.
func filterPendingBranches(branches []string) []string {
	var pending []string

	for _, branch := range branches {
		if isPendingBranch(branch) {
			pending = append(pending, branch)
		}
	}

	return pending
}
