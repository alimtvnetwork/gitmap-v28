# Generic Return Types — No interface{}/any/object Returns (AI Execution Prompt)

> **/goal** Eliminate loose, untyped return values (`interface{}`, `any`, `object`, `unknown`) across all function signatures, enforcing compile-time type safety via parametric generics and monadic Result wrappers.
> **/learn** Understand how untyped returns force downstream callers into error-prone runtime type assertions; master generic functions, Result wrappers, and reusable named type aliases across Go, TypeScript, C#, and Rust.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all `interface{}`, `any`, and `object` return signatures with generic types (`T`) or concrete domain types.
- [ ] `/learn` Never use unions or dynamic return types based on runtime flags; split into distinct, explicitly typed methods instead.
- [ ] `/goal` Create named type aliases (e.g. `type UserResult = apperror.Result[User]`) when generic wrappers appear more than once.
- [ ] `/learn` Verify zero untyped return violations and 100% type preservation across all packages via automated guideline linters.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Cross-Language Overview](./readme.md)
> **Version:** 1.0.0
> **Updated:** 2026-04-02
> **AI Confidence:** Production-Ready
> **Ambiguity:** None

## Keywords

`generic-return` · `result-type` · `type-safety` · `no-any` · `no-interface` · `method-design`

---

## Rule

**🔴 CODE RED:** When a method returns different types based on context, use generic Result types or generics — never `interface{}`, `any`, `object`, or `unknown`.

Returning untyped values forces callers to cast, which defeats compile-time safety and creates runtime panics. Generic wrappers preserve type information through the entire call chain.

---

## The Problem

```
// Every caller must guess the type and cast
value := cache.Get("user")        // returns interface{}
user := value.(User)              // runtime panic if wrong

result := service.Process(input)  // returns any
data := result.(OrderData)        // no compiler help
```

The compiler cannot verify correctness. Bugs surface at runtime, not build time.

---

## The Rule: Generic Return Types

### Go

```go
// ❌ FORBIDDEN — interface{} return
func (c *Cache) Get(key string) interface{} {
    return c.store[key]
}

// ❌ FORBIDDEN — any return (Go 1.18+)
func (s *Service) Process(input Input) any {
    if input.IsOrder {
        return processOrder(input)
    }
    return processRefund(input)
}

// ✅ REQUIRED — generic function
func Get[T any](c *Cache, key string) (T, bool) {
    val, ok := c.store[key]
    if !ok {
        var zero T
        return zero, false
    }
    return val.(T), true
}

// ✅ REQUIRED — Result wrapper (project pattern)
func (s *Service) ProcessOrder(input Input) apperror.Result[OrderData] {
    // returns typed Result — caller uses .Value() after .HasError() check
}

func (s *Service) ProcessRefund(input Input) apperror.Result[RefundData] {
    // separate method for different return type
}
```

### TypeScript

```typescript
// ❌ FORBIDDEN — any/unknown return
function fetchData(endpoint: string): Promise<any> {
    return axios.get(endpoint).then(r => r.data);
}

// ❌ FORBIDDEN — union that forces narrowing everywhere
function getItem(id: string): User | Order | Product { /* ... */ }

// ✅ REQUIRED — generic function
async function fetchData<T>(endpoint: string): Promise<T> {
    const response = await axios.get<T>(endpoint);
    return response.data;
}

// ✅ REQUIRED — separate typed methods
async function fetchUser(id: string): Promise<User> { /* ... */ }
async function fetchOrder(id: string): Promise<Order> { /* ... */ }
```

### C#

```csharp
// ❌ FORBIDDEN — object return
public object GetValue(string key) {
    return _store[key];
}

// ✅ REQUIRED — generic method
public T GetValue<T>(string key) {
    return (T)_store[key];
}

// ✅ REQUIRED — generic Result wrapper
public Result<T> Process<T>(Request request) where T : class {
    // typed result
}
```

### PHP

```php
// ❌ FORBIDDEN — mixed return
function getData(string $key): mixed {
    return $this->store[$key];
}

// ✅ REQUIRED — typed return with PHPDoc generics
/** @template T
 *  @param class-string<T> $type
 *  @return T */
function getData(string $key, string $type): object {
    $value = $this->store[$key];
    if (!$value instanceof $type) {
        throw new \InvalidArgumentException("Expected {$type}");
    }
    return $value;
}

// ✅ REQUIRED — separate typed methods
function getUser(string $id): User { /* ... */ }
function getOrder(string $id): Order { /* ... */ }
```

