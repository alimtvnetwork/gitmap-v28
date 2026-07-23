# Dependency Management

Standards for version pinning, audit workflows, upgrade strategies, and license compliance across all languages.

---

## 1. Version Pinning

### Principles

- Pin every dependency to an exact version in production.
- Use lock files to guarantee reproducible builds.
- Never use floating ranges (`latest`, `*`, `>=`) in production manifests.

### Lock Files

| Language | Lock File | Install Command |
|---|---|---|
| Go | `go.sum` | `go mod download` |
| TypeScript | `package-lock.json` / `bun.lock` | `npm ci` / `bun install --frozen-lockfile` |
| Python | `requirements.txt` (pinned) | `pip install -r requirements.txt` |

### Go

```go
// go.mod — pin exact versions
require (
    github.com/spf13/cobra v1.8.0
    github.com/mattn/go-sqlite3 v1.14.22
)
```

### TypeScript

```json
// package.json — use exact versions
{
  "dependencies": {
    "react": "18.3.1",
    "zod": "3.25.76"
  }
}
```

### Compromised Packages — Do Not Upgrade

Some package versions have been compromised through supply-chain attacks.
Pin to the last known safe version and **never upgrade past it** without
a verified security review.

| Package | Bad Versions | Safe Version | Notes |
|---|---|---|---|
| `axios` | `1.14.1` | `1.14.0` | Supply-chain compromise — do not install `1.14.1` |
| `axios` | `0.30.4` | `0.30.3` | Supply-chain compromise — do not install `0.30.4` |

> **WARNING**: Versions `axios@1.14.1` and `axios@0.30.4` were published
> as part of a supply-chain attack. Always pin to `1.14.0` or `0.30.3`
> respectively. Use exact versions (no `^` or `~`) to prevent automatic
> upgrades to the compromised releases.

### Rules

| Rule | Detail |
|---|---|
| Exact versions in manifests | No `^`, `~`, or `>=` in production |
| Lock files are committed | Always version-control lock files |
| CI uses frozen installs | `npm ci`, `--frozen-lockfile`, `go mod download` |
| No manual lock file edits | Only tools should modify lock files |
| Pin CI action versions | Use full SHA or exact tag, never `@latest` or `@main` |

### CI Tool Versions

Pin all CI tool installs to exact version tags — `go install tool@latest` is prohibited:

| Tool | Pinned Version | Used In |
|------|---------------|---------|
| `golangci-lint` | `v1.64.8` | `setup.sh`, `ci.yml` |
| `govulncheck` | `v1.1.4` | `ci.yml`, `vulncheck.yml` |

---

## 2. Audit Workflows

### Regular Audits

Run dependency audits on every CI build and on a weekly schedule:

```bash
# Go
go list -m -json all | go-licenses check .
govulncheck ./...

# TypeScript
npm audit --audit-level=high
```

### Automated Audit Pipeline

```yaml
# .github/workflows/audit.yml
name: Dependency Audit
on:
  schedule:
    - cron: "0 9 * * 1"  # Weekly Monday 9 AM
  push:
    paths:
      - "go.sum"
      - "package-lock.json"

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - run: npm audit --audit-level=high
      - run: govulncheck ./...
```

### Vulnerability Response

| Severity | Response Time | Action |
|---|---|---|
| Critical | 24 hours | Patch or remove immediately |
| High | 72 hours | Patch in next release |
| Medium | 2 weeks | Schedule upgrade |
| Low | Next upgrade cycle | Bundle with routine updates |

### Rules

| Rule | Detail |
|---|---|
| Block CI on critical/high | Builds fail if unresolved critical or high vulnerabilities exist |
| Track advisories | Subscribe to security advisories for key dependencies |
| Document exceptions | If a vulnerability is accepted, document the reason and review date |
| Audit transitive deps | Vulnerabilities in indirect dependencies count equally |

---

## 3. Upgrade Strategies

### Routine Upgrades

Perform dependency upgrades on a regular cadence:

| Cadence | Scope |
|---|---|
| Weekly | Patch versions (bug fixes, security) |
| Monthly | Minor versions (new features, non-breaking) |
| Quarterly | Major versions (breaking changes) |

