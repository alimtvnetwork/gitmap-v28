# gitmap rerun

Replay active or recent prompts into Antigravity workspaces with optional IDE restart, image re-injection, and queued completion verification prefix checks.

## Synopsis

```bash
gitmap rerun [sequence|project|all|queue] [flags]
```

## Aliases

- `rr`
- `rra` (alias for `gitmap rerun all`)
- `rrq` (alias for `gitmap rerun queue`)
- `gitmap agy rerun`

## Description

`gitmap rerun` automates the extraction and re-injection of previous prompts into Antigravity workspaces. It discovers active or recently modified workspaces, formats prompt template wrappers, handles clipboard copying, and coordinates with running Antigravity IDE instances.

When targeting projects, `gitmap rerun` safely matches against sequence numbers (ordered by recent activity and pinned rank), project names, or repository workspace paths. Help keywords (`help`, `info`, `man`) are guarded to prevent false-positive substring matches against project names.

## Subcommands

- `gitmap rerun help`: Display comprehensive interactive help and usage guide.

## Arguments

- `<sequence>`: Numeric 1-based index (e.g. `1`, `2`, `3`) corresponding to active project order. `1` defaults to the current or most recently active workspace.
- `<project>`: Exact or flexible name/slug matching an Antigravity project (e.g. `gitmap`, `my-project`).
- `all`: Replay active prompts concurrently across all running Antigravity projects.
- `queue`: Re-inject queued prompts with completion verification prefixes.

## Flags & Options

| Flag | Default | Description |
|------|---------|-------------|
| `-p`, `--prompt <template>` | `is-done` | Prefix prompt template name or ID |
| `-d`, `--dry-run` | `false` | Preview constructed prompt without executing replay |
| `-r`, `--restart` | `false` | Restart Antigravity IDE before replaying prompt |
| `--no-restart` | `false` | Direct prompt injection without closing or restarting IDE |
| `-n`, `--new-conversation` | `true` | Create a new conversation session for prompt replay |
| `-m`, `--model <model>` | `""` | Target model for new conversation (`flash_lite`, `flash`, `pro`) |
| `-P`, `--project <id\|seq>` | `""` | Target project by sequence number, ID, or slug |
| `-c`, `--conversation <id>` | `""` | Target specific conversation ID |
| `-a`, `--all` | `false` | Rerun active prompts across all active projects |
| `-q`, `--queue` | `true` | Re-inject queued prompts with completion prefix |
| `--prefix <text>` | Default | Custom prefix added to queued prompts |
| `--no-clipboard` | `false` | Do not copy constructed prompt to system clipboard |

## Examples

### Example 1: Replay last prompt in current project

```bash
gitmap rerun
```

### Example 2: Replay prompt into project sequence #2

```bash
gitmap rerun 2
```

### Example 3: Preview prompt replay without sending

```bash
gitmap rr 1 --dry-run
```

### Example 4: Restart IDE and replay prompt into specific project

```bash
gitmap rerun gitmap -r
```

### Example 5: Replay prompts across all running projects with Pro model

```bash
gitmap rerun all --model pro
```

### Example 6: Show rerun help

```bash
gitmap rerun help
```
