package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AccountSwitchState persists active account, threshold, and latest switch telemetry.
type AccountSwitchState struct {
	ThresholdPct       int      `json:"thresholdPct"`
	ActiveAccount      string   `json:"activeAccount"`
	ActiveCreditPct    int      `json:"activeCreditPct"`
	PreviousAccount    string   `json:"previousAccount,omitempty"`
	LastBatchId        string   `json:"lastBatchId,omitempty"`
	BackedProjectCount int      `json:"backedProjectCount"`
	BackedProjectNames []string `json:"backedProjectNames,omitempty"`
	RestoredPrompts    int      `json:"restoredPrompts"`
	IsRefreshed        bool     `json:"isRefreshed"`
	IsLockVerified     bool     `json:"isLockVerified"`
	IsFastForwarded    bool     `json:"isFastForwarded"`
	IsLivenessVerified bool     `json:"isLivenessVerified"`
	IsTelegramNotified bool     `json:"isTelegramNotified"`
	IsEmailNotified    bool     `json:"isEmailNotified"`
	UpdatedAt          string   `json:"updatedAt"`
}

func resolveAccountSwitchStatePath() string {
	return filepath.Join(store.BinaryDataDir(), "account-switch", "state.json")
}

func loadAccountSwitchState() AccountSwitchState {
	data, err := os.ReadFile(resolveAccountSwitchStatePath())
	if err != nil {
		return defaultAccountSwitchState()
	}
	var st AccountSwitchState
	if json.Unmarshal(data, &st) != nil {
		return defaultAccountSwitchState()
	}
	if st.ThresholdPct <= 0 {
		st.ThresholdPct = DefaultAccountSwitchThreshold
	}
	return st
}

func defaultAccountSwitchState() AccountSwitchState {
	return AccountSwitchState{
		ThresholdPct:    DefaultAccountSwitchThreshold,
		ActiveAccount:   "primary-dev@antigravity.dev",
		ActiveCreditPct: 12,
		UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
}

func saveAccountSwitchState(st AccountSwitchState) error {
	path := resolveAccountSwitchStatePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir account-switch state")
	}
	st.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	raw, _ := json.MarshalIndent(st, "", "  ")
	return os.WriteFile(path, raw, 0644)
}

// RunSetAccountSwitchThreshold updates and persists the account switch threshold percentage.
func RunSetAccountSwitchThreshold(pct int, isJSON bool) error {
	if pct < 1 || pct > 100 {
		return apperror.NewSimple("threshold percentage must be between 1 and 100", "E400")
	}
	st := loadAccountSwitchState()
	st.ThresholdPct = pct
	if err := saveAccountSwitchState(st); err != nil {
		return err
	}
	if isJSON {
		return printJSON(st)
	}
	fmt.Printf("%s✔ Account switch threshold set to %d%%%s\n", constants.ColorGreen, pct, constants.ColorReset)
	return nil
}

// RunAccountSwitchStatus prints the current account switch configuration and telemetry.
func RunAccountSwitchStatus(isJSON bool) error {
	st := loadAccountSwitchState()
	if isJSON {
		return printJSON(st)
	}
	renderAccountSwitchStatusBox(st)
	return nil
}

func renderAccountSwitchStatusBox(st AccountSwitchState) {
	fmt.Printf("\n  %s● Antigravity Account Switch Status%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Active Account:     %s (%d%% credit)\n", st.ActiveAccount, st.ActiveCreditPct)
	fmt.Printf("  Switch Threshold:   %d%% (default: %d%%)\n", st.ThresholdPct, DefaultAccountSwitchThreshold)
	fmt.Printf("  Last Backup Batch:  %s (%d project(s): %s)\n",
		st.LastBatchId, st.BackedProjectCount, formatNoticeProjectsList(st.BackedProjectNames))
	fmt.Printf("  Liveness Verified:  %v | Fast-Forward: %v\n\n", st.IsLivenessVerified, st.IsFastForwarded)
}
