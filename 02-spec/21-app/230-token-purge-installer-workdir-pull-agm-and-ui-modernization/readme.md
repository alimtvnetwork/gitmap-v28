# Spec 230: Security Token Purge, Installer Navigation, Pull Remediation, AGM Linux Update, Multi-Instance API & UI Modernization

> **Module:** `02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization`  
> **Status:** Active (In Progress)  
> **Version:** `6.482.0`  
> **Target Package:** `alimtvnetwork/gitmap-v28`  
> **Author:** MD Alim Ul Karim (Riseup Asia LLC)  

---

## 1. Executive Summary

Spec 230 addresses 8 critical deliverables:
1. **Refresh Token Security Purge:** Machine-wide audit and safe removal of exposed OAuth refresh tokens in unencrypted backup JSON files.
2. **GitMap CD Workdir Expansion:** Supporting `$work`, `$def`, and multiple persistent work directories in `gitmap cd` and shell completion wrappers.
3. **Pull Auto-Remediation:** Enhancing `gitmap pull` (`gitmap pa`) with automatic fast-forward fallback, safe merge conflict abortion, and actionable remediation prompts.
4. **AGM Linux Update Engine:** Silencing verbose stdout on update success, displaying detailed bounded stack traces on failure, and offering interactive resolution for batch failures.
5. **Interim Testing & Telemetry Audit:** Targeted unit verification of recent commits across `cli/cmdpull/`, `cli/cmdpipeline/`, and `cli/scanner/`.
6. **Antigravity Multi-Instance API:** Instance-filtered querying of running prompts, in-queue prompts, and recent projects via structured JSON flags and REST API endpoints.
7. **Settings UI Modernization:** Refreshing the settings UI with dark theme design tokens and compiling an authoritative index of reference UI/UX and CSS3 repositories.
8. **Codebase Image Audit & Git Scrub:** Cataloging 161 repository images, planning email redaction, removing text-only screenshots, and planning commit history cleanup.

---

## 2. Specification File Inventory

| File | Description | Ownership |
|:---|:---|:---:|
| [`00-master-audit-ledger.md`](00-master-audit-ledger.md) | Central audit ledger & task matrix | Lead |
| [`01-architecture-spec.md`](01-architecture-spec.md) | Security, installer workdir, pull remediation & AGM update spec | Worker 01 |
| [`02-component-and-cli-spec.md`](02-component-and-cli-spec.md) | Multi-instance API, settings UI design system & image audit | Worker 02 |
| [`99-consistency-report.md`](99-consistency-report.md) | Spec validation, link checks & naming compliance | Lead |
| [`readme.md`](readme.md) | Canonical spec module entry point | Lead |

---

## 3. Cross References

- Parent Plan: [Plan 230](../../../.ai-memory/plans/pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)
- Subtasks Directory: [Plan 230 Subtasks](../../../.ai-memory/plans/subtasks/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/)
- Design System: [02-spec/07-design-system/](../../07-design-system/readme.md)
- Coding Guidelines: [02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md](../../17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md)
