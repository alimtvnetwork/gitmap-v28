# OWASP Top 10 Mitigation (AI Execution Prompt)

> **/goal** Defend applications against OWASP Top 10 vulnerabilities via parameterized queries, automated static analysis, XSS sanitization, and secure design patterns.
> **/learn** Never concatenate dynamic SQL queries, implement framework-native auto-escaping, enforce DOMPurify for raw HTML, and conduct proactive threat modeling.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce parameterized queries, prepared statements, and ORMs across all data access layers to eliminate SQL injection.
- [ ] `/learn` Never use raw un-sanitized HTML injection; require DOMPurify or framework auto-escaping against XSS.
- [ ] `/goal` Integrate static application security testing (SAST) tools in CI/CD pipelines to flag OWASP Top 10 risks.
- [ ] `/learn` Mandate multi-factor authentication (MFA) and threat modeling during initial feature architecture.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Mandatory Checks

All applications must actively defend against the OWASP Top 10 vulnerabilities. CI/CD pipelines must include static analysis tools (e.g., SonarQube, Semgrep) configured to flag these issues.

## 2. Key Mitigations

- **Injection:** Always use ORMs, parameterized queries, or prepared statements. Never concatenate strings to build SQL queries.

```typescript
// ❌ FORBIDDEN: Raw SQL concatenation vulnerable to injection
const query = `SELECT * FROM users WHERE email = '${userInput}'`;

// ✅ REQUIRED: Parameterized query
const query = "SELECT * FROM users WHERE email = $1";
await db.query(query, [userInput]);
```

- **Broken Authentication:** Implement multi-factor authentication (MFA) and strict password complexity rules.
- **Cross-Site Scripting (XSS):** Rely on modern framework auto-escaping (React, Vue). Never use dangerously set inner HTML without a strict sanitizer (e.g., DOMPurify).

```typescript
// ❌ FORBIDDEN: Unsanitized inner HTML
element.innerHTML = rawUntrustedHtml;

// ✅ REQUIRED: Sanitized HTML via DOMPurify
element.innerHTML = DOMPurify.sanitize(rawUntrustedHtml);
```

- **Insecure Design:** Threat modeling must be conducted during the planning phase of any major new feature.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-SEC-003: OWASP Top 10 Mitigation Controls & Input Validation

**Given** Web application handlers, API endpoints, and data access layers
**When** Audited against OWASP Top 10 defense guidelines and static security scans
**Then** Parameterized queries eliminate injection, user inputs are sanitized, and SAST tools report zero high/critical vulnerabilities

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/04-owasp-top-10.md --check-only
```
**Expected:** exit 0. Zero violations.
