# Boolean Flag Method Splitting (AI Execution Prompt)

> **/goal** Eliminate all boolean flag parameters that alter method behavior by splitting them into dedicated, self-documenting methods that express explicit caller intent.
> **/learn** Master the single-responsibility principle for functions, identify anti-patterns of hidden branching caused by boolean arguments, extract shared initialization/teardown into private helpers, and understand exemptions (options structs, thin wrappers, state toggles).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Identify and eliminate all boolean parameters that branch execution paths inside functions and methods.
- [ ] `/learn` Split flagged methods into two distinct, descriptive methods (e.g. `formatUserSummary` and `formatUserDetailed`).
- [ ] `/goal` Extract any common setup, validation, or teardown logic into private non-exported helper functions.
- [ ] `/learn` Restrict boolean parameters strictly to options/config structs, standard library pass-throughs, or explicit state setters (`setEnabled`).

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Cross-Language Overview](./readme.md)
> **Version:** 1.0.0
> **Updated:** 2026-04-02
> **AI Confidence:** Production-Ready
> **Ambiguity:** None

## Keywords

`boolean-flag` · `method-splitting` · `single-responsibility` · `function-design` · `clean-code`

---

## Rule

**🔴 CODE RED:** If a method's behavior changes based on a boolean parameter, split it into two named methods that express each intent explicitly.

Boolean flags hide branching logic inside function calls. The caller cannot understand what `true` or `false` means without reading the implementation. Two named methods make intent obvious at every call site.

---

## The Problem

```
// ❌ FORBIDDEN — What does `true` mean here?
processOrder(order, true)
processOrder(order, false)

sendEmail(user, true)
sendEmail(user, false)
```

The reader must open the function to understand the flag. This violates self-documenting code principles and makes code reviews slower.

---

## The Rule: Split Into Two Methods

Every boolean flag parameter that changes method behavior must be replaced with two methods whose names describe the behavior.

### Go

```go
// ❌ FORBIDDEN — boolean flag hides intent
func ProcessOrder(order Order, isPriority bool) error {
    if isPriority {
        // priority logic
    } else {
        // standard logic
    }
}

// ✅ REQUIRED — two methods, intent is clear
func ProcessPriorityOrder(order Order) error {
    // priority logic
}

func ProcessStandardOrder(order Order) error {
    // standard logic
}
```

### TypeScript

```typescript
// ❌ FORBIDDEN
function formatUser(user: User, isDetailed: boolean): string {
    if (isDetailed) {
        return `${user.name} (${user.email}, ${user.role})`;
    }
    return user.name;
}

// ✅ REQUIRED
function formatUserSummary(user: User): string {
    return user.name;
}

function formatUserDetailed(user: User): string {
    return `${user.name} (${user.email}, ${user.role})`;
}
```

### PHP

```php
// ❌ FORBIDDEN
function syncPlugin(Plugin $plugin, bool $isForced): void {
    if ($isForced) {
        // force sync logic
    } else {
        // incremental sync logic
    }
}

// ✅ REQUIRED
function syncPluginIncremental(Plugin $plugin): void {
    // incremental sync logic
}

function syncPluginForced(Plugin $plugin): void {
    // force sync logic
}
```

### Rust

```rust
// ❌ FORBIDDEN
fn write_log(entry: &LogEntry, is_verbose: bool) {
    if is_verbose {
        // verbose output
    } else {
        // compact output
    }
}

// ✅ REQUIRED
fn write_log_compact(entry: &LogEntry) {
    // compact output
}

fn write_log_verbose(entry: &LogEntry) {
    // verbose output
}
```

### C#

```csharp
// ❌ FORBIDDEN
public void SaveDocument(Document doc, bool isDraft)
{
    if (isDraft) { /* draft logic */ }
    else { /* publish logic */ }
}

// ✅ REQUIRED
public void SaveDraft(Document doc)
{
    // draft logic
}

public void PublishDocument(Document doc)
{
    // publish logic
}
```

---

## When Shared Logic Exists

If both paths share setup or teardown, extract the shared logic into a private helper:

```go
// ✅ Shared logic extracted
func ProcessPriorityOrder(order Order) error {
    validateOrder(order)        // shared
    applyPriorityDiscount(order) // unique
    return finalizeOrder(order)  // shared
}

func ProcessStandardOrder(order Order) error {
    validateOrder(order)         // shared
    applyStandardPricing(order)  // unique
    return finalizeOrder(order)  // shared
}

// Private shared helpers
func validateOrder(order Order) { /* ... */ }
func finalizeOrder(order Order) error { /* ... */ }
```

---

## Exemptions

| Case | Reason |
|------|--------|
| **Options/config structs** | Booleans inside an options struct are acceptable — caller sees named fields (`Config{Verbose: true}`) |
| **Standard library wrappers** | Thin wrappers around stdlib that pass through bool params (e.g., `os.OpenFile` flags) |
| **Toggle methods** | Methods that flip state (`SetEnabled(bool)`) where the name already describes intent |

---

## Cross-References

- [Boolean Principles](./02-boolean-principles/readme.md) — P5: No boolean parameters
- [Function Naming](./10-function-naming.md) — naming conventions for split methods
- [Cyclomatic Complexity](./06-cyclomatic-complexity.md) — flag splitting reduces branching
- [SOLID Principles](./23-solid-principles.md) — Single Responsibility applied to methods
- [Nesting Resolution](./20-nesting-resolution-patterns.md) — related pattern: flatten `if/else`

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-BOOL-024: Boolean Flag Method Splitting

**Given** Function or method declarations across Go, TypeScript, PHP, Rust, or C#.
**When** Code guideline linters or CI autofixers inspect function signatures for boolean parameters.
**Then** Zero boolean flag parameters altering control flow are permitted; behavior-altering branches are split into separate named methods with deterministic compliance.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
