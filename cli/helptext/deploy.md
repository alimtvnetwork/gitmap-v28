# gitmap deploy

Smart file and folder deployment engine across cluster nodes with bidirectional and unidirectional synchronization modes.

## Synopsis

    gitmap deploy <target> <source> <destination> [flags]
    gitmap deploy-right <target> <source> <destination> [flags]
    gitmap deploy-left <target> <source> <destination> [flags]
    gitmap ssh deploy <target> <source> <destination> [flags]

## Description

`gitmap deploy` provides high-speed, deterministic file and directory distribution across GitMap cluster nodes. It connects via encrypted SSH streaming transport, detects existing files, computes modification timestamps (`mtime`), and supports parallel file replacement for multi-file folder hierarchies.

The command seamlessly handles single files or entire directory trees, automatically creates missing remote directories, and adapts path separators across Windows (`\`) and Linux/macOS (`/`) remote hosts.

## Positional Arguments

| Argument | Description |
| :--- | :--- |
| `<target>` | Remote node identifier: alias (`w1`), IP (`192.168.1.3`), sequence index (`1`), or host ID (`host-1`). |
| `<source>` | Local file or folder path. Relative paths resolve against the local working directory (`os.Getwd()`). |
| `<destination>` | Remote file or folder destination path. Relative paths resolve against the target node's default work directory (`D:/work` on Windows, `~` on Linux). |

## Flags

| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--overwrite` | `-o` | false | Unconditionally replace existing remote files without prompting. |
| `--skip` | `-s` | false | Skip any file that already exists on the destination; write only missing files. |
| `--sync` | *(none)* | false | Bi-directional timestamp sync: transfer newer file by `mtime` (left to right or right to left). |
| `--sync-right` | *(none)* | false | Unidirectional push (left to right): copy only if local is newer than remote; never write local. |
| `--sync-left` | *(none)* | false | Unidirectional pull (right to left): pull newer remote files; keep local files intact on conflict. |
| `--parallel` | `-p` | 4 | Concurrency level (number of parallel worker goroutines, 1–16) for directory transfers. |
| `--json` | `-j` | false | Output machine-readable structured JSON telemetry and suppress terminal progress bars. |
| `--dry-run` | `-n` | false | Preview deployment actions and transfer manifest without writing any files. |
| `--help` | `-h` | false | Display help and usage information. |

## Subcommand Shortcuts

- **`gitmap deploy-right <target> <source> <dest>`**: Equivalent to `gitmap deploy` defaulting to `--sync-right`. Ideal for publishing local builds, release packages, or configs to remote workers.
- **`gitmap deploy-left <target> <source> <dest>`**: Equivalent to `gitmap deploy` defaulting to `--sync-left`. Ideal for pulling logs, telemetry dumps, or remote build artifacts without risking overwrites of local modifications.

## Conflict Resolution & Exit Codes

When neither `--overwrite`, `--skip`, nor `--sync*` is specified:
- In interactive terminal mode, the user is prompted: `[y/N/a(ll)/s(kip all)]`.
- In non-interactive or JSON mode, conflicts are detected and reported with exit code `3`.

| Exit Code | Meaning |
| :--- | :--- |
| `0` | Success (all transfers completed successfully). |
| `1` | Fatal error (network failure, invalid arguments, host unreachable). |
| `2` | Conflict skipped (files existed and were skipped due to `--skip`). |
| `3` | Partial or conflict detected (unresolved conflict requiring explicit sync flag). |

## Examples

### Deploy a single executable to worker-1 with overwrite
```bash
gitmap deploy w1 ./cli/gitmap.exe D:/work/gitmap/cli/gitmap.exe -o
```

### Sync local build directory to remote node using 8 parallel workers
```bash
gitmap deploy 192.168.1.5 ./dist dist --sync-right -p 8
```

### Pull remote log files to local workspace safely
```bash
gitmap deploy-left w2 logs/server.log ./logs/remote-server.log
```

### Automated pipeline deployment with JSON telemetry
```bash
gitmap deploy-right prod-1 ./release /opt/app --parallel=6 --json
```

## See Also

- `gitmap ssh ls` — list available cluster nodes and sequence numbers
- `gitmap ssh deploy bin` — deploy GitMap CLI binary across fleet
- `gitmap ssh deploy keys` — deploy authorized SSH keys across fleet
