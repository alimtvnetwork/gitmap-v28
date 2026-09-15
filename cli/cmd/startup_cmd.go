package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunStartupCmd routes gitmap startup subcommands.
func RunStartupCmd(args []string) error {
	checkHelp(constants.CmdStartup, args)
	if len(args) == 0 {
		return runStartupLs([]string{})
	}

	sub := strings.ToLower(args[0])
	tail := args[1:]

	return routeStartupSubcommand(sub, tail).AsError()
}

func routeStartupSubcommand(sub string, tail []string) result.ErrorWrapper {
	switch sub {
	case "ls", "list":
		return result.MatchWrapper(runStartupLs(tail))
	case "add":
		return result.MatchWrapper(runStartupAddCmd(tail))
	case "rm", "remove", "delete":
		return result.MatchWrapper(runStartupRm(tail))
	case "run", "exec":
		return result.MatchWrapper(runStartupExec(tail))
	case "logs", "log":
		return result.MatchWrapper(runStartupLogs(tail))
	case "help", "-h", "--help":
		printStartupHelp()
		return result.SuccessWrapper()
	default:
		return result.FailureWrapper(apperror.New("startup", "E_UNKNOWN_SUBCMD", map[string]any{
			"msg": fmt.Sprintf("Unknown startup subcommand: %q (expected ls, add, rm, run, logs)", sub),
		}))
	}
}

func runStartupLs(args []string) error {
	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	items, err := db.ListStartupItems()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Printf("%s No startup items configured.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		fmt.Println("  Use 'gitmap startup add <script|macro|exe>' to add one.")

		return nil
	}

	printStartupTable(items)

	return nil
}

func printStartupTable(items []store.StartupItemRecord) {
	fmt.Printf("\n  %-20s %-12s %-14s %-8s %s\n", "NAME", "FREQUENCY", "TARGET TYPE", "ACTIVE", "TARGET PATH")
	fmt.Println("  ─────────────────────────────────────────────────────────────────────────────")
	for _, it := range items {
		activeStr := "yes"
		if it.IsActive == false {
			activeStr = "no"
		}
		fmt.Printf("  %-20s %-12s %-14s %-8s %s\n", it.Name, it.RunFrequency, it.TargetType, activeStr, it.TargetPath)
	}
	fmt.Println()
}

func runStartupAddCmd(args []string) error {
	opts, err := parseStartupAddArgs(args)
	if err != nil {
		return err
	}

	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec := buildStartupRecord(opts)
	if err := db.SaveStartupItem(rec); err != nil {
		return err
	}

	registerNativeOSAutostart(rec)
	fmt.Printf("%s Registered startup item %q (%s, %s)\n",
		constants.ColorGreen+"✓"+constants.ColorReset, rec.Name, rec.TargetType, rec.RunFrequency)

	return nil
}

func runStartupRm(args []string) error {
	if len(args) == 0 {
		return apperror.New("startup", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify startup item name to remove: gitmap startup rm <name>",
		})
	}

	name := args[0]
	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.DeleteStartupItem(name); err != nil {
		return err
	}

	unregisterNativeOSAutostart(name)
	fmt.Printf("%s Removed startup item %q\n", constants.ColorGreen+"✓"+constants.ColorReset, name)

	return nil
}

func printStartupHelp() {
	fmt.Println("Usage: gitmap startup <subcommand> [args...]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  ls, list                      List all registered startup items")
	fmt.Println("  add <path|macro> [flags]      Register a new startup script, binary, or macro")
	fmt.Println("  rm <name>                     Remove a registered startup item")
	fmt.Println("  run <name>                    Run a registered startup item on demand")
	fmt.Println("  logs <name>                   View execution logs for a startup item")
	fmt.Println("\nFlags for 'add':")
	fmt.Println("  --frequency <everytime|once-a-week>  Execution frequency (default: everytime)")
	fmt.Println("  --name <name>                        Item identifier (default: file/macro name)")
	fmt.Println("  --icon <path>                        Icon file path (.ico / .png)")
	fmt.Println("  --args <arguments>                   Arguments passed to script or binary")
	fmt.Println("  --desc <description>                 Helpful description")
}
