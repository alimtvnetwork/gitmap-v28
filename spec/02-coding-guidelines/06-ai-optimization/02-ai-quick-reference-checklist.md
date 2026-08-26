# AI Quick Reference Checklist

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Purpose:** Condensed, machine-parsable checklist for AI to validate generated code in <30 seconds

---

## Pre-Output Validation

Check every generated code block against these rules before outputting.

---

### Critical Repository Rules (1 check)

- [ ] NEVER output generated code files (`*.generated.*`, ORM models, etc.) or test artifacts for Git check-in. They belong in `.gitignore`.

---

### Naming (5 checks)

- [ ] Variables: camelCase (`pluginSlug`, not `plugin_slug`)
- [ ] Types/Classes/Components: PascalCase (`SnapshotManager`)
- [ ] JSON/API keys: PascalCase (`"PluginSlug"`, not `"pluginSlug"`)
- [ ] Abbreviations: First-letter-only (`Id`, `Url`, `Api`, `Http`, `Json`)
- [ ] Booleans: `is`/`has` prefix (99%), `should` rare, no `can`/`was`/`will`, no negative words

### Structure (7 checks)

- [ ] Zero nested `if` — use early return pattern
- [ ] Function body ≤15 lines (error handling lines exempt)
- [ ] Function params ≤3 (use options object for 4+)
- [ ] No boolean flag params — split into two named methods (e.g., `SaveDraft()` / `PublishDocument()`)
- [ ] No `any`/`interface{}`/`object`/`unknown` returns — use generics or typed Result wrappers
- [ ] File size ≤300 lines
- [ ] Blank line before `return`/`throw` when preceded by statements

### Go-Specific (8 checks)

- [ ] File names: PascalCase, no underscores, named after primary struct/definition (`PluginService.go`)
- [ ] Enums: `byte` type, `Invalid` zero value, `iota`
- [ ] Returns: single `Result[T]`, never `(T, error)`
- [ ] Errors: `apperror.Wrap()`, never `fmt.Errorf()`
- [ ] No `any`/`interface{}` in business logic
- [ ] No type assertions — use concrete structs or `typecast.CastOrFail[T]()`
- [ ] No explicit `json:` tags (only `,omitempty` or `-`)
- [ ] Getters: `Field()` not `GetField()`; max 1 `defer` per function

### PHP-Specific (4 checks)

- [ ] No `\Throwable` — use `use` import
- [ ] Enum comparison via `isEqual()`, not `===`
- [ ] Blank line before `if` when preceded by statements
- [ ] Log keys: camelCase; DB keys: PascalCase

### TypeScript-Specific (5 checks)

- [ ] No `any` — use explicit types
- [ ] Enum: PascalCase + `Type` suffix (`StatusType`)
- [ ] No magic strings — use enum constants
- [ ] `isDefined()` / `isDefinedAndValid()` — no raw null checks
- [ ] 🔴 **Promise.all for independent async calls** — sequential `await` on independent promises is auto-rejection

### Error Handling (4 checks)

- [ ] `HasError()`/`hasError()` before `.Value()`/`.value()`
- [ ] No silent error swallowing — log, return, or propagate
- [ ] Struct error fields: `*AppError` (Go), `Throwable` (PHP)
- [ ] No `_ = riskyOperation()` — handle every error

### Database (3 checks)

- [ ] Table names: PascalCase (`UserProfiles`)
- [ ] Column names: PascalCase (`PluginSlug`, `CreatedAt`)
- [ ] Index names: `Idx` prefix + PascalCase (`IdxTransactions_CreatedAt`)

### Mutation & Immutability (3 checks)

- [ ] Variables assigned once — no reassignment after init
- [ ] No post-construction mutation — pass all values to constructor
- [ ] Concurrent mutation uses mutex locks (provide locked + unlocked versions)

### Null Pointer Safety (4 checks)

