package constants

// Clean-code / coding-guidelines installer.
// One-liner published at the URL below installs the alimtvnetwork
// coding-guidelines (v15) into the current directory via PowerShell IRM | IEX.
//
// The four CLI aliases (clean-code, code-guide, cg, cc) all dispatch to the
// same flow — see gitmap/cmd/installcleancode.go.
const (
	DefaultCleanCodeURL = "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v15/main/install.ps1"
)

// Clean-code installer messages.
const (
	MsgCleanCodeRunning = "  Installing coding guidelines from %s\n"
	MsgCleanCodeDone    = "  OK Coding guidelines installed.\n"
	MsgCleanCodeNoPwsh  = "  ✗ PowerShell not found on PATH. Install PowerShell 7+ or run the one-liner manually:\n      irm %s | iex\n"
	ErrCleanCodeFailed  = "  ✗ Coding guidelines install failed: %v\n"
	MsgCleanCodeNonWin  = "  Note: this installer is PowerShell-based; on non-Windows it requires PowerShell 7+ (pwsh).\n"
)

// gitmap:cmd top-level
// Clean-code installer alias tokens exposed to shell tab-completion. These
// are not standalone top-level commands — they are argument values to
// `gitmap install` (e.g. `gitmap install clean-code`, `gitmap i cc`). We
// expose them through the completion marker block so users get tab-complete
// hints when typing `gitmap install <TAB>`.
//
// NOTE: `cg` is intentionally OMITTED here because it is already owned by
// `CmdChangelogGenAlias` as a top-level command. Adding it twice would
// silently win the dedupe in the generator but mask the conflict at the
// dispatch layer. Users can still run `gitmap install cg` — the install
// command parses its own positional argument and routes it to the
// clean-code installer via cleanCodeAliases (see installcleancode.go).
const (
	CmdInstallCleanCode      = "clean-code"
	CmdInstallCleanCodeGuide = "code-guide"
	CmdInstallCleanCodeCC    = "cc"
)
