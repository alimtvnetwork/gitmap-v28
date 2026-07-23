package constants

// gitmap:cmd top-level
// CLI commands.
const (
	CmdScan               = "scan"
	CmdScanAlias          = "s"
	CmdClone              = "clone"
	CmdCloneAlias         = "c"
	CmdUpdate             = "update"
	CmdUpdateRunner       = "update-runner"  // gitmap:cmd skip
	CmdUpdateCleanup      = "update-cleanup" // gitmap:cmd skip
	CmdInstalledDir       = "installed-dir"  // gitmap:cmd skip
	CmdInstalledDirAlias  = "id"
	CmdVersion            = "version"
	CmdVersionAlias       = "v"
	CmdHelp               = "help"
	CmdDesktopSync        = "desktop-sync"
	CmdDesktopSyncAlias   = "ds"
	CmdGitHubDesktop      = "github-desktop"
	CmdGitHubDesktopAlias = "gd"
	CmdPull               = "pull"
	CmdPullAlias          = "p"
	CmdPullAll            = "pull-all"
	CmdPullAllAlias       = "pa"
	// CmdPush is the cwd-scoped `git push` wrapper added in v5.29.0.
	// Alias is `ph` (NOT `p` — collides with CmdPullAlias). Supports
	// shared `--ssh` / `--https` transport flags.
	CmdPush        = "push"
	CmdPushAlias   = "ph"
	CmdRescan      = "rescan"
	CmdRescanAlias = "rsc"
	// CmdRescanSubtree narrowly re-runs `gitmap scan` against the
	// absolutePath of an at-cap row from a previous scan output. Default
	// MaxDepth is bumped to RescanSubtreeDefaultMaxDepth so users get a
	// deeper view in one step; override with --max-depth.
	CmdRescanSubtree       = "rescan-subtree"
	CmdRescanSubtreeAlias  = "rss"
	CmdSetup               = "setup"
	CmdStatus              = "status"
	CmdStatusAlias         = "st"
	CmdExec                = "exec"
	CmdExecAlias           = "x"
	CmdRelease             = "release"
	CmdReleaseShort        = "r"
	CmdReleaseBranch       = "release-branch"
	CmdReleaseBranchAlias  = "rb"
	CmdReleasePending      = "release-pending"
	CmdReleasePendingAlias = "rp"
	CmdReleaseUndo         = "release-undo"
	CmdReleaseUndoAlias    = "ru"
	CmdChangelog           = "changelog"
	CmdChangelogAlias      = "cl"
	CmdChangelogMD         = "changelog.md" // gitmap:cmd skip
	CmdDoctor              = "doctor"
	CmdBackup              = "backup"
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
	CmdPR                = "pull-requests"
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
	CmdGroupCreate           = "create" // gitmap:cmd skip
	CmdGroupAdd              = "add"    // gitmap:cmd skip
	CmdGroupRemove           = "remove" // gitmap:cmd skip
	CmdGroupList             = "list"   // gitmap:cmd skip
	CmdGroupShow             = "show"   // gitmap:cmd skip
	CmdGroupDelete           = "delete" // gitmap:cmd skip
	CmdDBReset               = "db-reset"
	CmdReset                 = "reset"
	CmdListVersions          = "list-versions"
	CmdListVersionsAlias     = "lv"
	CmdRevert                = "revert"
	CmdRevertRunner          = "revert-runner" // gitmap:cmd skip
	CmdRm                    = "rm"
	CmdRmAlias               = "remove"
	CmdRmAlias2              = "del"
	CmdListReleases          = "list-releases"
	CmdListReleasesAlias     = "lr"
	CmdReleases              = "releases" // v3.20.0: alias of list-releases, intended for --all-repos batch view
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
	CmdReleaseSelfAlias2     = "rs"
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
	CmdSync                  = "sync"
	CmdSyncAlias             = "sy"
	// CmdCommons (v6.76.0) is a shortcut for `gitmap sync all` — adds
	// or dedupe-merges the curated .gitignore, .gitattributes,
	// .prettierignore, .prettierrc baselines + `git lfs install`.
	CmdCommons      = "commons"
	CmdCommonsAlias = "co"


	CmdSetSourceRepo         = "set-source-repo" // gitmap:cmd skip
	CmdSelfInstall           = "self-install"
	CmdSelfUninstall         = "self-uninstall"
	CmdSf                    = "sf"
	CmdProbe                 = "probe"
	CmdCode                  = "code"
	CmdVSCodePMPath          = "vscode-pm-path"
	CmdVSCodePMPathAlias     = "vpath"
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
	// spec/05-coding-guidelines/21-golden-fixture-regeneration.md.
	CmdRegoldens      = "regoldens"
	CmdRegoldensAlias = "rg"
	// CmdAuditLegacy scans the workspace for forbidden legacy strings
	// (e.g. gitmap-v27/v6/v7) and exits non-zero on any hit. Used as a // gitmap-legacy-ref-allow
	// regression guard for remixes / rename commits. Help: audit-legacy.md.
	CmdAuditLegacy       = "audit-legacy"
	CmdAuditLegacyAlias  = "aul"
	CmdAuditLegacyAlias2 = "audit"
	// CmdFixRepo rewrites prior `{base}-vN` tokens in tracked files to
	// the current version, mirroring fix-repo.ps1 / fix-repo.sh exit
	// codes and config schema. Spec: spec/04-generic-cli/27-fix-repo-command.md.
	CmdFixRepo      = "fix-repo"
	CmdFixRepoAlias = "fr"
	// CmdUndo restores the latest `gitmap fix-repo` backup snapshot
	// for the current repo + current version (v5.40.0+). Backups live
	// at `<repoRoot>/.gitmap/backup/<repo>/v<N>/fix-repo/<UTC-ts>/`.
	// Spec: spec/04-generic-cli/27-fix-repo-command.md §"Backup & undo".
	CmdUndo      = "undo"
	CmdUndoAlias = "ud"
	// CmdReleasePull is sugar for `release` that first runs
	// `git pull --ff-only` in the current repo, then delegates to
	// the regular release pipeline. Hard-fails on non-fast-forward
	// so we never tag on top of a divergent tree.
	//
	// Renamed v5.6.0: canonical name is now `pull-release` (verb-first,
	// reads as "pull then release"). Short alias is `pr`. The legacy
	// names `release-pull`, `relp`, and `rlp` remain wired as aliases
	// for backward compatibility and will not be removed.
	CmdReleasePull       = "pull-release"
	CmdReleasePullAlias  = "pr"
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
	// spec/01-app/116-bulk-visibility-mapub-mapri.md.
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
	// spec/01-app/116-bulk-visibility-mapub-mapri.md.
	CmdMakeLastPublic  = "make-last-public"
	CmdMakeLastPrivate = "make-last-private"
	CmdMLPUB           = "MLPUB"
	CmdMLPRI           = "MLPRI"
	// CmdVisibilityUndo reverses the most recent successful bulk
	// `make-all-public` / `make-all-private` (or prior undo) by
	// reading the persisted MakeAllVisibilityResult rows and
	// re-applying each repo's PrevVisibility. Logged as a new run
	// with CommandKind=VisibilityUndo. Spec:
	// spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
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
	// mirror-clone sandbox. Spec: spec/04-generic-cli/16-history-rewrite.md.
	// The user's working repo is never rewritten in place.
	CmdHistoryPurge      = "history-purge"
	CmdHistoryPurgeAlias = "hp"
	CmdHistoryPin        = "history-pin"
	CmdHistoryPinAlias   = "hpin"
	// CmdCommitIn / CmdCommitInAlias — chronological multi-source
	// commit replay into a single destination repo. See
	// spec/03-commit-in/. Both tokens resolve to the same handler;
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
	// spec/01-vscode-project-manager-sync/04-tag-resync.md.
	// Help: vscode-pm-sync.md.
	CmdVSCodePMSync      = "vscode-pm-sync"
	CmdVSCodePMSyncAlias = "vpm"
	// CmdPullReleaseCD is a multi-repo, one-shot pull-release runner.
	// Accepts a comma-separated list of `<name-or-url> <version>` pairs;
	// chdirs into each (cloning URLs first) and runs `gitmap pull-release
	// <version> -y` as an isolated subprocess. Spec:
	// spec/01-app/112-pull-release-cd.md.
	CmdPullReleaseCD      = "pull-release-cd"
	CmdPullReleaseCDAlias = "prc"
)

