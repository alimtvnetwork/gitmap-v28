# 07 — Generic Release Pipeline Specification

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
