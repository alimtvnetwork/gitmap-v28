package cmdos

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/netip"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type mockIPDriver struct{}

func (m *mockIPDriver) Name() string { return "mock" }

func (m *mockIPDriver) ApplyConfig(ctx context.Context, opts netip.ChangeOptions) *apperror.AppError {
	return nil
}

func (m *mockIPDriver) RevertConfig(ctx context.Context, snap netip.Snapshot) *apperror.AppError {
	return nil
}

func (m *mockIPDriver) ListInterfaces(ctx context.Context) netip.InterfaceSliceResult {
	return result.OkSlice([]netip.InterfaceInfo{})
}

func (m *mockIPDriver) GetInterface(ctx context.Context, name string) netip.InterfaceResult {
	return result.Fail[netip.InterfaceInfo](apperror.NewSimple("not found", "NOT_FOUND"))
}

func setupMockIPManager() func() {
	orig := defaultNetIPManagerFactory
	defaultNetIPManagerFactory = func() *netip.Manager {
		return netip.NewManager(&mockIPDriver{}, nil)
	}
	return func() {
		defaultNetIPManagerFactory = orig
	}
}

func TestRunOSIP_Help(t *testing.T) {
	cleanup := setupMockIPManager()
	defer cleanup()

	helpCases := [][]string{{"--help"}, {"-h"}, {"help"}}
	for _, args := range helpCases {
		if err := runOSIP(args); err != nil {
			t.Fatalf("expected nil for help args %v, got: %v", args, err)
		}
	}
}

func TestRunOSIP_InvalidCommand(t *testing.T) {
	cleanup := setupMockIPManager()
	defer cleanup()

	if err := runOSIP([]string{"unknown-subcommand-xyz"}); err == nil {
		t.Fatal("expected error on unknown subcommand")
	}
}

func TestRunOSIP_ShowMocked(t *testing.T) {
	cleanup := setupMockIPManager()
	defer cleanup()

	if err := runOSIP([]string{"show"}); err != nil {
		t.Fatalf("unexpected error on show: %v", err)
	}
}
