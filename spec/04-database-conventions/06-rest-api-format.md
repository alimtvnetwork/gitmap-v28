# REST API Response Format

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

All REST API responses MUST use **PascalCase** for JSON keys, matching the database column naming convention. This ensures a single casing standard from database → ORM → API response → frontend TypeScript types.

All responses MUST use the **Universal Response Envelope** defined in the [Response Envelope Specification](../03-error-manage/02-error-architecture/05-response-envelope/04-response-envelope-reference.md).

---

## Golden Rule

> **Every JSON key in a REST API response MUST be PascalCase.** No camelCase, no snake_case, no kebab-case.

---

## 1. Response Key Format

```json
// ✅ CORRECT — PascalCase keys
{
    "TransactionId": 42,
    "PluginSlug": "my-awesome-plugin",
    "Amount": 29.99,
    "StatusName": "Pending",
    "IsActive": true,
    "HasLicense": false,
    "CreatedAt": "2026-04-02T10:30:00Z"
}

// ❌ WRONG — camelCase
{
    "TransactionId": 42,
    "PluginSlug": "my-awesome-plugin",
    "IsActive": true
}

// ❌ WRONG — snake_case
{
    "transaction_id": 42,
    "plugin_slug": "my-awesome-plugin",
    "is_active": true
}
```

---

## 2. Full REST API Samples (Using Universal Envelope)

All examples below conform to the [Universal Response Envelope](../03-error-manage/02-error-architecture/05-response-envelope/04-response-envelope-reference.md).

### 2.1 List Resources (Paginated)

```
GET /api/v1/transactions?StatusName=Pending&IsActive=1&page=2&perPage=10
```

**Response:**

```json
{
    "Status": {
        "IsSuccess": true,
        "IsFailed": false,
        "Code": 200,
        "Message": "OK",
        "Timestamp": "2026-04-02T10:30:00Z"
    },
    "Attributes": {
        "RequestedAt": "http://localhost:8080/api/v1/transactions",
        "RequestDelegatedAt": "",
        "HasAnyErrors": false,
        "IsSingle": false,
        "IsMultiple": true,
        "IsEmpty": false,
        "TotalRecords": 47,
        "PerPage": 10,
        "TotalPages": 5,
        "CurrentPage": 2
    },
    "Results": [
        {
            "TransactionId": 1,
            "PluginSlug": "my-plugin",
            "Amount": 29.99,
            "StatusName": "Pending",
            "FileTypeName": "Plugin",
            "AgentSiteName": "Example Site",
            "IsActive": true,
            "HasLicense": true,
            "CreatedAt": "2026-04-02T10:30:00Z"
        },
        {
            "TransactionId": 2,
            "PluginSlug": "another-plugin",
            "Amount": 9.99,
            "StatusName": "Pending",
            "FileTypeName": "Theme",
            "AgentSiteName": "Other Site",
            "IsActive": true,
            "HasLicense": false,
            "CreatedAt": "2026-04-01T08:15:00Z"
        }
    ],
    "Navigation": {
        "NextPage": "http://localhost:8080/api/v1/transactions?page=3&perPage=10",
        "PrevPage": "http://localhost:8080/api/v1/transactions?page=1&perPage=10",
        "CloserLinks": [
            "http://localhost:8080/api/v1/transactions?page=3&perPage=10",
            "http://localhost:8080/api/v1/transactions?page=4&perPage=10",
            "http://localhost:8080/api/v1/transactions?page=5&perPage=10"
        ]
    }
}
```

### 2.2 Get Single Resource

```
GET /api/v1/transactions/42
```

**Response:**

```json
{
    "Status": {
        "IsSuccess": true,
        "IsFailed": false,
        "Code": 200,
        "Message": "OK",
        "Timestamp": "2026-04-02T11:00:00Z"
    },
    "Attributes": {
        "RequestedAt": "http://localhost:8080/api/v1/transactions/42",
        "RequestDelegatedAt": "",
        "HasAnyErrors": false,
        "IsSingle": true,
        "IsMultiple": false,
        "IsEmpty": false,
        "TotalRecords": 1,
        "PerPage": 0,
        "TotalPages": 0,
        "CurrentPage": 0
    },
    "Results": [
        {
            "TransactionId": 42,
            "PluginSlug": "my-plugin",
            "Amount": 29.99,
            "StatusTypeId": 1,
            "StatusName": "Pending",
            "FileTypeId": 1,
            "FileTypeName": "Plugin",
            "AgentSiteId": 5,
            "AgentSiteName": "Example Site",
            "IsActive": true,
            "HasLicense": true,
            "CreatedAt": "2026-04-02T10:30:00Z",
            "UpdatedAt": "2026-04-02T11:00:00Z"
        }
    ]
}
```

