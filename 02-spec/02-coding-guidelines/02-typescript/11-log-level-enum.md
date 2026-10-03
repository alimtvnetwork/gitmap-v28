# TypeScript LogLevel Enum — `src/lib/enums/log-level.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string log severity levels (`'debug'`, `'info'`, `'warn'`, `'error'`, `'fatal'`) across logging, error display, and theme maps by standardizing on `LogLevel`.
> **/learn** Master log severity typing: `export enum LogLevel` with uppercase string values, typed `Record<LogLevel, string>` color theme mappings, and log entry interfaces.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw string log level comparisons (`entry.level === 'error'`) with `LogLevel.Error`.
- [ ] `/learn` Never use raw strings in theme mapping dictionary keys; use `Record<LogLevel, string>` with enum member computed keys (`[LogLevel.Error]`).
- [ ] `/goal` Strongly type all `LogEntry` interfaces using `LogLevel` instead of string union literals.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

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
| `02-spec/03-error-manage/02-error-architecture/04-error-modal/04-color-themes.md` | `LogLevel.Error`, `LogLevel.Warn`, `LogLevel.Info`, `LogLevel.Debug` color mappings |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/loglevel` | `src/lib/enums/log-level.ts` |
| Type | `byte` iota | String enum |
| Values | `Debug`, `Info`, `Warn`, `Error`, `Fatal` | Same |

---

## Cross-References

- [ConnectionStatus Enum](./02-connection-status-enum.md) — Sibling enum spec
- [HttpMethod Enum](./06-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent spec

---

*LogLevel enum v1.0.0 — 2026-03-31*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-011: LogLevelEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All log level comparisons, interface models, and theme mappings strictly utilize `LogLevel` enum constants, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
