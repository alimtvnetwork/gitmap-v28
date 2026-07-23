// Package cmd — chromeprofile_merge.go: `gitmap chrome-profile-merge`
// (alias `cpm`). Merges selected slices of a source Chrome profile
// INTO a destination profile without clobbering destination values.
//
// Default policy (interactive): for every conflicting key/bookmark,
// prompt user to [k]eep, [o]verwrite, [a]ll-keep, [A]ll-overwrite,
// [q]uit. `--yes` auto-keeps destination on conflict; `--force`
// auto-overwrites destination with source.
//
// `--what` selects the slices to merge:
//
//	all         (default) — settings + bookmarks + extensions
//	settings    — Preferences + Secure Preferences (top-level keys)
//	bookmarks   — Bookmarks file (by GUID/url under bookmark_bar/other/synced)
//	extensions  — Extensions/ subdir (per-extension folder add-only)
package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type mergePolicy struct {
	autoKeep      bool
	autoOverwrite bool
	dryRun        bool
	reader        *bufio.Reader
}

type mergeStats struct {
	added, skipped, overwrote int
}

// runChromeProfileMerge implements `gitmap cpm`.
func runChromeProfileMerge(args []string) {
	checkHelp(constants.CmdChromeProfileMerge, args)
	fs := flag.NewFlagSet(constants.CmdChromeProfileMerge, flag.ExitOnError)
	what := fs.String("what", constants.ChromeMergeWhatAll, "settings|bookmarks|extensions|all")
	yes := fs.Bool("yes", false, "auto-keep destination on every conflict")
	fs.BoolVar(yes, "y", false, "auto-keep destination on every conflict (short)")
	force := fs.Bool("force", false, "auto-overwrite destination on every conflict")
	dryRun := fs.Bool("dry-run", false, "print plan but do not write")
	_ = fs.Parse(args)
	pos := fs.Args()
	if len(pos) < 2 {
		fmt.Fprint(os.Stderr, constants.ErrChromeMergeUsage)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	if !isKnownMergeWhat(*what) {
		fmt.Fprintf(os.Stderr, constants.ErrChromeMergeUnknown, *what)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	src, ok := resolveChromeProfile(pos[0])
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, pos[0], src.Path)
		printAvailableChromeProfilesWithDisplay()
		os.Exit(constants.ExitChromeProfileNotFound)
	}
	dst, ok := resolveChromeProfile(pos[1])
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, pos[1], dst.Path)
		printAvailableChromeProfilesWithDisplay()
		os.Exit(constants.ExitChromeProfileNotFound)
	}
	executeChromeProfileMerge(src, dst, *what, mergePolicy{
		autoKeep: *yes, autoOverwrite: *force, dryRun: *dryRun, reader: bufio.NewReader(os.Stdin),
	})
}

func executeChromeProfileMerge(src, dst chromeProfileResolution, what string, pol mergePolicy) {
	fmt.Printf(constants.MsgChromeMergeStart,
		chromeProfileSummary(src), chromeProfileSummary(dst), what, src.Path, dst.Path)
	if pol.dryRun {
		fmt.Print(constants.MsgChromeMergeDryRun)
	}
	var total mergeStats
	if what == constants.ChromeMergeWhatAll || what == constants.ChromeMergeWhatSettings {
		total = addStats(total, mergeChromeSettings(src.Path, dst.Path, &pol))
	}
	if what == constants.ChromeMergeWhatAll || what == constants.ChromeMergeWhatBookmarks {
		total = addStats(total, mergeChromeBookmarks(src.Path, dst.Path, &pol))
	}
	if what == constants.ChromeMergeWhatAll || what == constants.ChromeMergeWhatExtensions {
		total = addStats(total, mergeChromeExtensions(src.Path, dst.Path, &pol))
	}
	fmt.Printf(constants.MsgChromeMergeSummary, total.added, total.skipped, total.overwrote)
}

func addStats(a, b mergeStats) mergeStats {
	return mergeStats{a.added + b.added, a.skipped + b.skipped, a.overwrote + b.overwrote}
}

func isKnownMergeWhat(w string) bool {
	switch w {
	case constants.ChromeMergeWhatAll, constants.ChromeMergeWhatSettings,
		constants.ChromeMergeWhatBookmarks, constants.ChromeMergeWhatExtensions:
		return true
	}
	return false
}

