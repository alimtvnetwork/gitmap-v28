package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runOSThemeCommand(args []string) error {
	if len(args) == 0 {
		return handleThemeStatus()
	}

	subCmd := strings.ToLower(args[0])
	switch subCmd {
	case "help", "-h", "--help":
		printThemeUsage()

		return nil
	case "status", "st", "info":
		return handleThemeStatus()
	case "dark", "d":
		return handleThemeSet(ThemeModeDark)
	case "light", "l":
		return handleThemeSet(ThemeModeLight)
	default:
		printThemeUsage()

		return apperror.NewSimple("unknown theme mode: "+args[0], "E_INVALID_THEME_MODE")
	}
}

func handleThemeStatus() error {
	engine := newPlatformThemeEngine()
	mode, err := engine.GetTheme()
	if err != nil {
		return err
	}

	fmt.Printf("▶ Current Desktop Theme: %s\n", mode)

	return nil
}

func handleThemeSet(mode ThemeModeType) error {
	engine := newPlatformThemeEngine()
	if err := engine.SetTheme(mode); err != nil {
		return err
	}

	fmt.Printf("✔ Desktop theme set to: %s successfully.\n", mode)

	return nil
}

func printThemeUsage() {
	fmt.Println("Usage: gitmap os theme <mode>")
	fmt.Println("       gitmap os theme dark     (Set dark mode appearance)")
	fmt.Println("       gitmap os theme light    (Set light mode appearance)")
	fmt.Println("       gitmap os theme status   (Inspect current theme appearance)")
}
