// Package cmdagy — agy_wpr_task.go enqueues WPR actions into TaskHistory, TaskQueue, Scheduler, and SUG.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RecordWPRActionToTaskHistory records a WPR operation into SectionTasksDB TaskHistory and TaskQueue.
func RecordWPRActionToTaskHistory(action, target string, opts WPROptions) error {
	fwdPayload, _ := json.Marshal(opts)
	invOpts := buildInverseWPROptions(action, opts)
	invPayload, _ := json.Marshal(invOpts)

	err := recordAgyTask("wpr_"+action, target, string(fwdPayload), string(invPayload))
	if err != nil {
		return err
	}

	_ = enqueueWPRTaskInQueue(action, target, string(fwdPayload), string(invPayload))
	_ = syncWPRWithScheduler(action, target, opts)
	_ = syncWPRWithSUG(action, target)

	return nil
}

func buildInverseWPROptions(action string, opts WPROptions) WPROptions {
	inv := opts
	switch action {
	case "start":
		inv.Subcommand = "disable"
	case "disable":
		inv.Subcommand = "start"
	case "shutdown":
		inv.Subcommand = "start"
	case "switch-account", "fast-forward":
		inv.Subcommand = "switch-account"
	default:
		inv.Subcommand = "status"
	}

	return inv
}

func enqueueWPRTaskInQueue(action, target, fwd, inv string) error {
	db, err := store.OpenSectionTasksDB("agy", "")
	if err != nil {
		return err
	}
	defer db.Close()

	queueId := fmt.Sprintf("wpr-q-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO TaskQueue (QueueId, Section, Action, Target, ForwardPayload, InversePayload, Status, CreatedAt, UpdatedAt)
		VALUES (?, 'agy', ?, ?, ?, ?, 'completed', ?, ?)`

	_ = store.ExecWrapper(db.Conn(), query, queueId, "wpr_"+action, target, fwd, inv, now, now)

	return nil
}

func syncWPRWithScheduler(action, target string, opts WPROptions) error {
	slug := store.SanitizeSlug(target)
	if slug == "" || slug == "default" {
		slug = "global-wpr"
	}

	schedDB, err := store.OpenScheduleSplitDB("wpr-" + slug)
	if err != nil {
		return err
	}
	defer schedDB.Close()

	interval := opts.IntervalStr
	if interval == "" {
		interval = "2m"
	}

	isEnabled := action == "start" || action == "all"
	cfg := store.ScheduleConfig{
		Name:        "Watch Prompts Running (" + target + ")",
		Slug:        "wpr-" + slug,
		CommandLine: "gitmap wpr start " + target,
		IntervalVal: interval,
		IsEnabled:   isEnabled,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	return schedDB.SaveConfig(cfg)
}

func syncWPRWithSUG(action, target string) error {
	if action != "start" && action != "all" {
		return nil
	}

	cfg := loadSUGConfig()
	if target != "" && target != "all" {
		resolved, isValid := ValidateSUGTarget(target)
		if isValid {
			appendNewSUGTargets(&cfg, []string{resolved})
			_ = saveSUGConfig(cfg)
		}
	}

	return nil
}
