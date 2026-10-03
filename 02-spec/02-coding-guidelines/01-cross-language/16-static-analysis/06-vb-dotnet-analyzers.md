# VB.NET — .NET Analyzers + StyleCop Enforcement (AI Execution Prompt)

> **/goal** Configure and enforce .NET Analyzers, StyleCop.Analyzers, Roslynator, and SonarAnalyzer for VB.NET projects to automate coding guidelines in CI/CD pipelines.
> **/learn** Enforce `Option Strict On`, max 15 lines per function, max 3 parameters, cognitive complexity ≤ 10, zero nested if statements, and ban magic values and late binding.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce `Option Strict On` and zero late binding across all VB.NET project configurations.
- [ ] `/learn` Map SonarAnalyzer (S134, S138, S107, S3776) and Roslynator (RCS1004) rules to block nested `If`, oversized functions (>15 lines), and redundant `Else` clauses.
- [ ] `/goal` Configure `.editorconfig` and `.globalconfig` to treat coding guideline analyzer warnings as build errors (`TreatWarningsAsErrors = true`).
- [ ] `/learn` Ban magic strings (S1192) and magic numbers (S109) by extracting them into named `Const` identifiers.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`vb.net` · `visual-basic` · `static-analysis` · `stylecop` · `roslyn` · `net-analyzers` · `roslynator`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `readme.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Maps cross-language coding guidelines to **.NET Analyzers**, **StyleCop.Analyzers**, **Roslynator**, and **SonarAnalyzer.VisualBasic** for VB.NET projects. Mirrors the [C# StyleCop spec](./04-csharp-stylecop.md) with VB.NET-specific diagnostics.

---

## Guideline → Analyzer Rule Mapping

| # | Guideline | Analyzer | Rule ID | Severity | Notes |
|---|-----------|----------|---------|----------|-------|
| 1 | Zero nested `if` | SonarAnalyzer | S134 | Error | `max = 1` nesting depth |
| 2 | No `Else` after `Return` | Roslynator | RCS1004 | Error | Remove redundant `Else` clause |
| 3 | Boolean naming (`Is/Has/Can/Should/Was/Will`) | .NET Analyzers | CA1716, CA1720 | Warning | Partial; custom Roslyn analyzer recommended |
| 4 | No magic strings | SonarAnalyzer | S1192 | Error | Threshold: 3 occurrences → must be `Const` |
| 5 | No magic numbers | SonarAnalyzer | S109 | Error | Extract to named `Const` |
| 6 | Max 15-line functions | SonarAnalyzer | S138 | Error | `max = 15` |
| 7 | Max 3 parameters | SonarAnalyzer | S107 | Error | `max = 3` |
| 8 | Cognitive complexity ≤ 10 | SonarAnalyzer | S3776 | Error | Threshold: 10 |
| 9 | DRY — no duplicate code | SonarAnalyzer | S1871, S4144 | Error | Identical branches / functions |
| 10 | No `Option Strict Off` | Compiler | `Option Strict On` | Error | Must be project-wide |
| 11 | No late binding | .NET Analyzers | BC42016, BC42017 | Error | Implicit conversion warnings as errors |
| 12 | Blank line before `Return` | StyleCop | SA1513 | Warning | Blank line before closing brace / return |
| 13 | Single responsibility — small types | SonarAnalyzer | S2436 | Warning | Too many type parameters |
| 14 | No unused variables | .NET Analyzers | IDE0059 | Error | Unnecessary assignment |
| 15 | No unused imports | .NET Analyzers | IDE0005 | Error | Remove unnecessary `Imports` |
| 16 | Use `String.IsNullOrEmpty` | .NET Analyzers | CA1820 | Error | Don't compare string length to zero |
| 17 | Dispose pattern | .NET Analyzers | CA1816, CA2000 | Error | Proper `IDisposable` usage |

---

## SonarQube Rule Mapping (sonar-vbnet)

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S4144 | Identical functions | DRY |
| S1066 | Collapsible `If` | Zero nesting |
| S1192 | Duplicated string literals | No magic strings |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |
| S107 | Too many parameters | Max 3 parameters |
| S109 | Magic numbers | No magic numbers |
| S1126 | Return boolean directly | Boolean principles |

---

## VB.NET-Specific Rules

| Rule | Analyzer | Rule ID | Severity | Notes |
|------|----------|---------|----------|-------|
| `Option Strict On` required | Compiler | — | Error | Project property; no late binding |
| `Option Explicit On` required | Compiler | — | Error | All variables must be declared |
| `Option Infer On` allowed | Compiler | — | — | Type inference is acceptable |
| No `On Error Resume Next` | SonarAnalyzer | S2359 | Error | Use `Try/Catch` structured error handling |
| No `GoTo` | SonarAnalyzer | S907 | Error | Structured control flow only |
| Use `Is Nothing` / `IsNot Nothing` | .NET Analyzers | — | Warning | Prefer over `= Nothing` |
| Use string interpolation | .NET Analyzers | IDE0071 | Warning | Prefer `$"..."` over `String.Format` |

---

## Reference Configuration

### `.editorconfig`

```ini

