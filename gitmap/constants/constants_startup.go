// Package constants — Linux/Unix startup-management commands and the
// marker convention used to scope list/remove operations to entries
// gitmap itself created (never third-party autostart files).
package constants

// gitmap:cmd top-level
// Startup CLI commands. The pair `startup-list` / `startup-remove`
// is intentionally limited in scope on Linux/Unix: it ONLY enumerates
// and deletes entries gitmap created itself, identified by the
// `X-Gitmap-Managed=true` key inside the .desktop file. This matches
// the user requirement that no third-party autostart entry is ever
// touched, even if its filename happens to start with "gitmap".
//
// Aliases mirror the install/uninstall short-form pattern so users
// who already know `in`/`un` can guess `sl`/`sr` without lookup.
const (
	CmdStartupAdd         = "startup-add"
	CmdStartupAddAlias    = "sa"
	CmdStartupList        = "startup-list"
	CmdStartupListAlias   = "sl"
	CmdStartupRemove      = "startup-remove"
	CmdStartupRemoveAlias = "sr"
)

// Startup help text. Format mirrors HelpInstall / HelpUninstall so
// the rendered `gitmap help` table stays visually consistent.
const (
	HelpStartupAdd    = "  startup-add (sa)          Create a Linux/Unix or macOS autostart entry pointing at gitmap"
	HelpStartupList   = "  startup-list (sl)         List Linux/Unix autostart entries created by gitmap"
	HelpStartupRemove = "  startup-remove (sr) <name> Remove a gitmap-managed autostart entry by name"
)

// Startup .desktop file convention.
//
// StartupDesktopExt is the file extension XDG autostart entries MUST
// use; the directory scanner only considers files ending in this
// suffix to avoid matching swap / backup files.
//
// StartupMarkerKey is the .desktop key whose presence with value
// "true" marks an entry as gitmap-managed. We deliberately use a
// custom `X-` namespaced key (per the freedesktop.org spec — keys
// prefixed with `X-` are reserved for application-specific use) so a
// future tool with the same key cannot collide with us.
//
// StartupFilePrefix is the filename prefix every gitmap-created
// .desktop file gets. It's a SECONDARY safeguard — list/remove still
// require StartupMarkerKey to be present, so renaming a third-party
// file to start with "gitmap-" does NOT make us touch it.
const (
	StartupDesktopExt = ".desktop"
	StartupMarkerKey  = "X-Gitmap-Managed"
	StartupMarkerVal  = "true"
	StartupFilePrefix = "gitmap-"
)

// macOS LaunchAgent convention. The marker uses a hyphen-free key
// (`XGitmapManaged`) because plist keys are XML element text and
// keeping them simple avoids any escaping ambiguity in tests / hand
// edits. The on-disk filename prefix is `gitmap.` (reverse-DNS-ish)
// so .plist files match the LaunchAgent convention while still
// being cheap to pre-filter — same two-gate safety as Linux: the
// filename prefix is a hint, the in-file marker is the proof.
const (
	StartupPlistExt    = ".plist"
	StartupPlistMarker = "XGitmapManaged"
	StartupPlistPrefix = "gitmap."
)

