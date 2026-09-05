package regexnew

import (
	"sync"
)

var (
	regexMutex       = sync.Mutex{}
	lazyRegexLock    = sync.Mutex{}
	lazyRegexOnceMap = lazyRegexMap{
		items: make(
			map[string]*LazyRegex,
			DefaultCapacity),
	}

	New = newCreator{}
)
