package mapper

import "github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

// BuildOptions bundles every configurable knob BuildRecords* exposes
// so adding a new option later doesn't grow the helper signatures
// further. The two thin wrappers in mapper.go preserve the legacy
// positional signatures for callers that don't need the new fields.
type BuildOptions struct {
	// Mode is "https" or "ssh" — selects which clone URL is recorded.
	Mode string
	// DefaultNote is written to ScanRecord.Notes when the repo has a
	// remote URL. Empty by default.
	DefaultNote string
	// RelRoot, when non-empty, rewrites every RelativePath against
	// this absolute, cleaned path. Repos outside relRoot fall back to
	// the scanner-computed RelativePath with a stderr warning. See
	// relativePathFor for the precise contract.
	RelRoot string
	// DefaultBranch is the fallback branch name passed to
	// gitutil.DetectBranchWithDefault. Empty string means "use the
	// built-in constants.DefaultBranch" (preserves legacy behavior).
	// CLI surface: `gitmap scan --default-branch <name>`.
	DefaultBranch string
}

// resolveDefaultBranch picks the fallback branch name passed to
// gitutil.DetectBranchWithDefault. An empty override (the zero value
// of BuildOptions.DefaultBranch) resolves to constants.DefaultBranch
// so legacy callers see identical behavior to the pre-flag impl.
func resolveDefaultBranch(override string) string {
	if override == "" {
		return constants.DefaultBranch
	}

	return override
}
