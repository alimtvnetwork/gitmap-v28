package store

import "time"

// SSHHost represents an SSH host entry in the database.
type SSHHost struct {
	ID                string    `json:"id" db:"id"`
	Alias             string    `json:"alias" db:"alias"`
	IP                string    `json:"ip" db:"ip"`
	Username          string    `json:"username" db:"username"`
	Port              int       `json:"port,omitempty" db:"port"`
	EncryptedPassword string    `json:"encrypted_password,omitempty" db:"encrypted_password"`
	ClusterRole       string    `json:"cluster_role,omitempty" db:"cluster_role"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
