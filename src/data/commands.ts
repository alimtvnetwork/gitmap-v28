export interface CommandSeeAlso {
  name: string;
  description: string;
  url?: string;
}

export interface CommandFlag {
  flag: string;
  description: string;
}

export interface CommandExample {
  command: string;
  description?: string;
}

export interface CommandDef {
  name: string;
  alias?: string;
  description: string;
  usage?: string;
  flags?: CommandFlag[];
  examples?: CommandExample[];
  category: string;
  seeAlso?: CommandSeeAlso[];
}

export interface CommandCategory {
  key: string;
  label: string;
  description: string;
  icon?: string;
}

export const Categories: CommandCategory[] = [
  { key: "scanning", label: "Scanning & Discovery", description: "Find and catalog Git repositories on disk", icon: "🔍" },
  { key: "cloning", label: "Cloning & Pulling", description: "Clone, pull, and sync repositories", icon: "📥" },
  { key: "monitoring", label: "Monitoring & Status", description: "Track repo state and run batch git commands", icon: "📡" },
  { key: "release", label: "Release & Versioning", description: "Create releases, tags, and branches", icon: "🚀" },
  { key: "changelog", label: "Changelog & Tags", description: "View release notes, list tags, and manage metadata", icon: "📋" },
  { key: "navigation", label: "Navigation & Groups", description: "Move between repos and organize them into groups", icon: "🧭" },
  { key: "history", label: "History & Analytics", description: "Audit trail, usage stats, and dashboards", icon: "📊" },
  { key: "detection", label: "Project Detection", description: "Query repos by detected language or framework", icon: "🔬" },
  { key: "data", label: "Data & Profiles", description: "Export, import, bookmark, and manage profiles", icon: "💾" },
  { key: "tools", label: "Tools & Setup", description: "Setup, diagnostics, updates, SSH, and shell completions", icon: "🔧" },
  { key: "movemerge", label: "Move & Merge", description: "Move folders or merge file trees across local and remote endpoints", icon: "🔀" },
  { key: "chrome", label: "Chrome Profiles", description: "Copy, merge, export, import, list, and delete Chrome browser profiles", icon: "🌐" },
  { key: "visibility-bulk", label: "Bulk Visibility", description: "Flip many repos public/private at once with audit trail", icon: "🔐" },
  { key: "maintenance", label: "Maintenance & Backups", description: "Backups, SSH health, doctor, changelog regen, and file deletion", icon: "🧹" },
];

