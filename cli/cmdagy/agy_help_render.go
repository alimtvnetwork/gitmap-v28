package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderAgyHelp renders the rich two-column help menu for the agy command suite.
func RenderAgyHelp() {
	termhelp.RenderMenu(buildAgyHelpMenu())
}
