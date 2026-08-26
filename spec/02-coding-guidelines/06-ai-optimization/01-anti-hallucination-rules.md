# Anti-Hallucination Rules

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Purpose:** Explicit rules that prevent AI from generating non-compliant code

---

## How to Use

AI agents **MUST** check every generated code block against these rules before output. Each rule has:
- **Rule ID** — for referencing in reviews
- **❌ Forbidden** — what to never generate
- **✅ Required** — what to generate instead
- **📖 Why** — rationale for the rule

---

## Category 0: Repository Protection

### AH-00: No Generated Code or Artifacts

? **Never generate/commit:** Generated code (`*.generated.*`, ORM models, gRPC clients), test results, test reports, or compiled binaries.
? **Always generate:** `.gitignore` entries to exclude them from the repository.
?? Auto-generated code clutters the repository, creates diff noise, and should be produced dynamically at build time or in CI.

---

## Category 1: Naming Hallucinations

### AH-N1: Abbreviation Casing

❌ **Never generate:** `ID`, `URL`, `API`, `JSON`, `HTTP`, `HTML`, `SQL`, `DB`, `YAML`, `XML`, `CSS`  
✅ **Always generate:** `Id`, `Url`, `Api`, `Json`, `Http`, `Html`, `Sql`, `Db`, `Yaml`, `Xml`, `Css`  
📖 Abbreviations are treated as regular words — only first letter capitalized.

**Exemption:** Go standard library interfaces (`MarshalJSON`, `UnmarshalJSON`, `Error`, `String`) retain original spelling.

### AH-N2: Variable Casing

❌ **Never generate:** `plugin_slug`, `user_id`, `created_at` (snake_case)  
✅ **Always generate:** `pluginSlug`, `userId`, `createdAt` (camelCase)  
📖 Zero underscore policy for all logic-level identifiers.

### AH-N3: JSON/API Key Casing

❌ **Never generate:** `"userId"`, `"createdAt"`, `"plugin_slug"` (camelCase/snake_case)  
✅ **Always generate:** `"UserId"`, `"CreatedAt"`, `"PluginSlug"` (PascalCase)  
📖 All keys in our control use PascalCase.

### AH-N4: Go Getter Naming

❌ **Never generate:** `GetField()`, `GetName()`, `GetValue()`  
✅ **Always generate:** `Field()`, `Name()`, `Value()` (for getters); `SetField()` (for setters)  
📖 Go convention — getters don't use `Get` prefix.

### AH-N5: Enum Zero Value

❌ **Never generate:** `Unknown = iota`, `None = iota`, `Default = iota`  
✅ **Always generate:** `Invalid Variant = iota`  
📖 Zero value MUST always be `Invalid`.

### AH-N6: Boolean Naming

❌ **Never generate:** `active`, `loaded`, `ready`, `blocked` (no prefix)  
✅ **Always generate:** `isActive`, `isLoaded`, `isReady`, `isBlocked` (with prefix)  
📖 Every boolean must start with `is` or `has` (99%). Use `should` only for recommendations. Never use `can`, `was`, or `will`.

### AH-N7: Negative Boolean Names

❌ **Never generate:** `isNotReady`, `hasNoPermission`, `isNotBlocked`  
✅ **Always generate:** `isPending`, `isUnauthorized`, `isActive`  
📖 No negative words (`not`, `no`, `non`) in boolean names.

---

## Category 2: Type Safety Hallucinations

### AH-T1: No `any` / `interface{}`

❌ **Never generate:** `any`, `interface{}`, `map[string]any` in business logic  
✅ **Always generate:** Concrete typed structs  
📖 Type erasure defeats compile-time safety.

### AH-T2: No `fmt.Errorf()`

❌ **Never generate:** `fmt.Errorf("failed: %w", err)`  
✅ **Always generate:** `apperror.Wrap(err, apperror.ErrCode, "message")`  
📖 `fmt.Errorf` loses stack trace and structured error context.

### AH-T3: No Multi-Return Go Functions

❌ **Never generate:** `func F() (T, error)`, `func F() (T, bool, error)`  
✅ **Always generate:** `func F() apperror.Result[T]`  
📖 Single return value — Result pattern handles errors.

**Exemption:** Go's comma-ok pattern for map lookups and type assertions within `apperror` package.

### AH-T4: No Type Assertions in Business Logic

❌ **Never generate:** `msg["text"].(string)`, `cached.(float64)`, `payload.(*Type)`  
✅ **Always generate:** Concrete struct deserialization or `typecast.CastOrFail[T]()`  
📖 Type assertions panic at runtime and signal untyped data models.

### AH-T5: No Explicit Go JSON Tags

❌ **Never generate:** `json:"fieldName"`, `yaml:"fieldName"`  
✅ **Always generate:** No tag (PascalCase is default) or `json:",omitempty"` / `json:"-"`  
📖 Go JSON encoder uses PascalCase field names automatically.

### AH-T6: No Untyped Returns — Use Generics

