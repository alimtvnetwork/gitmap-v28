# API Design

Universal guidelines for REST conventions, request/response formatting, versioning, and pagination patterns.

---

## 1. URL Structure

Use lowercase, hyphen-separated, plural resource names:

| Pattern | Example |
|---|---|
| List resources | `GET /api/v1/repositories` |
| Get single resource | `GET /api/v1/repositories/:id` |
| Create resource | `POST /api/v1/repositories` |
| Update resource | `PATCH /api/v1/repositories/:id` |
| Delete resource | `DELETE /api/v1/repositories/:id` |
| Nested resource | `GET /api/v1/repositories/:id/releases` |

### Rules

- No verbs in URLs — use HTTP methods to express actions
- No trailing slashes
- Maximum two levels of nesting; beyond that, use query filters
- Use path parameters for identity, query parameters for filtering

---

## 2. HTTP Methods & Status Codes

### Methods

| Method | Use | Idempotent |
|---|---|---|
| `GET` | Read | Yes |
| `POST` | Create | No |
| `PUT` | Full replace | Yes |
| `PATCH` | Partial update | Yes |
| `DELETE` | Remove | Yes |

### Status Codes

| Code | Meaning | When to Use |
|---|---|---|
| `200` | OK | Successful read or update |
| `201` | Created | Successful `POST` that creates a resource |
| `204` | No Content | Successful `DELETE` with no response body |
| `400` | Bad Request | Validation failure or malformed input |
| `401` | Unauthorized | Missing or invalid authentication |
| `403` | Forbidden | Authenticated but insufficient permissions |
| `404` | Not Found | Resource does not exist |
| `409` | Conflict | Duplicate or state conflict |
| `422` | Unprocessable Entity | Valid JSON but semantic errors |
| `429` | Too Many Requests | Rate limit exceeded |
| `500` | Internal Server Error | Unhandled server failure |

---

## 3. Request Format

### Body Conventions

- Use `camelCase` for JSON field names
- Send only the fields being changed on `PATCH`
- Use ISO 8601 for dates: `"2026-03-31T14:30:00Z"`
- Use strings for IDs (UUIDs), not integers

```json
{
  "name": "my-repo",
  "description": "A sample repository",
  "isPrivate": false,
  "createdAt": "2026-03-31T14:30:00Z"
}
```

### Query Parameters

| Purpose | Pattern | Example |
|---|---|---|
| Filter | `?status=active` | `GET /repositories?status=active` |
| Search | `?q=term` | `GET /repositories?q=gitmap-v28` |
| Sort | `?sort=name&order=asc` | `GET /repositories?sort=createdAt&order=desc` |
| Pagination | `?page=2&limit=25` | See Section 6 |

---

## 4. Response Format

### Success Envelope

Single resource:

```json
{
  "data": {
    "id": "abc-123",
    "name": "my-repo",
    "createdAt": "2026-03-31T14:30:00Z"
  }
}
```

Collection:

```json
{
  "data": [
    { "id": "abc-123", "name": "my-repo" }
  ],
  "pagination": {
    "page": 1,
    "limit": 25,
    "total": 142,
    "totalPages": 6
  }
}
```

### Error Envelope

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Name is required",
    "details": [
      { "field": "name", "reason": "must not be empty" }
    ]
  }
}
```

### Rules

- Always wrap responses in `data` or `error` — never return raw arrays or values
- Include `pagination` metadata on all list endpoints
- Use machine-readable `code` strings, not just human messages
- Never expose stack traces or internal paths in error responses

---

## 5. Versioning

### URL Path Versioning (Preferred)

```
GET /api/v1/repositories
GET /api/v2/repositories
```

### Rules

- Version the entire API, not individual endpoints
- Maintain backward compatibility within a major version
- Deprecate old versions with a sunset header: `Sunset: 2026-12-31`
- Document breaking changes in a changelog

### What Counts as Breaking

| Breaking | Non-Breaking |
|---|---|
| Removing a field | Adding a new optional field |
| Renaming a field | Adding a new endpoint |
| Changing a field type | Adding a new query parameter |
| Changing status code semantics | Adding a new enum value |

---

## 6. Pagination

### Offset-Based (Simple)

```
GET /repositories?page=2&limit=25
```

Response includes:

```json
{
  "pagination": {
    "page": 2,
    "limit": 25,
    "total": 142,
    "totalPages": 6
  }
}
```

### Cursor-Based (Large Datasets)

Use when total count is expensive or data changes frequently:

```
GET /repositories?cursor=eyJpZCI6MTAwfQ&limit=25
```

```json
{
  "pagination": {
    "nextCursor": "eyJpZCI6MTI1fQ",
    "hasMore": true,
    "limit": 25
  }
}
```

### Defaults

- Default `limit`: 25
- Maximum `limit`: 100
- Default `page`: 1
- Return `total` and `totalPages` for offset pagination

---

## 7. Rate Limiting

Include rate limit headers on every response:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1711900800
```

Return `429 Too Many Requests` with a `Retry-After` header when exceeded.

---

## 8. Authentication & Authorization Headers

```
Authorization: Bearer <token>
```

- Use `Bearer` tokens (JWT or opaque)
- Never send credentials in query parameters
- Return `401` for missing/invalid tokens, `403` for insufficient permissions
- See [Security & Secrets](./08-security-secrets.md) for credential handling rules

---

## References

- [Security & Secrets](./08-security-secrets.md)
- [Error Handling Patterns](./04-error-handling.md)
- [Logging & Observability](./07-logging-observability.md)
- [Performance & Optimization](./09-performance-optimization.md)

---

**Contributors**: Alim Ul Karim · [Riseup Labs](https://riseuplabs.com)
