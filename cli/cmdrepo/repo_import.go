package cmdrepo

// repo_import.go — `gitmap repo import <file> --into <dir>` / `gitmap repoi`:
// consume a `repo export` manifest: clone, restore remotes, register in the
// gitmap DB (program 275, proposal P-06). Side-effect command: supports
// --dry-run (shared-patterns spec §5). One concern per file (≤300 lines).

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// repo import messages.
const (
	repoImportName      = "repo import"
	repoImportFlagInto  = "into"
	repoImportFlagDry   = "dry-run"
	repoImportUsageLine = "Usage: gitmap repo import <file> --into <dir> [--dry-run]\n       gitmap repoi <file> --into <dir> [same flags]\n"
	repoImportErrFile   = "missing required argument: <file> (manifest JSON from `gitmap repo export`)"
	repoImportErrInto   = "missing required flag: --into <dir> (destination directory for the clone)"
	repoImportErrParse  = "could not parse %s as a repo manifest: %v"
	repoImportErrURL    = "manifest %s has no clone URL (not a repo export manifest?)"
	repoImportErrExists = "destination %s already exists — remove it or pick another --into dir"

	repoImportDryRunFmt = "[repo import] dry-run — no changes made\n  manifest:   %s\n  would run:  git clone %s %s\n  remotes:    restore %d (%s)\n  register:   upsert %q into the gitmap database\n"
	repoImportResultFmt = "cloned → %s · registered in gitmap.db · remotes restored (%d) · scan updated\n"
)

// repoImportOpts carries the parsed flags for `repo import`.
type repoImportOpts struct {
	file  string
	into  string
	isDry bool
}

// RunRepoImport implements `gitmap repo import` / `gitmap repoi` for rootdata wiring.
func RunRepoImport(args []string) error {
	if appErr := runRepoImport(args); appErr != nil {
		return appErr
	}
	return nil
}

// runRepoImport executes the import flow: parse, read, clone, restore, register.
func runRepoImport(args []string) *appfault.AppError {
	if isRepoHelpArg(args) {
		fmt.Print(repoImportHelpText)
		return nil
	}
	opts, appErr := parseRepoImportFlags(args)
	if appErr != nil {
		return appErr
	}
	manifest, appErr := readRepoManifest(opts.file)
	if appErr != nil {
		return appErr
	}
	dest := filepath.Join(opts.into, manifestRepoDirName(manifest))
	if opts.isDry {
		printRepoImportDryRun(opts, manifest, dest)
		return nil
	}
	if appErr := cloneRepoImport(manifest.CloneURL, dest); appErr != nil {
		return appErr
	}
	restored, appErr := restoreRepoImportRemotes(dest, manifest.Remotes)
	if appErr != nil {
		return appErr
	}
	if appErr := registerRepoImport(dest, manifest); appErr != nil {
		return appErr
	}
	fmt.Printf(repoImportResultFmt, dest, restored)
	return nil
}

// parseRepoImportFlags parses the `repo import` flag set + required <file>.
func parseRepoImportFlags(args []string) (repoImportOpts, *appfault.AppError) {
	opts := repoImportOpts{}
	fs := flag.NewFlagSet(repoImportName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.into, repoImportFlagInto, "", "destination directory for the clone (required)")
	fs.BoolVar(&opts.isDry, repoImportFlagDry, false, "print what would happen, change nothing")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(repoImportUsageLine + "error: " + err.Error())
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Print(repoImportUsageLine)
		return opts, appfault.NewValidationError(repoImportErrFile)
	}
	opts.file = rest[0]
	if strings.TrimSpace(opts.into) == "" {
		fmt.Print(repoImportUsageLine)
		return opts, appfault.NewValidationError(repoImportErrInto)
	}
	return opts, nil
}

// readRepoManifest reads and validates a repo export manifest file.
func readRepoManifest(path string) (repoManifest, *appfault.AppError) {
	var manifest repoManifest
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, appfault.NewNotFoundError(fmt.Sprintf("repo import: could not read %s: %v", path, err))
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest, appfault.NewValidationError(fmt.Sprintf(repoImportErrParse, path, err))
	}
	if manifest.CloneURL == "" {
		manifest.CloneURL = manifest.Record.URL
	}
	if strings.TrimSpace(manifest.CloneURL) == "" {
		return manifest, appfault.NewValidationError(fmt.Sprintf(repoImportErrURL, path))
	}
	if manifest.Remotes == nil {
		manifest.Remotes = map[string]string{}
	}
	return manifest, nil
}