// vscode-pm-sync flag names + descriptions.
const (
	FlagCloneYes      = "yes"
	FlagCloneYesShort = "y"
	FlagDescCloneYes  = "Auto-accept new SSH host keys for SSH clone URLs " +
		"(sets StrictHostKeyChecking=accept-new; changed keys still fail)"

	FlagVSCodePMSyncDryRun     = "dry-run"
	FlagDescVSCodePMSyncDryRun = "Preview projects.json changes without writing"
	// FlagVSCodePMSyncProjectsJSON overrides the resolved
	// projects.json location for a single run. Useful when VS Code
	// is installed in a non-standard location (portable build,
	// remote-dev container, second user profile) or when scripting
	// against a backup copy. Path must be absolute. When omitted the
	// command falls back to the default discovery chain.
	FlagVSCodePMSyncProjectsJSON     = "projects-json"
	FlagDescVSCodePMSyncProjectsJSON = "Absolute path to projects.json (overrides VS Code user-data discovery)"
	// FlagVSCodePMSyncTag REPLACES the auto-detected tag set for
	// every entry with the user-supplied list. Repeatable AND
	// accepts a comma-separated value: `--tag a --tag b,c` produces
	// {a,b,c}. Differs from the global `--vscode-tag` (which only
	// ADDS to detection); use this when you want exact control over
	// what lands in projects.json. The "gitmap" brand tag is NOT
	// auto-prepended in this mode — pass it explicitly if you want
	// to keep it.
	FlagVSCodePMSyncTag     = "tag"
	FlagDescVSCodePMSyncTag = "Replace detected tags with this name (repeatable; accepts comma-list). Use instead of --vscode-tag for exact control"
	// FlagVSCodePMSyncMode picks the tag-merge strategy used when the
	// detected set is reconciled with the existing on-disk tags.
	// Values: union (default, additive), replace (detector wins),
	// intersection (set-AND + brand pin — only tags present in BOTH
	// sources survive, but the "gitmap" brand tag is always preserved
	// so re-syncs never silently strip our own marker).
	FlagVSCodePMSyncMode     = "mode"
	FlagDescVSCodePMSyncMode = "Tag merge strategy: union (default) | replace | intersection"
	// VSCodePMSyncMode* are the canonical CLI string values for the
	// --mode flag. Centralized so the parser, the helptext, and the
	// MergeMode.String() implementation can never drift.
	VSCodePMSyncModeUnion        = "union"
	VSCodePMSyncModeReplace      = "replace"
	VSCodePMSyncModeIntersection = "intersection"
)

