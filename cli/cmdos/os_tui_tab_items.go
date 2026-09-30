package cmdos

import "runtime"

func buildTweaksTab() OSTUITab {
	isWin := runtime.GOOS == "windows"

	return OSTUITab{
		ID:          "tweaks",
		Title:       "Tweaks",
		ShortcutKey: "1",
		Items: []OSTUIItem{
			{ID: "tweak-telemetry", Title: "Disable Windows Telemetry", Description: "Turn off DiagTrack and data telemetry", Badge: "Privacy", IsSupported: isWin},
			{ID: "tweak-activity", Title: "Disable Activity Feed", Description: "Prevent user activity publication and upload", Badge: "Privacy", IsSupported: isWin},
			{ID: "tweak-search", Title: "Disable Bing Start Search", Description: "Keep search purely local to device", Badge: "Clean", IsSupported: isWin},
			{ID: "tweak-context", Title: "Classic Context Menu", Description: "Restore Windows 10 style full context menus", Badge: "Desktop", IsSupported: isWin},
			{ID: "tweak-power", Title: "Ultimate Performance", Description: "Activate ultimate performance power scheme", Badge: "Speed", IsSupported: isWin},
			{ID: "tweak-theme", Title: "Dark Mode Theme", Description: "Switch desktop appearance to dark theme", Badge: "Theme", IsSupported: true},
		},
	}
}

func buildAutoLoginTab() OSTUITab {
	return OSTUITab{
		ID:          "autologin",
		Title:       "Auto-Login",
		ShortcutKey: "2",
		Items: []OSTUIItem{
			{ID: "al-status", Title: "Inspect Auto-Login State", Description: "Display current active auto-login configuration", Badge: "Status", IsSupported: true},
			{ID: "al-disable", Title: "Disable Auto-Login", Description: "Remove auto-login credentials and revert to lockscreen", Badge: "Security", IsSupported: true},
		},
	}
}

func buildDisplayDMTab() OSTUITab {
	isLinux := runtime.GOOS == "linux"

	return OSTUITab{
		ID:          "display-dm",
		Title:       "Display/DM",
		ShortcutKey: "3",
		Items: []OSTUIItem{
			{ID: "disp-never", Title: "Never Sleep Display", Description: "Set AC and DC screen timeouts to 0 (always on)", Badge: "Display", IsSupported: true},
			{ID: "disp-reset", Title: "Reset Display Timeout", Description: "Revert screen blanking to default 15 minutes", Badge: "Display", IsSupported: true},
			{ID: "dm-wayland", Title: "Toggle Wayland Session", Description: "Switch between Wayland and X11 in GDM3 config", Badge: "Linux", IsSupported: isLinux},
			{ID: "dm-restart", Title: "Restart Display Manager", Description: "Restart active display-manager service", Badge: "Service", IsSupported: isLinux},
		},
	}
}
