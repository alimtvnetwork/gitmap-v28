package cmdssh

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestParseAddPassParams(t *testing.T) {
	_, err := parseAddPassParams([]string{})
	if err == nil {
		t.Fatal("expected error on empty args")
	}

	p1, err := parseAddPassParams([]string{"alim@192.168.1.14"})
	if err != nil || p1.targetRaw != "alim@192.168.1.14" || p1.password != "" {
		t.Fatalf("unexpected result for 1 arg: %+v", p1)
	}

	p2, err := parseAddPassParams([]string{"alim@192.168.1.14", "secret"})
	if err != nil || p2.password != "secret" || p2.alias != "" {
		t.Fatalf("unexpected result for 2 args: %+v", p2)
	}

	p3, err := parseAddPassParams([]string{"alim@192.168.1.14", "secret", "mybox"})
	if err != nil || p3.password != "secret" || p3.alias != "mybox" {
		t.Fatalf("unexpected result for 3 args: %+v", p3)
	}
}

func TestExecuteEnrollWithPassCLI(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		ctx := context.Background()
		args := []string{"alim@192.168.1.14", "P@ss1234", "devbox"}
		if err := executeEnrollWithPassCLI(ctx, args); err != nil {
			t.Fatalf("executeEnrollWithPassCLI failed: %v", err)
		}

		host, err := store.GetHostByAlias(ctx, "devbox", db.SQL())
		if err != nil {
			t.Fatalf("failed to retrieve enrolled host: %v", err)
		}

		if host.Username != "alim" || host.IP != "192.168.1.14" {
			t.Fatalf("host metadata mismatch: %+v", host)
		}

		if host.EncryptedPassword == "" {
			t.Fatal("expected EncryptedPassword to be populated")
		}

		decrypted, err := DecryptSSHPassword(host.EncryptedPassword)
		if err != nil {
			t.Fatalf("failed to decrypt password: %v", err)
		}
		if decrypted != "P@ss1234" {
			t.Fatalf("decrypted password mismatch: got %q, want P@ss1234", decrypted)
		}
	})
}

func TestSJAddWithPassCmd_Execute(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		cmd := SJAddWithPassCmd
		args := []string{"root@10.0.0.15", "RootPass2026", "gateway"}
		if err := cmd.RunE(cmd, args); err != nil {
			t.Fatalf("SJAddWithPassCmd RunE failed: %v", err)
		}

		host, err := store.GetHostByAlias(context.Background(), "gateway", db.SQL())
		if err != nil {
			t.Fatalf("host not found: %v", err)
		}
		if host.EncryptedPassword == "" {
			t.Fatal("expected encrypted password")
		}
	})
}
