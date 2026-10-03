# PascalCase Key Naming Standard (AI Execution Prompt)

> **/goal** Enforce PascalCase key naming across all serialized formats (JSON, YAML, PHP array keys, Go log key constants, database columns) across the entire project.
> **/learn** Master the unified PascalCase serialization standard, understand exemptions for external third-party protocols/APIs, and eliminate camelCase or snake_case key drift.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce PascalCase formatting on all JSON request/response keys, YAML configuration keys, and database column names.
- [ ] `/learn` Ban camelCase and snake_case in project-controlled data structures, PHP array keys, and Go log key constants.
- [ ] `/goal` Allow exemptions strictly for uncontrollable third-party APIs, Prometheus metrics, or external protocol specifications.
- [ ] `/learn` Verify zero key naming violations across Go, PHP, and TypeScript codebases via automated guideline linters.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0
> **Updated:** 2026-03-09
> **Applies to:** Go, PHP, TypeScript — all serialized keys in this project

## 1. Rule

**ALL string keys** across the project **MUST** use PascalCase. This includes:

| Key Type | ❌ Wrong | ✅ Correct |
|----------|----------|-----------|
| JSON response keys | `"userId"`, `"createdAt"` | `"UserId"`, `"CreatedAt"` |
| JSON request keys | `"pageSize"`, `"totalItems"` | `"PageSize"`, `"TotalItems"` |
| Log context keys | `"errorCode"`, `"stackTrace"` | `"ErrorCode"`, `"StackTrace"` |
| Config keys | `"readTimeout"`, `"cacheSize"` | `"ReadTimeout"`, `"CacheSize"` |
| PHP array keys | `$data['pluginVersion']` | `$data['PluginVersion']` |
| PHP seed data keys | `'examId'`, `'wpId'` | `'ExamId'`, `'WpId'` |
| Go log key constants | `LogKeyUserId = "userId"` | `LogKeyUserId = "UserId"` |
| WebSocket message types | `"streamStart"` | `"StreamStart"` |
| YAML tags | `yaml:"cacheSize"` | `yaml:"CacheSize"` |
| Database column names | `user_id` | `UserId` |

## 2. Go Log Key Constants

```go
// ❌ WRONG — camelCase values
const (
    LogKeyUserId    = "userId"
    LogKeyErrorCode = "errorCode"
    LogKeyStackTrace = "stackTrace"
)

// ✅ CORRECT — PascalCase values
const (
    LogKeyUserId     = "UserId"
    LogKeyRequestId  = "RequestId"
    LogKeyErrorCode  = "ErrorCode"
    LogKeyStackTrace = "StackTrace"
)
```

## 3. Go Struct Serialization

Go structs serialize to PascalCase by default — **omit explicit JSON tags** unless `omitempty` or `json:"-"` is needed:

```go
// ✅ CORRECT — implicit PascalCase serialization
type User struct {
    Id        string
    SessionId string
    CreatedAt time.Time
}
// Serializes to: {"Id":"...","SessionId":"...","CreatedAt":"..."}
```

## 4. API Response / Request Keys

```json
// ❌ WRONG — camelCase
{
  "success": true,
  "data": {
    "id": "uuid",
    "createdAt": "2026-01-28T12:00:00Z"
  },
  "meta": {
    "requestId": "req_abc123",
    "pagination": {
      "pageSize": 20,
      "totalItems": 150
    }
  }
}

// ✅ CORRECT — PascalCase
{
  "Success": true,
  "Data": {
    "Id": "uuid",
    "CreatedAt": "2026-01-28T12:00:00Z"
  },
  "Meta": {
    "RequestId": "req_abc123",
    "Pagination": {
      "PageSize": 20,
      "TotalItems": 150
    }
  }
}
```

## 5. PHP Array Keys

```php
// ❌ WRONG — camelCase
$response = [
    'pluginVersion' => '1.0.0',
    'examId' => 42,
];

// ✅ CORRECT — PascalCase
$response = [
    'PluginVersion' => '1.0.0',
    'ExamId' => 42,
];
```

## 6. Exemptions

| Exemption | Reason | Example |
|-----------|--------|---------|
| External API parameters | Cannot control third-party naming | WordPress REST API `post_title` |
| Prometheus metrics | Prometheus convention requires snake_case | `requests_total` |
| Go standard library interfaces | Required by Go stdlib | `MarshalJSON()` |
| Native WordPress hooks | WordPress core naming convention | `plugins_api`, `admin_init` |
| Protocol-driven values | External protocol defines naming | HTTP headers `Content-Type` |

**Key principle:** If the key is **in our control**, it is PascalCase. If an **external system** defines it and we cannot change it, use the external convention.

## 7. Abbreviation Handling

Abbreviations follow the same rule as identifiers — capitalize only the first letter:

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `"userID"` | `"UserId"` |
| `"apiURL"` | `"ApiUrl"` |
| `"httpMethod"` | `"HttpMethod"` |

See [Master Coding Guidelines §1.2](./15-master-coding-guidelines/02-naming-and-database.md#12--abbreviation-standard-all-languages) for the full abbreviation table.

## 8. Cross-References

- [Master Coding Guidelines §1.1](./15-master-coding-guidelines/readme.md) — JSON / API keys row
- Go Backend Prompt — PascalCase mandate <!-- external: 02-spec/02-spec-management-software/12-prompts/01-coding-guideline/01-backend-go.md -->
- Coding Standards Foundation §10 — Log key constants <!-- external: 02-spec/01-general-spec/01-foundation/01-coding-standards-foundation.md -->
- [Database Naming](./07-database-naming.md) — Column naming

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-NAME-011: PascalCase Key Naming Standard

**Given** Serialized data structures, JSON payloads, YAML configurations, PHP array keys, and log constants in Go, PHP, or TypeScript.
**When** Codebases and data payloads are evaluated by CI/CD linters and guideline auditors.
**Then** All internal keys use PascalCase with zero unauthorized camelCase or snake_case occurrences, achieving deterministic compliance with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
