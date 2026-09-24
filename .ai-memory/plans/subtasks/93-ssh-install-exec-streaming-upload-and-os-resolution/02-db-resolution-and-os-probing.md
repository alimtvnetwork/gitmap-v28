# Subtask 02: Dynamic Remote OS Probing & Global DB Fallback Resolution
Traceability ID: Task-02
Spec Reference: [02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md](../../../02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md)
Target Files: [cli/cmdssh/ssh_install_exec.go, cli/db/sshconnection.go, cli/store/location.go, cli/cmdssh/ssh_filter.go]
Action: Always probe remote OS upon connection (probeRemoteOSType(client)) in ssh_install_exec.go, persist updated OS back to DB, fix mergeSSHHostsConnections to avoid blind "linux" defaulting, and add AppData gitmap-cli/data fallback in store.OpenDefault().
Acceptance Criteria:
- Node w1 (Windows Server) is detected as "windows" and targeted with Windows paths (C:\Windows\Temp).
- Invoking gitmap ssh install-exec --dry-run from arbitrary folders (Downloads) discovers all 5 nodes.
Targeted Verification: go test ./cli/cmdssh -run TestFilter
