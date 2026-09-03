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

## Examples

```bash
gitmap agy ls
gitmap agy ls show-projects-with-empty-conversations
```
