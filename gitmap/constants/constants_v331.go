// Package constants — constants_v331.go centralizes the new strings introduced
// in v3.31.0 for: cross-dir release (`r <repo> <ver>`), cross-dir clone-next
// (`cn <repo> <ver>`), `has-change` (`hc`), and the SSH existing-key-on-disk
// fix. Kept in a single file (rather than scattered across domain files) to
// keep the v3.31.0 surface area easy to audit and unwind.
package constants

// gitmap:cmd top-level
// New CLI commands.
const (
	CmdHasChange      = "has-change"
	CmdHasChangeAlias = "hc"
)

// has-change flag names + descriptions.
const (
	FlagHCMode      = "mode"
	FlagHCAll       = "all"
	FlagHCFetch     = "fetch"
	FlagDescHCMode  = "Dimension to check: dirty (default), ahead, or behind"
	FlagDescHCAll   = "Print all three dimensions as structured output"
	FlagDescHCFetch = "Run 'git fetch' before checking ahead/behind (default true)"
)

// has-change mode values + literals.
const (
	HCModeDirty  = "dirty"
	HCModeAhead  = "ahead"
	HCModeBehind = "behind"
	HCTrue       = "true"
	HCFalse      = "false"
)

// has-change messages + errors.
const (
	MsgHCAllFmt        = "dirty=%s ahead=%s behind=%s\n"
	MsgHCAllNoUpstream = "dirty=%s ahead=n/a behind=n/a (no upstream)\n"
	ErrHCUsage         = "Usage: gitmap has-change <repo> [--mode dirty|ahead|behind] [--all]"
	ErrHCBadMode       = "  ✗ Unknown mode %q. Use one of: dirty, ahead, behind.\n"
	WarnHCFetchFailed  = "  ⚠ git fetch failed in %s: %v (ahead/behind may be stale)\n"
)

// Cross-dir release (`r <repo> <ver>`) messages + errors.
const (
	MsgRRStartingFmt     = "  → Releasing %s at %s (version %s)...\n"
	MsgRRFetchingFmt     = "  📡 Fetching remote refs in %s...\n"
	MsgRRRebasingFmt     = "  🔁 Pull --rebase in %s...\n"
	MsgRRReturnedFmt     = "  ↩ Returned to %s\n"
	ErrRRFetchFailedFmt  = "  ✗ git fetch failed in %s: %v"
	ErrRRRebaseFailedFmt = "  ✗ git pull --rebase failed in %s: %v\n  Resolve the conflict, then re-run the release."
)

// Cross-dir clone-next (`cn <repo> <ver>`) messages.
const (
	MsgCNXStartingFmt = "  → clone-next for %s at %s (version %s)...\n"
	MsgCNXReturnedFmt = "  ↩ Returned to %s\n"
)

// Folder-arg clone-next (`cn vX <folder>` and `cn <folder>`) — v3.117.0.
// See spec/01-app/111-cn-folder-arg.md for the full classification table.
const (
	// CloneNextDefaultVersionArg is what the single-positional folder
	// form defaults to when the user omits a version. v++ is the
	// established muscle-memory shortcut for "next sibling version".
	CloneNextDefaultVersionArg = "v++"

	// ErrCNFolderNotFound fires when a path-shaped token (slash, ~,
	// or otherwise hint-bearing) does not resolve to an existing
	// directory. Includes the original token (not the resolved abs)
	// so the user sees what they typed verbatim.
	ErrCNFolderNotFound = "  ✗ cn: folder not found or not a directory: %s\n"

	// ErrCNAmbiguousBothVersions fires when both positional args
	// match the version pattern (e.g. `cn v++ v15`). The user almost
	// certainly meant one of them as a folder; refuse instead of
	// silently picking one.
	ErrCNAmbiguousBothVersions = "  ✗ cn: ambiguous arguments — both look like version strings (use 'cn vN' for in-place, 'cn vN <folder>' for cross-dir)"

	// ErrCNAmbiguousBothFolders fires when neither positional arg
	// matches the version pattern AND the second isn't folder-shaped.
	// Most often a typo where the version was misspelled.
	ErrCNAmbiguousBothFolders = "  ✗ cn: ambiguous arguments — neither looks like a version (use vN, v+N, or v++)"
)

// gitmap inject (alias: inj) — v3.119.0. Registers an existing folder
// with GitHub Desktop + VS Code, and (when a remote origin exists)
// upserts into the SQLite database. See gitmap/cmd/inject.go and
// gitmap/helptext/inject.md.
const (
	MsgInjectStart     = "→ injecting %s (%s)\n"
	MsgInjectDone      = "✓ injected %s\n"
	MsgInjectNoRemote  = "  ↳ no remote origin configured for %s — skipping database upsert\n"
	MsgInjectDBOK      = "  ↳ database: upserted %s (%s)\n"
	ErrInjectResolve   = "  ✗ inject: %v\n"
	WarnInjectDBOpen   = "  ⚠ inject: could not open database: %v\n"
	WarnInjectDBUpsert = "  ⚠ inject: could not save repo to database: %v\n"
)

// SSH existing-key-on-disk fix messages.
const (
	MsgSSHExistsOnDisk = "\n  ℹ SSH key already exists on disk: %s\n  Reusing existing key (no regeneration needed).\n\n"
	MsgSSHForceHint    = "\n  💡 Pass --force to back up and regenerate this key.\n"
	MsgSSHBackedUp     = "  💾 Backed up existing key: %s.bak.<timestamp>\n"
	ErrSSHBackup       = "  ✗ Could not back up existing SSH key: %v\n"
)
