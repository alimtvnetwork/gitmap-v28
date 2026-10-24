package constants

// Command names and aliases for cmddocs.
const (
	CmdDocs      = "docs"
	CmdDocsAlias = "d"
)

// URLs, hosts, and default network endpoints.
const (
	DocsURL           = "https://gitmap.dev/docs"
	DefaultServerHost = "127.0.0.1"
	DefaultServerPort = 8080
)

// Flag names used for CLI parsing in cmddocs.
const (
	FlagBrowser   = "browser"
	FlagTerminal  = "terminal"
	FlagServer    = "server"
	FlagCategory  = "category"
	FlagFormat    = "format"
	FlagPort      = "port"
	FlagHelp      = "help"
	FlagHelpShort = "h"
)

// Flag CLI argument strings.
const (
	FlagArgBrowser   = "--browser"
	FlagArgTerminal  = "--terminal"
	FlagArgServer    = "--server"
	FlagArgCategory  = "--category"
	FlagArgFormat    = "--format"
	FlagArgPort      = "--port"
	FlagArgHelp      = "--help"
	FlagArgHelpShort = "-h"
)

// Flag descriptions for help text and flag parsing.
const (
	FlagDescBrowser  = "Open documentation in default web browser"
	FlagDescTerminal = "Render documentation in terminal"
	FlagDescServer   = "Start local documentation HTTP server"
	FlagDescCategory = "Documentation category filter"
	FlagDescFormat   = "Documentation output format"
	FlagDescPort     = "HTTP server port number"
	FlagDescHelp     = "Show documentation help text"
)

// OS identifiers and execution tokens for platform-specific browser launching.
const (
	OSWindows       = "windows"
	OSDarwin        = "darwin"
	OSLinux         = "linux"
	CmdWindowsShell = "cmd"
	CmdArgSlashC    = "/c"
	CmdArgStart     = "start"
	CmdOpen         = "open"
	CmdXdgOpen      = "xdg-open"
)

// Output messages and status banners.
const (
	MsgDocsOpened         = "  ✓ Opened %s\n"
	MsgDocsServerStarting = "  Starting documentation server on http://%s:%d\n"
)

// Error messages and format strings.
const (
	ErrDocsOpen       = "  ✗ Failed to open browser: %v\n"
	ErrDocsOpenFailed = "failed to open docs browser"
)

// Error codes for cmddocs execution failures.
const (
	ErrCodeDocsLaunchFailed = "EDOCS01"
	ErrCodeDocsInvalidArgs   = "EDOCS02"
	ErrCodeDocsServerFailed = "EDOCS03"
)

// Help documentation strings and usage specifications.
const (
	HelpDocs  = "  docs (d)            Open documentation website in default browser"
	HelpUsage = "gitmap docs [--browser] [--terminal] [--server] [--category <name>] [--format <fmt>]"
)

// Exit codes for cmddocs execution.
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
)

// AllDocsFlags returns all supported flag names for cmddocs.
func AllDocsFlags() []string {
	return []string{
		FlagBrowser,
		FlagTerminal,
		FlagServer,
		FlagCategory,
		FlagFormat,
		FlagPort,
		FlagHelp,
	}
}
