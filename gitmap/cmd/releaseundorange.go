// Package cmd — `gitmap release-undo --range` extension.
//
// Roll back several contiguous release tags at once. The range is
// inclusive on both ends and must be in `vX.Y.Z..vX.Y.Z` form.
// Each version is processed sequentially via the existing single-tag
// release-undo pipeline, so failures stop the run and leave earlier
// successes intact (idempotent — safe to re-run).
package cmd

import (
	"fmt"
	"io"
	"strings"
)

// ReleaseUndoRangeOptions configures a multi-tag undo.
type ReleaseUndoRangeOptions struct {
	Range       string // e.g. "v6.60.0..v6.65.0"
	KeepRemote  bool
	KeepSidecar bool
	Yes         bool
	DryRun      bool
	Stdout      io.Writer
	UndoOne     func(version string) error // injected for testability
}

// RunReleaseUndoRange parses the range and undoes each version in order.
func RunReleaseUndoRange(opts ReleaseUndoRangeOptions) int {
	versions, err := expandReleaseRange(opts.Range)
	if err != nil {
		fmt.Fprintf(opts.Stdout, "release-undo: %v\n", err)
		return 2
	}
	fmt.Fprintf(opts.Stdout, "release-undo: %d versions in range %s\n", len(versions), opts.Range)
	for _, v := range versions {
		if opts.DryRun {
			fmt.Fprintf(opts.Stdout, "[dry-run] would undo %s\n", v)
			continue
		}
		if err := opts.UndoOne(v); err != nil {
			fmt.Fprintf(opts.Stdout, "release-undo: %s failed: %v\n", v, err)
			return 1
		}
		fmt.Fprintf(opts.Stdout, "release-undo: %s done\n", v)
	}
	return 0
}

// expandReleaseRange turns "v6.60.0..v6.65.0" into the inclusive list
// of patch-level tags between the two endpoints. Endpoints must share
// the same major+minor; cross-minor ranges are rejected to avoid
// accidentally deleting hundreds of tags.
func expandReleaseRange(spec string) ([]string, error) {
	lo, hi, ok := strings.Cut(spec, "..")
	if !ok {
		return nil, fmt.Errorf("range must be of form vX.Y.Z..vX.Y.Z, got %q", spec)
	}
	a, b := splitSemver(strings.TrimPrefix(lo, "v")), splitSemver(strings.TrimPrefix(hi, "v"))
	if a[0] != b[0] || a[1] != b[1] {
		return nil, fmt.Errorf("range endpoints must share major.minor (%s vs %s)", lo, hi)
	}
	if a[2] > b[2] {
		return nil, fmt.Errorf("range endpoints reversed (%s > %s)", lo, hi)
	}
	out := make([]string, 0, b[2]-a[2]+1)
	for p := a[2]; p <= b[2]; p++ {
		out = append(out, fmt.Sprintf("v%d.%d.%d", a[0], a[1], p))
	}
	return out, nil
}