- [ ] Check `err`/error before accessing value — always
- [ ] Never chain method calls on unchecked returns
- [ ] Nil-check pointer before dereference
- [ ] Check nil AND `len()` before slice/array index access

### Lazy Evaluation (3 checks)

- [ ] Lazy fields: non-exported field + public getter method — never direct access
- [ ] If dependency is lazy, dependent field MUST also be lazy
- [ ] Concurrent lazy access → mutex lock wrapper

### Regex (3 checks)

- [ ] Regex is last resort — prefer `strings.Contains()`, `HasPrefix()`, `Split()`
- [ ] Go: compile at package level `var re = regexp.MustCompile(...)` — never inside functions
- [ ] Never use regex in loops without reviewer approval

### Nesting & Newlines (3 checks)

- [ ] Zero nested `if` — extract to function, inverse logic, or named booleans
- [ ] No blank line after opening brace; no double blank lines anywhere
- [ ] Use `constants.NewLineUnix` (`"\n"`) — not raw `"\r\n"`

### Boolean Complexity (3 checks)

- [ ] Max 2 operands per boolean expression — decompose 3+ into named variables
- [ ] Never mix `&&` and `||` in one expression — split into named booleans
- [ ] Never mix negative (`!`, `== nil`) and positive checks — use early return for negatives

### Variable Naming (3 checks)

- [ ] Singular for single items, plural for collections (`user` vs `users`)
- [ ] Loop variable = singular of collection (`for user of users`)
- [ ] Maps use `Map` suffix or `By[Key]` pattern (`usersById`, `configMap`)

### SOLID Principles (3 checks)

- [ ] Single Responsibility — one reason to change per class/module/function
- [ ] Interface Segregation — no client forced to depend on methods it doesn't use
- [ ] Dependency Inversion — depend on abstractions, not concretions

### Struct & File Limits (2 checks)

- [ ] Struct/class ≤10 fields — split into composed sub-structs if exceeded
- [ ] No `else` after `return`/`throw`/`break`/`continue` — code after block is implicit else

### Defer — Go (2 checks)

- [ ] Max one `defer` per function
- [ ] Multiple defers needed → extract into separate single-defer functions

### C#-Specific (6 checks)

- [ ] Interfaces prefixed with `I` (`IUserRepository`, not `UserRepository`)
- [ ] Private fields use `_camelCase` (`_logger`, `_connectionString`)
- [ ] No `.Result` or `.GetAwaiter().GetResult()` — async all the way
- [ ] Independent async calls use `Task.WhenAll()`, not sequential `await`
- [ ] Pattern matching over type casts (`if (obj is User user)`, not `(User)obj`)
- [ ] Records for immutable DTOs (`record UserDto(string Name, string Email)`)

### 🔴 Caching (6 checks — CODE RED)

- [ ] Never cache errors as success — `catch` blocks must `cache.delete()`, not `cache.set()`
- [ ] Every `cache.set()` has explicit TTL — unbounded caches are prohibited
- [ ] Create/update/delete mutations immediately invalidate related cache entries
- [ ] Cache keys use deterministic, stable inputs — no `Date.now()` or random values
- [ ] Cache entries are typed — no `any` or untyped objects in cache values
- [ ] React Query: explicit `staleTime` set (not default `0`) + `invalidateQueries` after mutations

---

## Total: 78 checks

**Pass criteria:** All 78 checks must pass before outputting code.

---

## Cross-References

- [Anti-Hallucination Rules](./01-anti-hallucination-rules.md) — Detailed forbidden/required patterns
- [Common AI Mistakes](./03-common-ai-mistakes.md) — Real mistake examples
- [Master Coding Guidelines](../01-cross-language/15-master-coding-guidelines/00-overview.md) — Full checklist at bottom
- [Condensed Master Guidelines](./04-condensed-master-guidelines.md) — Sub-200-line AI context version

---

*AI quick reference checklist v2.1.0 — 2026-03-31*
