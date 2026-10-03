# C# Coding Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all C# coding guideline specifications in `07-csharp/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-CS-[NUM]`), PascalCase conventions, boolean flag splitting, async patterns, exception handling, pattern matching, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `07-csharp/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths across all guideline links.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. C# Criteria Inventory (`AC-CG-CS-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-CS-001` | C# Coding Guidelines Index & Module Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| `AC-CG-CS-002` | C# PascalCase/camelCase and Positive Booleans | [`02-naming-and-conventions.md`](02-naming-and-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| `AC-CG-CS-003` | C# Method Signatures, Parameter Limits and Guard Clauses | [`03-method-design.md`](03-method-design.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| `AC-CG-CS-004` | C# Structured Exceptions, Specific Catches and Result Types | [`04-error-handling.md`](04-error-handling.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| `AC-CG-CS-005` | C# Nullable Reference Types and Pattern Matching | [`05-type-safety.md`](05-type-safety.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| `AC-CG-CS-REG-001` | C# Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-CS-001: C# Coding Guidelines Index & Module Conformance

- [ ] Root `readme.md` provides navigation, scoring table, document inventory, and cross-references.
- [ ] All files in `07-csharp/` adhere to active prompt anatomy, checklist headers, and acceptance criteria blocks.
- [ ] Internal relative markdown links resolve accurately with zero broken links.

**Given** C# source code and specifications across modules and projects.
**When** Guideline linters audit the codebase for C# coding standards compliance.
**Then** All C# specifications and source files adhere to PascalCase naming, method sizing, async patterns, structured error handling, and type safety standards with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CS-002: C# PascalCase/camelCase and Positive Booleans

- [ ] All classes, structs, and records use PascalCase (`SnapshotManager`, `UserProfile`).
- [ ] Interfaces are prefixed with `I` in PascalCase (`IUserRepository`, `ILogger`).
- [ ] Methods and properties use PascalCase (`ProcessUpload()`, `GetActiveUsers()`, `IsActive`, `PluginSlug`).
- [ ] Local variables and parameters use camelCase (`pluginSlug`, `userId`, `siteId`).
- [ ] Private instance fields use `_camelCase` (`_logger`, `_connectionString`).
- [ ] Abbreviations treat acronyms as words with first-letter capitalization (`UserId` not `UserID`, `ApiClient` not `APIClient`, `GetUrl` not `GetURL`).
- [ ] All boolean properties, variables, and parameters use positive prefixes (`Is` and `Has` only).
- [ ] Negative boolean names are prohibited (`IsPending` not `IsNotReady`, `IsUnauthorized` not `HasNoPermission`).
- [ ] File names match the primary type declared inside using `{PascalCase}.cs` with strictly one primary type per file.
- [ ] Namespaces reflect folder hierarchy and use PascalCase (`RiseupAsia.Services`, `RiseupAsia.Domain.Models`).

**Given** C# source files containing type declarations, properties, methods, fields, and variables.
**When** Codebases are audited against C# naming conventions.
**Then** All classes, structs, records, interfaces, methods, and properties follow PascalCase, locals and parameters follow camelCase, abbreviations use word casing (`UserId`, `GetUrl`), booleans use affirmative prefixes with zero negative names, and file names match primary types with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CS-003: C# Method Signatures, Parameter Limits and Guard Clauses

- [ ] Zero boolean flag parameters that branch method logic — split into separate intention-revealing methods (`SaveDraft(doc)` vs `PublishDocument(doc)`).
- [ ] Shared setup/validation between split methods is extracted into private helper methods.
- [ ] Method bodies are capped at 15 lines (excluding error handling blocks).
- [ ] Method parameters are capped at at most 3 — 4 or more parameters must be encapsulated in an options class or record.
- [ ] Pure non-blocking async patterns throughout: zero calls to `.Result` or `.GetAwaiter().GetResult()`.
- [ ] Independent concurrent tasks use `Task.WhenAll` instead of sequential `await`.
- [ ] All asynchronous methods are suffixed with `Async` (`GetUsersAsync()`, `SaveDocumentAsync()`).
- [ ] LINQ is preferred over manual loops for collection transformations (`Select`, `Where`, `Any`).
- [ ] Complex LINQ predicates are extracted to named static private methods; nested LINQ deeper than 2 levels is prohibited.