### Rust

```rust
// Rust enforces this at the language level — no untyped returns possible
// Use generics or enums with pattern matching

// ✅ Generic function
fn get_value<T: DeserializeOwned>(key: &str) -> Result<T, AppError> {
    let raw = store.get(key)?;
    serde_json::from_str(raw).map_err(|e| AppError::new("parse_failed", e))
}

// ✅ Enum for variant returns (compiler-enforced exhaustive matching)
enum ProcessResult {
    Order(OrderData),
    Refund(RefundData),
}
```

---

## Decision Guide

| Situation | Solution |
|-----------|----------|
| Same logic, different output type | Generic function `Fn[T]()` |
| Different logic per type | Separate named methods (see [24-boolean-flag-methods](./24-boolean-flag-methods.md)) |
| Known set of variant types | Enum/union with exhaustive matching |
| External API boundary | Deserialize into concrete type immediately at boundary |

---

## Best Practice: Concrete Type Aliases

When using generic types repeatedly, **create a named type alias** for each concrete instantiation. This eliminates repeated generic syntax, improves readability, and provides a single place to update if the underlying generic changes.

### Go

```go
// ❌ FORBIDDEN — Repeated generic syntax everywhere
func GetUser(ctx context.Context, id int64) apperror.Result[User] { ... }
func GetOrder(ctx context.Context, id int64) apperror.Result[Order] { ... }
func ListUsers(ctx context.Context) apperror.Result[[]User] { ... }

// ✅ REQUIRED — Concrete type aliases
type UserResult = apperror.Result[User]
type OrderResult = apperror.Result[Order]
type UserListResult = apperror.Result[[]User]

func GetUser(ctx context.Context, id int64) UserResult { ... }
func GetOrder(ctx context.Context, id int64) OrderResult { ... }
func ListUsers(ctx context.Context) UserListResult { ... }
```

### TypeScript

```typescript
// ❌ FORBIDDEN — Verbose generics repeated across the codebase
function fetchUser(id: string): Promise<ApiResponse<User>> { ... }
function fetchOrder(id: string): Promise<ApiResponse<Order>> { ... }

// ✅ REQUIRED — Named type aliases
type UserResponse = ApiResponse<User>;
type OrderResponse = ApiResponse<Order>;

function fetchUser(id: string): Promise<UserResponse> { ... }
function fetchOrder(id: string): Promise<OrderResponse> { ... }
```

### C# / PHP / Rust

```csharp
// C# — using alias (C# 12+)
using UserResult = Result<User>;
using OrderResult = Result<Order>;
```

```rust
// Rust — type alias
type UserResult = Result<User, AppError>;
type OrderResult = Result<Order, AppError>;
```

> **Rule of thumb:** If a generic instantiation appears more than once, create a named alias.

---

## Exemptions

| Case | Reason |
|------|--------|
| **Serialization boundaries** | `json.Unmarshal`, `sql.Scan` — cast at the boundary with `// EXEMPTED` annotation |
| **Plugin/extension systems** | Dynamic dispatch where types are unknown at compile time |
| **Reflection-based frameworks** | DI containers, ORM internals — exempted at framework boundary only |

---

## Cross-References

- [Strict Typing](./13-strict-typing.md) — all parameters and returns must be explicitly typed
- [Casting Elimination Patterns](./04-casting-elimination-patterns.md) — centralize casts at boundaries
- [Boolean Flag Methods](./24-boolean-flag-methods.md) — split methods instead of returning different types
- [AppError Result Types](../../03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/04-result-types.md) — Go Result[T] pattern

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TYPE-025: Generic Return Types and Elimination of Untyped Values

**Given** Function and method signatures returning dynamic or polymorphic values across polyglot codebases.
**When** Guidelines/linters audit the codebase for loose `interface{}`, `any`, `object`, or `unknown` returns.
**Then** All functions return strongly-typed generic structures or named Result aliases with zero downcasting requirements, achieving exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
