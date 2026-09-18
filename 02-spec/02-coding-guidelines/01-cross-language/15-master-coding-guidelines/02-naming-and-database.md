# Master Coding Guidelines — Naming conventions, database naming, file naming

> **Parent:** [Master Coding Guidelines](./01-index.md)
> **Version:** 2.1.0
> **Updated:** 2026-03-31

---

## How to Use This Document

This is the **master reference**. Every rule here is enforced across all languages. Language-specific details are in:

- [PHP Standards](../../04-php/01-index.md)
- [Go Standards](../../03-golang/04-golang-standards-reference/01-index.md)
- [TypeScript Standards](../../02-typescript/09-typescript-standards-reference.md)
- [Database Naming](../07-database-naming.md)
- [Boolean Principles](../02-boolean-principles/01-index.md)
- [No-Negatives](../12-no-negatives.md)
- [Test Naming & Structure](../14-test-naming-and-structure.md)
- [Lazy Evaluation Patterns](../16-lazy-evaluation-patterns.md)
- [Regex Usage Guidelines](../17-regex-usage-guidelines.md)
- [Code Mutation Avoidance](../18-code-mutation-avoidance.md)
- [Null Pointer Safety](../19-null-pointer-safety.md)
- [Nesting Resolution Patterns](../20-nesting-resolution-patterns.md)
- [Newline Styling Examples](../21-newline-styling-examples.md)
- [Defer Rules (Go)](../../03-golang/05-defer-rules.md)

---


---

## 1. Naming Conventions

### 1.1 — Universal Rules

| Element | Convention | PHP Example | Go Example | TS Example |
|---------|-----------|-------------|------------|------------|
| Class / Struct | PascalCase | `SnapshotManager` | `SnapshotManager` | `SnapshotManager` |
| Enum type name | PascalCase + `Type` suffix | `StatusType` | `status.Variant` (package-scoped) | `StatusType` |
| Enum case / constant | PascalCase | `StatusType::Success` | `status.Success` | `StatusType.Success` |
| Method (exported) | camelCase (PHP) / PascalCase (Go) | `processUpload()` | `ProcessUpload()` | `processUpload()` |
| Variable | camelCase | `$pluginSlug` | `pluginSlug` | `pluginSlug` |
| Boolean variable | `is`/`has` + camelCase (99%), `should` rare | `$isActive` | `isActive` | `isActive` |
| File name (source) | PascalCase | `SnapshotManager.php` | `SnapshotManager.go` | `SnapshotManager.tsx` |
| Directory (Go pkg) | snake_case | — | `site_health/` | — |
| Abbreviations | First letter only caps | `$postId`, `$fileUrl` | `postId`, `fileUrl` | `postId`, `fileUrl` |
| JSON / API keys | PascalCase | `"PluginSlug"` | `"PluginSlug"` | `"PluginSlug"` |

### 1.2 — Abbreviation Standard (ALL LANGUAGES)

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `ID` | `Id` |
| `URL` | `Url` |
| `MD5` | `Md5` |
| `JSON` | `Json` |
| `API` | `Api` |
| `IP` | `Ip` |
| `SQL` | `Sql` |
| `HTTP` | `Http` |
| `HTML` | `Html` |
| `YAML` | `Yaml` |
| `XML` | `Xml` |
| `CSS` | `Css` |
| `DB` | `Db` |

> **Go Standard Library Exemptions:** `MarshalJSON()`, `UnmarshalJSON()`, `Error() string`, and `String() string` are required by Go's standard library interfaces. These method names are **exempt** and MUST retain their original spelling. Go standard library struct fields accessed on framework types (e.g., `r.URL.Path` from `net/http`) are also exempt — you cannot rename them. All **custom** identifiers (struct fields, variables, function names, parameters, enum constants, `variantLabels` values) follow the table above.

### 1.3 — Source File Naming: PascalCase (ALL LANGUAGES)

Source code files that define a **single primary type** (struct, class, component, enum) MUST be named in PascalCase to match the definition name.

| Scenario | File Name | Why |
|----------|-----------|-----|
| Go struct `SiteManager` | `SiteManager.go` | Matches primary type |
| PHP class `SnapshotManager` | `SnapshotManager.php` | Matches primary class |
| TS component `UserProfile` | `UserProfile.tsx` | Matches primary component |
| Go enum package `content_type` | `Variant.go` | Matches Go enum convention (type is `Variant`) |
| PHP enum `StatusType` | `StatusType.php` | Matches enum type name |
| Go test file for `SiteManager` | `SiteManager_test.go` | Matches source + `_test` suffix |

**Exemptions (keep lowercase):**

| File | Reason |
|------|--------|
| `main.go`, `index.ts`, `index.php` | Package/module entry points |
| `helpers.go`, `utils.ts` | Multi-purpose utility files with no single primary type |
| `routes.go`, `middleware.go` | Infrastructure files not tied to one definition |
| Go package directories | Stay `snake_case` per Go convention (`site_health/`) |

**Spec/documentation files** are **NOT** affected — they remain lowercase kebab-case with numeric prefixes (e.g., `02-architecture.md`). See File Naming Convention <!-- external: 02-spec/02-spec-management-software/02-instructions/01-file-naming-convention.md -->.

```
❌ WRONG — Go file with snake_case for a single-type file
snapshot_manager.go   →  contains type SnapshotManager

✅ CORRECT — PascalCase matches the definition
SnapshotManager.go   →  contains type SnapshotManager
```

```
❌ WRONG — TypeScript component with kebab-case
user-profile.tsx     →  contains function UserProfile

✅ CORRECT — PascalCase matches the component
UserProfile.tsx      →  contains function UserProfile
```

### 1.4 — Zero Underscore Policy

**Snake_case is prohibited** for all logic-level identifiers across PHP, Go, and TypeScript. This includes:

- Variables, method names, properties, parameters
- Log context array keys (PHP): use camelCase (`'postId'`, not `'post_id'`)
- Internal array keys used in code logic

**Exemptions** (persistence-level only):

- WordPress hooks, capabilities, option keys, core table/column names
- Database migration rename maps (old→new mappings)
- PHP superglobals (`$_GET`, `$_POST`)
- HTML form `name` attributes and URL query parameters
- WP-Cron arguments, manifest JSON keys
- Go `runtime.GOOS` comparisons (`"windows"`, `"darwin"`)

---


---

## 2. Database Naming — PascalCase

> Full reference: [database-naming.md](../07-database-naming.md)

### Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Custom table names | PascalCase | `Transactions`, `AgentSites` |
| Custom column names | PascalCase | `PluginSlug`, `CreatedAt` |
| Index names | `Idx` prefix + PascalCase | `IdxTransactions_CreatedAt` |
| WordPress core tables | snake_case (EXEMPT) | `wp_posts`, `wp_options` |

### Common Mistakes

```php
// ❌ MISTAKE: Using camelCase or snake_case for DB columns
$record = array(
    'pluginSlug' => $slug,       // Wrong — camelCase
    'created_at' => $now,        // Wrong — snake_case
);

// ✅ CORRECT: PascalCase matches the schema
$record = array(
    'PluginSlug' => $slug,
    'CreatedAt'  => $now,
);
```

```go
// ❌ MISTAKE: snake_case in SQL and struct tags
const query = `SELECT plugin_slug FROM transactions`
type Tx struct {
    PluginSlug string `db:"plugin_slug"`
}

// ✅ CORRECT: PascalCase db tag, no redundant json tag
const query = `SELECT PluginSlug FROM Transactions`
type Tx struct {
    PluginSlug string `db:"PluginSlug"`
}
```

---
