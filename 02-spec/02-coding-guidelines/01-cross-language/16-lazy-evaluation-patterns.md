# Lazy Evaluation Patterns (AI Execution Prompt)

> **/goal** Enforce lazy evaluation patterns for expensive computations, static shared data, and optional fields across codebases to eliminate eager initialization bottlenecks and reduce memory overhead.
> **/learn** Master lazy evaluation triggers (heavy computations, invariant results, multi-caller lookups), non-exported field backing with getter methods, thread-safe synchronization (`sync.Mutex`), and avoid eager anti-patterns.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Identify expensive computations, static lookups, and optional payload fields that qualify for deferred lazy evaluation.
- [ ] `/learn` Ensure lazy fields remain non-exported with thread-safe getter methods returning cached pointers or collections.
- [ ] `/goal` Eliminate direct access to uninitialized private backing fields and replace eager constructors with deferred resolution.
- [ ] `/learn` Verify that dynamic, per-request, or parameterized results are never cached under lazy evaluation patterns.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

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
- [Master Coding Guidelines](./15-master-coding-guidelines/readme.md) — §7 Type Safety

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ARCH-016: Lazy Evaluation and Deferred Initialization

**Given** Data models and service structures containing expensive computations, static lookups, or optional fields.
**When** Code guideline linters or CI autofixers scan struct definitions and initialization patterns.
**Then** Expensive or optional fields utilize deferred lazy getters with non-exported backing fields and thread safety, ensuring zero eager initialization bottlenecks with deterministic compliance.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.

---

*Lazy evaluation patterns — consolidated from pre-code review guides.*
