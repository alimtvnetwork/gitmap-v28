// Package cmd — chrome_diff.go: extension + bookmark diff between two
// Chrome profiles. Lists added/removed only (not modified).
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func runChromeDiff(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "chrome diff: ERROR usage: gitmap chrome diff <A> <B>")
		os.Exit(2)
	}
	a, okA := resolveChromeProfile(args[0])
	b, okB := resolveChromeProfile(args[1])
	if !okA || !okB {
		fmt.Fprintln(os.Stderr, "chrome diff: ERROR one or both profiles not found")
		printAvailableChromeProfilesWithDisplay()
		os.Exit(1)
	}
	fmt.Printf("\n\033[1;96m▸ chrome diff\033[0m  \033[1m%s\033[0m ↔ \033[1m%s\033[0m\n",
		chromeProfileSummary(a), chromeProfileSummary(b))

	extA, extB := listChromeExtensions(a.Path), listChromeExtensions(b.Path)
	printSetDiff("extensions", extA, extB)

	bmA, bmB := collectBookmarkURLs(a.Path), collectBookmarkURLs(b.Path)
	printSetDiff("bookmarks", bmA, bmB)
}

func listChromeExtensions(profile string) map[string]bool {
	out := map[string]bool{}
	entries, err := os.ReadDir(filepath.Join(profile, "Extensions"))
	if err != nil {
		return out
	}
	for _, e := range entries {
		if e.IsDir() {
			out[e.Name()] = true
		}
	}
	return out
}

func collectBookmarkURLs(profile string) map[string]bool {
	out := map[string]bool{}
	raw, err := os.ReadFile(filepath.Join(profile, "Bookmarks")) //nolint:gosec
	if err != nil {
		return out
	}
	var doc struct {
		Roots map[string]json.RawMessage `json:"roots"`
	}
	if json.Unmarshal(raw, &doc) != nil {
		return out
	}
	for _, r := range doc.Roots {
		walkBookmarkURLs(r, out)
	}
	return out
}

func walkBookmarkURLs(raw json.RawMessage, sink map[string]bool) {
	var node struct {
		URL      string            `json:"url"`
		Children []json.RawMessage `json:"children"`
	}
	if json.Unmarshal(raw, &node) != nil {
		return
	}
	if node.URL != "" {
		sink[node.URL] = true
	}
	for _, c := range node.Children {
		walkBookmarkURLs(c, sink)
	}
}

func printSetDiff(label string, a, b map[string]bool) {
	addOnly, bOnly := setSubtract(b, a), setSubtract(a, b)
	fmt.Printf("\n\033[1;94m%s\033[0m  A=%d B=%d  only-A=%d only-B=%d\n",
		label, len(a), len(b), len(bOnly), len(addOnly))
	for _, k := range sortedKeys(bOnly) {
		fmt.Printf("  \033[1;91m- A\033[0m %s\n", k)
	}
	for _, k := range sortedKeys(addOnly) {
		fmt.Printf("  \033[1;92m+ B\033[0m %s\n", k)
	}
}

func setSubtract(a, b map[string]bool) map[string]bool {
	out := map[string]bool{}
	for k := range a {
		if !b[k] {
			out[k] = true
		}
	}
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
