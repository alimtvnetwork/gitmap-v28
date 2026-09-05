package regexnew

import "regexp"

// Create compiles or retrieves a cached regex using the single global pattern cache.
func Create(regularExpressionPattern string) (*regexp.Regexp, error) {
	lz := New.Lazy(regularExpressionPattern)
	return lz.Compile()
}

// CreateLock calls Create protected by regexMutex.
func CreateLock(regularExpressionPattern string) (*regexp.Regexp, error) {
	lz := New.LazyLock(regularExpressionPattern)
	return lz.Compile()
}

// CreateLockIf calls Create with mutex locking if isLock is true.
func CreateLockIf(isLock bool, regularExpressionSyntax string) (*regexp.Regexp, error) {
	if isLock {
		return CreateLock(regularExpressionSyntax)
	}

	return Create(regularExpressionSyntax)
}

// CreateApplicableLock compiles under lock and returns applicability boolean.
func CreateApplicableLock(regularExpressionPattern string) (
	regEx *regexp.Regexp,
	err error,
	isApplicable bool,
) {
	regexMutex.Lock()
	defer regexMutex.Unlock()

	lz := New.LazyLock(regularExpressionPattern)
	regEx, err = lz.Compile()
	isApplicable = err == nil && regEx != nil

	return regEx, err, isApplicable
}

// CreateMust compiles the regex or panics if invalid, caching on success.
func CreateMust(regularExpressionSyntax string) *regexp.Regexp {
	lz := New.Lazy(regularExpressionSyntax)
	return lz.CompileMust()
}

// CreateMustLockIf calls CreateMust with conditional mutex locking.
func CreateMustLockIf(isLock bool, regularExpressionSyntax string) *regexp.Regexp {
	if isLock {
		return NewMustLock(regularExpressionSyntax)
	}

	return CreateMust(regularExpressionSyntax)
}

// NewMustLock compiles or retrieves a cached regex under regexMutex, panicking on error.
func NewMustLock(regularExpressionSyntax string) *regexp.Regexp {
	lz := New.LazyLock(regularExpressionSyntax)
	return lz.CompileMust()
}
