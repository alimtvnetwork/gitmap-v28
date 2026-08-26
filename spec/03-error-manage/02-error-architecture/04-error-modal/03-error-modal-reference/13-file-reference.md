# File Reference

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

| File | Lines | Purpose |
|------|-------|---------|
| `src/stores/errorStore.ts` | 674 | Central error store (Zustand): CapturedError type, capture pipeline, queue management |
| `src/components/errors/GlobalErrorModal.tsx` | 239 | Modal shell: header, section toggle, footer actions |
| `src/components/errors/BackendSection.tsx` | 721 | Backend tabs: Overview, Log, Execution, Stack, Session, Request, Traversal |
| `src/components/errors/FrontendSection.tsx` | 328 | Frontend tabs: Overview, Stack, Context, Fixes |
| `src/components/errors/RequestDetails.tsx` | 176 | Request chain visualization (React → Go → Delegated, 3-hop) |
| `src/components/errors/TraversalDetails.tsx` | 149 | Envelope traversal: endpoint flow, methods stack, delegated server details, error stacks |
| `src/components/errors/SessionLogsTab.tsx` | 443 | Session diagnostics: logs, request, response, stack traces |
| `src/components/errors/ErrorModalActions.tsx` | 194 | Download and Copy dropdown menus |
| `src/components/errors/errorReportGenerator.ts` | 185 | Pure function: Markdown report generation (compact + full) + suggested fixes |
| `src/components/errors/errorLogAdapter.ts` | 22 | Maps backend `ErrorLog` → `CapturedError` for ErrorDetailModal report generation |
| `src/components/errors/ErrorModalTypes.ts` | 26 | Shared types: PHPStackFrame, AppInfo, SectionCommonProps, DelegatedRequestServer |
| `src/components/errors/ErrorDetailModal.tsx` | — | Standalone error detail viewer (Split Button copy + DownloadDropdown) |
| `src/components/errors/ErrorHistoryDrawer.tsx` | — | Error history browser drawer |
| `src/components/errors/ErrorQueueBadge.tsx` | — | Error queue indicator badge |
| `src/components/errors/AppErrorBoundary.tsx` | — | React error boundary wrapping the app |
| `src/hooks/useSessionDiagnostics.ts` | — | Hook for auto-fetching session-level diagnostics |
| `src/hooks/useClickTracker.ts` | — | Click path tracking for error context |
| `src/hooks/useExecutionLogger.ts` | — | React execution logger (debug mode) |
| `src/lib/api/envelope.ts` | — | Envelope parsing and error extraction (incl. DelegatedRequestServer) |

---

## Cross-References

- [Error Handling Cross-Stack Spec](../../01-error-handling-reference.md) — PHP, Go, frontend error chain + DelegatedRequestServer flow
- [Copy Format Samples](../01-copy-formats/00-overview.md) — Complete samples for all copy/export formats
- [React Components Reference](../02-react-components/00-overview.md) — Portable React code for rebuilding the modal
- [Response Envelope Schema](../../05-response-envelope/envelope.schema.json) — JSON Schema for envelope (incl. DelegatedRequestServer)
- [Session-Based Logging](../../07-logging-and-diagnostics/02-session-based-logging.md) — Backend session system
- [React Execution Logger](../../07-logging-and-diagnostics/01-react-execution-logger.md) — Frontend debug logger
- [TypeScript Standards](../../../../02-coding-guidelines/02-typescript/08-typescript-standards-reference.md) — Type safety rules

---

*File reference — updated: 2026-03-31*
