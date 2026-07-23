package constants

// Notes.
const (
	NoteNoRemote    = "no remote configured"
	UnknownRepoName = "unknown"
)

// GitHub Desktop.
const (
	GitHubDesktopBin     = "github"
	OSWindows            = "windows"
	MsgDesktopNotFound   = "GitHub Desktop CLI not found — skipping."
	MsgDesktopAdded      = "  ✓ Added to GitHub Desktop: %s\n"
	MsgDesktopFailed     = "  ✗ Failed to add %s: %v\n"
	MsgDesktopSummary    = "GitHub Desktop: %d added, %d failed\n"
	MsgGHDesktopRegister = "  Registering with GitHub Desktop: %s\n"
	MsgGHDesktopDone     = "  ✓ Registered with GitHub Desktop: %s\n"
	ErrGHDesktopCwd      = "  ✗ Could not determine current directory: %v\n"
	ErrGHDesktopNotRepo  = "  ✗ Not a git repository: %s\n     (run `gitmap gd <path>` to register a different folder)\n"
	ErrGHDesktopInvoke   = "  ✗ GitHub Desktop CLI failed: %v\n%s\n"
)

// Latest-branch display messages.
const (
	MsgLatestBranchFetching     = "  Fetching remotes..."
	MsgLatestBranchFetchWarning = "  Warning: fetch failed: %v\n"
	LBUnknownBranch             = "<unknown>"
	// `lb --switch` / `-s` and `branch default` post-report messages.
	MsgLatestBranchSwitching      = "  ▶ Switching to %s...\n"
	ErrLatestBranchSwitchNoTarget = "  ✗ No usable branch name to switch to (try --contains-fallback)"
	ErrLatestBranchSwitchFailed   = "  ✗ Failed to checkout %s: %v\n"
	MsgBranchDefaultSwitching     = "  ▶ Switching to default branch %s...\n"
	ErrBranchDefaultFailed        = "  ✗ Failed to checkout default branch %s: %v\n"
	ErrBranchUnknownSubcommand    = "  ✗ Unknown branch subcommand: %q (expected: default | def)\n"
	ErrBranchMissingSubcommand    = "  ✗ Missing subcommand for `gitmap branch` (expected: default | def)\n"
	ErrBranchNotRepo              = "  ✗ Not inside a git work tree\n"
)

// Generic error formats.
const (
	ErrGenericFmt = "Error: %v\n"
	ErrBareFmt    = "%v\n"
)

// OS platform constants.
const OSDarwin = "darwin"

// gitmap:cmd top-level
// OS file-explorer commands.
const (
	CmdExplorer = "explorer" // gitmap:cmd skip
	// CmdOpen lives in constants_cli.go (canonical CLI ID location).
	CmdXdgOpen      = "xdg-open" // gitmap:cmd skip
	CmdWindowsShell = "cmd"      // gitmap:cmd skip
	CmdArgSlashC    = "/c"       // gitmap:cmd skip
	CmdArgStart     = "start"    // gitmap:cmd skip
	CmdArgEmpty     = ""         // gitmap:cmd skip
)

// Desktop sync error messages.
const (
	ErrDesktopReadFailed  = "Error reading %s: %v\n"
	ErrDesktopParseFailed = "Error parsing JSON from %s: %v\n"
	ErrNoAbsPath          = "no absolute path"
)

// Command dispatch errors.
const (
	ErrUnknownCommand        = "Unknown command: %s\n"
	ErrUnknownCommandURLHint = "Unknown command: %s\n" +
		"\n" +
		"  This looks like a git URL. Newer gitmap versions auto-redirect\n" +
		"  bare-URL invocations to `gitmap clone`. Your installed binary\n" +
		"  appears to predate that shortcut (added in v3.81.0).\n" +
		"\n" +
		"  Fix one of two ways:\n" +
		"    1. Use the explicit form right now:\n" +
		"         gitmap clone %[1]s\n" +
		"    2. Update gitmap so the shortcut is built in:\n" +
		"         gitmap update\n" +
		"       (then re-open your terminal so PATH picks up the new binary)\n\n"
	ErrUnknownGroupSub = "Unknown group subcommand: %s\n"
)

