# Repo Aliases & Remote Release — `gitmap-v28 as` / `release-alias` / `release-alias-pull`

> **Status:** Implemented in v2.97.0; auto-stash semantics finalised in v2.99.0.
> **Related specs:**
> - [97-move-and-merge.md](97-move-and-merge.md) — sibling family of remote-aware commands
> - [05-cloner.md](05-cloner.md) — repo discovery primitives reused by `as`
> - [02-cli-interface.md](02-cli-interface.md) — global flag conventions (`-y`, `--dry-run`)
> - [16-cicd.md](../12-consolidated-guidelines/16-cicd.md) — release pipeline this command drives

## Overview

Three CLI verbs that decouple the "I am sitting in a repo" step from
the "I want to release a repo" step:

```
gitmap-v28 as          [alias-name]            # tag the CURRENT repo with an alias (run inside it)
gitmap-v28 release-alias <alias> <version>     # release a previously-aliased repo from anywhere
gitmap-v28 release-alias-pull <alias> <ver>    # pull --ff-only first, then release
```

The pair lets a developer register a repo once with `as`, then trigger
its full release pipeline (lint → test → tag → push → assets) from any
working directory — including from a CI runner that never `cd`s into the
repo.

| Verb | Aliases |
|------|---------|
| `as` | `s-alias` |
| `release-alias` | `ra` |
| `release-alias-pull` | `rap` |

`rap` is a thin alias for `ra --pull`. The flag is canonical;
the verb is sugar for users who want a single token.

---

## Command: `as`

Tags the **current** Git repository with a short alias and records it
in the active-profile SQLite database.

### Behaviour

1. Resolve the repo top level via `git rev-parse --show-toplevel`.
   Abort if the CWD is not inside a Git repository.
2. Build a `ScanRecord` for that single repo using the same
   `mapper.BuildRecords()` path that `gitmap-v28 scan` uses
   — guarantees the upserted row matches the schema other commands
   already understand.
3. Upsert the record into the `Repos` table (so `gitmap-v28 list`,
   `status`, `pull`, etc. immediately see it).
4. Map `alias-name → Repos.Id` in the alias store. When `alias-name`
   is omitted the repo folder basename is used.
5. Refuse to clobber an existing alias unless `--force` (`-f`) is
   passed.

### Usage

```
gitmap-v28 as                       # alias defaults to filepath.Base(repo-root)
gitmap-v28 as project-x             # explicit alias
gitmap-v28 as project-x --force     # overwrite an existing project-x alias
gitmap-v28 s-alias project-x        # long-form alias for the verb
```

### Flags

| Flag | Default | Meaning |
|------|---------|---------|
| `--force` / `-f` | false | Replace an alias that already points elsewhere. |

### Errors

| Condition | Exit | Message |
|-----------|------|---------|
| CWD not inside a Git repo | 1 | `error: 'as' must be run from inside a Git repo (cwd: <path>)` |
| `mapper.BuildRecords` returned no records | 1 | `error: could not resolve repo metadata for <path>: <reason>` |
| Alias exists, no `--force` | 1 | `error: alias '<name>' already maps to <other-path>; pass --force to overwrite` |
| More than one positional arg | 2 | `usage: gitmap-v28 as [alias-name] [--force]` |

---

## Command: `release-alias` (`ra`)

Releases the repo behind an alias using the existing `runRelease`
pipeline — without requiring the user to `cd` into the repo first.

### Step-by-step

1. Resolve `<alias>` → absolute path via the alias store. Abort with
   a hint if the alias is unknown.
2. `os.Chdir(target)` and `defer os.Chdir(originalDir)` so the rest
   of the pipeline runs in the repo root.
3. (Optional) `git pull --ff-only` if `--pull` is set or the verb
   was `rap`. Hard-fail on pull conflict — releasing on top of a
   non-fast-forward is unsafe.
4. **Auto-stash** dirty working trees (see semantics below) unless
   `--no-stash` is passed.
5. Invoke `runRelease(<version> [--dry-run])` — the same entry
   point used by `gitmap-v28 release` from inside the repo.
6. **Pop the auto-stash** on the way out, even on release failure.

### Usage

```
gitmap-v28 release-alias project-x v1.2.0
gitmap-v28 ra              project-x v1.2.0
gitmap-v28 ra              project-x v1.2.0 --pull --dry-run
gitmap-v28 release-alias-pull project-x v1.2.0     # equivalent to: ra ... --pull
gitmap-v28 rap             project-x v1.2.0
```

### Flags

| Flag | Default | Meaning |
|------|---------|---------|
| `--pull` | false (true when invoked as `rap`) | Run `git pull --ff-only` before releasing. |
| `--no-stash` | false | Abort instead of auto-stashing a dirty tree. |
| `--dry-run` | false | Forward `--dry-run` to `runRelease`; no commits, tags, or pushes. |

