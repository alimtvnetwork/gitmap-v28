# Root Cause Analysis (RCA-19): VS Code Startup Failure (ICU Data Error) and Project Manager JSON Audit

## 1. Symptom
Launching VS Code (via `code`, `code --version`, or double-clicking the application icon) failed immediately with exit code -2147483645 (`0x80000003` / `STATUS_BREAKPOINT`) and logged:
```text
[ERROR:base\i18n\icu_util.cc:232] Invalid file descriptor to ICU data received.
```
VS Code workbench failed to open completely.

## 2. Root Cause
1. **VS Code Startup Failure:**
   An interrupted/aborted background automatic update of VS Code (from build commit `7debcd0e2a` v1.138.0 to build commit `2242ebbb54` v1.139.0) occurred while `Code.exe` was running and locked.
   - The installer extracted the versioned resource folder `2242ebbb54` (per VS Code's `win32VersionedUpdate` architecture) but could not overwrite `Code.exe` in the root.
   - When launched, `Code.exe` (built for commit `7debcd0e2a`) looked for its versioned resource assets and ICU data (`icudtl.dat`) in folder `7debcd0e2a`. Because that folder was removed/replaced with `2242ebbb54`, Chromium's `LazyInitIcuDataFile()` failed to obtain a valid file handle, causing `LoadIcuData()` to fail with `Invalid file descriptor to ICU data received` and trigger a fatal `CHECK(result)` abort.
2. **Project Manager JSON (`projects.json`):**
   The Project Manager JSON (`projects.json`) in `%APPDATA%\Code\User\globalStorage\alefragnani.project-manager\projects.json` was intact and valid JSON (62 registered projects), but lacked mirroring to the legacy `%APPDATA%\Code\User\projects.json` path checked by older versions of the extension and external scripts.

## 3. Resolution
1. **Restored Version Asset Tree & Binary Alignment:**
   - Synchronized the matching commit resource directories (`7debcd0e2a` and `2242ebbb54`) into `C:\Program Files\Microsoft VS Code` and `%LOCALAPPDATA%\Programs\Microsoft VS Code`.
   - Verified `code --version` outputs:
     ```text
     1.139.0
     2242ebbb54efeeb0129e08e919e7e8d43033cd83
     x64
     ```
   - Confirmed VS Code launches cleanly without error.
2. **Created Multi-Machine Repair PowerShell Script:**
   - Authored [repair-vscode-and-projects.ps1](tools/repair-vscode-and-projects.ps1) in `tools/`.
   - The script automates:
     - Process termination of stuck `Code.exe` / `inno_updater.exe` / `code-tunnel` instances.
     - Deep inspection and validation of `projects.json` (checking syntax, missing keys, stripped fields, non-existent directories).
     - Automated backup of any corrupt JSON before resetting/rebuilding.
     - Dual-sync between modern `globalStorage` and legacy `User` paths.
     - Detection of versioned update folder mismatches and automatic repair.
     - Non-destructive `winget` re-application if binaries remain corrupt.
     - GitMap Project Manager integration validation (`gitmap vscode find-duplicates`).

## 4. Prevention & Learnings
- When diagnosing VS Code launch crashes on Windows with `icu_util.cc:232`, check the commit hash of `Code.exe` against the subfolder names in the installation directory. If the commit folder matching `Code.exe` is absent, an incomplete background update occurred due to a locked file.
- Keep both modern and legacy `projects.json` locations mirrored to prevent extension version incompatibilities across different machines.
