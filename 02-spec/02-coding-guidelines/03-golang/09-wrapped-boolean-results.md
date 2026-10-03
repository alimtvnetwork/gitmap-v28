# Golang Single Return Parameter & Wrapped Results (AI Execution Prompt)

> **/goal** Enforce single return parameter design across Go functions by replacing multi-value tuples `(T, error)` and raw `bool` returns with the standard monadic `Result[T]` container and `*appfault.AppError`.
> **/learn** Master the anti-patterns of raw boolean status returns (`(bool, error)`), understand mutually exclusive binary state flags (`IsSuccess` vs `IsFailed`), enforce constructor initialization (`NewSuccess`, `NewFailure`), and eliminate multi-value return clutter.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace multi-value return tuples with a single generic `Result[T]` wrapper or dedicated response struct.
- [ ] `/learn` Never return raw booleans (`(bool, error)`) to indicate binary execution state or success/failure flags.
- [ ] `/goal` Construct result instances strictly through `NewSuccess` or `NewFailure` helpers so `IsSuccess` and `IsFailed` cannot be mismatched.
- [ ] `/learn` Verify that errors are bound as structured `*appfault.AppError` and checked implicitly without explicit boolean comparisons.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Golang Overview](./readme.md)
> **Version:** 1.0.0

---

## The Rule: Single Return Parameter

In Go, we **strictly recommend returning a single parameter**.
If a function needs to return multiple values (e.g., a value and an error, or multiple pieces of data), they **must be bundled into a struct** or a generic wrapper.
The standard generic `Result` wrapper (described below) is the required pattern, as it bundles the `Data`, an `AppError`, and the `Status` together into a single return object.

## The Rule: No Raw Boolean Returns

In Go, functions that conceptually return a boolean to indicate success, failure, or a binary state **MUST NOT** return a raw `bool`.

Instead, they must return a generic wrapped result object that exposes a **status flag with two mutually exclusive parts**: a positive state (`IsSuccess`) and a negative state (`IsFailed`).

## The Wrapped Result Pattern

The wrapped result ensures that state is explicit and impossible to set to an invalid combination (e.g., both true or both false).

### Data Structure

```go
package result

// Result wraps a generic payload with explicit boolean state flags.
type Result[T any] struct {
	IsSuccess bool
	IsFailed  bool
	Data      T
	AppError  error // Or a custom AppError struct
}
```

### Initialization / Constructors (Crucial Step)

The user (developer) **never sets both flags manually**. The developer sets only one flag via a constructor method, and the other is automatically inferred/set.

```go
// NewSuccess creates a Result where IsSuccess is automatically true and IsFailed is false.
func NewSuccess[T any](data T) Result[T] {
	return Result[T]{
		IsSuccess: true,
		IsFailed:  false, // Automatically set as the inverse
		Data:      data,
	}
}

// NewFailure creates a Result where IsFailed is automatically true and IsSuccess is false.
func NewFailure[T any](err error) Result[T] {
	return Result[T]{
		IsSuccess: false, // Automatically set as the inverse
		IsFailed:  true,
		AppError: err,
	}
}
```

### Usage Example

**❌ FORBIDDEN: Raw Boolean Return**
```go
func ProcessPayment(amount int) (bool, error) {
    if amount <= 0 {
        return false, errors.New("invalid amount")
    }

    return true, nil
}
```

**✅ REQUIRED: Wrapped Result**
```go
func ProcessPayment(amount int) result.Result[PaymentReceipt] {
    if amount <= 0 {
        return result.NewFailure[PaymentReceipt](errors.New("invalid amount"))
    }

    receipt := PaymentReceipt{Amount: amount, Status: "cleared"}

    return result.NewSuccess[PaymentReceipt](receipt)
}

// Checking the status:
paymentStatus := ProcessPayment(100)

if paymentStatus.IsFailed {
    log.Println("Payment failed:", paymentStatus.AppError)
} else if paymentStatus.IsSuccess {
    log.Println("Payment cleared:", paymentStatus.Data)
}
```

## Why?

1. **Clarity**: It forces explicit checking of `.IsSuccess` or `.IsFailed` instead of a cryptic `if ok { ... }`.
2. **Extensibility**: You can add metadata, logging contexts, or payload data (`T`) without changing the function signature.
3. **Safety**: By forcing the use of `NewSuccess()` or `NewFailure()`, it is structurally impossible for a developer to accidentally set `IsSuccess = true` and `IsFailed = true` at the same time.
4. **Single Return:** Bundling data and errors into a single struct (like `Result`) strictly adheres to the single return parameter rule, keeping function signatures clean and predictable.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-009: Wrapped Result Monad and Single Return Parameter Standards

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Functions adhere to single return parameter signatures, raw boolean returns are replaced by monadic `Result[T]` containers, and error states bind to `*appfault.AppError`.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.
