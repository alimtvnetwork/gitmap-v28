# Remove Missing Projects (`gitmap agy remove-missing-projects`)

Scans all Antigravity project entries in `~/.gemini/config/projects/` and prunes stale project configuration files whose target workspace directory paths no longer exist on disk (such as deleted temporary folders or relocated directories).

## Synopses & Aliases

```bash
gitmap agy remove-missing-projects [flags]
gitmap agy remove-misisng-projects [flags]   # Typo-tolerant alias
gitmap agy rm-missing [flags]
gitmap agy clean-missing [flags]
```

## Flags

| Flag | Shorthand | Type | Default | Description |
|---|---|---|---|---|
| `--except` | `-e` | string | `""` | Exclude project IDs, names, slugs, or short prefix starts with |
| `--dry-run` | `-d` | bool | `false` | Preview missing projects targeted for deletion without deleting |
| `--yes` | `-y` | bool | `false` | Confirm removal without interactive prompt |

## Examples

### 1. Preview missing projects (dry run)
```bash
gitmap agy remove-missing-projects --dry-run
```
**Alias:**
```bash
gitmap agy rm-missing -d
```

### 2. Remove missing projects with confirmation prompt
```bash
gitmap agy remove-missing-projects
```

### 3. Non-interactive automatic cleanup
```bash
gitmap agy remove-missing-projects -y
```

### 4. Remove missing while exempting specific projects by prefix or name
```bash
gitmap agy remove-missing-projects --except "1a9408cc, repo, temp-work"
```

### 5. Typo-tolerant usage
```bash
gitmap agy remove-misisng-projects -y
```
