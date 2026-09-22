// Package cmdagy — agy_rm_help.go renders rich guidance for removing Antigravity projects.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderAgyRmHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║        ANTIGRAVITY PROJECT REMOVAL (gitmap agy rm)       ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %sIMPORTANT: Project files, source repositories, and git history on disk%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  %sare strictly PRESERVED. Only Antigravity workspace metadata is removed.%s\n\n", constants.ColorYellow, constants.ColorReset)

	fmt.Printf("  %sOPERATIONS PERFORMED:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    1. Unregisters project configuration from .gemini/config/projects/\n")
	fmt.Printf("    2. Preserves repository files and directories untouched on disk\n")
	fmt.Printf("    3. Records snapshot into TaskHistory with immediate undo guidance\n\n")

	fmt.Printf("  %sCOMMAND USAGE:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy rm <target...>%s (aliases: %sremove%s, %sdelete%s, %sdel%s)\n", constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy rm --folder <dir>%s (or: %sgitmap agy rm folder <dir>%s)\n\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  %sTARGET SPECIFIERS (COMMA-SEPARATED SUPPORTED):%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • Sequence: %sgitmap agy rm 001%s or %sgitmap agy rm 001,002%s\n", constants.ColorYellow, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • Slug:     %sgitmap agy rm gitmap%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • ID:       %sgitmap agy rm 418bd745%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • Combined: %sgitmap agy rm 001,gitmap,418bd745%s\n\n", constants.ColorYellow, constants.ColorReset)
}
