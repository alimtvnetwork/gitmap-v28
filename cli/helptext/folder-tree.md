# gitmap folder-tree

Scan, visualize, export, and recreate directory hierarchy and repository structures across any filesystem location.

## Synopsis

```bash
gitmap folder-tree [path] [flags]
gitmap folder-tree ls [path] [flags]
gitmap folder-tree view [path] [flags]
gitmap folder-tree export [path] [flags]
gitmap folder-tree import <file> [target] [flags]
gitmap ft [subcommand] [flags]
```

## Aliases

- `folder-tree`
- `foldertree`
- `ft`

## Description

`gitmap folder-tree` (`ft`) scans any local directory path (no prior database indexing or GitMap scan required) and formats the hierarchy with sequence numbers for AI-driven workflows and human readability. Git repositories are automatically discovered with active branch tags. Structures can be visualized as an emoji tree (`--tree`) or a numbered list with single-line gaps (`--preview`), exported to JSON, YAML, tree, or preview format, and reconstituted into directories and 0-byte placeholder files using `import`.

## Subcommands

- `ls` / `view [path]`: Render directory hierarchy with Git detection and sequence numbers.
- `export [path]`: Export hierarchy to JSON, YAML, tree text, or one-line gap preview.
- `import <file> [target-dir]`: Recreate directory skeleton and empty placeholder files from file.
- `help`: Display comprehensive boxed ASCII terminal usage guide.

## Flags & Options

| Flag | Default | Description |
|------|---------|-------------|
| `--tree` | `true` | Emoji tree view with branch connectors |
| `--preview`, `-p`, `--gap` | `false` | Numbered sequence with one-line gap (name + path) |
| `--dirs-only`, `-d` | `false` | Display or export directories only, omit files |
| `--files`, `-f` | `true` | Include files in tree rendering |
| `--git-only` | `false` | Filter to display only Git repositories |
| `--depth <N>`, `-L <N>` | `3` | Maximum scan recursion depth (0 for infinite) |
| `--format <fmt>`, `-fmt` | `json` | Export format: `json`, `yaml`, `tree`, `preview` |
| `--output <file>`, `-o` | | Write exported output directly to specified file |
| `--no-numbers` | `false` | Suppress sequence numbering in output |
| `--all`, `-a` | `false` | Include hidden items (.git internals, node_modules) |
| `--dry-run`, `-n` | `false` | Preview import actions without modifying disk |
| `-h`, `--help` | `false` | Display help menu |

## Examples

### Example 1: View folder tree with depth limit

```bash
gitmap ft --depth 2
```

### Example 2: Numbered prompt preview format with one-line gap

```bash
gitmap ft --preview --depth 1
```

### Example 3: Export directory skeleton to JSON

```bash
gitmap ft export . --format json -o skeleton.json
```

### Example 4: Recreate directory scaffolding from file

```bash
gitmap ft import skeleton.json ./my-new-project
```
