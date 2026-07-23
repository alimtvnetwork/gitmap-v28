# Memory: features/cli-commands
Updated: now

The CLI supports 60 subcommands with aliases: 'scan' (s), 'clone' (c), 'clone-next' (cn), 'pull' (p), 'rescan' (rsc), 'setup', 'status' (st), 'exec' (x), 'desktop-sync' (ds), 'release' (r), 'release-self' (rs/rself), 'release-branch' (rb), 'release-pending' (rp), 'latest-branch' (lb), 'list' (ls), 'group' (g), 'multi-group' (mg), 'db-reset', 'version' (v), 'changelog' (cl), 'list-versions' (lv), 'list-releases' (lr), 'revert', 'doctor', 'update', 'seo-write' (sw), 'amend' (am), 'amend-list' (al), 'history' (hi), 'history-reset' (hr), 'stats' (ss), 'bookmark' (bk), 'export' (ex), 'import' (im), 'profile' (pf), 'cd' (go), 'watch' (w), 'diff-profiles' (dp), 'gomod' (gm), 'go-repos' (gr), 'node-repos' (nr), 'react-repos' (rr), 'cpp-repos' (cr), 'csharp-repos' (csr), 'alias' (a), 'zip-group' (z), 'completion' (cmp), 'interactive' (i), 'clear-release-json' (crj), 'update-cleanup', 'has-any-updates' (hau/hac), 'docs' (d), 'changelog-generate' (cg), 'ssh', 'prune' (pr), 'temp-release' (tr), 'dashboard' (db), 'task' (tk), 'env' (ev), and 'install' (in). Current version: v2.48.1.

The release workflow re-runs legacy directory migration after returning to the original branch, ensuring old `.release/` files are merged into `.gitmap/release/` and removed before auto-commit.

## Batch Operations

The `pull` and `exec` commands support `--stop-on-fail` to halt batch operations after the first failure. Failed items are tracked with `FailWithError` and reported via `PrintFailureReport`. Partial failures exit with code 3 (`ExitPartialFailure`).

## task (tk) — File Sync Watch

Named, persistent file-sync tasks with one-way timestamp-based synchronization. Source-to-destination folder sync with configurable interval (default 5s), parallel goroutines, and .gitignore-based filtering. Tasks stored in `.gitmap/tasks.json`. Spec: `spec/01-app/79-task-watch.md`.

## env (ev) — Environment Variable Management

Cross-platform persistent environment variable and PATH management. Windows: User/System level via setx/registry. Unix: writes to shell profiles (.bashrc/.zshrc). Tracks managed variables in `.gitmap/env-registry.json`. Spec: `spec/01-app/80-env.md`.

## install (in) — Developer Tool Installer

Automated installation of dev tools (VS Code, Node.js, Go, Git, Python, etc.) using platform package managers (Chocolatey/Winget on Windows, apt/brew on Linux/macOS). Spec: `spec/01-app/81-install.md`.

## temp-release (tr)

Lightweight temporary branch creation from recent commits. Creates branches from SHAs without checkout or tags. Supports batch creation with version pattern (`$$` placeholder), auto-increment sequencing, listing, and removal (single/range/all) with confirmation prompts. Tracked in `TempReleases` SQLite table. Spec: `spec/01-app/55-temp-release.md`.

## interactive (i)

Full-screen TUI with 9 views: Repos, Actions, Groups, Status, Releases, Temp Releases, Zip Groups, Aliases, Logs. See `features/interactive-tui.md` for details.

## clone-next (cn)

Clone the next or a specific versioned iteration of the current repo into the parent directory. Parses `-vN` suffix from folder name and remote URL, increments (`v++`) or jumps (`vN`), clones, registers with GitHub Desktop, and optionally removes the old folder. Flags: `--delete`, `--keep`, `--no-desktop`, `--ssh-key`, `--verbose`. Spec: `spec/01-app/59-clone-next.md`.

## release-self (rs / rself)

Release gitmap itself from any directory. Resolves the source repo via `os.Executable()` + symlink resolution + `.git` root walk; on failure, falls back to `source_repo_path` stored in the SQLite Settings table. Skips directory switch if already in the source repo. Auto-fallback: `gitmap release` outside a Git repo triggers self-release automatically. Full flag parity with `release`. Spec: `spec/01-app/60-release-self.md`.
