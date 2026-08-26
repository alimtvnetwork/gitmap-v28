# 10 — Cross-Spec Blind-AI Readiness Sweep (Non-spec/19 Folders)

> **Scope:** All `spec/` folders **except** `spec/19-main-worker-service/` (which is closed at 96/100 by audit-09).
> **Folders covered:** `00`, `01–18`, `21–24`, plus root meta files (`_template.md`, `99-consistency-report.md`, `authoring-guideline/`, `health-dashboard.md`, `spec-index.md`).
> **Audit date:** 2026-05-07
> **Project version at audit:** **5.28.0**
> **Persona:** same hostile-literal blind-AI used in audits `07-/08-/09-`.
> **Method:** ripgrep sweep for ambiguity markers (`SHOULD`, `MAY`, `TBD`, `TODO`, `FIXME`, `OQ-N`, "ambiguous"), then triage by category.

---

## TL;DR

| Folder cluster                              | Files scanned | Hits | True ambiguity? | Blind-AI risk |
| ------------------------------------------- | ------------: | ---: | :-------------- | :------------ |
| `01-spec-authoring-guide`                   | 17            | 6    | No — meta/scaffolding | None |
| `02-coding-guidelines`                      | 197           | 4    | No — lint-rule examples | None |
| `03-error-manage`                           | 118           | 5    | **Yes (1 file)** | **Low — F6/F8/F13 SHOULD→MUST** |
| `04-database-conventions`                   | 9             | 3    | No — RFC-2119 reference docs | None |
| `05-split-db-architecture`                  | 14            | 0    | — | None |
| `06-seedable-config-architecture`           | 19            | 9    | No — all OQs have `Inferred:` resolutions | None |
| `07-design-system` … `13-generic-cli`       | mixed         | 6    | No — examples and templates | None |
| `14-update`                                 | 44            | 8    | **Yes (2 files)** | **Low — installer-behavior + version-probe** |
| `15-distribution-and-runner`                | 7             | 1    | **Yes (1 line)**  | Trivial |
| `16-generic-release`, `17-consolidated-…`, `18-wp-plugin-how-to` | mixed | 7 | No — lint examples + worked examples | None |
| `21-app` … `24-app-ui-design-system`        | 4             | 0    | — | None |
| **Totals**                                  | ~600          | **51** | **4 files contain operative ambiguity** | **Low overall** |

**Headline:** The non-spec/19 corpus is dramatically cleaner than spec/19 was at audit-07. 47 of 51 hits are **false positives** (resolved OQs with `Inferred:` defaults, RFC-2119 reference text, lint-rule code samples, or template scaffolding). Only **4 files** contain actionable ambiguity for a blind-AI reader, and all are **Low** risk.

---

## 1. False-positive categories (47 hits, no action)

| Category | Examples | Why benign |
| --- | --- | --- |
| **Resolved OQs** | `06-seedable-config/02-features/{07,08,09,10}-*.md`, `14-update/28-worker-push-instruction.md` | Each OQ ships an `Inferred: …` line giving the deterministic default. A literal AI reads the inference as the rule. |
| **Lint-rule examples** | `17-consolidated-guidelines/27-linter-authoring-guide.md` (lines 361–424) | Documenting *how* to detect stale TODOs — not stale TODOs themselves. |
| **RFC-2119 reference docs** | `01-spec-authoring-guide/{04,07,10,11,17}.md`, `04-database-conventions/{00,01,02}.md` | The keyword `SHOULD` appears in normative-language definitions, not as an operative requirement. |
| **Code-style allowances** | `02-coding-guidelines/02-typescript/08-…`, `01-cross-language/16-static-analysis/09-…`, `04-code-style/06-…` | Permit `// TODO(TICKET-N)` comments — a coding rule, not a spec gap. |
| **Template scaffolding** | `_template.md` (TODO placeholders for activation) | Authoring template; never read by implementers. |

---

## 2. Actionable findings (4 files)

### 2.1 `spec/03-error-manage/02-error-architecture/07-logging-and-diagnostics/02-session-based-logging.md`

