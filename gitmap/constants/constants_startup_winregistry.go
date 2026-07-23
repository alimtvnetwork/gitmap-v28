package constants

// Windows-specific startup constants. Kept in their own file so the
// cross-platform startup constants in constants_startup.go don't
// have to grow Windows-only fields, and so a future change to the
// registry path / .lnk folder shape touches one obviously-named
// file. None of these are referenced on Linux/macOS code paths —
// they exist as plain string constants on every OS so tests and
// rendering helpers compile cross-platform without build tags.

// Windows backends. The user picks one via `--backend=<value>` on
// `gitmap startup-add` (default: registry). `gitmap startup-list`
// and `gitmap startup-remove` enumerate / search BOTH backends so
// users don't have to remember which one a given entry lives in.
const (
	// StartupBackendRegistry writes ONE direct value to HKCU\
	// Software\Microsoft\Windows\CurrentVersion\Run and tracks
	// ownership in a SEPARATE scope under HKCU\Software\Gitmap\
	// StartupRegistry\<name>. The Run key never carries a sibling
	// marker — Windows would dispatch it as an autostart command
	// at every login.
	StartupBackendRegistry = "registry"
	// StartupBackendStartupFolder writes a .lnk file into
	// %APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup,
	// with a sibling tracking subkey under HKCU\Software\Gitmap\
	// StartupFolder\<name>.
	StartupBackendStartupFolder = "startup-folder"
	// StartupBackendRegistryHKLM mirrors StartupBackendRegistry but
	// targets the MACHINE-WIDE Run key under HKEY_LOCAL_MACHINE
	// instead of the per-user HKEY_CURRENT_USER. The autostart value
	// fires for EVERY interactive user that logs into the machine,
	// so writes require administrator privileges (UAC elevation).
	// Tracking metadata is stored under
	// HKLM\Software\Gitmap\StartupRegistry\<name> so a non-admin
	// reader can still discover ownership without touching HKCU.
	// Reads (list) work without elevation; add/remove require admin
	// and surface ErrStartupHKLMNotAdmin up-front when the current
	// process token is not elevated.
	StartupBackendRegistryHKLM = "registry-hklm"
)

// Registry paths. The same RUN-key relative path
// (`Software\Microsoft\Windows\CurrentVersion\Run`) is used under
// BOTH HKCU (per-user, default) and HKLM (machine-wide, opt-in via
// `--backend=registry-hklm`). Tracking metadata mirrors the same
// shape under both hives so a non-admin reader can discover
// ownership of HKCU entries without touching HKLM and vice versa.
const (
	// RegRunKeyPath is the canonical autostart Run key, used under
	// HKCU for the default registry backend AND under HKLM for the
	// machine-wide registry-hklm backend. Values placed here
	// execute once at the next interactive login. We deliberately
	// avoid RunOnce (single-execution) and RunOnceEx (chained
	// execution) — both have surprising auto-deletion semantics
	// that conflict with idempotent `gitmap startup-add` re-runs.
	RegRunKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
	// RegGitmapRoot is the parent for all gitmap tracking metadata
	// (registry-backend AND startup-folder-backend). The same
	// relative path is used under HKCU and HKLM. Two leaf subkeys
	// live underneath each hive:
	//   StartupRegistry\<name> → tracks Run-key entries
	//   StartupFolder\<name>   → tracks .lnk Startup-folder entries
	// Each leaf carries metadata values (CreatedAt, Exec, Source)
	// that the Linux marker line and macOS plist key encode inline.
	RegGitmapRoot          = `Software\Gitmap`
	RegGitmapRegistrySub   = `Software\Gitmap\StartupRegistry`
	RegGitmapStartupFolder = `Software\Gitmap\StartupFolder`
	// RegMarkerSiblingSuffix is DEPRECATED and no longer written.
	// The direct-value Run-key model (gitmap v3.175.0+) keeps the
	// ownership marker out-of-band under <hive>\Software\Gitmap so
	// the Run key contains only real autostart commands. Kept as
	// a constant so external cleanup scripts that historically
	// removed `<name>.gitmap-managed` companions can still
	// reference the suffix without code drift.
	RegMarkerSiblingSuffix = `.gitmap-managed`
	// Tracking-subkey value names. Stored as REG_SZ so `reg query`
	// shows them readable; CreatedAt is RFC3339 UTC. WorkingDir is
	// only written when the user passed --working-dir; an empty
	// value is omitted entirely so `reg query` output stays tidy
	// for the common no-cwd case.
	RegTrackKeyExec       = "Exec"
	RegTrackKeyCreatedAt  = "CreatedAt"
	RegTrackKeySource     = "Source" // "registry" | "registry-hklm" | "startup-folder"
	RegTrackKeyWorkingDir = "WorkingDir"
)

// Windows file naming. The Run-key value name and the .lnk filename
// both use the same `gitmap-<name>` form so List can recognize them
// from the name alone (cheap pre-filter), then re-check the
// tracking subkey for definitive ownership. Same two-gate model as
// Linux/macOS.
const (
	StartupWinValuePrefix = "gitmap-"
	StartupLnkExt         = ".lnk"
)

// Flag values + descriptions for the --backend flag on
// `gitmap startup-add` / `startup-remove` / `startup-list`. Listed
// at the package's existing flag block for discoverability via
// `gitmap startup-add --help`.
const (
	FlagStartupAddBackend     = "backend"
	FlagDescStartupAddBackend = "Windows backend: registry (default, HKCU per-user), " +
		"registry-hklm (HKLM machine-wide; requires admin), or startup-folder"
	ErrStartupAddBadBackend = "startup-add: unknown --backend %q " +
		"(expected: registry, registry-hklm, startup-folder)"
)

// Windows user-visible messages. Plain ASCII glyphs to match the
// existing PowerShell-encoding constraint (mem://constraints/
// powershell-encoding) and the bullet/arrow style from
// constants_startup.go.
const (
	MsgStartupListHeaderWindows  = "Windows autostart entries managed by gitmap:\n"
	MsgStartupListBackendSection = "  [%s]\n"
	ErrStartupRegistryOpen       = "open registry key %s: %v"
	ErrStartupRegistryWrite      = "write registry value %s: %v"
	ErrStartupRegistryRead       = "read registry value %s: %v"
	ErrStartupShortcutCreate     = "create shortcut %s: %v"
	ErrStartupPowerShellMissing  = "powershell.exe not found on PATH " +
		"(required for --backend=startup-folder; install Windows " +
		"PowerShell or use --backend=registry)"
	// ErrStartupHKLMNotAdmin is surfaced BEFORE any registry write
	// attempt when --backend=registry-hklm is requested but the
	// current process token is not elevated. Keeping the check
	// up-front (instead of relying on the registry ACL to refuse
	// SET_VALUE) means the user sees a friendly, actionable
	// message instead of a raw "Access is denied" Win32 error and
	// no half-written tracking metadata is ever left behind.
	ErrStartupHKLMNotAdmin = "startup-add: --backend=registry-hklm requires " +
		"administrator privileges (re-run from an elevated shell, e.g. " +
		"`Run as administrator` from the Start menu, or use the per-user " +
		"`--backend=registry` default)"
	StartupFolderRelative = `Microsoft\Windows\Start Menu\Programs\Startup`
)