// Docs command.
const (
	DocsURL       = "https://gitmap.dev/docs"
	MsgDocsOpened = "  ✓ Opened %s\n"
	ErrDocsOpen   = "  ✗ Failed to open browser: %v\n"
)

// Version display.
const MsgVersionFmt = "gitmap v%s\n"

// CLI messages.
const (
	MsgFoundRepos         = "Found %d repositories.\n"
	MsgCSVWritten         = "  📊 CSV          %s\n"
	MsgJSONWritten        = "  🧬 JSON         %s\n"
	MsgTextWritten        = "  📝 Text list    %s\n"
	MsgStructureWritten   = "  🌳 Structure    %s\n"
	MsgCloneScript        = "  🪄 Clone PS1    %s\n"
	MsgDirectClone        = "  ⚡ HTTPS PS1    %s\n"
	MsgDirectCloneSSH     = "  🔐 SSH PS1      %s\n"
	MsgDesktopScript      = "  🖥️  Desktop PS1  %s\n"
	MsgCloneComplete      = "\nClone complete: %d succeeded, %d failed\n"
	MsgAutoSafePull       = "Existing repos detected — safe-pull enabled automatically.\n"
	MsgOpenedFolder       = "Opened output folder: %s\n"
	MsgVerboseLogFile     = "Verbose log: %s\n"
	MsgDesktopSyncStart   = "\n  Syncing repos to GitHub Desktop from %s...\n"
	MsgDesktopSyncSkipped = "  ⊘ Skipped (already exists): %s\n"
	MsgDesktopSyncAdded   = "  ✓ Added to GitHub Desktop: %s\n"
	MsgDesktopSyncFailed  = "  ✗ Failed: %s — %v\n"
	MsgDesktopSyncDone    = "\n  GitHub Desktop sync: %d added, %d skipped, %d failed\n"
	MsgNoOutputDir        = "Error: .gitmap/output/ not found in current directory.\nRun 'gitmap scan' first to generate output files."
	MsgNoJSONFile         = "Error: %s not found.\nRun 'gitmap scan' first to generate the JSON output."
	MsgFailedClones       = "\nFailed clones:"
	MsgFailedEntry        = "  - %s (%s): %s\n"
	MsgPullStarting       = "\n  Pulling %s (%s)...\n"
	MsgPullSuccess        = "  ✓ %s is up to date.\n"
	MsgPullFailed         = "  ✗ Pull failed for %s: %s\n"
	MsgPullAvailable      = "\nAvailable repos:"
	MsgPullListEntry      = "  - %s\n"
	WarnVerboseLogFailed  = "Warning: could not create verbose log: %v\n"
	MsgRescanReplay       = "\n  Rescanning with cached flags (dir: %s)...\n"
	MsgScanCacheSaved     = "  💾 Cache        %s\n"
	MsgDBUpsertDone       = "  ✅ %d repositories upserted into database\n"
	MsgDBUpsertFailed     = "Warning: database upsert failed: %v\n"
	// Section headers for the post-scan summary.
	// MsgSectionArtifacts takes the common output base directory so it can be
	// printed once instead of repeated on every file line.
	MsgSectionArtifacts = "\n📦 Output Artifacts\n" + MsgSectionRule + "\n📂 Base: %s\n\n"
	MsgSectionDatabase  = "\n🗄️  Database\n" + MsgSectionRule + "\n"
	MsgSectionProjects  = "\n🔍 Project Detection\n" + MsgSectionRule + "\n"
	MsgSectionDone      = "\n🎉 Scan complete.\n"
	MsgSectionRule      = "────────────────────────────────────────────"
	MsgScanFolderTagged = "  🏷️  Tagged %d repo(s) with scan folder #%d\n"
	MsgUpdateStarting   = "\n  Updating gitmap from source repo...\n"
	MsgUpdateRepoPath   = "  → Repo path: %s\n"
	MsgUpdateVersion    = "\n  ✓ Updated to gitmap v%s\n"
)

