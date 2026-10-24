// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// scheduleSendEmailFlags holds the parsed flags for schedule-send-email.
type scheduleSendEmailFlags struct {
	to       string
	task     string
	name     string
	isList   bool
	isRemove bool
}

// RunScheduleSendEmail implements `gitmap space schedule-send-email`:
// registers (--to/--task [+ --name]), lists (--list), or removes
// (--remove --name) email triggers that fire when a gitmap task completes.
func RunScheduleSendEmail(args []string) error {
	if hasHelpFlag(args) {
		printScheduleSendEmailHelp()

		return nil
	}

	flags, err := parseScheduleSendEmailArgs(args)
	if err != nil {
		return err
	}

	return dispatchScheduleSendEmail(flags)
}

func dispatchScheduleSendEmail(flags scheduleSendEmailFlags) error {
	if flags.isList {
		return listScheduleSendEmailTriggers()
	}

	if flags.isRemove {
		return removeScheduleSendEmailTrigger(flags.name)
	}

	return registerScheduleSendEmailTrigger(flags)
}

// parseScheduleSendEmailArgs parses the flags position-agnostically,
// matching the space common flag style. Unknown flags are errors.
func parseScheduleSendEmailArgs(args []string) (scheduleSendEmailFlags, error) {
	flags := scheduleSendEmailFlags{}
	for i := 0; i < len(args); i++ {
		nextIdx, err := applyScheduleSendEmailArg(args, i, &flags)
		if err != nil {
			return scheduleSendEmailFlags{}, err
		}
		i = nextIdx
	}

	return flags, nil
}

func applyScheduleSendEmailArg(args []string, i int, flags *scheduleSendEmailFlags) (int, error) {
	arg := args[i]
	switch {
	case arg == "--list":
		flags.isList = true
		return i, nil
	case arg == "--remove":
		flags.isRemove = true
		return i, nil
	case arg == "--to" || strings.HasPrefix(arg, "--to="):
		flags.to, i = nextScheduleFlagValue(args, i, arg, "--to")
		return i, nil
	case arg == "--task" || strings.HasPrefix(arg, "--task="):
		flags.task, i = nextScheduleFlagValue(args, i, arg, "--task")
		return i, nil
	case arg == "--name" || strings.HasPrefix(arg, "--name="):
		flags.name, i = nextScheduleFlagValue(args, i, arg, "--name")
		return i, nil
	default:
		return i, apperror.NewValidationError(fmt.Sprintf("unknown flag %q: usage: %s", arg, scheduleSendEmailUsageLine))
	}
}

// nextScheduleFlagValue reads a `--flag value` or `--flag=value` pair,
// advancing the loop index past a consumed positional value.
func nextScheduleFlagValue(args []string, i int, arg, flag string) (string, int) {
	if strings.HasPrefix(arg, flag+"=") {
		return strings.TrimPrefix(arg, flag+"="), i
	}

	if i+1 < len(args) {
		return args[i+1], i + 1
	}

	return "", i
}

// registerScheduleSendEmailTrigger validates and stores one trigger.
func registerScheduleSendEmailTrigger(flags scheduleSendEmailFlags) error {
	task, valErr := validateScheduleTriggerFlags(flags)
	if valErr != nil {
		return valErr
	}

	trigger := EmailTrigger{
		Name: scheduleTriggerName(flags.name, task),
		To:   strings.TrimSpace(flags.to),
		Task: task,
	}

	return persistScheduleTrigger(trigger)
}

func validateScheduleTriggerFlags(flags scheduleSendEmailFlags) (EmailTaskType, error) {
	if strings.TrimSpace(flags.to) == "" {
		return "", apperror.NewValidationError(fmt.Sprintf("missing --to <email>: usage: %s", scheduleSendEmailUsageLine))
	}

	if strings.TrimSpace(flags.task) == "" {
		return "", apperror.NewValidationError(fmt.Sprintf("missing --task <task>: usage: %s", scheduleSendEmailUsageLine))
	}

	task, parseErr := ParseEmailTaskType(flags.task)
	if parseErr != nil {
		return "", parseErr
	}

	return task, nil
}

func persistScheduleTrigger(trigger EmailTrigger) error {
	db, dbErr := openTriggerStore()
	if dbErr != nil {
		return apperror.WrapSimple(dbErr, "open settings db:")
	}
	defer db.Close()

	if regErr := RegisterEmailTrigger(db, trigger); regErr != nil {
		return regErr
	}

	fmt.Printf("  ✓ email trigger %q registered: notify %s when %s completes\n", trigger.Name, trigger.To, trigger.Task)

	return nil
}

// scheduleTriggerName defaults the trigger label to the task name.
func scheduleTriggerName(name string, task EmailTaskType) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) > 0 {
		return trimmed
	}

	return string(task)
}

// listScheduleSendEmailTriggers prints all registered triggers in a formatted table.
func listScheduleSendEmailTriggers() error {
	db, err := openTriggerStore()
	if err != nil {
		return apperror.WrapSimple(err, "open settings db:")
	}
	defer db.Close()

	triggers, loadErr := LoadEmailTriggers(db)
	if loadErr != nil {
		return loadErr
	}

	if len(triggers) == 0 {
		fmt.Println("  no email triggers registered")

		return nil
	}

	printRegisteredTriggersTable(triggers)

	return nil
}

func printRegisteredTriggersTable(triggers []EmailTrigger) {
	cfg := termout.TableConfig{
		HeaderColor: constants.ColorYellow,
		BorderColor: constants.ColorDim,
		Columns: []termout.Column{
			{Title: "Trigger Name", MinWidth: 16},
			{Title: "Watched Task", MinWidth: 14},
			{Title: "Recipient", MinWidth: 26},
		},
		Rows: buildTriggerTableRows(triggers),
	}
	termout.PrintTable(cfg)
}

func buildTriggerTableRows(triggers []EmailTrigger) []termout.Row {
	rows := make([]termout.Row, 0, len(triggers))
	for _, trigger := range triggers {
		rows = append(rows, termout.Row{
			Cells: []string{trigger.Name, string(trigger.Task), trigger.To},
			Color: constants.ColorWhite,
		})
	}

	return rows
}

// removeScheduleSendEmailTrigger deletes the named trigger.
func removeScheduleSendEmailTrigger(name string) error {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return apperror.NewValidationError(fmt.Sprintf("missing --name <label>: usage: %s", scheduleSendEmailUsageLine))
	}

	db, err := openTriggerStore()
	if err != nil {
		return apperror.WrapSimple(err, "open settings db:")
	}
	defer db.Close()

	if remErr := RemoveEmailTrigger(db, trimmed); remErr != nil {
		return remErr
	}

	fmt.Printf("  ✓ email trigger %q removed\n", trimmed)

	return nil
}
