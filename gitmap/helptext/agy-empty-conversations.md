# gitmap agy remove-projects-with-empty-conversations

Audit and prune Antigravity projects that have no active conversations or only aborted empty sessions.

## Usage

```bash
gitmap agy ls show-projects-with-empty-conversations
gitmap agy remove-projects-with-empty-conversations [flags]
```

## Aliases

- `rm-empty-conversations`
- `clean-empty-conversations`
- `prune-empty-conversations`

## Flags

| Flag | Default | Description |
|---|---|---|
| `--except`, `-e` | `""` | Comma-separated list or file path (.csv/.txt) of project IDs, names, paths, or aliases to preserve |
| `--dry-run`, `-d` | `false` | Preview projects targeted for removal without deleting files |
| `--yes`, `-y` | `false` | Confirm removal without interactive prompt |

## Examples

### List Projects with Empty Conversations

```bash
gitmap agy ls show-projects-with-empty-conversations
```

**Output:**

```text
  ── Antigravity Projects with Empty Conversations ──

  Found 50 project(s) with empty or zero conversations (out of 59 total):

  PROJECT ID                             NAME                     WORKSPACE PATH                             CONVS  STATUS
  ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  098304c1-b79b-4c7e-a223-b4b318b8a6cd   prompts-connect-v3       D:\wp-work\riseup-asia\02-prompts\prom...  0      No Convs
```

### Dry-Run Removal with Exceptions

```bash
gitmap agy remove-projects-with-empty-conversations --except "prompts-connect-v3, extendcore" --dry-run
```

**Output:**

```text
  Targeting 48 project(s) with empty conversations for removal:
    • 1f908b80-7c2c-479a-86d7-021f568fa58a letsmarknow-ui-v2      D:\wp-work\riseup-asia\letsmarkn...

ℹ [dry-run] 48 project(s) would be removed. Remaining: 12
```

### Remove All Empty Projects Non-Interactively

```bash
gitmap agy remove-projects-with-empty-conversations --yes
```

**Output:**

```text
✓ Successfully removed 50 Antigravity project(s) with empty conversations. Remaining active: 9
```

## See Also

- `gitmap agy optimize-projects`: Deduplicate Antigravity projects sharing the same path
- `gitmap agy clear`: Clear all projects with `--except` whitelist
