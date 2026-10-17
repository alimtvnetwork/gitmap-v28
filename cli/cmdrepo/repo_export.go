package cmdrepo

// repo_export.go — `gitmap repo export <name> --out <file>` / `gitmap repox`:
// portable manifest JSON (clone URL, remotes, branches), reusing the `scan
// export` record shape (program 275, proposal P-05). One concern, ≤300 lines.

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// repo export messages.
const (
	repoExportName       = "repo export"
	repoExportFlagOut    = "out"
	repoExportUsageLine  = "Usage: gitmap repo export <name> --out <file>\n       gitmap repox <name> --out <file>\n"
	repoExportErrName    = "missing required argument: <name> (repo slug, name, or path, e.g. `gitmap repo export story-telling-v1 --out ./story-telling-v1.gitmap.json`)"
	repoExportErrOut     = "missing required flag: --out <file> (where to write the manifest JSON)"
	repoExportErrFound   = "no repository matched %q"
	repoExportErrMulti   = "multiple repositories matched %q; be more specific:\n%s"
	repoExportResultFmt  = "exported 1 repo · %d remotes · %d branches · written %s (%s)\n"
	repoExportFormatName = "gitmap-repo-manifest"
	repoExportFormatVer  = 1
)

// repoManifestRecord is the scan-export record shape: model.ScanRecord plus
// the `url` key `gitmap clone-from` consumes.
type repoManifestRecord struct {
	model.ScanRecord
	URL string `json:"url"`
}

// repoManifest is the portable single-repo manifest: written by repo export,
// consumed by repo import.
type repoManifest struct {
	Format     string             `json:"format"`
	Version    int                `json:"version"`
	ExportedAt string             `json:"exportedAt"`
	CloneURL   string             `json:"cloneUrl"`
	Remotes    map[string]string  `json:"remotes"`
	Branches   []string           `json:"branches"`
	Record     repoManifestRecord `json:"record"`
}

// repoExportOpts carries the parsed flags for `repo export`.
type repoExportOpts struct {
	name string
	out  string
}

// RunRepoExport implements `gitmap repo export` / `gitmap repox` for rootdata wiring.
func RunRepoExport(args []string) error {
	if appErr := runRepoExport(args); appErr != nil {
		return appErr
	}
	return nil
}

// runRepoExport executes the export flow: parse, find, gather, write.
func runRepoExport(args []string) *appfault.AppError {
	if isRepoHelpArg(args) {
		fmt.Print(repoExportHelpText)
		return nil
	}
	opts, appErr := parseRepoExportFlags(args)
	if appErr != nil {
		return appErr
	}
	db, appErr := openRepoExportDB()
	if appErr != nil {
		return appErr
	}
	defer db.Close()
	rec, appErr := findRepoByName(db, opts.name)
	if appErr != nil {
		return appErr
	}
	remotes, appErr := gitRepoRemotes(rec.AbsolutePath)
	if appErr != nil {
		return appErr
	}
	branches, appErr := gitRepoBranches(rec.AbsolutePath)
	if appErr != nil {
		return appErr
	}
	manifest := buildRepoManifest(rec, remotes, branches)
	size, appErr := writeRepoManifest(opts.out, manifest)
	if appErr != nil {
		return appErr
	}
	fmt.Printf(repoExportResultFmt, len(remotes), len(branches), opts.out, formatByteSize(size))
	return nil
}

// parseRepoExportFlags parses the `repo export` flag set + required <name>.
func parseRepoExportFlags(args []string) (repoExportOpts, *appfault.AppError) {
	opts := repoExportOpts{}
	fs := flag.NewFlagSet(repoExportName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.out, repoExportFlagOut, "", "output manifest file (required)")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(repoExportUsageLine + "error: " + err.Error())
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Print(repoExportUsageLine)
		return opts, appfault.NewValidationError(repoExportErrName)
	}
	opts.name = rest[0]
	if strings.TrimSpace(opts.out) == "" {
		fmt.Print(repoExportUsageLine)
		return opts, appfault.NewValidationError(repoExportErrOut)
	}
	return opts, nil
}

// openRepoExportDB opens the default gitmap database.
func openRepoExportDB() (*store.DB, *appfault.AppError) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "repo export: could not open the gitmap database")
	}
	return db, nil
}

