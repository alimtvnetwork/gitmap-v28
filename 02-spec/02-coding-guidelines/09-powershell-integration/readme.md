# PowerShell Integration (AI Execution Prompt)

> **/goal** Establish standards, scripting conventions, and architectural best practices for cross-platform PowerShell scripting and CI/CD automation.
> **/learn** Master cross-platform PowerShell core principles (pwsh 7+), strict parameter validation, PascalCase cmdlet naming, structured error records, and zero-leak environment variables.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Ensure all PowerShell scripts target cross-platform PowerShell Core (`pwsh`) with strict strict-mode and error action preferences (`$ErrorActionPreference = 'Stop'`).
- [ ] `/learn` Avoid Windows-only cmdlets, hardcoded backslashes, or assuming registry/WMI availability on Linux/macOS.
- [ ] `/goal` Enforce approved verb-noun cmdlet and function naming conventions with typed parameter blocks (`[CmdletBinding()]`).
- [ ] `/learn` Validate PowerShell integration specifications using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/09-powershell-integration --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

PowerShell integration guidelines, scripting conventions, and best practices for cross-platform automation within the project ecosystem.

---

## Contents

_No content yet. Add PowerShell-related specs as numbered files within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Guidelines | [../01-cross-language/readme.md](../01-cross-language/readme.md) |
| Coding Guidelines Spec | [../readme.md](../readme.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PWSH-000: PowerShell Integration Specification Conformance

**Given** Polyglot development guidelines and language standards.
**When** Audited against this language specification.
**Then** Zero non-compliant conventions or syntax patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/09-powershell-integration --check-only
```
**Expected:** exit 0. Zero violations.
