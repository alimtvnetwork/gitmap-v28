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
