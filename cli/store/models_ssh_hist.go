package store

import (
	"time"
)

// SSHHistory represents a record of an SSH login.
type SSHHistory struct {
	ID       string    `json:"id" db:"id"`
	HostIP   string    `json:"host_ip" db:"host_ip"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	User     string    `json:"user" db:"user"`
}
