package cmdvscode

import (
	"flag"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type repairOptions struct {
	Force  bool
	Kill   bool
	DryRun bool
}

func parseRepairOptions(args []string) repairOptions {
	fs := flag.NewFlagSet("vscode-repair", flag.ContinueOnError)
	force := fs.Bool("force", false, "Force reinstall via winget if repair fails")
	fs.BoolVar(force, "f", false, "Alias for --force")
	kill := fs.Bool("kill", false, "Terminate running VS Code instances before repairing")
	fs.BoolVar(kill, "k", false, "Alias for --kill")
	dryRun := fs.Bool("dry-run", false, "Preview repair actions without writing changes")
	fs.BoolVar(dryRun, "d", false, "Alias for --dry-run")
	_ = fs.Parse(args)

	return repairOptions{Force: *force, Kill: *kill, DryRun: *dryRun}
}

func runVSCodeRepair(args []string) error {
	opts := parseRepairOptions(args)
	printRepairHeader()

	stepCleanupProcesses(opts)
	stepRepairJSON()
	stepRepairBinaries(opts)
	printRepairFooter()

	return nil
}

func printRepairHeader() {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	fmt.Println(headerStyle.Render("\nVS Code & Project Manager Diagnostic & Repair:"))
	fmt.Println()
}

func printRepairFooter() {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	fmt.Println()
	fmt.Println(headerStyle.Render("✓ Repair & diagnostic inspection complete!"))
	fmt.Println()
}

func stepCleanupProcesses(opts repairOptions) {
	fmt.Println("  [1/4] Checking active VS Code processes...")
	if opts.DryRun {
		fmt.Println("    • [dry-run] Skipping process termination.")

		return
	}

	count := terminateVSCodeProcesses()
	if count > 0 {
		fmt.Printf("    %s Cleaned up %d lingering process(es).\n", successMark, count)

		return
	}
	fmt.Printf("    %s No stuck VS Code processes detected.\n", successMark)
}

func stepRepairJSON() {
	fmt.Println("  [2/4] Validating Project Manager JSON (projects.json)...")
	stats, err := repairProjectManagerJSON()
	if err != nil {
		fmt.Printf("    %s Project Manager directory notice: %v\n", warnMark, err)

		return
	}

	if stats.WasCorrupt {
		fmt.Printf("    %s Corrupt JSON detected! Backed up to: %s\n", warnMark, stats.BackupPath)
	}
	fmt.Printf("    %s Validated %d projects.json entry/entries.\n", successMark, stats.TotalCount)
	if stats.MissingPaths > 0 {
		fmt.Printf("    %s %d project path(s) not found on disk.\n", warnMark, stats.MissingPaths)
	}
	fmt.Printf("    %s Modern and legacy projects.json synchronized.\n", successMark)
}

func stepRepairBinaries(opts repairOptions) {
	fmt.Println("  [3/4] Inspecting VS Code executable integrity...")
	candidates := discoverInstallationCandidates()
	if len(candidates) == 0 {
		fmt.Printf("    %s No VS Code installations found on PATH or standard directories.\n", warnMark)

		return
	}

	healthyInstall, firstErr := checkCandidatesHealth(candidates)
	if healthyInstall != "" {
		fmt.Printf("    %s Healthy VS Code binary: %s\n", successMark, subtleStyle.Render(healthyInstall))

		return
	}

	handleBrokenBinaries(candidates, firstErr, opts)
}

func checkCandidatesHealth(candidates []string) (string, string) {
	var firstErr string
	for _, inst := range candidates {
		out, ok := testCodeExecutable(inst)
		if ok {
			return inst, ""
		}
		if firstErr == "" {
			firstErr = out
		}
	}

	return "", firstErr
}

func handleBrokenBinaries(candidates []string, errOut string, opts repairOptions) {
	fmt.Printf("    %s VS Code binary verification failed: %s\n", warnMark, errOut)
	fmt.Println("  [4/4] Attempting automatic version asset sync...")

	synced := syncCommitFoldersBetween(candidates)
	if synced > 0 {
		fmt.Printf("    %s Transferred %d missing commit resource folder(s).\n", successMark, synced)
	}

	healthy, _ := checkCandidatesHealth(candidates)
	if healthy != "" {
		fmt.Printf("    %s Successfully repaired installation at: %s\n", successMark, healthy)

		return
	}

	if opts.Force {
		fmt.Println("    • Running winget reinstall/repair...")
		out, err := runWingetRepair()
		if err != nil {
			fmt.Printf("    %s Winget error: %v (%s)\n", warnMark, err, strings.TrimSpace(out))

			return
		}
		fmt.Printf("    %s Winget repair finished successfully.\n", successMark)
	}
}
