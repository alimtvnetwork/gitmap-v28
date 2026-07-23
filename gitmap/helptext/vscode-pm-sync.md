# gitmap vscode-pm-sync

Re-tag every entry in the alefragnani.project-manager `projects.json`
file using the same auto-tag detector that runs after every `gitmap
clone`. Useful when you've added new repos, installed Docker / Cargo /
Python markers in existing folders, or just want to make sure the
`gitmap` brand tag is applied everywhere.

Alias: **`vpm`**.

## Usage

    gitmap vscode-pm-sync           # re-tag every entry, write changes
    gitmap vpm                      # short alias
    gitmap vpm --dry-run            # preview only, no write

## What it does

1. Resolves the active VS Code user-data root (`%APPDATA%\Code` on
   Windows, `~/Library/Application Support/Code` on macOS,
   `$XDG_CONFIG_HOME/Code` or `~/.config/Code` on Linux).
2. Reads every entry currently in
   `User/globalStorage/alefragnani.project-manager/projects.json`.
3. For each entry whose `rootPath` still exists on disk, re-runs the
   auto-tag detector (`.git`, `package.json`, `go.mod`, `pyproject.toml`,
   `Cargo.toml`, `Dockerfile`, ...) and **prepends the `gitmap` brand
   tag**.
4. UNIONs the freshly-detected tags into each entry's existing `tags`
   array. **User-added tags are never removed.**
5. Writes the result atomically (sibling `.tmp` then rename).

Foreign entries (anything not in your gitmap database) are walked
exactly the same way — this command operates on whatever already lives
in `projects.json`, regardless of who put it there.

## Example

    PS D:\projects> gitmap vpm
    → vscode-pm-sync: re-tagging projects.json entries at C:\Users\me\AppData\Roaming\Code\User\globalStorage\alefragnani.project-manager\projects.json
      ✓ projects.json synced: 0 added, 14 updated, 7 unchanged (21 total)
      ✓ scanned 21 entries, 0 skipped (rootPath missing on disk)

## Soft-fail policy

If VS Code or the project-manager extension is not installed (CI,
headless box, fresh dev VM) the command prints a single diagnostic
line to stderr and exits **0**. It will never break a CI pipeline
just because VS Code is absent.

## Customizing tags

The same `--vscode-tag` / `--vscode-tag-skip` / `--vscode-tag-marker`
flags accepted by `gitmap clone` work here too. They are stripped
from argv before flag parsing and persisted into env vars so the
detector picks them up on every entry.

    gitmap vpm --vscode-tag work --vscode-tag-skip docker
    gitmap vpm --vscode-tag-marker Gemfile=ruby

To opt the entire entry set out of the brand tag:

    gitmap vpm --vscode-tag-skip gitmap

## Flags

| Flag                | Purpose                                                                                       |
|---------------------|-----------------------------------------------------------------------------------------------|
| `--dry-run`         | Preview entry counts without writing any change.                                              |
| `--projects-json`   | Absolute path to a `projects.json` (overrides VS Code user-data discovery).                   |
| `--tag <name>`      | Replace detected tags with this name (repeatable; accepts comma-list).                        |
| `--mode <strategy>` | Tag merge strategy: `union` (default) | `replace` | `intersection`. See **Modes** below.    |

## Modes

| Mode           | Final tag set                                  | Use when                                                                       |
|----------------|------------------------------------------------|--------------------------------------------------------------------------------|
| `union`        | existing ∪ detected (dedup'd, default)         | You want gitmap to ADD tags but never remove user-added ones. The default.     |
| `replace`      | detected verbatim (brand survives via detector)| You want every entry's tag set to match exactly what the detector produces.    |
| `intersection` | (existing ∩ detected) ∪ {`gitmap`}             | You want to PRUNE stale tags but never lose the `gitmap` brand pin.            |

The `gitmap` brand tag is preserved under every mode:
- `union` — added by the detector, kept by the union.
- `replace` — added by the detector (which always pre-pends it).
- `intersection` — pinned explicitly even when the strict intersection is empty.

```
gitmap vpm                       # default — additive union
gitmap vpm --mode replace        # detector wins outright
gitmap vpm --mode intersection   # prune stale tags, keep brand
```

## Exit codes

| Code | Meaning                                                       |
|------|---------------------------------------------------------------|
| 0    | Success, OR soft-skip (VS Code / extension not installed).    |
| 1    | Hard failure (parse error on hand-edited projects.json, etc). |

## See also

- `gitmap code` — register a single repo with the extension and open it.
- `gitmap vscode-pm-path` / `vpath` — print the resolved projects.json path.
- `gitmap vscode-workspace` / `vsws` — emit a multi-root `.code-workspace` file.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter vscode-pm-sync
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap vscode-pm-sync
```
