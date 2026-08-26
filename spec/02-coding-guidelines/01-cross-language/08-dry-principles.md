# DRY Principles — Coding Guidelines

> **Version:** 1.0.0  
> **Updated:** 2026-03-09  
> **Applies to:** All code (Go, TypeScript, PHP, PowerShell)

---

## Core Principle

**Don't Repeat Yourself (DRY):** Every piece of knowledge must have a single, unambiguous, authoritative representation within a system.

---

## Rules

### 1. Single Source of Truth

Every concept, constant, type, or configuration must be defined in exactly one place.

| What | Where | Anti-pattern |
|------|-------|--------------|
| API endpoint paths | `constants.php` (PHP), route registration (Go) | Hardcoded strings in handlers |
| Response envelope types | `envelope.schema.json` → per-language types | Independent type definitions that drift |
| Error codes | Centralized constants file | Magic numbers in business logic |
| Feature flags | Config file or database | Scattered `if` checks with hardcoded booleans |

### 2. Extract, Don't Copy

When you find yourself writing the same logic a second time, **extract it immediately**.

```
❌ Copy-paste a function and modify slightly
✅ Extract shared logic into a helper, parameterize differences
```

**Extraction targets:**
- **3+ lines** of identical logic → Extract to function
- **2+ components** sharing state logic → Extract to custom hook (React) or service (Go)
- **2+ endpoints** sharing validation → Extract to middleware or shared validator

### 3. Composition Over Inheritance

Build complex behavior by composing small, focused units:

```typescript
// ✅ Composition: useApiQuery wraps useQuery + error handling
const { data } = useApiQuery<Plugin[]>('/plugins', ['plugins']);

// ❌ Inheritance: extending a BaseComponent for each entity
class PluginsPage extends DataFetchingComponent { ... }
```

### 4. Factory Patterns for Repetitive Structures

When multiple entities follow the same pattern, use factories:

```typescript
// ✅ Factory hook — one definition, many consumers
export function useApiQuery<T>(endpoint: string, queryKey: string[]) { ... }

// ❌ Per-entity hooks with duplicated useQuery + error handling
export function usePlugins() { const { data } = useQuery(...); /* error handling */ }
export function useSites() { const { data } = useQuery(...); /* same error handling */ }
```

### 5. Modular Decomposition

No single file should exceed **300 lines**. When it does, decompose:

| Before | After |
|--------|-------|
| `GlobalErrorModal.tsx` (800 lines) | `modal/BackendSection.tsx`, `modal/FrontendSection.tsx`, `modal/TraversalDetails.tsx`, etc. |
| `api.ts` (600 lines) | `api/types.ts`, `api/envelope.ts`, `api/client.ts`, `api/methods.ts` |

### 6. Cross-Stack Contracts

When multiple languages implement the same structure, use a **machine-readable schema** as the source of truth:

```
envelope.schema.json (JSON Schema Draft 2020-12)
    ├── Go:  envelope.go (structs with omitempty)
    ├── TS:  types.ts (interfaces with ? optionals)
    └── PHP: EnvelopeBuilder.php (builder methods)
```

Each implementation references the schema version in a comment.

---

## How to Write DRY Code

### Before Writing New Code

1. **Search first** — Does this logic already exist somewhere?
2. **Check for patterns** — Are there similar implementations for other entities?
3. **Identify the abstraction** — What's the reusable part vs. the specific part?

### During Code Review

Ask these questions:
- "Have I seen this pattern before in the codebase?"
- "Would a new team member know where to find this logic?"
- "If this requirement changes, how many files would I need to update?"

### Red Flags

| Smell | Likely Violation |
|-------|-----------------|
| Same error handling in 3+ places | Missing centralized error handler |
| Same API call shape in 3+ hooks | Missing factory hook |
| Same constant string in 2+ files | Missing constants file |
| Same validation logic in 2+ handlers | Missing middleware/validator |
| Copy-pasting a component and tweaking props | Missing variant or composition |

---

## Project-Specific DRY Patterns

| Pattern | Implementation | Location |
|---------|---------------|----------|
| API data fetching | `useApiQuery` / `useApiQueryPaginated` | `src/hooks/useApiQuery.ts` |
| Error reporting | Global error store | `src/stores/errorStore.ts` |
| Envelope parsing | `parseEnvelope` utility | `src/lib/api/envelope.ts` |
| PHP snapshot creation | `SnapshotFactory` | `includes/Snapshot/SnapshotFactory.php` |
| PHP logging context | Centralized enrichment in logger | `includes/Logging/FileLogger.php` |
| Go namespace resolution | `resolveNamespace` helper | `internal/wordpress/namespace.go` |
| Go PHP stack parsing | `php_stack.go` helper | `internal/wordpress/php_stack.go` |

---

## Cross-References

- [DRY Refactoring Summary](./09-dry-refactoring-summary.md) — Complete 10-phase history
- [Response Envelope Schema](../../03-error-manage/02-error-architecture/05-response-envelope/envelope.schema.json) — Cross-stack type contract
- [TypeScript Standards](../02-typescript/08-typescript-standards-reference.md) — TS-specific rules
- [Golang Standards](../03-golang/04-golang-standards-reference/00-overview.md) — Go-specific rules
- [PHP Standards](../04-php/07-php-standards-reference/00-overview.md) — PHP-specific rules

---

*DRY principles specification created: 2026-02-09*
