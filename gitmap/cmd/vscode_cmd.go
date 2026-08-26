package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

// runVSCode handles `gitmap vscode` CLI commands.
func runVSCode(args []string) {
	if len(args) == 0 {
		printVSCodeUsage()
		os.Exit(1)
	}
	dispatchVSCodeAction(args)
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
	default:
		printVSCodeUsage()
		os.Exit(1)
	}
}

func handleVSCodeAdd(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap vscode add <path>")
		os.Exit(1)
	}
	_ = runVSCodeAdd(args[1])
}

func handleVSCodeRm(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap vscode rm <path or name>")
		os.Exit(1)
	}
	_ = runVSCodeRm(args[1])
}

func printVSCodeUsage() {
	fmt.Println("Usage: gitmap vscode [ls|add <path>|rm <path or name>]")
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

func printVSCodeEntries(entries []vscodepm.Entry) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	fmt.Println(titleStyle.Render("VS Code Project Manager Entries:"))
	fmt.Printf("  %-35s  %s\n", "NAME", "ROOT PATH")
	fmt.Println("  --------------------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Printf("  %-35s  %s\n", e.Name, e.RootPath)
	}
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
		return "", fmt.Errorf("error resolving path: %v", err)
	}
	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("invalid directory: %s", absPath)
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
	if err != nil || targetPath == "" {
		if targetPath == "" {
			fmt.Printf("Project not found: %s\n", target)
		}
		return err
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
