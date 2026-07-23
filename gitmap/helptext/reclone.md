# gitmap reclone

Re-run `git clone` against the JSON / CSV / text artifacts produced by
`gitmap scan`, honoring the recorded folder structure and a user-
selected SSH/HTTPS mode.

> **`reclone` vs `clone`** — different commands, different inputs.
>
> | Command | Input | Use when |
> |---|---|---|
> | `gitmap clone <url> [folder]` | A single repo URL | You want to clone (or re-clone) one repo from a URL. |
> | `gitmap reclone <file>` | A `gitmap scan` artifact (JSON / CSV / TXT) | You want to round-trip an entire previously-scanned tree at its recorded relative paths. |
>
> If you're generating clone *commands* without running them, that's
> still `gitmap scan` — `reclone` is the side that consumes those
> artifacts and actually re-creates the tree.

## Synopsis

```
gitmap reclone                                            # auto-pickup .gitmap/output/gitmap.json (then .csv)
gitmap reclone  <file>                                    # dry-run (default)
gitmap reclone  <file> --execute                          # actually clone
gitmap reclone  --manifest <path>                         # explicit manifest (JSON or CSV)
gitmap reclone  --manifest <path> --execute               # explicit + execute
gitmap reclone  --scan-root <dir> --execute               # auto-pickup from <dir>/.gitmap/output/
gitmap reclone  --execute                                 # auto-pickup + execute
gitmap reclone  <file> --mode ssh --execute               # use SSH URLs
gitmap rec      <file> --execute                          # short alias
gitmap clone-now <file> --execute                         # legacy alias (kept forever)
gitmap cnow     <file> --execute                          # legacy short alias
gitmap relclone <file> --execute                          # legacy alias
gitmap rc       <file> --execute                          # legacy short alias

# Single-repo mode (overlay) — wipe + re-clone ONE git repo from its origin
gitmap reclone                                            # inside a git repo: prompts, then deletes + re-clones from remote.origin.url
gitmap reclone -y                                         # skip the y/N prompt
gitmap rc      ./some/repo                                # target a sibling repo by path
gitmap rc      ./some/repo -y                             # skip prompt + target sibling
```

## Single-repo mode (overlay)

When invoked WITHOUT a manifest — i.e. zero positionals + the cwd
contains `.git/`, OR exactly one positional that is a directory
containing `.git/` — `reclone` switches to single-repo mode:

1. Reads `remote.origin.url` from the target.
2. Confirms destructive removal interactively (`y/N`). Pass `-y` /
   `--yes` to skip. Non-TTY callers without `-y` are refused
   (exit `1`) — protects `yes | gitmap rc` in CI from nuking the
   wrong tree.
3. Releases the Windows cwd handle (so re-cloning from inside the
   target works on Windows), `os.RemoveAll`s the folder, then
   `git clone <origin>` into the same parent under the same name.
4. Emits a shell handoff so the wrapper cd's back into the fresh
   checkout.

Manifest mode is unchanged: any other arg shape (a `<file>` path,
`--manifest`, `--scan-root`, etc.) falls through to the manifest
pipeline byte-for-byte.


## Source resolution

`reclone` picks the input file in this priority order:

1. `--manifest <path>`  — explicit, highest priority. JSON or CSV;
   format is auto-detected from the extension (override with `--format`).
2. Positional `<file>`  — legacy form, kept for back-compat.
3. Auto-pickup          — searches `<root>/.gitmap/output/gitmap.json`
   then `<root>/.gitmap/output/gitmap.csv`. `<root>` defaults to the
   current directory and can be redirected with `--scan-root <dir>`.

Passing **both** `--manifest` and a positional `<file>` is a usage error
(exit `2`) so the chosen artifact is unambiguous. `--scan-root` is only
consulted by the auto-pickup branch — it is silently ignored when an
explicit path is supplied.


## Auto-pickup

When `<file>` and `--manifest` are both omitted, `reclone` looks for
a scan artifact under:

1. `<scan-root>/.gitmap/output/gitmap.json`  (preferred — richest schema)
2. `<scan-root>/.gitmap/output/gitmap.csv`   (fallback)

`<scan-root>` is the current directory by default, or the value of
`--scan-root <dir>` when supplied. The first match is used and its
path is echoed to stderr so the run stays reproducible. If neither
file exists, `reclone` exits with code `2` and tells you to run
`gitmap scan` first (or pass `--manifest` / a positional path).
Auto-pickup never walks parent or sibling directories.

## Pre-execute summary

When `--execute` is passed, `reclone` prints a one-screen summary
to stderr **before** any `git clone` runs (and before the safety
prompt). It shows:

- the resolved source manifest, format, mode, on-exists policy, and cwd
- row totals: `N total (X new, Y already exist)` so you see the
  blast radius at a glance
- a sorted, indented tree of destination `RelativePath`s — capped
  at 40 lines with an "... and N more" footer for big round-trips

Pass `--no-summary` to suppress it (e.g. when a wrapper has already
printed a richer dry-run preview).

## Safety prompt (existing destinations)

Before any `git clone` runs under `--execute`, `reclone` checks
whether any planned `RelativePath` already exists under `--cwd`
(default: current directory). If at least one does:

