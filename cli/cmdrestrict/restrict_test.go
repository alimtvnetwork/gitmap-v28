package cmdrestrict

import (
	"strings"
	"testing"
	"time"
)

func TestAddEditStatus(t *testing.T) {
	home := t.TempDir()
	restrictionsDirOverride = home
	t.Cleanup(func() { restrictionsDirOverride = "" })
	errAdd := Run([]string{"add", "--name", "morning", "--email", "local@local", "--tz", "Asia/Singapore", "--blocked-after", "08:00", "--weekend-open", "--permit-stay"})
	if errAdd != nil {
		t.Fatal(errAdd)
	}
	errEdit := Run([]string{"edit", "--name", "morning", "--email", "next@local"})
	if errEdit != nil {
		t.Fatal(errEdit)
	}
	file, errLoad := loadRestrictions()
	if errLoad != nil {
		t.Fatal(errLoad)
	}
	if file.Items[0].Email != "next@local" {
		t.Fatalf("email = %s", file.Items[0].Email)
	}
	if errStatus := printStatus(); errStatus != nil {
		t.Fatal(errStatus)
	}
}

func TestWindowAfterCutoffIsRestricted(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		t.Fatal(err)
	}
	friday := time.Date(2026, 10, 30, 10, 1, 0, 0, loc)
	item := defaultRestriction()
	label := windowLabel(friday, loc, item)
	if !strings.Contains(label, "restricted") {
		t.Fatalf("label = %s", label)
	}
}
