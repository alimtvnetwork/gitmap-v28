# 10 — API Design

REST conventions, request/response formatting, versioning, and pagination.

## URL Structure

Lowercase, hyphen-separated, plural resource names. No verbs in URLs. No trailing slashes. Max two levels of nesting.

## HTTP Methods

| Method | Use | Idempotent |
|--------|-----|------------|
| GET | Read | Yes |
| POST | Create | No |
| PUT | Full replace | Yes |
| PATCH | Partial update | Yes |
| DELETE | Remove | Yes |

## Response Format

Always wrap in `data` or `error` envelope. Include `pagination` on list endpoints. Use machine-readable error codes.

## Versioning

URL path versioning: `/api/v1/`. Maintain backward compatibility within a major version. Document breaking changes.

## Pagination

Offset-based for simple cases. Cursor-based for large datasets. Default limit: 25, max: 100.

## Rate Limiting

Include `X-RateLimit-*` headers. Return `429` with `Retry-After` when exceeded.

## Authentication

`Authorization: Bearer <token>`. Never send credentials in query parameters.

---

Source: `spec/05-coding-guidelines/10-api-design.md`
