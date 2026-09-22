package cmdos

import (
	"testing"
)

type mockAutoLoginEngine struct {
	lastConfig AutoLoginConfig
	isDisabled bool
	mockStatus AutoLoginStatus
}

func (m *mockAutoLoginEngine) Configure(cfg AutoLoginConfig) error {
	m.lastConfig = cfg
	m.isDisabled = false
	return nil
}

func (m *mockAutoLoginEngine) Disable() error {
	m.isDisabled = true
	return nil
}

func (m *mockAutoLoginEngine) Status() (AutoLoginStatus, error) {
	return m.mockStatus, nil
}

func TestAutoLoginMockConfigure(t *testing.T) {
	mock := &mockAutoLoginEngine{}
	prev := GetAutoLoginEngine()
	defer SetAutoLoginEngine(prev)
	SetAutoLoginEngine(mock)

	cfg := AutoLoginConfig{
		Username:  "testuser",
		Domain:    "WORKGROUP",
		Password:  "secret",
		IsEnabled: true,
	}
	if err := GetAutoLoginEngine().Configure(cfg); err != nil {
		t.Fatalf("unexpected configure error: %v", err)
	}
	if mock.lastConfig.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", mock.lastConfig.Username)
	}
	if mock.lastConfig.Domain != "WORKGROUP" {
		t.Errorf("expected domain WORKGROUP, got %s", mock.lastConfig.Domain)
	}
}

func TestAutoLoginMockDisableAndStatus(t *testing.T) {
	mock := &mockAutoLoginEngine{
		mockStatus: AutoLoginStatus{
			IsEnabled:      true,
			Username:       "alice",
			Domain:         ".",
			HasPassword:    true,
			DisplayManager: "MockDM",
		},
	}
	prev := GetAutoLoginEngine()
	defer SetAutoLoginEngine(prev)
	SetAutoLoginEngine(mock)

	st, err := GetAutoLoginEngine().Status()
	if err != nil {
		t.Fatalf("unexpected status error: %v", err)
	}
	if !st.IsEnabled || st.Username != "alice" {
		t.Errorf("unexpected status: %+v", st)
	}

	if err := GetAutoLoginEngine().Disable(); err != nil {
		t.Fatalf("unexpected disable error: %v", err)
	}
	if !mock.isDisabled {
		t.Error("expected mock engine to be disabled")
	}
}
