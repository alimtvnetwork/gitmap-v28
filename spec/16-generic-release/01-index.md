# Generic Release Pipeline Specification

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 16 Generic Release.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `16-generic-release/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


> **Version:** 1.0.0
> **Updated:** 2026-04-20
> **Status:** Active
> **Imported from:** sibling reference implementation `02-spec/16-generic-release`
>
> **Related local specs:**
>
> - [`../12-cicd-pipeline-workflows/05-release-pipeline.md`](../12-cicd-pipeline-workflows/05-release-pipeline.md) — this repo's concrete release workflow (consumes the generic contract below)
> - [`../12-cicd-pipeline-workflows/17-release-pipeline-issues-rca.md`](../12-cicd-pipeline-workflows/17-release-pipeline-issues-rca.md) — local RCA ledger of release-pipeline failures
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
| [02-cross-compilation.md](02-cross-compilation.md) | Building static binaries for 6+ platform targets |
| [05-release-pipeline.md](../12-cicd-pipeline-workflows/05-release-pipeline.md) | CI/CD workflow structure, triggers, and stages |
| [04-install-scripts.md](04-install-scripts.md) | Generating version-pinned PowerShell and Bash installers |
| [05-checksums-verification.md](05-checksums-verification.md) | SHA-256 checksum generation and verification |
| [06-release-assets.md](06-release-assets.md) | Asset naming, compression, and packaging conventions |
| [07-release-metadata.md](07-release-metadata.md) | Version resolution, tagging, and changelog extraction |
| [08-known-issues-and-fixes.md](08-known-issues-and-fixes.md) | Post-mortem catalog: every release-pipeline failure with root cause, fix, and prevention rule |
| [09-version-pinned-release-installers.md](09-version-pinned-release-installers.md) | **Authoritative contract** for the per-release `install.sh` / `install.ps1` assets — spec-first ordering, embedded version, no "latest" probe |

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

_Auto-generated section — see `02-spec/16-generic-release/97-acceptance-criteria.md` for the full criteria index._

### AC-REL-001: Generic-release conformance: Index

**Given** Inspect a release artifact bundle for required assets and checksums.
**When** Run the verification command shown below.
**Then** SHA-256 checksums verify; `release-metadata.json` matches the package version; install scripts pin the exact release tag.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
