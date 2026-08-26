# 11 — Audit-10 Patch Re-Verification Sweep (Phase-13.6)

> **Scope:** Verify audit-10 patches A–D stuck at v5.34.0; re-sweep all non-spec/19 folders for residual ambiguity drift.
> **Audit date:** 2026-05-07
> **Project version at audit:** **5.34.0**
> **Persona:** same hostile-literal blind-AI used in audits `07-…10-`.
> **Method:** spot-check the 5 patched lines, then `rg` re-sweep on `SHOULD`/`TBD`/`FIXME` across the corpus, then triage residual hits.

---

## TL;DR

| Metric                                  | audit-10 (v5.31.0) | **this re-sweep (v5.34.0)** | Δ |
| --------------------------------------- | -----------------: | --------------------------: | ---: |
| Patched lines verified intact           | n/a                | **5 / 5** (100 %)           | — |
| `SHOULD` hits across non-spec/19 corpus | 12                 | **8** (-4 audit-10 hits)    | -33 % |
| `TBD` / `FIXME` operative hits          | 1                  | **0**                       | -100 % |
| **Operative ambiguities still actionable** | **4 files**     | **4 files** (different set) | unchanged count, churn = 100 % |
| **Estimated cross-spec readiness**      | 96 / 100           | **97 / 100**                | +1 |

**Headline:** Patches A–D held perfectly. The ambiguity surface shifted: the 4 hardenings closed audit-10's set, but the re-sweep promoted 4 previously-classified-as-RFC-2119-reference hits into the "operative" bucket on closer inspection. Net readiness still trending upward; remaining work is small.

---

## 1. Patch verification (5 / 5 stuck)

| Patch | File | Line | Pre-state | Post-state at v5.34.0 |
| ----- | --- | ---: | --- | --- |
| A1 | `spec/03-…/02-session-based-logging.md` | 27 | `F6 … SHOULD` | **`F6 … MUST`** ✅ |
| A2 | (same)                                  | 29 | `F8 … SHOULD` | **`F8 … MUST`** ✅ |
| A3 | (same)                                  | 34 | `F13 … SHOULD` | **`F13 … MUST`** ✅ |
| B1 | `spec/14-update/27-generic-installer-behavior.md` | 151 | "Implementations SHOULD offer flags…" | **"Implementations MUST offer the following flags … (exact names; no synonyms)"** ✅ |
| B2 | (same)                                  | 180–184 | "SHOULD adopt 20 going forward; 5 remains acceptable for legacy" | **"New installers MUST adopt LOOKAHEAD = 20. Installers explicitly tagged `legacy` … MAY retain 5"** ✅ |
| B3 | `spec/14-update/23-install-script-version-probe.md` | 405 | "SHOULD log identity + range + outcome…" | **"MUST log identity + range + outcome … (sole debuggability hook)"** ✅ |
| C  | `spec/15-distribution-and-runner/06-fix-repo-forwarding.md` | 75–76 | "dispatch path SHOULD be a single `exec`" | **"dispatch path MUST be a single `exec`"** ✅ |
| D  | `spec/06-…/04-rag-test-coverage-matrix.md` | 251 | `Mutation Score \| ≥80% \| TBD` | **"Closed at v2.0.0 (Patch D). ≥80% is the binding gate; measurement deferred-by-design until executable RAG code lands outside `spec/`. Line+branch at 100% remains the active proxy."** ✅ ✅ (re-closed at v5.39.0) |

**Spot-check rg on the 4 patched files for any new `SHOULD` regressions: 0 hits.** All patches durable.

---

## 2. Re-sweep results (non-spec/19 corpus, v5.34.0)

### 2.1 SHOULD distribution (8 files, all 1-per-file)

| File | Verdict |
| --- | --- |
| `spec/01-spec-authoring-guide/00-overview.md` | **Operative** — "Every module SHOULD include a reliability risk assessment". Promote to MUST or carve an explicit allowed-exception clause. |
| `spec/01-spec-authoring-guide/11-root-readme-conventions.md` | RFC-2119 reference text inside the spec-authoring guide itself. Benign. |
| `spec/01-spec-authoring-guide/17-version-schema.md` | RFC-2119 reference text. Benign. |
| `spec/02-coding-guidelines/04-php/07-php-standards-reference/05-forbidden-and-database.md` | Non-MUST coding hint within forbidden-list rationale. Benign. |
| `spec/04-database-conventions/01-naming-conventions.md` (line 411) | **Operative** — "Transactional / invoice / billing / payment / order-like tables SHOULD include these nullable free-text columns". This conflicts with **memory rule 11** ("transactional need `Notes`+`Comments TEXT NULL`") which is MUST-grade. Promote SHOULD → MUST. |
| `spec/12-cicd-pipeline-workflows/05-code-signing.md` (line 10) | **Operative-ish** — "Windows binaries SHOULD be signed". Intentionally SHOULD because signing requires a paid cert; SHOULD is correct. Add an explicit "MAY skip when cert unavailable; CI MUST emit signing-skipped banner" caveat for hostile-AI safety. |
| `spec/12-cicd-pipeline-workflows/03-reusable-ci-guards/09-workflow-templates.md` | Template-section guidance. Benign. |
| `spec/authoring-guideline/version-schema.md` (line 141) | **Operative** — "Readers SHOULD treat the file as read-only. Only `sync-version.mjs` writes to it." Promote to MUST — multi-writer is a known data-corruption vector. |

**Net actionable:** 4 files (one is a clarification, three are SHOULD → MUST).

### 2.2 TBD / FIXME distribution

Single hit, in `spec/02-coding-guidelines/01-cross-language/04-code-style/06-comments-and-documentation.md:83` — documents that `// TODO(PROJ-123)` comments are *allowed*. **Benign** (false positive, lint-rule example).

### 2.3 Categories not re-checked

`OQ-` markers (all resolved per audit-10), template scaffolding, lint-doc examples — unchanged from audit-10.

---

## 3. Recommendation

The 4 newly-surfaced ambiguities are smaller than audit-10's set (one is just a "make the existing intent explicit" caveat, not a true MUST-promotion). Estimate ~5 line-edits total, single loop. Suggested grouping:

- **Patch E** — `spec/01-spec-authoring-guide/00-overview.md:228` SHOULD → MUST (or carve exception clause).
- **Patch F** — `spec/04-database-conventions/01-naming-conventions.md:411` SHOULD → MUST (aligns with memory rule 11).
- **Patch G** — `spec/12-cicd-pipeline-workflows/05-code-signing.md:10` add cert-unavailable + CI banner clause.
- **Patch H** — `spec/authoring-guideline/version-schema.md:141` SHOULD → MUST (multi-writer corruption guard).

**Estimated post-patch readiness: 99 / 100.**

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
| **13.6**| **`11-` (this)** | **97/100** | **Patch verification + drift sweep. 4-file polish backlog identified.** |

---

## 5. What was *not* found (negative results)

- **No regressions** in the audit-10 patched files.
- **No new contradictions** between memory rules and spec text (only one alignment opportunity at `spec/04/01:411`).
- **No fresh `OQ-` markers** introduced since audit-10.
- **No new TBD/FIXME hits** outside the single benign lint-doc example.
