// Package cmdchrome provides CLI commands and utilities for Chrome profiles,
// extensions, bookmarks, flags, and observe sessions.
package cmdchrome

var (
	// InstallChromeLinuxFn hook for linux chrome deb installation.
	InstallChromeLinuxFn func(isDryRun bool) error

	// InstallToolFn hook for platform-generic tool installation.
	InstallToolFn func(tool string, isDryRun bool)

	// RunFindDuplicatesFn hook for duplicates finding.
	RunFindDuplicatesFn func(category string, args []string) error

	// CheckHelpFn hook for help text verification.
	CheckHelpFn func(command string, args []string)
)

// RunChrome executes the root chrome command logic.
func RunChrome(args []string) error {
	return runChrome(args)
}

// FindChromeBinaryPath returns path to local chrome binary.
func FindChromeBinaryPath() (string, error) {
	return findChromeBinaryPath()
}

// ReadChromeBackup extracts a tar.gz archive into dstRoot.
func ReadChromeBackup(src, dstRoot string) (int, error) {
	return readChromeBackup(src, dstRoot)
}
