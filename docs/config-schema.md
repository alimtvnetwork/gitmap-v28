# `data/config.json` — JSON Config Schema

This page documents every key gitmap reads from `data/config.json` and the
exact defaults baked into the binary.

> **Source of truth:** [`gitmap/model/record.go`](../gitmap/model/record.go)
> (`type Config struct` + `func DefaultConfig`).
> If anything here drifts from those Go definitions, the Go code wins —
> please open an issue or PR to update this doc.

---

## Lookup order & merge behavior

1. The binary starts with the built-in defaults from `model.DefaultConfig()`.
2. If `data/config.json` exists next to the binary (or at the path given via
   `--config <path>`), every key it defines **overrides** the matching default.
   Keys you omit fall through to the default — partial files are valid.
3. Three CLI flags can override the loaded values one more time at run-time:
   `--mode`, `--output`, and `--output-path`. Flags only override when set
   to a non-empty value (see `MergeWithFlags` in [`gitmap/config/config.go`](../gitmap/config/config.go)).

If `data/config.json` is missing, gitmap silently uses defaults — it is
**not** an error.

---

## Top-level keys

| JSON key            | Type        | Default              | Allowed values                  | Effect |
|---------------------|-------------|----------------------|---------------------------------|--------|
| `defaultMode`       | string      | `"https"`            | `"https"`, `"ssh"`              | URL flavor written to scan output. Overridable per-run with `--mode`. |
| `defaultOutput`     | string      | `"terminal"`         | `"terminal"`, `"csv"`, `"json"` | Output format used by `gitmap scan` when `--output` is not passed. |
| `outputDir`         | string      | `".gitmap/output"`   | any path (absolute or relative) | Directory where `gitmap.csv` / `gitmap.json` are written. Overridable with `--output-path`. |
| `excludeDirs`       | string[]    | `[]`                 | any directory base names        | Directory names skipped during the scan walk (e.g. `"node_modules"`, `"vendor"`). Matched as basename only — not paths. |
| `notes`             | string      | `""`                 | free-form                       | Operator-only notes. Not parsed; useful for documenting why a config exists. |
| `dashboardRefresh`  | integer     | `30`                 | seconds (`> 0`)                 | Auto-refresh interval used by `gitmap dashboard` (`db`). Values `<= 0` fall back to the default. |
| `release`           | object      | see below            | see below                       | Release pipeline configuration consumed by `gitmap release` / `release-branch`. |

---

## `release` sub-object

| JSON key      | Type           | Default | Effect |
|---------------|----------------|---------|--------|
| `targets`     | object[]       | `[]`    | Cross-compile matrix. Each entry is a `{goos, goarch}` pair (see below). |
| `checksums`   | boolean        | `false` | If `true`, gitmap writes a `checksums.txt` (SHA-256) alongside the artifacts. |
| `compress`    | boolean        | `false` | If `true`, each per-target binary is bundled into a `.tar.gz` (Linux/macOS) or `.zip` (Windows). |

### `release.targets[]` entries

| JSON key    | Type   | Example     |
|-------------|--------|-------------|
| `goos`      | string | `"linux"`, `"darwin"`, `"windows"` |
| `goarch`    | string | `"amd64"`, `"arm64"` |

---

## Minimal example

A no-frills personal config:

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": ".gitmap/output",
  "excludeDirs": [".cache", "node_modules", "vendor", ".venv"],
  "dashboardRefresh": 30,
  "notes": "personal laptop — HTTPS only, no PATs needed"
}
```

## CI-flavoured example (SSH + JSON + custom artifact dir)

```json
{
  "defaultMode": "ssh",
  "defaultOutput": "json",
  "outputDir": "./build/scan-results",
  "excludeDirs": ["node_modules", "vendor", "dist", ".terraform"],
  "dashboardRefresh": 60,
  "notes": "CI runner — SSH key mounted at /root/.ssh"
}
```

## Full example with `release` matrix

```json
{
  "defaultMode": "https",
  "defaultOutput": "json",
  "outputDir": ".gitmap/output",
  "excludeDirs": ["node_modules", "vendor"],
  "dashboardRefresh": 30,
  "notes": "release-pipeline machine",
  "release": {
    "checksums": true,
    "compress":  true,
    "targets": [
      { "goos": "linux",   "goarch": "amd64" },
      { "goos": "linux",   "goarch": "arm64" },
      { "goos": "darwin",  "goarch": "amd64" },
      { "goos": "darwin",  "goarch": "arm64" },
      { "goos": "windows", "goarch": "amd64" }
    ]
  }
}
```

---

## Key & value rules

* **Unknown keys are ignored** — gitmap will not error on extra fields, so a
  JSON config that's slightly ahead of the binary still loads. Stick to the
  documented keys to avoid silent surprise.
* **Keys are case-sensitive** and use `camelCase` exactly as listed above.
  `defaultMode` works; `DefaultMode` and `default_mode` do not.
* **Empty strings count as omitted** for `defaultMode`, `defaultOutput`,
  and `outputDir`: at the CLI-flag merge stage an empty value preserves the
  loaded default rather than blanking it out.
* **Relative paths** are resolved relative to the **current working
  directory** of the gitmap process — not relative to `data/config.json`.
  Use absolute paths in CI to avoid surprises.

---

## Where this file lives

* **Default location:** `./data/config.json` relative to the gitmap binary.
* **Override:** `gitmap scan --config /path/to/your.json` (recorded in the
  rescan cache, so `gitmap rescan` replays the same path automatically).
* The companion files in the same directory — `data/git-setup.json` and
  `data/seo-templates.json` — have **separate schemas** and are not covered
  here.
