// Package cmd — `gitmap size`: per-repo .git size report with --prune
// to run `git gc --aggressive` on the worst offenders. v6.71.0 adds
// parallel `.git` sizing and --format=json|csv export.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

type repoSize struct {
	path string
	size int64
}

// runSize executes `gitmap size`.
func runSize(args []string) {
	checkHelp("size", args)
	fs := flag.NewFlagSet("size", flag.ContinueOnError)
	root := fs.String("root", ".", "scan root directory")
	topN := fs.Int("top", 0, "show only top N largest (0 = all)")
	prune := fs.Bool("prune", false, "run `git gc --aggressive` on each listed repo")
	dryRun := fs.Bool("dry-run", false, "with --prune: list gc invocations without running them")
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
	sizes := mapReposParallel(repos, func(r string) (repoSize, bool) {
		return repoSize{path: r, size: dirSize(filepath.Join(r, ".git"))}, true
	})
	sort.Slice(sizes, func(i, j int) bool { return sizes[i].size > sizes[j].size })
	if *topN > 0 && len(sizes) > *topN {
		sizes = sizes[:*topN]
	}
	emitSize(sizes, fmtKind)
	if *prune {
		runAggressiveGC(sizes, *dryRun)
	}
}

// emitSize dispatches to the requested output format.
func emitSize(sizes []repoSize, f hygieneFormat) {
	switch f {
	case hygieneFormatJSON:
		type row struct {
			Path  string `json:"path"`
			Bytes int64  `json:"bytes"`
		}
		out := make([]row, 0, len(sizes))
		for _, s := range sizes {
			out = append(out, row{Path: s.path, Bytes: s.size})
		}
		emitJSON(out)
	case hygieneFormatCSV:
		rows := make([][]string, 0, len(sizes))
		for _, s := range sizes {
			rows = append(rows, []string{s.path, fmt.Sprintf("%d", s.size)})
		}
		emitCSV([]string{"path", "bytes"}, rows)
	default:
		printSizeReport(sizes)
	}
}

// printSizeReport renders the size table sorted desc.
func printSizeReport(sizes []repoSize) {
	if len(sizes) == 0 {
		fmt.Fprint(os.Stdout, "\n  no repos found\n\n")

		return
	}
	var total int64
	for _, s := range sizes {
		total += s.size
	}
	fmt.Fprintf(os.Stdout, "\n  \033[36m%d repo(s)\033[0m  total .git = %s\n\n", len(sizes), humanBytes(total))
	for _, s := range sizes {
		fmt.Fprintf(os.Stdout, "  \033[33m%10s\033[0m  %s\n", humanBytes(s.size), s.path)
	}
	fmt.Fprintln(os.Stdout, "")
}

// runAggressiveGC invokes `git gc --aggressive --prune=now` on each repo.
func runAggressiveGC(sizes []repoSize, dryRun bool) {
	for _, s := range sizes {
		if dryRun {
			fmt.Fprintf(os.Stdout, "  \033[33mwould gc\033[0m %s\n", s.path)

			continue
		}
		fmt.Fprintf(os.Stdout, "  \033[36mgc\033[0m %s ...\n", s.path)
		cmd := exec.Command("git", "-C", s.path, "gc", "--aggressive", "--prune=now")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[31mfailed\033[0m %s: %v\n", s.path, err)

			continue
		}
		after := dirSize(filepath.Join(s.path, ".git"))
		fmt.Fprintf(os.Stdout, "  \033[32mdone\033[0m %s  %s -> %s\n", s.path, humanBytes(s.size), humanBytes(after))
	}
}
