# Antivirus Investigation & EDR Heuristics Mitigation Specification

## 1. Executive Summary & Problem Definition

In production and developer environments, Windows Defender, Microsoft Defender for Endpoint (EDR), and behavioral analysis engines have flagged the compiled `gitmap.exe` binary with high-severity threat classifications, including:
- `Trojan:Win32/AutoLogon`
- `Trojan:Win32/Wacatac.B!ml`
- `Behavior:Win32/TamperWinDefService`

This specification provides the exhaustive 4-part Root Cause Analysis (RCA), static signature analysis, behavioral heuristic breakdowns, and actionable code mitigations to permanently eradicate antivirus false-positive triggers from the GitMap codebase.

---

## 2. Comprehensive 4-Part Root Cause Analysis (RCA)

### 2.1 Part 1: Symptoms & Detection Signatures
- **Symptom 1 (Signature Match):** Windows Defender halts binary execution immediately upon compilation or download, isolating `gitmap.exe` into quarantine.
- **Symptom 2 (Behavioral Alert):** Executing administrative CLI subcommands triggers proactive blocking by heuristic detection engines.
- **Symptom 3 (Telemetry Alert):** Automated CI/CD virus scanner steps report high risk heuristics on unsigned Windows PE x86_64 artifacts.

---

### 2.2 Part 2: Detailed Root Causes (Code Audits)

