package regexnew

import (
	"regexp"
	"sync"
	"testing"
)

func TestLazyRegex_NilSafety(t *testing.T) {
	var nilRegex *LazyRegex

	if !nilRegex.IsNull() {
		t.Errorf("expected IsNull() to return true for nil")
	}

	if nilRegex.IsDefined() {
		t.Errorf("expected IsDefined() to return false for nil")
	}

	if !nilRegex.IsUndefined() {
		t.Errorf("expected IsUndefined() to return true for nil")
	}

	if nilRegex.IsApplicable() {
		t.Errorf("expected IsApplicable() to return false for nil")
	}

	if nilRegex.IsCompiled() {
		t.Errorf("expected IsCompiled() to return false for nil")
	}

	if !nilRegex.HasError() {
		t.Errorf("expected HasError() to return true for nil")
	}

	if !nilRegex.HasAnyIssues() {
		t.Errorf("expected HasAnyIssues() to return true for nil")
	}

	if !nilRegex.IsInvalid() {
		t.Errorf("expected IsInvalid() to return true for nil")
	}

	if nilRegex.Pattern() != "" {
		t.Errorf("expected Pattern() to return empty string for nil")
	}

	if nilRegex.String() != "" {
		t.Errorf("expected String() to return empty string for nil")
	}

	if nilRegex.FullString() != "" {
		t.Errorf("expected FullString() to return empty string for nil")
	}

	if nilRegex.IsMatch("sample") {
		t.Errorf("expected IsMatch() to return false for nil")
	}

	if nilRegex.IsMatchBytes([]byte("sample")) {
		t.Errorf("expected IsMatchBytes() to return false for nil")
	}

	if !nilRegex.IsFailedMatch("sample") {
		t.Errorf("expected IsFailedMatch() to return true for nil")
	}

	if !nilRegex.IsFailedMatchBytes([]byte("sample")) {
		t.Errorf("expected IsFailedMatchBytes() to return true for nil")
	}

	if nilRegex.FindString("sample") != "" {
		t.Errorf("expected FindString() to return empty string for nil")
	}

	if nilRegex.FindStringSubmatch("sample") != nil {
		t.Errorf("expected FindStringSubmatch() to return nil for nil")
	}

	if nilRegex.ReplaceAllString("sample", "x") != "sample" {
		t.Errorf("expected ReplaceAllString() to return src for nil")
	}

	firstMatch, isInvalid := nilRegex.FirstMatchLine("sample")
	if !isInvalid || firstMatch != "" {
		t.Errorf("expected FirstMatchLine() to report invalid for nil")
	}
}

func TestLazyRegex_LifecycleAndCompilation(t *testing.T) {
	lz := New.LazyLock(`^[a-z]+-\d+$`)

	if !lz.IsDefined() {
		t.Fatalf("expected IsDefined to be true")
	}

	if lz.Pattern() != `^[a-z]+-\d+$` {
		t.Errorf("unexpected pattern: %s", lz.Pattern())
	}

	compiled, err := lz.Compile()
	if err != nil {
		t.Fatalf("unexpected compilation error: %v", err)
	}

	if compiled == nil {
		t.Fatalf("expected non-nil compiled regex")
	}

	if !lz.IsCompiled() {
		t.Errorf("expected IsCompiled to be true")
	}

	if !lz.IsApplicable() {
		t.Errorf("expected IsApplicable to be true")
	}

	if lz.HasError() {
		t.Errorf("expected HasError to be false")
	}

	secondCompiled, err := lz.Compile()
	if err != nil {
		t.Fatalf("second compile error: %v", err)
	}

	if compiled != secondCompiled {
		t.Errorf("expected identical compiled regex pointer on subsequent calls")
	}
}

