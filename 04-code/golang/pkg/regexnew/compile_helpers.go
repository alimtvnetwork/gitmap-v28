package regexnew

import "regexp"

// Create creates regex if not already present in the global regexMaps cache.
// If compilation returns an error, it is returned and not stored in the cache.
func Create(regularExpressionPattern string) (*regexp.Regexp, error) {
	regex, has := regexMaps[regularExpressionPattern]
	if has {
		return regex, nil
	}

	newRegex, err := regexp.Compile(regularExpressionPattern)
	if err == nil {
		regexMaps[regularExpressionPattern] = newRegex
	}

	return newRegex, err
}

// CreateLock calls Create protected by regexMutex.
func CreateLock(regularExpressionPattern string) (*regexp.Regexp, error) {
	regexMutex.Lock()
	defer regexMutex.Unlock()

	return Create(regularExpressionPattern)
}

// CreateLockIf calls Create with mutex locking if isLock is true.
func CreateLockIf(isLock bool, regularExpressionSyntax string) (*regexp.Regexp, error) {
	if isLock {
		regexMutex.Lock()
		defer regexMutex.Unlock()
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

	regex, err := Create(regularExpressionPattern)
	applicable := err == nil && regex != nil

	return regex, err, applicable
}

// CreateMust compiles the regex or panics if invalid, caching on success.
func CreateMust(regularExpressionSyntax string) *regexp.Regexp {
	regex, has := regexMaps[regularExpressionSyntax]
	if has {
		return regex
	}

	newRegex := regexp.MustCompile(regularExpressionSyntax)
	regexMaps[regularExpressionSyntax] = newRegex

	return newRegex
}

// CreateMustLockIf calls CreateMust with conditional mutex locking.
func CreateMustLockIf(isLock bool, regularExpressionSyntax string) *regexp.Regexp {
	if isLock {
		regexMutex.Lock()
		defer regexMutex.Unlock()
	}

	return CreateMust(regularExpressionSyntax)
}

// NewMustLock compiles or retrieves a cached regex under regexMutex, panicking on error.
func NewMustLock(regularExpressionSyntax string) *regexp.Regexp {
	regexMutex.Lock()
	defer regexMutex.Unlock()

	return CreateMust(regularExpressionSyntax)
}
