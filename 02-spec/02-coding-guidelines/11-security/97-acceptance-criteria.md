# Security Guidelines — Acceptance Criteria (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all security coding guideline specifications in `11-security/`.
> **/learn** Enforce the canonical security criteria taxonomy (`AC-CG-SEC-[NUM]` and `AC-CG-SEC-AXIOS-[NUM]`), zero hardcoded secrets, Argon2id hashing, short-lived JWT lifetimes, exact Axios dependency pinning, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `11-security/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths across all guideline links.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Security Criteria Inventory (`AC-CG-SEC-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-SEC-000` | Security Guidelines Directory Index & Policy Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/readme.md --check-only` |
| `AC-CG-SEC-001` | JWT Token Lifecycle & HttpOnly Cookie Storage | [`02-jwt-standards.md`](02-jwt-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/02-jwt-standards.md --check-only` |
| `AC-CG-SEC-002` | Encryption at Rest & Key Derivation Standards | [`03-encryption-standards.md`](03-encryption-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/03-encryption-standards.md --check-only` |
| `AC-CG-SEC-003` | OWASP Top 10 Mitigation Controls & Input Validation | [`04-owasp-top-10.md`](04-owasp-top-10.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/04-owasp-top-10.md --check-only` |
| `AC-CG-SEC-004` | Zero Secrets in Source Control & Environment Vaulting | [`05-secret-management.md`](05-secret-management.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/05-secret-management.md --check-only` |
| `AC-CG-SEC-AXIOS-000` | Axios Client Security Overview | [`01-axios-version-control/readme.md`](01-axios-version-control/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/readme.md --check-only` |
| `AC-CG-SEC-AXIOS-001` | Strict Axios Version Pinning & Dependency Locking | [`01-axios-version-control/02-implementation-rules.md`](01-axios-version-control/02-implementation-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/02-implementation-rules.md --check-only` |
| `AC-CG-SEC-AXIOS-002` | CVE Remediation & Supply Chain Security Verification | [`01-axios-version-control/03-security-notes.md`](01-axios-version-control/03-security-notes.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/03-security-notes.md --check-only` |
| `AC-CG-SEC-REG-001` | Security Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-SEC-000: Security Guidelines Directory Index & Policy Conformance

- [ ] Directory entry point provides scope, keywords, scoring table, categories, and cross-references.
- [ ] Active AI prompt header, actionable checklist, and acceptance criteria sections are present.
- [ ] Internal relative markdown links resolve accurately with zero broken links.

**Given** The security specifications directory under `02-spec/02-coding-guidelines/11-security/`
**When** Running guideline and structure validation across all security policy files
**Then** All security guidelines comply with repository standards and pass automated verification

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/readme.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-001: JWT Token Lifecycle & HttpOnly Cookie Storage

- [ ] Access tokens enforce ephemeral lifecycles (maximum 15 minutes) with in-memory storage.
- [ ] Refresh tokens enforce mandatory rotation upon usage to prevent token replay attacks.
- [ ] Refresh tokens are stored exclusively in `HttpOnly`, `Secure`, `SameSite=Strict` cookies and are inaccessible to JavaScript.
- [ ] JWT payloads strictly exclude PII, passwords, and private credentials.

**Given** Authentication services handling JWT generation and validation
**When** Evaluating token lifetimes, storage mechanisms, and payload claims
**Then** Access tokens expire within 15 minutes, refresh tokens are stored in HttpOnly cookies, and payloads contain no PII

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/02-jwt-standards.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-002: Encryption at Rest & Key Derivation Standards

- [ ] AES-256 encryption at rest is mandatory for databases containing PII, PHI, or financial data.
- [ ] Sensitive database attributes enforce application-level column encryption prior to storage.
- [ ] External and internal service communication strictly mandates TLS 1.2 or higher.
- [ ] Password and secret hashing enforces Argon2id and completely bans MD5, SHA-1, and plain bcrypt.

**Given** Data persistence, transmission, and cryptographic hashing implementations
**When** Audited against encryption and hashing requirements
**Then** AES-256 is used at rest, TLS 1.2+ is enforced in transit, and Argon2id is used for password hashing

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/03-encryption-standards.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-003: OWASP Top 10 Mitigation Controls & Input Validation

- [ ] Data queries strictly utilize ORMs, parameterized statements, or prepared queries with zero string concatenation.
- [ ] Authentication layers enforce multi-factor authentication (MFA) and strong credential requirements.
- [ ] Frontend user input rendering enforces DOMPurify or framework auto-escaping against XSS attacks.
- [ ] CI/CD pipelines integrate automated SAST scanning tools to flag OWASP Top 10 issues.

**Given** Web application handlers, API endpoints, and data access layers
**When** Audited against OWASP Top 10 defense guidelines and static security scans
**Then** Parameterized queries eliminate injection, user inputs are sanitized, and SAST tools report zero high/critical vulnerabilities

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/04-owasp-top-10.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-004: Zero Secrets in Source Control & Environment Vaulting

- [ ] Source trees contain zero hardcoded secrets, API tokens, passwords, or committed `.env` files.
- [ ] Runtime environments dynamically inject secrets via secure vault providers (HashiCorp Vault, AWS KMS/Secrets Manager).
- [ ] Secrets and credentials offload cleanly to `repo-secrets` using GitMap CLI commands (`gitmap rs`).
- [ ] Temporary test scripts and harnesses offload cleanly to `repo-cache` using GitMap CLI commands (`gitmap rc`).

**Given** Application repositories, configuration files, and script automation
**When** Audited for secret segregation and special repository routing
**Then** Zero credentials exist in source trees and secrets are routed to `repo-secrets` via GitMap

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/05-secret-management.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-AXIOS-000: Axios Client Security Overview

- [ ] Axios dependency is pinned to an exact version string without range characters (`^`, `~`, `>=`, `*`).
- [ ] Known vulnerable versions (`1.14.1`, `0.30.4`) are permanently blocked across all deployment environments.
- [ ] Preferred approved version is `1.14.0` (with `0.30.3` for legacy systems).
- [ ] Automated dependency update tools are prevented from modifying the Axios version.

**Given** Project package dependencies and Axios version policy definitions
**When** Audited for version pinning and vulnerability mitigation
**Then** Axios is pinned to approved releases, range specifiers are eliminated, and automated upgrade tools ignore Axios

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/readme.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-AXIOS-001: Strict Axios Version Pinning & Dependency Locking

- [ ] `package.json` and lockfiles reflect identical exact Axios version declarations.
- [ ] Dependabot and Renovate configuration files explicitly ignore or exclude Axios.
- [ ] CI/CD pipeline check script validates Axios version and rejects range symbols and blocked versions.
- [ ] Version upgrades strictly follow the 6-step manual review and verification procedure.

**Given** Node.js package manifests, lock files, and CI dependency check scripts
**When** Validating Axios package entries and automated upgrade tool exclusions
**Then** Axios is declared without range specifiers, lock file reflects exact version, and CI rejects non-pinned entries

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/02-implementation-rules.md --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-SEC-AXIOS-002: CVE Remediation & Supply Chain Security Verification

- [ ] Security advisory impact assessment documents all affected tiers (Production, Staging, Dev, CI/CD).
- [ ] Blocked versions `1.14.1` and `0.30.4` are tracked with an explicit audit trail.
- [ ] Security monitoring schedule defines monthly vulnerability audits and pre-release lockfile checks.
- [ ] Rapid response procedure requires safe version remediation within 24 hours of CVE notification.

**Given** Security advisories, vulnerability tracking registries, and review audit logs
**When** Assessing known Axios CVEs and scheduled advisory monitoring channels
**Then** Blocked versions are quarantined across all tiers and response procedures execute within 24 hours

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/03-security-notes.md --check-only
```
**Expected:** exit 0. Zero violations.

---

## Verification & Acceptance Criteria

### AC-CG-SEC-REG-001: Security Acceptance Criteria Registry Conformance

**Given** The security acceptance criteria registry in `02-spec/02-coding-guidelines/11-security/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All security specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only
```
**Expected:** exit 0. Zero violations.
