# gitmap merge

Comprehensive multi-repository merge and commit transfer suite.

## Usage

```bash
gitmap merge-right <left> <right> [flags]  # One-way merge writing only into RIGHT
gitmap merge-left <left> <right> [flags]   # One-way merge writing only into LEFT
gitmap merge-both <left> <right> [flags]   # Bidirectional merge between LEFT and RIGHT
gitmap merge-in <target> <inputs...>       # Chronological multi-source commit replay into target
gitmap commit-in <target> <inputs...>      # Canonical commit replay engine
gitmap commit-both <left> <right> [flags]  # Bidirectional commit replay
gitmap commit-left <left> <right> [flags]  # Replay commits from right repo onto left repo
gitmap commit-right <left> <right> [flags] # Replay commits from left repo onto right repo
gitmap merge write --title <title>         # Open a GitHub pull request from CLI
```

## Subcommands

- `gitmap merge right <left> <right>`: Subcommand alias for `merge-right` (`mr`)
- `gitmap merge left <left> <right>`: Subcommand alias for `merge-left` (`ml`)
- `gitmap merge both <left> <right>`: Subcommand alias for `merge-both` (`mb`)
- `gitmap merge in <target> <inputs...>`: Subcommand alias for `commit-in` (`cin`)
- `gitmap merge write`: Open GitHub pull request from current branch

## Flags

| Flag | Default | Description |
|---|---|---|
| `-y`, `--yes`, `-a` | `false` | Bypass prompts; auto-accept default conflict resolution |
| `--prefer-left` | `false` | LEFT always wins on conflict (overwrite RIGHT) |
| `--prefer-right` | `false` | RIGHT always wins on conflict (overwrite LEFT) |
| `--prefer-newer` | `false` | Newer file modification time wins |
| `--no-push` | `false` | Skip git push on remote URL endpoints |
| `--dry-run` | `false` | Preview all merge actions without touching filesystem |
| `-h`, `--help` | `false` | Display merge suite help menu |

## Examples

```bash
# Publish local repository changes to remote repository without altering local
gitmap merge-right ./local https://github.com/owner/repo -y

# Pull upstream changes into local repository with upstream winning conflicts
gitmap merge-left ./local https://github.com/owner/upstream -y

# Bidirectionally merge two repositories
gitmap merge-both ./repo-a ./repo-b -y

# Replay all versioned sibling histories chronologically into canonical target
gitmap merge-in ./canonical all --default

# Open a GitHub pull request from CLI
gitmap merge write --title "Feature: add terminal help parity" --draft
```
