# 06 — RSA Credential Vault and Authentication

- **Subsystem:** Security Architecture & Credential Encryption
- **Status:** Authoritative Reference

## 1. Masked Password Capture
- Raw terminal mode captures sensitive input with bullet masking (`*`) and zero echo.
- Interactive user consent protocol prompts for explicit confirmation before key deployment.

## 2. RSA-OAEP Salt Credential Vault
- Local vault at `~/.gitmap/credentials.vault` stores secrets encrypted with RSA-2048 and OAEP padding.
- Key derivation uses PBKDF2 with SHA-256 and 100,000 salt iterations.
- Memory zeroization wipes sensitive byte buffers immediately after cryptographic use.

## 3. Zero Plaintext Invariants
- Passwords and private keys are never logged to console, written to JSON, or exposed in traces.