// Audit-legacy defaults + flag names. Patterns are comma-separated regexes.
// Defaults catch the gitmap-v27/v6/v7 → v8 rename regression class. // gitmap-legacy-ref-allow
const (
	FlagAuditLegacyPatterns     = "patterns"
	FlagDescAuditLegacyPatterns = "Comma-separated regex patterns to flag (default: gitmap-v[567])"
	FlagAuditLegacyPath         = "path"
	FlagDescAuditLegacyPath     = "Root directory to scan (default: current working directory)"
	FlagAuditLegacyJSON         = "json"
	FlagDescAuditLegacyJSON     = "Emit machine-readable JSON instead of human text"
	FlagAuditLegacyReport       = "report"
	FlagDescAuditLegacyReport   = "Write a Markdown audit report to this path (use '' for default .gitmap/audit-legacy-report.md)"
	FlagAuditLegacyDiffs        = "diffs"
	FlagDescAuditLegacyDiffs    = "Also write per-file unified diffs (legacy → v8) into the report's diffs/ folder and link them"
	DefaultAuditLegacyPatterns  = "gitmap-v[567]\\b"
	DefaultAuditLegacyReport    = ".gitmap/audit-legacy-report.md"
	DefaultAuditLegacyDiffsDir  = "diffs"
	DefaultAuditLegacyReplace   = "gitmap-v27"
	MsgAuditLegacyClean         = "  ✓ audit-legacy: no forbidden legacy strings found in %s\n"
	MsgAuditLegacyHeader        = "  ✗ audit-legacy: %d match(es) across %d file(s) for pattern(s) %q\n"
	MsgAuditLegacyHit           = "    %s:%d: %s\n"
	MsgAuditLegacyReportWrote   = "  📄 audit-legacy: report written to %s\n"
	MsgAuditLegacyDiffsWrote    = "  📎 audit-legacy: %d per-file diff(s) written to %s\n"
	ErrAuditLegacyWalk          = "audit-legacy: walk failed for %s: %v\n"
	ErrAuditLegacyRegex         = "audit-legacy: invalid regex %q: %v\n"
	ErrAuditLegacyReportWrite   = "audit-legacy: failed to write report to %s: %v\n"
	ErrAuditLegacyDiffWrite     = "audit-legacy: failed to write diff for %s: %v\n"
)

