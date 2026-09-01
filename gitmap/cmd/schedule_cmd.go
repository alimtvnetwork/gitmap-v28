// Package cmd — schedule_cmd.go: manages scheduled tasks, macro schedules,
// OS startup tasks, isolated split SQLite databases, run history, and logs.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/osutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runSchedule handles gitmap schedule commands.
func runSchedule(args []string) error {
	checkHelp("schedule", args)
	if len(args) == 0 {
		return runScheduleList(nil)
	}
	return dispatchScheduleSubcommand(args[0], args[1:])
}

func dispatchScheduleSubcommand(sub string, rest []string) error {
	switch sub {
	case "add", "create", "new":
		return runScheduleAdd(rest)
	case "list", "ls":
		return runScheduleList(rest)
	case "status":
		return runScheduleStatus(rest)
	case "export", "export-all":
		return runScheduleExport(rest)
	case "import", "import-all":
		return runScheduleImport(rest)
	case "enable":
		return runScheduleSetEnabled(rest, true)
	case "disable":
		return runScheduleSetEnabled(rest, false)
	case "logs", "log", "history":
		return runScheduleLogs(rest)
	case "reset":
		return runScheduleReset(rest)
	case "reset-all":
		return runScheduleResetAll(rest)
	case "run", "exec":
		return runScheduleRun(rest)
	case "test":
		return runScheduleTest(rest)
	case "rm", "delete", "del":
		return runScheduleDelete(rest)
	case "startup":
		return runScheduleStartup(rest)
	case "restart":
		return runScheduleRestart()
	case "shutdown":
		return runScheduleShutdown()
	default:
		return runScheduleAdd(append([]string{sub}, rest...))
	}
}

type scheduleAddOpts struct {
	Name      string
	MacroName string
	Commands  []string
	Interval  string
	Delay     string
	IsStartup bool
}

func parseScheduleAddOpts(args []string) scheduleAddOpts {
	var opts scheduleAddOpts
	for i := 0; i < len(args); i++ {
		a := args[i]
		if matchScheduleFlag(a, &opts, &i, args) {
			continue
		}
		if !strings.HasPrefix(a, "-") && opts.Name == "" {
			opts.Name = a
			continue
		}
		if !strings.HasPrefix(a, "-") {
			opts.Commands = append(opts.Commands, a)
		}
	}
	return opts
}

func matchScheduleFlag(a string, opts *scheduleAddOpts, idx *int, args []string) bool {
	if matchFlagWithVal(a, "--macro", "-m") {
		opts.MacroName = extractFlagValue(idx, args)
		return true
	}
	if matchFlagWithVal(a, "--every", "--interval", "-i") {
		opts.Interval = extractFlagValue(idx, args)
		return true
	}
	if matchFlagWithVal(a, "--delay", "--sleep", "-d") {
		opts.Delay = extractFlagValue(idx, args)
		return true
	}
	if matchScheduleTimeUnits(a, opts, idx, args) {
		return true
	}
	if a == "--startup" {
		opts.IsStartup = true
		return true
	}
	return false
}

func matchScheduleTimeUnits(a string, opts *scheduleAddOpts, idx *int, args []string) bool {
	if matchFlagWithVal(a, "--day", "--days") {
		opts.Interval = extractFlagValue(idx, args) + "d"
		return true
	}
	if matchFlagWithVal(a, "--hour", "--hours") {
		opts.Interval = extractFlagValue(idx, args) + "h"
		return true
	}
	if matchFlagWithVal(a, "--minute", "--minutes", "--min") {
		opts.Interval = extractFlagValue(idx, args) + "m"
		return true
	}
	if matchFlagWithVal(a, "--second", "--seconds", "--sec") {
		opts.Interval = extractFlagValue(idx, args) + "s"
		return true
	}
	return false
}

func runScheduleAdd(args []string) error {
	opts := parseScheduleAddOpts(args)
	if opts.Name == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule add <name> [commands...] [--macro <name>] [--every <1d|2h|30m|15s>] [--delay <10s>] [--startup]\n")
		return apperror.NewSimple("schedule name required", "E6001")
	}
	if opts.Interval == "" {
		opts.Interval = "1h"
	}
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return saveScheduleTask(db, opts)
}

