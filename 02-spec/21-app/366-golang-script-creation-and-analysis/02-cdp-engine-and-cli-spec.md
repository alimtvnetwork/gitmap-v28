# Spec 366: Golang Script Creation and Analysis — CDP Engine & CLI Specification

- **Task Code:** `366-golang-script-creation-and-analysis`
- **Slug:** `golang-script-creation-and-analysis`
- **Status:** `in-progress`
- **Date:** `31-Oct-26`

---

## 1. Package Organization (`cli/cmdcdp/`)

The CDP inspection suite is implemented as a self-contained Go package under `cli/cmdcdp/` adhering to GitMap's strict 15-line function limit, zero nested conditional branching, and positive boolean naming conventions.

| File | Purpose | Max Function Lines |
|------|---------|---------------------|
| `types.go` | DTOs, data structs, error codes, and constants | N/A (declarations) |
| `client.go` | RFC 6455 WebSocket client using pure standard library | <= 15 lines |
| `discovery.go` | HTTP target and port resolution (`/json/list`) | <= 15 lines |
| `eval.go` | `Runtime.evaluate` JSON-RPC execution and unpacking | <= 15 lines |
| `probe.go` | DOM traversal, input inspection, and CSS selector queries | <= 15 lines |
| `render.go` | Tokyo Night Lipgloss terminal cards and JSON serializers | <= 15 lines |
| `cdp_cmd.go` | Cobra CLI command definitions, flag wiring, and dispatch | <= 15 lines |

---

## 2. Go Type Definitions & Interface Contracts

All domain models and data transfer structs are defined in `cli/cmdcdp/types.go`.

### 2.1 Target and Protocol Structs

```go
package cmdcdp

import (
	"context"
	"encoding/json"
)

// CDPTarget represents a single debuggable target from /json/list.
type CDPTarget struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	DevToolsFrontendURL  string `json:"devtoolsFrontendUrl"`
}

// CDPRequest represents a JSON-RPC 2.0 outgoing command.
type CDPRequest struct {
	ID     int64          `json:"id"`
	Method string         `json:"method"`
	Params map[string]any `json:"params,omitempty"`
}

// CDPResponse represents a JSON-RPC 2.0 incoming response.
type CDPResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *CDPError       `json:"error,omitempty"`
}

// CDPError represents a JSON-RPC error envelope.
type CDPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}
```

### 2.2 DOM Inspection and Evaluation Models

```go
// SelectorHit records an individual matching DOM element discovered during probe.
type SelectorHit struct {
	Selector    string `json:"selector"`
	Tag         string `json:"tag"`
	ID          string `json:"id,omitempty"`
	ClassName   string `json:"className,omitempty"`
	Value       string `json:"value,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	IsVisible   bool   `json:"isVisible"`
	IsFocused   bool   `json:"isFocused"`
}

// DOMProbeResult aggregates all interactive elements and selector hits on the target.
type DOMProbeResult struct {
	TargetURL    string        `json:"targetUrl"`
	Title        string        `json:"title"`
	InputCount   int           `json:"inputCount"`
	Hits         []SelectorHit `json:"hits"`
	IsSuccess    bool          `json:"isSuccess"`
	ErrorMessage string        `json:"errorMessage,omitempty"`
}

// EvaluateResult holds the outcome of an arbitrary JavaScript expression execution.
type EvaluateResult struct {
	Type             string `json:"type"`
	Value            any    `json:"value"`
	Description      string `json:"description,omitempty"`
	IsException      bool   `json:"isException"`
	ExceptionDetails string `json:"exceptionDetails,omitempty"`
}
```

### 2.3 Client Interface Contract

```go
// Client defines the core capabilities required to communicate with a CDP endpoint.
type Client interface {
	Connect(ctx context.Context, wsURL string) error
	Call(ctx context.Context, method string, params map[string]any) ([]byte, error)
	Close() error
}
```

---

## 3. CLI Command Specifications & Syntax

The command hierarchy is exposed under the root `gitmap cdp` umbrella.

### 3.1 `gitmap cdp targets`

Discovers all available browser tabs, web workers, and IDE panels running CDP endpoints.

- **Syntax:**
  ```bash
  gitmap cdp targets [--port <port>] [--timeout <sec>] [--json]
  ```
- **Flags:**
  - `--port, -p`: Optional DevTools port (defaults to probing `9222`, `9229`, then scanning local Chrome processes).
  - `--timeout, -t`: Timeout in seconds (default: `5`).
  - `--json`: Output raw JSON target array without terminal styling.
- **Example Human Output (Tokyo Night):**
  ```text
  ╭── Chrome DevTools Targets (Port 9222) ─────────────────────────╮
  │ #  Type  Target ID   Title                         URL         │
  │ 1  page  3658C5BD..  Antigravity Visualizer        http://...  │
  │ 2  page  9F10A24C..  New Tab                       chrome://.. │
  ╰────────────────────────────────────────────────────────────────╯
  ```

### 3.2 `gitmap cdp probe`

Inspects the active browser DOM, enumerates form controls (`<input>`, `<textarea>`, `<button>`), and tests specified CSS selectors.

- **Syntax:**
  ```bash
  gitmap cdp probe [--port <port>] [--target <targetId>] [--selector <css>] [--json]
  ```
- **Flags:**
  - `--port, -p`: Remote debugging port (default: auto-discover).
  - `--target`: Target ID substring or index (default: first active `page` target).
  - `--selector, -s`: Optional CSS selector to query (e.g. `button.submit-btn`, `#username`).
  - `--json`: Output machine-readable `DOMProbeResult` JSON.
