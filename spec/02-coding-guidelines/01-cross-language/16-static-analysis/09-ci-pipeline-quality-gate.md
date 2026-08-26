# Unified CI Pipeline & Quality Gate Specification

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`ci` · `pipeline` · `quality-gate` · `sonarqube` · `github-actions` · `gitlab-ci` · `linter` · `static-analysis` · `continuous-integration`

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

Defines a **standard CI pipeline structure** and **quality gate** that applies to all 8 supported languages. Every repository MUST run the language-appropriate linter, type checker, formatter, and SonarQube analysis as blocking PR checks. This document is the single source of truth for CI enforcement.

---

## 1. Universal Quality Gate

Every PR MUST pass **all** of these checks before merge. No exceptions without a documented exemption in the PR description.

### 1.1 Mandatory Thresholds

| Metric | Threshold | SonarQube Rule | Enforcement |
|--------|-----------|----------------|-------------|
| Cognitive complexity per function | ≤ 10 | S3776 | Linter + SonarQube |
| Max function length | ≤ 15 lines | S138 | Linter |
| Max parameters per function | ≤ 3 | S107 | Linter |
| Max nesting depth | ≤ 1 (zero nested `if`) | S134 | Linter |
| Duplicated lines density | ≤ 3% | — | SonarQube |
| Code smells rating | A | — | SonarQube |
| No new bugs | 0 | — | SonarQube |
| No new vulnerabilities | 0 | — | SonarQube |
| No new security hotspots (unreviewed) | 0 | — | SonarQube |
| Test coverage on new code | ≥ 80% | — | SonarQube |

### 1.2 Blocking vs Warning

| Severity | Effect | Examples |
|----------|--------|---------|
| **Error** (blocking) | PR cannot merge | `no-any`, zero nesting, max params, max lines, unwrap/panic |
| **Warning** (non-blocking) | Must be reviewed; tracked in SonarQube | Boolean naming, blank line before return, documentation |

---

## 2. Pipeline Stages

Every CI pipeline MUST execute these stages **in order**. A failure in any stage stops the pipeline.

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Format  │───▶│   Lint   │───▶│  Type    │───▶│  Test    │───▶│  Sonar   │
│  Check   │    │  Check   │    │  Check   │    │  + Cover │    │  Scan    │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

| Stage | Purpose | Failure = |
|-------|---------|----------|
| **1. Format Check** | Verify code formatting matches standard | PR blocked |
| **2. Lint Check** | Run primary linter with zero warnings | PR blocked |
| **3. Type Check** | Run type checker in strict mode | PR blocked |
| **4. Test + Coverage** | Run tests, generate coverage report | PR blocked |
| **5. SonarQube Scan** | Quality gate analysis on new code | PR blocked |

---

## 3. Language-Specific CI Commands

### 3.1 TypeScript (Frontend)

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `npx prettier --check .` | — |
| Lint | `npx eslint . --max-warnings 0` | [TS ESLint](../../02-typescript/11-eslint-enforcement.md) |
| Type | `npx tsc --noEmit --strict` | — |
| Test | `npx vitest run --coverage` | — |

### 3.2 Go

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `gofmt -l . \| grep . && exit 1 \|\| true` | — |
| Lint | `golangci-lint run --timeout 5m` | [Go golangci-lint](./02-go-golangci-lint.md) |
| Type | `go vet ./...` | — |
| Test | `go test -race -coverprofile=coverage.out ./...` | — |

### 3.3 PHP

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `php-cs-fixer fix --dry-run --diff` | — |
| Lint | `phpcs --standard=phpcs.xml src/` | [PHP PHPCS](./03-php-phpcs-phpstan.md) |
| Type | `phpstan analyse --level=9` | [PHP PHPStan](./03-php-phpcs-phpstan.md) |
| Test | `phpunit --coverage-clover coverage.xml` | — |

