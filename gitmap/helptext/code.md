# gitmap code

Register a path with the VS Code Project Manager extension and open it in
VS Code. Supports a multi-root **`paths`** subcommand (v3.39.0+).

## Usage

    gitmap code [alias] [path] [extraPath...]
    gitmap code paths add  <alias> <path>
    gitmap code paths rm   <alias> <path>
    gitmap code paths list <alias>

| Form                                              | Behavior                                                                          |
|---------------------------------------------------|-----------------------------------------------------------------------------------|
| `gitmap code`                                     | Use git repo root (or CWD). Alias = folder basename.                              |
| `gitmap code myalias`                             | Same path resolution; alias overridden to `myalias`.                              |
| `gitmap code myalias D:\anywhere`                 | Use any path. Alias = `myalias`.                                                  |
| `gitmap code myalias D:\root D:\extra1 D:\extra2` | Register root + variadic extras (additive — never clobbers user-added paths).     |
| `gitmap code paths add webapp D:\frontend`        | Attach an extra folder to the existing `webapp` entry.                            |
| `gitmap code paths rm  webapp D:\frontend`        | Detach an extra folder. Forces a clean overwrite (does NOT re-union it back in).  |
| `gitmap code paths list webapp`                   | Print rootPath + every attached extra path for the alias.                         |

## What it does

1. Resolves the absolute `rootPath`.
2. Upserts the `VSCodeProject` row keyed by `rootPath` (case-insensitive).
3. UNIONs any `extraPath...` into the DB-side `Paths` column (JSON-encoded).
4. **Detects auto-tags** (v3.40.0+): inspects the rootPath for top-level
   markers (`.git`, `package.json`, `go.mod`, `pyproject.toml`,
   `Cargo.toml`, `Dockerfile`, ...) and UNIONs them into the entry's
   `tags` array. User-edited tags are never removed.
5. Atomically syncs the entry into the alefragnani.project-manager
   `projects.json` file. Foreign entries plus user-edited
   `enabled` / `profile` are preserved verbatim. Paths and tags added in
   the VS Code UI are also preserved — gitmap only ever adds (or, via
   `paths rm`, explicitly removes a single path entry).
6. Launches VS Code on the resolved root path. The `paths` subcommand
   skips this step.

## Multi-root semantics

- Adding via `paths add` or variadic `gitmap code <alias> <root> <extra...>`
  is **additive**: existing extras (DB-managed or UI-added) are kept.
- `paths rm` is the only way to drop a path — it overwrites `paths` so
  the deletion sticks across re-syncs.
- `gitmap as <newalias>` only renames `name`. Extras, tags, and other
  user fields are left exactly as set.

## Auto-tag detection (v3.40.0+)

| Marker file / dir                                   | Tag      |
|-----------------------------------------------------|----------|
| `.git`                                              | `git`    |
| `package.json`                                      | `node`   |
| `go.mod`                                            | `go`     |
| `pyproject.toml` / `requirements.txt`               | `python` |
| `Cargo.toml`                                        | `rust`   |
| `Dockerfile` / `compose.yaml` / `docker-compose.yml`| `docker` |

Detection is shallow (top-level only) and read-only. To skip auto-tagging
entirely during a bulk scan, use `gitmap scan --no-auto-tags`.

## Examples

    gitmap code                                  # rootPath = repo root
    gitmap code backend                          # alias overridden
    gitmap code docs ~/Documents/spec            # any path
    gitmap code mono ~/work/main \
                     ~/work/main/frontend \
                     ~/work/main/backend         # root + extras

    gitmap code paths add  mono ~/work/main/scripts
    gitmap code paths list mono
    gitmap code paths rm   mono ~/work/main/scripts

## See also

- `gitmap as` — register a short alias (mirrors to projects.json).
- `gitmap scan` — bulk-syncs every discovered repo into projects.json.
- Spec: `spec/01-vscode-project-manager-sync/README.md`

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter code
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
