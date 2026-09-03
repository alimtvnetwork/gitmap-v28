# `gitmap agy clear`

Clear missing or stale Antigravity workspace projects.

## Usage

```bash
gitmap agy clear [flags]
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--except`, `-e` | `""` | Comma-separated list or file path of IDs, names, or paths to preserve |
| `--missing`, `-m` | `true` | Only remove projects whose folders no longer exist on disk |
| `--yes`, `-y` | `false` | Skip confirmation prompt |

## Flag Examples

### 1. Clear Missing Only (Default Behavior)
Removes only project files whose associated folders no longer exist on disk:
```bash
gitmap agy clear
```

### 2. Clear All Projects Except Whitelisted (`--except`, `-e`)
```bash
gitmap agy clear --except "gitmap-v28, ai-runner"
```

### 3. Clear Using Whitelist CSV File
```bash
gitmap agy clear --except "my-whitelist.csv"
```

### 4. Non-Interactive Headless Mode (`--yes`, `-y`)
```bash
gitmap agy clear --yes
```
