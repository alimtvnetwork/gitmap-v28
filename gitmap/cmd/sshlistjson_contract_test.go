package cmd

// JSON contract tests for `gitmap ssh list --json`.
//
// ssh-list emits an array of model.SSHKey. The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, name, privatePath, publicKey, fingerprint, email, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run SSHListJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestSSHListJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestSSHListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "ssh_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeSSHListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalSSHKey builds a deterministic single row.
func canonicalSSHKey() model.SSHKey {
	return model.SSHKey{
		ID:          3,
		Name:        "deploy-key",
		PrivatePath: "/home/user/.ssh/deploy-key",
		PublicKey:   "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI...",
		Fingerprint: "SHA256:abcd1234efgh5678ijkl9012mnop3456qrst7890",
		Email:       "deploy@example.com",
		CreatedAt:   "2025-01-15T08:30:00Z",
	}
}

// TestSSHListJSONContract_CanonicalRow_KeyOrder asserts the key
// order of the emitted object matches the schema declaration.
func TestSSHListJSONContract_CanonicalRow_KeyOrder(t *testing.T) {
	records := []model.SSHKey{canonicalSSHKey()}
	var buf bytes.Buffer
	if err := encodeSSHListJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "ssh-list")
}
