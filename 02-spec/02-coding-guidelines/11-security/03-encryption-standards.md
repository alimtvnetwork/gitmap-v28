# Encryption Standards

## 1. Encryption at Rest

- All databases storing PII (Personally Identifiable Information), PHI (Protected Health Information), or financial data must have encryption at rest enabled using AES-256 (e.g., AWS KMS, Azure Key Vault).
- Column-level encryption is required for highly sensitive fields (e.g., SSN, credit card tokens) before they are written to the database.

## 2. Encryption in Transit

- All external and internal service-to-service communication MUST occur over TLS 1.2 or higher (HTTPS/gRPC over TLS). No plaintext HTTP allowed in production.

## 3. Hashing Algorithms

- Never use MD5, SHA-1, or bcrypt for new password storage.
- **Mandatory Algorithm:** Use Argon2id for all password and secret hashing. It provides superior resistance to GPU cracking attacks.
