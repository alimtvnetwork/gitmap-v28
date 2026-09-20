package cmdssh

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestRunSSHPassCLI_EmptyArgs(t *testing.T) {
	err := RunSSHPassCLI([]string{})
	isNil := err == nil
	if isNil == false {
		t.Fatalf("expected nil error on empty args (prints help), got: %v", err)
	}
}

func TestRunSSHPassCLI_Help(t *testing.T) {
	err := RunSSHPassCLI([]string{"--help"})
	isNil := err == nil
	if isNil == false {
		t.Fatalf("expected nil error on --help, got: %v", err)
	}
}

func TestRunSSHPassCLI_ShowMissingTarget(t *testing.T) {
	err := RunSSHPassCLI([]string{"show"})
	isNil := err == nil
	if isNil {
		t.Fatal("expected validation error on missing show target, got nil")
	}
}

func TestRunSSHPassCLI_ShowSuccess(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		ctx := context.Background()
		joinArgs := []string{"testuser@192.168.1.99", "SecretPass123!", "box99"}
		_ = executeEnrollWithPassCLI(ctx, joinArgs)

		err := RunSSHPassCLI([]string{"show", "box99"})
		isNil := err == nil
		if isNil == false {
			t.Fatalf("expected nil error on show existing pass, got: %v", err)
		}
	})
}

func TestRunSSHPassCLI_ListSuccess(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		ctx := context.Background()
		joinArgs := []string{"testuser@192.168.1.99", "SecretPass123!", "box99"}
		_ = executeEnrollWithPassCLI(ctx, joinArgs)

		err := RunSSHPassCLI([]string{"ls"})
		isNil := err == nil
		if isNil == false {
			t.Fatalf("expected nil error on pass ls, got: %v", err)
		}
	})
}