#### Cause A: Plaintext Credentials in Windows Registry (Winlogon AutoAdminLogon)
- **Target File:** `cli/cmdos/os_autologin_windows.go` (lines 11–37)
- **Problematic Code Pattern:**
  ```go
  const winlogonRegPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`
  // ...
  k.SetStringValue("AutoAdminLogon", "1")
  k.SetStringValue("DefaultUserName", cfg.Username)
  k.SetStringValue("DefaultDomainName", cfg.Domain)
  k.SetStringValue("DefaultPassword", cfg.Password)
  ```
- **EDR / Defender Heuristic Mechanism:**
  Security scanners maintain strict byte-sequence and AST signatures for strings matching `Winlogon` coupled with `DefaultPassword` and `AutoAdminLogon`. Malware families (specifically credentials droppers, worms, and persistence kits) write plaintext passwords into `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon` to establish persistent unattended system access. Statically embedding these registry keys and functions inside an executable immediately flags the binary as `Trojan:Win32/AutoLogon`.

#### Cause B: EDR & Telemetry Service Tampering (`sc.exe` & `DiagTrack`)
- **Target File:** `cli/cmdos/os_tweak_windows_privacy.go` (lines 27–50)
- **Problematic Code Pattern:**
  ```go
  // In applyTelemetry and toggleDiagTrackService
  exec.Command("sc.exe", "config", "DiagTrack", "start=", "disabled").Run()
  exec.Command("sc.exe", "stop", "DiagTrack").Run()
  exec.Command("setx.exe", "POWERSHELL_TELEMETRY_OPTOUT", "1", "/M").Run()
  k.SetDWordValue("AllowTelemetry", 0)
  ```
- **EDR / Defender Heuristic Mechanism:**
  Disabling `DiagTrack` (Windows Connected User Experiences and Telemetry) via `sc.exe config ... disabled` and `sc.exe stop` is a primary indicator of defense evasion (MITRE ATT&CK T1562.001 - Impair Defenses: Disable or Modify Tools). Ransomware loaders systematically execute this exact command sequence prior to payload staging. AV engines treat any non-operating-system binary that invokes `sc.exe stop DiagTrack` as high-risk defense evasion malware.

#### Cause C: Dynamic Win32 Callback Trampolines & Focus Stealing
- **Target File:** `cli/cmdagy/agy_window_focus_windows.go` (lines 11–46)
- **Problematic Code Pattern:**
  ```go
  var (
      user32 = syscall.NewLazyDLL("user32.dll")
      procEnumWindows = user32.NewProc("EnumWindows")
      procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
  )
  // ...
  cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr { ... })
  procEnumWindows.Call(cb, 0)
  ```
- **EDR / Defender Heuristic Mechanism:**
  The Go runtime's `syscall.NewCallback` constructs an executable memory page (trampoline) at runtime to translate OS callback conventions into Go goroutines. Security engines monitoring PE imports and memory allocation observe:
  1. `user32.dll` loaded dynamically via `LoadLibrary`/`GetProcAddress`.
  2. Enumeration of all top-level windows (`EnumWindows`).
  3. Window hierarchy traversal to locate process IDs and forcibly bring windows to the foreground (`SetForegroundWindow`).
  In unsigned binaries, this behavioral pattern strongly correlates with spyware, screen monitors, and keyloggers (`Trojan:Win32/Wacatac.B!ml`).

#### Cause D: Unsigned PE Artifacts & Lack of SmartScreen Reputation
- **Mechanism:**
  Official GitMap releases and local development builds have historically been distributed as unsigned Portable Executable (PE) binaries. In modern Windows security architecture, unsigned executables with zero SmartScreen reputation that perform low-level registry and Win32 operations are subjected to aggressive machine-learning heuristic scoring. The combined presence of registry writes and dynamic Win32 trampolines causes the aggregate risk score to exceed the blocking threshold.

#### Cause E: Adjacent `macro-ahk` Automation Correlation
- **Mechanism:**
  In automated testing and desktop developer environments, GitMap operates alongside macro automation tools and AutoHotkey scripts (`macro-ahk` workflows). Antivirus process monitors analyze parent-child process chains and adjacent script execution. When an unsigned tool spawns or interacts with script-based keystroke/window injectors, behavioral heuristics elevate the threat score of both the binary and the parent terminal session.

---

### 2.3 Part 3: Architecture Risk Assessment & Impact
1. **Developer Friction:** Terminal installations blocked, requiring manual Defender exclusions.
2. **Distribution Risk:** Automated enterprise deployments rejected by corporate security policies.
3. **Product Integrity:** GitMap is a high-speed repository management companion; inclusion of OS-tweak utilities that mimic malware mechanisms directly compromises product integrity.

---

### 2.4 Part 4: Concrete Code Mitigations & Eradication Plan

#### Mitigation 1: Decouple & Deprecate Plaintext Winlogon Writes
- In `cli/cmdos/os_autologin_windows.go`:
  - Completely deprecate setting `DefaultPassword` in plaintext registry.
  - Implement a secure credential gate. If autologin is requested, require explicit opt-in confirmation (`--force-insecure-registry`), or rely on native Windows Credential Manager / DPAPI APIs.
  - Ensure the string constant `DefaultPassword` is not linked or statically exposed in standard build artifacts.
  - Enforce positive boolean validation: `hasUser`, `hasDomain`, `isValidConfig`.

#### Mitigation 2: Neutralize Telemetry Service Tampering
- In `cli/cmdos/os_tweak_windows_privacy.go`:
  - Cease executing `sc.exe config DiagTrack start= disabled` and `sc.exe stop DiagTrack` from general utility commands.
  - Scope all OS tweak actions to non-destructive registry tweaks under user control.
  - Remove administrative service manipulation commands that trigger MITRE ATT&CK T1562.001 defense evasion alerts.
  - Use positive boolean flags: `isEnabled`, `hasTelemetryActive`.

#### Mitigation 3: Refactor Window Focus Discovery
- In `cli/cmdagy/agy_window_focus_windows.go`:
  - Replace dynamic `syscall.NewCallback` window enumeration trampolines with direct process querying or safe window handle lookup when possible.
  - Ensure all Win32 interop cleanly handles process termination without lingering hooks.
  - Gate Antigravity focus management behind explicit user configuration (`hasAntigravityFocusEnabled`).

#### Mitigation 4: Build Sanitization & Authenticode Signing
- Add linker flags to production releases: `-ldflags="-s -w"` to strip debug symbols, file system paths, and internal package references from the final PE binary.
- Integrate automated code signing in CI/CD pipeline using standard Windows Authenticode certificates.
- Submit signed release binaries to Microsoft Security Intelligence for proactive whitelist verification prior to public tag release.

---

## 3. Implementation Verification Checklist

- [ ] Audit `cli/cmdos/os_autologin_windows.go` to ensure zero plaintext password writes into Winlogon registry.
- [ ] Audit `cli/cmdos/os_tweak_windows_privacy.go` to eliminate `sc.exe` service tampering invocations.
- [ ] Validate `cli/cmdagy/agy_window_focus_windows.go` for memory safety and positive boolean naming.
- [ ] Verify that GitMap binary compiles with clean symbols (`-ldflags="-s -w"`).
- [ ] Ensure all modified files adhere to strict relative paths and affirmative booleans (`isValid`, `isEnabled`, `hasTarget`).
