// Package cmd — `gitmap backup` ls/prune (v6.57.0).
//
// Scans `<cwd>/.gitmap/backup/` (the tree used by fix-repo / undo)
// and exposes two read-only-friendly operations:
//
//	gitmap backup ls                          # group by repo, count + size
//	gitmap backup prune --keep=N              # keep newest N per repo
//	gitmap backup prune --older-than=DAYS     # drop entries older than DAYS
//	gitmap backup prune --dry-run             # print what would be deleted
//
// Bounded-retention answers item #2 of the v6.55.0 suggestions list.
// Self-contained: walks the fix-repo timestamp layout documented in
// constants_undo.go, never touches release/ or release-assets/.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// backupSnapshot is one `.gitmap/backup/<repo>/v<N>/fix-repo/<ts>/` entry.
type backupSnapshot struct {
	repo string
	full string // absolute path
	ts   time.Time
	size int64
}

// runBackup dispatches `gitmap backup <sub>`.
func runBackup(args []string) {
	checkHelp("backup", args)
	if len(args) == 0 {
		runBackupLs(nil)

		return
	}
	switch args[0] {
	case constants.SubCmdBackupLs, constants.SubCmdBackupList:
		runBackupLs(args[1:])
	case constants.SubCmdBackupPrune:
		runBackupPrune(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "gitmap backup: unknown subcommand %q (want: ls | prune)\n", args[0])
		os.Exit(2)
	}
}

// runBackupLs prints a per-repo summary of backup snapshots.
func runBackupLs(_ []string) {
	root, err := backupRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gitmap backup ls: %v\n", err)
		os.Exit(1)
	}
	snaps, err := collectSnapshots(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gitmap backup ls: %v\n", err)
		os.Exit(1)
	}
	if len(snaps) == 0 {
		fmt.Fprintf(os.Stdout, "\n  no backups under %s\n\n", root)

		return
	}
	printBackupTable(root, snaps)
}

// runBackupPrune deletes snapshots per --keep / --older-than flags.
func runBackupPrune(args []string) {
	fs := flag.NewFlagSet("backup prune", flag.ContinueOnError)
	keep := fs.Int("keep", 0, "keep newest N snapshots per repo (0 = unlimited)")
	olderDays := fs.Int("older-than", 0, "delete snapshots older than N days (0 = no age filter)")
	dryRun := fs.Bool("dry-run", false, "print what would be deleted; do not touch disk")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if *keep == 0 && *olderDays == 0 {
		fmt.Fprintln(os.Stderr, "gitmap backup prune: pass --keep=N and/or --older-than=DAYS")
		os.Exit(2)
	}
	root, err := backupRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gitmap backup prune: %v\n", err)
		os.Exit(1)
	}
	snaps, err := collectSnapshots(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gitmap backup prune: %v\n", err)
		os.Exit(1)
	}
	victims := selectPruneVictims(snaps, *keep, *olderDays)
	applyPrune(victims, *dryRun)
}

// backupRoot returns the absolute `<cwd>/.gitmap/backup` path.
func backupRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, constants.GitMapDir, "backup"), nil
}

// collectSnapshots walks the fix-repo backup layout and returns
// every leaf timestamp directory it can identify.
func collectSnapshots(root string) ([]backupSnapshot, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, nil
	}
	var out []backupSnapshot
	repos, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		if !repo.IsDir() {
			continue
		}
		out = append(out, walkRepoBackups(filepath.Join(root, repo.Name()), repo.Name())...)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ts.After(out[j].ts) })

	return out, nil
}

// walkRepoBackups descends `<root>/<repo>/v*/fix-repo/<ts>/`.
func walkRepoBackups(repoDir, repo string) []backupSnapshot {
	var out []backupSnapshot
	versions, _ := os.ReadDir(repoDir)
	for _, v := range versions {
		if !v.IsDir() {
			continue
		}
		fixDir := filepath.Join(repoDir, v.Name(), "fix-repo")
		stamps, _ := os.ReadDir(fixDir)
		for _, s := range stamps {
			if !s.IsDir() {
				continue
			}
			full := filepath.Join(fixDir, s.Name())
			ts := parseBackupTimestamp(s.Name())
			out = append(out, backupSnapshot{repo: repo, full: full, ts: ts, size: dirSize(full)})
		}
	}

	return out
}

