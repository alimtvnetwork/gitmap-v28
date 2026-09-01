# Specification 17 — Chapter 5: VS Code & GitHub Desktop Lifecycle Integration

## 1. VS Code Project Manager Sync

File: `~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json` (or Windows `%APPDATA%\Code\User\globalStorage\alefragnani.project-manager\projects.json`)

### Operations:

1. **`UpdateRootPath(oldPath, newPath, newName string)`**:
   - Reads `projects.json`.
   - Locates entry with `rootPath == oldPath` (using case-insensitive normalized comparison).
   - Updates `rootPath = newPath` and `name = newName`.
   - Atomically writes updated `projects.json`.
2. **`RemoveEntry(targetPath string)`**:
   - Removes entry whose `rootPath == targetPath`.
   - Atomically writes updated `projects.json`.

## 2. GitHub Desktop Tracking Lifecycle

1. **`UpdateRepoPath(oldPath, newPath string)`**:
   - Registers new repository location with GitHub Desktop CLI (`github-desktop <newPath>`).
2. **`RemoveRepo(repoPath string)`**:
   - Informs GitHub Desktop integration layer.

## 3. Interactive Prompt & `-y` Flag

- When moving or deleting, user is prompted:
  `Update VS Code Project Manager and GitHub Desktop? [Y/n]`
- With `-y` / `--yes`, auto-applies without blocking.
