package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// SUGRuntimeState captures live background monitoring status.
type SUGRuntimeState struct {
	PID            int      `json:"pid"`
	Status         string   `json:"status"` // "running" or "idle"
	StartedAt      string   `json:"startedAt"`
	IntervalSec    int      `json:"intervalSeconds"`
	ProjectTargets []string `json:"projectTargets"`
}

func resolveSUGStatePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".gitmap", "sug_state.json")
	}
	dir := filepath.Join(home, ".gemini", "antigravity")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "sug_state.json")
}

// RecordSUGWatchStarted records the active watch loop process and timestamp.
func RecordSUGWatchStarted(projects []string, interval time.Duration) {
	st := SUGRuntimeState{
		PID:            os.Getpid(),
		Status:         "running",
		StartedAt:      time.Now().UTC().Format(time.RFC3339),
		IntervalSec:    int(interval.Seconds()),
		ProjectTargets: projects,
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err == nil {
		_ = os.WriteFile(resolveSUGStatePath(), data, 0644)
	}
}

// RecordSUGWatchStopped marks the state as idle.
func RecordSUGWatchStopped() {
	st := SUGRuntimeState{
		PID:            0,
		Status:         "idle",
		StartedAt:      "",
		IntervalSec:    0,
		ProjectTargets: []string{},
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err == nil {
		_ = os.WriteFile(resolveSUGStatePath(), data, 0644)
	}
}

// LoadSUGRuntimeStatus evaluates live process state.
func LoadSUGRuntimeStatus() (SUGRuntimeState, bool) {
	path := resolveSUGStatePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return SUGRuntimeState{Status: "idle"}, false
	}
	var st SUGRuntimeState
	if err := json.Unmarshal(data, &st); err != nil {
		return SUGRuntimeState{Status: "idle"}, false
	}
	if st.PID <= 0 || st.Status != "running" {
		return st, false
	}
	if isProcessAlive(st.PID) {
		return st, true
	}
	RecordSUGWatchStopped()
	st.Status = "idle"
	return st, false
}

func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	if runtime.GOOS == "windows" {
		return isWindowsProcessAlive(pid)
	}
	return isUnixProcessAlive(pid)
}

func isWindowsProcessAlive(pid int) bool {
	out, err := exec.Command("tasklist", "/NH", "/FI", fmt.Sprintf("PID eq %d", pid)).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), strconv.Itoa(pid))
}

func isUnixProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(os.Signal(nil)) == nil
}

