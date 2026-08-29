package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
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
func runRm(args []string) error {
	checkHelp("rm", args)
	yes, dbOnly, rest := parseRmFlags(args)
	targets := expandRmTargets(rest)
	if len(targets) == 0 {
		fmt.Fprint(os.Stderr, rmUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "rm: open db:")
	}
	defer db.Close()

	matches, missing := ResolveMultiRepos(db, targets)
	var finalMissing []string
	for _, m := range missing {
		abs, err := filepath.Abs(m)
		if err == nil && fsutil.DirExists(abs) {
			matches = append(matches, model.ScanRecord{
				AbsolutePath: abs,
				Slug:         filepath.Base(abs) + " (untracked)",
			})
			continue
		}
		finalMissing = append(finalMissing, m)
	}

	for _, m := range finalMissing {
		fmt.Fprintf(os.Stderr, "rm: no repo matched %q\n", m)
		PrintRepoSuggestions(db, m)
	}
	if len(matches) == 0 {
		return apperror.NewSimple("fatal error", "E9000")
	}
	if removeRmMatches(db, matches, yes, dbOnly) {
		os.Exit(0)
	}
	return apperror.NewSimple("fatal error", "E9000")
}

func parseRmFlags(args []string) (bool, bool, []string) {
	yes, dbOnly := false, false
	var out []string
	for _, a := range args {
		switch a {
		case "-y", "--yes":
			yes = true
		case "--db-only":
			dbOnly = true
		default:
			out = append(out, a)
		}
	}
	return yes, dbOnly, out
}

func expandRmTargets(args []string) []string {
	var out []string
	for _, a := range args {
		for _, p := range strings.Split(a, ",") {
			if p = strings.TrimSpace(p); p != "" {
				p = strings.TrimRight(p, "/\\")
				out = append(out, p)
			}
		}
	}
	return out
}

func removeRmMatches(db *store.DB, matches []model.ScanRecord, yes, dbOnly bool) bool {
	ok := true
	reader := bufio.NewReader(os.Stdin)
	for _, r := range matches {
		if !yes && !confirmRemove(reader, r, dbOnly) {
			fmt.Printf("skip: %s\n", r.Slug)
			continue
		}
		if err := removeRepoFully(db, r, dbOnly); err != nil {
			fmt.Fprintf(os.Stderr, "rm: %s: %v\n", r.Slug, err)
			ok = false
			continue
		}
		fmt.Printf("removed: %s (%s)\n", r.Slug, r.AbsolutePath)
	}
	return ok
}

func confirmRemove(r *bufio.Reader, rec model.ScanRecord, dbOnly bool) bool {
	action := "Delete folder and untrack"
	if dbOnly {
		action = "Untrack from database"
	}
	fmt.Printf("%s %s\n  %s ? [y/N] ", action, rec.Slug, rec.AbsolutePath)
	line, _ := r.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes"
}

func removeRepoFully(db *store.DB, r model.ScanRecord, dbOnly bool) error {
	if err := removeRepoDisk(r.AbsolutePath, dbOnly); err != nil {
		return err
	}
	if _, err := db.DeleteByPath(r.AbsolutePath); err != nil {
		return fmt.Errorf("db delete: %w", err)
	}
	_ = vscodepm.RemoveEntry(r.AbsolutePath)
	_ = desktop.RemoveRepo(r.AbsolutePath)
	return nil
}

func removeRepoDisk(absPath string, dbOnly bool) error {
	if dbOnly {
		return nil
	}
	if err := fsutil.SafeRemoveAll(absPath); err != nil {
		return fmt.Errorf("remove dir: %w", err)
	}
	return nil
}
