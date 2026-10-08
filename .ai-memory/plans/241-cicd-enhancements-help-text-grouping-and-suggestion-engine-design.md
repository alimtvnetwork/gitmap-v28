# Master Execution Plan: CI/CD Enhancements, Muse Help Text Semantic Clustering & Suggestion Engine Architecture

> **Plan Version:** 1.0.0  
> **Status:** Active  
> **Parent Task:** Task-241  
> **Associated Specs:**  
> - `02-spec/21-app/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/01-architecture-spec.md`  
> - `02-spec/21-app/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/02-component-and-suggestion-spec.md`  

---

## 1. Goal & Objectives

1. **CI/CD Pipeline & Runner Enhancements:**
   - Decouple desktop and OS commands (`os dock`, `os update`, registry operations) from physical environments using headless mock adapters (`RegistryAccessor`, `LinuxDesktopRunner`), enabling end-to-end CI validation in headless GitHub Actions via `xvfb-run` and `dbus-run-session`.
   - Implement staged-only pre-commit fast-gating (`03-ai-scripts/50-fastgate.py`), dropping local commit validation time from 25s to <1.5s by linting only changed files.
   - Cache Split-DB pipeline telemetry (`.gitmap/repodb/pipeline.db`) across GitHub Actions workflow runs to detect historical duration regressions (>50%) and flaky tests.
   - Introduce duration-weighted dynamic test matrix sharding driven by `.ai-memory/test-inventory.json`, balancing CI test execution across equal 60s parallel shards.
   - Closed-loop AI telemetry integration: stream `gitmap pe -t` step failures directly into `store.AiAnalysisSplitDB` with status `failed`, auto-triggering `gitmap ai train` upon successful remediation.
   - Go test precompilation (`go test -c`) in a single build gate job to eliminate duplicate Go compilation overhead across test shards.

2. **Muse's Help Text Enhancement & Code Reduction (~4,300 LOC eliminated):**
   - Consolidate GitMap's 24 sprawling command groups into 5 intuitive semantic clusters:
     1. *Core & Repo Operations:* `scan`, `list`, `clone`, `pull`, `status`, `reconcile`, `stash`, `wip`, `group`, `cd`.
     2. *Release & Commits:* `release`, `changelog`, `list-versions`, `commit`, `cpf`, `cpb`, `cpr`, `cpar`.
     3. *Fleet & Remote SSH:* `ssh`, `nodes`, `cluster`, `sc`, `deploy`, `vhost`, `nginx`.
     4. *AI & Automation:* `agy`, `agm`, `aum`, `ai`, `macro`, `pipeline` (`pl`, `pe`), `rerun`, `sug`.
     5. *System, OS & Developer Tooling:* `os`, `apps`, `install`, `uninstall`, `storage`, `clean`, `vscode`, `sync`.
   - Replace 36+ repetitive Go files (`*_help_menu.go`, ~3,560 LOC) with a data-driven markdown help parser (`termhelp.FromMarkdown(mdBytes)`), utilizing existing files in `cli/helptext/`.
   - Net elimination of approximately 4,300 lines of boilerplate code across `cli/`.

3. **Centralized Suggestion Type Object & Automated Engine Design:**
   - Define canonical `suggestion.Suggestion` and `suggestion.SuggestionGroup` domain types in `cli/suggestion/`.
   - Build a multi-tier matching engine:
     * Tier 1: Canonical Command & Alias Registry (`gitmap s` -> `status`, `gitmap pe` -> `pull-error`).
     * Tier 2: Levenshtein distance & QGram fuzzy matching (edit distance <= 3).
     * Tier 3: Subsequence / Prefix matcher.
     * Tier 4: Semantic Context & Synonym Dictionary (`docker` -> `cluster`, `rm` -> `clean`).
   - Pluggable renderers: Catppuccin terminal box (`RenderFormatBox`), compact inline (`RenderFormatCompact`), and machine-readable JSON (`RenderFormatJSON`).
   - Global dispatch interceptor in `cli/cmd/root.go:runDispatch` and `appfault.AppError` integration (`WithSuggestions(...)`), eliminating ad-hoc `fmt.Printf("Did you mean: ...")` prints throughout the codebase.

---

## 2. User Request (Verbatim)

```text
Please read the recent codes and recent changes for the `gitmap` and share what can we introduce or enhance in the CI/CD. And also Muse, another tool by Facebook, actually told me that we should be able to enhance the help text by reducing it or grouping it. Tell me more about it, what can we do, and how much code we can reduce in the code base. For example, the suggestions that we give, we should have a suggestion type object where we just put the suggestion, let's say command or things. There should be a system where we put this stuff and it would automatically put the suggestion print rather than doing it manually. So create a system design plan that we can improve inside the code base. First read what is pending as a code and what we have done recently, read and then understand what can we improve. Is it clear?
```

---

## 3. Subtask Breakdown

- **Subtask 01 (`01-cicd-pipeline-and-desktop-mock-enhancements.md`):** CI/CD Pipeline Architecture, Desktop Mock Adapters, Staged-Only Fast-Gate & Telemetry Caching.
- **Subtask 02 (`02-muse-help-text-clustering-and-code-reduction.md`):** Muse Help Text 5 Semantic Clusters, Markdown-Driven Rendering & ~4,300 LOC Code Reduction.
- **Subtask 03 (`03-suggestion-type-object-and-engine-design.md`):** Suggestion Domain Model, Multi-Tier Matcher, and Pluggable Catppuccin Box Renderers.
- **Subtask 04 (`04-error-interceptor-and-cli-integration.md`):** Global CLI Dispatch Interception, AppError Wrapper Integration & Automated Failure Suggestions.

---

## 4. Execution Tracking Ledger

| Subtask ID | File | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| Subtask-01 | `01-cicd-pipeline-and-desktop-mock-enhancements.md` | Spec Author 01 | PENDING | Documenting 6 CI/CD enhancement blueprints, mock adapters, and pipeline telemetry |
| Subtask-02 | `02-muse-help-text-clustering-and-code-reduction.md` | Spec Author 01 | PENDING | Formulating 5 semantic clusters, data-driven Markdown migration, and ~4,300 LOC reduction |
| Subtask-03 | `03-suggestion-type-object-and-engine-design.md` | Spec Author 02 | PENDING | Specifying `suggestion.Suggestion` struct, multi-tier matcher, and Catppuccin box renderer |
| Subtask-04 | `04-error-interceptor-and-cli-integration.md` | Spec Author 02 | PENDING | Hooking `root.go` dispatch, `AppError` integration, and eliminating ad-hoc print calls |

---

## 5. Non-Negotiable Operational Invariants

1. **Relative Paths Exclusively:** Every path reference in specifications, plans, and code uses strictly relative repository paths (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero absolute paths or drive letters.
2. **GitMap Search Primacy:** Code and symbol discovery strictly through `gitmap aum search`, `gitmap find`, `gitmap lf`, and `gitmap cat`. Zero usage of `rg`, `ripgrep`, `grep`, `git grep`, or `Select-String`.
3. **Coding Guidelines Adherence:** Positive booleans with `is/has` prefixes, zero double negatives, maximum function length of 15 lines, and zero nested `if` statements (nesting depth <= 1).
4. **Task & Budget Integrity:** All subtasks managed through GitMap's SQLite task subsystem (`gitmap task`).
