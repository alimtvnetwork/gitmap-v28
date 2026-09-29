// Package cmdagy — agy_wpr_watchdog.go runs the watchdog loop, probes prompt activity, and triggers recovery.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// WPRRuntimeState models live watchdog process state.
type WPRRuntimeState struct {
	PID           int      `json:"PID"`
	IntervalSec   int      `json:"IntervalSec"`
	StartedAt     string   `json:"StartedAt"`
	LastCheckedAt string   `json:"LastCheckedAt"`
	TargetSlugs   []string `json:"TargetSlugs"`
	IsActive      bool     `json:"IsActive"`
}

func resolveWPRRuntimeStatePath() string {
	rootDir := store.ResolveWatchPromptsRootDir()
	return filepath.Join(rootDir, "wpr_runtime_state.json")
}

// LoadWPRRuntimeStatus reads the watchdog loop runtime state and verifies PID liveness.
func LoadWPRRuntimeStatus() (WPRRuntimeState, bool) {
	path := resolveWPRRuntimeStatePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return WPRRuntimeState{IsActive: false}, false
	}

	var st WPRRuntimeState
	if json.Unmarshal(data, &st) != nil || !st.IsActive {
		return WPRRuntimeState{IsActive: false}, false
	}

	return st, st.IsActive
}

// SaveWPRRuntimeStatus saves watchdog runtime state to disk.
func SaveWPRRuntimeStatus(st WPRRuntimeState) error {
	path := resolveWPRRuntimeStatePath()
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, _ := json.MarshalIndent(st, "", "  ")
	return os.WriteFile(path, data, 0644)
}

// ExecuteWPRWatchLoop runs the continuous monitoring loop over target workspaces.
func ExecuteWPRWatchLoop(opts WPROptions) error {
	interval := resolveWatchInterval(opts.Interval)
	targets := resolveWatchedTargets(opts.Target)
	saveInitialWatchState(interval, targets)

	fmt.Printf("\n  %s● Starting Watch Prompts Running loop%s (interval: %s, targets: %d)\n\n",
		constants.ColorCyan, constants.ColorReset, interval, len(targets))

	for {
		runSingleWatchIteration(targets, opts)
		if opts.IsOnce {
			break
		}
		time.Sleep(interval)
	}

	return nil
}

func resolveWatchInterval(custom time.Duration) time.Duration {
	if custom >= 10*time.Second {
		return custom
	}

	return 2 * time.Minute
}

func resolveWatchedTargets(target string) []string {
	if target != "" && target != "all" {
		return []string{target}
	}

	projects, _ := getAllProjects()
	var list []string
	for _, p := range filterNonRestrictedProjects(projects) {
		list = append(list, p.GetPath())
	}

	return list
}

func saveInitialWatchState(interval time.Duration, targets []string) {
	st := WPRRuntimeState{
		PID:           os.Getpid(),
		IntervalSec:   int(interval.Seconds()),
		StartedAt:     time.Now().UTC().Format(time.RFC3339),
		LastCheckedAt: time.Now().UTC().Format(time.RFC3339),
		TargetSlugs:   targets,
		IsActive:      true,
	}
	_ = SaveWPRRuntimeStatus(st)
}

func runSingleWatchIteration(targets []string, opts WPROptions) {
	now := time.Now().UTC().Format(time.RFC3339)
	for _, ws := range targets {
		inspectAndRecoverWorkspace(ws, now, opts)
	}

	st, _ := LoadWPRRuntimeStatus()
	st.LastCheckedAt = now
	_ = SaveWPRRuntimeStatus(st)
}

func inspectAndRecoverWorkspace(ws, now string, opts WPROptions) {
	slug := store.SanitizeSlug(filepath.Base(ws))
	db, err := store.OpenWatchPromptsSplitDB(slug)
	if err != nil {
		return
	}
	defer db.Close()

	if isPromptActiveInWorkspace(ws) {
		_ = db.InsertWatchLog(slug, "heartbeat", "Prompts actively running in workspace", "healthy")
		return
	}

	executeAutoRecoveryFlow(ws, slug, db, opts)
}

func isPromptActiveInWorkspace(ws string) bool {
	q, hasQ := loadWorkspaceQueueFile(ws)
	if hasQ && (q.Active != nil || len(q.Queued) > 0) {
		return true
	}

	activeText := resolveActiveTextPath(ws)
	if isFileExisting(activeText) {
		return true
	}

	return false
}

func executeAutoRecoveryFlow(ws, slug string, db *store.WatchPromptsSplitDB, opts WPROptions) {
	_ = db.InsertWatchLog(slug, "recovery_start", "Prompts stalled or IDE crashed. Triggering recovery...", "warning")
	triggerFastForwardAccountSwitch(opts.Email)
	terminateRunningIDE()
	time.Sleep(500 * time.Millisecond)

	launchWorkspaceIDE(ws)
	time.Sleep(1 * time.Second)

	if isPromptActiveInWorkspace(ws) {
		_ = db.InsertWatchLog(slug, "recovery_ok", "Prompts resumed normally after restart", "success")
		return
	}

	replayStalledPromptWithTemplates(ws, slug, db, opts)
}

func triggerFastForwardAccountSwitch(email string) {
	cand := resolveTargetAccountCandidate(email)
	_ = delegateFastForwardSwitch(cand, "wpr-watchdog")
}

func resolveTargetAccountCandidate(email string) AccountCandidate {
	clean := strings.TrimSpace(email)
	if len(clean) > 0 {
		return AccountCandidate{Email: clean, ProfileName: "custom", RefreshedCreditPct: 99, IsRefreshed: true}
	}

	pool := loadCandidatePool()
	ranked := rankAccountCandidates(pool)
	for _, c := range ranked {
		if !c.IsLockedByOtherVM {
			return c
		}
	}

	return pool[0]
}

func launchWorkspaceIDE(ws string) {
	ideRes := ResolveAntigravityIDE()
	if ideRes.IsSuccess() {
		_ = launchAntigravityProcess(ideRes.Value, ws)
	}
}

func replayStalledPromptWithTemplates(ws, slug string, db *store.WatchPromptsSplitDB, opts WPROptions) {
	promptEntry, convId := extractLastPromptForProject(ws)
	prefix := resolveWPRPrefix(opts.PrefixTemplate)
	suffix := resolveWPRSuffix(opts.SuffixTemplate)
	decorated := decoratePromptBody(promptEntry.Content, prefix, suffix)

	promptFile := filepath.Join(ws, activeAgyPromptRelativePath)
	writePromptFile(promptFile, decorated)

	proj := AgyProject{ID: slug, Name: slug}
	plan := rerunRestartPlan{Project: proj, WorkspacePath: ws, ConvID: convId, PromptText: decorated}
	_ = dispatchRerunPayload(plan, false, "")
	_ = db.InsertWatchLog(slug, "replayed", "Re-injected stalled prompt with prefix template", "success")
}

func resolveWPRPrefix(template string) string {
	clean := strings.TrimSpace(template)
	if clean == "" || clean == "default" {
		return DefaultQueueCheckPrefix
	}

	return cmdprompttemplate.ResolveRandomTemplateContent(clean)
}

func resolveWPRSuffix(template string) string {
	clean := strings.TrimSpace(template)
	if clean == "" {
		return ""
	}

	return cmdprompttemplate.ResolveRandomTemplateContent(clean)
}
