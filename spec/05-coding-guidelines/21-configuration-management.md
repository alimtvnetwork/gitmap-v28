# Configuration Management

Universal guidelines for environment variables, feature flags, config validation, and secrets rotation across all languages.

---

## 1. Environment Variables

### Naming

| Rule | Example |
|---|---|
| `UPPER_SNAKE_CASE` | `DATABASE_URL`, `API_KEY` |
| Prefix with service name | `STRIPE_SECRET_KEY`, `GITHUB_TOKEN` |
| Boolean vars use `true`/`false` strings | `ENABLE_CACHE=true` |
| Document all required vars | `.env.example` with placeholders |

### Loading Order

```
Process env → .env file → hardcoded defaults
```

- Process-level environment variables always win.
- `.env` files are for local development only — never commit real values.
- Defaults must be safe (restrictive mode, low limits).

### Go

```go
func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}

mode := getEnv("APP_MODE", "readonly")
```

### TypeScript

```typescript
function getEnv(key: string, fallback: string): string {
  return process.env[key] ?? fallback;
}

const mode = getEnv("APP_MODE", "readonly");
```

### Rules

| Rule | Detail |
|---|---|
| Never log variable values | Use `test -n "$VAR"` or check length |
| Never hardcode secrets | Always read from environment |
| Validate at startup | Fail fast if required vars are missing |
| Type-convert explicitly | `parseInt(process.env.PORT \|\| "3000", 10)` |

---

## 2. Feature Flags

### Purpose

Feature flags decouple deployment from release. Code ships to production but new behavior activates only when the flag is enabled.

### Flag Types

| Type | Lifetime | Example |
|---|---|---|
| Release flag | Short (days–weeks) | `ENABLE_NEW_DASHBOARD` |
| Ops flag | Medium (weeks–months) | `ENABLE_VERBOSE_LOGGING` |
| Experiment flag | Short (A/B test duration) | `EXPERIMENT_CHECKOUT_V2` |
| Permission flag | Long (permanent) | `ENABLE_ADMIN_PANEL` |

### Implementation

```go
// flags.go
func IsEnabled(flag string) bool {
    return os.Getenv(flag) == "true"
}

if IsEnabled("ENABLE_NEW_SCANNER") {
    runNewScanner()
} else {
    runLegacyScanner()
}
```

```typescript
// flags.ts
export function isEnabled(flag: string): boolean {
  return process.env[flag] === "true";
}

if (isEnabled("ENABLE_NEW_DASHBOARD")) {
  renderNewDashboard();
} else {
  renderLegacyDashboard();
}
```

### Rules

| Rule | Detail |
|---|---|
| Default to `false` (off) | New flags must not change existing behavior |
| Clean up after rollout | Remove flag checks and dead code paths |
| No nested flags | One flag per decision point |
| Test both paths | Unit tests must cover flag on and flag off |
| Name with `ENABLE_` prefix | Makes intent clear |

### Flag Lifecycle

```
Define → Implement → Test both paths → Enable in staging → Enable in production → Remove flag + dead code
```

---

## 3. Config Validation

### Validate at Startup

Never allow invalid configuration to propagate. Validate all config values immediately after loading:

```go
func ValidateConfig(cfg Config) error {
    if cfg.Port < 1 || cfg.Port > 65535 {
        return fmt.Errorf("invalid port: %d (must be 1-65535)", cfg.Port)
    }
    if cfg.Mode != "https" && cfg.Mode != "ssh" {
        return fmt.Errorf("invalid mode: %q (must be 'https' or 'ssh')", cfg.Mode)
    }
    if cfg.MaxRetries < 0 {
        return fmt.Errorf("invalid max_retries: %d (must be >= 0)", cfg.MaxRetries)
    }
    return nil
}
```

```typescript
interface ConfigSchema {
  port: number;
  mode: "https" | "ssh";
  maxRetries: number;
}

function validateConfig(cfg: ConfigSchema): void {
  if (cfg.port < 1 || cfg.port > 65535) {
    throw new Error(`Invalid port: ${cfg.port}`);
  }
  if (!["https", "ssh"].includes(cfg.mode)) {
    throw new Error(`Invalid mode: ${cfg.mode}`);
  }
  if (cfg.maxRetries < 0) {
    throw new Error(`Invalid maxRetries: ${cfg.maxRetries}`);
  }
}
```