// mergeChromeSettings merges top-level keys from src/Preferences into
// dst/Preferences. Nested objects are walked one level deep.
func mergeChromeSettings(srcDir, dstDir string, pol *mergePolicy) mergeStats {
	fmt.Printf(constants.MsgChromeMergeStepHdr, "settings (Preferences)")
	return mergeJSONFile(
		filepath.Join(srcDir, constants.ChromePreferencesFile),
		filepath.Join(dstDir, constants.ChromePreferencesFile),
		pol,
	)
}

func mergeJSONFile(srcPath, dstPath string, pol *mergePolicy) mergeStats {
	var stats mergeStats
	srcRoot, err := readJSONObject(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  skip: %v\n", err)
		return stats
	}
	dstRoot, err := readJSONObject(dstPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "  skip: %v\n", err)
		return stats
	}
	if dstRoot == nil {
		dstRoot = map[string]any{}
	}
	stats = mergeMapInto(srcRoot, dstRoot, "", pol)
	if pol.dryRun {
		return stats
	}
	if err := writeJSONObject(dstPath, dstRoot); err != nil {
		fmt.Fprintf(os.Stderr, "  write failed: %v\n", err)
	}
	return stats
}

func mergeMapInto(src, dst map[string]any, prefix string, pol *mergePolicy) mergeStats {
	var s mergeStats
	for k, v := range src {
		key := joinKey(prefix, k)
		existing, present := dst[k]
		if !present {
			if pol.dryRun {
				fmt.Printf(constants.MsgChromeMergeDryAdd, key)
			} else {
				dst[k] = v
			}
			s.added++
			continue
		}
		decision := resolveMergeConflict(key, existing, v, pol)
		if decision == mergeOverwrite {
			if pol.dryRun {
				fmt.Printf(constants.MsgChromeMergeDryOver, key)
			} else {
				dst[k] = v
			}
			s.overwrote++
		} else {
			if pol.dryRun && !jsonEqual(existing, v) {
				fmt.Printf(constants.MsgChromeMergeDryKeep, key)
			}
			s.skipped++
		}
	}
	return s
}

const (
	mergeKeep      = 0
	mergeOverwrite = 1
)

func resolveMergeConflict(key string, existing, incoming any, pol *mergePolicy) int {
	if jsonEqual(existing, incoming) {
		return mergeKeep
	}
	if pol.autoOverwrite {
		return mergeOverwrite
	}
	if pol.autoKeep || pol.dryRun {
		return mergeKeep
	}
	return promptMergeDecision(key, pol)
}

func promptMergeDecision(key string, pol *mergePolicy) int {
	fmt.Printf(constants.MsgChromeMergePrompt, key)
	line, _ := pol.reader.ReadString('\n')
	switch strings.TrimSpace(line) {
	case "o":
		return mergeOverwrite
	case "A":
		pol.autoOverwrite = true
		return mergeOverwrite
	case "a":
		pol.autoKeep = true
		return mergeKeep
	case "q":
		fmt.Fprintln(os.Stderr, "  aborted by user")
		os.Exit(0)
	}
	return mergeKeep
}

func joinKey(prefix, k string) string {
	if prefix == "" {
		return k
	}
	return prefix + "." + k
}

func readJSONObject(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path) //nolint:gosec // user-supplied path
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

func writeJSONObject(path string, m map[string]any) error {
	out, err := json.MarshalIndent(m, "", constants.JSONIndent)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		return err
	}
	return os.WriteFile(path, out, constants.FilePermission)
}

func jsonEqual(a, b any) bool {
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	return string(ja) == string(jb)
}

// mergeChromeBookmarks merges the `roots.{bookmark_bar,other,synced}`
// children of src/Bookmarks into dst/Bookmarks. De-dupes by url (for
// leaves) and by name (for folders).
func mergeChromeBookmarks(srcDir, dstDir string, pol *mergePolicy) mergeStats {
	fmt.Printf(constants.MsgChromeMergeStepHdr, "bookmarks")
	var stats mergeStats
	srcPath := filepath.Join(srcDir, constants.ChromeBookmarksFile)
	dstPath := filepath.Join(dstDir, constants.ChromeBookmarksFile)
	src, err := readJSONObject(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  skip: %v\n", err)
		return stats
	}
	dst, err := readJSONObject(dstPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "  skip: %v\n", err)
		return stats
	}
	if dst == nil {
		dst = map[string]any{"roots": map[string]any{}}
	}
	stats = mergeBookmarkRoots(src, dst, pol)
	if !pol.dryRun {
		if err := writeJSONObject(dstPath, dst); err != nil {
			fmt.Fprintf(os.Stderr, "  write failed: %v\n", err)
		}
	}
	return stats
}

