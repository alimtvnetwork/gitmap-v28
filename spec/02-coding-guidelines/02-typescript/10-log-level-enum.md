# TypeScript LogLevel Enum — `src/lib/enums/log-level.ts`

> **Version**: 1.0.0  
> **Last updated**: 2026-03-31  
> **Tracks**: Magic string elimination in error modal color-themes spec

---

## Purpose

Typed enum for application log severity levels. Replaces `level === 'error'` magic strings in logging, error display, and color-theme specs.

---

## Reference Implementation

```typescript
// src/lib/enums/log-level.ts

export enum LogLevel {
  Debug = "DEBUG",
  Info = "INFO",
  Warn = "WARN",
  Error = "ERROR",
  Fatal = "FATAL",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (entry.level === 'error') { ... }

// ✅ CORRECT: Enum constant
if (entry.level === LogLevel.Error) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{level === 'warn' && <WarningIcon />}

// ✅ CORRECT
{level === LogLevel.Warn && <WarningIcon />}
```

### Color Theme Mapping

```typescript
// ❌ WRONG: Raw strings in map keys
const colorMap = {
  'error': 'red',
  'warn': 'yellow',
};

// ✅ CORRECT: Enum keys
const colorMap: Record<LogLevel, string> = {
  [LogLevel.Error]: 'var(--color-error)',
  [LogLevel.Warn]: 'var(--color-warning)',
  [LogLevel.Info]: 'var(--color-info)',
  [LogLevel.Debug]: 'var(--color-muted)',
  [LogLevel.Fatal]: 'var(--color-critical)',
};
```

### Type Definitions

```typescript
// ❌ WRONG
interface LogEntry {
  level: 'debug' | 'info' | 'warn' | 'error' | 'fatal';
}

// ✅ CORRECT
interface LogEntry {
  level: LogLevel;
}
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `spec/03-error-manage/02-error-architecture/04-error-modal/04-color-themes.md` | `LogLevel.Error`, `LogLevel.Warn`, `LogLevel.Info`, `LogLevel.Debug` color mappings |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/loglevel` | `src/lib/enums/log-level.ts` |
| Type | `byte` iota | String enum |
| Values | `Debug`, `Info`, `Warn`, `Error`, `Fatal` | Same |

---

## Cross-References

- [ConnectionStatus Enum](./01-connection-status-enum.md) — Sibling enum spec
- [HttpMethod Enum](./05-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./08-typescript-standards-reference.md) — Parent spec

---

*LogLevel enum v1.0.0 — 2026-03-31*
