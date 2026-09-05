package regexnew

import "regexp"

type newCreator struct {
	LazyRegex newLazyRegexCreator
}

// Lazy creates or retrieves a cached LazyRegex without locking.
// Recommended for package-level var initialization.
func (it newCreator) Lazy(pattern string) *LazyRegex {
	return it.LazyRegex.New(pattern)
}

// LazyLock creates or retrieves a cached LazyRegex with mutex locking.
// Recommended for runtime calls inside functions.
func (it newCreator) LazyLock(pattern string) *LazyRegex {
	return it.LazyRegex.NewLock(pattern)
}

// Default compiles or retrieves a pre-compiled *regexp.Regexp without locking.
func (it newCreator) Default(pattern string) (*regexp.Regexp, error) {
	return Create(pattern)
}

// DefaultLock compiles or retrieves a pre-compiled *regexp.Regexp with mutex locking.
func (it newCreator) DefaultLock(pattern string) (*regexp.Regexp, error) {
	return CreateLock(pattern)
}

// DefaultLockIf compiles or retrieves a pre-compiled *regexp.Regexp with conditional locking.
func (it newCreator) DefaultLockIf(
	isLock bool,
	pattern string,
) (*regexp.Regexp, error) {
	return CreateLockIf(isLock, pattern)
}

// DefaultApplicableLock compiles under lock and returns regex, error, and whether applicable.
func (it newCreator) DefaultApplicableLock(pattern string) (
	regEx *regexp.Regexp,
	err error,
	isApplicable bool,
) {
	return CreateApplicableLock(pattern)
}
