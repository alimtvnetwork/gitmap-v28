# Encryption Standards (AI Execution Prompt)

> **/goal** Enforce robust data protection with AES-256 encryption at rest, TLS 1.2+ transit encryption, and Argon2id cryptographic hashing.
> **/learn** Eliminate weak hashing algorithms (MD5, SHA-1, plain bcrypt), enforce column-level encryption for PII/PHI, and mandate end-to-end TLS for all internal and external communication.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Require AES-256 encryption at rest across all databases storing PII, PHI, or financial data.
- [ ] `/learn` Implement column-level encryption for highly sensitive attributes prior to database persistence.
- [ ] `/goal` Enforce TLS 1.2 or higher for all external endpoints and internal service-to-service communication.
- [ ] `/learn` Use Argon2id exclusively for all password and secret hashing to ensure maximum GPU cracking resistance.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Encryption at Rest

- All databases storing PII (Personally Identifiable Information), PHI (Protected Health Information), or financial data must have encryption at rest enabled using AES-256 (e.g., AWS KMS, Azure Key Vault).
- Column-level encryption is required for highly sensitive fields (e.g., SSN, credit card tokens) before they are written to the database.

## 2. Encryption in Transit

- All external and internal service-to-service communication MUST occur over TLS 1.2 or higher (HTTPS/gRPC over TLS). No plaintext HTTP allowed in production.

## 3. Hashing Algorithms

- Never use MD5, SHA-1, or bcrypt for new password storage.
- **Mandatory Algorithm:** Use Argon2id for all password and secret hashing. It provides superior resistance to GPU cracking attacks.

```go
// ❌ FORBIDDEN: Legacy weak hashing algorithms (MD5, SHA-1)
hasher := md5.New()
hasher.Write([]byte(password))

// ✅ REQUIRED: Argon2id cryptographic hashing
hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-SEC-002: Encryption at Rest & Key Derivation Standards

**Given** Data persistence, transmission, and cryptographic hashing implementations
**When** Audited against encryption and hashing requirements
**Then** AES-256 is used at rest, TLS 1.2+ is enforced in transit, and Argon2id is used for password hashing

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/03-encryption-standards.md --check-only
```
**Expected:** exit 0. Zero violations.