// manifestRepoDirName picks the clone directory name: URL base, else the
// record's repo name, else the slug.
func manifestRepoDirName(manifest repoManifest) string {
	base := strings.TrimSuffix(filepath.Base(strings.TrimSuffix(manifest.CloneURL, "/")), ".git")
	if base != "" && base != "/" && base != "." {
		return base
	}
	if manifest.Record.RepoName != "" {
		return manifest.Record.RepoName
	}
	return manifest.Record.Slug
}

// printRepoImportDryRun prints exactly what would happen (spec §5).
func printRepoImportDryRun(opts repoImportOpts, manifest repoManifest, dest string) {
	names := sortedKeys(manifest.Remotes)
	fmt.Printf(repoImportDryRunFmt, opts.file, manifest.CloneURL, dest, len(names), strings.Join(names, ", "), manifest.Record.Slug)
}

// cloneRepoImport clones the manifest's URL into dest.
func cloneRepoImport(cloneURL, dest string) *appfault.AppError {
	if _, err := os.Stat(dest); err == nil {
		return appfault.NewValidationError(fmt.Sprintf(repoImportErrExists, dest))
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("repo import: could not create parent dir for %s: %v", dest, err))
	}
	var stderr strings.Builder
	cmd := exec.Command("git", "clone", cloneURL, dest)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return apperror.NewExecutionError(fmt.Sprintf("repo import: git clone failed: %s", detail))
	}
	return nil
}

// restoreRepoImportRemotes ensures every manifest remote exists in the clone.
// Returns the number of remotes restored.
func restoreRepoImportRemotes(dest string, remotes map[string]string) (int, *appfault.AppError) {
	restored := 0
	for _, name := range sortedKeys(remotes) {
		url := remotes[name]
		if name == "origin" {
			if appErr := gitRemoteSetURL(dest, name, url); appErr != nil {
				return restored, appErr
			}
			restored++
			continue
		}
		if err := exec.Command("git", "-C", dest, "remote", "add", name, url).Run(); err != nil {
			if appErr := gitRemoteSetURL(dest, name, url); appErr != nil {
				return restored, appErr
			}
		}
		restored++
	}
	return restored, nil
}

// gitRemoteSetURL points an existing remote at a URL.
func gitRemoteSetURL(dest, name, url string) *appfault.AppError {
	if err := exec.Command("git", "-C", dest, "remote", "set-url", name, url).Run(); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("repo import: could not set remote %q to %s: %v", name, url, err))
	}
	return nil
}

// registerRepoImport upserts the cloned repo into the gitmap database.
func registerRepoImport(dest string, manifest repoManifest) *appfault.AppError {
	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "repo import: could not open the gitmap database")
	}
	defer db.Close()
	rec := manifest.Record.ScanRecord
	rec.AbsolutePath = dest
	if rec.Slug == "" {
		rec.Slug = manifestRepoDirName(manifest)
	}
	if rec.RepoName == "" {
		rec.RepoName = manifestRepoDirName(manifest)
	}
	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		return apperror.WrapSimple(err, "repo import: could not register the repo in gitmap.db")
	}
	return nil
}

// sortedKeys returns a manifest remote map's keys in sorted order.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

// repoImportHelpText is the full `--help` text for `repo import`.
const repoImportHelpText = `gitmap repo import — import a repo from an export manifest

USAGE
  gitmap repo import <file> --into <dir> [--dry-run]
  gitmap repoi <file> --into <dir> [same flags]   (alias; identical behavior)

FLAGS
  --into <dir>   Destination directory for the clone (required)
  --dry-run      Print what would happen, change nothing

EXAMPLE
  gitmap repo import ./story-telling-v1.gitmap.json --into ~/workspace/repos

NOTES
  - Clones the manifest's clone URL, restores every extra remote, and
    registers the repo in gitmap.db (the scan index).
  - The destination must not already exist.
  - Preview first: gitmap repo import <file> --into <dir> --dry-run
`

// RepoImport help-example data lives in cli/helpdoc/examples_275.go
// (single source of truth, program 275).
