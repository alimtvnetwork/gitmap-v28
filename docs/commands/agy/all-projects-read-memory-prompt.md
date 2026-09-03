# All Projects Read Memory Prompt (`gitmap agy all-projects-read-memory-prompt`)

Broadcasts the standardized **Read Memory protocol instruction** to all active Antigravity projects simultaneously, with flexible exception filtering by project ID, project name, slug, or short prefix starts with.

## Synopses & Aliases

```bash
gitmap agy all-projects-read-memory-prompt [flags]
gitmap agy aprmp [flags]
gitmap agy read-memory-all [flags]
```

## Flags

| Flag | Shorthand | Type | Default | Description |
|---|---|---|---|---|
| `--except` | `-e` | string | `""` | Exclude projects matching projectid, name, slug, or short prefix starts with |
| `--prompt` | `-p` | string | Standard Read Memory protocol prompt | Custom prompt text to broadcast |
| `--dry-run` | `-d` | bool | `false` | Preview targeted and excluded projects without dispatching |
| `--yes` | `-y` | bool | `false` | Send prompt without interactive confirmation prompt |

## Default Prompt Text

> `"Execute enhanced Read Memory protocol. Defensively load memory, specs, constraints, and pending plans before taking action."`

## Examples

### 1. Preview targeting across all active projects (dry run)
```bash
gitmap agy all-projects-read-memory-prompt --dry-run
```
**Alias:**
```bash
gitmap agy aprmp -d
```

### 2. Exclude projects using short prefix starts with, name, or slug
```bash
gitmap agy aprmp --except "wp-, prompts, 46d0" --dry-run
```
*Notes:*
- `wp-`: Excludes all projects whose name, slug, or ID starts with `wp-` (e.g. `wp-exam-v2`, `wp-link-manager-v5`).
- `prompts`: Excludes projects starting with or named `prompts` (e.g. `prompts`, `prompts-connect-v3`).
- `46d0`: Excludes projects whose ID starts with `46d0` (e.g. `gitmap-v28` with ID `46d05021...`).

### 3. Broadcast to all active projects with confirmation prompt
```bash
gitmap agy all-projects-read-memory-prompt
```
**Alias:**
```bash
gitmap agy aprmp
```

### 4. Non-interactive immediate broadcast with custom prompt
```bash
gitmap agy aprmp -p "Load pending plans in .lovable/plans/pending/ and execute self-loop." -y
```

### 5. Exclude using an external CSV or text list
```bash
gitmap agy aprmp --except "./exceptions.csv" -y
```
