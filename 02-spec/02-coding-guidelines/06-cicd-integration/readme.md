# CI/CD Integration — Coding-Guidelines Linter Pack (AI Execution Prompt)

> **/goal** Master and deploy the portable, language-agnostic CI/CD linter pack (`linters-cicd/`) to enforce repository CODE RED rules across any CI engine with zero external dependencies.
> **/learn** Internalize SARIF 2.1.0 reporting contracts, modular language plugin architectures, one-line CI platform templates, and strict exit code semantics.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Integrate portable `linters-cicd/` checks emitting compliant SARIF 2.1.0 reports for automated pull request annotations.
- [ ] `/learn` Never require local repository dependencies or external package installation for Phase 1 check execution.
- [ ] `/goal` Maintain isolated language plugins adhering to standard CLI flags (`--path`, `--format`, `--severity`).
- [ ] `/learn` Validate compliance across the directory using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0
> **Status:** Active (Phase 1 shipping)
> **Updated:** 2026-04-19
> **AI Confidence:** Production-Ready
> **Ambiguity:** None

---

## Keywords

`ci`, `cd`, `linter`, `sarif`, `github-actions`, `gitlab`, `azure-devops`,
`jenkins`, `bitbucket`, `coding-guidelines`, `code-red`

---

## Purpose

Ship a portable, language-agnostic linter pack — `linters-cicd/` — that any
CI/CD pipeline can integrate with one line to enforce the **CODE RED**
coding-guidelines rules in this repository.

The pack must:

1. Run on any CI without requiring this repo as a dependency.
2. Emit **SARIF 2.1.0** by default so GitHub, GitLab, and Azure DevOps
   render findings inline on pull requests.
3. Support a **plugin model** so adding a new language is one file plus a
   test fixture — never a refactor.
4. Be installable three ways: ZIP one-liner, GitHub composite Action, or
   manual checkout.

---

## Document Inventory

| File | Purpose |
|------|---------|
| `readme.md` | This file — high-level architecture |
| `02-sarif-contract.md` | SARIF 2.1.0 output schema each check must emit |
| `03-plugin-model.md` | How a new-language check is added |
| `04-language-roadmap.md` | Phase 1/2/3 rollout, current status per language |
| `05-ci-templates.md` | Inventory of CI platform templates |
| `06-distribution.md` | ZIP, composite Action, install.sh contract |
| `07-rules-mapping.md` | Each rule → spec source → check script → severity |
| `08-performance.md` | Middle-out probe order, parallel jobs, timeout budgets |
| [`08-fix-repo-and-installers/`](./08-fix-repo-and-installers/readme.md) | Canonical contract for `install.*`, `release-install.*`, `fix-repo.*`, `visibility-change.*` |
| `97-acceptance-criteria.md` | Testable AC for every public artifact |
| `98-faq.md` | Suppression, baselining, single-rule runs, version pinning |
| `99-troubleshooting.md` | python3 missing, tree-sitter wheel issues, SARIF >10 MB, false-positive triage, TOML parse errors |

---

## Architecture (3 layers)

### Layer 1 — Portable check scripts (`linters-cicd/checks/`)

Pure POSIX shell + Python 3 (no project-local paths, no external installs).
Each script:

- Accepts `--path <dir>`, `--format sarif|text`, `--severity error|warning`
- Exits `0` on no findings, `1` on findings, `2` on tool failure
- Writes SARIF to stdout or `--output <file>`

### Layer 2 — Pre-built linter configs (`linters-cicd/configs/`)

Drop-in configs that wrap the existing community linters where they
already cover a CODE RED rule (golangci-lint, ESLint, PHPCS).

### Layer 3 — CI templates (`linters-cicd/ci/`)

One ready-to-paste workflow per platform plus the GitHub composite Action
under `linters-cicd/action.yml`.

---

## Out of scope (v1)

- Auto-fix / codemod (read-only checks only)
- IDE plugins (already covered by `.cursorrules` + linter configs)
- Custom rule authoring DSL — rules are coded in Python, period

---

## Cross-References

- [Code Red Guidelines](../../17-consolidated-guidelines/readme.md)
- [CI/CD Pipeline Workflows](../../12-cicd-pipeline-workflows/readme.md)
- [Linter Scripts](../../../linter-scripts/) — existing in-repo guards

---

## Contributors

- **Md. Alim Ul Karim** — Creator & Lead Architect
- **Riseup Asia LLC** — Sponsor

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CI-001: CI/CD Integration Architecture Conformance

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.
