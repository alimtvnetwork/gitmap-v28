package lazyregex

import (
	"strings"
	"testing"
)

func TestMatchResult_Success(t *testing.T) {
	re := New(`(?P<role>\w+):(?P<ip>\d+\.\d+\.\d+\.\d+)`)
	rs := re.MatchResult("master:192.168.1.10")

	if rs.IsFailed() || rs.IsFailure() {
		t.Fatal("expected match to be successful")
	}

	if !rs.IsMatch() || !rs.HasMatch() {
		t.Fatal("expected IsMatch and HasMatch to be true")
	}

	if rs.AsError() != nil || rs.ErrOrNil() != nil || rs.Cause() != nil {
		t.Errorf("expected nil error on success, got: %v", rs.AsError())
	}

	if rs.First() != "master:192.168.1.10" {
		t.Errorf("expected first match 'master:192.168.1.10', got: %q", rs.First())
	}

	if rs.At(1) != "master" || rs.At(2) != "192.168.1.10" {
		t.Errorf("expected submatches 'master' and '192.168.1.10', got: %q, %q", rs.At(1), rs.At(2))
	}

	if rs.Last() != "192.168.1.10" {
		t.Errorf("expected last match '192.168.1.10', got: %q", rs.Last())
	}

	if rs.Count() != 3 {
		t.Errorf("expected 3 submatches, got: %d", rs.Count())
	}

	named := rs.Map()
	if named.Get("role") != "master" || named.Get("ip") != "192.168.1.10" {
		t.Errorf("expected named groups role=master ip=192.168.1.10, got: %v", named)
	}

	if rs.Group().GetNamed("role") != "master" {
		t.Errorf("expected group GetNamed 'master', got: %q", rs.Group().GetNamed("role"))
	}
}

func TestMatchResult_Failure(t *testing.T) {
	re := New(`^gitmap\s+cluster\s+(nodes|ls)`)
	content := "gitmap unknown-cmd-xyz"
	rs := re.MatchResult(content)

	if !rs.IsFailed() || !rs.IsFailure() {
		t.Fatal("expected match to be failed")
	}

	if rs.IsSuccess() || rs.IsMatch() {
		t.Fatal("expected match to not be successful")
	}

	appErr := rs.AppError()
	if appErr == nil {
		t.Fatal("expected non-nil AppError on failure")
	}

	errStr := rs.ErrorString()
	if !strings.Contains(errStr, "regex pattern does not match content") {
		t.Errorf("expected diagnostic error message, got: %q", errStr)
	}

	if !strings.Contains(errStr, content) {
		t.Errorf("expected error to cite comparing content %q, got: %q", content, errStr)
	}

	if rs.FirstOrDefault("fallback") != "fallback" {
		t.Errorf("expected default 'fallback', got: %q", rs.FirstOrDefault("fallback"))
	}

	if len(rs.Items()) != 0 {
		t.Errorf("expected empty items on failure, got: %v", rs.Items())
	}
}

func TestMatchAllResults(t *testing.T) {
	re := New(`\b\w+@\w+\.\w+\b`)
	content := "contact user@test.com and admin@gitmap.org for details"
	all := re.MatchAllResults(content)

	if all.IsFailed() || all.IsFailure() {
		t.Fatal("expected match all to be successful")
	}

	if !all.IsMatch() {
		t.Fatal("expected IsMatch to be true")
	}

	if all.Count() != 2 {
		t.Fatalf("expected 2 matches, got: %d", all.Count())
	}

	items := all.Items()
	if len(items) != 2 || items[0] != "user@test.com" || items[1] != "admin@gitmap.org" {
		t.Errorf("expected matched items, got: %v", items)
	}

	if all.First().First() != "user@test.com" {
		t.Errorf("expected first match 'user@test.com', got: %q", all.First().First())
	}

	if all.Last().First() != "admin@gitmap.org" {
		t.Errorf("expected last match 'admin@gitmap.org', got: %q", all.Last().First())
	}
}

func TestMatchError_Helper(t *testing.T) {
	pattern := `^ok$`
	errPass := MatchError(pattern, "ok")
	if errPass != nil {
		t.Errorf("expected nil error on match, got: %v", errPass)
	}

	errFail := MatchError(pattern, "not-ok")
	if errFail == nil {
		t.Fatal("expected error on mismatch, got nil")
	}

	if !strings.Contains(errFail.Error(), pattern) {
		t.Errorf("expected error to contain pattern, got: %v", errFail)
	}
}

func TestMatchResult_ErrorFormattingInspection(t *testing.T) {
	re := New(`(?P<verb>add|join|nodes)\s+(?P<target>[\w@\.:]+)`)
	content := "gitmap cluster run invalid-param"
	rs := re.MatchResult(content)

	if rs.IsSuccess() {
		t.Fatal("expected match failure")
	}

	t.Logf("\n--- VISUAL ERROR REPORT SAMPLE ---\n%s\n---------------------------------", rs.ErrorString())
}

func TestNilSafety(t *testing.T) {

	var nilRes *MatchResult
	if !nilRes.IsFailed() || nilRes.IsSuccess() || nilRes.IsMatch() {
		t.Error("nil MatchResult must be failed and not success")
	}

	if nilRes.First() != "" || nilRes.Count() != 0 || len(nilRes.Items()) != 0 {
		t.Error("nil MatchResult must return safe empty values")
	}

	var nilGroup *MatchGroup
	if nilGroup.First() != "" || nilGroup.Count() != 0 || nilGroup.GetNamed("x") != "" {
		t.Error("nil MatchGroup must return safe empty values")
	}

	var nilAll *MatchAllResult
	if !nilAll.IsFailed() || nilAll.Count() != 0 || len(nilAll.Items()) != 0 {
		t.Error("nil MatchAllResult must return safe empty values")
	}
}
