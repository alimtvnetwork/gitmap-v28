package cmdos

import (
	"strings"
	"testing"
)

type mockDockOperator struct {
	cfg     DockConfig
	lastSet DockPosition
}

func (m *mockDockOperator) GetDockConfig() (DockConfig, error) {
	return m.cfg, nil
}

func (m *mockDockOperator) SetDockPosition(pos DockPosition) error {
	m.lastSet = pos
	return nil
}

func TestNormalizeDockPosition(t *testing.T) {
	for _, raw := range []string{"bottom", "BOTTOM", "b", "down"} {
		if pos, err := NormalizeDockPosition(raw); err != nil || pos != DockPositionBottom {
			t.Fatalf("expected bottom for %q, got %q (err=%v)", raw, pos, err)
		}
	}
	for _, raw := range []string{"left", "right", "top", "center"} {
		if _, err := NormalizeDockPosition(raw); err != nil {
			t.Fatalf("expected valid pos for %q, err: %v", raw, err)
		}
	}
	if _, err := NormalizeDockPosition("diagonal"); err == nil {
		t.Fatalf("expected error for invalid pos diagonal")
	}
}

func TestParseDockCLIOptions(t *testing.T) {
	testParseEmptyAndBottom(t)
	testParseRemoteNodeFlags(t)
}

func testParseEmptyAndBottom(t *testing.T) {
	opts, _, err := parseDockCLIOptions([]string{})
	if err != nil || opts.HasPosition {
		t.Fatalf("expected query mode, err: %v", err)
	}
	opts, _, err = parseDockCLIOptions([]string{"bottom"})
	if err != nil || opts.Position != DockPositionBottom {
		t.Fatalf("expected bottom pos, err: %v", err)
	}
}

func testParseRemoteNodeFlags(t *testing.T) {
	opts, _, err := parseDockCLIOptions([]string{"left", "--node", "u1"})
	if err != nil || opts.TargetNode != "u1" || opts.Position != DockPositionLeft {
		t.Fatalf("expected u1/left, err: %v", err)
	}
	opts, _, err = parseDockCLIOptions([]string{"--node=u1", "top"})
	if err != nil || opts.TargetNode != "u1" || opts.Position != DockPositionTop {
		t.Fatalf("expected u1/top, err: %v", err)
	}
}

func TestRunOSDockCommand_MockExecution(t *testing.T) {
	mock := &mockDockOperator{cfg: DockConfig{Position: "bottom", OS: "linux"}}
	defaultDockOperator = mock
	defer func() { defaultDockOperator = nil }()

	if err := RunOSDockCommand([]string{}); err != nil {
		t.Fatalf("query mode failed: %v", err)
	}
	if err := RunOSDockCommand([]string{"top"}); err != nil {
		t.Fatalf("set mode failed: %v", err)
	}
	if mock.lastSet != DockPositionTop {
		t.Fatalf("expected lastSet=top, got: %q", mock.lastSet)
	}
}

func TestRunOSDockCommand_RemoteDelegation(t *testing.T) {
	var delegated []string
	RunSSHExecFn = func(args []string) error {
		delegated = args
		return nil
	}
	defer func() { RunSSHExecFn = nil }()

	if err := RunOSDockCommand([]string{"bottom", "--node", "u1"}); err != nil {
		t.Fatalf("remote delegation failed: %v", err)
	}
	if len(delegated) < 2 || delegated[0] != "u1" || !strings.Contains(delegated[1], "bottom") {
		t.Fatalf("unexpected delegated args: %v", delegated)
	}
}
