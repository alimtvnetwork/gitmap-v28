# gitmap swap-binary

Safely terminate locking processes and replace the active installed GitMap executable with a freshly compiled binary.

## Usage

```
gitmap swap-binary [--from <path>] [--to <path>] [--pid <pid>] [--kill] [--reopen|--restart] [--force] [--script-only] [--json]
```

## Aliases

`swap-bin`, `replace-binary`, `replace-bin`, `bin-swap`

## Flags

| Flag | Description |
|---|---|
| `--from <path>` | Path to the source newly-built binary (defaults to candidate scan `./cli/gitmap.exe`, `./bin/gitmap.exe`, `./gitmap.exe`). |
| `--to <path>`, `-t` | Target binary destination to replace (defaults to installed deployed binary or current running executable). |
| `--pid <pid>` | Target process PID to terminate before replacing binary file handles. |
| `--kill`, `-k` | Kill target processes locking the binary prior to file replacement. |
| `--reopen`, `--restart` | Reopen/restart the executable after swapping to verify version and execution. |
| `--force`, `-f` | Force immediate process termination and overwrite destination file. |
| `--script-only` | Output the generated standalone handoff script without executing it. |
| `--json` | Emit machine-readable JSON envelope containing swap execution results. |

## Exit Codes

| Code | Meaning |
|---|---|
| 0 | Binary swap succeeded, or script-only generated successfully. |
| 1 | Validation failure, path resolution error, or swap execution error. |

## Examples

### Windows

```powershell
# Auto-detect built binary in ./cli/gitmap.exe and swap into installed PATH
gitmap swap-binary

# Swap with explicit paths and verify version restart
gitmap swap-binary --from .\cli\gitmap.exe --to C:\tools\gitmap\gitmap.exe --restart

# Kill locking process by PID and overwrite
gitmap swap-binary --pid 14820 --kill --force

# Inspect generated PowerShell handoff script without executing
gitmap swap-binary --script-only
```

### Linux / macOS

```bash
# Auto-detect built binary in ./cli/gitmap and swap into /usr/local/bin/gitmap
gitmap swap-binary

# Script-only mode preview
gitmap swap-binary --from ./bin/gitmap --to /usr/local/bin/gitmap --script-only

# JSON output mode for programmatic orchestration
gitmap swap-binary --json
```