// List and group messages.
const (
	MsgListHeader       = "SLUG                 REPO NAME"
	MsgListSeparator    = "──────────────────────────────────────────"
	MsgListRowFmt       = "%-20s %s\n"
	MsgListVerboseFmt   = "%-20s %-20s %s\n"
	MsgListEmpty        = "No repos tracked. Run 'gitmap scan' first."
	MsgListDBPath       = "  → Database: %s\n"
	MsgGroupCreated     = "Group created: %s\n"
	MsgGroupDeleted     = "Group deleted: %s\n"
	MsgGroupAdded       = "Added %s to group %s\n"
	MsgGroupRemoved     = "Removed %s from group %s\n"
	MsgGroupHeader      = "GROUP           REPOS   DESCRIPTION"
	MsgGroupRowFmt      = "%-15s %-7d %s\n"
	MsgGroupShowHeader  = "Group: %s (%d repos)\n"
	MsgGroupShowRowFmt  = "  %-16s %s\n"
	MsgGroupEmpty       = "No groups defined. Use 'gitmap group create <name>' to create one."
	MsgGroupActivated   = "Active group set: %s\n"
	MsgGroupActiveShow  = "Active group: %s\n"
	MsgGroupNoActive    = "No active group. Use 'gitmap g <name>' to set one."
	MsgGroupCleared     = "Active group cleared."
	ErrGroupNameReq     = "Error: group name is required"
	ErrGroupUsage       = "Usage: gitmap group <create|add|remove|list|show|delete> [args]"
	ErrGroupSlugReq     = "Error: at least one slug is required"
	ErrListDBFailed     = "Error: could not open database: %v\nRun 'gitmap scan' first.\n"
	ErrNoDatabase       = "No database found. Run 'gitmap scan' first."
	MsgDBResetDone      = "Database reset: all repos and groups cleared.\n"
	ErrDBResetFailed    = "Error: database reset failed: %v\n"
	ErrDBResetNoConfirm = "Error: this will delete all tracked repos and groups.\nRun with --confirm to proceed: gitmap db-reset --confirm"
	MsgResetFileRemoved = "Removed database file: %s\n"
	MsgResetReseeded    = "Reseeded %s\n"
	MsgResetDone        = "Reset complete: database file deleted, schema rebuilt, seeds reapplied.\n"
	ErrResetNoConfirm   = "Error: this will permanently delete the database file and rebuild it from scratch.\nRun with --confirm to proceed: gitmap reset --confirm"
	ErrResetRemoveFile  = "Error: could not delete database file %s: %v\n"
	ErrResetReinit      = "Error: could not reinitialize database: %v\n"
)

// Latest-branch error messages.
const (
	ErrLatestBranchNotRepo   = "Error: not inside a Git repository."
	ErrLatestBranchNoRefs    = "Error: no remote-tracking branches found for remote '%s'.\n"
	ErrLatestBranchNoRefsAll = "Error: no remote-tracking branches found on any remote."
	ErrLatestBranchNoCommits = "Error: could not read commit info for remote branches."
	ErrLatestBranchNoMatch   = "Error: no branches matching filter '%s'.\n"
)