### Upgrade Process

```
1. Create a branch for the upgrade
2. Update dependency version
3. Run full test suite
4. Review changelog for breaking changes
5. Test in staging environment
6. Merge with documented rationale
```

### Go

```bash
# Check for available updates
go list -m -u all

# Update a specific dependency
go get github.com/spf13/cobra@v1.9.0
go mod tidy

# Update all patch versions
go get -u=patch ./...
go mod tidy
```

### TypeScript

```bash
# Check for outdated packages
npm outdated

# Update a specific package
npm install react@18.4.0

# Update all within semver range
npm update
```

### Major Version Upgrades

Major versions require additional care:

| Step | Detail |
|---|---|
| Read migration guide | Review the changelog and breaking changes |
| Check compatibility | Verify peer dependencies and related packages |
| Update incrementally | One major upgrade per PR |
| Run full regression | Execute all tests including integration |
| Monitor after deploy | Watch error rates for 48 hours post-deploy |

### Rules

| Rule | Detail |
|---|---|
| One dependency per PR | Major upgrades get isolated PRs |
| Changelog review is mandatory | Never upgrade blindly |
| Test before merging | Full test suite must pass |
| No vendored forks | Prefer upstream patches over local forks |
| Document breaking changes | Note migration steps in PR description |

---

## 4. License Compliance

### Approved Licenses

| Category | Licenses | Usage |
|---|---|---|
| Permissive (preferred) | MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, ISC | Unrestricted use |
| Weak copyleft (review required) | LGPL-2.1, LGPL-3.0, MPL-2.0 | Allowed with isolation |
| Strong copyleft (prohibited) | GPL-2.0, GPL-3.0, AGPL-3.0 | Not allowed in proprietary projects |
| No license | Unlicensed | Not allowed — treat as all rights reserved |

### Automated License Checking

```bash
# Go
go-licenses check ./... --allowed_licenses=MIT,Apache-2.0,BSD-2-Clause,BSD-3-Clause,ISC

# TypeScript
npx license-checker --onlyAllow "MIT;Apache-2.0;BSD-2-Clause;BSD-3-Clause;ISC"
```

### Go

```go
// Verify licenses in CI
// go-licenses csv ./... > licenses.csv
// Review and approve new entries before merging
```

### TypeScript

```typescript
// license-check.ts
import { exec } from "child_process";

const allowedLicenses = ["MIT", "Apache-2.0", "BSD-2-Clause", "BSD-3-Clause", "ISC"];

function checkLicenses(): void {
  exec("npx license-checker --json", (err, stdout) => {
    const packages = JSON.parse(stdout);
    for (const [pkg, info] of Object.entries(packages)) {
      const license = (info as { licenses: string }).licenses;
      if (!allowedLicenses.includes(license)) {
        console.error(`Unapproved license: ${pkg} uses ${license}`);
        process.exit(1);
      }
    }
  });
}
```

### Rules

| Rule | Detail |
|---|---|
| CI blocks unapproved licenses | Automated check on every dependency change |
| Review new dependencies | Every new dependency requires license verification |
| Maintain an allow-list | Centralized list of approved licenses |
| Document exceptions | If a non-standard license is accepted, record the reason |
| Check transitive licenses | Indirect dependencies must also comply |

---

## Constraints

| Constraint | Detail |
|---|---|
| Exact versions only | No floating ranges in production |
| Lock files committed | Always version-controlled |
| CI audits on every build | Block on critical and high vulnerabilities |
| One major upgrade per PR | Isolate breaking changes |
| License allow-list enforced | CI blocks unapproved licenses |
| Changelog review mandatory | Never upgrade without reading the changelog |
| No unlicensed dependencies | Treat missing license as all rights reserved |
| Weekly patch cadence | Security and bug-fix updates applied promptly |

---

## References

- [Security & Secrets](./08-security-secrets.md)
- [CI/CD Patterns](./17-cicd-patterns.md)
- [Code Review](./25-code-review.md)
- [Configuration Management](./21-configuration-management.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
