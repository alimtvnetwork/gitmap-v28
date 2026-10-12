# Specification: Scripts-Fixer & Antigravity Manager (AGM) 2.22.0 Parity

- **Spec ID:** `02-spec/21-app/root-cause-analysis-and-update-instructions/03-scripts-fixer-and-agm-parity.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Target Repositories:** `scripts-fixer`, `antigravity-manager`
- **Parent Task:** `root-cause-analysis-and-update-instructions`

---

## 1. Overview & Objective

To eliminate the 2.19.1 bottleneck across the developer ecosystem, companion tools must be synchronized to version `2.22.0` (Build ID `5446056071266304`) and equipped with dynamic manifest discovery.

---

## 2. Scripts-Fixer Toolchain Modifications

### 2.1 Target File: `scripts/os/ubuntu/install-antigravity.sh`
In `fetch_ide_url()`:
- Implement dynamic manifest check against `https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/latest-x64-linux.yml`.
- Update fallback static constants from `2.19.1-6046815158665216` to `2.22.0-5446056071266304`:
  ```bash
  fetch_ide_url() {
      local arch="$1"
      local manifest_url="https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/latest-x64-linux.yml"
      local live_info
      live_info=$(curl -fsSL --connect-timeout 3 "$manifest_url" 2>/dev/null || true)
      local live_ver
      live_ver=$(echo "$live_info" | grep -E '^version:' | awk '{print $2}' | tr -d '\r"' || true)
      local live_build
      live_build=$(echo "$live_info" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+-[0-9]+' | head -n1 | cut -d'-' -f2 || true)

      local ver="${live_ver:-2.22.0}"
      local bld="${live_build:-5446056071266304}"

      if [ "$arch" = "arm64" ] || [ "$arch" = "arm" ] || [ "$arch" = "aarch64" ]; then
          echo "https://storage.googleapis.com/antigravity-public/antigravity-hub/${ver}-${bld}/linux-arm/Antigravity.tar.gz"
          return
      fi
      echo "https://storage.googleapis.com/antigravity-public/antigravity-hub/${ver}-${bld}/linux-x64/Antigravity.tar.gz"
  }
  ```

### 2.2 Target File: `scripts-linux/69-install-antigravity/config.json`
Update default archive link:
```json
{
  "install": {
    "tarball": "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/linux-x64/Antigravity.tar.gz"
  },
  "notes": "Google Antigravity IDE and agy CLI companion installer"
}
```

---

## 3. Antigravity Manager (AGM) Modifications

### 3.1 Target File: `antigravity-manager/scripts/update-antigravity-ide.sh`
- Update lines 91-93 from `2.19.1-6046815158665216` to `2.22.0-5446056071266304`.
- Maintain backup directory creation (`~/.local/share/antigravity-ide.bak-<timestamp>`).
- Preserve launcher wrapper `antigravity-ide.run` and symlinks.

### 3.2 Target File: `cli/cmdinstall/agm_update.go`
- Ensure `gitmap update agm` and `gitmap agy update` cleanly distinguish between updating Antigravity Desktop IDE (`2.22.0`) and Antigravity Manager daemon (`4.89.x`).
