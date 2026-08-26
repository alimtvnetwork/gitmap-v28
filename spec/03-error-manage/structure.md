# 04-error-manage — Structure Summary

**Version:** 3.2.0  
**Generated:** 2026-03-31  
**Total Files:** 71 (62 Markdown + 9 JSON)  
**Health Score:** 100/100 (A+)

---

## Quick Navigation

```
04-error-manage/
├── 00-overview.md                        ← START HERE
├── 97-acceptance-criteria.md
├── 98-changelog.md
├── 99-consistency-report.md
├── structure.md                          ← THIS FILE
│
├── 01-error-resolution/                  (14 files)
│   ├── 00-overview.md
│   ├── 01-cross-reference-diagram.md
│   ├── 02-debugging-cheat-sheet.md
│   ├── 03-retrospectives/               (4 retrospectives)
│   │   ├── 01-health-endpoint-mismatch.md
│   │   ├── 02-retry-debounce-dedup-fixes.md
│   │   ├── 03-zip-finalization-before-return.md
│   │   └── 04-activation-endpoint-mismatch.md
│   ├── 04-verification-patterns/        (1 pattern)
│   │   └── 01-frontend-backend-sync.md
│   └── 05-debugging-guides/             (3 guides)
│       ├── 01-debugging-php.md
│       ├── 02-debugging-go.md
│       └── 03-debugging-typescript.md
│
├── 02-error-architecture/                (22 files + 6 JSON)
│   ├── 00-overview.md
│   ├── 01-error-handling-reference.md     ← 3-tier error flow
│   ├── 02-go-delegation-fix.md            ← DelegatedRequestServer
│   ├── 03-notification-colors.md          ← Toast color tokens
│   ├── 04-error-modal/                   (6 docs)
│   │   ├── 01-copy-formats.md
│   │   ├── 02-react-components.md
│   │   ├── 03-error-modal-reference.md
│   │   ├── 04-color-themes.md
│   │   ├── 05-error-history-persistence.md
│   │   └── 06-suppress-global-error.md
│   ├── 05-response-envelope/            (4 docs + 6 JSON samples + schema)
│   │   ├── 01-adr.md
│   │   ├── 02-changelog.md
│   │   ├── 03-configurability.md
│   │   ├── 04-response-envelope-reference.md
│   │   ├── envelope.schema.json
│   │   ├── envelope-{minimal,single,multiple,error,debug}.json
│   ├── 06-apperror-package/             (2 docs)
│   │   ├── 01-apperror-reference.md
│   │   └── 01-apperror-reference/
│   │       ├── 00-overview.md
│   │       ├── 01-overview-and-stack.md
│   │       ├── 02-apperror-struct.md
│   │       ├── 03-result-types.md
│   │       ├── 04-codes-and-policy.md
│   │       ├── 05-apperrtype-enums.md    ← NEW: Domain error type enums
│   │       ├── 05-usage-and-adapters.md
│   │       ├── 06-serialization-and-guards.md
│   │       └── 99-consistency-report.md
│   └── 07-logging-and-diagnostics/      (2 docs)
│       ├── 01-react-execution-logger.md
│       └── 02-session-based-logging.md
│
└── 03-error-code-registry/               (18 files + 3 JSON)
    ├── 00-overview.md
    ├── 01-registry.md                     ← Master error code list
    ├── 02-integration-guide.md
    ├── 03-collision-resolution-summary.md
    ├── 04-error-code-utilization-report.md
    ├── 05-overlap-validator.md
    ├── error-codes-master.json
    ├── 07-schemas/                       (2 JSON schemas)
    │   ├── error-code.schema.json
    │   └── error-codes-index.schema.json
    ├── 08-linter-scripts/                (empty — placeholder)
    └── 09-templates/                     (1 template)
        └── 01-error-codes-template.md
```

---

## Category Summary

| # | Category | Purpose | Docs | Assets |
|---|----------|---------|------|--------|
| 01 | **Error Resolution** | Debugging guides, retrospectives, verification patterns, cheat sheet | 14 | — |
| 02 | **Error Architecture** | 3-tier error flow, error modal, response envelope, apperror, logging | 22 | 6 JSON |
| 03 | **Error Code Registry** | Master registry, schemas, scripts, templates, collision resolution | 18 | 3 JSON |

---

## Three-Tier Architecture (Core Concept)

```
Tier 1: Delegated Server (PHP/other) → structured error responses, stack traces
Tier 2: Go Backend                   → apperror package, DelegatedRequestServer, session logging
Tier 3: Frontend (React)             → Error store, Global Error Modal, toast notifications
```

---

## Key Entry Points by Role

| I need to… | Start at |
|------------|----------|
| Debug an error | `01-error-resolution/02-debugging-cheat-sheet.md` |
| Understand the error flow | `02-error-architecture/01-error-handling-reference.md` |
| Build/modify the error modal | `02-error-architecture/04-error-modal/03-error-modal-reference.md` |
| Format API responses | `02-error-architecture/05-response-envelope/04-response-envelope-reference.md` |
| Add a new error code | `03-error-code-registry/01-registry.md` + `09-templates/01-error-codes-template.md` |
| Review a past incident | `01-error-resolution/03-retrospectives/` |
| Validate error code ranges | `03-error-code-registry/05-overlap-validator.md` |

---

## Consolidation Source

This spec supersedes three archived folders:

| Archived Source | Content Merged Into |
|----------------|-------------------|
