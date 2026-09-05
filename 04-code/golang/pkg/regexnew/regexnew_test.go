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

func TestLazyRegex_CountAndGroups(t *testing.T) {
	lr := New.LazyLock(`(?P<word>\w+)`)

	count := lr.Count("alpha beta gamma delta")
	if count != 4 {
		t.Errorf("expected Count to be 4, got %d", count)
	}

	if !lr.IsFound("sample") {
		t.Errorf("expected IsFound to be true")
	}

	group := lr.GroupBy("alpha beta")
	if !group.Has("word") {
		t.Fatalf("expected group to have key 'word'")
	}
	if group.Get("word") != "alpha" {
		t.Errorf("expected group word 'alpha', got '%s'", group.Get("word"))
	}
	if group.GetOrDefault("missing", "fallback") != "fallback" {
		t.Errorf("expected fallback value for missing key")
	}

	all := lr.FindAllGroups("alpha beta gamma")
	if all.Len() != 3 {
		t.Fatalf("expected 3 group maps, got %d", all.Len())
	}
	if all.First().Get("word") != "alpha" {
		t.Errorf("expected first group word 'alpha', got '%s'", all.First().Get("word"))
	}
	if all.At(1).Get("word") != "beta" {
		t.Errorf("expected second group word 'beta', got '%s'", all.At(1).Get("word"))
	}
	if all.Last().Get("word") != "gamma" {
		t.Errorf("expected last group word 'gamma', got '%s'", all.Last().Get("word"))
	}

	words := all.ValuesOf("word")
	if len(words) != 3 || words[0] != "alpha" || words[1] != "beta" || words[2] != "gamma" {
		t.Errorf("unexpected ValuesOf result: %v", words)
	}

	keys := all.AllKeys()
	if len(keys) != 1 || keys[0] != "word" {
		t.Errorf("unexpected AllKeys: %v", keys)
	}
}

func TestLazyRegex_CompileBuilder(t *testing.T) {
	invalid := New.LazyLock(`[invalid regex`)

	res := invalid.CompileBuilder()
	if res.IsSuccess() {
		t.Errorf("expected compilation to fail on invalid pattern")
	}
	if !res.IsFailed() {
		t.Errorf("expected IsFailed to be true")
	}
	if !res.HasError() {
		t.Errorf("expected HasError to be true")
	}
	if res.Regexp() != nil {
		t.Errorf("expected nil Regexp on failure")
	}
	if res.Builder() == nil {
		t.Fatalf("expected non-nil AppBuilder on invalid pattern")
	}
	if res.AppError() == nil {
		t.Fatalf("expected non-nil AppError on invalid pattern")
	}
	if res.AppError().Message() != "lazy regex compilation failed" {
		t.Errorf("unexpected message: %s", res.AppError().Message())
	}

	valid := New.LazyLock(`^\w+$`)
	validRes := valid.CompileBuilder()
	if !validRes.IsSuccess() {
		t.Errorf("expected valid regex to succeed")
	}
	if validRes.HasError() {
		t.Errorf("expected no error on valid regex")
	}
	if validRes.Regexp() == nil {
		t.Errorf("expected non-nil Regexp on success")
	}
}

func TestGroupMap_Operations(t *testing.T) {
	gm := NewGroupMap()
	if !gm.IsEmpty() || gm.HasItems() {
		t.Errorf("expected new GroupMap to be empty")
	}

	gm.Set("key1", "val1").Add("key2", "val2")
	if gm.Len() != 2 {
		t.Errorf("expected len 2, got %d", gm.Len())
	}
	if !gm.Has("key1") || !gm.HasKey("key2") {
		t.Errorf("expected keys to exist")
	}
	if gm.Get("key1") != "val1" {
		t.Errorf("expected val1, got %s", gm.Get("key1"))
	}

	clone := gm.Clone()
	clone.Remove("key1")
	if !gm.Has("key1") {
		t.Errorf("original map should retain key1 after clone removal")
	}
	if clone.Has("key1") {
		t.Errorf("clone should not have key1")
	}

	raw := gm.ToMap()
	if raw["key1"] != "val1" || raw["key2"] != "val2" {
		t.Errorf("unexpected raw map: %v", raw)
	}

	var nilMap *GroupMap
	if nilMap.Has("foo") || nilMap.Get("foo") != "" || nilMap.Len() != 0 || !nilMap.IsEmpty() {
		t.Errorf("nil GroupMap should be safe")
	}
}

func TestGroupList_Operations(t *testing.T) {
	gl := NewGroupList()
	if !gl.IsEmpty() || gl.HasItems() {
		t.Errorf("expected empty GroupList")
	}

	g1 := NewGroupMap().Set("name", "alice").Set("role", "admin")
	g2 := NewGroupMap().Set("name", "bob").Set("role", "user")
	gl.Add(g1).Add(g2)

	if gl.Len() != 2 {
		t.Errorf("expected 2 items, got %d", gl.Len())
	}

	found := gl.Find(func(g *GroupMap) bool {
		return g.Get("role") == "admin"
	})
	if found.Get("name") != "alice" {
		t.Errorf("expected alice for admin, got %s", found.Get("name"))
	}

	filtered := gl.Filter(func(g *GroupMap) bool {
		return g.Get("role") == "user"
	})
	if filtered.Len() != 1 || filtered.First().Get("name") != "bob" {
		t.Errorf("unexpected filtered list")
	}

	outOfBounds := gl.At(999)
	if outOfBounds == nil || outOfBounds.HasItems() {
		t.Errorf("out of bounds At() should return empty non-nil GroupMap")
	}

	var nilList *GroupList
	if nilList.Len() != 0 || !nilList.IsEmpty() || nilList.First() == nil {
		t.Errorf("nil GroupList should be safe")
	}
}
