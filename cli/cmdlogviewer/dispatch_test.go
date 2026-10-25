package cmdlogviewer

import (
	"testing"
)

func TestIsShowLogCommand(t *testing.T) {
	validCommands := []string{"show-log", "SHOW-LOG", "showlog", "log-ui"}
	for _, cmd := range validCommands {
		if !IsShowLogCommand(cmd) {
			t.Errorf("expected IsShowLogCommand(%q) to be true", cmd)
		}
	}

	invalidCommands := []string{"log", "show", "status", "compare", "ui"}
	for _, cmd := range invalidCommands {
		if IsShowLogCommand(cmd) {
			t.Errorf("expected IsShowLogCommand(%q) to be false", cmd)
		}
	}
}

func TestIsLogUIArgs(t *testing.T) {
	if !IsLogUIArgs([]string{"ui"}) {
		t.Errorf("expected IsLogUIArgs(['ui']) to be true")
	}

	if !IsLogUIArgs([]string{"UI", "some/path"}) {
		t.Errorf("expected IsLogUIArgs(['UI', 'path']) to be true")
	}

	if IsLogUIArgs([]string{"-n", "5"}) {
		t.Errorf("expected IsLogUIArgs(['-n', '5']) to be false")
	}

	if IsLogUIArgs([]string{}) {
		t.Errorf("expected IsLogUIArgs([]) to be false")
	}
}
