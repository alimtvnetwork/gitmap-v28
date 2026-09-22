package cmdos

import (
	"testing"
)

type mockTweakEngine struct {
	isClassicContext bool
	isClassicStart   bool
	isUltimatePower  bool
	isHibernate      bool
}

func (m *mockTweakEngine) SetContextMenu(isClassic bool) error {
	m.isClassicContext = isClassic
	return nil
}

func (m *mockTweakEngine) SetStartMenu(isClassic bool) error {
	m.isClassicStart = isClassic
	return nil
}

func (m *mockTweakEngine) SetPowerScheme(isUltimate bool) error {
	m.isUltimatePower = isUltimate
	return nil
}

func (m *mockTweakEngine) SetHibernate(isEnabled bool) error {
	m.isHibernate = isEnabled
	return nil
}

func (m *mockTweakEngine) GetStatus() (TweakStatus, error) {
	return TweakStatus{
		IsClassicContextMenu: m.isClassicContext,
		IsClassicStartMenu:   m.isClassicStart,
		IsUltimatePower:      m.isUltimatePower,
		IsHibernateEnabled:   m.isHibernate,
	}, nil
}

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
}