// findRepoByName resolves <name> to one scan record: exact slug, exact
// repo name, directory base name (case-insensitive), then a single
// unambiguous substring match.
func findRepoByName(db *store.DB, name string) (model.ScanRecord, *appfault.AppError) {
	var zero model.ScanRecord
	records, err := db.ListRepos()
	if err != nil {
		return zero, apperror.WrapSimple(err, "repo export: could not list repositories")
	}
	want := strings.ToLower(strings.TrimSpace(name))
	var partial []model.ScanRecord
	for _, r := range records {
		if strings.EqualFold(r.Slug, want) || strings.EqualFold(r.RepoName, want) {
			return r, nil
		}
		if strings.EqualFold(filepath.Base(r.AbsolutePath), want) {
			return r, nil
		}
		if strings.Contains(strings.ToLower(r.Slug), want) || strings.Contains(strings.ToLower(r.RepoName), want) {
			partial = append(partial, r)
		}
	}
	if len(partial) == 1 {
		return partial[0], nil
	}
	if len(partial) > 1 {
		var lines []string
		for _, r := range partial {
			lines = append(lines, "  "+r.Slug+" ("+r.AbsolutePath+")")
		}
		return zero, appfault.NewValidationError(fmt.Sprintf(repoExportErrMulti, name, strings.Join(lines, "\n")))
	}
	return zero, appfault.NewNotFoundError(fmt.Sprintf(repoExportErrFound, name))
}

// gitRepoRemotes returns name→fetch-URL for every remote of the repo dir.
func gitRepoRemotes(dir string) (map[string]string, *appfault.AppError) {
	out, err := exec.Command("git", "-C", dir, "remote", "-v").Output()
	if err != nil {
		return nil, apperror.NewExecutionError(fmt.Sprintf("repo export: cannot read git remotes in %s: %v", dir, err))
	}
	remotes := map[string]string{}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[2] != "(fetch)" {
			continue
		}
		if _, seen := remotes[fields[0]]; !seen {
			remotes[fields[0]] = fields[1]
		}
	}
	return remotes, nil
}

// gitRepoBranches lists local + remote branch short names for the repo dir.
func gitRepoBranches(dir string) ([]string, *appfault.AppError) {
	out, err := exec.Command("git", "-C", dir, "for-each-ref", "--format=%(refname:short)", "refs/heads", "refs/remotes").Output()
	if err != nil {
		return nil, apperror.NewExecutionError(fmt.Sprintf("repo export: cannot read git branches in %s: %v", dir, err))
	}
	var branches []string
	for _, line := range strings.Split(string(out), "\n") {
		if b := strings.TrimSpace(line); b != "" {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

// preferredManifestCloneURL picks the clone URL: https first, then ssh,
// then discovered (same priority as scan export).
func preferredManifestCloneURL(r model.ScanRecord) string {
	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}
	if r.SSHUrl != "" {
		return r.SSHUrl
	}
	return r.DiscoveredURL
}

// buildRepoManifest assembles the manifest from the record + git state.
func buildRepoManifest(rec model.ScanRecord, remotes map[string]string, branches []string) repoManifest {
	cloneURL := preferredManifestCloneURL(rec)
	if remotes == nil {
		remotes = map[string]string{}
	}
	if branches == nil {
		branches = []string{}
	}
	return repoManifest{
		Format:     repoExportFormatName,
		Version:    repoExportFormatVer,
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		CloneURL:   cloneURL,
		Remotes:    remotes,
		Branches:   branches,
		Record:     repoManifestRecord{ScanRecord: rec, URL: cloneURL},
	}
}

// writeRepoManifest writes the indented manifest JSON and returns its size.
func writeRepoManifest(path string, manifest repoManifest) (int64, *appfault.AppError) {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return 0, apperror.NewExecutionError(fmt.Sprintf("repo export: could not create output dir %s: %v", dir, err))
		}
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return 0, apperror.NewExecutionError(fmt.Sprintf("repo export: could not encode manifest: %v", err))
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return 0, apperror.NewExecutionError(fmt.Sprintf("repo export: could not write %s: %v", path, err))
	}
	return int64(len(data)), nil
}

// formatByteSize renders a byte count as "4.2 KB".
func formatByteSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.1f KB", float64(n)/1024)
}

// isRepoHelpArg reports whether args ask for help (cmdrepo local helper).
func isRepoHelpArg(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

// repoExportHelpText is the full `--help` text for `repo export`.
const repoExportHelpText = `gitmap repo export — export a repo's portable manifest

USAGE
  gitmap repo export <name> --out <file>
  gitmap repox <name> --out <file>   (alias; identical behavior)

FLAGS
  --out <file>   Manifest JSON to write (required)

EXAMPLE
  gitmap repo export story-telling-v1 --out ./story-telling-v1.gitmap.json

NOTES
  - The manifest carries the clone URL, all remotes, branch names, and the
    scan-export record shape — everything 'gitmap repo import' needs.
  - <name> matches slug, repo name, or directory base name. Read-only.
`

// RepoExport help-example data lives in cli/helpdoc/examples_275.go
// (single source of truth, program 275).