- **Interactive TTY**: lists up to 10 existing destinations + total
  count and prompts `Proceed with git clone against these
  destinations? [y/N]:`. Only `y` proceeds; anything else aborts
  with exit `2` and no side effects.
- **Non-TTY** (CI, piped, redirected stdin): refuses with exit `2`
  and tells you to pass `--yes`. There is no silent fallthrough —
  you must opt in explicitly.
- **`--yes` passed**: skips the prompt entirely.

The per-row `--on-exists` policy (`skip` / `update` / `force`)
still controls what actually happens to each existing directory;
this gate is a single high-level confirmation that fires *before*
any side effect, so an accidental `--on-exists force` against a
populated tree is impossible without explicit confirmation.

## Arguments

| Argument | Required | Description |
|---|---|---|
| `<file>` | no | Path to a `.json`, `.csv`, or `.txt` file produced by `gitmap scan` (typically under `.gitmap/output/`). When omitted, auto-pickup is used (see above). |

## Flags

| Flag | Default | Description |
|---|---|---|
| `--manifest` | (none) | Explicit path to the scan artifact (`.json` or `.csv`). Equivalent to the positional `<file>` argument; cannot be combined with one. |
| `--scan-root` | current dir | Directory whose `.gitmap/output/` is probed during auto-pickup. Lets you `reclone` a tree scanned elsewhere without `cd`. Ignored when `--manifest` or a positional `<file>` is given. |
| `--execute` | off | Actually run `git clone`. Without this flag, only the dry-run plan is printed. |
| `--yes` | off | Skip the pre-flight confirmation when destination folders already exist. **Required for non-interactive / CI runs** — without a TTY and without `--yes`, `reclone --execute` exits `2` rather than block on stdin. The `--on-exists` policy still applies per row. |
| `--no-summary` | off | Suppress the pre-execute summary (totals + destination folder tree) printed before the safety prompt. Per-row results still print. |
| `--quiet` | off | Suppress per-row progress lines. The end-of-batch summary still prints. |
| `--mode` | `https` | URL mode to clone with: `https` or `ssh`. Falls back to the other mode if the preferred URL is missing on a row. |
| `--format` | auto | Force input format: `json`, `csv`, or `text`. Default auto-detects from the file extension. |
| `--cwd` | current dir | Working directory `git clone` runs in. Use to re-create the tree under a fresh root. |
| `--on-exists` | `skip` | Behavior when target already exists: `skip` (no-op when repo+branch match), `update` (fetch + checkout to align), `force` (remove target and re-clone — destructive). |
| `--max-concurrency` | auto | Worker count for parallel re-clones. `0` = `runtime.NumCPU()`, `1` = sequential. |
| `--no-vscode-sync` | off | Skip syncing every successfully re-cloned repo into VS Code Project Manager `projects.json`. Default is to sync once per batch (one entry per resolved destination). Has no effect during dry-run (no `--execute`). Honored identically by the `clone-now` / `cnow` aliases. |

Per-row destinations are canonicalized via `filepath.Clean` +
`filepath.EvalSymlinks` (Windows 8.3 short names + symlinks → long
form). Manifest `RelativePath` joins go through
`model.CleanRelativePath` so a row authored on POSIX produces an
identical `AbsolutePath` on Windows. `EvalSymlinks` failures soft-
fail to the cleaned absolute path; the re-clone never errors over a
resolution issue. Full rules: `gitmap clone --help` "Windows path
canonicalization & EvalSymlinks soft-fail".

## Aliases

`reclone` is the canonical name. The following spellings dispatch to
the exact same command and flag set, kept for backward compatibility:

- `rec`
- `clone-now`, `cnow`
- `relclone`, `rc`

## Examples

```
# Round-trip a previously scanned tree under a fresh root.
gitmap reclone .gitmap/output/repos.json --cwd ./mirror --execute

# Re-align an existing checkout with the recorded URL/branch.
gitmap reclone .gitmap/output/repos.csv --on-exists update --execute

# Inspect what would happen, with no side effects.
gitmap reclone .gitmap/output/repos.json
```

## Exit codes

- `0` — dry-run completed, OR every row was ok/skipped on `--execute`.
- `1` — file open / parse error, OR any row failed on `--execute`.
- `2` — bad CLI usage (missing `<file>`, invalid flag value), the safety prompt was declined / refused (existing destinations + non-TTY without `--yes`), OR the manifest failed semantic validation (see below).

## Manifest validation

Every run — dry-run and `--execute` alike — gates on a row-level
validation pass that fires immediately after the manifest is parsed.
Each row must satisfy:

- `RepoName` is non-empty.
- At least one of `HTTPSUrl` / `SSHUrl` is set, AND the URL selected
  by `--mode` (with fallback to the other) is a valid git URL —
  either `scheme://host/path` (https, ssh, git, file, ...) or the
  scp-like `user@host:path` shorthand.
- `RelativePath` is non-empty, **relative** (not absolute), and does
  not escape the working directory via `..`.

When any row fails, `reclone` prints a per-row report to stderr (one
line per issue, with the 1-based row index, repo name, dest, and a
short reason) and exits `2` without attempting any clones. There is
no opt-out flag: a malformed manifest is a usage error, and silently
proceeding would either crash inside `git` or write to a path you
did not intend. Fix the manifest and re-run.


## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter reclone
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
