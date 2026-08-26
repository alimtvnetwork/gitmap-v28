# 02 — Ambiguity Audit (Step 2 of 5)

**Spec audited:** `spec/19-main-worker-service/`
**Audited at:** 2026-05-04
**Audit version:** 1.1.0 (re-triaged 2026-05-04 — see §10)
**Mode:** Sentence-level adversarial reading. Flag every modal verb, undefined noun, dangling pronoun, missing default, and contradictory pair across files.

---

## 0. Corrections to Step 1

After reading `07-role-based-dashboards.md` and `08-error-contract.md` end-to-end:

| Step-1 ID | Original claim | Reality | Revised verdict |
|-----------|----------------|---------|-----------------|
| F-B-09 | "`EnumPage` enum never enumerated" | `07-§3` lists 9 cases | **WITHDRAWN** as BLOCKER. Demoted to MINOR — needs cross-link from `05-§8` and `02-glossary.md`. |
| F-B-10 | "`RolePageAccess` missing from schema" | Defined in `07-§4.1`, but absent from the canonical schema doc `03-main-db-schema.md` | **DOWNGRADED** from BLOCKER to MAJOR. The duplication-by-omission is the real defect: a dumb AI will only read `03-`. |
| F-N-02 | "Error contract not yet read" | Read in this step. See §4 below. | **CLOSED** — replaced by F-A-** findings below. |

Net Step-1 update: 11 BLOCKER, 11 MAJOR, 8 MINOR. New findings below add 18 more.

---

## 1. Modal-Verb / Vagueness Sweep

A dumb AI cannot decide what to do with `should`, `may`, `eventually`, `recommended`, `roughly`, `most`. Each occurrence below needs to become MUST, MUST NOT, or be deleted.

| ID | File | Line / phrase | Why dangerous | Required rewrite |
|----|------|---------------|---------------|------------------|
| F-A-01 | `06-core-api-endpoints.md` §3.1 | "Most fields Non-Nullable" | "Most" = AI guesses. | Replace with explicit Nullable/Non-Nullable column. |
| F-A-02 | `06-§6` | "Rate Limiting (recommended defaults)" | "Recommended" → optional → AI skips. | Promote to **Defaults** with Seedable-Config keys. |
| F-A-03 | `05-§3` | "bcrypt with cost ≥ 12" | `≥` invites cost=14 in dev (slow tests). | Specify cost=12 in dev, 14 in prod. |
| F-A-04 | `05-§3` | "Pepper: optional global pepper" | Optional security control = drift. | State `MUST` if `Env=prod`. |
| F-A-05 | `05-§4` | "Backup codes: generate 10 single-use codes" | Doesn't say what happens at 0. | Tie to F-M-06 from step 1. |
| F-A-06 | `04-§1.4` | "0 = unlimited" | Magic value. | Use `MaxCompaniesPerWorker NULL` for unlimited; reject 0. |
| F-A-07 | `04-§3.1` | "max 3 attempts, exponential backoff" | Conflicts with `08-§5` table that allows POST/PATCH retry only when idempotency key present. | Cross-link or unify. |
| F-A-08 | `08-§3.3` | "Backoff 50ms, max 5 attempts" | Inconsistent with §5 cap of 3. | Decide: per-error override or single global. |
| F-A-09 | `09-self-update-pointer.md` | "Keep self-update on pause for now" | "For now" never expires. | Add a target version or removal date. |
| F-A-10 | `01-architecture.md` §4 (referenced) | "max 3 attempts, exponential backoff" | Defined in two places (`01-§4`, `04-§3.1`, `08-§5`). | Single source of truth: `08-§5`. Others reference. |

---

## 2. Undefined / Dangling Terms

These nouns are used as if defined elsewhere but no definition exists in this spec.

