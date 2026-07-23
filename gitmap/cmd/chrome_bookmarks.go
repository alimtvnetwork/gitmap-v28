// Package cmd — chrome_bookmarks.go: export a profile's Bookmarks
// file to md|html|json. Defaults to md on stdout.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type bookmarkItem struct {
	Title    string
	URL      string
	Folder   string
	Children []bookmarkItem
}

func runChromeExportBookmarks(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "chrome export-bookmarks: ERROR usage: gitmap chrome export-bookmarks <profile> [--format md|html|json] [--out <file>] [--root <bookmark_bar|other|synced>] [--folder <path/to/folder>] [--match <substr>] [--title <exact>]")
		os.Exit(2)
	}
	profile, ok := resolveChromeProfile(args[0])
	if !ok {
		fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR profile %q not found\n", args[0])
		printAvailableChromeProfilesWithDisplay()
		os.Exit(1)
	}
	format, outPath, rootName, folderPath, match, title := "md", "", "", "", "", ""
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--format", "-f":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		case "--out", "-o":
			if i+1 < len(args) {
				outPath = args[i+1]
				i++
			}
		case "--root", "-r":
			if i+1 < len(args) {
				rootName = args[i+1]
				i++
			}
		case "--folder":
			if i+1 < len(args) {
				folderPath = args[i+1]
				i++
			}
		case "--match":
			if i+1 < len(args) {
				match = args[i+1]
				i++
			}
		case "--title":
			if i+1 < len(args) {
				title = args[i+1]
				i++
			}
		}
	}
	roots := loadBookmarkRoots(profile.Path)
	if len(roots) == 0 {
		fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR no Bookmarks file found or it is empty/unreadable at %q\n  hint: open Chrome with this profile once so it writes %s\n", profile.Path, filepath.Join(profile.Path, "Bookmarks"))
		os.Exit(1)
	}
	available := availableRootNames(roots)
	if rootName != "" {
		filtered := filterBookmarkRoots(roots, rootName, "")
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR --root=%q did not match any top-level root\n  available roots: %s\n", rootName, strings.Join(available, ", "))
			os.Exit(1)
		}
		roots = filtered
	}
	if folderPath != "" {
		filtered := filterBookmarkRoots(roots, "", folderPath)
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR --folder=%q not found under root=%q\n  available top-level folders: %s\n  hint: paths are slash-delimited and case-insensitive (e.g. --folder \"Work/Docs\")\n", folderPath, fallback(rootName, "<all>"), strings.Join(topLevelFolderNames(roots), ", "))
			os.Exit(1)
		}
		roots = filtered
	}
	if match != "" || title != "" {
		filtered := filterBookmarksByTitle(roots, match, title)
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR no bookmarks matched --match=%q --title=%q within the selected subtree\n  hint: --match is a case-insensitive substring; --title is an exact title match\n", match, title)
			os.Exit(1)
		}
		roots = filtered
	}

	body, err := renderBookmarks(roots, format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR %v\n", err)
		os.Exit(1)
	}
	if outPath == "" {
		fmt.Print(body)
		return
	}
	if err := os.WriteFile(outPath, []byte(body), 0o644); err != nil { //nolint:gosec
		fmt.Fprintf(os.Stderr, "chrome export-bookmarks: ERROR write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[1;92m✓ wrote\033[0m %s (%d bytes)\n", outPath, len(body))
}

