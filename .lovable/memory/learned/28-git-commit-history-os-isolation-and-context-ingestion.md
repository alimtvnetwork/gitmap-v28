# 28 — Git Commit History, OS Test Isolation & Project Context Ingestion

- **Slug:** git-commit-history-os-isolation-and-context-ingestion
- **Date:** 2026-09-17
- **Version:** v6.257.0
- **Category:** learned
- **Status:** permanent

---

## 1. Executive Summary

In accordance with the enhanced Read Memory protocol (Prompt Version 2.1.0), the repository context, recent 10 git commits, architectural changes, coding guidelines, error management architecture, database conventions, and Root Cause Analysis (RCA) records were thoroughly loaded and ingested.

The repository is currently at version **v6.257.0**. The recent commits completed Plan 190 (Coding Guideline 24: Isolating Destructive OS & Heavy Unit Tests) and the release of v6.257.0 featuring the cluster triad, scheduled shutdown status/cancel, remote GitMap bootstrap, and OS AI cache cleaning.

---

## 2. Ingested Repository Surface

- **Git Commit History:** Last 10 commits (`a1f1c75109` through `2a16e95450`) analyzed across commit messages, stats, and diffs.
- **Memory Files:** 132 files indexed in `.lovable/memory/01-index.md` across architecture, constraints, features, tech, and workflow; 34 learned memory files in `.lovable/memory/learned/`.
- **Consolidated Guidelines:** 36 policy documents in `spec/17-consolidated-guidelines/`.
- **Spec Authoring:** 31 specification files in `spec/01-spec-authoring-guide/`.
- **Coding Guidelines:** 29 core documents in `spec/02-coding-guidelines/01-cross-language/` along with language-specific guides in Go, TypeScript, PHP, Rust, and Python.
- **Error Management:** 45 specification files and retrospective post-mortems in `spec/03-error-manage/`.
- **Database Conventions:** 8 foundational documents in `spec/04-database-conventions/`.
- **CI/CD Issues & RCA Records:** 49 failure post-mortems in `.lovable/cicd-issues/` and 14 issue files in `.lovable/issues/`.
- **Ambiguity Records:** 1 open ambiguity (`01-spec-gaps.md`), 1 resolved ambiguity (`02-macro-step-open-command-behavior.md`).
- **Plans & Subtasks:** 0 pending plans in `.lovable/plans/pending/`; Plan 190 completed across all 7 subtasks in `.lovable/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/`.

---

## 3. Last 10 Git Commits & Architectural Intent

| Commit | SHA | Author & Date | Architectural Intent |
|--------|-----|---------------|----------------------|
| 1 | `2a16e95450` | Md. Alim Ul Karim (2026-09-17) | Decouple service driver (`DefaultServiceDriverResolver`), osuser command runner (`defaultOSCommandRunner`), and netip manager (`defaultNetIPManagerFactory`) with hermetic unit tests. |
| 2 | `79e5aa5287` | Md. Alim Ul Karim (2026-09-17) | Expand prompt 24 with master pipeline and author `cg-isolate-os-tests` skill in `.agents/skills/`. |
| 3 | `5fa36f40d6` | Md. Alim Ul Karim (2026-09-17) | Implement Coding Guideline 24 isolating destructive OS operations (`CancelSchedulePowerCLI`, `ExecShutdown`, `os_ai_clean`, `osclean`, `os_update`, `os_full_upgrade`). |
| 4 | `1d38830cfc` | Md. Alim Ul Karim (2026-09-17) | Update test inventory execution timings in `.lovable/test-inventory.json`. |
| 5 | `c7470c0f83` | Md. Alim Ul Karim (2026-09-17) | Enrich servers-clients and ssh-join subcommands table with install and power operations. |
| 6 | `bd01ddca49` | Md. Alim Ul Karim (2026-09-17) | Release v6.257.0: cluster triad, schedule shutdown status, remote gitmap bootstrap, and os ai-clean. |
| 7 | `c37a2da248` | Md. Alim Ul Karim (2026-09-17) | Fix golangci-lint, darwin OS detection, and smoke test checks. |
| 8 | `9aad9f8b9f` | Md. Alim Ul Karim (2026-09-17) | Deeply expand cluster triad, schedule power status, and os ai-clean help text across markdown and Go helptext catalog. |
| 9 | `adeb37be5e` | Md. Alim Ul Karim (2026-09-17) | Add scheduled power status inspection (`schedule shutdown status`) and cancel abort commands (`schedule shutdown cancel`). |
| 10 | `a1f1c75109` | Md. Alim Ul Karim (2026-09-17) | Add copy-pasteable examples to schedule shutdown and restart help markdown. |

---

## 4. Coding Guideline 24: OS Test Isolation Standard

Unit tests must NEVER trigger real operating system modifications, power state alterations (shutdown, restart, logoff), process termination, package manager upgrades, or actual file deletions during test suite execution.

### Architectural Decoupling Hooks Established:
1. `cli/cmdschedule/schedule_os.go`: `DefaultOSActionExecutor OSActionExecutor` with `OSActionCancel`.
2. `cli/cluster/exec_lifecycle_test.go`: Hermetic `runCmdFunc` mock injection with `defer` restoration blocks.
3. `cli/cmdos/os_ai_clean.go`: `defaultFileRemover FileRemover = removeSingleFileSafely`.
4. `cli/osclean/clean.go`: `defaultTempDirResolver TempDirResolver = resolveTempDirectories`.
5. `cli/cmdos/os_update.go`: `defaultOSCommandRunner OSCommandRunner`.
6. `cli/cmdservice/driver.go`: `DefaultServiceDriverResolver ServiceDriverResolver`.
7. `cli/osuser/runner.go`: `defaultOSCommandRunner OSCommandRunner`.
8. `cli/cmdos/os_ip.go`: `defaultNetIPManagerFactory NetIPManagerFactory`.

---

## 5. Core Invariants & CODE RED Prohibitions

1. **Zero Raw System Destruction in Tests:** Always inject test hooks and assert commands/paths without real execution.
2. **Never Run Tests Without Owner Command:** Routine guideline turns must use `--no-tests`. Only CI/CD fix tasks and release orchestrations execute full test suites.
3. **Implicit Positive Booleans:** Always evaluate positive booleans implicitly (`if isValid {`). Never write `if isValid == true`.
4. **Boolean Prefixing:** Strict `is*` and `has*` prefixes only. All other prefixes (`can`, `should`, `was`) are banned.
5. **Id Acronym Normalization:** Strict `Id` (PascalCase) and `id` (camelCase). Total ban on all-caps `ID`.
6. **No Bare Void in Go:** Functions must return `Result[T]` or `*appfault.AppError`.
7. **Single Source of Truth for Versions:** `version.json` controls versioning across all tools and languages; manual edits during routine turns are prohibited.
8. **Relative Git Paths Only:** Total ban on `file:///` and machine-specific absolute paths.
9. **Strict Lowercase Readme:** Root readme must remain lowercase `readme.md`.
