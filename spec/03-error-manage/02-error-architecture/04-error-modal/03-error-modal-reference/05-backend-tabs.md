# Backend Section (Tabs)

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

## Tab: Overview

Displays the primary error information:
- Error message with code badge and timestamp
- Target site URL (if WordPress operation)
- API request method and endpoint
- **Backend Error** banner (from `envelopeErrors.BackendMessage`) — red-themed
- Timing section (`requestedAt`, `requestDelegatedAt`)
- Availability badges (session, stack traces, execution logs)

## Tab: Log

Fetches and displays `error.log.txt` from the backend:
- Auto-fetched when the modal opens on the Backend section
- Refresh, Copy, and Download buttons
- `ScrollArea` with monospace font, 400px height

## Tab: Execution

Two sections:
1. **Go Call Chain** — Table from `envelopeMethodsStack.Backend`
2. **Session Execution Logs** — Timeline from `backendLogs[]` with level-based coloring and step labels

## Tab: Stack

Multi-source stack traces:
1. **Go Backend** — From `envelopeErrors.Backend` (blue-themed)
2. **PHP Delegated** — From `envelopeErrors.DelegatedServiceErrorStack` (orange-themed)
3. **PHP Structured Frames** — Table from `phpStackFrames[]` (file, line, class::function)
4. **Session Diagnostics** — Auto-fetched Go and PHP frames from session API
5. **PHP stacktrace.txt** — Raw backtrace from session diagnostics

## Tab: Session

Full session-level diagnostics (only shown when `sessionId` exists):
- Sub-tabs: **Logs**, **Request**, **Response**, **Stack Trace**
- Fetches from `GET /api/v1/sessions/{id}/logs` and `GET /api/v1/sessions/{id}/diagnostics`
- Log rendering with color-coded levels and stage headers

## Tab: Request

**Request Chain Visualization** — see [07-request-chain.md](./07-request-chain.md) for full details.

Plus environment diagnostics: API base, VITE_API_URL, resolved origin, UI origin.

## Tab: Traversal

See [08-traversal-details.md](./08-traversal-details.md) for full details.

### Sections Rendered (Priority Order):

1. **Endpoint Flow** — React → Go → Delegated badges with full URLs (3-hop when DelegatedRequestServer present)
2. **Methods Stack Table** — Go call chain from envelope
3. **Delegated Server Details** — Purple-themed block with endpoint, method, status, stack trace, response (NEW v2.0.0)
4. **Delegated Service Error Stack** — Orange-themed PHP error lines (legacy, ScrollArea)
5. **Backend Trace** — Go stack trace lines (ScrollArea)

---

*Backend tabs — updated: 2026-03-31*
