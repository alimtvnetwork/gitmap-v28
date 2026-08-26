# Secret Management

## 1. Zero Hardcoded Secrets

- Code MUST NEVER contain hardcoded secrets, API keys, passwords, or tokens.
- `.env` files containing secrets MUST NOT be committed to version control. Ensure `.env` is in `.gitignore`.

## 2. Secret Vaults

- Production environments must retrieve secrets dynamically at runtime or during the deployment phase from a secure vault (e.g., HashiCorp Vault, AWS Secrets Manager, Azure Key Vault).

## 3. Secret Rotation

- Services must be designed to handle dynamic secret rotation without requiring a full application restart if possible, or gracefully restart when the orchestrator updates the injected secrets.
