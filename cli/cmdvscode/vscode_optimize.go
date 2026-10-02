// Package cmd — vscode_optimize.go handles optimize-projects and clear for VS Code.
package cmdvscode

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

type vscodeOptimizeFlags struct {
	Except []string
	DryRun bool
	Yes    bool
}

func parseVSCodeOptimizeFlags(args []string) vscodeOptimizeFlags {
	args = stripProjectSubArg(args)
	fs := flag.NewFlagSet("vscode-optimize", flag.ContinueOnError)
	exceptStr := fs.String("except", "", "Comma-separated list of IDs, names, or paths to exclude")
	fs.StringVar(exceptStr, "e", "", "Alias for --except")
	dryRun := fs.Bool("dry-run", false, "Preview actions without modifying projects.json")
	fs.BoolVar(dryRun, "d", false, "Alias for --dry-run")
	yes := fs.Bool("yes", false, "Confirm optimization without prompt")
	fs.BoolVar(yes, "y", false, "Alias for --yes")
	_ = fs.Parse(args)

	return makeVSCodeOptimizeFlags(*exceptStr, *dryRun, *yes)
}

func stripProjectSubArg(args []string) []string {
	if len(args) > 0 && (args[0] == "projects" || args[0] == "project") {
		return args[1:]
	}

	return args
}

func makeVSCodeOptimizeFlags(exceptStr string, dryRun, yes bool) vscodeOptimizeFlags {
	var excepts []string
	if exceptStr != "" {
		excepts = strings.Split(exceptStr, ",")
	}

	return vscodeOptimizeFlags{Except: excepts, DryRun: dryRun, Yes: yes}
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
	prefix := constants.ColorGreen + "✓" + constants.ColorReset
	if isDryRun {
		prefix = constants.ColorYellow + "ℹ [dry-run]" + constants.ColorReset
	}
	if s.Removed == 0 && len(s.Advices) == 0 {
		fmt.Printf("%s No duplicate or missing VS Code projects found. Total active: %d\n", prefix, s.Remaining)
		return
	}

	fmt.Printf("%s VS Code Project Manager Optimization Complete:\n", prefix)
	fmt.Printf("    • Duplicate entries merged/removed: %d\n", s.RemovedDuplicates)
	fmt.Printf("    • Missing/stale folders pruned:    %d\n", s.RemovedMissing)
	fmt.Printf("    • Total active projects retained:  %d\n", s.Remaining)
	printOptimizationAdvices(s.Advices)
}

func printOptimizationAdvices(advices []vscodepm.ProjectOptimizationAdvice) {
	if len(advices) == 0 {
		return
	}
	fmt.Printf("\n  %s📁 Duplicate Project Relocation Advice:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, adv := range advices {
		fmt.Printf("    • Project %q:\n", adv.ProjectName)
		fmt.Printf("      - Canonical: %s\n", adv.CanonicalPath)
		fmt.Printf("      - Duplicate: %s\n", adv.DuplicatePath)
		fmt.Printf("      - %s\n", adv.Advice)
	}
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
		fmt.Println("Clearance canceled. No changes made.")

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
	ans = strings.TrimSpace(ans)

	return strings.EqualFold(ans, "y") || strings.EqualFold(ans, "yes")
}

func printVSCodeEntries(entries []vscodepm.Entry) {
	fmt.Printf("%sVS Code Project Manager Entries:%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %-35s  %s\n", "NAME", "ROOT PATH")
	fmt.Println("  --------------------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Printf("  %-35s  %s\n", e.Name, e.RootPath)
	}
}
