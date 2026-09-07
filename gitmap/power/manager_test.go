package power

import (
	"strings"
	"testing"
)

func TestNewManager_Instantiation(t *testing.T) {
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if mgr == nil {
		t.Fatalf("Expected non-nil manager")
	}

	if mgr.Platform() == "" {
		t.Errorf("Expected non-empty platform")
	}
}

func TestSettings_ValidateAndSummary(t *testing.T) {
	s := Settings{
		Platform:              "windows",
		DisplayTimeoutMinutes: 10,
		SleepTimeoutMinutes:   30,
		IsNeverSleep:          false,
	}

	if err := s.Validate(); err != nil {
		t.Errorf("Expected valid settings: %v", err)
	}

	summary := s.Summary()
	if !strings.Contains(summary, "10m") || !strings.Contains(summary, "30m") {
		t.Errorf("Unexpected summary: %s", summary)
	}

	invalid := Settings{DisplayTimeoutMinutes: -5}
	if err := invalid.Validate(); err == nil {
		t.Errorf("Expected error for negative display timeout")
	}
}

func TestParsePowercfgSettingIndex(t *testing.T) {
	sampleOutput := `
    Power Setting GUID: 3c0bc021-c8a8-4e07-a973-6b14cbcb2b7e  (Turn off display after)
      GUID Alias: VIDEOIDLE
      Current AC Power Setting Index: 0x0000012c
      Current DC Power Setting Index: 0x0000012c

    Power Setting GUID: 29f6c1db-86da-48c5-9fdb-f2b67b1f44da  (Sleep after)
      GUID Alias: STANDBYIDLE
      Current AC Power Setting Index: 0x00000000
`

	sec, err := ParsePowercfgSettingIndex(sampleOutput, "VIDEOIDLE")
	if err != nil {
		t.Fatalf("ParsePowercfgSettingIndex failed: %v", err)
	}
	if sec != 300 {
		t.Errorf("Expected 300 seconds, got %d", sec)
	}

	mins := ConvertSecondsToMinutes(sec)
	if mins != 5 {
		t.Errorf("Expected 5 minutes, got %d", mins)
	}

	secZero, err := ParsePowercfgSettingIndex(sampleOutput, "STANDBYIDLE")
	if err != nil {
		t.Fatalf("ParsePowercfgSettingIndex failed: %v", err)
	}
	if secZero != 0 {
		t.Errorf("Expected 0 seconds, got %d", secZero)
	}
	if ConvertSecondsToMinutes(secZero) != 0 {
		t.Errorf("Expected 0 minutes, got %d", ConvertSecondsToMinutes(secZero))
	}
}

func TestParseGnomeHelpers(t *testing.T) {
	if sec := ParseGnomeSeconds("uint32 600"); sec != 600 {
		t.Errorf("Expected 600 seconds, got %d", sec)
	}

	if mins := ParseGnomeTimeoutMinutes("uint32 600"); mins != 10 {
		t.Errorf("Expected 10 minutes, got %d", mins)
	}

	if !ParseGnomeBoolean("true\n") {
		t.Errorf("Expected true for parse boolean")
	}

	if ParseGnomeBoolean("false\n") {
		t.Errorf("Expected false for parse boolean")
	}
}

func TestParsePmsetHelpers(t *testing.T) {
	sample := " displaysleep 15\n sleep 30\n"
	if mins := ParsePmsetValue(sample, "displaysleep"); mins != 15 {
		t.Errorf("Expected 15 minutes, got %d", mins)
	}

	settings := ParsePmsetSettings(sample)
	if settings.DisplayTimeoutMinutes != 15 || settings.SleepTimeoutMinutes != 30 {
		t.Errorf("Unexpected pmset settings: %+v", settings)
	}

	neverSettings := ParsePmsetSettings("displaysleep 0\nsleep 0\n")
	if !neverSettings.IsNeverSleep {
		t.Errorf("Expected IsNeverSleep to be true for 0/0")
	}
}

func TestMockRunner_DriverInteraction(t *testing.T) {
	var executedCmds []string
	restore := SetRunnerForTesting(func(name string, args ...string) ([]byte, error) {
		cmdLine := name + " " + strings.Join(args, " ")
		executedCmds = append(executedCmds, cmdLine)

		if strings.Contains(cmdLine, "SUB_VIDEO") {
			return []byte("GUID Alias: VIDEOIDLE\nCurrent AC Power Setting Index: 0x00000000\n"), nil
		}
		if strings.Contains(cmdLine, "SUB_SLEEP") {
			return []byte("GUID Alias: STANDBYIDLE\nCurrent AC Power Setting Index: 0x00000000\n"), nil
		}
		if strings.Contains(cmdLine, "pmset") {
			return []byte("displaysleep 0\nsleep 0\n"), nil
		}
		if strings.Contains(cmdLine, "gsettings") {
			return []byte("uint32 0\n"), nil
		}

		return []byte("OK"), nil
	})
	defer restore()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	status, err := mgr.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if !status.IsNeverSleep {
		t.Errorf("Expected IsNeverSleep to be true")
	}

	if err := mgr.SetNeverSleep(); err != nil {
		t.Fatalf("SetNeverSleep failed: %v", err)
	}

	if len(executedCmds) == 0 {
		t.Errorf("Expected commands to be recorded via mock runner")
	}
}
