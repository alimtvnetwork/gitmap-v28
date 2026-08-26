# Node.js — ESLint Enforcement (Server-Side)

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`nodejs` · `eslint` · `static-analysis` · `server-side` · `sonarjs` · `typescript-eslint`

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

Maps cross-language coding guidelines to **ESLint rules** for **Node.js server-side** code (Express, Fastify, NestJS, plain Node). Extends the [TypeScript ESLint spec](../../02-typescript/11-eslint-enforcement.md) with Node-specific rules for async patterns, security, and runtime concerns.

---

## Guideline → ESLint Rule Mapping

### Core Cross-Language Rules

| # | Guideline | ESLint Rule | Plugin | Severity | Notes |
|---|-----------|-------------|--------|----------|-------|
| 1 | Zero nested `if` | `max-depth` | core | `error` | Configure: `max: 1` |
| 2 | No else after return | `no-else-return` | core | `error` | `allowElseIf: false` |
| 3 | Boolean naming (`is/has/can/should/was/will`) | `@typescript-eslint/naming-convention` | typescript-eslint | `warn` | Filter: `booleanPrefix: ['is','has','can','should','was','will']` |
| 4 | No magic strings | `no-magic-numbers` | core | `error` | `ignore: [-1, 0, 1, 2]`, `ignoreEnums: true` |
| 5 | Max 15-line functions | `max-lines-per-function` | core | `error` | `max: 15`, `skipBlankLines: true`, `skipComments: true` |
| 6 | Max 3 parameters | `max-params` | core | `error` | `max: 3` |
| 7 | Cognitive complexity | `sonarjs/cognitive-complexity` | eslint-plugin-sonarjs | `error` | Threshold: `10` |
| 8 | DRY — no duplicate code | `sonarjs/no-identical-functions`, `sonarjs/no-duplicated-branches` | eslint-plugin-sonarjs | `error` | — |
| 9 | No `any` / loose types | `@typescript-eslint/no-explicit-any` | typescript-eslint | `error` | — |
| 10 | Blank line before return | `padding-line-between-statements` | core | `warn` | `{ blankLine: 'always', prev: '*', next: 'return' }` |
| 11 | `Promise.all` for independent calls | `sonarjs/no-redundant-jump` | eslint-plugin-sonarjs | `warn` | Manual review; no direct lint |

### Node.js-Specific Rules

| # | Guideline | ESLint Rule | Plugin | Severity | Notes |
|---|-----------|-------------|--------|----------|-------|
| 12 | No callback pattern — use async/await | `no-callback-in-promise`, `prefer-promise-reject-errors` | eslint-plugin-promise | `error` | Modern async only |
| 13 | Always handle promise rejections | `@typescript-eslint/no-floating-promises` | typescript-eslint | `error` | Must `await` or `.catch()` |
| 14 | No `require()` — use ES modules | `@typescript-eslint/no-require-imports` | typescript-eslint | `error` | ESM only |
| 15 | No `process.exit()` in libraries | `n/no-process-exit` | eslint-plugin-n | `error` | Only allowed in CLI entry points |
| 16 | No sync fs methods | `n/no-sync` | eslint-plugin-n | `error` | Use async `fs/promises` |
| 17 | Handle stream errors | `n/handle-callback-err` | eslint-plugin-n | `error` | All error-first callbacks |
| 18 | No deprecated APIs | `n/no-deprecated-api` | eslint-plugin-n | `error` | Auto-detects deprecated Node APIs |
| 19 | No `eval()` / `new Function()` | `no-eval`, `no-new-func` | core | `error` | Security: code injection risk |
| 20 | No dynamic `import()` with variables | `no-restricted-syntax` | core | `warn` | Custom selector for variable imports |
| 21 | Validate environment variables | — | — | — | Use `zod` or `envalid`; enforce via code review |

### Security Rules (eslint-plugin-security)

| # | Rule | Plugin | Severity | Notes |
|---|------|--------|----------|-------|
| 22 | `security/detect-object-injection` | eslint-plugin-security | `warn` | Bracket notation with user input |
| 23 | `security/detect-non-literal-regexp` | eslint-plugin-security | `error` | ReDoS prevention |
| 24 | `security/detect-non-literal-fs-filename` | eslint-plugin-security | `warn` | Path traversal prevention |
| 25 | `security/detect-eval-with-expression` | eslint-plugin-security | `error` | Dynamic eval prevention |
| 26 | `security/detect-child-process` | eslint-plugin-security | `warn` | Command injection prevention |
| 27 | `security/detect-unsafe-regex` | eslint-plugin-security | `error` | ReDoS detection |

