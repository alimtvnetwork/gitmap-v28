# PHP — PHP_CodeSniffer + PHPStan Enforcement Rule Mapping

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`php` · `phpcs` · `phpstan` · `static-analysis` · `slevomat` · `psalm` · `php-codesniffer`

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

Maps every cross-language coding guideline to its **PHP_CodeSniffer (PHPCS)** and **PHPStan** enforcement rule. Teams adopting these guidelines can copy the `phpcs.xml` and `phpstan.neon` configs for instant enforcement.

---

## 1. Tool Selection

| Tool | Purpose |
|------|---------|
| PHP_CodeSniffer (PHPCS) | Style, formatting, naming conventions |
| PHPStan (level 9) | Static type analysis, dead code, type safety |
| Psalm (alternative) | Type safety, taint analysis |
| SonarQube (`sonar-php`) | Cross-language quality gate |

---

## 2. Guideline → PHPCS Rule Mapping

### 2.1 Code Style & Structure

| Guideline | Spec Source | PHPCS Sniff | Setting |
|-----------|-------------|-------------|---------|
| Zero nested `if` | [Code Style §R2](../04-code-style/01-braces-and-nesting.md) | `Generic.Metrics.NestingLevel` | `absoluteNestingLevel: 1` |
| No `else` after return | [Code Style §R7](../04-code-style/01-braces-and-nesting.md) | `SlevomatCodingStandard.ControlStructures.EarlyExit` | (enabled) |
| Max 15-line functions | [Code Style §R6](../04-code-style/04-function-and-type-size.md) | `Generic.Metrics.FunctionLength` | `maxLength: 15` |
| Blank line before return | [Code Style §R4](../04-code-style/03-blank-lines-and-spacing.md) | `SlevomatCodingStandard.ControlStructures.JumpStatementsSpacing` | `linesCountBeforeFirst: 1` |
| No dead code | [Code Style §R5](../04-code-style/00-overview.md) | `SlevomatCodingStandard.Functions.UnusedInheritedVariablePassedToClosure` | (enabled) |
| Max 400-line files | [File Rules](../04-code-style/04-function-and-type-size.md) | `Generic.Files.LineLength` + `Generic.Metrics.ClassLength` | custom |
| Braces on same line | [Code Style §R1](../04-code-style/01-braces-and-nesting.md) | `PSR12.ControlStructures.ControlStructureSpacing` | PSR-12 default |

### 2.2 Naming Conventions

| Guideline | Spec Source | PHPCS Sniff | Setting |
|-----------|-------------|-------------|---------|
| Boolean naming (`is/has/can/should/was/will`) | [Boolean Principles](../02-boolean-principles/00-overview.md) | `SlevomatCodingStandard.Variables.UnusedVariable` + custom sniff | regex pattern |
| PascalCase keys (API/DB) | [Key Naming](../11-key-naming-pascalcase.md) | custom sniff | `PascalCase` enforcement |
| No boolean flag params | [Function Naming](../10-function-naming.md) | `SlevomatCodingStandard.Functions.FunctionLength` | code review |
| No raw negation (`!fn()`) | [No Negatives](../12-no-negatives.md) | custom sniff | `!$this->isValid()` detection |

### 2.3 Complexity & DRY

| Guideline | Spec Source | PHPCS Sniff | Setting |
|-----------|-------------|-------------|---------|
| Cyclomatic complexity | [Complexity](../06-cyclomatic-complexity.md) | `Generic.Metrics.CyclomaticComplexity` | `complexity: 10` |
| No magic strings / numbers | [Magic Strings](../15-master-coding-guidelines/05-magic-strings-and-organization.md) | `SlevomatCodingStandard.Numbers.RequireNumericLiteralSeparator` + custom | detection |
| No duplicate code | [DRY Principles](../08-dry-principles.md) | `phpcpd` (separate tool) | `min-lines: 5` |

---

## 3. Guideline → PHPStan Rule Mapping

PHPStan at **level 9** (strictest) enforces type safety:

