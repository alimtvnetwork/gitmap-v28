# Macro

Record, replay, and manage automated command sequences with environment variable expansion, interactive session control, and structured JSON/YAML reporting.

## Aliases

m, macros

## Usage

    gitmap macro record <name>
    gitmap macro run <name> [--json] [--yaml] [--file <path>] [--dry-run] [--verbose]
    gitmap macro list [--json] [--yaml] [--file <path>]
    gitmap macro show <name> [--json] [--yaml] [--file <path>]
    gitmap macro rm <name>

## Subcommands

| Subcommand | Description |
|---|---|
| record <name> | Start interactive recording session for shell commands |
| run <name> | Replay a recorded macro sequence with optional JSON/YAML export |
| list | List all saved macros |
| show <name> | View steps in a recorded macro |
| rm <name> | Delete a saved macro |

## Flags

| Flag | Description |
|---|---|
| `--json` | Output execution report in formatted JSON |
| `--yaml`, `-y` | Output execution report in formatted YAML |
| `--file <path>`, `-o <path>` | Save execution report to file (also prints confirmation banner) |
| `--dry-run` | Simulate macro execution without invoking commands |
| `--verbose`, `-v` | Show live command stdout and stderr |

## In-Session Recorder Commands

While recording inside `gitmap macro record <name>`:

| Command | Description |
|---|---|
| stop / exit / quit | Save recorded macro and exit |
| cancel / abort | Discard session without saving |
| undo | Remove the last recorded command |
| undo-steps <N> [-y] | Remove last N recorded commands |
| redo | Restore previously undone command |
| redo-steps <N> | Restore last N undone commands |
| list / steps | Show currently recorded steps in session |
| help / ? | Show in-session commands |

## Environment & Path Expansion

Paths and environment variables are automatically expanded during execution and recording:
- Windows variables: `%TEMP%`, `%USERPROFILE%`, `%APPDATA%`, `%LOCALAPPDATA%`
- Unix variables: `$HOME`, `$VAR`, `${VAR}`
- Tilde expansion: `~/projects`, `~`

## Examples

```bash
# Record a new macro session with environment variable expansion
gitmap macro record deploy-temp

# Inside session:
# [rec:deploy-temp 1]> cd %temp%
# [rec:deploy-temp 2]> ls
# [rec:deploy-temp 3]> undo-steps 1
# [rec:deploy-temp 3]> redo
# [rec:deploy-temp 3]> stop

# Replay the recorded macro
gitmap macro run deploy-temp

# Replay with structured JSON output
gitmap macro run deploy-temp --json

# Replay with YAML output and save to file
gitmap macro run deploy-temp --yaml --file "reports/deploy.yaml"

# Inspect macro steps as JSON
gitmap macro show deploy-temp --json

# List all macros as YAML
gitmap macro list --yaml
```
