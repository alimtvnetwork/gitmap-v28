package cmd

// JSON encoder for `gitmap ssh list --json`.
//
// Migrated off json.MarshalIndent onto gitmap/stablejson so key order
// becomes a compile-time decision rather than a reflection accident.
// Schema: spec/08-json-schemas/ssh-list.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// ssh-list wire keys. Names + order are the contract.
const (
	sshListKeyID          = "id"
	sshListKeyName        = "name"
	sshListKeyPrivatePath = "privatePath"
	sshListKeyPublicKey   = "publicKey"
	sshListKeyFingerprint = "fingerprint"
	sshListKeyEmail       = "email"
	sshListKeyCreatedAt   = "createdAt"
)

// encodeSSHListJSON writes keys as a stablejson 2-space-indented
// array. Empty input emits `[]\n`.
func encodeSSHListJSON(w io.Writer, keys []model.SSHKey) error {
	return stablejson.WriteArray(w, buildSSHListJSONItems(keys))
}

// buildSSHListJSONItems is the single source of (field name,
// field order, value) for ssh-list.
func buildSSHListJSONItems(keys []model.SSHKey) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(keys))
	for _, k := range keys {
		items = append(items, []stablejson.Field{
			{Key: sshListKeyID, Value: k.ID},
			{Key: sshListKeyName, Value: k.Name},
			{Key: sshListKeyPrivatePath, Value: k.PrivatePath},
			{Key: sshListKeyPublicKey, Value: k.PublicKey},
			{Key: sshListKeyFingerprint, Value: k.Fingerprint},
			{Key: sshListKeyEmail, Value: k.Email},
			{Key: sshListKeyCreatedAt, Value: k.CreatedAt},
		})
	}

	return items
}
