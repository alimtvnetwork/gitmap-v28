package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

// runVSCode handles `gitmap vscode` CLI commands.
func runVSCode(args []string) error {
	if len(args) == 0 {
		printVSCodeUsage()

		return apperror.NewWithDetails(
			"cmd.vscode",
			"E1023",
			"missing required vscode subcommand",
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	return dispatchVSCodeAction(args)
}

func dispatchVSCodeAction(args []string) error {
	sub := strings.ToLower(args[0])
	if isVSCodeProjectSubcommand(sub) {
		return routeVSCodeProjectAction(sub, args)
	}

	return routeVSCodeMaintenanceAction(sub, args)
}

func isVSCodeProjectSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "add" || sub == "add-project" || sub == "ap" || sub == "rm" || sub == "remove" || sub == "delete" || sub == "del"
}

func routeVSCodeProjectAction(sub string, args []string) error {
	switch sub {
	case "ls", "list":
		return runVSCodeLs()
	case "add", "add-project", "ap":
		handleVSCodeAdd(args)

		return nil
	default:
		handleVSCodeRm(args)

		return nil
	}
}

func routeVSCodeMaintenanceAction(sub string, args []string) error {
	switch sub {
	case "pap", "prompt-all-project", "plugins", "plugin":
		fmt.Printf("Feature [vscode %s] is not yet implemented\n", sub)

		return nil
	case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe", "dedup":
		return runVSCodeOptimize(args[1:])
	case "clear", "clean":
		return runVSCodeClear(args[1:])
	case "group", "groups", "grp":
		return runVSCodeGroup(args[1:])
	case "find-duplicates", "duplicates", "dups", "find-dups":
		return runFindDuplicates("vscode", args[1:])
	default:
		printVSCodeUsage()

		return apperror.NewWithDetails(
			"cmd.vscode.dispatch",
			"E1024",
			fmt.Sprintf("unknown vscode subcommand '%s'", sub),
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"subcommand": sub},
		)
	}
}

func handleVSCodeAdd(args []string) {
	if len(args) < 2 {
		printVSCodeUsage()
		err := apperror.NewWithDetails(
			"cmd.vscode.add",
			"E1025",
			"missing path argument for vscode add",
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return
	}
	for _, p := range strings.Split(args[1], ",") {
		if p = strings.TrimSpace(p); p != "" {
			_ = runVSCodeAdd(p)
		}
	}
}

func handleVSCodeRm(args []string) {
	if len(args) < 2 {
		printVSCodeUsage()
		err := apperror.NewWithDetails(
			"cmd.vscode.rm",
			"E1026",
			"missing target argument for vscode rm",
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return
	}
	for _, t := range strings.Split(args[1], ",") {
		if t = strings.TrimSpace(t); t != "" {
			_ = runVSCodeRm(t)
		}
	}
}

func printVSCodeUsage() {
	fmt.Println("Usage: gitmap vscode [ls|add <path>|rm <path>|optimize-projects|clear] [flags]")
}

func runVSCodeLs() error {
	entries, err := vscodepm.ListEntries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading projects.json: %v\n", err)
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No VS Code projects registered.")
		return nil
	}
	printVSCodeEntries(entries)
	return nil
}

func runVSCodeAdd(target string) error {
	absPath, err := resolveVSCodePath(target)
	if err != nil {
		return err
	}
	name := filepath.Base(absPath)
	pair := vscodepm.Pair{RootPath: absPath, Name: name}
	summary, err := vscodepm.Sync([]vscodepm.Pair{pair})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding project: %v\n", err)
		return err
	}
	reportVSCodeAdd(summary, name, absPath)
	return nil
}

func resolveVSCodePath(target string) (string, error) {
	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", apperror.Wrap(err, "ResolveVSCodePath", map[string]any{"target": target})
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return "", apperror.Wrap(err, "ResolveVSCodePath", map[string]any{"path": absPath})
	}
	if !info.IsDir() {
		return "", apperror.New("ResolveVSCodePath", "INVALID_DIR", map[string]any{"path": absPath})
	}
	return absPath, nil
}

func reportVSCodeAdd(summary vscodepm.SyncSummary, name, absPath string) {
	if summary.Added > 0 || summary.Updated > 0 {
		fmt.Printf("Added/Updated %s (%s) in projects.json\n", name, absPath)
		return
	}
	fmt.Printf("Project %s already exists in projects.json\n", absPath)
}

func runVSCodeRm(target string) error {
	targetPath, err := findVSCodeTarget(target)
	if err != nil {
		return err
	}
	if targetPath == "" {
		fmt.Printf("Project not found: %s\n", target)
		return nil
	}
	if err := vscodepm.RemoveEntry(targetPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing project: %v\n", err)
		return err
	}
	fmt.Printf("Removed %s from projects.json\n", targetPath)
	return nil
}

func findVSCodeTarget(target string) (string, error) {
	entries, err := vscodepm.ListEntries()
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		isMatch := strings.EqualFold(e.Name, target) || strings.EqualFold(e.RootPath, target)
		if isMatch {
			return e.RootPath, nil
		}
	}
	return "", nil
}
