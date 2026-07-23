# Dependency Management

## Overview

Standards for adding, updating, auditing, and vendoring third-party
dependencies to keep builds reproducible and supply chains secure.

---

## Adding Dependencies

### Evaluation Criteria

Before adding a new dependency, verify:

| Criterion          | Minimum Bar                              |
|--------------------|------------------------------------------|
| Maintenance        | Active commits within the last 6 months  |
| License            | Compatible open-source license (MIT, BSD, Apache 2.0) |
| Transitive deps    | Fewer is better — prefer zero-dep libs   |
| Binary size impact | Measure before and after                 |
| Alternatives       | Standard library solution preferred      |

### One Purpose Per Package

Each dependency solves exactly one problem. Do not add a large framework
when a focused library (or standard library) covers the need.

---

## Version Pinning

### Pin Exact Versions

Lock every dependency to an exact version in the manifest:

```json
// package.json — pin exact
"dependencies": {
  "zod": "3.25.76"
}
```

```go
// go.mod — pin exact
require modernc.org/sqlite v1.29.6
```

### CI Tool Versions

Pin CI tool installs to exact version tags — `@latest` is prohibited:

| Tool | Pinned Version | Used In |
|------|---------------|---------|
| `golangci-lint` | `v1.64.8` | `setup.sh`, `ci.yml` |
| `govulncheck` | `v1.1.4` | `ci.yml`, `vulncheck.yml` |

```bash
# Correct — pinned
go install golang.org/x/vuln/cmd/govulncheck@v1.1.4

# Wrong — non-reproducible
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### Lock Files

- Always commit lock files (`bun.lock`, `go.sum`, `package-lock.json`).
- Never manually edit lock files — regenerate via the package manager.
- CI builds must use frozen installs (`npm ci`, `bun install --frozen-lockfile`).

### Range Specifiers

Avoid caret (`^`) and tilde (`~`) ranges in production dependencies.
Use them only in library packages where flexibility is intentional.

---

## Update Cadence

### Schedule

| Category       | Frequency      | Action                          |
|----------------|----------------|---------------------------------|
| Security fixes | Immediately    | Patch and deploy same day       |
| Patch versions | Weekly         | Review and apply in batch       |
| Minor versions | Monthly        | Test in CI before merging       |
| Major versions | Quarterly      | Evaluate changelog, plan migration |

### Update Process

1. Run the audit tool (`npm audit`, `go list -m -u all`).
2. Update one dependency at a time for major bumps.
3. Run the full test suite after each update.
4. Review the changelog for breaking changes.
5. Commit the updated manifest and lock file together.

---

## Audit Process

### Automated Scanning

Run dependency audits in CI on every pull request:

```bash
# Node.js
npm audit --audit-level=high

# Go
go list -m -u all
govulncheck ./...
```

### Severity Response

| Severity | Response Time | Action                    |
|----------|---------------|---------------------------|
| Critical | 24 hours      | Patch, test, deploy       |
| High     | 72 hours      | Patch in next release     |
| Medium   | 2 weeks       | Evaluate and schedule     |
| Low      | Next cycle    | Bundle with routine updates |

### License Compliance

Maintain an allow-list of approved licenses. Flag any dependency using
a copyleft or unknown license for manual review before merging.

---

## Vendoring Strategies

### Go Projects

Vendor all dependencies for hermetic, offline-capable builds:

```bash
go mod vendor
go build -mod=vendor ./...
```

- Commit the `vendor/` directory.
- CI builds use `-mod=vendor` to enforce vendored deps.
- Run `go mod tidy` before vendoring to remove unused modules.

### Node.js Projects

Do not vendor `node_modules`. Rely on lock files and frozen installs:

```bash
npm ci              # deterministic from lock file
bun install --frozen-lockfile
```

### When to Vendor

| Scenario                          | Vendor? |
|-----------------------------------|---------|
| Go CLI distributed as binary      | Yes     |
| Node.js app with CI/CD pipeline   | No      |
| Air-gapped or offline builds      | Yes     |
| Library published to registry     | No      |

---

## Dependency Hygiene

### Remove Unused Dependencies

Audit for unused packages regularly:

```bash
# Node.js
npx depcheck

# Go
go mod tidy
```

### Minimize Transitive Dependencies

Prefer libraries with fewer transitive dependencies. A smaller
dependency tree reduces attack surface and build times.

### Fork Policy

Fork a dependency only as a last resort. Document the fork reason
and track upstream for future reconciliation.

---

## Constraints

- All production dependencies pinned to exact versions.
- Lock files committed and enforced in CI.
- Security vulnerabilities at high or critical patched within 48 hours.
- No dependency added without evaluating standard library alternatives.
- Go projects vendor all dependencies.
- Unused dependencies removed before each release.
