//go:build windows

package cmdos

import (
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

const (
	dataCollectionKey = `SOFTWARE\Policies\Microsoft\Windows\DataCollection`
	systemPoliciesKey = `SOFTWARE\Policies\Microsoft\Windows\System`
	searchPoliciesKey = `Software\Microsoft\Windows\CurrentVersion\Search`
)

func disabledToZero(isDisabled bool) uint32 {
	if isDisabled {
		return 0
	}
	return 1
}

func applyTelemetry(isDisabled bool) error {
	defer toggleDiagTrackService(isDisabled)

	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, dataCollectionKey, registry.SET_VALUE)
	if err != nil {
		return nil
	}
	defer k.Close()

	_ = k.SetDWordValue("AllowTelemetry", disabledToZero(isDisabled))

	return nil
}

func toggleDiagTrackService(isDisabled bool) {
	if isDisabled {
		_ = exec.Command("sc.exe", "config", "DiagTrack", "start=", "disabled").Run()
		_ = exec.Command("sc.exe", "stop", "DiagTrack").Run()
		_ = exec.Command("setx.exe", "POWERSHELL_TELEMETRY_OPTOUT", "1", "/M").Run()

		return
	}

	_ = exec.Command("sc.exe", "config", "DiagTrack", "start=", "auto").Run()
	_ = exec.Command("sc.exe", "start", "DiagTrack").Run()
	_ = exec.Command("setx.exe", "POWERSHELL_TELEMETRY_OPTOUT", "0", "/M").Run()
}

func applyActivityFeed(isDisabled bool) error {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, systemPoliciesKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	val := disabledToZero(isDisabled)
	_ = k.SetDWordValue("PublishUserActivities", val)
	_ = k.SetDWordValue("UploadUserActivities", val)

	return nil
}

func applyBingSearch(isDisabled bool) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, searchPoliciesKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	val := disabledToZero(isDisabled)
	_ = k.SetDWordValue("BingSearchEnabled", val)
	_ = k.SetDWordValue("CortanaConsent", val)

	return nil
}
