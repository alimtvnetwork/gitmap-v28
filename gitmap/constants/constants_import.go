package constants

// gitmap:cmd top-level
// Import CLI commands.
const (
	CmdImport      = "import"
	CmdImportAlias = "im"
)

// Import help text.
const HelpImport = "  import (im) [file]  Import database from a gitmap-export.json file (--confirm required)"

// Import messages.
const (
	MsgImportDone        = "Database imported from %s (%d repos, %d groups, %d releases, %d history, %d bookmarks)\n"
	MsgImportFailed      = "import failed: %v\n"
	MsgImportReadFailed  = "failed to read import file: %v\n"
	MsgImportParseFailed = "failed to parse import file: %v\n"
	ErrImportNoConfirm   = "import requires --confirm flag (existing data will be merged)\n"
	MsgImportSkipGroup   = "skipped group %q: missing member repos\n"
)