// Startup user-visible messages. Plain ASCII arrows / glyphs to stay
// safe across terminals that don't render Unicode (matches the
// PowerShell-encoding constraint documented in mem://constraints/
// powershell-encoding even though this is a Linux command — keeps
// the message style uniform across the whole CLI).
const (
	MsgStartupListHeader    = "Linux/Unix autostart entries managed by gitmap (%s):\n"
	MsgStartupListEmpty     = "  (none — no gitmap-managed autostart entries found)\n"
	MsgStartupListRow       = "  • %s  →  %s\n"
	MsgStartupListFooter    = "\nTotal: %d entry(ies). Remove one with: gitmap startup-remove <name>\n"
	MsgStartupRemoveOK      = "✓ Removed gitmap-managed autostart entry: %s\n"
	MsgStartupRemoveNoOp    = "  (no-op) no gitmap-managed entry named %q found\n"
	MsgStartupRemoveNotOurs = "  (refused) %q exists but was NOT created by gitmap — skipping\n"
	MsgStartupRemoveBadName = "  (refused) name %q is empty or contains a path separator\n"
	// startup-remove --dry-run mirror messages. Same four outcomes,
	// each prefixed with `(dry-run)` so log-scrapers can tell a
	// preview from a real action without parsing flags.
	MsgStartupRemoveDryOK      = "  (dry-run) would remove gitmap-managed autostart entry: %s\n"
	MsgStartupRemoveDryNoOp    = "  (dry-run) no gitmap-managed entry named %q found — nothing to remove\n"
	MsgStartupRemoveDryNotOurs = "  (dry-run) %q exists but was NOT created by gitmap — would refuse\n"
	MsgStartupRemoveDryBadName = "  (dry-run) name %q is empty or contains a path separator — would refuse\n"
	// startup-add result messages. One line per outcome so log
	// scrapers can grep on the leading symbol/prefix.
	MsgStartupAddCreated     = "✓ Created gitmap-managed autostart entry: %s\n"
	MsgStartupAddOverwritten = "✓ Overwrote gitmap-managed autostart entry: %s\n"
	MsgStartupAddExists      = "  (exists) gitmap-managed entry already at %s — pass --force to overwrite\n"
	MsgStartupAddRefused     = "  (refused) %q exists but was NOT created by gitmap — refusing to overwrite\n"
	MsgStartupAddBadName     = "  (refused) name %q is empty or contains a path separator\n"
	ErrStartupResolveDir     = "could not resolve autostart directory: %v"
	ErrStartupReadDir        = "could not read autostart directory %s: %v"
	ErrStartupRemoveUsage    = "usage: gitmap startup-remove <name>"
	ErrStartupAddMissingExec = "startup-add: --exec is required " +
		"(or run from an installed gitmap binary so we can auto-detect it)"
	ErrStartupUnsupportedOS = "startup commands are not supported on this OS " +
		"(Linux/Unix XDG autostart and macOS LaunchAgents are supported; " +
		"on Windows use the Windows startup commands)"
)

// startup-add CLI flag descriptions. Kept here (not in
// constants_cli.go) so all startup-related strings live together and
// the flag parser file imports just one constants block.
const (
	FlagStartupAddName           = "name"
	FlagStartupAddExec           = "exec"
	FlagStartupAddDisplay        = "display-name"
	FlagStartupAddComment        = "comment"
	FlagStartupAddNoDisplay      = "no-display"
	FlagStartupAddForce          = "force"
	FlagStartupAddWorkingDir     = "working-dir"
	FlagDescStartupAddName       = "Logical name for the entry (filename becomes gitmap-<name>.desktop). Required."
	FlagDescStartupAddExec       = "Command to run at login (default: path to the running gitmap binary)"
	FlagDescStartupAddDisplay    = "Override the Name= field shown in desktop session managers"
	FlagDescStartupAddComment    = "Optional Comment= field text"
	FlagDescStartupAddNoDisplay  = "Set NoDisplay=true so the entry stays out of app menus"
	FlagDescStartupAddForce      = "Overwrite an existing gitmap-managed entry (never overwrites third-party files)"
	FlagDescStartupAddWorkingDir = "Working directory the entry runs in " +
		"(Linux Path=, macOS WorkingDirectory, Windows tracking-subkey WorkingDir)"
)

// startup-add / startup-remove --output flag. Shared so both
// commands accept the same value set ("terminal" preserves the
// legacy human-readable lines; "json" emits the consistent
// startupStatus object — see cmd/startupstatusjson.go for the
// full schema). Default is "terminal" so existing scripts keep
// working byte-for-byte.
const (
	FlagStartupOutput         = "output"
	FlagDescStartupOutput     = "Output mode: terminal (default human lines) or json (status object)"
	ErrStartupBadOutput       = "%s: unknown --output %q (expected: terminal, json)"
	FlagStartupJSONIndent     = "json-indent"
	FlagDescStartupJSONIndent = "Spaces per indent level for --output=json (0 = minified, default 2). Ignored for --output=terminal."
)

