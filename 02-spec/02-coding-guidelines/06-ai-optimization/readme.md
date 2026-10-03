# AI Optimization (AI Execution Prompt)

> **/goal** Provide the primary entry point and high-authority index for all AI optimization specifications, preventing hallucinations and ensuring machine-generated code complies with project conventions.
> **/learn** Master the core hierarchy of anti-hallucination rules, pre-output validation checklists, mistake catalogs, and condensed master guidelines for LLM context windows.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Consult and enforce the AI optimization catalog before generating or refactoring codebase logic.
- [ ] `/learn` Verify pre-output compliance using the quick-reference checklist and anti-hallucination constraints.
- [ ] `/goal` Guarantee zero generated code artifacts, test logs, or build binaries are staged into version control.
- [ ] `/learn` Validate AI optimization specifications using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

AI-specific guidelines designed to prevent hallucination and ensure AI-generated code meets all project standards. Contains explicit forbidden patterns, quick validation checklists, and common mistake catalogs.

---

## Files

| # | File | Description |
|---|------|-------------|
| 01 | [01-anti-hallucination-rules.md](./02-anti-hallucination-rules.md) | 30+ explicit "never generate X" rules with forbidden/required patterns |
| 02 | [02-ai-quick-reference-checklist.md](./03-ai-quick-reference-checklist.md) | 50-check pre-output validation checklist |
| 03 | [03-common-ai-mistakes.md](./04-common-ai-mistakes.md) | Top 15 real mistakes AI makes, with before/after corrections |
| 04 | [04-condensed-master-guidelines.md](./05-condensed-master-guidelines.md) | Sub-200-line distillation of master guidelines for AI context windows |
| 05 | [05-enum-naming-quick-reference.md](./07-enum-naming-quick-reference.md) | Cross-language enum naming rules: Go, TypeScript, PHP — declaration, naming, usage, validation checklist |

---

## How AI Should Use This Section

1. **Before generating code:** Scan the quick-reference checklist (02)
2. **During generation:** Apply anti-hallucination rules (01) to each code block
3. **After generation:** Verify against common mistakes (03)
4. **For enum code:** Consult the enum naming quick reference (05)

---

*AI optimization overview for coding guidelines.*

## Document Inventory

| File |
|------|
| 97-acceptance-criteria.md |
| 99-consistency-report.md |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-ai-optimization/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-AI-001: AI Optimization Index & Architecture Conformance

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.
