# Generic Return Types — No interface{}/any/object Returns

> **Parent:** [Cross-Language Overview](./00-overview.md)  
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
// ❌ BAD — interface{} return
func (c *Cache) Get(key string) interface{} {
    return c.store[key]
}

// ❌ BAD — any return (Go 1.18+)
func (s *Service) Process(input Input) any {
    if input.IsOrder {
        return processOrder(input)
    }
    return processRefund(input)
}

// ✅ GOOD — generic function
func Get[T any](c *Cache, key string) (T, bool) {
    val, ok := c.store[key]
    if !ok {
        var zero T
        return zero, false
    }
    return val.(T), true
}

// ✅ GOOD — Result wrapper (project pattern)
func (s *Service) ProcessOrder(input Input) apperror.Result[OrderData] {
    // returns typed Result — caller uses .Value() after .HasError() check
}

func (s *Service) ProcessRefund(input Input) apperror.Result[RefundData] {
    // separate method for different return type
}
```

### TypeScript

```typescript
// ❌ BAD — any/unknown return
function fetchData(endpoint: string): Promise<any> {
    return axios.get(endpoint).then(r => r.data);
}

// ❌ BAD — union that forces narrowing everywhere
function getItem(id: string): User | Order | Product { /* ... */ }

// ✅ GOOD — generic function
async function fetchData<T>(endpoint: string): Promise<T> {
    const response = await axios.get<T>(endpoint);
    return response.data;
}

// ✅ GOOD — separate typed methods
async function fetchUser(id: string): Promise<User> { /* ... */ }
async function fetchOrder(id: string): Promise<Order> { /* ... */ }
```

### C#

```csharp
// ❌ BAD — object return
public object GetValue(string key) {
    return _store[key];
}

// ✅ GOOD — generic method
public T GetValue<T>(string key) {
    return (T)_store[key];
}

// ✅ GOOD — generic Result wrapper
public Result<T> Process<T>(Request request) where T : class {
    // typed result
}
```

### PHP

```php
// ❌ BAD — mixed return
function getData(string $key): mixed {
    return $this->store[$key];
}

// ✅ GOOD — typed return with PHPDoc generics
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

// ✅ GOOD — separate typed methods
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
// ❌ BAD — Repeated generic syntax everywhere
func GetUser(ctx context.Context, id int64) apperror.Result[User] { ... }
func GetOrder(ctx context.Context, id int64) apperror.Result[Order] { ... }
func ListUsers(ctx context.Context) apperror.Result[[]User] { ... }

// ✅ GOOD — Concrete type aliases
type UserResult = apperror.Result[User]
type OrderResult = apperror.Result[Order]
type UserListResult = apperror.Result[[]User]

func GetUser(ctx context.Context, id int64) UserResult { ... }
func GetOrder(ctx context.Context, id int64) OrderResult { ... }
func ListUsers(ctx context.Context) UserListResult { ... }
```

### TypeScript

```typescript
// ❌ BAD — Verbose generics repeated across the codebase
function fetchUser(id: string): Promise<ApiResponse<User>> { ... }
function fetchOrder(id: string): Promise<ApiResponse<Order>> { ... }

// ✅ GOOD — Named type aliases
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
- [Casting Elimination Patterns](./03-casting-elimination-patterns.md) — centralize casts at boundaries
- [Boolean Flag Methods](./24-boolean-flag-methods.md) — split methods instead of returning different types
- [AppError Result Types](../../03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/03-result-types.md) — Go Result[T] pattern

---
