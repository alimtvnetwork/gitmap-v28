---
name: clone-pick
description: gitmap clone-pick (cpk) does sparse-checkout of selected paths from a git repo with optional --ask tree picker (clone-once: picker's metadata clone is promoted to final dest), auto-saves selection to CloneInteractiveSelection table, supports --replay <id|name>
type: feature
---

# gitmap clone-pick / cpk (spec 100, --replay v4.19.0, --ask v4.20.0, clone-once v4.21.0, windowed scroller v4.22.0)

## What it does
Sparse-checkout a subset of a git repo into the current dir (or `--dest`).
Auto-saves every run to a new SQLite table for `--replay`.

## Surface
- `gitmap clone-pick <url> <p1,p2,...> [--ask] [--name X] [--branch B] [--mode https|ssh] [--depth N] [--cone] [--dest D] [--keep-git] [--dry-run] [--quiet] [--force]`
- `gitmap cpk <url> <paths> ...` (short alias — NOT `ci` because that collides with CI/CD muscle memory)
- `gitmap clone-pick --replay <id|name>` (no `<paths>` needed)

## Mechanism
1. `git clone --filter=blob:none --no-checkout [--branch] [--depth N] <url> <dest>`
2. `git sparse-checkout init [--cone]`
3. `git sparse-checkout set <paths...>`
4. `git checkout`
5. If `--keep-git=false` → `rm -rf .git`

Cone mode is auto-flipped off when any path contains glob chars or a file extension after `/`.

## --ask picker
- bubbletea TUI built like `tui/browser.go`
- Source: `git ls-tree -r --name-only HEAD` against the partial clone in a temp dir (the same temp clone is reused for the final checkout — cloned once, not twice)
- User-supplied `<paths>` pre-checked
- Auto-greyed: `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`, `__pycache__/` (overridable via `clonePick.autoExclude` in config)
- Keys: ↑/↓/j/k move, pgup/pgdn (ctrl+b/ctrl+f) page, g/G (home/end) jump, space toggle, a all, n none, s/enter save & clone, q/ctrl+c quit (exit 130)
- Windowed scroller (v4.22.0+): only the visible terminal window is rendered (height − 3 chrome rows); reacts to tea.WindowSizeMsg; default 20 rows pre-resize. Header shows "selected/total, rows N-M".

## Persistence: CloneInteractiveSelection table
Columns: SelectionId PK, Name (optional, unique-non-empty enforced in store layer), RepoCanonicalId, RepoUrl, Mode, Branch, Depth, Cone, KeepGit, DestDir, PathsCsv (sorted+normalised), UsedAsk, CreatedAt. NO FK to Repo (picked repo may not be in any local scan). Indexed by RepoCanonicalId and by Name (partial index where Name <> '').

## Replay rules
- Numeric `--replay` → SELECT by SelectionId
- Non-numeric → SELECT by Name (case-sensitive)
- Replay does NOT insert duplicate row; updates CreatedAt
- `--dry-run` never writes to DB

## Exit codes
- 0 success / dry-run ok
- 1 runtime (git/fs/replay-not-found)
- 2 bad CLI usage
- 130 user cancelled picker (`q`)

## Where it lives
- `gitmap/clonepick/` (parse, plan, sparse, picker, persist, render)
- `gitmap/cmd/clonepick.go` (dispatcher entry, registered in `rootcore.go` coreDispatchEntries)
- `gitmap/constants/constants_clonepick.go` (flags, messages, autoExclude defaults)
- `gitmap/helptext/clone-pick.md`
- `gitmap/store/cloneinteractiveselection.go` + entry in `Migrate()` statements list

## Why sparse-checkout over tarball/copy
Works with any host (not GitHub-only), single git invocation tree, leaves a real `.git` so `git pull` works, cone mode is O(matched paths).

## Shell handoff
Calls `WriteShellHandoff(dest)` after success (skipped when dest == "." or `--dry-run`). Same pattern as `gitmap clone <url>`.

## Spec ref
`spec/01-app/100-clone-pick.md`
