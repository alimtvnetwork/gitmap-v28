package constants

// gitmap:cmd top-level
// Profile CLI commands.
const (
	CmdProfile       = "profile"
	CmdProfileAlias  = "pf"
	CmdProfiles      = "profiles"
	CmdProfilesAlias = "profs"
)

// gitmap:cmd top-level
// Profile subcommands.
const (
	CmdProfileCreate      = "create"       // gitmap:cmd skip
	CmdProfileList        = "list"         // gitmap:cmd skip
	CmdProfileSwitch      = "switch"       // gitmap:cmd skip
	CmdProfileDelete      = "delete"       // gitmap:cmd skip
	CmdProfileShow        = "show"         // gitmap:cmd skip
	CmdProfileImport      = "import"       // gitmap:cmd skip
	CmdProfileImportAll   = "import-all"   // gitmap:cmd skip
	CmdProfileExport      = "export"       // gitmap:cmd skip
	CmdProfileExportAll   = "export-all"   // gitmap:cmd skip
	CmdProfileInspect     = "inspect"      // gitmap:cmd skip
	CmdProfilePreview     = "preview"      // gitmap:cmd skip
	CmdProfileCheck       = "check"        // gitmap:cmd skip
	CmdProfileImportCheck = "import-check" // gitmap:cmd skip
)

// Profile help text.
const HelpProfile = "  profile (pf) <sub>  Manage profiles (create, list, switch, delete, show, import, export, inspect)"

// Profile file and defaults.
const (
	ProfileConfigFile  = "profiles.json"
	DefaultProfileName = "default"
	ProfileDBPrefix    = "gitmap-"
)

// Profile messages.
const (
	MsgProfileCreated       = "Profile created: %s\n"
	MsgProfileSwitched      = "Switched to profile: %s\n"
	MsgProfileDeleted       = "Profile deleted: %s\n"
	MsgProfileActive        = "Active profile: %s\n"
	MsgProfileColumns       = "PROFILE              STATUS"
	MsgProfileRowFmt        = "%-20s %s\n"
	MsgProfileActiveTag     = "(active)"
	MsgProfileEmpty         = "No profiles found.\n"
	ErrProfileUsage         = "usage: gitmap profile <create|list|switch|delete|show|import|export|inspect> [args]\n"
	ErrProfileCreateUsage   = "usage: gitmap profile create <name>\n"
	ErrProfileSwitchUsage   = "usage: gitmap profile switch <name>\n"
	ErrProfileDeleteUsage   = "usage: gitmap profile delete <name>\n"
	ErrProfileNotFound      = "profile not found: %s\n"
	ErrProfileExists        = "profile already exists: %s\n"
	ErrProfileDeleteActive  = "cannot delete the active profile (switch first)\n"
	ErrProfileDeleteDefault = "cannot delete the default profile\n"
	ErrProfileConfig        = "failed to manage profile config: %v\n"
)