- **Example Human Output (Tokyo Night):**
  ```text
  ╭── DOM Probe: Antigravity Visualizer ────────────────────────────╮
  │ URL: http://localhost:3000/dashboard                           │
  │ Total Interactive Inputs: 4                                    │
  ├────────────────────────────────────────────────────────────────┤
  │ Target Elements:                                               │
  │  • input#search [type=text] (Visible: yes, Focused: no)        │
  │    placeholder: "Search repositories..."                       │
  │  • input#branch-filter [type=text] (Visible: yes, Focused: yes)│
  │    value: "main"                                               │
  │  • button.btn-primary (Visible: yes)                           │
  │    text: "Execute Scan"                                        │
  ╰────────────────────────────────────────────────────────────────╯
  ```

### 3.3 `gitmap cdp eval`

Executes a JavaScript expression inside the remote browser context and returns the result.

- **Syntax:**
  ```bash
  gitmap cdp eval "<js-expression>" [--port <port>] [--target <targetId>] [--json]
  ```
- **Flags:**
  - `--port, -p`: Remote debugging port (default: auto-discover).
  - `--target`: Target ID substring or index (default: first active `page` target).
  - `--json`: Output raw JSON containing `{ "type": "string", "value": ... }`.
- **Example Invocations:**
  ```bash
  # Check document ready state
  gitmap cdp eval "document.readyState"

  # Query value of active input
  gitmap cdp eval "document.activeElement.value" --json

  # Get window dimensions
  gitmap cdp eval "({ width: window.innerWidth, height: window.innerHeight })" --json
  ```

---

## 4. Coding Guidelines Compliance

### 4.1 Strict 15-Line Function Limit & Guard Clauses

Every function must adhere to <= 15 lines max (<= 8 lines preferred). Nested `if` blocks are banned in favor of guard clauses and early returns.

**Example: Frame Reader Decomposition (`client.go`)**

```go
func (c *wsClient) ReadFrame(ctx context.Context) ([]byte, byte, error) {
	header, err := c.readHeader(ctx)
	if err != nil {
		return nil, 0, err
	}

	payload, err := c.readPayload(ctx, header.PayloadLen)
	if err != nil {
		return nil, 0, err
	}

	return payload, header.Opcode, nil
}

func (c *wsClient) readHeader(ctx context.Context) (frameHeader, error) {
	var buf [2]byte
	if err := c.readFull(ctx, buf[:]); err != nil {
		return frameHeader{}, err
	}

	return parseFrameHeader(buf[0], buf[1]), nil
}
```

### 4.2 Positive Booleans Only

- Allowed: `isVisible`, `isFocused`, `isSuccess`, `hasMask`, `hasMoreFrames`.
- Strictly Banned: `notVisible`, `unfocused`, `isNotSuccess`, `noMask`.

### 4.3 Leaf Error Sovereignty

Leaf functions return typed errors (`error` or `*apperror.AppError`) directly. They must NEVER call `cliexit.HandleError` or `os.Exit`. Outer command handlers in `cdp_cmd.go` catch errors and route them to root dispatch for database audit recording.

---

## 5. Tokyo Night Lipgloss Styling Specification

Human-readable CLI rendering follows GitMap's Tokyo Night design system:

| Token | Hex Value | Application |
|-------|-----------|-------------|
| Background | `#1a1b26` | Card background |
| Foreground | `#c0caf5` | Default body text |
| Cyan Accent | `#7dcfff` | Card headers, selector tags |
| Purple Accent | `#bb9af7` | URLs, command suggestions |
| Success Green | `#9ece6a` | Success states, visible elements |
| Warning Yellow | `#e0af68` | Focus indicators, warnings |
| Error Red | `#f7768e` | Exceptions, failed evaluations |

---

## 6. Verification & Automated Test Matrix

| Test Name | Target File | Verification Purpose |
|-----------|-------------|----------------------|
| `TestRFC6455MaskingKeyGeneration` | `client_test.go` | Verifies masking key randomness via `crypto/rand` |
| `TestRFC6455ClientFrameEncoding` | `client_test.go` | Verifies XOR payload masking conforms to RFC 6455 Sec 5.3 |
| `TestLargePayloadDefragmentation` | `client_test.go` | Validates receiving and concatenating frames > 64 KB without truncation |
| `TestTargetDiscoveryFilter` | `discovery_test.go` | Validates that `/json/list` parses target types and isolates `page` |
| `TestDOMProbeResultJSONSerialization` | `probe_test.go` | Verifies `--json` output contract against `DOMProbeResult` schema |
| `TestEvalExceptionPropagation` | `eval_test.go` | Verifies runtime JS exceptions return structured `EvaluateResult.IsException` |
