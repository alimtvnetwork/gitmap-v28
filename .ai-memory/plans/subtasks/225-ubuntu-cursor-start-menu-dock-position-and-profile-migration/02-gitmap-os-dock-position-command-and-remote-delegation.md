# Subtask Plan 02: Cross-Platform GitMap OS Dock & Start Menu Position Command with Remote Delegation

- **Spec Reference:** [02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/01-architecture-spec.md](../../../../02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/01-architecture-spec.md)
- **Status:** Queued
- **Target Area:** `cli/cmdos/os_dock_types.go`, `cli/cmdos/os_dock_cmd.go`, `cli/cmdos/os_dock_linux.go`, `cli/cmdos/os_dock_windows.go`, `cli/cmdos/os_dock_darwin.go`, `cli/cmdos/os_dock_other.go`, `cli/cmdos/os_dock_test.go`, `cli/cmdos/os.go`, `cli/cmdos/os_help_modern.go`

---

## 1. Objective

Implement a cross-platform command in GitMap OS to query and configure workstation dock, panel, and start menu positioning with remote fleet node delegation:
1. Command syntax: `gitmap os dock [bottom|left|right|top] [--node <alias>]` with aliases `panel`, `start-menu`, `dock-position`, and `taskbar`.
2. Query Mode: If zero position arguments are given, query and display current dock/taskbar configuration.
3. Mutation Mode: If a position argument is given, validate, normalize, and update dock position.
4. Remote Delegation: If `--node <alias>` is specified, delegate execution via SSH to target node using `cmdssh.RunSSHExec`.
5. Linux Engine: Control GNOME Shell dock orientation via `org.gnome.shell.extensions.dash-to-dock dock-position` (`'TOP'`, `'RIGHT'`, `'BOTTOM'`, `'LEFT'`) with automatic live DBus session bus resolution.
6. Windows Engine: Query and configure Windows taskbar alignment / screen position via registry (`TaskbarAl` / `StuckRects3`).
7. Help & Registry: Register commands in `cli/cmdos/os.go` and document in `cli/cmdos/os_help_modern.go`.

---

## 2. Technical Contracts & Subsystem Design

### 2.1 Core Types & Interfaces (`cli/cmdos/os_dock_types.go`)

```go
package cmdos

// DockPosition represents cardinal screen edge orientation for desktop dock.
type DockPosition string

const (
	DockPositionBottom DockPosition = "bottom"
	DockPositionLeft   DockPosition = "left"
	DockPositionRight  DockPosition = "right"
	DockPositionTop    DockPosition = "top"
	DockPositionCenter DockPosition = "center"
)

// DockConfig holds status metadata about the desktop dock.
type DockConfig struct {
	Platform           string         `json:"platform"`
	Engine             string         `json:"engine"`
	Position           DockPosition   `json:"position"`
	SupportedPositions []DockPosition `json:"supportedPositions"`
	CanMutate          bool           `json:"canMutate"`
	Details            string         `json:"details,omitempty"`
}

// DockOperator defines platform-specific querying and setting operations.
type DockOperator interface {
	GetDockConfig() (*DockConfig, error)
	SetDockPosition(pos DockPosition) error
}

// DockCLIOptions holds parsed CLI parameters.
type DockCLIOptions struct {
	TargetNode string
	Position   DockPosition
	IsJSON     bool
	IsHelp     bool
}
```

### 2.2 CLI Parsing & Remote Delegation Engine (`cli/cmdos/os_dock_cmd.go`)

1. **Option Extraction:**
   - Scan `args` for `--node <alias>` or `--node=<alias>`.
   - Scan for `--json` and help flags (`-h`, `--help`).
   - Extract position parameter (`bottom`, `left`, `right`, `top`, `center`).
2. **Remote Node Branch:**
   ```go
   if opts.TargetNode != "" {
       return delegateRemoteDock(opts)
   }
   ```
   - Build remote command string:
     - Query: `gitmap os dock`
     - Mutation: `gitmap os dock <position>`
     - If `--json`: append ` --json`
   - Print notice:
     `● Delegating dock configuration to remote node '<alias>'...`
   - Execute:
     `return cmdssh.RunSSHExec([]string{opts.TargetNode, remoteCmd})`
3. **Local Dispatch Branch:**
   - Query: `handleDockStatus(engine, opts.IsJSON)`
   - Mutation: `handleDockSet(engine, opts.Position)`

