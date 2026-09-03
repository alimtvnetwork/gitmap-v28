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
		err := apperror.NewWithDetails(
			"cmd.vscode",
			"E1023",
			"missing required vscode subcommand",
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
	}
	dispatchVSCodeAction(args)
	return nil
}

func dispatchVSCodeAction(args []string) {
	switch strings.ToLower(args[0]) {
	case "ls", "list":
		_ = runVSCodeLs()
	case "pap", "prompt-all-project":
		fmt.Println("Feature [vscode pap] is not yet implemented")
	case "plugins", "plugin":
		fmt.Println("Feature [vscode plugins] is not yet implemented")
	case "add", "add-project", "ap":
		handleVSCodeAdd(args)
	case "rm", "remove", "delete", "del":
		handleVSCodeRm(args)
	case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe", "dedup":
		_ = runVSCodeOptimize(args[1:])
	case "clear", "clean":
		_ = runVSCodeClear(args[1:])
	case "group", "groups", "grp":
		_ = runVSCodeGroup(args[1:])
	case "find-duplicates", "duplicates", "dups", "find-dups":
		_ = runFindDuplicates("vscode", args[1:])
	default:
		printVSCodeUsage()
		err := apperror.NewWithDetails(
			"cmd.vscode.dispatch",
			"E1024",
			fmt.Sprintf("unknown vscode subcommand '%s'", args[0]),
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"subcommand": args[0]},
		)
		cliexit.HandleError(err, 1)
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
