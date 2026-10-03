# TypeScript ExecutionStatus Enum — `src/lib/enums/execution-status.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string execution status literals (`'running'`, `'idle'`, `'paused'`, `'completed'`, `'failed'`, `'canceled'`) across pipeline, automation, and frontend specs by standardizing on `ExecutionStatus`.
> **/learn** Master execution lifecycle typing: `export enum ExecutionStatus` with uppercase string values, terminal state helper sets, and typed execution result interfaces.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw string execution state comparisons (`status === 'running'`) with `ExecutionStatus.Running`.
- [ ] `/learn` Never use loose union types (`'pending' | 'running' | 'success' | 'failed'`) for execution status in models or interfaces; use `ExecutionStatus`.
- [ ] `/goal` Use `TERMINAL_STATES` set lookup with `ExecutionStatus` constants for terminal execution checks.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 1.0.0
> **Last updated**: 2026-02-27
> **Tracks**: Issue #10 (`02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

---

## Purpose

Typed enum for all execution lifecycle states — pipeline runs, instruction execution, build processes, code generation chains, and any task that transitions through idle → running → terminal states. Eliminates `status === 'running'` magic strings across frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/execution-status.ts

export enum ExecutionStatus {
  Idle = "IDLE",
  Running = "RUNNING",
  Paused = "PAUSED",
  Completed = "COMPLETED",
  Failed = "FAILED",
  canceled = "canceled",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (execution.status === 'running') { ... }

// ✅ CORRECT: Enum constant
if (execution.status === ExecutionStatus.Running) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{chain.status !== 'running' && <ChainSummary chain={chain} />}

// ✅ CORRECT
{chain.status !== ExecutionStatus.Running && <ChainSummary chain={chain} />}
```

### Type Definitions

```typescript
// ❌ WRONG: Union of magic strings
interface ExecutionResult {
  status: 'pending' | 'running' | 'success' | 'failed';
}

// ✅ CORRECT: Enum-typed
interface ExecutionResult {
  status: ExecutionStatus;
}
```

### Terminal State Helpers

```typescript
const TERMINAL_STATES = new Set([
  ExecutionStatus.Completed,
  ExecutionStatus.Failed,
  ExecutionStatus.canceled,
]);

function isTerminal(status: ExecutionStatus): boolean {
  return TERMINAL_STATES.has(status);
}
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `05-features/06-ai-integration/08-ai-chat-ui.md` | `slot.status === 'loading'/'idle'` |
| `05-features/25-ai-enhancements/05-03-message-display.md` | `execution.status === 'success'/'failed'` |
| `05-features/09-knowledge-memory/11-knowledge-memory-ui.md` | `job.status === 'running'/'completed'/'failed'` |
| `05-features/27-automation-pipeline/10-react-flow-canvas.md` | `executionState?.status === 'RUNNING'` |
| `05-features/25-ai-enhancements/03-02-plan-execution.md` | Plan execution status checks |
| `20-shared-cli-frontend/15-hooks-library.md` | `Status === 'running'/'error'` |
| `09-gsearch-cli/02-frontend/05-ui-patterns.md` | Execution status checks |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/executionstatus` | `src/lib/enums/execution-status.ts` |
| Type | `byte` iota | String enum |
| Values | `Idle`, `Running`, `Paused`, `Completed`, `Failed`, `canceled` | Same |

---

## Cross-References

- Issue #10 — Domain Status Magic Strings <!-- external: 02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->
- [HttpMethod Enum](./06-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent spec
- [Master Coding Guidelines §8](../01-cross-language/15-master-coding-guidelines/readme.md) — Magic strings zero tolerance

---

*ExecutionStatus enum v1.0.0 — 2026-02-27*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-004: ExecutionStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All task and pipeline execution lifecycle states strictly utilize `ExecutionStatus` enum constants, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
