# Master Coding Guidelines — Magic strings, file organization, array keys (AI Execution Prompt)

> **/goal** Eliminate all magic string literals through typed enums and constants, enforce clean file and function size limits, and standardize PHP array key casing.
> **/learn** Enforce zero-tolerance for raw string comparisons in status and domain helpers, adhere to 300-line file and 15-line function limits in Go, and ensure log context keys use camelCase while database keys use PascalCase.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw string literals in comparison helpers and domain status checks with typed enum values and constants.
- [ ] `/learn` Prohibit `===` comparisons against raw status or categorical string literals across PHP, TypeScript, and Go.
- [ ] `/goal` Restrict Go source files to 300 lines (hard limit 400) and function bodies to 15 lines max, splitting into focused auxiliary files.
- [ ] `/learn` Standardize array keys in PHP to camelCase for log context and PascalCase for database records.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Master Coding Guidelines](./readme.md)
> **Version:** 2.1.0
> **Updated:** 2026-03-31

---

## 8. Magic Strings — Zero Tolerance

All repeated strings must be captured in enums or typed constants:

| Category | PHP Solution | Go Solution |
|----------|-------------|-------------|
| Hook names | `HookType::RestApiInit->value` | N/A |
| Capabilities | `CapabilityType::ManageOptions->value` | N/A |
| Table names | `TableType::Transactions->value` | Typed const |
| Error codes | `ErrorType::DATABASE_ERROR` | `apperror.ErrDatabase` |
| HTTP methods | `HttpMethodType::Post->value` | `http.MethodPost` |
| Log levels | `LogLevelType::Error->value` | `loglevel.Error.Lower()` |
| Status values | `StatusType::Success->value` | `status.Success.String()` |

### 8.1 — Comparison Helper Arguments

**Comparison helpers (`hasMismatch`, `isEqual`, `isMatch`) must always use enum constants — never raw string literals.**

```php
// ❌ WRONG: Magic string in comparison
if (hasMismatch($participant->status, 'active')) { ... }
if (isEqual($env, 'development')) { ... }

// ✅ CORRECT: Enum constant in comparison
if (hasMismatch($participant->status, ParticipantStatus::Active->value)) { ... }
if (isEqual($env, AppEnvType::Development->value)) { ... }
```

> **See also:** Issue #09 — Magic String Enum Comparisons <!-- external: 02-spec/23-how-app-issues-track/09-magic-string-enum-comparison.md -->

### 8.2 — Domain Status Comparisons

**All `===` comparisons against domain status, phase, role, or categorical string literals are PROHIBITED.** Use the corresponding enum constant instead.

```typescript
// ❌ FORBIDDEN: Raw string literal in status comparison
if (status === 'active') { ... }
if (execution.status === 'running') { ... }
if (connection.status === 'connected') { ... }
if (message.status === 'streaming') { ... }

// ✅ REQUIRED: Enum constant
if (status === EntityStatus.Active) { ... }
if (execution.status === ExecutionStatus.Running) { ... }
if (connection.status === ConnectionStatus.Connected) { ... }
if (message.status === MessageStatus.Streaming) { ... }
```

```go
// ❌ FORBIDDEN
if participant.Status == "active" { ... }

// ✅ REQUIRED
if participant.Status == entitystatus.Active { ... }
```

**Standardized domain enums:** `ExecutionStatus`, `ConnectionStatus`, `ExportStatus`, `MessageStatus`, `EntityStatus` — see [TypeScript Enum Inventory](../../02-typescript/09-typescript-standards-reference.md#enum-inventory).

**Exempt patterns** (no enum required):

- Framework/runtime APIs: `process.env.NODE_ENV === 'production'`
- Browser Web APIs: `mediaRecorder.state === 'recording'`
- Language operators: `typeof x === 'string'`
- Library internals: React Query `event.type === 'updated'`

> **See also:** Issue #10 — Domain Status Magic Strings <!-- external: 02-spec/23-how-app-issues-track/10-domain-status-magic-strings.md -->

---

---

## 9. File & Function Organization

### PHP

- PSR-4 autoloading: file name = class name
- One class/enum per file
- Traits must declare all their own `use` imports

### Go

- File target: 300 lines (hard limit 400)
- Function body: max 15 lines
- Split large files: `_crud.go`, `_helpers.go`, `_validation.go`
- Import order: stdlib → internal → third-party (3 groups, blank-line separated)

---

---

## 10. Array Key Conventions (PHP-Specific)

| Context | Convention | Example |
|---------|-----------|---------|
| Log context keys | camelCase | `'postId'`, `'masterDir'`, `'agentId'` |
| DB column keys | PascalCase | `'PluginSlug'`, `'CreatedAt'` |
| API response keys | Via `ResponseKeyType` enum | `ResponseKeyType::SnapshotId->value` |
| Persistence keys | Exempt (native casing) | `'schema_version'`, WP options |

### Common Mistakes

```php
// ❌ MISTAKE: snake_case in log context
$this->fileLogger->info('Post created', array('post_id' => $postId));

// ✅ CORRECT: camelCase
$this->fileLogger->info('Post created', array('postId' => $postId));
```

```php
// ❌ MISTAKE: camelCase for DB columns
$this->db->insert(TableType::Transactions->value, array('pluginSlug' => $slug));

// ✅ CORRECT: PascalCase matches schema
$this->db->insert(TableType::Transactions->value, array('PluginSlug' => $slug));
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-MASTER-MAGIC: Zero Magic Strings, File Sizing & Key Casing Organization

**Given** Application source code across PHP, Go, and TypeScript repositories.
**When** Guideline audit linters check for raw string comparisons, file sizes, and array key patterns.
**Then** All magic strings are represented by typed enums or constants, file and function lengths comply with size caps, and array keys use appropriate casing with zero violations.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
