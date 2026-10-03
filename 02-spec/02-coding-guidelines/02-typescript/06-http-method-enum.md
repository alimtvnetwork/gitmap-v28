# TypeScript HttpMethod Enum — `src/lib/enums/http-method-type.ts` (AI Execution Prompt)

> **/goal** Eliminate all magic string HTTP method literals (`"GET"`, `"POST"`, `"PUT"`, `"PATCH"`, `"DELETE"`, `"HEAD"`, `"OPTIONS"`) in `fetch()` calls and endpoint configs by standardizing on `HttpMethod`.
> **/learn** Master frontend HTTP method typing: `export enum HttpMethod` with kebab-case `-type.ts` file convention, typed webhook/endpoint config arrays, and full parity with Go `httpmethodtype.Variant`.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw HTTP method string literals in `fetch()` and API clients with `HttpMethod` enum members (e.g., `HttpMethod.Post`).
- [ ] `/learn` Never use string union types (`"POST" | "PUT"`) for HTTP methods in configuration interfaces; use `HttpMethod.Post | HttpMethod.Put`.
- [ ] `/goal` Standardize endpoint configuration tables and route definitions to use `HttpMethod` constants.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 2.0.0
> **Last updated**: 2026-02-28
> **Parity with**: [Go HttpMethod Enum](../03-golang/03-httpmethod-enum.md)

---

## Purpose

Frontend equivalent of the Go `httpmethod.Variant` enum. Replaces all magic string HTTP method literals (`"GET"`, `"POST"`, etc.) in `fetch()` calls and endpoint configuration across all frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/http-method-type.ts

export enum HttpMethod {
  Get = "GET",
  Head = "HEAD",
  Post = "POST",
  Put = "PUT",
  Patch = "PATCH",
  Delete = "DELETE",
  Options = "OPTIONS",
}
```

> **Convention:** TypeScript enum files use `-type` suffix in kebab-case (e.g., `http-method-type.ts`, `execution-status-type.ts`). The enum name itself remains PascalCase without suffix since the `enum` keyword already signals the construct.

---

## Usage Patterns

### fetch() Calls

```typescript
// ❌ WRONG: Magic string
const resp = await fetch("/api/v1/sites/validate", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ Url: url }),
});

// ✅ CORRECT: Enum constant (imported from http-method-type.ts)
const resp = await fetch("/api/v1/sites/validate", {
  method: HttpMethod.Post,
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ Url: url }),
});
```

### Endpoint Configuration Arrays

```typescript
// ❌ WRONG: Magic strings in config
const endpoints = [
  { method: "POST", path: "/search", name: "Execute Search" },
  { method: "GET", path: "/search/engines", name: "List Engines" },
];

// ✅ CORRECT: Enum constants
const endpoints = [
  { method: HttpMethod.Post, path: "/search", name: "Execute Search" },
  { method: HttpMethod.Get, path: "/search/engines", name: "List Engines" },
];
```

### Type Constraints

```typescript
// ❌ WRONG: Union of magic strings
interface WebhookConfig {
  readonly method: "POST" | "PUT";
}

// ✅ CORRECT: Enum-typed constraint
interface WebhookConfig {
  readonly method: HttpMethod.Post | HttpMethod.Put;
}
```

---

## Cross-Language Parity

| Feature | Go (`httpmethodtype.Variant`) | TypeScript (`HttpMethod`) |
|---------|--------------------------|---------------------------|
| Package | `pkg/enums/httpmethodtype` | `src/lib/enums/http-method-type.ts` |
| Type | `byte` iota | String enum |
| Values | `Get`, `Post`, `Put`, `Patch`, `Delete`, `Head`, `Options` | Same |
| String output | `.String()` → `"GET"` | Direct value `"GET"` |
| Parse | `httpmethodtype.Parse("GET")` | N/A (enum is the string) |

---

## Cross-References

- [Go HttpMethod Enum](../03-golang/03-httpmethod-enum.md) — Backend parity spec
- [TypeScript Standards](./09-typescript-standards-reference.md) — Parent TS spec
- [Master Coding Guidelines §8](../01-cross-language/15-master-coding-guidelines/readme.md) — Magic strings zero tolerance
- Enum Consumer Checklist — Cross-language sync process <!-- external: 02-spec/02-spec-management-software/18-enum-consumer-checklist.md -->

---

*TypeScript HttpMethod enum v1.0.0 — 2026-02-27*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-006: HttpMethodEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All HTTP request method assignments and configuration tables strictly utilize `HttpMethod` enum constants, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
