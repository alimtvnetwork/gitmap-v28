// Package cmdagy — agy_wpr_help.go displays the help menu and operational synopsis for WPR.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RenderWPRHelp prints the complete two-column boxed terminal UI guide for watch-prompts-running.
func RenderWPRHelp() {
	fmt.Println()
	fmt.Printf("  %s🚀 GitMap Watch Prompts Running (wpr) Suite%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap wpr [subcommand] [flags]")
	fmt.Println("    gitmap agy wpr [subcommand] [flags]")
	fmt.Println("    gitmap watch-prompts-running [subcommand] [flags]")
	fmt.Println()
	fmt.Println("  Subcommands:")
	fmt.Println("    all                      Snapshot all active/queued prompts to Split-DB and summarize")
	fmt.Println("    ls, list                 Inspect watched projects, watch loop state, and intervals")
	fmt.Println("    start [target]           Start watch loop and auto-recovery for target project or all")
	fmt.Println("    disable [target]         Pause monitoring loop without deleting watched projects")
	fmt.Println("    shutdown, sd             Stop monitoring loop, kill IDE process, and shutdown")
	fmt.Println("    restart [target]         Restart IDE process and re-inject recent active prompts")
	fmt.Println("    switch-account, sa [em]  Switch account (optional email), backup prompts, restart IDE")
	fmt.Println("    fast-forward, ff [email] Fast-forward account switch with prompt preservation")
	fmt.Println("    logs, log [target]       View watcher event history and auto-recovery log entries")
	fmt.Println("    status, st               Inspect watchdog state and cached remote machine identities")
	fmt.Println("    remove, rm <target|all>  Remove specific project or all projects from watch list")
	fmt.Println("    deploy <alias|all> <prj> Deploy Split-DB and enqueue WPR task via SSH/SCP stream")
	fmt.Println("    help                     Show this command synopsis and usage guide")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    -t, --time <duration>    Polling/schedule interval (e.g. 2m, 5m, default 2m)")
	fmt.Println("    -p, --prefix <template>  Prefix template prepended to re-injected prompts (default: 'default')")
	fmt.Println("    -s, --suffix <template>  Suffix template appended to re-injected prompts")
	fmt.Println("    --ssh                    Inspect, list, status, logs, or deploy across cluster SSH fleet")
	fmt.Println("    -j, --json               Emit machine-readable JSON output")
	fmt.Println("    -n, --dry-run            Simulate operations without launching IDE or modifying state")
	fmt.Println()
	fmt.Println("  Architecture:")
	fmt.Println("    • Split-DB: Persists to 'watch-prompts/<repo-slug>/sql.db' beside the CLI 'data' folder.")
	fmt.Println("    • Recent Prompts: Old backup records are pruned on every new backup snapshot.")
	fmt.Println("    • Media Preservation: Images and screenshots are extracted and preserved as file paths.")
	fmt.Println("    • Task History: Every AGY action is logged in TaskHistory for undo/redo and audit.")
	fmt.Println("    • Machine Cache: Remote machine identities are cached locally in GitMap home.")
	fmt.Println()
}
