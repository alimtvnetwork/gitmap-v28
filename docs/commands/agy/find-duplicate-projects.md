# Find Duplicate Projects (`gitmap agy find-duplicate-projects`)

Scans all Antigravity project definitions in `~/.gemini/config/projects/` to identify projects pointing to identical directories or duplicate project names, and outputs immediate copy-pasteable cure commands.

## Synopses & Aliases

```bash
gitmap agy find-duplicate-projects [flags]
gitmap agy fdp [flags]
gitmap agy find-duplicates [flags]
```

## Flags

| Flag | Shorthand | Type | Default | Description |
|---|---|---|---|---|
| `--except` | `-e` | string | `""` | Exclude project IDs, names, slugs, or paths from duplicate inspection |
| `--json` | — | bool | `false` | Output duplicate groupings as structured JSON |

## Examples

### 1. Find all duplicate projects (standard invocation)
```bash
gitmap agy find-duplicate-projects
```
**Alias:**
```bash
gitmap agy fdp
```

### 2. Exclude specific projects from duplicate scan
```bash
gitmap agy fdp --except "gitmap, coding-guidelines, 46d0"
```

### 3. Exclude via CSV or text file
```bash
gitmap agy fdp --except "./allowed-projects.csv"
```

### 4. Output structured JSON for automation
```bash
gitmap agy fdp --json
```

## Remediation Workflow

When duplicates are found, resolve them using:
- **Cure all duplicates:** `gitmap agy cure-duplicate-projects` (alias: `cdp`)
- **Preview before deletion:** `gitmap agy cure-duplicate-projects --dry-run`
- **Cure with exclusions:** `gitmap agy cure-duplicate-projects --except "<id, name, prefix>"`
