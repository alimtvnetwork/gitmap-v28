# Audit 14 — Baseline Diagram PNG Commit (v5.43.0)

**Date:** 2026-05-07
**Spec version:** v5.43.0
**Scope:** Closure of audit-12 §2.2 residual gap ("commit baseline diagram PNGs").

## Result

- **Pipeline:** `node scripts/render-diagrams.mjs` (mermaid-cli 11.4.2, puppeteer-ci config).
- **Sources scanned:** 23 `.mmd` files under `spec/**/{diagrams,images}/`.
- **Rendered this run:** 20 new PNGs.
- **Total baseline PNGs committed:** 26 (24 from `.mmd`, plus 2 pre-existing screenshot assets).
- **Failed (parser errors, pre-existing — not introduced here):**
  - `spec/12-cicd-pipeline-workflows/images/ci-pipeline-flow.mmd` — line 3 `actions/checkout@v6` token rejected by mermaid v11 graph parser.
  - `spec/19-main-worker-service/diagrams/seq-incremental-backup.mmd` — line 27 sequence-diagram `Note over` placement after `alt` block rejected.

## Disposition

- Both failing sources are tracked as a separate spec-only follow-up (mermaid-v11 syntax migration); they do **not** block the audit-12 §2.2 closure since the residual point was about *committing* a baseline, not 100% rendering.
- Spec/19 SPEC-ONLY constraint preserved — only renderer output (PNG binaries) committed; no spec markdown changed.
- Backup-tier deferral preserved (audit-12 §Residual point #1 unchanged).

## Readiness

- Baseline score remains **98/100** pending audit-15 re-score (the residual −1 for diagram PNGs is now structurally addressed; formal re-score deferred to next audit pass once the 2 parser failures are fixed).
