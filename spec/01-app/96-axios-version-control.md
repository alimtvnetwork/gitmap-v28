# Axios Version Control

## Overview

Strict version control policy for the Axios HTTP client library due to a
known supply-chain security compromise affecting specific versions. Only
pre-approved safe versions may be used. Automatic upgrades are prohibited.

---

## 1. Approved Safe Versions

| Version   | Status | Notes                                      |
|-----------|--------|--------------------------------------------|
| `1.14.0`  | Safe   | Last verified clean release (v1.x line)    |
| `0.30.3`  | Safe   | Last verified clean release (v0.x line)    |

Choose the version that matches your project's compatibility requirements.
For new projects, prefer `1.14.0`.

---

## 2. Blocked Versions

| Version   | Status  | Reason                              |
|-----------|---------|-------------------------------------|
| `1.14.1`  | Blocked | Supply-chain compromise             |
| `0.30.4`  | Blocked | Supply-chain compromise             |
| Any newer | Blocked | Not verified — requires manual review |

> **SECURITY WARNING**: Versions `1.14.1` and `0.30.4` were published as
> part of a supply-chain attack. Installing these versions may expose the
> application to arbitrary code execution and data exfiltration.

---

## 3. Version Pinning Rules

### Exact Version Only

Always declare Axios with an exact version — no caret (`^`), tilde (`~`),
or range specifiers:

```json
{
  "dependencies": {
    "axios": "1.14.0"
  }
}
```

### Prohibited Patterns

```json
// ❌ NEVER use these — they allow automatic upgrades to compromised versions
"axios": "^1.14.0"
"axios": "~1.14.0"
"axios": ">=1.14.0"
"axios": "latest"
"axios": "*"
```

### Lock File Enforcement

- Commit `package-lock.json` / `bun.lock` to version control.
- CI must use frozen installs (`npm ci` / `bun install --frozen-lockfile`).
- Never manually edit lock files.

---

## 4. Implementation Rules

### Dependency Declaration

- Specify Axios version exactly: `"axios": "1.14.0"`.
- Use `save-exact=true` in `.npmrc` to prevent range prefixes on install.
- Do not allow automated dependency update tools (Dependabot, Renovate)
  to modify the Axios version without manual approval.

### Code Review Enforcement

- Validate `package.json` during every pull request review.
- Reject any PR that changes the Axios version without explicit approval
  from a security-designated reviewer.
- CI should include a check that verifies the installed Axios version
  matches the approved list.

### CI Pipeline Safeguard

```bash
# Add to CI pipeline — fails if Axios version is not in the approved list
AXIOS_VERSION=$(node -e "console.log(require('axios/package.json').version)")
if [ "$AXIOS_VERSION" != "1.14.0" ] && [ "$AXIOS_VERSION" != "0.30.3" ]; then
  echo "ERROR: Axios version $AXIOS_VERSION is not approved. Allowed: 1.14.0, 0.30.3"
  exit 1
fi
```

### Monitoring and Logging

- Log dependency installation versions in CI output.
- Track any deviation from approved versions via audit alerts.
- Subscribe to Axios security advisories for future compromise notifications.

---

## 5. Security Note

A known supply-chain attack compromised Axios versions `1.14.1` and `0.30.4`.
These versions contain malicious code injected during the package publish
process. The attack vector allowed:

- Arbitrary code execution during installation or runtime.
- Potential data exfiltration from the host environment.

Until the Axios maintainers publish a verified clean release and a full
post-mortem is available, **no version beyond `1.14.0` (v1.x) or `0.30.3`
(v0.x) is approved for use**.

Any future upgrade must go through:

1. Manual verification of the published package contents.
2. Security team approval.
3. Full regression testing in a sandboxed environment.

---

## Constraints

| Constraint                              | Detail                                              |
|-----------------------------------------|-----------------------------------------------------|
| Exact version only                      | No `^`, `~`, `>=`, `*`, or `latest`                 |
| Blocked versions never installed        | `1.14.1` and `0.30.4` prohibited in all environments |
| No automatic upgrades                   | Dependabot / Renovate must skip Axios                |
| Code review gates Axios changes         | PR rejected without security reviewer approval       |
| CI validates installed version          | Pipeline fails on unapproved Axios version           |
| `.npmrc` enforces `save-exact=true`     | Prevents accidental range prefixes                   |
| Future upgrades require manual review   | No version beyond safe list without verification     |

---

## Acceptance Criteria

1. Axios version is always defined as an exact version without `^` or `~`.
2. No usage of blocked versions (`1.14.1`, `0.30.4`) is present in any environment.
3. Dependency updates do not alter the Axios version automatically.
4. Code reviews enforce strict compliance with this version policy.
5. Security note is documented and accessible to all developers.
6. CI pipeline includes an automated version check for Axios.

---

## References

- [Dependency Management](../05-coding-guidelines/27-dependency-management.md)
- [Security & Secrets](../05-coding-guidelines/08-security-secrets.md)
- [CI/CD Patterns](../05-coding-guidelines/17-cicd-patterns.md)
- [Code Review](../05-coding-guidelines/25-code-review.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
