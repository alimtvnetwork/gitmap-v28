# gitmap multiclone

Batch clone multiple repositories from multiline code blocks, Markdown pastes, files, or stdin pipes with automatic deduplication, inline description stripping, and VS Code Project Manager sync.

## Aliases

- `mc`
- `mutliclone`

## Usage

```bash
gitmap multiclone <urls|paste|file> [target-dir] [flags]
gitmap mc <paste> [target-dir] [flags]
gitmap mutliclone <paste> [target-dir] [flags]
```

## Description

`gitmap multiclone` (short alias `mc`) simplifies bulk repository cloning:

1. **Multiline Markdown Pastes**: Accepts text blocks copied directly from web pages, markdown documents, or chat terminals wrapped in triple backticks (```` ``` ````) or raw lines.
2. **Flexible Formats**: Recognizes full Git URLs (`https://`, `http://`, `git@`, `ssh://`), Markdown links (`[name](url)`), and `owner/repo` shorthands (e.g. `ChrisTitusTech/linutil`), auto-expanding them to Git clone URLs.
3. **Description Stripping**: Automatically strips descriptive text appended to repository titles (such as `owner/repo: Brief overview`).
4. **Intelligent Deduplication**: Deduplicates repeated entries case-insensitively, normalizing trailing `.git` suffixes and URL variants.
5. **Target Directory Support**: Specify an optional target folder (`-d <dir>` or as an argument) to clone all repositories into `<dir>/<repo-name>`.
6. **VS Code Project Manager Sync**: Synchronizes all successfully cloned repositories to VS Code Project Manager.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-d, --dir, --target <dir>` | current directory | Target folder to place cloned repositories |
| `-f, --file <path>` | none | Path to file containing repositories to clone |
| `-j, --concurrency <N>` | 1 | Number of concurrent clone operations |
| `--dry-run` | false | Preview parsed repositories and destination paths without cloning |
| `--no-replace` | false | Do not replace or overwrite existing destination folders |
| `--no-vscode-sync` | false | Skip registering cloned repositories in VS Code Project Manager |
| `--desktop` | false | Register cloned repositories in GitHub Desktop |
| `-h, --help` | false | Display the rich two-column styled help menu |

## Examples

### 1. Paste multiline markdown block
```bash
gitmap mc "```
ChrisTitusTech/ChrisTitusTech
https://github.com/ChrisTitusTech/ChrisTitusTech

ChrisTitusTech/linutil: Chris Titus Tech's Linux Toolbox
https://github.com/ChrisTitusTech/linutil
```"
```

### 2. Specify target destination directory
```bash
gitmap mc -d ./community-tools "```
ChrisTitusTech/winutil
ChrisTitusTech/image-tools
```"
```

### 3. Read repository list from text file
```bash
gitmap mc -f repos.txt -d ./repos
```

### 4. Pipe via standard input
```bash
cat repos.txt | gitmap mc
```

### 5. Dry-run preview
```bash
gitmap mc repos.txt --dry-run
```
