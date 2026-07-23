# Post-Install Shell Activation — Generic CLI Spec

> **Related specs:**
> - [11-build-deploy.md](11-build-deploy.md) — install/deploy step that places the binary on PATH
> - [19-shell-completion.md](19-shell-completion.md) — completion install uses the same profile-injection pattern
> - [13-checklist.md](13-checklist.md) — implementation phases that include setup
> - App reference: [spec/01-app/31-cd.md](../01-app/31-cd.md) — gitmap-v28 navigation helper that consumes this contract
> - Issue references: [spec/02-app-issues/22-installer-path-not-active-after-install.md](../02-app-issues/22-installer-path-not-active-after-install.md), [24-cd-command-does-not-change-shell-directory.md](../02-app-issues/24-cd-command-does-not-change-shell-directory.md), [25-powershell-cd-wrapper-not-loaded.md](../02-app-issues/25-powershell-cd-wrapper-not-loaded.md)
> - Consolidated guideline: [spec/12-consolidated-guidelines/19-post-install-shell-activation.md](../12-consolidated-guidelines/19-post-install-shell-activation.md)

## Overview

After `setup` (or the bootstrap installer) runs, the user MUST be able
to invoke the CLI **and any of its shell-integrated subcommands** in
the **current terminal session** without restarting it. This spec is
split into focused sub-documents so each concern stays under 200 lines
and is independently auditable by both humans and AI agents.

---

## Sub-Documents

| File | Scope |
|------|-------|
| [21-post-install-shell-activation/01-contract.md](21-post-install-shell-activation/01-contract.md) | Purpose, required behaviours (PIA-1..PIA-7), activation flow, in-session activation, shell detection, stderr warnings, global constraints. |
| [21-post-install-shell-activation/02-snippets.md](21-post-install-shell-activation/02-snippets.md) | Per-shell profile snippet bodies (PowerShell, Bash, Zsh, Fish), marker conventions, cross-platform parity table. |
| [21-post-install-shell-activation/03-doctor.md](21-post-install-shell-activation/03-doctor.md) | `doctor` wrapper check, three-state detection algorithm, implementation checklist for new CLIs. |
| [21-post-install-shell-activation/04-idempotency.md](21-post-install-shell-activation/04-idempotency.md) | Rewrite & removal algorithms, version-bump rules, testing requirements. |

---

## Quick Reference

- **Detection variable:** `<TOOL>_WRAPPER=1` (uppercased tool name).
- **Marker prefix:** `# <tool> shell wrapper v<N>`.
- **Marker suffix:** `# <tool> shell wrapper v<N> end`.
- **In-session activation:** PowerShell dot-sources `$PROFILE`; Unix
  shells print a one-liner (child cannot source parent).
- **Doctor states:** `LOADED`, `INSTALLED_BUT_NOT_LOADED`,
  `NOT_INSTALLED`.

---

## Why This Spec Exists

These bugs in the gitmap-v28 project triggered this generic spec:

| Issue | Root Cause |
|-------|------------|
| [22-installer-path-not-active-after-install](../02-app-issues/22-installer-path-not-active-after-install.md) | Installer wrote to PATH but never told the user to reload, and never auto-activated. |
| [24-cd-command-does-not-change-shell-directory](../02-app-issues/24-cd-command-does-not-change-shell-directory.md) | The `cd` subcommand silently no-op'd because the wrapper function was not loaded. |
| [25-powershell-cd-wrapper-not-loaded](../02-app-issues/25-powershell-cd-wrapper-not-loaded.md) | Same as 24 but on Windows — wrapper installed in `$PROFILE` but the running session never sourced it. |

By following the contract above, every new CLI in this framework
inherits a deterministic, AI-implementable post-install activation
flow on day one.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
