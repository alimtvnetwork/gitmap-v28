package startup

// Pin the registry-hklm flag value through ParseBackend / String so
// a future rename / typo cannot silently downgrade the machine-wide
// backend to the per-user default. OS-agnostic — the enum mapping
// itself doesn't touch the registry, so this runs on every CI runner.

import "testing"

func TestParseBackend_RegistryHKLMRoundTrip(t *testing.T) {
	b, err := ParseBackend("registry-hklm")
	if err != nil {
		t.Fatalf("ParseBackend(registry-hklm) err = %v", err)
	}
	if b != BackendRegistryHKLM {
		t.Fatalf("got enum %v, want BackendRegistryHKLM", b)
	}
	if got := b.String(); got != "registry-hklm" {
		t.Errorf("String() = %q, want %q", got, "registry-hklm")
	}
}

func TestParseBackend_UnknownStillRejected(t *testing.T) {
	if _, err := ParseBackend("hklm"); err == nil {
		t.Fatal("ParseBackend(hklm) err = nil, want non-nil")
	}
}
