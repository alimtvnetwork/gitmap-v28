# gitmap clone-sync

Clone one or more repositories and seamlessly synchronize them into GitHub Desktop, VS Code Project Manager, and Google Antigravity workspaces.

## Alias

cs

## Usage

    gitmap clone-sync <url1> [url2...] [flags]

## Description

The `clone-sync` (alias: `cs`) command functions exactly like the standard `clone` command, but automatically initiates the unified workspace synchronization engine immediately after a successful clone. 

This ensures that the repository is instantly available across all your favorite IDEs and the Antigravity agent context.

## Integration Highlights

1.  **Antigravity**: A workspace configuration JSON is placed into `~/.gemini/config/projects/` giving autonomous agents immediate access to the repo.
2.  **GitHub Desktop**: The repo is instantly appended to the application's tracked repositories database.
3.  **VS Code**: The repo is recorded in Project Manager so it appears in your side-panel.

## Examples

Clone a single repository and synchronize it:
```bash
gitmap clone-sync https://github.com/example/repo
```

Clone multiple repositories and synchronize them sequentially:
```bash
gitmap cs https://github.com/example/repo1 https://github.com/example/repo2
```
