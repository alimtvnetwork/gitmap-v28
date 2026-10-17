# gitmap agent reset

Reset agent split databases and telemetry while preserving task run directories.

## Usage

```bash
gitmap agent reset [flags]
```

## Description

`gitmap agent reset` resets the agent SQLite split databases (`agent-task.db`, `ai_agents.db`) in the target directory (default: `.ai-memory/temp-agents/`). It marks subtasks back to `PENDING` and purges execution action logs without deleting the parent directory hierarchy.

To display this help text without performing a reset:
```bash
gitmap agent reset --help
```

## Flags

| Flag | Shorthand | Type | Description | Default |
|------|-----------|------|-------------|---------|
| `--dir` | `-d` | string | Target run directory | `.ai-memory/temp-agents/` |
| `--task-id` | `-t` | string | Target parent task ID or slug | `""` |
| `--yes`, `--force` | `-y`, `-f` | bool | Confirm reset without prompt | `false` |
| `--json` | `-j` | bool | Output response in structured JSON | `false` |
| `--help` | `-h` | bool | Display help for reset command | `false` |

## Examples

### View Reset Help
```bash
gitmap agent reset --help
```

### Reset Default Temp-Agent Database
```bash
gitmap agent reset -y
```

### Reset Specific Task Run
```bash
gitmap agent reset -t task-20261017-sample -y
```

### Reset Custom Workspace Run Directory
```bash
gitmap agent reset -d .ai-memory/temp-agents/my-run/ -y
```
