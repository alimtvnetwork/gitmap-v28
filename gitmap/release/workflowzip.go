package release

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// buildZipGroupAssets creates archives from persistent zip groups.
func buildZipGroupAssets(opts Options) []string {
	if len(opts.ZipGroups) == 0 {
		return nil
	}

	fmt.Printf(constants.MsgZGProcessing, len(opts.ZipGroups))

	if verbose.IsEnabled() {
		for _, g := range opts.ZipGroups {
			verbose.Get().Log("zip-group: processing group %q", g)
		}
	}

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Cannot open DB for zip groups: %v\n", err)

		return nil
	}
	defer db.Close()

	stagingDir, err := EnsureStagingDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGStagingDir, constants.AssetsStagingDir, err)

		return nil
	}

	archives := BuildZipGroupArchives(db, opts.ZipGroups, stagingDir)

	if verbose.IsEnabled() {
		verbose.Get().Log("zip-group: %d group(s) produced %d archive(s)", len(opts.ZipGroups), len(archives))
	}

	if len(archives) == 0 {
		fmt.Printf(constants.MsgZGNoArchives, len(opts.ZipGroups))
	}

	collectZipChecksums(archives)

	return archives
}

// buildAdHocZipAssets creates archives from ad-hoc -Z paths.
func buildAdHocZipAssets(opts Options) []string {
	if len(opts.ZipItems) == 0 {
		return nil
	}

	if verbose.IsEnabled() {
		bundleLabel := opts.BundleName
		if bundleLabel == "" {
			bundleLabel = "(individual)"
		}
		verbose.Get().Log("ad-hoc-zip: %d item(s), bundle=%s", len(opts.ZipItems), bundleLabel)
		for _, item := range opts.ZipItems {
			verbose.Get().Log("ad-hoc-zip: item %s", item)
		}
	}

	stagingDir, err := EnsureStagingDir()
	if err != nil {
		return nil
	}

	archives := BuildAdHocArchive(opts.ZipItems, opts.BundleName, stagingDir)

	if verbose.IsEnabled() {
		verbose.Get().Log("ad-hoc-zip: produced %d archive(s)", len(archives))
	}

	collectZipChecksums(archives)

	return archives
}

// collectZipChecksums computes SHA-1 for each archive and stores in lastZipChecksums.
func collectZipChecksums(archives []string) {
	if len(archives) == 0 {
		return
	}

	if lastZipChecksums == nil {
		lastZipChecksums = make(map[string]string)
	}

	for _, archivePath := range archives {
		hash, err := sha1File(archivePath)
		if err != nil {
			continue
		}

		lastZipChecksums[filepath.Base(archivePath)] = hash
	}
}
