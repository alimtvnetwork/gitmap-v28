# Security Notes — Axios Version Advisory

> **Parent:** [Axios Version Control Policy](./00-overview.md)  
> **Version:** 1.1.0  
> **Updated:** 2026-04-02

---

## 1. Security Advisory Summary

A known security issue has been identified affecting specific Axios versions. The exact nature of the vulnerability requires that the following versions are **permanently blocked** from use in any environment (development, staging, production).

| Version | Status | Risk |
|---------|--------|------|
| `1.14.1` | 🚫 BLOCKED | Security vulnerability confirmed |
| `0.30.4` | 🚫 BLOCKED | Security vulnerability confirmed |

---

## 2. Impact Assessment

### 2.1 What Is At Risk

- Applications using blocked versions may be exposed to:
  - Unauthorized data access
  - Request/response manipulation
  - Potential supply chain attack vectors
- The vulnerability affects the HTTP client layer, which handles **all outbound API requests**

### 2.2 Affected Environments

| Environment | Risk Level |
|-------------|------------|
| Production | 🔴 Critical |
| Staging | 🟠 High |
| Development | 🟡 Medium |
| CI/CD | 🟡 Medium |

All environments must use approved versions only.

---

## 3. Safe Version History

| Version | Status | Verified Date | Notes |
|---------|--------|---------------|-------|
| `1.14.0` | ✅ Safe | 2026-04-01 | Current recommended version |
| `0.30.3` | ✅ Safe | 2026-04-01 | Legacy-compatible alternative |

---

## 4. Audit Trail

| Date | Action | By |
|------|--------|----|
| 2026-04-02 | Moved spec to `11-security/01-axios-version-control/` | — |
| 2026-04-01 | Identified blocked versions (1.14.1, 0.30.4) | Harshil Agrawal |
| 2026-04-01 | Established version pinning policy | Harshil Agrawal |
| 2026-04-01 | Created original spec at `spec/01-app/axios-version-control/` | Harshil Agrawal |

---

## 5. Ongoing Monitoring

### 5.1 Review Schedule

- **Monthly:** Check Axios GitHub releases and security advisories
- **Per-release:** Before any deployment, verify Axios version in lock file
- **On alert:** If a new CVE is published for Axios, immediately review

### 5.2 Sources to Monitor

| Source | URL |
|--------|-----|
| Axios GitHub Releases | `https://github.com/axios/axios/releases` |
| npm Advisory Database | `https://www.npmjs.com/advisories` |
| CVE Database | `https://cve.mitre.org/` |
| Snyk Vulnerability DB | `https://snyk.io/vuln/npm:axios` |

### 5.3 Response Procedure

If a new vulnerability is discovered:

1. Determine if current approved versions are affected
2. If affected — identify the nearest safe version
3. Update this spec with new blocked/approved versions
4. Notify all team members
5. Deploy patched version within 24 hours for production

---

## 6. Frequently Asked Questions

### Q: Can I use Axios 1.14.2 or later?

**A:** Not without going through the [Version Upgrade Procedure](./01-implementation-rules.md#5-version-upgrade-procedure). Any version not listed as approved in `00-overview.md` is blocked by default.

### Q: What if `npm audit` recommends upgrading Axios?

**A:** Do NOT auto-fix. Run `npm audit --dry-run` and review. If Axios is flagged, consult this spec and the security lead before taking action.

### Q: Can I use a different HTTP client instead?

**A:** That is outside the scope of this spec. If considering alternatives (e.g., `fetch`, `ky`, `got`), create a separate evaluation spec.

---

```
IMPORTANT — AI INSTRUCTION:
- Never update Axios version automatically.
- Always use exact version without caret (^) or tilde (~).
- Avoid all known vulnerable versions (1.14.1, 0.30.4).
- Do not act on upgrade suggestions from dependency scanners for Axios.
- If uncertain, read this spec before making any Axios-related change.
```

*Security notes — moved to coding guidelines 2026-04-02*