export const commands: CommandDef[] = [
  // ═══════════════════════════════════════════
  // Scanning & Discovery
  // ═══════════════════════════════════════════
  {
    category: "scanning",
    name: "scan", alias: "s", description: "Scan directory tree for Git repositories and generate output files",
    usage: "gitmap scan [dir] [--output csv|json|terminal] [--mode ssh|https]",
    flags: [
      { flag: "--config <path>", description: "Config file path" },
      { flag: "--mode ssh|https", description: "Clone URL style (default: https)" },
      { flag: "--output csv|json|terminal", description: "Output format (default: terminal)" },
      { flag: "--output-path <dir>", description: "Output directory" },
      { flag: "--github-desktop", description: "Register repos in GitHub Desktop" },
      { flag: "--open", description: "Open output folder after scan" },
      { flag: "--quiet", description: "Suppress clone help section" },
      { flag: "--no-vscode-sync", description: "Skip syncing scanned repos into VS Code Project Manager projects.json (default: sync ON)" },
      { flag: "--no-auto-tags", description: "Skip auto-derived tags (git/node/go/python/rust/docker) when syncing (default: tags ON)" },
    ],
    examples: [
      { command: "gitmap scan ~/projects", description: "Scan all repos under ~/projects (auto-syncs into VS Code Project Manager + auto-tags)" },
      { command: "gitmap s --output json --mode ssh", description: "JSON output with SSH clone URLs" },
      { command: "gitmap scan C:\\dev --github-desktop --open", description: "Scan, register in GitHub Desktop, open folder" },
      { command: "gitmap s --output csv --output-path ./backup", description: "CSV output to custom directory" },
      { command: "gitmap scan ~/projects --no-vscode-sync", description: "Scan without touching VS Code Project Manager" },
      { command: "gitmap scan ~/projects --no-auto-tags", description: "Sync without auto-derived language/tooling tags" },
    ],
    seeAlso: [
      { name: "rescan", description: "Re-scan using cached parameters" },
      { name: "clone", description: "Clone repos from scan output" },
      { name: "status", description: "View repo statuses after scanning" },
      { name: "desktop-sync", description: "Sync scanned repos to GitHub Desktop" },
      { name: "export", description: "Export scanned data" },
    ],
  },
  {
    category: "scanning",
    name: "rescan", alias: "rsc", description: "Re-run the last scan using cached parameters (same dir, flags, output)",
    usage: "gitmap rescan",
    examples: [
      { command: "gitmap rescan", description: "Repeat last scan with identical settings" },
      { command: "gitmap rsc", description: "Alias shorthand" },
    ],
    seeAlso: [
      { name: "scan", description: "Initial directory scan" },
      { name: "status", description: "View repo statuses" },
      { name: "clone", description: "Clone from scan output" },
    ],
  },
  {
    category: "scanning",
    name: "scan all", alias: "scan a", description: "Re-scan every previously-scanned root folder in parallel (planned v3.33.0)",
    usage: "gitmap scan all [--workers N] [--prune-missing] [--mode ssh|https] [--output csv|json|terminal]",
    flags: [
      { flag: "--workers <n>", description: "Parallel workers (1–16, default 4)" },
      { flag: "--prune-missing", description: "Auto-remove missing roots from DB without prompting" },
      { flag: "--mode ssh|https", description: "Forwarded to each per-root scan" },
      { flag: "--output csv|json|terminal", description: "Forwarded to each per-root scan" },
      { flag: "--github-desktop", description: "Forwarded to each per-root scan" },
      { flag: "--quiet", description: "Suppress per-root clone help section" },
    ],
    examples: [
      { command: "gitmap scan all", description: "Re-scan every root from the ScanFolder table (4 parallel workers)" },
      { command: "gitmap scan a", description: "Short alias" },
      { command: "gitmap scan all --workers 8", description: "Bump parallelism to 8" },
      { command: "gitmap scan all --prune-missing", description: "Cron-friendly: auto-prune unreachable roots" },
    ],
    seeAlso: [
      { name: "Spec: scan all", description: "Full specification document", url: "/scan-all" },
      { name: "scan", description: "Single-root scan (the source that populates ScanFolder)" },
      { name: "rescan", description: "Repeat the most recent single-root scan" },
      { name: "sf list", description: "Inspect the ScanFolder table directly" },
      { name: "find-next", description: "List repos with new versions available" },
    ],
  },
  {
    category: "scanning",
    name: "scan gd", alias: "scan github-desktop", description: "Register every repo under the current scan root in GitHub Desktop, sequentially + idempotent (planned v3.35.0)",
    usage: "gitmap scan gd  (or: gitmap scan github-desktop)",
    flags: [],
    examples: [
      { command: "cd D:\\projects && gitmap scan gd", description: "Register every repo under D:\\projects in Desktop" },
      { command: "gitmap scan github-desktop", description: "Long form" },
      { command: "gitmap s gd", description: "Inherits the existing s alias for scan" },
    ],
    seeAlso: [
      { name: "scan", description: "Use --github-desktop on a fresh scan to discover + register in one pass" },
      { name: "desktop-sync", description: "Legacy: register every tracked repo globally (no scan-root scoping)" },
      { name: "scan all", description: "Bulk re-scan every known root first" },
      { name: "sf list", description: "Inspect the ScanFolder table" },
    ],
  },
  {
    category: "scanning",
    name: "desktop-sync", alias: "ds", description: "Register all tracked repos with GitHub Desktop",
    usage: "gitmap desktop-sync",
    examples: [
      { command: "gitmap desktop-sync", description: "Sync all tracked repos to GitHub Desktop" },
      { command: "gitmap ds", description: "Alias shorthand" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories first" },
      { name: "clone", description: "Clone repos from scan output" },
      { name: "list", description: "List tracked repos" },
    ],
  },
  {
    category: "scanning",
    name: "probe", alias: "pb", description: "Check tracked repos for newer remote tags via ls-remote → shallow-clone fallback. With --depth N, walks up to N newer versions and verifies each via shallow clone (planned v3.36.0)",
    usage: "gitmap probe [<path>|--all] [--depth N] [--json]",
    flags: [
      { flag: "--all", description: "Probe every repo in the database" },
      { flag: "--depth <n>", description: "Walk up to N newer tags, shallow-verify each (default 1, max 10) — planned v3.36.0" },
      { flag: "--json", description: "Emit results as JSON for CI consumption" },
    ],
    examples: [
      { command: "gitmap probe --all", description: "Probe every tracked repo for the single newest tag (current behavior)" },
      { command: "gitmap probe --all --depth 5", description: "Walk up to 5 newer versions per repo, shallow-clone each to verify" },
      { command: "gitmap probe E:\\src\\my-repo --depth 3", description: "Single repo, 3-deep walk" },
      { command: "gitmap probe --all --depth 5 --json > probes.json", description: "Machine-readable upgrade-path report" },
    ],
    seeAlso: [
      { name: "find-next", description: "Read VersionProbe results without re-probing" },
      { name: "scan", description: "Populate the database that probe reads from" },
      { name: "scan all", description: "Bulk re-scan known roots before probing" },
    ],
  },
  {
    category: "scanning",
    name: "find-next", alias: "fn", description: "Read-only: list every repo whose latest VersionProbe row has IsAvailable=1, sorted newest first",
    usage: "gitmap find-next [--scan-folder <id>] [--include-intermediate] [--json]",
    flags: [
      { flag: "--scan-folder <id>", description: "Restrict to one ScanFolder (look up via 'gitmap sf list')" },
      { flag: "--include-intermediate", description: "Show every verified version from the latest --depth walk, not just the newest (planned v3.36.0)" },
      { flag: "--json", description: "Emit []FindNextRow as indented JSON" },
    ],
    examples: [
      { command: "gitmap find-next", description: "Every repo with an available update across the whole DB" },
      { command: "gitmap fn --scan-folder 3", description: "Only repos under ScanFolder id=3" },
      { command: "gitmap fn --json", description: "JSON output for CI" },
    ],
    seeAlso: [
      { name: "probe", description: "Run the probe to refresh data find-next reads" },
      { name: "sf list", description: "Find a ScanFolder id" },
      { name: "pull all", description: "Update the repos that find-next surfaces" },
    ],
  },

  // ═══════════════════════════════════════════
  // Cloning & Pulling
  // ═══════════════════════════════════════════
  {
    category: "cloning",
    name: "clone", alias: "c", description: "Re-clone repos from a structured scan output file (JSON, CSV, or text)",
    usage: "gitmap clone <source|json|csv|text> [--target-dir <dir>] [--safe-pull]",
    flags: [
      { flag: "--target-dir <dir>", description: "Base directory for cloned repos" },
      { flag: "--safe-pull", description: "Pull existing repos with retry + diagnostics" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap clone json --target-dir ./projects", description: "Clone all repos from JSON to ./projects" },
      { command: "gitmap c csv", description: "Clone from CSV scan output" },
      { command: "gitmap clone json --safe-pull", description: "Pull existing, clone missing" },
      { command: "gitmap c text --target-dir D:\\repos --verbose", description: "Clone from text with logging" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories to create clone source" },
      { name: "pull", description: "Pull updates for existing repos" },
      { name: "desktop-sync", description: "Sync cloned repos to GitHub Desktop" },
    ],
  },
  {
    category: "cloning",
    name: "clone-next", alias: "cn", description: "Clone the next versioned iteration of the current repo (e.g. v11 → v12)",
    usage: "gitmap clone-next <v++|vN> [--delete] [--keep] [--no-desktop]",
    flags: [
      { flag: "--delete", description: "Auto-remove current folder after successful clone" },
      { flag: "--keep", description: "Keep current folder without prompting" },
      { flag: "--no-desktop", description: "Skip GitHub Desktop registration" },
      { flag: "--ssh-key <name>", description: "Use a named SSH key for the clone" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap cn v++", description: "Increment version by one (v11 → v12)" },
      { command: "gitmap cn v15 --delete", description: "Jump to v15, auto-remove current folder, navigate to new" },
      { command: "gitmap clone-next v++ --keep", description: "Increment version, keep current folder" },
      { command: "gitmap cn v++ --ssh-key work", description: "Clone next version using work SSH key" },
    ],
    seeAlso: [
      { name: "clone", description: "Clone repos from structured file" },
      { name: "desktop-sync", description: "Sync repos to GitHub Desktop" },
      { name: "ssh", description: "Manage named SSH keys" },
    ],
  },
  {
    category: "cloning",
    name: "pull", alias: "p", description: "Pull latest changes for specific repos, groups, or all tracked repos",
    usage: "gitmap pull <repo-name> [--group <name>] [--all] [--verbose]",
    flags: [
      { flag: "--group <name>", description: "Pull all repos in a group" },
      { flag: "--all", description: "Pull all tracked repos" },
      { flag: "--verbose", description: "Enable verbose logging" },
    ],
    examples: [
      { command: "gitmap pull my-api-service", description: "Pull a single repo by exact name" },
      { command: "gitmap p my-api", description: "Partial match — finds my-api-service" },
      { command: "gitmap pull --group backend", description: "Pull all repos in the backend group" },
      { command: "gitmap pull --all", description: "Pull every tracked repo" },
      { command: "gitmap p --all --verbose", description: "Pull all with detailed logging" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories to populate the database" },
      { name: "clone", description: "Clone repos from structured file" },
      { name: "status", description: "View repo statuses" },
      { name: "group", description: "Manage repo groups for targeted pulls" },
    ],
  },
  {
    category: "cloning",
    name: "pull all", alias: "pull a", description: "Pull every repo under the CWD scan root in parallel; run.ps1/run.sh replaces git pull (planned v3.34.0)",
    usage: "gitmap pull all [--workers N] [--script-timeout <dur>] [--verbose]",
    flags: [
      { flag: "--workers <n>", description: "Parallel workers (1–16, default 4)" },
      { flag: "--script-timeout <dur>", description: "Per-script timeout (default 10m)" },
      { flag: "--verbose", description: "Forwarded to each per-repo pull" },
    ],
    examples: [
      { command: "cd D:\\projects && gitmap pull all", description: "Pull every repo registered under D:\\projects" },
      { command: "gitmap p all", description: "Short form using existing pull alias" },
      { command: "gitmap pull a --workers 8", description: "Bump parallelism" },
      { command: "gitmap pull all --script-timeout 30m", description: "Allow longer-running run.ps1/run.sh builds" },
    ],
    seeAlso: [
      { name: "pull", description: "Pull a single repo / group / global --all" },
      { name: "scan all", description: "Bulk-rescan all known roots first" },
      { name: "scan", description: "Register the current dir as a scan root" },
      { name: "sf list", description: "Inspect ScanFolder table" },
    ],
  },

  // ═══════════════════════════════════════════
  // Monitoring & Status
  // ═══════════════════════════════════════════
  {
    category: "monitoring",
    name: "status", alias: "st", description: "Show a one-shot status dashboard for tracked repos (dirty/clean, ahead/behind)",
    usage: "gitmap status [--group <name>] [--all]",
    flags: [
      { flag: "--group <name>", description: "Show status for repos in a specific group" },
      { flag: "--all", description: "Show status for every repo in the database" },
    ],
    examples: [
      { command: "gitmap status", description: "Status dashboard for default group" },
      { command: "gitmap st --group backend", description: "Status for backend group only" },
      { command: "gitmap status --all", description: "Status for every tracked repo" },
    ],
    seeAlso: [
      { name: "watch", description: "Live-refresh status dashboard", url: "/watch" },
      { name: "scan", description: "Scan directories to populate data" },
      { name: "exec", description: "Run git commands across repos" },
      { name: "group", description: "Filter status by group" },
    ],
  },
  {
    category: "monitoring",
    name: "watch", alias: "w", description: "Live-refresh status dashboard with configurable interval and group filter",
    usage: "gitmap watch [--interval <seconds>] [--group <name>] [--no-fetch] [--json]",
    flags: [
      { flag: "--interval <seconds>", description: "Refresh interval (default: 30, min: 5)" },
      { flag: "--group <name>", description: "Monitor only repos in a group" },
      { flag: "--no-fetch", description: "Skip git fetch on each refresh" },
      { flag: "--json", description: "Output single snapshot as JSON (no loop)" },
    ],
    examples: [
      { command: "gitmap watch", description: "Live dashboard, 30s refresh" },
      { command: "gitmap w --interval 10 --group frontend", description: "Fast refresh for frontend group" },
      { command: "gitmap watch --no-fetch --json", description: "Snapshot without fetching, as JSON" },
    ],
    seeAlso: [
      { name: "Spec: watch", description: "Live monitor documentation", url: "/watch" },
      { name: "status", description: "One-time status snapshot" },
      { name: "exec", description: "Run git commands across repos" },
      { name: "group", description: "Filter by group" },
    ],
  },
  {
    category: "monitoring",
    name: "has-any-updates", alias: "hau / hac", description: "Check if the remote has new commits you haven't pulled yet",
    usage: "gitmap has-any-updates\ngitmap hau\ngitmap hac",
    examples: [
      { command: "gitmap hau", description: "Quick check — are there unpulled remote commits?" },
      { command: "gitmap hac", description: "Same check (alternate alias)" },
    ],
    seeAlso: [
      { name: "status", description: "Show repo status dashboard" },
      { name: "pull", description: "Pull a specific repo" },
      { name: "watch", description: "Live-refresh dashboard", url: "/watch" },
    ],
  },
  {
    category: "monitoring",
    name: "exec", alias: "x", description: "Run any git command across all tracked repos simultaneously",
    usage: "gitmap exec <git-args...>",
    examples: [
      { command: "gitmap exec fetch --prune", description: "Fetch and prune all repos" },
      { command: "gitmap x remote -v", description: "Show remotes for every tracked repo" },
      { command: "gitmap exec status --short", description: "Quick git status across all repos" },
      { command: "gitmap x branch --list", description: "List branches in every repo" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories to populate the database" },
      { name: "pull", description: "Pull repos (built-in alternative)" },
      { name: "status", description: "View repo statuses" },
    ],
  },
  {
    category: "monitoring",
    name: "latest-branch", alias: "lb", description: "Find the most recently updated remote branches with filtering and sorting",
    usage: "gitmap latest-branch [--top N] [--format json|csv|terminal] [--filter <pattern>]",
    flags: [
      { flag: "--remote <name>", description: "Remote to filter against (default: origin)" },
      { flag: "--all-remotes", description: "Include branches from all remotes" },
      { flag: "--contains-fallback", description: "Fall back to --contains if --points-at is empty" },
      { flag: "--top <n>", description: "Show top N most recently updated branches" },
      { flag: "--format <fmt>", description: "Output format: terminal, json, csv" },
      { flag: "--json", description: "Shorthand for --format json" },
      { flag: "--no-fetch", description: "Skip git fetch, use existing refs" },
      { flag: "--sort <order>", description: "Sort: date (descending) or name (A-Z)" },
      { flag: "--filter <pattern>", description: "Filter branches by glob or substring" },
    ],
    examples: [
      { command: "gitmap lb", description: "Show the single most recent branch" },
      { command: "gitmap lb --top 5 --json", description: "Top 5 branches as JSON" },
      { command: "gitmap lb --filter 'feature/*'", description: "Only feature branches" },
      { command: "gitmap lb 3 --no-fetch --json", description: "Fast: skip fetch, top 3 as JSON" },
      { command: "gitmap lb --all-remotes --top 10", description: "Top 10 across all remotes" },
    ],
    seeAlso: [
      { name: "status", description: "View repo statuses" },
      { name: "watch", description: "Live-refresh dashboard", url: "/watch" },
      { name: "release-branch", description: "Create a release branch" },
    ],
  },

  // ═══════════════════════════════════════════
  // Release & Versioning
  // ═══════════════════════════════════════════
  {
    category: "release",
    name: "release", alias: "r", description: "Create a release: branch, tag, push, and optionally attach compiled assets",
    usage: "gitmap release [version] [--bump major|minor|patch] [--draft] [--dry-run]",
    flags: [
      { flag: "--assets <path>", description: "Attach files to release" },
      { flag: "--commit <sha>", description: "Release from specific commit" },
      { flag: "--branch <name>", description: "Release from branch" },
      { flag: "--bump major|minor|patch", description: "Auto-increment version" },
      { flag: "--draft", description: "Create unpublished draft" },
      { flag: "--dry-run", description: "Preview without executing" },
      { flag: "--compress", description: "Wrap assets in .zip (Windows) or .tar.gz" },
      { flag: "--checksums", description: "Generate SHA256 checksums.txt" },
      { flag: "--no-assets", description: "Skip Go binary cross-compilation" },
      { flag: "--targets <list>", description: "Cross-compile targets (e.g. windows/amd64,linux/arm64)" },
      { flag: "--list-targets", description: "Print resolved target matrix and exit" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap release --bump patch", description: "Patch release with auto-incremented version" },
      { command: "gitmap r v2.5.0 --draft", description: "Draft release for v2.5.0" },
      { command: "gitmap release --bump minor --dry-run", description: "Preview a minor release" },
      { command: "gitmap r --bump patch --compress --checksums", description: "Release with compressed assets and checksums" },
      { command: "gitmap release --targets windows/amd64,linux/arm64", description: "Cross-compile for specific platforms" },
    ],
    seeAlso: [
      { name: "Spec: release", description: "Full release workflow documentation", url: "/release" },
      { name: "release-self", description: "Release gitmap itself from any directory", url: "/release-self" },
      { name: "release-branch", description: "Create branch without tagging" },
      { name: "release-pending", description: "Show unreleased commits" },
      { name: "changelog", description: "View release notes" },
      { name: "list-versions", description: "List available tags" },
    ],
  },
  {
    category: "release",
    name: "release-undo", alias: "ru", description: "Undo the most recent (or a specific) release — deletes the local tag, the matching remote tag, and the .gitmap/release/<version>.json sidecar. Safe to run repeatedly; missing pieces are skipped with a warning.",
    usage: "gitmap release-undo [version] [--keep-remote] [--keep-sidecar] [-y]",
    flags: [
      { flag: "[version]", description: "Specific release to undo (e.g. v6.65.0). Defaults to the latest release if omitted." },
      { flag: "--keep-remote", description: "Skip deleting the matching remote tag (origin)" },
      { flag: "--keep-sidecar", description: "Keep .gitmap/release/<version>.json on disk" },
      { flag: "-y, --yes", description: "Skip the confirmation prompt" },
    ],
    examples: [
      { command: "gitmap release-undo", description: "Undo the most recent release (prompts before deleting)" },
      { command: "gitmap ru v6.65.0 -y", description: "Undo a specific release with no prompt" },
      { command: "gitmap ru --keep-remote", description: "Drop the local tag + sidecar but leave the remote tag intact" },
      { command: "gitmap ru v6.64.0 --keep-sidecar", description: "Remove tags but preserve the JSON release record" },
    ],
    seeAlso: [
      { name: "release", description: "Create a release (the operation this undoes)", url: "/release" },
      { name: "changelog", description: "View release notes" },
      { name: "list-versions", description: "List available tags" },
    ],
  },
  {
    category: "release",
    name: "commit-in", alias: "cin", description: "Replay commits from one or more source repos/URLs into a destination repo with full author, date, and tag fidelity. Resumable via .gitmap/commit-in/state.json — re-running picks up where a crash left off.",
    usage: "gitmap commit-in <destination> <source1,source2,...> [--profile <name>] [--save-profile <name>] [--set-default] [--resume]",
    flags: [
      { flag: "<destination>", description: "Target folder (auto-created and `git init`-ed if missing)" },
      { flag: "<sources>", description: "Comma-separated list of local paths or git URLs to replay" },
      { flag: "--profile <name>", description: "Load saved profile from .gitmap/commit-in/profiles/<name>.json" },
      { flag: "--save-profile <name>", description: "Persist the current invocation as a reusable profile" },
      { flag: "--set-default", description: "Mark the saved profile as the default" },
      { flag: "--resume", description: "Continue from the last checkpoint in state.json" },
      { flag: "--dry-run", description: "Plan the replay without writing commits" },
    ],
    examples: [
      { command: "gitmap commit-in ./my-project https://github.com/me/my-project-archive.git", description: "Pull history from a single remote into a local folder" },
      { command: "gitmap cin ./brand-new ./src-a,./src-b --save-profile MyMerge --set-default", description: "Merge two sources and save the recipe as a profile" },
      { command: "gitmap cin ./gitmap --profile Default", description: "Reuse the default saved profile" },
      { command: "gitmap cin ./dest ./src --resume", description: "Resume an interrupted replay from the last checkpoint" },
    ],
    seeAlso: [
      { name: "Spec: commit-in", description: "Full commit-in walkthrough with examples", url: "/commit-in" },
      { name: "release", description: "Tag and publish after replay" },
    ],
  },
  {

    category: "release",
    name: "release-self", alias: "rs / rself", description: "Release gitmap itself from any directory (uses embedded repo path)",
    usage: "gitmap release-self [version] [--bump major|minor|patch] [--draft] [--dry-run]",
    flags: [
      { flag: "--assets <path>", description: "Attach files to release" },
      { flag: "--commit <sha>", description: "Release from specific commit" },
      { flag: "--branch <name>", description: "Release from branch" },
      { flag: "--bump major|minor|patch", description: "Auto-increment version" },
      { flag: "--draft", description: "Create unpublished draft" },
      { flag: "--dry-run", description: "Preview without executing" },
      { flag: "--compress", description: "Wrap assets in .zip or .tar.gz" },
      { flag: "--checksums", description: "Generate SHA256 checksums.txt" },
      { flag: "--no-assets", description: "Skip Go binary cross-compilation" },
      { flag: "--targets <list>", description: "Cross-compile targets" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap rs --bump patch", description: "Self-release with patch bump from any directory" },
      { command: "gitmap release-self v2.46.0 --dry-run", description: "Preview self-release without executing" },
      { command: "gitmap rs --bump minor --draft", description: "Draft minor self-release" },
      { command: "gitmap rs --bump patch --compress --checksums", description: "Full self-release with assets" },
    ],
    seeAlso: [
      { name: "Spec: release-self", description: "Full release-self documentation", url: "/release-self" },
      { name: "release", description: "Standard release workflow", url: "/release" },
      { name: "release-branch", description: "Complete from existing branch" },
    ],
  },
  {
    category: "release",
    name: "release-branch", alias: "rb", description: "Complete a release from an existing release/* branch (no new branch created)",
    usage: "gitmap release-branch [version] [--bump major|minor|patch] [--draft] [--verbose]",
    flags: [
      { flag: "--assets <path>", description: "Directory or file to attach" },
      { flag: "--draft", description: "Create an unpublished draft release" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap release-branch release/v1.2.0", description: "Complete release from existing branch" },
      { command: "gitmap rb release/v1.2.0 --draft", description: "Complete as draft release" },
      { command: "gitmap rb release/v1.2.0 --assets ./dist", description: "Attach dist folder to release" },
    ],
    seeAlso: [
      { name: "release", description: "Full release with tag and push", url: "/release" },
      { name: "release-pending", description: "Show unreleased commits" },
      { name: "latest-branch", description: "Find most recent branch" },
    ],
  },
  {
    category: "release",
    name: "release-pending", alias: "rp", description: "Find and release untagged release branches and metadata-only versions",
    usage: "gitmap release-pending [--assets <path>] [--draft] [--dry-run] [--verbose]",
    flags: [
      { flag: "--assets <path>", description: "Directory or file to attach" },
      { flag: "--draft", description: "Mark release metadata as draft" },
      { flag: "--dry-run", description: "Preview steps without executing" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap release-pending", description: "Release all pending versions" },
      { command: "gitmap rp --dry-run", description: "Preview what would be released" },
      { command: "gitmap rp --draft --verbose", description: "Release as drafts with logging" },
    ],
    seeAlso: [
      { name: "release", description: "Create a release", url: "/release" },
      { name: "release-branch", description: "Complete from existing branch" },
      { name: "clear-release-json", description: "Remove a release metadata file" },
      { name: "changelog", description: "View release notes" },
    ],
  },
  {
    category: "release",
    name: "temp-release", alias: "tr", description: "Create lightweight temp branches from recent commits (no tags, no metadata)",
    usage: "gitmap temp-release <count> <version-pattern> [-s N]",
    flags: [
      { flag: "-s, --start", description: "Starting sequence number (default: auto-increment)" },
      { flag: "--dry-run", description: "Preview branch names without creating" },
      { flag: "--json", description: "JSON output for list subcommand" },
      { flag: "--verbose", description: "Detailed logging" },
    ],
    examples: [
      { command: "gitmap tr 10 v1.$$ -s 5", description: "Create 10 branches: v1.05 through v1.14" },
      { command: "gitmap tr 1 v1.$$", description: "Create 1 branch, auto-increment from last" },
      { command: "gitmap tr list", description: "List all temp-release branches" },
      { command: "gitmap tr remove v1.05 to v1.10", description: "Remove a range of temp-release branches" },
      { command: "gitmap tr 5 v2.$$ --dry-run", description: "Preview 5 branch names without creating" },
    ],
    seeAlso: [
      { name: "release", description: "Full release with tags and metadata" },
      { name: "prune", description: "Delete stale release branches" },
      { name: "release-branch", description: "Complete release from existing branch" },
      { name: "temp-release", description: "Dedicated docs page", url: "/temp-release" },
    ],
  },
  {
    category: "release",
    name: "prune", alias: "prn", description: "Delete stale release/* branches that already have a matching tag (alias changed v5.6.0 — was `pr`, now `prn`; `pr` belongs to pull-release)",
    usage: "gitmap prune [flags]",
    flags: [
      { flag: "--dry-run", description: "List stale branches without deleting" },
      { flag: "--confirm", description: "Skip interactive confirmation prompt" },
      { flag: "--remote", description: "Also delete remote release branches" },
    ],
    examples: [
      { command: "gitmap prune --dry-run", description: "Preview which branches would be deleted" },
      { command: "gitmap prune --confirm", description: "Delete stale branches without prompting" },
      { command: "gitmap prune --confirm --remote", description: "Delete locally and remotely" },
      { command: "gitmap pr --dry-run", description: "Alias shorthand preview" },
    ],
    seeAlso: [
      { name: "release", description: "Create release branches and tags" },
      { name: "clear-release-json", description: "Remove release metadata files" },
      { name: "list-releases", description: "Show stored releases from database" },
    ],
  },
  {
    category: "release",
    name: "release-alias", alias: "ra", description: "Release a repo by its registered alias from anywhere on disk (auto-stash + chdir + release)",
    usage: "gitmap release-alias <alias> <version> [--pull] [--no-stash] [--dry-run]",
    flags: [
      { flag: "--pull", description: "Run git pull --ff-only inside the target repo before releasing" },
      { flag: "--no-stash", description: "Abort if working tree is dirty (skip auto-stash)" },
      { flag: "--dry-run", description: "Forwarded to gitmap release — preview only" },
    ],
    examples: [
      { command: "gitmap release-alias my-api v1.4.0", description: "Release the repo registered as 'my-api'" },
      { command: "gitmap ra my-api v1.4.0 --pull", description: "Pull --ff-only first, then release" },
      { command: "gitmap ra backend v0.9.0 --dry-run", description: "Preview the release pipeline" },
    ],
    seeAlso: [
      { name: "as", description: "Register the current repo as an alias first" },
      { name: "release-alias-pull", description: "Equivalent to release-alias --pull" },
      { name: "release", description: "The underlying release workflow" },
    ],
  },
  {
    category: "release",
    name: "release-alias-pull", alias: "rap", description: "Pull-then-release shortcut (release-alias with --pull always implied)",
    usage: "gitmap release-alias-pull <alias> <version> [--no-stash] [--dry-run]",
    flags: [
      { flag: "--no-stash", description: "Abort if working tree is dirty (skip auto-stash)" },
      { flag: "--dry-run", description: "Forwarded to gitmap release — preview only" },
    ],
    examples: [
      { command: "gitmap rap my-api v1.4.0", description: "Pull then release in one shot" },
      { command: "gitmap rap backend v0.9.0 --dry-run", description: "Preview pull + release" },
    ],
    seeAlso: [
      { name: "release-alias", description: "Same command without forced --pull" },
      { name: "as", description: "Register an alias for the current repo" },
    ],
  },
  {
    category: "release",
    name: "pull-release", alias: "pr", description: "Sugar for `release` that first runs `git pull` in the current repo (your choice of mode), then delegates to the standard release pipeline. Default mode is fast-forward only so we never tag on top of a divergent tree. Renamed in v5.6.0 — legacy `release-pull` / `relp` / `rlp` still work.",
    usage: "gitmap pull-release [--ff-only | --rebase | --merge] [--dry-run] [--verbose] [version] [release flags...]",
    flags: [
      { flag: "--ff-only (default)", description: "Pull mode: fast-forward only. Hard-fails on any divergent history." },
      { flag: "--rebase", description: "Pull mode: rebase local commits onto upstream. Runs `git rebase --abort` and exits non-zero on conflict." },
      { flag: "--merge", description: "Pull mode: classic merge. Passes `--no-rebase` so it overrides any user-level pull.rebase=true config." },
      { flag: "--dry-run (pull-release)", description: "Print the `git pull` command that would run, then skip the pull and forward to release. (Release's own --dry-run is forwarded separately.)" },
      { flag: "--verbose (pull-release)", description: "Echo the git invocation to stderr before running it." },
      { flag: "[version]", description: "Forwarded to `gitmap release` (e.g. v1.4.0). Omit for auto-bump prompt." },
      { flag: "--bump <kind>", description: "Forwarded to release: patch | minor | major" },
      { flag: "--bin", description: "Forwarded to release: build & attach binaries" },
      { flag: "--draft", description: "Forwarded to release: create draft release" },
      { flag: "-y, --yes", description: "Forwarded to release: skip confirmation prompts" },
    ],
    examples: [
      { command: "gitmap pr v1.4.0", description: "Default fast-forward pull, then release v1.4.0" },
      { command: "gitmap pr --rebase v1.4.0", description: "Rebase local commits onto upstream, then release" },
      { command: "gitmap pr --merge v2.0.0 --bin --draft", description: "Classic merge (forces --no-rebase), build binaries, push as draft" },
      { command: "gitmap pr --rebase --dry-run", description: "Show the git pull command, skip it, forward to release" },
      { command: "gitmap release-pull v1.4.0", description: "Legacy long form still routes to pull-release" },
    ],
    seeAlso: [
      { name: "release", description: "The underlying release workflow", url: "/release" },
      { name: "release-alias-pull", description: "Same idea, but for a registered alias from any directory" },
      { name: "fix-repo", description: "Rewrite stale {base}-vN tokens after a version bump" },
    ],
  },
  {
    category: "release",
    name: "replace", alias: "rpl", description: "Repo-wide find/replace across every text file. Two modes: literal text swap, or version-suffix bump driven by the git remote URL.",
    usage: 'gitmap replace "<old>" "<new>"   |   gitmap replace -N | --audit | all',
    flags: [
      { flag: "-N (e.g. -2, -3, -5)", description: "Version mode: bump v(current-N)..v(current-1) → vCurrent" },
      { flag: "all", description: "Version mode: bump every prior version v1..v(current-1) → vCurrent" },
      { flag: "--audit", description: "Report-only scan; never writes (overrides --yes)" },
      { flag: "-y, --yes", description: "Skip the y/N confirmation prompt" },
      { flag: "--dry-run", description: "Print summary only, never write" },
      { flag: "-q, --quiet", description: "Suppress per-file diff lines" },
      { flag: "--ext <list>", description: "Comma-separated extension allow-list (e.g. .go,.md). Leading dot optional." },
      { flag: "--ext-case <mode>", description: "sensitive | insensitive (default insensitive). Sensitive matches filenames byte-exact." },
    ],
    examples: [
      { command: 'gitmap replace "old-name" "new-name"', description: "Literal repo-wide replace with confirmation prompt" },
      { command: "gitmap replace -3", description: "Bump v(K-3)..v(K-1) → vK using current repo's remote suffix" },
      { command: "gitmap replace all", description: "Bump every prior version v1..v(K-1) → vK" },
      { command: "gitmap replace --audit", description: "Report every file:line that would change; write nothing" },
      { command: 'gitmap rpl "old" "new" --ext .go,.md -y', description: "Restrict to .go/.md and skip prompt" },
    ],
    seeAlso: [
      { name: "fix-repo", description: "Stricter, git-tracked-only version-suffix rewriter with binary exit codes" },
      { name: "release-self", description: "Bump gitmap's own version" },
      { name: "clone-next", description: "Clone the next versioned repo iteration" },
    ],
  },
  {
    category: "scanning",
    name: "list --update", alias: "lu", description: "List every scanned repo whose local tag or HEAD is behind the latest known remote release. Read-only companion to update apply / update all. Also available as `gitmap update list`.",
    usage: "gitmap list --update [--json] [--source release|import|git-tag] [--limit N] [--stale-days N] [--only-behind] [--path-only]",
    flags: [
      { flag: "--json", description: "Emit a JSON array (schema: spec/08-json-schemas/list-update.schema.json)." },
      { flag: "--source <kind>", description: "Filter by origin of the `latest` value: release, import, or git-tag." },
      { flag: "--limit <n>", description: "Show only the top N upgradable repos, sorted by `behind` DESC." },
      { flag: "--stale-days <n>", description: "Re-probe rows older than N days before listing (0 = never)." },
      { flag: "--only-behind", description: "Hide up-to-date rows in the table view." },
      { flag: "--path-only", description: "Print one absolute repo path per line; script-friendly." },
    ],
    examples: [
      { command: "gitmap list --update", description: "Columned table of upgradable repos with next-step hints." },
      { command: "gitmap lu --path-only", description: "One absolute path per line for piping into xargs / while read." },
      { command: "gitmap lu --source release --limit 5 --json", description: "Top 5 release-tracked repos as JSON." },
    ],
    seeAlso: [
      { name: "scan", description: "Populate the database this command reads from" },
      { name: "update apply", description: "Upgrade one repo to its latest tag" },
      { name: "update all", description: "Upgrade every upgradable repo" },
      { name: "hd", description: "One-screen dashboard including this count" },
      { name: "find-next", description: "Legacy view of repos with new versions available" },
    ],
  },
  {
    category: "tools",
    name: "update apply", alias: "ua", description: "Upgrade one scanned repo to the latest release tag known to gitmap. On failure enqueues a PendingTask of the new `Upgrade` type so `do-pending` can replay it verbatim.",
    usage: "gitmap update apply <repo> [--yes] [--dry-run] [--tag vX.Y.Z] [--strategy fetch-checkout|source-release] [--stash] [--json]",
    flags: [
      { flag: "--yes, -y", description: "Skip the confirmation prompt (required for non-interactive use)." },
      { flag: "--dry-run", description: "Print the upgrade plan without touching the working tree." },
      { flag: "--tag <vX.Y.Z>", description: "Target a specific tag instead of the latest known." },
      { flag: "--strategy <mode>", description: "`fetch-checkout` (default) or `source-release` for repos with a release-self flow." },
      { flag: "--stash", description: "`git stash` local changes before checkout, restore on success." },
      { flag: "--json", description: "Emit a structured result object (schema: spec/08-json-schemas/update-apply.schema.json)." },
    ],
    examples: [
      { command: "gitmap update apply acme/api", description: "Interactive upgrade, prompts before running." },
      { command: "gitmap ua acme/api -y --stash", description: "Non-interactive upgrade, stash local edits first." },
      { command: "gitmap ua acme/api --dry-run", description: "Preview the fetch + checkout plan without changes." },
    ],
    seeAlso: [
      { name: "list --update", description: "Discover upgradable repos first" },
      { name: "update all", description: "Upgrade every upgradable repo in one command" },
      { name: "do-pending", description: "Retry Upgrade tasks that failed" },
    ],
  },
  {
    category: "tools",
    name: "update all", alias: "uall", description: "Upgrade every repo returned by `list --update`. Runs `update apply` per row through a worker pool, continues past per-repo failures, and enqueues an Upgrade PendingTask for each failure.",
    usage: "gitmap update all [--yes] [--dry-run] [--only release|import|git-tag] [--parallel N] [--strategy fetch-checkout|source-release] [--stash] [--stop-on-error] [--json]",
    flags: [
      { flag: "--yes, -y", description: "Skip the single top-level confirmation prompt." },
      { flag: "--dry-run", description: "Print the full plan and exit; nothing is fetched or checked out." },
      { flag: "--only <kind>", description: "Restrict to one source: release, import, or git-tag." },
      { flag: "--parallel <n>", description: "Concurrent workers, clamped to [1, 16]. Default 4." },
      { flag: "--strategy <mode>", description: "Passed through to each `update apply` call." },
      { flag: "--stash", description: "Passed through to each `update apply` call." },
      { flag: "--stop-on-error", description: "Abort the remaining queue after the first failure." },
      { flag: "--json", description: "Emit a JSON summary envelope wrapping per-repo results." },
    ],
    examples: [
      { command: "gitmap update all", description: "Prompted upgrade of every upgradable repo." },
      { command: "gitmap uall --only release --dry-run", description: "Preview only release-tracked upgrades." },
      { command: "gitmap uall -y --parallel 8 --json", description: "Batch mode: 8 concurrent workers, JSON summary." },
    ],
    seeAlso: [
      { name: "list --update", description: "Preview what would be upgraded" },
      { name: "update apply", description: "Upgrade one repo at a time" },
      { name: "pending", description: "Inspect failed upgrades" },
      { name: "do-pending", description: "Retry failed upgrades" },
    ],
  },
  {
    category: "scanning",
    name: "hd", alias: "help-dashboard", description: "One-screen operator dashboard: pending tasks, upgradable repos, last scan, current gitmap version, and the tail of completed tasks. Read-only alias for `gitmap help --dashboard`.",
    usage: "gitmap hd [--json] [--refresh] [--tail N]",
    flags: [
      { flag: "--json", description: "Emit the dashboard payload as JSON (schema: spec/08-json-schemas/hd.schema.json)." },
      { flag: "--refresh", description: "Re-probe upgrades before rendering (equivalent to `list --update --stale-days 0`)." },
      { flag: "--tail <n>", description: "How many recent completed tasks to show. Default 5." },
    ],
    examples: [
      { command: "gitmap hd", description: "Render the boxed dashboard for the current install." },
      { command: "gitmap hd --refresh", description: "Force a fresh probe pass before rendering." },
      { command: "gitmap hd --json --tail 20", description: "Structured payload with the last 20 completed tasks." },
    ],
    seeAlso: [
      { name: "stats", description: "Deeper aggregate metrics, now including the Upgrades block" },
      { name: "list --update", description: "Full upgradable list" },
      { name: "pending", description: "Full pending queue" },
      { name: "history", description: "Complete command history" },
    ],
  },
  {
    category: "tools",
    name: "fix-repo", alias: "fr", description: "Rewrite prior {base}-vN versioned-repo-name tokens in every git-tracked text file to the current version. Go-native re-implementation of fix-repo.ps1/.sh with byte-identical exit codes and config schema.",
    usage: "gitmap fix-repo [-2 | -3 | -5 | --all] [--dry-run] [--verbose] [--config <path>]",
    flags: [
      { flag: "-2 (default)", description: "Rewrite the last 2 prior versions (v(K-2)..v(K-1) → vK)" },
      { flag: "-3", description: "Widen window to last 3 prior versions" },
      { flag: "-5", description: "Widen window to last 5 prior versions" },
      { flag: "--all", description: "Rewrite every prior version v1..v(K-1) → vK" },
      { flag: "--dry-run", description: "Preview only; no file is written. PowerShell alias: -DryRun" },
      { flag: "--verbose", description: "Print every modified file with replacement count. PowerShell alias: -Verbose" },
      { flag: "--config <path>", description: "Override fix-repo.config.json location. PowerShell alias: -Config <path>" },
    ],
    examples: [
      { command: "gitmap fix-repo", description: "Default: rewrite last 2 prior versions, write changes" },
      { command: "gitmap fr -3 --dry-run", description: "Preview a 3-version-wide rewrite without touching disk" },
      { command: "gitmap fix-repo --all --verbose", description: "Rewrite every prior version, log each modified file" },
      { command: "gitmap fr --config ./fix-repo.config.json", description: "Use a custom ignore-rules config file" },
    ],
    seeAlso: [
      { name: "replace", description: "Broader find/replace (literal or version) across all text files" },
      { name: "pull-release", description: "Pull-then-release the current repo after fixing tokens" },
      { name: "clone-next", description: "Clone the next versioned iteration of the current repo" },
    ],
  },
  {
    category: "tools",
    name: "make-public", description: "Make the current repo public on GitHub or GitLab. Wraps the host CLI (gh or glab) so authentication uses your existing `gh auth login` / `glab auth login` session — no tokens stored in gitmap. Prompts for explicit `yes` confirmation before exposing a private repo (skip with --yes for CI). Idempotent: exits 0 if the repo is already public.",
    usage: "gitmap make-public [--yes] [--dry-run] [--verbose]",
    flags: [
      { flag: "--yes, -y", description: "Skip the private→public confirmation prompt (required for non-interactive use)." },
      { flag: "--dry-run", description: "Print the provider command that would run; do not invoke it." },
      { flag: "--verbose", description: "Echo every shell command to stderr before running it." },
    ],
    examples: [
      { command: "gitmap make-public", description: "Interactive — prompts for `yes` before exposing the repo" },
      { command: "gitmap make-public --yes", description: "Non-interactive (CI / scripts)" },
      { command: "gitmap make-public --dry-run", description: "Preview without touching the API" },
      { command: "gitmap make-public --yes --verbose", description: "Debug auth or argv issues" },
    ],
    seeAlso: [
      { name: "make-private", description: "The opposite direction — hide a public repo (no confirmation needed)" },
    ],
  },
  {
    category: "tools",
    name: "make-private", description: "Make the current repo private on GitHub or GitLab. Wraps the host CLI (gh or glab) so authentication uses your existing `gh auth login` / `glab auth login` session — no tokens stored in gitmap. No confirmation prompt: hiding a public repo is the safe direction. Idempotent: exits 0 if the repo is already private.",
    usage: "gitmap make-private [--dry-run] [--verbose]",
    flags: [
      { flag: "--dry-run", description: "Print the provider command that would run; do not invoke it." },
      { flag: "--verbose", description: "Echo every shell command to stderr before running it." },
      { flag: "--yes, -y", description: "Accepted but ignored — kept for symmetry with make-public so the same script works for both." },
    ],
    examples: [
      { command: "gitmap make-private", description: "Standard — flips visibility to private and verifies" },
      { command: "gitmap make-private --dry-run", description: "Preview without touching the API" },
      { command: "gitmap make-private --verbose", description: "Debug auth or argv issues" },
    ],
    seeAlso: [
      { name: "make-public", description: "The opposite direction — expose a private repo (with confirmation)" },
    ],
  },
  {
    category: "tools",
    name: "make-all-public", alias: "MAPUB", description: "Bulk-flip every matching repo on an owner to public in one call. Resolves the owner from `<owner-or-url>`, lists every repo under the owner (≤1000 per call) via `gh`/`glab`, matches comma-separated wildcard patterns (`prefix*`, `*contains*`, `prefix*suffix`, `!negation`) and prints a numbered match table. Prompts for confirmation (`y`/`n` or exclusion expressions like `1,3-5`) unless `-Y` is passed. Each flip is persisted in the audit DB so `vu`/`vr`/`vish` can replay or inspect the run later. Idempotent: already-public repos are skipped.",
    usage: "gitmap make-all-public <owner-or-url> <patterns> [-Y|--yes] [--verbose]",
    flags: [
      { flag: "-Y, --yes", description: "Skip the interactive confirmation prompt." },
      { flag: "--verbose", description: "Echo every shell command to stderr before running it." },
    ],
    examples: [
      { command: "gitmap make-all-public alice \"demo-*\"", description: "Flip every `demo-*` repo on owner `alice` to public" },
      { command: "gitmap make-all-public alice \"demo-*,proto-*,!proto-secret\"", description: "Multi-pattern with negation exclusion" },
      { command: "gitmap MAPUB https://github.com/alice \"*\" -Y --verbose", description: "Shorthand alias, every repo, no prompt" },
    ],
    seeAlso: [
      { name: "make-all-private", description: "Opposite direction, same machinery" },
      { name: "visibility-undo", description: "Reverse the most recent run (alias `vu`)" },
      { name: "visibility-history", description: "List past runs and their IDs (alias `vish`)" },
    ],
  },
  {
    category: "tools",
    name: "make-all-private", alias: "MAPRI", description: "Bulk-flip every matching repo on an owner to private in one call. Mirror image of `make-all-public` — same flags, exit codes, audit trail, and undo/redo/history wiring. Resolves the owner, lists every repo under it, matches comma-separated wildcard patterns (`prefix*`, `*contains*`, `prefix*suffix`, `!negation`), prints a numbered match table, and prompts for confirmation (`y`/`n` or exclusion ranges like `1,3-5`) unless `-Y`. Idempotent: already-private repos are skipped.",
    usage: "gitmap make-all-private <owner-or-url> <patterns> [-Y|--yes] [--verbose]",
    flags: [
      { flag: "-Y, --yes", description: "Skip the interactive confirmation prompt." },
      { flag: "--verbose", description: "Echo every shell command to stderr before running it." },
    ],
    examples: [
      { command: "gitmap make-all-private alice \"archive-*\"", description: "Hide every `archive-*` repo on owner `alice`" },
      { command: "gitmap make-all-private alice \"old-*,!old-keep\" -Y", description: "Negation excludes `old-keep`, no prompt" },
      { command: "gitmap MAPRI https://gitlab.com/alice \"*\" --verbose", description: "Shorthand alias, every repo, debug shell" },
    ],
    seeAlso: [
      { name: "make-all-public", description: "Opposite direction, same machinery" },
      { name: "visibility-undo", description: "Reverse the most recent run (alias `vu`)" },
      { name: "visibility-history", description: "List past runs and their IDs (alias `vish`)" },
    ],
  },
  {
    category: "tools",
    name: "clone-fix-repo", alias: "cfr", description: "Clone a repo, then immediately run `fix-repo --all` inside the new folder. One-shot replacement for the manual sequence `gitmap clone <url>` → `cd <folder>` → `gitmap fix-repo --all`. Versioned URLs auto-flatten to the base name (e.g. `myrepo-v13.git` → `myrepo/`). Order-independent pre-URL modifiers `cg` (install Coding Guidelines v24 + auto-commit/push) and `p` (promote to public) may appear in any order before the URL. Internally re-execs the gitmap binary for each step so exit codes propagate as-is; the pipeline halts on first failure.",
    usage: "gitmap clone-fix-repo [cg] [p] <url> [folder] [flags]",
    flags: [
      { flag: "--ssh", description: "Force SSH transport before clone (converts https:// URLs)" },
      { flag: "--https", description: "Force HTTPS transport before clone (converts SSH URLs)" },
      { flag: "--no-vscode-sync", description: "Skip writing the resolved folder into VS Code Project Manager" },
      { flag: "--require-version", description: "Fail (exit 4) when the cloned repo has no `-vN` suffix instead of skipping fix-repo" },
      { flag: "--dry-run", description: "Preview only, no clone or chdir" },
      { flag: "--no-commit", description: "(cg only) Stage guideline files but do not commit" },
      { flag: "--no-push", description: "(cg only) Commit locally but do not push (also implicit when no upstream)" },
    ],
    examples: [
      { command: "gitmap clone-fix-repo https://github.com/acme/myrepo-v13.git", description: "Clone (auto-flatten to `myrepo/`) and rewrite all prior `myrepo-vN` tokens" },
      { command: "gitmap cfr git@github.com:acme/myrepo-v13.git myrepo-fresh", description: "SSH clone into an explicit folder, then fix-repo" },
      { command: "gitmap cfr cg https://github.com/acme/myrepo-v13.git", description: "Clone + fix + install Coding Guidelines v24 (auto-commit + push)" },
      { command: "gitmap cfr cg https://github.com/acme/myrepo-v13.git --no-commit", description: "Clone + fix + install guidelines; leave files staged for review" },
      { command: "gitmap cfr cg https://github.com/acme/myrepo-v13.git --no-push", description: "Clone + fix + install + commit; skip the git push step" },
      { command: "gitmap cfr p cg https://github.com/acme/myrepo-v13.git", description: "Order-independent: clone + fix + make-public + install guidelines" },
    ],
    seeAlso: [
      { name: "clone-fix-repo-pub", description: "Same pipeline plus `make-public --yes` at the end" },
      { name: "clone", description: "The underlying clone step on its own" },
      { name: "fix-repo", description: "The underlying rewrite step on its own" },
    ],
  },
  {
    category: "tools",
    name: "clone-fix-repo-pub", alias: "cfrp", description: "Clone a repo, run `fix-repo --all`, then flip its visibility to public on GitHub or GitLab, all in one shot. Equivalent to `gitmap clone <url>` → `cd <folder>` → `gitmap fix-repo --all` → `gitmap make-public --yes`. Requires `gh` or `glab` installed and authenticated. Accepts the same pre-URL `cg` modifier as `cfr` to also install Coding Guidelines v24 (with auto-commit + push). Each step's exit code is propagated; the pipeline halts on the first non-zero exit.",
    usage: "gitmap clone-fix-repo-pub [cg] <url> [folder] [flags]",
    flags: [
      { flag: "--ssh", description: "Force SSH transport before clone" },
      { flag: "--https", description: "Force HTTPS transport before clone" },
      { flag: "--no-commit", description: "(cg only) Stage guideline files but do not commit" },
      { flag: "--no-push", description: "(cg only) Commit locally but do not push" },
    ],
    examples: [
      { command: "gitmap clone-fix-repo-pub https://github.com/acme/myrepo-v13.git", description: "Clone, fix prior version tokens, then expose publicly" },
      { command: "gitmap cfrp git@github.com:acme/myrepo-v13.git myrepo-fresh", description: "Same pipeline with explicit destination folder" },
      { command: "gitmap cfrp cg https://github.com/acme/myrepo-v13.git", description: "Clone + fix + make-public + install Coding Guidelines v24 (auto-commit + push)" },
      { command: "gitmap cfrp cg https://github.com/acme/myrepo-v13.git --no-push", description: "Same, but commit guidelines locally without pushing" },
    ],
    seeAlso: [
      { name: "clone-fix-repo", description: "Same pipeline without the visibility flip" },
      { name: "make-public", description: "The visibility step on its own (with confirmation)" },
      { name: "fix-repo", description: "The underlying rewrite step on its own" },
    ],
  },

  {
    category: "changelog",
    name: "changelog", alias: "cl", description: "View release notes from CHANGELOG.md with filtering and version lookup",
    usage: "gitmap changelog [version] [--latest] [--limit N] [--open] [--source <type>]",
    flags: [
      { flag: "--latest", description: "Show only the most recent version" },
      { flag: "--limit <n>", description: "Max number of versions to display (default: 5)" },
      { flag: "--open", description: "Open CHANGELOG.md in default application" },
      { flag: "--source <type>", description: "Filter by source: release or import" },
    ],
    examples: [
      { command: "gitmap changelog", description: "Show last 5 versions" },
      { command: "gitmap cl --latest", description: "Most recent version only" },
      { command: "gitmap changelog v2.3.0", description: "Notes for a specific version" },
      { command: "gitmap cl --source release --limit 10", description: "Last 10 release-sourced entries" },
      { command: "gitmap cl --open", description: "Open CHANGELOG.md in your editor" },
    ],
    seeAlso: [
      { name: "release", description: "Create a release", url: "/release" },
      { name: "list-versions", description: "List available tags" },
      { name: "list-releases", description: "List stored release metadata" },
    ],
  },
  {
    category: "changelog",
    name: "changelog-generate", alias: "cg", description: "Auto-generate CHANGELOG.md entries from commits between Git tags",
    usage: "gitmap changelog-generate [--from <tag>] [--to <tag>]",
    flags: [
      { flag: "--from <tag>", description: "Start tag (default: second-latest)" },
      { flag: "--to <tag>", description: "End tag (default: latest)" },
    ],
    examples: [
      { command: "gitmap changelog-generate", description: "Generate entries between last two tags" },
      { command: "gitmap cg --from v2.40.0 --to v2.45.0", description: "Generate entries for a specific range" },
      { command: "gitmap cg --from v1.0.0", description: "Everything from v1.0.0 to latest" },
    ],
    seeAlso: [
      { name: "changelog", description: "View generated changelog" },
      { name: "release", description: "Create a release" },
      { name: "list-versions", description: "List available tags" },
    ],
  },
  {
    category: "changelog",
    name: "list-versions", alias: "lv", description: "List all Git release tags with optional notes from CHANGELOG.md",
    usage: "gitmap list-versions [--json] [--limit N]",
    flags: [
      { flag: "--json", description: "Output as structured JSON" },
      { flag: "--limit N", description: "Show only the top N versions (0 = all)" },
    ],
    examples: [
      { command: "gitmap list-versions", description: "List all release tags" },
      { command: "gitmap lv --limit 10", description: "Last 10 versions" },
      { command: "gitmap lv --json", description: "JSON output for scripting" },
    ],
    seeAlso: [
      { name: "list-releases", description: "List stored release metadata" },
      { name: "changelog", description: "View release notes" },
      { name: "release", description: "Create a release", url: "/release" },
      { name: "revert", description: "Revert to a specific version" },
    ],
  },
  {
    category: "changelog",
    name: "list-releases", alias: "lr", description: "List release metadata records stored in the database",
    usage: "gitmap list-releases [--json] [--source manual|scan]",
    flags: [
      { flag: "--json", description: "Output as structured JSON" },
      { flag: "--source <type>", description: "Filter by release source" },
    ],
    examples: [
      { command: "gitmap list-releases", description: "List all stored releases" },
      { command: "gitmap lr --json", description: "JSON output" },
      { command: "gitmap lr --source manual", description: "Only manually created releases" },
    ],
    seeAlso: [
      { name: "list-versions", description: "List Git tags" },
      { name: "release", description: "Create a release", url: "/release" },
      { name: "changelog", description: "View release notes" },
    ],
  },
  {
    category: "changelog",
    name: "clear-release-json", alias: "crj", description: "Remove a .gitmap/release/vX.Y.Z.json metadata file",
    usage: "gitmap clear-release-json <version> [--dry-run]",
    flags: [
      { flag: "--dry-run", description: "Preview which file would be removed without deleting" },
    ],
    examples: [
      { command: "gitmap clear-release-json v2.20.0", description: "Remove v2.20.0 release metadata" },
      { command: "gitmap crj v1.0.0 --dry-run", description: "Preview removal without deleting" },
      { command: "gitmap crj v1.0.0", description: "Remove using alias" },
    ],
    seeAlso: [
      { name: "release", description: "Create a release", url: "/release" },
      { name: "list-releases", description: "Show stored releases" },
      { name: "db-reset", description: "Reset entire database" },
      { name: "Spec: clear-release-json", description: "Full specification", url: "/clear-release-json" },
    ],
  },
  {
    category: "changelog",
    name: "revert", alias: undefined, description: "Revert to a specific release version by checking out its tag",
    usage: "gitmap revert <version>",
    flags: [
      { flag: "<version>", description: "Release tag to revert to (auto-prefixed with 'v' if missing)" },
    ],
    examples: [
      { command: "gitmap revert v2.9.0", description: "Checkout tag v2.9.0" },
      { command: "gitmap revert 2.8.0", description: "Version auto-prefixed with v" },
    ],
    seeAlso: [
      { name: "list-versions", description: "List available versions" },
      { name: "update", description: "Self-update to latest" },
      { name: "update-cleanup", description: "Remove artifacts after revert" },
      { name: "changelog", description: "View release notes" },
    ],
  },

  // ═══════════════════════════════════════════
  // Navigation & Groups
  // ═══════════════════════════════════════════
  {
    category: "navigation",
    name: "cd", alias: "go", description: "Navigate your shell to a tracked repo directory (supports interactive picker)",
    usage: "gitmap cd <repo-name|repos> [--group <name>] [--pick]",
    examples: [
      { command: "gitmap cd myrepo", description: "Jump to myrepo's directory" },
      { command: "gitmap cd repos", description: "Interactive repo picker (fzf-style)" },
      { command: "gitmap cd repos --group work", description: "Pick from work group only" },
      { command: "gitmap go myrepo", description: "Alias shorthand" },
    ],
    seeAlso: [
      { name: "list", description: "List all tracked repos with slugs" },
      { name: "scan", description: "Scan directories to populate database" },
      { name: "group", description: "Manage repo groups" },
      { name: "bookmark", description: "Save commands for re-execution", url: "/bookmarks" },
    ],
  },
  {
    category: "navigation",
    name: "list", alias: "ls", description: "Show all tracked repos with slugs, supports filtering by project type or group",
    usage: "gitmap list [--group <name>] [--verbose]\ngitmap ls go|node|react|cpp|csharp\ngitmap ls groups",
    flags: [
      { flag: "--group <name>", description: "Filter by group name" },
      { flag: "--verbose", description: "Show full paths and URLs" },
    ],
    examples: [
      { command: "gitmap list", description: "List all tracked repos" },
      { command: "gitmap ls go", description: "List only Go projects" },
      { command: "gitmap ls node", description: "List only Node.js projects" },
      { command: "gitmap ls react", description: "List only React projects" },
      { command: "gitmap ls groups", description: "List all defined groups" },
      { command: "gitmap ls --group backend --verbose", description: "Verbose list for backend group" },
    ],
    seeAlso: [
      { name: "cd", description: "Navigate to a tracked repo" },
      { name: "group", description: "Manage and activate repo groups" },
      { name: "multi-group", description: "Select multiple groups" },
      { name: "status", description: "View repo statuses" },
    ],
  },
  {
    category: "navigation",
    name: "group", alias: "g", description: "Create, manage, and activate repo groups for batch operations",
    usage: "gitmap group <create|add|remove|list|show|delete|pull|status|exec|clear> [args]\ngitmap g <name>    Activate a group\ngitmap g           Show active group",
    flags: [
      { flag: "--description <text>", description: "Group description (for create)" },
      { flag: "--color <name>", description: "Terminal color for the group (for create)" },
    ],
    examples: [
      { command: "gitmap group create backend --description \"Backend services\"", description: "Create a named group" },
      { command: "gitmap group add backend my-api my-worker", description: "Add repos to backend group" },
      { command: "gitmap g backend", description: "Activate backend as the current group" },
      { command: "gitmap g pull", description: "Pull all repos in the active group" },
      { command: "gitmap g status", description: "Status for active group" },
      { command: "gitmap g exec fetch --prune", description: "Run git fetch across active group" },
      { command: "gitmap group list", description: "Show all groups" },
      { command: "gitmap g clear", description: "Clear active group" },
    ],
    seeAlso: [
      { name: "multi-group", description: "Select multiple groups" },
      { name: "list", description: "List all tracked repos" },
      { name: "pull", description: "Pull repos by group" },
      { name: "status", description: "Filter status by group" },
    ],
  },
  {
    category: "navigation",
    name: "multi-group", alias: "mg", description: "Select and operate on multiple groups at once",
    usage: "gitmap multi-group <group1,group2,...|clear|pull|status|exec>",
    examples: [
      { command: "gitmap mg backend,frontend", description: "Select multiple groups" },
      { command: "gitmap mg pull", description: "Pull repos from all selected groups" },
      { command: "gitmap mg status", description: "Status for all selected groups" },
      { command: "gitmap mg exec fetch --prune", description: "Fetch across all selected groups" },
      { command: "gitmap mg clear", description: "Clear multi-group selection" },
    ],
    seeAlso: [
      { name: "group", description: "Manage and activate single groups" },
      { name: "pull", description: "Pull repos" },
      { name: "status", description: "View repo statuses" },
    ],
  },
  {
    category: "navigation",
    name: "as", alias: "s-alias", description: "Register the current Git repo in SQLite + map a short name to it (run from inside the repo). Mirrors the alias to VS Code Project Manager projects.json when present.",
    usage: "gitmap as [alias-name] [--force]",
    flags: [
      { flag: "--force, -f", description: "Overwrite an existing alias that points to a different repo" },
    ],
    examples: [
      { command: "gitmap as", description: "Use the repo folder basename as the alias" },
      { command: "gitmap as backend", description: "Register as 'backend' (also renames the matching projects.json entry)" },
      { command: "gitmap as backend -f", description: "Overwrite an existing 'backend' alias" },
    ],
    seeAlso: [
      { name: "alias list", description: "Show every registered alias" },
      { name: "release-alias", description: "Release this repo by alias from anywhere" },
      { name: "scan", description: "Bulk-discover and register repos under a directory" },
      { name: "code", description: "Register a path with VS Code Project Manager and open it" },
    ],
  },
  {
    category: "navigation",
    name: "code", description: "Register the current repo (or any path) with the alefragnani.project-manager VS Code extension and open VS Code on it. Supports multi-root extras (v3.39.0+) and auto-derived tags (v3.40.0+).",
    usage: "gitmap code [alias] [path] [extraPath...]\ngitmap code paths add|rm|list <alias> [path]",
    examples: [
      { command: "gitmap code", description: "Register the git repo root (or CWD) — alias defaults to folder basename, tags auto-detected" },
      { command: "gitmap code backend", description: "Override the alias to 'backend' for the resolved path" },
      { command: "gitmap code docs ~/Documents/spec", description: "Register any path (no git requirement) with alias 'docs'" },
      { command: "gitmap code mono ~/work/main ~/work/main/frontend ~/work/main/backend", description: "Register a multi-root entry: root + variadic extras (additive, never clobbers UI-added paths)" },
      { command: "gitmap code paths add mono ~/work/main/scripts", description: "Attach an extra folder to an existing entry" },
      { command: "gitmap code paths list mono", description: "Show rootPath + every attached extra path" },
      { command: "gitmap code paths rm mono ~/work/main/scripts", description: "Detach an extra folder (overwrites — actually sticks across re-syncs)" },
    ],
    seeAlso: [
      { name: "as", description: "Register an alias for the current repo (mirrors to projects.json too)" },
      { name: "scan", description: "Bulk-syncs every discovered repo into projects.json" },
      { name: "cd", description: "Jump to a tracked repo directory in your shell" },
    ],
  },

  // ═══════════════════════════════════════════
  // History & Analytics
  // ═══════════════════════════════════════════
  {
    category: "history",
    name: "history", alias: "hi", description: "Browse the full CLI command execution history with timestamps and durations",
    usage: "gitmap history [--limit N] [--json]",
    flags: [
      { flag: "--limit N", description: "Number of entries to show" },
      { flag: "--json", description: "Output as structured JSON" },
    ],
    examples: [
      { command: "gitmap history", description: "Show recent command history" },
      { command: "gitmap hi --limit 20", description: "Last 20 commands" },
      { command: "gitmap hi --json", description: "JSON output for scripting" },
    ],
    seeAlso: [
      { name: "Spec: history", description: "Full history documentation", url: "/history" },
      { name: "history-reset", description: "Clear command history" },
      { name: "stats", description: "View aggregated usage metrics", url: "/stats" },
      { name: "bookmark", description: "Save commands for re-execution", url: "/bookmarks" },
    ],
  },
  {
    category: "history",
    name: "history-reset", alias: "hr", description: "Clear all command execution history (requires --confirm)",
    usage: "gitmap history-reset --confirm",
    flags: [
      { flag: "--confirm", description: "Required flag to confirm destructive reset" },
    ],
    examples: [
      { command: "gitmap history-reset --confirm", description: "Clear all command history" },
      { command: "gitmap hr --confirm", description: "Alias shorthand" },
    ],
    seeAlso: [
      { name: "history", description: "View command history", url: "/history" },
      { name: "db-reset", description: "Reset entire database" },
    ],
  },
  {
    category: "history",
    name: "stats", alias: "ss", description: "Show aggregated usage and performance metrics for all commands",
    usage: "gitmap stats [--command <name>] [--json]",
    flags: [
      { flag: "--command <name>", description: "Show stats for a specific command only" },
      { flag: "--json", description: "Output as JSON" },
    ],
    examples: [
      { command: "gitmap stats", description: "Show usage stats for all commands" },
      { command: "gitmap stats --command scan", description: "Stats for scan only" },
      { command: "gitmap ss --json", description: "JSON output for dashboards" },
      { command: "gitmap stats --command release", description: "How many releases have you done?" },
    ],
    seeAlso: [
      { name: "Spec: stats", description: "Full stats documentation", url: "/stats" },
      { name: "history", description: "View command history", url: "/history" },
      { name: "dashboard", description: "Interactive HTML dashboard" },
    ],
  },
  {
    category: "history",
    name: "dashboard", alias: "db", description: "Generate an interactive HTML dashboard with charts, tables, and heatmap for a repo",
    usage: "gitmap dashboard [flags]",
    flags: [
      { flag: "--limit <n>", description: "Maximum number of commits to include" },
      { flag: "--since <date>", description: "Only include commits after this date (YYYY-MM-DD)" },
      { flag: "--no-merges", description: "Exclude merge commits from the output" },
      { flag: "--out-dir <path>", description: "Output directory for dashboard files" },
      { flag: "--open", description: "Open the generated dashboard in the default browser" },
    ],
    examples: [
      { command: "gitmap dashboard --open", description: "Generate and open dashboard in browser" },
      { command: "gitmap db --limit 100 --open", description: "Last 100 commits, open immediately" },
      { command: "gitmap dashboard --since 2025-01-01 --no-merges", description: "2025 commits, no merges" },
      { command: "gitmap db --out-dir ./reports", description: "Save dashboard to custom directory" },
    ],
    seeAlso: [
      { name: "stats", description: "Aggregated command usage statistics" },
      { name: "history", description: "Command execution history" },
      { name: "dashboard", description: "Dedicated docs page", url: "/dashboard" },
    ],
  },
  {
    category: "history",
    name: "amend", alias: "am", description: "Rewrite commit author info (name/email) across a branch",
    usage: "gitmap amend [commit-hash] --name <name> --email <email> [--branch <branch>]",
    flags: [
      { flag: "--name <name>", description: "New author name for rewritten commits" },
      { flag: "--email <email>", description: "New author email for rewritten commits" },
      { flag: "--branch <branch>", description: "Target branch (default: current branch)" },
    ],
    examples: [
      { command: "gitmap amend --name \"John\" --email \"john@example.com\"", description: "Rewrite all commits on current branch" },
      { command: "gitmap amend abc123 --name \"John\" --email \"john@example.com\"", description: "Rewrite from commit abc123 onwards" },
      { command: "gitmap amend --name \"Bot\" --email \"bot@ci.com\" --branch main", description: "Rewrite all commits on main" },
    ],
    seeAlso: [
      { name: "amend-list", description: "List previous amendments" },
      { name: "history", description: "View command history", url: "/history" },
    ],
  },
  {
    category: "history",
    name: "amend-list", alias: "al", description: "List all previous author amendment audit records",
    usage: "gitmap amend-list [--json] [--limit <n>]",
    flags: [
      { flag: "--json", description: "Output in JSON format" },
      { flag: "--limit <n>", description: "Limit number of results" },
    ],
    examples: [
      { command: "gitmap amend-list", description: "Show all amendment records" },
      { command: "gitmap amend-list --json", description: "JSON output" },
      { command: "gitmap al --limit 5", description: "Last 5 amendments" },
    ],
    seeAlso: [
      { name: "amend", description: "Rewrite commit author info" },
      { name: "history", description: "View command history", url: "/history" },
    ],
  },
  {
    category: "history",
    name: "version-history", alias: "vh", description: "Show version transitions for the current repo (e.g. v11 → v12 → v13) recorded by clone-next flatten",
    usage: "gitmap version-history [--limit N] [--json]",
    flags: [
      { flag: "--limit <n>", description: "Show only the most recent N transitions" },
      { flag: "--json", description: "Emit machine-readable JSON instead of the human table" },
    ],
    examples: [
      { command: "gitmap version-history", description: "Print every recorded version transition for this repo" },
      { command: "gitmap vh --limit 5", description: "Last 5 transitions only" },
      { command: "gitmap vh --json", description: "JSON output for scripts and dashboards" },
    ],
    seeAlso: [
      { name: "clone-next", description: "Source of new version transitions (flatten-by-default since v2.75.0)" },
      { name: "Spec: version-history", description: "Full version-history docs page", url: "/version-history" },
      { name: "history", description: "Generic command-execution history", url: "/history" },
    ],
  },

  // ═══════════════════════════════════════════
  // Project Detection
  // ═══════════════════════════════════════════
  {
    category: "detection",
    name: "go-repos", alias: "gr", description: "List all detected Go projects (repos with go.mod)",
    usage: "gitmap go-repos [--json]",
    examples: [
      { command: "gitmap go-repos", description: "List all Go projects" },
      { command: "gitmap gr --json", description: "JSON output" },
    ],
    seeAlso: [
      { name: "node-repos", description: "List Node.js projects" },
      { name: "react-repos", description: "List React projects" },
      { name: "scan", description: "Scan directories first" },
      { name: "gomod", description: "Rename Go module paths", url: "/gomod" },
    ],
  },
  {
    category: "detection",
    name: "node-repos", alias: "nr", description: "List all detected Node.js projects (repos with package.json)",
    usage: "gitmap node-repos [--json]",
    examples: [
      { command: "gitmap node-repos", description: "List all Node.js projects" },
      { command: "gitmap nr --json", description: "JSON output" },
    ],
    seeAlso: [
      { name: "react-repos", description: "List React projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },
  {
    category: "detection",
    name: "react-repos", alias: "rr", description: "List all detected React projects (Node.js repos with react dependency)",
    usage: "gitmap react-repos [--json]",
    examples: [
      { command: "gitmap react-repos", description: "List all React projects" },
      { command: "gitmap rr --json", description: "JSON output" },
    ],
    seeAlso: [
      { name: "node-repos", description: "List Node.js projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },
  {
    category: "detection",
    name: "cpp-repos", alias: "cr", description: "List all detected C++ projects (repos with CMakeLists.txt or .vcxproj)",
    usage: "gitmap cpp-repos [--json]",
    examples: [
      { command: "gitmap cpp-repos", description: "List all C++ projects" },
      { command: "gitmap cr --json", description: "JSON output" },
    ],
    seeAlso: [
      { name: "csharp-repos", description: "List C# projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },
  {
    category: "detection",
    name: "csharp-repos", alias: "csr", description: "List all detected C# projects (repos with .csproj or .sln)",
    usage: "gitmap csharp-repos [--json]",
    examples: [
      { command: "gitmap csharp-repos", description: "List all C# projects" },
      { command: "gitmap csr --json", description: "JSON output" },
    ],
    seeAlso: [
      { name: "cpp-repos", description: "List C++ projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },

  // ═══════════════════════════════════════════
  // Data & Profiles
  // ═══════════════════════════════════════════
  {
    category: "data",
    name: "export", alias: "ex", description: "Export the full database to a portable JSON file",
    usage: "gitmap export [file]",
    examples: [
      { command: "gitmap export", description: "Export to default gitmap-export.json" },
      { command: "gitmap ex backup.json", description: "Export to custom filename" },
      { command: "gitmap ex ~/Desktop/repos-backup.json", description: "Export to specific path" },
    ],
    seeAlso: [
      { name: "Spec: export", description: "Full export specification", url: "/export" },
      { name: "import", description: "Import repos from file", url: "/import" },
      { name: "profile", description: "Manage database profiles" },
      { name: "scan", description: "Scan directories to populate data" },
    ],
  },
  {
    category: "data",
    name: "import", alias: "im", description: "Import repos from a previously exported JSON file",
    usage: "gitmap import [file] --confirm",
    flags: [
      { flag: "--confirm", description: "Confirm the import (required, prevents accidents)" },
    ],
    examples: [
      { command: "gitmap import --confirm", description: "Import from default gitmap-export.json" },
      { command: "gitmap im backup.json --confirm", description: "Import from custom file" },
      { command: "gitmap im ~/Desktop/repos-backup.json --confirm", description: "Import from specific path" },
    ],
    seeAlso: [
      { name: "Spec: import", description: "Full import specification", url: "/import" },
      { name: "export", description: "Export database to file", url: "/export" },
      { name: "scan", description: "Scan directories" },
      { name: "profile", description: "Manage database profiles" },
    ],
  },
  {
    category: "data",
    name: "profile", alias: "pf", description: "Create, switch, and manage isolated database profiles",
    usage: "gitmap profile <create|list|switch|delete|show> [name]",
    examples: [
      { command: "gitmap profile create work", description: "Create a new profile called 'work'" },
      { command: "gitmap pf list", description: "List all available profiles" },
      { command: "gitmap profile switch work", description: "Switch to the work profile" },
      { command: "gitmap profile show", description: "Show the currently active profile" },
      { command: "gitmap pf delete old-profile", description: "Delete a profile" },
    ],
    seeAlso: [
      { name: "Spec: profile", description: "Full profile specification", url: "/profile" },
      { name: "diff-profiles", description: "Compare repos across profiles" },
      { name: "export", description: "Export database", url: "/export" },
      { name: "import", description: "Import repos", url: "/import" },
    ],
  },
  {
    category: "data",
    name: "diff-profiles", alias: "dp", description: "Compare repos between two profiles and show added/removed/changed",
    usage: "gitmap diff-profiles <profileA> <profileB> [--all] [--json]",
    flags: [
      { flag: "--all", description: "Include identical repos in the output" },
      { flag: "--json", description: "Output as structured JSON" },
    ],
    examples: [
      { command: "gitmap diff-profiles default work", description: "Compare default and work profiles" },
      { command: "gitmap dp work personal --json", description: "JSON diff output" },
      { command: "gitmap dp home office --all", description: "Full comparison including identical repos" },
    ],
    seeAlso: [
      { name: "Spec: diff-profiles", description: "Full diff-profiles specification", url: "/diff-profiles" },
      { name: "profile", description: "Manage database profiles", url: "/profile" },
      { name: "list", description: "List tracked repos" },
      { name: "export", description: "Export database", url: "/export" },
    ],
  },
  {
    category: "data",
    name: "bookmark", alias: "bk", description: "Save, list, replay, and delete bookmarked commands",
    usage: "gitmap bookmark <save|list|run|delete> [args]",
    examples: [
      { command: "gitmap bookmark save ssh-scan scan --mode ssh", description: "Save a scan command as 'ssh-scan'" },
      { command: "gitmap bk list", description: "List all saved bookmarks" },
      { command: "gitmap bookmark run ssh-scan", description: "Replay the 'ssh-scan' bookmark" },
      { command: "gitmap bk delete ssh-scan", description: "Remove a bookmark" },
    ],
    seeAlso: [
      { name: "Spec: bookmarks", description: "Full bookmarks documentation", url: "/bookmarks" },
      { name: "history", description: "View command execution history", url: "/history" },
      { name: "scan", description: "Scan directories (common bookmark target)" },
    ],
  },
  {
    category: "data",
    name: "db-reset", alias: undefined, description: "Completely reset the local SQLite database (requires --confirm)",
    usage: "gitmap db-reset --confirm",
    flags: [
      { flag: "--confirm", description: "Required flag to confirm destructive reset" },
    ],
    examples: [
      { command: "gitmap db-reset --confirm", description: "Reset the database" },
    ],
    seeAlso: [
      { name: "history-reset", description: "Clear command history only" },
      { name: "scan", description: "Re-scan after reset" },
      { name: "setup", description: "Re-run setup wizard" },
    ],
  },

  // ═══════════════════════════════════════════
  // Tools & Setup
  // ═══════════════════════════════════════════
  {
    category: "tools",
    name: "add ignore", alias: undefined, description: "Merge a curated .gitignore template into a gitmap-managed marker block inside ./.gitignore. Always merges the `common` template plus any langs you list. Idempotent — re-runs with the same set are a byte-stable no-op.",
    usage: "gitmap add ignore [langs...] [--dry-run]",
    flags: [
      { flag: "--dry-run", description: "Preview the merged .gitignore block without writing anything" },
    ],
    examples: [
      { command: "gitmap add ignore go", description: "New Go repo — writes common + go blocks into ./.gitignore" },
      { command: "gitmap add ignore go node python", description: "Multi-language monorepo — union of common + all three langs (deduped)" },
      { command: "gitmap add ignore rust --dry-run", description: "Preview the exact block that would be written; touches no files" },
      { command: "gitmap add ignore go node && gitmap add ignore go node", description: "Second run reports 'unchanged' — idempotent by design" },
    ],
    seeAlso: [
      { name: "add attributes", description: "Same flow for .gitattributes" },
      { name: "add lfs-install", description: "LFS hooks + curated binary patterns" },
      { name: "templates list", description: "Discover which langs are available (overlay + embed)" },
      { name: "templates diff", description: "Preview drift between on-disk blocks and curated templates" },
    ],
  },
  {
    category: "tools",
    name: "add attributes", alias: undefined, description: "Merge a curated .gitattributes template into a gitmap-managed marker block inside ./.gitattributes. Always merges `common` (text eol=lf normalization) plus any langs you list. Idempotent.",
    usage: "gitmap add attributes [langs...] [--dry-run]",
    flags: [
      { flag: "--dry-run", description: "Preview the merged .gitattributes block without writing anything" },
    ],
    examples: [
      { command: "gitmap add attributes node", description: "TypeScript/Node repo — writes common + node blocks into ./.gitattributes" },
      { command: "gitmap add attributes go && gitmap add lfs-install", description: "Pair with LFS: two separate marker blocks refresh independently" },
      { command: "gitmap add attributes rust --dry-run", description: "Preview the block before committing" },
    ],
    seeAlso: [
      { name: "add ignore", description: "Same flow for .gitignore" },
      { name: "add lfs-install", description: "LFS hooks + curated binary patterns" },
      { name: "templates show", description: "Print the exact curated bytes for one lang" },
    ],
  },
  {
    category: "tools",
    name: "add lfs-install", alias: undefined, description: "Run `git lfs install --local` to wire per-repo LFS hooks, then merge the curated lfs/common .gitattributes block. Idempotent — the marker block interoperates with `templates init --lfs`.",
    usage: "gitmap add lfs-install [--dry-run]",
    flags: [
      { flag: "--dry-run", description: "Preview the merged lfs/common block without writing anything or running git" },
    ],
    examples: [
      { command: "gitmap add lfs-install", description: "One-shot: install LFS hooks + merge the binary-pattern block" },
      { command: "gitmap add lfs-install --dry-run", description: "Preview what would land in .gitattributes; does not shell out to git" },
    ],
    seeAlso: [
      { name: "add attributes", description: "Merge language attributes alongside the lfs/common block" },
      { name: "templates init", description: "Pass --lfs during scaffolding for the same block, no shell-out" },
      { name: "lfs-common", description: "Per-pattern `git lfs track` (no template)" },
    ],
  },
  {
    category: "tools",
    name: "sync", alias: "sy", description: "Union-merge curated defaults for .gitignore, .gitattributes, .prettierignore, and .prettierrc into the current directory without gitmap marker blocks. Existing content is preserved; only missing lines/keys are appended. Use this when the target file is hand-maintained.",
    usage: "gitmap sync <ignore|attributes|prettier-ignore|prettier-rc|all> [--dry-run] [--force]",
    flags: [
      { flag: "--dry-run, -n", description: "Print planned additions without touching disk" },
      { flag: "--force, -f", description: "For prettier-rc only: overwrite conflicting keys instead of preserving them" },
    ],
    examples: [
      { command: "gitmap sync ignore", description: "Append missing curated ignore entries to a hand-written .gitignore" },
      { command: "gitmap sync all --dry-run", description: "Preview additions for all four targets in one pass" },
      { command: "gitmap sync prettier-rc --force", description: "Reset prettier keys (printWidth, semi, etc.) to curated defaults" },
      { command: "gitmap sy attributes", description: "Short alias, union-merges line-ending + LFS defaults" },
    ],
    seeAlso: [
      { name: "add ignore", description: "Marker-block managed .gitignore (rewritten on each run)" },
      { name: "add attributes", description: "Marker-block managed .gitattributes" },
      { name: "templates init", description: "Scaffold from scratch when the file doesn't exist yet" },
    ],
  },
  {
    category: "tools",
    name: "commons", alias: "co", description: "One-shot shortcut for 'sync all'. Adds or dedupe-merges the curated .gitignore, .gitattributes, .prettierignore, .prettierrc baselines and runs 'git lfs install --local' + the lfs/common .gitattributes block. Existing lines/keys are preserved.",
    usage: "gitmap commons [--dry-run] [--force]",
    flags: [
      { flag: "--dry-run, -n", description: "Print planned additions without touching disk" },
      { flag: "--force, -f", description: "For .prettierrc only: overwrite conflicting keys instead of preserving them" },
    ],
    examples: [
      { command: "gitmap commons", description: "Add or merge all curated hygiene files in one pass" },
      { command: "gitmap commons --dry-run", description: "Preview every planned addition without writing" },
      { command: "gitmap co -n", description: "Short alias, dry-run preview" },
      { command: "gitmap commons --force", description: "Also overwrite conflicting .prettierrc keys with curated values" },
    ],
    seeAlso: [
      { name: "sync", description: "Per-target union merge (ignore | attributes | prettier-ignore | prettier-rc | all)" },
      { name: "add ignore", description: "Marker-block managed .gitignore (rewritten on each run)" },
      { name: "templates init", description: "Scaffold from scratch when the files don't exist yet" },
    ],
  },
  {
    category: "tools",
    name: "templates list", alias: "tpl tl", description: "List every available template with its KIND, LANG, SOURCE (user/embed), and PATH",
    usage: "gitmap templates list",
    flags: [],
    examples: [
      { command: "gitmap templates list", description: "Print the full table — embed entries plus any user overlays" },
      { command: "gitmap tpl tl", description: "Same, using the short aliases" },
    ],
    seeAlso: [
      { name: "templates show", description: "Print a single template's bytes to stdout" },
      { name: "add lfs-install", description: "Use the lfs/common template to populate .gitattributes" },
    ],
  },
  {
    category: "tools",
    name: "templates show", alias: "tpl ts", description: "Print one template (overlay > embed) to stdout, audit-trail header included",
    usage: "gitmap templates show <kind> <lang>",
    flags: [],
    examples: [
      { command: "gitmap templates show ignore go", description: "Resolve and print the Go .gitignore template" },
      { command: "gitmap tpl ts attributes common", description: "Same, short aliases — print the common .gitattributes" },
      { command: "gitmap templates show lfs common > .gitattributes.curated", description: "Diff your overlay against the curated embed" },
    ],
    seeAlso: [
      { name: "templates list", description: "Discover what kind/lang pairs are available" },
      { name: "add lfs-install", description: "Apply lfs/common into .gitattributes via marker block" },
    ],
  },
  {
    category: "tools",
    name: "templates init", alias: "tpl ti", description: "Scaffold .gitignore + .gitattributes for one or more languages by merging the embedded templates into the current directory (idempotent, marker-block aware)",
    usage: "gitmap templates init <lang> [<lang>...] [--lfs] [--dry-run] [--force]",
    flags: [
      { flag: "--lfs", description: "Also merge lfs/common.gitattributes into .gitattributes" },
      { flag: "--dry-run", description: "Preview every block that would be written; touch nothing on disk" },
      { flag: "--force", description: "Replace pre-existing .gitignore/.gitattributes outright (discards hand edits OUTSIDE the gitmap marker block)" },
    ],
    examples: [
      { command: "gitmap templates init go", description: "Scaffold ignore + attributes for Go in the current directory" },
      { command: "gitmap templates init go node --lfs", description: "Multi-lang scaffold plus an LFS attributes block" },
      { command: "gitmap tpl ti python --dry-run", description: "Preview what Python scaffolding would write without touching disk" },
    ],
    seeAlso: [
      { name: "templates diff", description: "Preview drift between on-disk blocks and the curated templates" },
      { name: "templates list", description: "Discover what kind/lang pairs are available before running init" },
      { name: "add lfs-install", description: "Apply lfs/common into .gitattributes via marker block (init --lfs is the same operation, no shell-out)" },
    ],
  },
  {
    category: "tools",
    name: "templates diff", alias: "tpl td", description: "Preview what `add ignore`/`add attributes` would change without writing — marker-block aware, exit codes mirror diff(1)",
    usage: "gitmap templates diff [--lang <name>] [--kind ignore|attributes] [--cwd <path>]",
    flags: [
      { flag: "--lang <name>", description: "Limit to one language (default: every resolvable lang)" },
      { flag: "--kind ignore|attributes", description: "Limit to one kind (default: both)" },
      { flag: "--cwd <path>", description: "Run against a different working tree" },
    ],
    examples: [
      { command: "gitmap templates diff --lang go", description: "Show what would change for the Go .gitignore + .gitattributes blocks" },
      { command: "gitmap tpl td --kind ignore --lang python", description: "Compare only the Python .gitignore block" },
      { command: "gitmap templates diff --lang java || gitmap add ignore java", description: "Pre-commit pattern: only run `add` when diff reports drift (exit 1)" },
    ],
    seeAlso: [
      { name: "templates show", description: "Print the curated bytes you would be diffing against" },
      { name: "add ignore", description: "Apply an ignore template into the marker-block region" },
    ],
  },
  {
    category: "tools",
    name: "setup", alias: undefined, description: "Configure Git global settings and install shell tab-completion scripts",
    usage: "gitmap setup [--config <path>] [--dry-run]",
    flags: [
      { flag: "--config <path>", description: "Path to git-setup.json config file" },
      { flag: "--dry-run", description: "Preview changes without applying" },
    ],
    examples: [
      { command: "gitmap setup", description: "Run the interactive setup wizard" },
      { command: "gitmap setup --dry-run", description: "Preview what setup would change" },
      { command: "gitmap setup --config ./custom-setup.json", description: "Use custom config file" },
    ],
    seeAlso: [
      { name: "completion", description: "Generate completion scripts manually" },
      { name: "doctor", description: "Diagnose issues" },
      { name: "scan", description: "Scan directories after setup" },
      { name: "update", description: "Self-update to latest version" },
    ],
  },
  {
    category: "tools",
    name: "doctor", alias: undefined, description: "Run 11 health checks: binary, Git, Go, config, database, lock file, and network",
    usage: "gitmap doctor [--fix-path]",
    flags: [
      { flag: "--fix-path", description: "Attempt to fix PATH issues automatically" },
    ],
    examples: [
      { command: "gitmap doctor", description: "Run all 11 diagnostic checks" },
      { command: "gitmap doctor --fix-path", description: "Diagnose and auto-fix PATH issues" },
    ],
    seeAlso: [
      { name: "setup", description: "Re-run setup wizard" },
      { name: "update", description: "Self-update to latest version" },
      { name: "version", description: "Show current version" },
    ],
  },
  {
    category: "tools",
    name: "update", alias: undefined, description: "Self-update gitmap from its source repo (go install)",
    usage: "gitmap update [--verbose]",
    flags: [
      { flag: "--verbose", description: "Write detailed debug log during update" },
    ],
    examples: [
      { command: "gitmap update", description: "Self-update to latest version" },
      { command: "gitmap update --verbose", description: "Update with detailed logging" },
    ],
    seeAlso: [
      { name: "update-cleanup", description: "Remove update artifacts" },
      { name: "version", description: "Show current version" },
      { name: "doctor", description: "Diagnose issues after update" },
    ],
  },
  {
    category: "tools",
    name: "update-cleanup", alias: undefined, description: "Remove leftover temp binaries and .old backups from previous updates",
    usage: "gitmap update-cleanup",
    examples: [
      { command: "gitmap update-cleanup", description: "Remove temp binaries and .old backups" },
    ],
    seeAlso: [
      { name: "update", description: "Self-update to latest version" },
      { name: "revert", description: "Revert to a previous version" },
    ],
  },
  {
    category: "tools",
    name: "installed-dir", alias: "id", description: "Show the full path and directory of the active gitmap binary, resolving symlinks",
    usage: "gitmap installed-dir",
    examples: [
      { command: "gitmap installed-dir", description: "Show installed binary path and directory" },
      { command: "gitmap id", description: "Alias shorthand" },
    ],
    seeAlso: [
      { name: "update", description: "Self-update to latest version" },
      { name: "version", description: "Show current version" },
      { name: "doctor", description: "Diagnose PATH and version issues" },
    ],
  },
  {
    category: "tools",
    name: "version", alias: "v", description: "Print the installed gitmap version",
    usage: "gitmap version",
    examples: [
      { command: "gitmap version", description: "Print current version (e.g. v2.48.2)" },
      { command: "gitmap v", description: "Alias shorthand" },
    ],
    seeAlso: [
      { name: "update", description: "Self-update to latest version" },
      { name: "doctor", description: "Diagnose version issues" },
    ],
  },
  {
    category: "tools",
    name: "completion", alias: "cmp", description: "Generate or install shell tab-completion scripts for Bash, Zsh, or PowerShell",
    usage: "gitmap completion <powershell|bash|zsh> [--list-repos] [--list-groups] [--list-commands]",
    flags: [
      { flag: "--list-repos", description: "Print repo slugs, one per line (for script use)" },
      { flag: "--list-groups", description: "Print group names, one per line (for script use)" },
      { flag: "--list-commands", description: "Print all command names, one per line (for script use)" },
    ],
    examples: [
      { command: "gitmap completion powershell", description: "Print PowerShell completion script" },
      { command: "gitmap completion bash", description: "Print Bash completion script" },
      { command: "gitmap cmp zsh", description: "Print Zsh completion script" },
      { command: "gitmap completion --list-repos", description: "List repo slugs for scripting" },
    ],
    seeAlso: [
      { name: "setup", description: "Auto-installs completions during setup" },
      { name: "cd", description: "Navigate to repos using tab-completed slugs" },
      { name: "group", description: "Group names are also tab-completed" },
    ],
  },
  {
    category: "tools",
    name: "interactive", alias: "i", description: "Launch the full-screen interactive TUI with 9 views for browsing, actions, and management",
    usage: "gitmap interactive [--refresh <seconds>]",
    flags: [
      { flag: "--refresh <seconds>", description: "Dashboard auto-refresh interval (default: 30)" },
    ],
    examples: [
      { command: "gitmap i", description: "Launch the interactive TUI" },
      { command: "gitmap interactive --refresh 10", description: "Launch with 10s dashboard refresh" },
    ],
    seeAlso: [
      { name: "list", description: "Non-interactive repo listing" },
      { name: "status", description: "Non-interactive status dashboard" },
      { name: "group", description: "CLI group management" },
      { name: "exec", description: "CLI batch execution" },
    ],
  },
  {
    category: "tools",
    name: "docs", alias: "d", description: "Open the gitmap documentation website in your default browser",
    usage: "gitmap docs",
    examples: [
      { command: "gitmap docs", description: "Open docs in browser" },
      { command: "gitmap d", description: "Open docs (short alias)" },
    ],
    seeAlso: [
      { name: "version", description: "Show installed version" },
    ],
  },
  {
    category: "tools",
    name: "ssh", description: "Generate, list, and manage SSH keys for Git authentication",
    usage: "gitmap ssh [subcommand] [flags]",
    flags: [
      { flag: "--name, -n <label>", description: "Key label in database (default: 'default')" },
      { flag: "--path, -p <path>", description: "Private key file path (default: ~/.ssh/id_rsa)" },
      { flag: "--email, -e <email>", description: "Email comment (default: git global email)" },
      { flag: "--force, -f", description: "Skip prompt if key already exists" },
      { flag: "--files", description: "Also delete key files from disk (delete subcommand)" },
      { flag: "--ssh-key, -K <name>", description: "SSH key name to use for cloning (clone flag)" },
    ],
    examples: [
      { command: "gitmap ssh", description: "Generate default RSA-4096 key" },
      { command: "gitmap ssh --name work", description: "Generate a named key for work" },
      { command: "gitmap ssh cat", description: "Display the default public key" },
      { command: "gitmap ssh cat --name work", description: "Display a named public key" },
      { command: "gitmap ssh list", description: "List all stored SSH keys" },
      { command: "gitmap ssh delete --name work --files", description: "Delete key record and files" },
      { command: "gitmap clone repos.json --ssh-key work", description: "Clone using a specific SSH key" },
    ],
    seeAlso: [
      { name: "clone", description: "Clone repos with --ssh-key integration" },
      { name: "setup", description: "Configure Git global settings" },
    ],
  },
  {
    category: "tools",
    name: "fix-auth", alias: "fa",
    description: "One-shot fix for 'Permission denied to <wrong-user>' SSH push failures — generates a per-account ed25519 key, pins the current repo to it via core.sshCommand (IdentitiesOnly=yes), and copies the public key to your clipboard.",
    usage: "gitmap fix-auth --user <github-username> [--email <addr>] [-y] [-f]",
    flags: [
      { flag: "--user, -u <name>", description: "GitHub username (required); first positional arg also accepted" },
      { flag: "--email, -e <addr>", description: "Email comment on the key (default: git config user.email)" },
      { flag: "--yes, -y", description: "Skip overwrite confirmation prompt" },
      { flag: "--force, -f", description: "Regenerate the key even if it already exists" },
    ],
    examples: [
      { command: "gitmap fix-auth --user aukgit --email me@example.com", description: "Full fix inside the broken repo" },
      { command: "gitmap fa aukgit", description: "Positional alias — reuse existing key or generate" },
      { command: "gitmap fix-auth -u aukgit -f -y", description: "Force-rotate the per-account key without prompt" },
    ],
    seeAlso: [
      { name: "whoami", description: "Diagnose which identity/key Git & SSH will use" },
      { name: "ssh-bind", description: "Pin an existing key to the current repo" },
      { name: "ssh", description: "General SSH key generation and DB" },
    ],
  },
  {
    category: "tools",
    name: "whoami", alias: "who",
    description: "Print local vs global git identity, origin transport, the SSH principal GitHub sees, the offered key, and the list of ~/.ssh candidate keys — with copy-paste fix hints.",
    usage: "gitmap whoami",
    examples: [
      { command: "gitmap whoami", description: "Diagnose why the wrong GitHub account is being used" },
      { command: "gitmap who", description: "Short alias" },
    ],
    seeAlso: [
      { name: "fix-auth", description: "One-shot generate + pin + clipboard" },
      { name: "ssh-bind", description: "Pin an existing key to the current repo" },
    ],
  },
  {
    category: "tools",
    name: "ssh-bind", alias: "sb",
    description: "Pin an existing SSH private key to the current repo via core.sshCommand with IdentitiesOnly=yes — stops the agent from silently offering the wrong-account default key.",
    usage: "gitmap ssh-bind <key-filename-or-path>",
    examples: [
      { command: "gitmap ssh-bind id_ed25519_aukgit", description: "Bare filename resolved under ~/.ssh/" },
      { command: "gitmap sb ~/.ssh/id_rsa_work", description: "Tilde-relative path" },
      { command: "gitmap sb /etc/ssh/deploy_key_aukgit", description: "Absolute path" },
    ],
    seeAlso: [
      { name: "fix-auth", description: "Generate a per-account key AND bind it in one step" },
      { name: "whoami", description: "List candidate keys and detect the current mismatch" },
    ],
  },
  {
    category: "tools",
    name: "gomod", alias: "gm", description: "Rename Go module path across the entire repo with branch safety",
    usage: "gitmap gomod <new-module-path> [--ext *.go,*.md] [--dry-run] [--no-merge] [--no-tidy] [--verbose]",
    flags: [
      { flag: "--ext <exts>", description: "Comma-separated extensions to filter (e.g. *.go,*.md)" },
      { flag: "--dry-run", description: "Preview changes without modifying files" },
      { flag: "--no-merge", description: "Commit on feature branch but do not merge back" },
      { flag: "--no-tidy", description: "Skip go mod tidy after replacement" },
      { flag: "--verbose", description: "Print each file path as it is modified" },
    ],
    examples: [
      { command: 'gitmap gomod "github.com/new/name"', description: "Rename module path in all files" },
      { command: 'gitmap gomod "x/y" --ext "*.go,*.md"', description: "Only replace in .go and .md files" },
      { command: 'gitmap gomod "github.com/new/name" --dry-run', description: "Preview what would change" },
      { command: 'gitmap gomod "github.com/new/name" --no-merge --verbose', description: "Replace on feature branch with logging" },
    ],
    seeAlso: [
      { name: "Spec: gomod", description: "Go module rename documentation", url: "/gomod" },
      { name: "go-repos", description: "List detected Go projects" },
      { name: "scan", description: "Scan directories" },
    ],
  },
  {
    category: "tools",
    name: "seo-write", alias: "sw", description: "Auto-generate and commit SEO-optimized messages with templates or CSV",
    usage: "gitmap seo-write [--url <url>] [--csv <path>] [--dry-run]",
    flags: [
      { flag: "--csv <path>", description: "Read title/description pairs from a CSV file" },
      { flag: "--url <url>", description: "Target website URL (required in template mode)" },
      { flag: "--service <name>", description: "Service name for {service} placeholder" },
      { flag: "--area <name>", description: "Area/location for {area} placeholder" },
      { flag: "--company <name>", description: "Company name for {company} placeholder" },
      { flag: "--phone <number>", description: "Phone number for {phone} placeholder" },
      { flag: "--email <address>", description: "Email for {email} placeholder" },
      { flag: "--address <text>", description: "Address for {address} placeholder" },
      { flag: "--max-commits <n>", description: "Stop after N commits (0 = unlimited)" },
      { flag: "--interval <min-max>", description: "Random delay range in seconds (default: 60-120)" },
      { flag: "--dry-run", description: "Preview commit messages without executing" },
      { flag: "--template <path>", description: "Load templates from a custom JSON file" },
      { flag: "--create-template", description: "Write starter seo-templates.json to current dir" },
      { flag: "--author-name <name>", description: "Git author name for commits" },
      { flag: "--author-email <email>", description: "Git author email for commits" },
    ],
    examples: [
      { command: "gitmap sw --url example.com --service Plumbing --area London", description: "Template mode with placeholders" },
      { command: "gitmap seo-write --csv ./commits.csv", description: "CSV mode — read from file" },
      { command: "gitmap sw --url example.com --dry-run", description: "Preview without committing" },
      { command: "gitmap seo-write --create-template", description: "Generate starter template file" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories" },
      { name: "history", description: "View command history", url: "/history" },
    ],
  },
  // ═══════════════════════════════════════════
  // Tools & Setup (new commands)
  // ═══════════════════════════════════════════
  {
    category: "tools",
    name: "task", alias: "tk", description: "Manage named file-sync watch tasks for one-way folder synchronization",
    usage: "gitmap task <create|list|run|show|delete> [flags]",
    flags: [
      { flag: "--src <path>", description: "Source directory path (create)" },
      { flag: "--dest <path>", description: "Destination directory path (create)" },
      { flag: "--interval <seconds>", description: "Sync interval in seconds, minimum 2 (run, default: 5)" },
      { flag: "--verbose", description: "Show detailed sync output (run)" },
      { flag: "--dry-run", description: "Preview sync actions without copying (run)" },
    ],
    examples: [
      { command: "gitmap task create my-sync --src ./src --dest ./backup", description: "Create a sync task" },
      { command: "gitmap tk run my-sync --interval 10 --verbose", description: "Run with 10s interval and verbose output" },
      { command: "gitmap task list", description: "List all saved tasks" },
      { command: "gitmap task show my-sync", description: "Show task details" },
      { command: "gitmap task delete my-sync", description: "Remove a task" },
    ],
    seeAlso: [
      { name: "watch", description: "Live-refresh dashboard of repo status" },
      { name: "exec", description: "Run git commands across repos" },
    ],
  },
  {
    category: "tools",
    name: "env", alias: "ev", description: "Manage persistent environment variables and PATH entries across platforms",
    usage: "gitmap env <set|get|delete|list|path add|path remove|path list> [flags]",
    flags: [
      { flag: "--system", description: "Target system-level variables (Windows, requires admin)" },
      { flag: "--shell <name>", description: "Target shell profile: bash, zsh (Unix only)" },
      { flag: "--verbose", description: "Show detailed operation output" },
      { flag: "--dry-run", description: "Preview changes without applying" },
    ],
    examples: [
      { command: 'gitmap env set GOPATH "/home/user/go"', description: "Set a persistent variable" },
      { command: "gitmap ev path add /usr/local/go/bin", description: "Add directory to PATH" },
      { command: "gitmap env list", description: "List managed variables" },
      { command: "gitmap env path list", description: "List managed PATH entries" },
      { command: "gitmap env delete GOPATH --dry-run", description: "Preview variable removal" },
      { command: 'gitmap env set JAVA_HOME "/usr/lib/jvm/java-17" --shell zsh', description: "Write to .zshrc instead of auto-detected profile" },
      { command: "gitmap ev path add /opt/bin --shell bash", description: "Add PATH entry targeting .bashrc explicitly" },
    ],
    seeAlso: [
      { name: "install", description: "Install developer tools" },
      { name: "doctor", description: "Diagnose PATH and version issues" },
      { name: "setup", description: "Configure Git global settings" },
    ],
  },
  {
    category: "tools",
    name: "install", alias: "in", description: "Install a developer tool by name using the platform package manager",
    usage: "gitmap install <tool> [flags]",
    flags: [
      { flag: "--manager <name>", description: "Force package manager (choco, winget, apt, brew)" },
      { flag: "--version <ver>", description: "Install a specific version" },
      { flag: "--verbose", description: "Show full installer output" },
      { flag: "--dry-run", description: "Show install command without executing" },
      { flag: "--check", description: "Only check if tool is installed" },
      { flag: "--list", description: "List all supported tools" },
    ],
    examples: [
      { command: "gitmap install vscode", description: "Install VS Code" },
      { command: "gitmap in go --check", description: "Check if Go is installed" },
      { command: "gitmap install python --dry-run", description: "Preview install command" },
      { command: "gitmap install node --verbose", description: "Install with verbose output" },
      { command: "gitmap install npp", description: "NPP + Settings — Install Notepad++ with settings" },
      { command: "gitmap install npp-settings", description: "NPP Settings — Sync settings only" },
      { command: "gitmap install install-npp", description: "Install NPP — Install Notepad++ only (no settings)" },
      { command: "gitmap install --list", description: "List all supported tools" },
      { command: "gitmap install scripts", description: "Clone gitmap scripts to local folder (Win: D:\\, Linux: ~/Desktop/)" },
    ],
    seeAlso: [
      { name: "env", description: "Manage environment variables and PATH" },
      { name: "doctor", description: "Diagnose PATH and version issues" },
      { name: "setup", description: "Configure Git global settings" },
    ],
  },
  {
    category: "tools",
    name: "uninstall", alias: "un", description: "Remove a developer tool installed via `gitmap install`. Without a tool name, delegates to `gitmap self-uninstall` for removing the gitmap binary itself.",
    usage: "gitmap uninstall <tool> [--dry-run] [--force] [--purge]",
    flags: [
      { flag: "--dry-run", description: "Print the package-manager uninstall command without running it" },
      { flag: "--force", description: "Skip the y/N confirmation prompt and the 'not in DB' guard" },
      { flag: "--purge", description: "Remove user/config data too (apt purge, choco -x); ignored on managers without a purge mode" },
    ],
    examples: [
      { command: "gitmap uninstall vscode", description: "Prompt, then uninstall using the manager that originally installed it" },
      { command: "gitmap un node --dry-run", description: "Preview the uninstall command (manager auto-detected)" },
      { command: "gitmap uninstall python --force --purge", description: "No prompt, also remove user config (apt purge / choco -x)" },
      { command: "gitmap uninstall --confirm --keep-data", description: "No tool name → shortcut for `gitmap self-uninstall`; flags pass through verbatim" },
    ],
    seeAlso: [
      { name: "install", description: "Install a developer tool" },
      { name: "self-uninstall", description: "Remove the gitmap binary itself" },
      { name: "doctor", description: "Verify a tool was actually removed from PATH" },
    ],
  },
  {
    category: "tools",
    name: "pending", description: "List all pending tasks that have not yet completed successfully",
    usage: "gitmap pending",
    examples: [
      { command: "gitmap pending", description: "List all pending tasks with ID, type, path, and failure reason" },
    ],
    seeAlso: [
      { name: "do-pending", description: "Retry pending tasks" },
      { name: "clone-next", description: "Clone next versioned iteration" },
      { name: "task", description: "Manage file-sync watch tasks" },
    ],
  },
  {
    category: "tools",
    name: "do-pending", alias: "dp", description: "Retry all pending tasks or a specific task by ID",
    usage: "gitmap do-pending [task-id]",
    examples: [
      { command: "gitmap do-pending", description: "Retry all pending tasks" },
      { command: "gitmap dp", description: "Retry all using alias" },
      { command: "gitmap do-pending 2", description: "Retry a specific task by its ID" },
    ],
    seeAlso: [
      { name: "pending", description: "List all pending tasks" },
      { name: "clone-next", description: "Clone next versioned iteration" },
      { name: "task", description: "Manage file-sync watch tasks" },
    ],
  },
  {
    category: "tools",
    name: "llm-docs", alias: "ld", description: "Generate a consolidated LLM.md reference file for AI assistants to understand all commands",
    usage: "gitmap llm-docs",
    examples: [
      { command: "gitmap llm-docs", description: "Generate LLM.md in the current directory" },
      { command: "gitmap ld", description: "Generate using alias" },
    ],
    seeAlso: [
      { name: "help", description: "Show CLI help text" },
      { name: "docs", description: "Open documentation website" },
      { name: "version", description: "Show version number" },
    ],
  },

  // ═══════════════════════════════════════════
  // Move & Merge (spec/01-app/97-move-and-merge.md)
  // ═══════════════════════════════════════════
  {
    category: "movemerge",
    name: "mv", alias: "move", description: "Move LEFT folder's contents into RIGHT, then delete LEFT entirely (each side may be a local folder or a remote git URL)",
    usage: "gitmap mv LEFT RIGHT [--no-push] [--no-commit] [--force-folder] [--pull] [--init] [--dry-run]",
    examples: [
      { command: "gitmap mv ./gitmap-v28 ./gitmap-v28", description: "Move local folder into another local folder, deleting source" },
      { command: "gitmap mv ./local https://github.com/owner/repo", description: "Move local folder into a remote repo (clone, copy, commit, push)" },
      { command: "gitmap mv https://github.com/owner/repo:develop ./mirror", description: "Pin remote branch and move into local folder" },
      { command: "gitmap mv ./a ./b --dry-run", description: "Preview the move without writing anything" },
    ],
    seeAlso: [
      { name: "merge-right", description: "Safer copy-with-prompt variant (LEFT not deleted)" },
      { name: "merge-both", description: "Two-way merge instead of move" },
      { name: "diff", description: "Preview tree differences before moving" },
    ],
  },
  {
    category: "movemerge",
    name: "merge-both", alias: "mb", description: "Bidirectional file-level merge: each side gains the other's missing files; conflicts trigger an [L]eft/[R]ight/[S]kip/[A]ll-left/[B]all-right/[Q]uit prompt",
    usage: "gitmap merge-both LEFT RIGHT [-y] [--prefer-newer|--prefer-left|--prefer-right|--prefer-skip] [--no-push] [--no-commit] [--dry-run]",
    examples: [
      { command: "gitmap merge-both ./gitmap-v28 ./gitmap-v28", description: "Interactive two-way merge between two local folders" },
      { command: "gitmap mb ./local https://github.com/owner/repo -y", description: "Non-interactive (newer wins by default for merge-both); commits + pushes the URL side" },
      { command: "gitmap merge-both ./a ./b -y --prefer-left --dry-run", description: "Preview a LEFT-wins merge without writing" },
    ],
    seeAlso: [
      { name: "diff", description: "Recommended dry-run preview before merge-both" },
      { name: "merge-left", description: "One-way merge into LEFT only" },
      { name: "merge-right", description: "One-way merge into RIGHT only" },
      { name: "mv", description: "Move LEFT into RIGHT and delete LEFT" },
    ],
  },
  {
    category: "movemerge",
    name: "merge-left", alias: "ml", description: "One-way merge that writes only into LEFT; missing files copied from RIGHT, conflicts resolved into LEFT. RIGHT is never modified.",
    usage: "gitmap merge-left LEFT RIGHT [-y] [--prefer-right|--prefer-left|--prefer-newer|--prefer-skip] [--no-push] [--no-commit]",
    examples: [
      { command: "gitmap merge-left ./gitmap-v28 ./gitmap-v28", description: "Pull RIGHT's changes into LEFT (interactive prompt)" },
      { command: "gitmap ml ./local https://github.com/owner/upstream -y", description: "Non-interactive (RIGHT wins by default for merge-left)" },
      { command: "gitmap merge-left ./mine ./theirs -y --prefer-left", description: "Bypass + keep LEFT everywhere on conflict" },
    ],
    seeAlso: [
      { name: "merge-right", description: "Mirror operation: write into RIGHT only" },
      { name: "merge-both", description: "Two-way merge instead of one-way" },
      { name: "diff", description: "Preview tree differences first" },
    ],
  },
  {
    category: "movemerge",
    name: "merge-right", alias: "mr", description: "One-way merge that writes only into RIGHT; missing files copied from LEFT, conflicts resolved into RIGHT. LEFT is never modified.",
    usage: "gitmap merge-right LEFT RIGHT [-y] [--prefer-left|--prefer-right|--prefer-newer|--prefer-skip] [--no-push] [--no-commit]",
    examples: [
      { command: "gitmap merge-right ./local https://github.com/owner/repo -y", description: "Push LEFT's changes into a remote repo (LEFT wins by default for merge-right)" },
      { command: "gitmap mr ./local https://github.com/owner/repo:develop -y", description: "Push to a specific branch" },
      { command: "gitmap merge-right ./local ./mirror --no-push", description: "Stage changes locally without pushing" },
    ],
    seeAlso: [
      { name: "merge-left", description: "Mirror operation: write into LEFT only" },
      { name: "merge-both", description: "Two-way merge instead of one-way" },
      { name: "mv", description: "Move LEFT into RIGHT and delete LEFT" },
    ],
  },
  {
    category: "movemerge",
    name: "diff", alias: "df", description: "Read-only preview of what merge-both/merge-left/merge-right would change between two folders; lists missing-on-each-side and conflicting files. Writes nothing.",
    usage: "gitmap diff LEFT RIGHT [--json] [--only-conflicts] [--only-missing] [--include-identical] [--include-vcs] [--include-node-modules]",
    examples: [
      { command: "gitmap diff ./gitmap-v28 ./gitmap-v28", description: "Plain text diff of two local folders" },
      { command: "gitmap diff ./a ./b --only-conflicts", description: "Show only files that differ on both sides" },
      { command: "gitmap df ./a ./b --json", description: "Machine-readable {summary, entries} JSON" },
    ],
    seeAlso: [
      { name: "merge-both", description: "Apply a two-way merge after previewing" },
      { name: "merge-left", description: "Apply RIGHT's changes into LEFT" },
      { name: "merge-right", description: "Apply LEFT's changes into RIGHT" },
      { name: "mv", description: "Move LEFT into RIGHT and delete LEFT" },
    ],
  },
  // ═══════════════════════════════════════════
  // Chrome Profiles (v6.43.0+)
  // ═══════════════════════════════════════════
  {
    category: "chrome",
    name: "chrome-profile-copy", alias: "cpc", description: "Copy a Chrome profile directory to a new profile, scrub identity fields so the copy appears in the profile picker, and register it in Chrome's Local State. Handles locked LOCK files and shows human-readable display names.",
    usage: "gitmap cpc <source> <dest> [-y] [-r|--register-only]",
    flags: [
      { flag: "-y, --yes", description: "Skip confirmation prompts" },
      { flag: "-r, --register-only", description: "Skip the file copy and only update Local State registry" },
    ],
    examples: [
      { command: "gitmap cpc \"Profile 1\" \"Work\"", description: "Copy by directory or display name" },
      { command: "gitmap cpc \"Profile 1\" \"Work\" -r", description: "Register an already-copied profile only" },
    ],
  },
  {
    category: "chrome",
    name: "chrome-profile-merge", alias: "cpm", description: "Merge settings, bookmarks, and extensions from one Chrome profile into another. Supports --dry-run with a diff renderer showing + add / ~ overwrite / = keep entries.",
    usage: "gitmap cpm <source> <dest> [--dry-run] [-y]",
    examples: [
      { command: "gitmap cpm \"Profile 1\" \"Profile 2\" --dry-run", description: "Preview the merge without writing" },
      { command: "gitmap cpm \"Profile 1\" \"Profile 2\" -y", description: "Apply the merge" },
    ],
  },
  {
    category: "chrome",
    name: "chrome-profiles", alias: "cpl", description: "List every Chrome profile with directory name, display name, and last-active timestamp resolved from Chrome's Local State JSON.",
    usage: "gitmap chrome-profiles [--json]",
    examples: [{ command: "gitmap chrome-profiles" }],
  },
  // ═══════════════════════════════════════════
  // Bulk Visibility (v6.53.0+, v6.56.0+)
  // ═══════════════════════════════════════════
  {
    category: "visibility-bulk",
    name: "make-last-public", alias: "MLPUB", description: "Flip the latest versioned repo (matched via fuzzy base+version index) to public on GitHub or GitLab. Uses the OwnerRepoNameIndex SQLite cache for fast resolution.",
    usage: "gitmap make-last-public <owner-or-url> <base>",
    examples: [{ command: "gitmap MLPUB myowner macro-ahk", description: "Make the newest macro-ahk-vN repo public" }],
  },
  {
    category: "visibility-bulk",
    name: "make-last-private", alias: "MLPRI", description: "Flip the latest versioned repo (matched via fuzzy base+version index) to private on GitHub or GitLab. Mirror of make-last-public.",
    usage: "gitmap make-last-private <owner-or-url> <base>",
    examples: [{ command: "gitmap MLPRI myowner macro-ahk" }],
  },
  // ═══════════════════════════════════════════
  // Maintenance & Backups (v6.57.0+)
  // ═══════════════════════════════════════════
  {
    category: "maintenance",
    name: "backup", description: "Manage bounded backups under .gitmap/backup/. `ls` lists existing snapshots; `prune` deletes old ones (keeps the newest N or those younger than a cutoff).",
    usage: "gitmap backup ls | gitmap backup prune [--keep N] [--older-than 30d]",
    examples: [
      { command: "gitmap backup ls" },
      { command: "gitmap backup prune --keep 10", description: "Keep only the 10 newest snapshots" },
      { command: "gitmap backup prune --older-than 30d", description: "Delete snapshots older than 30 days" },
    ],
  },
  {
    category: "maintenance",
    name: "ssh status", description: "Report SSH key health: which key matches the current remote, presence in the agent, last-used timestamp, and warnings for keys missing from ssh-agent.",
    usage: "gitmap ssh status [--json]",
    examples: [{ command: "gitmap ssh status" }],
  },
  {
    category: "maintenance",
    name: "changelog-regen", description: "Regenerate CHANGELOG.md from the JSON release manifests under .gitmap/release/. Useful after manual edits or after restoring release metadata from backup.",
    usage: "gitmap changelog-regen [--dry-run]",
    examples: [{ command: "gitmap changelog-regen --dry-run" }],
  },
  {
    category: "data",
    name: "rm", alias: "remove / del", description: "Delete tracked repos from the gitmap database, with on-disk folder deletion support, glob patterns (e.g. `macro*`), batch operation, and -y for no prompting.",
    usage: "gitmap rm <pattern>... [--disk] [-y|--yes]",
    flags: [
      { flag: "--disk", description: "Also delete the on-disk folder" },
      { flag: "-y, --yes", description: "Skip confirmation prompts" },
    ],
    examples: [
      { command: "gitmap rm macro-ahk-v1", description: "Remove a single tracked repo" },
      { command: "gitmap rm \"macro*\" -y", description: "Glob batch delete without prompting" },
      { command: "gitmap del old-repo --disk -y", description: "Remove from DB and delete folder" },
    ],
  },
  {
    category: "maintenance",
    name: "doctor", description: "One-shot health check for every external dependency gitmap touches (git, ssh, chrome, PATH, sqlite, disk). Prints targeted fix recipes for each failed probe and exits non-zero when anything is broken — drop it in CI as a pre-flight gate.",
    usage: "gitmap doctor [--json]",
    examples: [
      { command: "gitmap doctor", description: "Human-readable report" },
      { command: "gitmap doctor --json | jq '.data.checks[] | select(.ok==false)'", description: "Machine-readable, filter to failed probes" },
    ],
  },
  {
    category: "maintenance",
    name: "self-update", description: "Probe the GitHub release feed for the newest tag and re-run `gitmap self-install` non-interactively when a newer version is available. Honors --dry-run and --force.",
    usage: "gitmap self-update [--dry-run] [--force]",
    flags: [
      { flag: "--dry-run", description: "Print the would-be install command without touching disk" },
      { flag: "--force", description: "Re-install even when already on latest" },
    ],
    examples: [
      { command: "gitmap self-update" },
      { command: "gitmap self-update --dry-run" },
    ],
  },
  {
    category: "release",
    name: "release-undo --range", description: "Roll back several contiguous release tags at once. Endpoints must share major.minor (e.g. v6.60.0..v6.65.0) to prevent accidental mass-deletion. Each version is processed via the single-tag pipeline, so failures stop the run with earlier successes intact.",
    usage: "gitmap release-undo --range vX.Y.A..vX.Y.B [--keep-remote] [--keep-sidecar] [-y] [--dry-run]",
    examples: [
      { command: "gitmap release-undo --range v6.60.0..v6.65.0 --dry-run", description: "Preview the rollback list" },
      { command: "gitmap ru --range v6.60.0..v6.62.0 -y", description: "Roll back three tags non-interactively" },
    ],
  },
];

// Global flags applied by every command via the shared root parser.
export const GlobalFlags: CommandFlag[] = [
  { flag: "--log-json", description: "Emit structured NDJSON logs (schema gitmap.log.v1) to stderr — pipe into jq, Loki, or any log shipper." },
  { flag: "--quiet", description: "Suppress non-error output." },
  { flag: "--no-color", description: "Disable ANSI colors." },
  { flag: "--json", description: "Emit a single jsonenv envelope (schema gitmap.v1) to stdout." },
];
