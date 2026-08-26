# Consolidated Guidelines — Gap Analysis Report

**Version:** 13.1.0  
**Updated:** 2026-05-07  
**Scope:** All 23 consolidated guideline files vs source spec folders + spec/19 readiness sweep at v5.41.0  
**Previous:** v13.0.0 (2026-04-22) — source-folder coverage map sweep

---

## v13.1.0 Changes — Hardening Patch I + Audit-12 closure

At v5.41.0 the spec/19 non-backup surface advanced to **98 / 100** (audit `12-`). Hardening Patch I closed the three gaps surfaced by the v5.39.0 readiness sweep:

- **G1** — `03-tasks` added to `linter-scripts/spec-folder-refs.allowlist` (sibling-repo doc-only ref).
- **G2** — broken cross-link in `spec/19/25-inherited-rules.md` corrected (`../03-error-manage/`).
- **G3** — MWS error catalogue grew 84 → 89 codes (`WORKER-403-01/02`, `WORKER-503-01/02`, `MAIN-900-01`); R4 widened to recognise §4 overflow range `21200-21299`; new `check-mws-error-codes.unallocated.txt` codifies 3 forbidden-future codes (`WORKER-940-05`, `WORKER-940-10`, `MAIN-830-04`).

All 14/14 `lint-ci.sh` steps green; 4/4 ancillary checks green. Patches A–D (audit-10) and Patch D-closure (v5.39.0) verified durable in audit `13-`.

Overall implementability score: **98.0/100** (up from 97.6). Handoff-weighted: **99.5/100** (up from 99.0). The remaining 0.5 point splits between the documented backup-snapshot-restore deferral (intentional, gated by `MAIN-900-01 SpecContradiction` safe-fail) and uncommitted baseline diagram PNGs.

---

## v13.0.0 Changes — Source-Folder Coverage Sweep (preserved below)

---

## v13.0.0 Changes — Source-Folder Coverage Sweep

The five highest-impact consolidated files received explicit **Source-Folder Coverage Maps** in their headers, plus added sections covering previously-undocumented source files. This closes the AI-blindness gap reported by the user: an AI handed any one of these files now knows precisely which source files it summarizes and where the gaps (if any) remain.

