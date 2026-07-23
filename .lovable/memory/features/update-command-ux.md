# Update Command UX

The `update` command employs a robust 4-tier repository resolution strategy to maintain source connectivity:

1. **`--repo-path` flag** — Immediate override; the user-supplied path is trusted without existence check and saved to the DB.
2. **Embedded `RepoPath` constant** — For source-linked builds; set at compile time via `-ldflags`.
3. **SQLite Settings table** — Persistent lookup of `source_repo_path` from the binary's `data/` directory (via `store.OpenDefault()`).
4. **Interactive user prompt** — Fallback when all other tiers fail. If the user provides a path that already exists, it is validated as a gitmap source root. **If the path does not exist**, the system clones the gitmap source repository into that directory via `git clone` (`cloneRepoInto()`), then re-validates the cloned path. Successfully resolved paths are saved to the DB so future runs do not re-prompt.

Every candidate path is validated to ensure it is a legitimate gitmap source root by verifying the existence of `run.ps1` and the `constants/constants.go` marker; the system can automatically traverse upward to find the root if a subfolder is provided.

If no repository can be resolved across all four tiers, the command delegates to the standalone `gitmap-updater` tool (if on PATH), or prints the `ErrNoRepoPath` recovery guide with manual fix options.

## Key Constants

| Constant | Purpose |
|----------|---------|
| `SourceRepoCloneURL` | GitHub URL used for clone-on-missing-path |
| `MsgUpdateCloning` | Status message when cloning begins |
| `MsgUpdateCloneOK` | Success message after clone |
| `ErrUpdateCloneFailed` | Error message if clone fails |
| `SettingSourceRepoPath` | SQLite key shared with `release-self` |

## Files

| File | Purpose |
|------|---------|
| `cmd/update.go` | `resolveRepoPath()` with 4-tier dispatch |
| `cmd/updaterepo.go` | Path helpers: prompt, clone, DB read/write |
| `constants/constants_update.go` | Clone and path recovery constants |
| `release/selfrelease_resolve.go` | Shared DB functions for `release-self` |
