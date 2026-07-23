# gitmap vscode-workspace

Generate a single VS Code multi-root **`.code-workspace`** file from
every repo gitmap has tracked, so one click in VS Code (File → Open
Workspace from File…) opens them all in one window.

Source-of-truth is the same Repo table that drives the Project
Manager sync (`gitmap scan`), so the workspace file always mirrors
the latest scan output without a separate plumbing path.

## Alias

vsws

## Usage

    gitmap vscode-workspace [--out <path>] [--relative] [--tag <name>] [--root-subdir <subdir>]
    gitmap vsws [--out <path>] [--relative] [--tag <name>] [--root-subdir <subdir>]

## Flags

| Flag | Default | Purpose |
|------|---------|---------|
| `--out <path>` | `./gitmap.code-workspace` | Output `.code-workspace` file path. Parent dirs are created if missing. Written atomically (temp + rename) so VS Code never sees a half-written file. |
| `--relative` | off | Emit folder paths relative to the workspace file's directory (using forward slashes, VS Code's preferred form). Default is absolute paths so the file is portable across cwds. |
| `--tag <name>` | _(none)_ | Include only repos whose auto-detected tag set contains the given tag (e.g. `go`, `node`, `git`). Tag detection matches what the Project Manager sync writes. |
| `--root-subdir <subdir>` | _(none)_ | Add `<repo>/<subdir>` as the workspace folder instead of the repo root. Repos that don't contain that subdir are skipped (with a notice on stderr). Useful for monorepos where the interesting code lives under a fixed subpath like `src/`, `app/`, or `packages/web/`. |

## Output schema

    {
      "folders": [
        { "name": "my-repo", "path": "/abs/path/to/my-repo" },
        { "name": "other-repo", "path": "/abs/path/to/other-repo" }
      ],
      "settings": {}
    }

Folders are de-duplicated by cleaned path and sorted by name for
diff stability across runs.

## Examples

### Example 1: Default — every tracked repo, absolute paths

    gitmap vsws

**Output:**

      ✓ wrote /home/jane/work/gitmap.code-workspace with 12 folder(s)

Open the file in VS Code: `code gitmap.code-workspace`.

### Example 2: Portable workspace next to your repos

    cd ~/work
    gitmap vsws --out workspace.code-workspace --relative

Folder paths become `./repo-a`, `./repo-b`, … so the workspace can
be checked into git or shared via Dropbox without breaking on a
different machine.

### Example 3: Only Go repos

    gitmap vsws --tag go --out go-only.code-workspace

Useful when you maintain hundreds of repos but only want one window
per language.

### Example 4: Drill into a fixed subdirectory of each repo

    gitmap vsws --root-subdir src --out src-workspace.code-workspace

Each folder in the resulting workspace points at `<repo>/src` rather
than the repo root. Repos without a `src/` directory are skipped and
reported on stderr.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Workspace written, OR no tracked repos (printed a hint and exited cleanly). |
| 1 | DB open / list failure, or filesystem write failure. |
| 2 | Unknown flag (standard Go `flag` package behavior). |

## See Also

- [scan](scan.md) — Populates the Repo table that this command reads.
- [code](code.md) — Open a single repo in VS Code with the Project Manager registered.
- [vscode-pm-path](vscode-pm-path.md) — Print where `projects.json` lives.
## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter vscode-workspace
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
