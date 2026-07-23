# PowerShell Build & Deploy Patterns

## Overview

This document describes reusable patterns for PowerShell build scripts
that manage the lifecycle of compiled CLI tools: pull, build, deploy,
and run.

This specification is split into focused sub-documents:

| Document | Topic |
|----------|-------|
| [02a-script-architecture.md](02a-script-architecture.md) | Entry point, step-based execution, configuration, error handling |
| [02b-logging-patterns.md](02b-logging-patterns.md) | Semantic logging functions and banner display |
| [02c-build-patterns.md](02c-build-patterns.md) | Embedded variables, version verification, data folder copy |
| [02d-deploy-patterns.md](02d-deploy-patterns.md) | Retry-on-lock, nested deploy structure, PATH integration |
| [02e-run-pattern.md](02e-run-pattern.md) | `-R` flag, argument forwarding, path resolution |
| [02f-self-update-orchestration.md](02f-self-update-orchestration.md) | Two-phase handoff, rename-first sync, validation rules |
| [02g-last-release-detection.md](02g-last-release-detection.md) | `Get-LastRelease.ps1`, three-tier fallback, integration points |

## Cross-References (Generic Specifications)

This document is an application-level summary. The following generic,
tool-agnostic specs provide detailed breakdowns of each mechanism:

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Build pipeline | [04-build-scripts.md](../08-generic-update/04-build-scripts.md) | `run.ps1` / `run.sh` full pipeline, config loading, ldflags, `--force-pull`, logging helpers |
| Deploy strategy | [03-rename-first-deploy.md](../08-generic-update/03-rename-first-deploy.md) | Rename-first flow, rollback, PATH sync, retry reduction (20→5) |
| Self-update orchestration | [05-handoff-mechanism.md](../08-generic-update/05-handoff-mechanism.md) | Copy-and-handoff, worker launch, UTF-8 BOM, binary-based fallback |
| Cleanup | [06-cleanup.md](../08-generic-update/06-cleanup.md) | `.old` lifecycle, `update-cleanup` command, temp directory hygiene |
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) | Cross-compilation, checksums, version-pinned install scripts |
| Install scripts | [03-install-scripts.md](../07-generic-release/03-install-scripts.md) | `install.ps1` / `install.sh` generation, SHA-256 verification |
| Release metadata | [06-release-metadata.md](../07-generic-release/06-release-metadata.md) | `releases.json` manifest, `baseUrl`, asset maps |

### Mapping: Sub-Documents → Generic Specs

| Sub-Document | Generic Equivalent |
|-------------|-------------------|
| 02a Script Architecture | `04-build-scripts.md` §PowerShell, §Config Loading |
| 02d Deploy Patterns | `03-rename-first-deploy.md` §Implementation |
| 02f Self-Update Orchestration | `05-handoff-mechanism.md` §Solution: Copy-and-Handoff |
| 02g Last Release Detection | `04-build-scripts.md` §Validation |
| 02a Error Handling | `03-rename-first-deploy.md` §Rollback |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
