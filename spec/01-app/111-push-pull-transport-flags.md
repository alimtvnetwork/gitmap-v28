# 111 — `gitmap push` / `gitmap pull` transport flags (`--ssh` / `--https`)

Status: implemented in v5.29.0
Related: spec/01-app/110-clone-ssh-flag.md

## Summary

Two new affordances:

1. **`gitmap push`** — a `git push` wrapper that runs in the current
   working directory, identical in spirit to the `gitmap pull` cwd
   short-circuit added in v5.28.0.
2. **`--ssh` / `--https` transport flags** on both `gitmap push` and
   `gitmap pull` that rewrite the `remote.origin.url` of the cwd repo
   to the requested transport BEFORE the git operation runs.

The conversion is **persistent** — it calls
`git remote set-url origin <converted>` so subsequent `git push` /
`git pull` invocations (with or without gitmap) keep using the new
transport. This matches the user's mental model: "switch this repo to
SSH and push".

## Flags

| Flag | Aliases | Meaning |
|------|---------|---------|
| `--ssh` | `-ssh`, `--sh` | Force `remote.origin.url` to `git@host:owner/repo.git`. |
| `--https` | `-https`, `--ht` | Force `remote.origin.url` to `https://host/owner/repo.git`. |

Mutually exclusive — when both are supplied, `--ssh` wins and a
one-line warning is written to stderr (mirrors `clone` semantics from
spec 110).

When neither flag is set, behavior is unchanged: `push` / `pull` run
against whatever `remote.origin.url` is currently set.

## Behavior contract

1. Detect git work tree via `git rev-parse --is-inside-work-tree`.
   When not a git repo, exit with a clear error.
2. Read `git config --get remote.origin.url`. When the remote is
   unset, exit with a clear error pointing at `git remote add origin`.
3. Run the converter from spec 110 (`ConvertURLToSSH` /
   `ConvertURLToHTTPS`). If the URL is unrecognised, print a warning
   and skip the rewrite (run `git push`/`pull` against the original
   URL — fail-open, not fail-closed).
4. If the converted URL differs from the current one, run
   `git remote set-url origin <converted>` and print
   `→ remote.origin.url: <old> → <new>`.
5. Stream `git push` / `git pull` with stdin/stdout/stderr forwarded;
   propagate the underlying exit code.
6. Any positional arguments after the flags are forwarded verbatim to
   `git push` / `git pull` (e.g. `gitmap push --ssh origin main`).

## Why persist the rewrite

The alternative — rewriting on every invocation without persisting —
would leave a foot-gun: any non-gitmap `git push` would silently use
the stale transport. Persisting via `git remote set-url` is a single
local config write, identical to what users would type manually, and
keeps `gitmap` honest as a thin convenience layer over git.

## Wiring

- New module `gitmap/cmd/remotetransport.go`:
  - `ApplyTransportFlag(dir string, useSSH, useHTTPS bool) (changed bool, err error)`
  - Internal helpers: `currentOriginURL(dir)`, `setOriginURL(dir, url)`.
- New module `gitmap/cmd/push.go`:
  - `runPush(args []string)` — parses flags via `parseTransportFlags`,
    short-circuits to `git push` in cwd, propagates exit code.
- Updated `gitmap/cmd/pull.go::runPullCWD`:
  - Accepts the same flags; applies transport rewrite before
    `git pull`. Existing non-cwd batch behavior unchanged.
- Updated dispatch table `gitmap/cmd/rootcore.go`:
  - Adds `{CmdPush, CmdPushAlias}` route.
- New constants in `gitmap/constants/constants_cli.go`:
  - `CmdPush = "push"`, `CmdPushAlias = "ph"` (NOT `p` — collides
    with `CmdPullAlias`).

## E2E test plan

`gitmap/cmd/pushpull_transport_e2e_test.go` exercises the full
loop using temp directories, a local bare repo as origin, and the
real `git` binary (skipped when `git` is missing from PATH):

1. **Push HTTPS → SSH rewrite persists**:
   - Create bare repo at `<tmp>/origin.git`.
   - Clone with synthetic HTTPS URL, then `git remote set-url origin https://github.com/acme/widgets.git`.
   - Call `ApplyTransportFlag(dir, true, false)`.
   - Assert `remote.origin.url` is now `git@github.com:acme/widgets.git`.
2. **Pull SSH → HTTPS rewrite persists** — symmetric.
3. **No-op when already correct** — `changed=false`, no config write.
4. **Mutual exclusion** — both flags set ⇒ SSH wins.
5. **Unrecognised URL** — `changed=false`, no error (fail-open).

The actual `git push` call is exercised by a small wrapper test that
runs against a local bare remote so no network is required.

## Help text

`gitmap/helptext/push.md` (new) and `gitmap/helptext/pull.md` (edit)
document the flags with realistic simulations following the
120-line / 3–8-line conventions.

## Out of scope

- Per-record transport override for manifest-mode pulls (`json` /
  `csv` / `text`). Tracked separately under spec 110's follow-ups.
- Rewriting remotes other than `origin` (e.g. `upstream`). May be
  added later via `--remote <name>`.