func saveScheduleTask(db *store.DB, opts scheduleAddOpts) error {
	slug := store.ScheduleSlug(opts.Name)
	cmdStr := strings.Join(opts.Commands, " && ")
	task := store.SchedulerTask{
		Name:        opts.Name,
		Slug:        slug,
		DBPath:      store.ScheduleDBPath(slug),
		MacroName:   opts.MacroName,
		CommandLine: cmdStr,
		IntervalVal: opts.Interval,
		DelayVal:    opts.Delay,
		IsEnabled:   true,
		IsScheduled: true,
		HasDelay:    opts.Delay != "",
		IsStartup:   opts.IsStartup,
	}
	if err := db.InsertSchedule(task); err != nil {
		return apperror.WrapSimple(err, "insert schedule in root db")
	}
	_ = syncScheduleSplitDBConfig(task)
	handleStartupRegistration(opts.Name, opts.IsStartup)
	printScheduleAddSuccess(task)
	return nil
}

func syncScheduleSplitDBConfig(t store.SchedulerTask) error {
	splitDB, err := store.OpenScheduleSplitDB(t.Slug)
	if err != nil {
		return err
	}
	defer splitDB.Close()
	return splitDB.SaveConfig(store.ScheduleConfig{
		Name:        t.Name,
		Slug:        t.Slug,
		MacroName:   t.MacroName,
		CommandLine: t.CommandLine,
		IntervalVal: t.IntervalVal,
		DelayVal:    t.DelayVal,
		IsEnabled:   t.IsEnabled,
		IsStartup:   t.IsStartup,
	})
}

func handleStartupRegistration(name string, isStartup bool) {
	if !isStartup {
		return
	}
	exePath, err := os.Executable()
	if err == nil {
		_ = osutil.AddToStartup(exePath + " schedule run " + name)
	}
}

func printScheduleAddSuccess(t store.SchedulerTask) {
	steps := []string{"Schedule task: " + t.Name, "Split DB: " + t.Slug + ".db"}
	if t.MacroName != "" {
		steps = append(steps, "Linked macro: "+t.MacroName)
	}
	if t.CommandLine != "" {
		steps = append(steps, "Command: "+t.CommandLine)
	}
	if t.DelayVal != "" {
		steps = append(steps, "Initial delay/sleep: "+t.DelayVal)
	}
	if t.IsStartup {
		steps = append(steps, "OS Startup: enabled")
	}
	printScheduleSummaryTree(t.Name, t.IntervalVal, "schedule", steps)
	fmt.Printf("✔ Scheduled task \033[1m%q\033[0m successfully (interval: %s, split db: %s)\n\n", t.Name, t.IntervalVal, t.DBPath)
}

func runScheduleSetEnabled(args []string, isEnabled bool) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule enable|disable <name>\n")
		return apperror.NewSimple("schedule name required", "E6007")
	}
	name := args[0]
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(name)
	if err != nil {
		return apperror.WrapSimple(err, "get schedule "+name)
	}
	_ = db.SetScheduleEnabled(name, isEnabled)
	t.IsEnabled = isEnabled
	_ = syncScheduleSplitDBConfig(*t)
	stateStr := "enabled"
	if !isEnabled {
		stateStr = "disabled"
	}
	fmt.Printf("✔ Schedule \033[1m%q\033[0m is now %s\n", name, stateStr)
	return nil
}

func runScheduleList(args []string) error {
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	tasks, err := db.ListSchedules()
	if err != nil {
		return apperror.WrapSimple(err, "list schedules")
	}
	opts := parseExecOptions(args)
	if opts.JSON || opts.YAML || len(opts.FilePath) > 0 {
		return outputStructuredData(tasks, opts)
	}
	renderScheduleTable(tasks)
	return nil
}

