package cmdsetup

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestGenerateLaravelAppKey(t *testing.T) {
	key1, err1 := GenerateLaravelAppKey()
	if err1 != nil {
		t.Fatalf("GenerateLaravelAppKey failed: %v", err1)
	}

	hasPrefix := strings.HasPrefix(key1, "base64:")
	if !hasPrefix {
		t.Fatalf("expected key to start with 'base64:', got %s", key1)
	}

	rawB64 := strings.TrimPrefix(key1, "base64:")
	decoded, decErr := base64.StdEncoding.DecodeString(rawB64)
	if decErr != nil {
		t.Fatalf("failed to decode base64 key: %v", decErr)
	}

	hasLen := len(decoded) == 32
	if !hasLen {
		t.Fatalf("expected 32 decoded bytes, got %d", len(decoded))
	}

	key2, _ := GenerateLaravelAppKey()
	isSame := key1 == key2
	if isSame {
		t.Fatal("two generated keys must be distinct")
	}
}

func TestMergeEnvContent(t *testing.T) {
	base := "# Laravel Config\nAPP_NAME=Laravel\nAPP_DEBUG=true\n"
	updates := map[string]string{
		"APP_NAME": "GitMap App",
		"APP_ENV":  "production",
	}

	order := []string{"APP_NAME", "APP_ENV"}
	out := MergeEnvContent(base, updates, order)
	hasComment := strings.Contains(out, "# Laravel Config")
	if !hasComment {
		t.Error("expected comment to be preserved")
	}

	hasUpdatedName := strings.Contains(out, `APP_NAME="GitMap App"`)
	if !hasUpdatedName {
		t.Errorf("expected quoted name update, got: %s", out)
	}

	hasAppEnv := strings.Contains(out, "APP_ENV=production")
	if !hasAppEnv {
		t.Errorf("expected new key APP_ENV, got: %s", out)
	}
}

func TestSynthesizeLaravelEnvAutoKey(t *testing.T) {
	base := "APP_NAME=Old\nAPP_KEY=\n"
	opts := LaravelEnvOptions{
		AppName:      "New App",
		DBConnection: "mysql",
		DBDatabase:   "laravel_db",
	}

	out, err := SynthesizeLaravelEnv(base, opts)
	if err != nil {
		t.Fatalf("SynthesizeLaravelEnv failed: %v", err)
	}

	hasKey := strings.Contains(out, "APP_KEY=base64:")
	if !hasKey {
		t.Errorf("expected generated APP_KEY, got: %s", out)
	}

	hasDB := strings.Contains(out, "DB_CONNECTION=mysql")
	if !hasDB {
		t.Errorf("expected DB_CONNECTION=mysql, got: %s", out)
	}
}
