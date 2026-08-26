# C# — StyleCop Analyzers + Roslyn Enforcement Rule Mapping

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`csharp` · `c#` · `stylecop` · `roslyn` · `static-analysis` · `roslynator` · `net-analyzers` · `editorconfig`

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

Maps every cross-language coding guideline to its **StyleCop Analyzers**, **Roslyn Analyzers**, and **.editorconfig** enforcement rule for C# projects.

---

## 1. Tool Selection

| Tool | Purpose |
|------|---------|
| StyleCop.Analyzers | Style, naming, formatting, documentation |
| Microsoft.CodeAnalysis.NetAnalyzers | Code quality, design, performance |
| Roslynator.Analyzers | 500+ refactoring analyzers |
| SonarAnalyzer.CSharp | SonarQube-equivalent in-IDE analysis |
| .editorconfig | IDE-enforced style rules |

### Installation

```xml
<!-- In .csproj -->
<ItemGroup>
    <PackageReference Include="StyleCop.Analyzers" Version="1.2.0-beta.*" PrivateAssets="all" />
    <PackageReference Include="Roslynator.Analyzers" Version="4.12.*" PrivateAssets="all" />
    <PackageReference Include="SonarAnalyzer.CSharp" Version="9.32.*" PrivateAssets="all" />
</ItemGroup>

<PropertyGroup>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>
    <AnalysisLevel>latest-all</AnalysisLevel>
    <Nullable>enable</Nullable>
</PropertyGroup>
```

---

## 2. Guideline → Analyzer Rule Mapping

### 2.1 Code Style & Structure

| Guideline | Spec Source | Analyzer | Rule ID | Setting |
|-----------|-------------|----------|---------|---------|
| Zero nested `if` | [Code Style §R2](../04-code-style/01-braces-and-nesting.md) | Roslynator | `RCS1208` | Reduce `if` nesting |
| No `else` after return | [Code Style §R7](../04-code-style/01-braces-and-nesting.md) | Roslynator | `RCS1004` | Remove braces from `if-else` |
| No `else` after return (alt) | [Code Style §R7](../04-code-style/01-braces-and-nesting.md) | SonarAnalyzer | `S1126` | Return boolean directly |
| Max 15-line functions | [Code Style §R6](../04-code-style/04-function-and-type-size.md) | SonarAnalyzer | `S138` | `max: 15` |
| Blank line before return | [Code Style §R4](../04-code-style/03-blank-lines-and-spacing.md) | StyleCop | `SA1513` | Closing brace followed by blank line |
| No dead code | [Code Style §R5](../04-code-style/00-overview.md) | NetAnalyzers | `IDE0051` | Remove unused private members |
| No dead code (alt) | [Code Style §R5](../04-code-style/00-overview.md) | NetAnalyzers | `IDE0052` | Remove unread private members |
| Braces required | [Code Style §R1](../04-code-style/01-braces-and-nesting.md) | StyleCop | `SA1503` | Braces should not be omitted |

### 2.2 Naming Conventions

| Guideline | Spec Source | Analyzer | Rule ID | Setting |
|-----------|-------------|----------|---------|---------|
| Boolean naming (`Is/Has/Can/Should/Was/Will`) | [Boolean Principles](../02-boolean-principles/00-overview.md) | custom `.editorconfig` | `dotnet_naming_rule` | prefix pattern |
| PascalCase for public members | [Key Naming](../11-key-naming-pascalcase.md) | StyleCop | `SA1300` | Element must begin with upper-case |
| PascalCase for constants | [Key Naming](../11-key-naming-pascalcase.md) | StyleCop | `SA1303` | Const field names must begin with upper-case |
| No boolean flag params | [Function Naming](../10-function-naming.md) | SonarAnalyzer | `S2360` | Optional parameters should not be used |
| No raw negation (`!Method()`) | [No Negatives](../12-no-negatives.md) | code review | — | Use positive counterpart |

### 2.3 Type Safety

| Guideline | Spec Source | Analyzer | Rule ID | Setting |
|-----------|-------------|----------|---------|---------|
| No `dynamic` type (equiv. of `any`) | [Strict Typing](../13-strict-typing.md) | NetAnalyzers | `CA1859` | Use concrete types |
| No `object` in APIs | [Strict Typing](../13-strict-typing.md) | SonarAnalyzer | `S3776` | Cognitive complexity |
| Nullable enabled | [Null Safety](../19-null-pointer-safety.md) | compiler | `<Nullable>enable</Nullable>` | NRT enforcement |
| No null-forgiving (`!`) | [Null Safety](../19-null-pointer-safety.md) | Roslynator | `RCS1249` | Unnecessary null-forgiving |
| Max 3 parameters | [Strict Typing](../13-strict-typing.md) | SonarAnalyzer | `S107` | `max: 3` |

