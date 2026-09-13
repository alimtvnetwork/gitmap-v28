package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/startup"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type startupAddOpts struct {
	target    string
	name      string
	frequency string
	icon      string
	args      string
	desc      string
}

func parseStartupAddArgs(args []string) (*startupAddOpts, error) {
	if len(args) == 0 {
		return nil, apperror.New("startup", "E_ARG_REQUIRED", map[string]any{
			"msg": "Target script, executable, or macro required: gitmap startup add <target> [flags]",
		})
	}

	opts := &startupAddOpts{target: args[0], frequency: "everytime"}
	parseStartupFlags(args[1:], opts)

	if opts.name == "" {
		opts.name = deriveStartupName(opts.target)
	}

	return opts, nil
}

func parseStartupFlags(flagArgs []string, opts *startupAddOpts) {
	for i := 0; i < len(flagArgs); i++ {
		arg := flagArgs[i]
		if parseStartupKV(arg, opts) {
			continue
		}

		if i+1 < len(flagArgs) && parseStartupVal(arg, flagArgs[i+1], opts) {
			i++
		}
	}
}

func parseStartupKV(arg string, opts *startupAddOpts) bool {
	switch {
	case strings.HasPrefix(arg, "--frequency="):
		opts.frequency = strings.TrimPrefix(arg, "--frequency=")
		return true
	case strings.HasPrefix(arg, "--name="):
		opts.name = strings.TrimPrefix(arg, "--name=")
		return true
	case strings.HasPrefix(arg, "--icon="):
		opts.icon = strings.TrimPrefix(arg, "--icon=")
		return true
	case strings.HasPrefix(arg, "--args="):
		opts.args = strings.TrimPrefix(arg, "--args=")
		return true
	case strings.HasPrefix(arg, "--desc="):
		opts.desc = strings.TrimPrefix(arg, "--desc=")
		return true
	}

	return false
}

func parseStartupVal(flag, val string, opts *startupAddOpts) bool {
	switch flag {
	case "--frequency", "-f":
		opts.frequency = val
		return true
	case "--name", "-n":
		opts.name = val
		return true
	case "--icon":
		opts.icon = val
		return true
	case "--args":
		opts.args = val
		return true
	case "--desc":
		opts.desc = val
		return true
	}

	return false
}

func deriveStartupName(target string) string {
	clean := strings.TrimPrefix(target, "macro:")
	base := filepath.Base(clean)
	ext := filepath.Ext(base)
	if ext != "" {
		return strings.TrimSuffix(base, ext)
	}

	return base
}

func buildStartupRecord(opts *startupAddOpts) *store.StartupItemRecord {
	targetType := detectTargetType(opts.target)

	return &store.StartupItemRecord{
		Name:         opts.name,
		TargetType:   targetType,
		TargetPath:   opts.target,
		CommandArgs:  opts.args,
		IconPath:     opts.icon,
		RunFrequency: opts.frequency,
		IsActive:     true,
		Description:  opts.desc,
	}
}

func detectTargetType(target string) string {
	if strings.HasPrefix(target, "macro:") {
		return "macro"
	}

	ext := strings.ToLower(filepath.Ext(target))
	switch ext {
	case ".ps1":
		return "powershell"
	case ".sh":
		return "bash"
	case ".ico", ".png":
		return "icon"
	case ".bat", ".cmd":
		return "batch"
	case ".exe":
		return "binary"
	default:
		return "command"
	}
}

func registerNativeOSAutostart(rec *store.StartupItemRecord) {
	selfExe, err := os.Executable()
	if err != nil {
		return
	}

	execCmd := fmt.Sprintf("%q startup run %s", selfExe, rec.Name)
	opts := startup.AddOptions{
		Name:        "gitmap-" + rec.Name,
		DisplayName: "Gitmap Startup: " + rec.Name,
		Exec:        execCmd,
		Force:       true,
		Comment:     rec.Description,
	}
	_, _ = startup.Add(opts)
}

func unregisterNativeOSAutostart(name string) {
	_, _ = startup.Remove("gitmap-" + name)
	_, _ = startup.Remove(name)
}