// CLI error messages.
const (
	ErrSourceRequired    = "Error: source file or URL is required"
	ErrCloneUsage        = "Usage: gitmap clone <url|source|json|csv|text> [folder] [--target-dir <dir>] [--safe-pull]"
	ErrShorthandNotFound = "Error: %s not found.\nRun 'gitmap scan' first to generate output files.\n"
	ErrConfigLoad        = "Error: failed to load config from %s: %v (operation: read)\n"
	ErrScanFailed        = "Error: scan failed on directory %s: %v (operation: resolve)\n"
	ErrScanDirNotFound   = "Error: scan target %s does not exist (resolved to %s)\n"
	ErrScanDirNotDir     = "Error: scan target %s is not a directory (resolved to %s)\n"
	// ErrScanRelativeRootInvalid is emitted when --relative-root cannot be
	// resolved to an existing directory or is not a parent of the scan dir.
	// Operation tag mirrors ErrScanFailed style for log-grep consistency.
	ErrScanRelativeRootInvalid = "Error: --relative-root %s is invalid: %v (operation: resolve)\n"
	// ErrScanRelativeRootNotAncestor fires when the user-supplied root
	// does not contain a discovered repo, which would yield a "../"
	// path in the output and break clone-script reproducibility.
	ErrScanRelativeRootNotAncestor = "Error: --relative-root %s is not an ancestor of repo %s (would produce %q)\n"
	MsgScanResolvedDir             = "  ↳ Resolved %q → %s\n"
	// MsgScanRelativeRoot announces the pinned RelativePath base when it
	// differs from the scan dir, so the user sees that output paths
	// will be computed against this directory rather than [dir].
	MsgScanRelativeRoot = "  ↳ RelativePath root pinned to %s\n"
	ErrCloneFailed      = "Error: clone failed for source file %s: %v (operation: read)\n"
	ErrOutputFailed     = "Error: output generation failed: %v\n"
	ErrCreateDir        = "Error: cannot create directory at %s: %v (operation: mkdir)\n"
	ErrCreateFile       = "Error: cannot create file at %s: %v (operation: write)\n"
	ErrNoRepoPath       = `
  ✗ Source repository path not found.

  This binary was installed without a linked source repo, so 'update'
  cannot locate the code to pull and rebuild.

  How to fix:

    Option 1 — Re-install via the one-liner (recommended):
      irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/install.ps1 | iex

    Option 2 — Clone the repo and build from source:
      git clone https://github.com/alimtvnetwork/gitmap-v27.git C:\gitmap-src
      cd C:\gitmap-src
      .\run.ps1

    Option 3 — Download the latest release manually:
      https://github.com/alimtvnetwork/gitmap-v27/releases/latest

  After building from source, 'gitmap update' will work automatically.
`
	ErrUpdateFailed             = "Update error: %v\n"
	ErrPullSlugRequired         = "Error: repo name is required"
	ErrPullUsage                = "Usage: gitmap pull <repo-name> [--verbose]"
	ErrPullLoadFailed           = "Error: could not load gitmap.json at %s: %v (operation: read)\n"
	ErrPullNotFound             = "Error: no repo found matching '%s'\n"
	ErrPullNotRepo              = "Error: %s is not a git repository\n"
	ErrRescanNoCache            = "Error: no previous scan found. Run 'gitmap scan' first.\n%v\n"
	ErrSetupLoadFailed          = "Error: could not load git-setup.json at %s: %v (operation: read)\n"
	ErrStatusLoadFailed         = "Error: could not load gitmap.json at %s for status: %v (operation: read)\nRun 'gitmap scan' first.\n"
	MsgStatusNoData             = "No tracked repos found.\nRun 'gitmap scan' in a directory containing your git repos first,\nor pass --all to query the database directly.\n"
	ErrExecUsage                = "Usage: gitmap exec <git-args...>\nExample: gitmap exec fetch --prune"
	ErrExecLoadFailed           = "Error: could not load gitmap.json at %s: %v (operation: read)\nRun 'gitmap scan' first.\n"
	ErrReleaseVersionRequired   = "Error: version is required.\nProvide a version argument, use --bump, or create a version.json file."
	ErrReleaseUsage             = "Usage: gitmap release [version] [--assets <path>] [--commit <sha>] [--branch <name>] [--bump major|minor|patch] [--draft] [--dry-run]"
	ErrReleaseBranchUsage       = "Usage: gitmap release-branch <release/vX.Y.Z> [--assets <path>] [--draft]"
	ErrReleaseAlreadyExists     = "Error: version %s is already released. See .gitmap/release/%s.json for details.\n"
	ErrReleaseTagExists         = "Error: tag %s already exists.\n"
	ErrReleaseBranchNotFound    = "Error: branch %s does not exist.\n"
	ErrReleaseCommitNotFound    = "Error: commit %s not found.\n"
	ErrReleaseInvalidVersion    = "Error: '%s' is not a valid version.\n"
	ErrReleaseBumpNoLatest      = "Error: no previous release found. Create an initial release before using --bump.\n"
	ErrReleaseBumpConflict      = "Error: --bump cannot be used with an explicit version argument.\n"
	ErrReleaseCommitBranch      = "Error: --commit and --branch are mutually exclusive.\n"
	ErrReleasePushFailed        = "Error: failed to push to remote: %v\n"
	ErrReleaseVersionLoad       = "Error: could not read version.json at %s: %v (operation: read)\n"
	ErrReleaseMetaWrite         = "Error: could not write release metadata at %s: %v (operation: write)\n"
	ErrChangelogRead            = "Error: could not read CHANGELOG.md at %s: %v (operation: read)\n"
	ErrChangelogVersionNotFound = "Error: version %s not found in CHANGELOG.md\n"
	ErrChangelogOpen            = "Error: could not open CHANGELOG.md at %s: %v (operation: open)\n"
)

