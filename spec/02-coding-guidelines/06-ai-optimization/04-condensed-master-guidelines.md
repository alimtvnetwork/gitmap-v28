# Condensed Master Coding Guidelines — AI Context Reference

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Last synced with master:** 2026-03-31  
**Purpose:** Sub-200-line distillation of [15-master-coding-guidelines.md](../01-cross-language/15-master-coding-guidelines/00-overview.md) for AI context windows. See also: [05-enum-naming-quick-reference.md](./05-enum-naming-quick-reference.md) for detailed enum rules.

---

## 0. Rule Zero - NEVER COMMIT GENERATED CODE OR ARTIFACTS

**NEVER** commit generated code (e.g., ORM models, gRPC clients, OpenAPI stubs), test results, test reports, or compiled binaries. They must be added to `.gitignore`. **This is the most critical repository rule.**

---

## 1. Naming

| Element | Convention | Examples |
|---------|-----------|----------|
| Class/Struct/Enum type | PascalCase + `Type` suffix (enums) | `SnapshotManager`, `StatusType` |
| Method | camelCase (PHP/TS) / PascalCase (Go exported) | `processUpload()` / `ProcessUpload()` |
| Variable | camelCase | `pluginSlug`, `postId` |
| Boolean | `is`/`has` + camelCase (99%), `should` rare | `isActive`, `hasPermission` |
| Source file | PascalCase (matches primary type) | `SnapshotManager.go`, `UserProfile.tsx` |
| DB tables/columns | PascalCase | `Transactions`, `PluginSlug`, `CreatedAt` |
| JSON/API keys | PascalCase | `"PluginSlug"`, `"SiteId"` |
| Go package dir | snake_case | `site_health/` |
| **Rust** fn/var/mod | snake_case (RFC 430 override) | `process_upload()`, `plugin_slug` |
| **Rust** struct/enum/trait | PascalCase | `SnapshotManager`, `StatusType` |
| **Rust** source file | snake_case | `snapshot_manager.rs` |

**Abbreviations:** First letter only caps — `Id` not `ID`, `Url` not `URL`, `Api` not `API`, `Json` not `JSON`.

**Zero underscore policy:** snake_case banned for all logic-level identifiers. **Rust exception:** Rust uses snake_case for functions, variables, and modules per community standards (RFC 430). Exempt also: WP hooks, DB migrations, PHP superglobals.

> **🔴 AI CRITICAL — Rust Hybrid Naming:** Rust identifiers follow snake_case, BUT database identifiers (tables, columns, views) and enum string values (serialized to JSON/DB/API) MUST always be PascalCase. Use `#[serde(rename_all = "PascalCase")]` for serialization and `#[sqlx(rename = "PascalCase")]` for DB mappings.

---

## 2. Booleans — Positive Logic

| Principle | Rule |
|-----------|------|
| P1 Prefix | Every boolean: `is`/`has` (99%), `should` only for recommendations |
| P2 No negatives | `not`/`no`/`non` banned → use positive synonyms (`isPending` not `isNotReady`) |
| P3 Named guards | Never `!fn()` → use semantic inverse (`isInvalid()` not `!isValid()`) |
| P4 Extract complex | 2+ operators → named boolean variable |
| P5 No bool params | Use separate methods or options objects. See [24-boolean-flag-methods.md](../01-cross-language/24-boolean-flag-methods.md) |
| P6 No mixed polarity | `isX && !isY` → single-intent name |

**Go exemptions:** `if !ok`, `if err != nil`, `if !strings.HasPrefix(...)`, handler guards.

**Null guards:** Use `isDefined()`/`isDefinedAndValid()` — never raw `!= null/nil`.

---

## 3. Enums

> **Full cross-language enum reference:** [05-enum-naming-quick-reference.md](./05-enum-naming-quick-reference.md)

| Language | Key Rules |
|----------|-----------|
| **Go** | `byte` type, `Invalid = iota` zero value, package-scoped constants (`ProductionStr` not `EnvironmentProductionStr`), `variantLabels` PascalCase. Methods: `String()`, `Label()`, `IsValid()`, `Is{Value}()`, `Parse()`, `MarshalJSON()`, `UnmarshalJSON()`. Package name = grouping (`environmenttype.Production`). |
| **TypeScript** | String enum, PascalCase members, UPPER_SNAKE values. No string unions. File: `src/lib/enums/{name}.ts`. |
| **PHP** | `Type` suffix, string-backed, `isEqual()`/`isOtherThan()`/`isAnyOf()` required. Parse via `::from()`/`::tryFrom()` — never manual match. Never `===`. |
| **All** | Enum constants for all comparisons — never raw string literals. Exhaustive `default` branch in every switch. |

---

## 4. Code Style

| Rule | Description |
|------|-------------|
| R1 | Always use braces — no single-line `if` |
| R2/R7 | **Zero nested `if`** — absolute ban |
| R3 | Extract complex conditions → named booleans |
| R4 | Blank line before `return` when preceded by statements |
| R5 | Blank line after `}` when followed by more code |
| R6 | Max **15 lines** per function body |
| R9 | >2 params/args → one per line |
| R12 | No empty line after opening brace |
| R13 | No empty line at start of file; no double blank lines |

**Nesting resolution:** (1) Extract to named function, (2) Inverse logic / early return, (3) Named boolean variables. Exception: one loop with single `if` inside is OK.

---

## 5. Error Handling

