# TypeScript ExportStatus Enum — `src/lib/enums/export-status.ts`

> **Version**: 1.0.0  
> **Last updated**: 2026-02-27  
> **Tracks**: Issue #10 (`spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

---

## Purpose

Typed enum for import/export operation lifecycle states. Replaces `exportStatus === 'completed'` magic strings in frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/export-status.ts

export enum ExportStatus {
  Pending = "PENDING",
  Processing = "PROCESSING",
  Completed = "COMPLETED",
  Failed = "FAILED",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (exportStatus === 'completed') { ... }

// ✅ CORRECT: Enum constant
if (exportStatus === ExportStatus.Completed) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{!isExporting && exportStatus !== 'completed' && <ExportForm />}

// ✅ CORRECT
{!isExporting && exportStatus !== ExportStatus.Completed && <ExportForm />}
```

### Type Definitions

```typescript
// ❌ WRONG
interface ExportState {
  status: 'pending' | 'processing' | 'completed' | 'failed';
}

// ✅ CORRECT
interface ExportState {
  status: ExportStatus;
}
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `05-features/03-project-management/02-import-export-ui.md` | `exportStatus === 'completed'/'failed'/'processing'` |
| `05-features/27-automation-pipeline/20-import-export.md` | Import/export status checks |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/exportstatus` | `src/lib/enums/export-status.ts` |
| Type | `byte` iota | String enum |
| Values | `Pending`, `Processing`, `Completed`, `Failed` | Same |

---

## Cross-References

- Issue #10 — Domain Status Magic Strings <!-- external: spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./05-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./08-typescript-standards-reference.md) — Parent spec

---

*ExportStatus enum v1.0.0 — 2026-02-27*
