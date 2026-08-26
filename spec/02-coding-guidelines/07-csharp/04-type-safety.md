# C# Type Safety

> **Parent:** [C# Coding Standards](./00-overview.md)  
> **Version:** 1.0.0  
> **Updated:** 2026-04-02

---

## Generics Over Object

```csharp
// ❌ BAD — loses type safety
public object GetValue(string key) { /* ... */ }
public void SetValue(string key, object value) { /* ... */ }

// ✅ GOOD — generic
public T GetValue<T>(string key) { /* ... */ }
public void SetValue<T>(string key, T value) { /* ... */ }
```

---

## Pattern Matching Over Type Casting

```csharp
// ❌ BAD — explicit cast can throw
var user = (User)obj;

// ❌ BAD — as + null check
var user = obj as User;
if (user != null) { /* ... */ }

// ✅ GOOD — pattern matching
if (obj is User user)
{
    // use user
}

// ✅ GOOD — switch expression
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
// ❌ BAD — mutable class for data transfer
public class UserDto
{
    public string Name { get; set; }
    public string Email { get; set; }
}

// ✅ GOOD — record for immutable data
public record UserDto(string Name, string Email);

// ✅ GOOD — record with init-only props for complex cases
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
// ❌ BAD — magic strings
if (status == "active") { /* ... */ }
var role = "admin";

// ✅ GOOD — enum or constants
if (status == StatusType.Active) { /* ... */ }
var role = RoleType.Admin;
```

---

## Cross-References

- [Strict Typing](../01-cross-language/13-strict-typing.md) — cross-language type safety rules
- [Casting Elimination](../01-cross-language/03-casting-elimination-patterns.md) — avoid type casts
- [Code Mutation Avoidance](../01-cross-language/18-code-mutation-avoidance.md) — immutability patterns

---
