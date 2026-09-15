package osfix

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// FixItem models a registered OS repair script or command.
type FixItem struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Command     string    `json:"command"`
	Platform    string    `json:"platform"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type (
	// FixResult encapsulates a single FixItem result.
	FixResult = result.Result[FixItem]

	// FixSliceResult encapsulates a slice of FixItem results.
	FixSliceResult = result.ResultSlice[FixItem]

	// FixBoolResult encapsulates a boolean result.
	FixBoolResult = result.BoolResult
)
