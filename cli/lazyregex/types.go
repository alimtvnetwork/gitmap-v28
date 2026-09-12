// Package lazyregex — types.go centralizes regex domain models, envelopes, and Result aliases.
package lazyregex

import (
	"regexp"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// LazyRegexp provides a thread-safe, lazily compiled regular expression
// that caches its compiled state within the instance itself.
type LazyRegexp struct {
	expression string
	regex      *regexp.Regexp
	compileErr error
	isCompiled bool
	locker     sync.Mutex
}

// CompileResult encapsulates the outcome of a lazy regex compilation,
// holding either the compiled regexp or structured AppError diagnostics.
type CompileResult struct {
	re       *regexp.Regexp
	appError *apperror.AppError
}

// RegexpResult wraps a compiled regular expression in a Result envelope.
type RegexpResult = result.Result[*regexp.Regexp]