❌ **Never generate:** `func F() interface{}`, `func F() any`, `object GetValue()`, `(): unknown`  
✅ **Always generate:** Generic return `func Get[T any]() T`, `T GetValue<T>()`, `<T>(): T`, or separate typed methods  
📖 Untyped returns force callers to cast, defeating compile-time safety. Use generics or Result[T] wrappers. See [25-generic-return-types.md](../01-cross-language/25-generic-return-types.md).

---

## Category 3: Boolean & Condition Hallucinations

### AH-B1: No Raw Negation on Function Calls

❌ **Never generate:** `if (!response.ok)`, `if !v.IsValid()`, `if (!file_exists($path))`  
✅ **Always generate:** `if (isResponseFailed(response))`, `if v.IsInvalid()`, `if (PathHelper::isFileMissing($path))`  
📖 Use semantic inverse methods/functions instead of `!`.

**Go exemptions:** `if !ok` (comma-ok), `if err != nil` (idiomatic), `if !strings.HasPrefix()` (stdlib).

### AH-B2: No Boolean Flag Parameters — Split Into Two Methods

❌ **Never generate:** `function process(data, true, false)`, `func Save(doc, isDraft bool)`  
✅ **Always generate:** Two named methods that express each intent explicitly:
- `processWithValidation(data)` / `processWithoutValidation(data)`
- `SaveDraft(doc)` / `PublishDocument(doc)`

📖 If a method branches on a boolean flag, split it into two methods. The method name must describe the behavior — no boolean parameter needed. Shared logic goes into private helpers. See [24-boolean-flag-methods.md](../01-cross-language/24-boolean-flag-methods.md).

**Exemptions:** Options/config structs with named fields, stdlib wrappers, toggle methods (`SetEnabled(bool)`).

### AH-B3: Max Two Conditions per Expression

❌ **Never generate:** `if (isA && isB && isC || isD)`  
✅ **Always generate:** Extract to named booleans first  
📖 3+ operands → extract to named boolean variables.

---

## Category 4: Structure Hallucinations

### AH-S1: Zero Nested `if`

❌ **Never generate:** `if (a) { if (b) { ... } }`  
✅ **Always generate:** Early returns to flatten: `if (!a) return; if (!b) return; ...`  
📖 Absolute ban on nested if statements.

### AH-S2: Max 15 Lines per Function Body

❌ **Never generate:** Function bodies exceeding 15 lines  
✅ **Always generate:** Decompose into helper functions  
📖 Error-handling lines are exempt from the count.

### AH-S3: Max 3 Parameters

❌ **Never generate:** `function f(a, b, c, d, e)`  
✅ **Always generate:** `function f(params: Params)` with typed options object  
📖 4+ parameters → use options/params object.

### AH-S4: Blank Line Before Return/Throw

❌ **Never generate:** Statement immediately followed by `return`/`throw`  
✅ **Always generate:** Blank line separating statements from `return`/`throw`  
📖 Exception: return/throw as sole statement needs no blank line.

---

## Category 5: Error Handling Hallucinations

### AH-E1: Guard Before Value Access

❌ **Never generate:** `result.Value()` or `$result->value()` without prior error check  
✅ **Always generate:** `if (result.HasError()) { ... }` before `.Value()`  
📖 Accessing value without guard may return zero/null silently.

### AH-E2: No Raw `error` in Struct Fields

❌ **Never generate (Go):** `Err error` in struct fields  
✅ **Always generate:** `AppError *apperror.AppError`  
📖 Raw `error` is not serializable — loses all diagnostic context.

### AH-E3: No Silent Error Swallowing

❌ **Never generate:** `_ = riskyOperation()`, empty catch blocks  
✅ **Always generate:** Log, return, or propagate every error  
📖 Every error must be explicitly handled.

### AH-E4: No Backslash PHP Imports

❌ **Never generate (PHP):** `\Throwable`, `\RuntimeException`, `\PDO`  
✅ **Always generate:** `use Throwable;` import at top, then `Throwable` unqualified  
📖 Consistency and readability.

---

## Category 6: Enum Hallucinations

### AH-EN1: Go Enum Type

❌ **Never generate:** `type Variant string`, `type Variant int`  
✅ **Always generate:** `type Variant byte`  
📖 `byte` for memory efficiency, type safety, and jump-table optimization.

### AH-EN2: PHP Enum Comparison

❌ **Never generate:** `$status === StatusType::Success`  
✅ **Always generate:** `$status->isEqual(StatusType::Success)`  
📖 Use comparison methods, not raw `===`.

### AH-EN3: No Magic Strings for Status

❌ **Never generate:** `if (status === 'active')`, `if provider == "serpapi"`  
✅ **Always generate:** `if (status === EntityStatus.Active)`, `if provider.IsSerpApi()`  
📖 All domain status/category comparisons via enum constants.

### AH-EN4: Single `variantLabels` Table (Go)

❌ **Never generate:** Dual tables (`variantStrings` + `variantLabels`)  
✅ **Always generate:** Single `variantLabels` array; `Label()` delegates to `String()`  
📖 One lookup table for all serialization, display, and comparison.

---

