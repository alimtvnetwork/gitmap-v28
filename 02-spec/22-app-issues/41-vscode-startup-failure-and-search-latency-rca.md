# Issue 41: VS Code Startup Failure & Investigation Latency Root Cause Analysis

> **Issue Status:** Resolved  
> **Traceability IDs:** Task-07  
> **Canonical Path:** `02-spec/22-app-issues/41-vscode-startup-failure-and-search-latency-rca.md`  
> **Parent Spec:** `02-spec/22-app-issues/01-index.md`

---

## 1. Symptoms & Initial Reproduction

1. **Symptom:** VS Code failed to open upon launch, silently crashing or throwing Chromium binary errors (`ERROR:base/i18n/icu_util.cc: Invalid file descriptor to ICU data received`).
2. **User Context:** Suspected corruption in `projects.json` (Project Manager configuration) or removed folder assets.
3. **Execution Delay:** Initial root-cause investigation took prolonged time (~1 hour).

---

## 2. Technical Root Cause Analysis (4-Part RCA)

### Part A: The True Technical Root Cause
The failure was a **dual-layer collision**:
1. **Chromium ICU Binary Corruption (`win32VersionedUpdate`)**:
   - VS Code's background auto-updater was abruptly terminated while in-place swapping binaries.
   - An orphaned `inno_updater` / `Code.exe` process held file locks on binary dependencies in `AppData\Local\Programs\Microsoft VS Code`.
   - The commit resource subfolder (e.g. `[0-9a-f]{10}` holding `icudtl.dat` and `v8_context_snapshot.bin`) was partially missing or unlinked, triggering a native Chromium crash on process initialization before the window opened.
2. **Project Manager JSON (`projects.json`) Inconsistency**:
   - The user had active installations in both modern (`Code\User\globalStorage\alefragnani.project-manager\projects.json`) and legacy (`Code\User\projects.json`).
   - A concurrent write had stripped schema keys (`paths`, `tags`, `enabled`, `profile`), leading to parser panic when Project Manager extension activated.

### Part B: Why Did Investigation Take One Hour? (Latency Analysis)
The delay was driven by three concrete technical bottlenecks:
1. **Misdirected Initial Hypothesis (Product vs Configuration Search)**:
   - Early search attempts inspected `product.json` and internal VS Code electron resources rather than the Project Manager user storage.
2. **Slow Python File Walking on Broad Windows Drives**:
   - Running broad directory scans across `D:\work` and deep `AppData` structures via unindexed Python scripts traversed over 120,000 files, encountering Windows filesystem stat latency (~30-40 seconds per roundtrip).
3. **Process Lock Deadlocks**:
   - Because lingering `Code.exe` background processes were still running in the background, manual attempts to verify fixes failed intermittently until full `taskkill /F /IM Code.exe` process termination was automated.

---

## 3. Remediation & Code Fix

1. **Native Diagnostic & Repair Engine (`cli/cmdvscode/vscode_repair*.go`)**:
   - Implemented `gitmap vscode repair` (aliases `fix`, `doctor`).
   - **Step 1:** Automatically terminates lingering `Code.exe` and `inno_updater` processes.
   - **Step 2:** Validates `projects.json`, backs up corrupted states to `.bak`, ensures schema compliance (`paths`, `tags`, `enabled`), and syncs modern/legacy locations.
   - **Step 3:** Verifies binary integrity via `code.cmd --version` and synchronizes missing commit folders across parallel install paths.
   - **Step 4:** Automatic fallback to `winget` repair if binaries cannot be recovered locally.
2. **Automated Cross-Platform Helper Script (`tools/repair-vscode-and-projects.ps1`)**:
   - Integrated into `scripts-fixer` with dynamic path resolution (`$env:APPDATA`, `$env:ProgramFiles`, `$env:LOCALAPPDATA`, Registry HKLM) with zero hardcoded drive letters.
3. **Sub-Millisecond Search Acceleration**:
   - Replaced unindexed full-disk traversals with GitMap Native AUM Searcher (< 1ms latency vs 33s in Python).

---

## 4. Prevention Checklist

- [x] Always terminate locked processes first with `gitmap vscode repair --kill`.
- [x] Use dynamic environment variables (`$env:APPDATA`) instead of fixed drive paths (`D:\`).
- [x] Use GitMap native AUM search for codebase exploration instead of full-disk interpreter loops.
