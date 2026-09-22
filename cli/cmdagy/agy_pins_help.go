// Package cmdagy — agy_pins_help.go renders rich guidance and documentation for pinned projects.
package cmdagy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var agyPinProjectsHelpCmd = &cobra.Command{
	Use:     "help",
	Aliases: []string{"man", "info"},
	Short:   "Detailed guide and documentation on Antigravity pinned projects",
	Run: func(cmd *cobra.Command, args []string) {
		renderAgyPinsHelp()
	},
}

func renderAgyPinsHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║        ANTIGRAVITY PINNED PROJECTS WORKFLOW GUIDE        ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Pinned projects act as bookmarks and protection guards for frequently accessed\n")
	fmt.Printf("  workspaces. Pinned projects are immune to bulk prune actions like 'agy clear'\n")
	fmt.Printf("  and are listed first with dedicated 3-digit sequence numbers (e.g. 001, 002).\n\n")

	fmt.Printf("  %sCOMMANDS:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins%s (or %sls%s)             List all pinned projects with 3-digit SEQ\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins add <target...>%s      Pin one or more projects (by path, ID, slug)\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins rm <target...>%s       Unpin by 3-digit SEQ (001), ID, slug, or list\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins edit <target> <name>%s Update display name of a pinned project\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins help%s                 Show this comprehensive help guide\n\n", constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  %sTARGET NOTATIONS (COMMA-SEPARATED SUPPORTED):%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • Sequence:  %sgitmap agy pins rm 001%s or %sgitmap agy pins rm 001,002%s\n", constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • Slug:      %sgitmap agy pins rm gitmap%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • ID:        %sgitmap agy pins rm 418bd745%s\n\n", constants.ColorYellow, constants.ColorReset)
}