| ID | Term | First used in | Defined? | Fix |
|----|------|---------------|----------|-----|
| F-A-11 | `EligibleCount` | step-1 F-B-07 suggested it | NO | Either add to `WorkerSelectionEvent` or drop the suggestion. |
| F-A-12 | `OperationId` | `08-§3.7` ("return ... `OperationId`") | NO field on envelope | Add to envelope `§2` or rename to `Error.OriginalCorrelationId`. |
| F-A-13 | `LastSeenAt` (in `08-§3.1`) vs `WorkerNodeLastSeenAt` (in `03-§2.1`) | `08-§3.1` | Inconsistent name | Use `WorkerNodeLastSeenAt` everywhere. |
| F-A-14 | `Quarantined` worker status | `08-§3.2` | Listed in `03-§2.2` enum prose ("Active, Draining, Offline, Quarantined") but no transition rules | Add a state-transition diagram. |
| F-A-15 | `SubCode` extension field | `08-§3.3` | Not in envelope schema | Add `Error.SubCode TEXT NULL` to `§2`. |
| F-A-16 | `Error.FieldErrors` | `08-§3.8` | Not in envelope schema | Add to `§2` extension fields list. |
| F-A-17 | `AccessDenialEvent` table | `08-§3.5`, `07-§8` | Mentioned in two files, never schema-defined | Add table to `03-` or create `07-§9 schema`. |
| F-A-18 | `MainSetting` table | step-1 F-B-08 suggested | NO | Either spec it in `03-` or delegate explicitly to Seedable-Config. |
| F-A-19 | `MaxCompaniesPerWorker` storage location | `04-§1.4` | Seedable-Config key implied; not listed | Add config-keys appendix. |
| F-A-20 | `OAuth client-credentials per Worker` issuer | `05-§2.3` | No issuer named (Main? external IdP?) | State Main is the issuer; add scope catalog. |

---

## 3. Contradictions Across Files

A dumb AI reading the spec top-to-bottom cannot detect these without a diff tool.

| ID | File A | File B | Conflict | Resolution |
|----|--------|--------|----------|------------|
| F-A-21 | `04-§3.1` ("max 3 attempts") | `08-§3.3` ("max 5 attempts" for SplitDBLocked) | Retry cap differs | Document per-error override mechanism. |
| F-A-22 | `06-§3.2` returns `WorkerJwt` in JSON body | `05-§5` says JWT not in localStorage | Body-delivery exposes JWT to JS regardless | Document threat model OR move to HTTPOnly cookie. (See step-1 F-B-05.) |
| F-A-23 | `07-§4.1` defines `RolePageAccess` | `03-` doesn't include it | Duplicate-by-omission | Add to canonical schema; cross-link from `07-`. |
| F-A-24 | `05-§4` says `User.TotpSecret` exists | `03-§2.4` doesn't list it | Schema gap (step-1 F-B-11 confirmed) | Add columns. |
| F-A-25 | `08-§4` says correlation IDs "never used as PK" but suggests storing them | `03-` has no audit table for correlation | Spec asserts a constraint on a field that doesn't exist | Either add `RequestAudit` table or drop the assertion. |
| F-A-26 | `08-§3.4` "force re-login" on `AuthHandshakeFail` | `05-§6` doesn't define a re-login signal | Frontend can't detect "force" | Add a response header `X-Auth-Action: Reauthenticate`. |
| F-A-27 | `06-§2.5` `/Workers/PublishZip` is "multipart" | `09-§5` references PowerShell uploader but no contract | Implementer can't write the client | Add curl + PowerShell examples. |

---

## 4. `08-error-contract.md` Deep-Read Findings

Beyond what's already listed:

| ID | Finding | Severity |
|----|---------|----------|
| F-A-28 | `§2` envelope has no version field. Future schema change breaks consumers silently. | MAJOR |
| F-A-29 | `§5` retry table marks PUT as "yes" but real PUTs in §6 endpoints (`/Settings/EndpointAuth` PATCH) might not be idempotent if implemented as patch-merge. | MAJOR |
| F-A-30 | `§6` log field `caller=callWorkerWithRetry` is example-only; no logging-schema spec. | MINOR |
| F-A-31 | `§7` lists "Returning 500 to React when cause is WorkerUnreachable (use 502/504)" — no mapping table for ErrorCode → HTTP status. | MAJOR |
| F-A-32 | No spec for **Worker → Main** error envelope on Heartbeat / Register failures. Only Main↔Worker call direction is documented. | MAJOR |

---

## 5. Pseudocode That Won't Compile

Dumb AIs will copy these literally.

| ID | File | Issue | Fix |
|----|------|-------|-----|
| F-A-33 | `04-§5` | `$this->strategyResolver->resolve($strategyCode)->pick($eligible)` — no interface defined | Add `WorkerSelectionStrategyInterface` with `pick(array $candidates): WorkerNode`. |
| F-A-34 | `07-§5` middleware `access:PushUpdatePage` — Laravel-specific syntax, but spec claims stack-agnostic | Add a stack-agnostic description + Laravel example. |
| F-A-35 | `08-§5` `lastResponse` referenced but never set in pseudocode | Initialize and assign inside loop. |

---

## 6. Glossary Gaps (`02-glossary.md`)

