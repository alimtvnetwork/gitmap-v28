# C# Method Design (AI Execution Prompt)

> **/goal** Architect clean, maintainable C# methods by eliminating boolean flag parameters, enforcing strict line count limits, applying pure async patterns, and utilizing readable LINQ.
> **/learn** Split methods branching on booleans into distinct intention-revealing operations, cap methods at 15 lines and 3 parameters, avoid blocking async calls, and extract complex LINQ expressions.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Eliminate boolean flag parameters that branch control flow by splitting into dedicated methods (`SaveDraft` vs `PublishDocument`).
- [ ] `/learn` Cap method bodies at 15 lines (excluding error handling) and restrict parameter counts to at most 3 (use options classes for 4+).
- [ ] `/goal` Ban blocking calls (`.Result`, `.GetAwaiter().GetResult()`) on async tasks; enforce async/await throughout.
- [ ] `/learn` Use `Task.WhenAll` for independent concurrent async operations rather than sequential awaits.
- [ ] `/goal` Suffix all asynchronous methods with `Async`.
- [ ] `/learn` Prefer LINQ over imperative loops for collections, and extract complex or nested predicates into named private methods.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [C# Coding Standards](./readme.md)
> **Version:** 1.0.0
> **Updated:** 2026-04-02

---

## Boolean Flag Splitting

🔴 **CODE RED:** If a method branches on a boolean parameter, split it into two named methods.

```csharp
// ❌ FORBIDDEN — boolean flag hides intent
public void SaveDocument(Document doc, bool isDraft)
{
    if (isDraft) { /* draft logic */ }
    else { /* publish logic */ }
}

// Caller: SaveDocument(doc, true)  — What does true mean?

// ✅ REQUIRED — two methods, intent is obvious
public void SaveDraft(Document doc)
{
    // draft logic
}

public void PublishDocument(Document doc)
{
    // publish logic
}
```

When both paths share setup/teardown, extract into private helpers:

```csharp
public void SaveDraft(Document doc)
{
    ValidateDocument(doc);     // shared
    StoreDraft(doc);           // unique
    NotifyAuthor(doc);         // shared
}

public void PublishDocument(Document doc)
{
    ValidateDocument(doc);     // shared
    StorePublished(doc);       // unique
    NotifyAuthor(doc);         // shared
}

private void ValidateDocument(Document doc) { /* ... */ }
private void NotifyAuthor(Document doc) { /* ... */ }
```

**Exemptions:** Options objects with named properties, toggle methods (`SetEnabled(bool)`).

> **Full rule:** [24-boolean-flag-methods.md](../01-cross-language/24-boolean-flag-methods.md)

---

## Function Size

- **Max 15 lines** per method body (error handling exempt)
- **Max 3 parameters** — use an options class for 4+
- **Single responsibility** — one method does one thing

```csharp
// ❌ FORBIDDEN — too many params
public void CreateUser(string name, string email, string role, bool isActive, int age)

// ✅ REQUIRED — options class
public void CreateUser(CreateUserOptions options)

public class CreateUserOptions
{
    public string Name { get; init; }
    public string Email { get; init; }
    public string Role { get; init; }
    public bool IsActive { get; init; }
    public int Age { get; init; }
}
```

---

## Async Patterns

```csharp
// ❌ FORBIDDEN — blocking async
var result = GetDataAsync().Result;
var data = GetDataAsync().GetAwaiter().GetResult();

// ✅ REQUIRED — async all the way
var result = await GetDataAsync();

// ❌ FORBIDDEN — sequential independent calls
var users = await GetUsersAsync();
var orders = await GetOrdersAsync();

// ✅ REQUIRED — parallel independent calls
var usersTask = GetUsersAsync();
var ordersTask = GetOrdersAsync();
await Task.WhenAll(usersTask, ordersTask);
var users = usersTask.Result;
var orders = ordersTask.Result;
```

**Naming:** Async methods must end with `Async` suffix: `GetUsersAsync()`, `SaveDocumentAsync()`.

---

## LINQ Usage

```csharp
// ❌ FORBIDDEN — manual loops for simple transforms
var names = new List<string>();
foreach (var user in users)
{
    names.Add(user.Name);
}

// ✅ REQUIRED — LINQ
var names = users.Select(u => u.Name).ToList();

// ❌ FORBIDDEN — nested LINQ (hard to read)
var result = items.Where(x => x.Orders.Any(o => o.Items.Any(i => i.Price > 100)));

// ✅ REQUIRED — extract to named method
var result = items.Where(HasExpensiveOrderItem);

private static bool HasExpensiveOrderItem(Item item)
{
    return item.Orders.Any(o => o.Items.Any(i => i.Price > 100));
}
```

---

## Cross-References

- [Boolean Flag Methods](../01-cross-language/24-boolean-flag-methods.md) — cross-language rule with C# examples
- [Cyclomatic Complexity](../01-cross-language/06-cyclomatic-complexity.md) — max complexity rules
- [Nesting Resolution](../01-cross-language/20-nesting-resolution-patterns.md) — flatten nested conditions

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CS-003: C# Method Signatures, Parameter Limits and Guard Clauses

**Given** C# methods, constructors, and async implementations.
**When** Method structures and signatures are analyzed against design guidelines.
**Then** Zero boolean flag parameters branch method logic, method bodies remain within 15 lines, parameter lists are capped at 3 or refactored into options objects, async methods use non-blocking patterns with `Async` suffix, and complex LINQ operations are cleanly factored out with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.
