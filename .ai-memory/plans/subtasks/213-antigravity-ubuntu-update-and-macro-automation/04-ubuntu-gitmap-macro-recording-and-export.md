# Subtask 213.4: Ubuntu GitMap Macro Recording, Execution & Export

- **Parent Plan:** [83-antigravity-ubuntu-update-and-macro-automation.md](../../83-antigravity-ubuntu-update-and-macro-automation.md)
- **Spec Reference:** [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md](../../../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md)
- **Status:** Ready
- **Target Node:** Ubuntu U1 (`192.168.1.22`, user `a`)
- **Target Files:**
  - Ubuntu Macro Store: `/home/a/.gitmap/macros/update-antigravity.json`
  - Windows Fleet Staging: `d:/work/repo-secrets/04-ubuntu-migration/update-antigravity.json`

---

## 1. Objective

Author, record, and deploy the GitMap macro `update-antigravity.json` on Ubuntu workstation `u1`. Validate its automated execution via `gitmap macro run update-antigravity --verbose`, and establish bi-directional fleet portability by exporting the macro to Windows staging directory `d:/work/repo-secrets/04-ubuntu-migration/` and importing it into the Windows GitMap macro catalog.

---

## 2. Macro Step Breakdown & JSON Payload

The macro encapsulates 7 deterministic steps designed to execute cleanly under unprivileged user `a` with non-interactive `sudo` privilege escalation for system-level actions:

| Step | Command | Working Dir | Timeout | Continue on Error | Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **1** | `pkill -f antigravity \|\| true` | `/home/a` | 15s | `true` | Terminate active IDE and language server processes |
| **2** | `curl -fsSL https://storage.googleapis.com/... -o /tmp/Antigravity.tar.gz` | `/tmp` | 300s | `false` | Fetch official 2.19.1 Linux x64 tarball |
| **3** | `if [ -d ...-ide ]; then rm -rf ...bak-2.13.0 && mv ... ...bak-2.13.0; fi && mkdir -p ... && tar -xzf ...` | `/home/a` | 120s | `false` | Backup old version and unpack new release cleanly |
| **4** | `echo a \| sudo -S chown root:root .../chrome-sandbox && echo a \| sudo -S chmod 4755 .../chrome-sandbox` | `/home/a/.local/share/antigravity-ide` | 30s | `false` | Grant setuid root permissions to Chromium sandbox |
| **5** | `echo a \| sudo -S ln -sf .../antigravity /usr/local/bin/antigravity` | `/usr/local/bin` | 15s | `false` | Reconcile global binary symlink in `/usr/local/bin` |
| **6** | `rm -f /tmp/Antigravity.tar.gz` | `/tmp` | 15s | `true` | Delete ephemeral download archive |
| **7** | `/usr/local/bin/antigravity --version` | `/home/a` | 30s | `false` | Assert upgraded binary reports `2.19.1` |

### Exact Macro File Content (`update-antigravity.json`)
```json
{
  "id": 1,
  "name": "update-antigravity",
  "description": "Automated update pipeline for Antigravity 2.19.1 on Ubuntu with SUID sandbox hardening",
  "created_at": "2026-10-04T12:00:00Z",
  "updated_at": "2026-10-04T12:00:00Z",
  "total_steps": 7,
  "tags": "antigravity,update,fleet,ubuntu,sandbox",
  "steps": [
    {
      "id": 1,
      "macro_id": 1,
      "step_num": 1,
      "command_line": "pkill -f antigravity || true",
      "working_dir": "/home/a",
      "continue_on_error": true,
      "timeout_seconds": 15
    },
    {
      "id": 2,
      "macro_id": 1,
      "step_num": 2,
      "command_line": "curl -fsSL https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz -o /tmp/Antigravity.tar.gz",
      "working_dir": "/tmp",
      "continue_on_error": false,
      "timeout_seconds": 300
    },
    {
      "id": 3,
      "macro_id": 1,
      "step_num": 3,
      "command_line": "if [ -d /home/a/.local/share/antigravity-ide ]; then rm -rf /home/a/.local/share/antigravity-ide.bak-2.13.0 && mv /home/a/.local/share/antigravity-ide /home/a/.local/share/antigravity-ide.bak-2.13.0; fi && mkdir -p /home/a/.local/share/antigravity-ide && tar -xzf /tmp/Antigravity.tar.gz -C /home/a/.local/share/antigravity-ide --strip-components=1",
      "working_dir": "/home/a",
      "continue_on_error": false,
      "timeout_seconds": 120
    },
    {
      "id": 4,
      "macro_id": 1,
      "step_num": 4,
      "command_line": "echo a | sudo -S chown root:root /home/a/.local/share/antigravity-ide/chrome-sandbox && echo a | sudo -S chmod 4755 /home/a/.local/share/antigravity-ide/chrome-sandbox",
      "working_dir": "/home/a/.local/share/antigravity-ide",
      "continue_on_error": false,
      "timeout_seconds": 30
    },
    {
      "id": 5,
      "macro_id": 1,
      "step_num": 5,
      "command_line": "echo a | sudo -S ln -sf /home/a/.local/share/antigravity-ide/antigravity /usr/local/bin/antigravity",
      "working_dir": "/usr/local/bin",
      "continue_on_error": false,
      "timeout_seconds": 15
    },
    {
      "id": 6,
      "macro_id": 1,
      "step_num": 6,
      "command_line": "rm -f /tmp/Antigravity.tar.gz",
      "working_dir": "/tmp",
      "continue_on_error": true,
      "timeout_seconds": 15
    },
    {
      "id": 7,
      "macro_id": 1,
      "step_num": 7,
      "command_line": "/usr/local/bin/antigravity --version",
      "working_dir": "/home/a",
      "continue_on_error": false,
      "timeout_seconds": 30
    }
  ]
}
```

