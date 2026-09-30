package cmdos

func buildDNSNetTab() OSTUITab {
	return OSTUITab{
		ID:          "dns-net",
		Title:       "DNS & Net",
		ShortcutKey: "4",
		Items: []OSTUIItem{
			{ID: "dns-cf", Title: "Cloudflare DNS (1.1.1.1)", Description: "Ultra-fast privacy-focused resolver", Badge: "Fast", IsSupported: true},
			{ID: "dns-google", Title: "Google DNS (8.8.8.8)", Description: "Reliable global Anycast resolver", Badge: "Global", IsSupported: true},
			{ID: "dns-quad9", Title: "Quad9 DNS (9.9.9.9)", Description: "Malware-blocking secure resolver", Badge: "Security", IsSupported: true},
			{ID: "dns-dhcp", Title: "Revert to Automatic DHCP", Description: "Clear custom nameservers on primary adapter", Badge: "Default", IsSupported: true},
		},
	}
}

func buildCleanStorageTab() OSTUITab {
	return OSTUITab{
		ID:          "clean-storage",
		Title:       "Clean & Storage",
		ShortcutKey: "5",
		Items: []OSTUIItem{
			{ID: "clean-dev", Title: "Clean Dev Caches", Description: "Purge Go build cache, npm/pnpm caches", Badge: "Dev", IsSupported: true},
			{ID: "clean-ai", Title: "Clean AI Caches & Logs", Description: "Purge Antigravity scratch, transcripts, and logs", Badge: "AI", IsSupported: true},
			{ID: "clean-term", Title: "Clean Terminal Temp", Description: "Truncate terminal outputs and scratch files", Badge: "Disk", IsSupported: true},
			{ID: "clean-sys", Title: "Clean System Packages", Description: "Purge unused package archives and vacuum journals", Badge: "System", IsSupported: true},
		},
	}
}

func buildUpdateTab() OSTUITab {
	return OSTUITab{
		ID:          "update",
		Title:       "Update",
		ShortcutKey: "6",
		Items: []OSTUIItem{
			{ID: "update-repos", Title: "Refresh Repositories", Description: "Update package indexes without upgrading packages", Badge: "Sync", IsSupported: true},
			{ID: "update-upgrade", Title: "Full System Upgrade", Description: "Upgrade all system packages across detected managers", Badge: "Upgrade", IsSupported: true},
			{ID: "update-mirrors", Title: "Fix Regional Mirrors", Description: "Switch mirrors to fast US canonical upstream", Badge: "Mirrors", IsSupported: true},
		},
	}
}
