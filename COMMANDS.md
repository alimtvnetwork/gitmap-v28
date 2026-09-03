# GitMap Complete Command & Subcommand Reference Guide

Welcome to the canonical reference manual for all **GitMap** (`v6.166.0`) commands and subcommands.

> 💡 **Quick Terminal Access:** You can run `gitmap <command> --help` or `gitmap help <command>` at any time from your shell for interactive manual pages, flags, and examples.

---

## 📑 Table of Contents

1. [Scanning & Repository Discovery](#1-scanning--repository-discovery)
2. [Search & File Discovery (Split DB Indexed)](#2-search--file-discovery-split-db-indexed)
3. [Cloning & Workspace Sync](#3-cloning--workspace-sync)
4. [Git Operations & Batch Execution](#4-git-operations--batch-execution)
5. [Database Architecture & Maintenance (`gitmap db` & `start-fresh`)](#5-database-architecture--maintenance-gitmap-db--start-fresh)
6. [Cross-Platform Deduplication (`gitmap find-duplicates`)](#6-cross-platform-deduplication-gitmap-find-duplicates)
7. [Antigravity (AGY) Integration & Subcommands](#7-antigravity-agy-integration--subcommands)
8. [VS Code Project Manager Integration (`gitmap vscode`)](#8-vs-code-project-manager-integration-gitmap-vscode)
9. [Chrome Browser Profile Suite (`gitmap chrome`)](#9-chrome-browser-profile-suite-gitmap-chrome)
10. [Release Management & Versioning (`gitmap release`)](#10-release-management--versioning-gitmap-release)
11. [History, Stats & Author Amendments](#11-history-stats--author-amendments)
12. [Commit Transfer & Replay Engine](#12-commit-transfer--replay-engine)
13. [Navigation, Groups & Aliases](#13-navigation-groups--aliases)
14. [Templates, Curated Configs & LFS](#14-templates-curated-configs--lfs)
15. [Data Portability, Profiles & Bookmarks](#15-data-portability-profiles--bookmarks)
16. [Release ZIP Archives (`gitmap zip-group`)](#16-release-zip-archives-gitmap-zip-group)
17. [Background Task Scheduler (`gitmap schedule`)](#17-background-task-scheduler-gitmap-schedule)
18. [Cluster & Distributed Delegation (`gitmap cluster`)](#18-cluster--distributed-delegation-gitmap-cluster)
19. [Tool Installers & Automation Macros](#19-tool-installers--automation-macros)
20. [CI/CD Pipeline Telemetry & Utilities](#20-cicd-pipeline-telemetry--utilities)

---

## 1. Scanning & Repository Discovery

Commands for discovering, tracking, and auditing Git repositories on local disks.

### `gitmap scan [dir]`
* **Alias:** `s`
* **Description:** Recursively scans a directory tree for Git repositories, recording clone URLs, active branches, remote URLs, and tags.
* **Flags:**
  * `--mode ssh|https`: Preferred clone URL style (default: `https`).
  * `--output terminal|csv|json`: Output formatting mode.
  * `--output-path <dir>`: Target output directory.
  * `--github-desktop`: Automatically registers discovered repositories with GitHub Desktop.
  * `--open`: Opens the output folder after scan completion.
  * `--quiet`: Suppresses non-essential hints for scripted environments.

### `gitmap rescan`
* **Alias:** `rsc`
* **Description:** Re-runs the scan across previously scanned roots using cached flags.

### `gitmap rescan-subtree <path>`
* **Alias:** `rss`
* **Description:** Narrowly re-scans a specific subfolder from a previous scan output.

### `gitmap list`
* **Alias:** `ls`
* **Description:** Displays a formatted table of all tracked repositories with slugs, paths, and branch states.
* **Flags:** `--compact`, `--groups`, `--filter <q>`, `--json`.

### `gitmap sf <subcommand>`
* **Description:** Manages the ScanFolder database records (the root directories tracked by GitMap).
* **Subcommands:**
  * `gitmap sf add <path>`: Registers a new persistent directory root for automatic scanning.
  * `gitmap sf list` (alias: `ls`): Displays all registered scan folder roots.
  * `gitmap sf rm <id|path>`: Untracks a scan folder root.

---

## 2. Search & File Discovery (Split DB Indexed)

High-speed repository and file search powered by per-repository Split SQLite databases.

| Command | Alias | Description |
|---|---|---|
| `gitmap find-files <name>` | `ff` | Find files matching exact filename across all repos (`-ext` filter available) |
| `gitmap find-files-any <str>` | `ffa` | Find files containing substring in filename across all repos |
| `gitmap find-files-startswith <pfx>` | `ffs` | Find files by filename prefix |
| `gitmap find-files-endswith <sfx>` | `ffe` | Find files by filename suffix |
| `gitmap find <wildcard*>` | `f` | Universal glob/wildcard search across indexed repositories |
| `gitmap list-files [pattern]` | `lf` | List indexed repository files |

---

## 3. Cloning & Workspace Sync

Commands for re-cloning repositories and synchronizing local working trees.

### `gitmap clone <source|json|csv|text>`
* **Alias:** `c`
* **Description:** Clones repositories listed in a scan output file.
* **Flags:**
  * `--target-dir <dir>`: Base directory for cloned repositories.
  * `--safe-pull`: Pulls existing repos with retry and diagnostic unlock.
  * `--verbose`: Writes timestamped diagnostic logs.

### `gitmap clone-sync <url> [urls...]`
* **Alias:** `cs`
* **Description:** Clones one or more repositories and immediately registers them in VS Code Workspaces and Antigravity.

### `gitmap clone-next [v++|vN]`
* **Alias:** `cn`
* **Description:** Clones the next versioned iteration of the current repo (e.g. `v11` $\rightarrow$ `v12`).
* **Flags:**
  * `--delete`: Automatically deletes the previous version directory after cloning.
  * `--keep`: Keeps current folder without prompting.
  * `--no-desktop`: Skips GitHub Desktop registration.
  * `--ssh-key <name>` (`-K`): Named SSH key to authenticate the clone.
  * `--create-remote`: Automatically creates the remote GitHub repo if missing (requires `GITHUB_TOKEN`).

### `gitmap desktop-sync`
* **Alias:** `ds`
* **Description:** Registers all tracked repositories with the local GitHub Desktop application.

### `gitmap github-desktop`
* **Alias:** `gd`
* **Description:** Registers the current repository with GitHub Desktop directly.

---

## 4. Git Operations & Batch Execution

Perform batch git commands, check synchronization status, and manage pull failures.

### `gitmap pull [name]`
* **Alias:** `p`
* **Description:** Pulls the latest commits for a specific repository, group, or all repositories with smart non-unlink failure diagnostics and tabular multi-repo status display.
* **Flags:**
  * `--all`: Pulls all tracked repositories in parallel.
  * `--group <name>`: Pulls all repositories assigned to a specific group.
  * `--verbose`: Displays detailed git stderr output.

### `gitmap pull-all`
* **Alias:** `pa`
* **Description:** Shorthand for `gitmap pull --all`.

### `gitmap fix [repo] [action]`
* **Description:** Applies automatic pull failure remediation (`stash`, `wip`, `discard`).

### `gitmap push`
* **Alias:** `ph`
* **Description:** Push current repository commits with support for `--ssh` and `--https` transport flags.

### `gitmap status`
* **Alias:** `st`
* **Description:** Renders a consolidated status table showing dirty/clean, ahead/behind, and stash counts across all repos.

### `gitmap watch`
* **Alias:** `w`
* **Description:** Launches a live auto-refreshing terminal dashboard of repository status.

### `gitmap exec <args...>`
* **Alias:** `x`
* **Description:** Executes an arbitrary Git command across all tracked repositories simultaneously.

### `gitmap has-any-updates`
* **Aliases:** `hau`, `hac`
* **Description:** Queries remote repositories to check if unpulled commits exist.

### `gitmap latest-branch`
* **Alias:** `lb`
* **Description:** Finds the most recently updated remote branch across tracking repositories.

### `gitmap lfs-common`
* **Alias:** `lfsc`
* **Description:** Automatically adds standard binary extensions to Git LFS tracking.

---

## 5. Database Architecture & Maintenance (`gitmap db` & `start-fresh`)

GitMap utilizes a hybrid database architecture: a centralized **Primary Master SQLite Database** (`bin/data/gitmap.db`) and per-repository **Split Databases** (`bin/data/repo_search/*.db`) to eliminate SQLite WAL write locks during parallel multi-goroutine scans.

### `gitmap db <subcommand>`
* **Subcommands:**
  * `gitmap db ls` (alias: `list`): Displays absolute paths, sizes, status, and architectural justification for master and split DBs.
  * `gitmap db repo-db list`: Lists per-repository split databases, row counts (`RepoFile`, `SearchCache`), file sizes, and tracking status.
  * `gitmap db sizes list` (alias: `sizes ls`): Renders a consolidated disk utilization summary.
  * `gitmap db reset` (alias: `clear`): Clears cached database records while prompting for confirmation (`--yes` to skip).
  * `gitmap db help`: Shows detailed database usage manual.

### `gitmap start-fresh`
* **Description:** Performs a complete, irreversible wipe of all master databases, journals (`-wal`, `-shm`), and split databases, then immediately executes schema rebuild migrations from scratch.
* **Flags:**
  * `--yes` (`-y`): Bypasses the irreversible transaction warning prompt.
  * `--force` (`-f`): Alias for `--yes`.

---

## 6. Cross-Platform Deduplication (`gitmap find-duplicates`)

Audits duplicate projects and repositories across developer tooling and outputs immediate, copy-pasteable remediation commands directly beneath the findings.

### `gitmap find-duplicates [platform]`
* **Aliases:** `dups`, `find-dups`
* **Platforms:**
  * `gitmap find-duplicates`: Runs audit across all 4 platforms simultaneously.
  * `gitmap find-duplicates agy` (or `gitmap agy find-duplicates`): Antigravity workspace JSON files.
  * `gitmap find-duplicates vscode` (or `gitmap vscode find-duplicates`): VS Code Project Manager.
  * `gitmap find-duplicates chrome` (or `gitmap chrome find-duplicates`): Chrome browser profiles.
  * `gitmap find-duplicates git` (or `gitmap git find-duplicates`): Tracked Git repositories.

---

## 7. Antigravity (AGY) Integration & Subcommands

Complete management suite for Google Antigravity workspace projects and agent sessions.

### `gitmap agy <subcommand>`
* **Aliases:** `ag`, `antigravity`
* **Subcommands:**
  * `gitmap agy ls`: Lists all registered Antigravity projects with IDs, names, paths, and update timestamps.
  * `gitmap agy ls show-projects-with-empty-conversations` (aliases: `show-proects-with-empty-conversations`, `empty-conversations`, `empty-convs`): Audits and lists all Antigravity projects having zero or aborted/empty conversation sessions.
  * `gitmap agy remove-projects-with-empty-conversations` (aliases: `rm-empty-conversations`, `clean-empty-conversations`, `prune-empty-conversations`): Deletes orphaned configuration files for empty-conversation projects.
    * Flags: `--except <spec>` (`-e`), `--dry-run` (`-d`), `--yes` (`-y`).
  * `gitmap agy optimize-projects` (alias: `--repeat-fix`): Deduplicates Antigravity projects sharing identical paths, preserving the newest file.
  * `gitmap agy clear`: Clears missing or stale Antigravity projects (`--except` supported).
  * `gitmap agy add <path>`: Manually registers a directory workspace into Antigravity.
  * `gitmap agy rm <id>`: Deletes a specific project JSON configuration by ID.
  * `gitmap agy open [slug|path]`: Opens Antigravity focused on a specific project.
  * `gitmap agy scan [path]`: Scans filesystem recursively and checks Antigravity registration status.
  * `gitmap agy export-projects [file]`: Exports all Antigravity project JSON files into a ZIP archive.
  * `gitmap agy import-projects <file>`: Restores Antigravity project configurations from a ZIP archive.
  * `gitmap agy sync [path]`: Synchronizes all tracked GitMap repositories into Antigravity workspaces.
  * `gitmap agy status`: Displays tabular Antigravity project metrics.
  * `gitmap agy stats`: Outputs aggregated statistics on Antigravity projects.
  * `gitmap agy prompt <project> <prompt>`: Sends a user prompt directly to an Antigravity project session.
  * `gitmap agy prompt-all-project <title> <prompt>` (alias: `pap`): Broadcasts a prompt across all Antigravity project workspaces.
  * `gitmap agy rw <path|slug>`: Enables rewrite-both mode for a project.
  * `gitmap agy plugin ls`: Lists installed Antigravity plugins.
  * `gitmap agy plugin install <slug>`: Installs an Antigravity plugin by slug.

---

## 8. VS Code Project Manager Integration (`gitmap vscode`)

### `gitmap vscode <subcommand>`
* **Aliases:** `vsc`, `code`
* **Subcommands:**
  * `gitmap vscode sync`: Synchronizes all tracked GitMap repos into `projects.json` for the VS Code Project Manager extension.
  * `gitmap vscode ls`: Lists all projects registered in VS Code Project Manager.
  * `gitmap vscode optimize-projects`: Deduplicates repeated project paths in `projects.json`.
  * `gitmap vscode clear`: Prunes stale or non-existent projects (`--except` supported).
  * `gitmap vscode paths add <alias> <path>`: Adds an explicit path mapping to VS Code registry.
  * `gitmap vscode paths rm <alias>`: Removes a path mapping.
  * `gitmap vscode paths list`: Lists registered path mappings.
  * `gitmap vscode install`: Installs VS Code CLI integration.

---

## 9. Chrome Browser Profile Suite (`gitmap chrome`)

### `gitmap chrome <subcommand>`
* **Aliases:** `cprof`, `chrome-profile`
* **Subcommands:**
  * `gitmap chrome open <profile>`: Launches Chrome with a specific profile.
  * `gitmap chrome observe`: Monitors active Chrome profile processes and locks.
  * `gitmap chrome list` (alias: `cpl`): Lists all detected Chrome profiles and directories.
  * `gitmap chrome copy <src> <dst>` (alias: `cpc`): Clones a live profile into an offline backup profile.
  * `gitmap chrome export <name>` (alias: `cpe`): Exports a profile snapshot (`.json`, `.zip`, `.sqlite`).
  * `gitmap chrome import <file>` (alias: `cpi`): Imports a Chrome profile from a snapshot archive.
  * `gitmap chrome delete <name>` (alias: `cpd`): Removes a profile and its stored artifacts from database.
  * `gitmap chrome merge <src> <dst>` (alias: `cpm`): Merges bookmarks, extensions, or preferences between two profiles.
  * `gitmap chrome optimize-projects`: Prunes redundant profile references.
  * `gitmap chrome reset <profile>`: Cleans cache and temporary session files for a profile.
  * `gitmap chrome extensions <profile>`: Lists installed extensions for a profile.
  * `gitmap chrome flags <profile>`: Inspects configured Chrome feature flags.

---

## 10. Release Management & Versioning (`gitmap release`)

Automate release ceremonies, Git tagging, changelog generation, and asset distribution.

### `gitmap release [ver]`
* **Alias:** `r`
* **Description:** Performs full release workflow: validates working tree, bumps version, generates changelog, builds cross-compiled binaries, creates Git tag, pushes release branch, and publishes to GitHub.
* **Flags:**
  * `--bump major|minor|patch`: Automatic SemVer increment.
  * `--draft`: Publishes as a draft release.
  * `--dry-run`: Previews release steps without modifying Git or pushing.
  * `--assets <path>`: Attaches files or directories to release.
  * `--bin` (`-b`): Cross-compiles Go binaries for attachment.
  * `--targets <list>`: Cross-compile matrix (e.g. `windows/amd64,linux/amd64`).
  * `--compress`: Wraps release binaries in `.zip` or `.tar.gz`.
  * `--checksums`: Generates `SHA256` checksums file.
  * `--yes` (`-y`): Auto-confirms release prompts.

### Other Release Subcommands
| Command | Alias | Description |
|---|---|---|
| `gitmap pull-release [ver]` | `pr` | Pulls with fast-forward/rebase, then executes release |
| `gitmap release-self` | `rs` | Releases GitMap itself from any directory |
| `gitmap release-branch` | `rb` | Finalizes a release from an existing release branch |
| `gitmap temp-release <count> <pattern>` | `tr` | Creates temporary test branches from recent commits |
| `gitmap revert <version>` | — | Reverts local and remote state to a specified version tag |
| `gitmap release-pending` | `rp` | Releases all pending branches that lack release tags |
| `gitmap clear-release-json <ver>` | `crj` | Deletes a `.gitmap/release/vX.Y.Z.json` metadata file |
| `gitmap prune` | `prn` | Deletes local release branches that have already been tagged |
| `gitmap changelog [ver]` | `cl` | Displays formatted release notes from changelog |
| `gitmap changelog-gen` | `cg` | Automatically generates markdown release notes between tags |
| `gitmap list-versions` | `lv` | Lists all release tags sorted by SemVer (highest first) |
| `gitmap list-releases` | `lr` | Shows release history from metadata files or database |

---

## 11. History, Stats & Author Amendments

| Command | Alias | Description |
|---|---|---|
| `gitmap history` | `hi` | Displays audit log of executed GitMap commands (`--limit N`, `--json`) |
| `gitmap history-reset` | `hr` | Clears command execution history (`--confirm` required) |
| `gitmap version-history` | `vh` | Displays version transition ledger for current repo |
| `gitmap stats` | `ss` | Shows aggregated command usage frequency and performance |
| `gitmap amend [hash]` | `am` | Rewrites commit author name and email (`--name`, `--email`, `--force-push`) |
| `gitmap amend-list` | `al` | Lists stored author amendments from database |

---

## 12. Commit Transfer & Replay Engine

Clean, idempotent commit replay across decoupled repositories without submodule lock-in.

| Command | Alias | Description |
|---|---|---|
| `gitmap commit-right` | `cmr` | Replays LEFT repository's commits onto RIGHT repository |
| `gitmap commit-left` | `cml` | Replays RIGHT repository's commits onto LEFT repository |
| `gitmap commit-both` | `cmb` | Bidirectional commit replay (sequential or date-interleaved with `--interleave`) |

---

## 13. Navigation, Groups & Aliases

### `gitmap cd <name>`
* **Alias:** `go`
* **Description:** Fast directory navigation to any tracked repository by slug or partial match.

### `gitmap group <subcommand>`
* **Alias:** `g`
* **Subcommands:**
  * `gitmap group create <name> [repos...]`: Creates a new repository group.
  * `gitmap group add <group> <repo>`: Adds a repository to a group.
  * `gitmap group remove <group> <repo>`: Removes a repository from a group.
  * `gitmap group list`: Lists all defined repository groups.
  * `gitmap group show <group>`: Shows member repositories of a group.
  * `gitmap group delete <group>`: Deletes a group.
  * `gitmap group <name>`: Activates a group for scoped batch operations (`pull`, `status`, `exec`).

### `gitmap alias <subcommand>`
* **Alias:** `a`
* **Subcommands:**
  * `gitmap alias set <alias> <repo>`: Sets a short navigation alias for a repository.
  * `gitmap alias rm <alias>` (alias: `remove`): Removes an alias.
  * `gitmap alias list`: Lists all active repository aliases.
  * `gitmap alias show <alias>`: Displays the repository target for an alias.
  * `gitmap alias suggest`: Suggests intuitive short aliases for untracked repos.

---

## 14. Templates, Curated Configs & LFS

Scaffold standardized configuration files with idempotent marker-block preservation.

### `gitmap add <template>`
* `gitmap add ignore [langs...]`: Merges curated `.gitignore` blocks into `./.gitignore`.
* `gitmap add attributes [langs...]`: Merges curated `.gitattributes` blocks into `./.gitattributes`.
* `gitmap add lfs-install`: Runs `git lfs install --local` and adds binary LFS attributes.

### `gitmap templates <subcommand>`
* **Alias:** `tpl`
* **Subcommands:**
  * `gitmap templates init [langs...]` (alias: `ti`): Scaffolds `.gitignore` and `.gitattributes`.
  * `gitmap templates list` (alias: `tl`): Lists all available templates and origins.
  * `gitmap templates show <name>` (alias: `ts`): Prints template contents to stdout.
  * `gitmap templates diff` (alias: `td`): Previews what template application would change.

### `gitmap sync <target>`
* **Alias:** `sy`
* **Targets:** `ignore`, `attributes`, `lfs-install`, `prettier-ignore`, `prettier-rc`, `all`.

---

## 15. Data Portability, Profiles & Bookmarks

### `gitmap profile <subcommand>`
* **Alias:** `pf`
* **Subcommands:**
  * `gitmap profile create <name>`: Creates a new isolated database profile.
  * `gitmap profile list` (aliases: `ls`, `status`): Lists available profiles.
  * `gitmap profile switch <name>`: Switches the active database profile.
  * `gitmap profile delete <name>`: Deletes a database profile.
  * `gitmap profile show`: Displays details of the current active profile.

### `gitmap bookmark <subcommand>`
* **Alias:** `bk`
* **Subcommands:**
  * `gitmap bookmark save <name> [command...]`: Saves a command and flag configuration.
  * `gitmap bookmark list`: Lists all saved bookmarks.
  * `gitmap bookmark run <name>`: Executes a saved bookmark.
  * `gitmap bookmark delete <name>`: Deletes a saved bookmark.

### Portability & Moving
* `gitmap export [file]` (alias: `ex`): Exports all database tables to portable JSON.
* `gitmap import <file>` (alias: `im`): Restores database from a JSON export.
* `gitmap mv <src> <dst>` (alias: `move`): Moves a repository folder and updates VS Code and Desktop links.
* `gitmap rm <name>` (aliases: `remove`, `del`): Untracks repository from GitMap database.

---

## 16. Release ZIP Archives (`gitmap zip-group`)

### `gitmap zip-group <subcommand>`
* **Alias:** `z`
* **Subcommands:**
  * `gitmap zip-group create <name> [paths...] [--archive <name>]`: Defines a named collection of files for release packaging.
  * `gitmap zip-group add <group> <path...>`: Adds files or folders to a zip group.
  * `gitmap zip-group remove <group> <path>`: Removes files from a zip group.
  * `gitmap zip-group list` (alias: `ls`): Lists all defined zip groups.
  * `gitmap zip-group show <group>`: Inspects files within a zip group.
  * `gitmap zip-group rename <group> --archive <name>`: Updates the target output archive filename.
  * `gitmap zip-group delete <group>`: Deletes a zip group definition.

---

## 17. Background Task Scheduler (`gitmap schedule`)

Automate background jobs, triggers, and recurring routines.

### `gitmap schedule <subcommand>`
* **Alias:** `sc`
* **Subcommands:**
  * `gitmap schedule add <name> <command> [cron]`: Enqueues a scheduled background task.
  * `gitmap schedule list` (alias: `ls`): Lists all active scheduled tasks.
  * `gitmap schedule status`: Displays runner daemon health and execution telemetry.
  * `gitmap schedule enable <id>`: Enables a paused task.
  * `gitmap schedule disable <id>`: Disables an active task.
  * `gitmap schedule run <id>`: Triggers immediate one-shot execution of a scheduled task.
  * `gitmap schedule logs <id>` (aliases: `log`, `history`): Inspects task execution logs.
  * `gitmap schedule rm <id>` (aliases: `delete`, `del`): Deletes a task.
  * `gitmap schedule export [file]`: Exports scheduled task configurations.
  * `gitmap schedule import <file>`: Imports scheduled task configurations.
  * `gitmap schedule startup`: Configures system autostart for the scheduler service.
  * `gitmap schedule restart`: Restarts background scheduler service.
  * `gitmap schedule shutdown`: Gracefully stops the scheduler service.

---

## 18. Cluster & Distributed Delegation (`gitmap cluster`)

Coordinate operations across multiple machines in a distributed network.

### `gitmap cluster <subcommand>`
* **Subcommands:**
  * `gitmap cluster status`: Shows cluster node connectivity and synchronization health.
  * `gitmap cluster nodes` (alias: `ls`): Lists all active nodes in the cluster.
  * `gitmap cluster history`: Displays distributed execution audit history.
  * `gitmap cluster set-password`: Sets cluster join and access passwords.
  * `gitmap cluster reset-password`: Resets cluster credentials.
  * `gitmap cluster remove <node-id>` (alias: `rm`): Disconnects and unregisters a node.
  * `gitmap cluster audit-clean`: Purges orphaned node communication logs.
  * `gitmap cluster stats`: Shows cluster resource utilization metrics.
* `gitmap serve` (alias: `sv`): Starts the orchestrator daemon and emits a join token.
* `gitmap servers-clients <command>`: Broadcasts command across both servers and client nodes.
* `gitmap clients <command>`: Broadcasts command across client worker nodes only.

---

## 19. Tool Installers & Automation Macros

### Developer Tool Installers (`gitmap installer`)
* **Aliases:** `installer`, `in`
* **Subcommands:**
  * `gitmap installer ls`: Lists available developer tool install recipes.
  * `gitmap installer install <slug>`: Runs verified installer for a tool.
  * `gitmap installer create <name>`: Authors a new tool installation recipe.
  * `gitmap installer export <slug>`: Exports installer script to a standalone file.
  * `gitmap installer import <file>`: Imports an installer recipe.
  * `gitmap installer update <slug>`: Updates tool installer to latest remote recipe.
  * `gitmap installer rm <slug>`: Deletes an installer recipe.
  * `gitmap install <tool>`: Direct shortcut to install a tool by name.
  * `gitmap uninstall <tool>`: Direct shortcut to uninstall a tool.

### Task Automation Macros (`gitmap macro`)
* **Alias:** `m`
* **Subcommands:**
  * `gitmap macro record <name>` (alias: `rec`): Interactively records terminal commands into a reusable macro.
  * `gitmap macro list` (alias: `ls`): Lists recorded macros.
  * `gitmap macro show <name>`: Displays the steps inside a macro.
  * `gitmap macro run <name>`: Executes a macro.
  * `gitmap macro run-until-succeed <name>` (alias: `retry`): Loops macro execution until exit code 0.
  * `gitmap macro rm <name>` (alias: `delete`): Deletes a recorded macro.

---

## 20. CI/CD Pipeline Telemetry & Utilities

### CI/CD Pipeline Telemetry
* `gitmap pipeline status`: Displays active CI/CD pipeline runs and steps.
* `gitmap pipeline eta`: Predicts completion time and remaining run duration for CI jobs.
* `gitmap pipeline logs`: Streams recent CI/CD pipeline step logs.

### System Utilities
* `gitmap doctor`: Runs diagnostic health checks on PATH, git binaries, and database connectivity (`--fix-path` supported).
* `gitmap update`: Self-updates GitMap binary to latest GitHub release.
* `gitmap update-cleanup`: Purges leftover `.old` executables and update temporary files.
* `gitmap version` (alias: `v`): Displays GitMap SemVer version (`v6.166.0`).
* `gitmap completion <shell>` (alias: `cmp`): Emits shell completion scripts (`bash`, `zsh`, `powershell`, `fish`).
* `gitmap interactive` (alias: `i`): Launches the full-screen terminal interactive TUI browser.
* `gitmap dashboard` (alias: `db`): Generates a standalone interactive HTML dashboard for a repository.
* `gitmap docs` (alias: `d`): Opens documentation website in the default browser.
* `gitmap help-dashboard` (alias: `hd`): Runs a local offline documentation server.
* `gitmap gomod <path>` (alias: `gm`): Safely renames Go module paths across the workspace.
* `gitmap seo-write` (alias: `sw`): Automated SEO commit message scheduler with templates.
* `gitmap llm-docs` (alias: `ld`): Emits a consolidated `LLM.md` context file for AI assistants.
* `gitmap fix-repo` (alias: `fr`): Rewrites prior `{base}-vN` version tokens (`-2`, `-3`, `-5`, `--all`, `--dry-run`).
* `gitmap clone-fix-repo <url>` (alias: `cfr`): Clones a repository and runs `fix-repo --all` in one step.
* `gitmap clone-fix-repo-pub <url>` (alias: `cfrp`): Clones, fixes version tokens, and makes the repository public.
* `gitmap make-public`: Sets current GitHub/GitLab repository to public visibility.
* `gitmap make-private`: Sets current GitHub/GitLab repository to private visibility.
* `gitmap open [target]` (alias: `o`): Opens current repo or target folder in OS file manager.
* `gitmap setup`: Configures git diff/merge tools and global GitMap preferences.
* `gitmap env <sub>`: Manages environment variables and PATH configurations.
* `gitmap cg <sub>`: Scaffolds Coding Guidelines (v24) into a repository.
* `gitmap user <sub>`: Manages cross-platform OS system users (Windows, Ubuntu, Debian, Fedora).
* `gitmap help [command]`: Displays the manual and flag reference for any command.
