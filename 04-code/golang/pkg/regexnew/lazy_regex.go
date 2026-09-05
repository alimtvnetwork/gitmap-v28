package regexnew

import (
	"errors"
	"regexp"
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// LazyRegex provides a lazy-compiled, thread-safe regular expression wrapper.
// Patterns are compiled at most once and cached globally.
type LazyRegex struct {
	mu           sync.Mutex
	isCompiled   bool
	isApplicable bool
	pattern      string
	regex        *regexp.Regexp
	compiledErr  error
	compiler     func(pattern string) (*regexp.Regexp, error)
}

// IsNull returns true if the receiver pointer is nil.
func (it *LazyRegex) IsNull() bool {
	return it == nil
}

// IsDefined returns true if the receiver is non-nil with pattern set.
func (it *LazyRegex) IsDefined() bool {
	return it != nil && it.pattern != ""
}

// IsUndefined returns true if the receiver is nil or missing pattern.
func (it *LazyRegex) IsUndefined() bool {
	return it == nil || it.pattern == ""
}

// IsApplicable compiles the regex if needed and returns true if compilation succeeded.
func (it *LazyRegex) IsApplicable() bool {
	if it == nil {
		return false
	}

	it.mu.Lock()
	if it.isApplicable {
		it.mu.Unlock()
		return true
	}
	it.mu.Unlock()

	if it.IsUndefined() {
		return false
	}

	_, _ = it.Compile()

	it.mu.Lock()
	defer it.mu.Unlock()
	return it.isApplicable
}

// Compile compiles the regular expression using the assigned compiler function or standard regexp.Compile.
func (it *LazyRegex) Compile() (*regexp.Regexp, error) {
	if it == nil {
		return nil, errors.New("nil LazyRegex cannot compile")
	}

	it.mu.Lock()
	defer it.mu.Unlock()

	if it.isCompiled {
		return it.regex, it.compiledErr
	}

	if it.pattern == "" {
		return nil, errors.New("lazy regex has empty pattern")
	}

	var (
		compiledRegex *regexp.Regexp
		regExErr      error
	)
	if it.compiler != nil {
		compiledRegex, regExErr = it.compiler(it.pattern)
	} else {
		compiledRegex, regExErr = regexp.Compile(it.pattern)
	}

	it.isApplicable = compiledRegex != nil && regExErr == nil
	it.regex = compiledRegex
	it.compiledErr = regExErr
	it.isCompiled = true

	return compiledRegex, regExErr
}

// CompileMust compiles the regular expression and panics on compilation error.
func (it *LazyRegex) CompileMust() *regexp.Regexp {
	regexCompiled, err := it.Compile()
	if err != nil {
		panic(err)
	}

	return regexCompiled
}

// IsCompiled reports whether compilation has already occurred.
func (it *LazyRegex) IsCompiled() bool {
	if it == nil {
		return false
	}

	it.mu.Lock()
	defer it.mu.Unlock()

	return it.isCompiled
}

// OnRequiredCompiled triggers compilation if not already compiled, returning any error.
func (it *LazyRegex) OnRequiredCompiled() error {
	if it == nil {
		return errors.New("nil LazyRegex cannot compile")
	}

	if it.IsCompiled() {
		return it.compiledErr
	}

	_, err := it.Compile()
	return err
}

// OnRequiredCompiledMust triggers compilation and panics on error.
func (it *LazyRegex) OnRequiredCompiledMust() {
	err := it.OnRequiredCompiled()
	if err != nil {
		panic(err)
	}
}

// HasError reports whether compilation produced an error.
func (it *LazyRegex) HasError() bool {
	if it == nil {
		return true
	}

	_ = it.OnRequiredCompiled()
	return it.compiledErr != nil
}

// HasAnyIssues reports whether the receiver is nil, undefined, or failed compilation.
func (it *LazyRegex) HasAnyIssues() bool {
	if it == nil {
		return true
	}

	return !it.IsApplicable()
}

// IsInvalid reports whether the receiver failed compilation or is nil.
func (it *LazyRegex) IsInvalid() bool {
	if it == nil {
		return true
	}

	return !it.IsApplicable()
}

// CompiledError returns the error produced during compilation.
func (it *LazyRegex) CompiledError() error {
	return it.OnRequiredCompiled()
}

// Error returns the error produced during compilation (implements error inspector).
func (it *LazyRegex) Error() error {
	return it.OnRequiredCompiled()
}

// MustBeSafe panics if compilation encountered an error.
func (it *LazyRegex) MustBeSafe() {
	compiledErr := it.CompiledError()
	if compiledErr != nil {
		panic(compiledErr)
	}
}

// String returns the raw regular expression pattern.
func (it *LazyRegex) String() string {
	if it == nil {
		return ""
	}

	return it.pattern
}

// FullString returns a formatted JSON representation of the LazyRegex state.
func (it *LazyRegex) FullString() string {
	if it == nil {
		return ""
	}

	isApplicable := it.IsApplicable()
	isCompiled := it.IsCompiled()
	compiledErr := it.CompiledError()

	var errVal any
	if compiledErr != nil {
		errVal = compiledErr.Error()
	}

	stateMap := map[string]any{
		"pattern":      it.Pattern(),
		"isCompiled":   isCompiled,
		"isApplicable": isApplicable,
		"error":        errVal,
	}

	return prettyJson(stateMap)
}

// Pattern returns the raw regular expression pattern.
func (it *LazyRegex) Pattern() string {
	if it == nil {
		return ""
	}

	return it.pattern
}

// MatchError returns nil on successful match, or a descriptive validation error.
func (it *LazyRegex) MatchError(matchingPattern string) error {
	if it == nil {
		return errors.New("nil LazyRegex cannot match")
	}

	regEx, compiledErr := it.Compile()
	if regEx != nil && regEx.MatchString(matchingPattern) {
		return nil
	}

	return regExMatchValidationError(
		it.pattern,
		matchingPattern,
		compiledErr,
		regEx)
}

// MatchUsingFuncError matches using a custom validation function.
func (it *LazyRegex) MatchUsingFuncError(
	comparing string,
	matchFunc RegexValidationFunc,
) error {
	if it == nil {
		return errors.New("nil LazyRegex cannot match")
	}

	regEx, compiledErr := it.Compile()
	if regEx != nil && matchFunc != nil && matchFunc(regEx, comparing) {
		return nil
	}

	return regExMatchValidationError(
		it.pattern,
		comparing,
		compiledErr,
		regEx)
}

// IsMatch reports whether comparing matches the compiled regular expression.
func (it *LazyRegex) IsMatch(comparing string) bool {
	if it == nil {
		return false
	}

	regEx, compiledErr := it.Compile()
	if regEx == nil || compiledErr != nil {
		return false
	}

	return regEx.MatchString(comparing)
}

// IsMatchBytes reports whether comparingBytes matches the regular expression.
func (it *LazyRegex) IsMatchBytes(comparingBytes []byte) bool {
	if it == nil {
		return false
	}

	regEx, compiledErr := it.Compile()
	if regEx == nil || compiledErr != nil {
		return false
	}

	return regEx.Match(comparingBytes)
}

// IsFailedMatch returns true if the string does not match or compilation failed.
func (it *LazyRegex) IsFailedMatch(comparing string) bool {
	return !it.IsMatch(comparing)
}

// IsFailedMatchBytes returns true if the byte slice does not match or compilation failed.
func (it *LazyRegex) IsFailedMatchBytes(comparingBytes []byte) bool {
	return !it.IsMatchBytes(comparingBytes)
}

// FirstMatchLine returns the first submatch found in content.
func (it *LazyRegex) FirstMatchLine(
	content string,
) (firstMatch string, isInvalidMatch bool) {
	if it == nil {
		return "", true
	}

	regEx, compiledErr := it.Compile()
	if regEx == nil || compiledErr != nil {
		return "", true
	}

	lines := regEx.FindStringSubmatch(content)
	if len(lines) > 0 {
		return lines[0], false
	}

	return "", true
}

// Re returns the underlying compiled *regexp.Regexp, panicking on error.
func (it *LazyRegex) Re() *regexp.Regexp {
	return it.CompileMust()
}

// FindString returns the leftmost match in s.
func (it *LazyRegex) FindString(s string) string {
	if it == nil {
		return ""
	}

	re := it.CompileMust()
	return re.FindString(s)
}

// FindStringSubmatch returns slice of leftmost matches in s.
func (it *LazyRegex) FindStringSubmatch(s string) []string {
	if it == nil {
		return nil
	}

	re := it.CompileMust()
	return re.FindStringSubmatch(s)
}

// ReplaceAllString returns a copy of src with all matches replaced by repl.
func (it *LazyRegex) ReplaceAllString(src, repl string) string {
	if it == nil {
		return src
	}

	re := it.CompileMust()
	return re.ReplaceAllString(src, repl)
}

// FindAllString returns a slice of all successive matches of the expression.
func (it *LazyRegex) FindAllString(s string, n int) []string {
	if it == nil {
		return nil
	}

	re, err := it.Compile()
	if err != nil || re == nil {
		return nil
	}

	return re.FindAllString(s, n)
}

// IsFound is a semantic alias for IsMatch.
func (it *LazyRegex) IsFound(comparing string) bool {
	return it.IsMatch(comparing)
}

// Count returns the number of non-overlapping matches of the expression in comparing.
func (it *LazyRegex) Count(comparing string) int {
	if it == nil {
		return 0
	}

	re, err := it.Compile()
	if err != nil || re == nil {
		return 0
	}

	matches := re.FindAllString(comparing, -1)
	return len(matches)
}

// GroupBy extracts named capture groups (?P<name>...) from the first match into a map.
func (it *LazyRegex) GroupBy(comparing string) map[string]string {
	result := make(map[string]string)
	if it == nil {
		return result
	}

	re, err := it.Compile()
	if err != nil || re == nil {
		return result
	}

	match := re.FindStringSubmatch(comparing)
	if len(match) == 0 {
		return result
	}

	names := re.SubexpNames()
	for i, name := range names {
		if name != "" && i < len(match) {
			result[name] = match[i]
		}
	}

	return result
}

// FindAllGroups extracts named capture groups across all non-overlapping matches in comparing.
func (it *LazyRegex) FindAllGroups(comparing string) []map[string]string {
	if it == nil {
		return nil
	}

	re, err := it.Compile()
	if err != nil || re == nil {
		return nil
	}

	matches := re.FindAllStringSubmatch(comparing, -1)
	if len(matches) == 0 {
		return nil
	}

	names := re.SubexpNames()
	results := make([]map[string]string, 0, len(matches))

	for _, match := range matches {
		groupMap := make(map[string]string)
		for i, name := range names {
			if name != "" && i < len(match) {
				groupMap[name] = match[i]
			}
		}
		results = append(results, groupMap)
	}

	return results
}

// CompileBuilder compiles the regex, returning a structured AppBuilder on error.
func (it *LazyRegex) CompileBuilder() (*regexp.Regexp, *appfault.AppBuilder) {
	if it == nil {
		builder := appfault.NewAppBuilder(errtype.Execution, "nil LazyRegex cannot compile")
		return nil, builder
	}

	regEx, err := it.Compile()
	if err != nil {
		builder := appfault.NewAppBuilder(errtype.Execution, "lazy regex compilation failed")
		builder.SetCause(err)
		builder.SetContext("pattern", it.pattern)
		return nil, builder
	}

	return regEx, nil
}