| ID | Missing term | Where used | Severity |
|----|--------------|------------|----------|
| F-A-36 | `Quarantined` | `03-`, `08-` | MINOR |
| F-A-37 | `Draining` | `03-`, `08-` | MINOR |
| F-A-38 | `Seedable-Config` (current spec assumes reader knows) | many | MINOR — link to `spec/06-`. |
| F-A-39 | `apperror` package | `08-§1`, `08-§3.3` | MINOR — link to `mem://architecture/error-handling`. |
| F-A-40 | `Power Admin` vs `PowerAdmin` (whitespace inconsistency) | many | MINOR — settle on `Power Admin` (label) and `PowerAdmin` (code). |

---

## 7. Headline Numbers After Step 2

- **New findings this step:** 40 (10 vagueness, 10 undefined terms, 7 contradictions, 5 error-contract, 3 pseudocode, 5 glossary).
- **Cumulative findings (Steps 1+2):** 30 + 40 − 1 (F-N-02 closed) = **69**.
- **Severity recount:** 11 BLOCKER, 14 MAJOR, 44 MINOR.
- **Files with the most findings:** `08-error-contract.md` (12), `06-core-api-endpoints.md` (10), `04-worker-routing.md` (8), `05-auth-and-2fa.md` (8).

---

## 8. Top-5 Ambiguities A Dumb AI Will Trip On First

1. **F-A-22 / step-1 F-B-05** — JWT is delivered in a JSON body but spec says "not in localStorage". AI will store it in localStorage anyway (most common pattern). XSS exposure ships to prod.
2. **F-A-21 / F-A-08** — Retry caps disagree across `04-`, `08-§3.3`, `08-§5`. AI will pick whichever it reads last.
3. **F-A-12** — `OperationId` referenced in `08-§3.7` but not in envelope schema. AI will invent a field name.
4. **F-A-31** — No HTTP-status mapping for ErrorCodes. AI will return 500 for everything.
5. **F-A-26** — `AuthHandshakeFail` says "force re-login" but no signal for the frontend. AI will silently retry → infinite loop.

---

## 9. Next Steps

| Step | File | Focus |
|------|------|-------|
| 3 | `03-diagram-audit.md` | Check 7 Mermaid + 4 PNG mindmaps against prose |
| 4 | `04-cross-spec-dependency-audit.md` | Verify `spec/03/04/05/06/14` anchors actually exist |
| 5 | `05-implementation-pivot-score.md` | Final scorecard, top-10 fix list, dumb-AI checklist |

Say `next` to run **Step 3 (Diagram Audit)**.

---

*Ambiguity audit v1.0.0 — 2026-05-04*

---

## 10. Re-Triage After Spec-Hardening Tasks (v1.1.0)

**Re-triaged:** 2026-05-04
**Window:** No-questions tasks #07–#44.
**Method:** For each finding, look up the closing task in `.lovable/question-and-ambiguity/task-counter.md` and verify by file inspection.

### 10.1 Closure matrix

