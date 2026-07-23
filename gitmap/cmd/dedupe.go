// Package cmd — `gitmap dedupe`: detect identical repos cloned under
// different folders by hashing each repo's HEAD tree SHA. v6.71.0
// adds parallel scanning and --format=json|csv export.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type headTreeEntry struct {
	path string
	sha  string
}

// runDedupe executes `gitmap dedupe`.
func runDedupe(args []string) {
	checkHelp("dedupe", args)
	fs := flag.NewFlagSet("dedupe", flag.ContinueOnError)
	root := fs.String("root", ".", "scan root directory")
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
	entries := mapReposParallel(repos, func(r string) (headTreeEntry, bool) {
		sha, ok := headTreeSHA(r)
		if !ok {
			return headTreeEntry{}, false
		}
		return headTreeEntry{path: r, sha: sha}, true
	})
	groups := map[string][]string{}
	for _, e := range entries {
		groups[e.sha] = append(groups[e.sha], e.path)
	}
	dupes := filterDuplicateGroups(groups)
	emitDedupe(dupes, fmtKind)
}

// headTreeSHA returns the tree SHA pointed to by HEAD for repo at dir.
func headTreeSHA(dir string) (string, bool) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD^{tree}").Output()
	if err != nil {
		return "", false
	}
	s := strings.TrimSpace(string(out))

	return s, s != ""
}

// filterDuplicateGroups keeps only groups with 2+ entries.
func filterDuplicateGroups(groups map[string][]string) map[string][]string {
	out := map[string][]string{}
	for k, v := range groups {
		if len(v) > 1 {
			out[k] = v
		}
	}

	return out
}

// emitDedupe dispatches to the requested output format.
func emitDedupe(dupes map[string][]string, f hygieneFormat) {
	keys := make([]string, 0, len(dupes))
	for k := range dupes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	switch f {
	case hygieneFormatJSON:
		type group struct {
			Tree  string   `json:"tree"`
			Paths []string `json:"paths"`
		}
		out := make([]group, 0, len(keys))
		for _, k := range keys {
			out = append(out, group{Tree: k, Paths: dupes[k]})
		}
		emitJSON(out)
	case hygieneFormatCSV:
		rows := [][]string{}
		for _, k := range keys {
			for _, p := range dupes[k] {
				rows = append(rows, []string{k, p})
			}
		}
		emitCSV([]string{"tree_sha", "path"}, rows)
	default:
		printDedupeReport(dupes)
	}
}

// printDedupeReport renders the duplicate-group table.
func printDedupeReport(dupes map[string][]string) {
	if len(dupes) == 0 {
		fmt.Fprint(os.Stdout, "\n  no duplicate repos found\n\n")

		return
	}
	keys := make([]string, 0, len(dupes))
	for k := range dupes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintf(os.Stdout, "\n  \033[36m%d duplicate group(s)\033[0m (identical HEAD tree)\n\n", len(dupes))
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "  \033[1mtree %s\033[0m\n", k[:12])
		for _, p := range dupes[k] {
			fmt.Fprintf(os.Stdout, "    • %s\n", p)
		}
	}
	fmt.Fprintln(os.Stdout, "")
}
