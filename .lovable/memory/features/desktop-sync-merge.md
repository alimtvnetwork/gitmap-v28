---
name: Desktop Sync = GitHub Desktop merge
description: As of v3.37.0 `gitmap desktop-sync` (ds) is an alias of `gitmap github-desktop` (gd). Both follow the same code path, neither requires a prior scan, and `.git` worktree files are detected as repos.
type: feature
---

# Feature: `ds` = `gd` merge (v3.37.0)

**Specs:** `spec/01-app/10-github-desktop.md`, `spec/01-app/11-desktop-sync.md`
**Site routes:** `/github-desktop`, `/desktop-sync`
**Triggered by:** real bug report — `gitmap ds` failed because `.gitmap/output/gitmap.json` didn't exist; `gitmap gd` failed inside a valid git repo because the detector only looked for `.git/` as a directory and missed `.git` files (worktrees / submodules).

## Behaviour summary

- **One command, two names.** `desktop-sync` (`ds`) is now a 4-line shim that calls into the `github-desktop` (`gd`) handler. Same flags, same output, same exit codes.
- **No scan dependency.** Neither command reads `.gitmap/output/gitmap.json`. Removed entirely from the resolution path.
- **GitHub Desktop install check is step 1.** Before any filesystem or DB work, `exec.LookPath("github")` runs. Missing → exit 1 with install hint.
- **`.git` detection fixed.** `git.IsGitRepo(path)` accepts both `.git/` directory **and** `.git` file (worktree / submodule pointer). Was the cause of the false "Not a git repository" error.
- **Resolution order (no args):** install check → CWD is git repo? register it → CWD inside scan root? bulk-register tracked repos under it → otherwise friendly hint, exit 3.
- **`--all` flag** registers every repo in the DB regardless of CWD.

## Exit codes

| 0 | Success / idempotent no-op |
| 1 | GitHub Desktop missing OR explicit path invalid OR one+ failures |
| 2 | Lock contention |
| 3 | Nothing registerable at CWD (not a repo, not a scan root) |
| 4 | Platform unsupported |

## Constants

`constants_cli.go`: `CmdDesktopSync = "desktop-sync"`, `CmdGithubDesktop = "github-desktop"`, aliases `"ds"`, `"gd"`.

`constants_messages.go`: `MsgDesktopNotInstalled`, `MsgDesktopNothingToRegister`, `MsgDesktopRegisteredFmt`, `MsgDesktopAlreadyRegistered`, `MsgDesktopSummaryFmt`.

## Why these decisions

- **Merge instead of deprecate** — chosen by user. Both names are in muscle memory; aliasing keeps everyone happy with zero migration cost.
- **Drop scan dependency entirely** — `ds` was originally a thin "read JSON, register all" command. With the DB now authoritative (since v3.0.0), the JSON middleman is dead weight and confused users who hadn't run `scan` yet.
- **`.git` file detection** — Git worktrees and submodules use a `.git` *file* (pointing at the real gitdir) instead of a `.git` *directory*. `os.Stat` + check `IsDir() || Mode().IsRegular()` is the standard fix.
