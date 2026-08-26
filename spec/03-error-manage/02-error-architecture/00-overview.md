# Error Architecture

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`error-architecture` · `error-handling` · `error-modal` · `response-envelope` · `apperror` · `logging` · `notifications` · `delegation`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Cross-stack error handling architecture spanning React → Go → Delegated Server (PHP/3rd-party). Covers the three-tier error flow, error modal specification, response envelope format, Go `apperror` package, logging/diagnostics, and notification color tokens.

---

## Document Inventory

### Root Files

| # | File | Purpose |
|---|------|---------|
| 01 | [01-error-handling-reference.md](./01-error-handling-reference.md) | Cross-stack 3-tier error flow architecture |
| 02 | [02-go-delegation-fix.md](./02-go-delegation-fix.md) | DelegatedRequestServer implementation pattern |
| 03 | [03-notification-colors.md](./03-notification-colors.md) | Toast/notification color tokens & error code mapping |
| — | 99-consistency-report.md | — |

### Subfolders

| # | Folder | Description | Files |
|---|--------|-------------|-------|
| 04 | [04-error-modal/](./04-error-modal/00-overview.md) | Frontend Global Error Modal specification | 6 |
| 05 | [05-response-envelope/](./05-response-envelope/00-overview.md) | Universal Response Envelope spec + schema | 4 + JSON samples |
| 06 | [06-apperror-package/](./06-apperror-package/00-overview.md) | Go structured error package specification | 1 |
| 07 | [07-logging-and-diagnostics/](./07-logging-and-diagnostics/00-overview.md) | React execution logger + session-based logging | 2 |

| — | 99-consistency-report.md | — |
---

## Three-Tier Architecture Summary

```
Tier 1: Delegated Server (PHP/other) → structured error responses, stack traces
Tier 2: Go Backend → apperror package, DelegatedRequestServer, session logging
Tier 3: Frontend (React) → Error store, Global Error Modal, toast notifications
```

---

## Cross-References

- [Parent Overview](../00-overview.md) — Error Management root
- [Error Resolution](../01-error-resolution/00-overview.md) — Debugging and diagnostics
- [Error Code Registry](../03-error-code-registry/00-overview.md) — Error code ranges
