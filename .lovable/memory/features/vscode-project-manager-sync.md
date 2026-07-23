---
name: VS Code Project Manager Sync
description: gitmap scan auto-syncs and `gitmap code` registers + opens repos in alefragnani.project-manager projects.json (multi-root paths v3.39.0, auto-tags v3.40.0)
type: feature
---

# VS Code Project Manager Sync (v3.40.0)

Spec: `/spec/01-vscode-project-manager-sync/README.md`
Sample fixture: `/spec/01-vscode-project-manager-sync/sample-projects.json`

## Source of truth

- gitmap SQLite DB (table `VSCodeProject`, schema v20) is the source of truth
  for `name`, `rootPath`, and `paths`.
- `tags` are NOT stored in SQLite — derived per-sync from rootPath
  filesystem markers (v3.40.0+).
- `projects.json` (the alefragnani.project-manager extension file) is a
  synced export.

## projects.json schema (locked)

`name` (string), `rootPath` (string, absolute, native separators),
`paths` ([]), `tags` ([]), `enabled` (bool), `profile` (string).

`enabled` and `profile` are preserved on upsert. `paths` and `tags` are
UNIONed with their on-disk values — gitmap only ever ADDS, never silently
removes user-added entries.

## File path resolution (DO NOT hardcode full path)

Resolve VS Code user-data root first, then append:
`User/globalStorage/alefragnani.project-manager/projects.json`.

| OS      | Root                                                              |
|---------|-------------------------------------------------------------------|
| Windows | `%APPDATA%\Code` → fallback `%USERPROFILE%\AppData\Roaming\Code`  |
| macOS   | `$HOME/Library/Application Support/Code`                          |
| Linux   | `$XDG_CONFIG_HOME/Code` → fallback `$HOME/.config/Code`           |

## Match key

`rootPath`. Case-insensitive on Windows, case-sensitive on Unix.
SQLite unique index uses `COLLATE NOCASE`.

## Auto-tag detection (v3.40.0)

| Marker (top-level only)                              | Tag      |
|------------------------------------------------------|----------|
| `.git`                                               | `git`    |
| `package.json`                                       | `node`   |
| `go.mod`                                             | `go`     |
| `pyproject.toml` / `requirements.txt`                | `python` |
| `Cargo.toml`                                         | `rust`   |
| `Dockerfile` / `compose.yaml` / `docker-compose.yml` | `docker` |

- Shallow + read-only + deterministic order (`constants.AutoTagOrder`).
- UNION with existing on-disk `tags`. Never deletes user-added tags.
- Opt-out per scan via `gitmap scan --no-auto-tags`.
- Marker registry: `gitmap/constants/constants_vscode_pm_autotags.go`.
- Detector: `gitmap/vscodepm/autotags.go`.

## Behavior rules

- `gitmap scan` syncs by default; `--no-vscode-sync` disables;
  `--no-auto-tags` keeps the sync but skips tag detection.
- `gitmap scan` **never** opens VS Code.
- `gitmap code [alias] [path] [extras...]` upserts + opens VS Code via
  `code "<rootPath>"`.
  - Variadic extras add multi-root paths (v3.39.0+).
- `gitmap code paths add|rm|list <alias> [path]` manages multi-root.
  `rm` calls `vscodepm.OverwritePaths` (bypasses path UNION) so removals
  actually stick.
- `gitmap as <newalias>` mirrors the DB name change to `projects.json`
  for the matching `rootPath`. Paths/tags/enabled/profile preserved.
- Foreign entries in `projects.json` are always preserved.
- Atomic writes only: write to `.tmp` then `os.Rename`.

## Command name

`gitmap` (single word). The string `git map` (with a space) must
not appear anywhere in code, help, or logs.

## Out of scope

Reverse sync, profile assignment, custom user-defined tag rules
(only the built-in marker set ships).
