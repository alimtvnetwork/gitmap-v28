# Specification: Root Cause Analysis — Antigravity IDE Version Disparity & Linux In-App Auto-Update Inoperability

- **Spec ID:** `02-spec/21-app/root-cause-analysis-and-update-instructions/01-root-cause-analysis.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Target Subsystem:** Desktop IDE Fleet Management, Antigravity Update Engine, Linux Deployment Runtime
- **Parent Task:** `root-cause-analysis-and-update-instructions`

---

## 1. Executive Summary & Problem Formulation

Google Antigravity Desktop IDE exhibits a persistent cross-platform version divergence:
- **Windows Workstations:** Running version `2.22.0` (Build ID `5446056071266304`), successfully updating via the in-app "Check for Updates" mechanism.
- **Ubuntu Linux Workstations:** Stranded on version `2.19.1` (Build ID `6046815158665216`), where clicking "Check for Updates" is completely inert and fails to download or apply any updates.

![Antigravity Version 2.22.0 on Windows](assets/screenshots/root-cause-analysis-and-update-instructions-01.png)

*Figure 1: Verified Windows environment displaying Version 2.22.0 with functional update backend.*

This document details the multi-layer root cause analysis explaining:
1. Why the in-app "Check for Updates" button fails on Linux.
2. How Windows automatically achieved version 2.22.0.
3. Google's actual Linux distribution infrastructure (GCS archives vs package managers).
4. Why automation toolchains (`gitmap`, `scripts-fixer`, `antigravity-manager`) remained pinned to 2.19.1.
5. The definitive architectural solution for automated, zero-downtime updates on Linux.

---

## 2. Technical Deep Dive: Why In-App "Check for Updates" Fails on Linux

### 2.1 Electron Auto-Updater Runtime Selection

Antigravity Desktop IDE is an Electron application bundling `electron-updater` (v6.8.9) configured with generic HTTP update provider targeting Google Cloud Run:
```yaml
provider: generic
url: https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/
updaterCacheDirName: antigravity-updater
```

In `resources/app.asar/node_modules/electron-updater/out/main.js`, the runtime updater class is instantiated conditionally by platform:
```javascript
function doLoadAutoUpdater() {
    if (process.platform === "win32") {
        return new (require("./NsisUpdater").NsisUpdater)();
    } else if (process.platform === "darwin") {
        return new (require("./MacUpdater").MacUpdater)();
    } else {
        _autoUpdater = new (require("./AppImageUpdater").AppImageUpdater)();
        try {
            const identity = path.join(process.resourcesPath, "package-type");
            if (!existsSync(identity)) {
                return _autoUpdater;
            }
            const fileType = readFileSync(identity).toString().trim();
            switch (fileType) {
                case "deb": return new (require("./DebUpdater").DebUpdater)();
                case "rpm": return new (require("./RpmUpdater").RpmUpdater)();
            }
        } catch (error) {}
        return _autoUpdater;
    }
}
```

### 2.2 The Silent Abort in `AppImageUpdater`

On Linux workstations, GitMap and `scripts-fixer` deploy Antigravity as an unpacked `.tar.gz` directory tree into `~/.local/share/antigravity-ide/`.
1. The file `resources/package-type` does not exist in the official tarball.
2. Electron falls back to instantiating `AppImageUpdater`.
3. In `AppImageUpdater.js`, the readiness check strictly inspects environment variables:
   ```javascript
   isUpdaterActive() {
       if (process.env["APPIMAGE"] == null && !this.forceDevUpdateConfig) {
           this._logger.warn("APPIMAGE env is not defined, current application is not an AppImage");
           return false;
       }
       return super.isUpdaterActive();
   }
   ```
4. Because the application was launched from `antigravity-ide.run` rather than an AppImage mount, `process.env["APPIMAGE"]` is `undefined`.
5. `isUpdaterActive()` returns `false`.
6. When the user clicks **Check for Updates**, `checkForUpdates()` executes:
   ```javascript
   checkForUpdates() {
       if (!this.isUpdaterActive()) {
           return Promise.resolve(null);
       }
       ...
   }
   ```
7. It exits immediately with `null` without making any network requests or emitting IPC events (`update-available`, `update-not-available`, `error`).
8. The UI receives zero feedback and remains completely silent.

### 2.3 Structural Inability of Unpacked Directories to Self-Update

Even if `isUpdaterActive()` was overridden:
- `AppImageUpdater` only supports replacing a single file by moving a new `.AppImage` over `process.env["APPIMAGE"]`.
- It possesses no extraction engine to untar runtime `.so` libraries, Electron helpers, and assets into an active folder while running.
- Therefore, standalone tarball installations on Linux **cannot natively self-update from within Electron**. External orchestration is mandatory.

---

## 3. Why Windows Successfully Updates to 2.22.0

On Windows:
1. `electron-updater` loads `NsisUpdater`.
2. `NsisUpdater.isUpdaterActive()` evaluates to `true` (only requiring `this.app.isPackaged == true`).
3. The updater queries `https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/latest-x64-win.yml`.
4. It discovers release `2.22.0` (Build ID `5446056071266304`).
5. It streams `Antigravity-x64.exe` to `%LOCALAPPDATA%\antigravity-updater\pending\`.
6. Upon app restart, the NSIS installer runs silently (`/S`) and updates files in `%LOCALAPPDATA%\Programs\Antigravity`.
7. This automated flow is built into Windows Electron distributions but does not exist for unpacked Linux archives.

---

## 4. Google Linux Distribution Model: Package Managers vs Tarballs

### 4.1 Absence of Official APT / Debian Repositories
Rigorous discovery across Google infrastructure, Cloud Run updater services, and public Google Cloud Storage buckets confirms:
- **Google does NOT maintain an APT repository** (`deb [arch=amd64] https://...` or PPA) for Antigravity IDE.
- There is no `apt-get install antigravity` or automated `apt update` channel available from Google.

