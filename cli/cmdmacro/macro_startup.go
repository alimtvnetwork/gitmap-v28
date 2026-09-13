package cmdmacro

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runMacroStartup manages macro registration with the startup subsystem.
func runMacroStartup(args []string) error {
	if len(args) == 0 {
		printMacroStartupHelp()

		return nil
	}

	sub := args[0]
	tail := args[1:]

	switch sub {
	case "add":
		return runMacroStartupAdd(tail)
	case "ls", "list":
		return runMacroStartupList()
	case "rm", "remove", "delete":
		return runMacroStartupRm(tail)
	case "help", "-h", "--help":
		printMacroStartupHelp()

		return nil
	default:
		return runMacroStartupAdd(args)
	}
}

func runMacroStartupAdd(args []string) error {
	if len(args) == 0 {
		return apperror.New("macro", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify macro name to add to startup: gitmap macro startup add <macro-name> [--frequency=<freq>]",
		})
	}

	macroName := args[0]
	if !macro.MacroExists(macroName) {
		return apperror.New("macro", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Macro %q does not exist. Create it first with 'gitmap macro add %s'", macroName, macroName),
		})
	}

	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec := &store.StartupItemRecord{
		Name:         macroName,
		TargetType:   "macro",
		TargetPath:   "macro:" + macroName,
		RunFrequency: "everytime",
		IsActive:     true,
		Description:  fmt.Sprintf("Auto-registered via macro startup for %s", macroName),
	}

	if err := db.SaveStartupItem(rec); err != nil {
		return err
	}

	fmt.Printf("%s Macro %q registered to run at OS startup (frequency: everytime)\n",
		constants.ColorGreen+"✓"+constants.ColorReset, macroName)

	return nil
}

func runMacroStartupList() error {
	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	items, err := db.ListStartupItems()
	if err != nil {
		return err
	}

	var macroItems []store.StartupItemRecord
	for _, it := range items {
		if it.TargetType == "macro" {
			macroItems = append(macroItems, it)
		}
	}

	if len(macroItems) == 0 {
		fmt.Printf("%s No macros configured for startup.\n", constants.ColorCyan+"ℹ"+constants.ColorReset)
		fmt.Println("  Use 'gitmap macro startup add <macro-name>' to register one.")

		return nil
	}

	fmt.Printf("\n  %-20s %-12s %-8s %s\n", "MACRO NAME", "FREQUENCY", "ACTIVE", "TARGET")
	fmt.Println("  ─────────────────────────────────────────────────────────────")
	for _, it := range macroItems {
		fmt.Printf("  %-20s %-12s %-8v %s\n", it.Name, it.RunFrequency, it.IsActive, it.TargetPath)
	}
	fmt.Println()

	return nil
}

func runMacroStartupRm(args []string) error {
	if len(args) == 0 {
		return apperror.New("macro", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify macro name to remove from startup: gitmap macro startup rm <name>",
		})
	}

	db, err := store.OpenStartupSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.DeleteStartupItem(args[0]); err != nil {
		return err
	}

	fmt.Printf("%s Macro %q removed from startup.\n", constants.ColorGreen+"✓"+constants.ColorReset, args[0])

	return nil
}

func printMacroStartupHelp() {
	fmt.Println("Usage: gitmap macro startup <subcommand> [macro-name] [flags]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  add <macro-name>              Register a macro to run on OS startup")
	fmt.Println("  ls, list                      List all macros registered to run on startup")
	fmt.Println("  rm <macro-name>               Unregister a macro from OS startup")
}
