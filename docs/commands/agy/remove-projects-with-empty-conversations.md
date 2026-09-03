# `gitmap agy remove-projects-with-empty-conversations`

Prune Antigravity projects that have no active conversation history or only aborted initialization sessions.

## Usage

```bash
gitmap agy remove-projects-with-empty-conversations [flags]
```

## Aliases

`rm-empty-conversations`, `clean-empty-conversations`, `prune-empty-conversations`

## Flags

| Flag | Default | Description |
|---|---|---|
| `--except`, `-e` | `""` | Comma-separated list or file path (`.csv`/`.txt`) of IDs, names, paths, or aliases to preserve |
| `--dry-run`, `-d` | `false` | Preview projects targeted for removal without deleting files |
| `--yes`, `-y` | `false` | Non-interactive bypass of confirmation prompt |

## Flag Examples

### 1. Interactive Confirmation (Default)
Prompts `Are you sure you want to remove N project(s)? [y/N]` before deleting:
```bash
gitmap agy remove-projects-with-empty-conversations
```

### 2. Preview Candidates Safely (`--dry-run`, `-d`)
Inspects what would be removed without touching any files on disk:
```bash
gitmap agy remove-projects-with-empty-conversations --dry-run
```

### 3. Exclude Specific Projects by Name, ID, or Path (`--except`, `-e`)
```bash
gitmap agy remove-projects-with-empty-conversations -e "gitmap-v28, 0349c4d0-5a91-4f3e-800f-81fd53fc724f"
```

### 4. Exclude Using a Whitelist CSV/Text File (`--except <file.csv>`)
Reads allowed project names or IDs directly from a text or CSV file:
```bash
gitmap agy remove-projects-with-empty-conversations --except "./whitelist-projects.csv" --dry-run
```

### 5. Non-Interactive Headless / CI Execution (`--yes`, `-y`)
```bash
gitmap agy remove-projects-with-empty-conversations --except "production-agent" --yes
```
