// Package cmdzsh provides Oh-My-Zsh themes catalog and terminal rendering.
package cmdzsh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// SupportedThemes lists the 41 supported Oh-My-Zsh themes from Adem's automation suite.
var SupportedThemes = []string{
	"amuse",
	"robbyrussell",
	"af-magic",
	"afowler",
	"agnoster",
	"alanpeabody",
	"apple",
	"arrow",
	"aussiegeek",
	"avit",
	"awesomepanda",
	"bira",
	"blinks",
	"bureau",
	"candy",
	"clean",
	"cloud",
	"crcandy",
	"crunch",
	"cypher",
	"dallas",
	"darkblood",
	"daveverwer",
	"dpoggi",
	"dst",
	"duellj",
	"edvardm",
	"fino-time",
	"fino",
	"fishy",
	"fletcherm",
	"fox",
	"frisk",
	"frontcube",
	"funky",
	"fwalch",
	"gnzh",
	"godzilla",
	"half-life",
	"intheloop",
	"itchy",
}

// IsSupportedTheme returns whether the specified theme is present in the 41 supported themes.
func IsSupportedTheme(name string) bool {
	normalized := strings.ToLower(strings.TrimSpace(name))
	for _, theme := range SupportedThemes {
		isMatch := theme == normalized || (theme == "fletcherm" && normalized == "flettcherm")
		if isMatch {
			return true
		}
	}

	return false
}

// FormatThemesCatalog formats the 41 supported Oh-My-Zsh themes into a readable string.
func FormatThemesCatalog() string {
	var sb strings.Builder
	sb.WriteString("=== Supported Oh-My-Zsh Themes (41) ===" + constants.NewLineUnix)
	sb.WriteString("Wiki: https://github.com/ohmyzsh/ohmyzsh/wiki/Themes" + constants.NewLineUnix)
	sb.WriteString(constants.NewLineUnix)
	for i, theme := range SupportedThemes {
		sb.WriteString(fmt.Sprintf("  %2d. %s%s", i+1, theme, constants.NewLineUnix))
	}

	return sb.String()
}

// DisplayThemesCatalog prints the formatted theme catalog to standard output.
func DisplayThemesCatalog() error {
	report := FormatThemesCatalog()
	fmt.Print(report)

	return nil
}

func dispatchThemes(args []string) error {
	return DisplayThemesCatalog()
}
