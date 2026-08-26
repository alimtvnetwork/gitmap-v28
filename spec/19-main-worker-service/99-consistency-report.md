# 99 — Consistency Report

> **⚠ STALE — historical snapshot.** This file's cross-link claims pre-date Phase 13.2 hardening (WorkerNode dedupe, endpoint-catalogue dedupe, FK repoint, tunable alias map, 11- stub). For current consistency status see `audit/07-blind-ai-readiness-2026-05-06-v2.md` and re-run `npm run spec:links:check`.

**Spec:** `19-main-worker-service`
**Version:** 1.1.0
**Verified:** 2026-05-04

Cross-link integrity, naming-convention compliance, and rule-inheritance audit for this spec folder.

---

## 1. Files in this spec

### 1.1 Root spec files

| File | Purpose | Version | Status |
|------|---------|---------|--------|
| `plan.md` | Phased roadmap + decisions | 1.0.0 | ✅ |
| `00-overview.md` | Entry point, document map | 1.1.0 | ✅ |
| `01-architecture.md` | Topology + comms contract | 1.1.0 | ✅ |
| `02-glossary.md` | Canonical terms (5 entries added: Quarantined, Draining, Seedable-Config, apperror, Power Admin↔PowerAdmin) | 1.1.0 | ✅ |
| `03-main-db-schema.md` | Main SQLite schema (+EnumPage, RolePageAccess, AccessDenialEvent, MainSetting, User TOTP cols, **+EndpointAuthAuditEvent §2.6.4 / EndpointAuthChangeKind §2.6.5 — FU-17**) | 1.3.0 | ✅ |
| `04-worker-routing.md` | Selection + cache + failover + strategy interfaces | 1.1.0 | ✅ |
| `05-auth-and-2fa.md` | Auth surfaces + 2FA + bcrypt-cost env pinning + backup-code regen + EndpointAuth audit cite | 1.1.0 | ✅ |
| `06-core-api-endpoints.md` | REST surface + per-field validation + LeastLoaded tiebreaker patch + §5.6 audit wired to `EndpointAuthAuditEvent` (FU-17) | 1.2.0 | ✅ |
| `07-role-based-dashboards.md` | Roles + EnumPage + stack-agnostic guard contract | 1.1.0 | ✅ |
| `08-error-contract.md` | Main↔Worker envelope + EnvelopeVersion/OperationId/SubCode/FieldErrors + Worker→Main §9 | 1.1.0 | ✅ |
| `09-self-update-pointer.md` | Pointer to `spec/14-update/` + bounded sunset + deletion checklist | 1.2.0 | ✅ |
| `10-worker-bootstrap-protocol.md` | 8-step boot, /Workers/Register, public-key fetch, version pin | 1.0.0 | ✅ |
| `11-split-db-tier-reconciliation.md` | Main=3-tier / Worker=4-tier mapping over spec/05's 6 tiers | 1.0.0 | ✅ |
| `12-jwt-delivery-contract.md` | JSON-body + in-memory storage, mandatory CSP, 9 CI tests | 1.0.0 | ✅ |
| `13-error-codes.md` | 31 codes (22 WORKER-* + 9 MAIN-*) + prefixed↔flat mapping + MWS prefix 21000-21199 + slot-overflow rule (`MAIN-400-10 EndpointAuthLocked` / `21170` per FU-18) | 1.1.0 | ✅ |
| `14-rbac-and-status-seed.md` | 3 Roles + 9 EnumPages + 19 RolePageAccess + 4 WorkerNodeStatus + 4 AuthMechanism | 1.0.0 | ✅ |
| `15-tunable-constants.md` | 30 numeric tunables (retry, idempotency, heartbeat, JWT, push-update, bootstrap) + config.seed.json binding | 1.1.0 | ✅ |
| `96-linter-audit.md` | Linter pipeline audit | 1.0.0 | ✅ |
| `97-acceptance-criteria.md` | AC-1..AC-9 mapping | 1.0.0 | ✅ |
| `98-changelog.md` | Version history | 1.1.0 | ✅ |
| `99-consistency-report.md` | This file | 1.1.0 | ✅ |
| `error-codes.json` | Machine-readable mirror of `13-error-codes.md` | — | ✅ |

### 1.2 Diagrams (`diagrams/`)

All 6 diagrams carry the **NON-AUTHORITATIVE PROJECTION** banner (v1.0.0); spec wins on conflict.

