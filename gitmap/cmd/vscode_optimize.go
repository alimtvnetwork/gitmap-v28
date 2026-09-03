// Package cmd — vscode_optimize.go handles optimize-projects and clear for VS Code.
package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

type vscodeOptimizeFlags struct {
	Except []string
	DryRun bool
	Yes    bool
}

func parseVSCodeOptimizeFlags(args []string) vscodeOptimizeFlags {
	fs := flag.NewFlagSet("vscode-optimize", flag.ContinueOnError)
	exceptStr := fs.String("except", "", "Comma-separated list of IDs, names, or paths to exclude")
	fs.StringVar(exceptStr, "e", "", "Alias for --except")
	dryRun := fs.Bool("dry-run", false, "Preview actions without modifying projects.json")
	fs.BoolVar(dryRun, "d", false, "Alias for --dry-run")
	yes := fs.Bool("yes", false, "Confirm optimization without prompt")
	fs.BoolVar(yes, "y", false, "Alias for --yes")
	_ = fs.Parse(args)

	var excepts []string
	if *exceptStr != "" {
		excepts = strings.Split(*exceptStr, ",")
	}
	return vscodeOptimizeFlags{Except: excepts, DryRun: *dryRun, Yes: *yes}
}

func runVSCodeOptimize(args []string) error {
	opts := parseVSCodeOptimizeFlags(args)
	summary, err := vscodepm.OptimizeProjects(opts.Except, opts.DryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error optimizing VS Code projects: %v\n", err)
		return err
	}
	printVSCodeOptimizeResult(summary, opts.DryRun)
	return nil
}

func printVSCodeOptimizeResult(s vscodepm.OptimizeSummary, isDryRun bool) {
	if isDryRun {
		fmt.Printf("%s [dry-run] %d duplicate VS Code project(s) would be merged. Total remaining: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, s.Removed, s.Remaining)
		return
	}
	if s.Removed == 0 {
		fmt.Printf("%s No duplicate VS Code projects found. Total active: %d\n",
			constants.ColorGreen+"✓"+constants.ColorReset, s.Remaining)
		return
	}
	fmt.Printf("%s Successfully merged and removed %d duplicate VS Code project(s). Total active: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, s.Removed, s.Remaining)
}

type vscodeClearFlags struct {
	Except      []string
	OnlyMissing bool
	DryRun      bool
	Yes         bool
}

func parseVSCodeClearFlags(args []string) vscodeClearFlags {
	fs := flag.NewFlagSet("vscode-clear", flag.ContinueOnError)
	exceptStr := fs.String("except", "", "Comma-separated list of names or paths to keep")
	fs.StringVar(exceptStr, "e", "", "Alias for --except")
	missing := fs.Bool("missing", false, "Clear only projects whose directories no longer exist")
	fs.BoolVar(missing, "m", false, "Alias for --missing")
	dryRun := fs.Bool("dry-run", false, "Preview clearance without writing")
	fs.BoolVar(dryRun, "d", false, "Alias for --dry-run")
	yes := fs.Bool("yes", false, "Confirm clearing without prompt")
	fs.BoolVar(yes, "y", false, "Alias for --yes")
	_ = fs.Parse(args)

	var excepts []string
	if *exceptStr != "" {
		excepts = strings.Split(*exceptStr, ",")
	}
	return vscodeClearFlags{Except: excepts, OnlyMissing: *missing, DryRun: *dryRun, Yes: *yes}
}

func runVSCodeClear(args []string) error {
	opts := parseVSCodeClearFlags(args)
	summary, targets, err := vscodepm.ClearProjectsWithTargets(opts.Except, opts.OnlyMissing, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inspecting VS Code projects: %v\n", err)
		return err
	}
	if len(targets) == 0 {
		fmt.Println("No projects to clear in VS Code Project Manager.")
		return nil
	}
	printVSCodeClearPreview(targets)
	if opts.DryRun {
		fmt.Printf("\n%s [dry-run] %d VS Code project(s) would be cleared. Remaining: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets), summary.Remaining)
		return nil
	}
	if !opts.Yes && !askVSCodeClearConfirmation(len(targets)) {
		fmt.Println("Clearance cancelled. No changes made.")
		return nil
	}
	finalSummary, _, writeErr := vscodepm.ClearProjectsWithTargets(opts.Except, opts.OnlyMissing, false)
	if writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error clearing VS Code projects: %v\n", writeErr)
		return writeErr
	}
	fmt.Printf("\n%s Successfully removed %d VS Code project(s). Remaining: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, finalSummary.Removed, finalSummary.Remaining)
	return nil
}

func printVSCodeClearPreview(targets []vscodepm.Entry) {
	fmt.Printf("\n  %sTargeting %d VS Code project(s) to clear:%s\n\n",
		constants.ColorYellow, len(targets), constants.ColorReset)
	fmt.Printf("    %-6s %-26s %-20s %s\n", "ID", "NAME", "SLUG", "ROOT PATH")
	fmt.Printf("    %s\n", strings.Repeat("─", 88))
	for i, e := range targets {
		slug := filepath.Base(e.RootPath)
		fmt.Printf("    %-6d %-26s %-20s %s\n", i+1, e.Name, slug, e.RootPath)
	}
	fmt.Printf("\n    %sTip: Exclude items using: --except \"<id, name, slug, or starts-with text>\"%s\n",
		constants.ColorDim, constants.ColorReset)
}

func askVSCodeClearConfirmation(count int) bool {
	fmt.Printf("\n  %sAre you sure you want to remove these %d project(s)? [y/N]: %s",
		constants.ColorYellow, count, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y" || ans == "yes"
}

func printVSCodeEntries(entries []vscodepm.Entry) {
	fmt.Printf("%sVS Code Project Manager Entries:%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %-35s  %s\n", "NAME", "ROOT PATH")
	fmt.Println("  --------------------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Printf("  %-35s  %s\n", e.Name, e.RootPath)
	}
}
