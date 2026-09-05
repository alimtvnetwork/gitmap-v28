package lazyregex

import (
	"regexp"
	"sync"
	"testing"
)

func TestLazyRegexp_Compilation(t *testing.T) {
	lr := New("a(b)c")

	if lr.re != nil {
		t.Errorf("expected re to be nil before compilation")
	}

	re := lr.Re()
	if re == nil {
		t.Fatalf("expected re to be non-nil after Re()")
	}

	if re.String() != "a(b)c" {
		t.Errorf("expected string 'a(b)c', got '%s'", re.String())
	}

	// second call should return same instance
	re2 := lr.Re()
	if re != re2 {
		t.Errorf("expected same regexp instance on subsequent calls")
	}
}

func TestLazyRegexp_ThreadSafety(t *testing.T) {
	lr := New("foo|bar")
	var wg sync.WaitGroup

	numGoroutines := 100
	wg.Add(numGoroutines)

	instances := make([]*regexp.Regexp, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			instances[idx] = lr.Re()
		}(i)
	}

	wg.Wait()

	first := instances[0]
	if first == nil {
		t.Fatalf("expected non-nil regexp")
	}

	for i := 1; i < numGoroutines; i++ {
		if instances[i] != first {
			t.Errorf("expected instance %d to be same as first instance", i)
		}
	}
}

func TestLazyRegexp_Wrappers(t *testing.T) {
	lr := New("(?i)hello (world)")

	if !lr.MatchString("Hello World") {
		t.Errorf("expected MatchString to return true")
	}

	if lr.MatchString("Goodbye") {
		t.Errorf("expected MatchString to return false")
	}

	if str := lr.FindString("say Hello World!"); str != "Hello World" {
		t.Errorf("FindString unexpected result: %s", str)
	}

	sub := lr.FindStringSubmatch("say Hello World!")
	if len(sub) != 2 || sub[0] != "Hello World" || sub[1] != "World" {
		t.Errorf("FindStringSubmatch unexpected result: %v", sub)
	}

	repl := lr.ReplaceAllString("say Hello World!", "bye")
	if repl != "say bye!" {
		t.Errorf("ReplaceAllString unexpected result: %s", repl)
	}
}

func TestLazyRegexp_MustCompilePanic(t *testing.T) {
	lr := New("[invalid regex")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected Re() to panic with invalid regex")
		}
	}()

	lr.Re()
}

func TestLazyRegexp_GlobalMapCaching(t *testing.T) {
	ClearCache()
	defer ClearCache()

	expr := `^test-cache-pattern-\d+$`

	first := New(expr)
	second := New(expr)

	if first != second {
		t.Errorf("expected New to return identical instance for same pattern")
	}

	re1 := first.Re()
	re2 := second.Re()

	if re1 != re2 {
		t.Errorf("expected Re() to return identical compiled *regexp.Regexp instance")
	}

	if CacheLen() != 1 {
		t.Errorf("expected CacheLen to be 1, got %d", CacheLen())
	}
}

func TestLazyRegexp_CompiledFlagAndCount(t *testing.T) {
	ClearCache()
	defer ClearCache()

	lr := New(`\b\w+\b`)
	if lr.IsCompiled() {
		t.Errorf("expected IsCompiled to be false before compilation")
	}

	count := lr.Count("one two three four")
	if count != 4 {
		t.Errorf("expected Count to be 4, got %d", count)
	}

	if !lr.IsCompiled() {
		t.Errorf("expected IsCompiled to be true after Count()")
	}

	if !lr.IsFound("hello world") {
		t.Errorf("expected IsFound to be true")
	}

	if lr.IsFound("") {
		t.Errorf("expected IsFound on empty string to be false")
	}
}

func TestLazyRegexp_GroupBy(t *testing.T) {
	ClearCache()
	defer ClearCache()

	emailRegex := New(`(?P<user>[a-zA-Z0-9._%+-]+)@(?P<domain>[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)

	groups := emailRegex.GroupBy("Contact us at support@example.com for help")
	if groups["user"] != "support" {
		t.Errorf("expected user 'support', got '%s'", groups["user"])
	}
	if groups["domain"] != "example.com" {
		t.Errorf("expected domain 'example.com', got '%s'", groups["domain"])
	}

	allGroups := emailRegex.FindAllGroups("first@a.com and second@b.org")
	if len(allGroups) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(allGroups))
	}
	if allGroups[0]["user"] != "first" || allGroups[1]["user"] != "second" {
		t.Errorf("unexpected allGroups user values: %v", allGroups)
	}
}

func TestLazyRegexp_CompileAppError(t *testing.T) {
	ClearCache()
	defer ClearCache()

	invalid := New(`[unclosed bracket`)
	re, appErr := invalid.CompileAppError()

	if re != nil {
		t.Errorf("expected re to be nil for invalid pattern")
	}
	if appErr == nil {
		t.Fatalf("expected appErr to be non-nil for invalid pattern")
	}
	if appErr.Op != "lazyregex.Compile" {
		t.Errorf("expected Op lazyregex.Compile, got %s", appErr.Op)
	}
	if appErr.Cause == nil {
		t.Errorf("expected non-nil Cause on wrapped AppError")
	}
}
