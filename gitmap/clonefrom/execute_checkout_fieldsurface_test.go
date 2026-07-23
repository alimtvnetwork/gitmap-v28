package clonefrom

// Compile-time guard for the field surface the checkout-execution
// tests in execute_checkout_test.go depend on. If any of these
// fields is renamed or removed on Result / Row, this file fails
// to BUILD before `go test ./clonefrom/...` ever runs — surfacing
// the breakage as a clean compile error pointing at the missing
// selector instead of an opaque "unknown field" deep inside a
// table-driven test.
//
// Why a separate file (not just inline asserts in the tests)?
//   - The tests cover several scenarios; a rename would produce
//     N near-identical errors instead of one. Centralizing the
//     selectors makes the failure point obvious.
//   - This file is `_test.go`, so it ships zero bytes in the
//     production binary while still participating in `go vet`
//     and `go build ./...` under -tags=test.
//   - Pairs with the broader reflect-based drift guard in
//     result_schema_drift_test.go: that file pins the REPORT
//     schema; this one pins the EXECUTION-test contract. They
//     fail independently so a regression names the right caller.

import (
	"testing"
	"time"
)

// resultFieldsUsedByCheckoutTests is a no-op compile-time witness:
// every field selector below MUST resolve at build time. Adding a
// new field selector here promotes a test-only assumption into a
// build-time invariant — do that whenever a checkout test starts
// reading a new Result/Row field.
//
// The function is intentionally never called. `var _ = ...` keeps
// the Go compiler from complaining about the unused declaration
// while still type-checking the body.
var _ = resultFieldsUsedByCheckoutTests

func resultFieldsUsedByCheckoutTests() {
	var r Result
	// Fields read by TestExecute_SkipCheckout_NoWorkingTree and
	// TestExecute_ForceCheckout_BranchMissingFails when asserting
	// per-row outcome.
	_ = r.Status   // string — compared against constants.CloneFromStatus*
	_ = r.Detail   // string — used in failure messages + prefix matching
	_ = r.Dest     // string — resolved dest path
	_ = r.Duration // time.Duration — surfaced in progress writer
	// Embedded Row fields read by buildGitArgs / EffectiveCheckout
	// table-driven cases.
	_ = r.Row.URL
	_ = r.Row.Branch
	_ = r.Row.Depth
	_ = r.Row.Dest
	_ = r.Row.Checkout
	// Pin the concrete Duration type — an accidental switch to
	// int64-nanoseconds would silently break the progress writer
	// and the JSON report's `duration_seconds` column.
	var _ time.Duration = r.Duration
}

// TestResultFieldSurface_CheckoutTests is a runtime mirror of the
// compile-time witness above. It runs in <1ms and exists so the
// failure mode is visible in `go test` output (not just in
// `go build`) when CI surfaces test results but not build logs.
func TestResultFieldSurface_CheckoutTests(t *testing.T) {
	r := Result{
		Row: Row{
			URL: "u", Branch: "b", Depth: 1, Dest: "d", Checkout: "auto",
		},
		Dest:     "resolved",
		Status:   "ok",
		Detail:   "",
		Duration: time.Millisecond,
	}
	if r.Status != "ok" || r.Dest != "resolved" || r.Row.URL != "u" {
		t.Fatalf("Result/Row field round-trip failed: %+v", r)
	}
}
