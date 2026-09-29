// Package cmdinstaller — installer_pin.go manages version pinning for installers.
package cmdinstaller

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// installerPinCmd represents the 'gitmap installer pin' subcommand.
var installerPinCmd = &cobra.Command{
	Use:   "pin <name> <version>",
	Short: "Pin an installer or application to a specific version",
	Long:  "Pins the target installer to a specific version, preventing unwanted auto-upgrades.",
	RunE:  runInstallerPin,
}

// installerUnpinCmd represents the 'gitmap installer unpin' subcommand.
var installerUnpinCmd = &cobra.Command{
	Use:   "unpin <name>",
	Short: "Unpin an installer or application version",
	Long:  "Removes version pinning for the target installer, allowing latest upgrades.",
	RunE:  runInstallerUnpin,
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerPinCmd)
		installerCmd.AddCommand(installerUnpinCmd)
	}
}

// runInstallerPin handles the 'gitmap installer pin <name> <version>' command.
func runInstallerPin(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return apperror.NewValidationError("missing required arguments: expected <name> <version>")
	}

	name := strings.TrimSpace(args[0])
	version := strings.TrimSpace(args[1])
	return executeInstallerPin(context.Background(), name, version)
}

// runInstallerUnpin handles the 'gitmap installer unpin <name>' command.
func runInstallerUnpin(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return apperror.NewValidationError("missing required argument: expected <name>")
	}

	name := strings.TrimSpace(args[0])
	return executeInstallerUnpin(context.Background(), name)
}

func executeInstallerPin(ctx context.Context, name, version string) error {
	db, errDB := openAndMigrateInstallerDB()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if errPin := db.PinInstaller(name, version); errPin != nil {
		return errPin
	}

	fmt.Printf("✔ Successfully pinned %s to version %s\n", name, version)
	return nil
}

func executeInstallerUnpin(ctx context.Context, name string) error {
	db, errDB := openAndMigrateInstallerDB()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if errUnpin := db.UnpinInstaller(name); errUnpin != nil {
		return errUnpin
	}

	fmt.Printf("✔ Successfully unpinned %s\n", name)
	return nil
}

// PinVersion is an exported helper for pinning versions from any package (e.g. cmdupdate or cmd).
func PinVersion(name, version string) error {
	return executeInstallerPin(context.Background(), name, version)
}

// UnpinVersion is an exported helper for unpinning versions from any package.
func UnpinVersion(name string) error {
	return executeInstallerUnpin(context.Background(), name)
}

// GetPinnedVersion returns the pinned version for a target name or empty string.
func GetPinnedVersion(name string) (string, error) {
	db, errDB := openAndMigrateInstallerDB()
	if errDB != nil {
		return "", errDB
	}
	defer db.Close()

	return db.GetPinnedVersion(name)
}

// ListPinnedVersions returns all pinned versions.
func ListPinnedVersions() (map[string]string, error) {
	db, errDB := openAndMigrateInstallerDB()
	if errDB != nil {
		return nil, errDB
	}
	defer db.Close()

	return db.ListPinnedVersions()
}
