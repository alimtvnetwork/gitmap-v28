# Runner Contract — `fix-repo` Sub-command Forwarding

**Version:** 1.0.0
**Updated:** 2026-04-27
**Status:** Normative
**Companion to:** [02-runner-contract.md](./02-runner-contract.md)
**Inner scripts:** `fix-repo.sh`, `fix-repo.ps1`
**Inner spec:** [`spec-authoring/22-fix-repo/01-spec.md`](../../spec-authoring/22-fix-repo/01-spec.md)

---

## 1. Purpose

Define — exactly once, in one place — how the repo-root runners (`run.sh`,
`run.ps1`) MUST forward arguments to the inner `fix-repo` scripts. This
document is the single source of truth for:

- the user-visible CLI surface of `./run.{sh,ps1} fix-repo …`
- the flag mapping between Bash and PowerShell invocations
- the argument-forwarding semantics (verbatim, no re-quoting, no mutation)
- the conformance tests that gate every PR

Any divergence between `run.sh` and `run.ps1` for `fix-repo` is a **bug**.
Any drift between this document and the actual code is a **bug**.

---

## 2. Sub-command surface

| Bash invocation                  | PowerShell invocation             | Effect                                                |
|----------------------------------|-----------------------------------|-------------------------------------------------------|
| `./run.sh fix-repo`              | `./run.ps1 fix-repo`              | Default mode. Rewrite last 2 prior versions.          |
| `./run.sh fix-repo --3`          | `./run.ps1 fix-repo -3`           | Rewrite last 3 prior versions.                        |
| `./run.sh fix-repo --5`          | `./run.ps1 fix-repo -5`           | Rewrite last 5 prior versions.                        |
| `./run.sh fix-repo --all`        | `./run.ps1 fix-repo -all`         | Rewrite every prior version (`v1`..`v(current-1)`).   |
| `./run.sh fix-repo --dry-run`    | `./run.ps1 fix-repo -DryRun`      | Report changes; do not write.                         |
| `./run.sh fix-repo --verbose`    | `./run.ps1 fix-repo -Verbose`     | List every modified file.                             |

Each shell uses its own native flag dialect. The runners are **transparent
dispatchers** — they pass the user's exact tokens to the inner script,
which is the only place that parses flags.

---

## 3. Flag dialect (per shell)

| Semantic     | Bash token (`fix-repo.sh`) | PowerShell token (`fix-repo.ps1`) |
|--------------|----------------------------|-----------------------------------|
| Mode: 2      | `--2`                      | `-2`                              |
| Mode: 3      | `--3`                      | `-3`                              |
| Mode: 5      | `--5`                      | `-5`                              |
| Mode: all    | `--all`                    | `-all` / `-All`                   |
| Dry run      | `--dry-run`                | `-DryRun`                         |
| Verbose      | `--verbose`                | `-Verbose`                        |
| (none)       | (defaults to `--2`)        | (defaults to `-2`)                |

There is **no flag translation** in the runners. The Bash runner forwards
Bash-dialect tokens; the PowerShell runner forwards PowerShell-dialect
tokens. If a token is not recognized by the inner script, the inner
script emits the error — never the runner.

---

## 4. Argument-forwarding contract

The runner MUST forward every argument **byte-for-byte** to the inner
script. Specifically:

1. **No re-quoting.** Use the shell's native argv-array forwarding
   (`"$@"` in Bash, `@args` in PowerShell). Never `eval`, never
   `printf %q`, never join-then-split.
2. **No mutation.** The runner MUST NOT inject, drop, reorder, or
   transform any argument the user typed after the `fix-repo`
   sub-command token.
3. **No intermediate copies.** The dispatch path MUST be a single
   `exec` (Bash) or `&` invocation (PowerShell) on the original argv.
   Function indirection is allowed only when the function does not
   touch the argv (i.e. it just forwards `"$@"` / `@args` unchanged).
4. **Quoted arguments preserved.** Spaces, glob characters, dashes,
   and embedded equals signs MUST survive forwarding intact.
