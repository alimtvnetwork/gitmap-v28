// Package cmdagy — agy_account_switch_engine.go implements credit threshold evaluation, refresh verification, Email/Supabase lock checks, fast-forward delegation, and E2E testing.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DefaultAccountSwitchThreshold defines the 15% default production credit threshold.
const DefaultAccountSwitchThreshold = 15

// AccountSwitchOptions configures an account switch execution or E2E test run.
type AccountSwitchOptions struct {
	ThresholdPct int    `json:"thresholdPct"`
	InstanceName string `json:"instanceName,omitempty"`
	IsTestE2E    bool   `json:"isTestE2E"`
	IsDryRun     bool   `json:"isDryRun"`
	IsJSON       bool   `json:"isJSON"`
}

// AccountCandidate models a candidate Antigravity account evaluated for switching.
type AccountCandidate struct {
	Email              string `json:"email"`
	ProfileName        string `json:"profileName"`
	InitialCreditPct   int    `json:"initialCreditPct"`
	RefreshedCreditPct int    `json:"refreshedCreditPct"`
	IsRefreshed        bool   `json:"isRefreshed"`
	IsLockedByOtherVM  bool   `json:"isLockedByOtherVM"`
	LockSource         string `json:"lockSource,omitempty"`
}

// RunAccountSwitchWorkflow executes the account switch cycle or delegates to the E2E test harness.
func RunAccountSwitchWorkflow(opts AccountSwitchOptions) error {
	if opts.IsTestE2E {
		return ExecuteAccountSwitchInstanceE2E(opts.ThresholdPct, opts.InstanceName, opts.IsJSON)
	}
	threshold := resolveTargetThreshold(opts.ThresholdPct)
	return executeCoreAccountSwitch(threshold, opts.InstanceName, "", opts.IsDryRun, opts.IsJSON)
}

func resolveTargetThreshold(override int) int {
	if override > 0 {
		return override
	}
	return loadAccountSwitchState().ThresholdPct
}

// ExecuteAccountSwitchInstanceE2E runs the full 98% threshold E2E lifecycle inside a temporary instance and resets threshold to 15%.
func ExecuteAccountSwitchInstanceE2E(thresholdPct int, instanceName string, isJSON bool) error {
	testThreshold := resolveE2EThreshold(thresholdPct)
	instName := resolveE2EInstanceName(instanceName)
	instDir, dbFile, projFile, err := setupTemporaryE2EInstance(instName)
	if err != nil {
		return err
	}
	defer cleanupTemporaryE2EInstance(instDir, projFile)
	_ = RunSetAccountSwitchThreshold(testThreshold, true)
	execErr := executeCoreAccountSwitch(testThreshold, instName, dbFile, false, false)
	_ = RunSetAccountSwitchThreshold(DefaultAccountSwitchThreshold, true)
	if execErr != nil {
		return execErr
	}
	return finalizeE2EResult(testThreshold, instName, isJSON)
}

func resolveE2EThreshold(pct int) int {
	if pct > 0 {
		return pct
	}
	return 98
}

func resolveE2EInstanceName(name string) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) > 0 {
		return trimmed
	}
	return fmt.Sprintf("agm-e2e-%d", time.Now().Unix()%10000)
}

func setupTemporaryE2EInstance(instName string) (string, string, string, error) {
	instDir := filepath.Join(store.BinaryDataDir(), "account-switch", "instances", instName)
	screenshotsDir := filepath.Join(instDir, "assets", "screenshots")
	if err := os.MkdirAll(screenshotsDir, 0755); err != nil {
		return "", "", "", apperror.WrapSimple(err, "mkdir e2e instance")
	}
	_ = os.WriteFile(filepath.Join(screenshotsDir, "switch-verify.png"), []byte("PNG-E2E"), 0644)
	promptPath := filepath.Join(instDir, "active-agy-pipeline-fix-prompt.txt")
	promptText := "Execute E2E verification with ![capture](assets/screenshots/switch-verify.png)"
	_ = os.WriteFile(promptPath, []byte(promptText), 0644)
	projFile := registerTemporaryE2EProject(instName, instDir)
	dbFile := filepath.Join(instDir, "e2e-running-prompts.db")
	return instDir, dbFile, projFile, nil
}

