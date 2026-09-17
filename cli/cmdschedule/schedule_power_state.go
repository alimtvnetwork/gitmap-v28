package cmdschedule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PowerScheduleState stores metadata for pending scheduled power actions.
type PowerScheduleState struct {
	Action          OSActionType `json:"action"`
	Status          string       `json:"status"`
	ScheduledAt     time.Time    `json:"scheduledAt"`
	TriggerAt       time.Time    `json:"triggerAt"`
	DurationSeconds int64        `json:"durationSeconds"`
	RawDuration     string       `json:"rawDuration"`
	Command         string       `json:"command"`
}

var stateFileOverride string

func resolvePowerStateFilePath() string {
	if stateFileOverride != "" {
		return stateFileOverride
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "gitmap_power_schedule.json")
	}
	return filepath.Join(home, ".gitmap", "power_schedule.json")
}

func ensurePowerStateDir(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "ensurePowerStateDir")
	}
	return nil
}

// SavePowerScheduleState persists the power schedule state to disk.
func SavePowerScheduleState(s PowerScheduleState) error {
	path := resolvePowerStateFilePath()
	if err := ensurePowerStateDir(path); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "SavePowerScheduleState_Marshal")
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return apperror.WrapSimple(err, "SavePowerScheduleState_Write")
	}
	return nil
}

func readPowerScheduleState(path string) (*PowerScheduleState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s PowerScheduleState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, apperror.WrapSimple(err, "readPowerScheduleState_Unmarshal")
	}
	return &s, nil
}

// GetActivePowerSchedule returns the currently armed power schedule or nil.
func GetActivePowerSchedule() (*PowerScheduleState, error) {
	path := resolvePowerStateFilePath()
	s, err := readPowerScheduleState(path)
	if err != nil {
		return nil, nil
	}
	if s.Status != "ARMED" {
		return nil, nil
	}
	if time.Now().After(s.TriggerAt) {
		s.Status = "COMPLETED"
		_ = SavePowerScheduleState(*s)
		return nil, nil
	}
	return s, nil
}

// CancelActivePowerSchedule marks the active power schedule as canceled.
func CancelActivePowerSchedule() error {
	path := resolvePowerStateFilePath()
	s, err := readPowerScheduleState(path)
	if err != nil {
		return nil
	}
	s.Status = "CANCELED"
	return SavePowerScheduleState(*s)
}

// FormatRemainingDuration converts duration until target into human string.
func FormatRemainingDuration(target time.Time) string {
	diff := time.Until(target)
	if diff <= 0 {
		return "0s"
	}
	h := int(diff.Hours())
	m := int(diff.Minutes()) % 60
	s := int(diff.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
