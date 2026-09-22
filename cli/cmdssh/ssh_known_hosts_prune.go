package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var offendingLineRegex = regexp.MustCompile(`(?i)Offending\s+(?:\S+\s+)?key\s+in\s+((?:[a-zA-Z]:)?[^:\r\n]+):(\d+)`)

// ParseOffendingKnownHostsLine extracts file path and line number from OpenSSH error output.
func ParseOffendingKnownHostsLine(output string) (string, int, bool) {
	m := offendingLineRegex.FindStringSubmatch(output)
	if len(m) < 3 {
		return "", 0, false
	}
	filePath := strings.TrimSpace(m[1])
	lineNum, err := strconv.Atoi(m[2])
	if err != nil || lineNum <= 0 {
		return "", 0, false
	}
	return filePath, lineNum, true
}

// isHostKeyChangedError checks if the OpenSSH stderr indicates a host key change or collision.
func isHostKeyChangedError(output string) bool {
	if strings.Contains(output, "REMOTE HOST IDENTIFICATION HAS CHANGED") {
		return true
	}
	if strings.Contains(output, "Host key verification failed") {
		return true
	}
	return strings.Contains(output, "Host key for") && strings.Contains(output, "has changed")
}

// RemoveKnownHostsLine deletes a specific 1-based line number from a known_hosts file.
func RemoveKnownHostsLine(path string, lineNum int) error {
	if path == "" || lineNum <= 0 {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return apperror.WrapSimple(err, "RemoveKnownHostsLine.ReadFile")
	}
	lines := strings.Split(string(data), "\n")
	if lineNum > len(lines) {
		return nil
	}
	lines = append(lines[:lineNum-1], lines[lineNum:]...)
	out := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(out), 0600)
}

func buildPruneTargets(cleanHost string, port int) []string {
	targets := []string{cleanHost}
	if port > 0 {
		targets = append(targets, fmt.Sprintf("[%s]:%d", cleanHost, port))
	}
	if port != 22 {
		targets = append(targets, fmt.Sprintf("[%s]:22", cleanHost))
	}
	return targets
}

func runKeygenPruneTarget(t, path string) {
	var cmd *exec.Cmd
	if path != "" {
		cmd = exec.Command("ssh-keygen", "-R", t, "-f", path)
	} else {
		cmd = exec.Command("ssh-keygen", "-R", t)
	}
	_ = cmd.Run()
}

// PruneHostFromKnownHosts removes all keys for target from ~/.ssh/known_hosts (both plain and hashed).
func PruneHostFromKnownHosts(target string, port int) error {
	cleanHost := stripHostPort(target)
	if cleanHost == "" {
		return nil
	}
	path, _ := DefaultKnownHostsPath()
	if path != "" {
		_, _ = RemoveFromKnownHostsFile(path, cleanHost)
	}
	targets := buildPruneTargets(cleanHost, port)
	for _, t := range targets {
		runKeygenPruneTarget(t, path)
	}
	return nil
}
