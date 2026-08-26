# C# Naming and Conventions

> **Parent:** [C# Coding Standards](./00-overview.md)  
> **Version:** 1.0.0  
> **Updated:** 2026-04-02

---

## Identifier Casing

| Item | Convention | Example |
|------|-----------|---------|
| Classes, Structs, Records | PascalCase | `SnapshotManager`, `UserProfile` |
| Interfaces | `I` + PascalCase | `IUserRepository`, `ILogger` |
| Methods | PascalCase | `ProcessUpload()`, `GetActiveUsers()` |
| Properties | PascalCase | `IsActive`, `PluginSlug` |
| Local variables | camelCase | `pluginSlug`, `userId` |
| Parameters | camelCase | `siteId`, `userName` |
| Constants | PascalCase | `MaxRetryCount`, `DefaultTimeout` |
| Private fields | `_` + camelCase | `_logger`, `_connectionString` |
| Enums | PascalCase, no `Type` suffix | `Status`, `HttpMethod` |
| Enum values | PascalCase | `Active`, `Pending`, `Invalid` |

---

## Abbreviation Casing

Follow the cross-language rule — abbreviations are treated as regular words:

| ❌ Forbidden | ✅ Required |
|-------------|-----------|
| `UserID` | `UserId` |
| `GetURL` | `GetUrl` |
| `APIClient` | `ApiClient` |
| `ParseJSON` | `ParseJson` |
| `HTTPMethod` | `HttpMethod` |

**Exemption:** Two-letter abbreviations in .NET BCL (`IO`, `DB`) follow Microsoft convention when used in framework-level code. In business logic, prefer `Id`, `Db`.

---

## Boolean Naming

Every boolean property and variable must use a prefix:

```csharp
// ❌ BAD
public bool Active { get; set; }
public bool Loaded { get; set; }
var ready = CheckStatus();

// ✅ GOOD
public bool IsActive { get; set; }
public bool IsLoaded { get; set; }
var isReady = CheckStatus();
```

Allowed prefixes: `Is`, `Has`, `Can`, `Should`, `Was`.

No negative names: `IsNotReady` → `IsPending`, `HasNoPermission` → `IsUnauthorized`.

---

## File Naming

| Type | Convention | Example |
|------|-----------|---------|
| Classes | `{PascalCase}.cs` | `SnapshotManager.cs` |
| Interfaces | `I{PascalCase}.cs` | `IUserRepository.cs` |
| Records | `{PascalCase}.cs` | `UserProfile.cs` |
| One type per file | Always | Match file name to primary type |

---

## Namespace Conventions

```csharp
// ❌ BAD — flat or inconsistent
namespace App;
namespace app.services;

// ✅ GOOD — matches folder structure, PascalCase
namespace RiseupAsia.Services;
namespace RiseupAsia.Domain.Models;
```

---

## Cross-References

- [Cross-Language Abbreviation Casing](../01-cross-language/04-code-style/00-overview.md)
- [Boolean Principles](../01-cross-language/02-boolean-principles/00-overview.md)
- [Variable Naming Conventions](../01-cross-language/22-variable-naming-conventions.md)

---
