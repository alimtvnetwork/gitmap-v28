# C# Coding Standards — Acceptance Criteria

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

---

## AC-01: Naming & Conventions

- [ ] All classes, structs, records use PascalCase (`SnapshotManager`, `UserProfile`)
- [ ] Interfaces prefixed with `I` (`IUserRepository`, `ILogger`)
- [ ] Methods use PascalCase (`ProcessUpload()`, `GetActiveUsers()`)
- [ ] Local variables and parameters use camelCase (`pluginSlug`, `userId`)
- [ ] Private fields use `_camelCase` (`_logger`, `_connectionString`)
- [ ] Abbreviations: first-letter-only caps (`UserId` not `UserID`, `ApiClient` not `APIClient`)
- [ ] All booleans prefixed with `Is`/`Has`/`Can`/`Should`/`Was` (`IsActive`, `HasPermission`)
- [ ] No negative boolean names (`IsPending` not `IsNotReady`)
- [ ] File names match primary type in PascalCase (`SnapshotManager.cs`)
- [ ] One type per file

## AC-02: Method Design

- [ ] No boolean flag parameters that branch method behavior — split into two named methods
- [ ] Method bodies ≤15 lines (error handling exempt)
- [ ] Method parameters ≤3 — use options class for 4+
- [ ] Shared logic between split methods extracted into private helpers
- [ ] Async methods suffixed with `Async` (`GetUsersAsync()`, `SaveDocumentAsync()`)

## AC-03: Boolean Flag Splitting

- [ ] `Save(doc, isDraft)` → `SaveDraft(doc)` + `PublishDocument(doc)`
- [ ] `Process(data, isVerbose)` → `ProcessCompact(data)` + `ProcessVerbose(data)`
- [ ] Options objects with named bool properties are exempt (`Config { Verbose = true }`)
- [ ] Toggle methods are exempt (`SetEnabled(bool)`)

## AC-04: Error Handling

- [ ] Catch specific exceptions, never bare `catch (Exception)`
- [ ] No silently swallowed exceptions — log and rethrow or handle explicitly
- [ ] Guard clauses with early return instead of nested `if` blocks
- [ ] `ArgumentNullException` with `nameof()` for null parameter guards
- [ ] Nullable reference types enabled project-wide (`<Nullable>enable</Nullable>`)

## AC-05: Type Safety

- [ ] No `object` returns — use generics (`T GetValue<T>(string key)`)
- [ ] No explicit type casts in business logic — use pattern matching (`if (obj is User user)`)
- [ ] Records used for immutable data transfer objects
- [ ] No magic strings — use enums or typed constants
- [ ] `switch` expressions with exhaustive matching and `_` default case

## AC-06: Async Patterns

- [ ] No `.Result` or `.GetAwaiter().GetResult()` — async all the way
- [ ] Independent async calls use `Task.WhenAll()` not sequential `await`
- [ ] Async method naming ends with `Async` suffix

## AC-07: LINQ Usage

- [ ] LINQ preferred over manual loops for transforms (`Select`, `Where`, `Any`)
- [ ] Complex LINQ predicates extracted to named methods
- [ ] No nested LINQ deeper than 2 levels

---

## Validation

All criteria are testable via code review or static analysis. Each AC maps to rules in:
- [01-naming-and-conventions.md](./01-naming-and-conventions.md)
- [02-method-design.md](./02-method-design.md)
- [03-error-handling.md](./03-error-handling.md)
- [04-type-safety.md](./04-type-safety.md)
- [Boolean Flag Methods (cross-language)](../01-cross-language/24-boolean-flag-methods.md)
- [Generic Return Types (cross-language)](../01-cross-language/25-generic-return-types.md)

---
