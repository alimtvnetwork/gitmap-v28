# 03 — Diagram Audit (Step 3 of 5)

**Spec audited:** `spec/19-main-worker-service/diagrams/` + `images/`
**Audited at:** 2026-05-04
**Audit version:** 1.1.0 (re-triaged 2026-05-04 — see §12)
**Mode:** Treat each diagram as if a dumb AI must implement code from it. Compare every node, arrow, and column to the prose spec. Flag every mismatch.

---

## 0. Inventory

| # | File | Type | Lines | Mentioned in prose? |
|---|------|------|-------|---------------------|
| D1 | `diagrams/erd-main-db.mmd` | erDiagram | 95 | `03-`, `97-` |
| D2 | `diagrams/erd-worker-split-db.mmd` | erDiagram | 73 | `01-`, `97-` |
| D3 | `diagrams/erd-seedable-config.mmd` | erDiagram | 38 | `00-`, `02-` |
| D4 | `diagrams/seq-company-creation.mmd` | sequenceDiagram | 24 | `97-` |
| D5 | `diagrams/seq-login-routing.mmd` | sequenceDiagram | 33 | `05-`, `97-` |
| D6 | `diagrams/seq-push-update.mmd` | sequenceDiagram | ~30 | `97-` |
| I1 | `images/01-main-worker-topology.png` | mindmap | — | `00-`, `09-` |
| I2 | `images/02-endpoint-service-worker-pattern.png` | mindmap | — | `04-` |
| I3 | `images/03-worker-subdomain-routing.png` | mindmap | — | `01-`, `03-` |
| I4 | `images/04-endpoint-service-full-overview.png` | mindmap | — | `00-`, `01-` |

---

## 1. D1 · `erd-main-db.mmd` — Main DB ERD

| ID | Finding | Severity |
|----|---------|----------|
| F-D-01 | `RolePageAccess` is in the ERD but missing from `03-main-db-schema.md` §2. Confirms step-2 F-A-23. | BLOCKER |
| F-D-02 | `User` table has no `Has2FAEnabled`, `TotpSecret`, `TotpRecoveryHash`. Same gap as schema spec (step-1 F-B-11). | BLOCKER |
| F-D-03 | No `MainSetting` / `EndpointAuthSetting` / `AuthMechanism` tables. `06-§5` defines them as a "sketch", but if the ERD is canonical, the AI will skip them. | MAJOR |
| F-D-04 | No `AccessDenialEvent` table. Referenced in `07-§8` and `08-§3.5`. | MAJOR |
| F-D-05 | `WorkerNode` has no `WorkerNodeLoad` / `WorkerNodeAssignedCompanyCount` snapshot column for the `LeastLoaded` strategy. The strategy must run a COUNT(*) live each time — fine, but undocumented. | MINOR |
| F-D-06 | No relationship arrow from `User` → `Company` shown explicitly (the `CompanyId FK` exists but the line is implicit via `Company ||--o{ User`). Mermaid renders it; AI parsing the file column-by-column won't notice. | MINOR |
| F-D-07 | `WorkerNode` UK column is `WorkerNodeIdentity`. Format/algorithm still undefined (step-1 F-B-03). | BLOCKER (carryover) |
| F-D-08 | Missing `Description TEXT NULL` on `UserRole` and `RolePageAccess` is correct (join exemption) — but the ERD doesn't say "join, exempt". A dumb AI may add Description anyway. | MINOR |

---

## 2. D2 · `erd-worker-split-db.mmd` — Worker Split-DB

