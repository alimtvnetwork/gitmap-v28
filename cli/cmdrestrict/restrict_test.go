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

func TestMachineModeAutoThenFlag(t *testing.T) {
	dir := t.TempDir()
	SetStoreDir(dir)
	t.Cleanup(func() { SetStoreDir("") })
	if !IsAutoEnforced() {
		t.Fatal("missing file must default to auto")
	}
	if err := SetMode(ModeFlag); err != nil {
		t.Fatal(err)
	}
	if IsAutoEnforced() {
		t.Fatal("flag mode must not auto-enforce")
	}
	if err := Run([]string{"mode", ModeAuto}); err != nil {
		t.Fatal(err)
	}
	if CurrentMode() != ModeAuto {
		t.Fatalf("mode = %s", CurrentMode())
	}
	attrs := map[string]string{settingsModeKey: ModeFlag}
	if err := ApplySettingsMode(attrs); err != nil {
		t.Fatal(err)
	}
	if CurrentMode() != ModeFlag {
		t.Fatalf("settings mode = %s", CurrentMode())
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
