# C# Error Handling (AI Execution Prompt)

> **/goal** Implement resilient, structured error handling in C# through specific exception catches, non-swallowed exceptions, early guard clauses, and nullable reference type enforcement.
> **/learn** Master catching domain-specific exceptions, logging before rethrowing with `throw;`, flattening nested conditionals using null and validity guard clauses, and leveraging C# nullable annotations.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Catch specific exception types (`HttpRequestException`, `JsonException`) and ban bare `catch (Exception)`.
- [ ] `/learn` Prohibit empty catch blocks or silently swallowed exceptions; always log with context and rethrow using `throw;`.
- [ ] `/goal` Implement early return guard clauses with `ArgumentNullException.ThrowIfNull` or `if (param is null) throw new ArgumentNullException(nameof(param))`.
- [ ] `/learn` Flatten nested `if` statements into sequential guard clauses to maintain cyclomatic complexity ≤5.
- [ ] `/goal` Enable and enforce nullable reference types (`string?`) with null-coalescing throw operators (`?? throw`).
- [ ] `/learn` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [C# Coding Standards](./readme.md)
> **Version:** 1.0.0
> **Updated:** 2026-04-02

---

## Exception Guidelines

### Catch Specific Exceptions

```csharp
// ❌ FORBIDDEN — catching everything
try { /* ... */ }
catch (Exception ex) { Log(ex); }

// ✅ REQUIRED — catch specific, rethrow unknown
try { /* ... */ }
catch (HttpRequestException ex) { HandleNetworkError(ex); }
catch (JsonException ex) { HandleParseError(ex); }
```

### Never Swallow Exceptions

```csharp
// ❌ FORBIDDEN — silent swallow
try { Process(); }
catch { }

// ✅ REQUIRED — log + handle or rethrow
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
// ❌ FORBIDDEN — nested
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

// ✅ REQUIRED — guard clauses
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
// ❌ FORBIDDEN — unchecked null
string name = user.Name; // could be null

// ✅ REQUIRED — explicit null handling
string name = user.Name ?? throw new InvalidOperationException("Name is required");

// ✅ REQUIRED — nullable annotation
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

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CS-004: C# Structured Exceptions, Specific Catches and Result Types

**Given** C# error handling logic, catch blocks, and validation routines.
**When** Codebases are audited for exception safety and guard clause usage.
**Then** Catch blocks target specific exception types without silent swallows, parameters are validated using guard clauses and `nameof()`, nullable reference types are strictly annotated and handled, and nested branching is flattened with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.
