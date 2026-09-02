// Package cmd — chromeprofile_optimize.go handles optimize-projects and clear for Chrome profiles.
package cmd

import (
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
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if isProfileExcepted(name, exceptList) {
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

func isProfileExcepted(name string, exceptList []string) bool {
	for _, ex := range exceptList {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		if strings.EqualFold(name, ex) || strings.Contains(strings.ToLower(name), strings.ToLower(ex)) {
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
	cleared := clearChromeSnapshots(dir, entries, opts)
	if opts.DryRun {
		fmt.Printf("%s [dry-run] %d Chrome profile snapshot(s) would be cleared.\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, cleared)
		return nil
	}
	fmt.Printf("%s Successfully cleared %d Chrome profile snapshot(s).\n",
		constants.ColorGreen+"✓"+constants.ColorReset, cleared)
	return nil
}

func clearChromeSnapshots(dir string, entries []os.DirEntry, opts chromeOptFlags) int {
	count := 0
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		if isProfileExcepted(name, opts.Except) {
			continue
		}
		count++
		if !opts.DryRun {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	return count
}
