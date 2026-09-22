//go:build windows

package cmdos

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	guidUltimatePerformance = "e9a42b02-d5df-448d-aa00-03f14749eb61"
	guidBalanced            = "381b4222-f694-41f0-9685-ff5bb260df2e"
)

func applyPowerScheme(isUltimate bool) error {
	if isUltimate {
		return activateUltimatePerformance()
	}
	return activateBalancedPower()
}

func activateUltimatePerformance() error {
	_ = exec.Command("powercfg", "/duplicatescheme", guidUltimatePerformance).Run()
	cmd := exec.Command("powercfg", "/setactive", guidUltimatePerformance)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "failed to activate Ultimate Performance power scheme")
	}
	return nil
}

func activateBalancedPower() error {
	cmd := exec.Command("powercfg", "/setactive", guidBalanced)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "failed to activate Balanced power scheme")
	}
	return nil
}

func isUltimatePowerActive() bool {
	out, err := exec.Command("powercfg", "/getactivescheme").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), guidUltimatePerformance)
}

func applyHibernate(isEnabled bool) error {
	arg := "/hibernate"
	mode := "off"
	if isEnabled {
		mode = "on"
	}
	cmd := exec.Command("powercfg", arg, mode)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "failed to toggle hibernate (run as Administrator)")
	}
	return nil
}
