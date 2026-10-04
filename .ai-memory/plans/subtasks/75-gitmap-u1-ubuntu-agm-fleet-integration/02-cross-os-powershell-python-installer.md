# Subtask 75.2: Cross-OS PowerShell & Toolchain Provisioning in GitMap

- **Parent Plan:** [75-gitmap-u1-ubuntu-agm-fleet-integration.md](../../pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md)
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Status:** Completed
- **Target Area:** `cli/cmdinstall`, `cli/crossplatform`

## Objective
Enable GitMap to autonomously detect missing shell engines (PowerShell / Python) across operating systems, and provide an automated cross-OS PowerShell installer for Ubuntu so that `gitmap ps` works consistently across the fleet.

## Implementation Details
1. Check toolchain status on node `U1` (`pwsh`, `python3`, `go`, `git`).
2. Implement cross-OS PowerShell installation workflow in GitMap (`gitmap install pwsh` / standalone shell adapter):
   - Package manager route: Register Microsoft repository on Debian/Ubuntu and install `powershell`.
   - Portable binary route: Download standalone `.tar.gz` archive and extract to `~/.local/bin/pwsh` with zero root requirement.
3. Verify remote execution of PowerShell scripts via GitMap.
