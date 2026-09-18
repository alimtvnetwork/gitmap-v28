# 08 — Generic Self-Update Specification

## Purpose

This folder defines a **generic, reusable blueprint** for implementing
self-update functionality in CLI tools. It covers the full lifecycle:
detecting where the binary is installed, building a new version,
deploying it without file-lock errors, and cleaning up afterward.

Any AI or engineer reading these documents should be able to implement
a complete self-update system from scratch without ambiguity.

---

## Documents

| File | Topic |
|------|-------|
| [01-self-update-overview.md](01-self-update-overview.md) | The problem, approach, and platform differences |
| [02-deploy-path-resolution.md](02-deploy-path-resolution.md) | Deploy to running location, PATH registration, data co-location |
| [03-rename-first-deploy.md](03-rename-first-deploy.md) | Rename-first strategy to bypass file locks |
| [04-build-scripts.md](04-build-scripts.md) | `run.ps1` and `run.sh` patterns for build + deploy |
| [05-handoff-mechanism.md](05-handoff-mechanism.md) | Copy-and-handoff for Windows self-replacement |
| [06-cleanup.md](06-cleanup.md) | Post-update artifact removal |
| [07-console-safe-handoff.md](07-console-safe-handoff.md) | Prevent async handoff from breaking the console session |

---

## Unified Architecture Diagram

See how all specs connect: [`images/unified-architecture.mmd`](images/unified-architecture.mmd)

## Self-Update Flow Diagram

See the Mermaid diagram: [`images/self-update-flow.mmd`](images/self-update-flow.mmd)

```
<binary> update
    |
    +-- Resolve deploy target (running location > PATH > config)
    +-- Source repo available?
    |   +-- YES: Pull, build, deploy (rename-first)
    |   +-- NO:  Download pre-built binary from releases
    |
    +-- Deploy to running executable's directory
    |   +-- Rename existing binary to .old
    |   +-- Copy new binary
    |   +-- Copy data/ folder alongside binary
    |
    +-- Register directory in PATH (if first-time)
    +-- Reload terminal environment
    +-- Verify: <binary> version == expected
    +-- Cleanup temporary files
```

---

## Core Principles

1. **Deploy to running location** — The binary is always replaced
   in-place at the directory it is currently running from.
2. **Data co-location** — The `data/` folder lives alongside the
   binary at its physical location. Moving the binary moves data too.
3. **PATH auto-registration** — On first-time install, the deploy
   directory is added to PATH and the terminal is reloaded.
4. **Rename-first** — On Windows, rename the locked binary before
   replacing it. Never rely on overwrite retries alone.

## Placeholders

| Placeholder | Meaning |
|-------------|---------|
| `<binary>` | Your CLI binary name (e.g., `mytool`) |
| `<binary>.exe` | Windows binary with extension |
| `<deploy-dir>` | The directory where the binary is installed |
| `<repo-root>` | The root of the source repository |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
