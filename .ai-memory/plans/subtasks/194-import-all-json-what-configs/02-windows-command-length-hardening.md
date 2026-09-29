# Subtask 02: Windows Fleet Deploy Command-Line Length Hardening
- Windows `cmd.exe` 8,191-character buffer limit triggered `The command line is too long.` during fleet deployment of 6-node manifest.
- Updated `ssh_deploy_node_config.go` and `ssh_deploy_config_ssh.go` to stream base64 payloads >= 1,500 bytes directly to `.gitmap/fleet_config_deploy.json` via `streamDirectToRemote`.
- Remote node imports file with `gitmap ssh nodes import-json .gitmap/fleet_config_deploy.json` and safely cleans up in a defer.
- Fixed bidirectional SSH communication check in `sshjoin_enroll.go` to use PowerShell syntax on Windows targets.
- Verified live deployment across `main` and `w3` successfully without length errors.
- Status: Completed.
