package constants

// gitmap:cmd top-level
// CLI commands.
const (
	CmdScan                       = "scan"
	CmdScanAlias                  = "s"
	CmdClone                      = "clone"
	CmdCloneAlias                 = "c"
	CmdCloneSync                  = "clone-sync"
	CmdCloneSyncAlias             = "cs"
	CmdCloneOnlyMissing           = "clone-only-missing"
	CmdCloneOnlyMissingAlias      = "com"
	CmdUpdate                     = "update"
	CmdUpdateRunner               = "update-runner"  // gitmap:cmd skip
	CmdUpdateCleanup              = "update-cleanup" // gitmap:cmd skip
	CmdInstalledDir               = "installed-dir"  // gitmap:cmd skip
	CmdInstalledDirAlias          = "id"
	CmdVersion                    = "version"
	CmdVersionAlias               = "v"
	CmdHelp                       = "help"
	CmdDesktopSync                = "desktop-sync"
	CmdDesktopSyncAlias           = "ds"
	CmdGitHubDesktop              = "github-desktop"
	CmdGitHubDesktopAlias         = "gd"
	CmdPull                       = "pull"
	CmdPullAlias                  = "p"
	CmdPullAll                    = "pull-all"
	CmdPullAllAlias               = "pa"
	CmdPullAllEfficient           = "pull-all-efficient"
	CmdPullAllEfficientAlias      = "pae"
	CmdPullAE                     = "pull-ae"
	CmdPullAllEfficientTable      = "pull-all-efficient-table"
	CmdPullAllEfficientTableAlias = "paet"
	CmdReconcile                  = "reconcile"
	CmdReconcileAlias             = "recon"
	// CmdPush is the cwd-scoped `git push` wrapper added in v5.29.0.
	// Alias is `ph` (NOT `p` — collides with CmdPullAlias). Supports
	// shared `--ssh` / `--https` transport flags.
	CmdPush          = "push"
	CmdPushAlias     = "ph"
	CmdPushFix       = "push-fix"
	CmdPushFixAlias  = "pf"
	CmdPushFixSolid  = "pushfix"
	CmdPushFixInvert = "fix-push"
	CmdRescan        = "rescan"
	CmdRescanAlias   = "rsc"
	// CmdRescanSubtree narrowly re-runs `gitmap scan` against the
	// absolutePath of an at-cap row from a previous scan output. Default
	// MaxDepth is bumped to RescanSubtreeDefaultMaxDepth so users get a
	// deeper view in one step; override with --max-depth.
	CmdRescanSubtree        = "rescan-subtree"
	CmdRescanSubtreeAlias   = "rss"
	CmdSetup                = "setup"
	CmdStatus               = "status"
	CmdStatusAlias          = "st"
	CmdExec                 = "exec"
	CmdExecAlias            = "x"
	CmdRelease              = "release"
	CmdCreate               = "create"
	CmdCreateAlias          = "crt"
	CmdRepoCreate           = "repo-create"
	CmdRepoCreateAlias      = "repoc"
	CmdCreateRepo           = "create-repo"
	CmdCreateRepoAlias      = "crepo"
	CmdCreateLocalRepo      = "create-local-repo"
	CmdCreateLocalRepoAlias = "clr"
	CmdCreateRepoLocal      = "create-repo-local"
	CmdRepoCreateLocal      = "repo-create-local"
	SubCmdRepo              = "repo" // gitmap:cmd skip
	CmdFolder               = "folder"
	CmdLowercase            = "lowercase"
	CmdLowerCaseFix         = "lower-case-fix"
	CmdLowerCaseFixAlias    = "lowercase-fix"
	CmdLowerCaseFixShort    = "lcf"
	CmdLowercaseReadme      = "lowercase-readme"
	CmdReadmeLower          = "readme-lower"
	CmdFixSeqFiles          = "fix-seq-files"
	CmdFixSeqFilesAlias     = "fsf"
	CmdSequence             = "sequence"
	CmdSequenceAlias        = "seq"
	CmdGitRm                = "git-rm"
	CmdIgnore               = "ignore"
	CmdIgnoreRm             = "ignore-rm"
	// CmdCommit stages all changes and creates a flat git commit without forced prefixes.
	CmdCommit       = "commit"
	CmdCommitAlias  = "cm"
	CmdCommitAlias2 = "commit-all"
	CmdCommitAlias3 = "ca"
	// CmdCommitPush stages all changes, commits with a message, and pushes.
	CmdCommitPush      = "commit-push"
	CmdCommitPushAlias = "cp"
	// CmdPullCommitPush pulls first (resolving conflicts), then stages, commits, pushes.
	CmdPullCommitPush      = "pull-commit-push"
	CmdPullCommitPushAlias = "pcp"
	// CmdCommitPushBug commits with a "Bug: " prefix and pushes.
	CmdCommitPushBug       = "commit-push-bug"
	CmdCommitPushBugAlias  = "cpb"
	CmdCommitPushBugAlias2 = "cbf"
	// CmdCommitPushFeature commits with a "Feature: " prefix and pushes.
	CmdCommitPushFeature      = "commit-push-feature"
	CmdCommitPushFeatureAlias = "cpf"
	// CmdCommitPushRefactor stages all changes, commits with a "REFACTOR: " prefix, and pushes.
	CmdCommitPushRefactor      = "commit-push-refactor"
	CmdCommitPushRefactorAlias = "cfr" // gitmap:cmd skip
	// CmdCommitPushChore stages all changes, commits with a "Chore: " prefix, and pushes.
	CmdCommitPushChore      = "commit-push-chore" // gitmap:cmd skip
	CmdCommitPushChoreAlias = "cpc"               // gitmap:cmd skip
	// CmdCommitPushRelease commits with a "Release: " prefix and pushes.
	CmdCommitPushRelease      = "commit-push-release"
	CmdCommitPushReleaseAlias = "cpr"
	// CmdPendingCommits discovers uncommitted changes and unpushed commits across repos.
	CmdPendingCommits      = "pending-commits"
	CmdPendingCommitsAlias = "pc"
	summaryCmdVal          = "summary"
	fullSummaryCmdVal      = "full-summary"
	fullStatusCmdVal       = "full-status"
	fsCmdVal               = "fs"
	fullSummaryPECmdVal    = "full-summary+pe"
	fullStatusPECmdVal     = "full-status+pe"
	fsPECmdVal             = "fs+pe"
	fspeCmdVal             = "fspe"
	mergeAICmdVal          = "merge-ai"

	// CmdSummary summarizes the last N releases, heated files, and changelog gist with Split-DB caching.
	CmdSummary = summaryCmdVal
	// CmdFullSummary displays the activity heatmap TreeView for repositories active in 48h or dirty.
	CmdFullSummary   = fullSummaryCmdVal
	CmdFullStatus    = fullStatusCmdVal
	CmdFs            = fsCmdVal
	CmdFullSummaryPE = fullSummaryPECmdVal
	CmdFullStatusPE  = fullStatusPECmdVal
	CmdFsPE          = fsPECmdVal
	CmdFspe          = fspeCmdVal
	// CmdMergeAI orchestrates multi-repo staging into a single commit with 01_/02_ collision handling.
	CmdMergeAI      = mergeAICmdVal
	CmdMergeAIAlias = "ma" // gitmap:cmd skip
	// CmdNewCommands discovers and catalogs commands added in recent git history.
	CmdNewCommands      = "new-commands" // gitmap:cmd skip
	CmdNewCommandsAlias = "nc"           // gitmap:cmd skip
	// CmdSends dispatches semantic commits across repositories.
	CmdSends = "sends"
	// CmdRmGit removes a commit by its last 4-digit SHA prefix.
	CmdRmGit      = "rm-git"
	CmdRmGitAlias = "rmg"
	// CmdGitReset resets the branch history to a specific SHA code.
	CmdGitReset                   = "git-reset"
	CmdGitResetAlias              = "reset-git"
	CmdAg                         = "ag"
	CmdAntigravity                = "antigravity"
	CmdShutdownUntil              = "shutdown-until"
	CmdShutdownUntilAlias         = "sug"
	CmdShutdownUntilGreen         = "shutdown-until-green"
	CmdWatchPromptsRunning        = "watch-prompts-running"
	CmdWatchPromptsRunningAlias   = "wpr"
	CmdPrompt                     = "prompt"
	CmdPromptAlias                = "prompts"
	CmdAdd                        = "add"
	CmdAppend                     = "append"
	CmdWrite                      = "write"
	CmdHead                       = "head"
	CmdTail                       = "tail"
	CmdFileSearch                 = "file-search"
	CmdInstallScript              = "install-script"
	CmdCargo                      = "cargo"     // gitmap:cmd skip
	CmdCargoAlias                 = "crg"       // gitmap:cmd skip
	CmdRust                       = "rust"      // gitmap:cmd skip
	CmdRustAlias                  = "rst"       // gitmap:cmd skip
	CmdRustAlias2                 = "rustc"     // gitmap:cmd skip
	CmdKill                       = "kill"      // gitmap:cmd skip
	CmdKillAlias                  = "terminate" // gitmap:cmd skip
	CmdReleaseShort               = "r"
	CmdLlm                        = "llm"
	CmdFind                       = "find"
	CmdFindFiles                  = "find-files"
	CmdFindFilesAny               = "find-files-any"
	CmdFindFilesStartsWith        = "find-files-startswith"
	CmdFindFilesEndsWith          = "find-files-endswith"
	CmdListFiles                  = "list-files"
	CmdFindRegex                  = "find-regex"
	CmdFindRead                   = "find-read"
	CmdFindReadJson               = "find-read-json"
	CmdFindRegexRead              = "find-regex-read"
	CmdFindRegexReadJson          = "find-regex-read-json"
	CmdFindHelp                   = "find-help"
	CmdSearchHelp                 = "search-help"
	CmdRegexHelp                  = "regex-help"
	CmdSearch                     = "search"
	CmdReplaceRegex               = "replace-regex"
	CmdRepoSearch                 = "repo-search"
	CmdRepoSearchAlias            = "rps"
	CmdRepoRegex                  = "repo-regex"
	CmdRepoRegexAlias             = "rpr"
	CmdRepoSearchJson             = "repo-search-json"
	CmdRepoSearchJsonAlias        = "rsj"
	CmdRepoSearchRegexJson        = "repo-search-regex-json"
	CmdSearchReplaceAll           = "search-replace-all"
	CmdReleaseBranch              = "release-branch"
	CmdReleaseBranchAlias         = "rb"
	CmdReleasePending             = "release-pending"
	CmdReleasePendingAlias        = "rp"
	CmdReleaseUndo                = "release-undo"
	CmdReleaseUndoAlias           = "ru"
	CmdReleaseRemoveBetween       = "release-remove-between"
	CmdReleaseRemoveBetweenAlias  = "rrb"
	CmdReleasePinNotToRemove      = "release-pin-not-to-remove"
	CmdReleasePinNotToRemoveAlias = "pnr"
	CmdChangelog                  = "changelog"
	CmdChangelogAlias             = "cl"
	CmdChangelogMD                = "changelog.md" // gitmap:cmd skip
	CmdDoctor                     = "doctor"
	CmdBackup                     = "backup"
	// `gitmap backup <subcommand>` lists or prunes the on-disk backup
	// tree at `.gitmap/backup/<repo>/v<N>/fix-repo/<UTC-ts>/`. v6.57.0.
	SubCmdBackupLs    = "ls"    // gitmap:cmd skip
	SubCmdBackupList  = "list"  // gitmap:cmd skip
	SubCmdBackupPrune = "prune" // gitmap:cmd skip
	// Repo hygiene commands (v6.68.0).
	CmdStale      = "stale"
	CmdStaleAlias = "sta"
	CmdOrphans    = "orphans"
	CmdDedupe     = "dedupe"
	CmdSize       = "size"
	// Release tooling + workflow + safety (v6.69.0+).
	CmdReleaseNotes      = "release-notes"
	CmdReleaseDry        = "release-dry"
	CmdTagRename         = "tag-rename"
	CmdRecent            = "recent"
	CmdRecentAlias       = "rct"
	CmdTodo              = "todo"
	CmdOpen              = "open"
	CmdOpenAlias         = "o"
	CmdPullRequests      = "pull-requests"
	CmdPRAlias           = "prs"
	CmdBlameStats        = "blame-stats"
	CmdSnapshot          = "snapshot"
	CmdRollback          = "rollback"
	CmdGuard             = "guard"
	CmdLatestBranch      = "latest-branch"
	CmdLatestBranchAlias = "lb"
	// `gitmap branch <subcommand>` namespaces general branch-management
	// helpers. First subcommand: `default` / `def` — checkout the repo's
	// default branch (origin/HEAD or constants.DefaultBranch fallback).
	// Designed for muscle-memory parity with `gitmap lb -s` (latest) so
	// users can flip between "freshest" and "canonical" with a single
	// short word. Subcommand strings live in constants_messages.go to
	// keep all `branch`-specific tokens together.
	CmdBranch                = "branch"
	CmdBranchAlias           = "b"
	CmdBranchSubDefault      = "default" // gitmap:cmd skip
	CmdBranchSubDefaultAlias = "def"     // gitmap:cmd skip
	CmdList                  = "list"
	CmdListAlias             = "ls"
	CmdGroup                 = "group"
	CmdGroupAlias            = "g"
	CmdWorkgroup             = "workgroup"
	CmdWorkgroupAlias        = "wg"
	SubCmdGroups             = "groups"
	CmdCluster               = "cluster"
	CmdClusterAlias          = "clst"
	CmdClusterStatus         = "status"          // gitmap:cmd skip
	CmdServersClients        = "servers-clients" // gitmap:cmd top-level
	CmdServersClientsAlias   = "servers-client"
	CmdSC                    = "sc"
	CmdClients               = "clients" // gitmap:cmd top-level
	CmdClientsAlias          = "client"
	CmdServersLS             = "servers ls"          // gitmap:cmd top-level
	CmdClientsLS             = "clients ls"          // gitmap:cmd top-level
	CmdSCLS                  = "servers-clients ls"  // gitmap:cmd top-level
	CmdSCCat                 = "servers-clients cat" // gitmap:cmd top-level
	CmdSwapBinary            = "swap-binary"
	CmdSwapBinaryAlias       = "swap-bin"
	CmdSwapBinaryAlias2      = "replace-binary"
	CmdSwapBinaryAlias3      = "replace-bin"
)
