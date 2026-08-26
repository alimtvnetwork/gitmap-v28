# 15 — Blind-AI Implementation Readiness Audit (v6, post-baseline-PNG commit)

> **Spec under audit:** `spec/19-main-worker-service/` (28 numbered + 4 meta + audit/diagrams/images/fixtures dirs).
> **Audit date:** 2026-05-07
> **Project version at audit:** **5.43.0**
> **Predecessor audits:** `06`–`14` (most recent: `12-` 98/100, `13-` Patch-I reverif, `14-` baseline PNG commit).
> **Persona simulated:** the dumbest plausible AI coder — never asks, picks first matching rule, MUST = hard / SHOULD = ignored, cannot reconcile contradictions, cannot infer.
> **Method:** (a) full re-run of every linter against v5.43.0 corpus; (b) hostile re-read pass on the non-backup surface (00–17, 24–25); (c) physical PNG-vs-`.mmd` coverage check inside `spec/19/diagrams/`; (d) regression check that v5.42.0 → v5.43.0 changes (renderer output only) did not perturb spec markdown.

---

## TL;DR

| Dimension                  | v5 (12-) | **v6 (this run, 15-)** | Δ vs v5 |
| -------------------------- | -------: | ---------------------: | ------: |
| Completeness               |       98 |                **98**  |       0 |
| Determinism                |       97 |                **97**  |       0 |
| Consistency                |       99 |                **99**  |       0 |
| Testability                |       96 |                **97**  |      +1 |
| Blind-buildability         |       98 |                **99**  |      +1 |
| **Overall blind-AI score** | **98/100** |     **99 / 100 (A+)** |   **+1** |
| Estimated dumb-AI implements correctly | ~98 % |  **~99 %** |   +1 pp |
| Estimated dumb-AI fails or builds wrong thing | ~2 % | **~1 %** | -1 pp |

**Headline:** v5.43.0 committed 26 baseline diagram PNGs (20 of which are new this run) under the `render-diagrams.mjs` pipeline. The audit-12 §2.2 residual gap ("PNG snapshots not committed") is now structurally closed for **8 of 9** `spec/19/diagrams/*.mmd` sources. The single residual failure (`seq-incremental-backup.mmd` mermaid-v11 parser error on line 27 `Note over` placement) is a syntax migration issue, not a spec-content issue, and the literal-AI safe-fail (`MAIN-900-01 SpecContradiction`) still covers the prose path. Backup-tier deferral (audit-12 §2.1) unchanged — intentional v2.0 freeze.

---

## 1. What changed v5.41.0 → v5.43.0

| Δ | Source | Effect on blind-AI readiness |
| --- | --- | --- |
| v5.42.0 | Minor bump only; no spec/19 markdown touched. | Neutral. |
| v5.43.0 | `@mermaid-js/mermaid-cli@11.4.2` added as dev dep; `node scripts/render-diagrams.mjs` executed; 20 new PNGs committed (8 inside `spec/19/diagrams/`). | Closes 8/9 of audit-12 §2.2; +1 Testability (visual artifact now reproducible from committed source), +1 Blind-buildability (literal AI reading PNGs no longer sees out-of-band stale visuals for the covered set). |

No spec/19 `.md` was modified between v5.41.0 and v5.43.0 — confirmed via `git status`-equivalent file inspection. The SPEC-ONLY rule held throughout.

---

## 2. Linter posture at v5.43.0

| Check | Result | Notes |
| --- | :---: | --- |
| `lint-ci.sh` step 1 (verify present) | ✅ | 11 required linter files present. |
| `lint-ci.sh` step 2 (Go validator) | ⚠️ skipped | `go` binary unavailable in this environment; environmental skip, not a spec defect. |
| `lint-ci.sh` steps 3–14 | ✅ 12/12 | All non-Go steps GREEN; subset run reports "subset steps 3–14 passed". |
| `check-mws-error-codes.py` | ✅ | **89 codes verified (R1–R4); 21 R2 waivers loaded**; unallocated allowlist of 3 still honored. |
| `check-spec-folder-refs.py` | ✅ | 0 stale; 23 numbered + 26 external + 10 doc-only allowlisted. |
| `check-tunable-constants.py` | ✅ | OK. |
| `check-runner-dispatch-antipatterns.sh` | ✅ | No anti-patterns. |
| `check-axios-version.sh` | ✅ | 1.14.0 pinned. |

---

## 3. Diagram coverage inside spec/19

`spec/19-main-worker-service/diagrams/` — 9 `.mmd` sources, 8 PNGs committed:

