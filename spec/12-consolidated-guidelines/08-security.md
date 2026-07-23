# 08 — Security & Secrets

Credential handling, environment variables, and safe defaults.

## Never Hardcode Secrets

API keys, tokens, passwords, connection strings must never appear in source code or version-controlled config files.

## Environment Variable Conventions

Use `UPPER_SNAKE_CASE`. Prefix with service name. Document all required vars. Never log values.

## Publishable vs Private Keys

Publishable/anon keys (Supabase anon, Stripe publishable) are OK in code. Private/secret keys (service_role, Stripe secret) must be environment variables only.

## Safe Defaults

Missing config defaults to restrictive mode. Deny by default. Minimum privilege. Fail closed.

## File Permissions

Private keys: `0o600` (owner read/write). Public keys: `0o644`.

## Logging Safety

Never include credentials in logs, error messages, or stack traces. Reference secrets by source, not value.

## Zip Extraction Security

1. **Path traversal prevention (G305)**: Validate extracted paths stay within target directory.
2. **Decompression bomb prevention (G110)**: Cap extraction with `io.LimitReader` (10 MB default).

## Client-Side Security

Never store private keys in localStorage/sessionStorage. Never check admin status client-side. Use httpOnly cookies.

---

Source: `spec/05-coding-guidelines/08-security-secrets.md`