# VB.NET Analyzer Severity Overrides

[*.vb]

# Zero nesting / early exit

dotnet_diagnostic.RCS1004.severity = error

# Unused code

dotnet_diagnostic.IDE0059.severity = error
dotnet_diagnostic.IDE0005.severity = error

# String handling

dotnet_diagnostic.CA1820.severity = error

# Dispose pattern

dotnet_diagnostic.CA1816.severity = error
dotnet_diagnostic.CA2000.severity = error

# Style

dotnet_diagnostic.SA1513.severity = warning
dotnet_diagnostic.IDE0071.severity = warning
```

### `.vbproj` — Project File

```xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OptionStrict>On</OptionStrict>
    <OptionExplicit>On</OptionExplicit>
    <OptionInfer>On</OptionInfer>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>
    <EnableNETAnalyzers>true</EnableNETAnalyzers>
    <AnalysisLevel>latest-all</AnalysisLevel>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="StyleCop.Analyzers" Version="1.2.*">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>analyzers</IncludeAssets>
    </PackageReference>
    <PackageReference Include="Roslynator.Analyzers" Version="4.*">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>analyzers</IncludeAssets>
    </PackageReference>
    <PackageReference Include="SonarAnalyzer.VisualBasic" Version="9.*">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>analyzers</IncludeAssets>
    </PackageReference>
  </ItemGroup>
</Project>
```

### `.globalconfig`

```ini
is_global = true

# SonarAnalyzer thresholds

dotnet_diagnostic.S138.max = 15
dotnet_diagnostic.S107.max = 3
dotnet_diagnostic.S3776.threshold = 10
dotnet_diagnostic.S134.max = 1
dotnet_diagnostic.S1192.threshold = 3

# SonarAnalyzer severity

dotnet_diagnostic.S138.severity = error
dotnet_diagnostic.S107.severity = error
dotnet_diagnostic.S3776.severity = error
dotnet_diagnostic.S134.severity = error
dotnet_diagnostic.S1192.severity = error
dotnet_diagnostic.S1871.severity = error
dotnet_diagnostic.S4144.severity = error
dotnet_diagnostic.S109.severity = error
dotnet_diagnostic.S2359.severity = error
dotnet_diagnostic.S907.severity = error
```

---

## Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | Set `Option Strict On` in all `.vbproj` files | 🔲 |
| 2 | Add `.editorconfig` severity overrides | 🔲 |
| 3 | Add `.globalconfig` with SonarAnalyzer thresholds | 🔲 |
| 4 | Install NuGet analyzer packages | 🔲 |
| 5 | Set `TreatWarningsAsErrors = true` | 🔲 |
| 6 | CI builds with `/warnaserror` flag | 🔲 |
| 7 | SonarQube sonar-vbnet plugin configured | 🔲 |
| 8 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./readme.md) — Parent document
- [C# StyleCop Enforcement](./04-csharp-stylecop.md) — Sibling .NET language
- [Cross-Language Code Style](../04-code-style/readme.md) — Source rules
- [Master Coding Guidelines](../15-master-coding-guidelines/readme.md) — Full checklist

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-LINT-VB: VB.NET Static Analysis & Analyzer Enforcement

**Given** VB.NET source files and project configurations (`.vb`, `.vbproj`).
**When** Linters and .NET analyzers audit the codebase during local builds or CI.
**Then** Option Strict is enabled, nesting depth is ≤ 1, functions do not exceed 15 lines, parameters do not exceed 3, and zero analyzer violations occur with `TreatWarningsAsErrors`.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/16-static-analysis --check-only
```
**Expected:** exit 0. Zero violations.

---

*VB.NET .NET Analyzers enforcement v1.0.0 — cross-language guideline mapping — 2026-04-01*
