# Feature Flags

Standards for flag lifecycle management, rollout strategies, cleanup policies, and testing with feature flags across all languages.

---

## 1. Flag Lifecycle

### Stages

```
Define → Implement → Test both paths → Stage rollout → Production rollout → Remove flag + dead code
```

### Flag Types

| Type | Lifetime | Example |
|---|---|---|
| Release flag | Short (days–weeks) | `ENABLE_NEW_DASHBOARD` |
| Ops flag | Medium (weeks–months) | `ENABLE_VERBOSE_LOGGING` |
| Experiment flag | Short (A/B test duration) | `EXPERIMENT_CHECKOUT_V2` |
| Permission flag | Long (permanent) | `ENABLE_ADMIN_PANEL` |

### Flag Registry

Maintain a central registry documenting every flag:

| Field | Detail |
|---|---|
| Name | `ENABLE_NEW_SCANNER` |
| Type | Release |
| Owner | Team or individual responsible |
| Created | Date flag was introduced |
| Expected removal | Target date for cleanup |
| Description | One-line purpose statement |
| Status | `active`, `rolling-out`, `deprecated`, `removed` |

### Go

```go
// flags/registry.go
type Flag struct {
    Name        string
    FlagType    string // "release", "ops", "experiment", "permission"
    Owner       string
    CreatedAt   string
    RemoveBy    string
    Description string
}

var Registry = []Flag{
    {
        Name:        "ENABLE_NEW_SCANNER",
        FlagType:    "release",
        Owner:       "scanner-team",
        CreatedAt:   "2026-01-15",
        RemoveBy:    "2026-02-15",
        Description: "Activates the v2 scanner pipeline",
    },
}
```

### TypeScript

```typescript
// flags/registry.ts
interface FeatureFlag {
  name: string;
  type: "release" | "ops" | "experiment" | "permission";
  owner: string;
  createdAt: string;
  removeBy: string;
  description: string;
}

export const flagRegistry: FeatureFlag[] = [
  {
    name: "ENABLE_NEW_SCANNER",
    type: "release",
    owner: "scanner-team",
    createdAt: "2026-01-15",
    removeBy: "2026-02-15",
    description: "Activates the v2 scanner pipeline",
  },
];
```

---

## 2. Rollout Strategies

### Staged Rollout

Never enable a flag globally on first deploy. Follow a controlled progression:

```
Off → Internal/dev → Canary (1–5%) → Partial (10–50%) → Full (100%)
```

### Percentage-Based Rollout

```go
func IsEnabledForUser(flag string, userID string, percentage int) bool {
    if percentage >= 100 {
        return true
    }
    if percentage <= 0 {
        return false
    }
    hash := fnv.New32a()
    hash.Write([]byte(flag + ":" + userID))
    bucket := int(hash.Sum32() % 100)
    return bucket < percentage
}
```

```typescript
function isEnabledForUser(flag: string, userId: string, percentage: number): boolean {
  if (percentage >= 100) return true;
  if (percentage <= 0) return false;

  let hash = 0;
  const key = `${flag}:${userId}`;
  for (let i = 0; i < key.length; i++) {
    hash = ((hash << 5) - hash + key.charCodeAt(i)) | 0;
  }
  const bucket = Math.abs(hash) % 100;
  return bucket < percentage;
}
```

### Environment-Based Rollout

```go
func IsEnabled(flag string) bool {
    return os.Getenv(flag) == "true"
}

// Per-environment config
// dev:     ENABLE_NEW_SCANNER=true
// staging: ENABLE_NEW_SCANNER=true
// prod:    ENABLE_NEW_SCANNER=false  (until verified in staging)
```

### Rules

| Rule | Detail |
|---|---|
| Start at 0% | New flags ship disabled |
| Monitor after each increase | Check error rates and latency before expanding |
| Rollback immediately on anomaly | Disable the flag, investigate later |
| Document rollout state | Track current percentage in the registry |
| Use consistent hashing | Same user always gets the same flag state |

---

## 3. Cleanup Policies

### When to Clean Up

A flag must be removed when:

- It has been at 100% in production for one full release cycle.
- The experiment has concluded and a decision has been made.
- The flag has passed its `removeBy` date.

### Cleanup Process

