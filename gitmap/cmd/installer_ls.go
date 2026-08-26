// Package cmd — installer_ls.go defines the installer list CLI subcommand.
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// installerLsCmd represents the 'gitmap installer ls' subcommand.
var installerLsCmd = &cobra.Command{
	Use:   "ls [os_filter]",
	Short: "List registered installer scripts",
	Long:  "Displays a tabular view of all installer scripts, optionally filtered by target OS.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerLs(cmd, args)
	},
}

// InstallerLsCmd is an exported alias for installerLsCmd.
var InstallerLsCmd = installerLsCmd

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerLsCmd)
	}
}

// executeInstallerLs retrieves and displays installer records formatted as a table.
func executeInstallerLs(ctx context.Context, db *store.DB, osFilter string) error {
	if db == nil {
		return apperror.New("executeInstallerLs", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "db cannot be nil"})
	}

	list, errList := db.ListInstallers()
	if errList != nil {
		return errList
	}

	filter := strings.ToLower(strings.TrimSpace(osFilter))
	fmt.Printf("%-20s %-15s %-10s %-10s %s\n", "NAME", "SLUG", "OS", "VERSION", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 75))

	count := 0
	for _, item := range list {
		if filter != "" && filter != "all" {
			if !strings.EqualFold(item.TargetOS, filter) && !strings.EqualFold(item.TargetOS, "all") {
				continue
			}
		}
		fmt.Printf("%-20s %-15s %-10s %-10s %s\n",
			item.Name, item.Slug, item.TargetOS, item.Version, item.Description)
		count++
	}

	if count == 0 {
		fmt.Println("No installer scripts found matching criteria.")
	}
	return nil
}

// runInstallerLs coordinates arguments and dispatches to executeInstallerLs.
func runInstallerLs(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerLs", "E_INSTALLER_NIL_COMMAND", map[string]any{"error": "command is nil"})
	}

	osFilter := ""
	if len(args) > 0 {
		osFilter = args[0]
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerLs", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerLs", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeInstallerLs(cmd.Context(), db, osFilter)
}
