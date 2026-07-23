package constants

// Coding Guidelines v24 installer (used by `cfr cg` / `cfrp cg`).
//
// Distinct from `install clean-code` (v15, PowerShell-only), the v24
// integration ships an OS-aware installer:
//
//	Windows  -> PowerShell one-liner: irm <URL> | iex
//	Unix     -> bash one-liner:       curl -fsSL <URL> | bash
//
// URLs are pinned here so the runner (gitmap/cmd/codingguidelines.go) and
// modifier parser (gitmap/cmd/clonefixrepo_modifiers.go) share a single
// source of truth. Do not hard-code these strings elsewhere; import from
// constants per project Core rule ("No magic strings.").
const (
	DefaultCodingGuidelinesURLWindows = "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/error-manage-install.ps1"
	DefaultCodingGuidelinesURLUnix    = "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh"
)

// Modifier tokens accepted by `cfr` / `cfrp` before the repo URL.
// Order-independent: `cfr cg <url>`, `cfr p cg <url>`, `cfr cg p <url>`
// all parse identically.
const (
	CfrModifierCodingGuidelines = "cg"
	CfrModifierPublic           = "p"
)

// Flag names that opt out of the auto-commit / auto-push step run
// after the v24 installer finishes. Both accept the standard
// `--no-commit` / `--no-push` long forms (single-dash also parsed
// by parseCloneFixRepoArgs, matching the rest of the cfr flags).
const (
	FlagCGNoCommit = "no-commit"
	FlagCGNoPush   = "no-push"
)

// Commit metadata used when auto-committing the installed guidelines.
const (
	CodingGuidelinesCommitMessage = "chore: install coding guidelines (v24)"
	CodingGuidelinesCommitAuthor  = "gitmap"
)

// Runner status + error messages. All output is directed to os.Stderr per
// the zero-swallow error policy; keep the format aligned with the rest of
// the CLI (leading two spaces, tag prefix, single trailing newline).
const (
	MsgCGRunningWindows = "  Installing coding guidelines (v24, Windows) from %s\n"
	MsgCGRunningUnix    = "  Installing coding guidelines (v24, Unix) from %s\n"
	MsgCGDone           = "  OK Coding guidelines (v24) installed.\n"
	MsgCGCommitted      = "  OK Committed coding-guidelines changes: %s\n"
	MsgCGPushed         = "  OK Pushed coding-guidelines commit to %s\n"
	MsgCGSkipCommit     = "  Note: --no-commit set; leaving guideline files uncommitted.\n"
	MsgCGSkipPush       = "  Note: --no-push set (or no upstream); push step skipped.\n"
	MsgCGNoChanges      = "  Note: installer produced no working-tree changes; nothing to commit.\n"

	ErrCGShellNotFoundWindows = "  ✗ PowerShell not found on PATH. Install PowerShell 7+ or run manually:\n      irm %s | iex\n"
	ErrCGShellNotFoundUnix    = "  ✗ bash or curl not found on PATH. Install both or run manually:\n      curl -fsSL %s | bash\n"
	ErrCGInstallFailed        = "  ✗ Coding guidelines (v24) install failed on %s: %v\n"
	ErrCGCompatPrepareFailed  = "  ✗ Could not prepare coding-guidelines compatibility retry: %v\n"
	ErrCGCommitFailed         = "  ✗ Failed to commit coding-guidelines changes: %v\n"
	ErrCGPushFailed           = "  ✗ Failed to push coding-guidelines commit: %v\n"
)
