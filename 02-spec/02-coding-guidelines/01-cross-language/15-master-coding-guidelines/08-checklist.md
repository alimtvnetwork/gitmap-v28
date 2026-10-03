# Master Coding Guidelines — Quick checklist for any code change (AI Execution Prompt)

> **/goal** Provide an authoritative pre-commit checklist synthesizing all universal coding standards across naming, booleans, enums, database, styling, errors, single return values, and architecture.
> **/learn** Verify every proposed code modification against non-negotiable cross-language rules before submitting pull requests or completing agent tasks.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Validate naming conventions (camelCase variables, PascalCase classes/enums/DB columns, single-capitalized abbreviations).
- [ ] `/learn` Enforce positive boolean naming (`is`/`has`), implicit truthiness evaluation, and zero mixed-polarity conditions.
- [ ] `/goal` Ensure Go functions return a single `Result[T]` or `*appfault.AppError`, and eliminate all runtime type casting in business logic.
- [ ] `/learn` Enforce vertical whitespace rules, mandatory braces, zero nested `if` statements, and single defer hygiene.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Master Coding Guidelines](./readme.md)
> **Version:** 2.1.0
> **Updated:** 2026-03-31

---

## Quick Checklist for Any Code Change

```
[ ] Naming: camelCase variables, PascalCase classes/enums/DB columns
[ ] JSON/API keys: PascalCase (e.g., "PluginSlug", "SiteId" — never "SITE_ID" or "siteId")
[ ] Abbreviations: Id (not ID), Url (not URL), Md5 (not MD5), Json (not JSON), Api (not API)
[ ] Null guards: isDefined()/isDefinedAndValid() — never raw != null/nil + isValid()
[ ] Null safety: check err before value, nil before dereference, len before index
[ ] Booleans: is/has prefix, no negative words, no raw ! on calls
[ ] Enums: Type suffix, isEqual() not ===, PascalCase case names
[ ] DB: PascalCase tables/columns, PascalCase array keys for inserts
[ ] Formatting: braces always, zero nesting, blank before return, 15-line max
[ ] Nesting: zero nested if — use early returns, named booleans, extracted functions
[ ] Newlines: blank before return (multi-line), no double blanks, no blank at function start
[ ] Errors: apperror.Wrap (Go), Throwable imported (PHP), no fmt.Errorf
[ ] Results: hasError()/isSafe() checked before .value()/.Value() — use .AppError() in Go (not .Error())
[ ] Single return: Go functions return ONE value (Result[T] or typed struct) — never (T, bool, error)
[ ] No casting: zero type assertions in Go business logic — use concrete structs
[ ] No magic strings: all via enums/typed constants
[ ] Mutation: assign once, no post-construction mutation, mutex for concurrent state
[ ] Regex: last resort, compile at package level (Go), no regex in loops
[ ] Lazy eval: non-exported field + getter, mutex for concurrent, cascade lazy dependencies
[ ] Defer: max one per function (Go), top or bottom placement only
[ ] Log keys: camelCase in PHP context arrays
[ ] Types: no any/interface{} (Go), native types + no redundant PHPDoc (PHP)
[ ] Tests: three-part naming, AAA pattern, table-driven for 3+ cases, t.Helper() in Go helpers
```

---

*Master coding guidelines v2.0.0 — 2026-03-31*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-MASTER-CHECKLIST: Comprehensive Pre-Commit Verification Checklist

**Given** Any proposed code changes across PHP, Go, TypeScript, and SQL codebases.
**When** Pre-commit verification checks and CI autofixers evaluate code conformance against master standards.
**Then** 100% of checklist items pass with zero violations, clean newline formatting, and full compliance across all architectural rules.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
