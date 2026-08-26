// Package installer — manager.go provides the central orchestrator for installer scripts.
package installer

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// Manager coordinates creating, versioning, exporting, importing, and executing installers.
type Manager struct {
	db *store.DB
}

// NewManager initializes a new installer Manager instance.
func NewManager(db *store.DB) (*Manager, error) {
	if db == nil {
		return nil, apperror.New("NewManager", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "db cannot be nil",
		})
	}
	return &Manager{db: db}, nil
}
