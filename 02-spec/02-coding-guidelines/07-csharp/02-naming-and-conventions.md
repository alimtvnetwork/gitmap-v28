# C# Naming and Conventions (AI Execution Prompt)

> **/goal** Enforce strict C# identifier naming standards, PascalCase casing for types and properties, camelCase for parameters and local variables, and positive boolean prefixes.
> **/learn** Master .NET naming conventions, abbreviation rules treating acronyms as words (`UserId`, `GetUrl`), one type per file in PascalCase, and eliminating negative boolean identifiers.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce PascalCase for classes, structs, records, methods, properties, and constants.
- [ ] `/learn` Prefix interfaces with `I` and prefix private instance fields with `_camelCase`.
- [ ] `/goal` Apply camelCase to method parameters and local variables without Hungarian notation or leading underscores.
- [ ] `/learn` Treat acronyms as words in PascalCase/camelCase (`UserId`, `GetUrl`, `ApiClient`, `ParseJson`).
- [ ] `/goal` Require positive boolean prefixes (`Is`, `Has`, `Can`, `Should`, `Was`) and eliminate negative boolean naming.
- [ ] `/learn` Ensure file names match their primary declared type in `{PascalCase}.cs` with strictly one primary type per file.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [C# Coding Standards](./readme.md)
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

| ❌ FORBIDDEN | ✅ REQUIRED |
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
// ❌ FORBIDDEN
public bool Active { get; set; }
public bool Loaded { get; set; }
var ready = CheckStatus();

// ✅ REQUIRED
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
// ❌ FORBIDDEN — flat or inconsistent
namespace App;
namespace app.services;

// ✅ REQUIRED — matches folder structure, PascalCase
namespace RiseupAsia.Services;
namespace RiseupAsia.Domain.Models;
```

---

## Cross-References

- [Cross-Language Abbreviation Casing](../01-cross-language/04-code-style/readme.md)
- [Boolean Principles](../01-cross-language/02-boolean-principles/readme.md)
- [Variable Naming Conventions](../01-cross-language/22-variable-naming-conventions.md)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CS-002: C# PascalCase/camelCase and Positive Booleans

**Given** C# source files containing type declarations, properties, methods, fields, and variables.
**When** Codebases are audited against C# naming conventions.
**Then** All classes, structs, records, interfaces, methods, and properties follow PascalCase, locals and parameters follow camelCase, abbreviations use word casing (`UserId`, `GetUrl`), booleans use affirmative prefixes with zero negative names, and file names match primary types with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.
