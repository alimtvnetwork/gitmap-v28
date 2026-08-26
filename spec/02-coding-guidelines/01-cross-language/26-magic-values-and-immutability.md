# Magic Values, Immutability & Class-First Design

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Applies to:** All languages (Go, TypeScript, PHP, Rust, C#)  
**Source:** Consolidated from coding guidelines reviews, `18-code-mutation-avoidance.md`, and real-world `riseup-asia-uploader` patterns

---

## 1. Overview

Three interrelated principles that prevent the most common source of production bugs:

1. **No magic strings or magic numbers** — use enums or typed constants
2. **Immutable by default** — assign once, never mutate (`const` over `let`/`var`)
3. **Class-first design** (TypeScript/JS) — prefer classes over loose exported functions

---

## 2. The Dark Side of Magic Strings & Magic Numbers

### What Are They?

A **magic string** is a raw string literal used directly in logic (comparisons, switches, assignments) instead of a named constant or enum. A **magic number** is the same problem with numeric literals.

```typescript
// ❌ Magic string — what does "active" mean? Where is it defined?
if (user.status === "active") { ... }

// ❌ Magic number — what is 86400? Why 3? What happens when it changes?
setTimeout(callback, 86400)
if (retryCount > 3) { ... }
```

### Why They Are Dangerous

| Danger | Explanation | Real-World Impact |
|--------|-------------|-------------------|
| **Silent typos** | `"actve"` compiles fine but never matches | User locked out with no error — discovered weeks later |
| **No IDE support** | Can't "Find All References" or refactor | Changing `"active"` to `"enabled"` requires grepping the entire codebase |
| **No type safety** | Any string is accepted — the compiler can't help | Wrong status string passes silently through 5 layers of code |
| **Duplicated knowledge** | The same string appears in 20 files | One file gets updated, 19 don't — inconsistent behavior |
| **Impossible to test exhaustively** | String space is infinite — you can't match every case | `switch` without `default` silently skips unknown values |
| **Hidden coupling** | Two systems agree on `"webhook_completed"` via copy-paste | Backend renames to `"webhook_done"`, frontend breaks silently |
| **Localization nightmare** | `if (errorMsg === "Not found")` breaks in French | Entire feature fails for non-English users |
| **Security risk** | `if (role === "admin")` — attackers can guess the string | Privilege escalation by sending `"admin"` in a JWT claim |

### The Compounding Effect

Magic values don't just cause one bug — they **compound silently**:

```
Day 1:    Developer A writes `if (status === "active")`
Day 30:   Developer B copies it to 5 more files
Day 90:   Business says "active" → "enabled"
Day 91:   Developer C renames it in 3 of 6 files
Day 92:   3 features break silently in production
Day 120:  Customer reports data loss — traced to Day 91
```

With an enum, this entire sequence is impossible — renaming the enum value causes a **compile-time error** in all 6 files instantly.

---

## 3. Rules

### Rule 1: No Magic Strings — Use Enums or Constants

Every string literal used in comparisons, switch statements, or assignments **must** be replaced with a named constant or enum value.

#### TypeScript

```typescript
// ❌ FORBIDDEN — magic strings
function getDiscount(tier: string): number {
  if (tier === "premium") return 0.2
  if (tier === "basic") return 0.1

  return 0
}

// ✅ CORRECT — enum
enum CustomerTier {
  Basic = "basic",
  Premium = "premium",
}

function getDiscount(tier: CustomerTier): number {
  if (tier === CustomerTier.Premium) return 0.2
  if (tier === CustomerTier.Basic) return 0.1

  return 0
}
```

#### Go

```go
// ❌ FORBIDDEN — magic string
if site.Status == "active" {
    enablePlugins(site)
}

// ✅ CORRECT — byte enum with iota
type SiteStatus byte
const (
    SiteStatusActive SiteStatus = iota + 1
    SiteStatusBlocked
    SiteStatusPending
)

if site.Status == SiteStatusActive {
    enablePlugins(site)
}
```

#### PHP

```php
// ❌ FORBIDDEN — magic string
if ($order->status === 'completed') { ... }

// ✅ CORRECT — backed enum
enum OrderStatus: string {
    case Pending = 'pending';
    case Completed = 'completed';
    case Cancelled = 'cancelled';
}

if ($order->status->isEqual(OrderStatus::Completed)) { ... }
```

#### C#

```csharp
// ❌ FORBIDDEN — magic string
if (user.Role == "admin") { ... }

// ✅ CORRECT — enum
public enum UserRole { User, Moderator, Admin }

if (user.Role == UserRole.Admin) { ... }
```

#### Rust

```rust
// ❌ FORBIDDEN — magic string
if status == "active" { ... }

// ✅ CORRECT — enum
#[derive(PartialEq)]
enum Status { Active, Blocked, Pending }

if status == Status::Active { ... }
```

### Rule 2: No Magic Numbers — Name Every Constant

```typescript
// ❌ FORBIDDEN — what do these numbers mean?
setTimeout(retry, 86400000)
if (password.length < 8) { ... }
const tax = price * 0.2

// ✅ CORRECT — self-documenting
const ONE_DAY_MS = 86_400_000
const MIN_PASSWORD_LENGTH = 8
const VAT_RATE = 0.2

setTimeout(retry, ONE_DAY_MS)
if (password.length < MIN_PASSWORD_LENGTH) { ... }
const tax = price * VAT_RATE
```

**Exemptions:** `0`, `1`, `-1`, `""`, `true`, `false`, `null`/`nil` are exempt from this rule.

### Rule 3: Immutable by Default — Assign Once

Every variable should be assigned **exactly once**. Use `const` (TypeScript/JS/C#), `:=` without reassignment (Go), or `let` (Rust, which is immutable by default).

```typescript
// ❌ FORBIDDEN — mutable variable reassigned
let discount = 0
if (isPremium) {
  discount = 0.2
} else {
  discount = 0.1
}

// ✅ CORRECT — single assignment with ternary or function
const discount = isPremium ? 0.2 : 0.1
```

```go
// ❌ FORBIDDEN — variable mutated after creation
result := NewResult()
result.SetName("updated")
result.Status = StatusActive

return result

// ✅ CORRECT — construct with all values
return Result{
    Name:   computedName,
    Status: StatusActive,
}
```

```typescript
// ❌ FORBIDDEN — accumulating via mutation
let items: string[] = []
items.push("a")
items.push("b")

// ✅ CORRECT — declare the full value
const items: readonly string[] = ["a", "b"]
```

#### When Mutation Is Unavoidable

| Case | Why Allowed | Rule |
|------|-------------|------|
| Loop accumulation (`append`/`push`) | Collecting iteration results | Keep mutation in one method only |
| Builder pattern | Accumulates instructions, builds once | Final `.build()` call produces immutable value |
| Lazy/cached evaluation | Value generated once, then never changes | Use mutex lock for concurrent access |
| React `useState` | Framework requires it | State updates via setter only, never direct mutation |

### Rule 4: Class-First Design (TypeScript/JavaScript)

Prefer classes over loose exported functions when there is **shared state, configuration, or dependencies**.

```typescript
// ❌ DISCOURAGED — loose exported functions with hidden coupling
let apiUrl = ""

export function setApiUrl(url: string) {
  apiUrl = url
}

export function fetchUser(id: string) {
  return fetch(`${apiUrl}/users/${id}`)
}

export function fetchOrder(id: string) {
  return fetch(`${apiUrl}/orders/${id}`)
}
```

```typescript
// ✅ PREFERRED — class encapsulates state and dependencies
class ApiClient {
  constructor(private readonly baseUrl: string) {}

  fetchUser(id: string): Promise<Response> {
    return fetch(`${this.baseUrl}/users/${id}`)
  }

  fetchOrder(id: string): Promise<Response> {
    return fetch(`${this.baseUrl}/orders/${id}`)
  }
}

// Usage — immutable, testable, no global state
const api = new ApiClient("https://api.example.com")
```

**When functions are fine:** Pure utility functions with no shared state (e.g., `formatDate()`, `slugify()`) can remain as standalone exports.

---

## 4. Code Mutation: The Hidden Danger

### What Goes Wrong

```typescript
// ❌ DANGEROUS — object mutation across functions
function processOrder(order: Order): Order {
  applyDiscount(order)     // mutates order.price
  calculateTax(order)      // mutates order.tax (depends on mutated price)
  applyShipping(order)     // mutates order.shipping

  return order
}
// Bug: if applyDiscount is called twice, tax is calculated on wrong price
// Bug: reordering the calls produces different results
// Bug: testing any single function requires setting up the full object
```

```typescript
// ✅ SAFE — each step produces a new immutable value
function processOrder(order: Readonly<Order>): ProcessedOrder {
  const discountedPrice = calculateDiscountedPrice(order.price, order.tier)
  const tax = calculateTax(discountedPrice)
  const shipping = calculateShipping(order.weight)

  return { ...order, price: discountedPrice, tax, shipping }
}
// Each function is pure — testable in isolation, order doesn't matter
```

### Go Example — Mutation vs. Constructor

```go
// ❌ DANGEROUS — post-construction mutation
func buildResponse(data []byte) *Response {
    resp := &Response{}
    resp.Status = StatusOk
    resp.Body = data
    resp.Timestamp = time.Now()
    resp.Headers["Content-Type"] = "application/json"  // magic string!

    return resp
}

// ✅ SAFE — all values set at construction, no magic strings
const ContentTypeJSON = "application/json"

func buildResponse(data []byte) *Response {
    return &Response{
        Status:    StatusOk,
        Body:      data,
        Timestamp: time.Now(),
        Headers:   map[string]string{"Content-Type": ContentTypeJSON},
    }
}
```

---

## 5. Summary Checklist

```
□ No string literals in comparisons — use enum or const
□ No numeric literals in logic — name every constant
□ Every variable assigned exactly once (const over let/var)
□ No object mutation after construction — use constructors or struct literals
□ TS/JS: Prefer class over loose exports when state is shared
□ Mutation limited to one method — never across multiple function calls
□ Exemptions documented: 0, 1, -1, "", true, false, null/nil, loops, builders
```

---

## 6. Cross-References

- [Code Mutation Avoidance](./18-code-mutation-avoidance.md) — Detailed mutation rules and mutex patterns
- [Strict Typing](./13-strict-typing.md) — No `any`/`interface{}`/`object`
- [Boolean Principles](./02-boolean-principles.md) — Named booleans prevent `if (flag === true)`
- [Master Coding Guidelines](./15-master-coding-guidelines/00-overview.md) — §7 Type Safety
- [Generic Return Types](./25-generic-return-types.md) — Typed returns eliminate `any`

---

*Magic values, immutability & class-first design — cross-language specification.*
