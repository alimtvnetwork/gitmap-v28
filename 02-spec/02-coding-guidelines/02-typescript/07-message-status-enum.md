# TypeScript MessageStatus Enum — `src/lib/enums/message-status.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string chat message lifecycle states (`'pending'`, `'streaming'`, `'completed'`, `'error'`) by standardizing on `MessageStatus`.
> **/learn** Master message lifecycle typing: `export enum MessageStatus` with uppercase string values, typed chat message interfaces, and conditional action rendering.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw string message status comparisons (`message.status === 'streaming'`) with `MessageStatus.Streaming`.
- [ ] `/learn` Never define string union types (`'pending' | 'streaming' | 'complete' | 'error'`) in chat state models; use `MessageStatus`.
- [ ] `/goal` Ensure chat UI conditional components (like retry buttons and loading spinners) check `MessageStatus` constants.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 1.0.0
> **Last updated**: 2026-02-27
> **Tracks**: Issue #10 (`02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

---

## Purpose

Typed enum for AI chat message lifecycle states — pending send, actively streaming, completed, or errored. Replaces `message.status === 'streaming'` magic strings in frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/message-status.ts

export enum MessageStatus {
  Pending = "PENDING",
  Streaming = "STREAMING",
  Completed = "COMPLETED",
  Error = "ERROR",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (message.status === 'streaming') { ... }

// ✅ CORRECT: Enum constant
if (message.status === MessageStatus.Streaming) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{message.status === 'error' && <RetryButton />}

// ✅ CORRECT
{message.status === MessageStatus.Error && <RetryButton />}
```

### Type Definitions

```typescript
// ❌ WRONG
interface ChatMessage {
  status: 'pending' | 'streaming' | 'complete' | 'error';
}

// ✅ CORRECT
interface ChatMessage {
  status: MessageStatus;
}
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `05-features/25-ai-enhancements/05-03-message-display.md` | `message.status === 'streaming'/'error'/'pending'` |
| `05-features/06-ai-integration/08-ai-chat-ui.md` | Chat message status checks |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/messagestatus` | `src/lib/enums/message-status.ts` |
| Type | `byte` iota | String enum |
| Values | `Pending`, `Streaming`, `Completed`, `Error` | Same |

---

## Cross-References

- Issue #10 — Domain Status Magic Strings <!-- external: 02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./06-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent spec

---

*MessageStatus enum v1.0.0 — 2026-02-27*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-007: MessageStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All chat message lifecycle states strictly utilize `MessageStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
