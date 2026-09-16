// Package cmdagy — agy_help.go renders colorful, aligned help for the agy command suite.
package cmdagy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderAgyHelp(cmd *cobra.Command, args []string) {
	if cmd != nil && cmd != AgyCmd {
		_ = cmd.Usage()

		return
	}
	fmt.Println()
	printAgyHelpBanner()
	printAgyHelpUsage()
	printAgyHelpProjectMgmt()
	printAgyHelpDiagnostics()
	printAgyHelpAutomation()
	printAgyHelpFooter()
}

func printAgyHelpBanner() {
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║             Antigravity CLI Management (gitmap agy)              ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func printAgyHelpUsage() {
	fmt.Printf("  %sUsage:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    %sgitmap agy [command] [flags]%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sagy [command] [flags]%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func printAgyHelpProjectMgmt() {
	fmt.Printf("  %sProject Management:%s\n", constants.ColorYellow, constants.ColorReset)
	printHelpLine("open [path]", "Open Antigravity Desktop IDE on path (default: .)")
	printHelpLine("ls", "List projects in status table grouped by root folder")
	printHelpLine("add <id> <name>", "Add an Antigravity project configuration")
	printHelpLine("rm <id>", "Remove project configuration (files on disk preserved)")
	printHelpLine("update <id>", "Update project updatedAt timestamp to now")
	printHelpLine("scan [path]", "Scan directory and register Antigravity projects")
	printHelpLine("reconcile (recon)", "Reconcile missing projects with active paths")
	printHelpLine("remove-missing", "Remove stale references to missing directories")
	printHelpLine("pin-projects (pins)", "Manage pinned priority projects")
	printHelpLine("group", "Group and categorize projects by tag or folder")
	fmt.Println()
}

func printAgyHelpDiagnostics() {
	fmt.Printf("  %sDiagnostics & Optimization:%s\n", constants.ColorYellow, constants.ColorReset)
	printHelpLine("find-duplicates (fdp)", "Detect projects sharing identical filesystem paths")
	printHelpLine("optimize-projects", "Deduplicate and keep newest project per path")
	printHelpLine("clean-cache (cc)", "Clean runtime cache and orphan project artifacts")
	printHelpLine("remove-empty-convs", "Purge projects with 0 conversation steps")
	printHelpLine("status", "Show Antigravity installation and profile status")
	printHelpLine("stats", "Display Antigravity workspace usage statistics")
	fmt.Println()
}

func printAgyHelpAutomation() {
	fmt.Printf("  %sProtocols & Automation:%s\n", constants.ColorYellow, constants.ColorReset)
	printHelpLine("read-all-projects-with-read-prompts (rprp)", "Discover repos, sync, and broadcast Read Memory")
	printHelpLine("all-projects-read-memory-prompt", "Broadcast Read Memory prompt to active projects")
	printHelpLine("fix-pipeline (fp)", "Diagnose and fix CI/CD pipeline issues")
	printHelpLine("prompt <text>", "Send prompt to active project session")
	printHelpLine("sync", "Sync Antigravity settings and project states")
	printHelpLine("export / import", "Export or import Antigravity projects JSON")
	printHelpLine("plugins", "Inspect, enable, or disable Antigravity plugins")
	fmt.Println()
}

func printHelpLine(cmd, desc string) {
	fmt.Printf("    %s%-32s%s %s%s%s\n",
		constants.ColorGreen, cmd, constants.ColorReset,
		constants.ColorWhite, desc, constants.ColorReset)
}

func printAgyHelpFooter() {
	fmt.Printf("  %sFlags:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    %s-h, --help%s   Show help for agy\n\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %s💡 Tip: Run 'gitmap agy <command> --help' for details on any subcommand.%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}
