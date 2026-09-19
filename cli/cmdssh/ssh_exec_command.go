package cmdssh

import (
	"strings"
)

func extractFirstToken(arg string) string {
	fields := strings.Fields(arg)
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}

func determineFallbackShell(osType string) string {
	if isWindowsOS(osType) {
		return "ps"
	}
	return "bash"
}

func determineSSHCommand(osType string, args []string) (string, string, bool) {
	if len(args) == 0 {
		return "", "", false
	}
	firstToken := extractFirstToken(args[0])
	if isGitmapCommand(firstToken) {
		return "", resolveGitmapCommandString(args), true
	}
	if isExplicitShell(firstToken) {
		return firstToken, extractShellCommandArgs(args), false
	}
	cmdStr := normalizeMultiCommands(strings.Join(args, " "), isWindowsOS(osType))
	return determineFallbackShell(osType), cmdStr, false
}

func normalizeMultiCommands(cmdStr string, isWindows bool) string {
	if !strings.Contains(cmdStr, ",") {
		return cmdStr
	}
	parts := strings.Split(cmdStr, ",")
	var cleanParts []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			cleanParts = append(cleanParts, trimmed)
		}
	}
	if len(cleanParts) <= 1 {
		return cmdStr
	}
	if isWindows {
		return strings.Join(cleanParts, "; ")
	}
	return strings.Join(cleanParts, " && ")
}

func isGitmapCoreCommand(cmd string) bool {
	switch cmd {
	case "gitmap", "status", "pipeline", "pipe", "pl", "clone", "pull", "sync", "push", "clean", "log", "branch", "diff":
		return true
	default:
		return false
	}
}

func isGitmapSystemCommand(cmd string) bool {
	switch cmd {
	case "storage", "macro", "install", "update", "setup", "chrome", "vscode", "vsc", "vhost", "zip", "service", "os":
		return true
	case "schedule", "schedules", "scheduled", "cron", "crontab", "restore-db", "restoredb":
		return true
	default:
		return false
	}
}

func isGitmapAgyOrRemote(cmd string) bool {
	switch cmd {
	case "agy", "ag", "antigravity", "aef", "fix-pipeline", "fixpipeline", "pipeline-fix":
		return true
	case "prompts-template", "prompt-template", "prompts-templates", "prompt-templates", "prompt_templates", "pt":
		return true
	case "ssh", "se", "sj", "cluster", "sc", "mkdir", "cat", "prompt", "prompts", "pmt", "agm", "ip":
		return true
	default:
		return false
	}
}

func isGitmapCommand(first string) bool {
	low := strings.ToLower(first)
	return isGitmapCoreCommand(low) || isGitmapSystemCommand(low) || isGitmapAgyOrRemote(low)
}

func resolveGitmapCommandString(args []string) string {
	if len(args) == 0 {
		return "gitmap"
	}
	joined := strings.Join(args, " ")
	if strings.HasPrefix(joined, "gitmap ") || joined == "gitmap" {
		return joined
	}

	return "gitmap " + joined
}

func isExplicitShell(s string) bool {
	return s == "ps" || s == "cmd" || s == "bash" || s == "sh"
}

func extractShellCommandArgs(args []string) string {
	if len(args) > 1 {
		return strings.Join(args[1:], " ")
	}
	if len(args) == 0 {
		return ""
	}
	fields := strings.Fields(args[0])
	if len(fields) > 1 {
		return strings.Join(fields[1:], " ")
	}

	return ""
}