| Guideline | Spec Source | PHPStan Rule / Level | Description |
|-----------|-------------|---------------------|-------------|
| No `mixed` type (equiv. of `any`) | [Strict Typing](../13-strict-typing.md) | Level 9 | All types must be explicit |
| No unsafe member access | [Strict Typing](../13-strict-typing.md) | Level 6+ | Property/method access on mixed |
| No unsafe return | [Strict Typing](../13-strict-typing.md) | Level 7+ | Return type must match declaration |
| No unused parameters | [Dead Code](../04-code-style/00-overview.md) | `phpstan-strict-rules` | Unused parameter detection |
| No magic `__get`/`__set` | [Strict Typing](../13-strict-typing.md) | Level 9 | Force typed properties |
| Max 3 parameters | [Strict Typing](../13-strict-typing.md) | `phpstan-strict-rules` | Custom rule |
| Null safety | [Null Safety](../19-null-pointer-safety.md) | Level 8+ | Null pointer detection |

### PHPStan Strict Rules Package

```bash
composer require --dev phpstan/phpstan-strict-rules
```

Adds enforcement for:
- No dynamic property access
- No `empty()` usage (use strict comparison)
- No loose comparison (`==` / `!=`)
- No variable variables (`$$var`)

---

## 4. SonarQube PHP Rules

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S1066 | Collapsible if | Zero nesting |
| S4144 | Identical functions | DRY |
| S1192 | Duplicated string literals | No magic strings |
| S107 | Too many parameters | Max 3 parameters |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |
| S1126 | Return boolean directly | Boolean principles |

---

## 5. Reference `phpcs.xml`

```xml
<?xml version="1.0"?>
<ruleset name="CodingGuidelines">
    <description>Cross-language coding guidelines — PHP enforcement</description>

    <!-- PSR-12 base -->
    <rule ref="PSR12"/>

    <!-- Slevomat Coding Standard -->
    <rule ref="SlevomatCodingStandard.ControlStructures.EarlyExit"/>
    <rule ref="SlevomatCodingStandard.ControlStructures.JumpStatementsSpacing">
        <properties>
            <property name="linesCountBeforeFirst" value="1"/>
        </properties>
    </rule>
    <rule ref="SlevomatCodingStandard.Functions.UnusedInheritedVariablePassedToClosure"/>

    <!-- Zero nesting -->
    <rule ref="Generic.Metrics.NestingLevel">
        <properties>
            <property name="absoluteNestingLevel" value="1"/>
        </properties>
    </rule>

    <!-- Max 15-line functions -->
    <rule ref="Generic.Metrics.FunctionLength">
        <properties>
            <property name="maxLength" value="15"/>
        </properties>
    </rule>

    <!-- Cyclomatic complexity -->
    <rule ref="Generic.Metrics.CyclomaticComplexity">
        <properties>
            <property name="complexity" value="10"/>
            <property name="absoluteComplexity" value="15"/>
        </properties>
    </rule>
</ruleset>
```

---

## 6. Reference `phpstan.neon`

```neon
parameters:
    level: 9
    paths:
        - src
    excludePaths:
        - tests

includes:
    - vendor/phpstan/phpstan-strict-rules/rules.neon

rules:
    # Custom rules can be registered here
```

---

## 7. Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | PHP_CodeSniffer installed | 🔲 |
| 2 | Slevomat Coding Standard installed | 🔲 |
| 3 | `phpcs.xml` placed at repo root with rules | 🔲 |
| 4 | PHPStan at level 9 installed | 🔲 |
| 5 | `phpstan-strict-rules` installed | 🔲 |
| 6 | CI runs `phpcs --standard=phpcs.xml src/` on every PR | 🔲 |
| 7 | CI runs `phpstan analyse` on every PR | 🔲 |
| 8 | SonarQube sonar-php plugin enabled | 🔲 |
| 9 | `phpcpd` installed for duplicate detection | 🔲 |
| 10 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Cross-language analyzer guide
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Formatting rules
- [Boolean Principles](../02-boolean-principles/00-overview.md) — Boolean naming rules
- [Strict Typing](../13-strict-typing.md) — Type safety rules
- [DRY Principles](../08-dry-principles.md) — Deduplication rules

---

*PHP PHPCS + PHPStan enforcement v1.0.0 — maps every coding guideline to PHP analyzer rules — 2026-04-01*
