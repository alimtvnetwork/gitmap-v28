# Subtask 01: Machine OS Version Detection Implementation Plan

> **Subtask ID:** `Subtask-01`  
> **Parent Plan:** `.ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md`  
> **Target Path:** `.ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/01-os-version-detection-spec.md`  
> **Status:** `PENDING`  
> **Assigned Subsystem:** `cli/cmdos/`  
> **Related Architecture Spec:** `02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/01-architecture-spec.md`  
> **Owned Implementation Files:**  
> - `cli/cmdos/os_info_types.go`  
> - `cli/cmdos/os_info_windows.go`  
> - `cli/cmdos/os_info_other.go`  
> - `cli/cmdos/os_machine_alias.go`  
> - `cli/cmdos/os_info_test.go`  

---

## 1. Objectives & Scope

1. **Structured `OSVersionInfo` Data Type:**
   Replace raw, unparsed OS string heuristics with a first-class `OSVersionInfo` struct in `cli/cmdos/os_info_types.go`, capturing operating system category, affirmative booleans (`IsWindowsServer`, `IsWindowsPro`, `IsWindowsHome`), version breakdown (`Major`, `Minor`, `Build`), marketing release, edition, and kernel version.
2. **Deep Windows Registry Probing (`cli/cmdos/os_info_windows.go`):**
   Implement native registry queries against `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion` extracting `InstallationType`, `EditionID`, `ProductName`, `DisplayVersion`, `CurrentBuildNumber`, and `UBR` (Update Build Revision). Incorporate build thresholding (`Build >= 22000`) to accurately differentiate Windows 11 from Windows 10.
3. **Cross-Platform Parity for Unix & Darwin (`cli/cmdos/os_info_other.go`):**
   Parse Linux `/etc/os-release` and macOS `sw_vers` + `uname -r` into structured `OSVersionInfo` fields with safe fallbacks and zero panics.
4. **Dual Output Rendering Pipeline (`cli/cmdos/os_machine_alias.go`):**
   Upgrade `gitmap machine` to render a clean, colorized Unicode tree hierarchy (`├──`, `└──`) with edition badges (`[Server]`, `[Pro]`, `[Home]`), and update `gitmap machine --json` to include the nested `osVersionDetails` object while maintaining 100% backward compatibility for all existing top-level fields.
5. **Comprehensive Unit Testing (`cli/cmdos/os_info_test.go`):**
   Add unit tests covering registry value parsing, Windows 11 normalization, Linux os-release parsing, affirmative boolean assignment, and JSON serialization.

---

## 2. Detailed Technical Deliverables

### 2.1 Struct Definitions (`cli/cmdos/os_info_types.go`)
Extend the type definitions to expose `OSVersionInfo`:

```go
package cmdos

// OSVersionInfo provides typed, structured operating system and kernel metrics.
type OSVersionInfo struct {
	Type            string `json:"type"`                      // "windows", "linux", "darwin", "freebsd"
	IsWindowsServer bool   `json:"isWindowsServer"`           // Affirmative: true if Windows Server edition
	IsWindowsPro    bool   `json:"isWindowsPro"`              // Affirmative: true if Windows Pro/Enterprise
	IsWindowsHome   bool   `json:"isWindowsHome"`             // Affirmative: true if Windows Home/Core
	Version         string `json:"version"`                   // Canonical semver/composite string (e.g. "10.0.26100.2161")
	Major           int    `json:"major"`                     // Major release number (e.g. 10, 24, 15)
	Minor           int    `json:"minor"`                     // Minor release number (e.g. 0, 4, 3)
	Build           int    `json:"build"`                     // Build number (e.g. 26100)
	Release         string `json:"release"`                   // Marketing / code release (e.g. "24H2", "Noble", "Sequoia")
	Edition         string `json:"edition"`                   // Human-readable edition (e.g. "Windows 11 Pro", "Ubuntu 24.04 LTS")
	KernelVersion   string `json:"kernelVersion"`             // Underling kernel string (e.g. "10.0.26100.2161", "6.8.0-45-generic")
}

// Update LocalOSProbe to carry structured OSVersionInfo
type LocalOSProbe struct {
	OSType       string
	OSGroup      string
	OSVersion    string
	BuildVersion string
	Kernel       string
	VersionInfo  OSVersionInfo
}
```

Update `MachineIdentity` in `cli/cmdos/os_machine_alias.go`:
```go
type MachineIdentity struct {
	Sequence         int           `json:"sequence"`
	NodeID           string        `json:"nodeId"`
	IPAddress        string        `json:"ipAddress"`
	Alias            string        `json:"alias"`
	MachineName      string        `json:"machineName"`
	OSHostname       string        `json:"osHostname"`
	PreviousName     string        `json:"previousName,omitempty"`
	PreviousAlias    string        `json:"previousAlias,omitempty"`
	OSPlatform       string        `json:"osPlatform"`
	OSVersion        string        `json:"osVersion,omitempty"`
	Architecture     string        `json:"architecture,omitempty"`
	CPUCores         int           `json:"cpuCores,omitempty"`
	CurrentUser      string        `json:"currentUser,omitempty"`
	GitMapVersion    string        `json:"gitmapVersion,omitempty"`
	Scope            string        `json:"scope"`
	OSVersionDetails OSVersionInfo `json:"osVersionDetails"`
}
```

### 2.2 Windows Registry Prober (`cli/cmdos/os_info_windows.go`)
Refactor `probeLocalPlatformOS()` to query:
1. `registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)`
2. Extract:
   - `InstallationType` (string)
   - `EditionID` (string)
   - `ProductName` (string)
   - `DisplayVersion` (string)
   - `CurrentBuildNumber` (string)
   - `UBR` (uint64/uint32)