---

## SonarQube Rule Mapping (sonar-javascript)

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S4144 | Identical functions | DRY |
| S1066 | Collapsible if | Zero nesting |
| S1192 | Duplicated string literals | No magic strings |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |
| S107 | Too many parameters | Max 3 parameters |
| S4823 | `process.exit()` usage | No process.exit in libraries |
| S1126 | Return boolean directly | Boolean principles |
| S2245 | Pseudo-random values | Use `crypto.randomUUID()` |

---

## Reference Configuration

### `eslint.config.mjs` (Flat Config)

```js
import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import sonarjs from 'eslint-plugin-sonarjs';
import security from 'eslint-plugin-security';
import nodePlugin from 'eslint-plugin-n';
import promisePlugin from 'eslint-plugin-promise';

export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  sonarjs.configs.recommended,
  security.configs.recommended,
  nodePlugin.configs['flat/recommended'],
  promisePlugin.configs['flat/recommended'],
  {
    languageOptions: {
      parserOptions: {
        project: './tsconfig.json',
      },
    },
    rules: {
      // Cross-language rules
      'max-depth': ['error', { max: 1 }],
      'no-else-return': ['error', { allowElseIf: false }],
      'max-lines-per-function': ['error', {
        max: 15, skipBlankLines: true, skipComments: true,
      }],
      'max-params': ['error', { max: 3 }],
      'no-magic-numbers': ['error', {
        ignore: [-1, 0, 1, 2], ignoreEnums: true,
      }],
      'padding-line-between-statements': ['warn',
        { blankLine: 'always', prev: '*', next: 'return' },
      ],

      // TypeScript strict
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/no-floating-promises': 'error',
      '@typescript-eslint/no-require-imports': 'error',

      // SonarJS
      'sonarjs/cognitive-complexity': ['error', 10],
      'sonarjs/no-identical-functions': 'error',
      'sonarjs/no-duplicated-branches': 'error',

      // Node-specific
      'n/no-process-exit': 'error',
      'n/no-sync': 'error',
      'n/no-deprecated-api': 'error',

      // Security
      'no-eval': 'error',
      'no-new-func': 'error',

      // Promise
      'no-callback-in-promise': 'error',
      'prefer-promise-reject-errors': 'error',
    },
  },
);
```

### Required Packages

```bash
npm install -D eslint @eslint/js typescript-eslint \
  eslint-plugin-sonarjs eslint-plugin-security \
  eslint-plugin-n eslint-plugin-promise
```

### CI Command

```bash
npx eslint . --max-warnings 0
```

---

## Differences from TypeScript (Frontend) ESLint

| Concern | Frontend (TS) | Server (Node.js) |
|---------|--------------|-------------------|
| Module system | ESM (Vite/bundler) | ESM (`"type": "module"`) |
| Security plugin | Not required | `eslint-plugin-security` required |
| Node plugin | Not needed | `eslint-plugin-n` required |
| `process.exit` | N/A | Denied |
| Sync I/O | N/A | Denied (`n/no-sync`) |
| Promise handling | Same | `@typescript-eslint/no-floating-promises` critical |
| `require()` | Bundler blocks | Must lint-deny explicitly |

---

## Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | Install ESLint + all plugins | 🔲 |
| 2 | Add `eslint.config.mjs` with rules above | 🔲 |
| 3 | Set `"type": "module"` in `package.json` | 🔲 |
| 4 | CI runs `npx eslint . --max-warnings 0` | 🔲 |
| 5 | SonarQube sonar-javascript configured | 🔲 |
| 6 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Parent document
- [TypeScript ESLint Enforcement](../../02-typescript/11-eslint-enforcement.md) — Frontend sibling
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Source rules
- [Master Coding Guidelines](../15-master-coding-guidelines/00-overview.md) — Full checklist
- [Promise/Await Patterns](../../02-typescript/09-promise-await-patterns.md) — Async guidelines

---

*Node.js ESLint enforcement v1.0.0 — cross-language guideline mapping — 2026-04-01*
