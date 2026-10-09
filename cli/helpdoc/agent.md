# AI Agent Task Orchestrator & Telemetry

Manage AI agent tasks, subtasks, split-database telemetry, forensic inspection, and web UI.

## Alias

ai-agents, agents

## Usage

    gitmap agent <subcommand> [flags]

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| task | t | Manage parent agent tasks, initialization, and status rollups |
| subtask | st, sub | Manage atomic subtasks, worker claiming, and evidence recording |
| log | — | Record granular in-flight telemetry into agent split database |
| crashed | crash | Detect crashed or abandoned agents and display autopsy report |
| diagnose | diag, autopsy | Perform deep forensic inspection of agent tasks and telemetry |
| clear | — | Clear or archive completed agent task runs |
| reset | — | Reset agent databases and telemetry |
| temp-clear | clean-temp, purge-temp | Delete ephemeral .ai-memory/temp-agents/ directory |
| ui | — | Launch interactive web UI visualizer for agent tasks |

## Prerequisites

- Active workspace with `.ai-memory/` directory
- SQLite split-db support enabled

## Examples

### View agent status and help

```bash
gitmap agent --help
```

### Inspect crashed agent runs

```bash
gitmap agent crashed
```

### Perform forensic diagnosis

```bash
gitmap agent diagnose
```

### Clean temporary agent directory

```bash
gitmap agent temp-clear
```
