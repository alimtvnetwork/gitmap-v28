# GitMap v6.450.0

## What's Changed in v6.450.0

- **Auto Fast-Forward Merge Fallback on Pull:** Safe pull now automatically falls back to clean auto-merge (`--no-rebase --no-edit --autostash`) when local and remote branches have diverged, preventing batch pull failures across diverged repositories.
- **Merge Abort on Conflict:** If non-fast-forward merge encounters genuine merge conflicts, GitMap immediately aborts the in-progress merge (`git merge --abort`) to ensure the working tree is never left in an uncommitted conflicted state.
- **Pipeline Cancel Filtering:** Excluded superseded/cancelled workflow runs from false-positive failure flags in pipeline error logging, preventing spurious CI error reports when newer commits cancel in-flight jobs.
- **Split-DB GitIgnore Cache Integration:** Centralized SQLite Split-DB gitignore cache in `BinaryDataDir` with case-insensitive `COLLATE NOCASE` matching and dynamic TTL settings support.
- **Force Flag Bypass:** Added `--force` / `-f` cache bypass across `gitmap ignore`, `gitmap fix-ignore-all`, and `gitmap pull` routines.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.450.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.450.0/install.sh | sh
```
