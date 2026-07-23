# Desktop Sync (`ds` = `gd`)

> **As of v3.37.0** `desktop-sync` (`ds`) is an **alias** of `github-desktop` (`gd`).
> Both commands do the same thing, follow the same code path, and accept the
> same arguments. The two names are kept only because both are in muscle memory.

## What it does

Registers git repositories with **GitHub Desktop** so they appear in its
sidebar. It does **not** depend on a previous `gitmap scan`. It does **not**
read `.gitmap/output/gitmap.json`. It just talks to GitHub Desktop directly.

## Usage

```bash
gitmap gd                       # register CWD (or every DB-tracked repo under CWD's scan root)
gitmap ds                       # identical — alias of gd
gitmap gd D:\path\to\repo       # register an explicit folder
gitmap gd --all                 # register every repo currently in the gitmap DB
```

## Resolution order (cwd mode, no path arg)

1. **Is GitHub Desktop installed?** Look up the `github` CLI on `PATH`.
   If missing → exit 1 with install hint. **This check happens first**, before
   any filesystem or DB work.
2. **Is CWD itself a git repo?** Detect via either:
   - `.git/` directory present, **or**
   - `.git` file present (worktree / submodule pointer — the previous
     implementation missed this case and produced false "Not a git
     repository" errors on perfectly valid worktrees).
   If yes → register CWD only, exit 0.
3. **Is CWD a registered ScanFolder root (or descendant of one)?**
   If yes → bulk-register every repo from the DB whose path is under that
   scan root. No re-walk, no JSON file required.
4. **Neither?** Print a friendly hint:
   ```
   ✗ Nothing to register here.
     Run `gitmap scan .` to track this folder, or
     pass an explicit path: `gitmap gd <path>`.
   ```
   Exit 3.

## Prerequisites

- **GitHub Desktop installed** with the `github` CLI on `PATH` (Windows
  installers add this automatically).
- A git repository (CWD, explicit path, or DB-tracked repos under CWD).
- **No prior `gitmap scan` is required** for the cwd-single-repo and explicit-path modes.

## Behaviour matrix

| Scenario | Result |
|----------|--------|
| GitHub Desktop missing | Exit 1, print install hint, do nothing else |
| CWD is a git repo (`.git` dir or file) | Register CWD, done |
| CWD is a scan root with N tracked repos | Sequentially register all N |
| Explicit path arg, valid repo | Register that path |
| Explicit path arg, not a repo | Exit 1 with clear error |
| `--all` flag | Register every repo in DB regardless of CWD |
| One repo fails mid-batch | Continue, summarise at end, exit 1 if any failed |

## Output

```
gitmap ds
  ✓ Registered: macro-ahk
GitHub Desktop: 1 added · 0 skipped · 0 failed
```

```
gitmap gd            # inside a scan root
  [1/14] my-api ............ ✓
  [2/14] web-app ........... ✓
  [3/14] billing-svc ....... already registered
  ...
GitHub Desktop: 12 added · 2 skipped · 0 failed
```

## Exit codes

| 0 | All registrations succeeded (or were idempotent no-ops) |
| 1 | GitHub Desktop missing, explicit path invalid, or one+ failures |
| 2 | Lock contention |
| 3 | CWD has no repo and is not a scan root |
| 4 | Platform unsupported (non-Windows / non-macOS) |

## What changed from pre-v3.37.0

| Before | After |
|--------|-------|
| `ds` required `.gitmap/output/gitmap.json` to exist | No JSON dependency |
| `ds` and `gd` were two different code paths | One code path, `ds` is an alias |
| `gd` only checked `.git` as a directory | Detects `.git` file (worktrees) too |
| GitHub Desktop install check happened mid-flow | Install check is the first step |

## Implementation notes

| File | Responsibility |
|------|----------------|
| `cmd/githubdesktop.go` | Single command handler; `cmd/desktopsync.go` is now a 4-line shim that calls into it |
| `desktop/desktop.go` | `IsInstalled()`, `RegisterRepo(path)`, batch helpers |
| `git/detect.go` | `IsGitRepo(path)` checks both `.git/` dir and `.git` file |
| `constants/constants_cli.go` | `CmdDesktopSync = "desktop-sync"`, `CmdGithubDesktop = "github-desktop"`, aliases `ds`/`gd` |

## See also

- [scan gd (spec 102)](/scan-gd) — bulk register from DB without re-walking
- [scan all (spec 100)](/scan-all)
- [pull all (spec 101)](/pull-all)