### Validation Rules

| Rule | Detail |
|---|---|
| Fail fast | Crash on invalid config — don't silently use bad values |
| Validate types | Ensure numbers are numbers, booleans are booleans |
| Validate ranges | Ports, timeouts, retry counts have valid bounds |
| Validate enums | Mode, format, and level fields match allowed values |
| Validate paths | Check existence for required files/directories |
| Log validated config | Print sanitized config (no secrets) at startup for debugging |

### Schema Validation (TypeScript)

Use `zod` or equivalent for structured validation:

```typescript
import { z } from "zod";

const ConfigSchema = z.object({
  port: z.number().min(1).max(65535),
  mode: z.enum(["https", "ssh"]),
  maxRetries: z.number().min(0).default(3),
  outputDir: z.string().min(1),
});

type Config = z.infer<typeof ConfigSchema>;

const config = ConfigSchema.parse(rawConfig); // throws on invalid
```

---

## 4. Secrets Rotation

### Principles

- Design systems to support key rotation without downtime.
- Accept both old and new keys during a rotation window.
- Never cache secrets indefinitely — re-read from environment on each use or at startup.
- Use short-lived tokens where possible (OAuth2 access tokens, JWTs).

### Rotation Pattern

```
Generate new key → Deploy new key alongside old → Verify new key works → Remove old key → Update documentation
```

### Implementation

```go
// Support dual keys during rotation
func authenticate(token string) bool {
    primaryKey := os.Getenv("API_KEY")
    rotationKey := os.Getenv("API_KEY_PREVIOUS")

    if token == primaryKey {
        return true
    }
    if rotationKey != "" && token == rotationKey {
        log.Println("WARN: request authenticated with rotation key")
        return true
    }
    return false
}
```

```typescript
function authenticate(token: string): boolean {
  const primaryKey = process.env.API_KEY;
  const rotationKey = process.env.API_KEY_PREVIOUS;

  if (token === primaryKey) return true;
  if (rotationKey && token === rotationKey) {
    console.warn("authenticated with rotation key");
    return true;
  }
  return false;
}
```

### Rules

| Rule | Detail |
|---|---|
| Always support dual keys | `API_KEY` + `API_KEY_PREVIOUS` during rotation |
| Log rotation key usage | Track when old keys are still in use |
| Set rotation reminders | Calendar/CI alerts for key expiry |
| Automate where possible | Use secrets managers with auto-rotation |
| Never store in code | Environment variables or secrets manager only |

---

## 5. Config File Standards

### Format

- Use JSON for machine-readable config (`.json`).
- Use camelCase for field names.
- Arrays default to `[]`, never `null`.
- Strings default to `""`, never `null`.
- No deeply nested objects unless absolutely necessary.

### Three-Layer Merge

```
Hardcoded defaults → Config file → CLI flags / env vars (highest priority)
```

See [Configuration Pattern](../03-general/05-config-pattern.md) for full merge logic.

### Documentation

Every config file must have a corresponding documentation section listing:

| Required | Detail |
|---|---|
| Field name | Exact JSON key |
| Type | `string`, `number`, `bool`, `[]string`, etc. |
| Default | Value used when omitted |
| Description | One-line explanation |
| Valid values | Enum options or range constraints |

---

## Constraints

| Constraint | Detail |
|---|---|
| No secrets in config files | Use environment variables exclusively |
| No `null` values | Use empty strings or empty arrays |
| Validate before use | Never trust raw config input |
| Feature flags default to off | Opt-in, never opt-out |
| Clean up stale flags | Remove within one release cycle after full rollout |
| Config files are version-controlled | But `.env` files with real values are not |

---

## References

- [Security & Secrets](./08-security-secrets.md)
- [Configuration Pattern](../03-general/05-config-pattern.md)
- [Config Spec](../01-app/06-config.md)
- [Error Handling Patterns](./04-error-handling.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
