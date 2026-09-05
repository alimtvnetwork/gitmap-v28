package lazyregex

import (
	"regexp"
	"sync"
)

var (
	globalLock sync.Mutex
	globalMap  = make(map[string]*LazyRegexp, 64)
	regexMap   = make(map[string]*regexp.Regexp, 64)
)

// LazyRegexp is a lazily compiled regular expression with global deduplication.
type LazyRegexp struct {
	expr string
	re   *regexp.Regexp
	once sync.Once
}

// New creates or retrieves a globally cached LazyRegexp for the given expression.
func New(expr string) *LazyRegexp {
	globalLock.Lock()
	defer globalLock.Unlock()

	item, exists := globalMap[expr]
	if exists {
		return item
	}

	item = &LazyRegexp{expr: expr}
	globalMap[expr] = item

	return item
}

// NewLock creates or retrieves a cached LazyRegexp with mutex locking (alias for New).
func NewLock(expr string) *LazyRegexp {
	return New(expr)
}

// Re returns the underlying compiled *regexp.Regexp.
// It compiles the regular expression on the first call and caches the compiled instance globally.
func (r *LazyRegexp) Re() *regexp.Regexp {
	if r == nil {
		return nil
	}

	r.once.Do(func() {
		globalLock.Lock()
		compiled, exists := regexMap[r.expr]
		if exists {
			globalLock.Unlock()
			r.re = compiled
			return
		}
		globalLock.Unlock()

		compiled = regexp.MustCompile(r.expr)

		globalLock.Lock()
		regexMap[r.expr] = compiled
		globalLock.Unlock()

		r.re = compiled
	})

	return r.re
}

// MatchString reports whether the string s contains any match of the regular expression.
func (r *LazyRegexp) MatchString(s string) bool {
	if r == nil {
		return false
	}

	return r.Re().MatchString(s)
}

// FindString returns a string holding the text of the leftmost match in s of the regular expression.
func (r *LazyRegexp) FindString(s string) string {
	if r == nil {
		return ""
	}

	return r.Re().FindString(s)
}

// FindStringSubmatch returns a slice of strings holding the text of the leftmost match of the regular expression in s.
func (r *LazyRegexp) FindStringSubmatch(s string) []string {
	if r == nil {
		return nil
	}

	return r.Re().FindStringSubmatch(s)
}

// ReplaceAllString returns a copy of src, replacing matches of the Regexp with the replacement string repl.
func (r *LazyRegexp) ReplaceAllString(src, repl string) string {
	if r == nil {
		return src
	}

	return r.Re().ReplaceAllString(src, repl)
}

// CacheLen returns the number of uniquely cached regular expression patterns in the global registry.
func CacheLen() int {
	globalLock.Lock()
	defer globalLock.Unlock()

	return len(globalMap)
}

// ClearCache flushes the global pattern and compiled regex registries (primarily for testing).
func ClearCache() {
	globalLock.Lock()
	defer globalLock.Unlock()

	globalMap = make(map[string]*LazyRegexp, 64)
	regexMap = make(map[string]*regexp.Regexp, 64)
}
