# Lazy Evaluation Patterns

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Applies to:** Go (primary), general principle cross-language  
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`

---

## 1. Principle

Lazy evaluation means **don't execute until needed**. It acts as caching — generate once, serve the same data to many callers.

---

## 2. When to Apply Lazy Evaluation

### ✅ APPLY when:

| Condition | Example |
|-----------|---------|
| Heavy lifting / expensive computation | Database query, file parsing |
| Static value that doesn't change | Configuration, computed constants |
| Same data requested by many callers | Shared lookup tables |
| 60%+ cases don't need the value at all | Optional expensive fields |
| Builder pattern — instructions collected, executed once | `Builder.Build()` |

### ❌ DO NOT apply when:

| Condition | Reason |
|-----------|--------|
| Data changes per request | Cannot cache varying results |
| Function requires parameters to produce different results | Not cacheable |
| Value is cheap to compute | Overhead of lazy machinery not justified |

---

## 3. Go Implementation Pattern

### Step 1: Make Field Non-Exported

```go
type Group struct {
    GroupId   string
    GroupName string
    members   *UsersCollection  // non-exported — lazy
    sync.Mutex                  // lock for async safety
}
```

### Step 2: Expose via Getter Method

```go
func (g *Group) Members() *UsersCollection {
    if g.members == nil {
        g.members = generateMembers(g.GroupId)
    }

    return g.members
}
```

### Step 3: Always Use Method, Never Direct Field

```go
// ✅ CORRECT — uses getter
func (g *Group) IsMembersEmpty() bool {
    return g.Members().Length() == 0
}

// ❌ WRONG — accesses field directly, bypasses lazy init
func (g *Group) IsMembersEmpty() bool {
    return g.members.Length() == 0
}
```

### Step 4: Lock for Concurrent Access

```go
func (g *Group) MembersLock() *UsersCollection {
    g.Lock()
    defer g.Unlock()

    return g.Members()
}
```

---

## 4. Critical Rule: Lazy Field Dependencies

**If a required field is lazy, the current field MUST also be lazy.**

```go
// ❌ DANGEROUS — eager field depends on lazy field
type Report struct {
    data    *Data           // lazy
    Summary string          // eager — but computed from data!
}

// ✅ CORRECT — both lazy
type Report struct {
    data    *Data           // lazy
    summary string          // also lazy — depends on data
}

func (r *Report) Summary() string {
    if r.summary == "" {
        r.summary = computeSummary(r.Data())
    }

    return r.summary
}
```

---

## 5. Anti-Patterns

### ❌ Eager Initialization of Expensive Fields

```go
// WRONG — always computes even if never accessed
type Group struct {
    Members *UsersCollection  // exported, always initialized
}

func NewGroup(id string) *Group {
    return &Group{
        Members: loadAllMembers(id),  // expensive!
    }
}
```

### ❌ Direct Field Access Bypassing Lazy Getter

```go
// WRONG — skips nil check, will panic
if g.members.Length() > 0 { ... }

// CORRECT — uses lazy getter
if g.Members().Length() > 0 { ... }
```

---

## 6. Cross-References

- [Code Mutation Avoidance](./18-code-mutation-avoidance.md) — Lazy fields are an exempted mutation case
- [Cyclomatic Complexity](./06-cyclomatic-complexity.md) — Lazy getters keep callers simple
- [Master Coding Guidelines](./15-master-coding-guidelines/00-overview.md) — §7 Type Safety

---

*Lazy evaluation patterns — consolidated from pre-code review guides.*
