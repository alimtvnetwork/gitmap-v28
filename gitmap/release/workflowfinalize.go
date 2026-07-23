package release

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// LastMeta holds the most recent release metadata after Execute completes.
var LastMeta *ReleaseMeta

// lastZipChecksums stores SHA-1 hashes of zip group archives built during the
// current release, keyed by archive filename. Populated by buildZipGroupAssets
// and buildAdHocZipAssets, consumed by buildReleaseMeta.
var lastZipChecksums map[string]string

// pushAndFinalize pushes to remote and writes metadata.
func pushAndFinalize(v Version, branchName, tag, _ string, opts Options) error {
	lastZipChecksums = nil

	err := PushBranchAndTag(branchName, tag)
	if err != nil {
		return fmt.Errorf(constants.ErrReleasePushFailed, err)
	}
	fmt.Print(constants.MsgReleasePushed)

	assets := CollectAssets(opts.Assets)

	// Cross-compile Go binaries if applicable.
	goAssets := buildGoAssetsIfApplicable(v, opts)
	assets = append(assets, goAssets...)

	// Build zip group archives (persistent groups from DB).
	zipGroupAssets := buildZipGroupAssets(opts)
	assets = append(assets, zipGroupAssets...)

	// Build ad-hoc zip archives (-Z / --bundle).
	adHocAssets := buildAdHocZipAssets(opts)
	assets = append(assets, adHocAssets...)

	// Bundle docs-site for help-dashboard command, and bake per-version
	// snapshots of the release-version installer scripts (spec 105).
	// The release-version.{sh,ps1} snapshots hard-code REPO=alimtvnetwork/
	// gitmap-v28 + BINARY_NAME=gitmap, so they're only valid for the
	// gitmap source repo. For any other repo, attaching them would mislead
	// users into downloading gitmap binaries when they wanted the host
	// project. Gate by the same ShouldPrintInstallHint check used for the
	// release-body install snippet (spec/02-app-issues/27).
	if stagingDir, stagingErr := EnsureStagingDir(); stagingErr == nil {
		if docsSiteAsset := buildDocsSiteAsset(stagingDir); len(docsSiteAsset) > 0 {
			assets = append(assets, docsSiteAsset)
		}

		if ShouldPrintInstallHint(getRemoteURL()) {
			snapshotAssets := buildReleaseVersionSnapshots(v.String(), stagingDir)
			assets = append(assets, snapshotAssets...)
		}
	}

	if opts.Compress && len(assets) > 0 {
		compressed, compErr := CompressAssets(assets)
		if compErr == nil && len(compressed) > 0 {
			for _, a := range compressed {
				fmt.Printf(constants.MsgCompressArchive, filepath.Base(a), filepath.Base(a))
			}

			assets = compressed
		}
	}

	if opts.Checksums && len(assets) > 0 {
		checksumPath, csErr := GenerateChecksums(assets)
		if csErr == nil && len(checksumPath) > 0 {
			fmt.Printf(constants.MsgChecksumGenerated, constants.ChecksumsFile)
			assets = append(assets, checksumPath)
		}
	}

	for _, a := range assets {
		fmt.Printf(constants.MsgReleaseAttach, a)
	}

	// Upload to GitHub if token is available.
	uploadToGitHub(v, assets, opts)

	fmt.Printf(constants.MsgReleaseComplete, v.String())
	printInstallHint(v)

	return nil
}

// writeMetadata persists release info and updates latest.
func writeMetadata(v Version, branchName, tag, sourceName string, assets []string, opts Options) error {
	commit, commitErr := CurrentCommitSHA()
	if commitErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine current commit SHA: %v\n", commitErr)
	}
	meta := buildReleaseMeta(v, branchName, tag, sourceName, commit, assets, opts)

	metaPath := constants.DefaultReleaseDir + "/" + v.String() + constants.ExtJSON

	if verbose.IsEnabled() {
		verbose.Get().Log("metadata: writing %s", metaPath)
	}

	err := WriteReleaseMeta(meta)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseMetaWrite, metaPath, err)
	}
	fmt.Printf(constants.MsgReleaseMeta, metaPath)

	LastMeta = &meta

	return updateLatestIfStable(v)
}

// buildReleaseMeta constructs the metadata struct for a release.
func buildReleaseMeta(v Version, branchName, tag, sourceName, commit string, assets []string, opts Options) ReleaseMeta {
	assetPaths := make([]string, len(assets))
	copy(assetPaths, assets)

	zipGroups := collectZipGroupNames(opts)

	var checksums map[string]string
	if len(lastZipChecksums) > 0 {
		checksums = make(map[string]string, len(lastZipChecksums))
		for k, v := range lastZipChecksums {
			checksums[k] = v
		}
	}

	return ReleaseMeta{
		Version:           v.CoreString(),
		Branch:            branchName,
		SourceBranch:      sourceName,
		Commit:            commit,
		Tag:               tag,
		Assets:            assetPaths,
		Changelog:         loadChangelogNotes(v.String()),
		Notes:             opts.Notes,
		ZipGroups:         zipGroups,
		ZipGroupChecksums: checksums,
		IsDraft:           opts.IsDraft,
		IsPreRelease:      v.IsPreRelease(),
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
		IsLatest:          false,
	}
}

// collectZipGroupNames merges persistent group names and ad-hoc bundle name.
func collectZipGroupNames(opts Options) []string {
	var names []string
	names = append(names, opts.ZipGroups...)

	if len(opts.BundleName) > 0 {
		names = append(names, opts.BundleName)
	}

	if len(names) == 0 {
		return nil
	}

	return names
}

// loadChangelogNotes reads changelog notes for a version, returning nil on error.
func loadChangelogNotes(version string) []string {
	entries, err := ReadChangelog()
	if err != nil {
		return nil
	}

	entry, found := FindChangelogEntry(entries, version)
	if found {
		return entry.Notes
	}

	return nil
}

// updateLatestIfStable marks the release as latest if stable.
func updateLatestIfStable(v Version) error {
	if v.IsPreRelease() {
		if verbose.IsEnabled() {
			verbose.Get().Log("metadata: skipping latest.json (pre-release %s)", v.String())
		}
		fmt.Printf(constants.MsgReleaseComplete, v.String())
		printInstallHint(v)

		return nil
	}

	if LastMeta != nil {
		LastMeta.IsLatest = true
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("metadata: updating latest.json to %s", v.String())
	}

	err := WriteLatest(v)
	if err != nil {
		return err
	}

	fmt.Printf(constants.MsgReleaseLatest, v.String())
	fmt.Printf(constants.MsgReleaseComplete, v.String())
	printInstallHint(v)

	return nil
}
