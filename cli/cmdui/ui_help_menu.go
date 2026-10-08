package cmdui

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderUIHelp displays the styled two-column UI help menu.
func RenderUIHelp() {
	termhelp.RenderMenu(buildUIHelpMenu())
}

func buildUIHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Interactive Web Dashboard (gitmap ui)",
		UsageLines: []string{
			"gitmap ui [page] [flags]",
			"gitmap ui settings --port 8080",
			"gitmap ui cluster",
		},
		Sections: []termhelp.HelpSection{
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

func buildUIFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-p, --port <port>", Description: "HTTP listen port for Web UI server (default: 8080)"},
		{Command: "--host <addr>", Description: "Bind IP address (default: 127.0.0.1)"},
		{Command: "--no-browser", Description: "Do not automatically launch system default web browser"},
		{Command: "-h, --help", Description: "Show this UI help menu"},
	}
}

func buildUIPagesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Available Dashboard Pages",
		Entries: []termhelp.CommandEntry{
			{Command: "settings", Description: "CLI configuration, theme picker, and path settings"},
			{Command: "dashboard", Description: "Workspace health, commit graphs, and repo inventory"},
			{Command: "cluster", Description: "SSH fleet status, node topology, and telemetry"},
			{Command: "templates", Description: "Templates and substitution variable management"},
			{Command: "repos", Description: "Interactive multi-repository table and branch viewer"},
		},
	}
}

func buildUIOptionsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Server Architecture",
		Entries: []termhelp.CommandEntry{
			{Command: "Embedded Assets", Description: "Zero-dependency single binary serving reactive web assets"},
			{Command: "REST API", Description: "Native Go JSON endpoints with SSE live telemetry"},
		},
	}
}
