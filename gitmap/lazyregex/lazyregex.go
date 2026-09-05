package lazyregex

import (
	"errors"
	"regexp"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

var (
	globalLock sync.Mutex
	globalMap  = make(map[string]*LazyRegexp, 64)
)

// LazyRegexp provides a thread-safe, lazily compiled regular expression
// that caches its compiled state within the instance itself.
type LazyRegexp struct {
	expr       string
	re         *regexp.Regexp
	compileErr error
	isCompiled bool
	mu         sync.Mutex
}

// New creates or retrieves a globally cached LazyRegexp for the given expression.
// Each expression maps to exactly one instance; compilation is lazy and stored in the instance.
func New(expr string) *LazyRegexp {
	globalLock.Lock()
	defer globalLock.Unlock()

	item, exists := globalMap[expr]
	if exists {
		return item
	}

	item = &LazyRegexp{expr: expr}
	globalMap[expr] = item

	return item
}

// NewLock creates or retrieves a cached LazyRegexp with mutex locking (alias for New).
func NewLock(expr string) *LazyRegexp {
	return New(expr)
}

// Compile compiles the regular expression on demand, setting the isCompiled flag.
// Subsequent calls return the cached *regexp.Regexp without recompilation.
func (r *LazyRegexp) Compile() (*regexp.Regexp, error) {
	if r == nil {
		return nil, errors.New("nil LazyRegexp cannot compile")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isCompiled {
		return r.re, r.compileErr
	}

	compiled, err := regexp.Compile(r.expr)
	r.re = compiled
	r.compileErr = err
	r.isCompiled = true

	return compiled, err
}

// CompileMust compiles the regular expression and panics on compilation error.
func (r *LazyRegexp) CompileMust() *regexp.Regexp {
	compiled, err := r.Compile()
	if err != nil {
		panic(err)
	}

	return compiled
}

// CompileAppError compiles the regex and returns a wrapped CompileResult with typed AppError on failure.
func (r *LazyRegexp) CompileAppError() *CompileResult {
	if r == nil {
		appErr := apperror.NewWithDetails("lazyregex.Compile", "NIL_RECEIVER", "nil LazyRegexp cannot compile", "lazyregex", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
		return NewCompileFailure(appErr)
	}

	compiled, err := r.Compile()
	if err != nil {
		appErr := apperror.Wrap(err, "lazyregex.Compile", map[string]any{
			"pattern": r.expr,
		})
		return NewCompileFailure(appErr)
	}

	return NewCompileSuccess(compiled)
}

// CompileResult compiles the regex and returns a wrapped CompileResult.
func (r *LazyRegexp) CompileResult() *CompileResult {
	return r.CompileAppError()
}

// Re returns the underlying compiled *regexp.Regexp, panicking on compilation error.
// Kept for backward compatibility with existing callers.
func (r *LazyRegexp) Re() *regexp.Regexp {
	return r.CompileMust()
}

// IsCompiled reports whether compilation has already been executed.
func (r *LazyRegexp) IsCompiled() bool {
	if r == nil {
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.isCompiled
}

// String returns the raw regular expression pattern string.
func (r *LazyRegexp) String() string {
	if r == nil {
		return ""
	}

	return r.expr
}

// Pattern returns the raw regular expression pattern string.
func (r *LazyRegexp) Pattern() string {
	if r == nil {
		return ""
	}

	return r.expr
}

// IsMatch reports whether the string s matches the regular expression without panicking.
func (r *LazyRegexp) IsMatch(s string) bool {
	if r == nil {
		return false
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return false
	}

	return re.MatchString(s)
}

// IsFound is a semantic alias for IsMatch.
func (r *LazyRegexp) IsFound(s string) bool {
	return r.IsMatch(s)
}

// MatchString reports whether the string s contains any match of the regular expression.
func (r *LazyRegexp) MatchString(s string) bool {
	return r.IsMatch(s)
}

// Count returns the number of non-overlapping matches of the regular expression in s.
func (r *LazyRegexp) Count(s string) int {
	if r == nil {
		return 0
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return 0
	}

	matches := re.FindAllString(s, -1)
	return len(matches)
}

// GroupBy extracts named capture groups (?P<name>...) from the first match into a GroupMap.
func (r *LazyRegexp) GroupBy(s string) *GroupMap {
	result := NewGroupMap()
	if r == nil {
		return result
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return result
	}

	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return result
	}

	names := re.SubexpNames()
	for i, name := range names {
		if name == "" || i >= len(match) {
			continue
		}
		result.Set(name, match[i])
	}

	return result
}

// FindGroups is an alias for GroupBy.
func (r *LazyRegexp) FindGroups(s string) *GroupMap {
	return r.GroupBy(s)
}

// FindAllGroups extracts named capture groups across all non-overlapping matches in s into a GroupList.
func (r *LazyRegexp) FindAllGroups(s string) *GroupList {
	results := NewGroupList()
	if r == nil {
		return results
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return results
	}

	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return results
	}

	names := re.SubexpNames()
	for _, match := range matches {
		groupMap := NewGroupMap()
		for i, name := range names {
			if name == "" || i >= len(match) {
				continue
			}
			groupMap.Set(name, match[i])
		}
		results.Add(groupMap)
	}

	return results
}

// FindString returns the leftmost match in s of the regular expression.
func (r *LazyRegexp) FindString(s string) string {
	if r == nil {
		return ""
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return ""
	}

	return re.FindString(s)
}

// FindStringSubmatch returns a slice of strings holding the leftmost submatches in s.
func (r *LazyRegexp) FindStringSubmatch(s string) []string {
	if r == nil {
		return nil
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return nil
	}

	return re.FindStringSubmatch(s)
}

// FindAllString returns a slice of all successive matches of the expression.
func (r *LazyRegexp) FindAllString(s string, n int) []string {
	if r == nil {
		return nil
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return nil
	}

	return re.FindAllString(s, n)
}

// ReplaceAllString returns a copy of src with all matches replaced by repl.
func (r *LazyRegexp) ReplaceAllString(src, repl string) string {
	if r == nil {
		return src
	}

	re, err := r.Compile()
	if err != nil || re == nil {
		return src
	}

	return re.ReplaceAllString(src, repl)
}

// CacheLen returns the number of uniquely cached regular expression instances in the global registry.
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
