# Generic Release Pipeline Specification

> **Version:** 1.0.0  
> **Updated:** 2026-04-20  
> **Status:** Active  
> **Imported from:** sibling reference implementation `spec/16-generic-release`
>
> **Related local specs:**
> - [`../12-cicd-pipeline-workflows/02-release-pipeline.md`](../12-cicd-pipeline-workflows/02-release-pipeline.md) — this repo's concrete release workflow (consumes the generic contract below)
> - [`../12-cicd-pipeline-workflows/10-release-pipeline-issues-rca.md`](../12-cicd-pipeline-workflows/10-release-pipeline-issues-rca.md) — local RCA ledger of release-pipeline failures
> - [`../13-generic-cli/20-terminal-output-design.md`](../13-generic-cli/20-terminal-output-design.md) — terminal output contract used by install scripts
> - [`../13-generic-cli/21-post-install-shell-activation.md`](../13-generic-cli/21-post-install-shell-activation.md) — post-install PATH/profile/wrapper activation contract

## Purpose

This folder defines a **generic, reusable blueprint** for releasing
cross-compiled CLI binaries via CI/CD. It is tool-agnostic — replace
placeholder names with your actual binary name and repository URL.

Any AI or engineer reading these documents should be able to implement
a complete release pipeline from scratch without ambiguity.

---

## Documents

| File | Topic |
|------|-------|
| [01-cross-compilation.md](01-cross-compilation.md) | Building static binaries for 6+ platform targets |
| [02-release-pipeline.md](02-release-pipeline.md) | CI/CD workflow structure, triggers, and stages |
| [03-install-scripts.md](03-install-scripts.md) | Generating version-pinned PowerShell and Bash installers |
| [04-checksums-verification.md](04-checksums-verification.md) | SHA-256 checksum generation and verification |
| [05-release-assets.md](05-release-assets.md) | Asset naming, compression, and packaging conventions |
| [06-release-metadata.md](06-release-metadata.md) | Version resolution, tagging, and changelog extraction |
| [07-known-issues-and-fixes.md](07-known-issues-and-fixes.md) | Post-mortem catalog: every release-pipeline failure with root cause, fix, and prevention rule |
| [08-version-pinned-release-installers.md](08-version-pinned-release-installers.md) | **Authoritative contract** for the per-release `install.sh` / `install.ps1` assets — spec-first ordering, embedded version, no "latest" probe |

---

## Release Pipeline Diagram

See the Mermaid diagram: [`images/release-pipeline-flow.mmd`](images/release-pipeline-flow.mmd)

## Unified Architecture Diagram

See the Mermaid diagram: [`images/unified-architecture.mmd`](images/unified-architecture.mmd)

Shows how all six specs connect — from cross-compilation through packaging,
checksums, install scripts, and metadata into the final GitHub Release.

---

## Shared Conventions

- **Build once, package once** — binaries are compiled exactly once;
  all downstream steps (compress, checksum, publish) reuse the same
  artifacts and must never trigger a rebuild.
- **Pin all tool versions** — never use `@latest` or `@main` for
  CI actions or tool installs. Use exact version tags.
- **Static linking** — use `CGO_ENABLED=0` for Go binaries to produce
  fully static executables with no runtime dependencies.
- **Deterministic builds** — identical source + identical toolchain =
  identical output. Lock dependency versions via lock files.

## Placeholders

Throughout these documents:

| Placeholder | Meaning |
|-------------|---------|
| `<binary>` | Your CLI binary name (e.g., `mytool`) |
| `<repo>` | Your repository path (e.g., `github.com/org/repo`) |
| `<version>` | The release version (e.g., `v1.2.0`) |
| `<module>` | Your Go module path |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)

---

## Verification

_Auto-generated section — see `spec/16-generic-release/97-acceptance-criteria.md` for the full criteria index._

### AC-REL-000: Generic-release conformance: Overview

**Given** Inspect a release artifact bundle for required assets and checksums.  
**When** Run the verification command shown below.  
**Then** SHA-256 checksums verify; `release-metadata.json` matches the package version; install scripts pin the exact release tag.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
