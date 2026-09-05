# 06 — Macro Step Execution Engine & Cross-Platform Shell Open Behavior

**Status:** Resolved & Implemented  
**Date:** 2026-09-05  
**Category:** CLI / Macro Execution Engine / Cross-Platform Shims  
**Related Plan:** [`.lovable/plans/completed/64-macro-step-open-chrome-failure.md`](../../plans/completed/64-macro-step-open-chrome-failure.md)  

---

## 1. Verbatim Session Event & User Directives

During an interactive macro recording session on Windows PowerShell, the user added a 9-step macro named `"alim"`:

```text
PS M:\Work-files\export\export\chrome-ext> gitmap macro add "alim"

  ● Interactive Macro Builder: "alim"
  Enter commands one per line (empty line or 'done' to save, 'cancel' to abort):
  (Tip: type 'rec' or 'record' to launch live command recording session)

  Step 1> cd %temp%
  Step 2> echo PWD
  Step 3> gitmap cd macro
  Step 4> gitmap pull
  Step 5> cd d:\work
  Step 6> gitmap mkdir -p "sample"
  Step 7> gitmap rm "sample"
  Step 8> open chrome
  Step 9> open "linkedin.com"
  Step 10> done

✓ Created macro "alim" (9 step(s))
   1. cd %temp%
   2. echo PWD
   3. gitmap cd macro
   4. gitmap pull
   5. cd d:\work
   6. gitmap mkdir -p "sample"
   7. gitmap rm "sample"
   8. open chrome
   9. open "linkedin.com"

Run with: gitmap macro run alim (or gitmap alim)
```

The user then invoked the macro via `gitmap alim`:

```text
PS M:\Work-files\export\export\chrome-ext> gitmap alim
  alim
  ├── cd %temp%
  ├── echo PWD
  ├── gitmap cd macro
  ├── gitmap pull
  ├── cd d:\work
  ├── gitmap mkdir -p "sample"
  ├── gitmap rm "sample"
  ├── open chrome
  └── open "linkedin.com"

  ▶ Executing Macro: "alim" (9 steps)

  [ 1/9] ➜ cd C:\Users\Alim\AppData\Local\Temp ...   ➜ 📁 Directory: C:\Users\Alim\AppData\Local\Temp
✔ ok (0.0s)
  [ 2/9] ➜ echo PWD ... ✔ ok (0.2s)
  [ 3/9] ➜ gitmap cd macro ... ✔ ok (0.3s)
  [ 4/9] ➜ gitmap pull ... ✔ ok (11.7s)
  [ 5/9] ➜ cd d:\work ...   ➜ 📁 Directory: d:\work
✔ ok (0.0s)
  [ 6/9] ➜ gitmap mkdir -p "sample" ... ✔ ok (0.3s)
  [ 7/9] ➜ gitmap rm "sample" ... ✔ ok (0.2s)
  [ 8/9] ➜ open chrome ... ✖ failed (0.7s)
  ✖ Step 8 failed: exit status 1
gitmap: [E9000:EXECUTION] macro.Execute:
PS M:\Work-files\export\export\chrome-ext>
```

### User Directive:
> *"Write this as pending task I will explain later"*

---

## 2. Technical Audit & Root Cause Analysis

### A. The Windows `open` Dilemma
1. In macOS (`darwin`), `/usr/bin/open` is a standard system utility that opens files, URLs, and applications (e.g., `open -a "Google Chrome"`, `open "https://example.com"`).
2. On Windows (`windows`), `open` does not exist as an executable binary in `%PATH%`. Opening files or URLs is typically achieved via:
   - `explorer.exe <target>`
   - `cmd.exe /c start "" <target>`
   - PowerShell `Start-Process <target>`
   - ShellExecute Win32 API (`rundll32.exe url.dll,FileProtocolHandler <url>`)
3. When `gitmap macro` executes a line that begins with `open ...`:
   - It delegates execution to the system shell or an `exec.Command("open", ...)`.
   - On Windows, this lookup fails immediately with `exit status 1` (`'open' is not recognized as an internal or external command, operable program or batch file.`), returning error code `E9000:EXECUTION`.