// Regoldens flag names + descriptions. The two env vars driven by
// this command (GITMAP_UPDATE_GOLDEN, GITMAP_ALLOW_GOLDEN_UPDATE)
// are the canonical two-key safety gate from the goldenguard
// package; the values themselves live there to keep one source of
// truth (do NOT duplicate "1" here).
const (
	FlagRegoldensPattern         = "pattern"
	FlagDescRegoldensPattern     = "Test name pattern passed to `go test -run` (required)"
	FlagRegoldensPackage         = "package"
	FlagDescRegoldensPackage     = "Go package selector to test (default: ./...)"
	FlagRegoldensSkipVerify      = "skip-verify"
	FlagDescRegoldensSkipVerify  = "Skip the second verification pass that re-runs without the gate env vars"
	FlagRegoldensDryRun          = "dry-run"
	FlagDescRegoldensDryRun      = "Print the `go test` invocations that would run, then exit 0"
	FlagRegoldensDiff            = "diff"
	FlagDescRegoldensDiff        = "After pass 1, print a summary of touched testdata/ golden files. Modes: 'short' = one terse line per file (status + path + +/- counts); 'full' = adds rename details and aggregate totals. Usage: --diff=short or --diff=full"
	FlagRegoldensDeterminism     = "determinism"
	FlagDescRegoldensDeterminism = "Run a determinism pre-check pass BEFORE pass 1: invokes `go test` with the per-test trigger but WITHOUT the allow-update gate, so any writer using goldenguard.AllowUpdateAfterDeterminism is exercised but no fixtures are written. Fails fast (exit 1) on the first non-deterministic writer; otherwise pass 1 proceeds normally."
	RegoldensDiffModeShort       = "short"
	RegoldensDiffModeFull        = "full"
	RegoldensDefaultPackageGlob  = "./..."
	// RegoldensNonDetMarker is the substring goldenguard prints when
	// a writer fails AssertWriterDeterministic. The pre-check pass
	// scans `go test` output for this marker to distinguish a real
	// non-determinism failure from the EXPECTED gate-block failure.
	// Keep in sync with goldenguard/determinism.go nonDeterministicWriterMsgFmt.
	RegoldensNonDetMarker = "is non-deterministic"
)