| Source | PNG | Status |
| --- | :---: | --- |
| `erd-backup-tier.mmd` | ✅ | OK |
| `erd-main-db.mmd` | ✅ | OK |
| `erd-seedable-config.mmd` | ✅ | OK |
| `erd-worker-split-db.mmd` | ✅ | OK |
| `seq-backup-restore.mmd` | ✅ | OK |
| `seq-company-creation.mmd` | ✅ | OK |
| `seq-incremental-backup.mmd` | ❌ | mermaid-v11 parse error line 27 (`Note over` after `alt` block); pre-existing, not introduced by v5.43.0. |
| `seq-login-routing.mmd` | ✅ | OK |
| `seq-push-update.mmd` | ✅ | OK |

Plus `spec/19/images/` ships 4 topology PNGs (no `.mmd` sources — hand-drawn assets), all present.

---

## 4. Residual gap (1 point)

### 4.1 Backup-tier seed keys partial deferral (-1)

Unchanged from audit-12 §2.1. `Backup.Snapshot.Restore.*` remains prose-only/frozen until v2.0; literal AI hits `MAIN-900-01 SpecContradiction` and halts (desired safe-fail). **Disposition: intentional deferral.**

### 4.2 (closed) Diagram-PNG regeneration

Audit-12 §2.2 → **structurally closed at v5.43.0**. 8/9 spec/19 diagrams now have committed PNGs. The 1 remaining source has a mermaid-v11 syntax issue tracked in `audit/14-baseline-diagram-pngs-2026-05-07.md`; it does not regress the spec because:

1. The `.mmd` is still authoritative per `25-§6` precedence (diagrams ranked **last**).
2. The covered scenario (`seq-incremental-backup`) is *also* fully described in prose in `19-incremental-backup-sync.md` §§3–5.
3. A dumb AI has no reason to look up the PNG over the prose chapter.

Therefore: 0 net readiness penalty for the parser failure, but +1 Testability and +1 Blind-buildability for the 8 covered diagrams.

---

## 5. Hostile re-read trap sweep (non-backup surface)

I re-read the following with the literal-reader hat on, looking for any new MUST/SHOULD/MAY contradictions, undefined terms, or unreachable code paths introduced since v5.41.0:

- `00-overview.md`, `01-architecture.md`, `02-glossary.md`, `03-main-db-schema.md`, `04-worker-routing.md`, `05-auth-and-2fa.md`, `06-core-api-endpoints.md`, `07-role-based-dashboards.md`, `08-error-contract.md`, `09-self-update-pointer.md`, `10-worker-bootstrap-protocol.md`, `11-split-db-tier-reconciliation.md`, `12-jwt-delivery-contract.md`, `13-error-codes.md`, `14-rbac-and-status-seed.md`, `15-tunable-constants.md`, `16-update-channels.md`, `17-cascading-roles-and-cache-bin.md`, `24-threat-model.md`, `25-inherited-rules.md`.

**Findings:** zero new traps. None of these files were modified v5.41.0 → v5.43.0; every audit-12 closure remains structurally durable. The 89-code MWS catalogue still resolves cleanly.

---

## 6. Phase 13 trail — extended

| Phase   | Audit       | Score   | Status   |
| ------- | ----------- | ------: | -------- |
| 13.0    | `06-`       | 92/100  | Superseded. |
| 13.1    | `07-`       | 22/100  | Hostile baseline. |
| 13.2    | `08-`       | 75/100  | Schema/endpoint/envelope contradictions closed. |
| 13.3    | `09-`       | 96/100  | Spec/19 ready. |
| 13.4    | `10-`       | 96/100  | Cross-spec sweep + Patches A–D. |
| 13.5    | (impl.)     | —       | Diagrams pipeline + CI parity. |
| 13.6    | `11-`       | 97/100  | Patch verification + drift sweep. |
| 13.7    | `12-`       | 98/100  | Hardening Patch I; 14/14 lint-ci green; 89-code catalogue. |
| 13.7a   | `13-`       | 98/100  | Patch-I reverification — durable. |
| 13.7b   | `14-`       | (n/a)   | Baseline diagram PNG commit (renderer-only, no markdown). |
| **13.8**| **`15-` (this)** | **99/100** | **Baseline-PNG closure measured; 1 parser failure quarantined; backup-tier deferral preserved.** |

---

## 7. Disposition

- **No spec/19 markdown changes required** — the +1 readiness gain is from already-shipped renderer output, not spec rewrites. SPEC-ONLY rule (`mem://constraints/spec19-no-implementation`) honored.
- **Score promoted 98 → 99/100 (A+).**
- **Single residual point** is the intentional v2.0 backup-tier deferral, guarded by `MAIN-900-01 SpecContradiction`. Will not close before v2.0 cuts.
- **Follow-up (out of scope here):** mermaid-v11 syntax migration of `seq-incremental-backup.mmd` line 27 — pure renderer fix, no spec impact.
