---
name: release-alias
description: gitmap as / release-alias (ra) / release-alias-pull (rap) workflow with auto-stash labeled by alias-version-unixts and label-match pop for concurrent safety
type: feature
---

# Release Alias Workflow (v3.0.0)

## Commands

| Command | Alias | Purpose |
|---------|-------|---------|
| `gitmap as [alias-name] [--force\|-f]` | `s-alias` | Tag the **current** Git repo with a short alias, persisted in the active-profile SQLite DB. Defaults alias to repo folder basename. Refuses to clobber existing alias unless `--force`. |
| `gitmap release-alias <alias> <version>` | `ra` | Release a previously-aliased repo from **any** CWD. Resolves alias → absolute path, `os.Chdir`s in, runs `runRelease`, restores CWD via `defer`. Forwards `--dry-run`. |
| `gitmap release-alias-pull <alias> <version>` | `rap` | Sugar for `release-alias --pull`. Runs `git pull --ff-only` first; hard-fails on non-fast-forward. |
| `gitmap db-migrate` | `dbm` | Idempotent schema migration. Auto-invoked at end of `gitmap update`. |

## Auto-stash semantics

- Dirty trees are auto-stashed before release: `git stash push --include-untracked -m "gitmap-release-alias autostash <alias>-<version>-<unix-ts>"`.
- Pop runs in `defer` so it always fires (even when `runRelease` aborts).
- Pop locates the stash by **label match against `git stash list`** (not by `stash@{0}`) — a concurrent `git stash` from another process never causes the wrong entry to be popped.
- A failed pop **warns only** — the user's tree is recoverable via `git stash list` / `git stash apply`.
- Bypass with `--no-stash` (intended for CI runners that start clean and want loud failure on unexpected dirt).

## Files

- `gitmap/cmd/{as.go, asops.go, releasealias.go, releasealias_git.go, dbmigrate.go}`
- `gitmap/constants/{constants_as.go, constants_releasealias.go, constants_dbmigrate.go}`
- `gitmap/store/migrations.go` (shared helpers: `columnExists`, `tableExists`, `isBenignAlterError`, `logMigrationFailure`)
- `gitmap/helptext/{as.md, release-alias.md, release-alias-pull.md, db-migrate.md}`
- `spec/01-app/98-as-and-release-alias.md`

## Exit codes

- `0` success
- `1` not a git repo / alias not found / dirty tree with `--no-stash` / non-fast-forward pull
