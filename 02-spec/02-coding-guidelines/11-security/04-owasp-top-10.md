# OWASP Top 10 Mitigation

## 1. Mandatory Checks

All applications must actively defend against the OWASP Top 10 vulnerabilities. CI/CD pipelines must include static analysis tools (e.g., SonarQube, Semgrep) configured to flag these issues.

## 2. Key Mitigations

- **Injection:** Always use ORMs, parameterized queries, or prepared statements. Never concatenate strings to build SQL queries.
- **Broken Authentication:** Implement multi-factor authentication (MFA) and strict password complexity rules.
- **Cross-Site Scripting (XSS):** Rely on modern framework auto-escaping (React, Vue). Never use dangerously set inner HTML without a strict sanitizer (e.g., DOMPurify).
- **Insecure Design:** Threat modeling must be conducted during the planning phase of any major new feature.