// List-versions error messages.
const (
	ErrListVersionsNoTags = "Error: no version tags found. Create a release first."
)

// List-releases messages.
const (
	MsgListReleasesEmpty     = "No releases found."
	MsgListReleasesHeader    = "Releases (%d found)\n"
	MsgListReleasesSeparator = "────────────────────────────────────────────────────────────────────────"
	MsgListReleasesColumns   = "  VERSION    TAG          BRANCH              DRAFT  LATEST  SOURCE   DATE"
	MsgListReleasesRowFmt    = "  %-10s %-12s %-19s %-6s %-7s %-8s %s\n"
	ErrListReleasesFailed    = "Error: failed to load releases: %v\n"
	MsgYes                   = "yes"
	MsgNo                    = "no"
)

// List-releases --all-repos messages (v3.20.0). Wider table — adds REPO column.
const (
	MsgListReleasesAllReposEmpty     = "No releases recorded across any repo. Run `gitmap release` in any repo first."
	MsgListReleasesAllReposHeader    = "Releases across all repos (%d found)\n"
	MsgListReleasesAllReposSeparator = "────────────────────────────────────────────────────────────────────────────────────"
	MsgListReleasesAllReposColumns   = "  REPO                 VERSION    TAG          BRANCH              LATEST  SOURCE   DATE"
	MsgListReleasesAllReposRowFmt    = "  %-20s %-10s %-12s %-19s %-7s %-8s %s\n"
)

// Release import messages.
const (
	MsgReleasesImported   = "Releases imported: %d from .gitmap/release/\n"
	WarnReleaseImportSkip = "Warning: skipping %s: %v\n"
	ReleaseGlob           = "v*.json"
)

// Pending metadata discovery messages.
const (
	MsgPendingMetaFound     = "  → Found %d unreleased version(s) from .gitmap/release/ metadata\n"
	MsgPendingMetaRelease   = "  → Creating release from metadata: %s (commit: %s)\n"
	WarnPendingMetaNoCommit = "  ⚠ Skipping %s: commit %s not found in repository\n"
	WarnPendingMetaNoSHA    = "  ⚠ Skipping %s: no commit SHA in metadata\n"
)

// Clear release JSON messages.
const (
	MsgClearReleaseDone     = "  ✓ Removed .gitmap/release/%s.json\n"
	MsgClearReleaseDryRun   = "  [dry-run] Would remove %s\n"
	ErrClearReleaseUsage    = "Usage: gitmap clear-release-json <version> [--dry-run]\nExample: gitmap clear-release-json v2.20.0"
	ErrClearReleaseNotFound = "Error: no release file found for %s\n"
	ErrClearReleaseFailed   = "Error: could not remove release file at %s: %v (operation: delete)\n"
)

