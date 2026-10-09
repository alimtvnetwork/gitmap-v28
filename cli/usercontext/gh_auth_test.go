package usercontext

import (
	"context"
	"errors"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

func TestDetectGhAuth_AuthorizedViaStatus(t *testing.T) {
	restore := setupTestMocks(
		func() (string, error) { return "/mock/gh", nil },
		func(_ context.Context, name string, args ...string) ([]byte, error) {
			if len(args) >= 2 && args[0] == "auth" && args[1] == "status" {
				return []byte("Logged in to github.com account octocat (keyring)\nToken scopes: 'repo'"), nil
			}
			return nil, errors.New("command failed")
		},
		func() (string, secrets.SourceType, error) { return "", secrets.SourceNone, secrets.ErrNoToken },
	)
	defer restore()

	info := DetectGhAuth()
	if info.Status != StatusAuthorized {
		t.Fatalf("expected status %s, got %s", StatusAuthorized, info.Status)
	}
	if info.Username != "octocat" {
		t.Fatalf("expected username octocat, got %s", info.Username)
	}
}

func TestDetectGhAuth_AuthorizedViaApi(t *testing.T) {
	restore := setupTestMocks(
		func() (string, error) { return "/mock/gh", nil },
		func(_ context.Context, name string, args ...string) ([]byte, error) {
			if len(args) >= 1 && args[0] == "api" {
				return []byte("monalisa\n"), nil
			}
			return nil, errors.New("command failed")
		},
		func() (string, secrets.SourceType, error) { return "", secrets.SourceNone, secrets.ErrNoToken },
	)
	defer restore()

	info := DetectGhAuth()
	if info.Status != StatusAuthorized {
		t.Fatalf("expected status %s, got %s", StatusAuthorized, info.Status)
	}
	if info.Username != "monalisa" {
		t.Fatalf("expected username monalisa, got %s", info.Username)
	}
}

func TestDetectGhAuth_FallbackToken(t *testing.T) {
	restore := setupTestMocks(
		func() (string, error) { return "/mock/gh", nil },
		func(_ context.Context, name string, args ...string) ([]byte, error) {
			if len(args) >= 2 && args[0] == "config" && args[1] == "user.name" {
				return []byte("gituser\n"), nil
			}
			return nil, errors.New("command failed")
		},
		func() (string, secrets.SourceType, error) {
			return "ghp_mock_token", secrets.SourceGitCredential, nil
		},
	)
	defer restore()

	info := DetectGhAuth()
	if info.Status != StatusAuthorized {
		t.Fatalf("expected status %s, got %s", StatusAuthorized, info.Status)
	}
	if info.Username != "gituser" {
		t.Fatalf("expected username gituser, got %s", info.Username)
	}
	if info.Source != string(secrets.SourceGitCredential) {
		t.Fatalf("expected source Git Credential Manager, got %s", info.Source)
	}
}

func TestDetectGhAuth_Unauthenticated(t *testing.T) {
	restore := setupTestMocks(
		func() (string, error) { return "/mock/gh", nil },
		func(_ context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("command failed")
		},
		func() (string, secrets.SourceType, error) { return "", secrets.SourceNone, secrets.ErrNoToken },
	)
	defer restore()

	info := DetectGhAuth()
	if info.Status != StatusUnauthenticated {
		t.Fatalf("expected status %s, got %s", StatusUnauthenticated, info.Status)
	}
}

func TestDetectGhAuth_ToolMissing(t *testing.T) {
	restore := setupTestMocks(
		func() (string, error) { return "", errors.New("not found") },
		func(_ context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("command failed")
		},
		func() (string, secrets.SourceType, error) { return "", secrets.SourceNone, secrets.ErrNoToken },
	)
	defer restore()

	info := DetectGhAuth()
	if info.Status != StatusToolMissing {
		t.Fatalf("expected status %s, got %s", StatusToolMissing, info.Status)
	}
}

type lookPathFunc func() (string, error)
type cmdOutputFunc func(context.Context, string, ...string) ([]byte, error)
type tokenResolveFunc func() (string, secrets.SourceType, error)

func setupTestMocks(lp lookPathFunc, co cmdOutputFunc, tr tokenResolveFunc) func() {
	origLookPath := lookPath
	origRunCmd := runCommandOutput
	origTokenResolver := tokenResolver

	lookPath = func(file string) (string, error) { return lp() }
	runCommandOutput = co
	tokenResolver = tr

	return func() {
		lookPath = origLookPath
		runCommandOutput = origRunCmd
		tokenResolver = origTokenResolver
	}
}