### 2.4 Complexity & DRY

| Guideline | Spec Source | Analyzer | Rule ID | Setting |
|-----------|-------------|----------|---------|---------|
| Cyclomatic complexity | [Complexity](../06-cyclomatic-complexity.md) | SonarAnalyzer | `S3776` | `threshold: 10` |
| No magic strings / numbers | [Magic Strings](../15-master-coding-guidelines/05-magic-strings-and-organization.md) | SonarAnalyzer | `S1192` | `threshold: 3` |
| No duplicate code | [DRY Principles](../08-dry-principles.md) | SonarAnalyzer | `S4144` | Identical functions |
| No identical branches | [DRY Principles](../08-dry-principles.md) | SonarAnalyzer | `S1871` | Identical if/else branches |

### 2.5 Documentation

| Guideline | Spec Source | Analyzer | Rule ID | Setting |
|-----------|-------------|----------|---------|---------|
| XML doc on public members | [Code Style §R8](../04-code-style/06-comments-and-documentation.md) | StyleCop | `SA1600` | Elements should be documented |
| Summary required | [Code Style §R8](../04-code-style/06-comments-and-documentation.md) | StyleCop | `SA1604` | Element documentation must have summary |

---

## 3. SonarQube C# Rules

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S1066 | Collapsible if | Zero nesting |
| S1126 | Return boolean directly | Boolean principles |
| S4144 | Identical functions | DRY |
| S1192 | Duplicated string literals | No magic strings |
| S107 | Too many parameters | Max 3 parameters |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |

---

## 4. Reference `.editorconfig`

```ini
[*.cs]

# Naming: Boolean properties must use Is/Has/Can/Should/Was/Will prefix

# (enforced via dotnet_naming_rule — requires custom symbol group)

# Braces

csharp_prefer_braces = true:error

# Expression-bodied members (short methods)

csharp_style_expression_bodied_methods = when_on_single_line:suggestion

# Null checks

csharp_style_conditional_delegate_call = true:error
dotnet_style_coalesce_expression = true:error
dotnet_style_null_propagation = true:error

# Var usage

csharp_style_var_for_built_in_types = false:warning
csharp_style_var_when_type_is_apparent = true:suggestion

# Unused

dotnet_code_quality_unused_parameters = all:error
dotnet_remove_unnecessary_suppression_exclusions = none

# Severity overrides

dotnet_diagnostic.SA1503.severity = error
dotnet_diagnostic.SA1513.severity = warning
dotnet_diagnostic.SA1600.severity = warning
dotnet_diagnostic.IDE0051.severity = error
dotnet_diagnostic.IDE0052.severity = error
dotnet_diagnostic.CA1859.severity = error
dotnet_diagnostic.RCS1208.severity = error
dotnet_diagnostic.RCS1004.severity = error
```

---

## 5. Reference `.globalconfig`

For repo-wide rule severity (complements `.editorconfig`):

```ini
is_global = true

# Zero nesting

dotnet_diagnostic.S134.severity = error

# Max function length

dotnet_diagnostic.S138.severity = error

# Cognitive complexity

dotnet_diagnostic.S3776.severity = error

# Duplicated strings

dotnet_diagnostic.S1192.severity = warning

# Too many parameters

dotnet_diagnostic.S107.severity = error
```

---

## 6. Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | StyleCop.Analyzers NuGet added to all projects | 🔲 |
| 2 | Roslynator.Analyzers NuGet added | 🔲 |
| 3 | SonarAnalyzer.CSharp NuGet added | 🔲 |
| 4 | `TreatWarningsAsErrors = true` in `Directory.Build.props` | 🔲 |
| 5 | `EnforceCodeStyleInBuild = true` | 🔲 |
| 6 | `AnalysisLevel = latest-all` | 🔲 |
| 7 | `Nullable = enable` globally | 🔲 |
| 8 | `.editorconfig` with severity overrides at repo root | 🔲 |
| 9 | CI runs `dotnet build` with `/warnaserror` on every PR | 🔲 |
| 10 | SonarQube sonar-csharp quality gate configured | 🔲 |
| 11 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Cross-language analyzer guide
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Formatting rules
- [Boolean Principles](../02-boolean-principles/00-overview.md) — Boolean naming rules
- [Strict Typing](../13-strict-typing.md) — Type safety rules
- [Null Safety](../19-null-pointer-safety.md) — Null/nil safety guards
- [DRY Principles](../08-dry-principles.md) — Deduplication rules

---

*C# StyleCop + Roslyn enforcement v1.0.0 — maps every coding guideline to C# analyzer rules — 2026-04-01*
