# 12 — Blind-AI Implementation Readiness Audit (v5, post-Hardening Patch I)

> **Spec under audit:** `spec/19-main-worker-service/` (46 numbered + meta files).
> **Audit date:** 2026-05-07
> **Project version at audit:** **5.41.0**
> **Predecessor audits:** `06-` (92), `07-` (22), `08-` (75), `09-` (96), `10-` (96 cross-spec), `11-` (97 patch reverif), **this `12-`**.
> **Persona simulated:** the dumbest plausible AI coder — never asks, picks first matching rule, treats MUST as hard / SHOULD as ignored, cannot reconcile contradictions, cannot infer.
> **Method:** (a) full re-run of `lint-ci.sh` + `check-mws-error-codes.py` + `check-spec-folder-refs.py` + `check-tunable-constants.py` against v5.41.0; (b) hostile re-read of every chapter in non-backup surface (00–17, 24–25); (c) literal-reader trap sweep on chapters that gained content since audit-09.

---

## TL;DR

| Dimension                  | v4 (09-) | v5 cross-spec (10-) | **v5 (this run, 12-)** | Δ vs 09 |
| -------------------------- | -------: | ------------------: | ---------------------: | ------: |
| Completeness               |       96 |                  96 |                **98** |     +2 |
| Determinism                |       95 |                  95 |                **97** |     +2 |
| Consistency                |       97 |                  97 |                **99** |     +2 |
| Testability                |       94 |                  94 |                **96** |     +2 |
| Blind-buildability         |       96 |                  96 |                **98** |     +2 |
| **Overall blind-AI score** | **96/100** | **96/100** |        **98 / 100 (A+)** | **+2** |
| Estimated dumb-AI implements correctly | ~96 % | ~96 % | **~98 %** | +2 pp |
| Estimated dumb-AI fails or builds wrong thing | ~4 % | ~4 % | **~2 %** | -2 pp |

**Headline:** Hardening Patch I (G1+G2+G3) closed every red linter step. The MWS error catalogue grew from 84 → 89 codes (added `WORKER-403-01/02`, `WORKER-503-01/02`, `MAIN-900-01`), an explicit "documented-as-unallocated" allowlist now codifies the future-work labels (`WORKER-940-05/10`, `MAIN-830-04`) without forcing premature catalogue rows, and the §4 overflow range (`21200-21299`) is now first-class in both prose and linter R4. The residual 2 points are the same intentional v2.0.0 backup-tier deferrals previously documented.

---

## 1. What v5.40.0 → v5.41.0 fixed (mapped to audit-11 backlog + readiness sweep gaps)

| Gap (source) | Issue | Fix at v5.41.0 | Pts recovered |
| --- | --- | --- | ---: |
| G1 (sweep) | `audit/08-…v3.md` cited non-existent `spec/03-tasks/`; check-spec-folder-refs RED. | Added `03-tasks` to `[external]` allowlist (sibling-repo doc-only ref). | +0.5 |
| G2 (sweep) | `25-inherited-rules.md` link `[../03-error-manage/](../../03-error-manage/)` resolved one level too deep. | Fixed to `../03-error-manage/`. | +0.5 |
| G3 (sweep) | 16 MWS error-code violations: 8 missing (`MAIN-830-04`, `MAIN-900-01`, `WORKER-403-01/02`, `WORKER-503-01/02`, `WORKER-940-05/10`) + 8 R4 range mismatches (`WORKER-21200-21207`). | (a) Catalogued 5 real codes in §2.5/§2.9/§3.12. (b) Added `check-mws-error-codes.unallocated.txt` for the 3 forbidden-future codes; linter R1 honors it. (c) Header §1 + linter R4 widened to recognise `21200-21299` overflow per §4. | +1.0 |
| Patch E/F/H (audit-11) | SHOULD→MUST polish in 4 non-spec/19 files. | Verified all closed by earlier loops; backlog retired. | +0.0 |
| Patch G (audit-11) | code-signing.md cert-unavailable + CI banner clause. | Verified in place at `12-cicd-pipeline-workflows/05-code-signing.md:10`. | +0.0 |

Subtotal recovered: **+2 raw**, capped at the deferral ceiling (backup-tier T3).

---

## 2. Residual gap (2 points)

### 2.1 Backup-tier seed keys partial deferral (-1)

v5.38.0 materialised the 28 `MainWorker.Backup.*` tunables behind T3 silent-waiver; the cache-bin & rotation timers fully covered. One residual sub-bin (`Backup.Snapshot.Restore.*`) remains prose-only because the partial-restore feature is explicitly forbidden until v2.0 (per `24-threat-model.md` §1 and `97-acceptance-criteria.md` AC). A literal AI building this path will hit `MAIN-900-01 SpecContradiction` (now formally catalogued in §3.12) and halt — the desired safe-fail.

**Disposition:** intentional deferral.

### 2.2 Diagram-PNG regeneration (-1)

`.mmd` sources are authoritative; rendered `.png` snapshots are produced by the docs-diagrams CI job per push but are not committed to `diagrams/`. A literal AI reading PNGs out-of-band (rare path) could see stale visuals. Mitigation: `25-§6` precedence rules rank diagrams **last**.

**Disposition:** documented and bounded; closure deferred to optional baseline-PNG commit loop.

---

## 3. Linter posture at v5.41.0

| Check | v5.39.0 | v5.41.0 | Notes |
| --- | :---: | :---: | --- |
| `lint-ci.sh` (14 steps) | 12/14 | **14/14** | Steps #5 (spec-folder-refs) + #11 (mws-error-codes) flipped GREEN. |
| `check-mws-error-codes.py` | 16 fail | **0 fail** | 89 codes verified (R1–R4); 21 R2 waivers + 3 unallocated entries loaded. |
| `check-spec-folder-refs.py` | 1 stale | **0 stale** | 23 numbered + 26 external + 10 doc-only allowlisted. |
| `check-tunable-constants.py` | OK | OK | Unchanged. |
| `check-runner-dispatch-antipatterns.sh` | PASS | PASS | Unchanged. |

---

## 4. Phase 13 trail — extended

| Phase   | Audit       | Score   | Status   |
| ------- | ----------- | ------: | -------- |
| 13.0    | `06-`       | 92/100  | Superseded. |
| 13.1    | `07-`       | 22/100  | Hostile baseline. |
| 13.2    | `08-`       | 75/100  | Schema/endpoint/envelope contradictions closed. |
| 13.3    | `09-`       | 96/100  | Spec/19 ready. |
| 13.4    | `10-`       | 96/100  | Cross-spec sweep + Patches A–D applied. |
| 13.5    | (impl.)     | —       | Diagrams pipeline + CI parity. |
| 13.6    | `11-`       | 97/100  | Patch verification + drift sweep. |
| **13.7**| **`12-` (this)** | **98/100** | **Hardening Patch I; 14/14 lint-ci green; MWS catalogue at 89 codes.** |

---

## 5. What was *not* found (negative results)

- **No new SHOULD/TBD/FIXME hits** outside the benign single lint-doc example flagged in audit-11.
- **No regressions** in audit-10 patches A–D.
- **No new orphan codes**: every catalogued code resolves to ≥1 real reference or an explicit R2 waiver.
- **No flat-bijection violations** (R3 clean across all 89 codes).
- **No spec/19 implementation drift**: chapters remain markdown-only; the Phase-14 kickoff scaffold attempt at v5.40.0 was reverted before this audit and the constraint memory broadened to forbid recurrence.
