# Avoid 05: Hardcoded Paths and Secrets Exposure

- **Category:** Anti-Pattern & Strict Constraint
- **Status:** Mandatory Enforcement

## 1. Core Rule
NEVER hardcode absolute filesystem paths (such as root drive paths) or commit sensitive credentials (passwords, tokens, keys) in repository files.

## 2. Rationale
- Absolute paths break portability across developer machines, CI runners, and operating systems.
- Plaintext secrets represent severe security vulnerabilities.

## 3. Enforcement
- All paths must use relative Git paths or resolve dynamically via `os.UserHomeDir`.
- Credentials must be encrypted in `credentials.vault` or supplied via environment variables.
