# PHP Standards Reference (AI Execution Prompt)

> **/goal** Enforce and navigate the comprehensive PHP coding standards reference across naming, enums, constants, initialization, booleans, code style, and database patterns.
> **/learn** Master the decomposed modular structure under `07-php-standards-reference/` to ensure zero monolithic code files, 100% adherence to modular standards, and clean cross-language consistency.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Direct all comprehensive PHP standards audits to the modular specifications in `07-php-standards-reference/`.
- [ ] `/learn` Enforce naming conventions, structured error envelopes, and result responses from `02-naming-and-errors.md`.
- [ ] `/learn` Enforce centralized identifiers in `constants.php` and native backed enums in `03-constants-and-deps.md`.
- [ ] `/learn` Adhere to constructor rules, positive booleans, and `isDefined` guards from `04-initialization-and-booleans.md`.
- [ ] `/learn` Enforce PSR-12 braces, vertical spacing, early returns, and function size limits from `05-code-style.md`.
- [ ] `/learn` Enforce total ban on forbidden globals, raw SQL concatenation, and unescaped queries from `06-forbidden-and-database.md`.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Modular Standards Decomposition

> ⚠️ **This specification has been modularized.** All content now lives in [`07-php-standards-reference/`](./07-php-standards-reference/).

| Document | Focus Area | Key Requirements |
|:---|:---|:---|
| [`02-naming-and-errors.md`](./07-php-standards-reference/02-naming-and-errors.md) | Naming & Error Handling | PascalCase classes/enums, camelCase methods/variables, `ResultHelper` envelopes |
| [`03-constants-and-deps.md`](./07-php-standards-reference/03-constants-and-deps.md) | Constants & Dependencies | Centralized `constants.php`, backed enums in `includes/Enums/`, path constants |
| [`04-initialization-and-booleans.md`](./07-php-standards-reference/04-initialization-and-booleans.md) | Constructors & Booleans | Explicit property initialization, positive boolean prefixes (`is`/`has`), `isDefined` guards |
| [`05-code-style.md`](./07-php-standards-reference/05-code-style.md) | Formatting & Code Style | PSR-12 braces, blank line before `if`/`return`, early return flattening, 8–15 line function cap |
| [`06-forbidden-and-database.md`](./07-php-standards-reference/06-forbidden-and-database.md) | Prohibited Patterns & DB | Total ban on `eval`, superglobals directly, unescaped queries; `$wpdb->prepare` required |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PHP-006: PHP Standards Reference Conformance

**Given** PHP codebases and specification references within the repository.
**When** Code and documentation are validated against the modular PHP standards reference.
**Then** All implementation files conform to the decomposed standards in `07-php-standards-reference/` with zero monolithic violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.
