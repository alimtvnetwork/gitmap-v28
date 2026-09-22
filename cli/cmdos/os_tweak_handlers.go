package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func handleTweakContextMenu(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("context-menu requires 'classic' or 'modern'", "E_TWEAK_PARAM_MISSING")
	}
	action := strings.ToLower(args[0])
	isClassic := action == "classic" || action == "win10" || action == "on"
	if err := GetTweakEngine().SetContextMenu(isClassic); err != nil {
		return err
	}
	mode := "modern (Windows 11)"
	if isClassic {
		mode = "classic (Windows 10)"
	}
	fmt.Printf("✔ Right-click context menu set to: %s\n", mode)
	return nil
}

func handleTweakStartMenu(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("start-menu requires 'classic' or 'default'", "E_TWEAK_PARAM_MISSING")
	}
	action := strings.ToLower(args[0])
	isClassic := action == "classic" || action == "on"
	if err := GetTweakEngine().SetStartMenu(isClassic); err != nil {
		return err
	}
	mode := "default"
	if isClassic {
		mode = "classic layout"
	}
	fmt.Printf("✔ Start menu layout set to: %s\n", mode)
	return nil
}

func handleTweakPower(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("power requires 'ultimate' or 'balanced'", "E_TWEAK_PARAM_MISSING")
	}
	action := strings.ToLower(args[0])
	isUltimate := action == "ultimate" || action == "high" || action == "perf"
	if err := GetTweakEngine().SetPowerScheme(isUltimate); err != nil {
		return err
	}
	mode := "Balanced"
	if isUltimate {
		mode = "Ultimate Performance"
	}
	fmt.Printf("✔ Power scheme set to: %s\n", mode)
	return nil
}

func handleTweakHibernate(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("hibernate requires 'on' or 'off'", "E_TWEAK_PARAM_MISSING")
	}
	action := strings.ToLower(args[0])
	isEnabled := action == "on" || action == "enable" || action == "true"
	if err := GetTweakEngine().SetHibernate(isEnabled); err != nil {
		return err
	}
	mode := "disabled (hiberfil.sys deleted)"
	if isEnabled {
		mode = "enabled"
	}
	fmt.Printf("✔ Hibernation %s\n", mode)
	return nil
}
