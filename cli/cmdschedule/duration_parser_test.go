package cmdschedule

import (
	"testing"
	"time"
)

func assertDurationEquals(t *testing.T, raw string, expected time.Duration) {
	res := ParseScheduleDuration(raw)
	if res.IsFailure() {
		t.Fatalf("ParseScheduleDuration(%q) failed: %v", raw, res.AppError())
	}
	if res.Value != expected {
		t.Errorf("ParseScheduleDuration(%q) = %v, expected %v", raw, res.Value, expected)
	}
}

func assertDurationFails(t *testing.T, raw string) {
	res := ParseScheduleDuration(raw)
	if res.IsSuccess() {
		t.Fatalf("ParseScheduleDuration(%q) expected failure, got %v", raw, res.Value)
	}
	if isExpectedCode := res.IsErrorCode("E_INVALID_DURATION"); !isExpectedCode {
		t.Errorf("ParseScheduleDuration(%q) expected E_INVALID_DURATION, got %v", raw, res.AppError())
	}
}

func TestParseScheduleDurationColon(t *testing.T) {
	assertDurationEquals(t, "1:45hr", 6300*time.Second)
	assertDurationEquals(t, "1:45h", 6300*time.Second)
	assertDurationEquals(t, "1:45", 6300*time.Second)
	assertDurationEquals(t, "01:30:00", 5400*time.Second)
	assertDurationEquals(t, "1:30m", 90*time.Second)
}

func TestParseScheduleDurationDays(t *testing.T) {
	assertDurationEquals(t, "1day", 86400*time.Second)
	assertDurationEquals(t, "1d", 86400*time.Second)
	assertDurationEquals(t, "2days", 172800*time.Second)
	assertDurationEquals(t, "2d", 172800*time.Second)
}

func TestParseScheduleDurationHoursAndMinutes(t *testing.T) {
	assertDurationEquals(t, "2h", 7200*time.Second)
	assertDurationEquals(t, "2hr", 7200*time.Second)
	assertDurationEquals(t, "120m", 7200*time.Second)
	assertDurationEquals(t, "45min", 2700*time.Second)
	assertDurationEquals(t, "45mins", 2700*time.Second)
}

func TestParseScheduleDurationSecondsAndZero(t *testing.T) {
	assertDurationEquals(t, "1s", time.Second)
	assertDurationEquals(t, "2s", 2*time.Second)
	assertDurationEquals(t, "10sec", 10*time.Second)
	assertDurationEquals(t, "0", 0)
	assertDurationEquals(t, "now", 0)
	assertDurationEquals(t, "", 0)
	assertDurationEquals(t, "60", 60*time.Second)
}

func TestParseScheduleDurationInvalid(t *testing.T) {
	assertDurationFails(t, "invalid")
	assertDurationFails(t, "-5s")
	assertDurationFails(t, "foo:bar")
	assertDurationFails(t, "1:2:3:4")
}
