package appfault

import (
	"fmt"
)

// isActualError checks if an AppError is non-nil and represents an actual failure.
func isActualError(err *AppError) bool {
	if err == nil {
		return false
	}

	return err.HasError()
}

// Merge combines two AppError instances into one immutable AppError.
// It checks whether each error actually exists:
//   - If neither exists, it returns nil.
//   - If only one exists, it returns that error unchanged.
//   - If both exist, it merges prev into next, preserving next's failure identity while
//     capturing prev's stack trace, caller, message, and accumulated loop count inside
//     the context dictionary ("FirstErrorStackTrace", "PriorStackTrace", "LoopCount", "StackTraceHistory").
func Merge(prev *AppError, next *AppError) *AppError {
	prevExists := isActualError(prev)
	nextExists := isActualError(next)

	if !prevExists && !nextExists {
		return nil
	}

	if !prevExists {
		return next
	}

	if !nextExists {
		return prev
	}

	// Both errors exist: clone next to preserve its failure identity and immutability
	merged := next.clone()

	// Inherit cause from prev if next has no explicit cause
	if merged.cause == nil && prev.cause != nil {
		merged.cause = prev.cause
	}

	prevStack := prev.StackTrace().String()
	prevCaller := prev.Caller().String()

	// Determine loop iteration count
	loopCount := 2
	if prev.ctx.Has("LoopCount") {
		if rawCount, ok := prev.ctx["LoopCount"].(int); ok {
			loopCount = rawCount + 1
		}
	}

	merged.ctx.Set("LoopCount", loopCount)

	// Preserve the very first error's stack trace
	if prev.ctx.Has("FirstErrorStackTrace") {
		merged.ctx.Set("FirstErrorStackTrace", prev.ctx["FirstErrorStackTrace"])
		if prev.ctx.Has("FirstErrorCaller") {
			merged.ctx.Set("FirstErrorCaller", prev.ctx["FirstErrorCaller"])
		}

		if prev.ctx.Has("FirstErrorMessage") {
			merged.ctx.Set("FirstErrorMessage", prev.ctx["FirstErrorMessage"])
		}
	} else {
		merged.ctx.Set("FirstErrorStackTrace", prevStack)
		merged.ctx.Set("FirstErrorCaller", prevCaller)
		merged.ctx.Set("FirstErrorMessage", prev.Message())
		merged.ctx.Set("FirstErrorType", prev.Type().Name())
	}

	// Store immediate prior error diagnostics
	merged.ctx.Set("PriorStackTrace", prevStack)
	merged.ctx.Set("PriorCaller", prevCaller)
	merged.ctx.Set("PriorMessage", prev.Message())

	// Accumulate comprehensive stack trace history across looping attempts
	var history []string
	if rawHistory, ok := prev.ctx["StackTraceHistory"].([]string); ok {
		history = append([]string{}, rawHistory...)
	} else {
		history = append(history, fmt.Sprintf("[Attempt #1] [%s:%d] %s at %s\n%s",
			prev.Type().Name(), prev.Type().Code(), prev.Message(), prevCaller, prevStack))
	}

	currentStack := next.StackTrace().String()
	history = append(history, fmt.Sprintf("[Attempt #%d] [%s:%d] %s at %s\n%s",
		loopCount, next.Type().Name(), next.Type().Code(), next.Message(), next.Caller().String(), currentStack))
	merged.ctx.Set("StackTraceHistory", history)

	// Copy custom domain keys from previous error if not already present in next
	for k, v := range prev.ctx {
		if !merged.ctx.Has(k) {
			merged.ctx.Set(k, v)
		}
	}

	return merged
}

// Merge combines the receiver error with another error, returning a new immutable AppError.
func (e *AppError) Merge(other *AppError) *AppError {
	return Merge(e, other)
}

// LoopCount returns the number of loop iterations or chained errors if merged, or 1 if unmerged.
func (e *AppError) LoopCount() int {
	if e == nil || e.ctx == nil {
		return 1
	}

	if val, ok := e.ctx["LoopCount"].(int); ok {
		return val
	}

	return 1
}

// FirstErrorStackTrace returns the stack trace of the first failure in a merged chain.
func (e *AppError) FirstErrorStackTrace() string {
	if e == nil || e.ctx == nil {
		return ""
	}

	if val, ok := e.ctx["FirstErrorStackTrace"].(string); ok {
		return val
	}

	return ""
}

// PriorStackTrace returns the stack trace of the immediately preceding failure in a merged chain.
func (e *AppError) PriorStackTrace() string {
	if e == nil || e.ctx == nil {
		return ""
	}

	if val, ok := e.ctx["PriorStackTrace"].(string); ok {
		return val
	}

	return ""
}

// StackTraceHistory returns the list of all captured stack traces across loop attempts.
func (e *AppError) StackTraceHistory() []string {
	if e == nil || e.ctx == nil {
		return nil
	}

	if val, ok := e.ctx["StackTraceHistory"].([]string); ok {
		return val
	}

	return nil
}
