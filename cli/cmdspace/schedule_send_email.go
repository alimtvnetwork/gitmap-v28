// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
		arg := args[i]
		switch {
		case arg == "--list":
			flags.isList = true
		case arg == "--remove":
			flags.isRemove = true
		case arg == "--to" || strings.HasPrefix(arg, "--to="):
			flags.to, i = nextScheduleFlagValue(args, i, arg, "--to")
		case arg == "--task" || strings.HasPrefix(arg, "--task="):
			flags.task, i = nextScheduleFlagValue(args, i, arg, "--task")
		case arg == "--name" || strings.HasPrefix(arg, "--name="):
			flags.name, i = nextScheduleFlagValue(args, i, arg, "--name")
		default:
			return scheduleSendEmailFlags{}, fmt.Errorf("unknown flag %q: usage: %s", arg, scheduleSendEmailUsageLine)
		}
	}

	return flags, nil
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
	if strings.TrimSpace(flags.to) == "" {
		return fmt.Errorf("missing --to <email>: usage: %s", scheduleSendEmailUsageLine)
	}

	if strings.TrimSpace(flags.task) == "" {
		return fmt.Errorf("missing --task <task>: usage: %s", scheduleSendEmailUsageLine)
	}

	task, err := ParseEmailTaskType(flags.task)
	if err != nil {
		return err
	}

	trigger := EmailTrigger{
		Name: scheduleTriggerName(flags.name, task),
		To:   strings.TrimSpace(flags.to),
		Task: task,
	}

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

// listScheduleSendEmailTriggers prints all registered triggers.
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

	for _, trigger := range triggers {
		fmt.Printf("  • %s: %s → %s\n", trigger.Name, trigger.Task, trigger.To)
	}

	return nil
}

// removeScheduleSendEmailTrigger deletes the named trigger.
func removeScheduleSendEmailTrigger(name string) error {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return fmt.Errorf("missing --name <label>: usage: %s", scheduleSendEmailUsageLine)
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