### 3.4 C#

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `dotnet format --verify-no-changes` | — |
| Lint | `dotnet build /warnaserror` | [C# StyleCop](./04-csharp-stylecop.md) |
| Type | (included in `dotnet build`) | — |
| Test | `dotnet test --collect:"XPlat Code Coverage"` | — |

### 3.5 Rust

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `cargo fmt --all -- --check` | [Rust Clippy](./05-rust-clippy.md) |
| Lint | `cargo clippy --all-targets --all-features -- -D warnings` | [Rust Clippy](./05-rust-clippy.md) |
| Type | (included in `cargo clippy`) | — |
| Test | `cargo tarpaulin --out xml` | — |

### 3.6 VB.NET

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `dotnet format --verify-no-changes` | — |
| Lint | `dotnet build /warnaserror` | [VB.NET Analyzers](./06-vb-dotnet-analyzers.md) |
| Type | (included in `dotnet build`) | — |
| Test | `dotnet test --collect:"XPlat Code Coverage"` | — |

### 3.7 Node.js (Server)

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `npx prettier --check .` | — |
| Lint | `npx eslint . --max-warnings 0` | [Node.js ESLint](./07-nodejs-eslint.md) |
| Type | `npx tsc --noEmit --strict` | — |
| Test | `npx vitest run --coverage` or `npx jest --coverage` | — |

### 3.8 Python

| Stage | Command | Spec Reference |
|-------|---------|----------------|
| Format | `ruff format --check .` | [Python Ruff](./08-python-ruff.md) |
| Lint | `ruff check . --output-format=github` | [Python Ruff](./08-python-ruff.md) |
| Type | `mypy . --strict` | [Python Ruff](./08-python-ruff.md) |
| Test | `pytest --cov --cov-report=xml` | — |

---

## 4. SonarQube Configuration

### 4.1 `sonar-project.properties`

```properties

# Project identification

sonar.projectKey=<project-key>
sonar.projectName=<project-name>
sonar.projectVersion=1.0

# Source configuration

sonar.sources=src
sonar.tests=tests
sonar.sourceEncoding=UTF-8

# Language-specific coverage reports

# TypeScript / Node.js

sonar.javascript.lcov.reportPaths=coverage/lcov.info

# Go

sonar.go.coverage.reportPaths=coverage.out

# PHP

sonar.php.coverage.reportPaths=coverage.xml

# C# / VB.NET

sonar.cs.opencover.reportsPaths=**/coverage.opencover.xml
sonar.vbnet.opencover.reportsPaths=**/coverage.opencover.xml

# Python

sonar.python.coverage.reportPaths=coverage.xml

# Rust (community plugin)

# sonar.rust.lcov.reportPaths=lcov.info

# Quality gate thresholds (custom)

sonar.qualitygate.wait=true
```

### 4.2 Quality Gate Definition

Create this quality gate in SonarQube and set as default:

| Condition | Metric | Operator | Value |
|-----------|--------|----------|-------|
| New bugs | Bugs on new code | > | 0 |
| New vulnerabilities | Vulnerabilities on new code | > | 0 |
| New code smells rating | Maintainability on new code | worse than | A |
| New coverage | Coverage on new code | < | 80% |
| New duplicated lines | Duplicated lines on new code | > | 3% |
| New security hotspots | Security hotspots reviewed on new code | < | 100% |

### 4.3 SonarQube Rule Profile

Enable these rules in **every** language profile:

| Rule | Description | Threshold |
|------|-------------|-----------|
| S3776 | Cognitive Complexity | 10 |
| S138 | Function too long | 15 lines |
| S107 | Too many parameters | 3 |
| S134 | Nesting depth | 1 |
| S1871 | Identical branches | — |
| S4144 | Identical functions | — |
| S1066 | Collapsible if | — |
| S1126 | Return boolean directly | — |
| S1192 | Duplicated string literals | 3 occurrences |

---

## 5. Reference GitHub Actions Workflow

```yaml

# .github/workflows/quality-gate.yml

name: Quality Gate

on:
  pull_request:
    branches: [main, develop]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Required for SonarQube

      # ── Stage 1: Format ──
      - name: Format check
        run: |
          # Replace with language-specific command from §3

      # ── Stage 2: Lint ──
      - name: Lint check
        run: |
          # Replace with language-specific command from §3

      # ── Stage 3: Type check ──
      - name: Type check
        run: |
          # Replace with language-specific command from §3

      # ── Stage 4: Test + Coverage ──
      - name: Test with coverage
        run: |
          # Replace with language-specific command from §3

      # ── Stage 5: SonarQube ──
      - name: SonarQube scan
        uses: SonarSource/sonarqube-scan-action@v3
        env:
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
          SONAR_HOST_URL: ${{ secrets.SONAR_HOST_URL }}

      - name: SonarQube quality gate
        uses: SonarSource/sonarqube-quality-gate-action@v1
        timeout-minutes: 5
        env:
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

---

## 6. Reference GitLab CI Configuration

```yaml

# .gitlab-ci.yml

stages:
  - format
  - lint
  - type-check
  - test
  - sonar

format:
  stage: format
  script:
    # Replace with language-specific command from §3

lint:
  stage: lint
  script:
    # Replace with language-specific command from §3

type-check:
  stage: type-check
  script:
    # Replace with language-specific command from §3

test:
  stage: test
  script:
    # Replace with language-specific command from §3
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml

sonar:
  stage: sonar
  image: sonarsource/sonar-scanner-cli:latest
  variables:
    SONAR_USER_HOME: "${CI_PROJECT_DIR}/.sonar"
  cache:
    key: sonar
    paths:
      - .sonar/cache
  script:
    - sonar-scanner
  allow_failure: false
```

---

## 7. Mono-Repo Strategy

For repositories containing multiple languages, run language-specific checks in **parallel jobs**:

```yaml

# GitHub Actions — mono-repo

jobs:
  typescript:
    runs-on: ubuntu-latest
    steps:
      - run: npx eslint packages/frontend/ --max-warnings 0
      - run: npx tsc --noEmit --strict -p packages/frontend/tsconfig.json
      - run: npx vitest run --coverage --project packages/frontend

  go:
    runs-on: ubuntu-latest
    steps:
      - run: cd packages/api && golangci-lint run --timeout 5m
      - run: cd packages/api && go test -race -coverprofile=coverage.out ./...

  python:
    runs-on: ubuntu-latest
    steps:
      - run: ruff check packages/ml/ --output-format=github
      - run: mypy packages/ml/ --strict
      - run: cd packages/ml && pytest --cov --cov-report=xml

  sonar:
    needs: [typescript, go, python]
    runs-on: ubuntu-latest
    steps:
      - uses: SonarSource/sonarqube-scan-action@v3
```

---

## 8. Exemption Process

When a rule must be suppressed:

| Requirement | Description |
|-------------|-------------|
| **Inline justification** | Every suppression MUST include a comment explaining why |
| **PR description** | Exemptions MUST be listed in the PR description |
| **Time-boxed** | Suppressions MUST include a TODO with a ticket number for removal |
| **No blanket disables** | Never disable a rule for an entire file or project |

### Suppression Syntax per Language

| Language | Syntax | Example |
|----------|--------|---------|
| TypeScript / Node.js | `// eslint-disable-next-line rule-name -- reason` | `// eslint-disable-next-line max-params -- factory requires 4 deps` |
| Go | `//nolint:lintername // reason` | `//nolint:funlen // migration helper — TODO: PROJ-1234` |
| PHP | `// phpcs:ignore Sniff.Name -- reason` | `// phpcs:ignore Generic.Metrics.FunctionLength -- legacy` |
| C# / VB.NET | `#pragma warning disable RULE // reason` | `#pragma warning disable S138 // generated code` |
| Rust | `#[allow(clippy::lint_name)] // reason` | `#[allow(clippy::too_many_arguments)] // FFI boundary` |
| Python | `# noqa: RULE -- reason` | `# noqa: PLR0913 -- CLI entrypoint` |

---

## 9. Dashboard & Reporting

| Metric | Source | Frequency |
|--------|--------|-----------|
| Quality gate pass rate | SonarQube | Per PR |
| Code coverage trend | SonarQube | Weekly |
| Technical debt ratio | SonarQube | Sprint review |
| Lint warning trend | CI logs | Weekly |
| Suppression count | `grep -r "nolint\|noqa\|eslint-disable\|phpcs:ignore\|pragma warning disable\|allow(clippy" src/` | Monthly |

---

## 10. Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | CI pipeline created with all 5 stages | 🔲 |
| 2 | Language-specific linter configured per §3 | 🔲 |
| 3 | All linter warnings treated as errors (`--max-warnings 0` / `-D warnings` / `/warnaserror`) | 🔲 |
| 4 | SonarQube project created with `sonar-project.properties` | 🔲 |
| 5 | SonarQube quality gate created with thresholds from §4.2 | 🔲 |
| 6 | SonarQube rule profile configured with rules from §4.3 | 🔲 |
| 7 | Coverage reports wired to SonarQube | 🔲 |
| 8 | Quality gate set as required status check on `main` / `develop` | 🔲 |
| 9 | Exemption process documented in `CONTRIBUTING.md` | 🔲 |
| 10 | Team reviewed and approved pipeline + thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Parent document
- [TypeScript ESLint Enforcement](../../02-typescript/11-eslint-enforcement.md) — TS linter config
- [Go golangci-lint](./02-go-golangci-lint.md) — Go linter config
- [PHP PHPCS + PHPStan](./03-php-phpcs-phpstan.md) — PHP linter config
- [C# StyleCop + Roslyn](./04-csharp-stylecop.md) — C# linter config
- [Rust Clippy](./05-rust-clippy.md) — Rust linter config
- [VB.NET Analyzers](./06-vb-dotnet-analyzers.md) — VB.NET linter config
- [Node.js ESLint](./07-nodejs-eslint.md) — Node.js linter config
- [Python Ruff](./08-python-ruff.md) — Python linter config
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Source rules

---

*Unified CI Pipeline & Quality Gate v1.0.0 — standard enforcement across 8 languages — 2026-04-01*
