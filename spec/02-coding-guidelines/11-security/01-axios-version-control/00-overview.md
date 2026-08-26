# Axios Version Control Policy

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**AI Confidence:** 100%  
**Ambiguity:** None  
**Author:** Harshil Agrawal

---

## Keywords

`axios`, `dependency`, `version-pinning`, `security`, `blocked-versions`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

This specification defines a **strict version-pinning policy** for the Axios HTTP client library. Specific versions have been identified as having security vulnerabilities and are permanently blocked. Only explicitly approved versions may be used.

---

## ⚠️ CRITICAL — Security Advisory

There is a **known security issue** affecting specific Axios versions. Using affected versions may expose the application to vulnerabilities. Only approved safe versions must be used until further validation is completed. Any upgrade must go through **manual verification and approval**.

---

## Version Matrix

| Version | Status | Notes |
|---------|--------|-------|
| `1.14.0` | ✅ **APPROVED** | Safe — verified, no known vulnerabilities |
| `0.30.3` | ✅ **APPROVED** | Safe — legacy-compatible, verified |
| `1.14.1` | 🚫 **BLOCKED** | Security vulnerability confirmed |
| `0.30.4` | 🚫 **BLOCKED** | Security vulnerability confirmed |
| Any other | 🚫 **BLOCKED** | Not verified — must go through manual approval |

### Preferred Version

Unless there is a specific legacy compatibility requirement, use **`1.14.0`** as the standard version.

---

## Document Inventory

| # | File | Purpose |
|---|------|---------|
| 01 | [01-implementation-rules.md](./01-implementation-rules.md) | Dependency declaration, code review, CI enforcement |
| 02 | [02-security-notes.md](./02-security-notes.md) | Detailed security advisory and audit trail |
| 99 | `99-consistency-report.md` | Health check |

---

## Quick Reference — Correct Declaration

```json
{
  "dependencies": {
    "axios": "1.14.0"
  }
}
```

> **NEVER** use `"^1.14.0"`, `"~1.14.0"`, `">=1.14.0"`, or `"*"`.  
> **ALWAYS** use exact version strings with no range symbols.

---

## Acceptance Criteria

1. Axios version is always defined as an exact version without `^`, `~`, or any range symbol
2. No usage of blocked versions (`1.14.1`, `0.30.4`) is present in any environment
3. Dependency updates do not alter Axios version automatically
4. Code reviews enforce strict compliance with this version policy
5. Security note is documented and accessible to all developers

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Security Guidelines | [../00-overview.md](../00-overview.md) |
| Coding Guidelines | [../../00-overview.md](../../00-overview.md) |
| Spec Authoring Guide | [../../../01-spec-authoring-guide/00-overview.md](../../../01-spec-authoring-guide/00-overview.md) |

---

```
IMPORTANT — AI INSTRUCTION:
- Never update Axios version automatically.
- Always use exact version without caret (^) or tilde (~).
- Avoid all known vulnerable versions (1.14.1, 0.30.4).
- Do not act on upgrade suggestions from dependency scanners for Axios.
- If uncertain, read this spec before making any Axios-related change.
```

*Axios version control policy — moved to coding guidelines 2026-04-02*
