package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runOSAutoLogin(args []string) error {
	if len(args) == 0 {
		return handleAutoLoginInteractive()
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "help", "-h", "--help":
		printAutoLoginUsage()
		return nil
	case "status", "st", "info":
		return handleAutoLoginStatus()
	case "disable", "off", "remove", "rm":
		return handleAutoLoginDisable()
	case "enable", "on", "set":
		return handleAutoLoginEnable(args[1:])
	default:
		return handleAutoLoginFallback(args)
	}
}

func handleAutoLoginInteractive() error {
	cfg, err := promptAutoLoginConfig()
	if err != nil {
		return err
	}
	if len(cfg.Username) == 0 {
		return apperror.NewSimple("username is required", "E_AUTOLOGIN_USER_REQUIRED")
	}
	engine := GetAutoLoginEngine()
	if err := engine.Configure(cfg); err != nil {
		return err
	}
	fmt.Printf("✔ OS auto-login configured for user: %s (domain: %s)\n", cfg.Username, cfg.Domain)
	return nil
}

func handleAutoLoginDisable() error {
	engine := GetAutoLoginEngine()
	if err := engine.Disable(); err != nil {
		return err
	}
	fmt.Println("✔ OS auto-login disabled successfully.")
	return nil
}

func handleAutoLoginStatus() error {
	engine := GetAutoLoginEngine()
	st, err := engine.Status()
	if err != nil {
		return err
	}
	printAutoLoginStatusReport(st)
	return nil
}

func printAutoLoginStatusReport(st AutoLoginStatus) {
	fmt.Println("▶ OS Auto-Login Status")
	fmt.Printf("  • Enabled:         %t\n", st.IsEnabled)
	fmt.Printf("  • Username:        %s\n", st.Username)
	fmt.Printf("  • Domain:          %s\n", st.Domain)
	fmt.Printf("  • Password Stored: %t\n", st.HasPassword)
	fmt.Printf("  • Display Manager: %s\n", st.DisplayManager)
}
