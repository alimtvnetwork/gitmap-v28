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

	if it.isCompiled && it.regex != nil {
		return result.SuccessResult(it.regex)
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

	if it.isCompiled && it.regex != nil {
		return it.regex
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

	if it.isCompiled && it.regex != nil {
		return it.regex, nil
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
func (it *LazyRegexp) IsMatch(s string) bool {
	if it == nil {
		return false
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return false
	}

	return re.MatchString(s)
}

// IsFound is a semantic alias for IsMatch.
func (it *LazyRegexp) IsFound(s string) bool {
	return it.IsMatch(s)
}

// MatchString reports whether the string s contains any match of the regular expression.
func (it *LazyRegexp) MatchString(s string) bool {
	return it.IsMatch(s)
}

// Count returns the number of non-overlapping matches of the regular expression in s.
func (it *LazyRegexp) Count(s string) int {
	if it == nil {
		return 0
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return 0
	}

	matches := re.FindAllString(s, -1)

	return len(matches)
}

// GroupBy extracts named capture groups (?P<name>...) from the first match into a GroupMap.
func (it *LazyRegexp) GroupBy(s string) GroupMap {
	res := NewGroupMap()
	if it == nil {
		return res
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return res
	}

	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return res
	}

	names := re.SubexpNames()
	for i, name := range names {
		if name == "" || i >= len(match) {
			continue
		}

		res[name] = match[i]
	}

	return res
}

// FindGroups is an alias for GroupBy.
func (it *LazyRegexp) FindGroups(s string) GroupMap {
	return it.GroupBy(s)
}

// FindAllGroups extracts named capture groups across all non-overlapping matches in s into a GroupList.
func (it *LazyRegexp) FindAllGroups(s string) GroupList {
	results := NewGroupList()
	if it == nil {
		return results
	}

	re, err := it.compiledRegex()
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

			groupMap[name] = match[i]
		}

		results = append(results, groupMap)
	}

	return results
}

// FindString returns the leftmost match in s of the regular expression.
func (it *LazyRegexp) FindString(s string) string {
	if it == nil {
		return ""
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return ""
	}

	return re.FindString(s)
}

// FindStringSubmatch returns a slice of strings holding the leftmost submatches in s.
func (it *LazyRegexp) FindStringSubmatch(s string) []string {
	if it == nil {
		return nil
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return nil
	}

	return re.FindStringSubmatch(s)
}

// FindAllString returns a slice of all successive matches of the expression.
func (it *LazyRegexp) FindAllString(s string, n int) []string {
	if it == nil {
		return nil
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return nil
	}

	return re.FindAllString(s, n)
}

// ReplaceAllString returns a copy of src with all matches replaced by repl.
func (it *LazyRegexp) ReplaceAllString(src, repl string) string {
	if it == nil {
		return src
	}

	re, err := it.compiledRegex()
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
