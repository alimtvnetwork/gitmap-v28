// Package cloner — read-only audit for clone manifests.
//
// PlanCloneAudit parses a structured source file (CSV/JSON/text) and, for
// every record, computes the `git clone`/`git pull` command that
// `gitmap clone` would run — without actually running it. The result is a
// diff-style report:
//
//   - clone   path  (url, branch, strategy)   target missing
//     ~ pull    path  (url)                     existing git repo
//     = skip    path  (url)                     cache fingerprint matches
//     ! invalid path                            no clone URL on the record
//     ? conflict path                           target exists but is not a git repo
//
// Audit never touches the network and never writes anything outside stdout.
// It honors the same branch-selection strategy as the live clone path
// (pickCloneStrategy) so what you see is exactly what `gitmap clone` would
// execute.
package cloner

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// AuditActionType labels the planned action for a single record.
type AuditActionType string

const (
	AuditActionClone    AuditActionType = "clone"    // target missing — fresh clone
	AuditActionPull     AuditActionType = "pull"     // target is a git repo — would safe-pull
	AuditActionCached   AuditActionType = "cached"   // cache fingerprint matches dest
	AuditActionInvalid  AuditActionType = "invalid"  // record has no clone URL
	AuditActionConflict AuditActionType = "conflict" // target exists but is not a git repo
)

// AuditEntry is one planned-vs-actual row in the audit report.
type AuditEntry struct {
	Action       AuditActionType
	RelativePath string
	URL          string
	Branch       string
	UseBranch    bool
	Reason       string // strategy.reason or conflict explanation
	Command      string // exact git command line that would be run, "" for invalid/cached
}

// CloneAuditReport is the full per-record plan plus aggregate counters.
type CloneAuditReport struct {
	Source  string
	Target  string
	Entries []AuditEntry
	Counts  map[AuditActionType]int
}

// PlanCloneAudit loads the source manifest and produces a non-executing
// plan for every record. Returns an error only when the source file
// cannot be loaded; per-record problems are encoded as audit entries.
func PlanCloneAudit(sourcePath, targetDir string) (*CloneAuditReport, error) {
	records, err := loadRecords(sourcePath)
	if err != nil {
		return nil, err
	}

	cache := LoadCloneCache(targetDir)
	report := &CloneAuditReport{
		Source: sourcePath,
		Target: targetDir,
		Counts: make(map[AuditActionType]int),
	}
	for _, rec := range records {
		entry := planOne(rec, targetDir, cache)
		report.Entries = append(report.Entries, entry)
		report.Counts[entry.Action]++
	}

	return report, nil
}

// planOne computes the audit action for a single record. Pure function:
// no git invocations, no writes, no network.
func planOne(rec model.ScanRecord, targetDir string, cache *CloneCache) AuditEntry {
	dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
	url := pickURL(rec)

	if len(url) == 0 {
		return AuditEntry{
			Action:       AuditActionInvalid,
			RelativePath: rec.RelativePath,
			Reason:       "record has no HTTPSUrl or SSHUrl",
		}
	}

	strategy := pickCloneStrategy(rec)
	command := buildAuditCommand(url, dest, strategy)
	base := AuditEntry{
		RelativePath: rec.RelativePath,
		URL:          url,
		Branch:       strategy.branch,
		UseBranch:    strategy.useBranch,
		Reason:       strategy.reason,
		Command:      command,
	}

	classifyParams := ClassifyDestParams{
		Base:  base,
		Dest:  dest,
		Rec:   rec,
		Cache: cache,
	}

	return classifyDest(classifyParams)
}

// ClassifyDestParams encapsulates parameters for classifying a destination.
type ClassifyDestParams struct {
	Base  AuditEntry
	Dest  string
	Rec   model.ScanRecord
	Cache *CloneCache
}

// classifyDest decides clone/pull/cached/conflict by inspecting the
// destination on disk. Stat-only — never opens the repo.
func classifyDest(params ClassifyDestParams) AuditEntry {
	base := params.Base
	if !pathExists(params.Dest) {
		base.Action = AuditActionClone
		return base
	}
	if !IsGitRepo(params.Dest) {
		base.Action = AuditActionConflict
		base.Reason = "target path exists but is not a git repository"
		return base
	}
	if params.Cache != nil && params.Cache.IsUpToDate(params.Rec, params.Dest) {
		base.Action = AuditActionCached
		base.Reason = "clone-cache fingerprint matches local HEAD"
		return base
	}
	base.Action = AuditActionPull

	return base
}

// buildAuditCommand renders the exact `git clone` command that would be
// executed for a given (url, dest, strategy). Mirrors runClone() in
// cloner.go so the audit cannot drift from the live behavior.
func buildAuditCommand(url, dest string, strategy cloneStrategy) string {
	if strategy.useBranch {
		return fmt.Sprintf("%s %s %s %s %s %s",
			constants.GitBin, constants.GitClone,
			constants.GitBranchFlag, strategy.branch, url, dest)
	}

	return fmt.Sprintf("%s %s %s %s",
		constants.GitBin, constants.GitClone, url, dest)
}

// Print renders the report to w in the diff-style layout described in the
// package doc. Returns the number of bytes written or an error if the
// underlying writer fails — never silent.
func (r *CloneAuditReport) Print(w io.Writer) error {
	if _, err := fmt.Fprintf(w, constants.MsgCloneAuditHeader, r.Source, r.Target, len(r.Entries)); err != nil {
		return err
	}
	for _, e := range r.Entries {
		if err := printEntry(w, e); err != nil {
			return err
		}
	}

	return printSummary(w, r.Counts)
}

// printEntry prints a single audit row, prefixed by its diff marker.
func printEntry(w io.Writer, e AuditEntry) error {
	marker := actionMarker(e.Action)
	_, err := fmt.Fprintf(w, constants.MsgCloneAuditRow, marker, e.Action, e.RelativePath, e.URL, e.Reason)
	if err != nil {
		return err
	}
	if len(e.Command) == 0 {
		return nil
	}
	_, err = fmt.Fprintf(w, constants.MsgCloneAuditCmd, e.Command)

	return err
}

// printSummary prints the trailing aggregate counts.
func printSummary(w io.Writer, counts map[AuditActionType]int) error {
	_, err := fmt.Fprintf(w, constants.MsgCloneAuditSummary,
		counts[AuditActionClone], counts[AuditActionPull],
		counts[AuditActionCached], counts[AuditActionConflict],
		counts[AuditActionInvalid])

	return err
}

// actionMarker maps each action to its diff-style single-character marker.
func actionMarker(a AuditActionType) string {
	switch a {
	case AuditActionClone:
		return "+"
	case AuditActionPull:
		return "~"
	case AuditActionCached:
		return "="
	case AuditActionInvalid:
		return "!"
	case AuditActionConflict:
		return "?"
	default:
		return " "
	}
}
