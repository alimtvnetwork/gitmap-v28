# GitMap v6.606.0

## What's Changed in v6.606.0

- **Single-Repo Pull Isolation (`gitmap pull`)**: Fixed `gitmap pull` without arguments inside a repository directory (`cwd`) to pull exclusively the active repository instead of defaulting to `pull-all` across all cataloged workspaces.
- **Dynamic Invocation Target Scope**: Enhanced `printPullInvocationHeader` to dynamically resolve target scope `(repo-name)` from working directory or arguments, reserving `(all workspaces)` strictly for explicit `gitmap pull-all`, `pa`, or non-git fallback directory contexts.
- **Direct Filesystem Path Resolution**: Enhanced target slug resolver to automatically detect and pull direct filesystem directory paths matching existing Git repositories without requiring prior SQLite catalog enrollment.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.606.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.606.0/install.sh | sh
```
