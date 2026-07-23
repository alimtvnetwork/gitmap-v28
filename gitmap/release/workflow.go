// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/localdirs"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// Options holds all parameters for a release operation.
// v15: boolean fields use the IsX prefix convention.
type Options struct {
	Version       string
	Assets        string
	Commit        string
	Branch        string
	Bump          string
	Notes         string
	Targets       string
	ConfigTargets []model.ReleaseTarget
	ZipGroups     []string
	ZipItems      []string
	BundleName    string
	IsDraft       bool
	DryRun        bool
	Verbose       bool
	Compress      bool
	Checksums     bool
	Bin           bool
	NoCommit      bool
	SkipMeta      bool
	Yes           bool
}

// Result holds the outcome of a release operation.
type Result struct {
	Version    Version
	BranchName string
	Tag        string
	Commit     string
	Source     string
	Assets     []string
}

// Execute runs the full release workflow.
func Execute(opts Options) error {
	EnsureGitignore()

	// Early check: if no version/bump provided and we're on a release/* branch,
	// extract the version from the branch name and delegate.
	if len(opts.Version) == 0 && len(opts.Bump) == 0 {
		if delegated, delegateErr := tryDelegateFromCurrentBranch(opts); delegated {
			return delegateErr
		}
	}

	version, err := resolveVersion(opts)
	if err != nil {
		return err
	}

	// If version was specified and matches the current release branch, delegate.
	if delegated, delegateErr := tryDelegateFromBranch(version, opts); delegated {
		return delegateErr
	}

	err = checkDuplicate(version)
	if err != nil {
		return err
	}

	sourceRef, sourceName, err := ResolveSourceRef(opts.Commit, opts.Branch)
	if err != nil {
		return err
	}

	if len(opts.Notes) > 0 {
		fmt.Printf(constants.MsgReleaseNotes, opts.Notes)
	}

	return performRelease(version, sourceRef, sourceName, opts)
}

// tryDelegateFromBranch checks if we should delegate to ExecuteFromBranch.
// This handles the case where the user is already on a release/* branch
// and the tag hasn't been created yet (pending release).
func tryDelegateFromBranch(v Version, opts Options) (bool, error) {
	currentBranch, err := CurrentBranchName()
	if err != nil {
		return false, nil
	}

	branchName := constants.ReleaseBranchPrefix + v.String()

	// Check if we're on this release branch (or any release/* branch matching the version).
	if currentBranch != branchName {
		return false, nil
	}

	// Branch exists but tag doesn't — this is a pending release.
	tagExists := TagExistsLocally(v.String()) || TagExistsRemote(v.String())
	if tagExists {
		return false, nil
	}

	fmt.Printf(constants.MsgReleaseBranchPending, branchName)

	return true, ExecuteFromBranch(branchName, opts.Assets, opts.Notes, opts.IsDraft, opts.DryRun, opts.NoCommit, opts.Yes)
}

// tryDelegateFromCurrentBranch checks if we're on a release/* branch
// with no tag when no version was explicitly provided.
func tryDelegateFromCurrentBranch(opts Options) (bool, error) {
	currentBranch, err := CurrentBranchName()
	if err != nil {
		return false, nil
	}

	if !strings.HasPrefix(currentBranch, constants.ReleaseBranchPrefix) {
		return false, nil
	}

	v, err := extractVersionFromBranch(currentBranch)
	if err != nil {
		return false, nil
	}

	tagExists := TagExistsLocally(v.String()) || TagExistsRemote(v.String())
	if tagExists {
		return false, nil
	}

	fmt.Printf(constants.MsgReleaseBranchPending, currentBranch)

	return true, ExecuteFromBranch(currentBranch, opts.Assets, opts.Notes, opts.IsDraft, opts.DryRun, opts.NoCommit, opts.Yes)
}

// resolveVersion determines the version from CLI args, bump, or file.
func resolveVersion(opts Options) (Version, error) {
	if len(opts.Version) > 0 {
		v, err := Parse(opts.Version)
		if err != nil {
			return v, err
		}
		if verbose.IsEnabled() {
			verbose.Get().Log("version: resolved from CLI argument: %s", v.String())
		}
		return v, nil
	}
	if len(opts.Bump) > 0 {
		v, err := resolveBump(opts.Bump)
		if err != nil {
			return v, err
		}
		if verbose.IsEnabled() {
			verbose.Get().Log("version: resolved via --bump %s: %s", opts.Bump, v.String())
		}
		return v, nil
	}

	v, err := resolveFromFile()
	if err != nil {
		return v, err
	}
	if verbose.IsEnabled() {
		verbose.Get().Log("version: resolved from %s: %s", constants.DefaultVersionFile, v.String())
	}
	return v, nil
}

// performRelease executes the branch/tag/push/release steps.
func performRelease(v Version, sourceRef, sourceName string, opts Options) error {
	branchName := constants.ReleaseBranchPrefix + v.String()
	tag := v.String()

	originalBranch, branchErr := CurrentBranchName()
	if branchErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine current branch: %v\n", branchErr)
	}

	fmt.Printf(constants.MsgReleaseStart, tag)

	if opts.DryRun {
		return printDryRun(v, branchName, tag, sourceName, opts)
	}

	// Step 1: Create the release branch, tag, push, and finalize assets.
	err := executeSteps(v, branchName, tag, sourceRef, sourceName, opts)
	if err != nil {
		Rollback(branchName, tag, originalBranch)

		return err
	}

	// Step 2: Return to the original branch.
	err = returnToBranch(originalBranch)
	if err != nil {
		return err
	}

	// Step 3: Re-run legacy directory migration on the original branch.
	// Older branches may still track .release/, which checkout restores.
	localdirs.MigrateLegacyDirs()

	// Step 4: Write metadata JSON on the original branch (picked up by auto-commit).
	if !opts.SkipMeta {
		err = writeMetadata(v, branchName, tag, sourceName, nil, opts)
		if err != nil {
			return err
		}
	}

	// Step 5: Auto-commit the release metadata files.
	if !opts.NoCommit {
		AutoCommit(v.String(), false, opts.Yes)
	} else {
		fmt.Print(constants.MsgAutoCommitSkipped)
	}

	return nil
}

// executeSteps runs each release step in sequence.
func executeSteps(v Version, branchName, tag, sourceRef, sourceName string, opts Options) error {
	err := CreateBranch(branchName, sourceRef)
	if err != nil {
		return fmt.Errorf("create branch: %w", err)
	}
	fmt.Printf(constants.MsgReleaseBranch, branchName)

	err = CreateTag(tag, resolveTagMessage(tag, opts))
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	return pushAndFinalize(v, branchName, tag, sourceName, opts)
}

// resolveTagMessage returns the tag annotation message, using notes if provided.
func resolveTagMessage(tag string, opts Options) string {
	if len(opts.Notes) > 0 {
		return opts.Notes
	}

	return constants.ReleaseTagPrefix + tag
}
