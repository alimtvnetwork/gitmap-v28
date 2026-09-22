package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func handleTweakTelemetry(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("telemetry requires 'off' or 'on'", "E_TWEAK_PARAM_MISSING")
	}

	action := strings.ToLower(args[0])
	isDisabled := action == "off" || action == "disable" || action == "0"
	if err := GetTweakEngine().SetTelemetry(isDisabled); err != nil {
		return err
	}

	state := "enabled (default)"
	if isDisabled {
		state = "disabled (DiagTrack service stopped & opt-out set)"
	}
	fmt.Printf("✔ Windows Telemetry set to: %s\n", state)

	return nil
}

func handleTweakActivity(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("activity requires 'off' or 'on'", "E_TWEAK_PARAM_MISSING")
	}

	action := strings.ToLower(args[0])
	isDisabled := action == "off" || action == "disable" || action == "0"
	if err := GetTweakEngine().SetActivityFeed(isDisabled); err != nil {
		return err
	}

	state := "enabled (default)"
	if isDisabled {
		state = "disabled (activity publication & upload suppressed)"
	}
	fmt.Printf("✔ Activity History Feed set to: %s\n", state)

	return nil
}

func handleTweakSearch(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("search requires 'clean' or 'default'", "E_TWEAK_PARAM_MISSING")
	}

	action := strings.ToLower(args[0])
	isDisabled := action == "clean" || action == "off" || action == "disable"
	if err := GetTweakEngine().SetBingSearch(isDisabled); err != nil {
		return err
	}

	state := "default (Bing search enabled in Start Menu)"
	if isDisabled {
		state = "clean (Bing and web search removed from Start Menu)"
	}
	fmt.Printf("✔ Start Menu Search set to: %s\n", state)

	return nil
}
