# Specification 16 - Part 4: Help System Polish, DAC / Fix-Repo Redesign & Search Discoverability

## 1. DAC & Fix-Repo Section Redesign

### 1.1 Problem Statement

The current help text for `fix-repo` and advanced developer utilities prints a dense, unstructured block of flags with run-on explanations and unaligned examples directly in the main help view.

### 1.2 Structured Presentation Standard

All utility help sections (including `fix-repo`, `doctor`, `setup`, `gomod`) must follow a strict tabular structure:

```text
Fix-repo flags:
  -2 (default)        Rewrite last 2 prior versions (v(K-2)..v(K-1) -> vK)
  -3, -5              Widen rewrite window to last 3 or 5 versions
  --all               Rewrite every prior version v1..v(K-1) -> vK
  --dry-run           Preview changes without writing files
  --strict            Run `go test` on touched packages; exit 9 on failure
  --restrict <mode>   Narrow scope: no-version (nv) avoids bare-base replacement
  --verbose           Log each modified file and replacement count
```

## 2. Prominent Surfacing of Cluster & Server Subcommands

### 2.1 Cluster Command Hierarchy

The top-level `gitmap help` and interactive menus must prominently display the full cluster command hierarchy:

| Command | Shorthand | Role | Purpose |
|---|---|---|---|
| `gitmap serve` | `gitmap sv` | Orchestrator | Start cluster coordination daemon & token generator |
| `gitmap servers-clients` | `gitmap sc` | Broadcast | Run PowerShell/CMD/Git ops across ALL machines |
| `gitmap clients` | — | Broadcast | Run commands/reboots across CLIENT machines only |
| `gitmap cluster nodes` | `gitmap clst nodes` | Registry | List connected cluster machines, IPs, roles, status |
| `gitmap cluster history` | `gitmap clst hi` | Auditing | Inspect historical cluster run execution logs |
| `gitmap servers ls` | — | Query | List registered server nodes |
| `gitmap clients ls` | — | Query | List registered client nodes |
| `gitmap cluster export` | — | Backup | Export cluster node database to JSON/CSV |
| `gitmap cluster import` | — | Restore | Import cluster node database from JSON/CSV |

## 3. Multi-Keyword Filter & Discovery System

Searching `gitmap help --filter <q>` must return rich contextual groups for related workflows:

### 3.1 Searching `--filter cluster` or `--filter servers`

Returns the `Cluster & Network` group with direct usage lines and flags.

### 3.2 Searching `--filter join`

Returns the daemon join commands (`gitmap serve`, `gitmap-node-join`, join token flags, and cluster onboarding steps) along with runnable examples:
```text
  $ gitmap serve --port 9999
  $ gitmap-node-join --server 192.168.1.10:9999 --token <token>
  $ gitmap cluster nodes
```

## 4. Web Help Dashboard (`gitmap hd`) UI Expansion

The embedded documentation web application must provide a rich interactive reference:
- Dedicated **Cluster & Network** documentation section.
- Interactive command builder for multi-machine commands (`gitmap sc ps "..." --except 2`).
- Interactive macro execution tutorials.
- Search bar indexing commands, flags, subcommands, and markdown help guides.