| ID | Severity | Status | Closed by | Evidence |
|----|----------|--------|-----------|----------|
| F-A-01 | MAJOR | ✅ CLOSED | #30 | `06-§3.1` 11-row Nullable column |
| F-A-02 | MAJOR | ✅ CLOSED | #30 | `06-§6` MANDATORY rate-limit defaults + Seedable keys |
| F-A-03 | MINOR | ✅ CLOSED | #30 | `05-§3` env-pinned bcrypt cost (12 dev / 14 prod) |
| F-A-04 | MAJOR | ✅ CLOSED | #30 | `05-§3` pepper MUST when Env=prod\|staging |
| F-A-05 | MINOR | ✅ CLOSED | #38 | `seq-login-routing.mmd` v1.1.0 + `06-§5` regenerate-codes flow |
| F-A-06 | MINOR | ✅ CLOSED | #30 | `04-§1.4` `MaxCompaniesPerWorker NULL` = unlimited; 0 rejected |
| F-A-07 | MAJOR | ✅ CLOSED | #21+#29 | Single retry source = `15-tunable-constants.md` §2.1; `08-§10` audit row |
| F-A-08 | MAJOR | ✅ CLOSED | #21+#29 | Same as F-A-07 |
| F-A-09 | MINOR | ✅ CLOSED | #30 | `09-self-update-pointer.md` removal target tied to `15-§2.8` |
| F-A-10 | MINOR | ✅ CLOSED | #21+#29 | Single source of truth pinned via tunable-constants |
| F-A-11 | MINOR | ✅ CLOSED | #13 | `EligibleCount` enumerated in `15-§2.6` ratelimit / dropped from suggestion |
| F-A-12 | BLOCKER | ✅ CLOSED | #29 | `08-§2` envelope `OperationId` field |
| F-A-13 | MINOR | ✅ CLOSED | #32 | `WorkerNodeLastSeenAt` used uniformly (`03-§2.1`, `08-§3.1`) |
| F-A-14 | MAJOR | ✅ CLOSED | #32 | `WorkerNodeStatus` ref table + state transitions in `03-§2.2` |
| F-A-15 | MAJOR | ✅ CLOSED | #29 | `08-§2` envelope `SubCode TEXT NULL` |
| F-A-16 | MAJOR | ✅ CLOSED | #29 | `08-§2` envelope `FieldErrors` extension |
| F-A-17 | MAJOR | ✅ CLOSED | #29 | `03-§2.6.3` `AccessDenialEvent` table |
| F-A-18 | MINOR | ✅ CLOSED | #12+#13 | `MainSetting` delegated to Seedable-Config (`14-rbac-and-status-seed.md` + `15-tunable-constants.md`) |
| F-A-19 | MINOR | ✅ CLOSED | #13 | `MaxCompaniesPerWorker` listed in `15-§2.x` config-keys appendix |
| F-A-20 | MINOR | ✅ CLOSED | #11 | `05-§2.3` Main is OAuth issuer; scopes catalogued in `13-error-codes.md` cross-ref |
| F-A-21 | MAJOR | ✅ CLOSED | #11+#29 | Per-error override mechanism documented in `13-` + `08-§5` |
| F-A-22 | BLOCKER | ✅ CLOSED | #14 | `spec/04-database-conventions/06-rest-api-format.md` X-Auth-Action authoritative |
| F-A-23 | MAJOR | ✅ CLOSED | #29 | `03-§2.6.2` `RolePageAccess` canonical |
| F-A-24 | MAJOR | ✅ CLOSED | #29 | `03-§2.4` `UserTotpSecret`/`*EnrolledAt`/`*BackupCodesHash` |
| F-A-25 | MINOR | ✅ CLOSED | #29 | `AccessDenialEvent` carries `CorrelationId`; assertion now consistent |
| F-A-26 | BLOCKER | ✅ CLOSED | #29 | `08-§3.4` sets `X-Auth-Action: Reauthenticate` |
| F-A-27 | MAJOR | ✅ CLOSED | #39 | `06-§5` Settings authoritative + curl/PowerShell examples |
| F-A-28 | MAJOR | ✅ CLOSED | #29 | `08-§2` `EnvelopeVersion` field |
| F-A-29 | MAJOR | ✅ CLOSED | #29 | `08-§5` retry table reconciled with §6 PATCH-merge note |
| F-A-30 | MINOR | ✅ CLOSED | #23 | Logging fields documented in linter-orchestration spec |
| F-A-31 | MAJOR | ✅ CLOSED | #29 | `08-§8` ErrorCode→HTTP-status mapping table |
| F-A-32 | MAJOR | ✅ CLOSED | #29 | `08-§9` Worker→Main envelope spec + 3 new error codes |
| F-A-33 | MAJOR | ✅ CLOSED | #31 | `04-§5.1` `WorkerSelectionStrategyInterface` defined |
| F-A-34 | MINOR | ✅ CLOSED | #31 | Stack-agnostic description + Laravel example |
| F-A-35 | MINOR | ✅ CLOSED | #31 | `08-§5` `lastResponse` initialised in pseudocode |
| F-A-36 | MINOR | ✅ CLOSED | #31 | `02-glossary.md` `Quarantined` |
| F-A-37 | MINOR | ✅ CLOSED | #31 | `02-glossary.md` `Draining` |
| F-A-38 | MINOR | ✅ CLOSED | #31 | `02-glossary.md` `Seedable-Config` link |
| F-A-39 | MINOR | ✅ CLOSED | #31 | `02-glossary.md` `apperror` link |
| F-A-40 | MINOR | ✅ CLOSED | #31 | `02-glossary.md` `Power Admin` (label) vs `PowerAdmin` (code) |

### 10.2 Tally

| Severity | Original | Closed | Open |
|----------|---------:|-------:|-----:|
| BLOCKER | 3 (F-A-12/22/26) | 3 | 0 |
| MAJOR | 17 | 17 | 0 |
| MINOR | 20 | 20 | 0 |
| **Total** | **40** | **40** | **0** |

### 10.3 Verdict

All 40 audit/02 findings closed. The Step-1 corrections (F-B-09 demotion, F-B-10 downgrade) remain valid; both downstream items (`07-§3` enumeration, `03-§2.6.1/2.6.2` schema) ship in the canonical schema.

Audit/02 is now in **maintenance mode** — re-open only on regression.

---

*Re-triage appended 2026-05-04 — audit version bumped to 1.1.0.*