### 2.3 Create Resource

```
POST /api/v1/transactions
Content-Type: application/json
```

**Request body (PascalCase keys):**

```json
{
    "PluginSlug": "new-plugin",
    "Amount": 19.99,
    "StatusTypeId": 1,
    "FileTypeId": 2,
    "AgentSiteId": 3
}
```

**Response:**

```json
{
    "Status": {
        "IsSuccess": true,
        "IsFailed": false,
        "Code": 201,
        "Message": "Transaction created successfully",
        "Timestamp": "2026-04-02T12:00:00Z"
    },
    "Attributes": {
        "RequestedAt": "http://localhost:8080/api/v1/transactions",
        "RequestDelegatedAt": "",
        "HasAnyErrors": false,
        "IsSingle": true,
        "IsMultiple": false,
        "IsEmpty": false,
        "TotalRecords": 1,
        "PerPage": 0,
        "TotalPages": 0,
        "CurrentPage": 0
    },
    "Results": [
        {
            "TransactionId": 43,
            "PluginSlug": "new-plugin",
            "Amount": 19.99,
            "StatusName": "Pending",
            "FileTypeName": "Theme",
            "IsActive": true,
            "CreatedAt": "2026-04-02T12:00:00Z"
        }
    ]
}
```

### 2.4 Update Resource

```
PUT /api/v1/transactions/42
Content-Type: application/json
```

**Request body:**

```json
{
    "StatusTypeId": 2,
    "Amount": 34.99
}
```

**Response:**

```json
{
    "Status": {
        "IsSuccess": true,
        "IsFailed": false,
        "Code": 200,
        "Message": "Transaction updated successfully",
        "Timestamp": "2026-04-02T13:00:00Z"
    },
    "Attributes": {
        "RequestedAt": "http://localhost:8080/api/v1/transactions/42",
        "RequestDelegatedAt": "",
        "HasAnyErrors": false,
        "IsSingle": true,
        "IsMultiple": false,
        "IsEmpty": false,
        "TotalRecords": 1,
        "PerPage": 0,
        "TotalPages": 0,
        "CurrentPage": 0
    },
    "Results": [
        {
            "TransactionId": 42,
            "StatusName": "Complete",
            "Amount": 34.99,
            "UpdatedAt": "2026-04-02T13:00:00Z"
        }
    ]
}
```

### 2.5 Delete Resource

```
DELETE /api/v1/transactions/42
```

**Response:**

```json
{
    "Status": {
        "IsSuccess": true,
        "IsFailed": false,
        "Code": 200,
        "Message": "Transaction deleted successfully",
        "Timestamp": "2026-04-02T13:30:00Z"
    },
    "Attributes": {
        "RequestedAt": "http://localhost:8080/api/v1/transactions/42",
        "RequestDelegatedAt": "",
        "HasAnyErrors": false,
        "IsSingle": false,
        "IsMultiple": false,
        "IsEmpty": true,
        "TotalRecords": 0,
        "PerPage": 0,
        "TotalPages": 0,
        "CurrentPage": 0
    },
    "Results": []
}
```

### 2.6 Error Response

```json
{
    "Status": {
        "IsSuccess": false,
        "IsFailed": true,
        "Code": 404,
        "Message": "Transaction not found",
        "Timestamp": "2026-04-02T14:00:00Z"
    },
    "Attributes": {
        "RequestedAt": "http://localhost:8080/api/v1/transactions/999",
        "RequestDelegatedAt": "",
        "HasAnyErrors": true,
        "IsSingle": false,
        "IsMultiple": false,
        "IsEmpty": true,
        "TotalRecords": 0,
        "PerPage": 0,
        "TotalPages": 0,
        "CurrentPage": 0
    },
    "Results": [],
    "Errors": {
        "BackendMessage": "Transaction not found",
        "DelegatedServiceErrorStack": [],
        "Backend": [
            "handlers.go:92 handleGetTransaction",
            "service.go:45 FindTransactionById"
        ],
        "Frontend": []
    }
}
```

