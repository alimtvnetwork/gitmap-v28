// Package model defines the core data structures for gitmap.
package model

// SSHKey represents a stored SSH key pair.
type SSHKey struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	PrivatePath         string `json:"privatePath"`
	PublicKey           string `json:"publicKey"`
	Fingerprint         string `json:"fingerprint"`
	Email               string `json:"email"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey,omitempty"`
	KeyType             string `json:"keyType,omitempty"`
	Provider            string `json:"provider,omitempty"`
	IsDefault           bool   `json:"isDefault,omitempty"`
	CreatedAt           string `json:"createdAt"`
}
