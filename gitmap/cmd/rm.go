package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// rmUsage describes the `gitmap rm` command.
const rmUsage = `Usage: gitmap rm [-y|--yes] <target>[,<target>...] [<target>...]
       gitmap remove ...
       gitmap del ...

Targets may be:
  - a repo slug/name           (e.g. my-repo)
  - a path                     (./projects/foo, .\macro-ahk, /abs/path)
  - a glob over slug or path   (macro*, gitmap-*)
  - comma-joined combinations  (macro*,gitmap*)

Default: prompts before deleting each repo folder on disk.
With -y/--yes: deletes on-disk folder and DB row without prompting.

Examples:
  gitmap rm my-repo
  gitmap rm macro*
  gitmap rm macro*,gitmap*
  gitmap rm -y macro*
  gitmap rm ./projects/foo ../bar
`

// runRm handles `gitmap rm`. Supports globs, comma-joined targets,
// the -y/--yes auto-confirm flag, and removes the on-disk folder in
// addition to the DB row.
func runRm(args []string) {
	checkHelp("rm", args)
	yes, rest := extractYesFlag(args)
	targets := expandRmTargets(rest)
	if len(targets) == 0 {
		fmt.Fprint(os.Stderr, rmUsage)
		os.Exit(1)
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "rm: open db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	matches := resolveRmMatches(db, targets)
	if len(matches) == 0 {
		os.Exit(1)
	}
	if removeRmMatches(db, matches, yes) {
		os.Exit(0)
	}
	os.Exit(1)
}

// extractYesFlag pulls -y/--yes out of args and returns the remainder.
func extractYesFlag(args []string) (bool, []string) {
	yes := false
	out := make([]string, 0, len(args))
	for _, a := range args {
		if a == "-y" || a == "--yes" {
			yes = true

			continue
		}
		out = append(out, a)
	}

	return yes, out
}

// expandRmTargets splits comma-joined tokens into individual targets.
func expandRmTargets(args []string) []string {
	var out []string
	for _, a := range args {
		for _, p := range strings.Split(a, ",") {
			if p = strings.TrimSpace(p); p != "" {
				out = append(out, p)
			}
		}
	}

	return out
}

// resolveRmMatches turns every target (literal/path/glob) into repo rows.
// Missing targets emit a warning. Duplicate repos are de-duped by ID.
func resolveRmMatches(db *store.DB, targets []string) []model.ScanRecord {
	all, err := db.ListRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "rm: list: %v\n", err)

		return nil
	}
	seen := map[int64]bool{}
	var out []model.ScanRecord
	for _, t := range targets {
		hits := matchTarget(all, t)
		if len(hits) == 0 {
			fmt.Fprintf(os.Stderr, "rm: no repo matched %q\n", t)

			continue
		}
		for _, r := range hits {
			if !seen[r.ID] {
				seen[r.ID] = true
				out = append(out, r)
			}
		}
	}

	return out
}

// matchTarget returns every repo whose slug or absolute path matches the
// target. Globs (containing * ? [) match via filepath.Match against the
// slug and the basename of the absolute path. Non-glob targets match by
// exact slug or absolute path.
func matchTarget(all []model.ScanRecord, target string) []model.ScanRecord {
	if isGlob(target) {
		return matchGlob(all, target)
	}
	if abs, err := filepath.Abs(target); err == nil {
		for _, r := range all {
			if r.AbsolutePath == abs {
				return []model.ScanRecord{r}
			}
		}
	}
	var out []model.ScanRecord
	for _, r := range all {
		if r.Slug == target {
			out = append(out, r)
		}
	}

	return out
}

// isGlob reports whether s contains any filepath.Match metacharacter.
func isGlob(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

// matchGlob returns rows whose slug or path-basename matches the pattern.
func matchGlob(all []model.ScanRecord, pat string) []model.ScanRecord {
	var out []model.ScanRecord
	for _, r := range all {
		if globHit(pat, r.Slug) || globHit(pat, filepath.Base(r.AbsolutePath)) {
			out = append(out, r)
		}
	}

	return out
}

// globHit wraps filepath.Match and swallows ErrBadPattern as a no-match.
func globHit(pat, name string) bool {
	ok, err := filepath.Match(pat, name)

	return err == nil && ok
}

// removeRmMatches deletes each matched repo from disk (with prompt unless
// yes) and from the DB. Returns true when every match was removed.
func removeRmMatches(db *store.DB, matches []model.ScanRecord, yes bool) bool {
	ok := true
	reader := bufio.NewReader(os.Stdin)
	for _, r := range matches {
		if !yes && !confirmRemove(reader, r) {
			fmt.Printf("skip: %s\n", r.Slug)

			continue
		}
		if err := removeRepoFully(db, r); err != nil {
			fmt.Fprintf(os.Stderr, "rm: %s: %v\n", r.Slug, err)
			ok = false

			continue
		}
		fmt.Printf("removed: %s (%s)\n", r.Slug, r.AbsolutePath)
	}

	return ok
}

// confirmRemove prompts the user [y/N] for a single repo.
func confirmRemove(r *bufio.Reader, rec model.ScanRecord) bool {
	fmt.Printf("Delete folder and untrack %s\n  %s ? [y/N] ", rec.Slug, rec.AbsolutePath)
	line, _ := r.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))

	return ans == "y" || ans == "yes"
}

// removeRepoFully removes the on-disk folder and the DB row.
func removeRepoFully(db *store.DB, r model.ScanRecord) error {
	if r.AbsolutePath != "" {
		if err := os.RemoveAll(r.AbsolutePath); err != nil {
			return fmt.Errorf("remove dir: %w", err)
		}
	}
	if _, err := db.DeleteByPath(r.AbsolutePath); err != nil {
		return fmt.Errorf("db delete: %w", err)
	}

	return nil
}
