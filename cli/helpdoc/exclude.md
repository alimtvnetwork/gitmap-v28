# gitmap exclude

Manage the repository exclusion list — paths or patterns that gitmap skips
during scans and batch operations.

## Usage

```
gitmap exclude add <path-or-pattern> [...]
gitmap exclude rm <path-or-pattern> [...]
gitmap exclude list
gitmap exclude clean
gitmap exclude <path-or-pattern> [...]
```

Aliases: `repo-exclude`, `repo-exclusions`, `exclusions`.

Running `gitmap exclude` with no arguments prints the current exclusion list.
Bare arguments (not matching a subcommand) are treated as patterns to add.

## Subcommands

| Subcommand | Aliases | Behavior |
|------------|---------|----------|
| `add <patterns>` | — | Add paths or patterns to the exclusion list |
| `rm <patterns>` | `remove`, `delete` | Remove entries from the exclusion list |
| `list` | `ls`, `show` | Print the current exclusion list |
| `clean` | `purge` | Purge excluded, cache, and corrupt repos from the database (see `gitmap clean-repos`) |

## Examples

### Add exclusion patterns

```bash
gitmap exclude add "*/node_modules" "*/.git"
gitmap exclude "scratch-*"
```

### List and remove exclusions

```bash
gitmap exclude list
gitmap exclude rm "scratch-*"
```
