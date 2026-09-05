package regexnew

import "regexp"

type lazyRegexMap struct {
	items map[string]*LazyRegex
}

func (it *lazyRegexMap) IsEmpty() bool {
	return it == nil || len(it.items) == 0
}

func (it *lazyRegexMap) IsEmptyLock() bool {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	return it == nil || len(it.items) == 0
}

func (it *lazyRegexMap) HasAnyItem() bool {
	return it != nil && len(it.items) > 0
}

func (it *lazyRegexMap) HasAnyItemLock() bool {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	return it != nil && len(it.items) > 0
}

func (it *lazyRegexMap) Length() int {
	if it == nil {
		return 0
	}

	return len(it.items)
}

func (it *lazyRegexMap) LengthLock() int {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	return it.Length()
}

func (it *lazyRegexMap) Has(keyName string) bool {
	if it == nil || it.items == nil {
		return false
	}

	_, has := it.items[keyName]
	return has
}

func (it *lazyRegexMap) HasLock(keyName string) bool {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	return it.Has(keyName)
}

func (it *lazyRegexMap) CreateOrExisting(
	patternName string,
) (lazyRegex *LazyRegex, isExisting bool) {
	lazyRegEx, has := it.items[patternName]
	if has {
		return lazyRegEx, true
	}

	lazyRegex = it.createDefaultLazyRegex(patternName)
	it.items[patternName] = lazyRegex

	return lazyRegex, false
}

func (it *lazyRegexMap) CreateOrExistingLock(
	patternName string,
) (lazyRegex *LazyRegex, isExisting bool) {
	lazyRegexLock.Lock()
	defer lazyRegexLock.Unlock()

	return it.CreateOrExisting(patternName)
}

func (it *lazyRegexMap) CreateOrExistingLockIf(
	isLock bool,
	patternName string,
) (lazyRegex *LazyRegex, isExisting bool) {
	if isLock {
		lazyRegexLock.Lock()
		defer lazyRegexLock.Unlock()
	}

	return it.CreateOrExisting(patternName)
}

func (it *lazyRegexMap) createDefaultLazyRegex(
	patternName string,
) *LazyRegex {
	return &LazyRegex{
		pattern:  patternName,
		compiler: CreateLock,
	}
}

func (it *lazyRegexMap) createLazyRegex(
	patternName string,
	creatorFunc func(pattern string) (*regexp.Regexp, error),
) *LazyRegex {
	return &LazyRegex{
		pattern:  patternName,
		compiler: creatorFunc,
	}
}
