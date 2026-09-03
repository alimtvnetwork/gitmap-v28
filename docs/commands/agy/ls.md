# `gitmap agy ls`

List registered Antigravity projects and audit conversation activity.

## Usage

```bash
gitmap agy ls [flags]
gitmap agy ls show-projects-with-empty-conversations
```

## Subcommands & Aliases

- `show-projects-with-empty-conversations`: Filters projects whose conversation databases have zero user turns or failed initialization.
- Aliases: `show-proects-with-empty-conversations`, `empty-conversations`, `empty-convs`.

## Flags

| Flag | Default | Description |
|---|---|---|
| `--missing`, `-m` | `false` | Show only projects whose workspace paths are missing on disk |
| `--active`, `-a` | `false` | Show only active projects |
| `--sort`, `-s` | `name` | Sort by `name` or `time` |
| `--filter`, `-f` | `""` | Filter projects by name or path |
| `--json` | `false` | Emit results as structured JSON |

## Examples for All Flags

### 1. Default Listing

```bash
gitmap agy ls
```

### 2. Detect Missing / Stale Workspaces (`--missing`, `-m`)

```bash
gitmap agy ls --missing
```

### 3. Filter by Name or Slug (`--filter`, `-f`)

```bash
gitmap agy ls -f "riseup"
```

### 4. Sort by Modification Timestamp (`--sort`, `-s`)

```bash
gitmap agy ls --sort time
```

### 5. Show Only Active Projects (`--active`, `-a`)

```bash
gitmap agy ls --active
```

### 6. Emit Machine-Readable JSON (`--json`)

```bash
gitmap agy ls --json
```

### 7. Audit Empty / Aborted Conversation Projects

```bash
gitmap agy ls show-projects-with-empty-conversations
```
