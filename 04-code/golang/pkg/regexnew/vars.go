package regexnew

import (
	"regexp"
	"sync"
)

var (
	regexMutex    = sync.Mutex{}
	lazyRegexLock = sync.Mutex{}
	regexMaps     = make(
		map[string]*regexp.Regexp,
		DefaultCapacity)
	lazyRegexOnceMap = lazyRegexMap{
		items: make(
			map[string]*LazyRegex,
			DefaultCapacity),
	}

	New = newCreator{}
)
