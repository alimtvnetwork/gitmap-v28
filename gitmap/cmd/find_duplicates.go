package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runFindDuplicates dispatches duplicate discovery for a given platform or across all platforms.
func runFindDuplicates(platform string, args []string) error {
	resolved := resolveDuplicatePlatform(platform, args)
	switch resolved {
	case "agy", "ag", "antigravity":
		return runFindDuplicatesAgy()
	case "vscode", "vsc":
		return runFindDuplicatesVSCode()
	case "chrome", "chromeprofile", "chrome-profile":
		return runFindDuplicatesChrome()
	case "git", "repo", "clone":
		return runFindDuplicatesGit()
	case "all", "":
		return runFindDuplicatesAll()
	default:
		return printUnknownDupPlatform(resolved)
	}
}

func resolveDuplicatePlatform(platform string, args []string) string {
	if platform != "" {
		return strings.ToLower(strings.TrimSpace(platform))
	}
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return strings.ToLower(strings.TrimSpace(args[0]))
	}
	return "all"
}

func printUnknownDupPlatform(platform string) error {
	fmt.Printf(constants.ColorRed+"Unknown platform '%s' for find-duplicates."+constants.ColorReset+"\n", platform)
	fmt.Println("Available platforms: agy, vscode, chrome, git (or omit platform to check all)")
	return nil
}

func runFindDuplicatesAll() error {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "═════════════════════════════════════════════════════════════════════" + constants.ColorReset)
	fmt.Println("  " + constants.ColorMagenta + "       Gitmap Cross-Platform Duplicate Project & Repo Auditor        " + constants.ColorReset)
	fmt.Println("  " + constants.ColorMagenta + "═════════════════════════════════════════════════════════════════════" + constants.ColorReset)

	_ = runFindDuplicatesAgy()
	_ = runFindDuplicatesVSCode()
	_ = runFindDuplicatesChrome()
	_ = runFindDuplicatesGit()
	return nil
}
