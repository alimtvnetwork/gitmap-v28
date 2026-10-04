# Subtask 213.2: Remote SSH Update to 2.19 Pipeline

- **Parent Plan:** [83-antigravity-ubuntu-update-and-macro-automation.md](../../83-antigravity-ubuntu-update-and-macro-automation.md)
- **Spec Reference:** [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md](../../../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md)
- **Status:** Ready
- **Target Node:** Ubuntu U1 (`192.168.1.22`, user `a`)
- **Target Area:** `/home/a/.local/share/antigravity-ide/`, `/usr/local/bin/antigravity`, remote SSH runner

---

## 1. Objective

Execute an automated, zero-downtime remote upgrade pipeline over SSH to transition Ubuntu workstation `u1` from Antigravity **2.13.0 (Build 6362815968182272)** to **2.19.1 (Build 6046815158665216)**, configure SUID root permissions for the Chromium sandbox binary, and verify both headless CLI execution and user-session graphical launch.

---

## 2. Remote SSH Execution Script

The upgrade procedure is executed via an embedded remote script piped through SSH to eliminate Windows CRLF newline corruption:

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "=== [Phase 1] Terminating Running Antigravity Processes ==="
pkill -f antigravity || true
pkill -f antigravity-ide || true
sleep 1

echo "=== [Phase 2] Downloading Antigravity 2.19.1 Tarball ==="
DOWNLOAD_URL="https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz"
TEMP_ARCHIVE="/tmp/Antigravity.tar.gz"
curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_ARCHIVE"

echo "=== [Phase 3] Backing up Existing Installation and Extracting ==="
INSTALL_DIR="/home/a/.local/share/antigravity-ide"
BACKUP_DIR="/home/a/.local/share/antigravity-ide.bak-2.13.0"

if [ -d "$INSTALL_DIR" ]; then
    rm -rf "$BACKUP_DIR"
    mv "$INSTALL_DIR" "$BACKUP_DIR"
fi

mkdir -p "$INSTALL_DIR"
tar -xzf "$TEMP_ARCHIVE" -C "$INSTALL_DIR" --strip-components=1

echo "=== [Phase 4] Hardening SUID Root Permissions on chrome-sandbox ==="
sudo chown root:root "$INSTALL_DIR/chrome-sandbox"
sudo chmod 4755 "$INSTALL_DIR/chrome-sandbox"

echo "=== [Phase 5] Preserving and Linking Global Binary ==="
sudo ln -sf "$INSTALL_DIR/antigravity" /usr/local/bin/antigravity

echo "=== [Phase 6] Cleaning Temporary Artifacts ==="
rm -f "$TEMP_ARCHIVE"

echo "=== [Phase 7] Verifying Version Output ==="
INSTALLED_VER=$(/usr/local/bin/antigravity --version || true)
echo "Reported Version: $INSTALLED_VER"

if echo "$INSTALLED_VER" | grep -q "2.19.1"; then
    echo "SUCCESS: Antigravity upgraded to 2.19.1 successfully!"
else
    echo "ERROR: Version check did not match 2.19.1!" >&2
    exit 1
fi
```

---

## 3. Remote Invocation Protocol from Windows Host

To run this pipeline from Windows without intermediate file generation on node `u1`:

```powershell
$UpgradeScript = @'
set -euo pipefail
pkill -f antigravity || true
pkill -f antigravity-ide || true
curl -fsSL "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz" -o /tmp/Antigravity.tar.gz
if [ -d /home/a/.local/share/antigravity-ide ]; then
    rm -rf /home/a/.local/share/antigravity-ide.bak-2.13.0
    mv /home/a/.local/share/antigravity-ide /home/a/.local/share/antigravity-ide.bak-2.13.0
fi
mkdir -p /home/a/.local/share/antigravity-ide
tar -xzf /tmp/Antigravity.tar.gz -C /home/a/.local/share/antigravity-ide --strip-components=1
sudo chown root:root /home/a/.local/share/antigravity-ide/chrome-sandbox
sudo chmod 4755 /home/a/.local/share/antigravity-ide/chrome-sandbox
sudo ln -sf /home/a/.local/share/antigravity-ide/antigravity /usr/local/bin/antigravity
rm -f /tmp/Antigravity.tar.gz
/usr/local/bin/antigravity --version
'@

$UpgradeScript | ssh.exe -o BatchMode=yes u1 "tr -d '\r' | bash -s"
```

---

## 4. GUI Launch Validation

After verifying the CLI output, validate graphical launch in the active desktop session:

```bash
ssh u1 "systemd-run --user /usr/local/bin/antigravity /home/a/git-work/gitmap"
```

Confirm that the application starts cleanly without throwing sandbox abort errors:
- Inspect systemd journal logs: `journalctl --user-unit run-*.service -n 50 --no-pager`
- Verify process table: `pgrep -fl antigravity`

---

## 5. Rollback Procedures

If extraction or validation fails:
1. Stop any newly launched processes: `ssh u1 "pkill -f antigravity || true"`
2. Swap back the backup directory:
   ```bash
   ssh u1 "rm -rf /home/a/.local/share/antigravity-ide && mv /home/a/.local/share/antigravity-ide.bak-2.13.0 /home/a/.local/share/antigravity-ide"
   ssh u1 "sudo chown root:root /home/a/.local/share/antigravity-ide/chrome-sandbox && sudo chmod 4755 /home/a/.local/share/antigravity-ide/chrome-sandbox"
   ```
3. Verify restored version: `ssh u1 "antigravity --version"`

---

## 6. Acceptance Criteria

- [ ] Remote SSH script exits with return code `0`.
- [ ] `/usr/local/bin/antigravity --version` outputs string containing `2.19.1`.
- [ ] Permissions on `chrome-sandbox` are `-rwsr-xr-x` (`4755`) and owner is `root:root`.
- [ ] Backup directory `/home/a/.local/share/antigravity-ide.bak-2.13.0` exists and contains previous version binaries.
- [ ] File `/tmp/Antigravity.tar.gz` is completely deleted.
