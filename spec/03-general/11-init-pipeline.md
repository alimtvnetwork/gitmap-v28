# init.ps1 / init.sh — repo bootstrap pipeline

Status: normative
Owners: tooling
Related: `visibility-change.ps1`/`.sh`, `fix-repo.ps1`/`.sh`,
spec/04-generic-cli/27-fix-repo-command.md

## Purpose

`init.ps1` (Windows / PowerShell) and `init.sh` (POSIX) are
single-entry-point bootstrap scripts that prepare a freshly cloned
repository for downstream work. They wrap two existing scripts in a
fixed order and forward an auto-yes flag so the pipeline never blocks
on an interactive prompt.

The intent is: "I just cloned this repo. Make it usable in one
command." That means (a) the repo must be public so collaborators
can read it without an auth dance, and (b) any stale `{base}-vN`
tokens left over from a prior version must be rewritten to the
current version.

## Step order

1. **visibility-change** — force `public`, auto-yes.
   - PowerShell: `visibility-change.ps1 -Visible pub -Yes`
   - Shell:      `visibility-change.sh --visible pub --yes`
   - Idempotent: if already public, the underlying script exits 0
     with a `visibility: already public (...)` message.
2. **fix-repo --all** — rewrite every prior `{base}-vN` token to
   the current version.
   - PowerShell: `fix-repo.ps1 -All`
   - Shell:      `fix-repo.sh --all`

The order is **visibility first, fix-repo second**. Rationale: the
visibility flip is a remote-only API call; running it first surfaces
auth / CLI failures immediately, before fix-repo spends time walking
the working tree. fix-repo is a local rewrite and never depends on
remote state, so swapping the order would not improve correctness.

## Auto-yes contract

Both `visibility-change.ps1` (`-Yes`) and `visibility-change.sh`
(`--yes` / `-y`) already accept an auto-yes flag. `init` MUST always
forward it. The init scripts do not expose a flag to disable
auto-yes — interactive confirmation has no place in a bootstrap
pipeline. Users who want the prompt should call `visibility-change`
directly.

## Failure policy: best-effort

Both steps **always run**, even if the first step fails. The init
script:

- Captures each step's exit code.
- Prints a combined `==> init summary` block listing both codes.
- Exits 0 only when **both** steps succeeded.
- Otherwise exits with the **first non-zero** step exit code (so CI
  surfaces the earliest failure, while still showing the user the
  full picture for the second step).

Rationale: `fix-repo` is purely local and may still produce useful
rewrites even when the visibility flip failed (e.g. `gh` not
installed, network down). Running both gives the user the maximum
amount of work-done before they have to debug.

## Flags

| Flag | Behavior |
|------|----------|
| `-DryRun` / `--dry-run` | Forwarded to both inner scripts. No remote API calls; no file writes. |
| `-Help` / `-h` / `--help` | Print usage and exit 0. |

No other flags are accepted. The init scripts deliberately have a
tiny surface — anything more nuanced should be done by calling the
two underlying scripts directly.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Both steps succeeded. |
| 6 | Bad flag passed to `init` itself. |
| _other_ | First non-zero step exit code (forwarded as-is from `visibility-change` or `fix-repo`). |

See the inner scripts' specs for the full code tables:
- visibility-change exit codes: 0/2/3/4/5/6/7/8 (see
  `gitmap-v28/constants/constants_visibility.go`).
- fix-repo exit codes: 0/2/3/4/5/6/7/8/9 (see
  `spec/04-generic-cli/27-fix-repo-command.md`).

## Examples

```
# Standard bootstrap (Windows)
.\init.ps1

# Standard bootstrap (POSIX)
./init.sh

# Preview both steps without touching anything
./init.sh --dry-run
```

## Non-goals

- init does **not** clone. Use `gitmap-v28 clone-fix-repo-pub` (`cfrp`)
  for the clone+fix+publish pipeline that starts from a URL.
- init does **not** commit, push, or open PRs. It is purely a
  preparation step.
- init does **not** expose a `--no-yes` / interactive mode.
  Interactive confirmation belongs to `visibility-change` directly.