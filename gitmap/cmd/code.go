package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// runCode implements `gitmap code` and its `paths` subcommand.
//
//	gitmap code                              -> register CWD / git root, alias = basename
//	gitmap code <alias>                      -> override alias, same path resolution
//	gitmap code <alias> <root> [extra...]    -> any path + variadic extras (additive upsert)
//	gitmap code paths add <alias> <path>     -> add an extra root to an existing entry
//	gitmap code paths rm  <alias> <path>     -> remove an extra root
//	gitmap code paths list <alias>           -> print attached extras
//
// The first three forms launch VS Code on the resolved root. The `paths`
// subcommand never opens VS Code — it only mutates the registry.
func runCode(args []string) {
	checkHelp(constants.CmdCode, args)

	if len(args) > 0 && args[0] == "paths" {
		runCodePaths(args[1:])

		return
	}

	alias, rootPath, extras := parseCodeArgs(args)

	resolved, err := resolveCodeRootPath(rootPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if alias == "" {
		alias = filepath.Base(resolved)
	}

	upsertCodeEntry(resolved, alias)
	if len(extras) > 0 {
		appendCodePathsToDB(resolved, extras)
	}

	syncCodeEntry(resolved, alias, extras)
	openInVSCode(resolved)
}

// parseCodeArgs returns (alias, rootPath, extraPaths).
//
//	0 args  -> ("", "", nil)              auto-resolve, basename alias
//	1 arg   -> (args[0], "", nil)         alias override, auto-resolve
//	2 args  -> (args[0], args[1], nil)    alias + explicit root
//	3+ args -> (args[0], args[1], args[2:]) alias + root + extras
func parseCodeArgs(args []string) (alias, rootPath string, extras []string) {
	switch len(args) {
	case 0:
		return "", "", nil
	case 1:
		return args[0], "", nil
	case 2:
		return args[0], args[1], nil
	default:
		return args[0], args[1], args[2:]
	}
}

// resolveCodeRootPath picks the rootPath per the documented precedence.
func resolveCodeRootPath(pathArg string) (string, error) {
	if pathArg != "" {
		return absoluteExisting(pathArg)
	}

	if root, err := gitTopLevel(); err == nil {
		return root, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine current directory: %w", err)
	}

	return absoluteExisting(cwd)
}

// absoluteExisting cleans the path, returns it absolute, and verifies it
// exists. Non-existent paths are an error so we never write garbage rows.
func absoluteExisting(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("cannot resolve absolute path %q: %w", p, err)
	}

	if _, err := os.Stat(abs); err != nil {
		return "", fmt.Errorf("path does not exist: %s", abs)
	}

	return abs, nil
}

// upsertCodeEntry persists the row into the VSCodeProject table.
func upsertCodeEntry(rootPath, name string) {
	db, err := openCodeDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if err := db.UpsertVSCodeProject(rootPath, name); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// appendCodePathsToDB UNIONs `extras` into the DB-side Paths list of an
// already-upserted row. Extras are resolved to absolute existing paths
// first; missing paths exit non-zero so the user catches typos early.
func appendCodePathsToDB(rootPath string, extras []string) {
	resolvedExtras := resolveExtras(extras)

	db, err := openCodeDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer db.Close()

	row, err := db.FindVSCodeProjectByPath(rootPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	merged := mergeStringPaths(row.Paths, resolvedExtras)
	if len(merged) == len(row.Paths) {
		return
	}

	if err := db.SetVSCodeProjectPaths(rootPath, merged); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// resolveExtras absolutifies each extra path or exits on the first miss.
func resolveExtras(extras []string) []string {
	out := make([]string, 0, len(extras))

	for _, raw := range extras {
		abs, err := absoluteExisting(raw)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		out = append(out, abs)
	}

	return out
}

// mergeStringPaths returns the order-preserving union of `existing` and
// `incoming`. OS-aware key (case-insensitive on Windows).
func mergeStringPaths(existing, incoming []string) []string {
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	out := make([]string, 0, len(existing)+len(incoming))

	for _, p := range existing {
		key := pathKey(p)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}

	for _, p := range incoming {
		key := pathKey(p)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}

	return out
}

// pathKey mirrors vscodepm.normalizePath but lives in the cmd package so
// we don't import an unexported helper.
func pathKey(p string) string {
	if filepath.Separator == '\\' {
		return strings.ToLower(filepath.Clean(p))
	}

	return filepath.Clean(p)
}

// syncCodeEntry pushes the (rootPath, name, paths, auto-tags) tuple into
// projects.json. Auto-tags are derived from rootPath's filesystem markers
// and UNIONed with any user-edited tags. Soft-fails when VS Code or the
// extension is missing.
func syncCodeEntry(rootPath, name string, extras []string) {
	resolved := resolveExtras(extras)

	summary, err := vscodepm.Sync([]vscodepm.Pair{{
		RootPath: rootPath,
		Name:     name,
		Paths:    resolved,
		Tags:     vscodepm.DetectTagsCustom(rootPath),
	}})
	if err != nil {
		reportVSCodePMSoftError(err)

		return
	}

	fmt.Printf(constants.MsgVSCodePMSyncSummary,
		summary.Added, summary.Updated, summary.Unchanged, summary.Total)
}

// openCodeDB opens the default DB and runs Migrate(). Centralized so every
// `code` entry-point follows the same setup path.
func openCodeDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, fmt.Errorf(constants.MsgDBUpsertFailed, err)
	}

	if err := db.Migrate(); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(constants.MsgDBUpsertFailed, err)
	}

	return db, nil
}

