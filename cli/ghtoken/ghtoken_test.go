package ghtoken

import (
	"os"
	"testing"
)

func TestExtractPasswordFromCredentialOutput(t *testing.T) {
	sample := "protocol=https\nhost=github.com\nusername=testuser\npassword=ghp_sampleSecretToken123\n"
	tok, isDefined := extractPasswordFromCredentialOutput(sample)

	if !isDefined {
		t.Fatalf("expected isDefined to be true")
	}

	if tok != "ghp_sampleSecretToken123" {
		t.Fatalf("expected ghp_sampleSecretToken123, got %s", tok)
	}
}

func TestExtractPasswordFromCredentialOutputEmpty(t *testing.T) {
	sample := "protocol=https\nhost=github.com\nusername=testuser\npassword=\n"
	_, isDefined := extractPasswordFromCredentialOutput(sample)

	if isDefined {
		t.Fatalf("expected isDefined to be false for empty password")
	}
}

func TestExtractOauthTokenFromYAML(t *testing.T) {
	yaml := "github.com:\n    user: alimtvnetwork\n    oauth_token: \"gho_testTokenFromYaml456\"\n"
	tok, isDefined := extractOauthTokenFromYAML(yaml)

	if !isDefined {
		t.Fatalf("expected isDefined to be true")
	}

	if tok != "gho_testTokenFromYaml456" {
		t.Fatalf("expected gho_testTokenFromYaml456, got %s", tok)
	}
}

func TestResolveFromEnv(t *testing.T) {
	orig := os.Getenv("GH_TOKEN")
	defer os.Setenv("GH_TOKEN", orig)

	os.Setenv("GH_TOKEN", "ghp_mockResolveToken999")
	tok, src, err := Resolve()

	if err != nil {
		t.Fatalf("unexpected error resolving token: %v", err)
	}

	if tok != "ghp_mockResolveToken999" {
		t.Fatalf("expected mock token, got %s", tok)
	}

	if src != SourceEnvGHToken {
		t.Fatalf("expected SourceEnvGHToken, got %s", src)
	}
}

func TestResolveFromSystem(t *testing.T) {
	// Clears env to verify fallback to system registry or credential manager
	origGH := os.Getenv("GH_TOKEN")
	origGitHub := os.Getenv("GITHUB_TOKEN")
	defer func() {
		os.Setenv("GH_TOKEN", origGH)
		os.Setenv("GITHUB_TOKEN", origGitHub)
	}()

	os.Setenv("GH_TOKEN", "")
	os.Setenv("GITHUB_TOKEN", "")

	tok, src, err := Resolve()
	if err != nil {
		return
	}

	t.Logf("System resolved token from source: %s (len=%d)", src, len(tok))
	if len(tok) == 0 {
		t.Fatalf("expected non-empty token when err is nil")
	}
}
