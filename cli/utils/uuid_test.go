package utils

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
)

func TestUtilsUUID(t *testing.T) {
	uuid := GenerateUUID()
	if len(uuid) != uuidutil.UUIDLength {
		t.Fatalf("expected length %d, got %d", uuidutil.UUIDLength, len(uuid))
	}

	aliasUUID := NewUUID()
	if len(aliasUUID) != uuidutil.UUIDLength {
		t.Fatalf("expected length %d, got %d", uuidutil.UUIDLength, len(aliasUUID))
	}
}
