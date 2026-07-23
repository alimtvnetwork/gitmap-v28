# gitmap list --update

List every scanned repo that has an upgrade available (local tag or HEAD is
behind the latest known remote release). This is the read-only companion to
`gitmap update apply` / `gitmap update all`.

## Alias

lu

Also available as `gitmap update list`, which produces byte-identical output.

## Usage

    gitmap list --update [flags]
    gitmap lu [flags]
    gitmap update list [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Emit a JSON array instead of a table |
| --source \<kind\> | (all) | Filter by source: `release`, `import`, or `git-tag` |
| --limit \<n\> | 0 (all) | Show only the top N upgradable repos |
| --stale-days \<n\> | 0 | Rescan repos whose last probe is older than N days before listing |
| --only-behind | false | Hide up-to-date rows (default when flag is set); has no effect on `--json` filtering |
| --path-only | false | Print just one absolute repo path per line, script-friendly |

## Prerequisites

- At least one prior `gitmap scan` populated the local database
- Network access for the background version probe (otherwise the list uses
  cached probe results only)

## Output Columns

| Column | Description |
|--------|-------------|
| repo | Repo slug (`owner/name`) or nickname |
| current | Locally-checked-out tag, or `-` when detached / untagged |
| latest | Latest remote release tag known to gitmap |
| behind | How many releases the local checkout trails by |
| source | Where `latest` came from: `release` (self-hosted release JSON), `import` (imported release row), or `git-tag` (raw `git ls-remote --tags`) |
| path | Absolute path on disk |

Rows are sorted by `behind` descending, then repo slug ascending. Exit code is
`0` when the list rendered, `2` when the scan cache is missing.

## Examples

### Human-readable table

    $ gitmap list --update
      Upgradable repositories (3 of 142 scanned):

      REPO                      CURRENT   LATEST    BEHIND   SOURCE     PATH
      alimtvnetwork/gitmap-v28  v6.75.0   v6.80.0        5   release    /home/a/src/gitmap
      acme/api                  v1.11.0   v1.14.0        3   git-tag    /home/a/work/api
      acme/web                  -         v2.0.0         -   import     /home/a/work/web
      ─────────────────────────────────────────────────────────────────────────────
      Next steps:
        • gitmap update apply <repo>   Upgrade one repo
        • gitmap update all            Upgrade every row above (prompts for -y)

### Script-friendly path list

    $ gitmap lu --path-only
    /home/a/src/gitmap
    /home/a/work/api
    /home/a/work/web

### JSON output

    $ gitmap list --update --json
    [
      {"repo":"alimtvnetwork/gitmap-v28","current":"v6.75.0","latest":"v6.80.0","behind":5,"source":"release","path":"/home/a/src/gitmap"},
      {"repo":"acme/api","current":"v1.11.0","latest":"v1.14.0","behind":3,"source":"git-tag","path":"/home/a/work/api"}
    ]

### Only release-tracked repos

    $ gitmap lu --source release --limit 5

## Behavior Notes

1. `list --update` reads the same `VersionProbe` rows that `find-next`
   consumes, so probe freshness (see `probe --depth`) directly controls
   accuracy.
2. When `--stale-days N` is set, gitmap re-probes stale rows in-process
   before printing; without the flag the list is cache-only.
3. A repo whose remote is unreachable is omitted from the upgradable table
   and counted in `stats` under `unknown`.

## See Also

- [scan](scan.md) — Populate the database `list --update` reads from
- [stats](stats.md) — Aggregate view including the Upgrades block
- [hd](hd.md) — One-screen dashboard combining pending + upgrades
- [update-apply](update-apply.md) — Upgrade a single repo
- [update-all](update-all.md) — Upgrade every upgradable repo
- [find-next](find-next.md) — Legacy view: repos with new versions available
- [list-versions](list-versions.md) — All release tags in one repo

## Scripting (JSON)

    gitmap help --json --filter list-update

Schema: `spec/08-json-schemas/list-update.schema.json` (v6.80.0+).
