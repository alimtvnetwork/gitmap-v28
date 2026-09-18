# Audit 16 — Diagram PNG Regeneration Baseline (v5.45.0)

**Date:** 2026-05-07
**Spec version:** v5.45.0
**Scope:** Single-command rebaseline via `scripts/regen-diagrams-and-audit.mjs`.

## Result

- **Pipeline:** parse-validate (`scripts/validate-mermaid.mjs`) → render (`scripts/render-diagrams.mjs`) → drift-check (`--check`).
- **Sources scanned:** 24 `.mmd` files under `spec/**/{diagrams,images}/`.
- **Rendered this run:** 0
- **Skipped (already fresh):** 24
- **Failed:** 0
- **Total PNGs covered by .mmd sources:** 24
- **Drift-check after render:** PASS (no stale PNGs).

## Disposition

- Spec/19 SPEC-ONLY constraint preserved — only renderer output (PNG binaries) refreshed; no spec markdown changed by this script.
- Backup-tier deferral preserved (audit-12 §Residual point #1 unchanged).

## Reproducibility

```
npm run diagrams:rebaseline
```
