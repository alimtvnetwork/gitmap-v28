# TypeScript HttpMethod Enum — `src/lib/enums/http-method-type.ts`

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
- [TypeScript Standards](./08-typescript-standards-reference.md) — Parent TS spec
- [Master Coding Guidelines §8](../01-cross-language/15-master-coding-guidelines/00-overview.md) — Magic strings zero tolerance
- Enum Consumer Checklist — Cross-language sync process <!-- external: spec/02-spec-management-software/18-enum-consumer-checklist.md -->

---

*TypeScript HttpMethod enum v1.0.0 — 2026-02-27*
