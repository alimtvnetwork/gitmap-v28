package constants

// gitmap:cmd top-level
// Additional CLI commands (extended set).
const (
	CmdCompletion            = "completion"
	CmdCompletionAlias       = "cmp"
	CmdClearReleaseJSON      = "clear-release-json"
	CmdClearReleaseJSONAlias = "crj"
	CmdDocs                  = "docs"
	CmdDocsAlias             = "d"
	CmdCloneNext             = "clone-next"
	CmdCloneNextAlias        = "cn"
	CmdReleaseSelf           = "release-self"
	CmdReleaseSelfAlias      = "rself"
	CmdReleaseSelfAlias2     = "rel-self"
	CmdRepoSecrets           = "repo-secrets"
	CmdRepoSecretsAlias      = "rs"
	CmdRepoCache             = "repo-cache"
	CmdRepoCacheAlias        = "rc"
	CmdRepoStorageAlias      = "repo-storage"
	CmdExclude               = "exclude"
	CmdExcludeAlias          = "repo-exclude"
	CmdCleanRepos            = "clean-repos"
	CmdCleanReposAlias       = "repoclean"
	CmdHelpDashboard         = "help-dashboard"
	CmdHelpDashboardAlias    = "hd"
	CmdPending               = "pending"    // gitmap:cmd skip
	CmdDoPending             = "do-pending" // gitmap:cmd skip
	CmdDoPendingAlias        = "dp"         // gitmap:cmd skip
	CmdLLMDocs               = "llm-docs"
	CmdLLMDocsAlias          = "ld"
	CmdWhoAmI                = "whoami"
	CmdWhoAmIAlias           = "who"
	CmdSSHBind               = "ssh-bind"
	CmdSSHBindAlias          = "sb"
	CmdFixAuth               = "fix-auth"
	CmdFixAuthAlias          = "fa"
	CmdFixCredential         = "fix-credential"
	CmdFixCredentialAlias    = "fc"
	CmdSync                  = "sync"
	CmdSyncAlias             = "sy"
	// CmdCommons (v6.76.0) is a shortcut for `gitmap sync all` — adds
	// or dedupe-merges the curated .gitignore, .gitattributes,
	// .prettierignore, .prettierrc baselines + `git lfs install`.
	CmdCommons      = "commons"
	CmdCommonsAlias = "co"
	// CmdSpace is the top-level namespace for GitMap space operations.
	// Subcommands: `space common` applies the curated .gitignore,
	// .gitattributes, .prettierignore, .prettierrc baselines plus
	// `git lfs install --local` (same logic as `gitmap commons`);
	// `space backup-branch` snapshots HEAD into backup/<slug>.
	CmdSpace = "space"
	// CmdSpaceBackupBranch is the `space backup-branch` subcommand —
	// creates backup/<slug> from HEAD for the given task string.
	CmdSpaceBackupBranch = "backup-branch"
	// CmdSpacePRMerge is the `space PR-merge` subcommand (spec 271) —
	// merges the current branch into the target (default branch when
	// omitted); --squash embeds PR details plus a changed-files heatmap.
	CmdSpacePRMerge = "PR-merge"
	// CmdSpaceSendEmail is the `space send-email` subcommand — sends an
	// email via go-email-reader (when installed) or the built-in SMTP
	// sender with prompted authentication.
	CmdSpaceSendEmail = "send-email"
	// CmdSpaceScheduleSendEmail is the `space schedule-send-email`
	// subcommand — registers/lists/removes email triggers that fire
	// when a gitmap task (release/pull/clone/scan/push/backup) completes.
	CmdSpaceScheduleSendEmail = "schedule-send-email"
	// CmdSpec is the top-level namespace for spec-number issuance
	// (spec 252). Subcommand: `spec next` issues the next free spec
	// number for the current repository.
	CmdSpec = "spec" // gitmap:cmd skip
	// CmdSpecNext is the `spec next` subcommand — issues the next
	// free spec number and prints it on stdout.
	CmdSpecNext = "next" // gitmap:cmd skip
	// CmdPlans lists and summarizes completed and pending repository plans.
	CmdPlans             = "plans"
	CmdPlansAlias        = "plan"
	CmdRecreateRepo      = "recreate-repo"
	CmdRecreateRepoAlias = "recreate"
	CmdRestEnable        = "rest-enable"
	CmdVar               = "var"
	CmdVarAlias          = "variable"
	CmdMigrate           = "migrate"

	DefaultSequenceTTLSeconds = 60

	CmdSetSourceRepo     = "set-source-repo" // gitmap:cmd skip
	CmdSelfInstall       = "self-install"
	CmdSelfUninstall     = "self-uninstall"
	CmdSf                = "sf"
	CmdProbe             = "probe"
	CmdCode              = "code"
	CmdCodeAlias         = "vcode"
	CmdCodeAlias2        = "vscode"
	CmdVSCodePMPath      = "vscode-pm-path"
	CmdVSCodePMPathAlias = "vpath"
	// CmdVSCodeWorkspace emits a single `.code-workspace` file listing
	// every tracked repo as a folder, so one click in VS Code opens
	// them all in one window. Source is the same DB the Project
	// Manager sync reads (store.ListRepos), so the file always matches
	// the latest scan output. Help: vscode-workspace.md.
	CmdVSCodeWorkspace      = "vscode-workspace"
	CmdVSCodeWorkspaceAlias = "vsws"
	CmdLFSCommon            = "lfs-common"
	CmdLFSCommonAlias       = "lfsc"
	CmdReplace              = "replace"
	CmdReplaceAlias         = "rpl"
	CmdInject               = "inject"
	CmdInjectAlias          = "inj"
	// CmdOpen / CmdOpenAlias declared earlier in this file (line 75).
	// CmdRegoldens runs `go test -run <pattern>` in two passes: pass 1
	// regenerates golden fixtures with both gate env vars set; pass 2
	// re-runs the same tests with the env vars unset to verify the
	// regenerated bytes are deterministic. Spec:
	// 02-spec/05-coding-guidelines/21-golden-fixture-regeneration.md.
	CmdRegoldens      = "regoldens"
	CmdRegoldensAlias = "rg"
	// CmdAuditLegacy scans the workspace for forbidden legacy strings
	// (e.g. gitmap-v28/v6/v7) and exits non-zero on any hit. Used as a // gitmap-legacy-ref-allow
	// regression guard for remixes / rename commits. Help: audit-legacy.md.
	CmdAuditLegacy       = "audit-legacy"
	CmdAuditLegacyAlias  = "aul"
	CmdAuditLegacyAlias2 = "audit"
	// CmdFixRepo rewrites prior `{base}-vN` tokens in tracked files to
	// the current version, mirroring fix-repo.ps1 / fix-repo.sh exit
	// codes and config schema. Spec: 02-spec/04-generic-cli/27-fix-repo-command.md.
	CmdFixRepo      = "fix-repo"
	CmdFixRepoAlias = "fr"
	// CmdFixGit diagnoses and repairs Git commit, permission, lock, and index corruption issues.
	CmdFixGit      = "fix-git"
	CmdFixGitAlias = "fg"
	// CmdFixReleaseTags audits and cleans orphan or broken GitHub releases and tags lacking assets.
	CmdFixReleaseTags           = "fix-release-tags"
	CmdFixReleaseTagsAlias      = "frt"
	CmdFixReleaseTagsAliasSolid = "fixreleasetags"
	// CmdUndo restores the latest `gitmap fix-repo` backup snapshot
	// for the current repo + current version (v5.40.0+). Backups live
	// at `<repoRoot>/.gitmap/backup/<repo>/v<N>/fix-repo/<UTC-ts>/`.
	// Spec: 02-spec/04-generic-cli/27-fix-repo-command.md §"Backup & undo".
	CmdUndo               = "undo"
	CmdUndoAlias          = "ud"
	CmdVmware             = "vmware"
	CmdVmwareAlias        = "vm"
	SubCmdVmwareShared    = "shared" // gitmap:cmd skip
	SubCmdSharedEnable    = "enable" // gitmap:cmd skip
	SubCmdSharedStatus    = "status" // gitmap:cmd skip
	CmdNginx              = "nginx"
	CmdNginxAlias         = "ngx"
	SubCmdNginxStatus     = "status"  // gitmap:cmd skip
	SubCmdNginxTest       = "test"    // gitmap:cmd skip
	SubCmdNginxReload     = "reload"  // gitmap:cmd skip
	SubCmdNginxInstall    = "install" // gitmap:cmd skip
	CmdServerCmd          = "server-cmd"
	CmdServerCmds         = "server-cmds"
	CmdServerCmdAlias     = "scmd"
	CmdPower              = "power"
	CmdPowerAlias         = "pw"
	CmdPowerAlias2        = "pwr"
	SubCmdPowerStatus     = "status"      // gitmap:cmd skip
	SubCmdPowerNever      = "never-sleep" // gitmap:cmd skip
	SubCmdPowerSet        = "set"         // gitmap:cmd skip
	SubCmdPowerReset      = "reset"       // gitmap:cmd skip
	SubCmdPowerHistory    = "history"     // gitmap:cmd skip
	CmdOS                 = "os"
	CmdFixLink            = "fix-link"
	CmdFixLinkAlias       = "fixlink"
	CmdFixLinkAlias2      = "fl"
	SubCmdFixLink         = "fix-link" // gitmap:cmd skip
	SubCmdFixLinkAlias    = "fixlink"  // gitmap:cmd skip
	SubCmdFixLinkAlias2   = "fl"       // gitmap:cmd skip
	SubCmdOSDisplay       = "display"  // gitmap:cmd skip
	SubCmdOSDisplayAlias  = "disp"     // gitmap:cmd skip
	SubCmdOSDisplayAlias2 = "screen"   // gitmap:cmd skip
	SubCmdOSStatus        = "status"   // gitmap:cmd skip
	SubCmdOSIP            = "ip"       // gitmap:cmd skip
	SubCmdOSZsh           = "zsh"      // gitmap:cmd skip
	SubCmdOSUser          = "user"     // gitmap:cmd skip
	SubCmdOSHelp          = "help"     // gitmap:cmd skip
	CmdUser               = "user"     // gitmap:cmd skip
	CmdZsh                = "zsh"      // gitmap:cmd skip
	// CmdReleasePull is sugar for `release` that first runs
	// `git pull --ff-only` in the current repo, then delegates to
	// the regular release pipeline. Hard-fails on non-fast-forward
	// so we never tag on top of a divergent tree.
	//
	// Renamed v5.6.0: canonical name is now `pull-release` (verb-first,
	// reads as "pull then release"). Short aliases are `relp` and `rlp`.
	// Collision on `pr` eliminated in Spec 129 in favor of the PR replay engine.
	CmdReleasePull       = "pull-release"
	CmdReleasePullAlias  = "pr"           // gitmap:cmd skip
	CmdReleasePullAlias2 = "release-pull" // legacy long form
	CmdReleasePullAlias3 = "relp"         // legacy alias
	CmdReleasePullAlias4 = "rlp"          // legacy alias
	// CmdMakePublic / CmdMakePrivate toggle the current repo's
	// visibility on GitHub or GitLab via `gh` / `glab`. No alias —
	// the long forms are the spec-mandated UX. Spec:
	// spec-authoring/23-visibility-change/01-spec.md.
	CmdMakePublic  = "make-public"
	CmdMakePrivate = "make-private"
	// CmdMakeAllPublic / CmdMakeAllPrivate flip visibility for every
	// repo under an owner/org that matches a comma-separated wildcard
	// pattern list (exact, prefix*, *contains*, prefix*suffix). Short
	// uppercase aliases CmdMAPUB / CmdMAPRI are spec-mandated muscle-
	// memory shortcuts. Interactive numbered list with per-index
	// exclusion; -Y/--yes bypasses both prompts. Persists a GitMapRun
	// row plus one GitMapRepoResult per repo. Spec:
	// 02-spec/01-app/116-bulk-visibility-mapub-mapri.md.
	CmdMakeAllPublic  = "make-all-public"
	CmdMakeAllPrivate = "make-all-private"
	CmdMAPUB          = "MAPUB"
	CmdMAPRI          = "MAPRI"
	// Except-latest variants — same semantics as the base commands
	// but exclude the highest `-vN` sibling from each base group, so
	// the latest version stays on its prior visibility while every
	// earlier version is flipped. Short uppercase aliases mirror the
	// MAPUB / MAPRI muscle-memory shortcuts.
	CmdMakeAllPublicExceptLatest  = "make-all-public-except-latest"
	CmdMakeAllPrivateExceptLatest = "make-all-private-except-latest"
	CmdMAPUBXL                    = "MAPUBXL"
	CmdMAPRIXL                    = "MAPRIXL"
	// CmdMakeLastPublic / CmdMakeLastPrivate flip exactly one repo:
	// the highest `-vN` sibling under `<base>` for the given owner.
	// When `<base>` itself doesn't end in `-vN` we resolve via the
	// OwnerRepoNameIndex; if `<base>` is already an exact repo name
	// we apply directly. Honors -Y / --yes. Spec follow-up to
	// 02-spec/01-app/116-bulk-visibility-mapub-mapri.md.
	CmdMakeLastPublic  = "make-last-public"
	CmdMakeLastPrivate = "make-last-private"
	CmdMLPUB           = "MLPUB"
	CmdMLPRI           = "MLPRI"
	// CmdVisibilityUndo reverses the most recent successful bulk
	// `make-all-public` / `make-all-private` (or prior undo) by
	// reading the persisted MakeAllVisibilityResult rows and
	// re-applying each repo's PrevVisibility. Logged as a new run
	// with CommandKind=VisibilityUndo. Spec:
	// 02-spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
	CmdVisibilityUndo      = "visibility-undo"
	CmdVisibilityUndoAlias = "vu"
	// CmdVisibilityRedo reverses the most recent VisibilityUndo run
	// (i.e. re-applies the visibility state that the undo reverted).
	// Logged as a new run with CommandKind=VisibilityRedo. Accepts
	// `--run <id>` to target a specific historical undo run.
	CmdVisibilityRedo      = "visibility-redo"
	CmdVisibilityRedoAlias = "vr"
	// CmdVisibilityHistory lists the most recent make-all-* /
	// VisibilityUndo / VisibilityRedo runs (newest first) so users
	// can pick a `--run <id>` for vu/vr. Accepts `--limit N`.
	CmdVisibilityHistory      = "visibility-history"
	CmdVisibilityHistoryAlias = "vish"
	// CmdHistoryPurge / CmdHistoryPin wrap `git filter-repo` in a
	// mirror-clone sandbox. Spec: 02-spec/04-generic-cli/16-history-rewrite.md.
	// The user's working repo is never rewritten in place.
	CmdHistoryPurge      = "history-purge"
	CmdHistoryPurgeAlias = "hp"
	CmdHistoryPin        = "history-pin"
	CmdHistoryPinAlias   = "hpin"
	// CmdCommitIn / CmdCommitInAlias — chronological multi-source
	// commit replay into a single destination repo. See
	// 02-spec/03-commit-in/. Both tokens resolve to the same handler;
	// short alias is `cin`.
	CmdCommitIn      = "commit-in"
	CmdCommitInAlias = "cin"
	// CmdVSCodePMSync walks every entry already present in the
	// alefragnani.project-manager projects.json file and re-runs the
	// gitmap auto-tag detector against each rootPath, then merges the
	// detected tags (always including the "gitmap" brand tag) into the
	// existing tags array. User-added tags are preserved (union, never
	// truncate). Foreign entries that point at deleted folders are
	// skipped — their tags are left untouched. Spec:
	// 02-spec/01-vscode-project-manager-sync/04-tag-resync.md.
	// Help: vscode-pm-sync.md.
	CmdVSCodePMSync      = "vscode-pm-sync"
	CmdVSCodePMSyncAlias = "vpm"
	// CmdPullReleaseCD is a multi-repo, one-shot pull-release runner.
	// Accepts a comma-separated list of `<name-or-url> <version>` pairs;
	// chdirs into each (cloning URLs first) and runs `gitmap pull-release
	// <version> -y` as an isolated subprocess. Spec:
	// 02-spec/01-app/112-pull-release-cd.md.
	CmdPullReleaseCD      = "pull-release-cd"
	CmdPullReleaseCDAlias = "prc"
	// CmdExportConfig exports application and developer tool configurations to JSON.
	CmdExportConfig      = "export-config"
	CmdExportConfigAlias = "config-export"
	// CmdImportConfig imports application and developer tool configurations from JSON.
	CmdImportConfig       = "import-config"
	CmdImportConfigAlias  = "config-import"
	CmdImportConfigTypo   = "improt-config" // gitmap:cmd skip
	CmdAsync              = "async"
	CmdAsyncAlias         = "asyn"
	CmdStorage            = "storage"
	CmdStorageAlias       = "stor"
	CmdAi                 = "ai"
	CmdAiAlias            = "scripts"
	CmdAutomation         = "automation"
	CmdAutomationAlias    = "auto"
	CmdAutomationAumAlias = "aum"
	CmdAutomationPyAlias  = "py-auto"
	CmdAgent              = "agent"
	CmdAgentAlias         = "ai-agents"
	CmdAgentAlias2        = "agents"
	CmdChangeUser         = "change-user"  // gitmap:cmd skip
	CmdChangeUserAlias    = "cu"           // gitmap:cmd skip
	CmdChangeUserAlias2   = "chuser"       // gitmap:cmd skip
	CmdChangeLogin        = "change-login" // gitmap:cmd skip
	CmdChangeLoginAlias   = "cl"           // gitmap:cmd skip
	CmdChangeLoginAlias2  = "chlogin"      // gitmap:cmd skip
	CmdAutobot            = "autobot"      // gitmap:cmd skip
	CmdAutobotAlias       = "bot"          // gitmap:cmd skip
)