// Regoldens user-facing messages.
const (
	MsgRegoldensPass1Header     = "▸ Pass 1: regenerating goldens (GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1)\n"
	MsgRegoldensPass2Header     = "▸ Pass 2: verifying determinism (gate env vars cleared)\n"
	MsgRegoldensSkipVerify      = "▸ Pass 2 skipped (--skip-verify)\n"
	MsgRegoldensSuccess         = "✓ Goldens regenerated and verified for pattern %q in package %q\n"
	MsgRegoldensSuccessNoVeri   = "✓ Goldens regenerated for pattern %q in package %q (verification skipped)\n"
	MsgRegoldensDryRun          = "▸ Dry run — would execute:\n  %s\n  %s\n"
	MsgRegoldensDiffHeader      = "▸ Golden diff summary (testdata/**/*golden* files touched by pass 1) [mode=%s]:\n"
	MsgRegoldensDiffNoChanges   = "  (no testdata/**/*golden* files changed)\n"
	MsgRegoldensDiffLineShort   = "  %s  %s\n"
	MsgRegoldensDiffLineFull    = "  %s  %s  (+%d / -%d)\n"
	MsgRegoldensDiffRenameFull  = "      ↳ renamed from %s\n"
	MsgRegoldensDiffTotals      = "  ─ %d file(s) changed: %d added, %d modified, %d renamed, %d deleted (+%d / -%d total)\n"
	MsgRegoldensDiffSkipped     = "▸ Diff summary skipped: not a git working tree (or `git` not on PATH)\n"
	MsgRegoldensPass2NotRun     = "▣ Pass 2: did NOT run — pass 1 failed (exit %d); fixtures left in partial state, determinism unverified\n"
	MsgRegoldensPass2NotRunSkip = "▣ Pass 2: did NOT run — --skip-verify was set; determinism unverified\n"
	MsgRegoldensPass2Ran        = "▣ Pass 2: ran and PASSED — writer is deterministic; goldens are safe to commit\n"
	MsgRegoldensPrecheckHeader  = "▸ Pre-check: determinism (trigger ON, allow-update gate OFF — no fixture writes)\n"
	MsgRegoldensPrecheckPass    = "✓ Pre-check passed — no non-deterministic writers detected; proceeding to pass 1\n"
	ErrRegoldensMissingPat      = "regoldens: --pattern is required (e.g. --pattern TestCloneFromReportJSON_Golden)"
	ErrRegoldensPass1Failed     = "regoldens: pass 1 (regenerate) failed with exit code %d — fixtures may be partially written; inspect git status before committing"
	ErrRegoldensPass2Failed     = "regoldens: pass 2 (verify) failed with exit code %d — the writer is non-deterministic; fix the writer, do not re-run with --skip-verify"
	ErrRegoldensDiffMode        = "regoldens: --diff=%q is not valid; use --diff=short or --diff=full"
	ErrRegoldensPrecheckFailed  = "regoldens: determinism pre-check FAILED — at least one golden writer produced different bytes across runs; fix the writer (likely culprits: map iteration, time.Now, randomness, locale-dependent formatting). Pass 1 was NOT run; no fixtures were touched."
)

// Usage header — colorful banner shown above every help screen.
// Width is intentionally unboxed because ANSI escapes + dynamic version
// width make fixed-width box alignment unreliable across terminals.
const UsageHeaderFmt = "\n  " + ColorMagenta + "🗺  " + ColorWhite + "gitmap" + ColorReset +
	"  " + ColorCyan + "v%s" + ColorReset +
	"  " + ColorDim + "— Git repo discovery, cloning & release toolkit" + ColorReset + "\n" +
	"  " + ColorMagenta + "────────────────────────────────────────────────────────────" + ColorReset + "\n\n"