// Revert messages.
const (
	MsgRevertCheckout       = "  → Checking out %s...\n"
	MsgRevertStarting       = "\n  Building reverted version...\n"
	MsgRevertDone           = "\n  ✓ Revert complete.\n"
	ErrRevertUsage          = "Usage: gitmap revert <version>\nExample: gitmap revert v2.9.0"
	ErrRevertTagNotFound    = "Error: tag %s not found locally. Run 'git fetch --tags' first.\n"
	ErrRevertCheckoutFailed = "Error: git checkout failed: %v\n"
	ErrRevertFailed         = "Revert error: %v\n"
	RevertScriptLogExec     = "executing revert script: %s"
	RevertScriptLogExit     = "revert script exited: err=%v"
)

// Network / offline detection.
const (
	NetworkProto      = "tcp"
	NetworkCheckHost  = "github.com:443"
	NetworkTimeoutSec = 5
	ErrOffline        = "network unavailable — cannot reach github.com"
	MsgOfflineWarning = "\n  ⚠ Network unavailable — cannot reach github.com.\n"
	MsgOfflineHint    = "  Offline operations (scan, list, status, group) still work.\n"
)

// Legacy directory migration messages.
const (
	MsgMigrated         = "Migrated %s/ -> %s/\n"
	MsgMergedAndRemoved = "Merged %s/ into %s/ (%d files copied, %d skipped) and removed legacy folder\n"
	ErrMigrationFailed  = "Error: failed to migrate directory %s: %v (operation: move, reason: path is inaccessible)\n"
)

// Legacy ID migration messages.
const (
	MsgLegacyIDMigrationStart = "Migrating database from legacy UUID IDs to integer IDs..."
	MsgLegacyIDMigrationDone  = "Database migration complete. Group-repo associations have been reset."
)

// Direct URL clone messages.
const (
	MsgCloneURLCloning    = "Cloning %s into %s...\n"
	MsgCloneURLDone       = "Cloned %s successfully.\n"
	ErrCloneURLFailed     = "Error: clone failed for %s: %v (operation: git-clone)\n"
	MsgCloneDesktopPrompt = "Add to GitHub Desktop? (y/n): "
	ErrCloneURLExists     = "Error: target folder already exists: %s\n"
)

// Clone replace-existing-folder flow (spec/01-app/96-clone-replace-existing-folder.md).
const (
	MsgCloneReplaceFree       = "  [clone] target free, cloning directly into %s\n"
	MsgCloneReplaceExists     = "  [clone] target exists: %s\n"
	MsgCloneReplaceStrategy1  = "  [clone] strategy 1/2 — direct remove + clone"
	MsgCloneReplaceStrat1Fail = "  [clone] strategy 1/2 failed: %v\n"
	MsgCloneReplaceStrategy2  = "  [clone] strategy 2/2 — temp-clone then swap-in-place"
	MsgCloneReplaceTempClone  = "  [clone] cloning into %s\n"
	MsgCloneReplaceEmptying   = "  [clone] emptying target contents (%d entries) in %s\n"
	MsgCloneReplaceMoving     = "  [clone] moving %d entries from temp into target\n"
	MsgCloneReplaceSwapDone   = "  [clone] swap complete; target now points at fresh clone"
	WarnCloneReplaceEntryFail = "  [clone] could not remove %s: %v\n"
	FlagDescCloneNoReplace    = "Abort if the target folder already exists (disables replace)"
)

// VS Code integration messages.
const (
	VSCodeBin             = "code"
	VSCodeFlagReuseWindow = "--reuse-window"
	VSCodeFlagNewWindow   = "--new-window"
	MsgVSCodeOpening      = "  Opening %s in VS Code...\n"
	MsgVSCodeOpened       = "  VS Code opened."
	MsgVSCodeNotFound     = "  VS Code not found on PATH — skipping editor open.\n"
	ErrVSCodeOpenFailed   = "  Warning: could not open VS Code: %v\n"
	ErrVSCodeAdminLock    = "  Warning: VS Code may be running as administrator — could not open automatically.\n"
)
