// Package cmd — visibilityexceptlatest.go: splits a matched repo set
// so the newest `-vN` sibling per base group is held out for the
// INVERTED visibility flip. Repos that don't carry a `-vN` suffix are
// left in the main (target-visibility) bucket untouched.
//
// Behavior (v6.65.0+):
//
//	--except-latest no longer just "preserves" the latest version —
//	it flips it to the opposite of the requested target. Example:
//	  make-all-public  --except-latest → all → public, latest → PRIVATE
//	  make-all-private --except-latest → all → private, latest → PUBLIC
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §except-latest.
package cmd

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/visibility"
)

// versionSuffixRE captures the trailing `-v<digits>` segment.
var versionSuffixRE = regexp.MustCompile(`(?i)^(.*)-v(\d+)$`)

// splitExceptLatest separates `in` into (rest, latest). `latest` holds
// the highest -vN entry per base group (groups with a single versioned
// entry are NOT split out — there is nothing meaningful to invert).
// Unversioned repos always land in `rest`. Each split is logged via w.
func splitExceptLatest(in []visibility.MatchedRepo, w io.Writer, invertedTarget string) ([]visibility.MatchedRepo, []visibility.MatchedRepo) {
	type peak struct {
		idx int
		ver int
	}
	peaks := map[string]peak{}
	counts := map[string]int{}
	for i, m := range in {
		base, ver, ok := parseVersionedName(m.RepoName)
		if !ok {
			continue
		}
		counts[base]++
		if cur, seen := peaks[base]; !seen || ver > cur.ver {
			peaks[base] = peak{idx: i, ver: ver}
		}
	}

	pick := map[int]int{}
	for base, p := range peaks {
		if counts[base] < 2 {
			continue
		}
		pick[p.idx] = p.ver
	}

	rest := make([]visibility.MatchedRepo, 0, len(in))
	latest := make([]visibility.MatchedRepo, 0, len(pick))
	for i, m := range in {
		if v, ok := pick[i]; ok {
			fmt.Fprintf(w, constants.MsgBulkExceptInvertFmt, m.RepoName, v, invertedTarget)
			latest = append(latest, m)

			continue
		}
		rest = append(rest, m)
	}

	return rest, latest
}

// parseVersionedName returns (base, version, true) when `name` ends
// in `-v<digits>`, or ("", 0, false) otherwise.
func parseVersionedName(name string) (string, int, bool) {
	m := versionSuffixRE.FindStringSubmatch(name)
	if m == nil {
		return "", 0, false
	}
	v, err := strconv.Atoi(m[2])
	if err != nil {
		return "", 0, false
	}

	return m[1], v, true
}

// invertVisibility returns the opposite of constants.VisibilityPublic /
// VisibilityPrivate. Any other input passes through unchanged so the
// caller's downstream provider call still surfaces a clean error.
func invertVisibility(target string) string {
	switch target {
	case constants.VisibilityPublic:
		return constants.VisibilityPrivate
	case constants.VisibilityPrivate:
		return constants.VisibilityPublic
	}

	return target
}
