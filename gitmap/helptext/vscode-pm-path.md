# gitmap vscode-pm-path

Print the absolute path to the alefragnani.project-manager `projects.json`
file that gitmap will sync to.

## Alias

vpath

## Usage

    gitmap vscode-pm-path
    gitmap vpath

## What it does

1. Resolves the OS-specific VS Code user-data root:
   - **Windows**: `%APPDATA%\Code` → fallback `%USERPROFILE%\AppData\Roaming\Code`
   - **macOS**: `$HOME/Library/Application Support/Code`
   - **Linux**: `$XDG_CONFIG_HOME/Code` → fallback `$HOME/.config/Code`
2. Appends the extension-relative tail
   `User/globalStorage/alefragnani.project-manager/projects.json`.
3. Validates that both the user-data root and the extension storage dir
   exist on disk.
4. Prints the absolute path on success — or a soft-fail diagnostic on
   stderr (and exits non-zero) when something is missing.

## Exit codes

| Code | Meaning                                                        |
|------|----------------------------------------------------------------|
| 0    | Path resolved and validated. Printed to stdout.                |
| 1    | User-data root missing OR extension not installed OR I/O error.|

## Examples

### Example 1: Resolve and print the path (happy path)

    gitmap vscode-pm-path

**Output:**

    /home/jane/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json

### Example 2: VS Code not installed

    gitmap vpath

**Output (stderr):**

    vscode: user-data directory not found (is VS Code installed? checked APPDATA / HOME / XDG_CONFIG_HOME)

### Example 3: Use in a shell pipeline

    gitmap vpath | xargs jq '.[].name'

Streams every project name in `projects.json` straight to `jq`. Fails
fast (non-zero exit) when the file cannot be located.

## See Also

- [code](code.md) — Register a path with VS Code Project Manager and open it
- [scan](scan.md) — Bulk-syncs every discovered repo into `projects.json`
- [doctor](doctor.md) — Includes a `VS Code Project Manager` health check

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter vscode-pm-path
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
