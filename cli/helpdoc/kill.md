# Process Termination Utility

Safe, cross-platform process termination by Process ID (PID), process name, or TCP listening port with built-in safety guards.

## Synopsis

```bash
gitmap kill <pid|name|:port> [flags]
gitmap terminate <pid|name|:port> [flags]
```

## Description

`gitmap kill` terminates processes cleanly across Windows, Linux, and macOS without requiring OS-specific commands (`taskkill /F /PID ...` vs `kill -9`). It includes protective safety guards preventing accidental termination of system-critical processes (PID 0, 1, 4), the current process (`os.Getpid()`), and developer IDEs/shells (`code`, `antigravity`, `pwsh`, `bash`) unless explicitly forced with `--force`.

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Force immediate termination (`SIGKILL` or `/F`) |
| `--tree` | `-t` | Terminate process and all descendant child processes (`/T` on Windows) |
| `-h`, `--help` | | Display this help menu |

## Aliases

- `gitmap terminate`

## Examples

### Kill Process by PID

```bash
gitmap kill 12345
```

### Kill Process Tree by PID

```bash
gitmap kill 12345 --tree
```

### Terminate Process Listening on TCP Port

```bash
gitmap kill :8080
gitmap kill :3000 -f
```

### Terminate Lingering Process by Name

```bash
gitmap kill node.exe
gitmap kill cargo -f
```

## See Also

- `gitmap hd cargo` — Cargo toolchain runner.
- `gitmap hd rust` — Unified Rust compiler frontend.
