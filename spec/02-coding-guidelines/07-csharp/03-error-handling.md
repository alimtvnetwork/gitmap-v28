# C# Error Handling

> **Parent:** [C# Coding Standards](./00-overview.md)  
> **Version:** 1.0.0  
> **Updated:** 2026-04-02

---

## Exception Guidelines

### Catch Specific Exceptions

```csharp
// ❌ BAD — catching everything
try { /* ... */ }
catch (Exception ex) { Log(ex); }

// ✅ GOOD — catch specific, rethrow unknown
try { /* ... */ }
catch (HttpRequestException ex) { HandleNetworkError(ex); }
catch (JsonException ex) { HandleParseError(ex); }
```

### Never Swallow Exceptions

```csharp
// ❌ BAD — silent swallow
try { Process(); }
catch { }

// ✅ GOOD — log + handle or rethrow
try { Process(); }
catch (InvalidOperationException ex)
{
    _logger.LogError(ex, "Processing failed for {Id}", itemId);

    throw;
}
```

---

## Guard Clauses

Use early returns instead of nested `if`:

```csharp
// ❌ BAD — nested
public void ProcessOrder(Order order)
{
    if (order != null)
    {
        if (order.IsValid)
        {
            // process
        }
    }
}

// ✅ GOOD — guard clauses
public void ProcessOrder(Order order)
{
    if (order is null)
        throw new ArgumentNullException(nameof(order));

    if (!order.IsValid)

        return;

    // process
}
```

---

## Nullable Reference Types

Enable nullable reference types project-wide and use null guards:

```csharp
// ❌ BAD — unchecked null
string name = user.Name; // could be null

// ✅ GOOD — explicit null handling
string name = user.Name ?? throw new InvalidOperationException("Name is required");

// ✅ GOOD — nullable annotation
public string? GetMiddleName(User user)
{
    return user.MiddleName;
}
```

---

## Cross-References

- [Null Pointer Safety](../01-cross-language/19-null-pointer-safety.md) — cross-language null safety
- [Nesting Resolution Patterns](../01-cross-language/20-nesting-resolution-patterns.md) — guard clause patterns

---
