# C# Type Safety (AI Execution Prompt)

> **/goal** Maximize C# type safety and runtime reliability through generics, pattern matching over casting, immutable records, and elimination of magic strings.
> **/learn** Replace `object` returns with generic type parameters, leverage `is` type pattern matching and `switch` expressions, utilize `record` types for data transfer, and enforce typed enums.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace untyped `object` returns and parameters with strongly-typed generics (`T GetValue<T>(string key)`).
- [ ] `/learn` Eliminate explicit type casts (`(User)obj` and `as User`) in favor of pattern matching (`if (obj is User user)`).
- [ ] `/goal` Use switch expressions with exhaustive pattern matching and wildcard default discards (`_ => throw`).
- [ ] `/learn` Use `record` or `record struct` declarations for immutable data transfer objects (DTOs) and domain events.
- [ ] `/goal` Ban magic strings in business logic and state checks, replacing them with typed enums or constants.
- [ ] `/learn` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [C# Coding Standards](./readme.md)
> **Version:** 1.0.0
> **Updated:** 2026-04-02

---

## Generics Over Object

```csharp
// ❌ FORBIDDEN — loses type safety
public object GetValue(string key) { /* ... */ }
public void SetValue(string key, object value) { /* ... */ }

// ✅ REQUIRED — generic
public T GetValue<T>(string key) { /* ... */ }
public void SetValue<T>(string key, T value) { /* ... */ }
```

---

## Pattern Matching Over Type Casting

```csharp
// ❌ FORBIDDEN — explicit cast can throw
var user = (User)obj;

// ❌ FORBIDDEN — as + null check
var user = obj as User;
if (user != null) { /* ... */ }

// ✅ REQUIRED — pattern matching
if (obj is User user)
{
    // use user
}

// ✅ REQUIRED — switch expression
var result = shape switch
{
    Circle c => Math.PI * c.Radius * c.Radius,
    Rectangle r => r.Width * r.Height,
    _ => throw new InvalidOperationException($"Unknown shape: {shape.GetType().Name}")
};
```

---

## Records for Immutable Data

```csharp
// ❌ FORBIDDEN — mutable class for data transfer
public class UserDto
{
    public string Name { get; set; }
    public string Email { get; set; }
}

// ✅ REQUIRED — record for immutable data
public record UserDto(string Name, string Email);

// ✅ REQUIRED — record with init-only props for complex cases
public record OrderDto
{
    public string OrderId { get; init; }
    public decimal Total { get; init; }
    public IReadOnlyList<LineItem> Items { get; init; }
}
```

---

## No Magic Strings

```csharp
// ❌ FORBIDDEN — magic strings
if (status == "active") { /* ... */ }
var role = "admin";

// ✅ REQUIRED — enum or constants
if (status == StatusType.Active) { /* ... */ }
var role = RoleType.Admin;
```

---

## Cross-References

- [Strict Typing](../01-cross-language/13-strict-typing.md) — cross-language type safety rules
- [Casting Elimination](../01-cross-language/04-casting-elimination-patterns.md) — avoid type casts
- [Code Mutation Avoidance](../01-cross-language/18-code-mutation-avoidance.md) — immutability patterns

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CS-005: C# Nullable Reference Types and Pattern Matching

**Given** C# type declarations, casting logic, and data structures.
**When** Source code is audited for type safety and pattern matching standards.
**Then** Generics are utilized in place of `object`, casting is eliminated in favor of pattern matching `is` expressions and `switch` expressions, DTOs are declared as immutable records, magic strings are replaced by typed enums, and nullable reference types are honored with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.
