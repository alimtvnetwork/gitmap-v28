package lazyregex

import (
	"errors"
	"regexp"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	globalLock sync.Mutex
	globalMap  = make(map[string]*LazyRegexp, 64)
)

// New creates or retrieves a globally cached LazyRegexp for the given expression.
// Each expression maps to exactly one instance; compilation is lazy and stored in the instance.
func New(expression string) *LazyRegexp {
	globalLock.Lock()
	defer globalLock.Unlock()

	item, exists := globalMap[expression]
	if exists {
		return item
	}

	item = &LazyRegexp{expression: expression}
	globalMap[expression] = item

	return item
}

// NewLock creates or retrieves a cached LazyRegexp with mutex locking (alias for New).
func NewLock(expression string) *LazyRegexp {
	return New(expression)
}

// Compile compiles the regular expression on demand, setting the isCompiled flag.
// Subsequent calls return the cached *regexp.Regexp without recompilation.
// It checks first if an existing compiled regexp already exists and returns it immediately.
func (it *LazyRegexp) Compile() RegexpResult {
	if it == nil {
		return result.FailureResult[*regexp.Regexp](apperror.NewSimple("nil LazyRegexp cannot compile", "E9000"))
	}

	it.locker.Lock()
	defer it.locker.Unlock()

	if it.isCompiled && it.compileErr != nil {
		return result.FailureResult[*regexp.Regexp](apperror.WrapSimple(it.compileErr, "lazyregex.Compile"))
	}

	if it.isCompiled {
		return result.SuccessResult(it.regex)
	}

	compiled, err := regexp.Compile(it.expression)
	it.regex = compiled
	it.compileErr = err
	it.isCompiled = true

	if err != nil {
		appErr := apperror.WrapSimple(err, "lazyregex.Compile").
			WithContext("pattern", it.expression)

		return result.FailureResult[*regexp.Regexp](appErr)
	}

	return result.SuccessResult(compiled)
}

// CompileMust compiles the regular expression and panics on compilation error.
// It checks first if an existing compiled regexp already exists and returns it immediately.
func (it *LazyRegexp) CompileMust() *regexp.Regexp {
	if it == nil {
		return nil
	}

	res := it.Compile()
	res.HandleError()

	return res.Value
}

// CompileResult compiles the regex and returns a wrapped CompileResult.
// It checks first if an existing compiled regexp already exists and returns it immediately.
func (it *LazyRegexp) CompileResult() RegexpResult {
	return it.Compile()
}

// Re returns the underlying compiled *regexp.Regexp, panicking on compilation error.
// Kept for backward compatibility with existing callers.
func (it *LazyRegexp) Re() *regexp.Regexp {
	return it.CompileMust()
}

// Regex returns the underlying compiled *regexp.Regexp.
func (it *LazyRegexp) Regex() *regexp.Regexp {
	return it.CompileMust()
}

// Compiled returns the underlying compiled *regexp.Regexp.
func (it *LazyRegexp) Compiled() *regexp.Regexp {
	return it.CompileMust()
}

func (it *LazyRegexp) compiledRegex() (*regexp.Regexp, error) {
	if it == nil {
		return nil, errors.New("nil LazyRegexp")
	}

	res := it.Compile()
	if res.IsFailure() {
		return nil, res.AppError()
	}

	return res.Value, nil
}

// IsCompiled reports whether compilation has already been executed.
func (it *LazyRegexp) IsCompiled() bool {
	if it == nil {
		return false
	}

	it.locker.Lock()
	defer it.locker.Unlock()

	return it.isCompiled
}

// String returns the raw regular expression pattern string.
func (it *LazyRegexp) String() string {
	if it == nil {
		return ""
	}

	return it.expression
}

// Pattern returns the raw regular expression pattern string.
func (it *LazyRegexp) Pattern() string {
	if it == nil {
		return ""
	}

	return it.expression
}

// IsMatch reports whether the string s matches the regular expression without panicking.
func MatchPattern(pattern, comparing string) *MatchResult {
	return New(pattern).MatchResult(comparing)
}

// MatchError returns nil on match or a descriptive AppError on mismatch.
func MatchError(pattern, comparing string) error {
	return New(pattern).MatchError(comparing)
}

// Count returns the number of non-overlapping matches of the regular expression in s.
func CacheLen() int {
	globalLock.Lock()
	defer globalLock.Unlock()

	return len(globalMap)
}

// ClearCache flushes the global pattern registry (primarily for testing).
func ClearCache() {
	globalLock.Lock()
	defer globalLock.Unlock()

	globalMap = make(map[string]*LazyRegexp, 64)
}
