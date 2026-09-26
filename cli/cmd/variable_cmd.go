package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runVariableCmd dispatches gitmap var subcommands.
func runVariableCmd(args []string) error {
	if len(args) == 0 || isVarHelp(args[0]) {
		printVarHelp()

		return nil
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	switch sub {
	case "set":
		return runVarSet(rest)
	case "get":
		return runVarGet(rest)
	case "ls", "list":
		return runVarList(rest)
	case "rm", "delete", "del":
		return runVarDelete(rest)
	case "export":
		return runVarExport(rest)
	default:
		return apperror.NewValidationError("unknown var subcommand: " + sub)
	}
}

func isVarHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func extractScopeFlag(args []string) (string, []string) {
	scope := "global"
	var clean []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--scope" && i+1 < len(args) {
			scope = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(args[i], "--scope=") {
			scope = strings.TrimPrefix(args[i], "--scope=")
			continue
		}
		clean = append(clean, args[i])
	}

	return scope, clean
}

func runVarSet(args []string) error {
	scope, clean := extractScopeFlag(args)
	if len(clean) < 2 {
		return apperror.NewValidationError("usage: gitmap var set <key> <value> [--scope <subcommand>]")
	}
	key := clean[0]
	value := strings.Join(clean[1:], " ")

	if err := config.SetVariable(scope, key, value); err != nil {
		return apperror.WrapSimple(err, "save variable")
	}

	fmt.Printf("%s✔ Variable \033[1m$%s\033[0m set in [%s] scope%s\n", constants.ColorGreen, key, scope, constants.ColorReset)

	return nil
}

func runVarGet(args []string) error {
	scope, clean := extractScopeFlag(args)
	if len(clean) < 1 {
		return apperror.NewValidationError("usage: gitmap var get <key> [--scope <subcommand>]")
	}
	key := clean[0]
	val, found := config.GetVariable(scope, key)
	if !found {
		return apperror.NewNotFound("var.get", "E_NOT_FOUND", fmt.Sprintf("variable $%s not found in scope %s", key, scope))
	}

	fmt.Println(val)

	return nil
}

func runVarDelete(args []string) error {
	scope, clean := extractScopeFlag(args)
	if len(clean) < 1 {
		return apperror.NewValidationError("usage: gitmap var rm <key> [--scope <subcommand>]")
	}
	key := clean[0]
	if err := config.DeleteVariable(scope, key); err != nil {
		return apperror.WrapSimple(err, "delete variable")
	}

	fmt.Printf("%s✔ Variable \033[1m$%s\033[0m removed from [%s] scope%s\n", constants.ColorGreen, key, scope, constants.ColorReset)

	return nil
}

func runVarList(args []string) error {
	glob, scoped := config.ListVariables()
	isJSON := hasArgFlag(args, "--json")
	if isJSON {
		res := map[string]any{
			"global": glob,
			"scopes": scoped,
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))

		return nil
	}

	fmt.Printf("\n%s● GitMap Variables%s\n", constants.ColorCyan, constants.ColorReset)
	renderGlobalVars(glob)
	renderScopedVars(scoped)
	fmt.Println()

	return nil
}

func renderGlobalVars(glob map[string]string) {
	fmt.Println("  Global Scope:")
	if len(glob) == 0 {
		fmt.Println("    (no global variables defined)")

		return
	}
	keys := make([]string, 0, len(glob))
	for k := range glob {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("    $%s = %s\n", k, glob[k])
	}
}

func renderScopedVars(scoped map[string]map[string]string) {
	if len(scoped) == 0 {
		return
	}
	for sc, m := range scoped {
		fmt.Printf("  Scope [%s]:\n", sc)
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("    $%s = %s\n", k, m[k])
		}
	}
}

func runVarExport(args []string) error {
	if err := config.ExportVariablesToEnv(); err != nil {
		return apperror.WrapSimple(err, "export variables")
	}

	glob, scoped := config.ListVariables()
	fmt.Printf("%s✔ Exported GitMap variables to environment:%s\n", constants.ColorGreen, constants.ColorReset)
	for k, v := range glob {
		fmt.Printf("  export %s=%q\n", k, v)
	}
	for sc, m := range scoped {
		for k, v := range m {
			scopedKey := strings.ToUpper(sc) + "_" + k
			fmt.Printf("  export %s=%q\n", scopedKey, v)
		}
	}

	return nil
}

func printVarHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap var [command] [args]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Commands:" + constants.ColorReset)
	fmt.Printf("  %-24s %s\n", "set <key> <val>", "Set a variable (--scope <subcmd> for command scope)")
	fmt.Printf("  %-24s %s\n", "get <key>", "Retrieve a variable value")
	fmt.Printf("  %-24s %s\n", "ls, list", "List all stored variables (--json)")
	fmt.Printf("  %-24s %s\n", "rm, delete <key>", "Remove a variable")
	fmt.Printf("  %-24s %s\n", "export", "Export variables to OS environment")
	fmt.Println()
	fmt.Println("Variables can be referenced as $VAR or ${VAR} in prompt text, folders, and args.")
}
