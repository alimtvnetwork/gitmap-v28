// Package cmdinstaller provides CLI subcommands for managing, creating, exporting,
// importing, and versioning installation scripts.
package cmdinstaller

import (
	"archive/zip"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	// ResolveProfileTreeFn hook for external profile resolution.
	ResolveProfileTreeFn func(slug string) (any, bool)

	// PrintProfileTreeFn hook for external profile tree rendering.
	PrintProfileTreeFn func(profile any)

	// PrintProfileInstallSummaryFn hook for external profile summary rendering.
	PrintProfileInstallSummaryFn func(slug string)

	// RunInstallAddFn hook for interactive install add command execution.
	RunInstallAddFn func(args []string) error
)

// OpenAndMigrateInstallerDB opens default DB connection and ensures schema migration.
func OpenAndMigrateInstallerDB() (*store.DB, error) {
	return openAndMigrateInstallerDB()
}

// Slugify delegates to internal slugify.
func Slugify(name string) string {
	return slugify(name)
}

// WriteZipEntry writes a JSON-marshaled installer script into a zip archive.
func WriteZipEntry(zw *zip.Writer, script model.InstallerScript) error {
	return writeZipEntry(zw, script)
}

// PrintInstallerTree renders a recursive tree node hierarchy to stdout.
func PrintInstallerTree(root InstallerTreeNode, prefix string, isLast bool) {
	printInstallerTree(root, prefix, isLast)
}

// PrintInstallSummaryHeader displays the summary header with a green checkmark.
func PrintInstallSummaryHeader(slug string) {
	printInstallSummaryHeader(slug)
}

// RootCmd returns the root installer Cobra command.
func RootCmd() *cobra.Command {
	return installerCmd
}

// GetInstallerCmd returns the package installer command.
func GetInstallerCmd() *cobra.Command {
	return installerCmd
}

// GetInstallerAddCmd returns the installer add subcommand.
func GetInstallerAddCmd() *cobra.Command {
	return installerAddCmd
}
