package constants

// gitmap:cmd top-level
// Export CLI commands.
const (
	CmdExport      = "export"
	CmdExportAlias = "ex"
)

// Export help text.
const HelpExport = "  export (ex) [file]  Export full database as portable JSON (default: gitmap-export.json)"

// Export flag descriptions.
const FlagDescExportOut = "Output file path for the export"

// Export default file name.
const DefaultExportFile = "gitmap-export.json"

// Export messages.
const (
	MsgExportDone   = "Database exported to %s (%d repos, %d groups, %d releases, %d history, %d bookmarks)\n"
	MsgExportFailed = "export failed: %v\n"
)
