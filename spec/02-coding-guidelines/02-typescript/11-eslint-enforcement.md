# TypeScript ESLint Enforcement — Rule Mapping

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`typescript` · `eslint` · `static-analysis` · `sonarjs` · `typescript-eslint` · `custom-plugin`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Maps every TypeScript coding guideline rule to its ESLint enforcement. Includes custom plugin rules (in `eslint-plugins/coding-guidelines/`) and recommended third-party rules. Any team adopting these guidelines can copy the ESLint config and get instant enforcement.

---

## 1. Custom Plugin Rules — `coding-guidelines/*`

Source: `eslint-plugins/coding-guidelines/index.js`

| Rule | Severity | Spec Source | Description |
|------|----------|-------------|-------------|
| `coding-guidelines/no-nested-if` | `error` | [Code Style §R2](../01-cross-language/04-code-style/01-braces-and-nesting.md) | Zero nested `if` — flatten with early returns |
| `coding-guidelines/boolean-naming` | `error` | [Boolean Principles](../01-cross-language/02-boolean-principles/00-overview.md) | Boolean vars must use `is/has/can/should/was/will` prefix |
| `coding-guidelines/no-magic-strings` | `warn` | [TS Standards §3](./08-typescript-standards-reference.md) | No raw string literals in `===`/`!==`/`switch` — use enum/constant |
| `coding-guidelines/max-function-lines` | `error` | [Code Style §R6](../01-cross-language/04-code-style/04-function-and-type-size.md) | Max 15 lines per function body (non-blank, non-comment) |
| `coding-guidelines/promise-all-independent` | `error` | [Promise Patterns §3](./09-promise-await-patterns.md) | Sequential `await` on independent promises → use `Promise.all` |
| `coding-guidelines/blank-line-before-return` | `warn` | [Code Style §R4](../01-cross-language/04-code-style/03-blank-lines-and-spacing.md) | Blank line before `return`/`throw` when preceded by statements |
| `coding-guidelines/no-else-after-return` | `error` | [Code Style §R7](../01-cross-language/04-code-style/01-braces-and-nesting.md) | No `else` after `return`/`throw`/`continue`/`break` |

---

## 2. Recommended `@typescript-eslint` Rules

These enforce type safety rules from [TS Standards §1–§2](./08-typescript-standards-reference.md):

| Rule | Severity | Spec Source | Description |
|------|----------|-------------|-------------|
| `@typescript-eslint/no-explicit-any` | `error` | TS Standards §2.1 | `any` is prohibited everywhere |
| `@typescript-eslint/no-unsafe-assignment` | `error` | TS Standards §2.1 | No assigning `any`-typed values |
| `@typescript-eslint/no-unsafe-member-access` | `error` | TS Standards §2.1 | No member access on `any`-typed values |
| `@typescript-eslint/no-unsafe-call` | `error` | TS Standards §2.1 | No calling `any`-typed values |
| `@typescript-eslint/no-unsafe-return` | `error` | TS Standards §2.1 | No returning `any`-typed values |

> **Note:** The `no-unsafe-*` rules require `parserOptions.project` pointing to `tsconfig.json` for type-aware linting.

---

## 3. SonarQube Rule Mapping

SonarQube rules that overlap with our coding guidelines. Use [SonarQube for TypeScript](https://rules.sonarsource.com/typescript/) or the `eslint-plugin-sonarjs` ESLint plugin.

### Install

```bash
npm install --save-dev eslint-plugin-sonarjs
```

### Rule Mapping

| SonarQube Rule | ESLint Equivalent | Our Guideline |
|----------------|-------------------|---------------|
| S3776 — Cognitive Complexity | `sonarjs/cognitive-complexity` | Max 15-line functions + zero nesting |
| S1871 — Identical branches | `sonarjs/no-identical-functions` | DRY principles |
| S1126 — Return boolean directly | `sonarjs/prefer-immediate-return` | Boolean naming + no negations |
| S3923 — All branches identical | `sonarjs/no-all-duplicated-branches` | Code review enforcement |
| S1066 — Collapsible if | `sonarjs/no-collapsible-if` | Zero nesting rule |
| S1125 — Boolean literal in comparison | `sonarjs/no-gratuitous-expressions` | Boolean principles P1–P6 |
| S3358 — Nested ternary | `sonarjs/no-nested-template-literals` | Zero nesting |
| S1479 — Switch too many cases | `sonarjs/no-small-switch` | Extract to enum map pattern |
| S4144 — Identical functions | `sonarjs/no-identical-functions` | DRY enforcement |
| S1192 — String literals duplicated | `sonarjs/no-duplicate-string` | No magic strings |
| S107 — Too many parameters | `max-params` | Max 3 parameters |
| S138 — Function too long | `max-lines-per-function` | Max 15 lines |
| S134 — Nesting depth | `max-depth` / `coding-guidelines/no-nested-if` | Zero nesting |

### Recommended Config Addition

```javascript
// eslint.config.js — add to plugins and rules
import sonarjs from "eslint-plugin-sonarjs";

// In plugins:
"sonarjs": sonarjs,

// In rules:
"sonarjs/cognitive-complexity": ["error", 10],
"sonarjs/no-identical-functions": "error",
"sonarjs/no-collapsible-if": "error",
"sonarjs/no-all-duplicated-branches": "error",
"sonarjs/no-duplicate-string": ["warn", 3],
"sonarjs/prefer-immediate-return": "warn",
```

---

## 4. Full ESLint Config Reference

The canonical ESLint configuration lives at `eslint.config.js` in the project root. See that file for the active rule set.

---

## 5. Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | All custom rules in `eslint-plugins/coding-guidelines/` loaded | 🔲 |
| 2 | `@typescript-eslint/no-explicit-any` set to `error` | 🔲 |
| 3 | `eslint-plugin-sonarjs` installed and rules mapped | 🔲 |
| 4 | CI runs `eslint --max-warnings 0` on every PR | 🔲 |
| 5 | Test files excluded from `boolean-naming` and `no-magic-strings` | 🔲 |
| 6 | SonarQube sonar-typescript quality gate configured | 🔲 |
| 7 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [ESLint Config](../../../eslint.config.js) — Active rule set
- [Custom Plugin Source](../../../eslint-plugins/coding-guidelines/index.js) — Rule implementations
- [Promise Patterns](./09-promise-await-patterns.md) — Promise.all enforcement rationale
- [TS Standards Reference](./08-typescript-standards-reference.md) — Full TypeScript rules
- [Cross-Language Code Style](../01-cross-language/04-code-style/00-overview.md) — Formatting rules
- [Static Analysis Overview](../01-cross-language/16-static-analysis/00-overview.md) — Cross-language analyzer guide

---

*TypeScript ESLint enforcement v1.0.0 — maps every coding guideline to an enforced lint rule — 2026-04-01*