5. **Exit code propagation.** The runner MUST exit with the inner
   script's exit code. Bash uses `exec`; PowerShell uses `exit
   $LASTEXITCODE`.

### 4.1 Bash reference implementation

```bash
fix-repo)
  _assert_fix_repo_present
  shift
  _fix_repo_debug_preflight "$@"
  exec bash "$SCRIPT_DIR/fix-repo.sh" "$@"
  ;;
```

`exec` is required: it replaces the runner process with the inner
script, guaranteeing the exit code is the inner script's exit code and
that no further re-quoting can occur.

### 4.2 PowerShell reference implementation

```powershell
function Assert-FixRepoPresent {
    $inner = Join-Path $PSScriptRoot "fix-repo.ps1"
    if (Test-Path $inner) { return $inner }
    Write-Host "❌ Cannot find $inner" -ForegroundColor Red
    exit 1
}

# dispatch (inlined — no wrapper Invoke-FixRepo function):

"fix-repo" {
    $inner = Assert-FixRepoPresent
    Write-FixRepoDebugPreflight -Inner $inner -Argv $args
    & $inner @args; exit $LASTEXITCODE
}
```

`@args` (splatting) preserves argv as an array — equivalent to Bash's
`"$@"`. The dispatch is a single statement on the original `$args`:
the guard helper only resolves and validates the inner-script path,
and the preflight helper only writes diagnostics to stderr — neither
mutates the argument array. This mirrors the Bash implementation's
`exec bash "$SCRIPT_DIR/fix-repo.sh" "$@"` end-to-end forwarding.

### 4.3 `--debug` preflight

When the forwarded argv contains the literal token `--debug`, both
runners MUST emit a diagnostic block to **stderr** before invoking
the inner script. The block MUST include:

- the runner script path, script directory, and current working dir;
- the resolved inner-script path;
- `ARGC=<n>` (count of forwarded args);
- one `ARG[i]<<value>>` line per arg, in order.

The `--debug` token MUST NOT be consumed: the inner script MUST still
receive byte-identical argv (same count, same order, same values).
Runs without `--debug` MUST NOT print any preflight output.

---

## 5. Help text

Both runners MUST advertise `fix-repo` in their help output (`./run.sh
help`, `./run.ps1 help`) with at least:

```
fix-repo     rewrite prior versioned-repo-name tokens to current

Fix-repo flags forwarded to fix-repo.{sh,ps1}:
  --2 | --3 | --5 | --all   how many prior versions to rewrite (default: --2)
  --dry-run                 report changes; do not write
  --verbose                 list every modified file