func renderScheduleTable(tasks []store.SchedulerTask) {
	if len(tasks) == 0 {
		fmt.Println("  No scheduled tasks found. Create one with: gitmap schedule add <name> --every <interval>")
		return
	}
	fmt.Println()
	fmt.Printf("  %-18s %-10s %-10s %-20s %-8s %-6s %s\n", "NAME", "STATUS", "INTERVAL", "TARGET (MACRO/CMD)", "STARTUP", "RUNS", "SPLIT DB")
	fmt.Printf("  %s\n", strings.Repeat("─", 88))
	for _, t := range tasks {
		target := t.MacroName
		if target == "" {
			target = t.CommandLine
		}
		if len(target) > 18 {
			target = target[:15] + "..."
		}
		status := "\033[32menabled\033[0m"
		if !t.IsEnabled {
			status = "\033[31mdisabled\033[0m"
		}
		startup := "no"
		if t.IsStartup {
			startup = "yes"
		}
		fmt.Printf("  %-18s %-19s %-10s %-20s %-8s %-6d %s.db\n", t.Name, status, t.IntervalVal, target, startup, t.RunCount, t.Slug)
	}
	fmt.Println()
}

func runScheduleStatus(args []string) error {
	name, flagArgs := extractMacroNameAndFlags(args)
	if name == "" || name == "*" || name == "all" {
		return runScheduleList(args)
	}
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(name)
	if err != nil {
		return apperror.WrapSimple(err, "get schedule "+name)
	}
	opts := parseExecOptions(flagArgs)
	if opts.JSON || opts.YAML || len(opts.FilePath) > 0 {
		return outputStructuredData(t, opts)
	}
	renderSingleScheduleStatus(t)
	return nil
}

func renderSingleScheduleStatus(t *store.SchedulerTask) {
	fmt.Printf("\n  \033[1;96mScheduled Task Status:\033[0m \033[1m%q\033[0m\n", t.Name)
	fmt.Printf("    • Status:      %s\n", formatTaskEnabled(t.IsEnabled))
	fmt.Printf("    • Interval:    %s\n", t.IntervalVal)
	if t.DelayVal != "" {
		fmt.Printf("    • Delay:       %s\n", t.DelayVal)
	}
	if t.MacroName != "" {
		fmt.Printf("    • Macro:       %s\n", t.MacroName)
	}
	if t.CommandLine != "" {
		fmt.Printf("    • Command:     %s\n", t.CommandLine)
	}
	fmt.Printf("    • Startup:     %v\n", t.IsStartup)
	fmt.Printf("    • Total Runs:  %d\n", t.RunCount)
	fmt.Printf("    • Last Run:    %s\n", t.LastRunAt)
	fmt.Printf("    • Split DB:    %s\n\n", t.DBPath)
}

func formatTaskEnabled(isEnabled bool) string {
	if isEnabled {
		return "\033[32menabled\033[0m"
	}
	return "\033[31mdisabled\033[0m"
}

func runScheduleLogs(args []string) error {
	name, flagArgs := extractMacroNameAndFlags(args)
	if name == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule logs <name> [--limit <N>] [--json] [--yaml] [-f <path>]\n")
		return apperror.NewSimple("schedule name required", "E6008")
	}
	opts := parseExecOptions(flagArgs)
	limit := parseLogsLimit(flagArgs)
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(name)
	if err != nil {
		return apperror.WrapSimple(err, "get schedule "+name)
	}
	return renderScheduleLogsFromSplitDB(t, limit, opts)
}

func parseLogsLimit(args []string) int {
	for i := 0; i < len(args); i++ {
		if !matchFlagWithVal(args[i], "--limit", "-n") {
			continue
		}
		val, _ := strconv.Atoi(extractFlagValue(&i, args))
		if val > 0 {
			return val
		}
	}
	return 20
}

func renderScheduleLogsFromSplitDB(t *store.SchedulerTask, limit int, opts macro.ExecOptions) error {
	splitDB, err := store.OpenScheduleSplitDB(t.Slug)
	if err != nil {
		return apperror.WrapSimple(err, "open schedule split db")
	}
	defer splitDB.Close()
	runs, err := splitDB.GetRuns(limit)
	if err != nil {
		return apperror.WrapSimple(err, "fetch logs from split db")
	}
	if opts.JSON || opts.YAML || len(opts.FilePath) > 0 {
		return outputStructuredData(runs, opts)
	}
	renderScheduleRunsTable(t.Name, runs)
	return nil
}

