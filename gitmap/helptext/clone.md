# gitmap clone

Clone repositories from a structured output file (JSON, CSV, or text),
or clone a single repository directly from a Git URL.

## Alias

c

## Usage

    gitmap clone <source|json|csv|text|url> [folder] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --target-dir \<dir\> | current directory | Base directory for clones |
| --safe-pull | false | Pull existing repos with retry + diagnostics |
| --github-desktop | false | Auto-register with GitHub Desktop (no prompt) |
| --audit | false | Validate planned git clone commands and print a diff-style summary; never executes |
| --max-concurrency \<N\> | 0 (auto = NumCPU) | Run up to N clones in parallel. `1` = sequential, `0` = auto (one worker per CPU). Hierarchy is preserved at any N. |
| --default-branch \<name\> | (none) | Fallback branch name when HEAD/remote-tracking detection finds nothing. Rewrites `branchSource=detached \| unknown \| empty` rows so they go through the trusted `git clone -b <name>` path instead of relying on the remote's default HEAD. Empty preserves legacy behavior. |
| --verbose | false | Write detailed debug log |
| --output \<mode\> | (off) | `terminal` streams a standardized branch/from/to/command block to **stdout** immediately before each URL's `git clone`. Git progress and the per-repo summary stay on **stderr**. Pipe stdout to capture just the previews. |
| --no-vscode-sync | false | Skip writing every successfully-cloned repo into the VS Code Project Manager `projects.json`. Default is to sync (one entry per cloned repo, additive — never clobbers existing entries). Use this in CI / headless environments or when VS Code is not installed. |
| --ssh | false | Force every clone URL into `git@host:owner/repo.git` SSH-shorthand form before `git clone` runs. Converts `https://…` and `ssh://git@…` URLs automatically. Mutually exclusive with `--https`. |
| --https | false | Force every clone URL into `https://host/owner/repo.git` form. Converts SSH-shorthand and `ssh://…` URLs automatically. Useful in CI where the SSH agent isn't unlocked. |

## VS Code Project Manager sync

After every successful clone, gitmap upserts an entry into the
alefragnani.project-manager `projects.json` so the new folder appears
in VS Code's Project Manager sidebar immediately. URL-mode clones sync
the single resulting folder; manifest-mode clones (`json` / `csv` /
`text`) sync every repo that landed on disk in one batched pass. Pass
`--no-vscode-sync` to opt out — gitmap then prints
`• VS Code Project Manager sync skipped (--no-vscode-sync)` and exits
with the same code as a normal run. The flag is a no-op when the VS
Code extension directory is absent (a warning is logged, exit code is
unchanged).

### Windows path canonicalization & EvalSymlinks soft-fail