func TestLazyRegex_Matching(t *testing.T) {
	lz := New.LazyLock(`^(\w+):(\d+)$`)

	if !lz.IsMatch("port:8080") {
		t.Errorf("expected match for 'port:8080'")
	}

	if lz.IsMatch("invalid-string") {
		t.Errorf("did not expect match for 'invalid-string'")
	}

	if !lz.IsMatchBytes([]byte("port:8080")) {
		t.Errorf("expected byte match for 'port:8080'")
	}

	if !lz.IsFailedMatch("invalid-string") {
		t.Errorf("expected failed match for 'invalid-string'")
	}

	firstMatch, isInvalid := lz.FirstMatchLine("port:8080")
	if isInvalid || firstMatch != "port:8080" {
		t.Errorf("unexpected first match: %s, invalid: %v", firstMatch, isInvalid)
	}

	sub := lz.FindStringSubmatch("port:8080")
	if len(sub) != 3 || sub[1] != "port" || sub[2] != "8080" {
		t.Errorf("unexpected submatches: %v", sub)
	}

	repl := lz.ReplaceAllString("port:8080", "mapped")
	if repl != "mapped" {
		t.Errorf("unexpected replace: %s", repl)
	}

	err := lz.MatchError("port:8080")
	if err != nil {
		t.Errorf("expected nil error on valid match, got: %v", err)
	}

	errMismatch := lz.MatchError("invalid")
	if errMismatch == nil {
		t.Errorf("expected error on mismatch")
	}
}

func TestLazyRegex_InvalidPattern(t *testing.T) {
	lz := New.LazyLock(`[invalid`)

	if !lz.HasError() {
		t.Errorf("expected HasError to be true for invalid regex")
	}

	if lz.IsApplicable() {
		t.Errorf("expected IsApplicable to be false for invalid regex")
	}

	_, err := lz.Compile()
	if err == nil {
		t.Errorf("expected Compile to return error for invalid pattern")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected CompileMust to panic on invalid regex")
		}
	}()

	lz.CompileMust()
}

func TestNewCreator_GlobalDeduplication(t *testing.T) {
	pattern := `^global-dedup-(\d+)$`

	first := New.Lazy(pattern)
	second := New.Lazy(pattern)

	if first != second {
		t.Errorf("expected New.Lazy to return identical pointer for same pattern")
	}

	lockedFirst := New.LazyLock(pattern)
	if first != lockedFirst {
		t.Errorf("expected New.LazyLock to return identical pointer for existing pattern")
	}

	compiled1, _ := Create(pattern)
	compiled2, _ := CreateLock(pattern)

	if compiled1 != compiled2 {
		t.Errorf("expected Create and CreateLock to return identical pre-compiled regex")
	}
}

func TestNewCreator_Concurrency(t *testing.T) {
	pattern := `^concurrent-(\d+)$`
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	results := make([]*LazyRegex, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			lz := New.LazyLock(pattern)
			results[idx] = lz
			_ = lz.IsMatch("concurrent-42")
		}(i)
	}

	wg.Wait()

	base := results[0]
	for i := 1; i < numGoroutines; i++ {
		if results[i] != base {
			t.Fatalf("goroutine %d obtained different LazyRegex pointer", i)
		}
	}
}

func TestBatchCreators(t *testing.T) {
	p1 := `^p1$`
	p2 := `^p2$`

	first, second := New.LazyRegex.TwoLock(p1, p2)
	if first == nil || second == nil {
		t.Fatalf("expected non-nil instances from TwoLock")
	}

	if first.Pattern() != p1 || second.Pattern() != p2 {
		t.Errorf("TwoLock patterns mismatch")
	}

	many := New.LazyRegex.ManyUsingLock(`^m1$`, `^m2$`)
	if len(many) != 2 {
		t.Errorf("expected 2 items in ManyUsingLock")
	}

	all := New.LazyRegex.AllPatternsMap()
	if len(all) == 0 {
		t.Errorf("expected non-empty AllPatternsMap")
	}
}

func TestMatchHelpers(t *testing.T) {
	pattern := `^test-\d+$`

	if !IsMatchLock(pattern, "test-123") {
		t.Errorf("expected IsMatchLock to return true")
	}

	if !IsMatchFailed(pattern, "fail-xyz") {
		t.Errorf("expected IsMatchFailed to return true")
	}

	if MatchErrorLock(pattern, "test-123") != nil {
		t.Errorf("expected MatchErrorLock to return nil on match")
	}

	err := MatchUsingFuncErrorLock(pattern, "test-123", func(r *regexp.Regexp, s string) bool {
		return r.MatchString(s)
	})
	if err != nil {
		t.Errorf("unexpected error from MatchUsingFuncErrorLock: %v", err)
	}
}
