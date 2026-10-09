package cmdui

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderUIHelp displays the styled two-column UI help menu.
func RenderUIHelp() {
	termout.RenderMenu(buildUIHelpMenu())
}

func buildUIHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Interactive Web Dashboard (gitmap ui)",
		UsageLines: []string{
			"gitmap ui [page] [flags]",
			"gitmap ui settings --port 8080",
			"gitmap ui --port 9090 --no-browser",
		},
		Sections: []termout.HelpSection{
			buildUIPagesSection(),
			buildUIOptionsSection(),
		},
		FooterFlags: buildUIFooterFlags(),
		Tips: []string{
			"Launches a reactive local control panel for repos, cluster nodes, and config.",
			"Use '--no-browser' when running inside headless or SSH terminal sessions.",
		},
	}
}

func buildUIFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-p, --port <port>", Description: "HTTP listen port for Web UI server (default: 8080)"},
		{Command: "--host <addr>", Description: "Bind IP address (default: 127.0.0.1; non-loopback warns — exposes the embedded terminal)"},
		{Command: "--no-browser", Description: "Do not automatically launch system default web browser"},
		{Command: "-h, --help", Description: "Show this UI help menu"},
	}
}

func buildUIPagesSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Available Dashboard Pages",
		Entries: []termout.CommandEntry{
			{Command: "settings", Description: "CLI configuration, theme picker, and path settings"},
			{Command: "commitin", Description: "Multi-option commit engine with direction flows"},
			{Command: "ssh", Description: "SSH cluster nodes, fleet key deployment, and firewall helpers"},
			{Command: "macro", Description: "Macro automation builder"},
			{Command: "installer", Description: "Custom multi-OS installer manager"},
			{Command: "prompts", Description: "Prompts and AI instructions manager"},
			{Command: "import-export", Description: "System configuration and SSH fleet JSON import/export"},
			{Command: "schedules", Description: "Task and cron schedule manager"},
			{Command: "editor", Description: "Remote/local file editor"},
			{Command: "help", Description: "Built-in help panel and usage reference"},
		},
	}
}

func buildUIOptionsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Server Architecture",
		Entries: []termout.CommandEntry{
			{Command: "Embedded Assets", Description: "Zero-dependency single binary serving reactive web assets"},
			{Command: "REST API", Description: "Native Go JSON endpoints (poll-based live refresh)"},
		},
	}
}