Every `rootPath` written to `projects.json` is normalized via
`filepath.Clean` (collapses mixed `/` `\` separators, redundant `.`,
and trailing separators) and then `filepath.EvalSymlinks` (resolves
symlinks AND Windows 8.3 short names like `C:\PROGRA~1\…` to the
canonical long form). On Windows, the dedup key is also
case-insensitive — `C:\Foo` and `c:\foo` collapse to one entry.

**Soft-fail rule**: when `EvalSymlinks` errors (path not yet on
disk, permission denied, broken symlink), gitmap falls back to the
**cleaned absolute path** and still writes the `projects.json`
entry — never failing the clone over a symlink-resolution issue. The
next `gitmap scan` re-canonicalizes once the path is resolvable.

Pass `--debug-paths` to emit one `[debug-paths] in=… clean=…
resolved=…` line per canonicalize call on **stderr** (or set
`GITMAP_DEBUG_PATHS=1` directly in CI). See the README "Windows path
canonicalization & EvalSymlinks soft-fail" subsection for the full
failure-mode catalogue.

## Hierarchy preservation

Every record is cloned into `<target-dir>/<RelativePath>` exactly as
captured by `gitmap scan` — no flattening, no path rewriting. A scan
that recorded `group-a/sub/repo-x` reproduces `<target>/group-a/sub/repo-x`
on clone. This holds for every runner mode (sequential, parallel) and
for every input format (CSV, JSON, text).

## Parallel execution (`--max-concurrency`)

`gitmap clone` defaults to `--max-concurrency 0`, which auto-resolves
to `runtime.NumCPU()` workers at startup. Pass `1` to force the legacy
sequential runner (useful when scripts grep stderr for monotonically
ordered per-repo lines), or any positive N to pin the pool size:

    gitmap clone json --max-concurrency 8   # pin to 8 workers
    gitmap clone json --max-concurrency 1   # force sequential
    gitmap clone json                       # auto = NumCPU

When parallel mode is active gitmap prints a single header line
(`↪ parallel clone enabled: <N> workers`) before the per-repo lines
so you know which value was resolved. Progress lines arrive in
completion order rather than input order; the on-disk hierarchy is
unaffected. The clone-cache fingerprint, audit short-circuit, and
safe-pull retry behaviour all operate identically regardless of N.

The same `--max-concurrency` flag — same name, same semantics, same
auto-default — is also accepted by `gitmap clone-next`, `gitmap
clone-now` (`relclone`), and `gitmap clone-from`. The single-repo
`gitmap clone-pick` deliberately omits it (only one URL to clone).

## Audit mode

`gitmap clone --audit <source>` parses the manifest, computes the exact
`git clone` / `git pull` command that would run for every record, and
prints a diff-style report. It never invokes git, never writes outside
stdout, and works offline. Useful for reviewing a manifest before a
batch clone or for CI dry-runs against a generated `.gitmap/output/`.

Markers:

| Marker | Action   | Meaning                                          |
|--------|----------|--------------------------------------------------|
| `+`    | clone    | target missing — would run `git clone ...`       |
| `~`    | pull     | target is an existing git repo — would safe-pull |
| `=`    | cached   | clone-cache fingerprint matches local HEAD       |
| `?`    | conflict | target exists but is not a git repository        |
| `!`    | invalid  | record has no clone URL                          |

Audit requires a manifest source (`json`, `csv`, `text`, or a path) —
it cannot run against a single direct URL.

## Prerequisites

- For file-based clone: run `gitmap scan` first to generate output files
- For URL clone: just provide the HTTPS or SSH URL

## Branch selection strategy

`gitmap clone` decides whether to pass `-b <branch>` to `git clone` based on
each record's `branchSource` (captured during scan):

| branchSource     | Behavior                                                   |
|------------------|------------------------------------------------------------|
| `HEAD`           | Checkout the recorded branch (`-b <branch>`).              |
| `remote-tracking`| Checkout the recorded tracking branch (`-b <branch>`).     |
| `default`        | Checkout the recorded repo default (`-b <branch>`).        |
| `detached`       | Omit `-b`; let the remote's default HEAD decide.           |
| `unknown` / empty| Omit `-b`; let the remote's default HEAD decide.           |

This prevents "Remote branch not found" errors when a scan captured a
detached HEAD or a literal `HEAD` value that cannot be checked out.

Pass `--default-branch <name>` to override the last two rows: every
record that would otherwise omit `-b` is rewritten with
`branchSource=default` and `branch=<name>`, so the cloner emits
`git clone -b <name> <url> <dest>`. A breadcrumb (`default-branch
fallback applied: <name>`) is appended to the record's Notes so
`--audit` and `CloneResult.Notes` make the substitution traceable.
The flag uses the same constant as `gitmap scan --default-branch`
to guarantee identical wording across both surfaces.

## Idempotent clone cache

Repeated `gitmap clone` runs are idempotent. After each successful clone or
pull, gitmap writes a fingerprint (URL, branch, local HEAD SHA, remote HEAD
SHA) to `<target-dir>/.gitmap/clone-cache.json`. On the next run, repos
whose local HEAD still matches the cached SHA — and whose remote tip has
not advanced — are reported as `skipped (cached)` instead of being
re-cloned or re-pulled. If the remote is unreachable (offline), gitmap
trusts the cache as long as the local HEAD still matches.

Delete the cache file to force a full reclone of every entry.

## Examples

### Example 1: Clone from a direct URL (versioned — auto-flattened)

    gitmap clone https://github.com/alimtvnetwork/wp-onboarding-v13.git

**Output:**

    Cloning wp-onboarding-v13 into wp-onboarding...
    Cloned wp-onboarding-v13 successfully.
      + 1 repo(s) added to GitHub Desktop, 0 failed.
      Opening wp-onboarding in VS Code...
      VS Code opened.

### Example 2: Clone URL into a custom folder

    gitmap clone https://github.com/alimtvnetwork/wp-alim.git "my-project"

**Output:**

    Cloning wp-alim into my-project...
    Cloned wp-alim successfully.
      + 1 repo(s) added to GitHub Desktop, 0 failed.
      Opening my-project in VS Code...
      VS Code opened.

### Example 3: Clone from JSON output

    gitmap clone json --target-dir D:\projects

**Output:**

    Cloning from .gitmap/output/gitmap.json...
    [1/12] Cloning my-api... done
    [2/12] Cloning web-app... done
    ...
    Clone complete: 12 succeeded, 0 failed

### Example 4: Clone with safe-pull for existing repos

    gitmap c csv --safe-pull

**Output:**

    [1/8] my-api exists -> pulling... Already up to date.
    [2/8] web-app exists -> pulling... Updated (3 new commits)
    [3/8] Cloning billing-svc... done
    ...
    Clone complete: 8 succeeded, 0 failed

### Example 5: Clone from text file with verbose logging

    gitmap clone text --verbose

**Output:**

    [verbose] Log file: gitmap-debug-2025-03-10T14-30.log
    Cloning from .gitmap/output/gitmap.txt...
    [1/5] Cloning https://github.com/user/my-api.git... done
    ...
    Clone complete: 5 succeeded, 0 failed

### Example 6: Force every URL into SSH-shorthand (`--ssh`)

    gitmap clone https://github.com/alimtvnetwork/wp-onboarding.git --ssh

**Output:**

      ↪ --ssh rewrite: https://github.com/alimtvnetwork/wp-onboarding.git → git@github.com:alimtvnetwork/wp-onboarding.git
    Cloning wp-onboarding into wp-onboarding...
    Cloned wp-onboarding successfully.

## See Also

- [scan](scan.md) — Scan directories to generate output files
- [pull](pull.md) — Pull individual or grouped repos
- [desktop-sync](desktop-sync.md) — Sync repos to GitHub Desktop
- [clone-next](clone-next.md) — Clone next version of a repo

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter clone
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
