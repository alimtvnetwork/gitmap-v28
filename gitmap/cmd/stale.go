// Package cmd — `gitmap stale` (sta): list local repos with no commits
// in N days, with optional --archive to move them to .gitmap/archive/.
// v6.68.0.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type staleRepo struct {
	path string
	last time.Time
}

// runStale executes `gitmap stale`.
func runStale(args []string) {
	checkHelp("stale", args)
	fs := flag.NewFlagSet("stale", flag.ContinueOnError)
	days := fs.Int("days", 90, "report repos with no commits in the last N days")
	root := fs.String("root", ".", "scan root directory")
	archive := fs.Bool("archive", false, "move stale repos into .gitmap/archive/")
	dryRun := fs.Bool("dry-run", false, "preview archive moves without touching disk")
	format := fs.String("format", "table", "output format: table|json|csv")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	fmtKind, err := parseHygieneFormat(*format)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	repos := scanForReposParallel(*root)
	cutoff := time.Now().AddDate(0, 0, -*days)
	probed := mapReposParallel(repos, func(r string) (staleRepo, bool) {
		t, ok := lastCommitTime(r)
		if !ok || !t.Before(cutoff) {
			return staleRepo{}, false
		}
		return staleRepo{path: r, last: t}, true
	})
	sort.Slice(probed, func(i, j int) bool { return probed[i].last.Before(probed[j].last) })
	emitStale(probed, *days, fmtKind)
	if *archive {
		archiveStaleRepos(probed, *dryRun)
	}
}

// emitStale dispatches to the requested output format.
func emitStale(stale []staleRepo, days int, f hygieneFormat) {
	switch f {
	case hygieneFormatJSON:
		type row struct {
			Path    string `json:"path"`
			LastUTC string `json:"last_commit_utc"`
			AgeDays int    `json:"age_days"`
		}
		rows := make([]row, 0, len(stale))
		now := time.Now()
		for _, s := range stale {
			rows = append(rows, row{
				Path:    s.path,
				LastUTC: s.last.UTC().Format(time.RFC3339),
				AgeDays: int(now.Sub(s.last).Hours() / 24),
			})
		}
		emitJSON(rows)
	case hygieneFormatCSV:
		now := time.Now()
		rows := make([][]string, 0, len(stale))
		for _, s := range stale {
			age := int(now.Sub(s.last).Hours() / 24)
			rows = append(rows, []string{s.path, s.last.UTC().Format(time.RFC3339), fmt.Sprintf("%d", age)})
		}
		emitCSV([]string{"path", "last_commit_utc", "age_days"}, rows)
	default:
		printStaleTable(stale, days)
	}
}

// scanForRepos returns directories under root that contain a .git folder.
func scanForRepos(root string) []string {
	var out []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		full := filepath.Join(root, e.Name())
		if isGitRepo(full) {
			out = append(out, full)
		}
	}

	return out
}

// (isGitRepo lives in githubdesktop.go and is reused here.)

// lastCommitTime returns the last commit time for repo at dir.
func lastCommitTime(dir string) (time.Time, bool) {
	cmd := exec.Command("git", "-C", dir, "log", "-1", "--format=%ct")
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, false
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return time.Time{}, false
	}
	var sec int64
	if _, err := fmt.Sscanf(s, "%d", &sec); err != nil {
		return time.Time{}, false
	}

	return time.Unix(sec, 0), true
}

// printStaleTable renders the stale-repo list.
func printStaleTable(stale []staleRepo, days int) {
	if len(stale) == 0 {
		fmt.Fprintf(os.Stdout, "\n  no repos stale beyond %d days\n\n", days)

		return
	}
	fmt.Fprintf(os.Stdout, "\n  \033[36m%d stale repo(s)\033[0m (no commits in %d days)\n\n", len(stale), days)
	now := time.Now()
	for _, s := range stale {
		age := int(now.Sub(s.last).Hours() / 24)
		fmt.Fprintf(os.Stdout, "  \033[33m%4dd\033[0m  %s  (last: %s)\n", age, s.path, s.last.Format("2006-01-02"))
	}
	fmt.Fprintln(os.Stdout, "")
}

// archiveStaleRepos moves repos into .gitmap/archive/<basename>-<ts>/.
func archiveStaleRepos(stale []staleRepo, dryRun bool) {
	if len(stale) == 0 {
		return
	}
	stamp := time.Now().UTC().Format("20060102T150405Z")
	archiveDir := filepath.Join(".gitmap", "archive", stamp)
	for _, s := range stale {
		dest := filepath.Join(archiveDir, filepath.Base(s.path))
		if dryRun {
			fmt.Fprintf(os.Stdout, "  \033[33mwould archive\033[0m %s -> %s\n", s.path, dest)

			continue
		}
		if err := os.MkdirAll(archiveDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "  archive mkdir: %v\n", err)

			return
		}
		if err := os.Rename(s.path, dest); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[31mfailed\033[0m %s: %v\n", s.path, err)

			continue
		}
		fmt.Fprintf(os.Stdout, "  \033[32marchived\033[0m %s -> %s\n", s.path, dest)
	}
}
