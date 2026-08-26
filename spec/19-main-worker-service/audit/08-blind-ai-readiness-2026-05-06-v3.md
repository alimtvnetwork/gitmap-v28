# 08 — Blind-AI Implementation Readiness Audit (Post-Phase-13.2 Re-Run)

> **Spec under audit:** `spec/19-main-worker-service/` (29 numbered files, ~370 KB).
> **Audit date:** 2026-05-06
> **Project version at audit:** **5.22.0**
> **Predecessor audits:** `06-` (92/100, intent-based), `07-` (22/100, hostile-literal baseline), this `08-` (post-hardening literal re-run).
> **Persona simulated:** the dumbest plausible AI coder — never asks, picks first matching rule, treats MUST as hard / SHOULD as ignored, cannot reconcile contradictions, cannot infer.

---

## TL;DR

| Dimension                  | v2 (07-) | **v3 (this run)** | Δ |
| -------------------------- | -------: | ----------------: | ---: |
| Completeness               |       30 |            **78** | +48 |
| Determinism                |       10 |            **72** | +62 |
| Consistency                |        5 |            **80** | +75 |
| Testability                |       45 |            **70** | +25 |
| Blind-buildability         |       20 |            **74** | +54 |
| **Overall blind-AI score** | **22 / 100 (F)** | **75 / 100 (B)** | **+53** |
| Estimated dumb-AI implements correctly | ~20 % | **~75 %** | +55 pp |
| Estimated dumb-AI fails or builds wrong thing | ~80 % | **~25 %** | -55 pp |

**Headline:** Phase 13.2 hardening removed the hostile contradictions. The spec is now operationally tractable for a literal AI; the residual 25-point gap is mostly **diagram/test-vector debt** plus a few **prose-vs-JSON drift** spots in tunables, not structural ambiguity.

---

## 1. What Phase 13.2 fixed (and the score moved on)

| Fix applied (v5.22.0) | Audit line item closed | Pts recovered |
| --- | --- | ---: |
| `WorkerNode` schema dedup → single SoT in `03-§2.1` | "Schema contradicts itself" (was -18) | +18 |
| Endpoint catalogs in `04-`/`18-` removed → SoT `06-` | "POST /Workers/Register defined 4 ways" (was -14) | +14 |
| `EndpointAuthSetting` FK → `UserDirectory` | "FK references deleted `User` table" (was -6) | +6 |
| `15-§4.1` prose↔seed alias map (32 keys) + §4.2 cache tunables | "Config keys don't match JSON" (was -8) | +6 (partial; see §3) |
| `11-` stubbed as **applied**, not a TODO list | "Reconciliation listed as follow-ups" (was -9) | +8 |
| `08-§2` pinned as sole error envelope; `13-` disclaims | "Two error envelope shapes" (was -4) | +4 |
| `00-§6` Document Map covers all 25 + meta files | "Hidden files outside map" (was -3) | +3 |
| STALE banners on `96-`, `99-`, `audit/06-` | "Stale audits read as authoritative" (was -4) | +3 |

Subtotal recovered: **+62 raw points**, capped at +53 net by residual gaps below.

---

## 2. What still costs the dumb AI ~25 points

### 2.1 Diagram debt (-8)

`spec/19-main-worker-service/diagrams/` was not refreshed in Phase 13.2. The Mermaid/PNG diagrams still show:
- 3-tier DB instead of 4-tier (Root/Settings/App/Session per FU-1).
- `User` table boxes instead of `UserDirectory` / `AppUser` split.
- No backup-node lane in the routing sequence diagram.

A literal AI that prefers diagrams over prose (common failure mode) will replicate the stale model.

### 2.2 Test-vector / fixture gap (-7)

- `97-acceptance-criteria.md` lists ACs but provides **no canonical request/response JSON fixtures** for the 12 core endpoints in `06-`.
- `08-error-contract.md` defines the envelope shape but ships **no golden-file examples** keyed to `13-error-codes.md` codes.
- `12-jwt-delivery-contract.md` lacks a worked HS256/RS256 token example with claim ordering.