// runCodePaths dispatches the `code paths` subcommand.
func runCodePaths(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gitmap code paths <add|rm|list> <alias> [path]")
		os.Exit(2)
	}

	op := args[0]
	rest := args[1:]

	switch op {
	case "add":
		runCodePathsAdd(rest)
	case "rm", "remove":
		runCodePathsRm(rest)
	case "list", "ls":
		runCodePathsList(rest)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\nusage: gitmap code paths <add|rm|list> <alias> [path]\n", op)
		os.Exit(2)
	}
}

// runCodePathsAdd attaches one extra path to the entry matching <alias>.
func runCodePathsAdd(args []string) {
	alias, extra := requireAliasAndPath(args, "add")
	row := lookupAlias(alias)
	abs, err := absoluteExisting(extra)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	merged := mergeStringPaths(row.Paths, []string{abs})
	if len(merged) == len(row.Paths) {
		fmt.Printf(constants.MsgVSCodePMPathsExists, alias, abs)

		return
	}

	persistAliasPaths(row.RootPath, alias, merged)
	syncAliasEntry(row.RootPath, row.Name, merged)
	fmt.Printf(constants.MsgVSCodePMPathsAdded, alias, abs)
}

// runCodePathsRm detaches one extra path from the entry matching <alias>.
// projects.json is then re-synced so the user-visible UI reflects the drop.
func runCodePathsRm(args []string) {
	alias, extra := requireAliasAndPath(args, "rm")
	row := lookupAlias(alias)
	abs, err := filepath.Abs(extra)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	dropKey := pathKey(abs)
	pruned := make([]string, 0, len(row.Paths))
	dropped := false

	for _, p := range row.Paths {
		if pathKey(p) == dropKey {
			dropped = true

			continue
		}
		pruned = append(pruned, p)
	}

	if !dropped {
		fmt.Printf(constants.MsgVSCodePMPathsMissing, alias, abs)

		return
	}

	persistAliasPaths(row.RootPath, alias, pruned)
	overwriteAliasEntry(row.RootPath, row.Name, pruned)
	fmt.Printf(constants.MsgVSCodePMPathsRemoved, alias, abs)
}

// runCodePathsList prints the rootPath + all attached extras for <alias>.
func runCodePathsList(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gitmap code paths list <alias>")
		os.Exit(2)
	}

	alias := args[0]
	row := lookupAlias(alias)
	pathsCsv := strings.Join(row.Paths, ", ")
	if pathsCsv == "" {
		pathsCsv = "(none)"
	}

	fmt.Printf(constants.MsgVSCodePMPathsList, alias, row.Name, row.RootPath, pathsCsv)
	if len(row.Paths) == 0 {
		fmt.Print(constants.MsgVSCodePMPathsNone)
	}
}

// requireAliasAndPath enforces the two-arg contract for add/rm.
func requireAliasAndPath(args []string, op string) (string, string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: gitmap code paths %s <alias> <path>\n", op)
		os.Exit(2)
	}

	return args[0], args[1]
}

// lookupAlias loads the VSCodeProject row for the given alias name. Exits
// non-zero with an actionable hint when no row matches.
func lookupAlias(alias string) (row aliasRow) {
	db, err := openCodeDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer db.Close()

	found, err := db.FindVSCodeProjectByName(alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, constants.ErrVSCodePMAliasNotFound, alias, alias)
			fmt.Fprintln(os.Stderr)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return aliasRow{RootPath: found.RootPath, Name: found.Name, Paths: found.Paths}
}

// aliasRow is a small subset of model.VSCodeProject scoped to the fields
// the `paths` subcommand actually needs.
type aliasRow struct {
	RootPath string
	Name     string
	Paths    []string
}

// persistAliasPaths writes the new Paths slice to the DB.
func persistAliasPaths(rootPath, alias string, paths []string) {
	db, err := openCodeDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if err := db.SetVSCodeProjectPaths(rootPath, paths); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", alias, err)
		os.Exit(1)
	}
}

// syncAliasEntry pushes a UNION-style upsert (preserves user-added paths
// and tags). Auto-derived tags from the rootPath are added on every call.
func syncAliasEntry(rootPath, name string, paths []string) {
	if _, err := vscodepm.Sync([]vscodepm.Pair{{
		RootPath: rootPath, Name: name, Paths: paths,
		Tags: vscodepm.DetectTagsCustom(rootPath),
	}}); err != nil {
		reportVSCodePMSoftError(err)
	}
}

// overwriteAliasEntry forces the Paths field to the supplied slice (used by
// `rm` so the deleted path is actually removed from projects.json, not
// re-unioned back in).
func overwriteAliasEntry(rootPath, name string, paths []string) {
	if err := vscodepm.OverwritePaths(rootPath, name, paths); err != nil {
		reportVSCodePMSoftError(err)
	}
}
