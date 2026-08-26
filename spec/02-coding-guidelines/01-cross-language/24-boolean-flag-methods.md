# Boolean Flag Method Splitting

> **Parent:** [Cross-Language Overview](./00-overview.md)  
> **Version:** 1.0.0  
> **Updated:** 2026-04-02  
> **AI Confidence:** Production-Ready  
> **Ambiguity:** None

## Keywords

`boolean-flag` · `method-splitting` · `single-responsibility` · `function-design` · `clean-code`

---

## Rule

**🔴 CODE RED:** If a method's behavior changes based on a boolean parameter, split it into two named methods that express each intent explicitly.

Boolean flags hide branching logic inside function calls. The caller cannot understand what `true` or `false` means without reading the implementation. Two named methods make intent obvious at every call site.

---

## The Problem

```
// ❌ BAD — What does `true` mean here?
processOrder(order, true)
processOrder(order, false)

sendEmail(user, true)
sendEmail(user, false)
```

The reader must open the function to understand the flag. This violates self-documenting code principles and makes code reviews slower.

---

## The Rule: Split Into Two Methods

Every boolean flag parameter that changes method behavior must be replaced with two methods whose names describe the behavior.

### Go

```go
// ❌ BAD — boolean flag hides intent
func ProcessOrder(order Order, isPriority bool) error {
    if isPriority {
        // priority logic
    } else {
        // standard logic
    }
}

// ✅ GOOD — two methods, intent is clear
func ProcessPriorityOrder(order Order) error {
    // priority logic
}

func ProcessStandardOrder(order Order) error {
    // standard logic
}
```

### TypeScript

```typescript
// ❌ BAD
function formatUser(user: User, isDetailed: boolean): string {
    if (isDetailed) {
        return `${user.name} (${user.email}, ${user.role})`;
    }
    return user.name;
}

// ✅ GOOD
function formatUserSummary(user: User): string {
    return user.name;
}

function formatUserDetailed(user: User): string {
    return `${user.name} (${user.email}, ${user.role})`;
}
```

### PHP

```php
// ❌ BAD
function syncPlugin(Plugin $plugin, bool $isForced): void {
    if ($isForced) {
        // force sync logic
    } else {
        // incremental sync logic
    }
}

// ✅ GOOD
function syncPluginIncremental(Plugin $plugin): void {
    // incremental sync logic
}

function syncPluginForced(Plugin $plugin): void {
    // force sync logic
}
```

### Rust

```rust
// ❌ BAD
fn write_log(entry: &LogEntry, is_verbose: bool) {
    if is_verbose {
        // verbose output
    } else {
        // compact output
    }
}

// ✅ GOOD
fn write_log_compact(entry: &LogEntry) {
    // compact output
}

fn write_log_verbose(entry: &LogEntry) {
    // verbose output
}
```

### C#

```csharp
// ❌ BAD
public void SaveDocument(Document doc, bool isDraft)
{
    if (isDraft) { /* draft logic */ }
    else { /* publish logic */ }
}

// ✅ GOOD
public void SaveDraft(Document doc)
{
    // draft logic
}

public void PublishDocument(Document doc)
{
    // publish logic
}
```

---

## When Shared Logic Exists

If both paths share setup or teardown, extract the shared logic into a private helper:

```go
// ✅ Shared logic extracted
func ProcessPriorityOrder(order Order) error {
    validateOrder(order)        // shared
    applyPriorityDiscount(order) // unique
    return finalizeOrder(order)  // shared
}

func ProcessStandardOrder(order Order) error {
    validateOrder(order)         // shared
    applyStandardPricing(order)  // unique
    return finalizeOrder(order)  // shared
}

// Private shared helpers
func validateOrder(order Order) { /* ... */ }
func finalizeOrder(order Order) error { /* ... */ }
```

---

## Exemptions

| Case | Reason |
|------|--------|
| **Options/config structs** | Booleans inside an options struct are acceptable — caller sees named fields (`Config{Verbose: true}`) |
| **Standard library wrappers** | Thin wrappers around stdlib that pass through bool params (e.g., `os.OpenFile` flags) |
| **Toggle methods** | Methods that flip state (`SetEnabled(bool)`) where the name already describes intent |

---

## Cross-References

- [Boolean Principles](./02-boolean-principles/00-overview.md) — P5: No boolean parameters
- [Function Naming](./10-function-naming.md) — naming conventions for split methods
- [Cyclomatic Complexity](./06-cyclomatic-complexity.md) — flag splitting reduces branching
- [SOLID Principles](./23-solid-principles.md) — Single Responsibility applied to methods
- [Nesting Resolution](./20-nesting-resolution-patterns.md) — related pattern: flatten `if/else`

---
