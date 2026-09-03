# VS Code Project Manager Integration (`gitmap vscode`)

Synchronize GitMap repositories into the VS Code Project Manager extension (`projects.json`).

## Commands

### `gitmap vscode sync`
Scans tracked GitMap repositories and populates `projects.json` for Project Manager.
```bash
gitmap vscode sync
```

### `gitmap vscode ls`
Lists all projects registered in VS Code Project Manager.
```bash
gitmap vscode ls
```

### `gitmap vscode optimize-projects`
Removes duplicate project path entries from `projects.json`.
```bash
gitmap vscode optimize-projects
```

### `gitmap vscode clear`
Clears projects whose paths are no longer present on disk (`--except` supported).
```bash
gitmap vscode clear --except "my-project"
```

### `gitmap vscode paths <add|rm|list>`
Manually manages explicit path alias mappings in the registry.
```bash
gitmap vscode paths add my-alias D:\projects\my-repo
gitmap vscode paths list
gitmap vscode paths rm my-alias
```