### Errors

| Condition | Exit | Message |
|-----------|------|---------|
| Alias not found | 1 | `error: unknown alias '<name>'. Run 'gitmap-v28 as <name>' inside the repo first, or 'gitmap-v28 alias list' to see registered aliases.` |
| `os.Chdir` failed | 1 | `error: cannot chdir into '<path>': <reason>` |
| `git pull --ff-only` failed | 1 | `error: pull failed in '<path>': <reason>` |
| Wrong arg count | 2 | `usage: gitmap-v28 release-alias <alias> <version> [--pull] [--no-stash] [--dry-run]` |
| Inner `runRelease` failure | propagated | (handled by `runRelease`) |

---

## Auto-Stash Semantics

`release-alias` always wants a clean working tree before tagging. Rather
than refusing to release dirty repos, it auto-stashes the working tree
in a labelled stash, runs the release, then pops the stash on exit:

```
+-------------------+   isWorkingTreeDirty?   +------------------+
|  CWD = repo root  | ----------------------> | git status -s    |
+-------------------+                         +------------------+
         |                                            |
         | dirty                                      | clean
         v                                            v
+----------------------+                  +----------------------+
| git stash push       |                  |  no stash created    |
|   --include-untracked|                  +----------+-----------+
|   -m "<label>"       |                             |
+----------+-----------+                             |
           |                                         |
           v                                         |
+--------------------------+                         |
|  runRelease(version)     | <-----------------------+
+----------+---------------+
           |
           v   (defer)
+--------------------------+
| git stash pop <stash@N>  |  <-- located by label match in
+--------------------------+      `git stash list`
```

### Stash label format

```
gitmap-release-alias autostash <alias>-<version>-<unix-ts>
```

The unix-timestamp suffix guarantees uniqueness even when two parallel
release-alias runs target the same alias and version.

### Pop semantics

- `popAutoStash` is registered with `defer` BEFORE `runRelease` runs,
  so it always fires — including when the release pipeline aborts.
- The stash is located by label match against `git stash list` (not
  by stash index), so a concurrent `git stash` from another process
  does not cause us to pop the wrong entry.
- A failed pop **warns only** — the user's working tree is still
  recoverable via `git stash list` / `git stash apply`.

### Bypass

| Flag | Effect |
|------|--------|
| `--no-stash` | Skip the dirty-check; do not stash. The release will fail if the repo is dirty (the inner `runRelease` enforces clean). |

`--no-stash` is intended for CI runners that always start from a clean
checkout and want to fail loudly on unexpected dirt.

---

## Dispatcher Wiring

All three verbs route through the data-domain dispatcher chain in
`gitmap-v28/cmd/rootcore.go` and `gitmap-v28/cmd/rootrelease.go`:

```
main() -> Dispatch()
        -> tryDataCommands()                        // rootdata.go
            -> CmdAs / CmdAsAlias       -> runAs()
            -> CmdDBMigrate / CmdDBMigrateAlias -> runDBMigrate()
        -> tryReleaseCommands()                     // rootrelease.go
            -> CmdReleaseAlias / CmdRA  -> runReleaseAlias(args, false)
            -> CmdRAPull / CmdRAP       -> runReleaseAlias(args, true)
```

### Constants

| Constant | File | Value |
|----------|------|-------|
| `CmdAs` | `constants_as.go` | `"as"` |
| `CmdAsAlias` | `constants_as.go` | `"s-alias"` |
| `CmdReleaseAlias` | `constants_releasealias.go` | `"release-alias"` |
| `CmdRA` | `constants_releasealias.go` | `"ra"` |
| `CmdRAPull` | `constants_releasealias.go` | `"release-alias-pull"` |
| `CmdRAP` | `constants_releasealias.go` | `"rap"` |
| `FlagRAPull` | `constants_releasealias.go` | `"pull"` |
| `FlagRANoStash` | `constants_releasealias.go` | `"no-stash"` |
| `FlagRADryRun` | `constants_releasealias.go` | `"dry-run"` |
| `FlagAsForce` | `constants_as.go` | `"force"` |
| `FlagAsForceS` | `constants_as.go` | `"f"` |

### Completion

Both constants files carry the `// gitmap-v28:cmd top-level` marker so the
generator at `gitmap-v28/completion/internal/gencommands/main.go` picks
every command + alias up automatically. After adding a verb, run:

```
cd gitmap-v28 && go generate ./completion/...
```

The CI `generate-check` job (`.github/workflows/ci.yml`) fails the
build if `allcommands_generated.go` drifts.

---

## Files Involved

