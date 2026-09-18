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

	shell := "bash"
	if isWindowsOS(osType) {
		shell = "ps"
	}

	return shell, strings.Join(args, " "), false
}

func isGitmapCommand(first string) bool {
	switch first {
	case "gitmap", "status", "pipeline", "pipe", "pl", "clone", "pull", "sync", "push", "clean", "log", "branch", "diff":
		return true
	case "storage", "macro", "install", "update", "setup", "chrome", "vscode", "vsc", "vhost", "zip", "service":
		return true
	case "os", "schedule", "schedules", "agy", "ag", "antigravity", "aef", "fix-pipeline", "prompts-template", "pt":
		return true
	case "ssh", "se", "sj", "cluster", "sc", "mkdir", "cat", "prompt", "prompts", "pmt", "agm", "ip":
		return true
	default:
		return false
	}
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
	if len(args) == 1 {
		fields := strings.Fields(args[0])
		if len(fields) > 1 {
			return strings.Join(fields[1:], " ")
		}
	}

	return ""
}
