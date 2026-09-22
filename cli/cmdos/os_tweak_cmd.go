package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runOSTweak(args []string) error {
	if len(args) == 0 {
		return handleTweakStatus()
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "help", "-h", "--help":
		printTweakUsage()
		return nil
	case "status", "st", "info":
		return handleTweakStatus()
	case "context-menu", "right-click", "context":
		return handleTweakContextMenu(args[1:])
	case "start-menu", "startmenu", "start":
		return handleTweakStartMenu(args[1:])
	case "power", "power-scheme", "scheme":
		return handleTweakPower(args[1:])
	case "hibernate", "hibernation":
		return handleTweakHibernate(args[1:])
	case "telemetry", "telemetry-optout":
		return handleTweakTelemetry(args[1:])
	case "activity", "activity-feed", "activities":
		return handleTweakActivity(args[1:])
	case "search", "bing-search", "bing":
		return handleTweakSearch(args[1:])
	default:
		printTweakUsage()
		return apperror.NewSimple("unknown tweak category: "+args[0], "E_INVALID_TWEAK_CAT")
	}
}

func handleTweakStatus() error {
	engine := GetTweakEngine()
	st, err := engine.GetStatus()
	if err != nil {
		return err
	}
	printTweakStatusReport(st)
	return nil
}

func printTweakStatusReport(st TweakStatus) {
	fmt.Println("▶ Windows System Tweaks Status")
	fmt.Printf("  • Classic Context Menu: %t\n", st.IsClassicContextMenu)
	fmt.Printf("  • Classic Start Menu:   %t\n", st.IsClassicStartMenu)
	fmt.Printf("  • Ultimate Power:       %t\n", st.IsUltimatePower)
	fmt.Printf("  • Hibernation Active:   %t\n", st.IsHibernateEnabled)
	fmt.Printf("  • Telemetry Disabled:   %t\n", st.IsTelemetryDisabled)
	fmt.Printf("  • Activity Disabled:    %t\n", st.IsActivityFeedDisabled)
	fmt.Printf("  • Bing Search Disabled: %t\n", st.IsBingSearchDisabled)
}

func printTweakUsage() {
	fmt.Println("Usage: gitmap os tweak <category> <action>")
	fmt.Println("       gitmap os tweak status                       (Inspect tweak states)")
	fmt.Println("       gitmap os tweak context-menu classic|modern  (Toggle right-click menu)")
	fmt.Println("       gitmap os tweak start-menu classic|default   (Toggle Start Menu layout)")
	fmt.Println("       gitmap os tweak power ultimate|balanced      (Switch power scheme)")
	fmt.Println("       gitmap os tweak hibernate on|off             (Toggle hibernation)")
	fmt.Println("       gitmap os tweak telemetry off|on             (Toggle Windows telemetry)")
	fmt.Println("       gitmap os tweak activity off|on              (Toggle activity history feed)")
	fmt.Println("       gitmap os tweak search clean|default         (Toggle Start Menu web search)")
}
