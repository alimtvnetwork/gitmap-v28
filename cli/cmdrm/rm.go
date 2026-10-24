package cmdrm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// rmUsage describes the `gitmap rm` command.
const rmUsage = `Usage: gitmap rm [-y|--yes] <target>[,<target>...] [<target>...]
       gitmap rm <pattern> [--task <id>] [--reason <text>] [--undo] [--force]
       gitmap rm undo [<task_id>]
       gitmap rm list
       gitmap rm purge [--all] [--task <id>] [--older-than <duration>]
       gitmap remove ...
       gitmap del ...

Targets may be:
  - a repo slug/name           (e.g. my-repo)
  - a path                     (./projects/foo, .\macro-ahk, /abs/path)
  - a glob over slug or path   (macro*, gitmap-*)
  - comma-joined combinations  (macro*,gitmap*)
  - a file or file glob        (patch_*.py, temp_*.tmp)

Default: prompts before deleting each repo folder on disk.
With -y/--yes: deletes on-disk folder and DB row without prompting.

Examples:
  gitmap rm my-repo
  gitmap rm macro*
  gitmap rm "temp_test_*.tmp" --task "task-240-verify"
  gitmap rm undo task-240-verify
  gitmap rm list
  gitmap rm purge --all
`

// runRm handles `gitmap rm`. Supports globs, comma-joined targets,
// the -y/--yes auto-confirm flag, and removes the on-disk folder in
// addition to the DB row and .gitmap/output/gitmap.json entry.
func RunRm(args []string) error {
	checkHelp("rm", args)
	if isSafeRmDirective(args) {
		return RunSafeRmCLI(args)
	}

	yes, dbOnly, rest := parseRmFlags(args)
	targets := expandRmTargets(rest)
	if len(targets) == 0 {
		fmt.Fprint(os.Stderr, rmUsage)

		return apperror.NewValidationError("target repository required")
	}

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "rm: open db:")
	}

	defer db.Close()

	return executeRmTargets(db, targets, yes, dbOnly, args)
}

func executeRmTargets(db *store.DB, targets []string, isYes, isDbOnly bool, origArgs []string) error {
	matches, missing := resolveRmMatches(db, targets)
	if len(matches) == 0 && isSafeRmFallbackAllowed(targets) {
		return RunSafeRmCLI(origArgs)
	}

	reportMissingRmTargets(db, missing)
	if len(matches) == 0 {
		return nil
	}

	_ = removeRmMatches(db, matches, isYes, isDbOnly)
	return nil
}

func isSafeRmDirective(args []string) bool {
	if len(args) == 0 {
		return false
	}

	first := strings.ToLower(args[0])
	if first == "undo" || first == "list" || first == "purge" {
		return true
	}

	for _, a := range args {
		if a == "--undo" || a == "-u" || a == "--purge" || a == "--list" {
			return true
		}
		if a == "--task" || a == "-t" || a == "--reason" || a == "-r" {
			return true
		}
		if strings.HasPrefix(a, "--task=") || strings.HasPrefix(a, "--reason=") || strings.HasPrefix(a, "--older-than") {
			return true
		}
	}

	return false
}

func isSafeRmFallbackAllowed(targets []string) bool {
	expanded, err := ExpandFilePatterns(targets, ".")
	return err == nil && len(expanded) > 0
}

func resolveRmMatches(db *store.DB, targets []string) ([]model.ScanRecord, []string) {
	matches, missing := ResolveMultiRepos(db, targets)
	extraMatches, finalMissing := resolveDiskTargets(missing)
	matches = append(matches, extraMatches...)

	return matches, finalMissing
}

func resolveDiskTargets(missing []string) ([]model.ScanRecord, []string) {
	var matches []model.ScanRecord
	var stillMissing []string
	for _, m := range missing {
		rec, isFound := checkUntrackedDir(m)
		if isFound {
			matches = append(matches, rec)
			continue
		}

		stillMissing = append(stillMissing, m)
	}

	return matches, stillMissing
}

func checkUntrackedDir(target string) (model.ScanRecord, bool) {
	abs, err := filepath.Abs(target)
	isDir := err == nil && fsutil.DirExists(abs)
	if !isDir {
		return model.ScanRecord{}, false
	}

	return model.ScanRecord{
		AbsolutePath: abs,
		Slug:         filepath.Base(abs) + " (untracked)",
	}, true
}

func reportMissingRmTargets(db *store.DB, missing []string) {
	for _, m := range missing {
		fmt.Fprintf(os.Stderr, "rm: no repo matched %q\n", m)
		PrintRepoSuggestions(db, m)
	}
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
		tokens := strings.Split(a, ",")
		out = appendValidTargets(out, tokens)
	}

	return out
}

func appendValidTargets(out, tokens []string) []string {
	for _, p := range tokens {
		clean := strings.TrimRight(strings.TrimSpace(p), "/\\")
		if clean != "" {
			out = append(out, clean)
		}
	}

	return out
}

func removeRmMatches(db *store.DB, matches []model.ScanRecord, isYes, isDbOnly bool) bool {
	reader := bufio.NewReader(os.Stdin)
	isSuccess := true
	for _, r := range matches {
		if !processSingleRm(db, reader, r, isYes, isDbOnly) {
			isSuccess = false
		}
	}

	return isSuccess
}

func processSingleRm(db *store.DB, reader *bufio.Reader, r model.ScanRecord, isYes, isDbOnly bool) bool {
	isConfirmed := isYes || confirmRemove(reader, r, isDbOnly)
	if !isConfirmed {
		fmt.Printf("skip: %s\n", r.Slug)

		return true
	}

	if err := removeRepoFully(db, r, isDbOnly); err != nil {
		fmt.Fprintf(os.Stderr, "rm: %s: %v\n", r.Slug, err)

		return false
	}

	fmt.Printf("removed: %s (%s)\n", r.Slug, r.AbsolutePath)

	return true
}

func confirmRemove(r *bufio.Reader, rec model.ScanRecord, isDbOnly bool) bool {
	action := resolveConfirmAction(isDbOnly)
	fmt.Printf("%s %s\n  %s ? [y/N] ", action, rec.Slug, rec.AbsolutePath)
	line, _ := r.ReadString('\n')
	ans := strings.TrimSpace(line)

	return strings.EqualFold(ans, "y") || strings.EqualFold(ans, "yes")
}

func resolveConfirmAction(isDbOnly bool) string {
	if isDbOnly {
		return "Untrack from database"
	}

	return "Delete folder and untrack"
}

func removeRepoFully(db *store.DB, r model.ScanRecord, isDbOnly bool) *apperror.AppError {
	if err := removeRepoDisk(r.AbsolutePath, isDbOnly); err != nil {
		return err
	}

	if err := removeRepoDB(db, r); err != nil {
		return err
	}

	removeRepoFromJSON(r)
	_ = vscodepm.RemoveEntry(r.AbsolutePath)
	_ = desktop.RemoveRepo(r.AbsolutePath)

	return nil
}
