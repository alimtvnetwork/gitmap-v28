package timeutil

import (
	"testing"
	"time"
)

func TestResolveLocation(t *testing.T) {
	testValidLocations(t)
	testInvalidLocations(t)
}

func testValidLocations(t *testing.T) {
	locUTC, err := ResolveLocation("UTC")
	if err != nil || locUTC != time.UTC {
		t.Fatalf("expected UTC, got %v, err %v", locUTC, err)
	}

	locLocal, err := ResolveLocation("Local")
	if err != nil || locLocal != time.Local {
		t.Fatalf("expected Local, got %v, err %v", locLocal, err)
	}

	locEmpty, err := ResolveLocation("")
	if err != nil || locEmpty != time.Local {
		t.Fatalf("expected empty to resolve Local, got %v, err %v", locEmpty, err)
	}

	testOffsetsAndIANA(t)
}

func testOffsetsAndIANA(t *testing.T) {
	locPlus8, err := ResolveLocation("+08:00")
	if err != nil {
		t.Fatalf("failed to resolve +08:00: %v", err)
	}
	_, offset := time.Now().In(locPlus8).Zone()
	if offset != 8*3600 {
		t.Fatalf("expected offset 28800, got %d", offset)
	}

	locMinus4, err := ResolveLocation("-0400")
	if err != nil {
		t.Fatalf("failed to resolve -0400: %v", err)
	}
	_, offsetMinus := time.Now().In(locMinus4).Zone()
	if offsetMinus != -4*3600 {
		t.Fatalf("expected offset -14400, got %d", offsetMinus)
	}

	locNY, err := ResolveLocation("America/New_York")
	if err != nil || locNY == nil {
		t.Fatalf("failed to resolve America/New_York: %v", err)
	}
}

func testInvalidLocations(t *testing.T) {
	if _, err := ResolveLocation("NonExistent/Zone"); err == nil {
		t.Fatal("expected error for non-existent zone")
	}

	if _, err := ResolveLocation("+25:00"); err == nil {
		t.Fatal("expected error for invalid offset +25:00")
	}
}

func TestScheduleCommitWindows(t *testing.T) {
	loc := time.UTC
	wedMorning := time.Date(2026, 10, 28, 7, 30, 0, 0, loc)
	wedAfternoon := time.Date(2026, 10, 28, 14, 0, 0, 0, loc)
	saturdayAfternoon := time.Date(2026, 10, 31, 14, 0, 0, 0, loc)

	if !IsCommitWindowAllowed(wedMorning, loc) {
		t.Error("expected Wednesday 07:30 to be allowed")
	}
	if IsCommitWindowAllowed(wedAfternoon, loc) {
		t.Error("expected Wednesday 14:00 to be disallowed")
	}
	if !IsCommitWindowAllowed(saturdayAfternoon, loc) {
		t.Error("expected Saturday 14:00 to be allowed")
	}
}

func TestCalculateNextWindow(t *testing.T) {
	loc := time.UTC
	wedAfternoon := time.Date(2026, 10, 28, 14, 0, 0, 0, loc)
	nextFromWed := CalculateNextWindow(wedAfternoon, loc)
	expectedWedNext := time.Date(2026, 10, 29, 7, 30, 0, 0, loc)
	if !nextFromWed.Equal(expectedWedNext) {
		t.Errorf("expected %v, got %v", expectedWedNext, nextFromWed)
	}

	friAfternoon := time.Date(2026, 10, 30, 14, 0, 0, 0, loc)
	nextFromFri := CalculateNextWindow(friAfternoon, loc)
	expectedFriNext := time.Date(2026, 10, 31, 0, 0, 0, 0, loc)
	if !nextFromFri.Equal(expectedFriNext) {
		t.Errorf("expected %v, got %v", expectedFriNext, nextFromFri)
	}
}

func TestSimulateMorningWindow(t *testing.T) {
	loc := time.UTC
	wedAfternoon := time.Date(2026, 10, 28, 14, 0, 0, 0, loc)
	simulated := SimulateMorningWindow(wedAfternoon, loc)
	expected := time.Date(2026, 10, 28, 7, 30, 0, 0, loc)

	if !simulated.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, simulated)
	}
	if !IsCommitWindowAllowed(simulated, loc) {
		t.Error("simulated morning timestamp must be allowed")
	}
}
