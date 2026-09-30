# Subtask 04: Harmonize Machine Folders (`04-w1-machine` .. `07-final-network-machine`) SSH Node Manifests
Traceability ID: Task-01, Task-04
Spec Reference: [02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md](../../../02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md)
Target Files: D:\work\repo-secrets\07-final-network-machine\gitmap-ssh-nodes.json, D:\work\repo-secrets\07-final-network-machine\gitmap-ssh.json, D:\work\repo-secrets\04-w1-machine\gitmap-ssh-nodes.json, D:\work\repo-secrets\04-w1-machine\gitmap-ssh.json, D:\work\repo-secrets\05-w2-machine\gitmap-ssh-nodes.json, D:\work\repo-secrets\05-w2-machine\gitmap-ssh.json, D:\work\repo-secrets\06-w3-machine\gitmap-ssh-nodes.json, D:\work\repo-secrets\06-w3-machine\gitmap-ssh.json
Action:
1. Update `07-final-network-machine\gitmap-ssh-nodes.json`, `07-final-network-machine\gitmap-ssh.json`, `04-w1-machine\gitmap-ssh-nodes.json`, `04-w1-machine\gitmap-ssh.json`, `05-w2-machine\gitmap-ssh-nodes.json`, `05-w2-machine\gitmap-ssh.json`, `06-w3-machine\gitmap-ssh-nodes.json`, and `06-w3-machine\gitmap-ssh.json` to match `01-gitmap\04-ssh-nodes.json`:
   - `version: "2.0"`, lowerCamelCase keys (`workerId`, `ipAddress`, `authMethod`, `keyPath`, `encryptedPassword`, `firstRunAt`, `createdAt`).
   - Structured `attributes.workDirectory` object with `variables` (`workDir`, `userHome`, `sshDir`, `keyPath`).
   - Top-level `variables` (`workDir`, `adminUser`, `linuxUser`, `userHome`, `sshDir`, `keyPath`) so `"C:\\Users\\Administrator\\.ssh\\id_rsa"` is defined once and referenced via `"${keyPath}"` and `"${adminUser}"`.
   - Cross-platform `notes` explaining that `os: "windows"` does not restrict importing to Windows only.
   - Include `mainMachine` (`192.168.1.20`) and all 6 nodes (`w1`, `w2`, `w3`, `w4`, `u1`, `main`).
Acceptance Criteria:
- Zero repeated hardcoded `"C:\\Users\\Administrator\\.ssh\\id_rsa"` strings inside `nodes` or `connections` arrays across `04-w1-machine`..`07-final-network-machine`.
- Zero `snake_case` keys (`ip_address`, `auth_method`, `key_path`, `worker_id`) in those SSH node manifests.
Targeted Verification: pwsh -NoProfile -Command "Get-ChildItem D:\work\repo-secrets\0*-machine\gitmap-ssh*.json | ForEach-Object { Get-Content $_.FullName -Raw | ConvertFrom-Json | Out-Null }"
