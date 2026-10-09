package cmdapps

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// runAppsDispatch routes the "apps" subcommand.
func RunAppsDispatch(args []string) error {
	if len(args) == 0 {
		return printAppsHelp()
	}

	subcommand := strings.ToLower(args[0])
	tailArgs := args[1:]

	switch subcommand {
	case "list", "ls":
		return runAppsList(tailArgs)
	case "uninstall", "rm", "remove", "purge":
		return RunAppsUninstall(tailArgs)
	case "help", "-h", "--help":
		return printAppsHelp()
	default:
		// If first argument is not a subcommand, check if it's a flag for list (e.g. gitmap apps --json)
		if strings.HasPrefix(subcommand, "-") {
			return runAppsList(args)
		}
		return fmt.Errorf("unknown apps subcommand: %s. Run 'gitmap apps help' for usage", subcommand)
	}
}

func runAppsList(args []string) error {
	fs := flag.NewFlagSet("apps-list", flag.ContinueOnError)
	var isJson bool
	var isSystem bool
	var isUser bool
	var isAll bool
	var filterStr string

	fs.BoolVar(&isJson, "json", false, "Output structured JSON array")
	fs.BoolVar(&isSystem, "system", false, "Filter only system-wide applications")
	fs.BoolVar(&isUser, "user", false, "Filter only user-scoped applications")
	fs.BoolVar(&isAll, "all", false, "Include utilities and NoDisplay desktop entries")
	fs.BoolVar(&isAll, "a", false, "Alias for --all")
	fs.StringVar(&filterStr, "filter", "", "Fuzzy match name, executable, or package identifier")
	fs.StringVar(&filterStr, "f", "", "Alias for --filter")

	if err := fs.Parse(args); err != nil {
		return err
	}

	opts := ListOptions{
		IsJson:    isJson,
		IsSystem:  isSystem,
		IsUser:    isUser,
		IsAll:     isAll,
		FilterStr: filterStr,
	}

	resp, err := ListApps(opts)
	if err != nil && isJson {
		printAppsListErrorJSON(err)
		return nil
	}
	if err != nil {
		return err
	}

	if isJson {
		return printAppsListJSON(resp)
	}

	if len(resp.Data) == 0 {
		fmt.Println("No installed applications found matching criteria.")
		return nil
	}

	columns := []termout.Column{
		{Title: "APP ID", MinWidth: 20, Align: termout.AlignLeft},
		{Title: "NAME", MinWidth: 22, Align: termout.AlignLeft},
		{Title: "VERSION", MinWidth: 12, Align: termout.AlignLeft},
		{Title: "MANAGER", MinWidth: 10, Align: termout.AlignLeft},
		{Title: "SCOPE", MinWidth: 8, Align: termout.AlignLeft},
		{Title: "PACKAGE", MinWidth: 20, Align: termout.AlignLeft},
	}

	rows := make([]termout.Row, 0, len(resp.Data))
	for _, app := range resp.Data {
		version := app.Version
		if version == "" {
			version = "-"
		}
		pkg := app.Package
		if pkg == "" {
			pkg = "-"
		}
		rows = append(rows, termout.Row{
			Cells: []string{
				app.ID,
				app.Name,
				version,
				string(app.Manager),
				string(app.Scope),
				pkg,
			},
		})
	}

	tableCfg := termout.TableConfig{
		Columns:    columns,
		Rows:       rows,
		HasBorders: true,
	}
	termout.PrintTable(tableCfg)

	return nil
}

func RunAppsUninstall(args []string) error {
	fs := flag.NewFlagSet("apps-uninstall", flag.ContinueOnError)
	var isPurge bool
	var isForce bool
	var isDryRun bool
	var isJson bool

	fs.BoolVar(&isPurge, "purge", false, "Completely purge configuration, state, and cache directories")
	fs.BoolVar(&isForce, "force", false, "Bypass interactive confirmation prompt")
	fs.BoolVar(&isForce, "f", false, "Alias for --force")
	fs.BoolVar(&isForce, "yes", false, "Alias for --force")
	fs.BoolVar(&isForce, "y", false, "Alias for --force")
	fs.BoolVar(&isDryRun, "dry-run", false, "Preview uninstallation actions without modifying the filesystem")
	fs.BoolVar(&isDryRun, "n", false, "Alias for --dry-run")
	fs.BoolVar(&isJson, "json", false, "Output structured JSON result")

	if err := fs.Parse(args); err != nil {
		return err
	}

	posArgs := fs.Args()
	if len(posArgs) == 0 {
		return printAppsUninstallMissingArg(isJson)
	}

	target := posArgs[0]
	opts := UninstallOptions{
		IsPurge:  isPurge,
		IsForce:  isForce,
		IsDryRun: isDryRun,
		IsJson:   isJson,
	}

	resp, err := UninstallApp(target, opts)
	if isJson {
		return printAppsUninstallJSON(resp)
	}

	if err != nil {
		return fmt.Errorf("failed to uninstall %s: %w", target, err)
	}

	fmt.Printf("Successfully uninstalled application: %s\n", target)
	if len(resp.RemovedFiles) > 0 {
		fmt.Println("Removed files:")
		for _, f := range resp.RemovedFiles {
			fmt.Printf("  - %s\n", f)
		}
	}
	if len(resp.CachesReset) > 0 {
		fmt.Println("Caches refreshed:")
		for _, c := range resp.CachesReset {
			fmt.Printf("  - %s\n", c)
		}
	}

	return nil
}

func printAppsHelp() error {
	helpText := `gitmap apps - Discover, inspect, and cleanly uninstall desktop and system applications

Usage:
  gitmap apps [command] [flags]

Available Commands:
  list, ls         List all installed applications, desktop launchers, and CLI tools
  uninstall, rm    Surgically remove or purge an application and its desktop entries

Flags:
  --help, -h       Show help for apps command

Examples:
  # List all desktop applications in interactive table format:
  $ gitmap apps list

  # Discover applications matching 'antigravity' in structured JSON:
  $ gitmap apps list --filter antigravity --json

  # Completely purge an obsolete app and its launcher icon:
  $ gitmap apps uninstall antigravity-tools --purge --force

  # Remove an npm global CLI tool on Windows:
  $ gitmap apps uninstall clot --json
`
	fmt.Fprint(os.Stdout, helpText)
	return nil
}

func printAppsListErrorJSON(err error) {
	errResp := AppListResponse{
		Success: false,
		Error:   err.Error(),
	}
	data, _ := json.MarshalIndent(errResp, "", "  ")
	fmt.Println(string(data))
}

func printAppsListJSON(resp AppListResponse) error {
	data, marshalErr := json.MarshalIndent(resp, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	fmt.Println(string(data))
	return nil
}

func printAppsUninstallMissingArg(isJson bool) error {
	if isJson {
		errResp := AppUninstallResponse{
			Success: false,
			Error:   "missing application identifier argument to uninstall",
		}
		data, _ := json.MarshalIndent(errResp, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	return fmt.Errorf("missing application identifier. Usage: gitmap apps uninstall <app-id> [--purge] [--force]")
}

func printAppsUninstallJSON(resp AppUninstallResponse) error {
	data, marshalErr := json.MarshalIndent(resp, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	fmt.Println(string(data))
	return nil
}
