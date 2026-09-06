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
	locker       sync.Mutex
	isCompiled   bool
	isApplicable bool
	expression   string
	regex        *regexp.Regexp
	compiledErr  error
	compiler     func(pattern string) (*regexp.Regexp, error)
}

// IsNull returns true if the receiver pointer is nil.
func (it *LazyRegex) IsNull() bool {
	return it == nil
}

// IsDefined returns true if the receiver is non-nil with expression set.
func (it *LazyRegex) IsDefined() bool {
	return it != nil && it.expression != ""
}

// IsUndefined returns true if the receiver is nil or missing expression.
func (it *LazyRegex) IsUndefined() bool {
	return it == nil || it.expression == ""
}

// IsApplicable compiles the regex if needed and returns true if compilation succeeded.
func (it *LazyRegex) IsApplicable() bool {
	if it == nil {
		return false
	}

	it.locker.Lock()
	if it.isApplicable {
		it.locker.Unlock()
		return true
	}
	it.locker.Unlock()

	if it.IsUndefined() {
		return false
	}

	_ = it.Compile()

	it.locker.Lock()
	defer it.locker.Unlock()
	return it.isApplicable
}

// Compile compiles the regular expression using the assigned compiler function or standard regexp.Compile.
func (it *LazyRegex) Compile() appfault.Result[*regexp.Regexp] {
	if it == nil {
		return appfault.Fail[*regexp.Regexp](appfault.NewAppBuilder(errtype.Execution, "nil LazyRegex cannot compile").Build())
	}

	if it.isCompiled && it.regex != nil {
		return appfault.NewSuccess(it.regex)
	}

	it.locker.Lock()
	defer it.locker.Unlock()

	if it.isCompiled {
		if it.compiledErr != nil {
			return appfault.Fail[*regexp.Regexp](appfault.NewAppBuilder(errtype.Execution, "lazy regex compilation failed").SetCause(it.compiledErr).Build())
		}
		return appfault.NewSuccess(it.regex)
	}

	if it.expression == "" {
		return appfault.Fail[*regexp.Regexp](appfault.NewAppBuilder(errtype.Execution, "lazy regex has empty expression").Build())
	}

	var (
		compiledRegex *regexp.Regexp
		regExErr      error
	)
	if it.compiler != nil {
		compiledRegex, regExErr = it.compiler(it.expression)
	} else {
		compiledRegex, regExErr = regexp.Compile(it.expression)
	}

	it.isApplicable = compiledRegex != nil && regExErr == nil
	it.regex = compiledRegex
	it.compiledErr = regExErr
	it.isCompiled = true

	if regExErr != nil {
		builder := appfault.NewAppBuilder(errtype.Execution, "lazy regex compilation failed")
		builder.SetCause(regExErr)
		builder.SetContext("expression", it.expression)
		return appfault.Fail[*regexp.Regexp](builder.Build())
	}

	return appfault.NewSuccess(compiledRegex)
}

// CompileMust compiles the regular expression and panics on compilation error.
// It checks first if an existing compiled regexp already exists and returns it immediately.
func (it *LazyRegex) CompileMust() *regexp.Regexp {
	if it == nil {
		return nil
	}

	if it.isCompiled && it.regex != nil {
		return it.regex
	}

	res := it.Compile()
	if res.IsFailure() {
		res.HandleError()
	}

	return res.Value
}

// IsCompiled reports whether compilation has already occurred.
func (it *LazyRegex) IsCompiled() bool {
	if it == nil {
		return false
	}

	it.locker.Lock()
	defer it.locker.Unlock()

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

	err := it.Compile().Error()
	return err
}

// OnRequiredCompiledMust triggers compilation and panics on error.
func (it *LazyRegex) OnRequiredCompiledMust() {
	err := it.OnRequiredCompiled()
	if err != nil {
		if appErr, ok := err.(*appfault.AppError); ok {
			appErr.HandleError()
		}
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
		if appErr, ok := compiledErr.(*appfault.AppError); ok {
			appErr.HandleError()
		}
	}
}

// String returns the raw regular expression pattern.
func (it *LazyRegex) String() string {
	if it == nil {
		return ""
	}

	return it.expression
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

	return it.expression
}

// MatchError returns nil on successful match, or a descriptive validation error.
func (it *LazyRegex) MatchError(matchingPattern string) error {
	if it == nil {
		return errors.New("nil LazyRegex cannot match")
	}

	res := it.Compile()
	regEx := res.Value
	compiledErr := res.Error()
	if regEx != nil && regEx.MatchString(matchingPattern) {
		return nil
	}

	return regExMatchValidationError(
		it.expression,
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

	res := it.Compile()
	regEx := res.Value
	compiledErr := res.Error()
	if regEx != nil && matchFunc != nil && matchFunc(regEx, comparing) {
		return nil
	}

	return regExMatchValidationError(
		it.expression,
		comparing,
		compiledErr,
		regEx)
}

// IsMatch reports whether comparing matches the compiled regular expression.
func (it *LazyRegex) IsMatch(comparing string) bool {
	if it == nil {
		return false
	}

	res := it.Compile()
	regEx := res.Value
	compiledErr := res.Error()
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

	res := it.Compile()
	regEx := res.Value
	compiledErr := res.Error()
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

	res := it.Compile()
	regEx := res.Value
	compiledErr := res.Error()
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

// Regex returns the underlying compiled *regexp.Regexp.
func (it *LazyRegex) Regex() *regexp.Regexp {
	return it.CompileMust()
}

// Compiled returns the underlying compiled *regexp.Regexp.
func (it *LazyRegex) Compiled() *regexp.Regexp {
	return it.CompileMust()
}

func (it *LazyRegex) compiledRegex() (*regexp.Regexp, error) {
	if it == nil {
		return nil, errors.New("nil LazyRegex")
	}

	if it.isCompiled && it.regex != nil {
		return it.regex, nil
	}

		res := it.Compile()
	if res.AppError != nil {
		return res.Value, res.AppError
	}
	return res.Value, nil
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

	re, err := it.compiledRegex()
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

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return 0
	}

	matches := re.FindAllString(comparing, -1)
	return len(matches)
}

// GroupBy extracts named capture groups (?P<name>...) from the first match into a GroupMap.
func (it *LazyRegex) GroupBy(comparing string) GroupMap {
	result := NewGroupMap()
	if it == nil {
		return result
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return result
	}

	match := re.FindStringSubmatch(comparing)
	if len(match) == 0 {
		return result
	}

	names := re.SubexpNames()
	for i, name := range names {
		if name == "" || i >= len(match) {
			continue
		}
		result[name] = match[i]
	}

	return result
}

// FindGroups is an alias for GroupBy.
func (it *LazyRegex) FindGroups(comparing string) GroupMap {
	return it.GroupBy(comparing)
}

// FindAllGroups extracts named capture groups across all non-overlapping matches in comparing into a GroupList.
func (it *LazyRegex) FindAllGroups(comparing string) GroupList {
	results := NewGroupList()
	if it == nil {
		return results
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return results
	}

	matches := re.FindAllStringSubmatch(comparing, -1)
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

// CompileBuilder compiles the regex, returning a structured Result wrapper.
// It checks first if an existing compiled regex already exists.
func (it *LazyRegex) CompileBuilder() appfault.Result[*regexp.Regexp] {
	return it.Compile()
}

// CompileResult compiles the regex, returning a structured Result wrapper.
func (it *LazyRegex) CompileResult() appfault.Result[*regexp.Regexp] {
	return it.Compile()
}