```
1. Verify flag is at 100% in all environments
2. Remove flag checks from code
3. Remove dead code path (the "off" branch)
4. Remove flag from registry
5. Remove environment variable from all environments
6. Update documentation
```

### Automated Staleness Detection

```go
func FindStaleFlags(registry []Flag) []Flag {
    var stale []Flag
    now := time.Now()
    for _, f := range registry {
        removeBy, err := time.Parse("2006-01-02", f.RemoveBy)
        if err != nil {
            continue
        }
        if now.After(removeBy) {
            stale = append(stale, f)
        }
    }
    return stale
}
```

```typescript
function findStaleFlags(registry: FeatureFlag[]): FeatureFlag[] {
  const now = new Date();
  return registry.filter((flag) => {
    const removeBy = new Date(flag.removeBy);
    return now > removeBy;
  });
}
```

### Rules

| Rule | Detail |
|---|---|
| Max lifetime for release flags | 30 days after full rollout |
| Max lifetime for experiment flags | Duration of experiment + 7 days |
| CI warning on stale flags | Fail or warn if a flag is past its `removeBy` date |
| No nested flags | One flag per decision point; never `if flagA && flagB` |
| Track flag count | Alert if total active flags exceeds a threshold (e.g., 20) |

---

## 4. Testing with Flags

### Test Both Paths

Every feature flag must have tests covering both the enabled and disabled states:

```go
func TestNewScanner_Enabled(t *testing.T) {
    t.Setenv("ENABLE_NEW_SCANNER", "true")
    result := RunScanner()
    if result.Version != "v2" {
        t.Errorf("expected v2 scanner, got %s", result.Version)
    }
}

func TestNewScanner_Disabled(t *testing.T) {
    t.Setenv("ENABLE_NEW_SCANNER", "false")
    result := RunScanner()
    if result.Version != "v1" {
        t.Errorf("expected v1 scanner, got %s", result.Version)
    }
}
```

```typescript
describe("scanner with feature flag", () => {
  it("uses v2 scanner when enabled", () => {
    process.env.ENABLE_NEW_SCANNER = "true";
    const result = runScanner();
    expect(result.version).toBe("v2");
    delete process.env.ENABLE_NEW_SCANNER;
  });

  it("uses v1 scanner when disabled", () => {
    process.env.ENABLE_NEW_SCANNER = "false";
    const result = runScanner();
    expect(result.version).toBe("v1");
    delete process.env.ENABLE_NEW_SCANNER;
  });
});
```

### Flag Isolation in Tests

```go
// Use t.Setenv for automatic cleanup (Go 1.17+)
func TestFeature(t *testing.T) {
    t.Setenv("ENABLE_FEATURE_X", "true")
    // env var automatically restored after test
}
```

```typescript
// Helper for flag-scoped tests
function withFlag(flag: string, value: string, fn: () => void): void {
  const original = process.env[flag];
  process.env[flag] = value;
  try {
    fn();
  } finally {
    if (original === undefined) {
      delete process.env[flag];
    } else {
      process.env[flag] = original;
    }
  }
}

withFlag("ENABLE_FEATURE_X", "true", () => {
  // test code here
});
```

### Rules

| Rule | Detail |
|---|---|
| Test both states | Every flag has at least two tests: on and off |
| Isolate flag state | Clean up env vars after each test |
| No flag dependencies in tests | Each test sets its own flag state explicitly |
| Integration tests use production defaults | Run with flags in their default (off) state |
| Pre-cleanup verification | Run full test suite with flag removed before merging cleanup PR |

---

## Constraints

| Constraint | Detail |
|---|---|
| Flags default to `false` (off) | New flags must not change existing behavior |
| `ENABLE_` prefix required | Makes intent and searchability clear |
| No nested flag checks | One flag per decision point |
| Registry is mandatory | Every flag must be documented with owner and removal date |
| Stale flags block CI | Past-due flags generate warnings or failures |
| Clean up within one release cycle | After full rollout, remove flag and dead code promptly |
| No business logic in flag checks | Flags gate entry points, not internal logic |

---

## References

- [Configuration Management](./21-configuration-management.md)
- [Testing Patterns](./06-testing-patterns.md)
- [CI/CD Patterns](./17-cicd-patterns.md)
- [Security & Secrets](./08-security-secrets.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
