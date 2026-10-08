package cmdchromeprofile

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func importSingleSnapshotWithStepLogging(srcFile, explicitTarget string) error {
	if isDelegatedSnapshotFile(srcFile) {
		return handleDelegatedSnapshot(srcFile, explicitTarget)
	}

	meta, err := readSnapshotMetadata(srcFile)
	if err != nil {
		return err
	}

	exp := meta.Export
	logSnapshotInspectionStep(srcFile, meta)
	dest := resolveImportDestination(exp, explicitTarget, true)
	if err := os.MkdirAll(dest.Path, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", dest.Path, err)
	}

	if err := restoreSnapshotContent(exp, dest); err != nil {
		return err
	}

	registerSnapshotProfile(exp, dest)
	checkChromeRunningAdvisory()
	fmt.Printf("  \033[1;92m✓ Successfully imported\033[0m %s → %s (%q)\n\n", srcFile, dest.Dir, dest.DisplayName)

	return nil
}

func isDelegatedSnapshotFile(srcFile string) bool {
	return strings.HasSuffix(strings.ToLower(srcFile), constants.ExtZIP) || isDirectoryPath(srcFile)
}

func handleDelegatedSnapshot(srcFile, explicitTarget string) error {
	if strings.HasSuffix(strings.ToLower(srcFile), constants.ExtZIP) {
		return importZipSnapshotWithStepLogging(srcFile, explicitTarget)
	}

	return importDirectorySnapshotWithStepLogging(srcFile, explicitTarget)
}

func logSnapshotInspectionStep(srcFile string, meta *snapshotMetadata) {
	exp := meta.Export
	fmt.Printf("  \033[1;94m[Step 1/5]\033[0m Inspecting snapshot: %s\n", srcFile)
	fmt.Printf("        → Name: %q | Display: %q | Email: %q | Bookmarks: %d | Extensions: %d\n",
		exp.Name, exp.DisplayName, exp.Email, meta.BookmarksCount, meta.ExtensionsCount)
}

func restoreSnapshotContent(exp *chromeExport, dest importDestination) error {
	fmt.Printf("  \033[1;94m[Step 3/5]\033[0m Restoring Bookmarks, Preferences & Auth into %q...\n", dest.Dir)
	if err := restoreExportBaseFiles(exp, dest.Path); err != nil {
		return err
	}

	if err := patchImportedChromeProfilePreferences(dest.Path, dest.DisplayName); err != nil {
		fmt.Fprintf(os.Stderr, "        \033[1;93m⚠\033[0m Preferences patch notice: %v\n", err)
	}

	restoreExportAuthPayloads(exp, dest.Path)
	fmt.Printf("  \033[1;94m[Step 4/5]\033[0m Staging extensions (%d pending hints)...\n", len(exp.ExtensionIDs))

	return writePendingExtensions(dest.Path, exp.ExtensionIDs)
}

func registerSnapshotProfile(exp *chromeExport, dest importDestination) {
	fmt.Printf("  \033[1;94m[Step 5/5]\033[0m Registering %q in Chrome Local State...\n", dest.Dir)
	if err := registerChromeProfileWithFullSchemaAndGAIA(
		dest.Dir,
		dest.DisplayName,
		dest.Email,
		exp.GaiaID,
		exp.GaiaName,
		exp.GaiaGivenName,
	); err != nil {
		fmt.Fprintf(os.Stderr, "        \033[1;93m⚠\033[0m Warning: Local State registration notice: %v\n", err)
	}
}

func checkChromeRunningAdvisory() {
	isRunning, err := isChromeRunning(runtime.GOOS)
	if err != nil || !isRunning {
		return
	}

	fmt.Println()
	fmt.Println("        \033[1;93m⚠ Warning: Google Chrome is currently running!\033[0m")
	fmt.Println("          Chrome caches Local State in memory and will overwrite disk changes on exit.")
	fmt.Println("          To display this new profile in Chrome's Profile Picker:")
	fmt.Println("          1. Close Google Chrome completely.")
	fmt.Println("          2. Run: gitmap chrome profile reconcile")
	fmt.Println()
}
