# Acceptance Criteria: Static Analysis & Linter Enforcement (AI Execution Prompt)

> **/goal** Maintain and enforce centralized acceptance criteria for multi-language static analysis, linter configurations, and unified CI/CD quality gates across all 8 supported stacks.
> **/learn** Standardize thresholds (15-line functions, 3 params, complexity ≤ 10), SonarQube rule mappings, and automated lint verification across polyglot builds.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify every supported language has active static analysis rules matching canonical thresholds.
- [ ] `/learn` Map SonarQube rule IDs for every enforceable rule and ban unverified suppression comments.
- [ ] `/goal` Ensure unified CI/CD pipeline definitions fail builds on quality gate threshold breaches.
- [ ] `/learn` Verify zero broken references across language-specific linter documentation and matrix tables.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-10-02

---

## Must Have

- [ ] Every supported language (8) has a dedicated linter spec
- [ ] All specs enforce identical thresholds: 15-line functions, 3 params, complexity ≤ 10
- [ ] SonarQube rule IDs mapped for every enforceable rule per language
- [ ] Integration checklist uses standardized table format with 🔲 status
- [ ] CI pipeline spec defines a unified quality gate referencing all 8 languages
- [ ] Cross-language rule matrix covers all 7 SonarQube rules across all 8 languages
- [ ] TypeScript ESLint spec cross-referenced from overview (lives in `02-typescript/`)

## Should Have

- [ ] Each spec includes Keywords and Scoring sections
- [ ] Exemption/suppression syntax documented per language
- [ ] GitHub Actions and GitLab CI templates provided in CI spec

## Won't Have (This Version)

- Pre-commit hook configurations (deferred to future iteration)
- IDE-specific settings files (`.vscode/`, `.idea/`)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-STATIC-001: Static Analysis Multi-Language Conformance Registry

**Given** The static analysis specifications and rule mapping matrix in `01-cross-language/16-static-analysis/`.
**When** Automated linters audit rule coverage across all 8 supported languages.
**Then** All language specs enforce identical canonical size thresholds with zero unmapped rules.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/16-static-analysis --check-only
```
**Expected:** exit 0. Zero violations.
