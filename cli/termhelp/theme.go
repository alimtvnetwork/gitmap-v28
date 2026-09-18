package termhelp

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// DefaultSubcommandIndicator is the arrow symbol shown for subcommands.
const DefaultSubcommandIndicator = "➔"

// Theme holds color constants for help rendering.
type Theme struct {
	BannerColor  string
	HeaderColor  string
	CommandColor string
	DescColor    string
	HintColor    string
	TipColor     string
}

// DefaultTheme returns the standard Catppuccin / vivid dark theme.
func DefaultTheme() Theme {
	return Theme{
		BannerColor:  constants.ColorCyan,
		HeaderColor:  constants.ColorYellow,
		CommandColor: constants.ColorGreen,
		DescColor:    constants.ColorWhite,
		HintColor:    constants.ColorCyan,
		TipColor:     constants.ColorCyan,
	}
}
