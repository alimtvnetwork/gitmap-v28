# Nesting Resolution Patterns

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Applies to:** All languages  
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`

---

## 1. Principle

**Multiple nesting or branching makes code complex**, even if the logic is simple. Nesting increases cyclomatic complexity, complicates unit testing, and reduces maintainability. **Zero nested `if` is the target.**

**Exception:** One loop with a single `if` inside is acceptable. But not multiple `if` statements.

---

## 2. Three Resolution Methods

### Method 1: Extract to Named Function

```go
// ❌ NESTED — hard to read and test
func (wiki *WikiTimezone) Longitude() string {
    if len(wiki.LatLong) > 1 {
        if len(wiki.LatLong) > 11 {
            return string(wiki.LatLong[7:])
        } else {
            return string(wiki.LatLong[5:])
        }
    }

    return ""
}

// ✅ FLAT — extracted to named function
func isLatLongLengthGreaterThanEleven(wiki *WikiTimezone) bool {
    return len(wiki.LatLong) > 11
}

func (wiki *WikiTimezone) Longitude() string {
    if isLatLongLengthGreaterThanEleven(wiki) {
        return string(wiki.LatLong[7:])
    } else if len(wiki.LatLong) > 1 {
        return string(wiki.LatLong[5:])
    }

    return ""
}
```

### Method 2: Inverse Logic (Early Return)

```go
// ✅ FLAT — negative case exits first
func isEmpty(wiki *WikiTimezone) bool {
    return len(wiki.LatLong) <= 4
}

func (wiki *WikiTimezone) Longitude() string {
    if isEmpty(wiki) {
        return ""
    }

    if len(wiki.LatLong) > 11 {
        return string(wiki.LatLong[7:])
    }

    return string(wiki.LatLong[5:])
}
```

**Key insight:** Avoid `else` when not necessary. If the `if` block returns, the code after it is implicitly the else case.

### Method 3: Named Boolean Variables

```go
// ❌ NESTED — complex inline condition
if command != "" && (strings.HasPrefix(field, "-") || field == consts.NegateMatch) {
    // ...
}

// ✅ FLAT — decomposed into named booleans
hasValidCommand := command != ""
hasDashPrefix := strings.HasPrefix(field, "-")
hasNegateMatch := field == consts.NegateMatch
hasDashPrefixOrNegate := hasDashPrefix || hasNegateMatch

if hasValidCommand && hasDashPrefixOrNegate {
    // ...
}
```

---

## 3. Avoid Redundant `else` After Return — Absolute Rule

If an `if` block ends with `return`, `throw`, `break`, or `continue`, the code after the block is implicitly the "else" branch. Adding an explicit `else` is **forbidden** — it adds nesting and visual noise for no benefit.

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN — redundant else after return
func getValue(val []string, index int) string {
    if val[index] != "" {
        return val[index]
    } else {
        return "-"
    }
}

// ✅ REQUIRED — no else needed
func getValue(val []string, index int) string {
    if val[index] != "" {
        return val[index]
    }

    return "-"
}

// ❌ FORBIDDEN — else after early return
func process(data *Data) error {
    if data == nil {
        return ErrNilData
    } else {
        return execute(data)
    }
}

// ✅ REQUIRED — flat
func process(data *Data) error {
    if data == nil {
        return ErrNilData
    }

    return execute(data)
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
const getLabel = (status: string): string => {
    if (status === 'active') {
        return 'Active';
    } else {
        return 'Inactive';
    }

};

// ✅ REQUIRED
const getLabel = (status: string): string => {
    if (status === 'active') {
        return 'Active';
    }

    return 'Inactive';
};
```

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
if ($user === null) {
    return $this->error('User not found');
} else {
    return $this->success($user);
}

// ✅ REQUIRED
if ($user === null) {
    return $this->error('User not found');
}

return $this->success($user);
```

### When `else` IS Acceptable

`else` is only acceptable when **neither branch returns** — i.e., both branches assign a value and execution continues:

```go
// ✅ OK — else is needed because neither branch returns
var label string

if isAdmin {
    label = "Administrator"
} else {
    label = "User"
}

// ✅ BETTER — use ternary or conditional assignment when possible
// (In Go, the if/else form above is idiomatic since Go lacks ternary)
```

---

## 4. Complex Nesting Resolution (Level 4 → Level 1)

For deeply nested code, combine methods:

1. **Identify the exit conditions** — extract negative/guard cases first
2. **Name intermediate booleans** — replace inline compound conditions
3. **Extract helper functions** — each does one thing
4. **Use early returns** — flatten the remaining logic

See [How to Reduce Code Nested Branching in Go](https://hackmd.io/@akarimevatix/golang-code-refactoring-reduce-branching-v1) for a full walkthrough.

---

## 5. Cross-References

- [Code Style §R2](./04-code-style/00-overview.md) — Zero nested `if` rule
- [Cyclomatic Complexity](./06-cyclomatic-complexity.md) — Complexity limits
- [Boolean Principles](./02-boolean-principles/00-overview.md) — Named boolean extraction
- [Master Coding Guidelines §5](./15-master-coding-guidelines/00-overview.md) — Formatting rules

---

*Nesting resolution patterns — consolidated from pre-code review guides.*
