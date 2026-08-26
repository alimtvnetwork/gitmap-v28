# DRY Refactoring Project — Complete Summary

> **Status:** 100% Complete  
> **Duration:** 10 Phases  
> **Updated:** 2026-03-09  
**Version:** 3.2.0  

---

## Overview

The 10-phase DRY (Don't Repeat Yourself) refactoring initiative modernized the full stack — Go backend, React/TypeScript frontend, and PHP WordPress companion plugin — by eliminating duplication, establishing shared abstractions, and creating a canonical cross-stack contract via JSON Schema.

---

## Phase Summary

### Phase 1: Go Uploader Method Deduplication

**Goal:** Eliminate duplicated upload logic across multiple uploader methods.

- Consolidated `UploadPluginViaOnboard` and `UploadPlugin` into a single `UploadPlugin` method
- Modernized string utilities (removed redundant helper functions)
- Redirected all legacy Onboard API calls to Riseup Asia Uploader endpoints

**Files affected:** `backend/internal/wordpress/uploader.go`

---

### Phase 2: Shared Go Helpers

**Goal:** Extract common Go utilities into reusable helpers.

- Created `namespace.go` — centralized `resolveNamespace` utility for uploader version detection
- Created `php_stack.go` — shared PHP stack trace parser for formatting remote error frames
- Both helpers used by all WordPress-delegated operations

**Files affected:** `backend/internal/wordpress/namespace.go`, `backend/internal/wordpress/php_stack.go`

---

### Phase 3: Frontend Diagnostic Context Builders

**Goal:** Deduplicate error context construction in the React API client.

- Extracted diagnostic context builders from inline error handling
- Centralized call-chain construction and stack frame parsing
- Created reusable `buildDiagnosticContext()` utility

**Files affected:** `src/lib/api/` (client module)

---

### Phase 4: Frontend Error Store Consolidation

**Goal:** Unify error handling logic in the global error store.

- Consolidated scattered error handling patterns into the Zustand error store
- Standardized error enrichment (component → action → trigger context)
- Established single entry point for error reporting

**Files affected:** `src/stores/errorStore.ts`

---

### Phase 5: API Client Modularization

**Goal:** Break the monolithic `api.ts` into focused modules.

- Split into `src/lib/api/` directory:
  - `types.ts` — Shared interfaces (envelope types, API response shapes)
  - `envelope.ts` — Universal Response Envelope parsing
  - `client.ts` — Core HTTP logic, circuit-breaker, retry
  - `methods.ts` — Endpoint-specific implementations
  - `index.ts` — Barrel export for backward compatibility

**Files affected:** `src/lib/api/*`

---

### Phase 6: `useApiQuery` Factory Hook

**Goal:** Standardize all data fetching via a single factory.

- Created `useApiQuery` and `useApiQueryPaginated` in `src/hooks/useApiQuery.ts`
- Encapsulates `useQuery` + `requireSuccess` pattern
- Auto-handles error diagnostics, query key standardization, and pagination metadata
- All entity hooks (Plugins, Sites, Errors, Settings) migrated to use factory

**Files affected:** `src/hooks/useApiQuery.ts`, all entity hooks

---

### Phase 7: PHP `SnapshotFactory`

**Goal:** Deduplicate snapshot creation logic in the WordPress plugin.

- Extracted snapshot creation patterns into `SnapshotFactory` class (`RiseupAsia\Snapshot\SnapshotFactory`)
- Unified native SQLite, WP Reset, and UpdraftPlus snapshot strategies
- Standardized progress tracking and error handling across all snapshot types

**Files affected:** `wp-plugins/riseup-asia-uploader/includes/Snapshot/SnapshotFactory.php`

---

### Phase 8: PHP Logger Context Enrichment

**Goal:** Centralize diagnostic metadata injection in PHP logging.

- Extracted context enrichment (backtrace capture, request metadata, memory usage) into shared methods
- All `error()` and `log_exception()` calls automatically receive 6-frame backtrace + request context
- Eliminated manual context construction scattered across endpoint handlers

**Files affected:** `wp-plugins/riseup-asia-uploader/includes/Logging/FileLogger.php`

---

### Phase 9: `GlobalErrorModal.tsx` Decomposition

**Goal:** Break the monolithic error modal into maintainable sub-components.

- Decomposed into modular shell with focused sub-components:
  - `BackendSection` — Go-side error display
  - `FrontendSection` — React-side error display
  - `RequestDetails` — HTTP request/response visualization
  - `TraversalDetails` — React→Go→PHP request chain visualization
  - `ErrorModalActions` — Download bundle, copy, close actions
- Extracted pure logic to `errorReportGenerator.ts`

**Files affected:** `src/components/errors/GlobalErrorModal.tsx` → `src/components/errors/modal/*`

---

### Phase 10: Cross-Stack JSON Schema Alignment

**Goal:** Establish a machine-readable contract for the Universal Response Envelope.

- Created `spec/response-envelope/envelope.schema.json` (JSON Schema Draft 2020-12, v1.0.0)
- Pinned Go (`envelope.go`), TypeScript (`types.ts`), and PHP (`EnvelopeBuilder.php`) to schema v1.0.0
- Added `@schema` reference comments to all three implementations
- Validated all reference samples against the schema
- Documented decision in ADR #7

**Files affected:** `spec/response-envelope/envelope.schema.json`, `spec/response-envelope/adr.md`, `spec/response-envelope/changelog.md`

---

## Impact Summary

| Metric | Before | After |
|--------|--------|-------|
| Duplicated upload methods | 3 | 1 |
| API client files | 1 monolith | 5 focused modules |
| Error modal components | 1 (800+ lines) | 6 sub-components |
| Data fetching patterns | Ad-hoc per hook | `useApiQuery` factory |
| Cross-stack type contract | Prose README | Machine-readable JSON Schema |
| PHP snapshot patterns | Inline per handler | `SnapshotFactory` |
| PHP logger context | Manual per call | Auto-enriched centrally |

---

## Key Principles Established

1. **Single source of truth** — `envelope.schema.json` for cross-stack types
2. **Factory hooks** — `useApiQuery` for all frontend data fetching
3. **Modular decomposition** — No component or module > 300 lines
4. **Centralized error handling** — Single error store, single enrichment pipeline
5. **Lazy dependency resolution** — Breaks circular dependencies in Go and PHP

---

## Cross-References

- Response Envelope Spec <!-- external: spec/03-error-manage/01-error-resolution/09-response-envelope/04-response-envelope-reference.md -->
- Envelope JSON Schema <!-- external: spec/03-error-manage/01-error-resolution/09-response-envelope/envelope.schema.json -->
- ADR #7: JSON Schema Strategy <!-- external: spec/03-error-manage/01-error-resolution/09-response-envelope/01-adr.md -->
- Response Envelope Overview <!-- external: spec/03-error-manage/01-error-resolution/09-response-envelope/00-overview.md -->

---

*10-phase DRY refactoring completed 2026-02-09*
