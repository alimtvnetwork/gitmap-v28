# Repo Path Sync on Build/Deploy

Part of [Self-Update Overview](01-self-update-overview.md).

## Problem

When the source repository is moved or cloned to a new location, the CLI's
database still stores the old repo path. The next `tool update` command
attempts to use the stale path, causing a failure or prompting the user
unnecessarily.

## Solution: Post-Deploy Repo Path Sync

After every successful deploy, the build script (`run.ps1` / `run.sh`)
calls a hidden CLI command to persist the current repo root in the database.

```
# run.ps1 — after Deploy-Binary completes:
& <deployed-binary> set-source-repo <repo-root>
```

This ensures the `source_repo_path` setting in the database always reflects
the location from which the latest build was performed.

## CLI Command: `set-source-repo`

A hidden (unlisted in help) command that accepts a single path argument:

```
tool set-source-repo <path>
```

### Behavior

1. Validate the path is a valid source repo root (has `.git/`, correct `go.mod` module marker, and expected project structure).
2. Normalize to absolute path.
3. Persist to the Settings table (`source_repo_path` key).
4. Print confirmation: `Source repo path saved: <path>`.

### Error Cases

| Condition | Behavior |
|-----------|----------|
| No path argument | Print error, exit 1 |
| Path is not a valid source repo | Print error with path, exit 1 |
| DB write fails | Print warning to stderr (non-fatal in deploy context) |

## Build Script Integration

### PowerShell (`run.ps1`)

After the `Deploy-Binary` function completes successfully:

```powershell
# Sync source repo path in DB so "gitmap-v27 update" uses this repo location
$syncBinary = $destFile
if (-not (Test-Path $syncBinary)) { $syncBinary = $BinaryPath }
if (Test-Path $syncBinary) {
    try {
        & $syncBinary set-source-repo $RepoRoot 2>&1 | Out-Null
        Write-Info "Source repo path synced to DB: $RepoRoot"
    } catch {
        Write-Warn "Could not sync source repo path: $_"
    }
}
```

### Bash (`run.sh`)

After successful deploy:

```bash
"$deployed_binary" set-source-repo "$repo_root" 2>/dev/null || true
```

## Resolution Priority (Unchanged)

The `update` command resolves the repo path in this order:

1. `--repo-path` CLI flag (highest priority)
2. Embedded constant (compiled via `-ldflags`)
3. **Database `source_repo_path` setting** (updated by this mechanism)
4. Interactive prompt (fallback)

The database entry (tier 3) is the one kept current by the post-deploy sync.
When the repo moves, the next `run.ps1` deploy updates it automatically,
so tier 3 always reflects the latest build location.

## Critical Rules

- The `set-source-repo` command MUST validate the path before saving. Never persist an invalid or non-existent repo path.
- The build script MUST suppress output (`| Out-Null` / `2>/dev/null`) so the sync does not clutter normal deploy output.
- The sync MUST be non-fatal — a failure to update the DB should warn but not block the deploy.
- The command MUST NOT appear in help output or usage text (hidden/internal command).
- The command MUST NOT contain interactive prompts (`Read-Host`, `fmt.Scan`, etc.).

## Cross-References

- Repo path resolution: [01-self-update-overview.md](01-self-update-overview.md) §Repo Resolution
- Settings table schema: `constants/constants_settings.go` — `SettingSourceRepoPath`
- Build script pipeline: [04-build-scripts.md](04-build-scripts.md)
- Self-update orchestration: [../03-general/02f-self-update-orchestration.md](../03-general/02f-self-update-orchestration.md)
