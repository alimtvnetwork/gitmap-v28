package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func logImportProcessStart(alias, ip string) {
	fmt.Println()
	fmt.Printf("  %s● Importing all settings and data from [%s | %s]%s\n",
		constants.ColorCyan, alias, ip, constants.ColorReset)
}

func applyImportBundleWithDetailedLogs(bundle *GitmapExportBundle, c db.SSHConnection, opts SSHImportAllOptions) error {
	if opts.IsDryRun {
		fmt.Printf("  (dry-run) would import %d macro(s), config, and %d connection(s) from %s\n",
			len(bundle.Macros), len(bundle.Connections), c.Alias)
		return nil
	}

	cfgLog := importConfigWithLog(bundle.Config)
	macroCount, macroNames, mLog := importMacrosWithLog(bundle.Macros)
	connCount, cLog := importConnectionsWithLog(bundle.Connections)
	khCount, khLog := importKnownHostsWithLog(bundle.KnownHosts)

	printStepLog(cfgLog)
	printStepLog(mLog)
	printStepLog(cLog)
	printStepLog(khLog)

	printDetailedImportSummary(c.Alias, c.IPAddress, macroCount, macroNames, connCount, khCount)

	return nil
}

func printStepLog(msg string) {
	if msg != "" {
		fmt.Printf("    %s\n", msg)
	}
}

func importConfigWithLog(cfg map[string]any) string {
	if len(cfg) == 0 {
		return "• Config: (no remote config found)"
	}
	importConfigLocally(cfg)

	return "✓ Config: imported ~/.gitmap/config.json"
}

func importMacrosWithLog(macros []macro.Macro) (int, []string, string) {
	if len(macros) == 0 {
		return 0, nil, "• Macros: (no remote macros found)"
	}
	count := 0
	var names []string
	for _, m := range macros {
		if err := macro.SaveMacro(&m); err == nil {
			count++
			names = append(names, m.Name)
		}
	}
	namesStr := formatExportMacroNames(names)

	return count, names, fmt.Sprintf("✓ Macros: imported %d/%d macro(s) [%s]", count, len(macros), namesStr)
}

func importConnectionsWithLog(conns []db.SSHConnection) (int, string) {
	if len(conns) == 0 {
		return 0, "• SSH Registry: (no remote connections to import)"
	}
	count := importConnectionsLocally(conns)

	return count, fmt.Sprintf("✓ SSH Registry: imported %d connection(s) into SQLite DB", count)
}

func importKnownHostsWithLog(kh string) (int, string) {
	if kh == "" {
		return 0, "• Known Hosts: (no remote host keys to import)"
	}
	count := importKnownHostsLocally(kh)

	return count, fmt.Sprintf("✓ Known Hosts: imported %d host key(s) into ~/.ssh/known_hosts", count)
}

func printDetailedImportSummary(alias, ip string, macros int, names []string, conns, kh int) {
	fmt.Println()
	fmt.Printf("  %s✓ Import completed successfully from [%s | %s]%s\n",
		constants.ColorGreen, alias, ip, constants.ColorReset)
	if len(names) > 0 {
		for _, name := range names {
			fmt.Printf("    • Macro %q ready to run (gitmap macro run %s)\n", name, name)
		}
	}
	fmt.Println()
}
