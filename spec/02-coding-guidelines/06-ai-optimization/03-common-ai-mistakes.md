# Common AI Mistakes

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Purpose:** Real examples of mistakes AI commonly makes when generating code for this project

---

## Top 15 AI Mistakes (Ranked by Frequency)

---

### Mistake #1: camelCase JSON Keys

**Frequency:** Very High  
**Rule:** AH-N3

```json
// ❌ AI GENERATES THIS
{ "userId": "abc", "createdAt": "2026-01-01" }

// ✅ CORRECT
{ "UserId": "abc", "CreatedAt": "2026-01-01" }
```

---

### Mistake #2: Uppercase Abbreviations

**Frequency:** Very High  
**Rule:** AH-N1

```go
// ❌ AI GENERATES THIS
type APIResponse struct {
    UserID   string `json:"userId"`
    BaseURL  string `json:"baseUrl"`
}

// ✅ CORRECT
type ApiResponse struct {
    UserId  string
    BaseUrl string
}
```

---

### Mistake #3: Multi-Return Go Functions

**Frequency:** High  
**Rule:** AH-T3

```go
// ❌ AI GENERATES THIS
func GetUser(id string) (*User, error) {
    // ...
    return user, nil
}

// ✅ CORRECT
func GetUser(id string) apperror.Result[User] {
    // ...
    return apperror.Ok(user)
}
```

---

### Mistake #4: `fmt.Errorf` Instead of `apperror`

**Frequency:** High  
**Rule:** AH-T2

```go
// ❌ AI GENERATES THIS
return fmt.Errorf("failed to process: %w", err)

// ✅ CORRECT
return apperror.Wrap(err, apperror.ErrProcessFailed, "failed to process")
```

---

### Mistake #5: Nested `if` Statements

**Frequency:** High  
**Rule:** AH-S1

```go
// ❌ AI GENERATES THIS
if request != nil {
    if request.IsValid() {
        process(request)
    }
}

// ✅ CORRECT
if request == nil {
    return
}

if request.IsInvalid() {
    return
}

process(request)
```

---

### Mistake #6: Boolean Without Prefix

**Frequency:** High  
**Rule:** AH-N6

```typescript
// ❌ AI GENERATES THIS
const active = true;
const loading = false;
const valid = checkForm();

// ✅ CORRECT
const isActive = true;
const isLoading = false;
const isValid = checkForm();
```

---

### Mistake #7: Explicit Go JSON Tags

**Frequency:** Medium-High  
**Rule:** AH-T5

```go
// ❌ AI GENERATES THIS
type Config struct {
    MaxRetries int    `json:"MaxRetries"`
    BaseUrl    string `json:"BaseUrl"`
    Timeout    int    `json:"Timeout"`
}

// ✅ CORRECT — no tags needed, PascalCase is default
type Config struct {
    MaxRetries int
    BaseUrl    string
    Timeout    int
}
```

---

### Mistake #8: Using `any` / `interface{}`

**Frequency:** Medium  
**Rule:** AH-T1

```go
// ❌ AI GENERATES THIS
func ProcessData(data map[string]any) any {
    return data["result"]
}

// ✅ CORRECT
func ProcessData(data ProcessInput) apperror.Result[ProcessOutput] {
    return apperror.Ok(ProcessOutput{Result: data.Result})
}
```

---

### Mistake #9: Raw Negation on Function Calls

**Frequency:** Medium  
**Rule:** AH-B1

```php
// ❌ AI GENERATES THIS
if (!$order->isValid()) {
    return;
}

// ✅ CORRECT
if ($order->isInvalid()) {
    return;
}
```

---

### Mistake #10: Missing Blank Line Before Return

**Frequency:** Medium  
**Rule:** AH-S4

```go
// ❌ AI GENERATES THIS
func Process(data Input) Result {
    result := compute(data)

    return result
}

// ✅ CORRECT
func Process(data Input) Result {
    result := compute(data)

    return result
}
```

---

### Mistake #11: `\Throwable` in PHP

**Frequency:** Medium  
**Rule:** AH-E4

```php
// ❌ AI GENERATES THIS
try {
    $result = $this->process();
} catch (\Throwable $e) {
    Logger::error($e->getMessage());
}

// ✅ CORRECT
use Throwable;

try {
    $result = $this->process();
} catch (Throwable $e) {
    Logger::error($e->getMessage());
}
```

---

### Mistake #12: Magic String Status Comparisons

**Frequency:** Medium  
**Rule:** AH-EN3

```typescript
// ❌ AI GENERATES THIS
if (status === 'active') { ... }
if (connection.status === 'connected') { ... }

// ✅ CORRECT
if (status === EntityStatus.Active) { ... }
if (connection.status === ConnectionStatus.Connected) { ... }
```

---

### Mistake #13: String-Based Go Enums

