package usercontext

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestGetGitIdentity_Configured(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		key := args[len(args)-1]
		if key == "user.name" {
			return []byte("Jane Doe\n"), nil
		}
		if key == "user.email" {
			return []byte("jane@example.com\n"), nil
		}

		return nil, errors.New("unknown key")
	}

	identity, err := GetGitIdentity(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if identity.Name != "Jane Doe" {
		t.Errorf("expected Jane Doe, got %s", identity.Name)
	}

	if identity.Email != "jane@example.com" {
		t.Errorf("expected jane@example.com, got %s", identity.Email)
	}

	if !identity.IsConfigured {
		t.Errorf("expected IsConfigured to be true")
	}
}

func TestGetGitIdentity_Unconfigured(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		return nil, errors.New("exit status 1")
	}

	identity, err := GetGitIdentity(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if identity.Name != "" || identity.Email != "" {
		t.Errorf("expected empty name/email, got %s / %s", identity.Name, identity.Email)
	}

	if identity.IsConfigured {
		t.Errorf("expected IsConfigured to be false")
	}
}

func TestSetGitIdentity_Success(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	var calls [][]string
	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		calls = append(calls, args)

		return []byte(""), nil
	}

	err := SetGitIdentity("Alice", "alice@test.com", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}

	if calls[0][1] != "--global" || calls[0][3] != "Alice" {
		t.Errorf("unexpected call 0: %v", calls[0])
	}

	if calls[1][1] != "--global" || calls[1][3] != "alice@test.com" {
		t.Errorf("unexpected call 1: %v", calls[1])
	}
}

func TestSetGitIdentity_Error(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		return []byte("permission denied"), errors.New("exit 1")
	}

	err := SetGitIdentity("Bob", "bob@test.com", false)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "git config set failed") {
		t.Errorf("expected wrapped error message, got %v", err)
	}
}
