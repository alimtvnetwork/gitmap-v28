# Step 4 — Cross-Spec Dependency Audit

**Spec under audit:** `spec/19-main-worker-service/`
**Audit type:** Verify every external `spec/NN-*` reference resolves to real, consistent content.
**Audit version:** 1.1.0 (re-triaged 2026-05-04 — see §10)
**Implementer model assumed:** A "dumb" AI that follows pointers literally, will not patch contradictions, and will not invent missing artifacts.
**Severity scale:** BLOCKER · MAJOR · MINOR

---

## 0. Scope and method

Spec 19 references **5 external spec folders**:

| Folder | Referenced from spec/19 |
|--------|--------------------------|
| `spec/03-error-manage/` | `00-`, `01-§7`, `04-§3`, `07-§9`, `08-`, `plan.md` |
| `spec/04-database-conventions/` | `00-`, `01-§4`, `03-§1`, diagrams README |
| `spec/05-split-db-architecture/` | `00-`, `01-§2`, `03-§9`, `04-§4.2`, `97-AC`, ERD diagram |
| `spec/06-seedable-config-architecture/` | `00-`, `01-§5`, `02-glossary.md`, `03-§Seed` |
| `spec/14-update/` | `09-self-update-pointer.md` (entire file) |

Method:
1. Confirm each folder exists.
2. Confirm each **concept** spec/19 borrows from it actually appears there.
3. Flag contradictions, missing concepts, and dangling deep-references.

All 5 folders **exist on disk** (verified). The remainder of this audit checks **content consistency**, not file presence.

---

## 1. Findings: `spec/05-split-db-architecture/` — **CRITICAL DIVERGENCE**

| ID | Reference in spec/19 | Reality in spec/05 | Severity |
|---|---|---|---|
| F-X-01 | spec/19 calls the worker DB layout **"Root / App / Session"** (3 tiers) — see `00-overview.md`, `01-architecture.md §2`, `02-glossary.md`, `diagrams/erd-worker-split-db.mmd`. | spec/05 `01-fundamentals.md` defines **4 tiers**: **Root / Settings / App / Session**. The `Settings DB` tier is *missing* from every spec/19 mention. | **BLOCKER** — a dumb AI reading spec/19 will build a 3-tier DB and silently violate spec/05. |
| F-X-02 | spec/19 ERD `erd-worker-split-db.mmd` shows entities (`WorkerNode`, `Tenant`, etc.) without naming them as a "projection." | spec/05 owns the schema. spec/19's ERD therefore *competes* with spec/05 instead of pointing to it. | BLOCKER (already noted as F-D-09 in step 3; confirmed here). |
| F-X-03 | spec/19 `04-worker-routing.md §4.2` says "Copy split-DB rows (per `spec/05-`)." | spec/05 does **not** define a "copy rows on tenant migration" operation. There is no migration / row-copy section in `spec/05-split-db-architecture/02-features/` for tenant moves. | MAJOR — the cross-reference is *aspirational*; the target spec lacks the procedure. |
| F-X-04 | spec/19 `97-acceptance-criteria.md` AC: "Worker startup creates Root/App/Session DBs." | spec/05 expects 4 DBs (Root/Settings/App/Session). The AC will pass while the worker is *non-conformant* to spec/05. | BLOCKER. |

**Consequence:** Even with all of spec/19's ambiguities resolved, the dumb AI's worker will not match spec/05.

---

## 2. Findings: `spec/06-seedable-config-architecture/` — **CONCEPT GAPS**

