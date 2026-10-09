package secrets

import (
	"strings"
	"testing"
)

func TestSaltedCipher_Roundtrip(t *testing.T) {
	passwords := []string{
		"rtyrty123@",
		"a",
		"Administrator",
		"Super$ecret!99#%",
		"Simple123",
		"   spaces and symbols ~!@#$%^&*()_+ `-=[]\\{}|;':\",./<>?   ",
	}

	for _, pw := range passwords {
		enc := EncryptSalted(pw, "")
		if !strings.HasPrefix(enc, "salt:") {
			t.Errorf("expected salt: prefix, got %s", enc)
		}
		dec, err := DecryptSalted(enc)
		if err != nil {
			t.Fatalf("failed to decrypt %s: %v", pw, err)
		}
		if dec != pw {
			t.Errorf("decrypted mismatch: expected %q, got %q", pw, dec)
		}
	}
}

func TestSaltedCipher_FixedSalt(t *testing.T) {
	pw := "rtyrty123@"
	salt := "9f8b2c4e"
	enc := EncryptSalted(pw, salt)
	if !strings.HasPrefix(enc, "salt:9f8b2c4e:") {
		t.Errorf("expected salt:9f8b2c4e: prefix, got %s", enc)
	}
	dec, err := DecryptSalted(enc)
	if err != nil {
		t.Fatalf("failed to decrypt with fixed salt: %v", err)
	}
	if dec != pw {
		t.Errorf("expected %q, got %q", pw, dec)
	}
}

func TestCaesarCipher_Roundtrip(t *testing.T) {
	passwords := []string{
		"rtyrty123@",
		"a",
		"SecretPass2026!",
	}

	for _, pw := range passwords {
		// without salt
		enc := EncryptCaesar(pw, 13, "")
		if !strings.HasPrefix(enc, "caesar:13:") {
			t.Errorf("expected caesar:13: prefix, got %s", enc)
		}
		dec, err := DecryptCaesar(enc)
		if err != nil {
			t.Fatalf("failed to decrypt caesar: %v", err)
		}
		if dec != pw {
			t.Errorf("caesar mismatch: expected %q, got %q", pw, dec)
		}

		// with salt
		encWithSalt := EncryptCaesar(pw, 25, "customsalt")
		if !strings.HasPrefix(encWithSalt, "caesar:25:customsalt:") {
			t.Errorf("expected caesar:25:customsalt: prefix, got %s", encWithSalt)
		}
		decWithSalt, err := DecryptCaesar(encWithSalt)
		if err != nil {
			t.Fatalf("failed to decrypt caesar with salt: %v", err)
		}
		if decWithSalt != pw {
			t.Errorf("caesar with salt mismatch: expected %q, got %q", pw, decWithSalt)
		}
	}
}