A blind AI will invent fixtures and drift from intent on field ordering, timestamp format (epoch vs ISO), and null-vs-omitted.

### 2.3 Prose-vs-seed drift residue (-4)

`15-§4.1` alias map covers the 32 collisions caught in 07-, but spot-check finds **3 keys still only in prose**:
- `MainWorker.Cache.WorkerToCompanyTtlSeconds` (mentioned `01-§5`, missing from `config.seed.json`).
- `MainWorker.Routing.BackupPromotionGraceSeconds` (`18-§3`, not seeded).
- `MainWorker.Sync.IncrementalBatchMaxRows` (`19-§4`, not seeded).

### 2.4 Cross-spec inlining incomplete (-3)

`11-` now declares reconciliation **applied**, but several rules it inherits from `spec/01-app/` and `spec/03-tasks/` are still **referenced by link**, not inlined. A blind AI handed only the `spec/19-…` folder cannot resolve them.

### 2.5 Minor literal-reader traps (-3)

- `05-§2.1` uses "SHOULD" for the proxy-credentialed-body rule that the rest of the spec treats as MUST.
- `22-§6.4` (the new WAL prohibition) lives in a sub-section a literal AI may skip if it stops reading at §6.
- `17-cascading-roles-and-cache-bin.md` §3 still has one un-resolved OQ marker (`OQ-17-2`) without a default.

---

## 3. Sub-scores explained

| Dimension | Score | Why not 100 |
| --- | ---: | --- |
| Completeness | 78 | Missing canonical fixtures (§2.2), 3 unseeded keys (§2.3), unresolved OQ-17-2. |
| Determinism | 72 | One SHOULD-that-should-be-MUST (§2.5), diagrams drift from prose (§2.1). |
| Consistency | 80 | Schema/endpoint/error envelope all single-SoT now; only prose↔seed residue (§2.3) and diagram lag (§2.1) remain. |
| Testability | 70 | ACs exist but no machine-checkable golden files; no contract-test harness referenced. |
| Blind-buildability | 74 | Core API + DB + bootstrap now buildable end-to-end; fixtures gap forces guesswork at the wire layer. |

---

## 4. Estimated dumb-AI outcome (post-13.2)

- **Will build correctly (~75 %)**: schema, migrations, RBAC seed, routing strategies, register/heartbeat, error-envelope responses, JWT issuance contract shape, backup-node exclusion from routing, tier-3 split-DB layout.
- **Will build wrong or guess (~25 %)**: wire-format edge cases (timestamp encoding, null handling), backup promotion timing, three unseeded tunables defaulting to hard-coded literals, diagram-driven 3-tier regression, cross-spec inherited rules.

---

## 5. Path from 75 → 95+

Minimum edit set for **Phase 13.3** (estimated ~1 dev-session):

1. **Diagrams refresh** (`diagrams/*.mmd` + regenerated PNGs) → +8.
2. **Golden-file fixtures**: add `spec/19-main-worker-service/fixtures/` with one request/response per endpoint in `06-`, one per error code family in `13-`, one JWT example in `12-` → +7.
3. **Seed the 3 missing tunables** in `config.seed.json` and update `15-§4.1` alias map → +4.
4. **Inline cross-spec inherited rules** into `11-§6` (or a new `25-inherited-rules.md`) → +3.
5. **Tighten 3 literal-reader traps** (§2.5): promote SHOULD→MUST in `05-§2.1`, hoist `22-§6.4` to top-level §7, resolve OQ-17-2 with explicit default → +3.

Cap: 100. Realistic landing after Phase 13.3 = **95 ± 2**. Reaching 100 requires a contract-test harness wired into CI, which is out of spec scope.

---

## 6. Verdict

**Grade: B (75 / 100).** The spec moved from **operationally hostile** (F) to **operationally tractable** (B) in one hardening pass. The remaining gaps are documentation craftsmanship — fixtures, diagrams, three missed seeds — not structural ambiguity. A blind AI handed the current `spec/19-…/` folder will ship a working main-worker service with predictable, reviewable wire-format drift, instead of a broken core.

— end of audit 08 —
