// Package visibility — match.go: combines compiled patterns against an
// already-fetched list of repo names. Pure function, zero I/O — the
// provider CLI call lives in gitmap/cmd/visibilityownerlist.go.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §3.
package visibility

// MatchedRepo records which pattern first matched a repo. The "first
// matcher wins" rule keeps the matched-by column in the interactive
// table deterministic when several patterns overlap.
type MatchedRepo struct {
	RepoName       string
	MatchedPattern string // Pattern.Raw of the winning pattern
}

// MatchOwnerRepos walks `repos` in their input order, tests each name
// against every pattern in order, and returns one MatchedRepo per
// matched name. Dedupes by RepoName keeping the first pattern that
// matched. Returns an empty slice when nothing matches — never nil.
func MatchOwnerRepos(repos []string, patterns []Pattern) []MatchedRepo {
	out := make([]MatchedRepo, 0, len(repos))
	seen := make(map[string]bool, len(repos))
	for _, name := range repos {
		if seen[name] {
			continue
		}
		winner := firstMatchingPattern(name, patterns)
		if len(winner) == 0 {
			continue
		}
		seen[name] = true
		out = append(out, MatchedRepo{RepoName: name, MatchedPattern: winner})
	}

	return out
}

// firstMatchingPattern returns the Raw form of the first pattern that
// accepts `name`, or "" when none match.
func firstMatchingPattern(name string, patterns []Pattern) string {
	for _, p := range patterns {
		if p.Matches(name) {
			return p.Raw
		}
	}

	return ""
}