func runStartupExec(args []string) error {
	if len(args) == 0 {
		return apperror.New("startup", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify startup item name to run: gitmap startup run <name>",
		})
	}

	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec, err := db.GetStartupItem(args[0])
	if err != nil {
		return err
	}
	if rec == nil {
		return apperror.New("startup", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Startup item %q not found in registry.", args[0]),
		})
	}

	if shouldSkipWeeklyRun(db, rec) {
		fmt.Printf("Startup item %q (once-a-week) already ran recently. Skipping.\n", rec.Name)

		return nil
	}

	return executeStartupItemAndLog(db, rec)
}

func shouldSkipWeeklyRun(db *store.StartupSplitDB, rec *store.StartupItemRecord) bool {
	if rec.RunFrequency != "once-a-week" {
		return false
	}

	logs, err := db.ListStartupLogs(rec.StartupItemId, 1)
	if err != nil || len(logs) == 0 {
		return false
	}

	const weekSeconds = 7 * 86400
	hasRecentRun := time.Now().Unix()-logs[0].RunAt < weekSeconds && logs[0].IsSuccess

	return hasRecentRun
}

func executeStartupItemAndLog(db *store.StartupSplitDB, rec *store.StartupItemRecord) error {
	start := time.Now()
	exitCode, outSummary, execErr := invokeStartupTarget(rec)
	duration := time.Since(start).Milliseconds()

	logRec := &store.StartupLogRecord{
		StartupItemId: rec.StartupItemId,
		RunAt:         time.Now().Unix(),
		DurationMs:    duration,
		IsSuccess:     execErr == nil,
		ExitCode:      exitCode,
		OutputSummary: outSummary,
		Notes:         rec.RunFrequency,
	}
	_ = db.RecordStartupLog(logRec)
	if execErr != nil {
		return execErr
	}

	fmt.Printf("%s Startup item %q executed successfully (%d ms)\n",
		constants.ColorGreen+"✓"+constants.ColorReset, rec.Name, duration)

	return nil
}

func invokeStartupTarget(rec *store.StartupItemRecord) (int, string, error) {
	if rec.TargetType == "macro" || strings.HasPrefix(rec.TargetPath, "macro:") {
		macroName := strings.TrimPrefix(rec.TargetPath, "macro:")

		return invokeMacroTarget(macroName)
	}

	return invokeScriptTarget(rec)
}

func invokeScriptTarget(rec *store.StartupItemRecord) (int, string, error) {
	var cmd *exec.Cmd
	switch rec.TargetType {
	case "powershell":
		cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", rec.TargetPath)
	case "bash":
		cmd = exec.Command("bash", rec.TargetPath)
	default:
		if rec.CommandArgs != "" {
			cmd = exec.Command(rec.TargetPath, strings.Fields(rec.CommandArgs)...)
		} else {
			cmd = exec.Command(rec.TargetPath)
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 1, err.Error(), err
	}

	return 0, "OK", nil
}

func invokeMacroTarget(macroName string) (int, string, error) {
	selfExe, err := os.Executable()
	if err != nil {
		selfExe = "gitmap"
	}

	cmd := exec.Command(selfExe, "execute", macroName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		return 1, runErr.Error(), runErr
	}

	return 0, "macro finished", nil
}

func runStartupLogs(args []string) error {
	if len(args) == 0 {
		return apperror.New("startup", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify startup item name to view logs: gitmap startup logs <name>",
		})
	}

	name := args[0]
	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec, err := db.GetStartupItem(name)
	if err != nil || rec == nil {
		return apperror.New("startup", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Startup item %q not found.", name),
		})
	}

	logs, err := db.ListStartupLogs(rec.StartupItemId, 20)
	if err != nil {
		return err
	}

	printStartupLogs(rec.Name, logs)

	return nil
}

func printStartupLogs(name string, logs []store.StartupLogRecord) {
	fmt.Printf("\n  Recent Startup Logs for %s (%d entries):\n", name, len(logs))
	fmt.Printf("  %-20s %-10s %-10s %s\n", "RUN AT", "DURATION", "STATUS", "SUMMARY")
	fmt.Println("  ─────────────────────────────────────────────────────────────────────────────")
	for _, l := range logs {
		tStr := time.Unix(l.RunAt, 0).Format("2006-01-02 15:04:05")
		statusStr := "SUCCESS"
		if l.IsSuccess == false {
			statusStr = fmt.Sprintf("FAIL(%d)", l.ExitCode)
		}
		fmt.Printf("  %-20s %-10s %-10s %s\n", tStr, fmt.Sprintf("%dms", l.DurationMs), statusStr, l.OutputSummary)
	}
	fmt.Println()
}