**Frequency:** Low-Medium  
**Rule:** AH-EN1

```go
// ❌ AI GENERATES THIS
type Provider string

const (
    SerpApi Provider = "serpapi"
    Colly   Provider = "colly"
)

// ✅ CORRECT
type Variant byte

const (
    Invalid Variant = iota
    SerpApi
    Colly
)
```

---

### Mistake #14: Value Access Without Error Guard

**Frequency:** Low-Medium  
**Rule:** AH-E1

```go
// ❌ AI GENERATES THIS
result := svc.GetUser(ctx, id)
user := result.Value()  // may be zero if error!

// ✅ CORRECT
result := svc.GetUser(ctx, id)
if result.HasError() {
    return result
}

user := result.Value()
```

---

### Mistake #15: Functions with 4+ Parameters

**Frequency:** Low  
**Rule:** AH-S3

```typescript
// ❌ AI GENERATES THIS
function createUser(name: string, email: string, role: string, department: string): User { ... }

// ✅ CORRECT
interface CreateUserParams {
    name: string;
    email: string;
    role: string;
    department: string;
}

function createUser(params: CreateUserParams): User { ... }
```

---

## 🔴 Caching Mistakes (CODE RED)

---

### Mistake #16: Caching Errors as Success

**Frequency:** High  
**Rule:** AH-CA1

```typescript
// ❌ AI GENERATES THIS
try {
    const data = await fetchUsers();
    cache.set("users", data);
} catch {
    cache.set("users", []); // Silent failure cached as success!
}

// ✅ CORRECT
try {
    const data = await fetchUsers();
    cache.set("users", data, { ttl: 300_000 });
} catch (error) {
    cache.delete("users");
    logger.error("fetchUsers failed", { error });
    throw error;
}
```

---

### Mistake #17: Cache Without TTL

**Frequency:** High  
**Rule:** AH-CA2

```typescript
// ❌ AI GENERATES THIS
cache.set("config", configData);

// ✅ CORRECT
cache.set("config", configData, { ttl: 300_000 });
```

---

### Mistake #18: Missing Cache Invalidation After Mutation

**Frequency:** High  
**Rule:** AH-CA3

```typescript
// ❌ AI GENERATES THIS
const mutation = useMutation({
    mutationFn: updateUser,
    onSuccess: () => {
        toast.success("User updated");
    },
});

// ✅ CORRECT
const mutation = useMutation({
    mutationFn: updateUser,
    onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ["users"] });
        toast.success("User updated");
    },
});
```

---

### Mistake #19: Non-Deterministic Cache Keys

**Frequency:** Medium  
**Rule:** AH-CA4

```typescript
// ❌ AI GENERATES THIS
cache.set(`users-${Date.now()}`, data);

// ✅ CORRECT
cache.set(`users-${userId}-${role}`, data);
```

---

### Mistake #20: React Query Without Explicit staleTime

**Frequency:** High  
**Rule:** AH-CA6

```typescript
// ❌ AI GENERATES THIS — default staleTime: 0 causes refetch on every mount
const { data } = useQuery({
    queryKey: ["users"],
    queryFn: fetchUsers,
});

// ✅ CORRECT
const { data } = useQuery({
    queryKey: ["users"],
    queryFn: fetchUsers,
    staleTime: 5 * 60 * 1000,
});
```

---

## Pattern Recognition Guide

### Signs AI is About to Hallucinate

| Signal | Likely Mistake | Prevention |
|--------|---------------|------------|
| Generating Go struct with `json:` tags | Will use camelCase tags | Remove tags entirely |
| Generating Go function signature | Will use `(T, error)` return | Use `Result[T]` |
| Generating error handling in Go | Will use `fmt.Errorf` | Use `apperror.Wrap` |
| Generating JSON response body | Will use camelCase keys | Use PascalCase |
| Generating boolean variable | Will omit `is`/`has` prefix | Add prefix |
| Generating nested conditions | Will create nested `if` | Flatten with early returns |
| Generating Go enum | Will use `string` type | Use `byte` + `iota` |
| Generating cache `catch` block | Will cache empty array as success | Delete cache + rethrow |
| Generating `cache.set()` | Will omit TTL | Always include `{ ttl }` |
| Generating `useMutation` | Will skip `invalidateQueries` | Add invalidation in `onSuccess` |
| Generating `useQuery` | Will omit `staleTime` | Set explicit `staleTime` |

---

## Cross-References

- [Anti-Hallucination Rules](./01-anti-hallucination-rules.md) — Full rule catalog
- [AI Quick Reference Checklist](./02-ai-quick-reference-checklist.md) — Pre-output validation
- [Master Coding Guidelines](../01-cross-language/15-master-coding-guidelines/00-overview.md) — Complete reference

---

*Common AI mistakes v1.1.0 — 2026-04-04*
