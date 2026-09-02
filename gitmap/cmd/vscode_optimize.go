// Package cmd — vscode_optimize.go handles optimize-projects and clear for VS Code.
package cmd

import (
	"flag"
	"fmt"
	"os"
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
	summary, err := vscodepm.ClearProjects(opts.Except, opts.OnlyMissing, opts.DryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing VS Code projects: %v\n", err)
		return err
	}
	printVSCodeClearResult(summary, opts.DryRun, opts.OnlyMissing)
	return nil
}

func printVSCodeClearResult(s vscodepm.OptimizeSummary, isDryRun, isOnlyMissing bool) {
	action := "cleared"
	if isOnlyMissing {
		action = "missing project(s) cleared"
	}
	if isDryRun {
		fmt.Printf("%s [dry-run] %d VS Code project(s) would be %s. Remaining: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, s.Removed, action, s.Remaining)
		return
	}
	if s.Removed == 0 {
		fmt.Printf("%s No projects to clear in VS Code Project Manager.\n",
			constants.ColorGreen+"✓"+constants.ColorReset,
		)
		return
	}
	fmt.Printf("%s Successfully removed %d VS Code project(s). Remaining: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, s.Removed, s.Remaining)
}

func printVSCodeEntries(entries []vscodepm.Entry) {
	fmt.Printf("%sVS Code Project Manager Entries:%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %-35s  %s\n", "NAME", "ROOT PATH")
	fmt.Println("  --------------------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Printf("  %-35s  %s\n", e.Name, e.RootPath)
	}
}
