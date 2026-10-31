# Subtask 03: Native Go WebSocket CDP Client and Discovery Engine

- **Task Code:** `TASK-03`
- **Slug:** `golang-script-creation-and-analysis - Task 03`
- **Parent Plan:** [366-golang-script-creation-and-analysis.md](../../pending/366-golang-script-creation-and-analysis.md)
- **Target Path:** `cli/cmdcdp/`
- **Status:** `pending`
- **Assigned Worker:** CDP Go Engineer

---

## 1. Goal

Implement the native Go Chrome DevTools Protocol (CDP) WebSocket client, dynamic target discovery engine, and DOM evaluation subsystem in package `cli/cmdcdp/` using pure Go standard library packages with zero external dependencies.

---

## 2. In Scope

1. **Type Definitions (`cli/cmdcdp/types.go`):**
   - Define data models: `CDPTarget`, `CDPRequest`, `CDPResponse`, `DOMProbeResult`, `SelectorHit`, `EvaluatorOptions`.
   - Adhere strictly to positive boolean naming conventions (`hasTarget`, `isMasked`, `hasExtendedPayload`, `isJSONOutput`).

2. **Pure Go RFC 6455 WebSocket Client (`cli/cmdcdp/ws_client.go`):**
   - Connect via `net.DialTimeout` with configurable timeouts (default 5s).
   - Perform HTTP/1.1 Upgrade handshake sending random 16-byte base64 `Sec-WebSocket-Key` and validating `101 Switching Protocols`.
   - Implement client-to-server 4-byte XOR frame masking using `crypto/rand`.
   - Implement RFC 6455 frame framing supporting standard (<= 125 bytes), 16-bit extended (126), and 64-bit extended (127) payload lengths.
   - Implement frame unmasking and buffered frame reading for responses exceeding 64KB.
   - Handle WebSocket control frames (`0x8` Close, `0x9` Ping, `0xA` Pong).

3. **Dynamic Port & Target Discovery Engine (`cli/cmdcdp/discovery.go`):**
   - Query `http://127.0.0.1:<port>/json/list` to retrieve active inspectable targets.
   - Scan default candidate debug ports (`9222`, `49994`, `9229`) when no port is explicitly specified.
   - Filter and select the primary page target or target matching user-specified `--target <id>`.
   - Extract and sanitize the `webSocketDebuggerUrl`.

4. **DOM Probe & JavaScript Evaluator (`cli/cmdcdp/evaluator.go`):**
   - Construct JSON-RPC 2.0 payloads for `Runtime.evaluate`.
   - Embed the standard selector ladder:
     - `[aria-label="Message input"]`
     - `[role="combobox"][contenteditable="true"]`
     - `div.monaco-editor [contenteditable="true"]`
     - `div.interactive-input-editor [contenteditable="true"]`
     - `[role="textbox"][contenteditable="true"]`
     - `textarea.monaco-mouse-cursor-text`
     - `textarea[placeholder*="Ask"]`
     - `div[contenteditable="true"]`
     - `textarea`
     - `input`
   - Execute JavaScript evaluation, parse the returned remote object, and unmarshal into `DOMProbeResult`.

---

## 3. Out of Scope

- Cobra CLI command registration and flag binding (handled in Task 04).
- Terminal Lipgloss card formatting (handled in Task 04).
- Version bumping and release execution (handled in Task 06).

---

## 4. Architectural Invariants

- **Zero External Dependencies:** Only standard library imports (`net`, `crypto/rand`, `encoding/binary`, `encoding/json`, `fmt`, `io`, `strings`, `time`).
- **Leaf Error Sovereignty:** Leaf functions must return `error` or `*apperror.AppError` directly; never invoke exit handlers.
- **Function Limit:** Functions must adhere to <= 15 lines max with zero nested conditionals.
- **Positive Booleans Only:** All boolean fields and parameters must start with `is*` or `has*`.

---

## 5. Verification Checklist

- [ ] `cli/cmdcdp/types.go` defines all required structs with positive booleans.
- [ ] `cli/cmdcdp/ws_client.go` implements RFC 6455 handshake, masking, and extended payload decode.
- [ ] `cli/cmdcdp/discovery.go` discovers targets from `/json/list` and candidate ports.
- [ ] `cli/cmdcdp/evaluator.go` evaluates DOM selector ladder and arbitrary JS expressions.
- [ ] Functions stay within <= 15 lines and conditionals are fully flattened.
