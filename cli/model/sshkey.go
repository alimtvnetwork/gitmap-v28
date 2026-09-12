// Package model defines the core data structures for gitmap.
package model

// SSHKey represents a stored SSH key pair.
type SSHKey struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PrivatePath string `json:"privatePath"`
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	Email       string `json:"email"`
	CreatedAt   string `json:"createdAt"`
}
