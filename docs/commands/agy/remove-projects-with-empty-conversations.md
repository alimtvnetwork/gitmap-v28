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

## How to Proceed

1. **Audit First:**
   ```bash
   gitmap agy ls show-projects-with-empty-conversations
   ```
2. **Dry-Run with Whitelist:**
   ```bash
   gitmap agy remove-projects-with-empty-conversations --except "prompts-connect-v3, extendcore" --dry-run
   ```
3. **Execute Cleanup:**
   ```bash
   gitmap agy remove-projects-with-empty-conversations --yes
   ```
