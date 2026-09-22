// Package cmdagy — agy_rm_rejoin_help.go renders rich guidance for rm-rejoin-read commands.
package cmdagy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var agyRmRejoinHelpCmd = &cobra.Command{
	Use:     "help",
	Aliases: []string{"man", "info"},
	Short:   "Detailed guide on rm-rejoin-read and rm-rejoin-pin-read commands",
	Run: func(cmd *cobra.Command, args []string) {
		renderAgyRmRejoinHelp()
	},
}

func renderAgyRmRejoinHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║     ANTIGRAVITY RM-REJOIN-READ & RM-REJOIN-PIN-READ      ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  These commands perform clean workspace resets for Antigravity projects without\n")
	fmt.Printf("  affecting repository files or .git history on disk.\n\n")

	fmt.Printf("  %sOPERATIONS PERFORMED:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    1. Completely purges stale conversations from .gemini/antigravity/conversations/\n")
	fmt.Printf("    2. Deletes conversation brain artifacts from .gemini/antigravity/brain/\n")
	fmt.Printf("    3. Cleans conversation records from conversation_summaries.db\n")
	fmt.Printf("    4. Removes and rejoins the project configuration in .gemini/config/projects/\n")
	fmt.Printf("    5. If pin-read (rrpr): pins the project into bookmarks if not already pinned\n")
	fmt.Printf("    6. Injects the enhanced Read Memory protocol prompt into the session\n")
	fmt.Printf("    7. Auto-renames initial conversation title with the project slug/name\n")
	fmt.Printf("    8. Snapshots configuration and logs TaskHistory with undo guidance\n\n")

	fmt.Printf("  %sCOMMAND USAGE:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy rm-rejoin-read <target...>%s (alias: %srrr%s)\n", constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy rm-rejoin-pin-read <target...>%s (alias: %srrpr%s, %srrbr%s)\n\n", constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset)

	fmt.Printf("  %sTARGET SPECIFIERS (COMMA-SEPARATED SUPPORTED):%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • Sequence: %sgitmap agy rrr 001%s or %sgitmap agy rrr 001,002%s\n", constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • Slug:     %sgitmap agy rrr gitmap%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • ID:       %sgitmap agy rrr 418bd745%s\n\n", constants.ColorYellow, constants.ColorReset)
}