3. Calculate Build & Edition:
   - Parse `buildNum = strconv.Atoi(currentBuildNumber)`
   - Composite build: `fmt.Sprintf("%s.%d", currentBuildNumber, ubr)`
   - Windows 11 check: if `buildNum >= 22000`, ensure edition string reads `"Windows 11 ..."` rather than `"Windows 10 ..."`.
4. Derive affirmative booleans:
   - `IsWindowsServer = strings.EqualFold(installType, "Server") || strings.Contains(strings.ToLower(editionId), "server")`
   - `IsWindowsPro = strings.Contains(strings.ToLower(editionId), "pro") || strings.Contains(strings.ToLower(editionId), "enterprise")`
   - `IsWindowsHome = strings.EqualFold(editionId, "core") || strings.Contains(strings.ToLower(editionId), "home")`

### 2.3 Unix & Darwin Probers (`cli/cmdos/os_info_other.go`)
1. **Linux:**
   - Parse `/etc/os-release` for `ID`, `PRETTY_NAME`, `VERSION_ID`, `VERSION_CODENAME`.
   - Run `uname -r` for `KernelVersion`.
   - Parse `VERSION_ID` into `Major` and `Minor`.
   - Set `Type = "linux"`, affirmative Windows booleans to `false`.
2. **Darwin (macOS):**
   - Run `sw_vers -productVersion` and `sw_vers -buildVersion`.
   - Run `uname -r`.
   - Map Darwin major version to code release name (`15` = `"Sequoia"`, `14` = `"Sonoma"`, etc.).
   - Set `Type = "darwin"`, affirmative Windows booleans to `false`.

### 2.4 Terminal Tree & JSON Renderers (`cli/cmdos/os_machine_alias.go`)
1. **Terminal Tree View (`renderSingleMachineIdentity`):**
   Render hierarchical branch lines:
   ```
   ● Machine Identity:
     ├── Alias:             <Alias>
     ├── Hostname:          <MachineName>
     ├── IP Address:        <IPAddress>
     ├── Operating System:  <Edition> [<Badge>]
     │   ├── Edition:       <Edition> (<Arch>)
     │   ├── Version:       <Version> (Release <Release>)
     │   ├── Build:         <Build>
     │   └── Kernel:        <KernelVersion>
     ├── Architecture:      <Arch> (<NumCPU> Cores)
     ├── Current User:      <User>
     └── GitMap Version:    <Version>
   ```
2. **JSON View (`printIdentityJSON`):**
   Ensure `OSVersionDetails` serializes with camelCase keys (`isWindowsServer`, `isWindowsPro`, `isWindowsHome`, `kernelVersion`, etc.).

---

## 3. Acceptance Criteria & Invariants

- **`AC-01` (Windows Registry Probing):** Probing on Windows queries `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion` and extracts `DisplayVersion`, `CurrentBuildNumber`, `UBR`, and `EditionID` without spawning PowerShell or command prompt subprocesses.
- **`AC-02` (Windows 11 Normalization):** Machines running build numbers `>= 22000` report `"Windows 11"` in `Edition` and `OSVersionDetails.Edition`, even when the legacy `ProductName` registry value contains `"Windows 10"`.
- **`AC-03` (Affirmative Booleans):** The booleans `IsWindowsServer`, `IsWindowsPro`, and `IsWindowsHome` accurately reflect the probed OS edition and are strictly positive (never negative polarity).
- **`AC-04` (Terminal Tree Visual Quality):** Interactive invocation `gitmap machine` renders a structured Unicode tree hierarchy with proper indentation, color highlighting, and edition badge pills (`[Server]`, `[Pro]`, `[Home]`).
- **`AC-05` (Backward-Compatible JSON Schema):** Running `gitmap machine --json` returns a valid JSON array whose root object contains all previous keys (`ipAddress`, `alias`, `machineName`, `osPlatform`, `osVersion`, `architecture`, `currentUser`), plus the new `osVersionDetails` object.
- **`INV-01` (Affirmative Booleans):** No negative boolean flags (`IsNotServer`, `IsNonPro`).
- **`INV-02` (Backward Compatible JSON):** No breaking field renames or deletions.
- **`INV-03` (Strict Relative Paths):** Only relative Git paths permitted in code and tests.
- **`INV-04` (Safe Non-Panicking Probes):** Probes return safe fallbacks on errors, with zero crashes.

---

## 4. Verification Protocol & Test Execution

### 4.1 Unit Test Commands
Execute the Go test suite for the `cmdos` package:
```bash
go test -v ./cli/cmdos -run TestOSVersionInfo
go test -v ./cli/cmdos -run TestWindows11Normalization
go test -v ./cli/cmdos -run TestAffirmativeBooleans
go test -v ./cli/cmdos -run TestJSONBackwardCompatibility
```

### 4.2 Interactive CLI Verification
Run the compiled CLI to verify output formatting:
```bash
# 1. Verify terminal hierarchical tree rendering
gitmap machine

# 2. Verify machine alias shorthand
gitmap mac

# 3. Verify JSON output structure
gitmap machine --json
```

### 4.3 Evidence Checks
- [ ] `gitmap machine` outputs the tree with `● Machine Identity:` and nested `├──` / `└──` elements.
- [ ] On Windows 11, `gitmap machine` displays `Windows 11 Pro [Pro]` rather than `Windows 10`.
- [ ] `gitmap machine --json` validates against `jq` or JSON parsers and contains `"osVersionDetails": { "isWindowsPro": true, ... }`.
- [ ] All package unit tests pass with code 0.
