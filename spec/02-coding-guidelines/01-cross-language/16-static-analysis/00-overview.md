# Cross-Language Static Analysis & Linter Enforcement

**Version:** 3.2.0  
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`static-analysis` · `linters` · `sonarqube` · `stylecop` · `eslint` · `golangci-lint` · `phpcs` · `clippy` · `ruff`

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

Maps our cross-language coding guidelines to **static analysis tools and linter rules** for each supported language. Every guideline that can be machine-enforced MUST have a corresponding linter rule. This is the single source of truth for tool selection and rule mapping per language.

---

## Document Inventory

| # | File | Language | Analyzer | Status |
|---|------|----------|----------|--------|
| 01 | [11-eslint-enforcement.md](../../02-typescript/11-eslint-enforcement.md) | TypeScript | ESLint + custom plugin + SonarJS | ✅ Complete |
| 02 | [02-go-golangci-lint.md](./02-go-golangci-lint.md) | Go | golangci-lint | ✅ Complete |
| 03 | [03-php-phpcs-phpstan.md](./03-php-phpcs-phpstan.md) | PHP | PHP_CodeSniffer + PHPStan | ✅ Complete |
| 04 | [04-csharp-stylecop.md](./04-csharp-stylecop.md) | C# | StyleCop Analyzers + Roslyn | ✅ Complete |
| 05 | [05-rust-clippy.md](./05-rust-clippy.md) | Rust | Clippy | ✅ Complete |
| 06 | [06-vb-dotnet-analyzers.md](./06-vb-dotnet-analyzers.md) | VB.NET | .NET Analyzers + StyleCop | ✅ Complete |
| 07 | [07-nodejs-eslint.md](./07-nodejs-eslint.md) | Node.js | ESLint (server-side) | ✅ Complete |
| 08 | [08-python-ruff.md](./08-python-ruff.md) | Python | Ruff / Pylint / Flake8 | ✅ Complete |
| 09 | [09-ci-pipeline-quality-gate.md](./09-ci-pipeline-quality-gate.md) | All | CI Pipeline + SonarQube Quality Gate | ✅ Complete |
| 10 | [10-cross-language-rule-matrix.md](./10-cross-language-rule-matrix.md) | All | Side-by-side rule comparison matrix | ✅ Complete |
| — | 97-acceptance-criteria.md | — | — |
| — | 98-changelog.md | — | — |
| — | 99-consistency-report.md | — | — |

---

## Cross-Language Rule → Analyzer Mapping

These coding guidelines apply to **all** languages. Each language-specific doc maps them to concrete linter rules.

| Guideline | Spec Source | Enforcement Category |
|-----------|-------------|---------------------|
| Zero nested `if` | [Code Style §R2](../04-code-style/01-braces-and-nesting.md) | Complexity / nesting |
| Boolean naming (`is/has/can/should/was/will`) | [Boolean Principles](../02-boolean-principles/00-overview.md) | Naming convention |
| No magic strings | [Master §5](../15-master-coding-guidelines/05-magic-strings-and-organization.md) | Literal detection |
| Max 15-line functions | [Code Style §R6](../04-code-style/04-function-and-type-size.md) | Function size |
| No else after return | [Code Style §R7](../04-code-style/01-braces-and-nesting.md) | Control flow |
| Blank line before return | [Code Style §R4](../04-code-style/03-blank-lines-and-spacing.md) | Formatting |
| DRY — no duplicate code | [DRY Principles](../08-dry-principles.md) | Duplication detection |
| No `any` / loose types | [Strict Typing](../13-strict-typing.md) | Type safety |
| Promise.all for independent calls | [Promise Patterns](../../02-typescript/09-promise-await-patterns.md) | Async patterns (TS/JS/Node) |
| Single return value | [Master §4](../15-master-coding-guidelines/04-type-safety.md) | Return pattern (Go) |

---

## SonarQube Integration

SonarQube provides cross-language analysis. Rules that map to our guidelines:

| SonarQube Rule ID | Description | Our Guideline |
|-------------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S1066 | Collapsible if | Zero nesting |
| S1126 | Return boolean directly | Boolean principles |
| S4144 | Identical functions | DRY |
| S1192 | Duplicated string literals | No magic strings |
| S107 | Too many parameters | Max 3 parameters |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |

### SonarQube Quality Gate

```yaml

# sonar-project.properties

sonar.qualitygate.conditions:
  - metric: cognitive_complexity  threshold: 10
  - metric: duplicated_lines_density  threshold: 3
  - metric: code_smells  rating: A
```

---

## Tool Selection per Language

| Language | Primary Linter | Type Checker | Style Checker | SonarQube Plugin |
|----------|---------------|--------------|---------------|------------------|
| TypeScript | ESLint + custom plugin | tsc (`strict: true`) | ESLint formatting rules | sonar-typescript |
| Go | golangci-lint | go vet | golangci-lint (style linters) | sonar-go |
| PHP | PHP_CodeSniffer | PHPStan (level 9) | PHPCS (PSR-12 + custom) | sonar-php |
| C# | StyleCop Analyzers | Roslyn analyzers | StyleCop + .editorconfig | sonar-csharp |
| Rust | Clippy | rustc (built-in) | rustfmt | sonar-rust (community) |
| VB.NET | .NET Analyzers | Roslyn analyzers | StyleCop for VB | sonar-vbnet |
| Node.js | ESLint (same as TS) | tsc or JSDoc types | ESLint formatting rules | sonar-javascript |
| Python | Ruff (preferred) / Pylint | mypy or pyright | Ruff / Black | sonar-python |

| — | 97-acceptance-criteria.md | — | — |
| — | 98-changelog.md | — | — |
| — | 99-consistency-report.md | — | — |
---

## Cross-References

- [TypeScript ESLint Enforcement](../../02-typescript/11-eslint-enforcement.md) — Step 1 (complete)
- [CI Pipeline & Quality Gate](./09-ci-pipeline-quality-gate.md) — Unified CI pipeline
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Source rules
- [Master Coding Guidelines](../15-master-coding-guidelines/00-overview.md) — Full checklist
- [Boolean Principles](../02-boolean-principles/00-overview.md) — Boolean naming rules

---

*Static analysis overview v1.0.0 — cross-language linter enforcement mapping — 2026-04-01*
