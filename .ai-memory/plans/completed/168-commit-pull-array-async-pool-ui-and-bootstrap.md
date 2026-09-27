# Plan 168: Commit-Pull Array Async Pool, Interactive Web UI & Declarative Bootstrap

- **Status:** `completed`
- **Spec Reference:** [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)

---

## 1. Context & Architectural Overview

The multi-repository migration engine (`commit-pull` / `commit-in`) consolidates repository histories into a single target repository. During dry-run and input staging, sequential repository probing creates unnecessary latency.
This plan implements the **Array Async Pool Concept by Alim Ul Karim**, a lock-free parallel indexing and sequential terminal streaming pattern, introduces an interactive browser Web UI (`gitmap commit-pull ui`), creates declarative configuration scaffolding (`gitmap commit-pull bootstrap`), normalizes positive boolean conventions, enforces zero-allocation string prefix matching over regex, and embeds self-contained variables in SEO templates.

---

## 2. Deliverables Mapping

| Subtask | ID | Target Files | Description |
|---|---|---|---|
| `01-coding-guideline-array-async-pool.md` | Task-02 | `02-spec/02-coding-guidelines/14-array-async-pool.md`, `02-spec/02-coding-guidelines/readme.md`, `.ai-memory/coding-guidelines.md` | Author authoritative Coding Guideline 14 on Array Async Pool Concept by Alim Ul Karim with timing data comparison |
| `02-commit-pull-array-async-pool.md` | Task-03 | `cli/cmd/commitin/workspace/clone.go`, `cli/cmd/commitin/orchestrator/pipeline.go`, `cli/cmd/commitin/orchestrator/dryrun.go` | Implement pre-allocated Array Async Pool for concurrent repository probing and ordered ticker streaming |
| `03-booleans-and-fast-prefix-matching.md` | Task-04 | `cli/cmd/commitin/lineskipper/line_skipper.go`, `cli/cmd/commitin/parse_types.go`, `cli/cmd/commitin/config_json.go` | Normalize booleans (`is...`/`has...`) and prioritize zero-allocation `strings.HasPrefix` over slow regex in line skippers |
| `04-short-inputs-flags-and-bootstrap.md` | Task-05 | `cli/cmd/commitin.go`, `cli/cmd/commitin/config_json.go`, `cli/helptext/commit-pull.md`, `readme.md` | Add comma-separated short input syntax, terminal flags, `gitmap commit-pull bootstrap`, and enriched help text |
| `05-commit-pull-web-ui-and-seo-variables.md` | Task-06 | `cli/cmd/commitin_ui_server.go`, `cli/cmd/commitin.go`, `.ai-memory/temp/seo-templates.json`, `cli/store/templates_split_ops.go` | Implement dark-mode single-page Web Studio (`gitmap commit-pull ui`) and self-contained variables in SEO templates |
| `06-tempe2e-tests-and-release.md` | Task-07 | `cli/tests/e2e/test_gitmap_migration_tempe2e_test.go`, `cli/constants/constants.go` | Update isolated tempe2e tests, execute minor release bump `v6.348.0`, and verify 100% green CI/CD |