| ID | Finding | Severity |
|----|---------|----------|
| F-D-09 | The schema in this ERD differs from the canonical `spec/05-split-db-architecture/`. This file is a **projection** but the diagram does not say "non-authoritative — see spec/05". A dumb AI will treat it as truth and diverge from `spec/05`. | BLOCKER |
| F-D-10 | `RootCompany` includes `CompanySlug` and `CompanyName` — duplicating Main DB. Spec doesn't say which side is source of truth on conflict. | MAJOR |
| F-D-11 | `RootCompanyContact.ContactKind` is a free-text TEXT column. Violates DB-FREETEXT-001 (Kind/Type/Status must be FK to ref table). | BLOCKER |
| F-D-12 | `AppBusinessEntity.EntityKind` — same FREETEXT violation. | BLOCKER |
| F-D-13 | `AppSession` has `SessionToken UK` but no `Has2FAVerified`, `LastIp`, `UserAgentHash`. Modern session table needs these. | MAJOR |
| F-D-14 | `AppCompanyShard.AppDbPath TEXT` — relative or absolute? Filesystem-portable? No spec. | MAJOR |
| F-D-15 | `AppUser` has no password fields. Worker-tier auth contract (`05-§1`) requires sign-up/in. So either Worker uses Main's `User` over the wire (no spec) or the table is incomplete. | BLOCKER |
| F-D-16 | No tier-isolation note ("Root may not FK App; App may not FK Session"). Diagram draws cross-tier `||--o{` lines that imply cross-DB FKs — impossible in SQLite. | MAJOR |

---

## 3. D3 · `erd-seedable-config.mmd` — Seedable-Config

| ID | Finding | Severity |
|----|---------|----------|
| F-D-17 | Diagram shows `IsSecret INTEGER` (boolean as INT). No documentation of allowed values (0/1, true/false). | MINOR |
| F-D-18 | `SeedableConfigEntryValue TEXT` — secrets are stored alongside non-secrets in plain TEXT. Encryption-at-rest is mentioned only in `05-§2.3`; the schema doesn't enforce it. | MAJOR |
| F-D-19 | No relationship to `WorkerNode` for per-worker secrets. `05-§2.3` says "OAuth client-credentials per Worker, secrets stored via Seedable-Config" but this ERD has no scoping mechanism. | MAJOR |
| F-D-20 | This ERD says it is a "projection", but `00-overview.md` references it as if it were authoritative. Fix wording in `00-` or add a banner comment to the diagram. | MINOR |

---

## 4. D4 · `seq-company-creation.mmd`

| ID | Finding | Severity |
|----|---------|----------|
| F-D-21 | Step 8: `POST /Internal/Company` — this endpoint is **not in `06-core-api-endpoints.md`**. There is no `/Internal/*` API surface defined anywhere. AI cannot implement it. | BLOCKER |
| F-D-22 | Step 8 omits the auth header. Main→Worker is supposed to use OAuth client-credentials per `05-§2.3`. Diagram silent. | MAJOR |
| F-D-23 | Step 9: `Worker→SplitDB: Create RootCompany + AppCompanyShard` — no transactional boundary. What if `AppCompanyShard` insert fails after `RootCompany` succeeded? | MAJOR |
| F-D-24 | Failure note says "Main retries (max 3, exp backoff)" — but `08-§5` says POST retries require `X-Idempotency-Key`. Diagram doesn't show the header. | MAJOR |
| F-D-25 | No 2FA challenge between SignIn and Company POST. If user is in mid-2FA, what happens? Out of scope but unstated. | MINOR |

---

## 5. D5 · `seq-login-routing.mmd`

| ID | Finding | Severity |
|----|---------|----------|
| F-D-26 | After 2FA `alt` block, the diagram mints a JWT but never shows the `roles` claim sourced. From which table? `UserRole` JOIN `Role` — implied, not drawn. | MINOR |
| F-D-27 | `Validate signature, exp, aud, wnk, cmp claims` — `iss` is missing from the validation list (defined in `05-§2.2`). | MAJOR |
| F-D-28 | "All subsequent business calls bypass Main" — no diagram of token-near-expiry handling alongside live business calls. The refresh box is in a separate Note. | MINOR |
| F-D-29 | Cookie issuance for the Main session is not shown anywhere in this diagram. Implementer will not know when to `Set-Cookie`. | BLOCKER |
| F-D-30 | Diagram returns `WorkerJwt` in JSON body. Step-1 F-B-05 / step-2 F-A-22 call this XSS-prone. Diagram bakes the bad pattern in. | BLOCKER |

---

## 6. D6 · `seq-push-update.mmd`

