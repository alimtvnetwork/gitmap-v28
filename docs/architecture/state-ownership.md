# State Ownership: `.gitmap/` vs `gitmap.json` vs SQLite/split-DBs

One-page reference for which system owns what persistent state.
Researched from the codebase (spec 243, Wave D).

## `.gitmap/` — per-directory working state

`.gitmap/` lives at the root of a scanned working directory (or the repo
being operated on). It is **working state, not the source of truth** —
deleting it loses caches and derived artifacts, never the repos themselves.

| Path | Owner | Contents |
|---|---|---|
| `.gitmap/output/gitmap.json` | `cmdscan` (writer); many readers | Scan-result inventory: every discovered repo's clone URL, branch, and metadata. Regenerable via `gitmap scan`. |
| `.gitmap/data/<section>/<slug>/sql.db` | feature package (e.g. `pipelinedb`); path via `store` | Split SQLite DBs (see below). |
| `.gitmap/archive/` | `stale` command | Repos archived by `gitmap stale --archive`. |
| `.gitmap/config.json` | `config` package (candidate path) | Optional local config override. |

## `gitmap.json` — derived scan artifact

- Canonical location: `.gitmap/output/gitmap.json`
  (`constants.DefaultOutputDir = ".gitmap/output"`).
- **Writer:** the scan pipeline (`cmdscan`). **Readers:** `status`,
  `clone`, `pull`, `exec`, `github-desktop`, SSH inventory, etc.
- It is a **read-mostly derived artifact**: any command can regenerate it
  by re-running `gitmap scan`. Treat it as a cache of the last scan, not
  as authoritative state. (Legacy bare `output/gitmap.json` is only a
  fallback for backward compatibility.)

## SQLite databases — authoritative structured state

### Main DB: `gitmap.db`

- Location: the binary's data dir (`store.BinaryDataDir()` →
  `<binary-dir>/data/`; `constants.DBDir = "data"`,
  `constants.DBFile = "gitmap.db"`).
- **Owner:** the `store` package — schema definition, migrations, and
  connection configuration (`store.ConfigureSQLiteConn`).
- Holds: repo registry, scan records, settings, release cache, pending
  tasks, and other cross-command tables.

### Split DBs: `<base>/data/<section>/<slug>/sql.db`

- Path resolution: `store.ResolveSplitDbPath` /
  `store.ResolveSplitDbDir` (`cli/store/split_db_path.go`).
- Sections: `pipeline`, `installation`, `automation`, `startup`,
  `sites`, `schedule`, `ai-analysis` (`store.Section*`).
- Base dir: `<repo>/.gitmap/data` when operating inside a `.gitmap`
  working root, else the binary's data dir.
- **Owner:** the feature package per section — e.g. `pipelinedb`
  owns `pipeline/<slug>/sql.db` (`OpenPipelineSplitDb`,
  `ResolvePipelineDbPath`). The `store` package owns only the
  **path convention**, not the contents.
- Rationale: per-repo/per-feature isolation so one repo's pipeline
  history never contends with another's; `storage clean` can prune a
  single split DB.

### Satellite DBs

| File | Owner | Purpose |
|---|---|---|
| `gitmap-errors.db` | `store` | Failed/unknown command log. |
| `installation.db` | `cmdinstall` | Installer state (legacy name; new installs use split DBs). |

## `config.json` — user configuration

- Discovery order (`config.CandidateConfigPaths`): `./data/config.json`,
  `.gitmap/config.json`, `%LOCALAPPDATA%/gitmap-cli/data/config.json`,
  `~/.gitmap/config.json`. First hit wins; missing file falls back to
  defaults.
- **Owner:** the `config` package. User-editable; never written by
  commands except explicit `config`/`import-config` flows.

## Rule of thumb

- **`.gitmap/output/gitmap.json`** — scan cache. Safe to delete; rescan rebuilds it.
- **`.gitmap/data/.../sql.db`** — feature state. Delete only the section you mean to reset.
- **`gitmap.db`** — the registry. Back it up before surgery (`db-reset` is destructive).
- **`config.json`** — user intent. Commands read it; they don't rewrite it behind your back.
