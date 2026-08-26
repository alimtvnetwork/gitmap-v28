# Cross-Language Rule: Reduce Cyclomatic Complexity to Near-Zero

> **Version:** 1.0.0  
> **Updated:** 2026-03-09  
> **Applies to:** PHP, TypeScript, Go, C#, and any delegated language

---

## Overview

Cyclomatic complexity measures the number of linearly independent paths through a function. High complexity means more branches, more bugs, harder testing, and slower code review. This project enforces a **target complexity of 0–1** per function by using **early returns (guard clauses)** to eliminate nesting entirely.

---

## The Problem: Deeply Nested Code

The following pattern is **absolutely forbidden** in this project. Every nested `if` adds a branch, increases indentation, and forces readers to mentally track multiple conditions simultaneously.

### ❌ BAD — Deeply Nested (Cyclomatic Complexity: 5)

```csharp
public void Process(Order? order)
{
    if (order != null)
    {
        if (order.IsVerified)
        {
            if (order.Items.Count > 0)
            {
                if (order.Items.Count > 15)
                {
                    throw new Exception(
                        "The order " + order.Id + " has too many items");
                }

                if (order.Status != "ReadyToProcess")
                {
                    throw new Exception(
                        "The order " + order.Id + " isn't ready to process");
                }

                order.IsProcessed = true;
            }
        }
    }
}
```

**Problems:**
- 5 levels of nesting — reader must track all conditions mentally
- The "happy path" is buried at the deepest level
- Adding a new validation requires finding the right nesting level
- Testing requires complex setup to reach each branch

---

## The Solution: Guard Clauses With Early Returns

Invert every condition and return/throw early. The function body becomes a **flat sequence of guards** followed by the happy path at the bottom with zero indentation.

### ✅ GOOD — Flat Guard Clauses (Cyclomatic Complexity: 1)

```csharp
public void Process(Order? order)
{
    if (order == null) {
        return;
    }

    // Use positive naming: IsUnverified instead of !IsVerified
    // The Order class should provide both IsVerified and IsUnverified accessors
    if (order.IsUnverified) {
        return;
    }

    // order is valid and verified
    var isItemListEmpty = order.Items.Count <= 0;

    if (isItemListEmpty) {
        return;
    }

    // order has items
    var hasExcessiveItems = order.Items.Count > 15;

    if (hasExcessiveItems) {
        throw new Exception(
            "The order " + order.Id + " has too many items");
    }

    var isStatusMismatch = order.Status.HasMismatch(OrderStatusType.ReadyToProcess);

    if (isStatusMismatch) {
        throw new Exception(
            "The order " + order.Id + " is not ready to process");
    }

    // All guards passed — process the order
    order.IsProcessed = true;
}
```

**Benefits:**
- Zero nesting — every guard is at the same indentation level
- Happy path is at the bottom, clearly visible
- Adding a new guard is trivial — just add another `if (...) { return; }` block
- Each guard is independently testable
- Comments between guards describe what's been established so far

---

## The Pattern: Guard → Early Exit → Happy Path

```
function doSomething(input) {
    // Guard 1: Reject null/empty
    if (input is invalid) { return / throw }

    // Guard 2: Reject unauthorized
    if (input lacks permission) { return / throw }

    // Guard 3: Reject business rule violation
    if (input exceeds limit) { throw }

    // Guard 4: Reject wrong state
    if (input state is wrong) { throw }

    // ── All guards passed ──
    // Happy path executes here with zero nesting
    execute(input)
}
```

---

## Language-Specific Examples

### PHP

```php
// ❌ FORBIDDEN — Nested validation
public function processUpload(WP_REST_Request $request): WP_REST_Response {
    if ($request->has_param('file')) {
        $file = $request->get_param('file');
        if ($file['size'] > 0) {
            if ($file['size'] < MAX_UPLOAD_SIZE) {
                // process...
            }
        }
    }
}

// ✅ REQUIRED — Flat guard clauses
public function processUpload(WP_REST_Request $request): WP_REST_Response {
    // Use positive naming: isFileMissing instead of hasNoFile or !has_param
    $isUploadFileMissing = !$request->has_param('file');

    if ($isUploadFileMissing) {
        return $this->envelope->error('No file provided', 400);
    }

    $file = $request->get_param('file');
    $isFileEmpty = $file['size'] <= 0;

    if ($isFileEmpty) {
        return $this->envelope->error('Empty file', 400);
    }

    $isFileTooLarge = $file['size'] >= MAX_UPLOAD_SIZE;

    if ($isFileTooLarge) {
        return $this->envelope->error('File exceeds limit', 413);
    }

    // All guards passed — process the upload
    $result = $this->executeUpload($file);

    return $this->envelope->success($result);
}
```

