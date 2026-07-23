# Dependency Management

Dependencies must be pinned to exact versions without range operators (enforced by 'save-exact=true' in .npmrc and version-pinned tool installs). Audit workflows require automated vulnerability scanning (e.g., 'npm audit', 'govulncheck'). License compliance permits only permissive licenses like MIT and Apache-2.0, while strictly prohibiting strong copyleft licenses (GPL/AGPL) in proprietary code. Strict version control for Axios blocks compromised versions and mandates exact pinning to 1.14.0 or 0.30.3.

## Pinned Tool Versions

| Tool | Version | Used In |
|------|---------|---------|
| `golangci-lint` | `v1.64.8` | `setup.sh`, `ci.yml` |
| `govulncheck` | `v1.1.4` | `ci.yml`, `vulncheck.yml` |

All CI tool installs use exact version tags — `go install tool@latest` is prohibited for reproducibility. See `spec/05-coding-guidelines/17-cicd-patterns.md`.
