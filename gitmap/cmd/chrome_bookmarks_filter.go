// Package cmd — chrome_bookmarks_filter.go: --match / --title pruning.
// Walks the tree post-folder-filter and keeps only branches that contain
// at least one matching URL leaf; folders themselves are preserved as
// context so the rendered output retains hierarchy.
package cmd

import "strings"

// filterBookmarksByTitle prunes the tree, keeping leaves whose Title
// matches substring `match` (case-insensitive) and/or equals `exactTitle`
// (case-sensitive). Either may be empty; if both are empty the input is
// returned unchanged. Folders that end up with zero matching descendants
// are dropped.
func filterBookmarksByTitle(roots []bookmarkItem, match, exactTitle string) []bookmarkItem {
	if match == "" && exactTitle == "" {
		return roots
	}
	out := make([]bookmarkItem, 0, len(roots))
	for _, r := range roots {
		if pruned, ok := pruneBookmark(r, strings.ToLower(match), exactTitle); ok {
			out = append(out, pruned)
		}
	}
	return out
}

func pruneBookmark(n bookmarkItem, matchLower, exact string) (bookmarkItem, bool) {
	if n.URL != "" {
		if bookmarkMatches(n.Title, matchLower, exact) {
			return n, true
		}
		return bookmarkItem{}, false
	}
	kept := make([]bookmarkItem, 0, len(n.Children))
	for _, c := range n.Children {
		if p, ok := pruneBookmark(c, matchLower, exact); ok {
			kept = append(kept, p)
		}
	}
	if len(kept) == 0 {
		return bookmarkItem{}, false
	}
	n.Children = kept
	return n, true
}

func bookmarkMatches(title, matchLower, exact string) bool {
	if exact != "" && title == exact {
		return true
	}
	if matchLower != "" && strings.Contains(strings.ToLower(title), matchLower) {
		return true
	}
	return false
}
