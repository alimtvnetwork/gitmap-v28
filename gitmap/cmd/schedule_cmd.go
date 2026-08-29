package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runSchedule handles gitmap schedule commands.
func runSchedule(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		runScheduleAdd(args)
		return nil
	}
	dispatchSchedule(args[0], args)
	return nil
}

func dispatchSchedule(commandName string, args []string) {
	switch commandName {
	case "status":
		runScheduleStatus()
	case "restart":
		runScheduleRestart()
	case "shutdown":
		runScheduleShutdown()
	default:
		runScheduleAdd(args)
	}
}

// runScheduleStatus prints out current schedules.
func runScheduleStatus() error {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer db.Close()
	tasks, _ := db.ListSchedules()
	for _, t := range tasks {
		fmt.Printf("Task %s: interval=%s delay=%s\n", t.Name, t.IntervalVal, t.DelayVal)
	}
	return nil
}

// parseScheduleArgs parses CLI arguments.
func parseScheduleArgs(args []string) (string, string, string) {
	var name, interval, delay string
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--interval" && i+1 < len(args):
			interval, i = args[i+1], i+1
		case args[i] == "--delay" && i+1 < len(args):
			delay, i = args[i+1], i+1
		case !strings.HasPrefix(args[i], "--") && name == "":
			name = args[i]
		}
	}
	return name, interval, delay
}

// promptIfEmpty prompts for a value if it is empty.
func promptIfEmpty(prompt, val string) string {
	if val == "" {
		fmt.Print(prompt)
		reader := bufio.NewReader(os.Stdin)
		val, _ = reader.ReadString('\n')
		return strings.TrimSpace(val)
	}
	return val
}

// runScheduleAdd handles adding a new schedule.
func runScheduleAdd(args []string) error {
	name, interval, delay := parseScheduleArgs(args)
	name = promptIfEmpty("Name: ", name)
	interval = promptIfEmpty("Interval: ", interval)
	delay = promptIfEmpty("Delay: ", delay)

	isScheduled := interval != ""
	hasDelay := delay != ""

	steps := buildScheduleStepList(name, delay)
	printScheduleSummaryTree(name, interval, "default", steps)

	saveSchedule(name, interval, delay, isScheduled, hasDelay)
	return nil
}

// buildScheduleStepList creates step descriptions for schedule confirmation.
func buildScheduleStepList(taskName, delayVal string) []string {
	stepItems := []string{"Execute task: " + taskName}
	hasDelay := delayVal != ""
	if hasDelay {
		stepItems = append(stepItems, "Initial delay: "+delayVal)
	}
	return stepItems
}

// saveSchedule saves the parsed schedule to the DB.
func saveSchedule(n, i, d string, isScheduled, hasDelay bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()
	_ = db.InitSchedulerTable()
	if err := db.InsertSchedule(n, i, d, isScheduled, hasDelay); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Scheduled successfully")
	}
}