| File | Role |
|------|------|
| `gitmap-v28/cmd/as.go` | `runAs`, arg parsing, `git rev-parse` lookup. |
| `gitmap-v28/cmd/asops.go` | `upsertSingleRepo`, `registerAlias`. |
| `gitmap-v28/cmd/releasealias.go` | `runReleaseAlias`, arg parsing, dispatcher into `runRelease`. |
| `gitmap-v28/cmd/releasealias_git.go` | `runReleaseAliasPull`, `autoStashIfDirty`, `popAutoStash`, `findStashIndex`. |
| `gitmap-v28/cmd/rootdata.go` | Dispatch for `as`, `s-alias`, `db-migrate`, `dbm`. |
| `gitmap-v28/cmd/rootrelease.go` | Dispatch for `release-alias`, `ra`, `release-alias-pull`, `rap`. |
| `gitmap-v28/constants/constants_as.go` | Cmd / flag / message constants for `as`. |
| `gitmap-v28/constants/constants_releasealias.go` | Cmd / flag / message constants for `release-alias` family. |
| `gitmap-v28/helptext/as.md` | `gitmap-v28 as --help` content. |
| `gitmap-v28/helptext/release-alias.md` | `gitmap-v28 ra --help` content. |
| `gitmap-v28/helptext/release-alias-pull.md` | `gitmap-v28 rap --help` content. |

---

## Examples

```
# 1. Sit in the repo once, register the alias.
cd /code/project-x
gitmap-v28 as                                # alias = "project-x" (basename)
gitmap-v28 as px                             # explicit alias
gitmap-v28 alias list                        # confirm

# 2. From anywhere, release it.
cd ~                                     # any directory
gitmap-v28 ra px v1.2.0
gitmap-v28 ra px v1.2.0 --dry-run            # preview without tagging
gitmap-v28 rap px v1.2.0                     # pull --ff-only first

# 3. Dirty tree? auto-stash kicks in.
cd /code/project-x && echo dirt > scratch.txt
cd ~ && gitmap-v28 ra px v1.3.0
#   ▸ stashed: gitmap-release-alias autostash px-v1.3.0-1729400123
#   ▸ release v1.3.0 ... ok
#   ▸ popped stash: gitmap-release-alias autostash px-v1.3.0-1729400123

# 4. Refuse the auto-stash on a CI runner.
gitmap-v28 ra px v1.3.0 --no-stash           # exits non-zero if dirty

# 5. Reassign an alias that already exists.
cd /code/project-x-v2
gitmap-v28 as px --force
```

---

## Constraints

- `as` MUST refuse to run outside a Git repository — there is no
  meaningful repo to alias. Exit code 1, message includes the CWD.
- `release-alias` MUST chdir into the resolved path before invoking
  `runRelease`; the inner pipeline expects CWD == repo root.
- The `defer os.Chdir(originalDir)` MUST be registered before the
  `runRelease` call, so process state is restored even on panic.
- Auto-stash labels MUST include the unix timestamp suffix to survive
  parallel `release-alias` invocations on the same alias + version.
- `popAutoStash` MUST locate the stash by label match (not by
  `stash@{0}`), to avoid popping an unrelated stash created by
  the user or a sibling process during the release window.
- A failed `git stash pop` MUST warn but not exit non-zero — the
  release itself succeeded, and the stash is still recoverable.
- `--pull` and `--no-stash` MUST be honoured even when the verb form
  (`rap`) implies one of them; the flag is the canonical truth.

---

## Acceptance Checklist

- [x] `gitmap-v28 as` inside a repo registers the basename as alias.
- [x] `gitmap-v28 as <name>` registers an explicit alias.
- [x] `gitmap-v28 as <name> --force` overwrites an existing alias.
- [x] `gitmap-v28 as` outside a Git repo exits 1 with a CWD-aware message.
- [x] `gitmap-v28 ra <alias> <ver>` releases from anywhere.
- [x] `gitmap-v28 rap <alias> <ver>` pulls then releases (equivalent to `ra --pull`).
- [x] Dirty tree triggers labelled `git stash push --include-untracked`.
- [x] Stash is popped on the way out, even when `runRelease` aborts.
- [x] `--no-stash` skips stashing; release fails fast on dirty repo.
- [x] `--dry-run` is forwarded to `runRelease`.
- [x] Unknown alias exits 1 with a `gitmap-v28 as ...` hint.
- [x] `pull --ff-only` failure exits 1 — never releases on a non-FF tree.
- [x] All four verb tokens (`as`, `s-alias`, `release-alias`, `ra`,
      `release-alias-pull`, `rap`) appear in `allcommands_generated.go`.

> **Implementation:** v2.97.0 — `gitmap-v28/cmd/{as.go, asops.go,
> releasealias.go, releasealias_git.go}`, `constants/constants_as.go`,
> `constants/constants_releasealias.go`. Auto-stash defer-pop hardened
> in v2.98.0; helptext + dispatcher coverage finalised in v2.99.0.
