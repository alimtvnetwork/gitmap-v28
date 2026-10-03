# Security Guidelines (AI Execution Prompt)

> **/goal** Master and enforce repository-wide security guidelines, zero-trust token lifecycles, cryptographic standards, secret isolation, and dependency vulnerability defenses.
> **/learn** Internalize strict secret segregation (`repo-secrets`), mandatory Argon2id password hashing, AES-256 encryption at rest, short-lived JWT lifetimes with HttpOnly storage, and explicit dependency version pinning.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce zero secrets or credentials committed to standard source trees, routing them to `repo-secrets` via GitMap.
- [ ] `/learn` Never store sensitive tokens in client-accessible storage; mandate `HttpOnly` and `SameSite=Strict` cookies.
- [ ] `/goal` Require modern cryptographic standards (Argon2id hashing, AES-256 at rest, TLS 1.2+ in transit).
- [ ] `/learn` Verify security specifications and acceptance criteria pass validation using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-04-16
**AI Confidence:** High
**Ambiguity:** None

---

## Keywords

`security` · `dependency-pinning` · `vulnerability` · `version-control` · `axios` · `supply-chain` · `cve` · `audit`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `readme.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Central location for all **security-related coding guidelines**, policies, and advisory documentation. This module covers dependency security, version pinning policies, vulnerability tracking, and secure coding practices.

Any security discussion, advisory, or policy that affects how code is written or dependencies are managed belongs here.

---

## Categories

| # | Subfolder | Description | Files |
|---|-----------|-------------|-------|
| 01 | [01-axios-version-control/](./01-axios-version-control/readme.md) | Axios HTTP client version pinning policy and security advisory | 4 |

---

## When to Add Content Here

Add a new subfolder under `11-security/` when:

- A **dependency security vulnerability** is discovered and requires a pinning policy
- A **secure coding pattern** needs to be documented (e.g., input sanitization, auth token handling)
- A **supply chain security** concern arises (e.g., compromised packages)
- A **security audit** produces findings that should be codified as rules

### Subfolder Template

```
11-security/
└── NN-{topic-name}/
    ├── readme.md              ← Policy summary, version matrix
    ├── 01-implementation-rules.md  ← How to enforce the policy
    ├── 02-security-notes.md        ← Detailed advisory, audit trail
    └── 99-consistency-report.md    ← Health check
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Parent Overview | [../readme.md](../readme.md) |
| Cross-Language Guidelines | [../01-cross-language/readme.md](../01-cross-language/readme.md) |
| File & Folder Naming | [../08-file-folder-naming/readme.md](../08-file-folder-naming/readme.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-SEC-000: Security Guidelines Directory Index & Policy Conformance

**Given** The security specifications directory under `02-spec/02-coding-guidelines/11-security/`
**When** Running guideline and structure validation across all security policy files
**Then** All security guidelines comply with repository standards and pass automated verification

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/readme.md --check-only
```
**Expected:** exit 0. Zero violations.

---

*Security guidelines — single source of truth for all security-related coding policies.*
