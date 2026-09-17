package cmdservice

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mockDriver struct {
	services []ServiceInfo
	started  []string
	stopped  []string
	created  []string
	removed  []string
	isFail   bool
	err      error
}

func (m *mockDriver) ListServices() ([]ServiceInfo, error) {
	if m.isFail {
		return nil, m.err
	}
	return m.services, nil
}

func (m *mockDriver) GetService(name string) (*ServiceInfo, error) {
	if m.isFail {
		return nil, m.err
	}
	for _, s := range m.services {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, errors.New("service not found")
}

func (m *mockDriver) StartService(name string) error {
	if m.isFail {
		return m.err
	}
	m.started = append(m.started, name)
	return nil
}

func (m *mockDriver) StopService(name string) error {
	if m.isFail {
		return m.err
	}
	m.stopped = append(m.stopped, name)
	return nil
}

func (m *mockDriver) CreateService(name, execPath, description string) error {
	if m.isFail {
		return m.err
	}
	m.created = append(m.created, name)
	return nil
}

func (m *mockDriver) RemoveService(name string) error {
	if m.isFail {
		return m.err
	}
	m.removed = append(m.removed, name)
	return nil
}

func setupMockDriver(mock *mockDriver) func() {
	orig := DefaultServiceDriverResolver
	DefaultServiceDriverResolver = func() ServiceDriver {
		return mock
	}
	return func() {
		DefaultServiceDriverResolver = orig
	}
}

func createTestServices() []ServiceInfo {
	return []ServiceInfo{
		{Name: "svc-alpha", Status: "RUNNING", IsEnabled: true, Description: "Alpha Service"},
		{Name: "svc-beta", Status: "STOPPED", IsEnabled: false, Description: "Beta Service"},
	}
}

func TestRunService_Help(t *testing.T) {
	cleanup := setupMockDriver(&mockDriver{})
	defer cleanup()

	helpCases := [][]string{{}, {"--help"}, {"-h"}, {"help"}}
	for _, args := range helpCases {
		if err := RunService(args); err != nil {
			t.Fatalf("expected nil for help args %v, got: %v", args, err)
		}
	}
}

func TestRunService_List(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	if err := RunService([]string{"ls"}); err != nil {
		t.Fatalf("unexpected error for ls: %v", err)
	}
	if err := RunService([]string{"list", "--json"}); err != nil {
		t.Fatalf("unexpected error for list --json: %v", err)
	}
}

func TestRunService_Status(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	if err := RunService([]string{"status", "svc-alpha"}); err != nil {
		t.Fatalf("unexpected error for status: %v", err)
	}
	if err := RunService([]string{"status"}); err == nil {
		t.Fatal("expected error for status without service name")
	}
}

func TestRunService_Start(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	if err := RunService([]string{"start", "svc-beta"}); err != nil {
		t.Fatalf("unexpected error for start: %v", err)
	}
	if len(mock.started) != 1 || mock.started[0] != "svc-beta" {
		t.Fatalf("expected svc-beta started, got: %v", mock.started)
	}
	if err := RunService([]string{"start"}); err == nil {
		t.Fatal("expected error for start without name")
	}
}

func TestRunService_Stop(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	if err := RunService([]string{"stop", "svc-alpha"}); err != nil {
		t.Fatalf("unexpected error for stop: %v", err)
	}
	if len(mock.stopped) != 1 || mock.stopped[0] != "svc-alpha" {
		t.Fatalf("expected svc-alpha stopped, got: %v", mock.stopped)
	}
	if err := RunService([]string{"stop"}); err == nil {
		t.Fatal("expected error for stop without name")
	}
}

func TestRunService_Create(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	args := []string{"create", "new-svc", "--exec", "/usr/bin/app", "--desc", "Test App"}
	if err := RunService(args); err != nil {
		t.Fatalf("unexpected error for create: %v", err)
	}
	if len(mock.created) != 1 || mock.created[0] != "new-svc" {
		t.Fatalf("expected new-svc created, got: %v", mock.created)
	}
	if err := RunService([]string{"create"}); err == nil {
		t.Fatal("expected error for create without name")
	}
}

func TestRunService_Remove(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	if err := RunService([]string{"rm", "svc-alpha"}); err != nil {
		t.Fatalf("unexpected error for rm: %v", err)
	}
	if len(mock.removed) != 1 || mock.removed[0] != "svc-alpha" {
		t.Fatalf("expected svc-alpha removed, got: %v", mock.removed)
	}
	if err := RunService([]string{"rm"}); err == nil {
		t.Fatal("expected error for rm without name")
	}
}

func TestRunService_ExportAndImport(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	tempFile := filepath.Join(t.TempDir(), "services.json")
	if err := RunService([]string{"export", "--file", tempFile}); err != nil {
		t.Fatalf("unexpected export error: %v", err)
	}

	content, readErr := os.ReadFile(tempFile)
	if readErr != nil || len(content) == 0 {
		t.Fatalf("expected export file written, got err: %v", readErr)
	}

	if err := RunService([]string{"import", tempFile}); err != nil {
		t.Fatalf("unexpected import error: %v", err)
	}
	if err := RunService([]string{"import"}); err == nil {
		t.Fatal("expected error for import without file")
	}
}

func TestRunService_ExportYAML(t *testing.T) {
	mock := &mockDriver{services: createTestServices()}
	cleanup := setupMockDriver(mock)
	defer cleanup()

	tempFile := filepath.Join(t.TempDir(), "services.yaml")
	if err := RunService([]string{"export", "--file", tempFile, "--yaml"}); err != nil {
		t.Fatalf("unexpected export error: %v", err)
	}
	if err := RunService([]string{"import", tempFile}); err != nil {
		t.Fatalf("unexpected yaml import error: %v", err)
	}
}

func TestRunService_UnknownCommand(t *testing.T) {
	cleanup := setupMockDriver(&mockDriver{})
	defer cleanup()

	if err := RunService([]string{"unknown-cmd"}); err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestEnsureServiceDriver_Unsupported(t *testing.T) {
	orig := DefaultServiceDriverResolver
	defer func() { DefaultServiceDriverResolver = orig }()

	DefaultServiceDriverResolver = func() ServiceDriver {
		return nil
	}

	driver, appErr := EnsureServiceDriver()
	if driver != nil || appErr == nil {
		t.Fatal("expected nil driver and appErr for unsupported platform")
	}
}
