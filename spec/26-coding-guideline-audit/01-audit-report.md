# Master Coding Guideline Gap Audit Report

- **Audit Date:** 2026-08-29
- **Total Files Audited:** 3019
- **Overall Codebase Compliance Score:** 99.0 / 100 (🟢 Exemplary)
- **Total Violations Found:** Critical: 78 | Major: 101 | Minor: 930

## Executive Summary & Score Breakdown

| Module / Directory | Files | Critical (-10) | Major (-5) | Minor (-2) | Module Score | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `eslint.config.js` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix-repo.ps1` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `fix-repo.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `fix_cancelled.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix_envops.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix_lints2.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix_panics.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix_reinstall.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix_root.go` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `fix_task.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap-updater/cmd` | 6 | 0 | 0 | 2 | 99.3 / 100 | 🟢 Exemplary |
| `gitmap-updater/main.go` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/apperror` | 3 | 0 | 0 | 1 | 99.3 / 100 | 🟢 Exemplary |
| `gitmap/archive` | 6 | 1 | 0 | 4 | 97.0 / 100 | 🟢 Exemplary |
| `gitmap/cliexit` | 8 | 0 | 0 | 4 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/cloneconcurrency` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/clonefrom` | 40 | 0 | 0 | 22 | 98.9 / 100 | 🟢 Exemplary |
| `gitmap/clonenext` | 10 | 0 | 0 | 7 | 98.6 / 100 | 🟢 Exemplary |
| `gitmap/clonenow` | 22 | 1 | 0 | 14 | 98.3 / 100 | 🟢 Exemplary |
| `gitmap/clonepick` | 18 | 0 | 0 | 10 | 98.9 / 100 | 🟢 Exemplary |
| `gitmap/cloner` | 23 | 0 | 0 | 15 | 98.7 / 100 | 🟢 Exemplary |
| `gitmap/cluster` | 36 | 0 | 0 | 7 | 99.6 / 100 | 🟢 Exemplary |
| `gitmap/cmd` | 1643 | 28 | 0 | 425 | 99.3 / 100 | 🟢 Exemplary |
| `gitmap/committransfer` | 20 | 0 | 0 | 10 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/completion` | 15 | 2 | 0 | 8 | 97.6 / 100 | 🟢 Exemplary |
| `gitmap/config` | 6 | 0 | 0 | 6 | 98.0 / 100 | 🟢 Exemplary |
| `gitmap/constants` | 136 | 6 | 0 | 40 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/crypto` | 3 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/dashboard` | 9 | 0 | 0 | 5 | 98.9 / 100 | 🟢 Exemplary |
| `gitmap/db` | 11 | 0 | 0 | 5 | 99.1 / 100 | 🟢 Exemplary |
| `gitmap/desktop` | 5 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/detector` | 5 | 0 | 0 | 4 | 98.4 / 100 | 🟢 Exemplary |
| `gitmap/diff` | 4 | 0 | 0 | 2 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/downloaderconfig` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `gitmap/errreport` | 3 | 0 | 0 | 2 | 98.7 / 100 | 🟢 Exemplary |
| `gitmap/fix_dispatch.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/fixtureversion` | 7 | 0 | 0 | 5 | 98.6 / 100 | 🟢 Exemplary |
| `gitmap/formatter` | 18 | 0 | 0 | 11 | 98.8 / 100 | 🟢 Exemplary |
| `gitmap/fsutil` | 15 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/ghtoken` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/gitutil` | 24 | 1 | 0 | 2 | 99.4 / 100 | 🟢 Exemplary |
| `gitmap/glyphs` | 6 | 0 | 0 | 1 | 99.7 / 100 | 🟢 Exemplary |
| `gitmap/goldenguard` | 4 | 0 | 0 | 2 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/helptext` | 8 | 0 | 0 | 1 | 99.8 / 100 | 🟢 Exemplary |
| `gitmap/indexer` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `gitmap/installer` | 60 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/jsonenv` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/lazyregex` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/localdirs` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `gitmap/lockcheck` | 3 | 0 | 0 | 1 | 99.3 / 100 | 🟢 Exemplary |
| `gitmap/lockfile` | 2 | 0 | 0 | 2 | 98.0 / 100 | 🟢 Exemplary |
| `gitmap/logging` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/macro` | 5 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/main.go` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/mapper` | 10 | 0 | 0 | 2 | 99.6 / 100 | 🟢 Exemplary |
| `gitmap/model` | 38 | 0 | 0 | 1 | 99.9 / 100 | 🟢 Exemplary |
| `gitmap/movemerge` | 16 | 0 | 0 | 4 | 99.5 / 100 | 🟢 Exemplary |
| `gitmap/osuser` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/osutil` | 3 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/probe` | 5 | 0 | 0 | 3 | 98.8 / 100 | 🟢 Exemplary |
| `gitmap/release` | 59 | 0 | 0 | 33 | 98.9 / 100 | 🟢 Exemplary |
| `gitmap/render` | 11 | 0 | 0 | 6 | 98.9 / 100 | 🟢 Exemplary |
| `gitmap/repodb` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/result` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/scanner` | 8 | 1 | 0 | 2 | 98.2 / 100 | 🟢 Exemplary |
| `gitmap/scripts` | 8 | 4 | 0 | 2 | 94.5 / 100 | 🟡 Acceptable |
| `gitmap/searcher` | 4 | 0 | 0 | 2 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/setup` | 6 | 0 | 0 | 2 | 99.3 / 100 | 🟢 Exemplary |
| `gitmap/stablejson` | 3 | 0 | 0 | 3 | 98.0 / 100 | 🟢 Exemplary |
| `gitmap/startup` | 34 | 1 | 0 | 21 | 98.5 / 100 | 🟢 Exemplary |
| `gitmap/store` | 108 | 2 | 0 | 37 | 99.1 / 100 | 🟢 Exemplary |
| `gitmap/templates` | 15 | 0 | 0 | 7 | 99.1 / 100 | 🟢 Exemplary |
| `gitmap/tests` | 33 | 5 | 0 | 16 | 97.5 / 100 | 🟢 Exemplary |
| `gitmap/theme` | 4 | 0 | 0 | 2 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/transport` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/tui` | 20 | 0 | 0 | 10 | 99.0 / 100 | 🟢 Exemplary |
| `gitmap/txn` | 6 | 0 | 0 | 5 | 98.3 / 100 | 🟢 Exemplary |
| `gitmap/uipref` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/verbose` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/visibility` | 7 | 0 | 0 | 3 | 99.1 / 100 | 🟢 Exemplary |
| `gitmap/vscodepm` | 17 | 0 | 0 | 7 | 99.2 / 100 | 🟢 Exemplary |
| `gitmap/vscodeworkspace` | 3 | 0 | 0 | 1 | 99.3 / 100 | 🟢 Exemplary |
| `gitmap/worker` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `gitmap/workspacesync` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `init.ps1` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `init.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `install-quick.ps1` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `install-quick.sh` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `install.ps1` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `install.sh` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/allowlist-forbidden-string.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/check-axios-version.sh` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-boolean-guidelines.py` | 1 | 0 | 1 | 1 | 93.0 / 100 | 🟡 Acceptable |
| `linter-scripts/check-enum-and-boolean.py` | 1 | 0 | 2 | 1 | 88.0 / 100 | 🟡 Acceptable |
| `linter-scripts/check-error-management.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-file-sizes.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-forbidden-spec-paths.sh` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-forbidden-strings.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-function-lengths.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-markdown-headings.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-memory-mirror-drift.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-mws-error-codes.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-nested-ifs.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-newline-styling.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-placeholder-comments.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/check-prompts-loaded.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-readme-canonicals.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-readme-install-section.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/check-relative-paths.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-root-readme.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-runner-dispatch-antipatterns.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-spec-cross-links.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-spec-folder-refs.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/check-tunable-constants.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/forbidden-strings-summary.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/run.ps1` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/run.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `linter-scripts/suggest-spec-cross-link-fixes.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/tests` | 33 | 6 | 0 | 22 | 96.8 / 100 | 🟢 Exemplary |
| `linter-scripts/validate-guidelines.go` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/validate-guidelines.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `linter-scripts/validate-rename-intake.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `patch.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `postcss.config.js` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `remotion-demo/index.ts` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `remotion-demo/src` | 10 | 0 | 1 | 1 | 99.3 / 100 | 🟢 Exemplary |
| `run.ps1` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `run.sh` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `scripts/audit_codebase.py` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `scripts/build-stamp.ps1` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `scripts/build-stamp.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `scripts/changelog` | 22 | 0 | 0 | 3 | 99.7 / 100 | 🟢 Exemplary |
| `scripts/fix-repo` | 14 | 0 | 0 | 2 | 99.7 / 100 | 🟢 Exemplary |
| `scripts/format-go.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `scripts/generate_cg_audit_report.py` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `scripts/generate_plan_files.py` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `scripts/kubernetes` | 39 | 0 | 0 | 5 | 99.7 / 100 | 🟢 Exemplary |
| `scripts/misspell-local.ps1` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `scripts/misspell-local.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `scripts/preflight-ci.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `scripts/smoke-act.sh` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `scripts/visibility-change` | 4 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `setup.sh` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `spec/11-powershell-integration` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `src/App.tsx` | 1 | 0 | 1 | 0 | 95.0 / 100 | 🟢 Exemplary |
| `src/components` | 89 | 0 | 25 | 8 | 98.4 / 100 | 🟢 Exemplary |
| `src/constants` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `src/data` | 4 | 2 | 0 | 46 | 72.0 / 100 | 🟠 Needs Work |
| `src/hooks` | 3 | 0 | 0 | 2 | 98.7 / 100 | 🟢 Exemplary |
| `src/lib` | 5 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `src/main.tsx` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `src/pages` | 80 | 0 | 68 | 18 | 95.3 / 100 | 🟢 Exemplary |
| `src/test` | 8 | 0 | 3 | 2 | 97.6 / 100 | 🟢 Exemplary |
| `src/types` | 2 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `src/vite-env.d.ts` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `tailwind.config.ts` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `uninstall-quick.ps1` | 1 | 1 | 0 | 0 | 90.0 / 100 | 🟡 Acceptable |
| `uninstall-quick.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `visibility-change.ps1` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `visibility-change.sh` | 1 | 0 | 0 | 1 | 98.0 / 100 | 🟢 Exemplary |
| `visibility-change/Apply.ps1` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `visibility-change/Provider.ps1` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `visibility-change/apply.sh` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `visibility-change/provider.sh` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `vite.config.ts` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |
| `vitest.config.ts` | 1 | 0 | 0 | 0 | 100.0 / 100 | 🟢 Exemplary |

---

## Detailed Violation Ledger (Drop by Drop)

| Id | File Path | Line | Function / Component | Rule Code | Exact Snippet | Severity | Planned Remediation |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| V-001 | `fix-repo.ps1` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-002 | `fix-repo.sh` | 1 | `-` | `CG-SIZE-001` | 200 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-003 | `init.sh` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-004 | `install-quick.ps1` | 1 | `-` | `CG-SIZE-001` | 340 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-005 | `install-quick.sh` | 1 | `-` | `CG-SIZE-001` | 375 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-006 | `install.ps1` | 1 | `-` | `CG-SIZE-001` | 1439 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-007 | `install.sh` | 1 | `-` | `CG-SIZE-001` | 1627 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-008 | `run.ps1` | 1 | `-` | `CG-SIZE-001` | 1861 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-009 | `run.sh` | 1 | `-` | `CG-SIZE-001` | 1210 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-010 | `tailwind.config.ts` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-011 | `uninstall-quick.ps1` | 1 | `-` | `CG-SIZE-001` | 370 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-012 | `uninstall-quick.sh` | 1 | `-` | `CG-SIZE-001` | 280 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-013 | `visibility-change.ps1` | 1 | `-` | `CG-SIZE-001` | 169 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-014 | `visibility-change.sh` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-015 | `gitmap/apperror/apperror.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-016 | `gitmap/archive/archive.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-017 | `gitmap/archive/archive_test.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-018 | `gitmap/archive/create.go` | 1 | `-` | `CG-SIZE-001` | 290 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-019 | `gitmap/archive/extract.go` | 1 | `-` | `CG-SIZE-001` | 348 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-020 | `gitmap/archive/source.go` | 1 | `-` | `CG-SIZE-001` | 277 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-021 | `gitmap/cliexit/cliexit.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-022 | `gitmap/cliexit/kind.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-023 | `gitmap/cliexit/report.go` | 1 | `-` | `CG-SIZE-001` | 201 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-024 | `gitmap/cliexit/report_test.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-025 | `gitmap/clonefrom/depthflag_format_test.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-026 | `gitmap/clonefrom/execute.go` | 1 | `-` | `CG-SIZE-001` | 239 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-027 | `gitmap/clonefrom/execute_checkout_test.go` | 1 | `-` | `CG-SIZE-001` | 216 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-028 | `gitmap/clonefrom/execute_concurrent.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-029 | `gitmap/clonefrom/execute_concurrent_test.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-030 | `gitmap/clonefrom/execute_test.go` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-031 | `gitmap/clonefrom/jsonschema.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-032 | `gitmap/clonefrom/jsonschema_helpers.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-033 | `gitmap/clonefrom/jsonschema_test.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-034 | `gitmap/clonefrom/parse.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-035 | `gitmap/clonefrom/parsecsv.go` | 1 | `-` | `CG-SIZE-001` | 151 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-036 | `gitmap/clonefrom/parsecsv_columnerr_test.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-037 | `gitmap/clonefrom/parse_test.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-038 | `gitmap/clonefrom/render.go` | 1 | `-` | `CG-SIZE-001` | 205 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-039 | `gitmap/clonefrom/result_schema_drift_test.go` | 1 | `-` | `CG-SIZE-001` | 220 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-040 | `gitmap/clonefrom/summary.go` | 1 | `-` | `CG-SIZE-001` | 151 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-041 | `gitmap/clonefrom/summary_csvquoting_golden_test.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-042 | `gitmap/clonefrom/summary_golden_test.go` | 1 | `-` | `CG-SIZE-001` | 174 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-043 | `gitmap/clonefrom/summary_json.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-044 | `gitmap/clonefrom/summary_provenance_test.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-045 | `gitmap/clonefrom/summary_terminal.go` | 1 | `-` | `CG-SIZE-001` | 164 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-046 | `gitmap/clonefrom/validate.go` | 1 | `-` | `CG-SIZE-001` | 187 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-047 | `gitmap/clonenext/batch.go` | 1 | `-` | `CG-SIZE-001` | 234 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-048 | `gitmap/clonenext/batch_test.go` | 1 | `-` | `CG-SIZE-001` | 208 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-049 | `gitmap/clonenext/github.go` | 1 | `-` | `CG-SIZE-001` | 151 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-050 | `gitmap/clonenext/localstate.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-051 | `gitmap/clonenext/localstate_test.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-052 | `gitmap/clonenext/remoteupdate_test.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-053 | `gitmap/clonenext/version_test.go` | 1 | `-` | `CG-SIZE-001` | 101 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-054 | `gitmap/clonenow/clonenow.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-055 | `gitmap/clonenow/crossformat_golden_test.go` | 1 | `-` | `CG-SIZE-001` | 205 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-056 | `gitmap/clonenow/execute.go` | 1 | `-` | `CG-SIZE-001` | 193 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-057 | `gitmap/clonenow/execute_concurrent.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-058 | `gitmap/clonenow/execute_concurrent_test.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-059 | `gitmap/clonenow/execute_idempotent.go` | 1 | `-` | `CG-SIZE-001` | 323 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-060 | `gitmap/clonenow/execute_idempotent_test.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-061 | `gitmap/clonenow/execute_test.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-062 | `gitmap/clonenow/parse.go` | 1 | `-` | `CG-SIZE-001` | 290 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-063 | `gitmap/clonenow/parsetext.go` | 1 | `-` | `CG-SIZE-001` | 155 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-064 | `gitmap/clonenow/parse_schema.go` | 1 | `-` | `CG-SIZE-001` | 187 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-065 | `gitmap/clonenow/parse_schema_json.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-066 | `gitmap/clonenow/parse_schema_test.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-067 | `gitmap/clonenow/parse_test.go` | 1 | `-` | `CG-SIZE-001` | 253 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-068 | `gitmap/clonenow/render.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-069 | `gitmap/clonepick/clonepick.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-070 | `gitmap/clonepick/parse.go` | 1 | `-` | `CG-SIZE-001` | 237 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-071 | `gitmap/clonepick/picker.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-072 | `gitmap/clonepick/picker_test.go` | 1 | `-` | `CG-SIZE-001` | 130 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-073 | `gitmap/clonepick/picker_tree.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-074 | `gitmap/clonepick/picker_view.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-075 | `gitmap/clonepick/picker_window_test.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-076 | `gitmap/clonepick/promote.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-077 | `gitmap/clonepick/replay_test.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-078 | `gitmap/clonepick/sparse.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-079 | `gitmap/cloner/audit.go` | 1 | `-` | `CG-SIZE-001` | 208 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-080 | `gitmap/cloner/audit_test.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-081 | `gitmap/cloner/batchprogress.go` | 1 | `-` | `CG-SIZE-001` | 222 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-082 | `gitmap/cloner/batchprogress_test.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-083 | `gitmap/cloner/cache.go` | 1 | `-` | `CG-SIZE-001` | 226 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-084 | `gitmap/cloner/cloner.go` | 1 | `-` | `CG-SIZE-001` | 224 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-085 | `gitmap/cloner/concurrent.go` | 1 | `-` | `CG-SIZE-001` | 120 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-086 | `gitmap/cloner/concurrent_test.go` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-087 | `gitmap/cloner/progress.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-088 | `gitmap/cloner/pulldiag.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-089 | `gitmap/cloner/runners.go` | 1 | `-` | `CG-SIZE-001` | 162 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-090 | `gitmap/cloner/safe_pull.go` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-091 | `gitmap/cloner/safe_push.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-092 | `gitmap/cloner/strategy.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-093 | `gitmap/cloner/strategy_fallback_test.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-094 | `gitmap/cluster/dispatcher.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-095 | `gitmap/cluster/exec_lifecycle.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-096 | `gitmap/cluster/exec_proj.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-097 | `gitmap/cluster/exec_ps_test.go` | 1 | `-` | `CG-SIZE-001` | 151 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-098 | `gitmap/cluster/node_resolver.go` | 1 | `-` | `CG-SIZE-001` | 211 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-099 | `gitmap/cluster/node_resolver_test.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-100 | `gitmap/cluster/pool.go` | 1 | `-` | `CG-SIZE-001` | 243 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-101 | `gitmap/cmd/addignoreattrs.go` | 1 | `-` | `CG-SIZE-001` | 330 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-102 | `gitmap/cmd/addlfsinstall.go` | 1 | `-` | `CG-SIZE-001` | 197 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-103 | `gitmap/cmd/agy_cmd.go` | 1 | `-` | `CG-SIZE-001` | 396 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-104 | `gitmap/cmd/aliasops.go` | 1 | `-` | `CG-SIZE-001` | 183 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-105 | `gitmap/cmd/aliasresolve.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-106 | `gitmap/cmd/amend.go` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-107 | `gitmap/cmd/amendaudit.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-108 | `gitmap/cmd/amendauditrender.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-109 | `gitmap/cmd/amendexec.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-110 | `gitmap/cmd/amendexecprint.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-111 | `gitmap/cmd/amendlist.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-112 | `gitmap/cmd/as.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-113 | `gitmap/cmd/asops.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-114 | `gitmap/cmd/audit.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-115 | `gitmap/cmd/auditlegacy.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-116 | `gitmap/cmd/auditlegacy_diffs.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-117 | `gitmap/cmd/auditlegacy_report.go` | 1 | `-` | `CG-SIZE-001` | 210 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-118 | `gitmap/cmd/backup.go` | 1 | `-` | `CG-SIZE-001` | 289 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-119 | `gitmap/cmd/binarylocations.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-120 | `gitmap/cmd/cdops.go` | 1 | `-` | `CG-SIZE-001` | 216 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-121 | `gitmap/cmd/cfrppriorversion.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-122 | `gitmap/cmd/cg.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-123 | `gitmap/cmd/cg_version.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-124 | `gitmap/cmd/cg_worker.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-125 | `gitmap/cmd/changelog.go` | 1 | `-` | `CG-SIZE-001` | 205 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-126 | `gitmap/cmd/changelogprint.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-127 | `gitmap/cmd/changelogwrap.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-128 | `gitmap/cmd/changelog_regen.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-129 | `gitmap/cmd/chromeprofile.go` | 1 | `-` | `CG-SIZE-001` | 333 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-130 | `gitmap/cmd/chromeprofile_copy.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-131 | `gitmap/cmd/chromeprofile_copy_test.go` | 1 | `-` | `CG-SIZE-001` | 158 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-132 | `gitmap/cmd/chromeprofile_csv.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-133 | `gitmap/cmd/chromeprofile_csv_test.go` | 1 | `-` | `CG-SIZE-001` | 163 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-134 | `gitmap/cmd/chromeprofile_export.go` | 1 | `-` | `CG-SIZE-001` | 121 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-135 | `gitmap/cmd/chromeprofile_merge.go` | 1 | `-` | `CG-SIZE-001` | 458 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-136 | `gitmap/cmd/chromeprofile_paths.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-137 | `gitmap/cmd/chromeprofile_register.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-138 | `gitmap/cmd/chromeprofile_resolve.go` | 1 | `-` | `CG-SIZE-001` | 163 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-139 | `gitmap/cmd/chromeprofile_resolve_test.go` | 1 | `-` | `CG-SIZE-001` | 148 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-140 | `gitmap/cmd/chromeprofile_zip_export.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-141 | `gitmap/cmd/chrome_backup.go` | 1 | `-` | `CG-SIZE-001` | 284 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-142 | `gitmap/cmd/chrome_backup_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-143 | `gitmap/cmd/chrome_bookmarks.go` | 1 | `-` | `CG-SIZE-001` | 288 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-144 | `gitmap/cmd/chrome_diff.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-145 | `gitmap/cmd/chrome_manifest.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-146 | `gitmap/cmd/chrome_manifest_test.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-147 | `gitmap/cmd/cliexit_clone_test.go` | 1 | `-` | `CG-SIZE-001` | 165 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-148 | `gitmap/cmd/cliexit_context_test.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-149 | `gitmap/cmd/cliexit_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 222 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-150 | `gitmap/cmd/clone.go` | 1 | `-` | `CG-SIZE-001` | 614 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-151 | `gitmap/cmd/clonefixrepo.go` | 1 | `-` | `CG-SIZE-001` | 399 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-152 | `gitmap/cmd/clonefixrepofoldertransport.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-153 | `gitmap/cmd/clonefixrepoparallel.go` | 1 | `-` | `CG-SIZE-001` | 281 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-154 | `gitmap/cmd/clonefrom.go` | 1 | `-` | `CG-SIZE-001` | 219 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-155 | `gitmap/cmd/clonefrom_reports.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-156 | `gitmap/cmd/clonemulti.go` | 1 | `-` | `CG-SIZE-001` | 230 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-157 | `gitmap/cmd/clonenext.go` | 1 | `-` | `CG-SIZE-001` | 432 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-158 | `gitmap/cmd/clonenextbatch.go` | 1 | `-` | `CG-SIZE-001` | 189 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-159 | `gitmap/cmd/clonenextbatchconcurrent.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-160 | `gitmap/cmd/clonenextbatchconcurrent_e2e_csv_test.go` | 1 | `-` | `CG-SIZE-001` | 187 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-161 | `gitmap/cmd/clonenextbatchconcurrent_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 180 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-162 | `gitmap/cmd/clonenextflags.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-163 | `gitmap/cmd/clonenextfolderdispatch.go` | 1 | `-` | `CG-SIZE-001` | 189 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-164 | `gitmap/cmd/clonenextfolderdispatch_test.go` | 1 | `-` | `CG-SIZE-001` | 150 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-165 | `gitmap/cmd/clonenow.go` | 1 | `-` | `CG-SIZE-001` | 283 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-166 | `gitmap/cmd/clonepick.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-167 | `gitmap/cmd/clonepick_execute.go` | 1 | `-` | `CG-SIZE-001` | 111 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-168 | `gitmap/cmd/clonepick_flags.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-169 | `gitmap/cmd/clonepmsync.go` | 1 | `-` | `CG-SIZE-001` | 165 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-170 | `gitmap/cmd/clonepmsync_canonicalize_test.go` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-171 | `gitmap/cmd/clonepmsync_dedup_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-172 | `gitmap/cmd/clonepretty.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-173 | `gitmap/cmd/cloneprintargv.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-174 | `gitmap/cmd/clonereplace.go` | 1 | `-` | `CG-SIZE-001` | 167 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-175 | `gitmap/cmd/clonestream_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 228 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-176 | `gitmap/cmd/clonetermblock_golden_test.go` | 1 | `-` | `CG-SIZE-001` | 315 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-177 | `gitmap/cmd/clonetermplan.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-178 | `gitmap/cmd/clonetermrow.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-179 | `gitmap/cmd/clonetermstream.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-180 | `gitmap/cmd/clonetermverify.go` | 1 | `-` | `CG-SIZE-001` | 184 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-181 | `gitmap/cmd/clonetermverifyexit_test.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-182 | `gitmap/cmd/clonetermverifystate.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-183 | `gitmap/cmd/clonetermverify_test.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-184 | `gitmap/cmd/cloneurlconvert.go` | 1 | `-` | `CG-SIZE-001` | 169 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-185 | `gitmap/cmd/clonevscode.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-186 | `gitmap/cmd/clone_idempotent_test.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-187 | `gitmap/cmd/clone_test.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-188 | `gitmap/cmd/clustercommand.go` | 1 | `-` | `CG-SIZE-001` | 206 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-189 | `gitmap/cmd/clusterflags.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-190 | `gitmap/cmd/clustersubcmd.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-191 | `gitmap/cmd/cluster_ops.go` | 1 | `-` | `CG-SIZE-001` | 508 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-192 | `gitmap/cmd/cmdconstants_unique_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 202 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-193 | `gitmap/cmd/code.go` | 1 | `-` | `CG-SIZE-001` | 465 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-194 | `gitmap/cmd/codingguidelines.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-195 | `gitmap/cmd/codingguidelines_commit.go` | 1 | `-` | `CG-SIZE-001` | 162 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-196 | `gitmap/cmd/codingguidelines_test.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-197 | `gitmap/cmd/committransfer.go` | 1 | `-` | `CG-SIZE-001` | 244 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-198 | `gitmap/cmd/commit_push.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-199 | `gitmap/cmd/completion.go` | 1 | `-` | `CG-SIZE-001` | 167 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-200 | `gitmap/cmd/csvcrlf_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-201 | `gitmap/cmd/dashboard.go` | 1 | `-` | `CG-SIZE-001` | 124 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-202 | `gitmap/cmd/dedupe.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-203 | `gitmap/cmd/desktopsync.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-204 | `gitmap/cmd/diffprofiles.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-205 | `gitmap/cmd/diffprofilesops.go` | 1 | `-` | `CG-SIZE-001` | 150 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-206 | `gitmap/cmd/diffprofilesrender.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-207 | `gitmap/cmd/doctor.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-208 | `gitmap/cmd/doctorchecks.go` | 1 | `-` | `CG-SIZE-001` | 312 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-209 | `gitmap/cmd/doctordupbin.go` | 1 | `-` | `CG-SIZE-001` | 167 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-210 | `gitmap/cmd/doctorfixpath.go` | 1 | `-` | `CG-SIZE-001` | 193 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-211 | `gitmap/cmd/doctorsync.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-212 | `gitmap/cmd/doctorvalidate.go` | 1 | `-` | `CG-SIZE-001` | 120 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-213 | `gitmap/cmd/doctorversion.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-214 | `gitmap/cmd/doctor_extra.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-215 | `gitmap/cmd/doctor_fixrepo.go` | 1 | `-` | `CG-SIZE-001` | 267 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-216 | `gitmap/cmd/doctor_run.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-217 | `gitmap/cmd/dopendingretry.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-218 | `gitmap/cmd/downloaderconfig.go` | 1 | `-` | `CG-SIZE-001` | 177 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-219 | `gitmap/cmd/envops.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-220 | `gitmap/cmd/envplatform_unix.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-221 | `gitmap/cmd/envregistry.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-222 | `gitmap/cmd/escapecwd_test.go` | 1 | `-` | `CG-SIZE-001` | 109 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-223 | `gitmap/cmd/exec.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-224 | `gitmap/cmd/exec_pending_test.go` | 1 | `-` | `CG-SIZE-001` | 188 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-225 | `gitmap/cmd/exportrender.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-226 | `gitmap/cmd/fallback.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-227 | `gitmap/cmd/filemanipulator.go` | 1 | `-` | `CG-SIZE-001` | 405 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-228 | `gitmap/cmd/file_search.go` | 1 | `-` | `CG-SIZE-001` | 164 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-229 | `gitmap/cmd/findnextflags.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-230 | `gitmap/cmd/findnextflags_test.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-231 | `gitmap/cmd/find_entry.go` | 1 | `-` | `CG-SIZE-001` | 235 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-232 | `gitmap/cmd/fixauth.go` | 1 | `-` | `CG-SIZE-001` | 202 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-233 | `gitmap/cmd/fixrepo.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-234 | `gitmap/cmd/fixrepo_backup.go` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-235 | `gitmap/cmd/fixrepo_config.go` | 1 | `-` | `CG-SIZE-001` | 177 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-236 | `gitmap/cmd/fixrepo_flags.go` | 1 | `-` | `CG-SIZE-001` | 299 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-237 | `gitmap/cmd/fixrepo_gofmt.go` | 1 | `-` | `CG-SIZE-001` | 221 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-238 | `gitmap/cmd/fixrepo_gofmt_test.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-239 | `gitmap/cmd/fixrepo_identity.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-240 | `gitmap/cmd/fixrepo_rewrite.go` | 1 | `-` | `CG-SIZE-001` | 218 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-241 | `gitmap/cmd/fixrepo_rewrite_scan_test.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-242 | `gitmap/cmd/fixrepo_rewrite_v9tov12_test.go` | 1 | `-` | `CG-SIZE-001` | 285 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-243 | `gitmap/cmd/fixrepo_scan.go` | 1 | `-` | `CG-SIZE-001` | 252 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-244 | `gitmap/cmd/fixrepo_strict_packages_test.go` | 1 | `-` | `CG-SIZE-001` | 124 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-245 | `gitmap/cmd/flags_test.go` | 1 | `-` | `CG-SIZE-001` | 210 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-246 | `gitmap/cmd/gomod.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-247 | `gitmap/cmd/gomodbranch.go` | 1 | `-` | `CG-SIZE-001` | 146 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-248 | `gitmap/cmd/gomodreplace.go` | 1 | `-` | `CG-SIZE-001` | 198 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-249 | `gitmap/cmd/gomod_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 189 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-250 | `gitmap/cmd/gomod_test.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-251 | `gitmap/cmd/group.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-252 | `gitmap/cmd/haschange.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-253 | `gitmap/cmd/hd_terminal.go` | 1 | `-` | `CG-SIZE-001` | 197 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-254 | `gitmap/cmd/helpdashboard.go` | 1 | `-` | `CG-SIZE-001` | 237 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-255 | `gitmap/cmd/hints.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-256 | `gitmap/cmd/history.go` | 1 | `-` | `CG-SIZE-001` | 334 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-257 | `gitmap/cmd/historyrewrite_flags.go` | 1 | `-` | `CG-SIZE-001` | 120 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-258 | `gitmap/cmd/historyrewrite_pin.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-259 | `gitmap/cmd/historyrewrite_pin_test.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-260 | `gitmap/cmd/historyrewrite_sandbox.go` | 1 | `-` | `CG-SIZE-001` | 211 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-261 | `gitmap/cmd/hygiene_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 167 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-262 | `gitmap/cmd/hygiene_parallel.go` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-263 | `gitmap/cmd/inject.go` | 1 | `-` | `CG-SIZE-001` | 154 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-264 | `gitmap/cmd/inject_idempotency.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-265 | `gitmap/cmd/install.go` | 1 | `-` | `CG-SIZE-001` | 225 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-266 | `gitmap/cmd/installagmanager.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-267 | `gitmap/cmd/installctx.go` | 1 | `-` | `CG-SIZE-001` | 212 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-268 | `gitmap/cmd/installctxentries.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-269 | `gitmap/cmd/installctxentries_argv_test.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-270 | `gitmap/cmd/installctxflatten.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-271 | `gitmap/cmd/installctxlinux.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-272 | `gitmap/cmd/installctxlinuxthunar.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-273 | `gitmap/cmd/installctxmac.go` | 1 | `-` | `CG-SIZE-001` | 187 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-274 | `gitmap/cmd/installctxmenu.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-275 | `gitmap/cmd/installctx_argv_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 181 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-276 | `gitmap/cmd/installctx_darwin_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 248 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-277 | `gitmap/cmd/installctx_harness_test.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-278 | `gitmap/cmd/installctx_linux_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 340 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-279 | `gitmap/cmd/installctx_parity_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 215 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-280 | `gitmap/cmd/installctx_windows_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 267 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-281 | `gitmap/cmd/installer_create.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-282 | `gitmap/cmd/installer_create_test.go` | 1 | `-` | `CG-SIZE-001` | 253 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-283 | `gitmap/cmd/installer_export.go` | 1 | `-` | `CG-SIZE-001` | 224 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-284 | `gitmap/cmd/installer_history_tree.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-285 | `gitmap/cmd/installer_history_tree_test.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-286 | `gitmap/cmd/installer_import.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-287 | `gitmap/cmd/installer_install_win.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-288 | `gitmap/cmd/installer_ls.go` | 1 | `-` | `CG-SIZE-001` | 127 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-289 | `gitmap/cmd/installer_reset.go` | 1 | `-` | `CG-SIZE-001` | 121 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-290 | `gitmap/cmd/installer_revert.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-291 | `gitmap/cmd/installer_update.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-292 | `gitmap/cmd/installer_update_win.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-293 | `gitmap/cmd/installlist.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-294 | `gitmap/cmd/installnpp.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-295 | `gitmap/cmd/installnppextract.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-296 | `gitmap/cmd/installobs.go` | 1 | `-` | `CG-SIZE-001` | 317 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-297 | `gitmap/cmd/installscripts.go` | 1 | `-` | `CG-SIZE-001` | 155 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-298 | `gitmap/cmd/installtools.go` | 1 | `-` | `CG-SIZE-001` | 515 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-299 | `gitmap/cmd/installverify.go` | 1 | `-` | `CG-SIZE-001` | 204 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-300 | `gitmap/cmd/installvscode.go` | 1 | `-` | `CG-SIZE-001` | 175 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-301 | `gitmap/cmd/installwt.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-302 | `gitmap/cmd/install_profile_tree.go` | 1 | `-` | `CG-SIZE-001` | 162 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-303 | `gitmap/cmd/install_profile_tree_test.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-304 | `gitmap/cmd/install_unit_test.go` | 1 | `-` | `CG-SIZE-001` | 206 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-305 | `gitmap/cmd/jsoncontract_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-306 | `gitmap/cmd/jsonschema_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-307 | `gitmap/cmd/jsonsnapshot_helpers_failuremsg_test.go` | 1 | `-` | `CG-SIZE-001` | 200 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-308 | `gitmap/cmd/jsonsnapshot_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 141 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-309 | `gitmap/cmd/latestbranch.go` | 1 | `-` | `CG-SIZE-001` | 189 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-310 | `gitmap/cmd/lfscommon.go` | 1 | `-` | `CG-SIZE-001` | 261 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-311 | `gitmap/cmd/list.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-312 | `gitmap/cmd/listreleases.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-313 | `gitmap/cmd/listreleasesload.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-314 | `gitmap/cmd/listreleasesrender.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-315 | `gitmap/cmd/listreleases_jsonschema_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 215 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-316 | `gitmap/cmd/listversions.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-317 | `gitmap/cmd/llmdocs.go` | 1 | `-` | `CG-SIZE-001` | 177 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-318 | `gitmap/cmd/llmdocsgroups.go` | 1 | `-` | `CG-SIZE-001` | 225 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-319 | `gitmap/cmd/llmdocsjson_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-320 | `gitmap/cmd/llmdocsrender.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-321 | `gitmap/cmd/llmdocssections.go` | 1 | `-` | `CG-SIZE-001` | 141 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-322 | `gitmap/cmd/macro_cmd.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-323 | `gitmap/cmd/macro_tree_test.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-324 | `gitmap/cmd/move.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-325 | `gitmap/cmd/multigroup.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-326 | `gitmap/cmd/orphans.go` | 1 | `-` | `CG-SIZE-001` | 181 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-327 | `gitmap/cmd/outputfilenames_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 216 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-328 | `gitmap/cmd/pendingclear.go` | 1 | `-` | `CG-SIZE-001` | 358 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-329 | `gitmap/cmd/pendingtaskhelper.go` | 1 | `-` | `CG-SIZE-001` | 127 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-330 | `gitmap/cmd/pendingtaskhelper_test.go` | 1 | `-` | `CG-SIZE-001` | 260 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-331 | `gitmap/cmd/prettyflag.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-332 | `gitmap/cmd/prettyflag_test.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-333 | `gitmap/cmd/probe.go` | 1 | `-` | `CG-SIZE-001` | 170 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-334 | `gitmap/cmd/probeflags.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-335 | `gitmap/cmd/probereport.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-336 | `gitmap/cmd/profileops.go` | 1 | `-` | `CG-SIZE-001` | 101 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-337 | `gitmap/cmd/pull.go` | 1 | `-` | `CG-SIZE-001` | 516 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-338 | `gitmap/cmd/pullparallel.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-339 | `gitmap/cmd/pullreleasecd.go` | 1 | `-` | `CG-SIZE-001` | 231 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-340 | `gitmap/cmd/pullreleasecd_test.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-341 | `gitmap/cmd/push.go` | 1 | `-` | `CG-SIZE-001` | 302 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-342 | `gitmap/cmd/pushparallel.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-343 | `gitmap/cmd/pushpull_transport_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-344 | `gitmap/cmd/reclonetransport.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-345 | `gitmap/cmd/reclone_autopickup.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-346 | `gitmap/cmd/reclone_confirm.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-347 | `gitmap/cmd/reclone_summary.go` | 1 | `-` | `CG-SIZE-001` | 151 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-348 | `gitmap/cmd/reclone_validate.go` | 1 | `-` | `CG-SIZE-001` | 254 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-349 | `gitmap/cmd/regoldens.go` | 1 | `-` | `CG-SIZE-001` | 192 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-350 | `gitmap/cmd/regoldens_diff.go` | 1 | `-` | `CG-SIZE-001` | 201 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-351 | `gitmap/cmd/regoldens_diff_test.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-352 | `gitmap/cmd/regoldens_exec.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-353 | `gitmap/cmd/reinstall.go` | 1 | `-` | `CG-SIZE-001` | 183 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-354 | `gitmap/cmd/release.go` | 1 | `-` | `CG-SIZE-001` | 236 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-355 | `gitmap/cmd/releasealias.go` | 1 | `-` | `CG-SIZE-001` | 120 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-356 | `gitmap/cmd/releasealias_git.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-357 | `gitmap/cmd/releasepull.go` | 1 | `-` | `CG-SIZE-001` | 221 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-358 | `gitmap/cmd/releaserebase.go` | 1 | `-` | `CG-SIZE-001` | 183 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-359 | `gitmap/cmd/releasescan.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-360 | `gitmap/cmd/releaseundo.go` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-361 | `gitmap/cmd/release_notes_opts.go` | 1 | `-` | `CG-SIZE-001` | 280 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-362 | `gitmap/cmd/release_scan_commits.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-363 | `gitmap/cmd/replace.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-364 | `gitmap/cmd/replaceapply.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-365 | `gitmap/cmd/replaceaudit_test.go` | 1 | `-` | `CG-SIZE-001` | 124 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-366 | `gitmap/cmd/replaceextfilter_test.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-367 | `gitmap/cmd/replaceflags.go` | 1 | `-` | `CG-SIZE-001` | 190 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-368 | `gitmap/cmd/replaceversionparse_test.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-369 | `gitmap/cmd/replaceversionrun.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-370 | `gitmap/cmd/replacewalk.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-371 | `gitmap/cmd/replacewalk_test.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-372 | `gitmap/cmd/reporeclone.go` | 1 | `-` | `CG-SIZE-001` | 188 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-373 | `gitmap/cmd/reporeclone_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-374 | `gitmap/cmd/reporeclone_test.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-375 | `gitmap/cmd/rescan.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-376 | `gitmap/cmd/rescansubtree.go` | 1 | `-` | `CG-SIZE-001` | 242 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-377 | `gitmap/cmd/rescansubtree_test.go` | 1 | `-` | `CG-SIZE-001` | 207 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-378 | `gitmap/cmd/reset.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-379 | `gitmap/cmd/resolver.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-380 | `gitmap/cmd/revertscript.go` | 1 | `-` | `CG-SIZE-001` | 101 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-381 | `gitmap/cmd/reverttxn.go` | 1 | `-` | `CG-SIZE-001` | 249 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-382 | `gitmap/cmd/reverttxn_lastn.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-383 | `gitmap/cmd/rm.go` | 1 | `-` | `CG-SIZE-001` | 167 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-384 | `gitmap/cmd/root.go` | 1 | `-` | `CG-SIZE-001` | 416 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-385 | `gitmap/cmd/rootcore.go` | 1 | `-` | `CG-SIZE-001` | 180 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-386 | `gitmap/cmd/rootflags.go` | 1 | `-` | `CG-SIZE-001` | 429 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-387 | `gitmap/cmd/rootrelease_pullrelease_aliases_test.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-388 | `gitmap/cmd/roottooling.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-389 | `gitmap/cmd/rootusage.go` | 1 | `-` | `CG-SIZE-001` | 302 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-390 | `gitmap/cmd/rootusagecompact.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-391 | `gitmap/cmd/rootusagefilter.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-392 | `gitmap/cmd/rootusageflags.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-393 | `gitmap/cmd/rootusagefooter.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-394 | `gitmap/cmd/rootusage_groups.go` | 1 | `-` | `CG-SIZE-001` | 239 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-395 | `gitmap/cmd/rootutility.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-396 | `gitmap/cmd/root_url_shortcut_test.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-397 | `gitmap/cmd/safety_snapshot.go` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-398 | `gitmap/cmd/scan.go` | 1 | `-` | `CG-SIZE-001` | 277 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-399 | `gitmap/cmd/scanbackgroundprobe.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-400 | `gitmap/cmd/scanbenchmark.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-401 | `gitmap/cmd/scanimport.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-402 | `gitmap/cmd/scanoutput.go` | 1 | `-` | `CG-SIZE-001` | 204 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-403 | `gitmap/cmd/scanprogress.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-404 | `gitmap/cmd/scanprojectoutput.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-405 | `gitmap/cmd/scanprojects.go` | 1 | `-` | `CG-SIZE-001` | 165 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-406 | `gitmap/cmd/scanprojectsmeta.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-407 | `gitmap/cmd/scanproject_jsonschema_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-408 | `gitmap/cmd/scanresolve.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-409 | `gitmap/cmd/scan_export_clonefrom_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 190 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-410 | `gitmap/cmd/schedule_cmd.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-411 | `gitmap/cmd/schemaregistry_assert_test.go` | 1 | `-` | `CG-SIZE-001` | 187 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-412 | `gitmap/cmd/schemaregistry_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 203 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-413 | `gitmap/cmd/schemaregistry_test.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-414 | `gitmap/cmd/search_entry.go` | 1 | `-` | `CG-SIZE-001` | 196 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-415 | `gitmap/cmd/selfinstall.go` | 1 | `-` | `CG-SIZE-001` | 405 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-416 | `gitmap/cmd/selfuninstall.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-417 | `gitmap/cmd/selfuninstallhandoff.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-418 | `gitmap/cmd/selfuninstallparts.go` | 1 | `-` | `CG-SIZE-001` | 401 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-419 | `gitmap/cmd/selfupdate.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-420 | `gitmap/cmd/seowrite.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-421 | `gitmap/cmd/seowritegit.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-422 | `gitmap/cmd/seowriteloop.go` | 1 | `-` | `CG-SIZE-001` | 214 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-423 | `gitmap/cmd/seowritetemplate.go` | 1 | `-` | `CG-SIZE-001` | 202 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-424 | `gitmap/cmd/setup.go` | 1 | `-` | `CG-SIZE-001` | 202 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-425 | `gitmap/cmd/setupconfig.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-426 | `gitmap/cmd/setup_ubuntu.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-427 | `gitmap/cmd/sf.go` | 1 | `-` | `CG-SIZE-001` | 201 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-428 | `gitmap/cmd/size.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-429 | `gitmap/cmd/sshcat.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-430 | `gitmap/cmd/sshconfig.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-431 | `gitmap/cmd/sshexec.go` | 1 | `-` | `CG-SIZE-001` | 252 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-432 | `gitmap/cmd/sshexisting.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-433 | `gitmap/cmd/sshgen.go` | 1 | `-` | `CG-SIZE-001` | 295 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-434 | `gitmap/cmd/sshjoin.go` | 1 | `-` | `CG-SIZE-001` | 220 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-435 | `gitmap/cmd/sshjoin_auth_cmd.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-436 | `gitmap/cmd/sshjoin_auth_cmd_test.go` | 1 | `-` | `CG-SIZE-001` | 124 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-437 | `gitmap/cmd/sshstatus.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-438 | `gitmap/cmd/ssh_client_test.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-439 | `gitmap/cmd/stale.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-440 | `gitmap/cmd/startup.go` | 1 | `-` | `CG-SIZE-001` | 281 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-441 | `gitmap/cmd/startupadd.go` | 1 | `-` | `CG-SIZE-001` | 157 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-442 | `gitmap/cmd/startuplistcsv_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-443 | `gitmap/cmd/startuplistfilter.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-444 | `gitmap/cmd/startuplistfilter_test.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-445 | `gitmap/cmd/startuplistjsonl_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 191 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-446 | `gitmap/cmd/startuplistjson_determinism_test.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-447 | `gitmap/cmd/startuplistjson_indent_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 191 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-448 | `gitmap/cmd/startuplistjson_snapshot_test.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-449 | `gitmap/cmd/startuplistrender.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-450 | `gitmap/cmd/startuplisttable_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 166 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-451 | `gitmap/cmd/startuplist_jsonschema_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-452 | `gitmap/cmd/startupstatusjson.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-453 | `gitmap/cmd/startupstatusjson_test.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-454 | `gitmap/cmd/stats.go` | 1 | `-` | `CG-SIZE-001` | 148 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-455 | `gitmap/cmd/status.go` | 1 | `-` | `CG-SIZE-001` | 256 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-456 | `gitmap/cmd/statusformat.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-457 | `gitmap/cmd/statusprint.go` | 1 | `-` | `CG-SIZE-001` | 191 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-458 | `gitmap/cmd/sync.go` | 1 | `-` | `CG-SIZE-001` | 447 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-459 | `gitmap/cmd/taskfilter.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-460 | `gitmap/cmd/taskops.go` | 1 | `-` | `CG-SIZE-001` | 297 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-461 | `gitmap/cmd/tasksync.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-462 | `gitmap/cmd/task_unit_test.go` | 1 | `-` | `CG-SIZE-001` | 109 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-463 | `gitmap/cmd/templatescli.go` | 1 | `-` | `CG-SIZE-001` | 334 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-464 | `gitmap/cmd/templatescli_filter_test.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-465 | `gitmap/cmd/templatesdiff.go` | 1 | `-` | `CG-SIZE-001` | 180 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-466 | `gitmap/cmd/templatesinit.go` | 1 | `-` | `CG-SIZE-001` | 357 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-467 | `gitmap/cmd/templatesinit_test.go` | 1 | `-` | `CG-SIZE-001` | 272 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-468 | `gitmap/cmd/temprelease.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-469 | `gitmap/cmd/tempreleaseops.go` | 1 | `-` | `CG-SIZE-001` | 212 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-470 | `gitmap/cmd/tempreleaseremove.go` | 1 | `-` | `CG-SIZE-001` | 230 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-471 | `gitmap/cmd/undo.go` | 1 | `-` | `CG-SIZE-001` | 236 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-472 | `gitmap/cmd/uninstall.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-473 | `gitmap/cmd/unzipcompact.go` | 1 | `-` | `CG-SIZE-001` | 217 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-474 | `gitmap/cmd/update.go` | 1 | `-` | `CG-SIZE-001` | 372 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-475 | `gitmap/cmd/updatecleanup_extra.go` | 1 | `-` | `CG-SIZE-001` | 208 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-476 | `gitmap/cmd/updatecleanup_paths.go` | 1 | `-` | `CG-SIZE-001` | 233 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-477 | `gitmap/cmd/updatecleanup_remove.go` | 1 | `-` | `CG-SIZE-001` | 162 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-478 | `gitmap/cmd/updatedebugwindows.go` | 1 | `-` | `CG-SIZE-001` | 158 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-479 | `gitmap/cmd/updatedebugwindows_json.go` | 1 | `-` | `CG-SIZE-001` | 164 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-480 | `gitmap/cmd/updatedebugwindows_plan.go` | 1 | `-` | `CG-SIZE-001` | 200 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-481 | `gitmap/cmd/updatedebugwindows_source_test.go` | 1 | `-` | `CG-SIZE-001` | 190 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-482 | `gitmap/cmd/updatehandofflog.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-483 | `gitmap/cmd/updatehandoff_phase3.go` | 1 | `-` | `CG-SIZE-001` | 235 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-484 | `gitmap/cmd/updateprobe.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-485 | `gitmap/cmd/updateprobe_test.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-486 | `gitmap/cmd/updateremoteinstall.go` | 1 | `-` | `CG-SIZE-001` | 190 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-487 | `gitmap/cmd/updaterepo.go` | 1 | `-` | `CG-SIZE-001` | 212 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-488 | `gitmap/cmd/updatereport.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-489 | `gitmap/cmd/updatescript.go` | 1 | `-` | `CG-SIZE-001` | 257 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-490 | `gitmap/cmd/versionhistory.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-491 | `gitmap/cmd/visibility.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-492 | `gitmap/cmd/visibilityallbulk.go` | 1 | `-` | `CG-SIZE-001` | 284 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-493 | `gitmap/cmd/visibilityallbulkaudit.go` | 1 | `-` | `CG-SIZE-001` | 205 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-494 | `gitmap/cmd/visibilityapply.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-495 | `gitmap/cmd/visibilityapplyone.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-496 | `gitmap/cmd/visibilitybulk.go` | 1 | `-` | `CG-SIZE-001` | 192 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-497 | `gitmap/cmd/visibilitybulkprompt.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-498 | `gitmap/cmd/visibilityexceptlatest.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-499 | `gitmap/cmd/visibilitymakelast.go` | 1 | `-` | `CG-SIZE-001` | 178 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-500 | `gitmap/cmd/visibilityownerlistcache.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-501 | `gitmap/cmd/visibilityresolve.go` | 1 | `-` | `CG-SIZE-001` | 240 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-502 | `gitmap/cmd/visibilityresolveowner.go` | 1 | `-` | `CG-SIZE-001` | 204 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-503 | `gitmap/cmd/visibilityundo.go` | 1 | `-` | `CG-SIZE-001` | 272 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-504 | `gitmap/cmd/vscodecustomtags.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-505 | `gitmap/cmd/vscodepmsync.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-506 | `gitmap/cmd/vscodepmsync_dedupe_test.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-507 | `gitmap/cmd/vscodepmsync_flags.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-508 | `gitmap/cmd/vscodepmsync_flags_test.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-509 | `gitmap/cmd/vscodepmsync_mode_test.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-510 | `gitmap/cmd/vscodepmsync_pathtag_test.go` | 1 | `-` | `CG-SIZE-001` | 146 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-511 | `gitmap/cmd/vscodepmsync_test.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-512 | `gitmap/cmd/vscodepmsync_testhelper_test.go` | 1 | `-` | `CG-SIZE-001` | 170 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-513 | `gitmap/cmd/vscodeworkspace.go` | 1 | `-` | `CG-SIZE-001` | 215 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-514 | `gitmap/cmd/vscode_cmd.go` | 1 | `-` | `CG-SIZE-001` | 204 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-515 | `gitmap/cmd/watchformat.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-516 | `gitmap/cmd/watchops.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-517 | `gitmap/cmd/watchrender.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-518 | `gitmap/cmd/whoami.go` | 1 | `-` | `CG-SIZE-001` | 243 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-519 | `gitmap/cmd/workflow_open_pr.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-520 | `gitmap/cmd/zip.go` | 1 | `-` | `CG-SIZE-001` | 251 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-521 | `gitmap/cmd/zipgroupcreate.go` | 1 | `-` | `CG-SIZE-001` | 204 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-522 | `gitmap/cmd/zipgroupops.go` | 1 | `-` | `CG-SIZE-001` | 229 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-523 | `gitmap/cmd/zipgroupshow.go` | 1 | `-` | `CG-SIZE-001` | 202 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-524 | `gitmap/cmd/commitin/enums.go` | 1 | `-` | `CG-SIZE-001` | 231 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-525 | `gitmap/cmd/commitin/enums_test.go` | 1 | `-` | `CG-SIZE-001` | 208 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-526 | `gitmap/cmd/commitin/parse.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-527 | `gitmap/cmd/commitin/parse_test.go` | 1 | `-` | `CG-SIZE-001` | 236 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-528 | `gitmap/cmd/commitin/parse_validate.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-529 | `gitmap/cmd/commitin/checkpoint/checkpoint.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-530 | `gitmap/cmd/commitin/e2e/edge_cases_test.go` | 1 | `-` | `CG-SIZE-001` | 109 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-531 | `gitmap/cmd/commitin/e2e/repo.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-532 | `gitmap/cmd/commitin/e2e/sibling_test.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-533 | `gitmap/cmd/commitin/message/message_test.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-534 | `gitmap/cmd/commitin/orchestrator/commit.go` | 1 | `-` | `CG-SIZE-001` | 203 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-535 | `gitmap/cmd/commitin/orchestrator/pipeline.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-536 | `gitmap/cmd/commitin/orchestrator/run.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-537 | `gitmap/cmd/commitin/orchestrator/setup.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-538 | `gitmap/cmd/commitin/profile/json.go` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-539 | `gitmap/cmd/commitin/profile/profile_test.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-540 | `gitmap/cmd/commitin/profile/resolve.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-541 | `gitmap/cmd/commitin/replay/replay.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-542 | `gitmap/cmd/commitin/replay/replay_test.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-543 | `gitmap/cmd/commitin/replay/runner.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-544 | `gitmap/cmd/commitin/runlog/runlog_test.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-545 | `gitmap/cmd/commitin/runlog/tagreplay.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-546 | `gitmap/cmd/commitin/runlog/tagreplay_test.go` | 1 | `-` | `CG-SIZE-001` | 282 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-547 | `gitmap/cmd/commitin/walk/hydrate.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-548 | `gitmap/cmd/commitin/workspace/expand.go` | 1 | `-` | `CG-SIZE-001` | 170 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-549 | `gitmap/cmd/commitin/workspace/source.go` | 1 | `-` | `CG-SIZE-001` | 146 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-550 | `gitmap/cmd/commitin/workspace/workspace_test.go` | 1 | `-` | `CG-SIZE-001` | 240 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-551 | `gitmap/cmd/folder/folder.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-552 | `gitmap/cmd/gitmap-node-join/main.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-553 | `gitmap/cmd/gitrm/gitrm.go` | 1 | `-` | `CG-SIZE-001` | 155 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-554 | `gitmap/committransfer/count_parity_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-555 | `gitmap/committransfer/git.go` | 1 | `-` | `CG-SIZE-001` | 194 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-556 | `gitmap/committransfer/interleave.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-557 | `gitmap/committransfer/log.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-558 | `gitmap/committransfer/message.go` | 1 | `-` | `CG-SIZE-001` | 200 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-559 | `gitmap/committransfer/message_test.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-560 | `gitmap/committransfer/plan.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-561 | `gitmap/committransfer/plan_idempotence_test.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-562 | `gitmap/committransfer/replay.go` | 1 | `-` | `CG-SIZE-001` | 206 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-563 | `gitmap/committransfer/types.go` | 1 | `-` | `CG-SIZE-001` | 157 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-564 | `gitmap/completion/allcommands_generated.go` | 1 | `-` | `CG-SIZE-001` | 391 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-565 | `gitmap/completion/bash.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-566 | `gitmap/completion/cdfunction.go` | 1 | `-` | `CG-SIZE-001` | 229 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-567 | `gitmap/completion/cdfunction_test.go` | 1 | `-` | `CG-SIZE-001` | 228 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-568 | `gitmap/completion/completion_test.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-569 | `gitmap/completion/dynamic.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-570 | `gitmap/completion/install.go` | 1 | `-` | `CG-SIZE-001` | 253 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-571 | `gitmap/completion/powershell.go` | 1 | `-` | `CG-SIZE-001` | 229 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-572 | `gitmap/completion/zsh.go` | 1 | `-` | `CG-SIZE-001` | 162 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-573 | `gitmap/completion/internal/gencommands/main.go` | 1 | `-` | `CG-SIZE-001` | 333 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-574 | `gitmap/config/config.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-575 | `gitmap/config/config_test.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-576 | `gitmap/config/validate.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-577 | `gitmap/config/validate_shape.go` | 1 | `-` | `CG-SIZE-001` | 249 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-578 | `gitmap/config/validate_shape_test.go` | 1 | `-` | `CG-SIZE-001` | 210 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-579 | `gitmap/config/validate_test.go` | 1 | `-` | `CG-SIZE-001` | 163 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-580 | `gitmap/constants/cmd_constants_parity_test.go` | 1 | `-` | `CG-SIZE-001` | 254 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-581 | `gitmap/constants/cmd_constants_test.go` | 1 | `-` | `CG-SIZE-001` | 442 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-582 | `gitmap/constants/constants.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-583 | `gitmap/constants/constants_amend.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-584 | `gitmap/constants/constants_archive.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-585 | `gitmap/constants/constants_cd.go` | 1 | `-` | `CG-SIZE-001` | 260 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-586 | `gitmap/constants/constants_chromeprofile.go` | 1 | `-` | `CG-SIZE-001` | 202 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-587 | `gitmap/constants/constants_cli.go` | 1 | `-` | `CG-SIZE-001` | 710 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-588 | `gitmap/constants/constants_clonefrom.go` | 1 | `-` | `CG-SIZE-001` | 216 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-589 | `gitmap/constants/constants_clonenext.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-590 | `gitmap/constants/constants_clonenow.go` | 1 | `-` | `CG-SIZE-001` | 378 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-591 | `gitmap/constants/constants_clonepick.go` | 1 | `-` | `CG-SIZE-001` | 170 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-592 | `gitmap/constants/constants_clone_term.go` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-593 | `gitmap/constants/constants_commitin.go` | 1 | `-` | `CG-SIZE-001` | 227 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-594 | `gitmap/constants/constants_commitin_sql.go` | 1 | `-` | `CG-SIZE-001` | 238 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-595 | `gitmap/constants/constants_committransfer.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-596 | `gitmap/constants/constants_doctor.go` | 1 | `-` | `CG-SIZE-001` | 198 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-597 | `gitmap/constants/constants_find_next.go` | 1 | `-` | `CG-SIZE-001` | 121 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-598 | `gitmap/constants/constants_fixrepo.go` | 1 | `-` | `CG-SIZE-001` | 162 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-599 | `gitmap/constants/constants_git.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-600 | `gitmap/constants/constants_gomod.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-601 | `gitmap/constants/constants_helpgroups.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-602 | `gitmap/constants/constants_install.go` | 1 | `-` | `CG-SIZE-001` | 401 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-603 | `gitmap/constants/constants_messages.go` | 1 | `-` | `CG-SIZE-001` | 390 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-604 | `gitmap/constants/constants_pathsnippet.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-605 | `gitmap/constants/constants_probe.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-606 | `gitmap/constants/constants_project.go` | 1 | `-` | `CG-SIZE-001` | 215 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-607 | `gitmap/constants/constants_project_sql.go` | 1 | `-` | `CG-SIZE-001` | 226 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-608 | `gitmap/constants/constants_release.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-609 | `gitmap/constants/constants_scan_folder.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-610 | `gitmap/constants/constants_selfinstall.go` | 1 | `-` | `CG-SIZE-001` | 227 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-611 | `gitmap/constants/constants_seo.go` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-612 | `gitmap/constants/constants_ssh.go` | 1 | `-` | `CG-SIZE-001` | 186 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-613 | `gitmap/constants/constants_startup.go` | 1 | `-` | `CG-SIZE-001` | 223 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-614 | `gitmap/constants/constants_startup_winregistry.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-615 | `gitmap/constants/constants_store.go` | 1 | `-` | `CG-SIZE-001` | 263 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-616 | `gitmap/constants/constants_temprelease.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-617 | `gitmap/constants/constants_terminal.go` | 1 | `-` | `CG-SIZE-001` | 265 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-618 | `gitmap/constants/constants_transaction.go` | 1 | `-` | `CG-SIZE-001` | 148 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-619 | `gitmap/constants/constants_tui.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-620 | `gitmap/constants/constants_update.go` | 1 | `-` | `CG-SIZE-001` | 534 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-621 | `gitmap/constants/constants_v331.go` | 1 | `-` | `CG-SIZE-001` | 106 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-622 | `gitmap/constants/constants_visibility.go` | 1 | `-` | `CG-SIZE-001` | 250 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-623 | `gitmap/constants/constants_visibility_store_sql.go` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-624 | `gitmap/constants/constants_vscode_pm.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-625 | `gitmap/constants/constants_zipgroup.go` | 1 | `-` | `CG-SIZE-001` | 151 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-626 | `gitmap/dashboard/aggregate.go` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-627 | `gitmap/dashboard/collector.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-628 | `gitmap/dashboard/gitquery.go` | 1 | `-` | `CG-SIZE-001` | 130 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-629 | `gitmap/dashboard/parse.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-630 | `gitmap/dashboard/writer.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-631 | `gitmap/db/clusterexecresult.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-632 | `gitmap/db/clusternode.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-633 | `gitmap/db/clusterrun.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-634 | `gitmap/db/clusterrun_test.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-635 | `gitmap/db/enums.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-636 | `gitmap/detector/csharpparser.go` | 1 | `-` | `CG-SIZE-001` | 191 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-637 | `gitmap/detector/detector.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-638 | `gitmap/detector/goparser.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-639 | `gitmap/detector/rules.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-640 | `gitmap/diff/report.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-641 | `gitmap/diff/tree.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-642 | `gitmap/downloaderconfig/downloaderconfig.go` | 1 | `-` | `CG-SIZE-001` | 169 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-643 | `gitmap/errreport/errreport.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-644 | `gitmap/errreport/errreport_test.go` | 1 | `-` | `CG-SIZE-001` | 111 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-645 | `gitmap/fixtureversion/bump.go` | 1 | `-` | `CG-SIZE-001` | 184 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-646 | `gitmap/fixtureversion/bump_test.go` | 1 | `-` | `CG-SIZE-001` | 164 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-647 | `gitmap/fixtureversion/fixtureversion.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-648 | `gitmap/fixtureversion/hash_test.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-649 | `gitmap/fixtureversion/validate.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-650 | `gitmap/formatter/csv.go` | 1 | `-` | `CG-SIZE-001` | 157 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-651 | `gitmap/formatter/csv_header_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-652 | `gitmap/formatter/desktopscript.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-653 | `gitmap/formatter/formatter_test.go` | 1 | `-` | `CG-SIZE-001` | 111 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-654 | `gitmap/formatter/scangolden_contract_test.go` | 1 | `-` | `CG-SIZE-001` | 197 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-655 | `gitmap/formatter/structure.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-656 | `gitmap/formatter/terminal.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-657 | `gitmap/formatter/terminaltree.go` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-658 | `gitmap/formatter/validate.go` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-659 | `gitmap/formatter/validate_test.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-660 | `gitmap/formatter/validate_writers_test.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-661 | `gitmap/gitutil/gitutil.go` | 1 | `-` | `CG-SIZE-001` | 308 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-662 | `gitmap/gitutil/latestbranch.go` | 1 | `-` | `CG-SIZE-001` | 141 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-663 | `gitmap/gitutil/latestbranchresolve.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-664 | `gitmap/glyphs/install.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-665 | `gitmap/goldenguard/determinism.go` | 1 | `-` | `CG-SIZE-001` | 150 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-666 | `gitmap/goldenguard/determinism_test.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-667 | `gitmap/helptext/coverage_test.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-668 | `gitmap/indexer/walker.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-669 | `gitmap/localdirs/migrate.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-670 | `gitmap/lockcheck/lockcheck_windows.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-671 | `gitmap/lockfile/lockfile.go` | 1 | `-` | `CG-SIZE-001` | 164 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-672 | `gitmap/lockfile/lockfile_test.go` | 1 | `-` | `CG-SIZE-001` | 124 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-673 | `gitmap/mapper/mapper.go` | 1 | `-` | `CG-SIZE-001` | 232 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-674 | `gitmap/mapper/mapper_test.go` | 1 | `-` | `CG-SIZE-001` | 171 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-675 | `gitmap/model/record.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-676 | `gitmap/movemerge/conflict.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-677 | `gitmap/movemerge/integration_test.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-678 | `gitmap/movemerge/merge.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-679 | `gitmap/movemerge/resolve.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-680 | `gitmap/probe/background.go` | 1 | `-` | `CG-SIZE-001` | 249 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-681 | `gitmap/probe/background_test.go` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-682 | `gitmap/probe/probe.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-683 | `gitmap/release/assets.go` | 1 | `-` | `CG-SIZE-001` | 203 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-684 | `gitmap/release/assetsupload.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-685 | `gitmap/release/autocommit.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-686 | `gitmap/release/autocommitgit.go` | 1 | `-` | `CG-SIZE-001` | 189 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-687 | `gitmap/release/autocommit_test.go` | 1 | `-` | `CG-SIZE-001` | 166 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-688 | `gitmap/release/changelog.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-689 | `gitmap/release/changeloggen.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-690 | `gitmap/release/changelogparse.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-691 | `gitmap/release/compress.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-692 | `gitmap/release/gitops.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-693 | `gitmap/release/gitopsquery.go` | 1 | `-` | `CG-SIZE-001` | 142 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-694 | `gitmap/release/metadata.go` | 1 | `-` | `CG-SIZE-001` | 214 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-695 | `gitmap/release/metadata_test.go` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-696 | `gitmap/release/remoteorigin_test.go` | 1 | `-` | `CG-SIZE-001` | 163 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-697 | `gitmap/release/scan_executor.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-698 | `gitmap/release/selfrelease_resolve.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-699 | `gitmap/release/semver.go` | 1 | `-` | `CG-SIZE-001` | 166 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-700 | `gitmap/release/semver_test.go` | 1 | `-` | `CG-SIZE-001` | 223 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-701 | `gitmap/release/temprelease.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-702 | `gitmap/release/workflow.go` | 1 | `-` | `CG-SIZE-001` | 281 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-703 | `gitmap/release/workflowbranch.go` | 1 | `-` | `CG-SIZE-001` | 184 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-704 | `gitmap/release/workflowdocs.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-705 | `gitmap/release/workflowdryrun.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-706 | `gitmap/release/workflowfinalize.go` | 1 | `-` | `CG-SIZE-001` | 241 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-707 | `gitmap/release/workflowgithub.go` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-708 | `gitmap/release/workfloworder_test.go` | 1 | `-` | `CG-SIZE-001` | 190 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-709 | `gitmap/release/workflowpending.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-710 | `gitmap/release/workflowreleasescript.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-711 | `gitmap/release/workflowvalidate.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-712 | `gitmap/release/workflowzip.go` | 1 | `-` | `CG-SIZE-001` | 110 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-713 | `gitmap/release/workflow_test.go` | 1 | `-` | `CG-SIZE-001` | 155 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-714 | `gitmap/release/ziparchive.go` | 1 | `-` | `CG-SIZE-001` | 177 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-715 | `gitmap/release/zipio.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-716 | `gitmap/render/pretty.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-717 | `gitmap/render/prettypost.go` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-718 | `gitmap/render/pretty_emit.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-719 | `gitmap/render/pretty_parse.go` | 1 | `-` | `CG-SIZE-001` | 250 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-720 | `gitmap/render/repotermblock.go` | 1 | `-` | `CG-SIZE-001` | 181 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-721 | `gitmap/render/repotermblock_test.go` | 1 | `-` | `CG-SIZE-001` | 115 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-722 | `gitmap/scanner/progress_test.go` | 1 | `-` | `CG-SIZE-001` | 166 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-723 | `gitmap/scanner/scanner.go` | 1 | `-` | `CG-SIZE-001` | 470 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-724 | `gitmap/scanner/scanner_test.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-725 | `gitmap/scripts/Get-LastRelease.ps1` | 1 | `-` | `CG-SIZE-001` | 141 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-726 | `gitmap/scripts/install.ps1` | 1 | `-` | `CG-SIZE-001` | 1474 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-727 | `gitmap/scripts/install.sh` | 1 | `-` | `CG-SIZE-001` | 1691 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-728 | `gitmap/scripts/release-version.ps1` | 1 | `-` | `CG-SIZE-001` | 525 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-729 | `gitmap/scripts/release-version.sh` | 1 | `-` | `CG-SIZE-001` | 470 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-730 | `gitmap/scripts/uninstall.ps1` | 1 | `-` | `CG-SIZE-001` | 236 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-731 | `gitmap/searcher/db_search.go` | 1 | `-` | `CG-SIZE-001` | 119 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-732 | `gitmap/searcher/finder.go` | 1 | `-` | `CG-SIZE-001` | 174 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-733 | `gitmap/setup/pathsnippetwriter.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-734 | `gitmap/setup/setup.go` | 1 | `-` | `CG-SIZE-001` | 146 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-735 | `gitmap/stablejson/stablejson.go` | 1 | `-` | `CG-SIZE-001` | 211 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-736 | `gitmap/stablejson/stablejson_test.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-737 | `gitmap/stablejson/writers.go` | 1 | `-` | `CG-SIZE-001` | 125 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-738 | `gitmap/startup/add.go` | 1 | `-` | `CG-SIZE-001` | 207 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-739 | `gitmap/startup/addplist.go` | 1 | `-` | `CG-SIZE-001` | 205 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-740 | `gitmap/startup/add_darwin_test.go` | 1 | `-` | `CG-SIZE-001` | 204 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-741 | `gitmap/startup/add_test.go` | 1 | `-` | `CG-SIZE-001` | 163 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-742 | `gitmap/startup/desktop.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-743 | `gitmap/startup/lifecycle_integration_test.go` | 1 | `-` | `CG-SIZE-001` | 245 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-744 | `gitmap/startup/plist.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-745 | `gitmap/startup/plist_test.go` | 1 | `-` | `CG-SIZE-001` | 167 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-746 | `gitmap/startup/remove.go` | 1 | `-` | `CG-SIZE-001` | 196 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-747 | `gitmap/startup/startup.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-748 | `gitmap/startup/startup_test.go` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-749 | `gitmap/startup/winbackend.go` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-750 | `gitmap/startup/windows_lifecycle_test.go` | 1 | `-` | `CG-SIZE-001` | 358 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-751 | `gitmap/startup/windows_test.go` | 1 | `-` | `CG-SIZE-001` | 183 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-752 | `gitmap/startup/winregistry_hklm_windows.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-753 | `gitmap/startup/winregistry_remove_windows.go` | 1 | `-` | `CG-SIZE-001` | 203 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-754 | `gitmap/startup/winregistry_windows.go` | 1 | `-` | `CG-SIZE-001` | 263 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-755 | `gitmap/startup/winshortcut.go` | 1 | `-` | `CG-SIZE-001` | 185 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-756 | `gitmap/startup/winshortcut_linkinfo.go` | 1 | `-` | `CG-SIZE-001` | 170 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-757 | `gitmap/startup/winshortcut_linkinfo_safeuint32_test.go` | 1 | `-` | `CG-SIZE-001` | 114 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-758 | `gitmap/startup/winshortcut_writer.go` | 1 | `-` | `CG-SIZE-001` | 127 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-759 | `gitmap/startup/winshortcut_writer_test.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-760 | `gitmap/store/alias.go` | 1 | `-` | `CG-SIZE-001` | 179 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-761 | `gitmap/store/archive_history.go` | 1 | `-` | `CG-SIZE-001` | 132 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-762 | `gitmap/store/chromeprofile.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-763 | `gitmap/store/csharpmetadata.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-764 | `gitmap/store/downloader_seed.go` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-765 | `gitmap/store/erd_parity_test.go` | 1 | `-` | `CG-SIZE-001` | 153 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-766 | `gitmap/store/group.go` | 1 | `-` | `CG-SIZE-001` | 147 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-767 | `gitmap/store/import.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-768 | `gitmap/store/installedtool.go` | 1 | `-` | `CG-SIZE-001` | 188 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-769 | `gitmap/store/installer_create_test.go` | 1 | `-` | `CG-SIZE-001` | 180 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-770 | `gitmap/store/installer_get_test.go` | 1 | `-` | `CG-SIZE-001` | 152 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-771 | `gitmap/store/installer_list_test.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-772 | `gitmap/store/installer_reset_test.go` | 1 | `-` | `CG-SIZE-001` | 297 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-773 | `gitmap/store/installer_version_test.go` | 1 | `-` | `CG-SIZE-001` | 209 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-774 | `gitmap/store/makeallvisibility.go` | 1 | `-` | `CG-SIZE-001` | 130 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-775 | `gitmap/store/makeallvisibility_undo_test.go` | 1 | `-` | `CG-SIZE-001` | 127 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-776 | `gitmap/store/migrateids.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-777 | `gitmap/store/migrate_commitin_replaymap_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-778 | `gitmap/store/migrate_commitin_replaymap_test.go` | 1 | `-` | `CG-SIZE-001` | 213 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-779 | `gitmap/store/migrate_commitin_test.go` | 1 | `-` | `CG-SIZE-001` | 130 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-780 | `gitmap/store/migrate_v15phase2.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-781 | `gitmap/store/migrate_v15phase4.go` | 1 | `-` | `CG-SIZE-001` | 330 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-782 | `gitmap/store/migrate_v15rebuild.go` | 1 | `-` | `CG-SIZE-001` | 197 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-783 | `gitmap/store/migrate_v15repo.go` | 1 | `-` | `CG-SIZE-001` | 123 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-784 | `gitmap/store/migrations.go` | 1 | `-` | `CG-SIZE-001` | 214 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-785 | `gitmap/store/migrations_test.go` | 1 | `-` | `CG-SIZE-001` | 186 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-786 | `gitmap/store/pendingtask.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-787 | `gitmap/store/project.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-788 | `gitmap/store/release.go` | 1 | `-` | `CG-SIZE-001` | 143 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-789 | `gitmap/store/repo.go` | 1 | `-` | `CG-SIZE-001` | 146 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-790 | `gitmap/store/scan_folder.go` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-791 | `gitmap/store/sshkey.go` | 1 | `-` | `CG-SIZE-001` | 111 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-792 | `gitmap/store/ssh_hist_repo_test.go` | 1 | `-` | `CG-SIZE-001` | 104 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-793 | `gitmap/store/ssh_repo.go` | 1 | `-` | `CG-SIZE-001` | 108 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-794 | `gitmap/store/ssh_repo_test.go` | 1 | `-` | `CG-SIZE-001` | 241 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-795 | `gitmap/store/store.go` | 1 | `-` | `CG-SIZE-001` | 487 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-796 | `gitmap/store/transaction.go` | 1 | `-` | `CG-SIZE-001` | 240 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-797 | `gitmap/store/vscode_project.go` | 1 | `-` | `CG-SIZE-001` | 182 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-798 | `gitmap/store/zipgroup.go` | 1 | `-` | `CG-SIZE-001` | 196 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-799 | `gitmap/templates/corpus_parity_test.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-800 | `gitmap/templates/diff.go` | 1 | `-` | `CG-SIZE-001` | 174 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-801 | `gitmap/templates/diff_test.go` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-802 | `gitmap/templates/list.go` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-803 | `gitmap/templates/list_test.go` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-804 | `gitmap/templates/merge.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-805 | `gitmap/templates/merge_test.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-806 | `gitmap/tests/cmd_test/aliasresolve_test.go` | 1 | `-` | `CG-SIZE-001` | 349 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-807 | `gitmap/tests/cmd_test/amend_test.go` | 1 | `-` | `CG-SIZE-001` | 332 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-808 | `gitmap/tests/cmd_test/seowritecreate_test.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-809 | `gitmap/tests/cmd_test/seowritecsv_test.go` | 1 | `-` | `CG-SIZE-001` | 168 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-810 | `gitmap/tests/cmd_test/seowriteloop_test.go` | 1 | `-` | `CG-SIZE-001` | 242 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-811 | `gitmap/tests/cmd_test/seowritetemplate_test.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-812 | `gitmap/tests/cmd_test/seowrite_test.go` | 1 | `-` | `CG-SIZE-001` | 121 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-813 | `gitmap/tests/cmd_test/temprelease_test.go` | 1 | `-` | `CG-SIZE-001` | 101 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-814 | `gitmap/tests/constants_test/seo_constants_test.go` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-815 | `gitmap/tests/fixrepo_test/fixture_helpers_test.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-816 | `gitmap/tests/fixrepo_test/gofmt_e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-817 | `gitmap/tests/release_test/e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 281 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-818 | `gitmap/tests/release_test/edgecase_test.go` | 1 | `-` | `CG-SIZE-001` | 400 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-819 | `gitmap/tests/release_test/rollback_test.go` | 1 | `-` | `CG-SIZE-001` | 232 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-820 | `gitmap/tests/release_test/skipmeta_test.go` | 1 | `-` | `CG-SIZE-001` | 240 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-821 | `gitmap/tests/scanclone_test/e2e_test.go` | 1 | `-` | `CG-SIZE-001` | 224 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-822 | `gitmap/tests/store_test/alias_suggest_test.go` | 1 | `-` | `CG-SIZE-001` | 233 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-823 | `gitmap/tests/store_test/alias_test.go` | 1 | `-` | `CG-SIZE-001` | 306 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-824 | `gitmap/tests/store_test/location_test.go` | 1 | `-` | `CG-SIZE-001` | 138 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-825 | `gitmap/tests/store_test/pendingtask_test.go` | 1 | `-` | `CG-SIZE-001` | 361 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-826 | `gitmap/tests/store_test/template_test.go` | 1 | `-` | `CG-SIZE-001` | 195 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-827 | `gitmap/theme/filter.go` | 1 | `-` | `CG-SIZE-001` | 101 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-828 | `gitmap/theme/install.go` | 1 | `-` | `CG-SIZE-001` | 135 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-829 | `gitmap/tui/browser.go` | 1 | `-` | `CG-SIZE-001` | 178 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-830 | `gitmap/tui/dashboard.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-831 | `gitmap/tui/logs.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-832 | `gitmap/tui/releases.go` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-833 | `gitmap/tui/reltrigger.go` | 1 | `-` | `CG-SIZE-001` | 172 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-834 | `gitmap/tui/tempreleases.go` | 1 | `-` | `CG-SIZE-001` | 158 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-835 | `gitmap/tui/tui.go` | 1 | `-` | `CG-SIZE-001` | 165 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-836 | `gitmap/tui/tuiview.go` | 1 | `-` | `CG-SIZE-001` | 102 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-837 | `gitmap/tui/tui_test.go` | 1 | `-` | `CG-SIZE-001` | 254 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-838 | `gitmap/tui/zipgroups.go` | 1 | `-` | `CG-SIZE-001` | 133 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-839 | `gitmap/txn/action.go` | 1 | `-` | `CG-SIZE-001` | 176 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-840 | `gitmap/txn/action_test.go` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-841 | `gitmap/txn/journal.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-842 | `gitmap/txn/revert.go` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-843 | `gitmap/txn/snapshot.go` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-844 | `gitmap/visibility/exclude.go` | 1 | `-` | `CG-SIZE-001` | 126 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-845 | `gitmap/visibility/fuzzy.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-846 | `gitmap/visibility/pattern.go` | 1 | `-` | `CG-SIZE-001` | 131 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-847 | `gitmap/vscodepm/autotags_custom.go` | 1 | `-` | `CG-SIZE-001` | 197 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-848 | `gitmap/vscodepm/autotags_custom_test.go` | 1 | `-` | `CG-SIZE-001` | 166 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-849 | `gitmap/vscodepm/mergemode.go` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-850 | `gitmap/vscodepm/mergemode_test.go` | 1 | `-` | `CG-SIZE-001` | 140 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-851 | `gitmap/vscodepm/path.go` | 1 | `-` | `CG-SIZE-001` | 128 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-852 | `gitmap/vscodepm/path_test.go` | 1 | `-` | `CG-SIZE-001` | 159 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-853 | `gitmap/vscodepm/sync.go` | 1 | `-` | `CG-SIZE-001` | 120 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-854 | `gitmap/vscodeworkspace/build.go` | 1 | `-` | `CG-SIZE-001` | 109 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-855 | `gitmap/workspacesync/workspacesync.go` | 1 | `-` | `CG-SIZE-001` | 141 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-856 | `gitmap-updater/cmd/run.go` | 1 | `-` | `CG-SIZE-001` | 107 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-857 | `gitmap-updater/cmd/worker.go` | 1 | `-` | `CG-SIZE-001` | 118 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-858 | `linter-scripts/allowlist-forbidden-string.py` | 1 | `-` | `CG-SIZE-001` | 461 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-859 | `linter-scripts/check-boolean-guidelines.py` | 1 | `-` | `CG-SIZE-001` | 205 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-860 | `linter-scripts/check-boolean-guidelines.py` | 169 | `-` | `CG-BOOL-002` | `violations.append((line_num, f"Inverted success ch` | Major | Replace with affirmative naming (isFail, isMissing) |
| V-861 | `linter-scripts/check-enum-and-boolean.py` | 1 | `-` | `CG-SIZE-001` | 239 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-862 | `linter-scripts/check-enum-and-boolean.py` | 7 | `-` | `CG-BOOL-002` | `2. No inverted success checks (!isSuccess).` | Major | Replace with affirmative naming (isFail, isMissing) |
| V-863 | `linter-scripts/check-enum-and-boolean.py` | 193 | `-` | `CG-BOOL-002` | `violations.append(f"{filepath}:{idx}: Inverted suc` | Major | Replace with affirmative naming (isFail, isMissing) |
| V-864 | `linter-scripts/check-error-management.py` | 1 | `-` | `CG-SIZE-001` | 177 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-865 | `linter-scripts/check-file-sizes.py` | 1 | `-` | `CG-SIZE-001` | 246 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-866 | `linter-scripts/check-forbidden-strings.py` | 1 | `-` | `CG-SIZE-001` | 160 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-867 | `linter-scripts/check-function-lengths.py` | 1 | `-` | `CG-SIZE-001` | 258 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-868 | `linter-scripts/check-markdown-headings.py` | 1 | `-` | `CG-SIZE-001` | 177 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-869 | `linter-scripts/check-mws-error-codes.py` | 1 | `-` | `CG-SIZE-001` | 199 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-870 | `linter-scripts/check-nested-ifs.py` | 1 | `-` | `CG-SIZE-001` | 226 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-871 | `linter-scripts/check-placeholder-comments.py` | 1 | `-` | `CG-SIZE-001` | 3148 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-872 | `linter-scripts/check-prompts-loaded.py` | 1 | `-` | `CG-SIZE-001` | 134 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-873 | `linter-scripts/check-readme-canonicals.py` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-874 | `linter-scripts/check-readme-install-section.py` | 1 | `-` | `CG-SIZE-001` | 301 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-875 | `linter-scripts/check-root-readme.py` | 1 | `-` | `CG-SIZE-001` | 161 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-876 | `linter-scripts/check-runner-dispatch-antipatterns.sh` | 1 | `-` | `CG-SIZE-001` | 223 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-877 | `linter-scripts/check-spec-cross-links.py` | 1 | `-` | `CG-SIZE-001` | 280 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-878 | `linter-scripts/check-spec-folder-refs.py` | 1 | `-` | `CG-SIZE-001` | 253 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-879 | `linter-scripts/check-tunable-constants.py` | 1 | `-` | `CG-SIZE-001` | 439 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-880 | `linter-scripts/forbidden-strings-summary.py` | 1 | `-` | `CG-SIZE-001` | 398 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-881 | `linter-scripts/run.ps1` | 1 | `-` | `CG-SIZE-001` | 144 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-882 | `linter-scripts/run.sh` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-883 | `linter-scripts/suggest-spec-cross-link-fixes.py` | 1 | `-` | `CG-SIZE-001` | 359 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-884 | `linter-scripts/validate-guidelines.go` | 1 | `-` | `CG-SIZE-001` | 1045 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-885 | `linter-scripts/validate-guidelines.py` | 1 | `-` | `CG-SIZE-001` | 1222 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-886 | `linter-scripts/validate-rename-intake.py` | 1 | `-` | `CG-SIZE-001` | 376 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-887 | `linter-scripts/tests/check-file-sizes.test.py` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-888 | `linter-scripts/tests/test_audit_parse_helpers_trailing_summary.py` | 1 | `-` | `CG-SIZE-001` | 271 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-889 | `linter-scripts/tests/test_audit_reason_field_present.py` | 1 | `-` | `CG-SIZE-001` | 331 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-890 | `linter-scripts/tests/test_cache_segregation.py` | 1 | `-` | `CG-SIZE-001` | 225 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-891 | `linter-scripts/tests/test_case_insensitive_extensions.py` | 1 | `-` | `CG-SIZE-001` | 145 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-892 | `linter-scripts/tests/test_dedupe_changed_files_flag.py` | 1 | `-` | `CG-SIZE-001` | 188 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-893 | `linter-scripts/tests/test_diff_prev_shorthand.py` | 1 | `-` | `CG-SIZE-001` | 157 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-894 | `linter-scripts/tests/test_example_payload_coverage.py` | 1 | `-` | `CG-SIZE-001` | 139 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-895 | `linter-scripts/tests/test_ignored_deleted_audit_coverage.py` | 1 | `-` | `CG-SIZE-001` | 293 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-896 | `linter-scripts/tests/test_ignored_deleted_reason_cross_surface_parity.py` | 1 | `-` | `CG-SIZE-001` | 288 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-897 | `linter-scripts/tests/test_ignored_deleted_reason_exact_match.py` | 1 | `-` | `CG-SIZE-001` | 393 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-898 | `linter-scripts/tests/test_include_mdx_flag.py` | 1 | `-` | `CG-SIZE-001` | 207 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-899 | `linter-scripts/tests/test_include_txt_flag.py` | 1 | `-` | `CG-SIZE-001` | 194 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-900 | `linter-scripts/tests/test_list_changed_files_flag.py` | 1 | `-` | `CG-SIZE-001` | 218 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-901 | `linter-scripts/tests/test_list_changed_files_json_reason_schema.py` | 1 | `-` | `CG-SIZE-001` | 292 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-902 | `linter-scripts/tests/test_list_changed_files_verbose.py` | 1 | `-` | `CG-SIZE-001` | 344 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-903 | `linter-scripts/tests/test_only_changed_status_flag.py` | 1 | `-` | `CG-SIZE-001` | 219 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-904 | `linter-scripts/tests/test_placeholder_diff_parser.py` | 1 | `-` | `CG-SIZE-001` | 511 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-905 | `linter-scripts/tests/test_similarity_csv_export.py` | 1 | `-` | `CG-SIZE-001` | 156 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-906 | `linter-scripts/tests/test_similarity_csv_format.py` | 1 | `-` | `CG-SIZE-001` | 212 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-907 | `linter-scripts/tests/test_similarity_csv_stderr_parity.py` | 1 | `-` | `CG-SIZE-001` | 261 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-908 | `linter-scripts/tests/test_similarity_labels_flag.py` | 1 | `-` | `CG-SIZE-001` | 233 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-909 | `linter-scripts/tests/test_similarity_legend_flag.py` | 1 | `-` | `CG-SIZE-001` | 349 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-910 | `linter-scripts/tests/test_symlinked_uppercase_md.py` | 1 | `-` | `CG-SIZE-001` | 258 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-911 | `linter-scripts/tests/test_tunable_constants_t4.py` | 1 | `-` | `CG-SIZE-001` | 278 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-912 | `linter-scripts/tests/test_validate_rename_intake.py` | 1 | `-` | `CG-SIZE-001` | 311 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-913 | `linter-scripts/tests/test_with_similarity_flag.py` | 1 | `-` | `CG-SIZE-001` | 291 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-914 | `linter-scripts/tests/_audit_parse_helpers.py` | 1 | `-` | `CG-SIZE-001` | 218 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-915 | `remotion-demo/src/commands.ts` | 1 | `-` | `CG-SIZE-001` | 122 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-916 | `remotion-demo/src/Terminal.tsx` | 1 | `Terminal` | `CG-SIZE-002` | 227 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-917 | `scripts/audit_codebase.py` | 1 | `-` | `CG-SIZE-001` | 112 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-918 | `scripts/build-stamp.ps1` | 1 | `-` | `CG-SIZE-001` | 129 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-919 | `scripts/build-stamp.sh` | 1 | `-` | `CG-SIZE-001` | 158 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-920 | `scripts/format-go.sh` | 1 | `-` | `CG-SIZE-001` | 101 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-921 | `scripts/generate_cg_audit_report.py` | 1 | `-` | `CG-SIZE-001` | 372 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-922 | `scripts/misspell-local.sh` | 1 | `-` | `CG-SIZE-001` | 124 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-923 | `scripts/preflight-ci.sh` | 1 | `-` | `CG-SIZE-001` | 137 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-924 | `scripts/changelog/internal/gitlog/gitlog.go` | 1 | `-` | `CG-SIZE-001` | 175 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-925 | `scripts/changelog/internal/group/group.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-926 | `scripts/changelog/internal/runner/execute.go` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-927 | `scripts/fix-repo/Config.ps1` | 1 | `-` | `CG-SIZE-001` | 103 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-928 | `scripts/fix-repo/config.sh` | 1 | `-` | `CG-SIZE-001` | 116 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-929 | `scripts/kubernetes/02-ubuntu-install/01-zsh-theme-change-v2.sh` | 1 | `-` | `CG-SIZE-001` | 274 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-930 | `scripts/kubernetes/02-ubuntu-install/05-omy-zsh-only.sh` | 1 | `-` | `CG-SIZE-001` | 113 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-931 | `scripts/kubernetes/02-ubuntu-install/09-create-root-user-v2.sh` | 1 | `-` | `CG-SIZE-001` | 136 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-932 | `scripts/kubernetes/05-server-cmds/02-run-cmd-v2.sh` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-933 | `scripts/kubernetes/07-remote-commands/run-cmd.sh` | 1 | `-` | `CG-SIZE-001` | 121 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-934 | `spec/11-powershell-integration/templates/run.ps1` | 1 | `-` | `CG-SIZE-001` | 856 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-935 | `src/App.tsx` | 1 | `App` | `CG-SIZE-002` | 195 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-936 | `src/components/docs/CloneNextCommandBuilder.tsx` | 1 | `CloneNextCommandBuilder` | `CG-SIZE-002` | 317 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-937 | `src/components/docs/CloneNextCommandBuilder.tsx` | 78 | `-` | `CG-CASE-001` | `function buildOriginURL(s: BuilderState): string {` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-938 | `src/components/docs/CloneNextCommandBuilder.tsx` | 90 | `-` | `CG-CASE-001` | `function buildTargetURL(s: BuilderState): string {` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-939 | `src/components/docs/CloneNextCommandBuilder.tsx` | 94 | `-` | `CG-CASE-001` | `return buildOriginURL(s).replace(current, target);` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-940 | `src/components/docs/CloneNextCommandBuilder.tsx` | 113 | `-` | `CG-CASE-001` | `const url = buildTargetURL(s);` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-941 | `src/components/docs/CloneNextCommandBuilder.tsx` | 131 | `-` | `CG-CASE-001` | ``origin     : ${buildOriginURL(s)}`,` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-942 | `src/components/docs/CloneNextCommandBuilder.tsx` | 134 | `-` | `CG-CASE-001` | ``target url : ${buildTargetURL(s)}`,` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-943 | `src/components/docs/CodeBlock.tsx` | 1 | `CodeBlock` | `CG-SIZE-002` | 384 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-944 | `src/components/docs/CommandCard.tsx` | 1 | `CommandCard` | `CG-SIZE-002` | 112 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-945 | `src/components/docs/CommandPalette.tsx` | 1 | `CommandPalette` | `CG-SIZE-002` | 137 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-946 | `src/components/docs/CommitTransferPage.tsx` | 1 | `CommitTransferPage` | `CG-SIZE-002` | 304 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-947 | `src/components/docs/CommitTransferPage.tsx` | 204 | `-` | `CG-CASE-001` | `paths or <code>https://</code> / <code>git@</code>` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-948 | `src/components/docs/DocsSidebar.tsx` | 1 | `DocsSidebar` | `CG-SIZE-002` | 194 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-949 | `src/components/docs/DocsTooltip.tsx` | 1 | `DocsTooltip` | `CG-SIZE-002` | 116 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-950 | `src/components/docs/SpecPage.tsx` | 1 | `SpecPage` | `CG-SIZE-002` | 120 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-951 | `src/components/docs/TabOrderMap.tsx` | 1 | `TabOrderMap` | `CG-SIZE-002` | 491 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-952 | `src/components/docs/TerminalDemo.tsx` | 1 | `TerminalDemo` | `CG-SIZE-002` | 114 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-953 | `src/components/docs/troubleshooting/DiagnosticChecklist.tsx` | 1 | `DiagnosticChecklist` | `CG-SIZE-002` | 124 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-954 | `src/components/projects/ProjectDetailDialog.tsx` | 1 | `ProjectDetailDialog` | `CG-SIZE-002` | 188 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-955 | `src/components/spec/specData.ts` | 1 | `-` | `CG-SIZE-001` | 233 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-956 | `src/components/ui/alert-dialog.tsx` | 1 | `alert-dialog` | `CG-SIZE-002` | 105 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-957 | `src/components/ui/carousel.tsx` | 1 | `carousel` | `CG-SIZE-002` | 237 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-958 | `src/components/ui/chart.tsx` | 1 | `chart` | `CG-SIZE-002` | 320 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-959 | `src/components/ui/command.tsx` | 1 | `command` | `CG-SIZE-002` | 133 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-960 | `src/components/ui/context-menu.tsx` | 1 | `context-menu` | `CG-SIZE-002` | 179 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-961 | `src/components/ui/dropdown-menu.tsx` | 1 | `dropdown-menu` | `CG-SIZE-002` | 180 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-962 | `src/components/ui/form.tsx` | 1 | `form` | `CG-SIZE-002` | 133 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-963 | `src/components/ui/menubar.tsx` | 1 | `menubar` | `CG-SIZE-002` | 208 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-964 | `src/components/ui/navigation-menu.tsx` | 1 | `navigation-menu` | `CG-SIZE-002` | 121 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-965 | `src/components/ui/select.tsx` | 1 | `select` | `CG-SIZE-002` | 144 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-966 | `src/components/ui/sheet.tsx` | 1 | `sheet` | `CG-SIZE-002` | 108 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-967 | `src/components/ui/sidebar.tsx` | 1 | `sidebar` | `CG-SIZE-002` | 667 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-968 | `src/components/ui/toast.tsx` | 1 | `toast` | `CG-SIZE-002` | 112 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-969 | `src/data/changelog.ts` | 1 | `-` | `CG-SIZE-001` | 2292 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-970 | `src/data/changelog.ts` | 562 | `-` | `CG-CASE-001` | `subtitle: "`gitmap help --json` for scripting + ID` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-971 | `src/data/changelog.ts` | 577 | `-` | `CG-CASE-001` | `"**Help-file coverage audit.** Verified every prim` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-972 | `src/data/changelog.ts` | 622 | `-` | `CG-CASE-001` | `"**Critical fix.** Running `gitmap fix-repo` insid` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-973 | `src/data/changelog.ts` | 686 | `-` | `CG-CASE-001` | `"Test fix: registered `CmdPush` / `CmdPushAlias` i` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-974 | `src/data/changelog.ts` | 698 | `-` | `CG-CASE-001` | `"Added: when both flags are set, `--ssh` wins with` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-975 | `src/data/changelog.ts` | 736 | `-` | `CG-CASE-001` | `"Pinned: README pinned-version block + version mat` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-976 | `src/data/changelog.ts` | 739 | `-` | `CG-CASE-001` | `"Clarified: SSH-shorthand and `ssh://` URLs still ` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-977 | `src/data/changelog.ts` | 749 | `-` | `CG-CASE-001` | `"Clarified: SSH-shorthand and `ssh://` URLs alread` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-978 | `src/data/changelog.ts` | 757 | `-` | `CG-CASE-001` | `subtitle: "Root-level installer URLs — `/install.p` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-979 | `src/data/changelog.ts` | 761 | `-` | `CG-CASE-001` | `"Pinned: README pinned-version block + version mat` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-980 | `src/data/changelog.ts` | 770 | `-` | `CG-CASE-001` | `"Pinned: README pinned-version block + version mat` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-981 | `src/data/changelog.ts` | 794 | `-` | `CG-CASE-001` | `"Files: `gitmap-v28/cmd/cloneurlconvert.go` (new —` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-982 | `src/data/changelog.ts` | 833 | `-` | `CG-CASE-001` | `"Release body: `uploadToGitHub` in `gitmap-v28/rel` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-983 | `src/data/changelog.ts` | 844 | `-` | `CG-CASE-001` | `"New `gitmap install gitmap-oneliner` command prin` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-984 | `src/data/changelog.ts` | 846 | `-` | `CG-CASE-001` | `"New spec `spec/01-app/109-install-gitmap-oneliner` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-985 | `src/data/changelog.ts` | 867 | `-` | `CG-CASE-001` | `"Files: `README.md` (pinned-version block + 4 CI b` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-986 | `src/data/changelog.ts` | 994 | `-` | `CG-CASE-001` | `"New `gitmap clone-fix-repo` (alias `cfr`) chains ` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-987 | `src/data/changelog.ts` | 1210 | `-` | `CG-CASE-001` | `"The `shouldRewriteToClone` / `looksLikeURLToken` ` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-988 | `src/data/changelog.ts` | 1220 | `-` | `CG-CASE-001` | `"Added `gitmap zip` — accepts N heterogeneous sour` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-989 | `src/data/changelog.ts` | 1263 | `-` | `CG-CASE-001` | `"New `gitmap-v28/cmd/scanbenchmark.go` captures wa` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-990 | `src/data/changelog.ts` | 1284 | `-` | `CG-CASE-001` | `"Bash side (`install-quick.sh`): each probe runs i` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-991 | `src/data/changelog.ts` | 1308 | `-` | `CG-CASE-001` | `"Added 10 new files to the embedded corpus under `` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-992 | `src/data/changelog.ts` | 1323 | `-` | `CG-CASE-001` | `"Implemented `--force` as a pre-merge `os.Remove(t` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-993 | `src/data/changelog.ts` | 1577 | `-` | `CG-CASE-001` | `"`gitmap clone <url>` auto-flattens versioned URLs` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-994 | `src/data/changelog.ts` | 1638 | `-` | `CG-CASE-001` | `"Fixed `ShouldPrintInstallHint` not matching SSH r` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-995 | `src/data/changelog.ts` | 1730 | `-` | `CG-CASE-001` | `"Added `SourceRepoCloneURL`, `MsgUpdateCloning`, `` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-996 | `src/data/changelog.ts` | 1826 | `-` | `CG-CASE-001` | `"Interactive prompt to terminate locking processes` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-997 | `src/data/changelog.ts` | 1847 | `-` | `CG-CASE-001` | `"Repository renamed from `git-repo-navigator` to `` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-998 | `src/data/changelog.ts` | 1958 | `-` | `CG-CASE-001` | `"All DB query errors from legacy string-based IDs ` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-999 | `src/data/commands.ts` | 1 | `-` | `CG-SIZE-001` | 2486 lines (hard cap: 300 lines) | Critical | Decompose monolithic file into focused domain modules |
| V-1000 | `src/data/commands.ts` | 74 | `-` | `CG-CASE-001` | `{ command: "gitmap s --output json --mode ssh", de` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1001 | `src/data/commands.ts` | 459 | `-` | `CG-CASE-001` | `name: "commit-in", alias: "cin", description: "Rep` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1002 | `src/data/commands.ts` | 463 | `-` | `CG-CASE-001` | `{ flag: "<sources>", description: "Comma-separated` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1003 | `src/data/commands.ts` | 864 | `-` | `CG-CASE-001` | `{ name: "visibility-history", description: "List p` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1004 | `src/data/commands.ts` | 883 | `-` | `CG-CASE-001` | `{ name: "visibility-history", description: "List p` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1005 | `src/data/commands.ts` | 888 | `-` | `CG-CASE-001` | `name: "clone-fix-repo", alias: "cfr", description:` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1006 | `src/data/commands.ts` | 891 | `-` | `CG-CASE-001` | `{ flag: "--ssh", description: "Force SSH transport` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1007 | `src/data/commands.ts` | 892 | `-` | `CG-CASE-001` | `{ flag: "--https", description: "Force HTTPS trans` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1008 | `src/data/commands.ts` | 1081 | `-` | `CG-CASE-001` | `{ flag: "--verbose", description: "Show full paths` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1009 | `src/data/commands.ts` | 1636 | `-` | `CG-CASE-001` | `{ flag: "--force", description: "Replace pre-exist` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1010 | `src/data/commands.ts` | 2313 | `-` | `CG-CASE-001` | `{ flag: "--id <list>", description: "Target specif` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1011 | `src/data/commands.ts` | 2338 | `-` | `CG-CASE-001` | `{ flag: "--id <list>", description: "Run only on t` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1012 | `src/data/postMortems.ts` | 1 | `-` | `CG-SIZE-001` | 149 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-1013 | `src/data/postMortems.ts` | 36 | `-` | `CG-CASE-001` | `summary: "Update process entered an infinite retry` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1014 | `src/data/troubleshootingIssues.ts` | 1 | `-` | `CG-SIZE-001` | 287 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-1015 | `src/data/troubleshootingIssues.ts` | 162 | `-` | `CG-CASE-001` | `title: "Repos cloned with SSH URLs but org only al` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1016 | `src/data/troubleshootingIssues.ts` | 208 | `-` | `CG-CASE-001` | `fix: "Close the holder (VS Code, Explorer preview ` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1017 | `src/hooks/use-toast.ts` | 1 | `-` | `CG-SIZE-001` | 166 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-1018 | `src/hooks/useTheme.ts` | 1 | `-` | `CG-SIZE-001` | 121 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-1019 | `src/pages/Alias.tsx` | 1 | `Alias` | `CG-SIZE-002` | 263 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1020 | `src/pages/Architecture.tsx` | 1 | `Architecture` | `CG-SIZE-002` | 190 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1021 | `src/pages/Architecture.tsx` | 57 | `-` | `CG-CASE-001` | `<span className="text-muted-foreground">→ extracts` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1022 | `src/pages/BatchActions.tsx` | 1 | `BatchActions` | `CG-SIZE-002` | 332 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1023 | `src/pages/Bookmarks.tsx` | 1 | `Bookmarks` | `CG-SIZE-002` | 171 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1024 | `src/pages/Cd.tsx` | 1 | `Cd` | `CG-SIZE-002` | 223 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1025 | `src/pages/Changelog.tsx` | 1 | `Changelog` | `CG-SIZE-002` | 225 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1026 | `src/pages/ChangelogGenerate.tsx` | 1 | `ChangelogGenerate` | `CG-SIZE-002` | 178 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1027 | `src/pages/ChromeProfileSpec.tsx` | 1 | `ChromeProfileSpec` | `CG-SIZE-002` | 188 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1028 | `src/pages/ChromeProfileSpec.tsx` | 17 | `-` | `CG-CASE-001` | `{ name: "chrome-profile-import", alias: "cpi", des` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1029 | `src/pages/ClearReleaseJSON.tsx` | 1 | `ClearReleaseJSON` | `CG-SIZE-002` | 239 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1030 | `src/pages/CloneCommand.tsx` | 1 | `CloneCommand` | `CG-SIZE-002` | 124 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1031 | `src/pages/CloneFixRepo.tsx` | 1 | `CloneFixRepo` | `CG-SIZE-002` | 129 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1032 | `src/pages/CloneFixRepo.tsx` | 34 | `-` | `CG-CASE-001` | `{ icon: GitBranch, title: "Clones first", desc: "V` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1033 | `src/pages/CloneMultiSpec.tsx` | 6 | `-` | `CG-CASE-001` | `title="gitmap clone — Multiple URLs"` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1034 | `src/pages/CloneNext.tsx` | 1 | `CloneNext` | `CG-SIZE-002` | 371 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1035 | `src/pages/CloneNext.tsx` | 253 | `-` | `CG-CASE-001` | `"  PID     Process",` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1036 | `src/pages/CloneNext.tsx` | 259 | `-` | `CG-CASE-001` | `"✓ Terminated Code.exe (PID 14320)",` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1037 | `src/pages/CloneNext.tsx` | 260 | `-` | `CG-CASE-001` | `"✓ Terminated explorer.exe (PID 8412)",` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1038 | `src/pages/CloneNextCommand.tsx` | 1 | `CloneNextCommand` | `CG-SIZE-002` | 428 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1039 | `src/pages/CloneNextCommand.tsx` | 151 | `-` | `CG-CASE-001` | `— gitmap lists the offending PIDs (e.g. <code clas` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1040 | `src/pages/CloneOverview.tsx` | 1 | `CloneOverview` | `CG-SIZE-002` | 155 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1041 | `src/pages/Commands.tsx` | 1 | `Commands` | `CG-SIZE-002` | 255 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1042 | `src/pages/CommitIn.tsx` | 1 | `CommitIn` | `CG-SIZE-002` | 235 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1043 | `src/pages/CommitInExamples.tsx` | 1 | `CommitInExamples` | `CG-SIZE-002` | 194 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1044 | `src/pages/CommitInExamples.tsx` | 30 | `-` | `CG-CASE-001` | `remote URLs — each URL is shallow-cloned into{" "}` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1045 | `src/pages/Config.tsx` | 1 | `Config` | `CG-SIZE-002` | 199 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1046 | `src/pages/Dashboard.tsx` | 1 | `Dashboard` | `CG-SIZE-002` | 134 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1047 | `src/pages/DesignSystem.tsx` | 1 | `DesignSystem` | `CG-SIZE-002` | 398 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1048 | `src/pages/Diff.tsx` | 1 | `Diff` | `CG-SIZE-002` | 177 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1049 | `src/pages/DiffProfiles.tsx` | 1 | `DiffProfiles` | `CG-SIZE-002` | 214 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1050 | `src/pages/Doctor.tsx` | 1 | `Doctor` | `CG-SIZE-002` | 186 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1051 | `src/pages/Export.tsx` | 1 | `Export` | `CG-SIZE-002` | 205 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1052 | `src/pages/FixRepo.tsx` | 1 | `FixRepo` | `CG-SIZE-002` | 143 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1053 | `src/pages/FlagReference.tsx` | 1 | `FlagReference` | `CG-SIZE-002` | 175 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1054 | `src/pages/GenericCLI.tsx` | 1 | `GenericCLI` | `CG-SIZE-002` | 1108 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1055 | `src/pages/GenericCLI.tsx` | 591 | `-` | `CG-CASE-001` | `["Git operations", "Clone/pull commands, remote UR` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1056 | `src/pages/GettingStarted.tsx` | 1 | `GettingStarted` | `CG-SIZE-002` | 107 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1057 | `src/pages/GoMod.tsx` | 1 | `GoMod` | `CG-SIZE-002` | 264 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1058 | `src/pages/HelpDashboard.tsx` | 1 | `HelpDashboard` | `CG-SIZE-002` | 163 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1059 | `src/pages/HelpIndex.tsx` | 1 | `HelpIndex` | `CG-SIZE-002` | 251 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1060 | `src/pages/History.tsx` | 1 | `History` | `CG-SIZE-002` | 183 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1061 | `src/pages/HistoryRewrite.tsx` | 1 | `HistoryRewrite` | `CG-SIZE-002` | 191 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1062 | `src/pages/Import.tsx` | 1 | `Import` | `CG-SIZE-002` | 235 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1063 | `src/pages/Index.tsx` | 1 | `Index` | `CG-SIZE-002` | 153 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1064 | `src/pages/Install.tsx` | 1 | `Install` | `CG-SIZE-002` | 531 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1065 | `src/pages/Install.tsx` | 499 | `-` | `CG-CASE-001` | `["constants/constants_install.go", "Tool names, pa` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1066 | `src/pages/InstallGitmap.tsx` | 1 | `InstallGitmap` | `CG-SIZE-002` | 239 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1067 | `src/pages/InteractiveExamples.tsx` | 1 | `InteractiveExamples` | `CG-SIZE-002` | 125 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1068 | `src/pages/InteractiveTUI.tsx` | 1 | `InteractiveTUI` | `CG-SIZE-002` | 208 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1069 | `src/pages/MakeAllPrivate.tsx` | 1 | `MakeAllPrivate` | `CG-SIZE-002` | 140 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1070 | `src/pages/MakeAllPublic.tsx` | 1 | `MakeAllPublic` | `CG-SIZE-002` | 143 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1071 | `src/pages/Makefile.tsx` | 1 | `Makefile` | `CG-SIZE-002` | 110 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1072 | `src/pages/MakePublic.tsx` | 1 | `MakePublic` | `CG-SIZE-002` | 141 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1073 | `src/pages/MergeBoth.tsx` | 1 | `MergeBoth` | `CG-SIZE-002` | 178 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1074 | `src/pages/MergeLeft.tsx` | 1 | `MergeLeft` | `CG-SIZE-002` | 146 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1075 | `src/pages/MergeRight.tsx` | 1 | `MergeRight` | `CG-SIZE-002` | 146 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1076 | `src/pages/Mv.tsx` | 1 | `Mv` | `CG-SIZE-002` | 179 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1077 | `src/pages/PostMortems.tsx` | 1 | `PostMortems` | `CG-SIZE-002` | 118 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1078 | `src/pages/Profile.tsx` | 1 | `Profile` | `CG-SIZE-002` | 233 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1079 | `src/pages/ProjectDetection.tsx` | 1 | `ProjectDetection` | `CG-SIZE-002` | 450 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1080 | `src/pages/ProjectDetection.tsx` | 67 | `-` | `CG-CASE-001` | `["constants/constants_project.go", "Detection IDs,` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1081 | `src/pages/Projects.tsx` | 1 | `Projects` | `CG-SIZE-002` | 184 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1082 | `src/pages/Prune.tsx` | 1 | `Prune` | `CG-SIZE-002` | 155 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1083 | `src/pages/Release.tsx` | 1 | `Release` | `CG-SIZE-002` | 629 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1084 | `src/pages/ReleaseAlias.tsx` | 1 | `ReleaseAlias` | `CG-SIZE-002` | 166 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1085 | `src/pages/ReleaseSelf.tsx` | 1 | `ReleaseSelf` | `CG-SIZE-002` | 202 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1086 | `src/pages/ReleaseVersion.tsx` | 1 | `ReleaseVersion` | `CG-SIZE-002` | 220 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1087 | `src/pages/Replace.tsx` | 1 | `Replace` | `CG-SIZE-002` | 155 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1088 | `src/pages/ScanCloneFlags.tsx` | 1 | `ScanCloneFlags` | `CG-SIZE-002` | 268 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1089 | `src/pages/ScanCommand.tsx` | 1 | `ScanCommand` | `CG-SIZE-002` | 115 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1090 | `src/pages/ScanCommand.tsx` | 80 | `-` | `CG-CASE-001` | `<h3 className="text-base font-heading font-semibol` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1091 | `src/pages/ScanCommand.tsx` | 89 | `-` | `CG-CASE-001` | `✓ Clone URLs use SSH format (git@github.com:...)`}` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1092 | `src/pages/Setup.tsx` | 1 | `Setup` | `CG-SIZE-002` | 285 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1093 | `src/pages/SpecIndex.tsx` | 1 | `SpecIndex` | `CG-SIZE-002` | 163 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1094 | `src/pages/SSH.tsx` | 1 | `SSH` | `CG-SIZE-002` | 244 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1095 | `src/pages/Stats.tsx` | 1 | `Stats` | `CG-SIZE-002` | 167 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1096 | `src/pages/TempRelease.tsx` | 1 | `TempRelease` | `CG-SIZE-002` | 239 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1097 | `src/pages/Troubleshooting.tsx` | 1 | `Troubleshooting` | `CG-SIZE-002` | 176 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1098 | `src/pages/Troubleshooting.tsx` | 2 | `-` | `CG-CASE-001` | `import { useSearchParams, SetURLSearchParams } fro` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1099 | `src/pages/Troubleshooting.tsx` | 76 | `-` | `CG-CASE-001` | `searchParams: URLSearchParams,` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1100 | `src/pages/Troubleshooting.tsx` | 77 | `-` | `CG-CASE-001` | `setSearchParams: SetURLSearchParams,` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1101 | `src/pages/Troubleshooting.tsx` | 82 | `-` | `CG-CASE-001` | `const nextParams = new URLSearchParams(searchParam` | Minor | Change uppercase abbreviation to PascalCase (Id, Url, Api) |
| V-1102 | `src/pages/VersionHistory.tsx` | 1 | `VersionHistory` | `CG-SIZE-002` | 264 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1103 | `src/pages/Watch.tsx` | 1 | `Watch` | `CG-SIZE-002` | 301 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1104 | `src/pages/ZipGroup.tsx` | 1 | `ZipGroup` | `CG-SIZE-002` | 283 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1105 | `src/test/chip-contrast.test.tsx` | 1 | `chip-contrast.test` | `CG-SIZE-002` | 248 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1106 | `src/test/docs-tooltip-fallback-marker.test.tsx` | 1 | `docs-tooltip-fallback-marker.test` | `CG-SIZE-002` | 126 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1107 | `src/test/docs-tooltips.test.tsx` | 1 | `docs-tooltips.test` | `CG-SIZE-002` | 503 lines (.tsx cap: 100 lines) | Major | Split React component into subcomponents or hooks |
| V-1108 | `src/test/new-command-pages.test.ts` | 1 | `-` | `CG-SIZE-001` | 105 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |
| V-1109 | `src/test/release-version-snippets.test.ts` | 1 | `-` | `CG-SIZE-001` | 117 lines (standard cap: 100 lines) | Minor | Decompose into smaller sub-modules |

---

## Remediation Roadmap & Priority Matrix

1. **Sprint 1 (Critical Zero-Tolerance):** Eliminate all explicit boolean true checks, swallowed errors, and files > 300 lines.
2. **Sprint 2 (Major Architecture):** Refactor monolithic React page files (> 100 lines) and functions > 15 lines into subcomponents and hooks.
3. **Sprint 3 (Cleanliness & Hygiene):** Enforce `*Type` enums, PascalCase abbreviations (`Id`, `Url`), and return new lines (R13-R16).
