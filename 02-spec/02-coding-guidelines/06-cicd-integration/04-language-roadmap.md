# Language Roadmap (AI Execution Prompt)

> **/goal** Track, plan, and execute the phased rollout of language-specific linter plugins across Go, TypeScript, PHP, Python, and Rust.
> **/learn** Understand multi-phase promotion criteria, test fixture obligations, AST vs regex trade-offs, and rule parity invariants.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify Phase 1 language support (Go, TypeScript, PHP) across all 7 core CODE RED checks.
- [ ] `/learn` Never promote a planned language to shipping without corresponding positive and negative test fixtures.
- [ ] `/goal` Ensure language expansion maintains zero modifications to the core orchestrator or SARIF schemas.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0
> **Updated:** 2026-04-19

Phased rollout per user direction. Each phase is independently shippable
and adds **only** new plugins under `linters-cicd/checks/<rule>/<lang>.py`
plus registry entries.

---

## Phase 1 — Go + TypeScript + PHP (✅ shipping in v3.9.0)

All three currently-specced languages ship together so consumers with
mixed Go/TS/PHP repos (the common Riseup Asia stack) get full coverage
on day one — no waiting for a follow-up release.

| Check | Go | TypeScript | PHP |
|-------|----|------------|-----|
| nested-if (CODE-RED-001) | ✅ regex+AST hybrid | ✅ regex+AST hybrid | ✅ regex+`phply` hybrid |
| function-length (CODE-RED-004) | ✅ | ✅ | ✅ |
| file-length (CODE-RED-006) | ✅ | ✅ | ✅ |
| magic-strings (CODE-RED-003) | ✅ | ✅ | ✅ (WordPress-aware allowlist) |
| boolean-naming (CODE-RED-002) | ✅ | ✅ | ✅ |
| positive-conditions (CODE-RED-008) | ✅ | ✅ | ✅ |
| no-else-after-return (STYLE-002) | ✅ | ✅ | ✅ |

**Why these three first:**

- **Go + TS** are the languages used in this repo, so the checks self-test
  against the spec's own corpus.
- **PHP** rules already exist in `02-spec/02-coding-guidelines/04-php/` and
  the WordPress/Laravel projects under Riseup Asia LLC need parity with
  the Go/TS enforcement on day one. PHP plugins use `phply` for AST and
  regex fallbacks for legacy WordPress code that fails to parse.

---

## Phase 2 — Python + Rust (planned)

Python uses the standard library `ast` module — trivial. Rust uses
`tree-sitter-rust` Python bindings.

---

## Phase 4+ — On request

Any additional language (Java, Kotlin, Swift, C#, …) is added on user
request following [`02-plugin-model.md`](./03-plugin-model.md). The
orchestrator and SARIF contract remain unchanged.

---

## Promotion criteria (todo → shipping)

A language graduates from "planned" to "shipping" when:

1. All 7 Phase 1 checks have a working plugin with fixtures.
2. `linters-cicd/checks/<rule>/fixtures/<lang>/` has ≥ 1 bad and ≥ 1 good
   fixture per check.
3. CI runs `validate-sarif.py` on every emission.
4. The rule appears in `06-rules-mapping.md` with status `shipping`.

---

*Part of [CI/CD Integration](./readme.md)*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CI-004: Language Rollout Roadmap and Promotion Standards

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.
