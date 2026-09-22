package cmdssh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveSSHLogPath() string {
	dir := store.BinaryDataDir()
	logDir := filepath.Join(dir, "logs")
	_ = os.MkdirAll(logDir, 0755)

	return filepath.Join(logDir, "ssh_error.log")
}

func writeTempMirrorLog(content []byte) {
	tempDir := filepath.Join(".ai-memory", "temp")
	if info, err := os.Stat(tempDir); err == nil && info.IsDir() {
		_ = os.WriteFile(filepath.Join(tempDir, "ssh_last_error.log"), content, 0644)
	}
}

// PersistLog writes the execution trace to disk and mirrors to workspace temp if present.
func (t *SSHExecutionTrace) PersistLog() string {
	if t == nil {
		return ""
	}

	logPath := resolveSSHLogPath()
	t.LogPath = logPath

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return logPath
	}

	_ = os.WriteFile(logPath, data, 0644)
	writeTempMirrorLog(data)

	return logPath
}

// FormatTerminalReport renders human-readable diagnostic execution logs.
func (t *SSHExecutionTrace) FormatTerminalReport() string {
	if t == nil {
		return "No SSH execution trace recorded."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  ● SSH Operation: %s\n", t.Operation))
	sb.WriteString(fmt.Sprintf("    Target: %s\n", t.Target))
	sb.WriteString("  ● Steps Executed:\n")

	for _, step := range t.Steps {
		icon := "✓"
		if step.Status == "FAILED" {
			icon = "✗"
		} else if step.Status == "SKIPPED" {
			icon = "○"
		}
		sb.WriteString(fmt.Sprintf("    %s [%d] %s: %s\n", icon, step.Index, step.Name, step.Details))
		if step.InternalError != "" {
			sb.WriteString(fmt.Sprintf("        Internal Error: %s\n", step.InternalError))
		}
	}

	if t.InternalError != "" {
		sb.WriteString(fmt.Sprintf("  ● Root Internal Error: %s\n", t.InternalError))
	}
	if t.Suggestion != "" {
		sb.WriteString(fmt.Sprintf("  ● Actionable Hint: %s\n", t.Suggestion))
	}
	if t.LogPath != "" {
		sb.WriteString(fmt.Sprintf("  ● Diagnostic Log: %s\n", t.LogPath))
	}

	return sb.String()
}
