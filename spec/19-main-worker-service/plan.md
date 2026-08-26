# Plan — Main / Worker Service Architecture

**Spec ID:** 19-main-worker-service
**Status:** Spec authoring in progress
**Owner:** Md. Alim Ul Karim, Riseup Asia LLC
**Default stack:** Laravel (PHP). Stack-agnostic by design — see §Stack Flexibility in `00-overview.md`.

---

## Decisions locked in

| # | Question | Decision |
|---|----------|----------|
| Q1 | Tenant root model | **Company-as-root** (multi-tenant). User-as-root is a degenerate case (1:1). |
| Q2 | Spec folder slot | `spec/19-main-worker-service/` (slots 19–20 free, no conflict). |
| Q3 | Diagrams location | `spec/19-main-worker-service/diagrams/` (self-contained per spec-authoring guide). |
| Q4 | Error-manage integration | **(a)** Inline only main↔worker-specific patterns; reference `spec/03-error-manage/` for generics. |
| Q5 | 5 follow-up tasks | Each becomes a numbered spec file (`03-`…`07-`). Implementation waits for `next`. |

## Open questions (carried forward, NOT decided)

- **OQ-1 (from §Main Server Concept 3c):** Per-endpoint authentication-type configuration — flagged as `❓ Open Question` in `05-auth-and-2fa.md`.
- **OQ-2:** ✅ **Resolved Phase 12.3.** Worker-selection default = `LeastLoaded`. Pinned in `15-tunable-constants.md` §2.5; rationale + override guidance in `04-worker-routing.md` §1.5.

---

## Phases

Each phase is one user `next` cycle.

### ✅ Phase 1 — Foundation (this turn)

- `plan.md` (this file)
- `00-overview.md` — purpose, scope, stack flexibility, terminology pointers
- `01-architecture.md` — main/worker topology, request flow, boundaries
- `02-glossary.md` — Main Server, Worker Node, Tenant Root, Power Admin, etc.

### Phase 2 — Five deliverable specs (on `next`)

- `03-main-db-schema.md` — main-server SQLite schema (workers, tenant→worker map, version tracking)
- `04-worker-routing.md` — selection strategies, caching, failover
- `05-auth-and-2fa.md` — auth flows, 2FA, session, JWT/cookie, main↔worker auth handshake
- `06-core-api-endpoints.md` — `API/V1/Company`, status, version, push-update endpoints
- `07-role-based-dashboards.md` — Power Admin / Admin User / extensible roles, `User has access to {EnumPage}` pattern

### Phase 3 — Diagrams (on `next`)

- `diagrams/erd-main-db.mmd` — main-server SQLite ERD
- `diagrams/erd-worker-split-db.mmd` — worker-side split-DB ERD (Root / App / Session)
- `diagrams/erd-seedable-config.mmd` — Seedable-Config layout for both tiers
- `diagrams/seq-company-creation.mmd` — POST `/API/V1/Company` end-to-end
- `diagrams/seq-login-routing.mmd` — first-request resolution + caching
- `diagrams/seq-push-update.mmd` — main → worker push update
- `diagrams/readme.md` — index of diagrams

### Phase 4 — Error-manage integration + consistency (on `next`)

- `08-error-contract.md` — main↔worker correlation IDs, retry semantics, failure taxonomy (worker-unreachable / version-mismatch / split-DB-write-fail / auth-handshake-fail). References `spec/03-error-manage/`, does not duplicate.
- `09-self-update-pointer.md` — pointer-only doc per verbatim §Self-Update Mechanism (NOT implementation)
- `97-acceptance-criteria.md` — verbatim §Acceptance Criteria 1–9 mapped to deliverables
- `98-changelog.md`
- `99-consistency-report.md`

### Phase 5 — Version bump + sync (on `next`)

- Bump `package.json` minor (per standing rule: any codebase change → minor bump)
- `node scripts/sync-version.mjs`
- `node scripts/sync-spec-tree.mjs`
- Update `.lovable/plan.md` with v-bump entry

---

## Replacers (per verbatim §Important.3)

| Term in user input | Canonical term used in spec |
|--------------------|------------------------------|
| `CW configuration` | `Seedable-Config` |
| `git map`          | `gitmap`                     |

---

## Cross-spec references

- Split-DB mechanics → `spec/05-split-db-architecture/` (do NOT redefine here)
- Seedable-Config mechanics → `spec/06-seedable-config-architecture/`
- Error handling generics → `spec/03-error-manage/`
- Self-update mechanism → `spec/14-update/` (this spec only points)
- Coding rules → `.lovable/coding-guidelines/coding-guidelines.md`

---

*Plan v1.0.0 — 2026-05-04*
