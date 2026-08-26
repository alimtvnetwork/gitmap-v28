# Universal Response Envelope Specification


**Version:** 3.2.0  
**Last Updated:** 2026-04-16  

> **All systems (Go backend, PHP WordPress plugin, React frontend) MUST conform to this structure.**

## Top-Level Fields

Every API response is a JSON object with the following top-level keys. All JSON keys use **PascalCase**.

| Field | Type | Required | Description |
|---|---|---|---|
| `Status` | object | ✅ Always | Request outcome metadata |
| `Attributes` | object | ✅ Always | Response shape descriptors |
| `Results` | array | ✅ Always | Payload — always an array, even for single items |
| `Navigation` | object\|null | ⚙️ Conditional | Pagination links (only for paginated list responses) |
| `Errors` | object\|null | ⚙️ Conditional | Error details (only when errors exist AND reporting enabled) |
| `MethodsStack` | object\|null | ⚙️ Conditional | Debug call-chain trace (only when enabled in config) |

> **Go implementation note:** Use pointers (`*Navigation`, `*Errors`, `*MethodsStack`) with `omitempty` so absent sections are omitted from JSON entirely.

---

## Field Definitions

### Status

```json
{
  "Status": {
    "IsSuccess": true,
    "IsFailed": false,
    "Code": 200,
    "Message": "OK",
    "Timestamp": "2026-02-07T12:00:00Z"
  }
}
```

| Field | Type | Description |
|---|---|---|
| `IsSuccess` | bool | `true` if the request succeeded |
| `IsFailed` | bool | `true` if the request failed (inverse of IsSuccess) |
| `Code` | int | HTTP status code |
| `Message` | string | Human-readable status message |
| `Timestamp` | string | ISO 8601 timestamp of the response |

---

### Attributes

```json
{
  "Attributes": {
    "RequestedAt": "http://localhost:8080/api/v1/plugins",
    "RequestDelegatedAt": "https://example.com/wp-json/riseup-asia-uploader/v1/plugins",
    "HasAnyErrors": false,
    "IsSingle": false,
    "IsMultiple": true,
    "IsEmpty": false,
    "TotalRecords": 47,
    "PerPage": 10,
    "TotalPages": 5,
    "CurrentPage": 2
  }
}
```

| Field | Type | Description |
|---|---|---|
| `RequestedAt` | string | The Go backend endpoint that handled the request |
| `RequestDelegatedAt` | string | Downstream endpoint (PHP/other) if request was proxied; empty if not delegated |
| `HasAnyErrors` | bool | `true` if any errors occurred (even partial) |
| `IsSingle` | bool | `true` if `Results` contains exactly one item |
| `IsMultiple` | bool | `true` if `Results` contains multiple items |
| `IsEmpty` | bool | `true` when `Results` is empty (`TotalRecords` is 0 or `Results` array has no items) |
| `TotalRecords` | int | Total number of records across all pages |
| `PerPage` | int | Number of records per page |
| `TotalPages` | int | Total number of pages |
| `CurrentPage` | int | Current page number |

> For single-item responses, pagination fields (`TotalRecords`, `PerPage`, `TotalPages`, `CurrentPage`) may be omitted or zero.

---

### Results

Always an array. For single-item responses, the array has length 1.

```json
{
  "Results": [
    { "Id": 42, "Name": "My Plugin", "Slug": "my-plugin" }
  ]
}
```

---

### Navigation (conditional — pointer, omitempty)

Present only for paginated list responses. Contains URL strings for navigation links.

```json
{
  "Navigation": {
    "NextPage": "http://localhost:8080/api/v1/plugins?page=3&perPage=10",
    "PrevPage": "http://localhost:8080/api/v1/plugins?page=1&perPage=10",
    "CloserLinks": [
      "http://localhost:8080/api/v1/plugins?page=3&perPage=10",
      "http://localhost:8080/api/v1/plugins?page=4&perPage=10",
      "http://localhost:8080/api/v1/plugins?page=5&perPage=10"
    ]
  }
}
```

| Field | Type | Description |
|---|---|---|
| `NextPage` | string\|null | URL for the next page, null if on last page |
| `PrevPage` | string\|null | URL for the previous page, null if on first page |
| `CloserLinks` | string[] | Sliding window of nearby page URLs |

> Default `PerPage` value is seedable via backend config.

---

### Errors (conditional — pointer, omitempty)

Present only when `HasAnyErrors` is `true` AND error reporting is enabled in config.

```json
{
  "Errors": {
    "BackendMessage": "Failed to fetch plugin from remote site",
    "DelegatedServiceErrorStack": [
      "PHP Fatal error: Class 'PDO' not found in /var/www/plugin.php:42",
      "Stack trace:",
      "#0 /var/www/api.php(15): PluginManager->connect()"
    ],
    "Backend": [
      "handlers.go:85 handlePluginList",
      "service.go:120 FetchPlugins",
      "client.go:45 doRequest"
    ],
    "Frontend": []
  }
}
```

| Field | Type | Description |
|---|---|---|
| `BackendMessage` | string | Primary error message from the Go backend |
| `DelegatedServiceErrorStack` | string[] | Stack trace lines from downstream PHP/other services |
| `Backend` | string[] | Stack trace lines from the Go backend |
| `Frontend` | string[] | Reserved for frontend-injected error context |

---

### MethodsStack (conditional — pointer, omitempty)

Present only when debug/traversal mode is enabled in config.

```json
{
  "MethodsStack": {
    "Backend": [
      { "Method": "HandlePluginList", "File": "plugin_handlers.go", "LineNumber": 85 },
      { "Method": "FetchPlugins", "File": "plugin_service.go", "LineNumber": 120 }
    ],
    "Frontend": []
  }
}
```

| Field | Type | Description |
|---|---|---|
| `Backend` | array of `StackFrame` | Ordered list of method calls in the Go backend |
| `Frontend` | array of `StackFrame` | Reserved for frontend-injected call chain |

**StackFrame:**

| Field | Type |
|---|---|
| `Method` | string |
| `File` | string |
| `LineNumber` | int |

---

## Naming Conventions

### Boolean Fields

All boolean fields in API responses **MUST** use a prefix:
- **`Is`** — for state/identity booleans (e.g., `IsSuccess`, `IsFailed`, `IsSingle`, `IsMultiple`, `IsEmpty`, `IsActive`)
- **`Has`** — for presence booleans (e.g., `HasAnyErrors`)

**Do NOT** use bare boolean names like `success`, `active`, `enabled`. Use `isSuccess`, `isActive`, `isEnabled` instead.

This applies to:
- Envelope fields (`Status`, `Attributes`)
- Individual result object fields
- All backend response payloads (Go, PHP, WebSocket events)

---

## Rules

1. **Results is ALWAYS an array** — even for single items, deletions (empty array), or errors.
2. **PascalCase for ALL JSON keys** — no exceptions across Go, PHP, and frontend serialization.
3. **Conditional sections use pointers** — In Go, `Navigation`, `Errors`, and `MethodsStack` are pointer fields with `omitempty`. They are simply absent from the JSON when not applicable.
4. **Navigation is pagination-only** — It contains only page links. Errors and debug info are separate top-level fields.
5. **Errors and MethodsStack are top-level** — They are NOT nested inside Navigation.
6. **Boolean naming** — All boolean fields must use `Is` or `Has` prefix (see Naming Conventions above).
