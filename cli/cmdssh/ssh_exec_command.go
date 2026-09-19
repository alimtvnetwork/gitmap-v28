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
	if isIPCommand(args) {
		return resolveIPCommand(osType)
	}
	first := extractFirstToken(args[0])
	if isGitmapCommand(first) {
		return "", resolveGitmapCommandString(args), true
	}
	return resolveFallbackOrShell(osType, args, first)
}

func resolveFallbackOrShell(osType string, args []string, first string) (string, string, bool) {
	if isExplicitShell(first) {
		return first, extractShellCommandArgs(args), false
	}
	cmdStr := normalizeMultiCommands(strings.Join(args, " "), isWindowsOS(osType))
	return determineFallbackShell(osType), cmdStr, false
}

func isIPCommand(args []string) bool {
	return len(args) == 1 && strings.EqualFold(args[0], "ip")
}

func resolveIPCommand(osType string) (string, string, bool) {
	if isWindowsOS(osType) {
		return "ps", "Get-NetIPAddress -AddressFamily IPv4 | Select-Object -ExpandProperty IPAddress", false
	}

	return "bash", "ip -br a 2>/dev/null || ip a 2>/dev/null || hostname -I 2>/dev/null || ifconfig", false
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
	case "gitmap", "status", "st", "pipeline", "pipe", "pl", "clone", "pull", "sync", "push", "clean", "log", "branch", "diff":
		return true
	case "open", "o", "browse", "browse-url", "open-url", "pull-all", "clone-all", "sync-all", "push-all":
		return true
	case "clone-sync", "clone-only-missing", "clone-next", "clone-pick", "clone-from", "clone-now", "clone-reclone":
		return true
	case "reconcile", "latest-branch", "discard", "stash", "wip":
		return true
	default:
		return false
	}
}

func isGitmapSystemCommand(cmd string) bool {
	switch cmd {
	case "storage", "macro", "install", "reinstall", "uninstall", "update", "setup", "chrome", "vscode", "vsc", "vhost", "zip":
		return true
	case "service", "os", "power", "schedule", "schedules", "scheduled", "cron", "crontab":
		return true
	case "restore-db", "restoredb", "db", "db-reset", "db-migrate", "start-fresh", "backup", "stats", "task", "tasks", "watch":
		return true
	default:
		return false
	}
}

func isGitmapAgyOrRemote(cmd string) bool {
	switch cmd {
	case "agy", "ag", "antigravity", "aef", "fix-pipeline", "fixpipeline", "pipeline-fix", "pipeline-ai", "plai":
		return true
	case "prompts-template", "prompt-template", "prompts-templates", "prompt-templates", "prompt_templates", "pt":
		return true
	case "ssh", "se", "sj", "cluster", "sc", "mkdir", "cat", "prompt", "prompts", "pmt", "agm", "ip":
		return true
	default:
		return false
	}
}

func isGitmapUtilityCommand(cmd string) bool {
	switch cmd {
	case "doctor", "profile", "profiles", "config", "workdir", "cg", "codingguidelines", "coding-guidelines":
		return true
	case "ai", "cargo", "aum", "automation", "fix-auth", "fixauth", "ssh-bind", "help", "docs":
		return true
	default:
		return false
	}
}

func isGitmapCommand(first string) bool {
	low := strings.ToLower(first)
	if isGitmapCoreCommand(low) || isGitmapSystemCommand(low) {
		return true
	}
	return isGitmapAgyOrRemote(low) || isGitmapUtilityCommand(low)
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
