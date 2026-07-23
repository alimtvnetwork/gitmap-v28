# 89 — Deploy layout & binary readout

## Deploy folder convention (v3.6+)

The gitmap-v28 binary deploys to `<deployRoot>/gitmap-cli/gitmap.exe`, NOT `<deployRoot>/gitmap-v28/gitmap.exe`.

| Component | Value |
|---|---|
| Default deploy root (Windows) | `E:\bin-run` |
| Deploy subfolder | `gitmap-cli` |
| Binary name | `gitmap.exe` |
| Full default path | `E:\bin-run\gitmap-cli\gitmap.exe` |

The legacy subfolder name was `gitmap-v28`. That created visual collision with the binary name (`E:\gitmap-v28\gitmap.exe` looked like a typo) and confused users about whether they were looking at the deploy root, the app folder, or the binary itself. The rename is forward-compatible: a one-time migration in `run.ps1::Repair-DeployLayout` moves any legacy `<root>/gitmap-v28/gitmap.exe` → `<root>/gitmap-cli/gitmap.exe` and removes the empty legacy folder.

## Deploy target resolution

`run.ps1::Resolve-DeployTarget` priority:

1. `-DeployPath` CLI flag — explicit override always wins.
2. **PATH detection** — if `gitmap-v28` is already on PATH (`Get-Command gitmap-v28`), the deploy target is the parent of that binary's parent folder. This makes `run.ps1` "follow" the user's existing install regardless of what `powershell.json` says, so `git pull && .\run.ps1` always updates the binary the user is actually invoking.
3. `powershell.json` `deployPath` field (default `E:\bin-run`).

After every successful deploy, `Sync-ConfigDeployPath` rewrites `powershell.json` `deployPath` to match the actual install location, so the "Config binary:" readout stays in sync.

## Bare-invocation binary readout

Running `gitmap-v28` with no arguments prints a three-line readout BEFORE the usage text. The readout always prints (even when all three paths match) so users build a habit of recognising which binary they're hitting; CI scripts and pipelines that capture gitmap-v28 output can suppress it with `--no-banner` or by setting `GITMAP_QUIET=1`:

```
  Active binary:    E:\bin-run\gitmap-cli\gitmap.exe
  Deployed binary:  E:\bin-run\gitmap-cli\gitmap.exe
  Config binary:    E:\bin-run\gitmap-cli\gitmap.exe

  gitmap-v28 v3.6.0
  ...usage...
```

Definitions:

- **Active binary** — `os.Executable()` after `filepath.EvalSymlinks`. The file the OS actually loaded for this process.
- **Deployed binary** — `<powershell.json.deployPath>/gitmap-cli/gitmap.exe` if the file exists on disk; empty otherwise.
- **Config binary** — the literal path that `powershell.json` declares, whether or not the file exists. Represents config intent.

When all three match, the readout is informational. When they diverge, it pinpoints the exact source of "wrong version" or "stale binary" issues without requiring `gitmap-v28 doctor`.

## Legacy layout migration

When `run.ps1` runs and detects the legacy `<deployRoot>/gitmap-v28/gitmap.exe` layout, `Repair-DeployLayout` silently moves the binary to `<deployRoot>/gitmap-cli/gitmap.exe` and removes the empty legacy `gitmap-v28/` folder. No prompt, no user action required. Idempotent — re-runs are no-ops once migrated. The bare-invocation readout will then naturally show the new path next time the user invokes `gitmap-v28`.

- **Active binary** — `os.Executable()` after `filepath.EvalSymlinks`. The file the OS actually loaded for this process.
- **Deployed binary** — `<powershell.json.deployPath>/gitmap-cli/gitmap.exe` if the file exists on disk; empty otherwise.
- **Config binary** — the literal path that `powershell.json` declares, whether or not the file exists. Represents config intent.

When all three match, the readout is informational. When they diverge, it pinpoints the exact source of "wrong version" or "stale binary" issues without requiring `gitmap-v28 doctor`.

## Implementation

| File | Change |
|---|---|
| `gitmap-v28/constants/constants_doctor.go` | `GitMapSubdir = "gitmap-v28"` → `GitMapCliSubdir = "gitmap-cli"` |
| `gitmap-v28/cmd/root.go` | `Run()` calls `PrintBinaryLocations()` before `printUsage()` when `len(os.Args) < 2` |
| `gitmap-v28/cmd/binarylocations.go` (new) | Resolves and prints Active/Deployed/Config triplet |
| `gitmap-v28/constants/constants_update.go` | `Join-Path $cfg.deployPath "gitmap-v28\gitmap.exe"` → `"gitmap-cli\gitmap.exe"` |
| `gitmap-v28/cmd/doctorfixpath.go` | `filepath.Join(deployPath, constants.GitMapSubdir, binaryName)` uses new constant |
| `gitmap-v28/cmd/updatecleanup_paths.go::resolveConfigDeployAppDir` | Uses new constant |
| `run.ps1::Deploy-Binary` | `Join-Path $target "gitmap-v28"` → `Join-Path $target "gitmap-cli"` |
| `run.ps1::Repair-DeployLayout` | Migrate legacy `<root>/gitmap-v28/` → `<root>/gitmap-cli/` |