// startup-list CLI flag. Reuses the project-wide OutputTerminal/CSV/
// JSON constants for values; "table" is accepted as an alias for the
// default human-readable rendering ("terminal" works too).
const (
	FlagStartupListFormat     = "format"
	FlagDescStartupListFormat = "Output format: table (default), json, jsonl, or csv"
	StartupListFormatTable    = "table"
	// StartupListFormatJSONL is the jq/fluentd/BigQuery-friendly
	// line-oriented variant of --format=json: one compact object
	// per line, no array wrapper, empty list emits zero bytes.
	// Stable key order within each object is guaranteed via
	// stablejson.WriteJSONLines (same Field-slice contract as JSON).
	StartupListFormatJSONL  = "jsonl"
	ErrStartupListBadFormat = "startup-list: unknown --format %q (expected: table, json, jsonl, csv)"
	// CSV header row for `--format csv`. Kept here so the column
	// order is centralized and tests can reference it.
	StartupListCSVHeader = "name,path,exec"
	// Filter flags. --backend scopes the listing to one Windows
	// backend (registry or startup-folder) — Linux/macOS callers
	// can pass either value but only entries that came from that
	// backend will appear, which is always zero on those OSes
	// (their single canonical backend is neither). --name matches
	// against the entry's logical name (the same value passed to
	// startup-add --name) so a user can verify a specific entry
	// without grepping table output.
	FlagStartupListBackend     = "backend"
	FlagDescStartupListBackend = "Filter by backend: registry, registry-hklm, " +
		"or startup-folder (default: all backends)"
	FlagStartupListName     = "name"
	FlagDescStartupListName = "Filter by logical entry name (same form as " +
		"`startup-add --name`); exact match"
	ErrStartupListBadBackend = "startup-list: unknown --backend %q " +
		"(expected: registry, registry-hklm, startup-folder)"
)

// startup-list --json-indent flag. Controls whitespace in
// `--format=json` output without ever changing key order. Accepted
// values: integers 0..8 (0 = minified single-line; 2 = the long-
// standing pretty-printed default; 4/8 for editors that prefer
// wider indents). Negative or >8 values are rejected at parse time.
//
// Silently IGNORED for `--format=table`, `--format=csv`, and
// `--format=jsonl` (which is line-oriented and minified by design).
// Silent rather than error so shell scripts that always pass
// `--json-indent=N --format=$F` for varying $F don't have to branch.
const (
	FlagStartupListJSONIndent     = "json-indent"
	FlagDescStartupListJSONIndent = "Spaces per indent level for --format=json (0 = minified, default 2). Ignored for non-json formats."
	StartupListJSONIndentDefault  = 2
	StartupListJSONIndentMax      = 8
	ErrStartupListBadJSONIndent   = "startup-list: --json-indent %d out of range (expected: 0..8)"
)

// startup-remove CLI flags. --backend mirrors the same flag on
// startup-add so a user can scope a removal to one Windows backend
// (registry or startup-folder) instead of the legacy dual-backend
// fallback. Linux/macOS callers ignore --backend — there's only one
// autostart backend per OS.
const (
	FlagStartupRemoveDryRun      = "dry-run"
	FlagDescStartupRemoveDryRun  = "Show what would be deleted (or refused/no-op) without touching the filesystem"
	FlagStartupRemoveBackend     = "backend"
	FlagDescStartupRemoveBackend = "Windows backend to remove from: registry, " +
		"registry-hklm (requires admin), or startup-folder " +
		"(default: try all — registry, then registry-hklm, then startup-folder)"
)

// EnvStartupAutostartDir is an explicit override for AutostartDir().
// When set and non-empty, AutostartDir() returns this path verbatim
// on every OS (Linux, macOS, Windows), bypassing the per-OS default
// resolution (XDG_CONFIG_HOME/autostart, ~/Library/LaunchAgents,
// Registry/Startup folder respectively).
//
// Two intended uses:
//
//  1. Tests that need a writable, isolated autostart dir on every OS
//     without per-OS skips. Set via t.Setenv inside the test body.
//  2. Power users who keep their dotfiles under a non-default root
//     and want gitmap to install autostart entries alongside them
//     (e.g. inside a Nix profile or a portable USB-stick install).
const EnvStartupAutostartDir = "GITMAP_AUTOSTART_DIR"
