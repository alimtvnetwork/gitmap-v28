package uuidutil

import (
	"strings"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	if len(uuid) != UUIDLength {
		t.Fatalf("expected length %d, got %d for %q", UUIDLength, len(uuid), uuid)
	}

	// Verify hyphen positions
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		t.Errorf("invalid hyphen positions in UUID: %q", uuid)
	}

	// Verify version 4 character (at index 14)
	if uuid[14] != '4' {
		t.Errorf("expected version 4 character '4' at index 14, got %c in %q", uuid[14], uuid)
	}

	// Verify variant 1 character (at index 19 must be 8, 9, a, or b)
	variantChar := uuid[19]
	if variantChar != '8' && variantChar != '9' && variantChar != 'a' && variantChar != 'b' {
		t.Errorf("expected variant character 8/9/a/b at index 19, got %c in %q", variantChar, uuid)
	}
}

func TestGenerateUUIDUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		uuid := GenerateUUID()
		if seen[uuid] {
			t.Fatalf("duplicate UUID generated at iteration %d: %s", i, uuid)
		}
		seen[uuid] = true
	}
}

func TestFormatUUIDBytes(t *testing.T) {
	// Deterministic 16-byte raw vector: 00 01 02 03 04 05 46 07 88 09 0a 0b 0c 0d 0e 0f
	rawBytes := []byte{
		0x00, 0x01, 0x02, 0x03,
		0x04, 0x05,
		0x46, 0x07,
		0x88, 0x09,
		0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}

	formatted := FormatUUIDBytes(rawBytes)
	expected := "00010203-0405-4607-8809-0a0b0c0d0e0f"
	if !strings.EqualFold(formatted, expected) {
		t.Errorf("FormatUUIDBytes() = %q, want %q", formatted, expected)
	}
}

func TestNewAlias(t *testing.T) {
	uuid := New()
	if len(uuid) != UUIDLength {
		t.Fatalf("expected length %d, got %d", UUIDLength, len(uuid))
	}
}