func renderScheduleRunsTable(taskName string, runs []store.ScheduleRunRecord) {
	fmt.Printf("\n  \033[1;96mExecution Logs for Schedule:\033[0m \033[1m%q\033[0m (%d record(s))\n\n", taskName, len(runs))
	if len(runs) == 0 {
		fmt.Println("  (no execution logs recorded yet)")
		return
	}
	fmt.Printf("  %-6s %-19s %-12s %-10s %-8s %s\n", "RUN #", "STARTED AT", "USER", "DURATION", "STATUS", "EXIT")
	fmt.Printf("  %s\n", strings.Repeat("─", 72))
	for _, r := range runs {
		status := "\033[32msuccess\033[0m"
		if !r.IsSuccess {
			status = "\033[31mfailed\033[0m"
		}
		user := r.RunnerUser
		if user == "" {
			user = "system"
		}
		dur := fmt.Sprintf("%dms", r.DurationMS)
		fmt.Printf("  #%-5d %-19s %-12s %-10s %-17s %d\n", r.RunNumber, r.StartedAt, user, dur, status, r.ExitCode)
	}
	fmt.Println()
}

func runScheduleReset(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule reset <name>\n")
		return apperror.NewSimple("schedule name required", "E6009")
	}
	name := args[0]
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(name)
	if err != nil {
		return apperror.WrapSimple(err, "get schedule "+name)
	}
	splitDB, err := store.OpenScheduleSplitDB(t.Slug)
	if err == nil {
		_ = splitDB.ResetLogs()
		_ = splitDB.Close()
	}
	_ = db.UpdateScheduleRun(name, "")
	fmt.Printf("✔ Reset split database logs for schedule \033[1m%q\033[0m (%s.db)\n", name, t.Slug)
	return nil
}

func runScheduleResetAll(args []string) error {
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	tasks, _ := db.ListSchedules()
	for _, t := range tasks {
		splitDB, err := store.OpenScheduleSplitDB(t.Slug)
		if err == nil {
			_ = splitDB.ResetLogs()
			_ = splitDB.Close()
		}
		_ = db.UpdateScheduleRun(t.Name, "")
	}
	fmt.Printf("✔ Reset logs for all %d scheduled task split database(s)\n", len(tasks))
	return nil
}

func runScheduleRun(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule run <name>\n")
		return apperror.NewSimple("schedule name required", "E6002")
	}
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "get schedule "+args[0])
	}
	if !t.IsEnabled {
		fmt.Printf("⚠ Warning: schedule %q is currently disabled. Use 'gitmap schedule enable %s' to enable.\n", t.Name, t.Name)
	}
	return executeAndRecordScheduledTask(db, t, "manual")
}

func executeAndRecordScheduledTask(db *store.DB, t *store.SchedulerTask, triggerType string) error {
	fmt.Printf("\n\033[1;96m▸ Executing scheduled task:\033[0m \033[1m%q\033[0m (interval: %s)\n", t.Name, t.IntervalVal)
	applyScheduleDelay(t.DelayVal)
	start := time.Now()
	startedAt := start.Format("2006-01-02 15:04:05")
	currentUser := resolveCurrentUser()
	execErr, outputStr, exitCode := executeTaskTargetWithOutput(t)
	durMS := time.Since(start).Milliseconds()
	finishedAt := time.Now().Format("2006-01-02 15:04:05")

	recordTaskRunInSplitDB(t.Slug, t.RunCount+1, triggerType, currentUser, startedAt, finishedAt, durMS, execErr == nil, exitCode, outputStr, execErr)
	_ = db.UpdateScheduleRun(t.Name, finishedAt)
	if execErr != nil {
		fmt.Printf("\n\033[1;91m✖ Task %q failed:\033[0m %v\n\n", t.Name, execErr)
		return execErr
	}
	fmt.Printf("\n\033[1;92m✔ Task %q completed successfully\033[0m (logged to %s.db)\n\n", t.Name, t.Slug)
	return nil
}

func recordTaskRunInSplitDB(slug string, runNum int, triggerType, user, startAt, finishAt string, durMS int64, isSuccess bool, exitCode int, out string, execErr error) {
	splitDB, err := store.OpenScheduleSplitDB(slug)
	if err != nil {
		return
	}
	defer splitDB.Close()
	errMsg := ""
	if execErr != nil {
		errMsg = execErr.Error()
	}
	_ = splitDB.RecordRun(store.ScheduleRunRecord{
		RunNumber:   runNum,
		TriggerType: triggerType,
		RunnerUser:  user,
		StartedAt:   startAt,
		FinishedAt:  finishAt,
		DurationMS:  durMS,
		IsSuccess:   isSuccess,
		ExitCode:    exitCode,
		Output:      out,
		ErrorMsg:    errMsg,
	})
}

