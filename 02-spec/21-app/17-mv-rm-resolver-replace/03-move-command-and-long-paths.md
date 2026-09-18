# Specification 17 — Chapter 3: Move (`gitmap mv`) Command & Windows Long Paths

## 1. CLI Signature & Syntax

```bash
gitmap mv [options] <source> <destination>
```
Aliases: `gitmap move <source> <destination>`

### Arguments:

- `<source>`: Repo slug, alias, relative path (`.\prompt-architect`, `./subfolder`), or absolute path.
- `<destination>`: Target directory. Supports:
  - `..` (moves repository into the parent folder)
  - Relative directory (`../new-home`, `./archive/prompt-architect`)
  - Absolute directory (`D:\work\archive\prompt-architect`)

### Options:

- `-y`, `--yes`: Automatically confirm all prompts (including VS Code and GitHub Desktop updates).
- `--dry-run`: Simulate directory relocation and database updates without moving files.
- `--no-vscode`: Skip updating VS Code Project Manager configuration.
- `--no-desktop`: Skip updating GitHub Desktop tracking.

## 2. Windows Long Path Handling (`\\?\`)

On Windows hosts, paths exceeding `MAX_PATH` (260 characters) can cause silent filesystem failure.
The move engine encapsulates a Windows Long Path utility:
- Detects Windows platform (`runtime.GOOS == "windows"`).
- Automatically applies the extended-length path prefix (`\\?\` for local drives, `\\?\UNC\` for network shares) when path length exceeds 240 characters or when standard operations fail with path errors.
- Cleans and strips prefixes before writing paths into SQLite database to preserve portable storage.

## 3. Relocation Workflow Steps

1. **Resolve Source**: Run Unified Project Resolver on `<source>`.
2. **Resolve Destination**:
   - If destination is `..`, compute `filepath.Join(filepath.Dir(srcParent), repoDirName)`.
   - If destination is an existing directory, append the source repo folder name.
3. **Safety Preflight Checks**:
   - Source directory must exist.
   - Destination directory must NOT already exist.
   - Ensure destination is not inside source repository.
4. **Physical Filesystem Relocation**:
   - Attempt atomic `os.Rename(src, dest)`.
   - Fall back to cross-volume copy + verify + remove if moving across different drive letters (e.g. `C:\` to `D:\`).
5. **Database Atomic Update**:
   - Update `Repo.AbsolutePath` and `Repo.RepoName`.
   - Update `ScanFolderId` if moving to a new managed scan folder.
   - Update `Alias` records pointing to old path.
6. **External Integration Updates**:
   - Update VS Code Project Manager `projects.json`.
   - Update GitHub Desktop repository path.
