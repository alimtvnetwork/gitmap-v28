// Package cmd — schedule_cmd.go: manages scheduled tasks, macro schedules,
// OS startup tasks, intervals, on-the-fly commands, testing, and execution.
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
	case "list", "ls", "status":
		return runScheduleList(rest)
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
	cmdStr := strings.Join(opts.Commands, " && ")
	task := store.SchedulerTask{
		Name:        opts.Name,
		MacroName:   opts.MacroName,
		CommandLine: cmdStr,
		IntervalVal: opts.Interval,
		DelayVal:    opts.Delay,
		IsScheduled: true,
		HasDelay:    opts.Delay != "",
		IsStartup:   opts.IsStartup,
	}
	if err := db.InsertSchedule(task); err != nil {
		return apperror.WrapSimple(err, "insert schedule")
	}
	handleStartupRegistration(opts.Name, opts.IsStartup)
	printScheduleAddSuccess(task)
	return nil
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
	steps := []string{"Schedule task: " + t.Name}
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
	fmt.Printf("✔ Scheduled task \033[1m%q\033[0m successfully (interval: %s)\n\n", t.Name, t.IntervalVal)
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
	fmt.Printf("  %-20s %-12s %-24s %-10s %-8s %s\n", "NAME", "INTERVAL", "TARGET (MACRO/CMD)", "STARTUP", "RUNS", "LAST RUN")
	fmt.Printf("  %s\n", strings.Repeat("─", 88))
	for _, t := range tasks {
		target := t.MacroName
		if target == "" {
			target = t.CommandLine
		}
		if len(target) > 22 {
			target = target[:19] + "..."
		}
		startup := "no"
		if t.IsStartup {
			startup = "yes"
		}
		fmt.Printf("  %-20s %-12s %-24s %-10s %-8d %s\n", t.Name, t.IntervalVal, target, startup, t.RunCount, t.LastRunAt)
	}
	fmt.Println()
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
	return executeScheduledTaskNow(db, t)
}

func executeScheduledTaskNow(db *store.DB, t *store.SchedulerTask) error {
	fmt.Printf("\n\033[1;96m▸ Executing scheduled task:\033[0m \033[1m%q\033[0m (interval: %s)\n", t.Name, t.IntervalVal)
	applyScheduleDelay(t.DelayVal)
	execErr := executeTaskTarget(t)
	_ = db.UpdateScheduleRun(t.Name, time.Now().Format("2006-01-02 15:04:05"))
	if execErr != nil {
		fmt.Printf("\n\033[1;91m✖ Task %q failed:\033[0m %v\n\n", t.Name, execErr)
		return execErr
	}
	fmt.Printf("\n\033[1;92m✔ Task %q completed successfully\033[0m\n\n", t.Name)
	return nil
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

func executeTaskTarget(t *store.SchedulerTask) error {
	if t.MacroName != "" {
		return executeMacroTarget(t.MacroName)
	}
	if t.CommandLine != "" {
		return runShellCmdDirectOnly(t.CommandLine)
	}
	return apperror.NewSimple("no macro or command defined for task", "E6003")
}

func executeMacroTarget(name string) error {
	m, err := macro.LoadMacro(name)
	if err != nil {
		return err
	}
	return macro.Execute(context.Background(), m, macro.ExecOptions{})
}

func runShellCmdDirectOnly(cmdStr string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
		if err := executeScheduledTaskNow(db, t); err != nil {
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
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.DeleteSchedule(args[0]); err != nil {
		return apperror.WrapSimple(err, "delete schedule")
	}
	fmt.Printf("✔ Removed scheduled task %q\n", args[0])
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
