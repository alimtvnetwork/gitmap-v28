package constants

// Constants for parallel clone execution shared across the
// clone-family commands (`gitmap clone`, `clone-next`, `clone-now` /
// `relclone`, and `clone-from`). All of them surface the same
// `--max-concurrency N` flag with identical semantics:
//
//   - N == 0 → auto: resolve to runtime.NumCPU() at start time
//     (the default when the flag is omitted, since v3.101+).
//   - N == 1 → sequential, byte-for-byte compatible with the legacy
//     pre-v3.101 behavior; useful when scripts grep stderr for
//     stable per-row ordering.
//   - N >  1 → bounded worker pool. The on-disk hierarchy is
//     unaffected because every worker still uses the manifest's
//     RelativePath / Dest verbatim — only progress-line ORDER
//     becomes completion-driven.
//   - N <  0 → invalid input, exit 1 with ErrCloneMaxConcurrencyInvalid.
//
// The auto-resolve sentinel (0) is preferred over hard-coding
// runtime.NumCPU() at flag-default time so help text stays portable
// across machines (`--max-concurrency` describes auto, not "8").
//
// The single-repo `clone-pick` command intentionally does NOT
// surface this flag: there is exactly one URL to clone, so any
// worker count > 1 would be misleading.

// CloneFlagMaxConcurrency is the long-form flag name shared by
// every batch clone command.
const CloneFlagMaxConcurrency = "max-concurrency"

// CloneDefaultMaxConcurrency is the sentinel default. 0 means
// "auto" (resolved by ResolveCloneConcurrency to runtime.NumCPU at
// run time). The historical sequential behavior is still reachable
// with `--max-concurrency 1`.
const CloneDefaultMaxConcurrency = 0

// CloneAutoMaxConcurrency is the public name of the sentinel; use
// it at comparison sites instead of repeating the literal `0`.
const CloneAutoMaxConcurrency = 0

// FlagDescCloneMaxConcurrency is the help text shown by every
// `gitmap <clone-cmd> --help`. Kept short so it fits one line in
// `flag.PrintDefaults`.
const FlagDescCloneMaxConcurrency = "Run up to N clones in parallel (0 = auto / NumCPU, 1 = sequential). Hierarchy is preserved at any N."

// MsgCloneConcurrencyEnabledFmt is printed once before the first
// progress line when the parallel runner takes over. Single stable
// line that scripts can grep for. Resolved (effective) worker
// count is substituted, NOT the raw flag value.
const MsgCloneConcurrencyEnabledFmt = "  ↪ parallel clone enabled: %d workers\n"

// ErrCloneMaxConcurrencyInvalid is printed when the user supplies a
// non-positive integer to --max-concurrency. The CLI exits 1 to keep
// the contract: invalid input never silently degrades to a default.
const ErrCloneMaxConcurrencyInvalid = "clone --max-concurrency: must be a non-negative integer (got %d)\n"
