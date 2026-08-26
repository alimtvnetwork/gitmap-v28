# C# Method Design

> **Parent:** [C# Coding Standards](./00-overview.md)  
> **Version:** 1.0.0  
> **Updated:** 2026-04-02

---

## Boolean Flag Splitting

🔴 **CODE RED:** If a method branches on a boolean parameter, split it into two named methods.

```csharp
// ❌ BAD — boolean flag hides intent
public void SaveDocument(Document doc, bool isDraft)
{
    if (isDraft) { /* draft logic */ }
    else { /* publish logic */ }
}

// Caller: SaveDocument(doc, true)  — What does true mean?

// ✅ GOOD — two methods, intent is obvious
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
// ❌ BAD — too many params
public void CreateUser(string name, string email, string role, bool isActive, int age)

// ✅ GOOD — options class
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
// ❌ BAD — blocking async
var result = GetDataAsync().Result;
var data = GetDataAsync().GetAwaiter().GetResult();

// ✅ GOOD — async all the way
var result = await GetDataAsync();

// ❌ BAD — sequential independent calls
var users = await GetUsersAsync();
var orders = await GetOrdersAsync();

// ✅ GOOD — parallel independent calls
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
// ❌ BAD — manual loops for simple transforms
var names = new List<string>();
foreach (var user in users)
{
    names.Add(user.Name);
}

// ✅ GOOD — LINQ
var names = users.Select(u => u.Name).ToList();

// ❌ BAD — nested LINQ (hard to read)
var result = items.Where(x => x.Orders.Any(o => o.Items.Any(i => i.Price > 100)));

// ✅ GOOD — extract to named method
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
