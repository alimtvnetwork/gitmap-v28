package cmdcommit

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestBackdateUsesRestrictionTimezone(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		t.Fatal(err)
	}
	friday := time.Date(2026, time.October, 30, 10, 31, 0, 0, loc)
	t.Cleanup(func() {
		os.Unsetenv("GIT_AUTHOR_DATE")
		os.Unsetenv("GIT_COMMITTER_DATE")
		isMorningStampApplied = false
	})
	if errStamp := applySimulatedMorningWindow(friday, loc); errStamp != nil {
		t.Fatal(errStamp)
	}
	stamp := os.Getenv("GIT_AUTHOR_DATE")
	if !strings.Contains(stamp, "T07:30:00+08:00") {
		t.Fatalf("stamp = %s", stamp)
	}
	if os.Getenv("GIT_COMMITTER_DATE") != stamp {
		t.Fatalf("committer = %s", os.Getenv("GIT_COMMITTER_DATE"))
	}
	if !isMorningStampApplied {
		t.Fatal("stamp flag was not set")
	}
}
