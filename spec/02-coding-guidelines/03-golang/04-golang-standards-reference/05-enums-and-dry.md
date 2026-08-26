# Golang Coding Standards — Typed constants, enums, DRY enforcement

> **Parent:** [Golang Coding Standards](./00-overview.md)  
> **Version:** 3.7.0  
> **Updated:** 2026-03-31

---

## Typed Constants & Enums

> **Canonical source:** [Go Enum Specification](../01-enum-specification/00-overview.md) — core pattern, required methods, folder structure  
> **Cross-language reference:** [Enum Naming Quick Reference](../../06-ai-optimization/05-enum-naming-quick-reference.md) — Go, TypeScript, PHP comparison

All Go enum rules (byte type, `Invalid` zero value, `iota`, `variantLabels`, required methods, folder structure) are defined in the [Go Enum Specification](../01-enum-specification/00-overview.md). Do not duplicate here.

### Zero Magic Strings/Numbers

- All HTTP status codes → typed constants
- All error codes → `apperror` code constants
- All config keys → typed const block
- All status/event strings → typed byte-based enum constants
- All HTTP methods → `httpmethodtype.Variant` enum (see [HTTP Method Enum](../03-httpmethod-enum.md))
- All API operation names → per-domain operation enum

### API Client — No Tuple Returns, No Magic Strings

API client helper functions MUST return `apperror.Result[T]`, not `(*T, error)`. All struct literal fields for method and operation MUST use enum constants.

```go
// ❌ FORBIDDEN: magic strings + tuple return
func (c *Client) CleanupSnapshots(opts SnapshotCleanupOptions) (*SnapshotCleanupResult, error) {
    callInput := apiCallInput{
        Method:    "POST",
        Operation: "snapshot cleanup",
    }

    return doAPICall[SnapshotCleanupResult](c, callInput)
}

// ✅ REQUIRED: enum constants + Result[T] return
func (c *Client) CleanupSnapshots(opts SnapshotCleanupOptions) apperror.Result[SnapshotCleanupResult] {
    callInput := apiCallInput{
        Method:    httpmethod.Post,
        Operation: snapshotoperationtype.Cleanup,
    }

    return apiCallTo[SnapshotCleanupResult](c, callInput)
}
```

---

## DRY Enforcement

| Pattern | Solution |
|---------|----------|
| Repeated error handling | `apperror.Result[T]` or helper functions |
| Repeated JSON key access | Typed response structs |
| Repeated validation | `Validate()` method on input structs |
| Repeated DB patterns | `dbutil` generic wrappers |
| Repeated string constants | Typed const blocks with `Type` suffix |

---
