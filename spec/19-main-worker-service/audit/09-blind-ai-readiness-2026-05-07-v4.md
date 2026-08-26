# 09 — Blind-AI Implementation Readiness Audit (Phase-13.3 Wrap-Up)

> **Spec under audit:** `spec/19-main-worker-service/` (43 numbered + meta files, ~563 KB).
> **Audit date:** 2026-05-07
> **Project version at audit:** **5.27.0**
> **Predecessor audits:** `06-` (92, intent), `07-` (22, hostile-literal), `08-` (75, post-13.2), this `09-` (post-13.3 wrap-up).
> **Persona simulated:** the dumbest plausible AI coder — never asks, picks first matching rule, treats MUST as hard / SHOULD as ignored, cannot reconcile contradictions, cannot infer.

---

## TL;DR

| Dimension                  | v3 (08-) | **v4 (this run)** | Δ |
| -------------------------- | -------: | ----------------: | ---: |
| Completeness               |       78 |            **96** | +18 |
| Determinism                |       72 |            **95** | +23 |
| Consistency                |       80 |            **97** | +17 |
| Testability                |       70 |            **94** | +24 |
| Blind-buildability         |       74 |            **96** | +22 |
| **Overall blind-AI score** | **75 / 100 (B)** | **96 / 100 (A)** | **+21** |
| Estimated dumb-AI implements correctly | ~75 % | **~96 %** | +21 pp |
| Estimated dumb-AI fails or builds wrong thing | ~25 % | **~4 %** | -21 pp |

**Headline:** Phase 13.3 closed the diagram debt, shipped a full fixture suite, eliminated prose-vs-seed drift, inlined inherited rules into `25-`, and resolved every literal-reader trap surfaced by audit-08. The spec is now operationally tractable for a hostile-literal AI; the residual 4 points are intentional v2.0.0 deferrals (backup-tier seed keys behind a feature flag).

---

## 1. What Phase 13.3 fixed (mapped to audit-08 findings)

| Audit-08 §  | Issue (cost)                                              | Phase-13.3 fix                                                                                       | Pts recovered |
| ----------- | --------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ------------: |
| 2.1         | Diagram debt: 3-tier DB, `User` boxes, no backup lane (-8) | Refreshed `erd-main-db.mmd`, `erd-worker-split-db.mmd`, `seq-login-routing.mmd` to 4-tier + UserDirectory/AppUser split + backup lane | +8 |
| 2.2         | No canonical fixtures for 12 endpoints / errors / JWT (-7) | New `fixtures/` tree: 7 endpoint pairs, 7 error envelopes, RS256 JWT example + claim-ordering doc, plus `conventions.md` | +7 |
| 2.3         | 3 prose-only seed keys (-4)                                | Verified 2 were aliases; materialized `CacheRecentCompanyPerUserTtlSeconds` in `config.seed.json` v1.5.0; backup-tier keys deferred to v2.0.0 with linter waiver | +3 (1 deferred) |
| 2.4         | Inherited rules referenced by link only (-3)               | Created `25-inherited-rules.md` inlining operative subset of specs `03/04/05/06` + core memory; `00-§7` re-prioritized | +3 |
| 2.5         | Literal-reader traps (SHOULD→MUST, hidden §6.4, OQ-17-2) (-3) | Audited claims: `05-§2.1` already MUST (false positive), `OQ-17-2` does not exist (false positive), `22-§6` preamble added pointing to §6.4 WAL ban | +3 |

Subtotal recovered: **+24 raw**, net **+21** after applying the 1-point residual deferral cap.

---

## 2. Residual gap (4 points)

### 2.1 Backup-tier seed keys deferred to v2.0.0 (-3)

~28 `MainWorker.Backup.*` tunables remain prose-only in `15-` because the entire backup-promotion subsystem is feature-flagged off (per D9 — no auto-failover). Linter check T3 is waived for this namespace until the v2.0.0 seed bump. A blind AI building only the in-scope (non-backup) surface is unaffected; one building backup logic will (correctly) hit a `MAIN-900-01 SpecContradiction` and stop, which is the desired safe-fail behaviour.

**Disposition:** intentional deferral, not a defect. Will be closed in Phase 14.

### 2.2 Diagram-PNG regeneration (-1)

The `.mmd` source files are correct as of v5.26.0; rendered `.png` snapshots in `diagrams/` are not regenerated automatically by CI. A literal AI reading PNGs (rare but possible) could see stale visuals. Mitigation: `25-§6` precedence rules rank diagrams **last**, so prose wins.

**Disposition:** non-blocking; tracked for the docs-viewer build pipeline, not the spec.

---

## 3. Sub-scores explained

| Dimension       | Score | Justification |
| --------------- | ----: | ------------- |
| Completeness    | 96    | All 43 spec files mapped in `00-§6`; only deferred backup tunables prose-only. |
| Determinism     | 95    | Single SoT for schema, endpoints, error envelope, JWT; precedence chain in `25-§6`. |
| Consistency     | 97    | No contradictions remain between `03/04/06/15/18/22`; aliases formalized. |
| Testability     | 94    | Fixture suite + `conventions.md` covers happy path, error envelopes, and JWT shape. |
| Blind-buildability | 96 | Inherited rules inlined; precedence + escape hatch (`MAIN-900-01`) explicit. |

---

## 4. Phase 13 trail — closed

| Phase   | Audit       | Score   | Status   |
| ------- | ----------- | ------: | -------- |
| 13.0    | `06-`       | 92/100  | Superseded by hostile-literal re-baseline. |
| 13.1    | `07-`       | 22/100  | Hostile baseline; drove the hardening backlog. |
| 13.2    | `08-`       | 75/100  | Schema/endpoint/envelope contradictions closed. |
| **13.3**| **`09-` (this)** | **96/100** | **Wrap-up. Ready for Phase-14 implementation kickoff.** |

---

## 5. Recommendation

`spec/19-main-worker-service/` is **cleared for blind-AI implementation** of the in-scope (non-backup) surface at v5.27.0. Phase 14 may open. Backup-tier work remains gated behind the v2.0.0 seed bump and an explicit feature flag.

**No further audit cycles required for Phase 13.** Reopen only on material spec change (new endpoint, schema, or tier).
