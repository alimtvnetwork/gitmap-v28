# Error Modal — Frontend Specification (Index)

> **Parent:** [Error Modal Spec](../00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31  
> **Status:** Active  
> **Location:** `src/components/errors/`  
> **AI Confidence:** 95%  
> **Ambiguity Score:** 5%  
> **Purpose:** Comprehensive specification for the Global Error Modal — how errors are captured, enriched, displayed, and exported across the React → Go → Delegated Server request chain.

---

## File Index

| # | File | Section | Lines |
|---|------|---------|-------|
| 01 | [01-data-model.md](./01-data-model.md) | CapturedError interface + supporting types | ~130 |
| 02 | [02-capture-pipeline.md](./02-capture-pipeline.md) | Error capture: API client → store → modal | ~85 |
| 03 | [03-envelope-parsing.md](./03-envelope-parsing.md) | Envelope parsing, Errors/MethodsStack/Attributes mapping | ~85 |
| 04 | [04-modal-structure.md](./04-modal-structure.md) | Component hierarchy + visual layout diagrams | ~285 |
| 05 | [05-backend-tabs.md](./05-backend-tabs.md) | Backend section tabs: Overview, Log, Execution, Stack, Session, Request, Traversal | ~80 |
| 06 | [06-frontend-tabs.md](./06-frontend-tabs.md) | Frontend section tabs: Overview, Stack, Context, Fixes | ~25 |
| 07 | [07-request-chain.md](./07-request-chain.md) | Request chain visualization (3-hop React→Go→Delegated) | ~115 |
| 08 | [08-traversal-details.md](./08-traversal-details.md) | Traversal tab: endpoint flow, methods stack, delegated details | ~70 |
| 09 | [09-session-diagnostics.md](./09-session-diagnostics.md) | Session diagnostics auto-fetch + SessionDiagnostics shape | ~55 |
| 10 | [10-report-generation.md](./10-report-generation.md) | Error report generators (compact + full) + copy/download menus | ~75 |
| 11 | [11-queue-navigation.md](./11-queue-navigation.md) | Error queue navigation (multi-error support) | ~30 |
| 12 | [12-code-examples.md](./12-code-examples.md) | React code examples for error capture, modal, boundary | ~135 |
| 13 | [13-file-reference.md](./13-file-reference.md) | File reference table + cross-references | ~50 |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Error Capture Flow                              │
│                                                                         │
│  User Action ──▸ API Call ──▸ Go Backend ──▸ PHP (WordPress)            │
│       │              │              │              │                     │
│       │              │              │              ▼                     │
│       │              │              │     ┌──────────────────┐          │
│       │              │              │     │ PHP Error         │          │
│       │              │              │     │ - stackTrace      │          │
│       │              │              │     │ - stackTraceFrames│          │
│       │              │              │     │ - fatal-errors.log│          │
│       │              │              │     └────────┬─────────┘          │
│       │              │              │              │                     │
│       │              │              ▼              │                     │
│       │              │     ┌──────────────────┐    │                     │
│       │              │     │ Go Error Handler  │◀──┘                     │
│       │              │     │ - apperror.Wrap() │                         │
│       │              │     │ - session logger  │                         │
│       │              │     │ - envelope builder│                         │
│       │              │     └────────┬─────────┘                         │
│       │              │              │                                    │
│       │              ▼              │ Universal Response Envelope        │
│       │     ┌──────────────────┐    │                                    │
│       │     │ API Client       │◀───┘                                    │
│       │     │ - parseEnvelope()│                                         │
│       │     │ - extract errors │                                         │
│       │     └────────┬─────────┘                                        │
│       │              │                                                   │
│       ▼              ▼                                                   │
│  ┌──────────────────────────────────┐                                   │
│  │        Error Store (Zustand)      │                                   │
│  │  - buildCapturedError()           │                                   │
│  │  - commitErrorToStore()           │                                   │
│  │  - enrich with click path         │                                   │
│  │  - enrich with execution logs     │                                   │
│  └────────────────┬─────────────────┘                                   │
│                   │                                                      │
│                   ▼                                                      │
│  ┌──────────────────────────────────┐                                   │
│  │       Global Error Modal          │                                   │
│  │  ┌─────────┐ ┌─────────────────┐ │                                   │
│  │  │ Backend │ │    Frontend     │ │                                   │
│  │  │ Section │ │    Section      │ │                                   │
│  │  └─────────┘ └─────────────────┘ │                                   │
│  │  ┌─────────────────────────────┐ │                                   │
│  │  │ Download/Copy Actions       │ │                                   │
│  │  └─────────────────────────────┘ │                                   │
│  └──────────────────────────────────┘                                   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Document Inventory

| File |
|------|
| 99-consistency-report.md |


## Cross-References

- [Copy Format Samples](../01-copy-formats/00-overview.md) — Complete samples for all copy/export formats
- [React Components Reference](../02-react-components/00-overview.md) — Portable React code for rebuilding the modal
- [Color Themes](../04-color-themes/00-overview.md) — Color mapping for all error UI elements
- [Response Envelope Schema](../../05-response-envelope/envelope.schema.json) — JSON Schema for envelope
- [Error Handling Spec](../../01-error-handling-reference.md) — Cross-stack error architecture
- [Session-Based Logging](../../07-logging-and-diagnostics/02-session-based-logging.md) — Backend session system
- [React Execution Logger](../../07-logging-and-diagnostics/01-react-execution-logger.md) — Frontend debug logger

---

*Error Modal specification index — updated: 2026-03-31*
