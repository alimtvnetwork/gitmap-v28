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

## Examples

```bash
gitmap agy clear
gitmap agy clear --except "gitmap, coding-guidelines" --yes
```