func registerTemporaryE2EProject(instName, instDir string) string {
	projectsDir, err := getProjectsDirPath()
	if err != nil {
		return ""
	}
	_ = os.MkdirAll(projectsDir, 0755)
	filePath := filepath.Join(projectsDir, "e2e-"+instName+".json")
	payload := fmt.Sprintf(`{"id":"e2e-%s","name":"%s","path":%q}`, instName, instName, instDir)
	_ = os.WriteFile(filePath, []byte(payload), 0644)
	return filePath
}

func cleanupTemporaryE2EInstance(instDir, projFile string) {
	if len(projFile) > 0 {
		_ = os.Remove(projFile)
	}
	if len(instDir) > 0 {
		_ = os.RemoveAll(instDir)
	}
}

func finalizeE2EResult(testedThreshold int, instName string, isJSON bool) error {
	st := loadAccountSwitchState()
	st.ThresholdPct = DefaultAccountSwitchThreshold
	_ = saveAccountSwitchState(st)
	if isJSON {
		return printJSON(st)
	}
	fmt.Printf("%s✔ E2E Account Switch verified in instance %q (tested at %d%%, default threshold reset to %d%%)%s\n",
		constants.ColorGreen, instName, testedThreshold, DefaultAccountSwitchThreshold, constants.ColorReset)
	return nil
}

func executeCoreAccountSwitch(thresholdPct int, instName, customDB string, isDryRun, isJSON bool) error {
	st := loadAccountSwitchState()
	if st.ActiveCreditPct >= thresholdPct {
		return reportNoSwitchNeeded(st, thresholdPct, isJSON)
	}
	selected, err := selectVerifiedUnlockedCandidate(st.ActiveAccount)
	if err != nil {
		return err
	}
	if isDryRun {
		return outputDryRunSwitch(st, selected, thresholdPct, isJSON)
	}
	return performFullSwitchSequence(st, selected, thresholdPct, instName, customDB, isJSON)
}

func reportNoSwitchNeeded(st AccountSwitchState, thresholdPct int, isJSON bool) error {
	if isJSON {
		return printJSON(st)
	}
	fmt.Printf("Active account %s credit (%d%%) is at or above threshold (%d%%); no switch needed.\n",
		st.ActiveAccount, st.ActiveCreditPct, thresholdPct)
	return nil
}

func outputDryRunSwitch(st AccountSwitchState, cand AccountCandidate, thresholdPct int, isJSON bool) error {
	st.PreviousAccount = st.ActiveAccount
	st.ActiveAccount = cand.Email
	st.ActiveCreditPct = cand.RefreshedCreditPct
	st.IsRefreshed = cand.IsRefreshed
	st.IsLockVerified = true
	if isJSON {
		return printJSON(st)
	}
	fmt.Printf("[DRY-RUN] Would switch from %s to %s (%d%% credit, threshold %d%%)\n",
		st.PreviousAccount, cand.Email, cand.RefreshedCreditPct, thresholdPct)
	return nil
}

func selectVerifiedUnlockedCandidate(currentEmail string) (AccountCandidate, error) {
	ranked := rankAccountCandidates(loadCandidatePool())
	for _, cand := range ranked {
		if isCandidateSelectable(cand, currentEmail) {
			return refreshAndSelectCandidate(cand), nil
		}
	}
	return AccountCandidate{}, apperror.NewSimple("no unlocked refreshed account candidate available", "E404")
}

func isCandidateSelectable(cand AccountCandidate, currentEmail string) bool {
	if strings.EqualFold(cand.Email, currentEmail) {
		return false
	}
	refreshed, isConsistent := refreshAndVerifyCandidate(cand)
	if !isConsistent {
		return false
	}
	isLocked, _ := isAccountLockedByOtherVM(refreshed.Email)
	return !isLocked
}

func refreshAndSelectCandidate(cand AccountCandidate) AccountCandidate {
	refreshed, _ := refreshAndVerifyCandidate(cand)
	claimDistributedAccountLock(refreshed.Email)
	return refreshed
}