### B. Existing GitMap Capabilities
GitMap already has built-in URL, browser, and path opening implementations:
- `gitmap open <path|url>`
- `gitmap chrome open [profile] [url]`
- `browser.OpenURL(url)` / `exec.Command("rundll32", "url.dll,FileProtocolHandler", url)`

### C. Architectural Options for Macro Execution Engine
When the user provides their explanation, the macro execution engine in `gitmap/cmd/` can be extended with one of the following architectural patterns:
1. **Built-in `open` Command Interceptor / Shim:**
   If a macro step starts with `open <args>`:
   - On Windows: parse the target. If target is a known app (e.g. `chrome`), launch `chrome.exe` or `gitmap chrome open`. If target is a URL or file path, invoke Windows `explorer.exe` or `rundll32.exe url.dll,FileProtocolHandler`.
   - On Linux: route to `xdg-open <args>`.
   - On macOS: route to `/usr/bin/open <args>`.
2. **Implicit Delegation to `gitmap open`:**
   Expand bare `open <args>` into `gitmap open <args>`.
3. **Configurable Macro Aliases:**
   Allow users to define platform-specific aliases within macro definitions.

---

## 3. Chrome Profile Picker & Reconciliation Milestones Completed This Turn

Prior to the macro session, the Chrome profile picker visibility desynchronization was completely diagnosed, specified, implemented, tested, and pushed:

1. **Specification:**
   - Authored formal specification under `spec/25-chrome-profile-management/` (`00-overview.md`, `01-profile-registration-and-picker.md`, `97-acceptance-criteria.md`, `99-consistency-report.md`).
2. **Local State 13-Attribute Registration Schema:**
   - Implemented `registerChromeProfileWithFullSchema` and `applyChromeInfoEntryDefaults` in `gitmap/cmd/chromeprofile_register.go`.
   - Populated all required Chromium UI properties (`name`, `shortcut_name`, `user_name`, `avatar_icon`, `default_avatar_fill_color`, `default_avatar_stroke_color`, `profile_highlight_color`, `profile_color_seed`, `active_time`, `is_using_default_avatar`, `is_using_default_name`, `is_ephemeral`, `is_consented_primary_account`, `signin.with_credential_provider`).
3. **Process Concurrency Protection:**
   - Added `isChromeRunning()` check and user advisory warning that Chrome overwrites `Local State` from memory on exit.
4. **Preferences Sanitization:**
   - Stripped parental lock / `managed` keys while preserving user email identity in `patchImportedChromeProfilePreferences`.
5. **Orphan Reconcile Command:**
   - Added `gitmap chrome profile reconcile` (`gitmap/cmd/chromeprofile_reconcile.go`) to scan disk and re-register unlinked profiles.
6. **Commit & Push:**
   - Committed as `9f7f24ce` (`feat(chrome): add profile picker registration schema and reconcile engine`) and pushed to `main`.

---

## 4. Resolution & Implementation

1. **Root Cause:**
   On Windows, PowerShell commands invoked via `powershell -NoProfile -Command "open ..."` fail because `open` is not a native Windows binary or PowerShell cmdlet.
2. **Implementation:**
   - Authored `gitmap/macro/open.go` with `ParseOpenCommand` and `executeOpenStep`.
   - Built platform openers without nested ifs:
     - Windows: `launchChromeWindows` (checks App Paths / Program Files or `Start-Process chrome`), `launchURLWindows` (`rundll32 url.dll,FileProtocolHandler` with `Start-Process` fallback), `explorer.exe` for local paths.
     - macOS: `/usr/bin/open`.
     - Linux: `xdg-open` or `google-chrome`.
   - Intercepted `open` in `gitmap/macro/execute.go` within `executeSingleStep` (adjacent to `DirTracker.ProcessCd`).
3. **Verification:**
   - Comprehensive unit test suite in `gitmap/macro/open_test.go` verifying parsing, target normalization, and mock command dispatch.
   - All 16 CI/CD quality gates pass 100% green.