| ID | Reference in spec/19 | Reality in spec/06 | Severity |
|---|---|---|---|
| F-X-05 | spec/19 `03-main-db-schema.md §Seed` says "Seed data for `Role`, `WorkerNodeStatus`, `WorkerNodeKind`, `WorkerSelectionStrategy` ships via Seedable-Config." | spec/06 defines a *generic* seedable-config mechanism (SemVer + GORM merge of `config.seed.json`). It does **not** enumerate any of these four tables. spec/19 must either ship the seed file itself or spec/06 must register them. Neither happens. | MAJOR — dumb AI will emit no seed file, and the DB starts empty. |
| F-X-06 | spec/19 `07-role-based-dashboards.md` requires `EnumPage` + `RolePageAccess` to be *populated* somehow. | spec/06 does not mention `EnumPage` or `RolePageAccess`. spec/19 never says "seeded via spec/06." A dumb AI sees the table but ships no rows. | BLOCKER. |
| F-X-07 | spec/19 `01-architecture.md §6` says: "Both Main and Worker consume Seedable-Config." | spec/06 describes a single-process model. There is no rule for "two tiers consuming the *same* seed file with different subsets." | MAJOR — split-tier seed semantics are undefined. |

---

## 3. Findings: `spec/03-error-manage/` — **REGISTRY GAP**

| ID | Reference in spec/19 | Reality in spec/03 | Severity |
|---|---|---|---|
| F-X-08 | spec/19 `08-error-contract.md` defines an error taxonomy: `worker-unreachable`, `version-mismatch`, `split-DB-write-fail`, `auth-handshake-fail`. | spec/03 `03-error-code-registry/` is a *generic* `ErrorCodes` framework with templates. None of the four spec/19 error families are registered. | BLOCKER — dumb AI cannot `throw new AppError(ErrorCodes.WORKER_UNREACHABLE, ...)` because that constant does not exist. |
| F-X-09 | spec/19 references `spec/03-error-manage/` 11 times for "log, don't swallow." | spec/03's "log don't swallow" rule is in `02-error-architecture/`. The exact section / rule number is never cited. | MINOR — pointer is folder-level, fine for humans, weak for AI. |
| F-X-10 | spec/19 `08-` claims correlation-ID + idempotency-key headers. | spec/03 has no header-naming convention. spec/04 `06-rest-api-format.md` also has none (verified: 0 hits for `X-Correlation` / `Idempotency-Key`). | BLOCKER — header names invented inside spec/19 with no parent contract. |

---

## 4. Findings: `spec/04-database-conventions/` — **MOSTLY OK**

| ID | Reference in spec/19 | Reality in spec/04 | Severity |
|---|---|---|---|
| F-X-11 | spec/19 cites PascalCase, `{TableName}Id` PK, no UUIDs. | spec/04 `01-naming-conventions.md` and `02-schema-design.md` confirm all three. | ✅ — clean. |
| F-X-12 | spec/19 `01-architecture.md §correlation` adds: "no-UUID rule applies to PKs only; opaque request-scoped strings are allowed." | spec/04 does not explicitly grant this exception. | MINOR — spec/19 is *self-granting* an exception to spec/04. Should be lifted into spec/04. |
| F-X-13 | spec/19 borrows REST-API formatting from spec/04 `06-rest-api-format.md`. | spec/04 `06-` exists but does not define endpoints like `/api/Workers/{Id}` or `/api/Companies`. spec/19 invents these without referencing a spec/04 endpoint-naming rule. | MAJOR — REST URL patterns in spec/19 have no parent convention. |

---

## 5. Findings: `spec/14-update/` — **SEVERE CONTRACT GAP**

| ID | Reference in spec/19 | Reality in spec/14 | Severity |
|---|---|---|---|
| F-X-14 | spec/19 `09-self-update-pointer.md §3` says workers "Receive a JSON instruction document. Format spec deferred — will live in a sibling file under `spec/14-update/`." | spec/14 contains **27 numbered sub-docs** (`01-` through `27-`). **Zero** mention "JSON instruction document." (verified: 0 grep hits for `json instruction` or `instruction document`). | BLOCKER — the pointed-to format does not exist anywhere. |
| F-X-15 | spec/19 `97-acceptance-criteria.md` AC: "Endpoint, JSON instruction download, zip-based update flow documented." | spec/14 documents a *rename-first deployment* with `latest.json` + per-asset `.zip` artifacts (`13-release-assets.md`, `19-updater-binary.md`). The *worker self-update fan-out* described in spec/19 is a different, undocumented flow. | BLOCKER — spec/19's AC will be satisfied "by reference" but the referenced flow is not the worker-update flow. |
| F-X-16 | spec/19 `09-` mentions a "redirect-URL DB schema (do not invent here)." | spec/14 has no redirect-URL schema either (verified: 0 hits for `redirect`). | MAJOR — punt with no receiver. |
| F-X-17 | spec/19 `09-` says "If `spec/14-update/` and this pointer ever conflict, `spec/14-update/` wins." | The conflict is **already real today** (F-X-14, F-X-15). The "winner" rule means a dumb AI must implement *spec/14's existing flow* (CLI self-update via release binaries) for **server workers** — which is architecturally wrong (workers update via push, not via pulling release zips). | BLOCKER. |

