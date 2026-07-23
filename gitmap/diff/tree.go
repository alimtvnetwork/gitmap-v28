package diff

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// EntryKind classifies how a single relative path differs.
type EntryKind int

const (
	// MissingRight = present on LEFT, absent on RIGHT.
	MissingRight EntryKind = iota
	// MissingLeft = present on RIGHT, absent on LEFT.
	MissingLeft
	// Identical = present on both, byte-equal.
	Identical
	// Conflict = present on both, byte-different.
	Conflict
)

// Entry is one path classified for the diff report.
type Entry struct {
	RelPath    string    `json:"path"`
	Kind       EntryKind `json:"-"`
	KindLabel  string    `json:"kind"`
	LeftSize   int64     `json:"left_size,omitempty"`
	RightSize  int64     `json:"right_size,omitempty"`
	LeftMTime  int64     `json:"left_mtime,omitempty"`
	RightMTime int64     `json:"right_mtime,omitempty"`
}

// WalkOptions controls which files participate in the diff.
type WalkOptions struct {
	IncludeVCS         bool
	IncludeNodeModules bool
}

// DiffTrees walks LEFT and RIGHT and returns every classified entry,
// sorted ascending by relative path.
func DiffTrees(leftDir, rightDir string, opts WalkOptions) ([]Entry, error) {
	leftIdx, err := indexTree(leftDir, opts)
	if err != nil {
		return nil, err
	}
	rightIdx, err := indexTree(rightDir, opts)
	if err != nil {
		return nil, err
	}

	return classifyAll(leftIdx, rightIdx, leftDir, rightDir), nil
}

// indexTree returns rel-path -> os.FileInfo for every non-ignored file.
func indexTree(root string, opts WalkOptions) (map[string]os.FileInfo, error) {
	out := make(map[string]os.FileInfo)
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return relErr
		}
		if shouldIgnore(rel, opts) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
		if info.IsDir() {
			return nil
		}
		out[filepath.ToSlash(rel)] = info

		return nil
	})

	return out, walkErr
}

// shouldIgnore returns true when rel matches the default ignore list.
func shouldIgnore(rel string, opts WalkOptions) bool {
	base := filepath.Base(rel)
	if !opts.IncludeVCS && base == ".git" {
		return true
	}
	if !opts.IncludeNodeModules && base == "node_modules" {
		return true
	}
	relSlash := filepath.ToSlash(rel)

	return strings.HasPrefix(relSlash, ".gitmap/release-assets/")
}

// classifyAll merges the two indexes into a sorted Entry list.
func classifyAll(left, right map[string]os.FileInfo, leftDir, rightDir string) []Entry {
	seen := unionKeys(left, right)
	out := make([]Entry, 0, len(seen))
	for _, rel := range seen {
		out = append(out, classifyOne(rel, left[rel], right[rel], leftDir, rightDir))
	}

	return out
}

// unionKeys returns the sorted union of map keys.
func unionKeys(a, b map[string]os.FileInfo) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}

// classifyOne produces a single Entry for one relative path.
func classifyOne(rel string, l, r os.FileInfo, leftDir, rightDir string) Entry {
	entry := Entry{RelPath: rel}
	if l != nil {
		entry.LeftSize, entry.LeftMTime = l.Size(), l.ModTime().Unix()
	}
	if r != nil {
		entry.RightSize, entry.RightMTime = r.Size(), r.ModTime().Unix()
	}
	entry.Kind = pickKind(l, r, leftDir, rightDir, rel)
	entry.KindLabel = labelFor(entry.Kind)

	return entry
}

// pickKind decides MissingLeft / MissingRight / Identical / Conflict.
func pickKind(l, r os.FileInfo, leftDir, rightDir, rel string) EntryKind {
	switch {
	case l != nil && r == nil:
		return MissingRight
	case l == nil && r != nil:
		return MissingLeft
	case sameContent(filepath.Join(leftDir, rel), filepath.Join(rightDir, rel)):
		return Identical
	}

	return Conflict
}

// labelFor returns a stable string label used in JSON / text output.
func labelFor(k EntryKind) string {
	switch k {
	case MissingRight:
		return "missing_right"
	case MissingLeft:
		return "missing_left"
	case Identical:
		return "identical"
	case Conflict:
		return "conflict"
	}

	return "unknown"
}

// sameContent returns true when both files have equal SHA-256.
func sameContent(a, b string) bool {
	ha, errA := hashFile(a)
	hb, errB := hashFile(b)
	if errA != nil || errB != nil {
		return false
	}

	return ha == hb
}

// hashFile streams a file through SHA-256.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
