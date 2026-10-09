package cmdcreate

import (
	"testing"
)

func TestNormalizeCreateArgs(t *testing.T) {
	argsWithRepo := []string{"repo", "my-test-service", "--private"}
	norm1 := normalizeCreateArgs(argsWithRepo)
	if len(norm1) != 2 || norm1[0] != "my-test-service" {
		t.Fatalf("expected ['my-test-service', '--private'], got: %v", norm1)
	}

	argsWithoutRepo := []string{"direct-service", "--public"}
	norm2 := normalizeCreateArgs(argsWithoutRepo)
	if len(norm2) != 2 || norm2[0] != "direct-service" {
		t.Fatalf("expected ['direct-service', '--public'], got: %v", norm2)
	}
}