const (
	HelpUsage            = "Usage: gitmap <command> [flags]"
	HelpCommands         = "Commands:"
	HelpScan             = "  scan (s) [dir]      Scan directory for Git repos"
	HelpClone            = "  clone (c) <source|json|csv|text>  Re-clone from file (shorthands auto-resolve)"
	HelpUpdate           = "  update              Self-update from source repo"
	HelpUpdateCleanup    = "  update-cleanup      Remove leftover update temp files and .old backups"
	HelpInstalledDir     = "  installed-dir (id)  Show the active installed binary path"
	HelpVersion          = "  version (v)         Show version number"
	HelpDesktopSync      = "  desktop-sync (ds)   Sync repos to GitHub Desktop from output"
	HelpGitHubDesktop    = "  github-desktop (gd) Register current repo with GitHub Desktop (no scan needed)"
	HelpPull             = "  pull (p) <name>     Pull a specific repo by its name"
	HelpRescan           = "  rescan (rsc)        Re-run last scan with cached flags"
	HelpSetup            = "  setup               Configure Git diff/merge tool, aliases & core settings"
	HelpStatus           = "  status (st)         Show dirty/clean, ahead/behind, stash for all repos"
	HelpExec             = "  exec (x) <args...>  Run any git command across all repos"
	HelpRelease          = "  release (r) [ver]   Create release branch, tag, and push"
	HelpReleasePull      = "  pull-release (pr) [ver]    git pull (--ff-only|--rebase|--merge), then release the current repo"
	HelpReleaseBr        = "  release-branch (rb) Complete release from existing release branch"
	HelpReleasePend      = "  release-pending (rp) Release all pending branches without tags"
	HelpFixRepo          = "  fix-repo (fr)       Rewrite prior {base}-vN tokens (-2/-3/-5/--all/--dry-run)"
	HelpChangelog        = "  changelog (cl) [ver] Show concise release notes (use --open, --source)"
	HelpDoctor           = "  doctor [--fix-path] Diagnose PATH, deploy, and version issues"
	HelpLatestBr         = "  latest-branch (lb)  Find most recently updated remote branch"
	HelpList             = "  list (ls)           Show all tracked repos with slugs"
	HelpGroup            = "  group (g) <sub>     Manage repo groups / activate group for batch ops"
	HelpMultiGroup       = "  multi-group (mg)    Select multiple groups for batch operations"
	HelpSf               = "  sf <add|list|rm>    Manage scan folders (roots that gitmap scan tracks)"
	HelpRm               = "  rm (remove, del) <name-or-path>...  Untrack repos from the database only"
	HelpDBReset          = "  db-reset --confirm  Clear all tracked repos and groups from the database"
	HelpCompletion       = "  completion (cmp)    Generate shell tab-completion scripts"
	HelpClearReleaseJSON = "  clear-release-json (crj)  Remove a .gitmap/release/vX.Y.Z.json file"
	HelpDocs             = "  docs (d)            Open documentation website in browser"
	HelpHelpDash         = "  help-dashboard (hd) Serve the docs site locally in your browser"
	HelpCloneNext        = "  clone-next (cn)     Clone next versioned iteration of current repo"
	HelpReleaseSelf      = "  release-self (rs)   Release gitmap itself from any directory"
	HelpHelp             = "  help                Show this help message"
	HelpListVersions     = "  list-versions (lv)  Show all release tags, highest first (--limit N, --json, --source)"
	HelpListReleases     = "  list-releases (lr)  Show releases from .gitmap/release/ files or database (--limit N, --json, --source)"
	HelpRevert           = "  revert <version>    Revert to a specific release version"
	MsgHelpLFSCommon     = "  lfs-common (lfsc)   Track common binary file types with Git LFS in current repo"
)

// Help section headers and flag-line strings (HelpScanFlags, HelpCloneFlags,
// HelpReleaseFlags, ...) live in constants_helpsections.go.

