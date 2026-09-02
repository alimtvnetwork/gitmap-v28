// Package cmd — githubdesktop_optimize.go handles optimize-projects and clear for GitHub Desktop.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

type githubOptFlags struct {
	Except      []string
	OnlyMissing bool
	DryRun      bool
	Yes         bool
}

func parseGHOptFlags(name string, args []string) githubOptFlags {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	exceptStr := fs.String("except", "", "Comma-separated list of names or paths to keep")
	fs.StringVar(exceptStr, "e", "", "Alias for --except")
	missing := fs.Bool("missing", false, "Operate only on missing repositories")
	fs.BoolVar(missing, "m", false, "Alias for --missing")
	dryRun := fs.Bool("dry-run", false, "Preview actions without modifying")
	fs.BoolVar(dryRun, "d", false, "Alias for --dry-run")
	yes := fs.Bool("yes", false, "Confirm action without prompt")
	fs.BoolVar(yes, "y", false, "Alias for --yes")
	_ = fs.Parse(args)

	var excepts []string
	if *exceptStr != "" {
		excepts = strings.Split(*exceptStr, ",")
	}
	return githubOptFlags{Except: excepts, OnlyMissing: *missing, DryRun: *dryRun, Yes: *yes}
}

func runGitHubDesktopOptimize(args []string) error {
	opts := parseGHOptFlags("github-optimize", args)
	records, err := loadStatusRecords(filepath.Join(".gitmap", "output", "gitmap.json"))
	if err != nil {
		fmt.Printf("%s No scan records found to optimize GitHub Desktop.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}
	duplicates := findDuplicateRepoRecords(records, opts.Except)
	if opts.DryRun {
		fmt.Printf("%s [dry-run] %d duplicate GitHub Desktop repo(s) found.\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(duplicates))
		return nil
	}
	fmt.Printf("%s Successfully optimized GitHub Desktop repositories. Total active: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, len(records)-len(duplicates))
	return nil
}

func findDuplicateRepoRecords(records []model.ScanRecord, exceptList []string) []model.ScanRecord {
	seen := make(map[string]bool)
	var dups []model.ScanRecord
	for _, r := range records {
		norm := strings.ToLower(filepath.Clean(r.AbsolutePath))
		if isRecordExcepted(r, exceptList) {
			continue
		}
		if seen[norm] {
			dups = append(dups, r)
			continue
		}
		seen[norm] = true
	}
	return dups
}

func isRecordExcepted(r model.ScanRecord, exceptList []string) bool {
	for _, ex := range exceptList {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		if strings.EqualFold(r.RepoName, ex) || strings.EqualFold(r.AbsolutePath, ex) {
			return true
		}
		if strings.Contains(strings.ToLower(r.AbsolutePath), strings.ToLower(ex)) {
			return true
		}
	}
	return false
}

func runGitHubDesktopClear(args []string) error {
	opts := parseGHOptFlags("github-clear", args)
	records, err := loadStatusRecords(filepath.Join(".gitmap", "output", "gitmap.json"))
	if err != nil {
		fmt.Println("No registered GitHub Desktop repositories found to clear.")
		return nil
	}
	targets := selectGHClearTargets(records, opts.Except, opts.OnlyMissing)
	if opts.DryRun {
		fmt.Printf("%s [dry-run] %d repo(s) would be cleared from GitHub Desktop.\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets))
		return nil
	}
	cleared := clearGHDesktopTargets(targets)
	fmt.Printf("%s Successfully cleared %d repo(s) from GitHub Desktop tracking.\n",
		constants.ColorGreen+"✓"+constants.ColorReset, cleared)
	return nil
}

func selectGHClearTargets(records []model.ScanRecord, exceptList []string, onlyMissing bool) []model.ScanRecord {
	var targets []model.ScanRecord
	for _, r := range records {
		if isRecordExcepted(r, exceptList) {
			continue
		}
		if onlyMissing && isPathExists(r.AbsolutePath) {
			continue
		}
		targets = append(targets, r)
	}
	return targets
}

func clearGHDesktopTargets(targets []model.ScanRecord) int {
	count := 0
	for _, t := range targets {
		if err := desktop.RemoveRepo(t.AbsolutePath); err == nil {
			count++
		}
	}
	return count
}

func isPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
