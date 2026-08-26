# 13 — Hardening Patch I Verification (Phase-13.7)

> **Scope:** Verify Hardening Patch I (G1+G2+G3 from sweep `10-`) stuck at v5.41.0; re-confirm patches A–D from audit-10 still hold.
> **Audit date:** 2026-05-07
> **Project version at audit:** **5.41.0**
> **Method:** spot-check each patched line + re-run all CI linters; cross-check audit-12 score recovery.

---

## TL;DR

| Metric | audit-11 (v5.34.0) | audit-12 (v5.41.0) | **this re-verification** | Δ |
| --- | ---: | ---: | ---: | ---: |
| `lint-ci.sh` steps green | 12/14 | 14/14 | **14/14** | held |
| MWS-error-codes violations | 16 | 0 | **0** | held |
| Spec-folder stale refs | 1 | 0 | **0** | held |
| Cross-spec readiness | 97 / 100 | 98 / 100 | **98 / 100** | held |
| Patches A–D durable | 4/4 | 4/4 | **4/4** | held |

**Headline:** Patch I held cleanly. Patches A–D from audit-10 + Patch D closure from v5.39.0 also held. No regressions.

---

## 1. Patch I verification (3 / 3 stuck)

| Patch | File | Pre-state | Post-state at v5.41.0 |
| --- | --- | --- | --- |
| I-G1 | `linter-scripts/spec-folder-refs.allowlist:31-32` | `03-issues` only | **`03-issues` + `03-tasks`** ✅ |
| I-G2 | `spec/19-main-worker-service/25-inherited-rules.md:26` | link target was one directory too deep | corrected to sibling-relative path (`../03-error-manage/` from `25-…`) ✅ |
| I-G3a | `spec/19-main-worker-service/13-error-codes.md:7` | "Worker tier `21000-21099`" | **"Worker tier `21000-21099` (primary) + `21200-21299` (overflow)"** ✅ |
| I-G3b | `spec/19-main-worker-service/13-error-codes.md` §2.5/§2.9/§3.12 | 5 codes missing | **`WORKER-403-01/02` + `WORKER-503-01/02` + `MAIN-900-01` catalogued** ✅ |
| I-G3c | `linter-scripts/check-mws-error-codes.py` + `…unallocated.txt` | R4 range narrow; no unallocated allowlist | **R4 widened (21200-21299); 3 unallocated codes recognised** ✅ |

**Spot-check linter rerun:** 0 fail across all 4 ancillary checks + 14/14 lint-ci.

---

## 2. Re-sweep results (non-spec/19 corpus, v5.41.0)

### 2.1 SHOULD distribution

Same 8 hits as audit-11 §2.1 — **all 4 operative ones (Patches E/F/G/H) confirmed already closed by earlier loops**. No new SHOULD regressions.

### 2.2 TBD / FIXME

Single benign lint-doc example unchanged. No new hits.

### 2.3 Spec/19 implementation drift

v5.40.0 Phase-14 kickoff scaffold attempt was reverted same-loop; constraint memory `mem://constraints/spec19-no-implementation` broadened to forbid issue-tracking + app-side trackers + typed task indexes derived from spec/19. No code surface remains.

---

## 3. Recommendation

The 2-point residual gap (`§2.1 backup-snapshot-restore deferral`, `§2.2 diagram-PNG regeneration`) is bounded and intentional. **No further patches required this loop.** Suggested next milestones (outside this audit):

- **Audit-14** — re-baseline at the next v6.x kickoff or when the partial-restore feature is unfrozen.
- **Optional:** commit a baseline `diagrams/*.png` set to close §2.2 mechanically (still optional adoption — `25-§6` already neutralises drift impact).

**Estimated post-loop readiness: 98 / 100 — held.**

---

## 4. What was *not* found (negative results)

- No regressions in any audit-10 / audit-11 patch (A–H).
- No new contradictions between memory rules and spec text.
- No `OQ-` / FIXME / new TBD markers.
- No stale spec-folder references.
- No spec/19 implementation files.
