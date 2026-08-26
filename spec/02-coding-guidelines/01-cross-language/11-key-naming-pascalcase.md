# PascalCase Key Naming Standard

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

See [Master Coding Guidelines §1.2](./15-master-coding-guidelines/01-naming-and-database.md#12--abbreviation-standard-all-languages) for the full abbreviation table.

## 8. Cross-References

- [Master Coding Guidelines §1.1](./15-master-coding-guidelines/00-overview.md) — JSON / API keys row
- Go Backend Prompt — PascalCase mandate <!-- external: spec/02-spec-management-software/12-prompts/01-coding-guideline/01-backend-go.md -->
- Coding Standards Foundation §10 — Log key constants <!-- external: spec/01-general-spec/01-foundation/01-coding-standards-foundation.md -->
- [Database Naming](./07-database-naming.md) — Column naming
