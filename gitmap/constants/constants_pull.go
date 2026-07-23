package constants

// Pull-specific flags + messages added in Phase 2.5 (v3.10.0).
//
// `--parallel <N>` runs the pull batch through a worker pool of width N.
// `--only-available` intersects the target set with `gitmap find-next`
// before pulling, so we skip every repo whose latest VersionProbe row
// reports no new tag — turns "pull all" into "pull what's actually new".

// CLI flag descriptions (used by flag.NewFlagSet).
const (
	FlagDescPullParallel      = "Run up to N pulls concurrently (default 1; serial)"
	FlagDescPullOnlyAvailable = "Skip repos whose latest probe reports no new tag (run `gitmap probe --all` first)"
)

// User-facing pull messages (Phase 2.5).
const (
	MsgPullNoAvailable     = "No repos with available updates. Run `gitmap probe --all` first.\n"
	WarnPullFilterFallback = "  ⚠ --only-available filter unavailable (probe DB unreadable); falling back to full target set"
)
