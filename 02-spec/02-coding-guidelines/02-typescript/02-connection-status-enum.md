# TypeScript ConnectionStatus Enum — `src/lib/enums/connection-status.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string connection status comparisons (`'connected'`, `'disconnected'`, `'connecting'`, `'reconnecting'`, `'error'`) by standardizing on `ConnectionStatus` enum.
> **/learn** Master TypeScript enum architecture: `export enum ConnectionStatus` with uppercase string values (`"CONNECTED"`), typed component props, and state models.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw string connection status comparisons (`connection.status === 'connected'`) with `ConnectionStatus` enum constants.
- [ ] `/learn` Never use string union types (`'connected' | 'disconnected'`) for connection status; use `ConnectionStatus`.
- [ ] `/goal` Ensure conditional UI rendering gates use `ConnectionStatus.Disconnected` and related enum members.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 1.0.0
> **Last updated**: 2026-02-27
> **Tracks**: Issue #10 (`02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

---

## Purpose

Typed enum for WebSocket, SSE, and service connection lifecycle states. Replaces `connection.status === 'connected'` magic strings in frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/connection-status.ts

export enum ConnectionStatus {
  Connected = "CONNECTED",
  Disconnected = "DISCONNECTED",
  Connecting = "CONNECTING",
  Reconnecting = "RECONNECTING",
  Error = "ERROR",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (connection.status === 'connected') { ... }

// ✅ CORRECT: Enum constant
if (connection.status === ConnectionStatus.Connected) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{wsStatus === 'disconnected' && <ReconnectBanner />}

// ✅ CORRECT
{wsStatus === ConnectionStatus.Disconnected && <ReconnectBanner />}
```

### Type Definitions

```typescript
// ❌ WRONG
interface WebSocketState {
  status: 'connected' | 'disconnected' | 'connecting' | 'error';
}

// ✅ CORRECT
interface WebSocketState {
  status: ConnectionStatus;
}
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `05-features/05-voice-input/06-voice-session-manager.md` | `connection.status === 'connected'` |
| `01-general-spec/09-api-integration/02-websocket-patterns-api-integration.md` | WebSocket connection status checks |
| `05-features/27-automation-pipeline/24-collaboration.md` | Participant connection state |
| `08-roadmap-overview/05-gap-analysis.md` | Connection status references |
| `16-ai-transcribe-cli/02-frontend/01-testing-ui.md` | Recording connection status |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/connectionstatus` | `src/lib/enums/connection-status.ts` |
| Type | `byte` iota | String enum |
| Values | `Connected`, `Disconnected`, `Connecting`, `Reconnecting`, `Error` | Same |

---

## Cross-References

- Issue #10 — Domain Status Magic Strings <!-- external: 02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./06-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent spec

---

*ConnectionStatus enum v1.0.0 — 2026-02-27*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-002: ConnectionStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All connection status handling strictly utilizes `ConnectionStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