**Given** C# methods, constructors, and async implementations.
**When** Method structures and signatures are analyzed against design guidelines.
**Then** Zero boolean flag parameters branch method logic, method bodies remain within 15 lines, parameter lists are capped at 3 or refactored into options objects, async methods use non-blocking patterns with `Async` suffix, and complex LINQ operations are cleanly factored out with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CS-004: C# Structured Exceptions, Specific Catches and Result Types

- [ ] Catch specific exception types (`HttpRequestException`, `JsonException`); bare `catch (Exception)` or untyped `catch` is banned.
- [ ] Exceptions are never silently swallowed; all catches must log contextual details and rethrow using parameterless `throw;` or handle explicitly.
- [ ] Nested conditional branching is eliminated using early return guard clauses.
- [ ] Parameter null validation uses `ArgumentNullException.ThrowIfNull(param)` or `if (param is null) throw new ArgumentNullException(nameof(param))`.
- [ ] Nullable reference types are enabled project-wide (`<Nullable>enable</Nullable>`) and annotated (`string?`).
- [ ] Null-coalescing throw operators (`?? throw new ...`) guard against unexpected null states.

**Given** C# error handling logic, catch blocks, and validation routines.
**When** Codebases are audited for exception safety and guard clause usage.
**Then** Catch blocks target specific exception types without silent swallows, parameters are validated using guard clauses and `nameof()`, nullable reference types are strictly annotated and handled, and nested branching is flattened with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CS-005: C# Nullable Reference Types and Pattern Matching

- [ ] Untyped `object` returns and parameters are replaced with strongly-typed generics (`T GetValue<T>(string key)`).
- [ ] Explicit type casting (`(User)obj`, `obj as User`) is eliminated in favor of pattern matching (`if (obj is User user)`).
- [ ] Exhaustive `switch` expressions are preferred over branching statements, using wildcard discards (`_ => throw ...`) for unhandled cases.
- [ ] Records and record structs (`public record UserDto(...)`) are utilized for immutable data transfer objects.
- [ ] Magic strings in business logic and state checks are replaced with typed enums (`StatusType.Active`) or typed constants.

**Given** C# type declarations, casting logic, and data structures.
**When** Source code is audited for type safety and pattern matching standards.
**Then** Generics are utilized in place of `object`, casting is eliminated in favor of pattern matching `is` expressions and `switch` expressions, DTOs are declared as immutable records, magic strings are replaced by typed enums, and nullable reference types are honored with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.

---

## 3. Cross-References & Traceability

- [`readme.md`](readme.md) — C# Coding Standards Overview
- [`02-naming-and-conventions.md`](02-naming-and-conventions.md) — Identifier casing, abbreviation casing, and boolean prefixes
- [`03-method-design.md`](03-method-design.md) — Boolean flag splitting, method sizing, async patterns, and LINQ
- [`04-error-handling.md`](04-error-handling.md) — Specific exceptions, guard clauses, and null safety
- [`05-type-safety.md`](05-type-safety.md) — Generics, pattern matching, records, and enum typing
- [Cross-Language Acceptance Criteria Registry](../01-cross-language/97-acceptance-criteria.md) — Master criteria
- [Boolean Flag Methods (cross-language)](../01-cross-language/24-boolean-flag-methods.md) — Method splitting specification
- [Generic Return Types (cross-language)](../01-cross-language/25-generic-return-types.md) — Generic signatures

---

## Verification & Acceptance Criteria

### AC-CG-CS-REG-001: C# Acceptance Criteria Registry Conformance

**Given** The C# acceptance criteria registry in `02-spec/02-coding-guidelines/07-csharp/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All C# specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.
