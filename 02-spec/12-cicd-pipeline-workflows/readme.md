# 09 — Pipeline Specifications

Generic, portable documentation for the project's CI/CD pipeline architecture. These specs describe **what** each pipeline does, **why** each pattern exists, and **how** to implement it — in enough detail for any AI or engineer to reproduce the workflows from scratch.

---

## Documents

| Document | Purpose |
|----------|---------|
| [01-ci-pipeline.md](./01-ci-pipeline.md) | Continuous integration: lint, vulnerability scan, parallel tests, cross-compiled builds |
| [02-release-pipeline.md](./02-release-pipeline.md) | Release automation: version resolution, binary packaging, install scripts, GitHub releases |
| [03-vulnerability-scanning.md](./03-vulnerability-scanning.md) | Standalone vulnerability scanning: scheduled and manual |
| [04-installation-flow.md](./04-installation-flow.md) | End-to-end installation: one-liner scripts, terminal output, upgrade, uninstall |
| [05-changelog-integration.md](./05-changelog-integration.md) | Changelog format, CI extraction, release body assembly, terminal display |
| [06-version-and-help.md](./06-version-and-help.md) | Version display, help system, command-level docs, CI version verification |
| [07-environment-variable-setup.md](./07-environment-variable-setup.md) | `env` command: persistent variables, PATH registration, auto-home, drive setup |
| [08-terminal-output-standards.md](./08-terminal-output-standards.md) | Output formatting conventions: icons, tables, progress, errors, CI summaries |
| [09-binary-icon-branding.md](./09-binary-icon-branding.md) | Windows binary icon embedding via `go-winres`: icon, manifest, version info |
| [03-reusable-ci-guards/00-overview.md](./03-reusable-ci-guards/00-overview.md) | **Reusable CI guards** — six language-agnostic patterns (forbidden-name, grandfather baseline, collision audit, baseline-diff lint gate, lint suggestions, matrix test aggregator) with Go/Node/Python/Rust adaptations and an AI implementation guide |

---

## CI/CD Pipeline Diagram

See the Mermaid diagram: [`images/ci-pipeline-flow.mmd`](images/ci-pipeline-flow.mmd)

## Unified Architecture Diagram

See the Mermaid diagram: [`images/unified-architecture.mmd`](images/unified-architecture.mmd)

Shows how all nine pipeline specs connect — from CI validation through
release automation, installation, changelog, versioning, environment setup,
terminal standards, and binary branding.

---

## Quick Reference

### Pipeline Triggers

| Workflow | Trigger | Branch/Tag |
|----------|---------|------------|
| CI | Push, Pull Request | `main` |
| Release | Push | `release/**`, `v*` tags |
| Vulnerability Scan | Weekly schedule, Manual | Any (default branch) |

### Shared Conventions

- **Platform**: GitHub Actions
- **Runner**: `ubuntu-latest`
- **Language toolchain**: Go (version from `go.mod`)
- **Node.js compatibility**: `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true` environment variable
- **Action versions**: Pinned to exact tags (e.g., `@v6`), never `@latest` or `@main`
- **Tool versions**: Pinned to exact versions (e.g., `golangci-lint@v1.64.8`, `govulncheck@v1.1.4`)
- **Build mode**: Static linking (`CGO_ENABLED=0`) for all binaries
- **Cross-compilation targets**: `windows/amd64`, `windows/arm64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`

### Pinned Tool Versions

| Tool | Version | Used In |
|------|---------|---------|
| `golangci-lint` | `v1.64.8` | CI pipeline |
| `govulncheck` | `v1.1.4` | CI pipeline, Vulnerability scan |
| `actions/checkout` | `@v6` | All workflows |
| `actions/setup-go` | `@v6` | All workflows |
| `actions/cache` | `@v4` | CI pipeline |
| `actions/upload-artifact` | `@v4` | CI, Release |
| `actions/download-artifact` | `@v4` | CI pipeline |
| `softprops/action-gh-release` | `@v2` | Release pipeline |
| `golangci/golangci-lint-action` | `@v6` | CI pipeline |

### AI Handoff Checklist

When handing this project to any AI or engineer, they should read these specs in order:

1. **08** — Terminal output standards (understand the visual conventions first)
2. **06** — Version and help system (how commands present themselves)
3. **04** — Installation flow (how users get the tool)
4. **07** — Environment variable setup (how the tool configures the system)
5. **05** — Changelog integration (how changes are tracked and published)
6. **01** — CI pipeline (how code is validated)
7. **02** — Release pipeline (how releases are built and published)
8. **03** — Vulnerability scanning (security baseline)
