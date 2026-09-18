# Static Analysis & Security

Code quality and security are enforced through automated linting and scanning:

1. **golangci-lint** is configured with 28 enabled linters via `.golangci.yml` (including errcheck, gosec, and gocritic), with specific exclusions for `gosec` noise (G101, G104, G304, etc.) and linter-specific overrides for test files and the `cmd/` directory. Pinned to `v1.64.8` in `setup.sh` and CI workflows.

2. **govulncheck** identifies dependency vulnerabilities and runs on every CI build as well as a weekly schedule (Mondays at 9:00 UTC). Pinned to `v1.1.4` in `ci.yml` and `vulncheck.yml` for reproducible builds — `@latest` is prohibited. The CI pipeline treats unfixable standard library vulnerabilities as warnings while failing only on third-party package vulnerabilities (`packages you import`).

Dependencies are pinned to exact versions, and license compliance permits only permissive licenses like MIT and Apache-2.0.