### TypeScript

```typescript
// ❌ FORBIDDEN — Nested validation
const processOrder = (order: Order | null) => {
    if (order) {
        if (order.items.length > 0) {
            if (order.status === 'ready') {
                // process...
            }
        }
    }
};

// ✅ REQUIRED — Flat guard clauses
const processOrder = (order: Order | null) => {
    if (isNullish(order)) {
        return;
    }

    if (isArrayEmpty(order.items)) {
        return;
    }

    // Use positive naming: isStatusPending instead of negating 'ready'
    const isStatusPending = order.status !== 'ready';

    if (isStatusPending) {
        throw new Error(`Order ${order.id} is not ready`);
    }

    // All guards passed
    order.isProcessed = true;
};
```

### Go

```go
// ❌ FORBIDDEN — Nested validation
func ProcessOrder(order *Order) error {
    if order != nil {
        if order.IsVerified {
            if len(order.Items) > 0 {
                // process...
            }
        }
    }

    return nil
}

// ✅ REQUIRED — Flat guard clauses
func ProcessOrder(order *Order) error {
    if order == nil {
        return apperror.New(
            "E4001", "order is nil",
        )
    }

    // Use positive naming: IsUnverified instead of !IsVerified
    // The Order struct should provide both IsVerified() and IsUnverified() accessors
    if order.IsUnverified() {
        return apperror.New(
            "E4002", "order is not verified",
        )
    }

    isItemListEmpty := IsSliceEmpty(order.Items)

    if isItemListEmpty {
        return apperror.New(
            "E4003", "order has no items",
        )
    }

    hasExcessiveItems := len(order.Items) > MaxOrderItems

    if hasExcessiveItems {
        return apperror.New(
            "E4004", "order exceeds item limit",
        )
    }

    isStatusMismatch := order.Status.HasMismatch(StatusReadyToProcess)

    if isStatusMismatch {
        return apperror.New(
            "E4005", "order is not ready to process",
        )
    }

    // All guards passed
    order.IsProcessed = true

    return nil
}
```

---

## Complexity Scoring

| Complexity | Rating | Action Required |
|-----------|--------|-----------------|
| 0–1 | ✅ Ideal | No action needed |
| 2–3 | ⚠️ Acceptable | Consider refactoring |
| 4+ | ❌ Rejected | Must refactor before merge |

### How Guards Reduce Complexity

| Pattern | Before (Nested) | After (Guards) |
|---------|-----------------|----------------|
| Null check → validate → process | 3 branches, 3 levels deep | 3 guards, 0 nesting |
| Permission → role → action | 3 branches, 3 levels deep | 3 guards, 0 nesting |
| 5 validations → process | 5 branches, 5 levels deep | 5 guards, 0 nesting |

Each guard clause adds a branch but **removes a nesting level**, keeping the function body flat and linear.

---

## Relationship to Other Rules

This spec works in concert with:

| Rule | How It Helps |
|------|-------------|
| [Rule 2: Zero Nested `if`](./04-code-style/01-braces-and-nesting.md) | Enforces the same principle at the syntax level |
| [Rule 6: 15-Line Function Limit](./04-code-style/04-function-and-type-size.md) | Guards keep functions short — if too many guards, extract a `validate()` helper |
| [Boolean Principles](./02-boolean-principles/00-overview.md) | Guards use `is`/`has` named booleans for clarity |
| [No Raw Negations](./12-no-negatives.md) | Guards use positive names (`isFileMissing`) not `!fileExists()` |

---

## Checklist Summary (Copy for PRs)

```
[ ] No nested `if` blocks — all conditions are guard clauses
[ ] Happy path is at the bottom of the function with zero nesting
[ ] Each guard uses early return or throw
[ ] Complex conditions are extracted to named booleans
[ ] Function complexity is 0–1 (max 3 in exceptional cases)
[ ] Comments between guards describe the established state
```

---

## Cross-References

- [Code Style — Rules 2, 6, 7](./04-code-style/00-overview.md) — Zero nesting, 15-line limit
- [Boolean Principles](./02-boolean-principles/00-overview.md) — Named boolean variables
- [No Raw Negations](./12-no-negatives.md) — Positive guard function names
- [Function Naming](./10-function-naming.md) — Explicit method names over boolean flags

---

*Cyclomatic complexity specification v1.0.0 — 2026-02-17*