```

The help text source-of-truth lives in `scripts/runner-help.txt`
(Bash) and `scripts/runner-help.ps.txt` (PowerShell). Both files MUST
list the same flag set in the same order.

---

## 6. Conformance tests

Every PR that touches `run.sh`, `run.ps1`, `fix-repo.sh`, `fix-repo.ps1`,
or this document MUST keep the following tests green:

| Test                                               | Asserts                                                   |
|----------------------------------------------------|-----------------------------------------------------------|
| `tests/installer/check-fix-repo-runner-wiring.sh`  | Both runners advertise and dispatch the `fix-repo` token. |
| `tests/installer/check-fix-repo-arg-forwarding.sh` | `run.sh` forwards every flag verbatim (shim-based).       |
| `tests/installer/check-fix-repo-url-rewrite.sh`    | End-to-end inner-script behavior (URL rewrite + guard).   |

These run automatically inside `bash tests/installer/run-tests.sh` via
the dynamic `check-*.sh` discovery loop.

---

## 7. Forbidden patterns

The following are **bugs** if found in either runner's `fix-repo`
dispatch path. The CI guard
(`linter-scripts/check-runner-dispatch-antipatterns.sh`, wired as
step 13 of `scripts/lint-ci.sh`) enforces every entry below.

### 7.1 Bash (`run.sh`) — forbidden in the `fix-repo)` arm

| # | Pattern                                                     | Why it's a bug                                                |
|---|-------------------------------------------------------------|---------------------------------------------------------------|
| 1 | `"$*"` or bare `$*`                                         | Joins argv on `IFS`; collapses N args into one.               |
| 2 | bare `$@` (outside double quotes)                           | Re-splits each arg on `IFS`.                                  |
| 3 | `eval ... fix-repo.sh ...`                                  | Re-parses every argument; quote-unsafe.                       |
| 4 | `bash -c "..."` / `sh -c "..."` wrapper                     | Loses argv boundaries; needs manual re-quoting.               |
| 5 | `${@:N}` slicing                                            | Drops original argv; forward `"$@"` verbatim.                 |
| 6 | `printf '%q ' "$@"` / `printf '%s ' "$@"`                   | Rebuilds argv from a joined string; quote drift.              |
| 7 | `IFS=…` mutation in the dispatch arm                        | Alters how subsequent expansions split argv.                  |
| 8 | Command substitution that contains `"$@"` (`$(... "$@" ...)`) | Stringifies argv via subshell stdout.                       |
| 9 | `xargs` in the dispatch path                                | Reformats argv via stdin; loses quoting.                      |
|10 | Trailing `&` (background dispatch)                          | Detaches the child; runner cannot propagate the exit code.    |
|11 | `fix-repo.sh ... \| ...` (pipe after inner call)            | Masks inner exit code (last pipe stage wins by default).      |

**Required:** the arm MUST contain `exec ... fix-repo.sh "$@"`.

### 7.2 PowerShell (`run.ps1`) — forbidden in the `"fix-repo"` arm

| # | Pattern                                                     | Why it's a bug                                                |
|---|-------------------------------------------------------------|---------------------------------------------------------------|
| 1 | `$args -join …` / `-join $args`                             | Collapses argv array into a single string.                    |
| 2 | `"$args"` (interpolation in a double-quoted string)         | Flattens argv via `$OFS`; lossy. Use `@args` splatting.       |
| 3 | `[string]::Join(...)` on `$args`                            | Same array-to-string flattening.                              |
| 4 | `$args.ToString()` / `$args -as [string]` / `[string]$args` | Implicit array-to-string conversions; flatten argv.           |
| 5 | `$args[N..M]` slicing                                       | Drops original argv; forward `@args` verbatim.                |
| 6 | `Invoke-Expression …` (or alias `iex`) on argv              | Re-parses argv as a script; unsafe and quote-lossy.           |
| 7 | `cmd /c "..."` wrapper                                      | Re-parses argv under cmd.exe quoting rules.                   |
| 8 | `Start-Process …` for the inner script                      | Detaches the child; cannot propagate `$LASTEXITCODE`.         |
| 9 | `Start-Job …` for the inner script                          | Same — runs out of process; exit code lost.                   |

**Required:** the arm MUST invoke the inner script with `@args`
splatting AND end with `exit $LASTEXITCODE`.

### 7.3 Out-of-scope (intentionally NOT flagged)

- `eval`, `Invoke-Expression`, `Start-Process` etc. in **other**
  branches or helper functions. The guard scans only the `fix-repo`
  dispatch arm. False positives in unrelated code are by design out of
  scope; if those locations need rules, add a separate guard.
- Conditional flag rewriting (e.g. mapping `-DryRun` → `--dry-run`
  inside the runner). Not lint-checked, but still a spec violation
  per §4.2 — such mapping belongs in the inner script.

---

## 8. Change-control

Any change to the flag set, mapping, or forwarding semantics requires:

1. Updating §2, §3, or §4 of this document **first**.
2. Updating both runners and `fix-repo.{sh,ps1}` in the same PR.
3. Updating the help-text files in `scripts/runner-help{,.ps}.txt`.
4. Updating or adding conformance tests in `tests/installer/`.

The sub-command spec ([`spec-authoring/22-fix-repo/01-spec.md`](../../spec-authoring/22-fix-repo/01-spec.md))
remains the authority on inner-script behavior; this document is the
authority only on **how the runners reach it**.
