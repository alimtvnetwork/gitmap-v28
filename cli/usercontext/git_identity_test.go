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

	if identity.Name != "Jane Doe" || identity.Email != "jane@example.com" {
		t.Errorf("expected Jane Doe / jane@example.com, got %s / %s", identity.Name, identity.Email)
	}

	if !identity.IsConfigured || !identity.HasName || !identity.HasEmail {
		t.Errorf("expected fully configured identity, got %+v", identity)
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

	if identity.IsConfigured || identity.HasName || identity.HasEmail {
		t.Errorf("expected unconfigured identity, got %+v", identity)
	}
}

func TestGetGitIdentity_PartialConfigured(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	execGitCommand = mockSingleKeyGitCommand("user.name", "Only Name")
	ident, err := GetGitIdentity(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ident.HasName || ident.HasEmail || ident.IsConfigured {
		t.Errorf("unexpected identity state: %+v", ident)
	}
}

func mockSingleKeyGitCommand(matchedKey, value string) func(context.Context, ...string) ([]byte, error) {
	return func(ctx context.Context, args ...string) ([]byte, error) {
		if args[len(args)-1] == matchedKey {
			return []byte(value + "\n"), nil
		}

		return nil, errors.New("not set")
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

func TestUnsetGitConfigKey_Global(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	var captured []string
	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		captured = args

		return []byte(""), nil
	}

	err := UnsetGitConfigKey(context.Background(), "user.email", true)
	assertUnsetArgs(t, err, captured, []string{"config", "--global", "--unset", "user.email"})
}

func TestUnsetGitConfigKey_Local(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	var captured []string
	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		captured = args

		return []byte(""), nil
	}

	err := UnsetGitConfigKey(context.Background(), "user.email", false)
	assertUnsetArgs(t, err, captured, []string{"config", "--unset", "user.email"})
}

func assertUnsetArgs(t *testing.T, err error, actual, expected []string) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Join(actual, " ") != strings.Join(expected, " ") {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestUnsetGitConfigKey_ExitStatusFive(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		return []byte(""), errors.New("exit status 5")
	}

	if err := UnsetGitConfigKey(context.Background(), "user.email", true); err != nil {
		t.Fatalf("expected nil for exit status 5, got %v", err)
	}
}

func TestUnsetGitConfigKey_Error(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		return []byte("permission denied"), errors.New("exit status 1")
	}

	if err := UnsetGitConfigKey(context.Background(), "user.email", false); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSetGitIdentityWithOptions_Unset(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	var calls [][]string
	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		calls = append(calls, args)

		return []byte(""), nil
	}

	opts := GitIdentityOptions{Name: "Bot", IsUnsetEmail: true, IsGlobal: false}
	if err := SetGitIdentityWithOptions(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCallContains(t, calls, "--unset")
}

func TestSetGitIdentityWithOptions_EmptyEmail(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	var calls [][]string
	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		calls = append(calls, args)

		return []byte(""), nil
	}

	opts := GitIdentityOptions{Name: "Bot", HasEmptyEmail: true, IsGlobal: false}
	if err := SetGitIdentityWithOptions(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCallContains(t, calls, "user.email")
}

func TestSetGitIdentity_EmptyEmailUnsets(t *testing.T) {
	orig := execGitCommand
	defer func() { execGitCommand = orig }()

	var calls [][]string
	execGitCommand = func(ctx context.Context, args ...string) ([]byte, error) {
		calls = append(calls, args)

		return []byte(""), nil
	}

	if err := SetGitIdentity("Bot", "", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCallContains(t, calls, "--unset")
}

func assertCallContains(t *testing.T, calls [][]string, token string) {
	hasToken := false
	for _, call := range calls {
		if strings.Join(call, " ") != "" && strings.Contains(strings.Join(call, " "), token) {
			hasToken = true
			break
		}
	}

	if !hasToken {
		t.Errorf("expected at least one call to contain %q, calls: %v", token, calls)
	}
}
