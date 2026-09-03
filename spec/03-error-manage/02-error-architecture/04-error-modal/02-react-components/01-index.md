# Error Modal — Reusable React Components (Index)

. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


> **Parent:** [Error Modal Spec](../01-index.md)
> **Version:** 4.0.0
> **Updated:** 2026-04-01
> **AI Confidence:** 95%
> **Ambiguity Score:** 5%
> **Purpose:** Portable React code for rebuilding the Global Error Modal in any project.

---

## File Index

| # | File | Section | Description |
|---|------|---------|-------------|
| 01 | [01-typescript-interfaces.md](./02-typescript-interfaces.md) | TypeScript Interfaces | CapturedError, SessionDiagnostics, shared props |
| 02 | [02-error-store.md](./03-error-store.md) | Error Store (Zustand) | Store interface, key behaviors, stack trace parser |
| 03 | [03-api-types.md](./04-api-types.md) | API Types & Methods | Required API endpoints |
| 04 | [04-hooks.md](./05-hooks.md) | Hooks | useSessionDiagnostics |
| 05 | [05-component-hierarchy.md](./06-component-hierarchy.md) | Component Hierarchy | File structure + component props summary |
| 06 | [06-component-source.md](./07-component-source.md) | Component Source Code | All 7 major components with code patterns |
| 07 | [07-report-generator.md](./08-report-generator.md) | Error Report Generator | generateErrorReport + suggested fixes |
| 08 | [08-integration-guide.md](./09-integration-guide.md) | Integration Guide | Setup, React Query, utilities, adaptation |

---

## Architecture Overview

```
GlobalErrorModal (Dialog shell)
├── Header (error code, timestamp, queue navigation)
├── Section Toggle: Backend | Frontend
├── BackendSection (primary diagnostic view)
│   ├── Overview Tab
│   ├── Log Tab (error.log.txt viewer)
│   ├── Execution Tab (Go call chain + backend logs)
│   ├── Stack Tab (Go/PHP/Delegated stack frames)
│   ├── Session Tab (SessionLogsTab — 4 sub-tabs)
│   ├── Request Tab (RequestDetails — 3-hop chain)
│   └── Traversal Tab (TraversalDetails — endpoint flow)
├── FrontendSection
│   ├── Overview Tab (trigger, click path, call chain)
│   ├── Stack Tab (parsed/raw JS stack frames)
│   ├── Context Tab (JSON viewer)
│   └── Fixes Tab (suggested fixes by error code)
├── Footer
│   ├── DownloadDropdown (ZIP, error.log, log.txt, .md)
│   └── CopyDropdown (full report, with backend, logs)
```

**Dependencies:** React 18+, Zustand, Tailwind CSS, shadcn/ui (Dialog, Tabs, Badge, Button, ScrollArea, DropdownMenu), Lucide React icons.

---

## Document Inventory

| File |
|------|
| 99-consistency-report.md |


## Cross-References

- [Error Modal Spec](../03-error-modal-reference/01-index.md) — Full modal structure, data model, and UX specification
- [Copy Format Samples](../01-copy-formats/01-index.md) — Complete samples for all copy/export formats
- [Error Handling Spec](../../02-error-handling-reference.md) — Cross-stack error architecture
- [Response Envelope Schema](../../05-response-envelope/envelope.schema.json) — JSON Schema source of truth

---

*React components index — updated: 2026-03-31*