---

## 6. Findings: dangling cross-references in audit-internal context

| ID | Issue | Severity |
|---|---|---|
| F-X-18 | spec/19 references **5 external folders** but **0 deep anchors** (no `#section-id` links). A dumb AI must read entire folders to find the rule. | MINOR — accept cost; or add anchored links. |
| F-X-19 | None of spec/19's references include a **version pin** (e.g., "as of spec/05 changelog v1.4"). Specs evolve; this audit is a snapshot. | MINOR — recommend adding "consumed-version" notes. |
| F-X-20 | spec/19 `99-consistency-report.md` claims all 5 referenced folders "✅ exist," which is true but is being conflated with "✅ consistent." | MINOR — fix the wording in `99-`. |

---

## 7. Counts (this step)

- **BLOCKER:** 8 (F-X-01, 02, 04, 06, 08, 10, 14, 15, 17 → 9 actually)
- **MAJOR:** 5 (F-X-03, 05, 07, 13, 16)
- **MINOR:** 6 (F-X-09, 11✅ excluded, 12, 18, 19, 20)

(F-X-11 is a ✅, not a finding.)

**Cumulative across steps 1–4:**
- BLOCKER: 17 + 9 = **26**
- MAJOR: 22 + 5 = **27**
- MINOR: 70 + 6 = **76**
- TOTAL: **129 findings**

---

## 8. The 5 cross-spec gaps a dumb AI cannot bridge alone

1. **3-tier vs 4-tier split-DB** (F-X-01) — silent architecture violation.
2. **`EnumPage` / `RolePageAccess` have no seeder** (F-X-06) — RBAC ships empty.
3. **Worker error codes are not registered** in `spec/03` (F-X-08) — code won't compile against `AppError(ErrorCodes.X)`.
4. **JSON instruction document for worker push-updates does not exist** in `spec/14-update/` (F-X-14, F-X-15, F-X-17) — entire push-update path is undefined.
5. **REST endpoint URL conventions** are invented in spec/19 with no parent rule in `spec/04` (F-X-13) — every worker endpoint risks naming drift.

---

## 9. Remaining audit step

| # | File | Purpose |
|---|------|---------|
| 5 | `audit/05-implementation-pivot-score.md` | Final scorecard, "%-pivot-from-spec" estimate, top-10 prioritized fixes, AI-implementer checklist. |

Send `next` to run **Step 5 (Implementation Pivot Score)** — the final step.

---

## 10. Re-Triage After Spec-Hardening Tasks (v1.1.0)

**Re-triaged:** 2026-05-04
**Window:** No-questions tasks #07–#46.
**Method:** For each F-X finding, look up its closing task in `.lovable/question-and-ambiguity/task-counter.md` and confirm by file inspection.

### 10.1 Closure matrix

