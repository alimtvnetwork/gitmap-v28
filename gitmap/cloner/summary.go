// Package cloner — summary.go
//
// Result-shape helpers split out of cloner.go to keep that file focused
// on entry points + parsing. Both runners (sequential in runners.go and
// parallel in concurrent.go) share these helpers, so a single source of
// truth for "what counts as success / skipped" prevents the two paths
// from drifting.
package cloner

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// recordTag returns a short, log-friendly identifier for a record using
// the most specific field available. Used in error messages so users can
// locate the failing row in their clone manifest at a glance.
func recordTag(rec model.ScanRecord) string {
	switch {
	case len(rec.RepoName) > 0 && len(rec.RelativePath) > 0:
		return fmt.Sprintf("%s (%s)", rec.RepoName, rec.RelativePath)
	case len(rec.RepoName) > 0:
		return rec.RepoName
	case len(rec.RelativePath) > 0:
		return rec.RelativePath
	case len(rec.HTTPSUrl) > 0:
		return rec.HTTPSUrl
	case len(rec.SSHUrl) > 0:
		return rec.SSHUrl
	default:
		return "<unnamed record>"
	}
}

// pickURL selects the best available URL from a record, honoring the
// repo's identified Transport: SSH-origin repos are cloned/pulled over
// SSH (never silently fall back to HTTPS — that triggers browser auth
// on private remotes). Falls through to the other transport only when
// the preferred URL is empty.
func pickURL(rec model.ScanRecord) string {
	if rec.Transport == "ssh" && len(rec.SSHUrl) > 0 {
		return rec.SSHUrl
	}
	if len(rec.HTTPSUrl) > 0 {
		return rec.HTTPSUrl
	}

	return rec.SSHUrl
}

// updateSummary increments counters and collects results.
func updateSummary(s model.CloneSummary, r model.CloneResult) model.CloneSummary {
	if r.Success {
		s.Succeeded++
		s.Cloned = append(s.Cloned, r)

		return s
	}
	s.Failed++
	s.Errors = append(s.Errors, r)

	return s
}

// updateSummarySkipped records a cache-skipped repo: it counts toward
// Succeeded (the desired state is achieved) and is also tracked in
// Cloned + Skipped so downstream consumers (GitHub Desktop registration,
// reporting) treat it the same as a fresh clone.
func updateSummarySkipped(s model.CloneSummary, r model.CloneResult) model.CloneSummary {
	s.Succeeded++
	s.Cloned = append(s.Cloned, r)
	s.Skipped = append(s.Skipped, r)

	return s
}
