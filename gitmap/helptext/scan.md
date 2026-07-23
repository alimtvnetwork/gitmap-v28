# gitmap scan

Scan a directory tree for Git repositories and record them in the local database.

## Alias

s

## Usage

    gitmap scan [dir] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --config \<path\> | ./data/config.json | Config file path |
| --mode ssh\|https | https | Clone URL style |
| --output csv\|json\|terminal | terminal | Output format |
| --output-path \<dir\> | ./.gitmap/output | Output directory |
| --manifest \<dir\> | (alias) | Alias for `--output-path`; mirrors `gitmap reclone --manifest` so the scan→reclone round-trip uses one vocabulary. Honored only when `--output-path` is empty |
| --github-desktop | false | Add repos to GitHub Desktop |
| --open | false | Open output folder after scan |
| --quiet | false | Suppress clone help section |
| --no-vscode-sync | false | Skip syncing into VS Code Project Manager projects.json |
| --no-auto-tags | false | Skip auto-derived tags (git/node/go/...) when syncing |
| --workers \<n\> | 0 (auto) | Worker-pool size for the parallel directory walker. `0` picks `min(NumCPU, 16)`; explicit values are clamped into `[1, 16]` to stay under the per-process file-descriptor budget |
| --concurrency \<n\> | — | Deprecated alias for `--workers`; emits a one-line stderr notice. When both are passed, `--workers` wins |
| --relative-root \<dir\> | (scan dir) | Pin the base directory used to compute every output `RelativePath`. Makes CSV/JSON/text/structure/clone-script artifacts byte-stable across cwds. Repos outside the root keep the scanner-computed path and emit a stderr warning |
| --max-depth \<n\> | 4 | Hard cap on directory levels descended below the scan root. Scan root = depth 0, its children = depth 1, etc. Default 4 keeps walks bounded on huge trees; pass `-1` for unlimited (legacy behavior). Repos found earlier still stop their own subtree |
| --default-branch \<name\> | main | Fallback branch name when HEAD/remote-tracking detection finds nothing. Written into `ScanRecord.Branch` for repos where every live detection step (HEAD, remote-tracking, configured default) returned an empty value. The same flag is also accepted by `gitmap clone` to rebuild clone instructions for rows whose recorded Branch is empty / detached / unknown |
| --probe-workers \<n\> | 3 | Worker-pool size for the background version probe. `0` disables it. The auto-trigger ceiling (50 repos) is bypassed when this flag is explicitly set |
| --probe-concurrency \<n\> | — | Deprecated alias for `--probe-workers`; emits a one-line stderr notice |
| --probe-depth \<n\> | 1 | `--depth N` passed to the `git clone` shallow-clone fallback inside the background probe. No effect on the `ls-remote` fast path |
| --no-probe | false | Skip the background probe entirely (offline / air-gapped runs) |
| --no-probe-wait | false | Return as soon as scan finishes; let probes keep draining until process exit |

## Prerequisites

- None (this is typically the first command you run)

## Live progress indicator

While the walker runs, gitmap prints a single live status line to stderr
showing how many directories have been walked and how many Git repos
have been found so far:

    ⟳ Scanning — 1284 dirs · 37 repos

The line refreshes about ten times a second using a carriage return, so
it never grows in your scrollback. When the walk finishes, the live line
is replaced by a one-line summary:

    ✓ Walked 4271 directories · found 58 repositories

The indicator is suppressed automatically when `--quiet` is passed or
when stderr is not a terminal (CI, redirected output). In those cases
only the final summary line is emitted.

## Examples

### Relative path targets

The `[dir]` argument accepts any relative path. It is resolved against
your current working directory using `filepath.Abs`, then validated to
exist and to be a directory. When the resolved target differs from what
you typed, gitmap prints a one-line `↳ Resolved` hint to stderr so the
target is unambiguous.

    gitmap scan .          # scan the current directory
    gitmap scan ..         # scan the parent directory
    gitmap scan ../..      # scan two folders up
    gitmap scan ../../x    # scan the "x" folder two levels up
    gitmap scan ~/work     # "~" expands to your home directory

**Output (for `gitmap scan ../..`):**

    ↳ Resolved "../.." → /home/alim/projects
      ▶ gitmap scan v3.71.0 — /home/alim/projects
    ...

If the resolved path does not exist, gitmap exits with:

    Error: scan target ../../nope does not exist (resolved to /home/alim/nope)

## Depth column and `--max-depth`

Every CSV/JSON row carries a `depth` column: the directory level (relative
to the scan root) where the `.git` folder was discovered. The scan root
itself is `0`, its immediate children are `1`, and so on. The `--max-depth`
flag is the **inclusive** cap on how deep the walker may descend before it
stops enqueuing further subdirectories.

How to read at-cap rows in the CSV:

- `depth < max-depth` — the walker reached the repo with budget to spare.
  Anything nested under that repo was *not* skipped by the depth cap (the
  walker stops descending into a found repo for a different reason: repos
  own their subtree).
- `depth == max-depth` — the repo was found exactly at the boundary. The
  walker did **not** descend into its children, and any nested
  repositories below it were silently skipped because the cap was hit.
  These rows are the ones to re-run with a larger `--max-depth` (or
  `-1` for unlimited) if you suspect monorepos-of-monorepos at the edge.
- `depth > max-depth` — never produced by the walker. If you see one,
  treat it as a bug.

Bumping `--max-depth` does not change `depth` values for repos already
discovered; it only widens the search so deeper repos can appear with
their own (larger) `depth`. Lowering the cap below an existing row's
`depth` will make that row disappear from the next scan.

