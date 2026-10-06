# 01-cli-architecture: CLI Components, Shell Completion & Output Formatting Specification

- **Spec ID:** `01-cli-architecture/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Cobra Commands, Auto-Complete, ANSI Box Formatters, PowerShell Profile
- **Dependencies:** `cli/cmd`, `cli/completion`, `cli/termpad`, `cli/termtable`, `cli/jsonenv`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The CLI Architecture cluster encompasses four primary operational components:

```
01-cli-architecture/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Command Dispatcher** | `cli/cmd/root.go`, `cli/cmd/dispatch.go` | Top-level Cobra command tree, flags, and command execution lifecycle. |
| **Completion Subsystem** | `cli/completion/all_commands.go`, `cli/cmd/root_cobra_completion.go` | Dynamic command & argument suggestions across PowerShell, Bash, and Zsh. |
| **Terminal Box Renderer** | `cli/termpad/`, `cli/termtable/`, `cli/termhelp/` | Universal markdown box display, responsive padding, middle-ellipsized tables. |
| **JSON Envelope Engine** | `cli/jsonenv/envelope.go` | Typed JSON envelope formatting, machine parsing, and secret sanitization. |

---

## 2. Command Dispatcher & Cobra Integration

### 2.1 Universal Command Flags
Every Cobra command inherits standard persistent flags:
- `--json`, `-j` (`bool`): Directs output into `jsonenv.EnvelopeV2`.
- `--verbose`, `-v` (`bool`): Enables detailed execution telemetry.
- `--no-color` (`bool`): Strips ANSI escape sequences.
- `--config` (`string`): Path to custom configuration file.

### 2.2 Execution Error Wrapping
All command runners wrap return errors with structured `appfault.AppError`:

```go
func RunCommand(cmd *cobra.Command, args []string) error {
    res, err := executeSubsystem(args)
    if err != nil {
        return appfault.New("ERR_CLI_EXECUTION", "failed executing subsystem command", err)
    }
    return nil
}
```

---

## 3. Shell Tab Completion Engine

### 3.1 AllCommands Dynamic Registration
The completion engine maintains an authoritative index of all subcommands and aliases in `completion.AllCommands()`:

```go
package completion

type CommandMetadata struct {
    Name        string   `json:"name"`
    Aliases     []string `json:"aliases"`
    Description string   `json:"description"`
    Category    string   `json:"category"`
    IsVisible   bool     `json:"isVisible"`
}

func AllCommands() []CommandMetadata {
    // Returns full catalog of 500+ commands
}
```

### 3.2 PowerShell Profile Integration & Predictive Suggestions
The PowerShell configuration script (`scripts/install.ps1`) configures predictive autocompletion:
- Injects `Register-ArgumentCompleter -Native -CommandName gitmap -ScriptBlock { ... }`.
- Sets PSReadLine prediction source to `HistoryAndPlugin` with `ListView` rendering mode.
- Avoids pipeline stdout interception by writing directory change targets to `$env:GITMAP_HANDOFF_FILE`.

---

## 4. Markdown Box Display & ANSI Formatting

### 4.1 Box Rendering Protocol
Help text and tabular results are formatted via UTF-8 / ASCII bordered boxes:
- Border styles: Rounded (`╭──╮`), Sharp (`┌──┐`), Double (`╔══╗`), or ASCII (`+--+`) fallback for basic terminals.
- Header formatting: Bold cyan title aligned with uppercase section markers.
- Table columns: Adaptive column width calculation with middle-ellipsizing (`/very/long/.../path/to/repo`).

```text
┌─ GitMap Command Summary ───────────────────────────────────────────┐
│ Command       Alias     Description                                 │
├────────────────────────────────────────────────────────────────────┤
│ pull-all      pa, pat   Concurrently pulls all discovered git repos │
│ nodes         fleet     Inspects and orchestrates fleet machines   │
└────────────────────────────────────────────────────────────────────┘
```

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isCobraDispatchVerified: true
  isAllCommandsRegistered: true
  isBoxRendererValidated: true
  isPositiveBooleansUsed: true
```

- [x] All 500+ commands discoverable in shell tab completion.
- [x] Responsive box tables resize without line wrapping artifacts.
- [x] JSON envelope matches V2 contract.
