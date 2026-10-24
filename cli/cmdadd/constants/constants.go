package constants

// Command names and aliases for cmdadd.
const (
	CmdAdd      = "add"
	CmdAddAlias = "ad"
)

// Subcommand identifiers for cmdadd.
const (
	SubCmdRepo         = "repo"
	SubCmdFolder       = "folder"
	SubCmdLfs          = "lfs"
	SubCmdHook         = "hook"
	SubCmdCommonAttr   = "common-attr"
	SubCmdCommonIgnore = "common-ignore"
	SubCmdAttributes   = "attributes"
	SubCmdIgnore       = "ignore"
	SubCmdLfsInstall   = "lfs-install"
)

// File names and permission modes.
const (
	FileGitAttributes = ".gitattributes"
	FileGitIgnore     = ".gitignore"
	DefaultFileMode   = 0644
)

// Builtin file templates for standard configurations.
const (
	TemplateCommonAttr   = "* text=auto eol=lf\n*.png binary\n*.jpg binary\n*.jpeg binary\n*.gif binary\n*.ico binary\n"
	TemplateCommonIgnore = "node_modules/\n.DS_Store\nThumbs.db\n.idea/\n.vscode/\n*.log\n*.tmp\n*.bak\n.env\ndist/\nbuild/\ncoverage/\n"
)

// Operation and component names for telemetry and logging.
const (
	OpAddRun            = "AddRun"
	OpCmdAdd            = "CmdAdd"
	OpWriteCommonAttr   = "writeCommonAttr"
	OpWriteCommonIgnore = "writeCommonIgnore"
)

// Error codes, prefixes, and messages.
const (
	ErrCodeAddMissingArg    = "E_ADD_MISSING_ARG"
	ErrCodeAddUnknownTarget = "E_ADD_UNKNOWN_TARGET"
	ErrPrefixError          = "Error:"
	ErrMsgWriteAttrFailed   = "writeCommonAttr: failed to write"
	ErrMsgWriteIgnoreFailed = "writeCommonIgnore: failed to write"
)

// Success banners and output messages.
const (
	MsgCreatedGitAttributes = "Created .gitattributes"
	MsgCreatedGitIgnore     = "Created .gitignore"
)

// Key identifiers used in structured error payloads.
const (
	KeyTarget = "target"
)
