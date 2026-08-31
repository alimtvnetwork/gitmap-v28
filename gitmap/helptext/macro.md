# Macro

Record, replay, and manage automated command sequences with environment variable expansion and interactive session control.

## Aliases

m, macros

## Usage

    gitmap macro record <name>
    gitmap macro run <name> [--dry-run] [--verbose]
    gitmap macro list
    gitmap macro show <name>
    gitmap macro rm <name>

## Subcommands

| Subcommand | Description |
|---|---|
| record <name> | Start interactive recording session for shell commands |
| run <name> | Replay a recorded macro sequence |
| list | List all saved macros |
| show <name> | View steps in a recorded macro |
| rm <name> | Delete a saved macro |

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

# Inspect macro steps
gitmap macro show deploy-temp

# List all macros
gitmap macro list
```
