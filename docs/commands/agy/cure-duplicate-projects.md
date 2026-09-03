# Cure Duplicate Projects (`gitmap agy cure-duplicate-projects`)

Deduplicates Antigravity projects by retaining the newest / most recently updated project file per filesystem path and safely removing stale duplicate `.json` project files from `~/.gemini/config/projects/`.

## Synopses & Aliases

```bash
gitmap agy cure-duplicate-projects [flags]
gitmap agy cdp [flags]
gitmap agy optimize-projects [flags]
gitmap agy --repeat-fix [flags]
```

## Flags

| Flag | Shorthand | Type | Default | Description |
|---|---|---|---|---|
| `--except` | `-e` | string | `""` | Exclude project IDs, names, slugs, or paths from deletion |
| `--dry-run` | `-d` | bool | `false` | Preview duplicate removals without deleting any files |
| `--yes` | `-y` | bool | `false` | Confirm removal without interactive confirmation prompt |

## Examples

### 1. Preview deduplication safely (dry run)

```bash
gitmap agy cure-duplicate-projects --dry-run
```
**Alias:**
```bash
gitmap agy cdp -d
```

### 2. Cure all duplicates with interactive prompt

```bash
gitmap agy cure-duplicate-projects
```
**Alias:**
```bash
gitmap agy cdp
```

### 3. Non-interactive automatic cure

```bash
gitmap agy cdp -y
```

### 4. Cure duplicates while preserving specific projects

```bash
gitmap agy cdp --except "gitmap, prompts, f618"
```
*(Accepts project IDs, names, directory base slugs, or short prefix starts with).*
