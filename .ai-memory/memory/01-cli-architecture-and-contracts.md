# 01 — CLI Architecture, Contracts & Terminal Ecosystem

- **Domain:** Command-Line Interface, Output Contracts & Macro Automation
- **Authoritative Specification:** [01-cli-architecture](../../02-spec/21-app/01-cli-architecture/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. Cobra Command Hierarchy & Dispatch Engine

GitMap's CLI is built upon the Cobra framework with specialized enhancements for performance and user experience:

- **Command Tree:** Root command `gitmap` dispatches to domain-specific verbs (`clone`, `pull`, `commit`, `status`, `nodes`, `cluster`, `macro`, `pipeline`, `db`, `agy`).
- **Typo Tolerance & Suggestions:** When an unknown subcommand or flag is entered, Levenshtein distance matching identifies and suggests the closest valid command.
- **Fast Startup:** Subcommand registration is lazy; heavy dependencies (such as SSH connection pools or SQLite migration engines) are initialized only upon command execution.
- **Dual Invocation Modes:** All commands support interactive terminal execution (styled ANSI output) and non-interactive script execution (`--json` or `--quiet`).

---

## 2. Standardized JSON Envelope V2

When executed with `--json`, commands output a machine-readable payload adhering to the JSON Envelope V2 contract:

```go
type JSONEnvelopeV2 struct {
    IsSuccess bool        `json:"isSuccess"`
    Data      interface{} `json:"data,omitempty"`
    Error     *AppError   `json:"error,omitempty"`
    Meta      MetaPayload `json:"meta"`
}

type MetaPayload struct {
    DurationMs int64  `json:"durationMs"`
    Command    string `json:"command"`
    Version    string `json:"version"`
}
```

- **Invariant:** Successful execution emits `isSuccess: true` with populated `data` and null `error`.
- **Failure:** Emits `isSuccess: false` with null `data` and structured `*AppError` containing code, message, and caller.
- **Zero Raw Stdout in JSON Mode:** Warnings and progress meters are redirected to stderr to prevent corrupting parsed JSON stdout.

---

## 3. Terminal UI, Lipgloss Styling & Table Alignment

GitMap provides a rich, responsive terminal experience powered by Lipgloss:

- **Lipgloss Theme Tokens:** Unified styling tokens across colors, badges, borders, and margins ensure visual consistency.
- **Dry Table Renderer (`termtable`):** Dynamic column width calculation guarantees table borders align regardless of terminal window width.
- **Middle-Ellipsis Path Formatter:** Long filesystem paths exceeding available column space are truncated from the center (`/home/user/.../subdir/file.go`), preserving both root context and filename.
- **Status Badges:** Color-coded status markers (`[OK]` green, `[WARN]` yellow, `[FAIL]` red, `[INFO]` cyan).

---

## 4. StreamWriter Unbuffered Streaming Contract

Long-running commands (such as macro execution, cluster command delegation, and CI/CD log tailing) use unbuffered streaming interfaces:

```go
type StreamWriter interface {
    WriteLine(line string) error
    WriteChunk(p []byte) (int, error)
    Flush() error
}
```

- **Zero Latency:** Lines are flushed immediately to stdout without batch buffering delays.
- **Graceful Cancellation:** Listens to `context.Context` cancellation signals to terminate remote processes cleanly on `Ctrl+C`.

---

## 5. Macro Automation Engine

GitMap includes an embedded macro recording, serialization, and streaming replay engine:

- **Macro Definition:** Macros are serialized as JSON bundles (`.gmacro`), containing sequential execution steps, environment overrides, and conditional assertions.
- **Run-Until Assertions:** Supports regex-based assertion checks (`run_until: "exit code 0"` or pattern matches) to stop early or retry on transient errors.
- **Execution Idempotency:** Macro steps specify preconditions (`exists`, `absent`, `matches`, `process_running`) to ensure safe repeat runs without state pollution.
- **Interactive Macro Builder (`gitmap mb`):** CLI wizard allowing users to interactively record commands, inspect recorded arguments, and persist macros.
