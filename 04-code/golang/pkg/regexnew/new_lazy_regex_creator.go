package regexnew

type newLazyRegexCreator struct{}

// New creates or retrieves a cached LazyRegex without locking.
func (it newLazyRegexCreator) New(
	pattern string,
) *LazyRegex {
	lazyRegex, _ := lazyRegexOnceMap.CreateOrExisting(pattern)
	return lazyRegex
}

// NewLock creates or retrieves a cached LazyRegex with mutex locking.
func (it newLazyRegexCreator) NewLock(
	pattern string,
) *LazyRegex {
	lazyRegex, _ := lazyRegexOnceMap.CreateOrExistingLock(pattern)
	return lazyRegex
}

// TwoLock creates or retrieves two cached LazyRegex instances under a single lock.
func (it newLazyRegexCreator) TwoLock(
	pattern, secondPattern string,
) (first, second *LazyRegex) {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	first = it.New(pattern)
	second = it.New(secondPattern)

	return first, second
}

// ManyUsingLock creates or retrieves multiple cached LazyRegex instances under a single lock.
func (it newLazyRegexCreator) ManyUsingLock(
	patterns ...string,
) (patternsKeyAsMap map[string]*LazyRegex) {
	if len(patterns) == 0 {
		return map[string]*LazyRegex{}
	}

	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	patternsKeyAsMap = make(
		map[string]*LazyRegex,
		len(patterns))

	for _, pattern := range patterns {
		patternsKeyAsMap[pattern] = it.New(pattern)
	}

	return patternsKeyAsMap
}

// AllPatternsMap returns the underlying pattern map under lock.
func (it newLazyRegexCreator) AllPatternsMap() map[string]*LazyRegex {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	copyMap := make(map[string]*LazyRegex, len(lazyRegexOnceMap.items))
	for k, v := range lazyRegexOnceMap.items {
		copyMap[k] = v
	}

	return copyMap
}

// NewLockIf conditionally applies a lock when creating or retrieving a LazyRegex.
func (it newLazyRegexCreator) NewLockIf(
	isLock bool,
	pattern string,
) *LazyRegex {
	if isLock {
		return it.NewLock(pattern)
	}

	return it.New(pattern)
}