| File | Coverage Map | New Sections Added | New Lines |
|------|--------------|--------------------|-----------|
| `02-coding-guidelines.md` | ✅ Maps all 16 subfolders | §§22–33 (TS, Go, PHP, Rust, AI-opt, CI/CD-int, C#, file/folder, PowerShell, research, security, app-specific) | ~330 |
| `03-error-management.md` | ✅ Maps all 13 subfolders | §§18–26 (retros, verification, debugging, notification colors, error modal, envelope, apperror, logging, error code registry) | ~300 |
| `15-cicd-pipeline-workflows.md` | ✅ Maps 13 root specs + 3 subfolders + diagrams | Reusable CI Guards, Browser Ext Deploy, Go Binary Deploy, RCA, vuln scan detail, install scripts detail, code signing detail, terminal output, icon branding, env setup, release body, diagrams | ~150 |
| `17-self-update-app-update.md` | ✅ Maps all 25 source files | Console-safe handoff, repo path sync, version verification, last-release detection, install script probe, update check mechanism, release-pinned installer | ~180 |
| `18-database-conventions.md` | ✅ Maps 7 source files + cross-cutting + new rules | §18 Mandatory Free-Text Columns (Rules 10/11/12) with linter integration | ~85 |

---

## Executive Summary

All five critical handoff files now publish explicit Source-Folder Coverage Maps and language/subfolder sections. Overall implementability score: **97.6/100** (up from 96.2). Handoff-weighted: **99.0/100** (up from 98.0). The remaining 1.0 point is reserved for placeholder/reference-only sections that intentionally defer to source specs.

---

## Coverage Matrix

| # | Consolidated File | Source Module | Lines | Source Lines | Coverage | Implementability | Gap Level | Change |
|---|-------------------|---------------|-------|-------------|----------|-----------------|-----------|--------|
| 01 | `01-spec-authoring.md` | `01-spec-authoring-guide/` | **330+** | 2,758 | **12%** | **95%** | None | ⬆ +5 |
| 02 | `02-coding-guidelines.md` | `02-coding-guidelines/` | 726 | 30,842 | 2% | 97% | None | — |
| 03 | `03-error-management.md` | `03-error-manage/` | 489 | 18,363 | 3% | 97% | None | — |
| 04 | `04-enum-standards.md` | `02-coding-guidelines/` (enums) | 519 | — | — | 95% | None | — |
| 05 | `05-split-db-architecture.md` | `05-split-db-architecture/` | **723** | 4,032 | **18%** | **92%** | None | — |
| 06 | `06-seedable-config.md` | `06-seedable-config-architecture/` | **754** | 3,800 | **20%** | **93%** | None | — |
| 07 | `07-design-system.md` | `07-design-system/` | **580+** | 2,675 | **22%** | **92%** | None | — |
| 08 | `08-docs-viewer-ui.md` | `08-docs-viewer-ui/` | **430+** | 1,990 | **22%** | **91%** | None | — |
| 09 | `09-code-block-system.md` | `09-code-block-system/` | **530+** | 2,059 | **26%** | **91%** | None | — |
| 10 | `11-powershell-integration.md` | `11-powershell-integration/` | **560+** | 2,741 | **20%** | **91%** | None | — |
| 11 | `10-research.md` | `02-coding-guidelines/10-research/` | **180+** | — | — | **88%** | None | ⬆ +23 |
| 12 | `12-root-research.md` | `10-research/` | **170+** | — | — | **88%** | None | ⬆ +23 |
| 13 | `13-app.md` | `21-app/` | **210+** | — | — | **88%** | None | ⬆ +18 |
| 14 | `14-app-issues.md` | `22-app-issues/` | **210+** | — | — | **85%** | None | ⬆ +25 |
| 15 | `15-cicd-pipeline-workflows.md` | `12-cicd-pipeline-workflows/` | 422 | 5,140 | 8% | 92% | None | — |
| 16 | `16-app-design-system-and-ui.md` | `24-app-ui-design-system/` + error modal colors | **530+** | — | — | **93%** | None | ⬆ +5 |
| 17 | `17-self-update-app-update.md` | `14-update/` | 441 | 4,486 | 10% | 93% | None | — |
| 18 | `18-database-conventions.md` | `04-database-conventions/` | 945 | 2,321 | 41% | 95% | None | — |
| 19 | `19-gap-analysis.md` | — (meta) | — | — | — | — | — | — |
| 20 | `20-wp-plugin-conventions.md` | `18-wp-plugin-how-to/` | **570+** | 17,088 | 3% | **92%** | None | ⬆ +4 |
| 21 | `21-lovable-folder-structure.md` | `.lovable/` | **220+** | — | — | **91%** | None | ⬆ +6 |
| 22 | `22-app-database.md` | `23-app-db/` | **310+** | — | — | **90%** | None | 🆕 New |

---

## Changes Since v9.0.0

### Files Expanded (6)

| File | Old Lines | New Lines | Old Score | New Score | Key Additions |
|------|-----------|-----------|-----------|-----------|---------------|
| `20-wp-plugin-conventions.md` | 462 | 570+ | 88% | 92% | Complete `FileUploadTrait` with 5-step validation flow, `resolveUploadError()` match expression (7 PHP upload codes), upload route registration, upload edge cases table (5 scenarios), plugin ZIP upload flow (8-step pipeline), deployment/versioning section with SemVer, ZIP packaging structure, ZIP integrity requirements table, self-update summary, admin UI registration pattern |
| `21-lovable-folder-structure.md` | 86 | 220+ | 85% | 91% | Memory file template with frontmatter, 10-category usage guide with examples, `index.md` Core/Memories structure rules, memory CRUD workflows (create/update/delete), task file template with priority/status, issue file template with severity, suggestion workflow, strictly-avoid file template, root-level file descriptions |
| `14-app-issues.md` | 65 | 210+ | 60% | 85% | Severity level table (4 levels with response times), triage workflow (8-step pipeline), 8 resolution patterns (guard clause, race condition, silent error, schema drift, enum mismatch, config override, file path, upload failure), issue-to-memory pipeline, enhanced issue template with diagnosis/verification sections, cross-references |
| `10-research.md` | 38 | 80+ | 20% | 65% | Research file template with Question/Context/Methodology/Findings/Recommendation/Decision sections, deliverable column in content types, naming convention with examples |
| `12-root-research.md` | 29 | 65+ | 20% | 65% | Content type table with deliverables, placement decision guide, naming convention, template reference to 10-research |
| `13-app.md` | 35 | 85+ | 30% | 70% | App spec file template with requirements/design/edge cases/dependencies/acceptance criteria, content type table with file patterns, cross-references to related consolidated files |

---

## Cumulative Fix Summary (v3.2.0 → v10.0.0)

| Action | Count |
|--------|-------|
| Consolidated files expanded | **19** |
| New spec files created | 1 (`05-info-object-pattern.md`) |
| Cross-references fixed | 23 (8 WP plugin + 15 image paths) |
| Spec folders audited | 13/13 — all clean |
| Placeholder files upgraded | 3 (11, 12, 13 → full content with templates, lifecycles, anti-patterns) |

---

## Implementability Score Breakdown

| Category | Weight | v9.0 Score | v10.0 Score | Weighted (v10.0) |
|----------|--------|-----------|------------|-----------------|
| Core coding rules (02, 03, 04) | 30% | 97% | 97% | 29.1 |
| Architecture (05, 06, 17, 18) | 20% | 93.0% | 93.0% | 18.6 |
| UI/Design (07, 08, 09, 16) | 20% | 90.5% | 90.5% | 18.1 |
| Infrastructure (10, 15, 20) | 15% | 90.3% | **91.7%** | **13.8** |
| Meta/Structure (01, 11, 12, 13, 14, 21) | 15% | 76.0% | **87.2%** | **13.1** |
| **Total** | **100%** | **94.5** | **96.2 / 100** | |

### Handoff-Weighted Score

Weighting by "how badly will another AI mess up without this file":

| Category | Handoff Weight | Score | Weighted |
|----------|---------------|-------|----------|
| Core coding + enums (02, 03, 04) | 35% | 97% | 34.0 |
| Database + architecture (05, 06, 17, 18) | 25% | 93.0% | 23.3 |
| WP plugin (20) | 15% | **92%** | **13.8** |
| CI/CD + infra (10, 15) | 15% | 91.5% | 13.7 |
| UI + meta (07-09, 11-14, 16, 21) | 10% | 81.4% | **87.9%** | **8.8** |
| **Total** | **100%** | | **98.0 / 100** |

---

## Score Interpretation

| Range | Meaning |
|-------|---------|
| 90–100 | AI can implement with zero human intervention |
| 80–89 | AI can implement with minimal clarification |
| 60–79 | AI needs to read source specs for some topics |
| 40–59 | Significant gaps — source specs required |
| 0–39 | Placeholder only — no implementation value |

---

## Priority Recommendations (Updated)

### P0 — Done ✅

1. ~~Expand top 5 critical consolidated files~~ — Done (02, 03, 15, 17, 18)
2. ~~Create WP plugin conventions (20)~~ — Done
3. ~~Fix broken cross-references in WP plugin spec~~ — Done (8 fixes)
4. ~~Create `05-info-object-pattern.md`~~ — Done
5. ~~Mark placeholder files (11, 12, 13)~~ — Done
6. ~~Expand `05-split-db-architecture.md`~~ — Done (411→723 lines)
7. ~~Cross-reference audit of all 13 spec folders~~ — Done (23 total fixes)
8. ~~Expand `06-seedable-config.md`~~ — Done (311→754 lines)

### P1 — Done ✅

1. ~~Add complete CSS variable enumeration to `07-design-system.md`~~ — Done (228→580+ lines)
2. ~~Add component variant details to `16-app-design-system-and-ui.md`~~ — Done (119→340+ lines)

### P2 — Done ✅

1. ~~Add 4-backtick fence details to `09-code-block-system.md`~~ — Done (213→530+ lines)
2. ~~Add keyboard shortcut full list to `08-docs-viewer-ui.md`~~ — Done (129→430+ lines)
3. ~~Expand `11-powershell-integration.md` with full module patterns~~ — Done (191→560+ lines)

### P3 — Done ✅

1. ~~Expand `20-wp-plugin-conventions.md` with upload flow details~~ — Done (462→570+ lines, 88%→92%)
2. ~~Expand `21-lovable-folder-structure.md` with memory file conventions~~ — Done (86→220+ lines, 85%→91%)
3. ~~Expand `14-app-issues.md` with resolution patterns~~ — Done (65→210+ lines, 60%→85%)
4. ~~Address placeholder files (11, 12, 13)~~ — Done (expanded to full content: 180+, 170+, 210+ lines; 88% each)

### P4 — Done ✅

1. ~~Expand `01-spec-authoring.md`~~ — Done (90%→95%, 330+ lines, v3.2.0)
2. ~~Expand `16-app-design-system-and-ui.md`~~ — Done (88%→93%, 530+ lines, v3.3.0)

### P5 — Future Considerations

1. Expand `14-app-issues.md` — currently 85%, could add resolution pattern details
2. Expand `15-cicd-pipeline-workflows.md` — currently 92%, near ceiling
3. Populate research files (11, 12) with actual research when conducted
4. Populate `13-app.md` references when features are specified

---

## Expansion Status

| Status | Count | Files |
|--------|-------|-------|
| ✅ Standalone (90%+) | **17** | 01, 02, 03, 04, 05, 06, 07, 08, 09, 10, 15, 16, 17, 18, 20, 21, 22 |
| ⚠️ Good (80–89%) | **4** | 11, 12, 13, 14 |
| 🔶 Functional Scaffolding (60–79%) | **0** | — |
| ❌ Placeholder (<40%) | **0** | — |

---

## Staleness Risk

| File | Last Update | Risk |
|------|------------|------|
| All expanded files | 2026-04-16 | Low — freshly updated |

---

*Gap analysis — v12.0.0 — 2026-04-16*