**PHP:** `try/catch` with `Throwable` (imported, no leading `\`).
**Go:** `apperror.New()`/`apperror.Wrap()` only — never `fmt.Errorf()`. Service methods return `apperror.Result[T]` — never `(T, error)`.

**Result guard rule:** Always check `HasError()`/`hasError()` or `IsSafe()`/`isSafe()` before accessing `.Value()`/`.value()`. No silent failures.

**Struct error fields:** Use `*AppError` (Go) / `Throwable` (PHP) — never raw `error` or `string`.

---

## 6. Type Safety

**Go:** Zero `any`/`interface{}`/`map[string]any` in business logic. Zero type assertions — use concrete structs. Unavoidable casts → `typecast.CastOrFail[T]()`.

**Single return rule (Go):** Every function returns exactly ONE value (`Result[T]` or typed struct). Never `(T, bool, error)`.

**PHP:** Native type declarations on all params/returns/properties. Max 3 params. Remove redundant PHPDoc.

---

## 6b. Async / Promise Patterns (TypeScript) — 🔴 CODE RED

**Independent promises MUST run in parallel.** Sequential `await` on independent calls is an **automatic rejection**.

```typescript
// ❌ CODE RED — sequential on independent calls
const users = await fetchUsers();
const posts = await fetchPosts();

// ✅ REQUIRED — Promise.all
const [users, posts] = await Promise.all([fetchUsers(), fetchPosts()]);
```

**Decision:** Does call B need the result of call A? → sequential. Otherwise → `Promise.all`. Mixed deps → parallel the independent ones, then sequential for dependent.

**React Query:** Multiple `useQuery` hooks auto-parallelize — no action needed.

---

## 7. Magic Strings — Zero Tolerance

All repeated strings → enums or typed constants. Hook names, capabilities, table names, error codes, HTTP methods, log levels, status values — all via enum `.value` / `.String()`.

**Comparison helpers** (`hasMismatch`, `isEqual`) must always receive enum constants — never string literals.

---

## 8. Mutation Avoidance

- Assign variables **once** — prefer immutable
- No post-construction mutation — pass all values to constructor
- Concurrent mutation → mutex locks (provide locked + unlocked versions)
- **Exemptions:** lazy eval/caching, loop `append`, builder pattern

---

## 9. Lazy Evaluation

- Use for expensive/static fields accessed by many callers
- Non-exported field + getter method — never direct access
- If dependency is lazy, dependent field MUST also be lazy
- Concurrent access → mutex lock wrapper
- Do NOT use when data varies per request or is cheap

---

## 10. Null Pointer Safety

- Check `err` before value — always
- Never chain method calls on unchecked returns
- Nil-check pointer before dereference
- Check nil AND `len()` before slice index access
- Nil-check before creating pointer returns

---

## 11. Regex

- **Last resort** — prefer `strings.Contains()`, `HasPrefix()`, `Split()`
- Go: compile at package level `var re = regexp.MustCompile(...)` — never in functions
- Add sample data as comments above regex vars
- Never use regex in loops without reviewer approval

---

## 12. Defer (Go)

- Max **one** `defer` per function
- Place at top or bottom — never buried in middle
- Multiple defers needed → extract into separate single-defer functions

---

## 13. Newlines

- Blank before `return` in multi-line functions (not single-line)
- No blank between consecutive `}`
- No blank at function start; no double blanks
- Use `constants.NewLineUnix` (`"\n"`) in 90% of cases

---

## 14. File Organization

**Go:** 300-line file target (400 hard limit). 15-line function max. Split: `_crud.go`, `_helpers.go`. Import order: stdlib → internal → third-party.
**PHP:** PSR-4, one class/enum per file, traits declare own `use` imports.

---

## 15. Tests

Three-part naming: `Test{Unit}_{Scenario}_{ExpectedOutcome}`. Table-driven for 3+ cases. AAA pattern (Arrange/Act/Assert). Each test independently runnable. Go: `t.Helper()` required in helpers.

---

## 16. Caching — 🔴 CODE RED

| Rule | Description |
|------|-------------|
| C1 | **Never cache errors as success** — empty `catch { cache.set(key, []) }` is silent failure = automatic rejection |
| C2 | **Always set TTL** — unbounded caches serve stale data and grow without limit |
| C3 | **Invalidate on mutation** — after create/update/delete, invalidate or update cache immediately |
| C4 | **Deterministic keys** — stable inputs only, no timestamps or random values |
| C5 | **Typed entries** — never `cache.set(key, any)` |
| C6 | **React Query** — explicit `staleTime`, use `invalidateQueries()` after mutations |

**On failure:** invalidate/delete the cache entry — never store empty or stale values.

---

## Quick Checklist

```
[ ] Naming: camelCase vars, PascalCase classes/enums/DB, Id not ID
[ ] JSON/API keys: PascalCase
[ ] Null guards: isDefined() — never raw != null/nil
[ ] Null safety: err before value, nil before dereference, len before index
[ ] Booleans: is/has prefix, no negatives, no raw !
[ ] Enums: Type suffix, isEqual() not ===, PascalCase cases
[ ] DB: PascalCase tables/columns
[ ] Style: braces, zero nesting, blank before return, 15-line max
[ ] Errors: apperror.Wrap (Go), Throwable (PHP), no fmt.Errorf
[ ] Results: hasError()/isSafe() before .value()/.Value()
[ ] Single return: Go → Result[T] or typed struct only
[ ] No casting: concrete structs, typecast.CastOrFail[T]
[ ] No magic strings: enums/typed constants everywhere
[ ] Mutation: assign once, no post-construction mutation
[ ] Regex: last resort, package-level compile, no loops
[ ] Lazy eval: non-exported + getter, cascade lazy deps
[ ] Defer: max one per function (Go)
[ ] 🔴 Promise.all: independent async calls MUST be parallel — never sequential await
[ ] 🔴 Caching: never cache errors as success, always TTL, invalidate on mutation
[ ] Tests: three-part naming, AAA, table-driven 3+ cases
```

---

*Condensed master guidelines v1.2.0 — AI context optimized — 2026-04-04*
