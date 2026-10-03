# TypeScript ExportStatus Enum — `src/lib/enums/export-status.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string export and import status literals (`'pending'`, `'processing'`, `'completed'`, `'failed'`) by standardizing on `ExportStatus`.
> **/learn** Master import/export operation typing: `export enum ExportStatus` with uppercase string values, typed state interfaces, and conditional rendering guards.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw string export status checks (`exportStatus === 'completed'`) with `ExportStatus.Completed`.
- [ ] `/learn` Never use string union types (`'pending' | 'processing' | 'completed' | 'failed'`) for export state; use `ExportStatus`.
- [ ] `/goal` Ensure UI conditional rendering guards utilize `ExportStatus` enum constants.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 1.0.0
> **Last updated**: 2026-02-27
> **Tracks**: Issue #10 (`02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

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

- Issue #10 — Domain Status Magic Strings <!-- external: 02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./06-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent spec

---

*ExportStatus enum v1.0.0 — 2026-02-27*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-005: ExportStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All import and export operation states strictly utilize `ExportStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
