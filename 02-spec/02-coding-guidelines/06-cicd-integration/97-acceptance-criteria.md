# CI/CD Integration — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all CI/CD integration specifications in `06-cicd-integration/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-CICD-[NUM]` / `AC-CG-CI-[NUM]`), SARIF 2.1.0 output schemas, portable plugin contracts, multi-platform CI templates, versioned release distribution, rule mapping tiers, and performance budgets.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `06-cicd-integration/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths across all guideline links.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. CI/CD Integration Criteria Inventory (`AC-CG-CICD-` / `AC-CG-CI-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-CICD-001` (alias `AC-CG-CI-001`) | CI/CD Integration Architecture Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-002` (alias `AC-CG-CI-002`) | SARIF 2.1.0 Contract Conformance and Validation | [`02-sarif-contract.md`](02-sarif-contract.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-003` (alias `AC-CG-CI-003`) | Language Plugin Architecture and Registry Standards | [`03-plugin-model.md`](03-plugin-model.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-004` (alias `AC-CG-CI-004`) | Language Rollout Roadmap and Promotion Standards | [`04-language-roadmap.md`](04-language-roadmap.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-005` (alias `AC-CG-CI-005`) | Cross-Platform CI Templates and Invocation Standards | [`05-ci-templates.md`](05-ci-templates.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-006` (alias `AC-CG-CI-006`) | Distribution Packaging and Release Asset Governance | [`06-distribution.md`](06-distribution.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-007` (alias `AC-CG-CI-007`) | Rule Taxonomy, Severity Mapping, and Tier Coordination | [`07-rules-mapping.md`](07-rules-mapping.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-008` (alias `AC-CG-CI-008`) | Probe Ordering, Parallelism, and Timeout Budgets | [`08-performance.md`](08-performance.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-FAQ-001` | Linter Pack FAQ & Consumer Operations Conformance | [`98-faq.md`](98-faq.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-TROUBLE-001` | Linter Pack Operations Troubleshooting Conformance | [`99-troubleshooting.md`](99-troubleshooting.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| `AC-CG-CICD-REG-001` | CI/CD Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-CICD-001 / AC-CG-CI-001: CI/CD Integration Architecture Conformance

- [ ] Portable check scripts run under stock POSIX shell + Python 3 with zero external dependencies.
- [ ] Root `readme.md` provides complete module navigation, execution prompts, and cross-references.
- [ ] Emits SARIF 2.1.0 by default to render inline findings across GitHub, GitLab, and Azure DevOps.
- [ ] Standard exit code contract: `0` for clean, `1` for rule findings, and `2` for tool failure.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-002 / AC-CG-CI-002: SARIF 2.1.0 Contract Conformance and Validation

- [ ] Emitted SARIF complies with official schema `https://json.schemastore.org/sarif-2.1.0.json`.
- [ ] All `artifactLocation.uri` paths are strictly relative to scanned repository root.
- [ ] CODE RED rules map to `error`, STYLE rules map to `warning`, and informational rules map to `note`.
- [ ] Multiple check outputs merge cleanly into a single valid SARIF document.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-003 / AC-CG-CI-003: Language Plugin Architecture and Registry Standards

- [ ] Each language plugin resides under `linters-cicd/checks/<rule>/<language>.py` and registers in `registry.json`.
- [ ] Adding a new language requires zero modifications to `run-all.sh` or `action.yml`.
- [ ] Plugins accept standard CLI arguments `--path <dir>`, `--format sarif|text`, and `--severity error|warning`.
- [ ] Every plugin provides paired positive and negative test fixtures under `fixtures/<language>/`.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-004 / AC-CG-CI-004: Language Rollout Roadmap and Promotion Standards

