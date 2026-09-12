package cmdchromeprofile

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RunProfileCopy implements `gitmap chrome-profile-copy`.
func RunProfileCopy(args []string) error {
	return runChromeProfileCopy(args)
}

// RunProfileExport implements `gitmap chrome-profile-export`.
func RunProfileExport(args []string) error {
	return runChromeProfileExport(args)
}

// RunProfileImport implements `gitmap chrome-profile-import`.
func RunProfileImport(args []string) error {
	return runChromeProfileImport(args)
}

// RunProfileImportCheck implements `gitmap chrome-profile-import check`.
func RunProfileImportCheck(args []string) error {
	return runChromeProfileImportCheck(args)
}

// RunProfileList implements `gitmap chrome-profile-list`.
func RunProfileList(args []string) error {
	return runChromeProfileList(args)
}

// RunProfileDelete implements `gitmap chrome-profile-delete`.
func RunProfileDelete(args []string) error {
	return runChromeProfileDelete(args)
}

// RunProfileClear implements `gitmap chrome-profile clear`.
func RunProfileClear(args []string) error {
	return runChromeProfileClear(args)
}

// RunProfileOptimize implements `gitmap chrome-profile optimize-projects`.
func RunProfileOptimize(args []string) error {
	return runChromeProfileOptimize(args)
}

// RunProfileMerge implements `gitmap chrome-profile-merge`.
func RunProfileMerge(args []string) error {
	return runChromeProfileMerge(args)
}

// RunProfileReconcile implements `gitmap chrome-profile reconcile`.
func RunProfileReconcile(args []string) error {
	return runChromeProfileReconcile(args)
}

// RunProfileUndo handles undo for Chrome profile mutations.
func RunProfileUndo(profileName string) error {
	return runChromeProfileUndo(profileName)
}

// RunProfileRedo handles redo for Chrome profile mutations.
func RunProfileRedo(profileName string) error {
	return runChromeProfileRedo(profileName)
}

// RunProfileGroupDispatch dispatches Chrome profile group commands.
func RunProfileGroupDispatch(args []string) error {
	return runChromeGroupDispatch(args)
}

// RunCopyAll implements `gitmap cpa` / `gitmap chrome-profile-copy-all`.
func RunCopyAll(args []string) error {
	return runChromeCopyAll(args)
}

// RunExportAll implements `gitmap cpea` / `gitmap chrome-profile-export-all`.
func RunExportAll(args []string) error {
	return runChromeExportAll(args)
}

// RunImportAll implements `gitmap cpia` / `gitmap chrome-profile-import-all`.
func RunImportAll(args []string) error {
	return runChromeImportAll(args)
}

// RunFindDuplicates scans and reports duplicate Chrome profiles.
func RunFindDuplicates() error {
	return runFindDuplicatesChrome()
}

// UserDataDir returns the platform-specific Chrome User Data root.
func UserDataDir() string {
	return chromeUserDataDir()
}

// ProfilePath joins User Data root with a named profile directory.
func ProfilePath(name string) string {
	return chromeProfilePath(name)
}

// ProfilePathExists reports whether a profile directory exists on disk.
func ProfilePathExists(path string) bool {
	return chromeProfilePathExists(path)
}

// AvailableProfileNames returns profile-shaped directories in User Data root.
func AvailableProfileNames() []string {
	return availableChromeProfileNames()
}

// ProfileResolution represents a resolved Chrome profile directory and metadata.
type ProfileResolution struct {
	Input       string
	Path        string
	Dir         string
	DisplayName string
}

// ResolveProfile resolves a profile name or directory.
func ResolveProfile(name string) (ProfileResolution, bool) {
	res, ok := resolveChromeProfile(name)

	return ProfileResolution{
		Input:       res.Input,
		Path:        res.Path,
		Dir:         res.Dir,
		DisplayName: res.DisplayName,
	}, ok
}

// ResolveProfileDir resolves a profile directory name.
func ResolveProfileDir(name string) (string, bool) {
	return resolveChromeProfileDir(name)
}

// ProfileDisplayName returns the display name for a profile directory.
func ProfileDisplayName(dirName string) string {
	return chromeProfileDisplayName(dirName)
}

// ProfileSummary returns a formatted summary of a profile.
func ProfileSummary(p ProfileResolution) string {
	return chromeProfileSummary(chromeProfileResolution{
		Input:       p.Input,
		Path:        p.Path,
		Dir:         p.Dir,
		DisplayName: p.DisplayName,
	})
}

// PrintAvailableProfilesWithDisplay prints all detected Chrome profiles.
func PrintAvailableProfilesWithDisplay() {
	printAvailableChromeProfilesWithDisplay()
}

// SnapshotProfile creates a timestamped snapshot of a Chrome profile before mutation.
func SnapshotProfile(profileDir string, tag string) (string, *apperror.AppError) {
	return snapshotChromeProfile(profileDir, tag)
}

// CopyEntry copies a single file or directory tree, returning the count of files copied.
func CopyEntry(src, dst string) (int, error) {
	return copyEntry(src, dst)
}

// IsChromeRunning reports whether a Chrome process is currently running on the given OS.
func IsChromeRunning(goos string) (bool, error) {
	return isChromeRunning(goos)
}
