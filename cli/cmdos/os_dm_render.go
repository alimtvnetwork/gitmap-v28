package cmdos

import (
	"fmt"
)

func renderOSDMStatus(s DMStatus) {
	fmt.Println("▶ Display Manager & Session Status:")
	fmt.Printf("  • Display Manager: %s\n", s.Name)
	fmt.Printf("  • Service Status:  %s\n", s.ServiceStatus)
	fmt.Printf("  • Session Type:    %s\n", s.SessionType)

	waylandState := "Disabled (forcing X11)"
	if s.IsWaylandEnabled {
		waylandState = "Enabled"
	}
	fmt.Printf("  • Wayland Support: %s\n", waylandState)

	autoLoginState := "Disabled"
	if s.IsAutoLoginSet {
		autoLoginState = fmt.Sprintf("Enabled (user: %s)", s.AutoLoginUser)
	}
	fmt.Printf("  • Auto-Login:      %s\n", autoLoginState)
	fmt.Printf("  • Config File:     %s\n", s.ConfigFile)
}

func printOSDMHelp() {
	fmt.Println("Usage: gitmap os dm <subcommand> [flags]")
	fmt.Println("")
	fmt.Println("Subcommands:")
	fmt.Println("  status             Display current Display Manager & session status")
	fmt.Println("  wayland <on|off>   Enable or disable Wayland (forcing X11 on GDM3)")
	fmt.Println("  restart            Restart display-manager service")
	fmt.Println("  help               Show this help message")
}
