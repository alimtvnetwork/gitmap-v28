// Package cmdagy — agy_settings_help.go implements help rendering and flag parsing helpers for speed settings.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderAgySettingsHelp renders the boxed help menu for gitmap settings / gitmap agy settings.
func RenderAgySettingsHelp() {
	renderSettingsHelp()
}

func renderSettingsHelp() {
	termout.RenderMenu(termout.HelpMenu{
		Title: "GitMap & Antigravity Speed Settings (gitmap settings)",
		UsageLines: []string{
			"gitmap settings [--node-threshold 30] [--pae-window 48] [--lap-hours 24] [--threshold 15] [--alias <name>] [--json]",
			"gitmap settings set node_resource_threshold 30",
			"gitmap settings set pae_recency_window_hours 48",
			"gitmap settings set special_repos.secrets_name repo-secrets",
			"gitmap settings set special_repos.cache_name repo-cache",
			"gitmap settings set ai_auto.default_name \"gitmap automation bot - runs autonomously\"",
			"gitmap settings set ai_auto.default_email info@rgitmap.com",
			"gitmap settings set ai_auto.default_agent gemini",
			"gitmap agy settings [set <key> <value> | export <file> | import <file>]",
		},
		Sections: buildSettingsHelpSections(),
	})
}

func buildSettingsHelpSections() []termout.HelpSection {
	return []termout.HelpSection{
		{
			Title: "Speed Settings & Configuration Keys",
			Entries: []termout.CommandEntry{
				{Command: "--node-threshold <N>", Description: "Set CPU/RAM benchmark threshold percentage for nodes & pull (default: 30)"},
				{Command: "--pae-window <N>", Description: "Set commit recency lookback hours for gitmap pae (default: 48)"},
				{Command: "--lap-hours <N>", Description: "Set default lookback hours for last-active-projects (default: 24)"},
				{Command: "--threshold <N>", Description: "Set account-switch remaining credit threshold percentage (default: 15)"},
				{Command: "--alias <name>", Description: "Set machine network alias for SSH fleet recognition"},
				{Command: "set node_resource_threshold <N>", Description: "Set node resource benchmark threshold (default: 30%)"},
				{Command: "set pae_recency_window_hours <N>", Description: "Set PAE commit recency lookback window in hours (default: 48)"},
				{Command: "set special_repos.secrets_name repo-secrets", Description: "Configure special secrets repository name (shortcut: gitmap rs / gitmap cd rs)"},
				{Command: "set special_repos.cache_name repo-cache", Description: "Configure special cache repository name (shortcut: gitmap rc / gitmap cd rc)"},
				{Command: "set ai_auto.default_name <name>", Description: "Set default AI bot git user.name for login set global ai-auto"},
				{Command: "set ai_auto.default_email <email>", Description: "Set default AI bot git user.email (default: info@rgitmap.com)"},
				{Command: "set ai_auto.default_agent <agent>", Description: "Set fallback default AI agent name (e.g. gemini, muse, cursor-model)"},
				{Command: "export / import <file>", Description: "Export or import full Antigravity settings JSON"},
			},
		},
	}
}

func isTgEmailHelpArg(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))
	return low == "help" || low == "-h" || low == "--help"
}

func hasSpeedFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag || strings.HasPrefix(a, flag+"=") {
			return true
		}
	}
	return false
}

func extractSpeedFlag(args []string, primary, short string) string {
	for i, a := range args {
		val := matchSpeedFlagPair(args, i, a, primary, short)
		if len(val) > 0 {
			return val
		}
	}
	return ""
}

func matchSpeedFlagPair(args []string, i int, a, primary, short string) string {
	if (a == primary || a == short) && i+1 < len(args) {
		return args[i+1]
	}
	if strings.HasPrefix(a, primary+"=") {
		return strings.TrimPrefix(a, primary+"=")
	}
	return ""
}

func filterSpeedPositional(args []string) []string {
	var out []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			out = append(out, a)
		}
	}
	return out
}

func printSpeedJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal speed settings json")
	}
	fmt.Println(string(data))
	return nil
}
