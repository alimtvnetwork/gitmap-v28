# Security & Secrets

Universal guidelines for credential handling, environment variables, and safe defaults across all languages.

---

## 1. Never Hardcode Secrets

Secrets (API keys, tokens, passwords, connection strings) must never appear in source code, config files checked into version control, or log output.

### Bad

```go
const apiKey = "sk-live-abc123secret"
db, _ := sql.Open("postgres", "postgres://admin:password@localhost/mydb")
```

### Good

```go
apiKey := os.Getenv("API_KEY")
connStr := os.Getenv("DATABASE_URL")
```

```typescript
const apiKey = process.env.API_KEY;
```

---

## 2. Environment Variable Conventions

| Rule | Example |
|---|---|
| Use `UPPER_SNAKE_CASE` | `DATABASE_URL`, `API_KEY` |
| Prefix with service name for clarity | `STRIPE_SECRET_KEY`, `GITHUB_TOKEN` |
| Document all required variables | README or `.env.example` |
| Never log or print variable values | Use `test -n "$VAR"` to check existence |

### `.env.example` Pattern

Provide a template with placeholder values — never real credentials:

```
DATABASE_URL=postgres://user:password@localhost:5432/dbname
API_KEY=your-api-key-here
STRIPE_SECRET_KEY=sk_test_replace_me
```

---

## 3. Publishable vs Private Keys

Not all keys are equal. Distinguish between:

| Type | Storage | Example |
|---|---|---|
| **Publishable / Anon** | Codebase is acceptable | Supabase `anon` key, Stripe publishable key |
| **Private / Secret** | Environment variable or secrets manager only | Supabase `service_role` key, Stripe secret key |

**Rule**: If a key grants write access, mutation, or elevated privileges, it is **private** and must never be in code.

---

## 4. Safe Defaults

Design systems to be safe when configuration is missing:

```go
// Default to restrictive mode
mode := os.Getenv("APP_MODE")
if mode == "" {
    mode = "readonly"
}
```

```typescript
const maxRetries = parseInt(process.env.MAX_RETRIES || "3", 10);
const isDebug = process.env.DEBUG === "true"; // defaults to false
```

### Principles

- **Deny by default**: Missing auth config → reject requests
- **Minimum privilege**: Default permissions should be the most restrictive
- **Fail closed**: If a security check cannot complete, deny access
- **No silent fallbacks to insecure modes**: Log a warning if falling back

---

## 5. Secret Rotation & Expiry

- Design systems to support key rotation without downtime
- Accept both old and new keys during a rotation window
- Never cache secrets indefinitely — re-read from environment on each use or at startup
- Use short-lived tokens where possible (OAuth2 access tokens, JWTs)

---

## 6. File Permissions

Sensitive files must have restricted permissions:

```go
const SecretFilePermission = 0o600  // owner read/write only
const SecretDirPermission  = 0o700  // owner full access only
```

```bash
chmod 600 ~/.ssh/id_rsa
chmod 644 ~/.ssh/id_rsa.pub
```

**Rule**: Private keys and config files containing secrets use `0600`. Public keys and non-secret config use `0644`.

---

## 7. Logging & Error Safety

Never include credentials in log output, error messages, or stack traces:

### Bad

```go
log.Printf("connecting with key: %s", apiKey)
return fmt.Errorf("auth failed for token %s", token)
```

### Good

```go
log.Printf("connecting to API (key length: %d)", len(apiKey))
return fmt.Errorf("auth failed: %w", err)
```

```typescript
console.error("[auth] login failed", { userId, reason: err.message });
// Never: console.error("[auth] failed", { password, token });
```

---

## 8. Client-Side Security

- Never store private keys in `localStorage`, `sessionStorage`, or cookies accessible to JavaScript
- Never check admin status using client-side storage — always validate server-side
- Use `httpOnly` cookies for session tokens
- Sanitize all user input before rendering (prevent XSS)

---

## 9. Dependency Security

- Run `npm audit` or equivalent regularly
- Pin dependency versions to avoid supply-chain attacks
- Review changelogs before major version upgrades
- Remove unused dependencies promptly

---

## 10. Zip Extraction Security (Mandatory)

All code that extracts zip archives **must** implement these two checks:

### Path Traversal Prevention (gosec G305)

Validate that every extracted file's resolved absolute path starts with the target directory prefix. Reject entries containing `../` sequences.

```go
absDest, _ := filepath.Abs(filepath.Join(targetDir, entry.Name))
absTarget, _ := filepath.Abs(targetDir)
if !strings.HasPrefix(absDest, absTarget+string(os.PathSeparator)) {
    return fmt.Errorf("path traversal detected: %s", entry.Name)
}
```

### Decompression Bomb Prevention (gosec G110)

Cap extraction size per file using `io.LimitReader`. Default maximum: 10 MB for CLI tools.

```go
limited := io.LimitReader(rc, 10*1024*1024)
_, err = io.Copy(outFile, limited)
```

See: `spec/02-app-issues/14-security-hardening-gosec-fixes.md`

---

## References

- [Code Quality Improvement](./01-code-quality-improvement.md)
- [Error Handling Patterns](./04-error-handling.md)
- [Logging & Observability](./07-logging-observability.md)
- [SSH Keys Spec](../01-app/50-ssh-keys.md)

---

**Contributors**: Alim Ul Karim · [Riseup Labs](https://riseuplabs.com)