Three feature requirements (F6, F8, F13) tagged `SHOULD` in a table where peer rows are `MUST`:

| Line | Req | Current | Recommendation |
| ---: | --- | :-----: | --- |
| 27   | F6 — sessions filterable by method/path/status | SHOULD | **MUST** (peer F1–F5 are MUST; filtering is core UX) |
| 29   | F8 — sessions auto-expire after retention period | SHOULD | **MUST** (storage-bound; absence = unbounded growth) |
| 34   | F13 — delegated server stack traces captured | SHOULD | **MUST** (debuggability invariant for distributed flows) |

**Blind-AI risk:** A literal reader will skip these features entirely. **Severity: Low** (the surrounding spec implies them), but cheap to harden.

### 2.2 `spec/14-update/27-generic-installer-behavior.md`

| Line | Issue | Recommendation |
| ---: | --- | --- |
| 151  | "Implementations SHOULD offer flags to disable individual sources" — flag set unspecified | Pin the flag names: `--no-github`, `--no-mirror`, `--no-local`. Promote to MUST. |
| 182  | "SHOULD adopt 20 going forward; 5 remains acceptable for legacy" — split rule | Make MUST=20 for new installers; explicit MAY=5 only for installers tagged `legacy`. |

### 2.3 `spec/14-update/23-install-script-version-probe.md`

| Line | Issue | Recommendation |
| ---: | --- | --- |
| 405  | "SHOULD log identity + range + outcome on stdout" | Promote to MUST — log surface is the only debuggability hook for headless installs. |

### 2.4 `spec/15-distribution-and-runner/06-fix-repo-forwarding.md`

| Line | Issue | Recommendation |
| ---: | --- | --- |
| 75   | "The dispatch path SHOULD be a single …" — sentence is also truncated mid-rule | Restore the full rule and promote to MUST; truncation itself is a doc bug. |

### 2.5 Trivial cosmetic (no risk, optional)

`spec/06-seedable-config-architecture/02-features/04-rag-test-coverage-matrix.md:251` — Mutation Score target shown as `TBD`. Either fix a number or strike the row; non-blocking.

---

## 3. Summary scoring (mirrors audit-09 dimensions)

| Dimension          | Score | Notes |
| ------------------ | ----: | --- |
| Completeness       | 95    | Two installer flag sets need pinning. |
| Determinism        | 96    | Three SHOULD→MUST promotions in `03-error-manage`. |
| Consistency        | 98    | No cross-folder contradictions detected. |
| Testability        | 94    | One `TBD` mutation-score target. |
| Blind-buildability | 96    | One truncated rule in `15-§6-fix-repo-forwarding`. |
| **Overall**        | **96 / 100 (A)** | Matches spec/19's post-Phase-13.3 grade. |

---

## 4. Recommendation

The non-spec/19 corpus does **not** require a Phase-13-style hardening campaign. A **single targeted patch** addressing the 4 files above (≈ 8 line-edits total) would close the residual ambiguity. Suggested grouping:

- **Patch A** — `03-error-manage` SHOULD→MUST (F6, F8, F13). 3 lines.
- **Patch B** — `14-update` installer determinism (flag names + cap rule + log requirement). 3 lines across 2 files.
- **Patch C** — `15-distribution-and-runner/06` truncated rule restored. 1–2 lines.
- **Patch D (optional)** — `06-seedable-config` mutation-score `TBD` resolved. 1 line.

Total effort: ~1 loop. Recommend executing as Phase-14.0 (pre-implementation polish) before opening Phase-14 implementation work.

---

## 5. What was *not* found (negative results worth recording)

- **No schema contradictions** across `04/05/06/14`.
- **No endpoint duplication** (the issue that drove spec/19's audit-07).
- **No stale audits** outside `spec/19-…/audit/` (which are already STALE-banner'd).
- **No prose↔seed drift** in `06-seedable-config` — alias maps are current.
- **No naming convention drift** — PascalCase/PK rules consistently applied across `04` and folder examples.