---

## 3. Envelope Quick Reference

All responses use the [Universal Response Envelope](../03-error-manage/02-error-architecture/05-response-envelope/04-response-envelope-reference.md). Summary:

| Section | Type | Present | Description |
|---------|------|---------|-------------|
| `Status` | object | ✅ Always | `IsSuccess`, `IsFailed`, `Code`, `Message`, `Timestamp` |
| `Attributes` | object | ✅ Always | Shape flags (`IsSingle`, `IsMultiple`, `IsEmpty`), pagination, error flag |
| `Results` | array | ✅ Always | Payload — **always an array**, even for single items or deletes |
| `Navigation` | object\|null | ⚙️ Conditional | Pagination links (paginated lists only) |
| `Errors` | object\|null | ⚙️ Conditional | Error details (when `HasAnyErrors` is `true` AND reporting enabled) |
| `MethodsStack` | object\|null | ⚙️ Conditional | Debug call-chain trace (when enabled in config) |

> **Key rule:** `Results` is ALWAYS an array. Delete returns `[]`. Single item returns `[{...}]`.

---

## 4. Data Flow — PascalCase End-to-End

```
Database Column     →  ORM Struct/Model   →  API Response JSON  →  Frontend Type
───────────────────────────────────────────────────────────────────────────────
PluginSlug TEXT      →  PluginSlug string  →  "PluginSlug": "x"  →  PluginSlug: string
IsActive BOOLEAN     →  IsActive bool      →  "IsActive": true   →  IsActive: boolean
StatusTypeId INT     →  StatusTypeId int   →  "StatusTypeId": 1  →  StatusTypeId: number
CreatedAt TEXT       →  CreatedAt string   →  "CreatedAt": "..."  →  CreatedAt: string
```

> **No transformation layer needed.** PascalCase flows from DB to frontend without any key mapping or conversion.

---

## 5. Language Implementation

### Go — Default Behavior

Go marshals PascalCase struct fields to PascalCase JSON by default — no tags needed:

```go
type Transaction struct {
    TransactionId int64   `db:"TransactionId"`
    PluginSlug    string  `db:"PluginSlug"`
    IsActive      bool    `db:"IsActive"`
    Amount        float64 `db:"Amount"`
}

// json.Marshal(tx) produces:
// {"TransactionId":42,"PluginSlug":"my-plugin","IsActive":true,"Amount":29.99}
```

> **Do NOT add `json:"PluginSlug"` tags** — that would convert to camelCase.

### PHP — Array Keys

```php
// Response array keys are PascalCase
return [
    'Success' => true,
    'Data' => [
        'TransactionId' => $tx->TransactionId,
        'PluginSlug'    => $tx->PluginSlug,
        'IsActive'      => $tx->IsActive,
        'Amount'        => $tx->Amount,
    ],
];
```

### TypeScript — Interface

```typescript
// Envelope types
interface ApiResponse<T> {
    Status: {
        IsSuccess: boolean;
        IsFailed: boolean;
        Code: number;
        Message: string;
        Timestamp: string;
    };
    Attributes: {
        RequestedAt: string;
        RequestDelegatedAt: string;
        HasAnyErrors: boolean;
        IsSingle: boolean;
        IsMultiple: boolean;
        IsEmpty: boolean;
        TotalRecords: number;
        PerPage: number;
        TotalPages: number;
        CurrentPage: number;
    };
    Results: T[];
    Navigation?: {
        NextPage: string | null;
        PrevPage: string | null;
        CloserLinks: string[];
    };
    Errors?: {
        BackendMessage: string;
        DelegatedServiceErrorStack: string[];
        Backend: string[];
        Frontend: string[];
    };
}

// Domain type
interface Transaction {
    TransactionId: number;
    PluginSlug: string;
    IsActive: boolean;
    Amount: number;
    StatusName: string;
    CreatedAt: string;
}

// Usage
type TransactionListResponse = ApiResponse<Transaction>;
```

---

## 6. URL Paths vs JSON Keys