- [ ] Phase 1 covers Go, TypeScript, and PHP across all 7 core CODE RED checks.
- [ ] Language plugins use native AST walkers (Python stdlib `ast`, Tree-sitter, or `phply`) with regex fallbacks.
- [ ] Promotion from planned to shipping requires 100% fixture coverage and green `validate-sarif.py` runs.
- [ ] Subsequent phases (Python, Rust, C#) maintain exact orchestrator and SARIF contract compatibility.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-005 / AC-CG-CI-005: Cross-Platform CI Templates and Invocation Standards

- [ ] Ready-to-paste workflow templates provided for GitHub Actions, GitLab CI, Azure DevOps, Bitbucket, and Jenkins.
- [ ] Single-line integration supported via GitHub composite Action `linters-cicd/action.yml`.
- [ ] Workflows pass explicit target `--languages` to prevent auto-detection latency or errors.
- [ ] Pull request checks fail on non-zero exit codes and display annotations inline.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-006 / AC-CG-CI-006: Distribution Packaging and Release Asset Governance

- [ ] Universal standalone ZIP `coding-guidelines-linters-vX.Y.Z.zip` attached to every release tag.
- [ ] SHA-256 checksums recorded in `checksums.txt` and verified by `install.sh`.
- [ ] Release pipeline strictly adheres to zero GitHub Actions storage upload limits.
- [ ] Composite action references pinned version tags (`@v3.9.0`) for reproducible CI runs.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-007 / AC-CG-CI-007: Rule Taxonomy, Severity Mapping, and Tier Coordination

- [ ] Canonical mapping of every rule ID to spec source, check script, supported languages, and SARIF severity.
- [ ] Coordinated function length tier: CODE-RED-005 (strict-8 error) owns build failure, CODE-RED-004 (hard-15 error) acts as defense-in-depth.
- [ ] Database rules enforce non-negative boolean naming and required description/notes columns.
- [ ] Rule removals or breaking modifications require a major version bump and deprecation notice.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-008 / AC-CG-CI-008: Probe Ordering, Parallelism, and Timeout Budgets

- [ ] Middle-out probe ordering sorts candidate directories by byte weight to surface findings early.
- [ ] Parallel check execution amortizes interpreter startup across (rule, language) tuples.
- [ ] Hard timeout budgets enforced: 20s per check, 120s total run, and 2s per file parse.
- [ ] Run completes in < 30 seconds wall-clock on typical 50 kLOC repository runners.

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-FAQ-001: Linter Pack FAQ & Consumer Operations Conformance

- [ ] Inline suppression syntax mandates rule ID and reason text (`// codeguidelines:disable=RULE — reason`).
- [ ] Suppressions lacking justification trigger synthetic warning finding `STYLE-099`.
- [ ] Baseline workflows support `--baseline` and `--refresh-baseline` flags to isolate legacy violations.
- [ ] Filter flags (`--rules`, `--languages`, `--exclude-rules`) and configuration files (`.codeguidelines.toml`) operate deterministically.
- [ ] Version pinning instructions mandate exact semantic versions or commit SHAs with zero floating tags (`@latest`, `@main`).

**Given** CI/CD pipeline infrastructure, configuration files, and linter consumer queries.
**When** Audited against this FAQ specification.
**Then** All suppression patterns, baselining steps, and version pins conform to documented standards with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-TROUBLE-001: Linter Pack Operations Troubleshooting Conformance

- [ ] Python 3 runtime requirements (>= 3.10) and setup actions (`actions/setup-python@v5`) are documented across platforms.
- [ ] Native C dependency requirements for AST parsers (tree-sitter, python3-dev, build-essential) are explicitly detailed.
- [ ] SARIF schema validation and upload troubleshooting cover GitHub, GitLab, and Azure DevOps integration.
- [ ] Exit code semantics (0=pass, 1=violations, 2=configuration/runtime error) are verified with reproducible remediation steps.
- [ ] Configuration parsing (`.codeguidelines.toml`) errors provide actionable validation commands.

**Given** CI/CD pipeline infrastructure, runner execution environments, and diagnostic incident reports.
**When** Audited against this troubleshooting specification.
**Then** All operational failure modes have actionable diagnosis and remediation steps with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CICD-REG-001: CI/CD Acceptance Criteria Registry Conformance

- [ ] Master registry at `02-spec/02-coding-guidelines/97-acceptance-criteria.md` references the CI/CD Integration registry under section `AC-06B`.
- [ ] Every criterion in the CI/CD Integration registry maps 1:1 to an authoritative specification in `06-cicd-integration/`.
- [ ] Both master `AC-CG-CICD-*` taxonomy and local `AC-CG-CI-*` aliases are tracked and harmonized with zero divergence.
- [ ] All criteria follow the structured `Given / When / Then` contract.
- [ ] Verification commands execute cleanly with exit code 0.
- [ ] Strict relative paths are used throughout with zero absolute paths or `file:///` URIs.

**Given** The CI/CD Integration criteria registry at `06-cicd-integration/97-acceptance-criteria.md`.
**When** Evaluated for taxonomic integrity and traceability.
**Then** Every entry links 1:1 to authoritative specifications in `06-cicd-integration/`, features explicit `Given/When/Then` definitions, and defines functional verification commands with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.

---

*Part of [CI/CD Integration](./readme.md)*