| ID | Severity | Status | Closed by | Evidence |
|----|----------|--------|-----------|----------|
| F-X-01 | BLOCKER | ✅ CLOSED | #09 | `11-split-db-tier-reconciliation.md` pins Worker = 4 tiers (Root/Settings/App/Session) per spec/05 |
| F-X-02 | BLOCKER | ✅ CLOSED | #15 | "NON-AUTHORITATIVE PROJECTION" banner on `erd-worker-split-db.mmd` cites spec/05 as authority |
| F-X-03 | MAJOR | ✅ CLOSED | #09+#33 | `11-§5` row-copy procedure + cross-spec anchor verification |
| F-X-04 | BLOCKER | ✅ CLOSED | #09 | `97-acceptance-criteria.md` AC-2 updated to require 4-tier worker DB |
| F-X-05 | MAJOR | ✅ CLOSED | #12 | `14-rbac-and-status-seed.md` enumerates Role / WorkerNodeStatus / AuthMechanism rows via Tables-block schema (`07-reference-table-seeding.md`) |
| F-X-06 | BLOCKER | ✅ CLOSED | #12 | `14-rbac-and-status-seed.md` ships 9 EnumPage + 19 RolePageAccess seed rows |
| F-X-07 | MAJOR | ✅ CLOSED | #12 | `07-reference-table-seeding.md` Tables-block schema + `@-ref` resolver covers split-tier seed semantics |
| F-X-08 | BLOCKER | ✅ CLOSED | #11 | MWS prefix (21000-21199) registered in `spec/03/03-error-code-registry/01-registry.md`; 30 codes catalogued in `13-error-codes.md` |
| F-X-09 | MINOR | ✅ CLOSED | #33 | Cross-spec anchor verification swept all `spec/(03\|04\|05\|06\|14)*` references; deep anchors validated |
| F-X-10 | BLOCKER | ✅ CLOSED | #14 | `spec/04/06-rest-api-format.md` defines authoritative `X-Correlation-Id` / `X-Idempotency-Key` / `X-Auth-Action` conventions |
| F-X-11 | (✅ baseline) | — | — | Not a finding (baseline confirmation) |
| F-X-12 | MINOR | ✅ CLOSED | #33 | Opaque request-scoped string exception lifted into `spec/04/06-rest-api-format.md` |
| F-X-13 | MAJOR | ✅ CLOSED | #14+#39 | REST URL patterns documented in `06-core-api-endpoints.md` v1.1.0; cross-ref to `spec/04/06-` |
| F-X-14 | BLOCKER | ✅ CLOSED | #07 | `spec/14-update/28-worker-push-instruction.md` defines JID schema, transport, RenameFirst flow |
| F-X-15 | BLOCKER | ✅ CLOSED | #07 | Same — `28-worker-push-instruction.md` documents endpoint + JID download + zip flow |
| F-X-16 | MAJOR | ✅ CLOSED | #07 | `28-§7` defines `WorkerUpdateInstruction` table with redirect-URL fields |
| F-X-17 | BLOCKER | ✅ CLOSED | #07 | Conflict resolved by authoring `28-` — push-update flow now exists in spec/14 |
| F-X-18 | MINOR | ✅ CLOSED | #33 | Anchored deep-links added across spec/19 cross-references |
| F-X-19 | MINOR | ✅ CLOSED | #33 | "Consumed-version" notes added at cross-ref sites |
| F-X-20 | MINOR | ✅ CLOSED | #33 | `99-consistency-report.md` wording corrected — "exist" vs "consistent" disambiguated |

### 10.2 Tally

| Severity | Original | Closed | Open |
|----------|---------:|-------:|-----:|
| BLOCKER | 9 | 9 | 0 |
| MAJOR | 5 | 5 | 0 |
| MINOR | 6 | 6 | 0 |
| **Total** | **20** | **20** | **0** |

### 10.3 Verdict

All 20 F-X findings closed. The 5 cross-spec gaps from §8 are all resolved:

1. **3-tier vs 4-tier split-DB** → fixed by `11-split-db-tier-reconciliation.md` (#09).
2. **EnumPage / RolePageAccess seeder** → fixed by `14-rbac-and-status-seed.md` (#12).
3. **Worker error codes registered** → fixed by `13-error-codes.md` + spec/03 registry update (#11).
4. **JSON instruction document** → fixed by `spec/14-update/28-worker-push-instruction.md` (#07).
5. **REST endpoint URL conventions** → lifted into `spec/04/06-rest-api-format.md` (#14).

Audit/04 is now in **maintenance mode** — re-open only on regression.

---

*Re-triage appended 2026-05-04 — audit version bumped to 1.1.0.*