| Context | Convention | Example |
|---------|-----------|---------|
| URL paths (slugs) | **kebab-case lowercase** | `/api/v1/blog-posts/my-first-post` |
| Query parameters | **PascalCase** | `?StatusName=Pending&IsActive=1` |
| JSON request keys | **PascalCase** | `{"PluginSlug": "my-plugin"}` |
| JSON response keys | **PascalCase** | `{"TransactionId": 42}` |

> See [slug conventions](../02-coding-guidelines/01-cross-language/28-slug-conventions.md) for URL path rules.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| **Response Envelope Spec** | [04-response-envelope-reference.md](../03-error-manage/02-error-architecture/05-response-envelope/04-response-envelope-reference.md) |
| Envelope examples | [envelope-minimal.json](../03-error-manage/02-error-architecture/05-response-envelope/envelope-minimal.json), [envelope-single.json](../03-error-manage/02-error-architecture/05-response-envelope/envelope-single.json), [envelope-multiple.json](../03-error-manage/02-error-architecture/05-response-envelope/envelope-multiple.json) |
| Database naming | [./01-naming-conventions.md](./01-naming-conventions.md) |
| Schema design | [./02-schema-design.md](./02-schema-design.md) |
| Key naming PascalCase | [../02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md](../02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md) |
| Slug conventions | [../02-coding-guidelines/01-cross-language/28-slug-conventions.md](../02-coding-guidelines/01-cross-language/28-slug-conventions.md) |

---

## Required Request & Response Headers

**Added:** 2026-05-04 (resolves spec/19 audit F-X-10, F-A-22 — top-10 fix #8). Authoritative for ALL REST APIs across every spec, not just spec/19.

### Inbound headers (request)

| Header | Required when | Format | Notes |
|---|---|---|---|
| `X-Correlation-Id` | Always (every request) | ULID (26 chars) — server-generates if missing and echoes in response | Single id propagates across tier hops (Main → Worker, etc.) and into log lines. |
| `X-Idempotency-Key` | All `POST`, `PUT`, `PATCH` | ULID (26 chars), max 64 chars | Server stores result for `IdempotencyKeyTtlSeconds` (per `spec/19-main-worker-service/15-tunable-constants.md` §2.2). Replay returns cached result. |
| `X-Auth-Action` | When endpoint is multi-step (e.g. login flows) | PascalCase action name, max 64 chars | Distinguishes step-1 (`SignInBegin`) from step-2 (`Verify2FA`) under one correlation id. |
| `Authorization` | Per endpoint auth requirement | `Bearer <token>` | Token format defined per spec (e.g. spec/19/12 for Worker JWT). |
| `Content-Type` | Always for body-bearing methods | `application/json; charset=utf-8` | Reject other types with `400`. |

### Outbound headers (response)

| Header | Set when | Format |
|---|---|---|
| `X-Correlation-Id` | Always | Same value as inbound (or server-generated if inbound was missing). |
| `X-Idempotency-Key` | When inbound carried one | Echoed verbatim. |
| `X-Idempotency-Replay` | Only if response is a cached replay | `true`. Absent on first execution. |
| `Retry-After` | On 429 / 503 | Seconds (integer). |

### Validation rules

1. Missing `X-Idempotency-Key` on `POST/PUT/PATCH` → `400` with envelope error code `WORKER-300-01` (spec/19 tier) or analogous code in other specs.
2. Missing `X-Correlation-Id` → server generates, no error. NEVER reject solely for absent correlation id.
3. `X-Auth-Action` is REQUIRED on endpoints whose path matches `/Auth/*` (multi-step flows). Optional elsewhere.
4. Header names are case-insensitive on the wire (per RFC 9110); spec uses canonical casing for documentation.

### Why these three?

- `X-Correlation-Id` makes cross-tier debugging tractable.
- `X-Idempotency-Key` prevents double-charge / double-create on retries.
- `X-Auth-Action` lets server pin which step of a multi-step flow this request belongs to (resolves spec/19 audit gap where 2FA step-2 could not be distinguished from a login replay).

### Cross-references

- `spec/19-main-worker-service/06-core-api-endpoints.md` §1 — first introduced these headers.
- `spec/19-main-worker-service/15-tunable-constants.md` §2.2 — TTL + max-length values.
- `spec/19-main-worker-service/13-error-codes.md` — error codes for header-validation failures.
- `spec/03-error-manage/02-error-architecture/05-response-envelope/` — envelope shape that carries these headers' echoes.
