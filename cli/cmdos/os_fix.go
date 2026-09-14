package cmdos

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/osfix"
)

func runOSFix(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) || args[0] == "help" {
		printOSFixUsage()
		return nil
	}
	return dispatchOSFixSubcommand(strings.ToLower(args[0]), args[1:])
}

func dispatchOSFixSubcommand(sub string, args []string) error {
	switch sub {
	case "ls", "list":
		return runOSFixList()
	case "add":
		return runOSFixAdd(args)
	case "edit":
		return runOSFixEdit(args)
	case "rm", "del", "delete":
		return runOSFixDelete(args)
	case "run":
		return runOSFixExecute(args)
	case "export":
		return runOSFixExport(args)
	case "export-all":
		return runOSFixExportAll(args)
	case "import":
		return runOSFixImport(args)
	case "import-all":
		return runOSFixImportAll(args)
	default:
		return runOSFixExecute(append([]string{sub}, args...))
	}
}

func runOSFixList() error {
	res := osfix.ListFixes()
	if res.IsFailure() {
		return res.AppError()
	}
	if len(res.Value) == 0 {
		fmt.Println("No registered OS fixes found.")
		return nil
	}
	fmt.Printf("%-20s %-10s %s\n", "NAME", "PLATFORM", "COMMAND")
	fmt.Println(strings.Repeat("-", 60))
	for _, item := range res.Value {
		fmt.Printf("%-20s %-10s %s\n", item.Name, item.Platform, item.Command)
	}
	return nil
}

func runOSFixAdd(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("usage: gitmap os fix add <name> <command>", "E_MISSING_ARG")
	}
	item := osfix.FixItem{
		Name:     args[0],
		Command:  strings.Join(args[1:], " "),
		Platform: "all",
	}
	res := osfix.SaveFix(item)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Successfully registered fix: %s\n", item.Name)
	return nil
}

func runOSFixEdit(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("usage: gitmap os fix edit <name> <command>", "E_MISSING_ARG")
	}
	item := osfix.FixItem{
		Name:     args[0],
		Command:  strings.Join(args[1:], " "),
		Platform: "all",
	}
	res := osfix.SaveFix(item)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Successfully updated fix: %s\n", item.Name)
	return nil
}

func runOSFixDelete(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os fix rm <name>", "E_MISSING_ARG")
	}
	res := osfix.DeleteFix(args[0])
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Successfully removed fix: %s\n", args[0])
	return nil
}

func runOSFixExecute(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os fix run <name|commands>", "E_MISSING_ARG")
	}
	res := osfix.GetFix(args[0])
	if res.IsSuccess() {
		return runRegisteredFix(res.Value.Command)
	}
	return runAdhocFix(args)
}

func runRegisteredFix(command string) error {
	execRes := osfix.RunFixCommand(command)
	if execRes.IsFailure() {
		return execRes.AppError()
	}
	return nil
}

func runAdhocFix(args []string) error {
	command := strings.Join(args, " ")
	execRes := osfix.RunFixCommand(command)
	if execRes.IsFailure() {
		return execRes.AppError()
	}
	return nil
}

func runOSFixExport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os fix export <name> [file.json]", "E_MISSING_ARG")
	}
	res := osfix.GetFix(args[0])
	if res.IsFailure() {
		return res.AppError()
	}
	data, _ := json.MarshalIndent(res.Value, "", "  ")
	if len(args) >= 2 {
		return writeFixExportToFile(args[1], data, args[0])
	}
	fmt.Println(string(data))
	return nil
}

func writeFixExportToFile(destPath string, data []byte, name string) error {
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return apperror.WrapSimple(err, "write export file")
	}
	fmt.Printf("✔ Exported fix %s to %s\n", name, destPath)
	return nil
}

func runOSFixExportAll(args []string) error {
	res := osfix.ListFixes()
	if res.IsFailure() {
		return res.AppError()
	}
	data, _ := json.MarshalIndent(res.Value, "", "  ")
	if len(args) >= 1 {
		return writeAllFixesExportToFile(args[0], data, len(res.Value))
	}
	fmt.Println(string(data))
	return nil
}

func writeAllFixesExportToFile(destPath string, data []byte, count int) error {
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return apperror.WrapSimple(err, "write export-all file")
	}
	fmt.Printf("✔ Exported %d fixes to %s\n", count, destPath)
	return nil
}

func runOSFixImport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os fix import <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import file")
	}
	var item osfix.FixItem
	if unmarshalErr := json.Unmarshal(content, &item); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "unmarshal fix item")
	}
	saveRes := osfix.SaveFix(item)
	if saveRes.IsFailure() {
		return saveRes.AppError()
	}
	fmt.Printf("✔ Imported fix: %s\n", item.Name)
	return nil
}

func runOSFixImportAll(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os fix import-all <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import-all file")
	}
	var items []osfix.FixItem
	if unmarshalErr := json.Unmarshal(content, &items); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "unmarshal fixes")
	}
	for _, item := range items {
		_ = osfix.SaveFix(item)
	}
	fmt.Printf("✔ Imported %d fixes from %s\n", len(items), args[0])
	return nil
}

func printOSFixUsage() {
	fmt.Println("Usage: gitmap os fix <subcommand> [args]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  ls, list                          List all registered OS fixes")
	fmt.Println("  add <name> <command>              Register a new fix script/command")
	fmt.Println("  edit <name> <command>             Edit an existing registered fix")
	fmt.Println("  rm <name>                         Remove a registered fix")
	fmt.Println("  run <name|command...>             Run an existing fix or execute new commands")
	fmt.Println("  export <name> [file.json]         Export a fix definition")
	fmt.Println("  export-all [file.json]            Export all registered fixes")
	fmt.Println("  import <file.json>                Import a single fix")
	fmt.Println("  import-all <file.json>            Import all fixes from JSON")
}