| File | Purpose | Status |
|------|---------|--------|
| `diagrams/erd-main-db.mmd` | Main ERD (synced to schema v1.2.0; +EnumPage, +AccessDenialEvent, +User TOTP, RolePageAccess upgraded to FK) — banner v1.1.0 | ✅ |
| `diagrams/erd-worker-split-db.mmd` | Worker split-DB ERD (projection of spec/05) | ✅ |
| `diagrams/erd-seedable-config.mmd` | Seedable-Config ERD | ✅ |
| `diagrams/seq-company-creation.mmd` | Sequence | ✅ |
| `diagrams/seq-login-routing.mmd` | Sequence (banner v1.1.0 — synced X-Auth-Action: Reauthenticate + RegenerateBackupCodes per tasks #29/#30/#38) | ✅ |
| `diagrams/seq-push-update.mmd` | Sequence | ✅ |
| `diagrams/readme.md` | Diagrams index + per-file authority table | ✅ |

### 1.3 Author mindmap images (`images/`)

| File | Purpose | Status |
|------|---------|--------|
| `images/01-main-worker-topology.png` | Author mindmap | ✅ |
| `images/02-endpoint-service-worker-pattern.png` | Author mindmap | ✅ |
| `images/03-worker-subdomain-routing.png` | Author mindmap | ✅ |
| `images/04-endpoint-service-full-overview.png` | Author mindmap | ✅ |
| `images/readme.md` | Image index | ✅ |

### 1.4 Audit suite (`audit/`)

| File | Purpose | Status |
|------|---------|--------|
| `audit/01-completeness-audit.md` | Step 1/5 — 30 findings, **30/30 closed** (28 fixed + 2 deferred), §7 re-triage v1.1.0 | ✅ |
| `audit/02-ambiguity-audit.md` | Step 2/5 — 69 findings, ~30 closed via tasks #29–31 | ⚠ FU |
| `audit/03-diagram-audit.md` | Step 3/5 — 109 findings, banner cluster closed via task #15 | ⚠ FU |
| `audit/04-cross-spec-dependency-audit.md` | Step 4/5 — 20 findings, anchor sweep verified clean (task #33) | ✅ |
| `audit/05-implementation-pivot-score.md` | Step 5/5 — original pivot 66%; refresh pending | ⚠ FU |

**Total: 39 files** (22 root + 7 diagrams + 5 images + 5 audit).


---

## 2. Cross-spec references (must exist)

| Reference | Verified |
|-----------|----------|
| `spec/03-error-manage/` | ✅ exists |
| `spec/04-database-conventions/` | ✅ exists |
| `spec/05-split-db-architecture/` | ✅ exists |
| `spec/06-seedable-config-architecture/` | ✅ exists |
| `spec/14-update/` | ✅ exists |
| `.lovable/coding-guidelines/coding-guidelines.md` | ✅ exists |
| `mem://architecture/error-handling` | ✅ in memory index |
| `mem://architecture/caching-policy` | ✅ in memory index |
| `mem://features/self-update-architecture` | ✅ in memory index |
| `mem://project/author-attribution` | ✅ in memory index |

---

## 3. Naming convention compliance

| Rule | Status | Notes |
|------|--------|-------|
| PascalCase for tables, columns, JSON keys | ✅ | All schema in `03-`, `06-`, `07-` follows |
| PKs `{TableName}Id INTEGER AUTOINCREMENT` | ✅ | Every table |
| No UUIDs as PKs | ✅ | Correlation IDs are opaque request strings, NOT keys (`08-` §4 explicit) |
| Entity/ref tables include `Description TEXT NULL` | ✅ | All applicable tables |
| Transactional tables include `Notes` + `Comments TEXT NULL` | ✅ | `WorkerVersion`, `WorkerSelectionEvent`, `AppBusinessEntity`, `AppSession`, `SeedableConfigVersion` |
| Join tables exempt | ✅ | `UserRole`, `RolePageAccess` correctly omit `Description`/`Notes` |
| Type/Status/Category/Kind via join tables | ✅ | `WorkerNodeStatus`, `WorkerNodeKind`, `WorkerSelectionStrategy`, `Role`, `RootCompanyStatus`, `AuthMechanism` |
| Files numbered `NN-kebab-case.md` | ✅ | All conform |

---

## 4. Forbidden-term audit (per `02-glossary.md`)

| Forbidden | Replacement | Found in this spec? |
|-----------|-------------|----------------------|
| `CW configuration` | `Seedable-Config` | NO ✅ |
| `git map` | `gitmap` | NO ✅ (term not used; would use `gitmap` if needed) |
| `Master/Slave` | `Main/Worker` | NO ✅ |
| `CEO` for Md. Alim Ul Karim | `Chief Software Engineer` | N/A — author not referenced in this spec |

**Reproducible verification commands** (closes audit finding F-N-03; run from repo root):

```bash

# 4.1 — forbidden literals (case-insensitive). Excludes meta-docs that legitimately

# define/cite the replacements: this report, the glossary, plan, changelog, AC.

rg -i -n \
   --glob '!spec/19-main-worker-service/02-glossary.md' \
   --glob '!spec/19-main-worker-service/99-consistency-report.md' \
   --glob '!spec/19-main-worker-service/98-changelog.md' \
   --glob '!spec/19-main-worker-service/97-acceptance-criteria.md' \
   --glob '!spec/19-main-worker-service/plan.md' \
   -e 'CW configuration' -e 'git map' \
   -e '\bmaster/slave\b' -e '\bmaster-slave\b' \
   spec/19-main-worker-service/

# 4.2 — author-title compliance (must NOT appear). Same meta-doc exclusions.

rg -n \
   --glob '!spec/19-main-worker-service/02-glossary.md' \
   --glob '!spec/19-main-worker-service/99-consistency-report.md' \
   -e '\bCEO\b' -e 'Chief Executive Officer' \
   spec/19-main-worker-service/

# 4.3 — global enforcement (delegated to the canonical linter)

python3 linter-scripts/check-forbidden-strings.py
```

**Expected:** 4.1 and 4.2 print nothing and exit 1 (rg "no matches" convention); 4.3 exits 0. Anything else is a regression.

Last verified: 2026-05-04 — 4.1: 0 matches in non-meta files (`01-architecture.md` §6 still references "Seedable-Config" by its canonical name only); 4.2: 0 matches; 4.3: exit 0.

---

## 5. CODE RED compliance (rules applied to spec content)

| Rule | Applied? |
|------|----------|
| Function pseudocode ≤ 8 lines (best practice), ≤ 15 (cap) | ✅ All pseudocode in `04-`, `05-`, `07-`, `08-` ≤ 15 lines |
| Zero nested `if` in pseudocode | ✅ All examples use early-return guards |
| Positive boolean names (`is`/`has` prefix, no `!`) | ✅ `IsWorkerActive`, `HasCapacity`, `userHasAccessToPage`, `canRetry`, `isSuccessful`, `isPermanentFailure` |
| Max 2 operands per condition | ✅ All examples |
| No magic strings — Enums for Type/Status/Category/Kind | ✅ `EnumPage`, `WorkerNodeStatusCode`, `WorkerSelectionStrategyCode`, `Cadence` |

---

## 6. Inheritance audit (rules deferred, not duplicated)

| Topic | Authoritative source | This spec duplicates? |
|-------|----------------------|-----------------------|
| Split-DB internals | `spec/05-split-db-architecture/` | NO ✅ — only references and projection diagram |
| Seedable-Config internals | `spec/06-seedable-config-architecture/` | NO ✅ |
| Generic error rules | `spec/03-error-manage/` | NO ✅ — `08-` adds only Main↔Worker patterns |
| Self-update mechanism | `spec/14-update/` | NO ✅ — `09-` is pointer-only |
| DB schema conventions | `spec/04-database-conventions/` | NO ✅ — `03-` references and applies |
| Coding metrics | `.lovable/coding-guidelines/coding-guidelines.md` | NO ✅ — applied via examples |

---

## 7. Open Questions (carried, not closed)

| ID | Topic | Location | Status |
|----|-------|----------|--------|
| OQ-1 | Per-endpoint auth-mechanism overrides | `05-auth-and-2fa.md` §8, `06-core-api-endpoints.md` §5 | ✅ RESOLVED 2026-05-04 (task #39) — single-row whole-replace PATCH; 06 → v1.1.0 |
| OQ-2 | Default selection strategy (RoundRobin vs LeastLoaded) | `04-worker-routing.md` §1 | RESOLVED — default `LeastLoaded`; user can change via Settings |

---

## 8. Acceptance Criteria coverage

All 9 verbatim acceptance criteria mapped in `97-acceptance-criteria.md`. Coverage matrix verified — every sub-criterion has at least one defining spec file and at least one test condition.

---

## 9. Outstanding (next phase)

- **Phase 5** (separate task): bump `package.json` minor, run `node scripts/sync-version.mjs` and `node scripts/sync-spec-tree.mjs`, log entry in `.lovable/plan.md`.
- Future: implement Phase 1 of the deliverables when user issues `next` after Phase 5 (or explicitly says "implement").

---

*Consistency report v1.2.0 — 2026-05-04 (hardened §4 grep commands with meta-doc glob excludes; both 4.1 and 4.2 now exit 1/no-matches; rephrased `01-architecture.md` §6 to remove literal legacy term — closes F-N-03 fully)*
