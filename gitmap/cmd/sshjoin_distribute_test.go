package cmd

import (
	"context"
	"testing"
)

func TestSSHDistributeKeys(t *testing.T) {
	if err := DistributeKeysToHosts(context.Background(), nil, "root", 22); err == nil {
		t.Fatal("expected error on empty hosts")
	}

	if err := executeDistributeKeys([]string{}); err == nil {
		t.Fatal("expected error on empty args")
	}

	if err := executeSJBroadcast([]string{}); err == nil {
		t.Fatal("expected error on empty args")
	}

	if err := executeSJSyncProfile([]string{}); err == nil {
		t.Fatal("expected error on empty args")
	}
}
