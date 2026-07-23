# 31 — `gitmap-v27 open` (alias `op`)

> **Status:** specified — implemented in `gitmap-v27/cmd/open.go`.
> **Schema dependency:** v25 (idempotency stamps on `Repo`).

## Purpose

`gitmap-v27 open` is the "I just `cd`'d into this repo, get me into both
my GUI tools" command. It launches **GitHub Desktop AND VS Code** on
the current repo in a single invocation, replacing the previous
muscle-memory pair of `gitmap-v27 inject .` followed by `code .`.

`open` is intentionally a superset of the historical `inject` flow
when run on the cwd, except it never falls back to "no remote
configured, skipping" silently — non-repo folders still get Desktop
+ VS Code, just no DB upsert.

## Usage

```
gitmap-v27 open              # alias: op
gitmap-v27 open --force      # bypass idempotency stamps; alias: -f
```

The command takes **no positional arguments**. The target is always
resolved from cwd:

1. Try `git rev-parse --show-toplevel` (so a sub-folder still opens
   the repo root).
2. Fall back to `os.Getwd()` (so plain folders still work).

## Steps

1. **Best-effort DB upsert.** Same logic as `inject`: if `git remote
   get-url origin` returns a URL, upsert into `Repo`. Otherwise skip
   silently — Desktop + VS Code still proceed.
2. **GitHub Desktop registration.** Skipped when
   `Repo.LastInjectedDesktopAt` is non-empty AND `--force` is not set.
   When run, the column is stamped to `CURRENT_TIMESTAMP`.
3. **VS Code open.** Same idempotency contract as Desktop, against
   `Repo.LastInjectedVSCodeAt`.
4. **No shell handoff.** Unlike `inject`, `open` does NOT `cd` the
   parent shell anywhere — the user is already there.

## Idempotency

Both stamps live on the `Repo` table:

```
LastInjectedDesktopAt TEXT DEFAULT ''
LastInjectedVSCodeAt  TEXT DEFAULT ''
```

Schema v25 added them via additive `ALTER TABLE`. Empty string ==
"never injected", so legacy rows are correctly treated as
"do everything on the next run".

The check is **per tool, per repo**. If only Desktop has been
injected before, the next `open` will skip Desktop but still open
VS Code, and vice versa. `--force` zeros both gates and re-stamps
both timestamps after the side effects run.

## Output

```
Opening "img-pdf" (D:\wp-work\riseup-asia\img-pdf) in GitHub Desktop and VS Code...
  ↳ github-desktop: already injected (2026-05-15 09:42:11) — pass --force to re-register
  ↳ vscode:         already injected (2026-05-15 09:42:11) — pass --force to re-open
  ✓ open: "img-pdf" ready in both tools
```

With `--force`:

```
Opening "img-pdf" ...
  ⟳ --force: re-injecting img-pdf into both tools
  ✓ open: "img-pdf" ready in both tools
```

## Comparison with related commands

| Command | DB upsert | Desktop | VS Code | Shell `cd` | Idempotency |
|---|---|---|---|---|---|
| `gitmap-v27 inject [path]` | yes (if origin) | yes | yes | yes | per-tool stamp |
| `gitmap-v27 open` | yes (if origin) | yes | yes | no | per-tool stamp |
| `gitmap-v27 code` | no | no | yes | no | none |
| `gitmap-v27 github-desktop` | no | yes | no | no | none |

## Errors

| Condition | Exit | Message |
|---|---|---|
| `os.Getwd()` fails | 1 | `open: ERROR cannot determine current directory: %v` |
| Unknown flag | 2 | (flag.ExitOnError standard) |
| Desktop / VS Code missing | 0 | warning, continue |

## Cross-refs

- `gitmap-v27/cmd/open.go` — implementation.
- `gitmap-v27/cmd/inject.go` — sibling command, shares
  `inject_idempotency.go` helpers.
- `spec/04-generic-cli/29-inject.md` — `inject` semantics and
  resolution rules reused here.
- `mem://features/open-command` — design memory.
