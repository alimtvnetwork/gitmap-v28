package constants

// Subcommand name and alias constants.
const (
	CmdSync          = "sync"
	CmdSyncAlias     = "sy"
	CmdCommon        = "common"
	CmdCommonAlias   = "co"
	CmdAddLFSInstall = "add-lfs-install"
)

// CLI flag literals.
const (
	FlagProjects      = "projects"
	FlagProjectsShort = "p"
	FlagRepo          = "repo"
	FlagRepoShort     = "r"
	FlagWorkers       = "workers"
	FlagWorkersShort  = "w"
	FlagDryRun        = "dry-run"
	FlagDryRunShort   = "n"
	FlagNoPush        = "no-push"
	FlagNoRelease     = "no-release"
	FlagVerbose       = "verbose"
	FlagForce         = "force"
	FlagForceShort    = "f"
	FlagHelp          = "help"
	FlagHelpShort     = "h"
	FlagLongHelp      = "--help"
	FlagShortHelp     = "-h"
	FlagLongDryRun    = "--dry-run"
	FlagShortDryRun   = "-n"
	FlagLongForce     = "--force"
	FlagShortForce    = "-f"
	FlagNoRebase      = "--no-rebase"
	FlagPorcelain     = "--porcelain"
	FlagAllShort      = "-A"
	FlagMsgShort      = "-m"
	FlagBranchShort   = "-b"
	FlagLocal         = "--local"
)

// Target scope constants.
const (
	TargetAll            = "all"
	TargetIgnore         = "ignore"
	TargetAttributes     = "attributes"
	TargetLFSInstall     = "lfs-install"
	TargetPrettierIgnore = "prettier-ignore"
	TargetPrettierRC     = "prettier-rc"
)

// File and directory name tokens.
const (
	FileGitignore        = ".gitignore"
	FileGitattributes    = ".gitattributes"
	FilePrettierignore   = ".prettierignore"
	FilePrettierRC       = ".prettierrc"
	FileCredentialsJSON  = "credentials.json"
	FileGitmapConfigJSON = "gitmap.config.json"
	DirGit               = ".git"
	DirPrompts           = "01-prompts"
	DirSpec              = "02-spec"
	DirAiScripts         = "03-ai-scripts"
	DirAgentsSkills      = ".agents/skills"
	DirCursorSkills      = ".cursor/skills"
	DirAgentsScripts     = ".agents/scripts"
	DirArchive           = "06-archive"
	DirOldPrompts        = "06-old-prompts"
)

// Outcome status constants for RepoSyncResult.
const (
	StatusOk      = "OK"
	StatusFail    = "FAIL"
	StatusSkipped = "SKIPPED"
)

// Discrete action outcome messages.
const (
	ActionNoChanges          = "No Changes"
	ActionAlreadySynced      = "Already Synced"
	ActionDryRun             = "Dry Run"
	ActionMirrorFailed       = "Mirror Failed"
	ActionGitAddFailed       = "Git Add Failed"
	ActionCommitFailed       = "Commit Failed"
	ActionCommitted          = "Committed"
	ActionCommittedAndPushed = "Committed & Pushed"
	ActionTagFailed          = "Committed (Tag Failed)"
	ActionReleasedLocally    = "Released Locally"
	ActionReleasedAndPushed  = "Released & Pushed"
	ActionPushFailed         = "Push Failed"
	ActionSkippedExternal    = "Skipped (External Repo)"
)

// Git operation and environment defaults.
const (
	DefaultBranch        = "main"
	DefaultInitialTag    = "v0.1.0"
	DefaultFallbackTag   = "v0.1.1"
	DefaultCommitMsg     = "chore(sync): mirror canonical prompts, skills, specs"
	AddLFSInstallTag     = "lfs/common"
	DefaultWorkers       = 8
	MinWorkers           = 4
	MaxWorkers           = 16
	CommentAddedByCommon = "\n# added by gitmap common\n"
)

// Git CLI command names and arguments.
const (
	GitCmdPull     = "pull"
	GitCmdStatus   = "status"
	GitCmdAdd      = "add"
	GitCmdCommit   = "commit"
	GitCmdTag      = "tag"
	GitCmdPush     = "push"
	GitCmdCheckout = "checkout"
	GitRemoteOrigin = "origin"
)

// Application error code constants for cmdsync.
const (
	ErrCodeUnknownTarget     = "E1145"
	ErrCodeWriteLines        = "E1146"
	ErrCodeMarshalPrettierRC = "E1147"
	ErrCodeWritePrettierRC   = "E1148"
	ErrCodeParsePrettierRC   = "E1149"
)

// UI and summary layout labels.
const (
	SummaryHeader    = "MULTI-REPOSITORY SYNCHRONIZATION SUMMARY"
	SummaryColRepo   = "Repository"
	SummaryColStatus = "Status"
	SummaryColPreTag = "Pre-Tag"
	SummaryColPostTag = "Post-Tag"
	SummaryColFiles  = "Files"
	SummaryColAction = "Action"
	MenuTitle        = "Repository Baseline Sync (gitmap sync)"
	MenuSection      = "Configuration Sync Targets"
)