Quick rescan recipe for at-cap rows:

    # find candidates in the previous CSV (depth column == cap, default 4)
    awk -F, 'NR>1 && $NF==4' .gitmap/output/gitmap.csv

    # rescan with a wider cap
    gitmap scan ~/work --max-depth 8

    # …or remove the cap entirely
    gitmap scan ~/work --max-depth -1

### Example 1: Scan a directory

    gitmap scan D:\wp-work

**Output:**

    Scanning D:\wp-work...
    [1/42] github/user/my-api
    [2/42] github/user/web-app
    [3/42] github/org/billing-svc
    ...
    Found 42 repositories
    ✓ Output written to ./.gitmap/output/
    ✓ Database updated (42 repos)

### Example 2: Scan with JSON output and SSH URLs

    gitmap scan ~/work --output json --mode ssh

**Output:**

    Scanning ~/work...
    Found 18 repositories
    ✓ .gitmap/output/gitmap.json written
    ✓ .gitmap/output/gitmap.csv written
    ✓ Clone URLs use SSH format (git@github.com:...)

### Example 3: Scan and register with GitHub Desktop

    gitmap scan D:\repos --github-desktop

**Output:**

    Scanning D:\repos...
    Found 12 repositories
    ✓ Output written to ./.gitmap/output/
    Registering with GitHub Desktop...
    [1/12] my-api... added
    [2/12] web-app... already registered
    ✓ 12 repos synced to GitHub Desktop (10 new, 2 existing)

### Example 4: Scan current directory quietly

    gitmap s . --quiet --output csv

**Output:**

    Scanning current directory...
    Found 7 repositories
    ✓ .gitmap/output/gitmap.csv written

## End-to-End Examples

These walkthroughs string `--config`, `--mode`, and `--output` together
into the workflows users actually run. Each one starts from a clean shell
and ends with a verifiable artifact on disk.

### E2E 1: Custom config + JSON output, then re-clone elsewhere

Use a project-local config to exclude vendored folders, emit JSON, then
hand the result to `gitmap clone` to mirror the same hierarchy on a
different machine.

    # 1. Author a config that skips heavy vendored trees.
    cat > ./gitmap.config.json <<'JSON'
    {
      "excludeDirs": ["node_modules", "vendor", ".next", "dist"],
      "defaultMode": "https",
      "outputDir": ".gitmap/output"
    }
    JSON

    # 2. Scan with the custom config and JSON output.
    gitmap scan ~/work \
      --config ./gitmap.config.json \
      --mode https \
      --output json \
      --output-path ./.gitmap/output

    # 3. Re-clone the same tree on another host (preserves folder hierarchy).
    gitmap clone ./.gitmap/output/gitmap.json --target-dir ~/mirror

**Expected output (step 2):**

    Scanning ~/work...
    ⟳ Scanning — 4271 dirs · 58 repos
    ✓ Walked 4271 directories · found 58 repositories
    ✓ .gitmap/output/gitmap.json written
    ✓ .gitmap/output/gitmap.csv written
    ✓ Database updated (58 repos)

### E2E 2: SSH mode for CI, CSV for spreadsheets

Same scan, two consumers: a CI job that needs SSH clone URLs and a PM
who wants the inventory in Excel.

    # CI-friendly: SSH URLs, JSON for tooling.
    gitmap scan /var/repos --config /etc/gitmap/ci.json --mode ssh --output json

    # Human-friendly: CSV alongside, with a custom filename.
    gitmap scan /var/repos --config /etc/gitmap/ci.json --mode ssh \
      --output csv --output-path ./reports

**What you get:**

    /var/repos/.gitmap/output/gitmap.json    ← used by CI clone job
    ./reports/gitmap.csv                     ← opens directly in Excel/Sheets

Both files describe the **same repo set** (the `--mode` flag only changes
which clone URL column is populated as the primary), so you can swap
output formats without re-walking the tree.

### E2E 3: Terminal preview before committing to a config

When you don't yet know which directories to exclude, preview with the
default `--output terminal` first, then promote a config file once the
tree looks right.

    # Step A — preview only, no files written.
    gitmap scan D:\code --output terminal

    # Step B — you noticed `node_modules` noise. Author a config:
    echo '{"excludeDirs":["node_modules",".turbo"]}' > gitmap.config.json

    # Step C — re-scan with the new config and persist JSON + CSV.
    gitmap scan D:\code --config ./gitmap.config.json --output json

The terminal preview in Step A and the JSON in Step C share the same
record schema, so anything visible in the preview also lands in the
written artifacts.

### E2E 4: Override the config's mode at the CLI

Flags beat the config file (see `config.MergeWithFlags`). Useful when one
scan needs SSH but the team config defaults to HTTPS.

    # team config says "mode: https" — override for this run only.
    gitmap scan ~/work --config ./team.gitmap.json --mode ssh --output json

The written `gitmap.json` will carry SSH URLs even though `team.gitmap.json`
says HTTPS. The config file is **not** mutated.

## See Also

- [rescan](rescan.md) — Re-scan using cached parameters
- [clone](clone.md) — Clone repos from scan output
- [status](status.md) — View repo statuses after scanning
- [desktop-sync](desktop-sync.md) — Sync scanned repos to GitHub Desktop
- [export](export.md) — Export scanned data
- [list-update](list-update.md) — Which scanned repos have upgrades available (post-scan hint links here)
- [hd](hd.md) — One-screen dashboard of scan + upgrades + pending

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter scan
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
