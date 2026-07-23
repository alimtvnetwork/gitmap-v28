---
name: deploy-layout-and-binary-readout
description: Deploy folder is gitmap-cli (not gitmap). Bare `gitmap` must print Active/Deployed/Config binary triplet. Reuses existing install location via PATH lookup.
type: feature
---
# Deploy layout & binary readout

## Deploy folder naming
- **Default deploy subfolder is `gitmap-cli`** (not `gitmap`). Full default path on Windows: `E:\bin-run\gitmap-cli\gitmap.exe`.
- The legacy subfolder name was `gitmap`, which collided visually with the binary name and caused user confusion (`E:\gitmap\gitmap.exe` looked like a typo). New layout: `<deployRoot>\gitmap-cli\gitmap.exe`.
- The constant `constants.GitMapSubdir` (currently `"gitmap"`) MUST be renamed to `GitMapCliSubdir = "gitmap-cli"` and every reference (run.ps1 `Join-Path $target "gitmap"`, `constants_update.go` `Join-Path $cfg.deployPath "gitmap\gitmap.exe"`, `doctorfixpath.go` `filepath.Join(deployPath, constants.GitMapSubdir, binaryName)`, `updatecleanup_paths.go::resolveConfigDeployAppDir`) updated to match.
- Migration: when a legacy `<deployRoot>\gitmap\gitmap.exe` is detected on next run/update, move it to `<deployRoot>\gitmap-cli\gitmap.exe` and delete the empty `gitmap\` folder. Idempotent — no-op if already migrated.

## Deploy target resolution (already implemented in run.ps1::Resolve-DeployTarget)
Priority order, do not change:
1. `-DeployPath` CLI flag override.
2. **PATH lookup**: if `gitmap` is already on PATH, deploy to its parent's parent (so re-deploys land in the same place the user originally chose, regardless of what powershell.json says).
3. `powershell.json` `deployPath` (default `E:\bin-run`).

## Bare `gitmap` (no args) behavior
- Today: `len(os.Args) < 2` → calls `printUsage()` and exits 1. Does NOT print the binary location triplet.
- **Required**: bare `gitmap` should ALSO print the same Active/Deployed/Config binary triplet that `gitmap update` shows at completion, BEFORE the usage text. Format:

  ```
  Active binary:    <os.Executable() resolved>
  Deployed binary:  <powershell.json deployPath + GitMapCliSubdir + binaryName, if exists>
  Config binary:    <whatever powershell.json declares, even if missing>

  <usage text...>
  ```

- Rationale: users frequently confused about which binary they're hitting (PATH vs latest deploy vs config). Showing this on every bare invocation eliminates the "why does -version not match?" support churn.
- Implementation: extract the existing PowerShell readout block from `constants_update.go::UpdatePSPostSyncReport` into a Go helper `cmd/binarylocations.go::PrintBinaryLocations()`, call it from both `root.go::Run()` (when `len(os.Args) < 2`) and the post-update flow.

## Active vs Deployed vs Config
- **Active binary** = `os.Executable()` after `filepath.EvalSymlinks` — the file the OS actually loaded for this process.
- **Deployed binary** = `<powershell.json.deployPath>/<GitMapCliSubdir>/<binaryName>` if it exists on disk. Empty string otherwise.
- **Config binary** = literal value of `<powershell.json.deployPath>/<GitMapCliSubdir>/<binaryName>` whether or not the file exists — represents what the config *thinks* should be there.
- All three may match; the readout only gets noisy when they diverge (which is when the user needs to see it).
