package committransfer

import (
	"regexp"
	"strings"
	"time"
)

// DefaultDropPatterns mirrors spec §6.1.
var DefaultDropPatterns = []string{
	`^Merge branch`,
	`^Merge pull request`,
	`^Revert "`,
	`^fixup!`,
	`^squash!`,
	`^WIP$`,
}

// CleanResult is the output of the message pipeline for one commit.
type CleanResult struct {
	Final   string // empty when the commit should be skipped
	Skipped string // non-empty reason when Final == ""
}

// CleanMessage runs the full §6 pipeline (drop → strip → conventional →
// provenance → empty-guard) and returns the final commit message or a
// skip reason. footerArgs (shortSHA + when) are only used when
// MessagePolicy.Provenance is true.
func CleanMessage(subject, body string, p MessagePolicy, shortSHA string, when time.Time) CleanResult {
	if reason := matchDrop(subject, p.DropPatterns); reason != "" {
		return CleanResult{Skipped: "drop-pattern " + reason}
	}
	cleanedSubject := applyStrip(subject, p.StripPatterns)
	cleanedSubject = strings.TrimSpace(cleanedSubject)
	if cleanedSubject == "" {
		return CleanResult{Skipped: "cleaned-empty"}
	}
	if p.Conventional {
		cleanedSubject = normalizeConventional(cleanedSubject)
	}
	final := assembleMessage(cleanedSubject, body)
	if p.Provenance {
		final = appendProvenanceFooter(final, p, shortSHA, when)
	}

	return CleanResult{Final: final}
}

// matchDrop returns the first matching pattern (or ""). Bad regexes are
// silently ignored so a malformed user pattern can never crash a replay.
func matchDrop(subject string, patterns []string) string {
	for _, pat := range patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			continue
		}
		if re.MatchString(subject) {
			return pat
		}
	}

	return ""
}

// applyStrip runs every pattern as a regex replace-with-empty on subject.
func applyStrip(subject string, patterns []string) string {
	out := subject
	for _, pat := range patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			continue
		}
		out = re.ReplaceAllString(out, "")
	}

	return out
}

// conventionalPrefix matches `type` or `type(scope):` at the start.
var conventionalPrefix = regexp.MustCompile(
	`^(feat|fix|chore|docs|refactor|test|build|ci|perf|style|revert)(\([^)]+\))?:\s`)

// normalizeConventional applies the spec §6.3 heuristic. It does NOT
// inspect the diff (that's a future enhancement) — it only acts on the
// subject's leading verb.
func normalizeConventional(subject string) string {
	if conventionalPrefix.MatchString(subject) {
		return subject
	}
	lower := strings.ToLower(subject)
	switch {
	case startsWithAny(lower, "fix", "bugfix", "hotfix"):
		return "fix: " + stripLeadingVerb(subject)
	case startsWithAny(lower, "add", "introduce", "implement"):
		return "feat: " + stripLeadingVerb(subject)
	}

	return "chore: " + subject
}

// startsWithAny is a word-boundary "starts-with" helper for the
// heuristic prefix list.
func startsWithAny(lowerSubject string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(lowerSubject, p+" ") || lowerSubject == p {
			return true
		}
	}

	return false
}

// stripLeadingVerb drops the first word so "Add login form" → "login form"
// before "feat: " is glued on.
func stripLeadingVerb(subject string) string {
	parts := strings.SplitN(subject, " ", 2)
	if len(parts) < 2 {
		return subject
	}

	return parts[1]
}

// assembleMessage joins subject and body with a blank line between them
// (or just the subject when body is empty).
func assembleMessage(subject, body string) string {
	if body == "" {
		return subject
	}

	return subject + "\n\n" + body
}

// appendProvenanceFooter appends the §6.4 footer block.
func appendProvenanceFooter(msg string, p MessagePolicy, shortSHA string, when time.Time) string {
	footer := "gitmap-replay: from " + p.SourceDisplayName + " " + shortSHA +
		"\ngitmap-replay-cmd: " + p.CommandName +
		"\ngitmap-replay-at: " + when.Format(time.RFC3339)

	return msg + "\n\n" + footer
}

// AlreadyReplayed checks whether the recent target log contains a
// provenance footer that points at (sourceDisplay, shortSHA). Used by
// the planner when --force-replay is not set.
//
// Deprecated: this rebuilds the set on every call. Prefer BuildReplayedSet
// once per plan, then SetHasReplayed for O(1) lookups.
func AlreadyReplayed(recentLog, sourceDisplay, shortSHA string) bool {
	set := BuildReplayedSet(recentLog)

	return SetHasReplayed(set, sourceDisplay, shortSHA)
}

// BuildReplayedSet parses the concatenated target-log output (from
// recentLogSubjectsAndBodies) and returns a hash-set of all provenance
// (sourceDisplay, shortSHA) pairs.  Building the set once and re-using
// it avoids the O(N×M) substring scan when the target history is large
// (spec 114 Gap A).
func BuildReplayedSet(recentLog string) map[string]struct{} {
	set := make(map[string]struct{})
	prefix := "gitmap-replay: from "
	for _, line := range strings.Split(recentLog, "\n") {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		rest := strings.TrimPrefix(line, prefix)
		// rest = "<sourceDisplay> <shortSHA>"
		if idx := strings.LastIndex(rest, " "); idx > 0 {
			set[rest] = struct{}{}
		}
	}

	return set
}

// SetHasReplayed reports whether the pre-built set contains a provenance
// footer for (sourceDisplay, shortSHA).
func SetHasReplayed(set map[string]struct{}, sourceDisplay, shortSHA string) bool {
	_, ok := set[sourceDisplay+" "+shortSHA]

	return ok
}