### 2.3 Platform Implementations

#### Linux Platform (`cli/cmdos/os_dock_linux.go` with `//go:build linux`)
- Uses `gsettings` under active DBus session (`/run/user/<uid>/bus`).
- Schema: `org.gnome.shell.extensions.dash-to-dock`.
- Key: `dock-position`.
- Read:
  ```go
  cmd := exec.Command("gsettings", "get", "org.gnome.shell.extensions.dash-to-dock", "dock-position")
  // Output: 'LEFT' -> DockPositionLeft
  ```
- Write:
  ```go
  // Normalize pos to uppercase string: "BOTTOM", "LEFT", "RIGHT", "TOP"
  cmd := exec.Command("gsettings", "set", "org.gnome.shell.extensions.dash-to-dock", "dock-position", fmt.Sprintf("'%s'", upperVal))
  ```

#### Windows Platform (`cli/cmdos/os_dock_windows.go` with `//go:build windows`)
- Inspects Windows version.
- On Windows 11:
  - Reads/Writes `HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced\TaskbarAl`.
  - `0`: Left, `1`: Center.
- On Windows 10:
  - Reads/Writes `StuckRects3` byte 12 (0=Left, 1=Top, 2=Right, 3=Bottom).

#### Darwin & Other Platforms (`cli/cmdos/os_dock_darwin.go` & `os_dock_other.go`)
- macOS (`//go:build darwin`):
  - Uses `defaults read com.apple.dock orientation` and `defaults write com.apple.dock orientation <pos>` followed by `killall Dock`.
- Other (`//go:build !linux && !windows && !darwin`):
  - Returns `apperror.NewSimple("dock positioning is not supported on this operating system", "E_UNSUPPORTED_PLATFORM")`.

---

## 3. Subcommand Registration & Modern Help Update

### In `cli/cmdos/os.go`:
```go
case "dock", "panel", "start-menu", "dock-position", "taskbar":
    return runOSDockCommand(subArgs)
```

### In `cli/cmdos/os_help_modern.go`:
Under `buildOSSystemAndNetworkSection()`:
```go
{Command: "os dock [bottom|left|right|top] [--node]", Description: "Inspect or configure desktop dock / taskbar position (aliases: panel, start-menu)"},
```

---

## 4. Unit Testing Strategy (`cli/cmdos/os_dock_test.go`)

1. **`TestNormalizeDockPosition`:**
   - Assert `"bottom"`, `"BOTTOM"`, `"Bottom"` -> `DockPositionBottom`.
   - Assert `"left"` -> `DockPositionLeft`, `"right"` -> `DockPositionRight`, `"top"` -> `DockPositionTop`.
   - Assert invalid values (`"diagonal"`, `"floating"`) return validation error.
2. **`TestParseDockCLIOptions`:**
   - Test empty args -> query mode.
   - Test `["bottom"]` -> mutation mode with `DockPositionBottom`.
   - Test `["left", "--node", "u1"]` -> remote delegation to `"u1"`.
   - Test `["--node=u1", "top"]` -> remote delegation to `"u1"`.
3. **`TestMockDockOperator`:**
   - Test mock operator simulating successful query and set operations.
   - Ensure tests are isolated and do not invoke real OS desktop commands or power functions.

---

## 5. Verification Commands

Run the following checks to verify execution:

```bash
# 1. Local command help
gitmap os dock --help
gitmap os panel --help

# 2. Local status query
gitmap os dock
gitmap os dock --json

# 3. Local position update (Linux / Windows)
gitmap os dock bottom
gitmap os dock left

# 4. Remote delegation to node u1
gitmap os dock --node u1
gitmap os dock bottom --node u1
gitmap os dock left --node u1

# 5. Unit test execution
go test ./cli/cmdos/... -v -run "TestDock"
```

---

## 6. Acceptance Criteria

- [ ] `gitmap os dock`, `gitmap os panel`, `gitmap os start-menu` dispatches correctly without error.
- [ ] Querying without arguments renders current dock position and engine details cleanly.
- [ ] Setting valid positions (`bottom`, `left`, `right`, `top`) applies changes using platform engine.
- [ ] Passing `--node <alias>` delegates via `cmdssh.RunSSHExec` to the remote node over SSH.
- [ ] Comprehensive unit tests in `cli/cmdos/os_dock_test.go` achieve 100% pass rate with mocks.
