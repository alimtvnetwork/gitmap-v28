# TypeScript EntityStatus Enum — `src/lib/enums/entity-status.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string entity status literals (`'active'`, `'inactive'`, `'draft'`, `'archived'`) across frontend projects, plugins, and shared components by standardizing on `EntityStatus`.
> **/learn** Master entity lifecycle typing: `export enum EntityStatus` with PascalCase keys and uppercase string values, typed interface fields, and default status constants.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all magic string comparisons (`entity.status === 'active'`) with `EntityStatus.Active`.
- [ ] `/learn` Never define string union types (`'active' | 'inactive' | 'draft' | 'archived'`) for entity status fields in interfaces; use `EntityStatus`.
- [ ] `/goal` Standardize default entity status definitions to reference `EntityStatus.Active` instead of raw string literals.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 1.0.0
> **Last updated**: 2026-02-27
> **Tracks**: Issue #10 (`02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

---

## Purpose

Typed enum for general entity lifecycle states — projects, plugins, shares, resources, and any domain object that can be active, inactive, drafted, or archived. Replaces `entity.status === 'active'` magic strings in frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/entity-status.ts

export enum EntityStatus {
  Active = "ACTIVE",
  Inactive = "INACTIVE",
  Draft = "DRAFT",
  Archived = "ARCHIVED",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (plugin.status === 'active') { ... }

// ✅ CORRECT: Enum constant
if (plugin.status === EntityStatus.Active) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{share.status !== 'active' && 'opacity-60'}

// ✅ CORRECT
{share.status !== EntityStatus.Active && 'opacity-60'}
```

### Type Definitions

```typescript
// ❌ WRONG
interface Project {
  status: 'active' | 'inactive' | 'draft' | 'archived';
}

// ✅ CORRECT
interface Project {
  status: EntityStatus;
}
```

### Default Values

```typescript
// ❌ WRONG
const DEFAULT_STATUS = 'active';

// ✅ CORRECT
const DEFAULT_STATUS = EntityStatus.Active;
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `05-features/06-ai-integration/08-ai-chat-ui.md` | `slot.status === 'active'` |
| `05-features/25-ai-enhancements/06-04-sharing-ui.md` | `share.status !== 'active'` |
| `05-features/25-ai-enhancements/06-01-sharing-architecture.md` | `share.Status != "active"` |
| `13-wp-plugin/05-wp-plugin-publish/02-frontend/28-remote-plugins.md` | `plugin.status === 'active'` |
| `07-database-design/03b-seed-data.md` | `Status: "active"` seed values |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/entitystatus` | `src/lib/enums/entity-status.ts` |
| Type | `byte` iota | String enum |
| Values | `Active`, `Inactive`, `Draft`, `Archived` | Same |

---

## Cross-References

- Issue #10 — Domain Status Magic Strings <!-- external: 02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./06-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent spec

---

*EntityStatus enum v1.0.0 — 2026-02-27*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-003: EntityStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All general entity lifecycle handling strictly utilizes `EntityStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
