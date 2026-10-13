package constants

const (
	HelpUsage            = "Usage: gitmap <command> [flags]"
	HelpCommands         = "Commands:"
	HelpScan             = "  scan (s) [dir]      Scan directory for Git repos"
	HelpClone            = "  clone (c) <source|json|csv|text>  Re-clone from file (shorthands auto-resolve)"
	HelpCloneSync        = "  clone-sync (cs) <url> [urls]      Clone and immediately sync to Workspaces & Antigravity"
	HelpUpdate           = "  update              Self-update from source repo"
	HelpUpdateCleanup    = "  update-cleanup      Remove leftover update temp files and .old backups"
	HelpInstalledDir     = "  installed-dir (id)  Show the active installed binary path"
	HelpVersion          = "  version (v)         Show version number"
	HelpDesktopSync      = "  desktop-sync (ds)   Sync repos to GitHub Desktop from output"
	HelpGitHubDesktop    = "  github-desktop (gd)      Register current repo with GitHub Desktop (no scan needed)"
	HelpPull             = "  pull (p) <name>     Pull a specific repo by its name"
	HelpPullAll          = "  pull-all (pa)       Pull all tracked repos (equivalent to pull " + ColorCyan + "--all" + ColorReset + ")"
	HelpReconcile        = "  reconcile (recon)   Reconcile failed or dirty pull repositories repo-by-repo"
	HelpFix              = "  fix [repo] [action] Apply pull remediation to a repo (aliases: stash, wip, discard)"
	HelpLowerCaseFix     = "  lowercase (lcf)     Rename uppercase files to lowercase via safe 2-step git mv (" + ColorCyan + "--dry-run" + ColorReset + ", " + ColorCyan + "--no-commit" + ColorReset + ", " + ColorCyan + "--readme" + ColorReset + ")"
	HelpLowerCaseReadme  = "  lowercase-readme    Rename root README (e.g. README.md -> readme.md) via safe 2-step git mv (" + ColorCyan + "--dry-run" + ColorReset + ", " + ColorCyan + "--no-commit" + ColorReset + ")"
	HelpRescan           = "  rescan (rsc)        Re-run last scan with cached flags"
	HelpSetup            = "  setup               Configure Git diff/merge tool, aliases & core settings"
	HelpStatus           = "  status (st)         Show dirty/clean, ahead/behind, stash for all repos"
	HelpExec             = "  exec (x) <args...>  Run any git command across all repos"
	HelpRelease          = "  release (r) [ver]   Create release branch, tag, and push"
	HelpReleasePull      = "  pull-release (relp) [ver]  git pull (" + ColorCyan + "--ff-only" + ColorReset + "|" + ColorCyan + "--rebase" + ColorReset + "|" + ColorCyan + "--merge" + ColorReset + "), then release the current repo"
	HelpReleaseBr        = "  release-branch (rb)      Complete release from existing release branch"
	HelpReleasePend      = "  release-pending (rp)      Release all pending branches without tags"
	HelpFixReleaseTags   = "  fix-release-tags (frt) Delete orphan or broken GitHub releases and tags lacking assets or with failed CI/CD"
	HelpFixRepo          = "  fix-repo (fr)       Rewrite prior {base}-vN tokens (-2/-3/-5/" + ColorCyan + "--all" + ColorReset + "/" + ColorCyan + "--dry-run" + ColorReset + ")"
	HelpChangelog        = "  changelog (cl) [ver]      Show concise release notes (use " + ColorCyan + "--open" + ColorReset + ", " + ColorCyan + "--source" + ColorReset + ")"
	HelpDoctor           = "  doctor [" + ColorCyan + "--fix-path" + ColorReset + "]      Diagnose PATH, deploy, and version issues"
	HelpLatestBr         = "  latest-branch (lb)  Find most recently updated remote branch"
	HelpList             = "  list (ls)           Show all tracked repos with slugs"
	HelpGroup            = "  group (g) <sub>     Manage repo groups / activate group for batch ops"
	HelpMultiGroup       = "  multi-group (mg)    Select multiple groups for batch operations"
	HelpSf               = "  sf <add|list|rm>    Manage scan folders (roots that gitmap scan tracks)"
	HelpRm               = "  rm (remove, del) <name-or-path>...  Untrack repos from the database only"
	HelpDBReset          = "  db-reset " + ColorCyan + "--confirm" + ColorReset + "  Clear all tracked repos and groups from the database"
	HelpDB               = "  db <ls|repo-db|sizes|reset> Manage SQLite master and split databases"
	HelpStartFresh       = "  start-fresh                 Clear all databases & reinitialize from scratch"
	HelpFindDuplicates   = "  find-duplicates [platform]  Scan and remediate duplicate projects/repos across vscode/chrome/agy/git"
	HelpCompletion       = "  completion (cmp)    Generate shell tab-completion scripts"
	HelpClearReleaseJSON = "  clear-release-json (crj)  Remove a .gitmap/release/vX.Y.Z.json file"
	HelpDocs             = "  docs (d)            Open documentation website in browser"
	HelpHelpDash         = "  help-dashboard (hd)      Serve the docs site locally in your browser"
	HelpCloneNext        = "  clone-next (cn)     Clone next versioned iteration of current repo"
	HelpReleaseSelf      = "  release-self (rs)   Release gitmap itself from any directory"
	HelpHelp             = "  help                Show this help message"
	HelpListVersions     = "  list-versions (lv)  Show all release tags, highest first (" + ColorCyan + "--limit" + ColorReset + " N, " + ColorCyan + "--json" + ColorReset + ", " + ColorCyan + "--source" + ColorReset + ")"
	HelpListReleases     = "  list-releases (lr)  Show releases from .gitmap/release/ files or database (" + ColorCyan + "--limit" + ColorReset + " N, " + ColorCyan + "--json" + ColorReset + ", " + ColorCyan + "--source" + ColorReset + ")"
	HelpRevert           = "  revert <version>    Revert to a specific release version"
	MsgHelpLFSCommon     = "  lfs-common (lfsc)   Track common binary file types with Git LFS in current repo"
	HelpPower            = "  power (pw, pwr)     Inspect, set, and reset OS display and sleep timeouts"
	HelpOS               = "  os                  Manage OS configurations and diagnostics (" + ColorCyan + "display" + ColorReset + ", " + ColorCyan + "fix-link" + ColorReset + ", " + ColorCyan + "status" + ColorReset + ")"
	HelpFixLink          = "  fix-link (fixlink)  Inspect and repair broken symlinks and shared directories (" + ColorCyan + "--target" + ColorReset + ", " + ColorCyan + "--recursive" + ColorReset + ", " + ColorCyan + "--force" + ColorReset + ", " + ColorCyan + "--dry-run" + ColorReset + ", " + ColorCyan + "--json" + ColorReset + ")"
	HelpVmware           = "  vmware (vm)         Manage VMware shared folders, tools, and services (" + ColorCyan + "install" + ColorReset + ", " + ColorCyan + "shared" + ColorReset + ", " + ColorCyan + "status" + ColorReset + ")"
	HelpNginx            = "  nginx (ngx)         Manage Nginx HTTP server, virtual hosts, syntax test & reload"
	HelpAi               = "  ai (scripts)        Native AI scripts catalog, execution, and autofix engine"
	HelpPendingCommits   = "  pending-commits (pc) Discover uncommitted changes & unpushed commits with tree remediation"
	HelpNewCommands      = "  new-commands (nc)    Inspect last 100 commands added in recent commits with examples"
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
	FlagScanCompact     = "compact"
	FlagDescScanCompact = "Export gitmap.json to current directory and print concise clone commands"
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
	// MsgFixGitStateDeprecated is printed to os.Stderr when the old
	// git-state `fix` surface is invoked (now owned by the content fixer).
	// Per docs/deprecation-policy.md: stderr one-liner naming replacements.
	MsgFixGitStateDeprecated = "  · git-state `fix` is deprecated; use `gitmap stash` / `gitmap wip` / `gitmap discard` instead\n"
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
	FlagDescScanMaxDepth = "Max directory levels to descend below scan root (0 = default 16, negative = unlimited)"
	DefaultScanMaxDepth  = 0
	// RescanSubtreeDefaultMaxDepth is the cap `gitmap rescan-subtree`
	// applies when the user does not pass an explicit --max-depth.
	RescanSubtreeDefaultMaxDepth = 32
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
	FlagDescOpenChangelog    = "Open changelog.md with the default system app"
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
	FlagDescLBSwitch          = "Checkout the resolved latest branch after printing the report"
	FlagDescLBSwitchShort     = "Short alias for --switch"
	FlagDescGroup             = "Filter by group name"
	FlagDescAll               = "Run against all tracked repos from database"
	FlagDescListVerbose       = "Show full paths and URLs"
	FlagDescGroupDesc         = "Optional group description"
	FlagDescGroupColor        = "Terminal color for group display"
	FlagConfirm               = "confirm"
	FlagDescConfirm           = "Confirm destructive operation"
	FlagRescan                = "rescan"
	FlagDescRescan            = "Trigger a rescan immediately after reset"
	FlagDescSource            = "Filter by source: release or import"
	FlagSkipZsh               = "skip-zsh"
	FlagDescSkipZsh           = "Skip ZSH and Oh-My-Zsh setup on Ubuntu"
	EnvGitmapSkipZsh          = "GITMAP_SKIP_ZSH"
	FlagNoCache               = "no-cache"
	FlagDescNoCache           = "Bypass status cache and query git status directly from filesystem"
	FlagBackup                = "backup"
	FlagDescBackup            = "Serve status from secondary persistent backup DB"
	FlagServeBackup           = "serve-backup"
	FlagDescServeBackup       = "Serve status from secondary persistent backup DB (alias for --backup)"
	FlagDetail                = "detail"
	FlagCategory              = "category"
	FlagDescCategory          = "Filter new commands by functional category"
	FlagNewCommandsFilter     = "filter"
	FlagDescNewCommandsFilter = "Search commands by name, alias, description, or example"
)

// FlagDebugPaths exposes the canonicalize-path tracer used by the
// VS Code Project Manager sync helper. When set on `gitmap clone`,
// the flag exports GITMAP_DEBUG_PATHS=1 for the current process so
// every clone variant that routes through canonicalizePMPath emits
// a one-line stderr trace of (raw → cleaned → resolved) rootPaths.
// Diagnostic-only — Windows 8.3 short-name and symlink resolution is
// the primary debugging target. Off by default.
