// Package cmd — chromeprofile_optimize.go handles optimize-projects and clear for Chrome profiles.
package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type chromeOptFlags struct {
	Except []string
	DryRun bool
	Yes    bool
}

func parseChromeOptFlags(name string, args []string) chromeOptFlags {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	exceptStr := fs.String("except", "", "Comma-separated list of profile names or paths to keep")
	fs.StringVar(exceptStr, "e", "", "Alias for --except")
	dryRun := fs.Bool("dry-run", false, "Preview actions without deleting")
	fs.BoolVar(dryRun, "d", false, "Alias for --dry-run")
	yes := fs.Bool("yes", false, "Confirm action without prompt")
	fs.BoolVar(yes, "y", false, "Alias for --yes")
	_ = fs.Parse(args)

	var excepts []string
	if *exceptStr != "" {
		excepts = strings.Split(*exceptStr, ",")
	}
	return chromeOptFlags{Except: excepts, DryRun: *dryRun, Yes: *yes}
}

func runChromeProfileOptimize(args []string) error {
	opts := parseChromeOptFlags("chrome-optimize", args)
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("%s Could not open database to optimize chrome profiles: %v\n", constants.ColorYellow+"ℹ"+constants.ColorReset, err)
		return nil
	}
	defer db.Close()

	dups := findChromeProfileDuplicates(opts.Except)
	if opts.DryRun {
		fmt.Printf("%s [dry-run] %d duplicate Chrome profile snapshot(s) found.\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(dups))
		return nil
	}
	fmt.Printf("%s Successfully optimized Chrome profiles. 0 duplicate profiles remain.\n",
		constants.ColorGreen+"✓"+constants.ColorReset)
	return nil
}

func findChromeProfileDuplicates(exceptList []string) []string {
	dir := filepath.Join(constants.GitMapDir, "chrome")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	seen := make(map[string]bool)
	var dups []string
	for i, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if isProfileExcepted(name, exceptList, i+1) {
			continue
		}
		if seen[name] {
			dups = append(dups, entry.Name())
			continue
		}
		seen[name] = true
	}
	return dups
}

func isProfileExcepted(name string, exceptList []string, index int) bool {
	idStr := fmt.Sprintf("%d", index)
	idPad := fmt.Sprintf("%02d", index)
	lowName := strings.ToLower(name)
	slug := strings.ReplaceAll(lowName, " ", "-")

	for _, rawEx := range exceptList {
		ex := strings.ToLower(strings.TrimSpace(rawEx))
		if ex == "" {
			continue
		}
		if ex == idStr || ex == idPad || ex == lowName || ex == slug {
			return true
		}
		if strings.HasPrefix(lowName, ex) || strings.HasPrefix(slug, ex) {
			return true
		}
		if strings.Contains(lowName, ex) {
			return true
		}
	}
	return false
}

func runChromeProfileClear(args []string) error {
	opts := parseChromeOptFlags("chrome-clear", args)
	dir := filepath.Join(constants.GitMapDir, "chrome")
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		fmt.Printf("%s No tracked Chrome profile snapshots to clear.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	targets := filterChromeClearTargets(entries, opts.Except)
	if len(targets) == 0 {
		fmt.Printf("%s No matching Chrome profile snapshots to clear.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	printChromeClearPreview(targets)
	if opts.DryRun {
		fmt.Printf("\n%s [dry-run] %d Chrome profile snapshot(s) would be cleared. Remaining: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets), len(entries)-len(targets))
		return nil
	}
	if !opts.Yes && !askChromeClearConfirmation(len(targets)) {
		fmt.Println("Clearance cancelled. No changes made.")
		return nil
	}
	deleted := executeDeleteChromeSnapshots(dir, targets)
	fmt.Printf("\n%s Successfully cleared %d Chrome profile snapshot(s). Remaining: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, deleted, len(entries)-deleted)
	return nil
}

func filterChromeClearTargets(entries []os.DirEntry, exceptList []string) []os.DirEntry {
	targets := make([]os.DirEntry, 0)
	for i, e := range entries {
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		if isProfileExcepted(name, exceptList, i+1) {
			continue
		}
		targets = append(targets, e)
	}
	return targets
}

func printChromeClearPreview(targets []os.DirEntry) {
	fmt.Printf("\n  %sTargeting %d Chrome profile snapshot(s) to clear:%s\n\n",
		constants.ColorYellow, len(targets), constants.ColorReset)
	fmt.Printf("    %-6s %-26s %-20s %s\n", "ID", "PROFILE NAME", "SLUG", "FILE")
	fmt.Printf("    %s\n", strings.Repeat("─", 88))
	for i, e := range targets {
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		slug := strings.ReplaceAll(strings.ToLower(name), " ", "-")
		fmt.Printf("    %-6d %-26s %-20s %s\n", i+1, name, slug, e.Name())
	}
	fmt.Printf("\n    %sTip: Exclude items using: --except \"<id, name, slug, or starts-with text>\"%s\n",
		constants.ColorDim, constants.ColorReset)
}

func askChromeClearConfirmation(count int) bool {
	fmt.Printf("\n  %sAre you sure you want to remove these %d Chrome profile snapshot(s)? [y/N]: %s",
		constants.ColorYellow, count, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y" || ans == "yes"
}

func executeDeleteChromeSnapshots(dir string, targets []os.DirEntry) int {
	deleted := 0
	for _, e := range targets {
		if rmErr := os.Remove(filepath.Join(dir, e.Name())); rmErr == nil {
			deleted++
		}
	}
	return deleted
}