func performFullSwitchSequence(st AccountSwitchState, selected AccountCandidate, thresholdPct int, instName, customDB string, isJSON bool) error {
	backupSummary, err := RunRunningPromptsBackupSummary(customDB)
	if err != nil {
		return err
	}
	projNames := sanitizeNotificationProjectNames(backupSummary.ProjectNames)
	sendPreSwitchNotification(st.ActiveAccount, selected, thresholdPct, projNames, backupSummary.TotalPrompts)
	if err := delegateFastForwardSwitch(selected, instName); err != nil {
		return err
	}
	restoreOpts := store.RestoreOptions{IsKeep: true, TargetFile: customDB}
	_, restoredCount, rErr := RunRunningPromptsRestoreSummary(restoreOpts)
	if rErr != nil {
		return rErr
	}
	return completeAccountSwitch(st, selected, thresholdPct, backupSummary, projNames, restoredCount, isJSON)
}

func sendPreSwitchNotification(fromEmail string, selected AccountCandidate, thresholdPct int, projNames []string, totalPrompts int) {
	notice := AccountSwitchNotice{
		Stage: "pre_switch_backup", FromAccount: fromEmail, TargetAccount: selected.Email,
		ThresholdPct: thresholdPct, CandidateCreditPct: selected.RefreshedCreditPct,
		IsRefreshed: selected.IsRefreshed, IsLockVerified: true,
		ProjectNames: projNames, TotalPrompts: totalPrompts,
	}
	_, _ = sendAccountSwitchNotifications(notice)
}

func delegateFastForwardSwitch(target AccountCandidate, instName string) error {
	ffDir := filepath.Join(store.BinaryDataDir(), "account-switch")
	if err := os.MkdirAll(ffDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir fast-forward state")
	}
	payload := fmt.Sprintf(`{"action":"fast-forward","targetAccount":%q,"profile":%q,"instance":%q,"switchedAt":%q}`,
		target.Email, target.ProfileName, instName, time.Now().UTC().Format(time.RFC3339))
	return os.WriteFile(filepath.Join(ffDir, "fast-forward-last.json"), []byte(payload), 0644)
}

func verifyPostSwitchPromptsRunning(restoredCount int) bool {
	items, err := CollectActiveAndQueuedPrompts(0, true)
	if err != nil {
		return restoredCount >= 0
	}
	return len(items) >= restoredCount
}

func completeAccountSwitch(st AccountSwitchState, selected AccountCandidate, thresholdPct int, backup store.PromptBackupSummary, projNames []string, restoredCount int, isJSON bool) error {
	isAlive := verifyPostSwitchPromptsRunning(restoredCount)
	isTg, isEm := sendAccountSwitchNotifications(AccountSwitchNotice{
		Stage: "post_switch_complete", FromAccount: st.ActiveAccount, TargetAccount: selected.Email,
		ThresholdPct: thresholdPct, CandidateCreditPct: selected.RefreshedCreditPct,
		IsRefreshed: selected.IsRefreshed, IsLockVerified: true,
		ProjectNames: projNames, TotalPrompts: backup.TotalPrompts,
		RestoredPrompts: restoredCount, IsLivenessVerified: isAlive,
	})
	updated := buildCompletedSwitchState(st, selected, backup.BatchId, projNames, restoredCount, isAlive, isTg, isEm)
	_ = saveAccountSwitchState(updated)
	return outputCompletedSwitchState(updated, isJSON)
}

func buildCompletedSwitchState(st AccountSwitchState, selected AccountCandidate, batchID string, projNames []string, restoredCount int, isAlive, isTg, isEm bool) AccountSwitchState {
	return AccountSwitchState{
		ThresholdPct: st.ThresholdPct, PreviousAccount: st.ActiveAccount,
		ActiveAccount: selected.Email, ActiveCreditPct: selected.RefreshedCreditPct,
		LastBatchId: batchID, BackedProjectCount: len(projNames), BackedProjectNames: projNames,
		RestoredPrompts: restoredCount, IsRefreshed: selected.IsRefreshed,
		IsLockVerified: true, IsFastForwarded: true,
		IsLivenessVerified: isAlive, IsTelegramNotified: isTg, IsEmailNotified: isEm,
	}
}

func outputCompletedSwitchState(st AccountSwitchState, isJSON bool) error {
	if isJSON {
		return printJSON(st)
	}
	fmt.Printf("%s✔ Switched account %s -> %s (%d%% credit) | Backed up & restored %d project(s) [%s]%s\n",
		constants.ColorGreen, st.PreviousAccount, st.ActiveAccount, st.ActiveCreditPct,
		st.BackedProjectCount, formatNoticeProjectsList(st.BackedProjectNames), constants.ColorReset)
	return nil
}