---

## 3. Remote Staging & Live Macro Execution

### 3.1 Provisioning the Macro on Node U1
From Windows, stream the JSON directly into `/home/a/.gitmap/macros/update-antigravity.json`:

```powershell
$MacroJSON = Get-Content -Raw "d:\work\repo-secrets\04-ubuntu-migration\update-antigravity.json"
$BashProvision = @"
mkdir -p /home/a/.gitmap/macros
cat << 'EOF' > /home/a/.gitmap/macros/update-antigravity.json
$MacroJSON
EOF
chmod 644 /home/a/.gitmap/macros/update-antigravity.json
"@

$BashProvision | ssh.exe -o BatchMode=yes u1 "tr -d '\r' | bash -s"
```

### 3.2 Executing the Macro on Node U1
Execute the macro remotely via the GitMap CLI on U1:
```bash
ssh u1 "gitmap macro run update-antigravity --verbose"
```
The GitMap macro runner prints the step execution tree, evaluates each command sequentially, records task audit telemetry in `installation.db`, and outputs version confirmation.

---

## 4. Bi-Directional Export & Import Synchronization

To guarantee version-controlled persistence in git staging:

1. **Export from Ubuntu U1:**
   ```bash
   ssh u1 "gitmap macro export update-antigravity --out /tmp/update-antigravity.json"
   ```
2. **Transfer to Windows Fleet Vault:**
   ```powershell
   scp.exe u1:/tmp/update-antigravity.json d:\work\repo-secrets\04-ubuntu-migration\update-antigravity.json
   ```
3. **Import into Windows GitMap Macro Engine:**
   ```powershell
   gitmap macro import "d:\work\repo-secrets\04-ubuntu-migration\update-antigravity.json"
   ```
4. **Inspect Windows Catalog:**
   ```powershell
   gitmap macro list
   ```

---

## 5. Rollback Procedures

If any step in the macro fails:
1. Revert to backup directory:
   ```bash
   ssh u1 "if [ -d /home/a/.local/share/antigravity-ide.bak-2.13.0 ]; then rm -rf /home/a/.local/share/antigravity-ide && mv /home/a/.local/share/antigravity-ide.bak-2.13.0 /home/a/.local/share/antigravity-ide && echo a | sudo -S chown root:root /home/a/.local/share/antigravity-ide/chrome-sandbox && echo a | sudo -S chmod 4755 /home/a/.local/share/antigravity-ide/chrome-sandbox; fi"
   ```
2. Verify restored binary reports `2.13.0`:
   ```bash
   ssh u1 "/usr/local/bin/antigravity --version"
   ```

---

## 6. Acceptance Criteria

- [ ] `/home/a/.gitmap/macros/update-antigravity.json` exists on node `u1` and contains valid JSON with 7 steps.
- [ ] `gitmap macro list` on `u1` lists `update-antigravity` with 7 steps.
- [ ] `gitmap macro run update-antigravity` succeeds with exit code `0`.
- [ ] Exported JSON file is staged at `d:/work/repo-secrets/04-ubuntu-migration/update-antigravity.json`.
- [ ] `gitmap macro import` successfully loads the macro on the Windows workstation.
- [ ] Remote `chrome-sandbox` has ownership `root:root` and mode `4755`.
