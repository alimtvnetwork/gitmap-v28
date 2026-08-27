package lazyregex

import (
	"regexp"
	"sync"
)

// LazyRegexp is a lazily compiled regular expression.
type LazyRegexp struct {
	expr string
	re   *regexp.Regexp
	once sync.Once
}

// New creates a new LazyRegexp that will compile expr on first use.
func New(expr string) *LazyRegexp {
	return &LazyRegexp{expr: expr}
}

// Re returns the underlying compiled *regexp.Regexp.
// It compiles the regular expression on the first call.
func (r *LazyRegexp) Re() *regexp.Regexp {
	r.once.Do(func() {
		r.re = regexp.MustCompile(r.expr)
	})
	return r.re
}

// MatchString reports whether the string s contains any match of the regular expression.
func (r *LazyRegexp) MatchString(s string) bool {
	return r.Re().MatchString(s)
}

// FindString returns a string holding the text of the leftmost match in s of the regular expression.
func (r *LazyRegexp) FindString(s string) string {
	return r.Re().FindString(s)
}

// FindStringSubmatch returns a slice of strings holding the text of the leftmost match of the regular expression in s.
func (r *LazyRegexp) FindStringSubmatch(s string) []string {
	return r.Re().FindStringSubmatch(s)
}

// ReplaceAllString returns a copy of src, replacing matches of the Regexp with the replacement string repl.
func (r *LazyRegexp) ReplaceAllString(src, repl string) string {
	return r.Re().ReplaceAllString(src, repl)
}