func mergeBookmarkRoots(src, dst map[string]any, pol *mergePolicy) mergeStats {
	var s mergeStats
	srcRoots, _ := src["roots"].(map[string]any)
	dstRoots, _ := dst["roots"].(map[string]any)
	if dstRoots == nil {
		dstRoots = map[string]any{}
		dst["roots"] = dstRoots
	}
	for _, name := range []string{"bookmark_bar", "other", "synced"} {
		sn, _ := srcRoots[name].(map[string]any)
		dn, _ := dstRoots[name].(map[string]any)
		if sn == nil {
			continue
		}
		if dn == nil {
			dstRoots[name] = sn
			s.added++
			continue
		}
		s = addStats(s, mergeBookmarkFolder(sn, dn, name, pol))
	}
	return s
}

func mergeBookmarkFolder(src, dst map[string]any, label string, pol *mergePolicy) mergeStats {
	var s mergeStats
	srcChildren, _ := src["children"].([]any)
	dstChildren, _ := dst["children"].([]any)
	seen := bookmarkChildIndex(dstChildren)
	for _, c := range srcChildren {
		child, ok := c.(map[string]any)
		if !ok {
			continue
		}
		key := bookmarkKey(child)
		if _, dup := seen[key]; dup {
			if pol.autoOverwrite {
				s.overwrote++
			} else {
				s.skipped++
			}
			continue
		}
		dstChildren = append(dstChildren, child)
		seen[key] = struct{}{}
		s.added++
	}
	dst["children"] = dstChildren
	_ = label
	return s
}

func bookmarkChildIndex(children []any) map[string]struct{} {
	out := map[string]struct{}{}
	for _, c := range children {
		if m, ok := c.(map[string]any); ok {
			out[bookmarkKey(m)] = struct{}{}
		}
	}
	return out
}

func bookmarkKey(b map[string]any) string {
	if u, ok := b["url"].(string); ok && u != "" {
		return "u:" + u
	}
	if n, ok := b["name"].(string); ok {
		return "f:" + n
	}
	return ""
}

// mergeChromeExtensions copies any per-extension subdir under
// src/Extensions/ that isn't already present in dst/Extensions/.
// Existing destination extensions are left untouched (add-only).
func mergeChromeExtensions(srcDir, dstDir string, pol *mergePolicy) mergeStats {
	fmt.Printf(constants.MsgChromeMergeStepHdr, "extensions")
	var stats mergeStats
	srcRoot := filepath.Join(srcDir, "Extensions")
	dstRoot := filepath.Join(dstDir, "Extensions")
	entries, err := os.ReadDir(srcRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  skip: %v\n", err)
		return stats
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		stats = addStats(stats, mergeOneExtension(filepath.Join(srcRoot, e.Name()), filepath.Join(dstRoot, e.Name()), pol))
	}
	return stats
}

func mergeOneExtension(srcExt, dstExt string, pol *mergePolicy) mergeStats {
	var s mergeStats
	if _, err := os.Stat(dstExt); err == nil {
		if pol.autoOverwrite {
			if pol.dryRun {
				fmt.Printf(constants.MsgChromeMergeDryOver, "ext "+filepath.Base(dstExt))
			}
			s.overwrote++
		} else {
			if pol.dryRun {
				fmt.Printf(constants.MsgChromeMergeDryKeep, "ext "+filepath.Base(dstExt))
			}
			s.skipped++
		}
		return s
	}
	if pol.dryRun {
		fmt.Printf(constants.MsgChromeMergeDryAdd, "ext "+filepath.Base(dstExt))
		s.added++
		return s
	}
	if err := copyTree(srcExt, dstExt); err != nil {
		fmt.Fprintf(os.Stderr, "  copy %s failed: %v\n", filepath.Base(srcExt), err)
		return s
	}
	s.added++
	return s
}

// copyTree is a minimal recursive copy used only by extensions merge.
func copyTree(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if err := os.MkdirAll(dst, constants.DirPermission); err != nil {
			return err
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if err := copyTree(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	return copyOneFile(src, dst)
}

func copyOneFile(src, dst string) error {
	in, err := os.Open(src) //nolint:gosec // user-supplied path
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), constants.DirPermission); err != nil {
		return err
	}
	out, err := os.Create(dst) //nolint:gosec // user-supplied path
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