func resolveCurrentUser() string {
	u := os.Getenv("USERNAME")
	if u == "" {
		u = os.Getenv("USER")
	}
	if u == "" {
		u = "runner"
	}
	return u
}

func applyScheduleDelay(delayVal string) {
	if delayVal == "" {
		return
	}
	d := parseDurationArg(delayVal, 0)
	if d <= 0 {
		return
	}
	fmt.Printf("  ⏳ Applying delay of %v...\n", d)
	time.Sleep(d)
}

func executeTaskTargetWithOutput(t *store.SchedulerTask) (error, string, int) {
	if t.MacroName != "" {
		return executeMacroTargetWithOutput(t.MacroName)
	}
	if t.CommandLine != "" {
		return runShellCmdWithCapture(t.CommandLine)
	}
	return apperror.NewSimple("no macro or command defined for task", "E6003"), "", 1
}

func executeMacroTargetWithOutput(macroName string) (error, string, int) {
	m, err := macro.LoadMacro(macroName)
	if err != nil {
		return err, "", 1
	}
	execErr := macro.Execute(context.Background(), m, macro.ExecOptions{})
	if execErr != nil {
		return execErr, "", 1
	}
	return nil, "macro " + macroName + " executed", 0
}

func runShellCmdWithCapture(cmdStr string) (error, string, int) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Print(string(out))
		return nil, string(out), 0
	}
	exitCode := 1
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}
	fmt.Print(string(out))
	return err, string(out), exitCode
}

func runScheduleTest(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule test <name> [--delay 1s] [--times <N>]\n")
		return apperror.NewSimple("schedule name required", "E6004")
	}
	times := parseTestTimes(args)
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "get schedule "+args[0])
	}
	return runTestIterations(db, t, times)
}

func parseTestTimes(args []string) int {
	for i := 0; i < len(args); i++ {
		if !matchFlagWithVal(args[i], "--times", "-n") {
			continue
		}
		val, _ := strconv.Atoi(extractFlagValue(&i, args))
		if val > 0 {
			return val
		}
	}
	return 1
}

func runTestIterations(db *store.DB, t *store.SchedulerTask, times int) error {
	fmt.Printf("\n\033[1;93m🧪 Testing schedule %q (%d iteration(s))...\033[0m\n", t.Name, times)
	for i := 1; i <= times; i++ {
		fmt.Printf("\n--- Test Iteration #%d/%d ---", i, times)
		if err := executeAndRecordScheduledTask(db, t, "test"); err != nil {
			return err
		}
	}
	fmt.Printf("\033[1;92m✔ All %d test iteration(s) passed for %q\033[0m\n\n", times, t.Name)
	return nil
}

func runScheduleDelete(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule rm <name>\n")
		return apperror.NewSimple("schedule name required", "E6005")
	}
	name := args[0]
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	t, err := db.GetSchedule(name)
	if err == nil && t != nil {
		_ = store.DeleteScheduleSplitDB(t.Slug)
	}
	if err := db.DeleteSchedule(name); err != nil {
		return apperror.WrapSimple(err, "delete schedule")
	}
	fmt.Printf("✔ Removed scheduled task %q and deleted split database\n", name)
	return nil
}

func runScheduleStartup(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule startup <name> [--enable|--disable]\n")
		return apperror.NewSimple("schedule name required", "E6006")
	}
	exePath, _ := os.Executable()
	cmdStr := exePath + " schedule run " + args[0]
	if err := osutil.AddToStartup(cmdStr); err != nil {
		return apperror.WrapSimple(err, "register startup")
	}
	fmt.Printf("✔ Registered %q for OS startup\n", args[0])
	return nil
}

func openSchedulerDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open db")
	}
	if err := db.InitSchedulerTable(); err != nil {
		db.Close()
		return nil, apperror.WrapSimple(err, "init scheduler table")
	}
	return db, nil
}
