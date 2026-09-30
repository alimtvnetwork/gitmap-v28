package cmdos

import "strings"

func isTweakAction(id string) bool {
	return strings.HasPrefix(id, "tweak-")
}

func dispatchTweakAction(id, title string) OSTUIActionResult {
	eng := GetTweakEngine()
	switch id {
	case "tweak-telemetry":
		return actionResultFromError(title, eng.SetTelemetry(true), "telemetry disabled")
	case "tweak-activity":
		return actionResultFromError(title, eng.SetActivityFeed(true), "activity feed disabled")
	case "tweak-search":
		return actionResultFromError(title, eng.SetBingSearch(true), "bing search disabled")
	case "tweak-context":
		return actionResultFromError(title, eng.SetContextMenu(true), "classic context menu enabled")
	case "tweak-power":
		return actionResultFromError(title, eng.SetPowerScheme(true), "ultimate power scheme applied")
	case "tweak-theme":
		themeEng := newPlatformThemeEngine()
		return actionResultFromError(title, themeEng.SetTheme(ThemeModeDark), "dark theme enabled")
	default:
		return OSTUIActionResult{Title: title, IsSuccess: false, Message: "unknown tweak"}
	}
}

func isAutoLoginAction(id string) bool {
	return strings.HasPrefix(id, "al-")
}

func dispatchAutoLoginAction(id, title string) OSTUIActionResult {
	al := GetAutoLoginEngine()
	switch id {
	case "al-disable":
		return actionResultFromError(title, al.Disable(), "autologin disabled")
	case "al-status":
		st, err := al.Status()
		if err != nil {
			return actionResultFromError(title, err, "")
		}
		return OSTUIActionResult{Title: title, IsSuccess: true, Message: "active user: " + st.Username}
	default:
		return OSTUIActionResult{Title: title, IsSuccess: false, Message: "unknown autologin action"}
	}
}

func isDisplayDMAction(id string) bool {
	return strings.HasPrefix(id, "disp-") || strings.HasPrefix(id, "dm-")
}

func dispatchDisplayDMAction(id, title string) OSTUIActionResult {
	switch id {
	case "disp-never":
		return actionResultFromError(title, runPowerNeverSleep(), "display set to never sleep")
	case "disp-reset":
		return actionResultFromError(title, runPowerReset(), "display sleep reset to 15m")
	case "dm-wayland":
		mgr := newLinuxDMManager()
		return actionResultFromError(title, mgr.SetWayland(false), "wayland toggled")
	case "dm-restart":
		mgr := newLinuxDMManager()
		return actionResultFromError(title, mgr.RestartService(), "display-manager restarted")
	default:
		return OSTUIActionResult{Title: title, IsSuccess: false, Message: "unknown display action"}
	}
}
