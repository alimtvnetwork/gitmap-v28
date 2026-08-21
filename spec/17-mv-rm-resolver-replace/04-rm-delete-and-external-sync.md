# Specification 17 — Chapter 4: Enhanced `gitmap rm` (Remove & Delete)

## 1. CLI Signature & Syntax
```bash
gitmap rm [-y|--yes] [--db-only] [--no-vscode] [--no-desktop] <target>[,<target>...] [<target>...]
```
Aliases: `gitmap remove`, `gitmap del`, `gitmap delete`

## 2. Path-Aware Targeting
Resolves all target formats seamlessly:
- `gitmap rm .\prompt-architect` (Windows relative path)
- `gitmap rm ./prompt-architect/` (POSIX relative with trailing slash)
- `gitmap rm prompt-architect` (Slug or alias)
- `gitmap rm .` (Current working directory repo)
- `gitmap rm macro*` (Glob match across slug and folder name)

## 3. Removal & Cascade Execution
1. **Interactive Prompt**:
   - Unless `-y` / `--yes` is passed, prompts:
     `Delete folder and untrack prompt-architect (D:\work\prompt-architect)? [y/N]`
2. **Physical Disk Deletion**:
   - If `--db-only` is omitted, removes the folder recursively from disk with long-path safety.
3. **Database Cascade Untrack**:
   - Deletes from `Repo`, `GroupRepo`, `Alias`, `Bookmark`, `VersionProbe`, `DetectedProject`.
4. **External Integrations Cascade**:
   - Removes entry from VS Code `projects.json`.
   - Cleans tracking in GitHub Desktop.
