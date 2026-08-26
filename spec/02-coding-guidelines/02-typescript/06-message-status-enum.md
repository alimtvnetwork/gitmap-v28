# TypeScript MessageStatus Enum — `src/lib/enums/message-status.ts`

> **Version**: 1.0.0  
> **Last updated**: 2026-02-27  
> **Tracks**: Issue #10 (`spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

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

- Issue #10 — Domain Status Magic Strings <!-- external: spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./05-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./08-typescript-standards-reference.md) — Parent spec

---

*MessageStatus enum v1.0.0 — 2026-02-27*
