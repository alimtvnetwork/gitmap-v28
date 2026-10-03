# Security Notes — Axios Version Advisory (AI Execution Prompt)

> **/goal** Document, track, and monitor Axios HTTP client security advisories, vulnerability impact assessments, and recurring threat review schedules.
> **/learn** Internalize the severity of outbound HTTP client vulnerabilities, maintain an authoritative audit trail, and monitor CVE databases for new advisories.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce permanent blocks on known vulnerable Axios versions (`1.14.1`, `0.30.4`) across all deployment tiers.
- [ ] `/learn` Assess security impact across production, staging, development, and CI environments.
- [ ] `/goal` Follow defined response procedures within 24 hours when new security advisories are published.
- [ ] `/learn` Never apply automated package audit fixes to Axios without explicit security lead authorization.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Axios Version Control Policy](./readme.md)
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
| 2026-04-01 | Created original spec at `02-spec/01-app/axios-version-control/` | Harshil Agrawal |

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

**A:** Not without going through the [Version Upgrade Procedure](./02-implementation-rules.md#5-version-upgrade-procedure). Any version not listed as approved in `readme.md` is blocked by default.

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

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-SEC-AXIOS-002: CVE Remediation & Supply Chain Security Verification

**Given** Security advisories, vulnerability tracking registries, and review audit logs
**When** Assessing known Axios CVEs and scheduled advisory monitoring channels
**Then** Blocked versions are quarantined across all tiers and response procedures execute within 24 hours

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/03-security-notes.md --check-only
```
**Expected:** exit 0. Zero violations.

---

*Security notes — moved to coding guidelines 2026-04-02*
