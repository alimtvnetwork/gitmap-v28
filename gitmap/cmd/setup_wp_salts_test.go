package cmd

import (
	"strings"
	"testing"
)

func TestGenerateWpSalts(t *testing.T) {
	salts, err := GenerateWpSalts()
	if err != nil {
		t.Fatalf("GenerateWpSalts failed: %v", err)
	}

	keys := []string{
		salts.AuthKey, salts.SecureAuthKey, salts.LoggedInKey, salts.NonceKey,
		salts.AuthSalt, salts.SecureAuthSalt, salts.LoggedInSalt, salts.NonceSalt,
	}

	seen := make(map[string]bool)
	for _, k := range keys {
		hasLen := len(k) == WpSaltLength
		if !hasLen {
			t.Errorf("expected salt length %d, got %d", WpSaltLength, len(k))
		}

		isDup := seen[k]
		if isDup {
			t.Errorf("duplicate salt detected: %s", k)
		}

		seen[k] = true
	}
}

func TestGenerateRandomSaltInvalidLength(t *testing.T) {
	_, err := GenerateRandomSalt(0)
	if err == nil {
		t.Fatal("expected error for salt length <= 0")
	}
}

func TestFormatWpSaltsPHP(t *testing.T) {
	salts, err := GenerateWpSalts()
	if err != nil {
		t.Fatalf("GenerateWpSalts failed: %v", err)
	}

	out := FormatWpSaltsPHP(salts)
	hasAuthKey := strings.Contains(out, "define('AUTH_KEY', '")
	if !hasAuthKey {
		t.Error("expected PHP output to define AUTH_KEY")
	}

	hasNonceSalt := strings.Contains(out, "define('NONCE_SALT', '")
	if !hasNonceSalt {
		t.Error("expected PHP output to define NONCE_SALT")
	}
}

func TestWrapWpSaltsMarker(t *testing.T) {
	content := "define('FOO', 'bar');"
	wrapped := WrapWpSaltsMarker(content)
	hasBegin := strings.Contains(wrapped, "// >>> gitmap:wp-salts >>>")
	if !hasBegin {
		t.Error("expected begin marker")
	}

	hasEnd := strings.Contains(wrapped, "// <<< gitmap:wp-salts <<<")
	if !hasEnd {
		t.Error("expected end marker")
	}
}

func TestInjectWpSaltsExistingMarker(t *testing.T) {
	orig := "<?php\n// >>> gitmap:wp-salts >>>\nold_salt\n// <<< gitmap:wp-salts <<<\n"
	salts, err := GenerateWpSalts()
	if err != nil {
		t.Fatalf("GenerateWpSalts failed: %v", err)
	}

	injected := InjectWpSalts(orig, salts)
	hasOld := strings.Contains(injected, "old_salt")
	if hasOld {
		t.Error("expected old_salt to be replaced")
	}

	hasNew := strings.Contains(injected, salts.AuthKey)
	if !hasNew {
		t.Error("expected new salt to be injected")
	}
}

func TestInjectWpSaltsBeforeStopEditing(t *testing.T) {
	orig := "<?php\n/* That's all, stop editing! Happy publishing. */\n"
	salts, err := GenerateWpSalts()
	if err != nil {
		t.Fatalf("GenerateWpSalts failed: %v", err)
	}

	injected := InjectWpSalts(orig, salts)
	hasNew := strings.Contains(injected, salts.AuthKey)
	if !hasNew {
		t.Error("expected salts before stop editing comment")
	}
}