### 4.2 Verified Google Cloud Storage Artifacts
Google distributes releases directly via Cloud Storage under:
`https://storage.googleapis.com/antigravity-public/antigravity-hub/${VERSION}-${BUILD_ID}/${PLATFORM}/${ARTIFACT}`

For Version `2.22.0` (Build ID `5446056071266304`), the verified production URLs on GCS are:

| Platform | Format | Artifact URL | Verified Size |
| :--- | :--- | :--- | :--- |
| **Linux x64** | Tarball | `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/linux-x64/Antigravity.tar.gz` | 180,243,767 B |
| **Linux x64** | Debian Package | `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/linux-x64/Antigravity.deb` | 150,451,896 B |
| **Linux x64** | AppImage | `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/linux-x64/Antigravity.AppImage` | 189,298,175 B |
| **Linux ARM64** | Tarball | `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/linux-arm/Antigravity.tar.gz` | 178,116,821 B |
| **Windows x64** | NSIS Exe | `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/windows-x64/Antigravity-x64.exe` | 164,082,024 B |

### 4.3 Evaluation of Installation Strategies

| Method | User Elevation Required | In-App Auto-Update Works? | External Tooling Support | AppArmor Sandbox Compatibility | Recommendation |
| :--- | :---: | :---: | :---: | :---: | :--- |
| **Tarball (`.tar.gz`)** | No (`~/.local/`) | No (silent abort) | High (`gitmap`, `scripts-fixer`) | Perfect via `--no-sandbox` wrapper | **Primary Standard** |
| **Debian (`.deb`)** | Yes (`sudo dpkg`) | Partial (requires helper) | Low (prompts for root password) | Requires setuid `chrome-sandbox` | Secondary / System Admin |
| **AppImage** | No | Yes (requires FUSE) | Medium | Requires FUSE library installed | Not universally installed on minimal Ubuntu |
| **Snapcraft** | Yes (`sudo snap`) | Handled by `snapd` | Independent | Confined sandbox | Heavy startup latency |

**Conclusion:** The unpacked user tarball in `~/.local/share/antigravity-ide` remains the ideal installation architecture for non-root developer environments, but **updates must be driven externally** by `gitmap update ide`, `gitmap agy update`, or `scripts-fixer`.

---

## 5. Root Cause of Version Pinning in Tooling

Both `gitmap` and `scripts-fixer` hardcoded static version numbers:
1. In `cli/cmdinstall/installantigravity_types.go`:
   ```go
   AntigravityDefaultVersion = "2.19.1"
   AntigravityDefaultBuildID = "6046815158665216"
   ```
2. In `scripts/os/ubuntu/install-antigravity.sh`:
   ```bash
   echo "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz"
   ```
Neither tool implemented dynamic resolution against Google's Cloud Run update manifest (`https://antigravity-hub-auto-updater-974169037036.us-central1.run.app/manifest/latest-x64-linux.yml`). Consequently, running `gitmap agy update` continuously redownloaded version 2.19.1.