// parseBackupTimestamp accepts the two common forms used in the
// backup tree; on any parse failure, falls back to dir modtime.
func parseBackupTimestamp(name string) time.Time {
	for _, layout := range []string{"20060102T150405Z", "2006-01-02T15-04-05Z"} {
		if t, err := time.Parse(layout, name); err == nil {
			return t
		}
	}

	return time.Time{}
}

// dirSize returns the total byte count under path (best-effort).
func dirSize(path string) int64 {
	var total int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			total += info.Size()
		}

		return nil
	})

	return total
}

// selectPruneVictims applies --keep + --older-than per repo.
func selectPruneVictims(snaps []backupSnapshot, keep, olderDays int) []backupSnapshot {
	byRepo := map[string][]backupSnapshot{}
	for _, s := range snaps {
		byRepo[s.repo] = append(byRepo[s.repo], s)
	}
	cutoff := time.Time{}
	if olderDays > 0 {
		cutoff = time.Now().AddDate(0, 0, -olderDays)
	}
	var victims []backupSnapshot
	for _, list := range byRepo {
		sort.Slice(list, func(i, j int) bool { return list[i].ts.After(list[j].ts) })
		for idx, s := range list {
			if keep > 0 && idx >= keep {
				victims = append(victims, s)

				continue
			}
			if olderDays > 0 && !s.ts.IsZero() && s.ts.Before(cutoff) {
				victims = append(victims, s)
			}
		}
	}

	return victims
}

// applyPrune deletes (or just prints) the selected snapshots.
func applyPrune(victims []backupSnapshot, dryRun bool) {
	if len(victims) == 0 {
		fmt.Fprintln(os.Stdout, "  nothing to prune")

		return
	}
	var freed int64
	for _, v := range victims {
		freed += v.size
		if dryRun {
			fmt.Fprintf(os.Stdout, "  \033[33mwould delete\033[0m %s (%s)\n", v.full, humanBytes(v.size))

			continue
		}
		if err := os.RemoveAll(v.full); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[31mfailed\033[0m %s: %v\n", v.full, err)

			continue
		}
		fmt.Fprintf(os.Stdout, "  \033[32mdeleted\033[0m %s (%s)\n", v.full, humanBytes(v.size))
	}
	verb := "freed"
	if dryRun {
		verb = "would free"
	}
	fmt.Fprintf(os.Stdout, "\n  %s: %d snapshot(s), %s\n\n", verb, len(victims), humanBytes(freed))
}

// printBackupTable renders the `ls` summary grouped by repo.
func printBackupTable(root string, snaps []backupSnapshot) {
	fmt.Fprintf(os.Stdout, "\n  \033[36m%s\033[0m  (%d snapshot(s))\n\n", root, len(snaps))
	byRepo := map[string][]backupSnapshot{}
	for _, s := range snaps {
		byRepo[s.repo] = append(byRepo[s.repo], s)
	}
	repos := make([]string, 0, len(byRepo))
	for r := range byRepo {
		repos = append(repos, r)
	}
	sort.Strings(repos)
	for _, r := range repos {
		list := byRepo[r]
		var total int64
		for _, s := range list {
			total += s.size
		}
		fmt.Fprintf(os.Stdout, "  \033[1m%s\033[0m — %d snapshot(s), %s\n", r, len(list), humanBytes(total))
		for _, s := range list {
			ts := s.ts.Format(time.RFC3339)
			if s.ts.IsZero() {
				ts = "(no timestamp)"
			}
			fmt.Fprintf(os.Stdout, "    • %s  %s\n", ts, humanBytes(s.size))
		}
	}
	fmt.Fprintf(os.Stdout, "\n  prune: `gitmap backup prune --keep=5` or `--older-than=30 --dry-run`\n\n")
}

// humanBytes formats n as a short human-readable size.
func humanBytes(n int64) string {
	const k = 1024.0
	if n < int64(k) {
		return fmt.Sprintf("%d B", n)
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	f := float64(n) / k
	for _, u := range units {
		if f < k {
			return fmt.Sprintf("%.1f %s", f, u)
		}
		f /= k
	}

	return fmt.Sprintf("%.1f EB", f)
}
