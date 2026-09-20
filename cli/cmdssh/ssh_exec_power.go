package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func isPowerCommand(first string) bool {
	low := strings.ToLower(first)

	return low == "shutdown" || low == "reboot" || low == "poweroff" || low == "halt" || low == "restart"
}

func isRebootAction(first string, args []string) bool {
	low := strings.ToLower(first)
	if low == "reboot" || low == "restart" {
		return true
	}
	if low == "shutdown" {
		return hasRebootFlag(args)
	}

	return false
}

func hasRebootFlag(args []string) bool {
	for _, a := range args {
		if a == "-r" || a == "/r" || strings.HasPrefix(a, "-r") || strings.HasPrefix(a, "/r") {
			return true
		}
	}

	return false
}

func resolvePowerCommand(c db.SSHConnection, args []string) (string, string) {
	first := extractFirstToken(args[0])
	isReboot := isRebootAction(first, args)
	if isWindowsOS(c.OS) {
		return resolveWindowsPowerCmd(isReboot)
	}

	return resolveUnixPowerCmd(c, isReboot)
}

func resolveWindowsPowerCmd(isReboot bool) (string, string) {
	if isReboot {
		return "ps", "shutdown.exe /r /t 0 /f"
	}

	return "ps", "shutdown.exe /s /t 0 /f"
}

func resolveUnixPowerCmd(c db.SSHConnection, isReboot bool) (string, string) {
	enc := c.EncryptedPassword
	if enc == "" {
		enc = queryHostPasswordFromDB(c.Alias, c.IPAddress)
	}
	plainPass, err := decryptPasswordCandidate(enc)
	if err == nil && plainPass != "" {
		return "bash", buildUnixPassPowerCmd(plainPass, isReboot)
	}

	return "bash", buildUnixSudoPowerCmd(isReboot)
}

func buildUnixPassPowerCmd(pass string, isReboot bool) string {
	escapedPass := strings.ReplaceAll(pass, "'", "'\\''")
	if isReboot {
		return fmt.Sprintf("echo '%s' | sudo -S reboot 2>/dev/null || sudo -n reboot 2>/dev/null || reboot", escapedPass)
	}

	return fmt.Sprintf("echo '%s' | sudo -S poweroff 2>/dev/null || echo '%s' | sudo -S shutdown now 2>/dev/null || sudo -n poweroff 2>/dev/null || shutdown now", escapedPass, escapedPass)
}

func buildUnixSudoPowerCmd(isReboot bool) string {
	if isReboot {
		return "sudo -n reboot 2>/dev/null || sudo -n shutdown -r now 2>/dev/null || systemctl reboot 2>/dev/null || reboot"
	}

	return "sudo -n poweroff 2>/dev/null || sudo -n shutdown now 2>/dev/null || sudo -n systemctl poweroff 2>/dev/null || shutdown now"
}
