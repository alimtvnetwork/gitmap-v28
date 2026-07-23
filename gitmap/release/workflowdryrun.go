package release

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// printDryRun shows what would happen without executing.
func printDryRun(v Version, branchName, tag, sourceName string, opts Options) error {
	printDryRunSteps(branchName, tag, sourceName)
	printDryRunGoAssets(v, opts)
	printDryRunZipGroups(opts)
	printDryRunAssets(opts.Assets, opts.Compress, opts.Checksums)
	fmt.Printf(constants.MsgReleaseDryRun, "Switch back to "+sourceName)
	printDryRunMeta(v)
	if !opts.NoCommit {
		fmt.Print(constants.MsgAutoCommitScanning)
	} else {
		fmt.Print(constants.MsgAutoCommitSkipped)
	}
	fmt.Printf(constants.MsgReleaseComplete, v.String())

	return nil
}

// printDryRunGoAssets shows Go cross-compile plan in dry-run mode.
func printDryRunGoAssets(v Version, opts Options) {
	if !opts.Bin {
		return
	}

	if !DetectGoProject() {
		return
	}

	binName := resolveBinName()
	targets, err := ResolveTargets(opts.Targets, opts.ConfigTargets)
	if err != nil {
		return
	}

	names := DescribeTargets(binName, v.String(), targets)
	fmt.Printf(constants.MsgAssetDryRunHeader, len(names))

	for _, name := range names {
		fmt.Printf(constants.MsgAssetDryRunBinary, name)
	}

	fmt.Printf(constants.MsgAssetDryRunUpload, len(names))
}

// printDryRunSteps prints branch/tag/push dry-run lines.
func printDryRunSteps(branchName, tag, sourceName string) {
	fmt.Printf(constants.MsgReleaseDryRun, "Create branch "+branchName+" from "+sourceName)
	fmt.Printf(constants.MsgReleaseDryRun, "Create tag "+tag)
	fmt.Printf(constants.MsgReleaseDryRun, "Push branch and tag to origin")

	body := DetectChangelog()
	if len(body) > 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Use CHANGELOG.md as release body")
	}
	readme := DetectReadme()
	if len(readme) > 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Attach README.md")
	}
}

// printDryRunAssets prints asset attachments in dry-run mode.
func printDryRunAssets(assetsPath string, compress, checksums bool) {
	userAssets := CollectAssets(assetsPath)

	if compress && len(userAssets) > 0 {
		archiveNames := DescribeCompression(userAssets)
		for _, name := range archiveNames {
			fmt.Printf(constants.MsgReleaseDryRun, "Compress → "+name)
		}
	}

	for _, a := range userAssets {
		fmt.Printf(constants.MsgReleaseDryRun, "Attach "+a)
	}

	if checksums && len(userAssets) > 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Generate "+constants.ChecksumsFile+" (SHA256)")
	}
}

// printDryRunMeta prints metadata and latest marker in dry-run mode.
func printDryRunMeta(v Version) {
	fmt.Printf(constants.MsgReleaseDryRun, "Write metadata to "+constants.DefaultReleaseDir+"/"+v.String()+constants.ExtJSON)

	if len(v.PreRelease) == 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Mark "+v.String()+" as latest")
	}
}

// printDryRunZipGroups shows zip group plan in dry-run mode.
func printDryRunZipGroups(opts Options) {
	if len(opts.ZipGroups) > 0 {
		db, err := store.OpenDefault()
		if err == nil {
			defer db.Close()

			DryRunZipGroups(db, opts.ZipGroups)
		}
	}

	DryRunAdHoc(opts.ZipItems, opts.BundleName)
}

// returnToBranch switches back to the original branch after a release.
func returnToBranch(branch string) error {
	if len(branch) == 0 {
		return nil
	}

	err := CheckoutBranch(branch)
	if err != nil {
		return fmt.Errorf("switch back to %s: %w", branch, err)
	}

	fmt.Printf(constants.MsgReleaseSwitchedBack, branch)

	return nil
}
