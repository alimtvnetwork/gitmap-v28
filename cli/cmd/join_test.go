package cmd

import (
	"crypto/tls"
	"net"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
)

func TestRunJoin_ZeroArgsTriggersHelp(t *testing.T) {
	var isExitCalled bool
	prev := cliexit.SetExitFunc(func(code int) {
		isExitCalled = true
	})
	defer cliexit.SetExitFunc(prev)

	_ = runJoin([]string{})
	if !isExitCalled {
		t.Fatal("expected exit to be called on zero args")
	}
}

func verifyJoinHelpFlag(t *testing.T, flag string) {
	var isExitCalled bool
	prev := cliexit.SetExitFunc(func(code int) {
		isExitCalled = true
	})
	defer cliexit.SetExitFunc(prev)

	_ = runJoin([]string{flag})
	if !isExitCalled {
		t.Fatalf("expected exit called for help flag %s", flag)
	}
}

func TestRunJoin_HelpFlags(t *testing.T) {
	flags := []string{"-h", "--help", "help"}
	for _, flagName := range flags {
		verifyJoinHelpFlag(t, flagName)
	}
}

func TestRunJoin_MissingAddress(t *testing.T) {
	err := runJoin([]string{"--token", "my-token"})
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
	if !err.IsErrorCode("E1000") {
		t.Fatalf("expected code E1000, got %s", err.Code)
	}
	if err.Type != apperror.ErrorTypeValidation {
		t.Fatalf("expected validation error type, got %s", err.Type)
	}
}

func TestRunJoin_MissingToken(t *testing.T) {
	err := runJoin([]string{"127.0.0.1:9000"})
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
	if !err.IsErrorCode("E1000") {
		t.Fatalf("expected code E1000, got %s", err.Code)
	}
	if err.Type != apperror.ErrorTypeValidation {
		t.Fatalf("expected validation error type, got %s", err.Type)
	}
}

func TestRunJoin_InvalidFlag(t *testing.T) {
	err := runJoin([]string{"--nonexistent-flag"})
	if err == nil {
		t.Fatal("expected error on invalid flag, got nil")
	}
	if !err.IsErrorCode("E1000") {
		t.Fatalf("expected code E1000, got %s", err.Code)
	}
	if err.Type != apperror.ErrorTypeValidation {
		t.Fatalf("expected validation error type, got %s", err.Type)
	}
}

func TestRunJoin_HandshakeFailure(t *testing.T) {
	err := runJoin([]string{"127.0.0.1:1", "--token", "test-token"})
	if err == nil {
		t.Fatal("expected handshake failure on unreachable address, got nil")
	}
	if !err.IsErrorCode("E8005") {
		t.Fatalf("expected code E8005, got %s", err.Code)
	}
	if err.Type != apperror.ErrorTypeExecution {
		t.Fatalf("expected execution error type, got %s", err.Type)
	}
}

func startMockJoinServer(t *testing.T, token string) (net.Listener, func()) {
	tlsConf, err := cluster.GenerateTLSConfig()
	if err != nil {
		t.Fatalf("failed to generate TLS config: %v", err)
	}
	listener, err := tls.Listen("tcp", "127.0.0.1:0", tlsConf)
	if err != nil {
		t.Fatalf("failed to listen on random port: %v", err)
	}
	srv := cluster.NewServer(token, 10*time.Second)
	go srv.Serve(listener)
	return listener, func() { _ = listener.Close() }
}

func TestRunJoin_Success(t *testing.T) {
	token := "valid-join-token"
	listener, cleanup := startMockJoinServer(t, token)
	defer cleanup()

	err := runJoin([]string{listener.Addr().String(), "--token", token})
	if err != nil {
		t.Fatalf("expected successful join, got error: %v", err)
	}
}
