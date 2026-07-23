// Package vscodepm — mergemode.go: tag-merge strategy enum + dispatch.
//
// Until v4.36.0 the only merge strategy was UNION (existing ∪ incoming).
// v4.37.0 adds REPLACE and INTERSECTION via a `--mode` flag on
// `gitmap vscode-pm-sync`. To keep the rest of the package tiny and
// the merge call sites grep-able, all three strategies live in one
// place behind a single dispatcher.
//
// Modes:
//
//   - UNION         existing ∪ incoming, dedup'd, existing-order first.
//     The original v4.36.0 behavior. Default.
//
//   - REPLACE       incoming verbatim. User-added tags are dropped.
//     The detector (DetectTagsCustom) already pre-pends
//     the "gitmap" brand tag, so REPLACE naturally keeps
//     the brand without any special-case code path.
//
//   - INTERSECTION  (existing ∩ incoming) ∪ {AutoTagGitmap}.
//     Strict set-AND, but the brand tag is ALWAYS pinned
//     so re-syncs never silently strip the gitmap brand
//     from an entry whose detector run no longer fires
//     (e.g. .git/ folder moved). Matches the brand-
//     everywhere rule introduced in v4.36.0.
//
// Magic strings: the literal mode names live in `constants_cli.go`
// alongside the flag itself, never duplicated here.
package vscodepm

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// MergeMode picks the tag-merge strategy used by SyncMode / SyncAtMode.
// uint8 per the project-wide rule "smallest viable int type for ≤255 members".
type MergeMode uint8

const (
	// MergeModeUnion is the v4.36.0 default — existing ∪ incoming, dedup'd.
	MergeModeUnion MergeMode = iota
	// MergeModeReplace overwrites with incoming verbatim. Detector keeps brand.
	MergeModeReplace
	// MergeModeIntersection keeps only tags present in BOTH sets, plus brand pin.
	MergeModeIntersection
)

// String returns the canonical CLI literal for m. Matches the values
// accepted by --mode and the values listed in the helptext.
func (m MergeMode) String() string {
	switch m {
	case MergeModeReplace:
		return constants.VSCodePMSyncModeReplace
	case MergeModeIntersection:
		return constants.VSCodePMSyncModeIntersection
	default:
		return constants.VSCodePMSyncModeUnion
	}
}

// ParseMergeMode converts the CLI string into a MergeMode. Empty
// string is treated as the default (UNION) so callers don't have to
// special-case "flag absent". Unknown values return an error so we
// fail loud per the zero-swallow rule.
func ParseMergeMode(raw string) (MergeMode, error) {
	switch raw {
	case "", constants.VSCodePMSyncModeUnion:
		return MergeModeUnion, nil
	case constants.VSCodePMSyncModeReplace:
		return MergeModeReplace, nil
	case constants.VSCodePMSyncModeIntersection:
		return MergeModeIntersection, nil
	default:
		return MergeModeUnion, fmt.Errorf(
			constants.ErrVSCodePMSyncBadMode, raw,
			constants.VSCodePMSyncModeUnion,
			constants.VSCodePMSyncModeReplace,
			constants.VSCodePMSyncModeIntersection)
	}
}

// mergeTags dispatches to the right strategy. Lives next to the enum
// so adding a 4th mode is a one-file change.
func mergeTags(mode MergeMode, existing, incoming []string) []string {
	switch mode {
	case MergeModeReplace:
		return replaceTags(incoming)
	case MergeModeIntersection:
		return intersectTagsWithBrand(existing, incoming)
	default:
		return unionTags(existing, incoming)
	}
}

// replaceTags returns a dedup'd copy of incoming. Order preserved.
// Nil input -> empty slice (never nil) so the JSON encoder emits []
// instead of null, matching the rest of the package.
func replaceTags(incoming []string) []string {
	seen := make(map[string]struct{}, len(incoming))
	out := make([]string, 0, len(incoming))
	for _, t := range incoming {
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}

	return out
}

// intersectTagsWithBrand returns (existing ∩ incoming), then guarantees
// the gitmap brand tag is present. Existing-order is preserved so the
// VS Code UI doesn't reshuffle on every sync.
func intersectTagsWithBrand(existing, incoming []string) []string {
	incomingSet := make(map[string]struct{}, len(incoming))
	for _, t := range incoming {
		incomingSet[t] = struct{}{}
	}

	seen := make(map[string]struct{}, len(existing)+1)
	out := make([]string, 0, len(existing)+1)
	for _, t := range existing {
		if _, ok := incomingSet[t]; !ok {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}

	if _, has := seen[constants.AutoTagGitmap]; !has {
		out = append(out, constants.AutoTagGitmap)
	}

	return out
}
