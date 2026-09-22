package cmdos

import (
	"testing"
)

func TestTweakMockOperations(t *testing.T) {
	mock := &mockTweakEngine{}
	prev := GetTweakEngine()
	defer SetTweakEngine(prev)
	SetTweakEngine(mock)

	_ = runOSTweak([]string{"context-menu", "classic"})
	if !mock.isClassicContext {
		t.Error("expected classic context menu")
	}

	_ = runOSTweak([]string{"start-menu", "classic"})
	if !mock.isClassicStart {
		t.Error("expected classic start menu")
	}

	_ = runOSTweak([]string{"power", "ultimate"})
	if !mock.isUltimatePower {
		t.Error("expected ultimate power scheme")
	}

	_ = runOSTweak([]string{"hibernate", "off"})
	if mock.isHibernate {
		t.Error("expected hibernation to be false")
	}

	_ = runOSTweak([]string{"telemetry", "off"})
	if !mock.isTelemetry {
		t.Error("expected telemetry disabled to be true")
	}

	_ = runOSTweak([]string{"activity", "off"})
	if !mock.isActivity {
		t.Error("expected activity feed disabled to be true")
	}

	_ = runOSTweak([]string{"search", "clean"})
	if !mock.isBing {
		t.Error("expected search disabled to be true")
	}
}
