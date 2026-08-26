# Cross-Language Rule Matrix

> **Version:** 1.0.0  
> **Updated:** 2026-04-01  
> **AI Confidence:** 95%  
> **Ambiguity:** 5%  
> **Keywords:** rule matrix, SonarQube, linter, cross-language, comparison

---

## Purpose

Side-by-side mapping of every enforced rule across all 8 languages, with SonarQube IDs and per-language linter rule names.

---

## 1. Universal Thresholds

| Metric | Limit | Enforced In |
|--------|-------|-------------|
| Max function length | ≤ 15 lines | All 8 languages |
| Max parameters | ≤ 3 | All 8 languages |
| Cognitive complexity | ≤ 10 | All 8 languages |
| Max nesting depth | ≤ 1 | All 8 languages |

---

## 2. SonarQube Rule → Linter Rule Matrix

### 2.1 Function Length (SonarQube S138)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `max-lines-per-function` | `{ max: 15 }` |
| Go | golangci-lint | `funlen` | `lines: 15` |
| PHP | PHPCS | `Generic.Metrics.FunctionLength` | `maxLength=15` |
| C# | StyleCop | `SA1200` + Roslyn `CA1502` | `.editorconfig` |
| Rust | Clippy | `too_many_lines` | `max-fn-lines = 15` |
| VB.NET | .NET Analyzers | `CA1502` | `.editorconfig` |
| Node.js | ESLint | `max-lines-per-function` | `{ max: 15 }` |
| Python | Ruff | `PLR0915` (too-many-statements) | `max-statements = 15` |

### 2.2 Parameter Count (SonarQube S107)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `max-params` | `{ max: 3 }` |
| Go | golangci-lint | `funlen` (custom) | reviewer-enforced |
| PHP | PHPStan | Level 9 analysis | reviewer-enforced |
| C# | Roslyn | `CA1026` | `.editorconfig` |
| Rust | Clippy | `too_many_arguments` | `max-fn-params = 3` |
| VB.NET | .NET Analyzers | `CA1026` | `.editorconfig` |
| Node.js | ESLint | `max-params` | `{ max: 3 }` |
| Python | Ruff | `PLR0913` | `max-args = 3` |

### 2.3 Cognitive Complexity (SonarQube S3776)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `sonarjs/cognitive-complexity` | `10` |
| Go | golangci-lint | `gocognit` | `min-complexity: 10` |
| PHP | PHPStan | `phpstan-strict-rules` | level 9 |
| C# | Roslyn | `S3776` | SonarScanner |
| Rust | Clippy | `cognitive_complexity` | `threshold = 10` |
| VB.NET | .NET Analyzers | `S3776` | SonarScanner |
| Node.js | ESLint | `sonarjs/cognitive-complexity` | `10` |
| Python | Ruff | `C901` | `max-complexity = 10` |

### 2.4 Nesting Depth (SonarQube S134)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `max-depth` | `{ max: 1 }` |
| Go | golangci-lint | `nestif` | `min-complexity: 1` |
| PHP | PHPCS | `Generic.Metrics.NestingLevel` | `maxNestingLevel=1` |
| C# | StyleCop | custom analyzer | `.editorconfig` |
| Rust | Clippy | `manual` | reviewer-enforced |
| VB.NET | .NET Analyzers | `manual` | reviewer-enforced |
| Node.js | ESLint | `max-depth` | `{ max: 1 }` |
| Python | Ruff | `PLR1702` | `max-nested-blocks = 1` |

### 2.5 Boolean Return Simplification (SonarQube S1126)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `sonarjs/prefer-single-boolean-return` | `error` |
| Go | golangci-lint | `gosimple` (`S1008`) | default |
| PHP | PHPStan | level 9 | default |
| C# | Roslyn | `IDE0046` | `.editorconfig` |
| Rust | Clippy | `needless_bool` | default |
| VB.NET | .NET Analyzers | `IDE0046` | `.editorconfig` |
| Node.js | ESLint | `sonarjs/prefer-single-boolean-return` | `error` |
| Python | Ruff | `SIM103` | default |

### 2.6 Duplicate Functions (SonarQube S4144)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `sonarjs/no-identical-functions` | `error` |
| Go | golangci-lint | `dupl` | `threshold: 50` |
| PHP | PHPCS | `Generic.CodeAnalysis.DuplicateCode` | default |
| C# | Roslyn | `S4144` | SonarScanner |
| Rust | Clippy | `S4144` | SonarScanner |
| VB.NET | .NET Analyzers | `S4144` | SonarScanner |
| Node.js | ESLint | `sonarjs/no-identical-functions` | `error` |
| Python | Ruff | `PLR0801` (planned) | SonarScanner fallback |

### 2.7 Unused Code (SonarQube S1481 / S1144)

| Language | Linter | Rule Name | Config |
|----------|--------|-----------|--------|
| TypeScript | ESLint | `no-unused-vars` | `error` |
| Go | compiler | built-in | compile error |
| PHP | PHPStan | dead code detection | level 9 |
| C# | Roslyn | `CS0219` / `IDE0051` | `warning` |
| Rust | compiler | `#[warn(dead_code)]` | default |
| VB.NET | .NET Analyzers | `IDE0051` | `warning` |
| Node.js | ESLint | `no-unused-vars` | `error` |
| Python | Ruff | `F841` / `F811` | default |

---

## 3. Coverage Matrix Summary

| SonarQube Rule | TS | Go | PHP | C# | Rust | VB | Node | Python |
|----------------|----|----|-----|----|------|----|------|--------|
| S138 (length) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| S107 (params) | ✅ | ⚠️ | ⚠️ | ✅ | ✅ | ✅ | ✅ | ✅ |
| S3776 (complexity) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| S134 (nesting) | ✅ | ✅ | ✅ | ✅ | ⚠️ | ⚠️ | ✅ | ✅ |
| S1126 (bool return) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| S4144 (duplicates) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ |
| S1481 (unused) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

**Legend:** ✅ = native linter rule · ⚠️ = SonarQube fallback or reviewer-enforced

---

## Cross-References

- [CI Pipeline Quality Gate](./09-ci-pipeline-quality-gate.md)
- [Static Analysis Overview](./00-overview.md)
- [TypeScript ESLint Enforcement](../../02-typescript/11-eslint-enforcement.md)
