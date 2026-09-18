# Gold Standard Spec — Changelog

All notable changes to the WordPress plugin development specification.

---

## [1.3.0] — 2026-04-14

### Added

- Phase 20: End-to-End Walkthrough (`20-end-to-end-walkthrough.md`)
- Complete "Task Tracker" plugin built from scratch in 14 steps
- Phase coverage matrix mapping all 19 phases to walkthrough steps
- Final "Is My Plugin Gold Standard?" checklist

---

## [1.2.0] — 2026-04-14

### Priority 1 Fixes

| Phase | Section | Addition |
|-------|---------|----------|
| 04 | §4.8 | **Forbidden error patterns** — 8 anti-patterns with correct alternatives |
| 04 | §4.10 | **Stack trace transport format** — JSON structure for AJAX/REST delivery |
| 12 | §2.4 | **Dark mode token overrides** — `:root.dark-mode` CSS variable layer |
| 12 | §3.4 | **Slug substitution guide** — `{plugin-slug}` replacement rules for CSS classes, option names, text domains |
| 16 | §16.10 | **ErrorResponse class** — immutable value object with `toArray()` envelope |
| 16 | §16.11 | **AdminErrorAjaxTrait** — complete AJAX handler for error log viewer (read/clear/download) |
| 16 | §16.12 | **admin-errors.php template** — tabbed error viewer partial with JS integration points |

### Priority 2 Fixes

| Phase | Section | Addition |
|-------|---------|----------|
| 08 | §8.1 | **Admin menu error count badge** — `update-plugins` CSS class pattern with `wp_options` storage |
| 11 | §11.4 | **Complete partial example** — full data flow from orchestrator to partial with escaping rules |
| 16 | §16.13 | **ErrorSessions SQLite migration** — `DatabaseMigrationsErrorSessionsTrait` with `TableType::ErrorSessions` |
| 17 | §17.2 | **colors.json formal JSON Schema** — hex pattern validation + `ColorConfig` static-cache helper |

### Validation

- Cross-reference pass: **191 sections verified, 0 broken links**
- Implementability audit: overall AI success probability **83.0% → 89.1%** (+6.1pp)
- All 19 active phases now score ≥85% confidence

---

## [1.1.0] — 2026-04-09

### Added

- Phase 19: Micro-ORM & Root Database (`19-micro-orm-and-root-db.md`)
- Fluent query builder, `TypedQuery`, Go-style `DbResult`/`DbResultSet`/`DbExecResult` wrappers
- `FileCache` with scan/store trait decomposition
- Cross-plugin `RootDb` manifest pattern

---

## [1.0.0] — 2026-03-15

### Added

- Initial 18-phase Gold Standard specification (Phases 0–18)
- Phase 2 subfiles: enum architecture, metadata pattern, `SelfUpdateStatusType`, `ActionType`
- Complete reading order, cross-reference table, and phase index

---

*Update this file when spec phases are added, modified, or restructured.*
