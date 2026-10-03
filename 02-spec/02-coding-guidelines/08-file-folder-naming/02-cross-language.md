# File & Folder Naming — Cross-Language Rules (AI Execution Prompt)

> **/goal** Enforce universal cross-language file and folder naming invariants (zero spaces, no special characters, strict lowercase for system files, dot-prefixed hidden files).
> **/learn** Master cross-language naming rules, PowerShell `.ps1` kebab-case standards, leading zero numbering (`01-`, `02-`), and case-sensitivity safety.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Eliminate all spaces and special characters from repository filenames across all platforms.
- [ ] `/learn` Never use uppercase letters in markdown, system files, scripts, or documentation (strict lowercase mandate).
- [ ] `/goal` Enforce lowercase kebab-case for PowerShell scripts and numbered files (`01-file.ps1`).
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Universal Rules

These rules apply to **every language and framework** in the project.

### 1. No Spaces in Names

```
✅ user-profile.ts
✅ http_handler.go
❌ user profile.ts
❌ http handler.go
```

### 2. No Special Characters

Only alphanumeric characters, hyphens (`-`), underscores (`_`), and dots (`.`) are allowed.

```
✅ api-client.ts
❌ api@client.ts
❌ api client (v2).ts
```

### 3. Always Lowercase Folders

Folders/directories MUST be lowercase in all languages except C# (which uses PascalCase).

```
✅ src/components/
✅ internal/handlers/
❌ src/Components/
❌ Internal/Handlers/
```

### 4. File Extensions Must Match Language

| Language | Extension |
|----------|-----------|
| TypeScript | `.ts`, `.tsx` |
| JavaScript | `.js`, `.jsx` |
| Go | `.go` |
| PHP | `.php` |
| PowerShell | `.ps1`, `.psm1`, `.psd1` |
| Rust | `.rs` |
| C# | `.cs` |

### 5. Test Files Follow Source Naming

| Language | Pattern | Example |
|----------|---------|---------|
| Go | `*_test.go` | `handler_test.go` |
| TypeScript | `*.test.ts` or `*.spec.ts` | `handler.test.ts` |
| Rust | inline `#[cfg(test)]` or `tests/` folder | `tests/handler_test.rs` |
| PHP | `*Test.php` | `AdminSettingsTest.php` |
| C# | `*Tests.cs` | `UserServiceTests.cs` |

### 6. Config Files Are Lowercase

All config files use lowercase with dots or hyphens:

```
✅ .eslintrc.json
✅ tsconfig.json
✅ docker-compose.yml
❌ Docker-Compose.yml
```

### 7. README and LICENSE Are UPPERCASE

These are the only exceptions to the lowercase rule:

```
✅ readme.md
✅ LICENSE
✅ changelog.md
```

---

## Forbidden Patterns (All Languages)

| Pattern | Why |
|---------|-----|
| Spaces in filenames | Breaks CLI tools, requires escaping |
| Mixed case folders | Inconsistent across OS (macOS case-insensitive, Linux case-sensitive) |
| Numeric-only names | No semantic meaning (`1.go`, `2.ts`) |
| Trailing hyphens/underscores | `user-.ts`, `handler_.go` |
| Double separators | `user--profile.ts`, `http__handler.go` |

---

## PowerShell Naming Convention

PowerShell scripts and modules use **lowercase kebab-case** for file names, NOT PascalCase.

### File Naming Rules

| Rule | Convention | Example |
|------|-----------|---------|
| Script files | `lowercase-kebab-case.ps1` | `upload-plugin.ps1`, `run-validator.ps1` |
| Module files | `lowercase-kebab-case.psm1` | `site-health.psm1` |
| Manifest files | `lowercase-kebab-case.psd1` | `site-health.psd1` |

### Function & Cmdlet Names (Inside Scripts)

Functions and cmdlets inside `.ps1` files follow the standard PowerShell **Verb-Noun** pattern with **PascalCase**:

| Element | Convention | Example |
|---------|-----------|---------|
| Functions | `Verb-Noun` (PascalCase) | `Get-ServiceStatus`, `Set-PluginConfig` |
| Verbs | Use only [Approved Verbs](https://learn.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands) | `Get`, `Set`, `New`, `Remove`, `Invoke` |
| Nouns | Singular, specific | `ServiceStatus` (not `ServiceStatuses`) |

### Examples

```
✅ Correct file names:
  run.ps1
  upload-plugin.ps1
  validate-guidelines.ps1
  check-site-health.ps1

❌ Incorrect file names:
  UploadPlugin.ps1          # PascalCase — forbidden for file names
  Run-Validator.ps1         # Verb-Noun pattern is for functions, not file names
  upload_plugin.ps1         # Underscores — use hyphens
  Upload-Plugin.ps1         # Mixed case — keep lowercase

✅ Correct function names (inside scripts):
  function Get-ServiceStatus { ... }
  function Set-PluginConfig { ... }

❌ Incorrect function names:
  function get-service-status { ... }   # Functions must be PascalCase Verb-Noun
  function getServiceStatus { ... }     # camelCase not allowed for PowerShell functions
```

### Summary

| Element | Convention |
|---------|-----------|
| **File names** (`.ps1`, `.psm1`, `.psd1`) | `lowercase-kebab-case` |
| **Function names** (inside scripts) | `PascalCase Verb-Noun` |
| **Folders** | `lowercase` (universal rule) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/08-file-folder-naming/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-FILE-002: Universal Cross-Language Naming Invariants

**Given** Repository file and directory structures across polyglot stacks.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase or invalid naming patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.