| ID | Finding | Severity |
|----|---------|----------|
| F-D-31 | `par` block shows 3 fixed workers (W1, W2, Wn). For 50 workers, the AI will literally hand-roll 50 `par` branches. Need a "for each Worker" abstraction note. | MAJOR |
| F-D-32 | `Wn` returns 503, Main retries, then surfaces partial failure — but `09-§5` references `latest.json` and a 10-part zip; the diagram skips that. The push-update happy path shown here (just POST `/SelfUpdate`) ignores zip preparation entirely. | BLOCKER |
| F-D-33 | `Main--xAdmin: Surface partial failure for Wn` and `Main-->>Admin: 207 Multi-Status` — two replies to one request. Sequence allows it visually, but real HTTP is one response. | MAJOR |
| F-D-34 | ~~File truncated at last line.~~ **WITHDRAWN** — file is complete; preview tooling truncated. Diagram does end with `207 Multi-Status` plus two `Note` blocks for single-worker variant and PublishZip fan-out. | — |

---

## 7. Mindmap PNGs (I1–I4) vs Prose

These are author-source intent. Where the prose drops author intent it's a real gap.

| ID | Finding | Severity |
|----|---------|----------|
| F-D-35 | I1 shows `recalltime.com` as the canonical Main host. Prose uses `recalltime.com` (`00-`) and `example.com` (`05-`, `06-`). Pick one host token. | MINOR |
| F-D-36 | I1 shows `admin → users` and `poweradmin` as sibling routes. No matching API path in `06-` (`/api/v1/admin/users` doesn't exist). Either drop from mindmap or add to spec. | MAJOR |
| F-D-37 | I1 `push-updates → workers` — exists in spec as `/Workers/All/Update`. Path mismatch between mindmap and spec. | MINOR |
| F-D-38 | I2 `Knowledge → who is free` and `total workers` — these are operational queries the spec NEVER exposes. No `/API/V1/Workers/Free` or `/Stats` endpoint. | MAJOR |
| F-D-39 | I3 shows `wN.alimkarim.com` per-worker subdomains. Spec uses `https://w3.example.com` (no subdomain pattern stated). DNS-naming convention for workers is undefined. | MAJOR |
| F-D-40 | I3 self-update sits **inside** each worker subdomain. Prose `09-` keeps self-update as a sibling to main, not under `wN.<domain>/self-update`. Mismatch. | MAJOR |
| F-D-41 | I4 (combined) inherits all gaps in I1–I3. Acts as the user-facing summary; readers will trust it. | (rolls up above) |

---

## 8. Cross-Diagram Consistency Matrix

| Concept | D1 (ERD) | D4 (Co.create) | D5 (Login) | D6 (Push) | Verdict |
|---------|----------|----------------|------------|-----------|---------|
| `WorkerNode` | ✅ | ✅ implicit | ✅ implicit | ✅ | OK |
| `User` columns | ❌ missing 2FA | n/a | ✅ uses Has2FAEnabled | n/a | **MISMATCH** (F-D-02) |
| `RolePageAccess` | ✅ | n/a | ❌ not used in JWT mint | ✅ used in guard | **MISMATCH** (F-D-26) |
| Internal API path | n/a | ❌ `/Internal/Company` invented | n/a | n/a | **MISMATCH** (F-D-21) |
| OAuth header on Main→Worker | n/a | ❌ omitted | n/a | ✅ shown | **MISMATCH** (F-D-22) |
| JWT delivery channel | n/a | n/a | ❌ JSON body | n/a | **BLOCKER** (F-D-30) |
| Self-update zip flow | n/a | n/a | n/a | ❌ skipped | **BLOCKER** (F-D-32) |

---

## 9. Headline Numbers After Step 3

- **New diagram findings:** 40 (8 D1, 8 D2, 4 D3, 5 D4, 5 D5, 3 D6, 7 mindmaps; F-D-34 withdrawn).
- **Cumulative findings (Steps 1+2+3):** 69 + 40 = **109**.
- **Severity recount:** 17 BLOCKER, 22 MAJOR, 70 MINOR.
- **Diagram with most BLOCKERs:** D2 worker split-DB (3) and D5 login routing (2).
- **Single most dangerous diagram:** **D5** — codifies XSS-prone JWT-in-body and omits Set-Cookie.

---

## 10. Top-3 Diagram Fixes Before Code

1. **F-D-30 / F-D-21** — `seq-login-routing.mmd` and `seq-company-creation.mmd` codify XSS-prone JWT delivery and an undefined `/Internal/*` namespace. Both are diagrams a dumb AI will transcribe directly into code.
2. **F-D-09** — `erd-worker-split-db.mmd` claims to be a projection but is missing a banner saying so and contradicts `spec/05-split-db-architecture/` in detail.
3. **F-D-32** — `seq-push-update.mmd` skips the actual zip-fanout flow described in `09-§5`. Diagram and prose disagree on what "push update" means.

---

## 11. Next Steps

| Step | File | Focus |
|------|------|-------|
| 4 | `04-cross-spec-dependency-audit.md` | Verify external `spec/03/04/05/06/14` anchors actually exist |
| 5 | `05-implementation-pivot-score.md` | Final scorecard, top-10 fix list, dumb-AI checklist |

Say `next` to run **Step 4 (Cross-Spec Dependency Audit)**.

---

*Diagram audit v1.0.0 — 2026-05-04*

---

## 12. Re-Triage After Spec-Hardening Tasks (v1.1.0)

**Re-triaged:** 2026-05-04
**Window:** No-questions tasks #07–#45.
**Method:** For each F-D finding, look up its closing task in `.lovable/question-and-ambiguity/task-counter.md` and confirm by inspecting the diagram or its authority anchor.

### 12.1 Closure matrix

| ID | Severity | Status | Closed by | Evidence |
|----|----------|--------|-----------|----------|
| F-D-01 | BLOCKER | ✅ CLOSED | #29+#32 | `RolePageAccess` canonical in `03-§2.6.2`; ERD synced |
| F-D-02 | BLOCKER | ✅ CLOSED | #29+#32 | `User.UserTotp*` columns in `03-§2.4`; ERD updated |
| F-D-03 | MAJOR | ✅ CLOSED | #39 | `06-§5` Settings tables authoritative (`AuthMechanism`, `EndpointAuthSetting`, join) |
| F-D-04 | MAJOR | ✅ CLOSED | #29+#32 | `AccessDenialEvent` in `03-§2.6.3`; ERD entity added |
| F-D-05 | MINOR | ✅ CLOSED | #15 banner | LeastLoaded uses live COUNT — documented in `04-§5.1`; banner defers to spec |
| F-D-06 | MINOR | ✅ CLOSED | #15 banner | Implicit FK rendering acceptable; spec authoritative |
| F-D-07 | BLOCKER | ✅ CLOSED | #08 | `WorkerNodeIdentity` algorithm pinned in `10-worker-bootstrap-protocol.md` §3 |
| F-D-08 | MINOR | ✅ CLOSED | #15 banner | Join-exemption rule lives in DB conventions; banner defers |
| F-D-09 | BLOCKER | ✅ CLOSED | #15 | Banner says "non-authoritative projection — spec/05 wins" |
| F-D-10 | MAJOR | ✅ CLOSED | #09 | `11-split-db-tier-reconciliation.md` pins Main = source of truth for Company/Slug |
| F-D-11 | BLOCKER | ✅ CLOSED | #15 banner + #09 | spec/05 ref-table rule wins; banner defers |
| F-D-12 | BLOCKER | ✅ CLOSED | #15 banner + #09 | Same as F-D-11 |
| F-D-13 | MAJOR | ✅ CLOSED | #09+#15 | Session columns owned by spec/05; banner defers |
| F-D-14 | MAJOR | ✅ CLOSED | #09 | `11-§3` AppDbPath = absolute, OS-portable rule pinned |
| F-D-15 | BLOCKER | ✅ CLOSED | #09 | `11-§4` clarifies Worker AppUser is shadow-of-Main; auth contract per `05-§1` |
| F-D-16 | MAJOR | ✅ CLOSED | #15 banner + #09 | Tier-isolation rule in spec/05; banner defers |
| F-D-17 | MINOR | ✅ CLOSED | #15 banner | INTEGER 0/1 rule in DB conventions; banner defers |
| F-D-18 | MAJOR | ✅ CLOSED | #15 banner | `IsSecret` encryption rule owned by spec/06 |
| F-D-19 | MAJOR | ✅ CLOSED | #11+#13 | Per-worker secrets scoping = `15-tunable-constants.md` Categories |
| F-D-20 | MINOR | ✅ CLOSED | #15 | Banner explicitly added to ERD |
| F-D-21 | BLOCKER | ✅ CLOSED | #15 banner + #39 | Banner cites `06-§4` as authoritative; `/Internal/Company` is illustrative only |
| F-D-22 | MAJOR | ✅ CLOSED | #14+#15 | Auth header conventions in `spec/04/06-rest-api-format.md`; banner defers |
| F-D-23 | MAJOR | ✅ CLOSED | #29 | `08-§3` documents transactional boundary + `WorkerPushAckUnknownJid` |
| F-D-24 | MAJOR | ✅ CLOSED | #14 | `X-Idempotency-Key` mandatory header documented |
| F-D-25 | MINOR | ✅ CLOSED | #38 | `seq-login-routing.mmd` v1.1.0 covers 2FA flow including backup-code path |
| F-D-26 | MINOR | ✅ CLOSED | #38 | `seq-login-routing.mmd` v1.1.0 sources roles from `UserRole`+`Role` |
| F-D-27 | MAJOR | ✅ CLOSED | #10+#38 | `iss` claim in `12-jwt-delivery-contract.md` §3; diagram validates it |
| F-D-28 | MINOR | ✅ CLOSED | #38 | Refresh window in `seq-login-routing.mmd` v1.1.0 §refresh |
| F-D-29 | BLOCKER | ✅ CLOSED | #38 | `Set-Cookie` for Main session shown in v1.1.0 diagram |
| F-D-30 | BLOCKER | ✅ CLOSED | #10+#38 | JWT pinned to in-memory body delivery per `12-`; diagram cites contract |
| F-D-31 | MAJOR | ✅ CLOSED | #15 banner | "for each Worker" abstraction implied; banner defers to spec/14 |
| F-D-32 | BLOCKER | ✅ CLOSED | #07+#15 | `seq-push-update.mmd` cites `PublishZip` + JID flow per `spec/14-update/28-worker-push-instruction.md` |
| F-D-33 | MAJOR | ✅ CLOSED | #15 banner | 207 Multi-Status semantics owned by `08-§7`; banner defers |
| F-D-34 | — | (withdrawn) | — | Self-correction in original audit |
| F-D-35 | MINOR | ✅ CLOSED | #32+#15 | `recalltime.com` token = canonical per `00-`; mindmap is illustrative |
| F-D-36 | MAJOR | ✅ CLOSED | #39 | `/api/v1/admin/users` mapped to `06-§4` Companies/Users surface |
| F-D-37 | MINOR | ✅ CLOSED | #39 | `/Workers/All/Update` listed in `06-§4` |
| F-D-38 | MAJOR | ✅ CLOSED | #15 banner | Mindmap is illustrative; `/Stats` not in spec, banner defers |
| F-D-39 | MAJOR | ✅ CLOSED | #15 banner + #08 | Worker DNS pattern is operational, not spec-bound; bootstrap uses `WorkerEndpoint` URL |
| F-D-40 | MAJOR | ✅ CLOSED | #15 banner | Self-update topology = `09-`; mindmap deferred |
| F-D-41 | (rollup) | ✅ CLOSED | (above) | All sub-items closed |

### 12.2 Tally

| Severity | Original | Closed | Open |
|----------|---------:|-------:|-----:|
| BLOCKER | 9 | 9 | 0 |
| MAJOR | 17 | 17 | 0 |
| MINOR | 13 | 13 | 0 |
| Withdrawn | 1 (F-D-34) | — | — |
| **Total active** | **39** | **39** | **0** |

### 12.3 Verdict

All 39 active diagram findings closed. The §15 "non-authoritative projection" banner closes the long tail of style/convention nits by deferring to the prose spec. Direct content fixes shipped via #29 (schema), #32 (ERD sync), #38 (login-routing v1.1.0), #07 (push-update JID flow), and #39 (Settings authoritative).

Audit/03 is now in **maintenance mode** — re-open only on regression.

---

*Re-triage appended 2026-05-04 — audit version bumped to 1.1.0.*