// Flag descriptions.
const (
	FlagDescConfig     = "Path to config file"
	FlagDescMode       = "Clone URL style: https or ssh"
	FlagDescOutput     = "Output format: terminal, csv, json"
	FlagDescOutFile    = "Exact output file path"
	FlagDescOutputPath = "Output directory for CSV/JSON"
	// FlagScanManifest is a unified alias for `--output-path` on the
	// scan command, mirroring the same flag on `gitmap reclone` so
	// the scan→reclone round-trip uses one vocabulary. Resolution:
	// `--manifest` is honored only when `--output-path` is empty.
	FlagScanManifest     = "manifest"
	FlagDescScanManifest = "Alias for --output-path: directory where the scan " +
		"manifest (JSON + CSV) is written. Used by `gitmap reclone --manifest`."
	FlagDescTargetDir   = "Base directory for cloned repos"
	FlagDescSafePull    = "If repo exists, run safe git pull with retries and unlock diagnostics"
	FlagDescGHDesktop   = "Add discovered repos to GitHub Desktop"
	FlagDescOpen        = "Open output folder after scan completes"
	FlagDescQuiet       = "Suppress terminal clone help section"
	FlagDescVerbose     = "Write detailed stdout/stderr debug log to a timestamped file"
	FlagScanWorkers     = "workers"
	FlagDescScanWorkers = "Worker pool size for scan (0 = auto, capped at 16)"
	DefaultScanWorkers  = 0
	// FlagScanWorkersConcurrencyAlias is the deprecated long-form
	// spelling kept for muscle-memory compatibility. Mirrors the
	// --probe-workers / --probe-concurrency precedent in
	// constants_probe.go: --workers stays canonical; --concurrency
	// is honored but emits a one-line stderr deprecation notice.
	// When BOTH are passed, --workers wins (explicit canonical
	// flag overrides the deprecated alias) without a warning.
	FlagScanWorkersConcurrencyAlias     = "concurrency"
	FlagDescScanWorkersConcurrencyAlias = "(deprecated, use --workers) Worker pool size for scan (0 = auto, capped at 16)"
	// MsgScanWorkersConcurrencyAlias is printed to os.Stderr when
	// the deprecated --concurrency flag is used. Single-line,
	// prefixed with two-space + middle-dot to match the rest of
	// gitmap's stderr advisory format (see MsgScanProbeConcurrencyAlias).
	MsgScanWorkersConcurrencyAlias = "  · --concurrency is deprecated; use --workers instead\n"
	// FlagScanRelativeRoot lets the user pin the base path used to compute
	// each repo's RelativePath in the output (CSV/JSON/text/structure/
	// clone scripts). Without it, RelativePath is derived from the scan
	// dir, which means running `gitmap scan .` from different cwds
	// produces different paths for the same repos. With --relative-root,
	// every output row is computed against the supplied (absolute or
	// relative) directory so artifacts stay byte-stable across cwds.
	FlagScanRelativeRoot     = "relative-root"
	FlagDescScanRelativeRoot = "Pin the base directory used for output RelativePath (absolute or relative; must contain every scanned repo)"
	// FlagScanMaxDepth caps how many directory levels below the scan root
	// the walker may descend. Default 0 → resolves to scanner.DefaultMaxDepth
	// (4) inside the scanner. Negative → unbounded (legacy behavior).
	// Honored regardless of whether a `.git` was found on the path —
	// repos already stop their own subtree as before.
	FlagScanMaxDepth     = "max-depth"
	FlagDescScanMaxDepth = "Max directory levels to descend below scan root (0 = default 4, negative = unlimited)"
	DefaultScanMaxDepth  = 0
	// RescanSubtreeDefaultMaxDepth is the cap `gitmap rescan-subtree`
	// applies when the user does not pass an explicit --max-depth. It is
	// deliberately deeper than DefaultScanMaxDepth (which resolves to 4
	// inside the scanner) so the typical at-cap-rescan workflow finds
	// the previously hidden nested repos in one shot.
	RescanSubtreeDefaultMaxDepth = 8
	// FlagScanDefaultBranch overrides the fallback branch name written
	// to ScanRecord.Branch when none of the live detection steps in
	// gitutil.DetectBranchWithDefault returned a usable name. Without
	// the flag, the fallback is the package-default constants.DefaultBranch
	// ("main"). Useful for catalogs that target legacy infra still on
	// "master", or for forcing a project-specific convention without
	// touching the binary's compiled-in default.
	FlagScanDefaultBranch     = "default-branch"
	FlagDescScanDefaultBranch = "Fallback branch name when HEAD/remote-tracking detection finds nothing (default: main)"
	// FlagScanReportErrors enables a JSON failure report at command
	// exit for `gitmap scan` and `gitmap cn --all/--csv`. Bare boolean
	// — output path is fixed at `<binaryDir>/.gitmap/reports/errors-
	// <unixts>.json`. Honored by both `gitmap scan` (captures scanner
	// ReadDir + background ls-remote / shallow-clone failures) and
	// `gitmap cn --all/--csv` (captures per-repo clone failures).
	// Clean runs leave NO file on disk.
	//
	// NOTE: distinct from the existing `--report-errors` (with leading
	// dashes baked in) on the `gitmap update` command — that flag
	// takes a value (`json`) and writes a JSONL handoff trace. Naming
	// this one `--errors-report` keeps both UX surfaces intact.
	FlagScanReportErrors     = "errors-report"
	FlagDescScanReportErrors = "Write per-repo failures to .gitmap/reports/errors-<unixts>.json (only emitted when failures occur)"
	FlagDescSetupConfig      = "Path to git-setup.json config file"
	FlagDescDryRun           = "Preview changes without applying them"
	FlagDescAssets           = "Directory or file to attach to the release"
	FlagDescCommit           = "Create release from a specific commit"
	FlagDescRelBranch        = "Create release from latest commit of a branch"
	FlagDescBump             = "Auto-increment version: major, minor, or patch"
	FlagDescDraft            = "Create an unpublished draft release"
	FlagDescLatest           = "Show only the latest changelog entry"
	FlagDescLimit            = "Number of changelog versions to show"
	FlagDescOpenChangelog    = "Open CHANGELOG.md with the default system app"
	FlagDescLBRemote         = "Remote to filter branches against (default: origin)"
	FlagDescLBAllRemotes     = "Include branches from all remotes"
	FlagDescLBContains       = "Fall back to --contains if --points-at returns empty"
	FlagDescLBTop            = "Show top N most recently updated branches"
	FlagDescLBJSON           = "Output structured JSON instead of plain text (shorthand for --format json)"
	FlagDescLBFormat         = "Output format: terminal, json, csv (default: terminal)"
	FlagDescLBNoFetch        = "Skip git fetch (use existing remote refs)"
	FlagDescLBSort           = "Sort order: date (default, descending) or name (alphabetical)"
	FlagDescLBFilter         = "Filter branches by glob or substring pattern"
	// FlagDescLBSwitch / FlagDescLBSwitchShort document `gitmap lb
	// --switch` / `gitmap lb -s`. Both register the SAME effect: after
	// printing the latest-branch report, run `git checkout` against the
	// resolved branch name. Useful as a single-keystroke "jump to whatever
	// is freshest" workflow when reviewing PR work across many branches.
	FlagDescLBSwitch      = "Checkout the resolved latest branch after printing the report"
	FlagDescLBSwitchShort = "Short alias for --switch"
	FlagDescGroup         = "Filter by group name"
	FlagDescAll           = "Run against all tracked repos from database"
	FlagDescListVerbose   = "Show full paths and URLs"
	FlagDescGroupDesc     = "Optional group description"
	FlagDescGroupColor    = "Terminal color for group display"
	FlagDescConfirm       = "Confirm destructive operation"
	FlagDescSource        = "Filter by source: release or import"
)

// FlagDebugPaths exposes the canonicalize-path tracer used by the
// VS Code Project Manager sync helper. When set on `gitmap clone`,
// the flag exports GITMAP_DEBUG_PATHS=1 for the current process so
// every clone variant that routes through canonicalizePMPath emits
// a one-line stderr trace of (raw → cleaned → resolved) rootPaths.
// Diagnostic-only — Windows 8.3 short-name and symlink resolution is
// the primary debugging target. Off by default.
const (
	FlagDebugPaths     = "debug-paths"
	FlagDescDebugPaths = "trace VS Code Project Manager sync rootPath canonicalization (raw → cleaned → resolved) to stderr; useful for diagnosing Windows 8.3 short-name and symlink dedup"
)

// Release flag names — surfaced as Cmd*-style constants so callers
// (e.g. the right-click context menu in installctxentries.go and the
// release parser in release.go) reference the same single source of
// truth rather than literal "--bump" / "bump" strings.
const (
	FlagBump      = "bump"
	FlagBumpDash  = "--bump"
	FlagDryRunDsh = "--dry-run"
)
