# gitmap update

Update GitMap CLI locally or across an entire cluster fleet in parallel via encrypted SSH.

## Synopsis

    gitmap update [target|version] [flags]
    gitmap ua [flags]
    gitmap update-all [flags]
    gitmap update all [--target <node>] [--except <nodes>] [flags]
    gitmap update ls [flags]
    gitmap update agm [node] [flags]

## Aliases

- **`gitmap ua`**: High-speed parallel update of GitMap CLI across all registered cluster nodes.
- **`gitmap update-all` / `gitmap updateall`**: Hyphenated alias for fleet update.
- **`gitmap update all` / `gitmap update -all`**: Argument form for fleet update.

## Description

`gitmap update` manages the lifecycle of the GitMap CLI binary both locally and across multi-node SSH cluster fleets.

### 1. Local Self-Update
When run without cluster flags, `gitmap update` resolves the latest GitHub release, downloads the platform-specific installer, verifies checksums, performs atomic binary replacement, and cleans temporary files.

### 2. Fleet & Cluster Update (`ua` / `update all`)
When executed with `ua`, `update all`, or `--all`, GitMap broadcasts update commands to all active cluster nodes in parallel. Remote execution utilizes structured JSON communication over SSH channels, completely suppressing raw installer stdout, 404 probing warnings, and help text from polluting the terminal.

The resulting output is rendered in a high-quality, professional terminal table detailing each node's alias, IP address, status (`SUCCESS`, `FAILED`, `OFFLINE`), roundtrip duration, and installed version.

## Subcommands

| Subcommand | Description |
| :--- | :--- |
| `update` | Self-update local GitMap CLI to latest release. |
| `update all` (or `ua`) | Broadcast parallel updates across all cluster nodes. |
| `update <node>` | Update GitMap on a single remote node (by alias, IP, or sequence). |
| `update ls` | Query and list installed packages and versions across all fleet nodes in parallel. |
| `update agm` | Update Antigravity-Manager (`agm`) across cluster nodes. |
| `update cleanup` | Purge obsolete handoff binaries and temporary installer artifacts. |

## Flags

| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--all` | `-a`, `-all` | false | Target all active cluster nodes in parallel (default for `ua`). |
| `--target <node>` | `-t` | "" | Target specific node by alias (`w1`), IP (`192.168.1.3`), sequence (`1`), or ID. |
| `--except <nodes>` | `-e` | "" | Exclude specific nodes from fleet update (comma-separated list). |
| `--force` | `-f` | false | Reinstall binary even if already at latest release. |
| `--dry-run` | *(none)* | false | Simulate update plan and print node manifest without executing. |
| `--json` | `-j` | false | Emit structured JSON telemetry summary. |
| `--help` | `-h` | false | Display two-column interactive help menu. |

## Examples

### Example 1: Local Self-Update
```bash
gitmap update
```

### Example 2: Update All Cluster Nodes in Parallel
```bash
gitmap ua
```

**Output:**
```text
[FLEET UPDATE] Updating 'all' across 3 cluster node(s) in parallel...

  ✓ [w1|192.168.1.3] Version: gitmap v6.361.0 (1420ms)
  ✓ [w2|192.168.1.4] Version: gitmap v6.361.0 (1380ms)
  ● [w3|192.168.1.5] OFFLINE: machine is off or unreachable (15ms)

================================================================================
 SSH Fleet Update Summary [all]: Total: 3 | Succeeded: 2 | Failed: 0 | Offline: 1 | Excluded: 0
--------------------------------------------------------------------------------
 ALIAS          IP               STATUS     DURATION   DETAILS
 w1             192.168.1.3      SUCCESS    1420ms     Version: gitmap v6.361.0
 w2             192.168.1.4      SUCCESS    1380ms     Version: gitmap v6.361.0
 w3             192.168.1.5      OFFLINE    15ms       machine is off or unreachable
================================================================================
```

### Example 3: Update Fleet Except Node 2
```bash
gitmap update all --except=w2
```

### Example 4: Update a Single Remote Worker Node
```bash
gitmap update w1
```

### Example 5: Query Cluster Software Inventory
```bash
gitmap update ls
```

## See Also

- `gitmap ssh ls` — list available cluster nodes and sequence numbers
- `gitmap deploy` — smart file and folder deployment engine
- `gitmap install` — local runtime and tool profile installer
