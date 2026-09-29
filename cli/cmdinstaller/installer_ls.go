// Package cmd — installer_ls.go defines the installer list CLI subcommand.
package cmdinstaller

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var isInstallerTreeMode bool

// installerLsCmd represents the 'gitmap installer ls' subcommand.
var installerLsCmd = &cobra.Command{
	Use:   "ls [os_filter]",
	Short: "List registered installer scripts",
	Long:  "Displays a tabular view of all installer scripts, optionally filtered by target OS.",
	RunE:  runInstallerLs,
}

// InstallerLsCmd is an exported alias for installerLsCmd.
var InstallerLsCmd = installerLsCmd

func init() {
	if installerLsCmd != nil {
		installerLsCmd.Flags().BoolVarP(&isInstallerTreeMode, "tree", "t", false, "Display installer history in tree view")
	}

	if installerCmd != nil {
		installerCmd.AddCommand(installerLsCmd)
	}
}

// printInstallerTableHeader outputs the table headers and divider for installer ls.
func printInstallerTableHeader() {
	fmt.Printf("%-20s %-15s %-10s %-10s %-12s %s\n", "NAME", "SLUG", "OS", "VERSION", "PINNED", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 88))
}

// isSkippableInstallerRow checks if target OS matches the provided filter.
func isSkippableInstallerRow(targetOS, filter string) bool {
	if filter == "" || filter == "all" {
		return false
	}

	return !strings.EqualFold(targetOS, filter) && !strings.EqualFold(targetOS, "all")
}

func resolvePinnedVersion(slug, name string, pins map[string]string) (string, bool) {
	if p, ok := pins[slug]; ok && p != "" {
		return p, true
	}
	if p, ok := pins[name]; ok && p != "" {
		return p, true
	}
	return "-", false
}

// printInstallerRows iterates through scripts and prints matching table rows.
func printInstallerRows(scriptList []model.InstallerScript, osFilter string, pins map[string]string) (int, int) {
	filterKey := strings.ToLower(strings.TrimSpace(osFilter))
	matchedCount := 0
	pinnedCount := 0
	for _, scriptRecord := range scriptList {
		if isSkippableInstallerRow(scriptRecord.TargetOS, filterKey) {
			continue
		}

		pinned, isPinned := resolvePinnedVersion(strings.ToLower(scriptRecord.Slug), strings.ToLower(scriptRecord.Name), pins)
		if isPinned {
			pinnedCount++
		}

		fmt.Printf("%-20s %-15s %-10s %-10s %-12s %s\n",
			scriptRecord.Name, scriptRecord.Slug, scriptRecord.TargetOS, scriptRecord.Version, pinned, scriptRecord.Description)
		matchedCount++
	}

	return matchedCount, pinnedCount
}

func printPinnedSummary(total, pinnedCount int, pins map[string]string) {
	fmt.Println()
	fmt.Printf("Summary: %d registered installer(s), %d pinned\n", total, pinnedCount)
	if len(pins) == 0 {
		return
	}
	fmt.Print("Pinned Versions: ")
	var items []string
	for slug, ver := range pins {
		items = append(items, fmt.Sprintf("%s@%s", slug, ver))
	}
	fmt.Println(strings.Join(items, ", "))
}

// executeInstallerLs retrieves and displays installer records formatted as a table.
func executeInstallerLs(ctx context.Context, dbInstance *store.DB, osFilter string) error {
	if dbInstance == nil {
		return apperror.New("executeInstallerLs", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "db cannot be nil"})
	}

	scriptList, errList := dbInstance.ListInstallers()
	if errList != nil {
		return errList
	}

	pins, _ := dbInstance.ListPinnedVersions()
	printInstallerTableHeader()
	matchedCount, pinnedCount := printInstallerRows(scriptList, osFilter, pins)
	if matchedCount == 0 {
		fmt.Println("No installer scripts found matching criteria.")
	}

	printPinnedSummary(len(scriptList), pinnedCount, pins)
	return nil
}

// openAndMigrateInstallerDB opens default DB connection and ensures schema migration.
func openAndMigrateInstallerDB() (*store.DB, error) {
	dbInstance, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "openAndMigrateInstallerDB", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"

		return nil, appErr
	}

	if errMigrate := dbInstance.MigrateInstallers(); errMigrate != nil {
		dbInstance.Close()
		appErr := apperror.Wrap(errMigrate, "openAndMigrateInstallerDB", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"

		return nil, appErr
	}

	return dbInstance, nil
}

// extractOsFilter extracts the optional OS filter argument.
func extractOsFilter(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	return ""
}

// runInstallerLs coordinates arguments and dispatches to executeInstallerLs or history tree.
func runInstallerLs(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerLs", "E_INSTALLER_NIL_COMMAND", map[string]any{"error": "command is nil"})
	}

	dbInstance, errDB := openAndMigrateInstallerDB()
	if errDB != nil {
		return errDB
	}

	defer dbInstance.Close()
	if hasTree, _ := cmd.Flags().GetBool("tree"); hasTree || isInstallerTreeMode {
		printInstallerHistoryTree(dbInstance)

		return nil
	}

	return executeInstallerLs(cmd.Context(), dbInstance, extractOsFilter(args))
}
