package store

import "time"

// SSHHost represents an SSH host entry in the database.
type SSHHost struct {
	ID        string    `json:"id" db:"id"`
	Alias     string    `json:"alias" db:"alias"`
	IP        string    `json:"ip" db:"ip"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
