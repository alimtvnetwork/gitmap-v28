package cmdssh

import (
	"context"
	"strings"
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

func TestRunSSHPassCLI_EncryptAndDecrypt(t *testing.T) {
	plain := "rtyrty123@"
	salt := "9f8b2c4e"

	enc := EncryptSaltedPassword(plain, salt)
	if !strings.HasPrefix(enc, "salt:9f8b2c4e:") {
		t.Fatalf("expected salt prefix in encrypted output, got: %s", enc)
	}

	dec, err := DecryptSSHPassword(enc)
	if err != nil {
		t.Fatalf("failed to decrypt password: %v", err)
	}
	if dec != plain {
		t.Fatalf("expected decrypted password %q, got %q", plain, dec)
	}

	// test CLI handlers
	if err := RunSSHPassCLI([]string{"encrypt", plain, "--salt", salt}); err != nil {
		t.Fatalf("run encrypt failed: %v", err)
	}
	if err := RunSSHPassCLI([]string{"decrypt", enc}); err != nil {
		t.Fatalf("run decrypt failed: %v", err)
	}
}
