# 13 — Dependency Management

Adding, updating, auditing, and vendoring third-party dependencies.

## Adding Dependencies

Evaluate: active maintenance, compatible license (MIT/BSD/Apache), minimal transitive deps, standard library alternatives.

## Version Pinning

Pin every dependency to exact versions. No caret (`^`) or tilde (`~`) ranges in production. CI tool installs pinned: `golangci-lint@v1.64.8`, `govulncheck@v1.1.4`. `@latest` is prohibited.

## Lock Files

Always commit lock files. Never manually edit. CI uses frozen installs.

## Update Cadence

| Category | Frequency |
|----------|-----------|
| Security fixes | Immediately |
| Patch versions | Weekly |
| Minor versions | Monthly |
| Major versions | Quarterly |

## Audit Process

Run `npm audit` / `govulncheck` in CI on every PR. Critical: patch within 24h. High: within 72h.

## License Compliance

Allow-list: MIT, BSD, Apache 2.0. Flag copyleft or unknown licenses for manual review.

## Vendoring

Go projects vendor all dependencies. Node.js relies on lock files. Commit `vendor/` for Go.

## Hygiene

Remove unused dependencies before each release. Prefer zero-dep libraries.

---

Source: `spec/05-coding-guidelines/13-dependency-management.md`