## Category 7: C# Hallucinations

### AH-CS1: Interface Naming

❌ **Never generate:** `interface UserRepository`, `interface Logger`  
✅ **Always generate:** `interface IUserRepository`, `interface ILogger`  
📖 C# interfaces must be prefixed with `I`.

### AH-CS2: No Blocking Async

❌ **Never generate:** `.Result`, `.GetAwaiter().GetResult()`, `.Wait()`  
✅ **Always generate:** `await` — async all the way through the call chain  
📖 Blocking async causes deadlocks in ASP.NET and UI contexts.

### AH-CS3: No Object Returns

❌ **Never generate:** `public object GetValue()`, `public dynamic Process()`  
✅ **Always generate:** `public T GetValue<T>()`, `public Result<T> Process<T>()`  
📖 Untyped returns defeat compile-time safety. Use generics.

### AH-CS4: Pattern Matching Over Casts

❌ **Never generate:** `var user = (User)obj;`, `var user = obj as User;`  
✅ **Always generate:** `if (obj is User user) { ... }` or `switch` expression with patterns  
📖 Pattern matching is null-safe and compiler-verified.

### AH-CS5: Records for Immutable DTOs

❌ **Never generate:** Mutable class with `{ get; set; }` for data transfer  
✅ **Always generate:** `public record UserDto(string Name, string Email);` or `{ get; init; }`  
📖 Records provide value equality, immutability, and concise syntax.

### AH-CS6: Sequential Independent Async

❌ **Never generate:** `var a = await GetA(); var b = await GetB();` (when A and B are independent)  
✅ **Always generate:** `var aTask = GetA(); var bTask = GetB(); await Task.WhenAll(aTask, bTask);`  
📖 Same as TypeScript Promise.all rule — parallel independent calls.

---

## Category 8: Caching — 🔴 CODE RED

### AH-CA1: Caching Errors as Success

❌ **Never generate:**
```typescript
try {
  const data = await fetchUsers();
  cache.set("users", data);
} catch {
  cache.set("users", []); // Silent failure cached as success
}
```
✅ **Always generate:**
```typescript
try {
  const data = await fetchUsers();
  cache.set("users", data, { ttl: 300_000 });
} catch (error) {
  cache.delete("users");
  logger.error("fetchUsers failed", { error });
  throw error;
}
```
📖 Caching empty/stale data on failure hides errors and serves broken state to all consumers.

### AH-CA2: Unbounded Cache (No TTL)

❌ **Never generate:** `cache.set(key, value)` without expiration  
✅ **Always generate:** `cache.set(key, value, { ttl: 300_000 })` — explicit TTL on every entry  
📖 Unbounded caches grow without limit and serve stale data indefinitely.

### AH-CA3: Missing Mutation Invalidation

❌ **Never generate:** Mutation (create/update/delete) without cache invalidation  
✅ **Always generate:** `cache.delete(key)` or `invalidateQueries()` immediately after mutation  
📖 Stale cache after mutation causes UI to show outdated data until TTL expires.

### AH-CA4: Non-Deterministic Cache Keys

❌ **Never generate:** `cache.set(\`users-${Date.now()}\`, data)` — timestamp/random in key  
✅ **Always generate:** `cache.set(\`users-${userId}-${role}\`, data)` — stable, typed inputs only  
📖 Non-deterministic keys prevent cache hits and cause unbounded key growth.

### AH-CA5: Untyped Cache Entries

❌ **Never generate:** `cache.set(key, data as any)` or `cache.get(key)` without type  
✅ **Always generate:** Typed cache wrapper — `cache.set<User[]>(key, users)` / `cache.get<User[]>(key)`  
📖 Untyped cache entries bypass type safety and cause runtime errors on consumption.

### AH-CA6: React Query Without Explicit Staleness

❌ **Never generate:** `useQuery({ queryKey, queryFn })` with no `staleTime`  
✅ **Always generate:** `useQuery({ queryKey, queryFn, staleTime: 5 * 60 * 1000 })` — explicit staleness  
📖 Default `staleTime: 0` causes unnecessary refetches on every mount.

---

## Quick Count: 40 Rules

| Category | Count |
|----------|-------|
| Naming (AH-N) | 7 |
| Type Safety (AH-T) | 6 |
| Boolean/Condition (AH-B) | 3 |
| Structure (AH-S) | 4 |
| Error Handling (AH-E) | 4 |
| Enum (AH-EN) | 4 |
| C# (AH-CS) | 6 |
| Caching (AH-CA) | 6 |
| **Total** | **40** |

---

## Cross-References

- [AI Quick Reference Checklist](./02-ai-quick-reference-checklist.md) — Condensed validation checklist
- [Common AI Mistakes](./03-common-ai-mistakes.md) — Real before/after examples
- [Master Coding Guidelines](../01-cross-language/15-master-coding-guidelines/00-overview.md) — Full rule reference
- [Condensed Master Guidelines — Caching](./04-condensed-master-guidelines.md) — Section 16

---

*Anti-hallucination rules v1.1.0 — 2026-04-04*
