# Subtask 91-03: Remote OS Detection & Host Info (`gitmap os-info`)

Spec Reference: [02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md](../../../../02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md)
Parent Plan: [.ai-memory/plans/completed/91-ssh-join-common-scan-and-agy-rop.md](../../completed/91-ssh-join-common-scan-and-agy-rop.md)

## Objective
Enhance remote OS detection during SSH node join to capture deep OS metadata (OS type, OS version, OS architecture), persist into SQLite databases, and provide a standalone `gitmap os-info [--json]` command for local/remote consumption.

## Functional Requirements
1. **Deep Remote OS Detection:**
   - In `cli/cmdssh/sshjoin_enroll.go`, expand `probeRemoteOSDetails` to accurately probe:
     - Distro / OS Type: `windows`, `ubuntu`, `debian`, `centos`, `fedora`, `macos`, `linux`.
     - OS Version: `Windows 11 Pro 10.0.22631`, `Ubuntu 22.04.3 LTS`, etc.
     - Architecture: `amd64`, `arm64`, `x86_64`.
   - Update `db.SSHConnection` table schema and `store.SSHHost` table to ensure `os_type`, `os_version`, and `os_arch` are persisted on first join.
2. **Local & Remote OS Info Command:**
   - Implement `gitmap os-info [--json]`.
   - Output fields: `osType`, `osVersion`, `architecture`, `hostname`, `numCpu`, `kernel`, `platform`.
   - In `--json` mode, outputs compact or formatted JSON for programmatic integration.
   - In text mode, renders an aligned summary table.

## Files to Create/Modify
- `cli/cmdos/os_info_types.go` [NEW]
- `cli/cmdos/os_info.go` [NEW]
- `cli/cmdos/os_info_cmd.go` [NEW]
- `cli/cmdos/os_info_test.go` [NEW]
- `cli/cmdssh/sshjoin_enroll.go` [MODIFY]
- `cli/db/sshconnection.go` [MODIFY]
- `cli/cmd/root.go` & `cli/cmd/rootdispatch.go` [MODIFY]