// filterBookmarkRoots narrows the tree to a top-level root (bookmark_bar,
// other, synced) and/or a slash-delimited folder subtree. Matching is
// case-insensitive on folder names. Empty filters are no-ops.
func filterBookmarkRoots(roots []bookmarkItem, rootName, folderPath string) []bookmarkItem {
	out := roots
	if rootName != "" {
		filtered := make([]bookmarkItem, 0, 1)
		for _, r := range roots {
			if strings.EqualFold(r.Folder, rootName) {
				filtered = append(filtered, r)
			}
		}
		out = filtered
	}
	if folderPath == "" {
		return out
	}
	parts := []string{}
	for _, p := range strings.Split(folderPath, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	matched := []bookmarkItem{}
	for _, r := range out {
		if sub, ok := findBookmarkFolder(r, parts); ok {
			matched = append(matched, sub)
		}
	}
	return matched
}

// findBookmarkFolder walks `parts` (case-insensitive) into the tree and
// returns the matching subtree.
func findBookmarkFolder(node bookmarkItem, parts []string) (bookmarkItem, bool) {
	if len(parts) == 0 {
		return node, true
	}
	for _, c := range node.Children {
		if c.URL != "" {
			continue
		}
		name := c.Title
		if name == "" {
			name = c.Folder
		}
		if strings.EqualFold(name, parts[0]) {
			return findBookmarkFolder(c, parts[1:])
		}
	}
	return bookmarkItem{}, false
}

func loadBookmarkRoots(profile string) []bookmarkItem {
	raw, err := os.ReadFile(filepath.Join(profile, "Bookmarks")) //nolint:gosec
	if err != nil {
		return nil
	}
	var doc struct {
		Roots map[string]json.RawMessage `json:"roots"`
	}
	if json.Unmarshal(raw, &doc) != nil {
		return nil
	}
	out := []bookmarkItem{}
	for name, r := range doc.Roots {
		item := parseBookmarkNode(r)
		item.Folder = name
		out = append(out, item)
	}
	return out
}

func parseBookmarkNode(raw json.RawMessage) bookmarkItem {
	var n struct {
		Name     string            `json:"name"`
		URL      string            `json:"url"`
		Children []json.RawMessage `json:"children"`
	}
	if json.Unmarshal(raw, &n) != nil {
		return bookmarkItem{}
	}
	item := bookmarkItem{Title: n.Name, URL: n.URL}
	for _, c := range n.Children {
		item.Children = append(item.Children, parseBookmarkNode(c))
	}
	return item
}

func renderBookmarks(roots []bookmarkItem, format string) (string, error) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(roots, "", "  ")
		return string(b) + "\n", err
	case "html":
		var sb strings.Builder
		sb.WriteString("<!DOCTYPE html><html><body>\n")
		for _, r := range roots {
			renderBookmarkHTML(&sb, r, 0)
		}
		sb.WriteString("</body></html>\n")
		return sb.String(), nil
	default: // md
		var sb strings.Builder
		for _, r := range roots {
			renderBookmarkMD(&sb, r, 0)
		}
		return sb.String(), nil
	}
}

func renderBookmarkMD(sb *strings.Builder, n bookmarkItem, depth int) {
	indent := strings.Repeat("  ", depth)
	switch {
	case n.URL != "":
		fmt.Fprintf(sb, "%s- [%s](%s)\n", indent, fallback(n.Title, n.URL), n.URL)
	case n.Title != "" || n.Folder != "":
		fmt.Fprintf(sb, "%s- **%s/**\n", indent, fallback(n.Title, n.Folder))
	}
	for _, c := range n.Children {
		renderBookmarkMD(sb, c, depth+1)
	}
}

func renderBookmarkHTML(sb *strings.Builder, n bookmarkItem, depth int) {
	if n.URL != "" {
		fmt.Fprintf(sb, "<dt><a href=\"%s\">%s</a></dt>\n", n.URL, fallback(n.Title, n.URL))
		return
	}
	fmt.Fprintf(sb, "<dt><h3>%s</h3><dl>\n", fallback(n.Title, n.Folder))
	for _, c := range n.Children {
		renderBookmarkHTML(sb, c, depth+1)
	}
	sb.WriteString("</dl></dt>\n")
}

func fallback(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func availableRootNames(roots []bookmarkItem) []string {
	out := make([]string, 0, len(roots))
	for _, r := range roots {
		out = append(out, r.Folder)
	}
	if len(out) == 0 {
		return []string{"(none)"}
	}
	return out
}

func topLevelFolderNames(roots []bookmarkItem) []string {
	out := []string{}
	for _, r := range roots {
		for _, c := range r.Children {
			if c.URL == "" {
				out = append(out, fallback(c.Title, c.Folder))
			}
		}
	}
	if len(out) == 0 {
		return []string{"(none)"}
	}
	return out
}
