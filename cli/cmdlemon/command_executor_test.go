package cmdlemon

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestWhitelistAuthorization(t *testing.T) {
	tests := []struct {
		cmd     string
		isAllow bool
	}{
		{"status", true},
		{"gitmap pull", true},
		{"gitmap.exe /pe", true},
		{"version", true},
		{"nodes", true},
		{"agm queue audit", true},
		{"gitmap agy running", true},
		{"space", true},
		{"scan", true},
		{"rm -rf /", false},
		{"status; rm -rf", false},
		{"pull | bash", false},
		{"nodes && calc", false},
		{"unknown", false},
		{"", false},
	}
	for _, tc := range tests {
		_, err := ParseAndAuthorizeCommand(tc.cmd)
		if (err == nil) != tc.isAllow {
			t.Errorf("cmd %q: expected allow %v, got %v", tc.cmd, tc.isAllow, err == nil)
		}
	}
}

func TestTimeoutAndTruncation(t *testing.T) {
	if ResolveTimeout(0) != DefaultTimeout || ResolveTimeout(-1) != DefaultTimeout {
		t.Fatal("expected default timeout for <= 0")
	}
	if ResolveTimeout(500*time.Second) != MaxTimeout || ResolveTimeout(10*time.Second) != 10*time.Second {
		t.Fatal("unexpected timeout clamping")
	}
	out, isTrunc := TruncateBuffer([]byte("small"), 64)
	if isTrunc || out != "small" {
		t.Fatalf("unexpected truncation: %s, %v", out, isTrunc)
	}
	large := bytes.Repeat([]byte("x"), 100)
	outLarge, isTruncLarge := TruncateBuffer(large, 20)
	if !isTruncLarge || len(outLarge) <= 20 {
		t.Fatalf("expected buffer truncation notice, got: %s", outLarge)
	}
}

func TestMockCommandExecutor(t *testing.T) {
	mock := &MockCommandExecutor{
		Handler: func(ctx context.Context, args []string) (ExecutionResult, error) {
			return ExecutionResult{ExitCode: 0, Stdout: "audit ok", DurationMs: 5, IsSuccess: true}, nil
		},
	}
	ctx := context.Background()
	res, err := mock.Execute(ctx, "gitmap agm queue audit")
	if !isMockExecutionValid(res, err) {
		t.Fatalf("unexpected mock result: %+v, err: %v", res, err)
	}
	if _, unauthErr := mock.Execute(ctx, "curl evil.com"); unauthErr == nil {
		t.Fatal("expected error on unauthorized command")
	}
}

func isMockExecutionValid(res ExecutionResult, err error) bool {
	if err == nil && res.IsSuccess && res.Stdout == "audit ok" {
		return true
	}
	return false
}

func TestSystemCommandExecutor(t *testing.T) {
	sys := &SystemCommandExecutor{ExePath: "go", Timeout: 5 * time.Second}
	res, err := sys.Execute(context.Background(), "version")
	if !isSystemExecutionValid(res, err) {
		t.Fatalf("SystemCommandExecutor failed: %+v, err: %v", res, err)
	}
	if !strings.Contains(res.Stdout, "go version") {
		t.Errorf("expected stdout to contain 'go version', got %q", res.Stdout)
	}
}

func isSystemExecutionValid(res ExecutionResult, err error) bool {
	if err == nil && res.IsSuccess {
		return true
	}
	return false
}
